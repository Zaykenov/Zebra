using Microsoft.EntityFrameworkCore;

namespace MobileApi.Sources.db.zebra {
    public class ZebraApplicationContext : DbContext {
        public DbSet<ChecksModel> Checks { get; set; }
        public DbSet<ShopModel> Shops { get; set; }
        public DbSet<SkladModel> Sklads { get; set; }
        public DbSet<PostavkaModel> Postavkas { get; set; }
        public DbSet<ItemPostavkaModel> ItemPostavkas { get; set; }
        public DbSet<SchetModel> Schets { get; set; }
        public DbSet<TovarModel> Tovars { get; set; }
        public DbSet<CategoryTovarsModel> CategoryTovars { get; set; }
        public DbSet<IngredientsTechcartModel> IngredientsTechcarts { get; set; }
        public DbSet<IngredientModel> Ingredients { get; set; }
        public DbSet<TechcartModel> Techcarts { get; set; }
        public DbSet<NaborTechCartsModel> NaborTechCarts { get; set; }
        public DbSet<CheckTechcartModel> CheckTechcarts { get; set; }
        public DbSet<CheckExpenceIngredientModel> CheckExpenceIngredients { get; set; }
        public DbSet<CheckTovarModel> CheckTovars { get; set; }
        public DbSet<CheckModificatorModel> CheckModificators { get; set; }
        public DbSet<RemoveFromSkladModel> RemoveFromSklads { get; set; }
        public DbSet<RemoveFromSkladItemModel> RemoveFromSkladItems { get; set; }
        public DbSet<MobileUserModel> Users { get; set; }
        public DbSet<FeedbackModel> Feedbacks { get; set; }
        public DbSet<CalcCoinsTransactionModel> CalcCoinsTransactions { get; set; }
        public DbSet<CoinsTransactionModel> CoinsTransactions { get; set; }
        public DbSet<NaborModel> NaborMasters { get; set; }
        public DbSet<IngredientNaborsModel> IngredientNabors { get; set; }

        public ZebraApplicationContext(DbContextOptions<ZebraApplicationContext> options) : base(options) {
        }

        protected override void OnModelCreating(ModelBuilder modelBuilder) {
            modelBuilder.Entity<ChecksModel>().ToTable(t => t.ExcludeFromMigrations());
            modelBuilder.Entity<ShopModel>().ToTable(t => t.ExcludeFromMigrations());
            modelBuilder.Entity<SkladModel>().ToTable(t => t.ExcludeFromMigrations());
            modelBuilder.Entity<PostavkaModel>().ToTable(t => t.ExcludeFromMigrations());
            modelBuilder.Entity<ItemPostavkaModel>().ToTable(t => t.ExcludeFromMigrations());
            modelBuilder.Entity<SchetModel>().ToTable(t => t.ExcludeFromMigrations());
            modelBuilder.Entity<TovarModel>().ToTable(t => t.ExcludeFromMigrations());
            modelBuilder.Entity<CategoryTovarsModel>().ToTable(t => t.ExcludeFromMigrations());
            modelBuilder.Entity<IngredientsTechcartModel>().ToTable(t => t.ExcludeFromMigrations());
            modelBuilder.Entity<IngredientModel>().ToTable(t => t.ExcludeFromMigrations());
            modelBuilder.Entity<TechcartModel>().ToTable(t => t.ExcludeFromMigrations());
            modelBuilder.Entity<CheckTechcartModel>().ToTable(t => t.ExcludeFromMigrations());
            modelBuilder.Entity<CheckExpenceIngredientModel>().ToTable(t => t.ExcludeFromMigrations());
            modelBuilder.Entity<CheckTovarModel>().ToTable(t => t.ExcludeFromMigrations());
            modelBuilder.Entity<CheckModificatorModel>().ToTable(t => t.ExcludeFromMigrations());
            modelBuilder.Entity<RemoveFromSkladModel>().ToTable(t => t.ExcludeFromMigrations());
            modelBuilder.Entity<RemoveFromSkladItemModel>().ToTable(t => t.ExcludeFromMigrations());
            modelBuilder.Entity<NaborTechCartsModel>().ToTable(t => t.ExcludeFromMigrations()).HasNoKey();
            modelBuilder.Entity<IngredientNaborsModel>().ToTable(t => t.ExcludeFromMigrations()).HasNoKey();
            modelBuilder.Entity<NaborModel>().ToTable(t => t.ExcludeFromMigrations());
        }
    }
}