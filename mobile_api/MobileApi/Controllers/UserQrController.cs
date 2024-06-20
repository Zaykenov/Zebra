using Microsoft.AspNetCore.Mvc;
using Microsoft.Extensions.Caching.Memory;
using MobileApi.Sources.db.zebra;
using System.Text;

namespace MobileApi.Controllers.UserQr {
    [ApiController]
    [Route("user-qr")]
    public class UserQrController : ControllerBase {

        private readonly IUserQrService _service;

        public UserQrController(IUserQrService service) {
            _service = service;
        }

        [HttpGet]
        [Route("try-generate")]
        public async Task<TryGenerateQrRes> TryGenerateQr([FromQuery] Guid userId) {
            return await _service.TryGenerateQr(userId);
        }

        [HttpGet]
        [Route("try-get-user")]
        public async Task<TryGetUserByQrRes> TryGenerateQr([FromQuery] string qrContent) {
            return await _service.TryGetUserByQr(qrContent);
        }
    }

    public interface IUserQrService {
        public Task<TryGenerateQrRes> TryGenerateQr(Guid userId);
        public Task<TryGetUserByQrRes> TryGetUserByQr(string qrContent);
    }

    public record TryGenerateQrRes(TryGenerateQrResStatus Status, string QrContent, DateTime? Until);
    public record TryGetUserByQrRes(TryGetUserByQrResStatus Status, UserByQrModel User);
    public record UserByQrModel(
        Guid UserId,
        string Name,
        string? Email,
        decimal ZebraCoinBalance,
        decimal Discount
    );

    [Newtonsoft.Json.JsonConverter(typeof(Newtonsoft.Json.Converters.StringEnumConverter))]
    [System.Text.Json.Serialization.JsonConverter(typeof(System.Text.Json.Serialization.JsonStringEnumConverter))]
    public enum TryGenerateQrResStatus {
        UserNotExists = 1,
        Ok = 2,
    }


    [Newtonsoft.Json.JsonConverter(typeof(Newtonsoft.Json.Converters.StringEnumConverter))]
    [System.Text.Json.Serialization.JsonConverter(typeof(System.Text.Json.Serialization.JsonStringEnumConverter))]
    public enum TryGetUserByQrResStatus {
        QrNotFound = 1,
        Ok = 2,
    }

    public class UserQrService : IUserQrService {
        private readonly IMemoryCache _memoryCache;
        private readonly ZebraApplicationContext _db;

        public UserQrService(IMemoryCache memoryCache, ZebraApplicationContext db) {
            _memoryCache = memoryCache;
            _db = db;
        }

        public async Task<TryGenerateQrRes> TryGenerateQr(Guid userId) {
            var userOrNull = MobileUserModel.GetOrNull(userId, _db, _memoryCache);
            if (userOrNull == null) {
                return new TryGenerateQrRes(TryGenerateQrResStatus.UserNotExists, QrContent: null, Until: null);
            }

            var generateQrRes = _memoryCache.GetOrCreate($"qr-for-{userId}", (memory) => {
                memory.AbsoluteExpirationRelativeToNow = TimeSpan.FromMinutes(5);

                var code = randomNumbers(4);
                return new TryGenerateQrRes(TryGenerateQrResStatus.Ok, QrContent: code, Until: DateTime.Now.AddMinutes(5));
            })!;

            _memoryCache.Set($"qr_user_{generateQrRes.QrContent}", userOrNull, new MemoryCacheEntryOptions {
                AbsoluteExpiration = generateQrRes.Until
            });

            return generateQrRes;
        }

        public async Task<TryGetUserByQrRes> TryGetUserByQr(string qrContent) {
            var user = _memoryCache.Get<MobileUserModel>($"qr_user_{qrContent}");
            if(user == null) {
                return new TryGetUserByQrRes(TryGetUserByQrResStatus.QrNotFound, User: null);
            }

            return new TryGetUserByQrRes(TryGetUserByQrResStatus.Ok, new UserByQrModel(
                    user.Id,
                    user.Name,
                    user.Email,
                    user.ZebraCoinBalance,
                    user.Discount
            ));
        }


        private static string randomNumbers(int length) {
            var random = new Random();
            var sb = new StringBuilder();
            for (var i = 0; i < length; i++) {
                sb.Append(random.Next(0, 9));
            }

            return sb.ToString();
        }
    }
}
