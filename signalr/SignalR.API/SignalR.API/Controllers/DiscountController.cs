using Microsoft.AspNetCore.Mvc;
using Microsoft.AspNetCore.SignalR;
using Microsoft.Extensions.Caching.Memory;
using SignalR.API.Hubs;

namespace SignalR.API.Controllers {
    [ApiController]
    [Route("api/v1/discount")]
    public class DiscountController : ControllerBase {
        private readonly IHubContext<DiscountHub> _hub;
        private readonly IDictionary<TerminalId, ConnectionParameter> _connections;
        private readonly IMemoryCache _memoryCache;

        public DiscountController(IHubContext<DiscountHub> hub, IDictionary<TerminalId, ConnectionParameter> connections, IMemoryCache memoryCache) {
            _hub = hub;
            _connections = connections;
            _memoryCache = memoryCache;
        }

        [HttpPost]
        [Route("nfc")]
        public async Task<TryPostNfcDiscountRes> TryPostNfcDiscountAsync([FromBody] PostNfcDiscountReq request) {
            try {
                var terminalId = getDeviceId(request.NfcContent);

                if(!_connections.ContainsKey(terminalId)) {
                    throw new InvalidOperationException("terminal is not connected!");
                }

                await _hub.Clients.Group(terminalId.Name).SendAsync("ReceiveUserMobileCode", request.UserMobileCode);
                var key = "InCacheStorage";
                if(_memoryCache.Get(key) == null) {
                    _memoryCache.GetOrCreate(key, f => new List<PostNfcDiscountReq>() { request });
                } else {
                    var list = _memoryCache.Get<List<PostNfcDiscountReq>>(key);
                    if(list == null) {
                        list = new List<PostNfcDiscountReq> { request };
                    } else {
                        list.Add(request);
                    }
                    _memoryCache.Set(key, list);
                }

                
                return new TryPostNfcDiscountRes(TryPostNfcDiscountResStatus.Ok, "Success");
            } catch(Exception ex) {
                return new TryPostNfcDiscountRes(TryPostNfcDiscountResStatus.Fail, ex.Message);
            }
        }

        [HttpGet]
        [Route("ping")]
        public async Task<List<PostNfcDiscountReq>> Ping() {
            var key = "InCacheStorage";
            var list = _memoryCache.Get<List<PostNfcDiscountReq>>(key);
            return list;
        }

        [HttpGet]
        [Route("connections")]
        public async Task<JsonResult> Connections() {
            return new JsonResult(_connections.Keys);
        }

        private TerminalId getDeviceId(string nfcContent) {
            return nfcContent switch {
                "63d27c7d-5b31-411b-96db-181f1a964381" => new TerminalId("10"),
                _ => throw new NotImplementedException()
            };
        }
    }
}
