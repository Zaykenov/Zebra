using System.Text;
using ClickHouse.Ado;
using ClickHouseProxy;
using Microsoft.AspNetCore.Server.Kestrel.Core;
using Yarp.ReverseProxy.Transforms;

var builder = WebApplication.CreateBuilder(args);

builder.Services.Configure<KestrelServerOptions>(options => {
    options.Limits.MaxRequestBodySize = 100000000;
});

builder.Services.AddMemoryCache();
builder.Services.AddLogging();

builder.Services.AddHostedService<ClickHouseTimer>();

builder.Services.AddReverseProxy()
    .LoadFromConfig(builder.Configuration.GetSection("ReverseProxy"))
    .AddTransforms(async builderContext => {

        try {
            var clickHouseConnectionString = builder.Configuration["clickHouseConnectionString"];

            using (var clickHouseConnection = new ClickHouseConnection(clickHouseConnectionString)) {
                clickHouseConnection.Open();

                var command = clickHouseConnection.CreateCommand();
                command.CommandText = ClickHouseHelper.CreateTable;
                command.CommandTimeout = 3000;
                command.ExecuteNonQuery();
            }


            builderContext.AddRequestTransform(async x => {
                var body = "";
                var isImage = x.ProxyRequest?.Content?.Headers?.ContentType?.ToString()?.ToLower()?.Contains("image") ?? false;
                var hasAnyContent = (x.ProxyRequest?.Content?.Headers?.ContentLength ?? 0) > 0;

                var needToWriteBody = hasAnyContent && !isImage;

                if (needToWriteBody) {
                    using (var reader = new StreamReader(x.HttpContext.Request.Body)) {
                        body = await reader.ReadToEndAsync();
                        var bytes = Encoding.UTF8.GetBytes(body);
                        x.HttpContext.Request.Body = new MemoryStream(bytes);
                    };
                }
                
                var request = x.HttpContext.Request;
                var visitReqInfo = new VisitInfo(
                    Type: "Request",
                    VisitId: x.HttpContext.TraceIdentifier,
                    Brouser: request.Headers["User-Agent"],
                    Date: DateTime.Now,
                    ForwardedFor: ClickHouseHelper.GetForwardedIpOrNull(x.HttpContext),
                    RemoteAddr: x.HttpContext.Connection.RemoteIpAddress?.ToString() ?? "",
                    Referer: request.Headers["referer"],
                    Url: request.GetRawUrl(),
                    Method: request.Method,
                    Body: body,
                    StatusCode: null
                );

                ClickHouseTimer.NeedToInsertReqRes.Enqueue(visitReqInfo);

            });


            builderContext.AddResponseTransform(async x => {
                var body = "";
                var isImage = x.ProxyResponse?.Content?.Headers?.ContentType?.ToString()?.ToLower()?.Contains("image") ?? false;
                var hasAnyContent = (x.ProxyResponse?.Content?.Headers?.ContentLength ?? 0) > 0;

                var needToWriteBody = hasAnyContent && !isImage;

                if (needToWriteBody) {
                    var stream = await x.ProxyResponse.Content.ReadAsStreamAsync();
                    using (var reader = new StreamReader(stream)) {
                        body = await reader.ReadToEndAsync();

                        if (!string.IsNullOrEmpty(body)) {
                            x.SuppressResponseBody = false;
                            var bytes = Encoding.UTF8.GetBytes(body);
                            await x.HttpContext.Response.Body.WriteAsync(bytes);
                        }
                    }
                }

                var request = x.HttpContext.Request;

                var visitResInfo = new VisitInfo(
                    Type: "Response",
                    VisitId: x.HttpContext.TraceIdentifier,
                    Brouser: request.Headers["User-Agent"],
                    Date: DateTime.Now,
                    ForwardedFor: ClickHouseHelper.GetForwardedIpOrNull(x.HttpContext),
                    RemoteAddr: x.HttpContext.Connection.RemoteIpAddress?.ToString() ?? "",
                    Referer: request.Headers["referer"],
                    Url: request.GetRawUrl(),
                    Method: request.Method,
                    Body: body,
                    StatusCode: x.HttpContext.Response.StatusCode.ToString()
                );
                ClickHouseTimer.NeedToInsertReqRes.Enqueue(visitResInfo);
            });
        } catch (Exception ex) {
            var logger = builderContext.Services.GetRequiredService<ILogger<VisitInfo>>();
            logger.LogError(ex, "error when AddTransforms");
        }

        static async Task<string> getBodyAsString(Stream bodyStream) {
            StreamReader reader = new StreamReader(bodyStream);
            string bodyText = await reader.ReadToEndAsync();
            return bodyText;
        }
    });

var app = builder.Build();
app.MapReverseProxy();
app.Run();


