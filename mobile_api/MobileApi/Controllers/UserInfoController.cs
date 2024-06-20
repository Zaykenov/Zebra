using Microsoft.AspNetCore.Mvc;
using Microsoft.Extensions.Caching.Memory;
using MobileApi.Sources.db.zebra;

namespace MobileApi.Controllers.UserInfo {
    [ApiController]
    [Route("user-info")]
    public class UserInfoController : ControllerBase {

        private readonly IUserInfoService _service;

        public UserInfoController(IUserInfoService service) {
            _service = service;
        }

        [HttpGet]
        [Route("try-get-info")]
        public async Task<TryGetUserInfoRes> TryGetInfo([FromQuery] Guid userId) {
            return await _service.TryGetInfo(userId);
        }
    }

    public interface IUserInfoService {
        public Task<TryGetUserInfoRes> TryGetInfo(Guid userId);
    }

    public record TryGetUserInfoRes(TryGetUserInfoStatus Status, UserInfo? Info);
    public record UserInfo(decimal ZebraCoinBalance, decimal Discount);

    [Newtonsoft.Json.JsonConverter(typeof(Newtonsoft.Json.Converters.StringEnumConverter))]
    [System.Text.Json.Serialization.JsonConverter(typeof(System.Text.Json.Serialization.JsonStringEnumConverter))]
    public enum TryGetUserInfoStatus {
        UserNotExists = 1,
        Ok = 2,
    }

    public class UserInfoService : IUserInfoService {
        private readonly IMemoryCache _memoryCache;
        private readonly ZebraApplicationContext _db;

        public UserInfoService(IMemoryCache memoryCache, ZebraApplicationContext db) {
            _memoryCache = memoryCache;
            _db = db;
        }

        public async Task<TryGetUserInfoRes> TryGetInfo(Guid userId) {
            var user = MobileUserModel.GetOrNull(userId, _db, _memoryCache);

            if (user == null) {
                return new TryGetUserInfoRes(TryGetUserInfoStatus.UserNotExists, Info: null);
            }

            return _memoryCache.GetOrCreate($"user-info-for-{userId}", (memory) => {
                memory.AbsoluteExpirationRelativeToNow = TimeSpan.FromMinutes(10);
                var discount = user.Discount / 100; //наш бэк без этого не работает((
                return new TryGetUserInfoRes(TryGetUserInfoStatus.Ok, new UserInfo(user.ZebraCoinBalance, discount));
            })!;
        }
    }
}
