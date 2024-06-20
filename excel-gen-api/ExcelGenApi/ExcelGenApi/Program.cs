using Serilog;
using System.Reflection;

namespace ExcelGenApi
{
    public class Program {
        public static void Main(string[] args) {
            var app = CreateHostBuilder(args).Build();


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
