using Microsoft.AspNetCore.SignalR;

namespace SignalR.API.Hubs {
    public class DiscountHub : Hub {
        private readonly IDictionary<TerminalId, ConnectionParameter> _connections;
        public DiscountHub(IDictionary<TerminalId, ConnectionParameter> connections) {
            _connections = connections;
        }

        public async Task JoinZebraSignalRChannel(TerminalId terminalId) {

            if (!_connections.ContainsKey(terminalId)) {
                _connections.Add(terminalId, new ConnectionParameter());
                await Groups.AddToGroupAsync(Context.ConnectionId, terminalId.Name);
            }
        }


    }


    public record TerminalId(string Id) {
        public string Name => Id;
    };
    public class ConnectionParameter {
        public string TerminalId { get; set; }
    }
}
