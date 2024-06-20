using Microsoft.AspNetCore.Mvc;
using Microsoft.Extensions.Caching.Memory;
using MobileApi.Sources.db.zebra;
using System.Linq;
using System.Text;

namespace MobileApi.Controllers.UserOrders {
    [ApiController]
    [Route("user-orders")]
    public class UserOrdersController : ControllerBase {

        private readonly IUserOrdersService _service;

        public UserOrdersController(IUserOrdersService service) {
            _service = service;
        }

        [HttpGet]
        [Route("try-get-last-orders")]
        public async Task<TryGetLastOrdersRes> TryGetLastOrders([FromQuery] Guid userId) {
            return await _service.TryGetLastOrders(userId);
        }
    }

    public interface IUserOrdersService {
        public Task<TryGetLastOrdersRes> TryGetLastOrders(Guid userId);
    }

    public record TryGetLastOrdersRes(TryGetLastOrdersResStatus Status, LastOrder[] Orders);
    public record LastOrder(
        int Id,
        DateTime OpenedAt,
        decimal? Sum,
        string? Payment,

        Tovar[] Tovars,
        TechCart[] TechCarts,
        string? Comment,
        string? Feedback
    );
    public record Tovar(
        string Name,
        int Quantity,
        decimal Price
    );

    public record TechCart(
        string Name,
        int Quantity,
        decimal Price,
        string[] Modificators
    );


    [Newtonsoft.Json.JsonConverter(typeof(Newtonsoft.Json.Converters.StringEnumConverter))]
    [System.Text.Json.Serialization.JsonConverter(typeof(System.Text.Json.Serialization.JsonStringEnumConverter))]
    public enum TryGetLastOrdersResStatus {
        UserNotExists = 1,
        Ok = 2,
    }

    public class UserOrdersService : IUserOrdersService {
        private readonly IMemoryCache _memoryCache;
        private readonly ZebraApplicationContext _db;

        public UserOrdersService(IMemoryCache memoryCache, ZebraApplicationContext db) {
            _memoryCache = memoryCache;
            _db = db;
        }

        public async Task<TryGetLastOrdersRes> TryGetLastOrders(Guid userId) {
            var user = MobileUserModel.GetOrNull(userId, _db, _memoryCache);

            if (user == null) {
                return new TryGetLastOrdersRes(TryGetLastOrdersResStatus.UserNotExists, Orders: null);
            }

            var checkTovars = _db.Checks
                .Where(x =>
                        x.MobileUserId == userId.ToString()
                    && x.Status == "closed"
                    && x.CheckOpenedAt.AddHours(12).ToUniversalTime() > DateTime.Now.ToUniversalTime()
                )
                .Join(_db.CheckTovars,
                    check => check.Id,
                    tovar => tovar.CheckId,
                    (check, tovar) =>
                        new {
                            Check = check,
                            Tovar = tovar,
                        }
                    )
                .ToArray();

            var userTechCarts = _db.Checks
                .Where(x =>
                        x.MobileUserId == userId.ToString()
                    && x.Status == "closed"
                    && x.CheckOpenedAt.AddHours(12).ToUniversalTime() > DateTime.Now.ToUniversalTime()
                )
                .Join(_db.CheckTechcarts,
                    check => check.Id,
                    techcart => techcart.CheckId,
                    (check, techcart) =>
                        new {
                            Check = check,
                            Techcart = techcart,
                        }
                    );
            var userTechCartsModificators = _db.CheckExpenceIngredients
                .Where(x => x.Type == "modificator")
                .Join(userTechCarts,
                    checkExpenceIngredient => checkExpenceIngredient.CheckTechCartId,
                    userTechCart => userTechCart.Techcart.Id,
                    (checkExpenceIngredient, userTechCart) =>
                        new {
                            IngredientId = checkExpenceIngredient.IngredientId,
                            UserTechCart = userTechCart
                        }
                    )
                .Join(_db.Ingredients,
                    checkExpenceIngredient => checkExpenceIngredient.IngredientId,
                    ingredient => ingredient.Id,
                    (checkExpenceIngredient, ingredient) =>
                        new {
                            IngredientId = checkExpenceIngredient.IngredientId,
                            IngredientName = ingredient.Name,
                            UserTechCart = checkExpenceIngredient.UserTechCart
                        }
                    )
                .ToArray();

            var checkTechcarts = userTechCarts.ToArray().Select(userTechCart => new {
                Check = userTechCart.Check,
                Techcart = userTechCart.Techcart,
                Modificators = userTechCartsModificators.Where(x => x.UserTechCart.Techcart.Id == userTechCart.Techcart.Id).ToArray()
            }).ToArray();


            var allChecks = checkTovars.Select(x => x.Check).ToList().AddAndReturn(checkTechcarts.Select(x => x.Check).ToArray()).DistinctBy(x => x.Id).ToArray();

            var lastOrders = allChecks.Select(check => {
                return new LastOrder(
                    check.Id,
                    check.CheckOpenedAt,
                    check.CheckSum,
                    check.Payment,
                    checkTovars.Where(x => x.Check.Id == check.Id).Select(x => {
                        return new Tovar(x.Tovar.TovarName, x.Tovar.Quantity, x.Tovar.Price);
                    }).ToArray(),
                    checkTechcarts.Where(x => x.Check.Id == check.Id).Select(x => {
                        return new TechCart(x.Techcart.TechCartName, x.Techcart.Quantity, x.Techcart.Price, x.Modificators.Select(i => i.IngredientName).ToArray());
                    }).ToArray(),
                    Comment: null,
                    Feedback: null
                );
            }).ToArray();

            return _memoryCache.GetOrCreate($"last-orders-for-{userId}", (memory) => {
                memory.AbsoluteExpirationRelativeToNow = TimeSpan.FromSeconds(5);

                return new TryGetLastOrdersRes(TryGetLastOrdersResStatus.Ok, Orders: lastOrders.ToArray());
            })!;
        }


    }

    public static class ListEx {
        public static IEnumerable<T> AddAndReturn<T>(this IEnumerable<T> source, params T[] items) {
            return source.Concat(items);
        }
    }
}
