using ExcelGenApi.Controllers;
using ExcelGenApi.Interfaces;
using ExcelGenApi.Services;
using HealthChecks.UI.Client;
using Microsoft.AspNetCore.Authentication.JwtBearer;
using Microsoft.AspNetCore.Diagnostics.HealthChecks;
using Microsoft.Extensions.Caching.Memory;
using Microsoft.Extensions.Diagnostics.HealthChecks;
using Microsoft.IdentityModel.Tokens;
using Microsoft.OpenApi.Models;
using Newtonsoft.Json.Serialization;
using Prometheus;
using Serilog;
using System.Net;
using System.Reflection;
using System.Security.Claims;

namespace ExcelGenApi {
    public class Startup {
        private readonly string _apiDescription;

        public Startup(IConfiguration configuration, IHostEnvironment hostEnvironment) {
            Configuration = configuration;
            _apiDescription = getCurrentBuildInfo(hostEnvironment);

            if (configuration.GetSection("Serilog").Exists()) {
                Log.Logger = new LoggerConfiguration()
                    .ReadFrom.Configuration(configuration)
                    .CreateLogger();
            } else {
                Log.Logger = new LoggerConfiguration()
                    .WriteTo.Console(
                        restrictedToMinimumLevel: Serilog.Events.LogEventLevel.Warning
                    )
                    .CreateLogger();
            }
        }

        private static string getCurrentBuildInfo(IHostEnvironment hostEnvironment) {
            var assembly = Assembly.GetExecutingAssembly();
            var buildDate = File.GetCreationTime(assembly.Location);
            return $"[{assembly.GetName().FullName}, Creation Time={buildDate:O}; Environment={hostEnvironment.EnvironmentName}]";
        }

        public IConfiguration Configuration { get; }


        public void ConfigureServices(IServiceCollection services) {

            services
                .AddControllers()
                .AddNewtonsoftJson(options => {
                    options.SerializerSettings.ContractResolver = new DefaultContractResolver();
                });

            services.AddDataProtection();
            services.AddLogging();

            var identityServerUrl = Configuration["IdentityServerUrl"];
            services
                .AddAuthentication(JwtBearerDefaults.AuthenticationScheme)
                .AddJwtBearer(options => {
                    options.Authority = identityServerUrl;
                    options.TokenValidationParameters = new TokenValidationParameters {
                        ValidateAudience = false
                    };
                });

            services
                .AddAuthorization(options => {
                    options.AddPolicy("ExcelGenApiScope", policy => {
                        policy.RequireAuthenticatedUser();
                        policy.RequireClaim("scope", "ExcelGenApi");
                    });
                });

            services.AddMemoryCache();

            services.AddHttpClient();


            services.AddScoped<IExcelGenerator, ExcelGenerator>();

            services.AddHttpContextAccessor();
            services.AddHealthChecks();

            services.AddSwaggerGen(c => {
                c.SwaggerDoc("v1", new OpenApiInfo { Title = "ExcelGenApi", Version = "v1", Description = _apiDescription });
                c.IncludeXmlComments(Path.Combine(AppContext.BaseDirectory, "ExcelGenApi.xml"), true);

                c.UseAllOfToExtendReferenceSchemas();
            });
            services.AddCors(options =>
            {
                options.AddDefaultPolicy(builder =>
                {
                    builder.AllowAnyOrigin()
                        .AllowAnyMethod()
                        .AllowAnyHeader();
                });
            });
        }


        public void Configure(IApplicationBuilder app, IWebHostEnvironment env) {
            if (env.IsDevelopment()) {
                app.UseDeveloperExceptionPage();
            }

            app.UseHealthChecks("/healthz-check-ui-endpoint", new HealthCheckOptions {
                Predicate = _ => true,
                ResponseWriter = UIResponseWriter.WriteHealthCheckUIResponse
            });

            app.UseHealthChecksPrometheusExporter(
                "/healthz-check-prometheus-endpoint",
                options => options.ResultStatusCodes[HealthStatus.Unhealthy] = (int)HttpStatusCode.OK
            );
            app.UseCors();
            app.UseHttpsRedirection();

            app.UseRouting();
            app.UseHttpMetrics();

            app.UseStaticFiles();

            app.UseAuthentication();

            if (env.EnvironmentName == "Development" || Configuration.GetValue<bool?>("DisableAuthentication") == true) {
                app.Use((context, next) => {
                    context.User = new ClaimsPrincipal(new ClaimsIdentity(new List<Claim> {
                        new Claim(ClaimTypes.NameIdentifier, "000000000001"),
                        new Claim(ClaimTypes.Name, "Test user"),
                        new Claim(ClaimTypes.Email, "test@example.com"),
                        new Claim("scope", "MobileApi"),
                    }, JwtBearerDefaults.AuthenticationScheme));
                    var auth = context.User.Identity.IsAuthenticated;
                    return next();
                });
            }

            app.UseAuthorization();

            app.UseSwagger(new Swashbuckle.AspNetCore.Swagger.SwaggerOptions());
            app.UseSwaggerUI(c => {
                c.SwaggerEndpoint("/swagger/v1/swagger.json", "ExcelGenApi V1");
            });

            var trustedOriginsForWebWidgetsSection = Configuration.GetSection("trustedOriginsForWebWidgets");
            if (trustedOriginsForWebWidgetsSection.Exists()) {
                app.Use(async (context, next) => {
                    var urls = trustedOriginsForWebWidgetsSection.Get<string[]>();
                    context.Response.Headers.Add("Content-Security-Policy", $"frame-ancestors 'self' {string.Join(" ", urls)};");
                    await next(context);
                });
            }

            app.UseEndpoints(endpoints => {
                endpoints.MapControllers();
            });
        }
    }
}
