using Microsoft.AspNetCore.Mvc;
using Microsoft.EntityFrameworkCore;
using MobileApi.Sources.db.zebra;

namespace MobileApi.Controllers {
    [ApiController]
    [Route("shop")]
    public class ShopController : ControllerBase {
        private readonly IShopService _shopService;
        public ShopController(IShopService shopService) {
            _shopService = shopService;
        }
        [HttpGet]
        [Route("get-all")]
        public async Task<TryGetShopsRes> TryGetShops() {
            var response = await _shopService.TryGetShops();
            return response;
        }

        [HttpGet]
        [Route("get-tech-carts/{shopId}")]
        public async Task<TryGetShopFullInfoRes> TryGetShopFullInfo(int shopId) {
            var response = await _shopService.TryGetShopFullInfo(shopId);
            return response;
        }

        [HttpGet]
        [Route("get-tovars/{shopId}")]
        public async Task<TryGetTovarsRes> TryGetTovars(int shopId) {
            var response = await _shopService.TryGetTovars(shopId);
            return response;
        }

    }

    public interface IShopService {
        Task<TryGetShopsRes> TryGetShops();
        Task<TryGetShopFullInfoRes> TryGetShopFullInfo(int shopId);
        Task<TryGetTovarsRes> TryGetTovars(int shopId);
    }

    public class ShopService : IShopService {
        private readonly ZebraApplicationContext _db;
        public ShopService(ZebraApplicationContext db) {
            _db = db;
        }

        public async Task<TryGetShopFullInfoRes> TryGetShopFullInfo(int shopId) {
            try {
                var categories = await _db.CategoryTovars.ToDictionaryAsync(m => m.Id);

                var shopItemInfos = new List<ShopItemInfo>();
                var techCarts = await _db.Techcarts.Where(x => x.ShopId == shopId && x.IsDeleted == false).ToArrayAsync();
                foreach (var techCart in techCarts) {
                    var naborIds = await _db.NaborTechCarts
                        .Where(x => x.TechCartId == techCart.TechCartId && x.ShopId == shopId)
                        .Select(x => x.NaborId)
                        .ToArrayAsync();

                    var naborResult = await _db.NaborMasters
                        .Where(x => naborIds.Contains(x.Id))
                        .ToArrayAsync();

                    var nabors = new List<NaborInfo>();

                    foreach (var nabor in naborResult) {
                        var ingredientIdsAndPrices = await _db.IngredientNabors
                            .Where(x => x.NaborId == nabor.Id && x.ShopId == shopId)
                            .Select(x => new { Id = x.IngredientId, Price = x.Price })
                            .ToArrayAsync();

                        var ingredientIds = ingredientIdsAndPrices.Select(x => x.Id).ToArray();

                        var ingredients = await _db.Ingredients
                            .Where(x => ingredientIds.Contains(x.IngredientId) && x.ShopId == shopId)
                            .Select(x => new ModificatorsInfo(x.IngredientId, x.Name, x.Image))
                            .ToArrayAsync();

                        var price = ingredientIdsAndPrices.First().Price;
                        nabors.Add(new NaborInfo(nabor.Id, nabor.Name, "No Description", nabor.Min, nabor.Max, price, ingredients));
                    }

                    var shopItemInfo = new ShopItemInfo(
                        techCart.Id,
                        techCart.Name,
                        categories[techCart.CategoryId].Name,
                        "",
                        techCart.Price,
                        techCart.HasDiscount,
                        nabors.ToArray()
                    );
                    shopItemInfos.Add(shopItemInfo);
                }

                var tovars = await _db.Tovars
                    .Where(x => x.ShopId == shopId)
                    .Select(x => new ShopItemInfo(x.Id, x.Name, categories[x.CategoryId].Name, x.Image, x.Price, x.HasDiscount, null))
                    .ToArrayAsync();

                shopItemInfos.AddRange(tovars);

                var shop = await _db.Shops.FirstOrDefaultAsync(x => x.Id == shopId);
                if (shop == null) {
                    return new TryGetShopFullInfoRes(null, $"Shop with id {shopId} does not exist", 404);
                }
                var data = new TryGetShopFullInfoModel(shopId, shop.Name, shop.Address, shopItemInfos.ToArray());
                return new TryGetShopFullInfoRes(data, "Success", 200);
            } catch (Exception ex) {
                //_logger.LogError(ex.Message);
                return new TryGetShopFullInfoRes(null, "Failed to retrieve data from database", 500);

            }
        }

        public async Task<TryGetShopsRes> TryGetShops() {
            try {
                var dbResult = ShopsDb.GetShopList();
                return new TryGetShopsRes(dbResult, "Success", 200);
            } catch (Exception ex) {
                //_logger.LogError(ex.Message);
                return new TryGetShopsRes(null, "Failed to retrieve data from database", 500);
            }
        }

        public async Task<TryGetTovarsRes> TryGetTovars(int shopId) {
            try {
                var categories = await _db.CategoryTovars.ToDictionaryAsync(m => m.Id);
                var tovars = await _db.Tovars
                    .Where(x => x.ShopId == shopId)
                    .Select(x => new TovarInfo(x.Id, x.Name, categories[x.CategoryId].Name, x.Price, x.HasDiscount))
                    .ToArrayAsync();

                return new TryGetTovarsRes(tovars, "Success", 200);
            } catch (Exception ex) {
                return new TryGetTovarsRes(null, "Failed to retrieve data from database", 500);
            }
        }
    }
    #region TryGetShopsReqRes
    public record TryGetShopsRes(ShopViewModel[] Data, string Message, int StatusCode);
    public record ShopViewModel(int ShopId, string Name, string Address, decimal Latitude, decimal Longitude);
    #endregion

    #region TryGetShopFullInfoReqRes
    public record TryGetShopFullInfoRes(TryGetShopFullInfoModel Data, string Message, int StatusCode);
    public record TryGetShopFullInfoModel(int ShopId, string Name, string Address, ShopItemInfo[] Items);
    public record ShopItemInfo(
        int ItemId,
        string ItemName,
        string Category,
        string ImageUrl,
        decimal Price,
        bool HasDiscount,
        NaborInfo[] Nabors
    );

    public record NaborInfo(int NaborId, string Name, string Description, int Min, int Max, decimal Price, ModificatorsInfo[] Modificators);
    public record ModificatorsInfo(int IngredientId, string IngredientName, string Image);
    #endregion

    #region TryGetTovarsReqRes
    public record TryGetTovarsRes(TovarInfo[] Data, string Message, int StatusCode);
    public record TovarInfo(int TovarId, string TovarName, string Category, decimal Price, bool HasDiscount);
    #endregion

    public static class ShopsDb {
        public static ShopViewModel[] GetShopList() {
            return new ShopViewModel[] {
                new ShopViewModel(1,"Азия Парк","Проспект Кабанбай батыра, 21",(decimal)51.12782950360762, (decimal)71.41261632565774),
                new ShopViewModel(2,"Ландмарк","Достык, 8",(decimal)51.12681393763807, (decimal)71.42175590526928),
                new ShopViewModel(4,"Ауезова","Ауезова 17", (decimal)51.17021060759292, (decimal)71.42398235153077),
                new ShopViewModel(5,"Астана молл","Тауелсиздик 34/1", (decimal)51.14200122686502, (decimal)71.46459456742544),
                new ShopViewModel(6,"Бахус","Сатпаева 14", (decimal)51.146470520635766, (decimal)71.474207604058),
                new ShopViewModel(7,"Австрия","Калдаякова 6", (decimal)51.1180982528877, (decimal)71.4596148079073),
                new ShopViewModel(8,"Австрия","Калдаякова 6", (decimal)51.1180982528877, (decimal)71.4596148079073),
                new ShopViewModel(9,"Мега","Кабанбай Батыра 62", (decimal)51.088951684041824, (decimal)71.41044955449095),
                new ShopViewModel(10,"Иманова","Иманова 3/1В", (decimal)51.16434694349917, (decimal)71.4307974594781),
                new ShopViewModel(11,"МФЦА","Проспект Мангилик Ел, 55/18", (decimal)51.089341510350785, (decimal)71.41928232380667),
                new ShopViewModel(12,"Sat City","Асфендиярова 1", (decimal)51.1311450877771, (decimal)71.37639620453102),
                new ShopViewModel(13,"City Lake","Сыганак 6", (decimal)51.1304875204193, (decimal)71.37137486874916),
                new ShopViewModel(14,"Экспо","Мангилик ел 55/20", (decimal)51.08927223109697, (decimal)71.41934345281035),
                new ShopViewModel(15,"Манхэтан","Мухамедханова, 6", (decimal)51.13972253613088, (decimal)71.38633537651387),
                new ShopViewModel(17,"Зеленый Квартал","E900, 2", (decimal)51.13217791516735, (decimal)71.39303175168493),
                new ShopViewModel(18,"Coffeetime 37","Мангилик Ел 37", (decimal)51.103754264933876, (decimal)71.42889579785364),
                new ShopViewModel(19,"Coffeetime 38","Мангилик Ел 38", (decimal)51.10487076590403, (decimal)71.43137404108032),
                new ShopViewModel(20,"Money 39","Мангилик Ел 39", (decimal)51.10273608291485, (decimal)71.42867084426405),
                new ShopViewModel(21,"Money 13","Кабанбай батыра 13", (decimal)51.13782898604746, (decimal)71.41553203923677),
            };
        }
    }
}

/*
 Data:{
	ShopInfo: {
		Id:123,
		Name: asda,
	},
	Items:[
		{
			"Id":123,
			"Name":"ASD",
			"CategoryID":12,
			"Description":"asd",
			"Image":,
			"Discount":true,
			"Nabors":[
				{
					"NaborID":32,
					"Name":"Moloko",
					"Description":"asd",
					"Min":1,
					"Max":1,
					"Price":150,
					"Modificators":[
						{
							"IngredientID":345,
							"IngredientName":"moloko",
							"Image":,
						}
					]
				}
			]
		},
		{
			"Id":123,
			"Name":"ASD",
			"CategoryID":12,
			"Description":"asd",
			"Image":,
			"Discount":true,
			"Nabors":[
				{
					"NaborID":32,
					"Name":"Moloko",
					"Description":"asd",
					"Min":1,
					"Max":1,
					"Price":150,
					"Modificators":[
						{
							"IngredientID":345,
							"IngredientName":"moloko",
							"Image":,
						}
					]
				}
			]
		}
	]
}

 */
