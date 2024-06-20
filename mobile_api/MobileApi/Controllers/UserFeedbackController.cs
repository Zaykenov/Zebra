using Microsoft.AspNetCore.Mvc;
using Microsoft.Extensions.Caching.Memory;
using MobileApi.Sources.db.zebra;

namespace MobileApi.Controllers.UserFeedback {
    [ApiController]
    [Route("user-feedback")]
    public class UserFeedbackController : ControllerBase {

        private readonly IUserFeedbackService _service;

        public UserFeedbackController(IUserFeedbackService service) {
            _service = service;
        }

        [HttpPost]
        [Route("try-set-feedback")]
        public async Task<TrySetFeedbackRes> TrySetFeedback(UserFeedbackReq req) {
            return await _service.TrySetFeedback(req);
        }
    }

    public interface IUserFeedbackService {
        public Task<TrySetFeedbackRes> TrySetFeedback(UserFeedbackReq req);
    }

    public record UserFeedbackReq(
        int Checkid,
        Guid UserId,
        double ScoreQuality,
        double ScoreService,
        string FeedbackText,
        string CheckJson
    );

    public record TrySetFeedbackRes(TrySetFeedbackResStatus Status);

    [Newtonsoft.Json.JsonConverter(typeof(Newtonsoft.Json.Converters.StringEnumConverter))]
    [System.Text.Json.Serialization.JsonConverter(typeof(System.Text.Json.Serialization.JsonStringEnumConverter))]
    public enum TrySetFeedbackResStatus {
        UserNotExists = 1,
        Ok = 2,
    }

    public class UserFeedbackService : IUserFeedbackService {
        private readonly IMemoryCache _memoryCache;
        private readonly ZebraApplicationContext _db;

        public UserFeedbackService(IMemoryCache memoryCache, ZebraApplicationContext db) {
            _memoryCache = memoryCache;
            _db = db;
        }


        public async Task<TrySetFeedbackRes> TrySetFeedback(UserFeedbackReq req) {
            var user = MobileUserModel.GetOrNull(req.UserId, _db, _memoryCache);
            

            if (user == null) {
                return new TrySetFeedbackRes(TrySetFeedbackResStatus.UserNotExists);
            }

            var check = _db.Checks.First(x => x.Id == req.Checkid);

            var feedback = new FeedbackModel() {
                CheckId = req.Checkid,
                FeedbackDate = DateTime.UtcNow,
                UserId = req.UserId,
                ScoreQuality = req.ScoreQuality,
                ScoreService = req.ScoreService,
                FeedbackText = req.FeedbackText,
                ShopId = check.ShopId,
                WorkerId = check.WorkerId,
                CheckJson = req.CheckJson
            };
            await _db.Feedbacks.AddAsync(feedback);


            var needToAddCoins = 5;
            await _db.CoinsTransactions.AddAsync(new CoinsTransactionModel {
                Date = feedback.FeedbackDate,
                UserId = user.Id,
                ZebraCoins = needToAddCoins,
                TransactionType = "add",
                Note = "Комментарий"
            });

            var userFromDb = _db.Users.First(x => x.Id == req.UserId);
            userFromDb.ZebraCoinBalance += needToAddCoins;
            await _db.SaveChangesAsync();
            MobileUserModel.CleanUserCache(req.UserId, _memoryCache);

            return new TrySetFeedbackRes(TrySetFeedbackResStatus.Ok);
        }
    }
}
