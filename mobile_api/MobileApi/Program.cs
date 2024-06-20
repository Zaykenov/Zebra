using Microsoft.EntityFrameworkCore;
using MobileApi.Sources.db.zebra;
using Serilog;
using System.Reflection;

namespace MobileApi {
    public class Program {
        public static void Main(string[] args) {
            var app = CreateHostBuilder(args).Build();


            // run before migration:
            // 
            //  CREATE SCHEMA mobile_api;
            //  CREATE ROLE mobile_api_user NOSUPERUSER NOCREATEDB NOCREATEROLE NOINHERIT LOGIN NOREPLICATION NOBYPASSRLS PASSWORD '123456';
            //  GRANT ALL ON SCHEMA mobile_api TO mobile_api_user;
            //  GRANT SELECT ON ALL TABLES IN SCHEMA public TO mobile_api_user;
            //

            //run migration -> EntityFrameworkCore\ADD-MIGRATION email-nullable
            using (var scope = app.Services.CreateScope()) {
                var dbContext = scope.ServiceProvider.GetRequiredService<ZebraApplicationContext>();
                dbContext.Database.Migrate();
            }

            app.Run();
        }

        public static IHostBuilder CreateHostBuilder(string[] args) =>
            Host.CreateDefaultBuilder(args)
                .UseSerilog()
                .ConfigureAppConfiguration((context, configuration) => {
                    configuration
                        .AddJsonFile("appsettings.Development.localsettings.json", optional: true)
                        .AddUserSecrets(Assembly.GetEntryAssembly(), true)
                        .AddEnvironmentVariables()
                        .AddCommandLine(args);
                })
                .ConfigureWebHostDefaults(webBuilder => {
                    webBuilder
                        .UseStartup<Startup>();
                });
    }
}
