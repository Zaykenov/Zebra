using Microsoft.Extensions.Caching.Memory;
using SignalR.API.Hubs;

var builder = WebApplication.CreateBuilder(args);

// Add services to the container.

builder.Services.AddControllers();
// Learn more about configuring Swagger/OpenAPI at https://aka.ms/aspnetcore/swashbuckle
builder.Services.AddEndpointsApiExplorer();
builder.Services.AddSwaggerGen();

builder.Services.AddSignalR();

builder.Services.AddSingleton<IDictionary<TerminalId, ConnectionParameter>>(options => new Dictionary<TerminalId, ConnectionParameter>());
builder.Services.AddMemoryCache();


var app = builder.Build();

// Configure the HTTP request pipeline.
if (app.Environment.IsDevelopment()) {
    app.UseSwagger();
    app.UseSwaggerUI();
}

app.UseHttpsRedirection();

app.UseRouting();

if (app.Environment.IsDevelopment()) {
    app.UseCors(builder => {
        builder.WithOrigins("https://localhost:7268")
        .AllowAnyHeader().AllowAnyMethod().AllowCredentials();
    });
}


app.UseAuthorization();

app.MapControllers();

app.MapHub<DiscountHub>("/hub/discount");

app.Run();
