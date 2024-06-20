using System.Text;
using OfficeOpenXml;

namespace Migrator
{
    public static class Programm
    {

        public static void Main(string[] args)
        {
            var ingredientsCategoriesPath = @"C:\Users\user\Desktop\FromMarzhan\export_ingredients_categories.xlsx";
            var ingredientsPath = @"C:\Users\user\Desktop\FromMarzhan\export_ingredients.xlsx";
            var techCartsPath = @"C:\Users\user\Desktop\FromMarzhan\export_dishes_221215.xlsx";
            var techCartsModificatorsPath = @"C:\Users\user\Desktop\FromMarzhan\modificatory_1.xlsx";
            var outputSqlPath = @"C:\Users\user\Desktop\total.sql";

            var (migrateIngredientsCategoriesSql, ingredientsCategories) = migrateIngredientsCategories(ingredientsCategoriesPath);
            var (migrateIngredientsSql, ingredients) = migrateIngredients(ingredientsPath, ingredientsCategories);
            var (migrateTechCartsSql, techCarts) = migrateTechCarts(techCartsPath, ingredients, getTechCartCategory());
            var migrateTechCartsModificatorsSql = migrateTechCartsModificators(techCartsModificatorsPath, techCarts, ingredients);

            var totalSql = new StringBuilder()
                .AppendLine("BEGIN;")
                
                .AppendLine("DELETE FROM public.nabor_tech_carts;")
                .AppendLine("DELETE FROM public.ingredient_nabors;")
                .AppendLine("DELETE FROM public.nabors;")
                .AppendLine("DELETE FROM public.ingredient_tech_carts;")
                .AppendLine("DELETE FROM public.tech_carts;")
                .AppendLine("DELETE FROM public.ingredients;")
                .AppendLine("DELETE FROM public.category_ingredients;")
                .AppendLine("---------------------------------------")
                .Append(migrateIngredientsCategoriesSql.SqlRaw)
                .AppendLine("---------------------------------------")
                .Append(migrateIngredientsSql.SqlRaw)
                .AppendLine("---------------------------------------")
                .Append(migrateTechCartsSql.SqlRaw)
                .AppendLine("---------------------------------------")
                .Append(migrateTechCartsModificatorsSql.SqlRaw)
                .ToString();

            File.WriteAllText(outputSqlPath, totalSql);

            Console.WriteLine($"done! Check {outputSqlPath}");
           // Console.ReadLine();
        }

        public static (Sql, IngredientCategory[]) migrateIngredientsCategories(string path)
        {
            var ingredientCategories = new List<IngredientCategory>();

            ExcelPackage.LicenseContext = LicenseContext.NonCommercial;
            FileInfo existingFile = new FileInfo(path);
            using (ExcelPackage package = new ExcelPackage(existingFile))
            {
                ExcelWorksheet worksheet = package.Workbook.Worksheets[0];
                int colCount = worksheet.Dimension.End.Column;
                int rowCount = worksheet.Dimension.End.Row;
                for (int row = 3; row <= rowCount; row++)
                {

                    var name = worksheet.Cells[row, 2].Value.ToString().Trim();
                    ingredientCategories.Add(new IngredientCategory(row + 100, name));
                   
                }
            }

            var lastIngredient = ingredientCategories.OrderByDescending(x => x.Id).First();
            ingredientCategories.Add(new IngredientCategory(lastIngredient.Id + 1, "Без категории"));

            var sql = new StringBuilder();
            sql.AppendLine($"SELECT setval('public.category_ingredients_id_seq', {ingredientCategories.OrderByDescending(x => x.Id).First().Id}, true);");

            foreach (var i in ingredientCategories)
            {
                sql.AppendLine($"INSERT INTO public.category_ingredients(id, name, image, deleted) VALUES ({i.Id}, '{i.Name}', '', false);");
            }

            return (new Sql(sql.ToString()), ingredientCategories.ToArray());
        }

        public static (Sql, Ingredient[]) migrateIngredients(string path, IngredientCategory[] categories)
        {
            var ingredients = new List<Ingredient>();

            ExcelPackage.LicenseContext = LicenseContext.NonCommercial;
            FileInfo existingFile = new FileInfo(path);
            using (ExcelPackage package = new ExcelPackage(existingFile))
            {
                ExcelWorksheet worksheet = package.Workbook.Worksheets[0];
                int colCount = worksheet.Dimension.End.Column;
                int rowCount = worksheet.Dimension.End.Row;
                for (int row = 3; row <= rowCount; row++)
                {
                    
                    var name = worksheet.Cells[row, 2].Value.ToString().Trim();
                    var categoryName = worksheet.Cells[row, 3].Value.ToString().Trim();
                   
                    var categoryId = categories.First(x => x.Name == categoryName).Id;

                    var measure = worksheet.Cells[row, 4].Value.ToString().Trim();
                    var cost = (double) worksheet.Cells[row, 11].Value;
                  
                    ingredients.Add(new Ingredient(row + 100, name, categoryId, measure, cost));
                }
            }

            var sql = new StringBuilder();
            sql.AppendLine($"SELECT setval('public.ingredients_id_seq', {ingredients.OrderByDescending(x => x.Id).First().Id}, true);");

            foreach (var i in ingredients)
            {
                sql.AppendLine($"INSERT INTO public.ingredients(id, name, category, image, measure, cost, deleted) VALUES ({i.Id}, '{i.Name}', {i.CategoryId}, '', '{i.Measure}', {i.Cost}, false);");
            }

            return (new Sql(sql.ToString()), ingredients.ToArray());
        }
         
        public static (Sql, TechCart[]) migrateTechCarts(string path, Ingredient[] ingredients, TechCartCategory[] techCartCategories)
        {
            var techCarts = new List<TechCart>();

            ExcelPackage.LicenseContext = LicenseContext.NonCommercial;
            FileInfo existingFile = new FileInfo(path);
            using (ExcelPackage package = new ExcelPackage(existingFile))
            {
                ExcelWorksheet worksheet = package.Workbook.Worksheets[0];
                int colCount = worksheet.Dimension.End.Column;
                int rowCount = worksheet.Dimension.End.Row;
                for (int row = 3; row <= rowCount; row++)
                {
                    
                    var name = worksheet.Cells[row, 3].Value.ToString().Trim();
                    var categoryName = worksheet.Cells[row, 4].Value.ToString().Trim();
                    //var tax = worksheet.Cells[row, 5].Value.ToString();
                    var tax = "Фискальный налог";
                    var measure = "шт.";
                    var cost = (double) worksheet.Cells[row, 6].Value;
                    var price = (double) worksheet.Cells[row, 7].Value;
                    var profit = price - cost;
                    var marginRaw = worksheet.Cells[row, 8];
                    double margin = marginRaw.Value == null ? 0 : (double) marginRaw.Value;
                  

                    var techCart = techCarts.FirstOrDefault(x => x.Name == name);
                    if(techCart == null){
                        var category = techCartCategories.First(x => x.Name == categoryName);
                        techCart = new TechCart(row + 100, name, category.Id, tax, measure, cost, price, profit, margin, new IngredientAndBrutto[0]);
                    }


                    var newIngredientsList = techCart.Ingredients.ToList();

                  
                    var ingredientName = worksheet.Cells[row, 10].Value?.ToString().Trim() ?? "";
                    if(!string.IsNullOrEmpty(ingredientName)){
                        var ingredient = ingredients.First(x => x.Name == ingredientName);
                        var ingredientBrutto = double.Parse(worksheet.Cells[row, 11].Value.ToString().Split(" ")[0]);
                        var ingredientBruttoMeasure = worksheet.Cells[row, 11].Value.ToString().Split(" ")[1].Trim();
                        
                        if(ingredientBruttoMeasure != ingredient.Measure){
                            if(ingredientBruttoMeasure == "г" && ingredient.Measure == "кг"){
                                ingredientBrutto = ingredientBrutto / 1000;
                            }else{
                                if(ingredientBruttoMeasure == "мл" && ingredient.Measure == "л"){
                                    ingredientBrutto = ingredientBrutto / 1000;
                                }else{
                                    throw new NotImplementedException();
                                }
                            }
                        }

                        newIngredientsList.Add(new IngredientAndBrutto(ingredient, ingredientBrutto));
                    }

                   
                    var newTechCart = techCart with { Ingredients = newIngredientsList.ToArray() };
                    techCarts = techCarts.Where(x => x.Name != newTechCart.Name).ToList();
                    techCarts.Add(newTechCart);
                }
            }

            var sql = new StringBuilder();
            sql.AppendLine($"SELECT setval('public.tech_carts_id_seq', {techCarts.OrderByDescending(x => x.Id).First().Id}, true);");

            foreach (var i in techCarts)
            {
                sql.AppendLine($"INSERT INTO public.tech_carts(id, name, category, image, tax, measure, cost, price, profit, margin, deleted) VALUES ({i.Id}, '{i.Name}', {i.CategoryId}, '', '{i.Tax}', '{i.Measure}', {i.Cost}, {i.Price}, {i.Profit}, {i.Margin}, false);");
            
            }
            sql.AppendLine().AppendLine();
            foreach (var i in techCarts)
            {
                foreach(var ingredient in i.Ingredients){
                    sql.AppendLine($"INSERT INTO public.ingredient_tech_carts(tech_cart_id, ingredient_id, brutto) VALUES ({i.Id}, {ingredient.Ingredient.Id}, {ingredient.Brutto});");
                }
            }

            return (new Sql(sql.ToString()), techCarts.ToArray());
        }
   
        public static Sql migrateTechCartsModificators(string path, TechCart[] techCarts, Ingredient[] ingredients)
        {
            var techCartNabors = new List<TechCartNabor>();
            var allNabors = getNabors(path, ingredients);

            ExcelPackage.LicenseContext = LicenseContext.NonCommercial;
            FileInfo existingFile = new FileInfo(path);
            using (ExcelPackage package = new ExcelPackage(existingFile))
            {
                ExcelWorksheet worksheet = package.Workbook.Worksheets[1];
                int colCount = worksheet.Dimension.End.Column;
                int rowCount = worksheet.Dimension.End.Row;
                for (int row = 3; row <= rowCount; row++)
                {
                    
                    var techCartName = worksheet.Cells[row, 2].Value.ToString().Trim();
                    var naborName = worksheet.Cells[row, 3].Value.ToString().Trim();

                    var techCart = techCarts.First(x => x.Name == techCartName);

                    var techCartNabor = techCartNabors.FirstOrDefault(x => x.TechCart.Id == techCart.Id);
                    if(techCartNabor == null){
                        techCartNabor = new TechCartNabor(techCart, new Nabor[0]);
                    }

                    var nabor = allNabors.First(x => x.Name == naborName);
                    
                    var naborList = techCartNabor.Nabors.Where(x => x.Name != nabor.Name).ToList();
                    naborList.Add(nabor);
                    var newTechCartNabor = techCartNabor with { Nabors = naborList.ToArray() };

                    
                    var techCartNaborsList = techCartNabors.Where(x => x.TechCart.Id != newTechCartNabor.TechCart.Id).ToList();
                    techCartNaborsList.Add(newTechCartNabor);

                    techCartNabors = techCartNaborsList;
                }
            }

            var allNaborIngredients = allNabors.SelectMany(x => x.Ingredients).ToArray(); 

            var sql = new StringBuilder();
            sql.AppendLine($"SELECT setval('public.nabors_id_seq', {allNabors.OrderByDescending(x => x.Id).First().Id}, true);");
            sql.AppendLine($"SELECT setval('public.ingredient_nabors_id_seq', {allNaborIngredients.OrderByDescending(x => x.Id).First().Id}, true);");
            



            foreach(var nabor in techCartNabors.SelectMany(x => x.Nabors).DistinctBy(x => x.Id)){
                sql.AppendLine($"INSERT INTO public.nabors(id, name, min, max, deleted) VALUES({nabor.Id}, '{nabor.Name}', 0, 0, false);");

                foreach(var ingredient in nabor.Ingredients){
                    sql.AppendLine($@"INSERT INTO public.ingredient_nabors(id, nabor_id, ingredient_id, name, image, measure, brutto, price) 
                    VALUES({ingredient.Id}, {nabor.Id}, {ingredient.IngredientId}, '', '', '', {ingredient.Brutto}, {ingredient.Price});");
                }
            }

            
            foreach (var techCartNabor in techCartNabors)
            {
                foreach(var nabor in techCartNabor.Nabors){
                    sql.AppendLine($"INSERT INTO public.nabor_tech_carts(tech_cart_id, nabor_id) VALUES ({techCartNabor.TechCart.Id}, {nabor.Id});");
                }
            }

            return new Sql(sql.ToString());
        }
    
        private static Nabor[] getNabors(string path, Ingredient[] ingredients){
            var nabors = new List<Nabor>();

            ExcelPackage.LicenseContext = LicenseContext.NonCommercial;
            FileInfo existingFile = new FileInfo(path);
            using (ExcelPackage package = new ExcelPackage(existingFile))
            {
                ExcelWorksheet worksheet = package.Workbook.Worksheets[0];
                int colCount = worksheet.Dimension.End.Column;
                int rowCount = worksheet.Dimension.End.Row;
                for (int row = 3; row <= rowCount; row++)
                {
                    
                    var naborName = worksheet.Cells[row, 2].Value.ToString().Trim();
                    var ingredientName = worksheet.Cells[row, 3].Value.ToString().Trim();
                    var ingredientBrutto = double.Parse(worksheet.Cells[row, 4].Value.ToString());
                    var ingredientPrice = double.Parse(worksheet.Cells[row, 5].Value.ToString());
                    var ingredient = ingredients.First(x => x.Name == ingredientName);


                    var nabor = nabors.FirstOrDefault(x => x.Name == naborName);
                    if(nabor == null){
                        nabor = new Nabor(row + 100, naborName, Min:0, Max: 0, new NaborIngredient[0]);
                    }

                    var naborIngredient = nabor.Ingredients.FirstOrDefault(x => x.IngredientId == ingredient.Id);
                    if(naborIngredient == null){
                        naborIngredient = new NaborIngredient(row + 100, ingredient.Id, ingredientBrutto, ingredientPrice);
                    }

                    var naborIngredientList = nabor.Ingredients.Where(x => x.IngredientId !=  naborIngredient.IngredientId).ToList();
                    naborIngredientList.Add(naborIngredient);
                    var newNabor = nabor with { Ingredients = naborIngredientList.ToArray() };
                    
                    var naborList = nabors.Where(x => x.Name != newNabor.Name).ToList();
                    naborList.Add(newNabor);

                    nabors = naborList;
                }

                return nabors.ToArray();
            }
        }

        public static TechCartCategory[] getTechCartCategory(){
            return new []{
                new TechCartCategory(1, "Молочные коктейли", "image.png"),
                new TechCartCategory(2, "Кофе", "image.png"),
                new TechCartCategory(3, "Лимонады", "image.png"),
                new TechCartCategory(4, "Горячие напитки", "image.png"),
                new TechCartCategory(5, "Еда", "image.png"),
                new TechCartCategory(6, "Десерты", "image.png"),
                new TechCartCategory(7, "Напитки", "image.png"),
                new TechCartCategory(8, "Холодный кофе", "image.png"),
                new TechCartCategory(9, "Фирменные напитки", "image.png"),
                new TechCartCategory(10, "Матча", "image.png"),
                new TechCartCategory(11, "Сезонное Меню", "image.png"),
                new TechCartCategory(12, "Эссенции", "image.png"),
                new TechCartCategory(13, "Главный экран", "image.png")
            };
        }

    }

    public record Sql(string SqlRaw);
    public record IngredientCategory(int Id, string Name);

    public record Ingredient(int Id, string Name, int CategoryId, string Measure, double Cost);
    public record IngredientAndBrutto(Ingredient Ingredient, double Brutto) ;
    public record TechCart(int Id, string Name, int CategoryId, string Tax, string Measure, double Cost, double Price, double Profit, double Margin, IngredientAndBrutto[] Ingredients);

    public record TechCartCategory(int Id, string Name, string Image);

    public record TechCartNabor(TechCart TechCart, Nabor[] Nabors);
    public record Nabor(int Id, string Name, int Min, int Max, NaborIngredient[] Ingredients);
    public record NaborIngredient(int Id, int IngredientId, double Brutto, double Price);
}