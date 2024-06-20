package repository

import (
	"database/sql"
	"zebra/model"
	"zebra/utils"

	"gorm.io/gorm"
)

type TerminalDB struct {
	db     *sql.DB
	gormDB *gorm.DB
}

func NewTerminalDB(db *sql.DB, gormDB *gorm.DB) *TerminalDB {
	return &TerminalDB{db: db, gormDB: gormDB}
}

func (r *TerminalDB) GetAllProducts(filter *model.Filter) ([]*model.Product, error) {
	products := []*model.Product{}

	productsQuery := `
	SELECT id, name, category_name, category, image, price, discount,
	CASE WHEN source_table = 'tovars' THEN 'tovar'
		 WHEN source_table = 'tech_carts' THEN 'techCart'
	END AS type
	FROM (
	SELECT tovars.tovar_id AS id, tovars.name, category_tovars.name AS category_name, tovars.category, tovars.image, tovars.price, tovars.discount, tovars.deleted, 'tovars' AS source_table
	FROM tovars 
	INNER JOIN category_tovars ON tovars.category = category_tovars.id
	WHERE tovars.deleted IS FALSE AND category_tovars.deleted IS FALSE AND tovars.shop_id = $1 AND tovars.deleted is false
	UNION
	SELECT tech_carts.tech_cart_id AS id, tech_carts.name, category_tovars.name AS category_name, tech_carts.category, tech_carts.image, tech_carts.price, tech_carts.discount,tech_carts.deleted, 'tech_carts' AS source_table
	FROM tech_carts 
	INNER JOIN category_tovars ON tech_carts.category = category_tovars.id
	WHERE tech_carts.deleted IS FALSE AND category_tovars.deleted IS FALSE AND tech_carts.shop_id = $1 AND tech_carts.deleted is false
	) AS combined_query
	ORDER BY category;`

	productsQueryRows, err := r.gormDB.Raw(productsQuery, filter.BindShop).Rows()
	if err != nil {
		return nil, productsQueryRows.Err()
	}
	defer productsQueryRows.Close()
	for productsQueryRows.Next() {
		product := &model.Product{}
		if productQueryRowErr := productsQueryRows.Scan(&product.ID, &product.Name, &product.CategoryName, &product.Category, &product.Image, &product.Price, &product.Discount, &product.Type); productQueryRowErr != nil {
			continue
		}
		if product.Type == utils.TypeTechCart {
			nabors := []*model.NaborInfo{}
			naborOutput := []*model.NaborOutput{}
			nab := r.gormDB.Model(&model.NaborTechCart{}).Select("nabor_tech_carts.nabor_id as id, nabors.name as name, nabors.min as min, nabors.max as max").Joins("inner join nabors on nabors.nabor_id = nabor_tech_carts.nabor_id and nabors.shop_id = ?", filter.BindShop).Where("nabor_tech_carts.tech_cart_id = ? and nabor_tech_carts.shop_id = ?", product.ID, filter.BindShop).Scan(&nabors)
			if nab.Error != nil {
				return nil, nab.Error
			}

			for _, nabor := range nabors {
				naborIngredients := []*model.IngredientOutput{}
				nabingre := r.gormDB.Model(&model.IngredientNabor{}).Select("ingredients.ingredient_id as id, ingredients.name, ingredients.category as category_id, category_ingredients.name as category, ingredient_nabors.nabor_id, ingredient_nabors.brutto, ingredient_nabors.price, ingredients.measure").Joins("inner join ingredients on ingredients.ingredient_id = ingredient_nabors.ingredient_id inner join category_ingredients on category_ingredients.id = ingredients.category").Where("nabor_id = ? and  ingredients.shop_id = ? and ingredient_nabors.shop_id = ?", nabor.ID, filter.BindShop, filter.BindShop).Scan(&naborIngredients)
				if nabingre.Error != nil {
					return nil, nabingre.Error
				}

				if naborIngredients == nil {
					naborIngredients = []*model.IngredientOutput{}
				}

				naborOutput = append(naborOutput, &model.NaborOutput{
					ID:              nabor.ID,
					Name:            nabor.Name,
					Min:             nabor.Min,
					Max:             nabor.Max,
					NaborIngredient: naborIngredients,
				})
			}
			if naborOutput == nil {
				naborOutput = []*model.NaborOutput{}
			}
			product.Nabor = naborOutput
		}
		products = append(products, product)
	}
	return products, nil
}
