using System.Collections.Concurrent;
using ClickHouse.Ado;

namespace ClickHouseProxy {
    public class ClickHouseTimer : BaseTimerService {

        public static ConcurrentQueue<VisitInfo> NeedToInsertReqRes = new ConcurrentQueue<VisitInfo>();
        private readonly ILogger<ClickHouseTimer> _logger;
        private readonly IConfiguration _configuration;

        public ClickHouseTimer(ILogger<ClickHouseTimer> logger, IConfiguration configuration) : base("ClickHouseTimer", new(Delay: TimeSpan.FromSeconds(30)), logger) {
            _logger = logger;
            _configuration = configuration;
        }

        public override async Task DoAction() {
            if (NeedToInsertReqRes.Count > 0) {
                try {
                    var clickHouseConnectionString = _configuration["clickHouseConnectionString"];

                    using (var clickHouseConnection = new ClickHouseConnection(clickHouseConnectionString)) {
                        clickHouseConnection.Open();
                        var rows = NeedToInsertReqRes.ToArray();

                        var command = clickHouseConnection.CreateCommand();
                        command.CommandText = "INSERT INTO clickhouse_logs(type, id, brouser, timestamp, forwardedFor, remoteAddr, referer, url, method, body, statusCode) VALUES @bulk";
                        command.Parameters.Add(new ClickHouseParameter {
                            ParameterName = "bulk",
                            Value = rows
                            .Select(v =>
                                new object[] {
                                        v.Type,
                                        v.VisitId,
                                        v.Brouser,
                                        v.Date.ToUniversalTime(),
                                        v.ForwardedFor,
                                        v.RemoteAddr,
                                        v.Referer,
                                        v.Url,
                                        v.Method,
                                        v.Body,
                                        v.StatusCode,
                                }
                            )
                        });

                        command.CommandTimeout = 3000;
                        command.ExecuteNonQuery();
                        NeedToInsertReqRes.Clear();
                    }

                } catch (Exception ex) {
                    _logger.LogError(ex, "error when DoAction in ClickHouseTimer");
                }
            }
        }
    }


    public abstract class BaseTimerService : BackgroundService {
        private readonly string _timerName;
        private readonly BaseTimerServiceSettings _delaySettings;
        private readonly ILogger _logger;

        public BaseTimerService(string timerName, BaseTimerServiceSettings delaySettings, ILogger logger) {
            _timerName = timerName;
            _logger = logger;
            _delaySettings = delaySettings;
        }

        public abstract Task DoAction();

        protected override async Task ExecuteAsync(CancellationToken stoppingToken) {
            _logger.LogDebug($"{_timerName} is starting.");

            stoppingToken.Register(() =>
                _logger.LogDebug($" {_timerName} background task is stopping."));

            while (!stoppingToken.IsCancellationRequested) {
                _logger.LogDebug($"{_timerName} task doing background work.");

                try {
                    await DoAction();
                    await Task.Delay(_delaySettings.Delay, stoppingToken);
                } catch (TaskCanceledException exception) {
                    _logger.LogCritical(exception, $" {_timerName} TaskCanceledException Error", exception.Message);
                }
            }

            _logger.LogDebug($"{_timerName} background task is stopping.");
        }
    }
    public record BaseTimerServiceSettings(TimeSpan Delay);

}
