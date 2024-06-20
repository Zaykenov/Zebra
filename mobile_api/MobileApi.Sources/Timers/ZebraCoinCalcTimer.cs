using Microsoft.Extensions.Caching.Memory;
using Microsoft.Extensions.DependencyInjection;
using Microsoft.Extensions.Logging;
using MobileApi.Sources.db.zebra;

namespace MobileApi.Sources.Timers {
    public class ZebraCoinCalcTimer : BaseTimerService {
        private readonly IServiceScopeFactory _sp;
        private readonly IMemoryCache _memoryCache;
        public ZebraCoinCalcTimer(IServiceScopeFactory sp, IMemoryCache memoryCache, ILogger<ZebraCoinCalcTimer> logger) : base("ZebraCoinCalcTimer", new(Delay: TimeSpan.FromMinutes(1)), logger) {
            _sp = sp;
            _memoryCache = memoryCache;
        }

        public override async Task DoAction() {
            //await calcRegUsersCoins();
            //await calcCommentsCoins();

            using (var scope = _sp.CreateScope()) {
                var db = scope.ServiceProvider.GetRequiredService<ZebraApplicationContext>();

                var lastOperationDate = db.CalcCoinsTransactions.OrderByDescending(x => x.CalcDate).FirstOrDefault()?.CalcDate ?? DateTime.MinValue;

                var orderSum = 2500; //сумма чека, превыше которой нужно считать коины
                var checks = db.Checks
                    .Where(x =>
                        x.CheckOpenedAt > lastOperationDate && 
                        !string.IsNullOrEmpty(x.MobileUserId) &&
                        x.CheckSum > orderSum
                    )
                    .OrderBy(x => x.CheckOpenedAt)
                    .Take(30)
                    .ToArray();

                if(checks.Length == 0) {
                    return;
                }

                var operations = new List<CoinsTransactionModel>();
                var needToAddCoinsUserIdVsCoins = new Dictionary<Guid, decimal>();
                foreach (var check in checks) {
                    var needToAddZebraCoins = 10;
                    var userId = Guid.Parse(check.MobileUserId!);

                    operations.Add(new CoinsTransactionModel() {
                        Date = check.CheckOpenedAt,
                        TransactionType = "add",
                        ZebraCoins = needToAddZebraCoins,
                        UserId = userId,
                        Note = $"Заказ #{check.Id}"
                    });

                    if (!needToAddCoinsUserIdVsCoins.ContainsKey(userId)) {
                        needToAddCoinsUserIdVsCoins.Add(userId, needToAddZebraCoins);
                    } else {
                        needToAddCoinsUserIdVsCoins[userId] += needToAddZebraCoins;
                    }
                }

                foreach (var userIdVsNeedToAddCoins in needToAddCoinsUserIdVsCoins) {
                    var user = db.Users.FirstOrDefault(x => 
                        x.Status == MobileUserModel.UserStatus.Active && 
                        x.Id == userIdVsNeedToAddCoins.Key);

                    if (user == null || string.IsNullOrEmpty(user.Email)) {
                        operations = operations.Where(x => x.UserId != userIdVsNeedToAddCoins.Key).ToList();
                        continue;
                    }

                    user.ZebraCoinBalance += userIdVsNeedToAddCoins.Value;
                    MobileUserModel.CleanUserCache(user.Id, _memoryCache);
                }

                if(operations.Count > 0) {
                    db.CoinsTransactions.AddRange(operations);
                }

                db.CalcCoinsTransactions.RemoveRange(db.CalcCoinsTransactions.Where(x => x.CalcDate.Date < DateTime.UtcNow.Date));
                var newMigratedLastOperationDate = operations.OrderByDescending(x => x.Date).FirstOrDefault()?.Date ?? lastOperationDate;
                db.CalcCoinsTransactions.Add(new CalcCoinsTransactionModel { Id = Guid.NewGuid(), CalcDate = newMigratedLastOperationDate });

                await db.SaveChangesAsync();
            }

            Console.WriteLine(DateTime.Now.ToLongTimeString());
        }


        private async Task calcRegUsersCoins() {

            using (var scope = _sp.CreateScope()) {
                var db = scope.ServiceProvider.GetRequiredService<ZebraApplicationContext>();

                var regUsers = db.Users.Where(x => !string.IsNullOrEmpty(x.Email)).ToArray();
                   
                var operations = new List<CoinsTransactionModel>();
                var needToAddCoinsUserIdVsCoins = new Dictionary<Guid, decimal>();
                foreach (var regUser in regUsers) {
                    var needToAddZebraCoins = 10;
                    var userId = regUser.Id;

                    operations.Add(new CoinsTransactionModel() {
                        Date = regUser.RegDate,
                        ZebraCoins = needToAddZebraCoins,
                        UserId = userId,
                        TransactionType = "add",
                        Note = "Регистрация"
                    });

                    if (!needToAddCoinsUserIdVsCoins.ContainsKey(userId)) {
                        needToAddCoinsUserIdVsCoins.Add(userId, needToAddZebraCoins);
                    } else {
                        needToAddCoinsUserIdVsCoins[userId] += needToAddZebraCoins;
                    }
                }

                foreach (var userIdVsNeedToAddCoins in needToAddCoinsUserIdVsCoins) {
                    var user = db.Users.First(x => x.Id == userIdVsNeedToAddCoins.Key);
                    user.ZebraCoinBalance += userIdVsNeedToAddCoins.Value;
                }

                if (operations.Count > 0) {
                    db.CoinsTransactions.AddRange(operations);
                }

                await db.SaveChangesAsync();
            }
        }

        private async Task calcCommentsCoins() {
            using (var scope = _sp.CreateScope()) {
                var db = scope.ServiceProvider.GetRequiredService<ZebraApplicationContext>();

                var feedbacks = db.Feedbacks.ToArray();

                var operations = new List<CoinsTransactionModel>();
                var needToAddCoinsUserIdVsCoins = new Dictionary<Guid, decimal>();
                foreach (var feedback in feedbacks) {
                    var needToAddZebraCoins = 5;
                    var userId = feedback.UserId;

                    operations.Add(new CoinsTransactionModel() {
                        Date = feedback.FeedbackDate,
                        ZebraCoins = needToAddZebraCoins,
                        UserId = userId,
                        TransactionType = "add",
                        Note = "Комментарий"
                    });

                    if (!needToAddCoinsUserIdVsCoins.ContainsKey(userId)) {
                        needToAddCoinsUserIdVsCoins.Add(userId, needToAddZebraCoins);
                    } else {
                        needToAddCoinsUserIdVsCoins[userId] += needToAddZebraCoins;
                    }
                }

                foreach (var userIdVsNeedToAddCoins in needToAddCoinsUserIdVsCoins) {
                    var user = db.Users.FirstOrDefault(x => x.Id == userIdVsNeedToAddCoins.Key);

                    if (user == null || string.IsNullOrEmpty(user.Email)) {
                        operations = operations.Where(x => x.UserId != userIdVsNeedToAddCoins.Key).ToList();
                        continue;
                    }


                    user.ZebraCoinBalance += userIdVsNeedToAddCoins.Value;
                }

                if (operations.Count > 0) {
                    db.CoinsTransactions.AddRange(operations);
                }

                await db.SaveChangesAsync();
            }
        }
    }
}
