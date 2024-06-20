using Microsoft.Extensions.Hosting;
using Microsoft.Extensions.Logging;

namespace MobileApi.Sources.Timers {

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
