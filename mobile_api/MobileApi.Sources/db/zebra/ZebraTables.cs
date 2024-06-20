using System.ComponentModel.DataAnnotations;
using System.ComponentModel.DataAnnotations.Schema;
using Microsoft.EntityFrameworkCore;
using Microsoft.Extensions.Caching.Memory;
using IndexAttribute = Microsoft.EntityFrameworkCore.IndexAttribute;

namespace MobileApi.Sources.db.zebra {

    #region scheme public

    [Table("checks", Schema = "public")]
    public class ChecksModel {
        [Column("id")]
        public int Id { get; set; }

        [Column("opened_at")]
        public DateTime CheckOpenedAt { get; set; }

        [Column("sum")]
        public decimal CheckSum { get; set; }

        [Column("shop_id")]
        public int ShopId { get; set; }

        [Column("worker_id")]
        public int WorkerId { get; set; }

        [Column("status")]
        public string Status { get; set; } //opened or closed

        [Column("payment")]
        public string Payment { get; set; }
        [Column("discount_percent")]
        public decimal? DiscountPercent { get; set; }

        [Column("mobile_user_id")]
        public string? MobileUserId { get; set; }
    }


    [Table("shops", Schema = "public")]
    public class ShopModel {
        [Column("id")]
        public int Id { get; set; }

        [Column("name")]
        public string Name { get; set; }
        [Column("address")]
        public string Address { get; set; }
    }

    [Table("sklads", Schema = "public")]
    public class SkladModel {
        [Column("id")]
        public int Id { get; set; }

        [Column("name")]
        public string Name { get; set; }

        [Column("shop_id")]
        public int ShopId { get; set; }

        [Column("deleted")]
        public bool Deleted { get; set; }
    }

    [Table("postavkas", Schema = "public")]
    public class PostavkaModel {
        [Column("id")]
        public int Id { get; set; }

        [Column("sklad_id")]
        public int SkladId { get; set; }

        [Column("schet_id")]
        public int SchetId { get; set; }

        [Column("time")]
        public DateTime PostavkaDate { get; set; }

        [Column("sum")]
        public decimal PostavkaSum { get; set; }

        [Column("deleted")]
        public bool Deleted { get; set; }
    }


    [Table("item_postavkas", Schema = "public")]
    public class ItemPostavkaModel {
        [Column("id")]
        public int Id { get; set; }

        [Column("item_id")]
        public int ItemId { get; set; }

        [Column("postavka_id")]
        public int PostavkaId { get; set; }
        [Column("type")]
        public string Type { get; set; }

        [Column("quantity")]
        public int Quantity { get; set; }

        [Column("cost")]
        public decimal Cost { get; set; }
    }

    [Table("schets", Schema = "public")]
    public class SchetModel {
        [Column("id")]
        public int Id { get; set; }

        [Column("name")]
        public string Name { get; set; }

        [Column("type")]
        public string Type { get; set; }
    }

    [Table("tovars", Schema = "public")]
    public class TovarModel {
        [Column("id")]
        public int Id { get; set; }

        [Column("name")]
        public string Name { get; set; }

        [Column("price")]
        public decimal Price { get; set; }

        [Column("measure")]
        public string Measure { get; set; }
        [Column("image")]
        public string Image { get; set; }
        [Column("category")]
        public int CategoryId { get; set; }
        [Column("shop_id")]
        public int ShopId { get; set; }
        [Column("discount")]
        public bool HasDiscount { get; set; }
    }

    [Table("category_tovars", Schema = "public")]
    public class CategoryTovarsModel {
        [Key]
        [Column("id")]
        public int Id { get; set; }
        [Column("name")]
        public string Name { get; set; }

    }

    [Table("ingredients", Schema = "public")]
    public class IngredientModel {
        [Key]
        [Column("id")]
        public int Id { get; set; }

        [Column("name")]
        public string Name { get; set; }

        [Column("measure")]
        public string Measure { get; set; }
        [Column("image")]
        public string Image { get; set; }
        [Column("ingredient_id")]
        public int IngredientId { get; set; }
        [Column("shop_id")]
        public int ShopId { get; set; }
    }


    [Keyless]
    [Table("ingredient_tech_carts", Schema = "public")]
    public class IngredientsTechcartModel {
        [Column("tech_cart_id")]
        public int TechCartId { get; set; }

        [Column("ingredient_id")]
        public int IngredientId { get; set; }

        [Column("brutto")]
        public decimal Brutto { get; set; }

    }

    [Table("tech_carts", Schema = "public")]
    public class TechcartModel {
        [Key]
        [Column("id")]
        public int Id { get; set; }

        [Column("name")]
        public string Name { get; set; }
        [Column("category")]
        public int CategoryId { get; set; }
        [Column("discount")]
        public bool HasDiscount { get; set; }
        [Column("deleted")]
        public bool IsDeleted { get; set; }

        [Column("measure")]
        public string Measure { get; set; }

        [Column("price")]
        public decimal Price { get; set; }
        [Column("shop_id")]
        public int ShopId { get; set; }
        [Column("tech_cart_id")]
        public int TechCartId { get; set; }
    }

    [Table("nabor_tech_carts", Schema = "public")]
    public class NaborTechCartsModel {
        [Column("tech_cart_id")]
        public int TechCartId { get; set; }
        [Column("nabor_id")]
        public int NaborId { get; set; }
        [Column("shop_id")]
        public int ShopId { get; set; }
    }

    [Table("ingredient_nabors", Schema = "public")]
    public class IngredientNaborsModel {
        [Column("nabor_id")]
        public int NaborId { get; set; }
        [Column("ingredient_id")]
        public int IngredientId { get; set; }
        [Column("price")]
        public decimal Price { get; set; }
        [Column("shop_id")]
        public int ShopId { get; set; }
    }

    [Table("nabor_masters", Schema = "public")]
    public class NaborModel {
        [Key]
        [Column("id")]
        public int Id { get; set; }
        [Column("name")]
        public string Name { get; set; }
        [Column("min")]
        public int Min { get; set; }
        [Column("max")]
        public int Max { get; set; }
    }


    [Table("check_tech_carts", Schema = "public")]
    public class CheckTechcartModel {
        [Key]
        [Column("id")]
        public int Id { get; set; }

        [Column("check_id")]
        public int CheckId { get; set; }

        [Column("tech_cart_id")]
        public int TechCartId { get; set; }

        [Column("tech_cart_name")]
        public string TechCartName { get; set; }

        [Column("quantity")]
        public int Quantity { get; set; }

        [Column("price")]
        public decimal Price { get; set; }

    }

    [Table("expence_ingredients", Schema = "public")]
    public class CheckExpenceIngredientModel {
        [Key]
        [Column("id")]
        public int Id { get; set; }

        [Column("sklad_id")]
        public int SkladId { get; set; }

        [Column("ingredient_id")]
        public int IngredientId { get; set; }

        [Column("quantity")]
        public int Quantity { get; set; }

        [Column("cost")]
        public decimal Cost { get; set; }

        [Column("type")]
        public string Type { get; set; }

        [Column("price")]
        public decimal Price { get; set; }

        [Column("check_tech_cart_id")]
        public int CheckTechCartId { get; set; }

        [Column("status")]
        public string Status { get; set; }
    }


    [Table("check_tovars", Schema = "public")]
    public class CheckTovarModel {
        [Key]
        [Column("id")]
        public int Id { get; set; }

        [Column("check_id")]
        public int CheckId { get; set; }

        [Column("tovar_id")]
        public int TovarId { get; set; }

        [Column("tovar_name")]
        public string TovarName { get; set; }

        [Column("quantity")]
        public int Quantity { get; set; }

        [Column("price")]
        public decimal Price { get; set; }

    }

    [Table("check_modificators", Schema = "public")]
    public class CheckModificatorModel {
        [Key]
        [Column("id")]
        public int Id { get; set; }

        [Column("check_id")]
        public int CheckId { get; set; }

        [Column("name")]
        public int Name { get; set; }

        [Column("quantity")]
        public int Quantity { get; set; }

    }

    [Table("remove_from_sklads", Schema = "public")]
    public class RemoveFromSkladModel {
        [Key]
        [Column("id")]
        public int Id { get; set; }

        [Column("sklad_id")]
        public int SkladId { get; set; }

        [Column("reason")]
        public string Reason { get; set; }

        [Column("cost")]
        public decimal Cost { get; set; }
        [Column("time")]
        public DateTime Time { get; set; }
        [Column("status")]
        public string Status { get; set; }
    }

    [Table("remove_from_sklad_items", Schema = "public")]
    public class RemoveFromSkladItemModel {
        [Column("id")]
        public int Id { get; set; }

        [Column("remove_id")]
        public int RemoveId { get; set; }

        [Column("item_id")]
        public int ItemId { get; set; }

        [Column("type")]
        public string Type { get; set; }

        [Column("quantity")]
        public int Quantity { get; set; }
    }


    #endregion

    #region scheme mobile_api

    [Table("mobile_users", Schema = "mobile_api")]
    [Index(nameof(Email))]
    [Index(nameof(Status))]
    public class MobileUserModel {

        [Key]
        [DatabaseGenerated(DatabaseGeneratedOption.Identity)]
        [Column("id")]
        public Guid Id { get; set; }

        [StringLength(200)]
        [Column("email")]
        public string? Email { get; set; }

        [Column("name")]
        public string Name { get; set; }

        [Column("birth_date")]
        public DateTime? BirthDate { get; set; }

        [Column("reg_date")]
        public DateTime RegDate { get; set; }

        [Column("zebra_coin_balance")]
        public decimal ZebraCoinBalance { get; set; }

        [Column("discount")]
        public decimal Discount { get; set; }


        [Column("removed_date")]
        public DateTime? RemovedDate { get; set; }


        [Column("status")]
        public UserStatus Status { get; set; }

        public static MobileUserModel? GetOrNull(Guid id, ZebraApplicationContext db, IMemoryCache memoryCache) {
            return memoryCache.GetOrCreate($"user_{id}", (options) => {
                options.AbsoluteExpirationRelativeToNow = TimeSpan.FromMinutes(40);

                return db.Users.Where(x => x.Status == MobileUserModel.UserStatus.Active).FirstOrDefault(x => x.Id == id);
            });
        }

        public static void CleanUserCache(Guid id, IMemoryCache memoryCache) {
            memoryCache.Remove($"user_{id}");
        }

        public static bool UserExistsWithEmail(string email, ZebraApplicationContext db) {
            var userExists = db.Users.Where(x => x.Status == MobileUserModel.UserStatus.Active).Count(x => x.Email == email) > 0;
            return userExists;
        }

        public static async Task RemoveUser(Guid id, ZebraApplicationContext db, IMemoryCache memoryCache) {
            var user = db.Users.First(x => x.Id == id);
            user.Status = UserStatus.Removed;
            user.RemovedDate = DateTime.UtcNow;
            await db.SaveChangesAsync();
            memoryCache.Remove($"user_{id}");
        }

        public static async Task<Guid> InsertIncognitoUser(ZebraApplicationContext db, string name) {
            var mobileUser = new MobileUserModel() {
                Name = name,
                Email = null,
                BirthDate = null,
                RegDate = DateTime.UtcNow,
                Discount = 5,
                ZebraCoinBalance = 0,
                Status = MobileUserModel.UserStatus.Active
            };
            await db.Users.AddAsync(mobileUser);
            await db.SaveChangesAsync();

            return mobileUser.Id;
        }
        public static async Task UpdateIncognitoUser(ZebraApplicationContext db, Guid userId, string name) {
            var user = db.Users.First(x => x.Id == userId);
            user.Name = name;
            await db.SaveChangesAsync();
        }
        public static async Task<Guid> InsertUser(ZebraApplicationContext db, MobileUserModel user) {
            var needToAddCoins = 10;
            user.ZebraCoinBalance += needToAddCoins;
            await db.Users.AddAsync(user);
            await db.CoinsTransactions.AddAsync(new CoinsTransactionModel {
                Date = user.RegDate,
                UserId = user.Id,
                ZebraCoins = needToAddCoins,
                TransactionType = "add",
                Note = "Регистрация"
            });
            await db.SaveChangesAsync();

            return user.Id;
        }

        public static async Task<(Guid Id, string Name)> GetUserId(ZebraApplicationContext db, string email) {
            var user = db.Users.Where(x => x.Status == MobileUserModel.UserStatus.Active).First(x => x.Email == email);
            return (user.Id, user.Name);
        }


        public enum UserStatus {
            Active,
            Removed
        }
    }


    [Table("user_feedback", Schema = "mobile_api")]
    [Index(nameof(ShopId))]
    [Index(nameof(WorkerId))]
    [Index(nameof(FeedbackDate))]
    [Index(nameof(ScoreQuality))]
    [Index(nameof(ScoreService))]
    public class FeedbackModel {
        [DatabaseGenerated(DatabaseGeneratedOption.Identity)]
        [Key, Column(Order = 0)]
        public int Id { get; set; }

        [Column("user_id")]
        public Guid UserId { get; set; }

        [Column("feedback_date")]
        public DateTime FeedbackDate { get; set; }

        [Column("check_id")]
        public int CheckId { get; set; }

        [Column("score_quality")]
        public double ScoreQuality { get; set; }

        [Column("score_service")]
        public double ScoreService { get; set; }

        [Column("feedback_text")]
        public string? FeedbackText { get; set; }

        [Column("shop_id")]
        public int ShopId { get; set; }

        [Column("worker_id")]
        public int WorkerId { get; set; }

        [Column("check_json")]
        public string CheckJson { get; set; }
    }

    [Table("calc_coins_transaction", Schema = "mobile_api")]
    public class CalcCoinsTransactionModel {
        [Key]
        [Column("id")]
        public Guid Id { get; set; }
        [Column("calc_date")]
        public DateTime CalcDate { get; set; }
    }


    [Table("coins_transaction", Schema = "mobile_api")]
    public class CoinsTransactionModel {
        [Key]
        [Column("id")]
        public Guid Id { get; set; }

        [Column("user_id")]
        public Guid UserId { get; set; }

        [Column("date")]
        public DateTime Date { get; set; }

        [Column("transaction_type")] //add, remove
        public string TransactionType { get; set; }

        [Column("zebra_coins")]
        public decimal ZebraCoins { get; set; }

        [Column("note")] //add, remove
        public string Note { get; set; }
    }

    #endregion

}
