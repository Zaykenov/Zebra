using HealthChecks.UI.Client;
using Microsoft.AspNetCore.Authentication.JwtBearer;
using Microsoft.AspNetCore.Diagnostics.HealthChecks;
using Microsoft.EntityFrameworkCore;
using Microsoft.EntityFrameworkCore.Migrations;
using Microsoft.Extensions.Caching.Memory;
using Microsoft.Extensions.Diagnostics.HealthChecks;
using Microsoft.IdentityModel.Tokens;
using Microsoft.OpenApi.Models;
using MobileApi.Controllers;
using MobileApi.Controllers.Auth;
using MobileApi.Controllers.Registration;
using MobileApi.Controllers.UserFeedback;
using MobileApi.Controllers.UserInfo;
using MobileApi.Controllers.UserOrders;
using MobileApi.Controllers.UserQr;
using MobileApi.Sources.db.zebra;
using MobileApi.Sources.email;
using MobileApi.Sources.Timers;
using Newtonsoft.Json.Serialization;
using Prometheus;
using Serilog;
using System;
using System.Net;
using System.Reflection;
using System.Security.Claims;

namespace MobileApi {
    public class Startup {
        private readonly string _apiDescription;

        public Startup(IConfiguration configuration, IHostEnvironment hostEnvironment) {
            Configuration = configuration;
            _apiDescription = getCurrentBuildInfo(hostEnvironment);

            if (configuration.GetSection("Serilog").Exists()) {
                Log.Logger = new LoggerConfiguration()
                    .ReadFrom.Configuration(configuration)
                    .CreateLogger();
            }
            else {
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
                    options.AddPolicy("MobileApiScope", policy => {
                        policy.RequireAuthenticatedUser();
                        policy.RequireClaim("scope", "MobileApi");
                    });
                });

            services.AddMemoryCache();

            string zebraConnStr = Configuration.GetConnectionString("dbZebraConnStr")!;
            services.AddDbContext<ZebraApplicationContext>(options => options.UseNpgsql(zebraConnStr, x => x.MigrationsHistoryTable(HistoryRepository.DefaultTableName, "mobile_api")));

            services.AddHttpClient();


            services.Configure<EmailSenderConfig>(options => {
                options.UserName = Configuration["EmailConf:UserName"]!;
                options.UserPwd = Configuration["EmailConf:UserPwd"]!;
                options.RegisterContinueOnMobileLink = Configuration["EmailConf:RegisterContinueOnMobileLink"]!;
            });
            services.AddScoped<ISignInEmailSenderService, SignInEmailSenderService>();
            services.AddScoped<IRegistrationEmailSenderService, RegistrationEmailSenderService>();

            services.AddScoped<IAuthService, AuthService>();
            services.AddScoped<IRegistrationService, RegistrationService>();
            services.AddScoped<IUserQrService, UserQrService>();
            services.AddScoped<IUserInfoService, UserInfoService>();
            services.AddScoped<IUserOrdersService, UserOrdersService>();
            services.AddScoped<IUserFeedbackService, UserFeedbackService>();
            services.AddScoped<IShopService, ShopService>();
            services.AddHttpContextAccessor();
            services.AddHealthChecks();

            services.AddHostedService<ZebraCoinCalcTimer>();

            services.AddSwaggerGen(c => {
                c.SwaggerDoc("v1", new OpenApiInfo { Title = "MobileApi", Version = "v1", Description = _apiDescription });
                c.IncludeXmlComments(Path.Combine(AppContext.BaseDirectory, "MobileApi.xml"), true);

                c.UseAllOfToExtendReferenceSchemas();
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

            var accessControlAllowOriginsSection = Configuration.GetSection("AccessControlAllowOrigins");
            if (accessControlAllowOriginsSection.Exists()) {
                var origins = accessControlAllowOriginsSection.Get<string[]>();
                app.UseCors((options) => {
                    options.WithOrigins(origins);
                });
            }

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
                c.SwaggerEndpoint("/swagger/v1/swagger.json", "MobileApi V1");
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
