package repository

import (
	"context"
	"database/sql"
	"errors"
	"time"
	"zebra/model"
	"zebra/utils"

	"gorm.io/gorm"
)

const dbTimeout = time.Second * 3

type TovarDB struct {
	db     *sql.DB
	gormDB *gorm.DB
}

func NewTovarDB(db *sql.DB, gormDB *gorm.DB) *TovarDB {
	return &TovarDB{db: db, gormDB: gormDB}
}

func (r *TovarDB) AddTovars(tovars []*model.Tovar) error {

	err := r.gormDB.Create(tovars).Error
	if err != nil {
		return err
	}
	for _, tovar := range tovars {
		sklad := &model.Sklad{}
		err = r.gormDB.Model(&model.Sklad{}).Where("shop_id = ? and deleted = ?", tovar.ShopID, false).First(sklad).Error
		if err != nil {
			return err
		}
		skladTovar := &model.SkladTovar{
			TovarID:  tovar.TovarID,
			SkladID:  sklad.ID,
			Quantity: 0,
			Cost:     tovar.Cost,
		}
		err = r.gormDB.Create(skladTovar).Error
		if err != nil {
			return err
		}
	}
	return nil
}

func (r *TovarDB) AddTovar(tovar *model.Tovar) error {
	/*ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	stmt := `insert into tovar (name, category, image, tax, measure, cost, price, profit, margin)
	values ($1, $2, $3, $4, $5, $6, $7, $8, $9)`
	err := r.db.QueryRowContext(ctx, stmt,
		tovar.Name,
		tovar.Category,
		tovar.Image,
		tovar.Tax,
		tovar.Measure,
		tovar.Cost,
		tovar.Price,
		tovar.Profit,
		tovar.Margin,
	)

	if err != nil {
		return err.Err()
	}

	return nil*/
	err := r.gormDB.Transaction(func(tx *gorm.DB) error {
		err := r.gormDB.Create(tovar).Error
		if err != nil {
			return err
		}
		sklads := []*model.Sklad{}
		err = r.gormDB.Model(&model.Sklad{}).Where("deleted = ?", false).Find(&sklads).Error
		if err != nil {
			return err
		}
		for _, sklad := range sklads {
			skladTovar := &model.SkladTovar{
				TovarID:  tovar.ID,
				SkladID:  sklad.ID,
				Quantity: 0,
				Cost:     tovar.Cost,
			}
			err = r.gormDB.Create(skladTovar).Error
			if err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}

func (r *TovarDB) GetAllTovar(filter *model.Filter) ([]*model.TovarOutput, int64, error) {
	/*ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()
	tovars := []*model.TovarOutput{}
	stmt := `select tovar.id, tovar.name, tovar.image, category_tovar.name, tovar.tax, tovar.measure, tovar.cost, tovar.price, tovar.profit, tovar.margin from tovar inner join category_tovar on tovar.category = category_tovar.id where tovar.deleted = $1`

	row, err := r.db.QueryContext(ctx, stmt, false)
	if err != nil {
		return nil, err
	}

	for row.Next() {
		tovar := &model.TovarOutput{}
		err := row.Scan(
			&tovar.ID,
			&tovar.Name,
			&tovar.Image,
			&tovar.Category,
			&tovar.Tax,
			&tovar.Measure,
			&tovar.Cost,
			&tovar.Price,
			&tovar.Profit,
			&tovar.Margin,
		)
		if err != nil {
			return nil, err
		}
		tovars = append(tovars, tovar)
	}

	if len(tovars) == 0 {
		tovars = []*model.TovarOutput{}
	}

	return tovars, nil*/
	tovars := []*model.TovarOutput{}
	res := r.gormDB
	if len(filter.SkladID) > 0 {
		sklads := []*model.Sklad{}
		err := r.gormDB.Model(&model.Sklad{}).Where("id IN (?) and deleted = ?", filter.SkladID, false).Scan(&sklads).Error
		if err != nil {
			return nil, 0, err
		}
		filter.SkladID = nil
		for _, sklad := range sklads {
			filter.Shop = append(filter.Shop, sklad.ShopID)
		}
	}
	if filter.Measure == "" {
		res = r.gormDB.Model(&model.Tovar{}).Select("tovars.id, tovars.tovar_id, tovars.shop_id, shops.name as shop_name, tovars.name, tovars.image, category_tovars.name as category, tovars.tax, tovars.measure, AVG(sklad_tovars.cost) as cost, tovars.price, tovars.discount").Joins("inner join category_tovars on tovars.category = category_tovars.id inner join sklad_tovars on tovars.tovar_id = sklad_tovars.tovar_id inner join sklads on  sklads.id = sklad_tovars.sklad_id inner join shops on shops.id = tovars.shop_id").Where("tovars.deleted = ? and sklads.shop_id IN (?)", false, filter.AccessibleShops).Group("tovars.id, tovars.tovar_id, shops.name, category_tovars.name")
		if res.Error != nil {
			return nil, 0, res.Error
		}
	} else {
		res = r.gormDB.Model(&model.Tovar{}).Select("tovars.id, tovars.tovar_id, tovars.shop_id, shops.name as shop_name, tovars.name, tovars.image, category_tovars.name as category, tovars.tax, tovars.measure, AVG(sklad_tovars.cost) as cost, tovars.price, tovars.discount").Joins("inner join category_tovars on tovars.category = category_tovars.id inner join sklad_tovars on tovars.tovar_id = sklad_tovars.tovar_id inner join sklads on  sklads.id = sklad_tovars.sklad_id inner join shops on shops.id = tovars.shop_id").Where("tovars.deleted = ? and sklads.shop_id IN (?) and tovars.measure = ?", false, filter.AccessibleShops, filter.Measure).Group("tovars.id, tovars.tovar_id, shops.name, category_tovars.name").Scan(&tovars)
		if res.Error != nil {
			return nil, 0, res.Error
		}
	}

	newRes, count, err := filter.FilterResults(res, &model.Tovar{}, utils.DefaultPageSize, "", "", "")

	if err != nil {
		return nil, 0, err
	}

	if newRes.Scan(&tovars).Error != nil {
		return nil, 0, newRes.Error
	}

	for _, tovar := range tovars {
		skladTovars := &model.SkladTovar{}
		res := r.gormDB.Model(&model.SkladTovar{}).Joins("inner join sklads on sklads.id = sklad_tovars.sklad_id").Where("sklad_tovars.tovar_id = ? and sklads.shop_id = ?", tovar.TovarID, tovar.ShopID).Scan(&skladTovars)
		if res.Error != nil {
			return nil, 0, res.Error
		}
		tovar.Cost = skladTovars.Cost
		// var sum float32 = 0
		// var quantity float32 = 0
		// for _, val := range skladTovars {
		// 	if val.Quantity <= 0 {
		// 		continue
		// 	}
		// 	sum += (val.Cost * val.Quantity)
		// 	quantity += val.Quantity
		// }
		// if quantity > 0 {
		// 	tovar.Cost = sum / quantity
		// } else {
		// 	var cost float32
		// 	res := r.gormDB.Model(&model.ItemPostavka{}).Select("cost").Joins("inner join postavkas on postavkas.id = item_postavkas.postavka_id").Where("item_id = ? and item_postavkas.type = ? and postavkas.deleted = ?", tovar.TovarID, utils.TypeTovar, false).Order("postavkas.time desc").First(&cost)
		// 	if res.Error != nil {
		// 		if res.Error != gorm.ErrRecordNotFound { //???
		// 			return nil, 0, res.Error
		// 		}
		// 	}
		// 	if res.RowsAffected == 0 {
		// 		deletedCost := r.gormDB.Model(&model.ItemPostavka{}).Select("cost").Joins("inner join postavkas on postavkas.id = item_postavkas.postavka_id").Where("item_id = ? and item_postavkas.type = ?", tovar.TovarID, utils.TypeTovar).Order("postavkas.time desc").First(&cost)
		// 		if deletedCost.Error != nil {
		// 			if deletedCost.Error != gorm.ErrRecordNotFound {
		// 				return nil, 0, deletedCost.Error
		// 			}
		// 		}
		// 		if deletedCost.RowsAffected == 0 {
		// 			res := r.gormDB.Model(&model.SkladTovar{}).Select("AVG(sklad_tovars.cost)").Joins("inner join sklads on sklads.id = sklad_tovars.sklad_id").Where("sklad_tovars.tovar_id = ? and sklads.shop_id IN (?)", tovar.TovarID, filter.AccessibleShops).Scan(&cost)
		// 			if res.Error != nil {
		// 				return nil, 0, res.Error
		// 			}
		// 		}
		// 	}
		// 	tovar.Cost = cost
		// }
		tovar.Profit = tovar.Price - tovar.Cost
		if tovar.Cost != 0 {
			tovar.Margin = tovar.Profit * 100 / tovar.Cost
		}
	}

	return tovars, count, nil

}

func (r *TovarDB) GetTovar(id int, filter *model.Filter) (*model.TovarOutput, error) {
	/*ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()
	tovar := &model.TovarOutput{}
	stmt := `select tovar.id, tovar.name, tovar.image, category_tovar.name, tovar.tax, tovar.measure, tovar.cost, tovar.price, tovar.profit, tovar.margin from tovar inner join category_tovar on tovar.category = category_tovar.id where tovar.id = $1`

	row := r.db.QueryRowContext(ctx, stmt, id)
	if row.Err() != nil {
		return nil, row.Err()
	}

	err := row.Scan(
		&tovar.ID,
		&tovar.Name,
		&tovar.Image,
		&tovar.Category,
		&tovar.Tax,
		&tovar.Measure,
		&tovar.Cost,
		&tovar.Price,
		&tovar.Profit,
		&tovar.Margin,
	)
	if err != nil {
		return nil, err
	}

	return tovar, nil*/

	tovar := &model.TovarOutput{}
	err := r.gormDB.Model(&model.Tovar{}).Select("tovars.id, tovars.tovar_id, tovars.shop_id, tovars.name, tovars.image, tovars.category as category_id,  category_tovars.name as category, tovars.tax, tovars.measure, AVG(sklad_tovars.cost) as cost, tovars.price, tovars.discount, tovars.is_visible").Joins("inner join category_tovars on tovars.category = category_tovars.id inner join sklad_tovars on tovars.tovar_id = sklad_tovars.tovar_id inner join sklads on  sklads.id = sklad_tovars.sklad_id").Where("tovars.id = ? and sklads.shop_id IN (?)", id, filter.AccessibleShops).Group("tovars.id, category_tovars.name").Scan(&tovar).Error
	if err != nil {
		return nil, err
	}

	skladTovars := &model.SkladTovar{}
	res := r.gormDB.Model(&model.SkladTovar{}).Joins("inner join sklads on sklads.id = sklad_tovars.sklad_id").Where("sklad_tovars.tovar_id = ? and sklads.shop_id = ?", tovar.TovarID, tovar.ShopID).Scan(&skladTovars)
	if res.Error != nil {
		return nil, res.Error
	}
	tovar.Cost = skladTovars.Cost
	// var sum float32 = 0
	// var quantity float32 = 0
	// for _, val := range skladTovars {
	// 	if val.Quantity <= 0 {
	// 		continue
	// 	}
	// 	sum += (val.Cost * val.Quantity)
	// 	quantity += val.Quantity
	// }
	// if quantity > 0 {
	// 	tovar.Cost = sum / quantity
	// } else {
	// 	var cost float32
	// 	res := r.gormDB.Model(&model.ItemPostavka{}).Select("cost").Joins("inner join postavkas on postavkas.id = item_postavkas.postavka_id").Where("item_id = ? and item_postavkas.type = ? and postavkas.deleted = ?", tovar.TovarID, utils.TypeTovar, false).Order("postavkas.time desc").First(&cost)
	// 	if res.Error != nil {
	// 		if res.Error != gorm.ErrRecordNotFound {
	// 			return nil, res.Error
	// 		}
	// 	}
	// 	if res.RowsAffected == 0 {
	// 		deletedCost := r.gormDB.Model(&model.ItemPostavka{}).Select("cost").Joins("inner join postavkas on postavkas.id = item_postavkas.postavka_id").Where("item_id = ? and item_postavkas.type = ?", tovar.TovarID, utils.TypeTovar).Order("postavkas.time desc").First(&cost)
	// 		if deletedCost.Error != nil {
	// 			if deletedCost.Error != gorm.ErrRecordNotFound {
	// 				return nil, deletedCost.Error
	// 			}
	// 		}
	// 		if deletedCost.RowsAffected == 0 {
	// 			res := r.gormDB.Model(&model.SkladTovar{}).Select("AVG(sklad_tovars.cost)").Joins("inner join sklads on sklads.id = sklad_tovars.sklad_id").Where("sklad_tovars.tovar_id = ? and sklads.shop_id IN (?)", tovar.TovarID, filter.AccessibleShops).Scan(&cost)
	// 			if res.Error != nil {
	// 				return nil, res.Error
	// 			}
	// 		}
	// 	}
	// 	tovar.Cost = cost
	// }
	tovar.Profit = tovar.Price - tovar.Cost
	if tovar.Cost != 0 {
		tovar.Margin = tovar.Profit * 100 / tovar.Cost
	}
	return tovar, nil
}

func (r *TovarDB) UpdateTovar(tovar *model.ReqTovar) error {
	/*ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()
	stmt := `update tovar set name = $2, category = $3, tax = $4, measure = $5, cost = $6, price = $7, profit = $8, margin = $9, deleted = $10, image = $11 where id=$1`

	row := r.db.QueryRowContext(ctx, stmt,
		tovar.ID,
		tovar.Name,
		tovar.Category,
		tovar.Tax,
		tovar.Measure,
		tovar.Cost,
		tovar.Price,
		tovar.Profit,
		tovar.Margin,
		tovar.Deleted,
		tovar.Image,
	)

	if row.Err() != nil {
		return row.Err()
	}

	return nil*/
	oldItem := &model.Tovar{}
	err := r.gormDB.Model(oldItem).Where("id = ?", tovar.ID).First(oldItem).Error
	if err != nil {
		return err
	}
	for _, shop := range tovar.ShopID {
		res := r.gormDB.Model(model.Tovar{}).Where("tovar_id = ? and shop_id = ?", oldItem.TovarID, shop).Updates(model.Tovar{Name: tovar.Name, Category: tovar.Category, Image: tovar.Image, Tax: tovar.Tax, Measure: tovar.Measure, Price: tovar.Price, Deleted: false, Discount: tovar.Discount})
		if res.Error != nil {
			return res.Error
		}
		if res.RowsAffected == 0 {
			res := r.gormDB.Model(model.Tovar{}).Create(&model.Tovar{TovarID: oldItem.TovarID, ShopID: shop, Name: tovar.Name, Category: tovar.Category, Image: tovar.Image, Tax: tovar.Tax, Measure: tovar.Measure, Price: tovar.Price, Discount: tovar.Discount})
			if res.Error != nil {
				return res.Error
			}
			sklad := &model.Sklad{}
			err = r.gormDB.Model(&model.Sklad{}).Where("shop_id = ? and deleted = ?", shop, false).First(sklad).Error
			if err != nil {
				return err
			}
			skladTovar := &model.SkladTovar{
				TovarID:  oldItem.TovarID,
				SkladID:  sklad.ID,
				Quantity: 0,
				Cost:     tovar.Cost,
			}
			err = r.gormDB.Create(skladTovar).Error
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (r *TovarDB) DeleteTovar(id int) error {
	/*ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()
	stmt := `update tovar set deleted = $2 where id=$1`
	row := r.db.QueryRowContext(ctx, stmt, id, true)

	if row.Err() != nil {
		return row.Err()
	}

	return nil*/
	tovar := &model.Tovar{}
	err := r.gormDB.Model(tovar).Where("id = ?", id).Update("deleted", true).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *TovarDB) AddCategoryTovar(category *model.CategoryTovar) error {
	/*ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	stmt := `insert into category_tovar (name, image)
	values ($1, $2)`
	err := r.db.QueryRowContext(ctx, stmt,
		category.Name,
		category.Image,
	)

	if err.Err() != nil {
		return err.Err()
	}

	return nil*/
	err := r.gormDB.Create(category).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *TovarDB) GetAllCategoryTovar(filter *model.Filter) ([]*model.CategoryTovar, int64, error) {
	/*ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()
	categories := []*model.Category{}
	stmt := `select * from category_tovar where deleted = $1`

	row, err := r.db.QueryContext(ctx, stmt, false)
	if err != nil {
		return nil, err
	}

	for row.Next() {
		category := &model.Category{}
		err := row.Scan(
			&category.ID,
			&category.Name,
			&category.Image,
			&category.Deleted,
		)
		if err != nil {
			return nil, err
		}
		categories = append(categories, category)
	}

	if len(categories) == 0 {
		categories = []*model.Category{}
	}

	return categories, nil*/
	categories := []*model.CategoryTovar{}
	res := r.gormDB.Where("deleted = ?", false).Order("id asc").Find(&categories).Scan(&categories)
	if res.Error != nil {
		return nil, 0, res.Error
	}

	if len(categories) == 0 {
		categories = []*model.CategoryTovar{}
	}
	var count int64
	err := res.Count(&count).Error
	if err != nil {
		return nil, 0, err
	}

	return categories, count, nil
}

func (r *TovarDB) GetCategoryTovar(id int) (*model.CategoryTovar, error) {
	/*ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()
	category := &model.Category{}
	stmt := `select * from category_tovar where id=$1`

	row := r.db.QueryRowContext(ctx, stmt, id)
	if row.Err() != nil {
		return nil, row.Err()
	}

	err := row.Scan(
		&category.ID,
		&category.Name,
		&category.Image,
		&category.Deleted,
	)
	if err != nil {
		return nil, err
	}

	return category, nil*/
	category := &model.CategoryTovar{}
	err := r.gormDB.Where("id = ?", id).Find(&category).Error
	if err != nil {
		return nil, err
	}
	return category, nil
}

func (r *TovarDB) UpdateCategoryTovar(category *model.CategoryTovar) error {
	/*ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()
	stmt := `update category_tovar set name = $2, deleted = $3 where id=$1`

	row := r.db.QueryRowContext(ctx, stmt,
		category.ID,
		category.Name,
		&category.Deleted,
	)

	if row.Err() != nil {
		return row.Err()
	}

	return nil*/
	err := r.gormDB.Save(category).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *TovarDB) DeleteCategoryTovar(id int) error {
	/*ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()
	stmt := `update category_tovars set deleted = $2 where id=$1`
	row := r.db.QueryRowContext(ctx, stmt, id, true)

	if row.Err() != nil {
		return row.Err()
	}

	return nil*/
	category := &model.CategoryTovar{}
	err := r.gormDB.Model(category).Where("id = ?", id).Update("deleted", true).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *TovarDB) AddTechCarts(techCarts []*model.TechCart) error {
	for _, techCart := range techCarts {
		ingredients := techCart.Ingredients
		techCart.Ingredients = nil
		nabors := techCart.Nabor
		techCart.Nabor = nil
		err := r.gormDB.Create(techCart).Error
		if err != nil {
			return err
		}
		for _, ingredient := range ingredients {
			ingredient.TechCartID = techCart.TechCartID
			err := r.gormDB.Create(ingredient).Error
			if err != nil {
				return err
			}
		}
		for _, nabor := range nabors {
			nabor.TechCartID = techCart.TechCartID
			err := r.gormDB.Create(nabor).Error
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (r *TovarDB) AddTechCart(tech *model.TechCart) error {
	/*ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	for _, v := range tech.Ingredients {
		ingredients := model.Ingredient{}
		err := r.gormDB.Model(ingredients).Select("*").Where("ingredient_id = ?", v.IngredientID).Find(&ingredients).Error
		if err != nil {
			return err
		}
		tech.Cost += ingredients.Cost * v.Brutto
	}

	stmt := `insert into tech_carts (name, category, image, tax, measure, cost, price, discount)
	values ($1, $2, $3, $4, $5, $6, $7, $8) returning id`
	row := tx.QueryRowContext(ctx, stmt,
		tech.Name,
		tech.Category,
		tech.Image,
		tech.Tax,
		tech.Measure,
		tech.Cost,
		tech.Price,
		tech.Discount,
	)
	if row.Err() != nil {
		if err := tx.Rollback(); err != nil {
			return err
		}
		return row.Err()
	}
	row.Scan(&tech.ID)
	for _, ingredient := range tech.Ingredients {
		query := `insert into ingredient_tech_carts (tech_cart_id, ingredient_id, brutto)
		values ($1, $2, $3)`
		row := tx.QueryRowContext(ctx, query,
			tech.ID,
			ingredient.IngredientID,
			ingredient.Brutto,
		)
		if row.Err() != nil {
			if err := tx.Rollback(); err != nil {
				return err
			}
			return row.Err()
		}
	}

	for _, nabor := range tech.Nabor {
		query := `insert into nabor_tech_carts (tech_cart_id, nabor_id)
		values ($1, $2)`
		row := tx.QueryRowContext(ctx, query,
			tech.ID,
			nabor.NaborID,
		)
		if row.Err() != nil {
			if err := tx.Rollback(); err != nil {
				return err
			}
			return row.Err()
		}
	}
	err = tx.Commit()
	if err != nil {
		return err
	}
	return nil
	/*err := r.gormDB.Transaction(func(tx *gorm.DB) error {
		for _, ingredient := range tech.Ingredients {
			ingredientOld := &model.Ingredient{}
			err := tx.Where("id = ?", ingredient.IngredientID).Find(&ingredientOld).Error
			if err != nil {
				return err
			}
			tech.Cost = tech.Cost + ingredientOld.Cost*ingredient.Brutto
		}
		err := tx.Create(tech).Error
		if err != nil {
			return err
		}
		return nil
	},
	)
	if err != nil {
		return err
	}
	return nil*/
	return nil
}

func (r *TovarDB) GetAllTechCart(filter *model.Filter) ([]*model.TechCartResponse, int64, error) {
	/*ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()
	techs := []*model.TechCartOutput{}
	stmt := `select tech_carts.id, tech_carts.name, category_tovars.name, tech_carts.image, tech_carts.tax, tech_carts.measure, tech_carts.cost, tech_carts.price, tech_carts.profit, tech_carts.margin, (select json_agg(json_build_object('id',ingredients.id,'name',ingredients.name,'brutto', ingredient_tech_carts.brutto,'cost', ingredients.cost)) from ingredient_tech_carts inner join ingredients on ingredient_tech_carts.ingredient_id = ingredients.id where tech_cart_id = tech_carts.id) from tech_carts inner join category_tovars on tech_carts.category = category_tovars.id where tech_carts.deleted = $1 group by tech_carts.id,category_tovars.name `

	//select nabor.id, nabor.name, nabor.min, nabor.max, (select json_agg(json_build_object('id',ingredient.id,'name', ingredient.name,'cost', ingredient.cost,'measure', ingredient.measure))from ingredient_nabors inner join ingredient on ingredient_nabors.ingredient_id = ingredient.id where ingredient_nabors.nabor_id = nabor.id ) from nabor_tech_carts inner join nabor on nabor_tech_carts.nabor_id = nabor.id where nabor_tech_carts.tech_cart_id = 4
	row, err := r.db.QueryContext(ctx, stmt, false)
	if err != nil {
		return nil, err
	}

	for row.Next() {
		var ingre *model.IngredientNumOuputArray
		tech := &model.TechCartOutput{}
		err := row.Scan(
			&tech.ID,
			&tech.Name,
			&tech.Category,
			&tech.Image,
			&tech.Tax,
			&tech.Measure,
			&tech.Cost,
			&tech.Price,
			&tech.Profit,
			&tech.Margin,
			&ingre,
		)
		if ingre != nil {
			tech.Ingredients = *ingre
		}
		if err != nil {
			return nil, err
		}
		tech.Cost = 0
		for i := 0; i < len(tech.Ingredients); i++ {
			tech.Ingredients[i].Netto = tech.Ingredients[i].Brutto
			tech.Cost = tech.Ingredients[i].Cost*tech.Ingredients[i].Brutto + tech.Cost
		}
		tech.Profit = tech.Price - tech.Cost
		tech.Margin = int(tech.Price/tech.Cost*100) - 100
		if len(tech.Ingredients) == 0 {
			tech.Ingredients = []*model.IngredientNumOutput{}
		}
		query := `select nabors.id, nabors.name, nabors.min, nabors.max, (select json_agg(json_build_object('id',ingredients.id,'name', ingredients.name,'cost', ingredients.cost,'image', ingredients.image,'measure', ingredients.measure, 'brutto', ingredient_nabors.brutto, 'price', ingredient_nabors.price)) from ingredient_nabors inner join ingredients on ingredient_nabors.ingredient_id = ingredients.id where ingredient_nabors.nabor_id = nabors.id ) from nabor_tech_carts inner join nabors on nabor_tech_carts.nabor_id = nabors.id where nabor_tech_carts.tech_cart_id = $1`
		naborQuery, err := r.db.QueryContext(ctx, query, tech.ID)
		if err != nil {
			return nil, err
		}
		for naborQuery.Next() {
			nabor := &model.NaborOutput{}
			ingreNabor := &model.IngredientNaborOutput{}

			err := naborQuery.Scan(
				&nabor.ID,
				&nabor.Name,
				&nabor.Min,
				&nabor.Max,
				&ingreNabor,
			)
			if err != nil {
				return nil, err
			}

			if ingreNabor != nil {
				nabor.NaborIngredient = *ingreNabor
			}

			if len(nabor.NaborIngredient) <= 0 {
				nabor.NaborIngredient = []*model.IngredientOutput{}
			}
			tech.Nabors = append(tech.Nabors, nabor)
		}
		if len(tech.Nabors) <= 0 {
			tech.Nabors = []*model.NaborOutput{}
		}
		techs = append(techs, tech)

	}

	if len(techs) == 0 {
		techs = []*model.TechCartOutput{}
	}
	return techs, nil*/
	if len(filter.SkladID) > 0 {
		sklads := []*model.Sklad{}
		err := r.gormDB.Model(&model.Sklad{}).Where("id IN (?) and deleted = ?", filter.SkladID, false).Scan(&sklads).Error
		if err != nil {
			return nil, 0, err
		}
		filter.SkladID = nil
		for _, sklad := range sklads {
			filter.Shop = append(filter.Shop, sklad.ShopID)
		}
	}
	techs := []*model.TechCartResponse{}
	res := r.gormDB.Model(&model.TechCart{}).Select("tech_carts.id, tech_carts.tech_cart_id, tech_carts.shop_id, shops.name as shop_name, tech_carts.name, category_tovars.name as category,category_tovars.id as category_id, tech_carts.image, tech_carts.tax, tech_carts.measure, tech_carts.price, tech_carts.discount").Joins("inner join category_tovars on tech_carts.category = category_tovars.id inner join shops on shops.id = tech_carts.shop_id").Where("tech_carts.deleted = ?", false)
	if res.Error != nil {
		return nil, 0, res.Error
	}

	newRes, count, err := filter.FilterResults(res, model.TechCart{}, utils.DefaultPageSize, "", "", "")
	if err != nil {
		return nil, 0, err
	}
	if newRes.Scan(&techs).Error != nil {
		return nil, 0, newRes.Error
	}

	for _, tech := range techs {
		techIngredients := []*model.IngredientTechCart{}
		res := r.gormDB.Model(&model.IngredientTechCart{}).Select("*").Where("ingredient_tech_carts.tech_cart_id = ? and shop_id = ?", tech.TechCartID, tech.ShopID).Scan(&techIngredients)
		if res.Error != nil {
			return nil, 0, res.Error
		}

		var ingredientCost float32 = 0
		var ingredientSum float32 = 0
		for _, item := range techIngredients {
			skladIngredient := &model.SkladIngredient{}
			res := r.gormDB.Model(&model.SkladIngredient{}).Joins("inner join sklads on sklads.id = sklad_ingredients.sklad_id").Where("sklad_ingredients.ingredient_id = ? and sklads.shop_id = ?", item.IngredientID, tech.ShopID).Scan(&skladIngredient)
			if res.Error != nil {
				return nil, 0, res.Error
			}
			ingredientCost = skladIngredient.Cost
			// var sum float32 = 0
			// var quantity float32 = 0
			// for _, val := range skladIngredient {
			// 	if val.Quantity <= 0 {
			// 		continue
			// 	}
			// 	sum += (val.Cost * val.Quantity)
			// 	quantity += val.Quantity
			// }
			// if quantity > 0 {
			// 	ingredientCost = sum / quantity
			// } else {
			// 	var cost float32
			// 	res := r.gormDB.Model(&model.ItemPostavka{}).Select("cost").Joins("inner join postavkas on postavkas.id = item_postavkas.postavka_id").Where("item_id = ? and item_postavkas.type = ? and postavkas.deleted = ?", item.IngredientID, utils.TypeIngredient, false).Order("postavkas.time desc").First(&cost)
			// 	if res.Error != nil {
			// 		if res.Error != gorm.ErrRecordNotFound {
			// 			return nil, 0, res.Error
			// 		}
			// 	}
			// 	if res.RowsAffected == 0 {
			// 		deletedCost := r.gormDB.Model(&model.ItemPostavka{}).Select("cost").Joins("inner join postavkas on postavkas.id = item_postavkas.postavka_id").Where("item_id = ? and item_postavkas.type = ?", item.IngredientID, utils.TypeIngredient).Order("postavkas.time desc").First(&cost)
			// 		if deletedCost.Error != nil {
			// 			if deletedCost.Error != gorm.ErrRecordNotFound {
			// 				return nil, 0, deletedCost.Error
			// 			}
			// 		}
			// 		if deletedCost.RowsAffected == 0 {
			// 			res := r.gormDB.Model(&model.SkladIngredient{}).Select("COALESCE(AVG(sklad_ingredients.cost), 0.0)").Joins("inner join sklads on sklads.id = sklad_ingredients.sklad_id").Where("sklad_ingredients.ingredient_id = ? and sklads.shop_id IN (?)", item.IngredientID, filter.AccessibleShops).Scan(&cost)
			// 			if res.Error != nil {
			// 				return nil, 0, res.Error
			// 			}
			// 		}
			// 	}
			// 	ingredientCost = cost
			// }
			ingredientSum += ingredientCost * item.Brutto
		}
		tech.Cost = ingredientSum
		tech.Profit = tech.Price - tech.Cost
		if tech.Cost != 0 {
			tech.Margin = tech.Profit * 100 / tech.Cost
		}
	}
	return techs, count, nil
}

func (r *TovarDB) GetTechCart(id int, filter *model.Filter) (*model.TechCartOutput, error) {
	/*ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()
	tech := &model.TechCartOutput{}
	stmt := `select tech_carts.id, tech_carts.name, category_tovars.name, tech_carts.category, tech_carts.image, tech_carts.tax, tech_carts.measure, tech_carts.cost, tech_carts.price, tech_carts.profit, tech_carts.margin, tech_carts.discount from tech_carts inner join category_tovars on tech_carts.category = category_tovars.id where tech_carts.id=$1`
	row := r.db.QueryRowContext(ctx, stmt, id)
	if row.Err() != nil {
		return nil, row.Err()
	}

	err := row.Scan(
		&tech.ID,
		&tech.Name,
		&tech.Category,
		&tech.CategoryID,
		&tech.Image,
		&tech.Tax,
		&tech.Measure,
		&tech.Cost,
		&tech.Price,
		&tech.Profit,
		&tech.Margin,
		&tech.Discount,
	)
	tech.Cost = 0
	if err != nil {
		return nil, err
	}
	query := `select ingredients.id, ingredients.name, ingredients.measure, ingredients.image, ingredient_tech_carts.brutto, sklad_ingredients.cost from ingredient_tech_carts inner join ingredients on ingredient_tech_carts.ingredient_id = ingredients.id where tech_cart_id = $1`
	rows, err := r.db.QueryContext(ctx, query, tech.ID)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		ingredient := &model.IngredientNumOutput{}
		err := rows.Scan(
			&ingredient.ID,
			&ingredient.Name,
			&ingredient.Measure,
			&ingredient.Image,
			&ingredient.Brutto,
			&ingredient.Cost,
		)
		if err != nil {
			return nil, err
		}
		ingredient.Netto = ingredient.Brutto
		ingredient.Cost = ingredient.Brutto * ingredient.Cost
		tech.Cost = ingredient.Cost + tech.Cost

		tech.Ingredients = append(tech.Ingredients, ingredient)
	}
	query = `select nabors.id, nabors.name, nabors.min, nabors.max, (select json_agg(json_build_object('id',ingredients.id,'name', ingredients.name,'cost', ingredients.cost,'image', ingredients.image,'measure', ingredients.measure, 'brutto', ingredient_nabors.brutto, 'price', ingredient_nabors.price)) from ingredient_nabors inner join ingredients on ingredient_nabors.ingredient_id = ingredients.id where ingredient_nabors.nabor_id = nabors.id ) from nabor_tech_carts inner join nabors on nabor_tech_carts.nabor_id = nabors.id where nabor_tech_carts.tech_cart_id = $1`
	naborQuery, err := r.db.QueryContext(ctx, query, tech.ID)
	if err != nil {
		return nil, err
	}
	for naborQuery.Next() {
		nabor := &model.NaborOutput{}
		ingreNabor := &model.IngredientNaborOutput{}

		err := naborQuery.Scan(
			&nabor.ID,
			&nabor.Name,
			&nabor.Min,
			&nabor.Max,
			&ingreNabor,
		)
		if err != nil {
			return nil, err
		}

		if ingreNabor != nil {
			nabor.NaborIngredient = *ingreNabor
		}

		if len(nabor.NaborIngredient) <= 0 {
			nabor.NaborIngredient = []*model.IngredientOutput{}
		}
		tech.Nabors = append(tech.Nabors, nabor)
	}
	if len(tech.Nabors) <= 0 {
		tech.Nabors = []*model.NaborOutput{}
	}
	tech.Profit = tech.Price - tech.Cost
	tech.Margin = tech.Price/tech.Cost*100 - 100

	return tech, nil*/
	//techInfo := &model.TechCartInfo{}
	tech := &model.TechCartOutput{}
	res := r.gormDB.Model(&model.TechCart{}).Select("tech_carts.id, tech_carts.tech_cart_id, tech_carts.shop_id, tech_carts.name, category_tovars.name as category,category_tovars.id as category_id, tech_carts.image, tech_carts.tax, tech_carts.measure, tech_carts.price, tech_carts.discount").Joins("inner join category_tovars on tech_carts.category = category_tovars.id").Where("tech_carts.id = ?", id).Scan(&tech)
	if res.Error != nil {
		return nil, res.Error
	}

	ingredients := []*model.IngredientNumOutput{}
	ingre := r.gormDB.Model(&model.IngredientTechCart{}).Select("ingredients.ingredient_id as id, ingredient_tech_carts.ingredient_id, ingredients.name, ingredients.measure, ingredients.image, ingredient_tech_carts.brutto, AVG(sklad_ingredients.cost) as cost").Joins("inner join ingredients on ingredients.ingredient_id = ingredient_tech_carts.ingredient_id inner join sklad_ingredients on ingredients.ingredient_id = sklad_ingredients.ingredient_id inner join sklads on  sklads.id = sklad_ingredients.sklad_id").Where("ingredient_tech_carts.tech_cart_id = ? and sklads.shop_id IN (?) and ingredients.shop_id = ? and ingredient_tech_carts.shop_id = ?", tech.TechCartID, filter.AccessibleShops, tech.ShopID, tech.ShopID).Group("ingredients.id, ingredient_tech_carts.ingredient_id, ingredient_tech_carts.brutto").Scan(&ingredients)
	if ingre.Error != nil {
		return nil, ingre.Error
	}
	var ingredientCost float32 = 0
	var ingredientSum float32 = 0

	for _, ingredient := range ingredients {
		skladIngredient := &model.SkladIngredient{}
		res := r.gormDB.Model(&model.SkladIngredient{}).Joins("inner join sklads on sklads.id = sklad_ingredients.sklad_id").Where("sklad_ingredients.ingredient_id = ? and sklads.shop_id = ?", ingredient.IngredientID, tech.ShopID).Scan(&skladIngredient)
		if res.Error != nil {
			return nil, res.Error
		}
		ingredient.Cost = skladIngredient.Cost * ingredient.Brutto
		// var sum float32 = 0
		// var quantity float32 = 0
		// for _, val := range skladIngredient {
		// 	if val.Quantity <= 0 {
		// 		continue
		// 	}
		// 	sum += (val.Cost * val.Quantity)
		// 	quantity += val.Quantity
		// }
		// if quantity > 0 {
		// 	ingredientCost = sum / quantity
		// } else {
		// 	var cost float32
		// 	res := r.gormDB.Model(&model.ItemPostavka{}).Select("cost").Joins("inner join postavkas on postavkas.id = item_postavkas.postavka_id").Where("item_id = ? and item_postavkas.type = ? and postavkas.deleted = ?", ingredient.IngredientID, utils.TypeIngredient, false).Order("postavkas.time desc").First(&cost)
		// 	if res.Error != nil {
		// 		if res.Error != gorm.ErrRecordNotFound {
		// 			return nil, res.Error
		// 		}
		// 	}
		// 	if res.RowsAffected == 0 {
		// 		deletedCost := r.gormDB.Model(&model.ItemPostavka{}).Select("cost").Joins("inner join postavkas on postavkas.id = item_postavkas.postavka_id").Where("item_id = ? and item_postavkas.type = ?", ingredient.IngredientID, utils.TypeIngredient).Order("postavkas.time desc").First(&cost)
		// 		if deletedCost.Error != nil {
		// 			if deletedCost.Error != gorm.ErrRecordNotFound {
		// 				return nil, deletedCost.Error
		// 			}
		// 		}
		// 		if deletedCost.RowsAffected == 0 {
		// 			res := r.gormDB.Model(&model.SkladIngredient{}).Select("AVG(sklad_ingredients.cost)").Joins("inner join sklads on sklads.id = sklad_ingredients.sklad_id").Where("sklad_ingredients.ingredient_id = ? and sklads.shop_id IN (?)", ingredient.IngredientID, filter.AccessibleShops).Scan(&cost)
		// 			if res.Error != nil {
		// 				return nil, res.Error
		// 			}
		// 		}
		// 	}
		// 	ingredient.Cost = cost
		// }
		ingredientSum += ingredientCost
		ingredient.Netto = ingredient.Brutto
	}
	tech.Ingredients = ingredients
	tech.Cost = ingredientSum

	nabors := []*model.NaborInfo{}
	naborOutput := []*model.NaborOutput{}
	nab := r.gormDB.Model(&model.NaborTechCart{}).Select("nabor_tech_carts.nabor_id as id, nabors.name as name, nabors.min as min, nabors.max as max").Joins("inner join nabors on nabors.nabor_id = nabor_tech_carts.nabor_id and nabors.shop_id = ?", tech.ShopID).Where("nabor_tech_carts.tech_cart_id = ? and nabor_tech_carts.shop_id = ?", tech.TechCartID, tech.ShopID).Scan(&nabors)
	if nab.Error != nil {
		return nil, nab.Error
	}

	for _, nabor := range nabors {
		naborIngredients := []*model.IngredientOutput{}

		nabingre := r.gormDB.Model(&model.IngredientNabor{}).Select("ingredients.ingredient_id as id, ingredients.name, ingredients.category as category_id, category_ingredients.name as category, ingredient_nabors.nabor_id, ingredient_nabors.brutto, ingredient_nabors.price, ingredients.measure, AVG(sklad_ingredients.cost)").Joins("inner join ingredients on ingredients.ingredient_id = ingredient_nabors.ingredient_id inner join category_ingredients on category_ingredients.id = ingredients.category inner join sklad_ingredients on sklad_ingredients.ingredient_id = ingredients.ingredient_id inner join sklads on sklads.id = sklad_ingredients.sklad_id").Where("nabor_id = ? and sklads.shop_id IN (?) and ingredients.shop_id = ? and ingredient_nabors.shop_id = ?", nabor.ID, filter.AccessibleShops, tech.ShopID, tech.ShopID).Group("ingredients.id, category_ingredients.name, ingredient_nabors.nabor_id, ingredient_nabors.brutto, ingredient_nabors.price").Scan(&naborIngredients)
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

	tech.Nabors = naborOutput

	tech.Profit = tech.Price - tech.Cost
	if tech.Cost != 0 {
		tech.Margin = tech.Profit * 100 / tech.Cost
	} else {
		tech.Margin = 0
	}

	return tech, nil
}

func (r *TovarDB) UpdateTechCart(tech *model.ReqTechCart, role string) error {
	/*ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	for _, v := range tech.Ingredients {
		ingredients := model.Ingredient{}
		err := r.gormDB.Model(ingredients).Select("*").Where("ingredient_id = ?", v.IngredientID).Find(&ingredients).Error
		if err != nil {
			return err
		}
	}

	stmt := `update tech_carts set name=$1, category=$2, tax=$3, measure=$4, price=$5, image=$7, discount=$8 where id=$6`

	_, err = tx.ExecContext(ctx, stmt,
		tech.Name,
		tech.Category,
		tech.Tax,
		tech.Measure,
		tech.Price,
		tech.ID,
		tech.Image,
		tech.Discount,
	)
	if err != nil {
		if err := tx.Rollback(); err != nil {
			return err
		}
		return err
	}

	delete := `delete from ingredient_tech_carts where tech_cart_id=$1`
	row := tx.QueryRowContext(ctx, delete,
		tech.ID,
	)
	if row.Err() != nil {
		if err := tx.Rollback(); err != nil {
			return err
		}
		return row.Err()
	}
	for _, ingredient := range tech.Ingredients {
		query := `insert into ingredient_tech_carts (tech_cart_id, ingredient_id, brutto)
		values ($1, $2, $3)`
		err := tx.QueryRowContext(ctx, query,
			tech.ID,
			ingredient.IngredientID,
			ingredient.Brutto,
		)
		if err.Err() != nil {
			if err := tx.Rollback(); err != nil {
				return err
			}
			return err.Err()
		}
	}

	delete = `delete from nabor_tech_carts where tech_cart_id=$1`
	row = tx.QueryRowContext(ctx, delete,
		tech.ID,
	)
	if row.Err() != nil {
		if err := tx.Rollback(); err != nil {
			return err
		}
		return row.Err()
	}
	for _, nabor := range tech.Nabor {
		query := `insert into nabor_tech_carts (tech_cart_id, nabor_id)
		values ($1, $2)`
		err := tx.QueryRowContext(ctx, query,
			tech.ID,
			nabor.NaborID,
		)
		if err.Err() != nil {
			if err := tx.Rollback(); err != nil {
				return err
			}
			return err.Err()
		}
	}
	err = tx.Commit()
	if err != nil {
		return err
	}
	return nil*/

	err := r.gormDB.Transaction(func(tx *gorm.DB) error {
		oldTechCart := model.TechCart{}
		err := tx.Model(&model.TechCart{}).Where("id = ?", tech.ID).Find(&oldTechCart).Error
		if err != nil {
			return err
		}
		for _, shop := range tech.ShopID {
			techCart := &model.TechCart{
				TechCartID: oldTechCart.TechCartID,
				ShopID:     shop,
				Name:       tech.Name,
				Category:   tech.Category,
				Tax:        tech.Tax,
				Image:      tech.Image,
				Measure:    tech.Measure,
				Price:      tech.Price,
				Discount:   tech.Discount,
				Deleted:    false,
				IsVisible:  true,
			}
			res := tx.Model(&model.TechCart{}).Select("*").Omit("id").Where("tech_cart_id = ? and shop_id = ?", oldTechCart.TechCartID, shop).Updates(techCart)
			if res.Error != nil {
				return res.Error
			}
			if res.RowsAffected == 0 {
				techCart.TechCartID = oldTechCart.TechCartID
				techCart.ShopID = shop
				err := tx.Model(&model.TechCart{}).Create(techCart).Error
				if err != nil {
					return err
				}
			} else {
				if role == utils.MasterRole {
					res := tx.Where("tech_cart_id = ? and shop_id = ?", techCart.TechCartID, shop).Delete(&model.IngredientTechCart{})
					if res.Error != nil {
						return res.Error
					}
					res = tx.Where("tech_cart_id = ? and shop_id = ?", techCart.TechCartID, shop).Delete(&model.NaborTechCart{})
					if res.Error != nil {
						return res.Error
					}
				}
			}
			if role == utils.MasterRole {
				for _, ingredient := range tech.Ingredients {
					ingredientTechCart := &model.IngredientTechCart{
						TechCartID:   oldTechCart.TechCartID,
						ShopID:       shop,
						IngredientID: ingredient.IngredientID,
						Brutto:       ingredient.Brutto,
					}
					res := tx.Create(ingredientTechCart)
					if res.Error != nil {
						return res.Error
					}
				}
				for _, nabor := range tech.Nabor {
					naborTechCart := &model.NaborTechCart{
						TechCartID: oldTechCart.TechCartID,
						ShopID:     shop,
						NaborID:    nabor.NaborID,
					}
					res := tx.Create(naborTechCart)
					if res.Error != nil {
						return res.Error
					}
				}
			}
		}
		return nil
	})
	if err != nil {
		return err
	}

	return nil
}

func (r *TovarDB) DeleteTechCart(id int) error {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()
	stmt := `update tech_carts set deleted = $2 where id=$1`
	_, err := r.db.ExecContext(ctx, stmt, id, true)
	if err != nil {
		return err
	}

	return nil
}

func (r *TovarDB) GetTovarWithParams(sortParam, sklad, search string, category int) ([]*model.TovarOutput, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()
	tovars := []*model.TovarOutput{}
	stmt := ``
	var row *sql.Rows
	var err error
	if category != 0 {
		stmt = `select tovars.id, tovars.name, category_tovars.name, tovars.image, tovars.tax, tovars.measure, tovars.cost, tovars.price, tovars.profit, tovars.margin, tovars.discount from tovars inner join category_tovars on tovars.category = category_tovars.id where tovars.deleted = $1 and tovars.category = $2 `
		row, err = r.db.QueryContext(ctx, stmt, false, category)
	} else {
		stmt = `select tovars.id, tovars.name, category_tovars.name, tovars.image, tovars.tax, tovars.measure, tovars.cost, tovars.price, tovars.profit, tovars.margin, tovars.discount from tovars inner join category_tovars on tovars.category = category_tovars.id where tovars.deleted = $1 and tovars.name like $2 or category_tovars.name like $2 or tovars.tax like $2 or tovars.measure like $2 order by ` + `tovars.` + sortParam + `;`
		row, err = r.db.QueryContext(ctx, stmt, false, search)
	}
	if err != nil {
		return nil, err
	}

	for row.Next() {
		tovar := &model.TovarOutput{}
		err := row.Scan(
			&tovar.ID,
			&tovar.Name,
			&tovar.Category,
			&tovar.Image,
			&tovar.Tax,
			&tovar.Measure,
			&tovar.Cost,
			&tovar.Price,
			&tovar.Profit,
			&tovar.Margin,
			&tovar.Discount,
		)
		if err != nil {
			return nil, err
		}
		tovars = append(tovars, tovar)
	}
	if len(tovars) == 0 {
		tovars = []*model.TovarOutput{}
	}
	return tovars, nil
}

func (r *TovarDB) GetTechCartWithParams(sortParam, sklad string, category int) ([]*model.TechCartOutput, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()
	techCarts := []*model.TechCartOutput{}
	stmt := ``
	var row *sql.Rows
	var err error
	if category != 0 {
		stmt = `select tech_carts.id, tech_carts.name, category_tovars.name, tech_carts.image, tech_carts.tax, tech_carts.measure, tech_carts.cost, tech_carts.price, tech_carts.profit, tech_carts.margin, tech_carts.discount from tech_carts inner join category_tovars on tech_carts.category = category_tovars.id where tech_carts.deleted = $1 and tech_carts.category = $2`
		row, err = r.db.QueryContext(ctx, stmt, false, category)
	} else {
		stmt = `select tech_carts.id, tech_carts.name, category_tovars.name, tech_carts.image, tech_carts.tax, tech_carts.measure, tech_carts.cost, tech_carts.price, tech_carts.profit, tech_carts.margin, tech_carts.discount from tech_carts inner join category_tovars on tech_carts.category = category_tovars.id where tech_carts.deleted = $1`
		row, err = r.db.QueryContext(ctx, stmt, false)
	}

	if err != nil {
		return nil, err
	}

	for row.Next() {
		techCart := &model.TechCartOutput{}
		err := row.Scan(
			&techCart.ID,
			&techCart.Name,
			&techCart.Category,
			&techCart.Image,
			&techCart.Tax,
			&techCart.Measure,
			&techCart.Cost,
			&techCart.Price,
			&techCart.Profit,
			&techCart.Margin,
			&techCart.Discount,
		)
		if err != nil {
			return nil, err
		}
		techCarts = append(techCarts, techCart)
	}

	if len(techCarts) == 0 {
		techCarts = []*model.TechCartOutput{}
	}

	return techCarts, nil
}

func (r *TovarDB) GetTechCartNabor(id int) ([]*model.NaborOutput, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()
	nabors := []*model.NaborOutput{}
	stmt := `select nabors.id, nabors.name, nabors.min, nabors.max from nabor_tech_carts inner join nabors on nabor_tech_carts.nabor_id = nabors.id where nabor_tech_carts.tech_cart_id = $1`
	row, err := r.db.QueryContext(ctx, stmt, id)
	if err != nil {
		return nil, err
	}
	for row.Next() {
		nabor := &model.NaborOutput{}
		err := row.Scan(
			&nabor.ID,
			&nabor.Name,
			&nabor.Min,
			&nabor.Max,
		)
		if err != nil {
			return nil, err
		}
		query := `select ingredients.ingredient_id, ingredients.name, ingredients.image, ingredients.measure, ingredient_nabors.brutto, ingredient_nabors.price from ingredient_nabors inner join ingredients on ingredient_nabors.ingredient_id = ingredients.ingredient_id where ingredient_nabors.nabor_id = $1`
		naborRow, err := r.db.QueryContext(ctx, query, nabor.ID)
		if err != nil {
			return nil, err
		}
		ingredients := []*model.IngredientOutput{}
		for naborRow.Next() {
			ingredient := &model.IngredientOutput{}
			err := naborRow.Scan(
				&ingredient.ID,
				&ingredient.Name,
				&ingredient.Image,
				&ingredient.Measure,
				&ingredient.Brutto,
				&ingredient.Price,
			)
			if err != nil {
				return nil, err
			}
			ingredients = append(ingredients, ingredient)
		}
		nabor.NaborIngredient = ingredients
		nabors = append(nabors, nabor)
	}

	if len(nabors) == 0 {
		nabors = []*model.NaborOutput{}
	}

	return nabors, nil
}

func (r *TovarDB) GetEverything() ([]*model.ItemOutput, error) {
	items := []*model.ItemOutput{}
	ingredients := []*model.Ingredient{}
	err := r.gormDB.Model(&model.Ingredient{}).Scan(&ingredients).Error
	if err != nil {
		return nil, err
	}
	tovars := []*model.Tovar{}
	err = r.gormDB.Model(&model.Tovar{}).Scan(&tovars).Error
	if err != nil {
		return nil, err
	}
	techCarts := []*model.TechCart{}
	err = r.gormDB.Model(&model.TechCart{}).Scan(&techCarts).Error
	if err != nil {
		return nil, err
	}
	for _, ingredient := range ingredients {
		item := &model.ItemOutput{
			ID:      ingredient.IngredientID,
			Name:    ingredient.Name,
			Measure: ingredient.Measure,
			Cost:    ingredient.Cost,
			Type:    utils.TypeIngredient,
		}
		items = append(items, item)
	}
	for _, tovar := range tovars {
		item := &model.ItemOutput{
			ID:      tovar.TovarID,
			Name:    tovar.Name,
			Measure: tovar.Measure,
			Cost:    0, //????
			Type:    utils.TypeIngredient,
		}
		items = append(items, item)
	}
	for _, techCart := range techCarts {
		item := &model.ItemOutput{
			//ID:      techCart.TechCartID,
			ID:      techCart.ID,
			Name:    techCart.Name,
			Measure: techCart.Measure,
			Type:    utils.TypeIngredient,
		}
		items = append(items, item)
	}
	return items, nil
}

func (r *TovarDB) GetPureTovarByShopID(id, shopID int) (*model.Tovar, error) {
	tovar := &model.Tovar{}
	err := r.gormDB.Model(&model.Tovar{}).Where("tovar_id = ? and shop_id = ?", id, shopID).Scan(tovar).Error
	if err != nil {
		return nil, err
	}
	return tovar, nil
}

func (r *TovarDB) GetTovarsTovarIDByTovarsID(id int) (int, error) {
	tovar := &model.Tovar{}
	err := r.gormDB.Model(&model.Tovar{}).Where("id = ?", id).Scan(tovar).Error
	if err != nil {
		return 0, err
	}
	if tovar.TovarID == 0 {
		return id, nil
	}
	return tovar.TovarID, nil
}

func (r *TovarDB) GetTechCartsTechCartIDByTechCartsID(id int) (int, error) {
	techCart := &model.TechCart{}
	err := r.gormDB.Model(&model.TechCart{}).Where("id = ?", id).Scan(techCart).Error
	if err != nil {
		return 0, err
	}
	if techCart.TechCartID == 0 {
		return id, nil
	}
	return techCart.TechCartID, nil
}

func (r *TovarDB) GetDeletedTovar(filter *model.Filter) ([]*model.TovarOutput, int64, error) {
	tovars := []*model.TovarOutput{}
	res := r.gormDB

	if filter.Measure == "" {
		res = r.gormDB.Model(&model.Tovar{}).Select("tovars.id, tovars.tovar_id, tovars.shop_id, shops.name as shop_name, tovars.name, tovars.image, category_tovars.name as category, tovars.tax, tovars.measure, AVG(sklad_tovars.cost) as cost, tovars.price, tovars.discount").Joins("inner join category_tovars on tovars.category = category_tovars.id inner join sklad_tovars on tovars.tovar_id = sklad_tovars.tovar_id inner join sklads on  sklads.id = sklad_tovars.sklad_id inner join shops on shops.id = tovars.shop_id").Where("tovars.deleted = ? and sklads.shop_id IN (?)", true, filter.AccessibleShops).Group("tovars.id, tovars.tovar_id, shops.name, category_tovars.name")
		if res.Error != nil {
			return nil, 0, res.Error
		}
	} else {
		res = r.gormDB.Model(&model.Tovar{}).Select("tovars.id, tovars.tovar_id, tovars.shop_id, shops.name as shop_name, tovars.name, tovars.image, category_tovars.name as category, tovars.tax, tovars.measure, AVG(sklad_tovars.cost) as cost, tovars.price, tovars.discount").Joins("inner join category_tovars on tovars.category = category_tovars.id inner join sklad_tovars on tovars.tovar_id = sklad_tovars.tovar_id inner join sklads on  sklads.id = sklad_tovars.sklad_id inner join shops on shops.id = tovars.shop_id").Where("tovars.deleted = ? and sklads.shop_id IN (?) and tovars.measure = ?", true, filter.AccessibleShops, filter.Measure).Group("tovars.id, tovars.tovar_id, shops.name, category_tovars.name").Scan(&tovars)
		if res.Error != nil {
			return nil, 0, res.Error
		}
	}

	newRes, count, err := filter.FilterResults(res, &model.Tovar{}, utils.DefaultPageSize, "", "", "")

	if err != nil {
		return nil, 0, err
	}

	if err := newRes.Scan(&tovars).Error; err != nil {
		return nil, 0, err
	}

	for _, tovar := range tovars {
		skladTovars := []*model.SkladTovar{}
		res := r.gormDB.Model(&model.SkladTovar{}).Joins("inner join sklads on sklads.id = sklad_tovars.sklad_id").Where("sklad_tovars.tovar_id = ? and sklads.shop_id IN (?)", tovar.TovarID, filter.AccessibleShops).Scan(&skladTovars)
		if res.Error != nil {
			return nil, 0, res.Error
		}
		var sum float32 = 0
		var quantity float32 = 0
		for _, val := range skladTovars {
			if val.Quantity <= 0 {
				continue
			}
			sum += (val.Cost * val.Quantity)
			quantity += val.Quantity
		}
		if quantity > 0 {
			tovar.Cost = sum / quantity
		} else {
			var cost float32
			res := r.gormDB.Model(&model.ItemPostavka{}).Select("cost").Joins("inner join postavkas on postavkas.id = item_postavkas.postavka_id").Where("item_id = ? and item_postavkas.type = ? and postavkas.deleted = ?", tovar.TovarID, utils.TypeTovar, false).Order("postavkas.time desc").First(&cost)
			if res.Error != nil {
				if res.Error != gorm.ErrRecordNotFound { //???
					return nil, 0, res.Error
				}
			}
			if res.RowsAffected == 0 {
				deletedCost := r.gormDB.Model(&model.ItemPostavka{}).Select("cost").Joins("inner join postavkas on postavkas.id = item_postavkas.postavka_id").Where("item_id = ? and item_postavkas.type = ?", tovar.TovarID, utils.TypeTovar).Order("postavkas.time desc").First(&cost)
				if deletedCost.Error != nil {
					if deletedCost.Error != gorm.ErrRecordNotFound {
						return nil, 0, deletedCost.Error
					}
				}
				if deletedCost.RowsAffected == 0 {
					res := r.gormDB.Model(&model.SkladTovar{}).Select("AVG(sklad_tovars.cost)").Joins("inner join sklads on sklads.id = sklad_tovars.sklad_id").Where("sklad_tovars.tovar_id = ? and sklads.shop_id IN (?)", tovar.TovarID, filter.AccessibleShops).Scan(&cost)
					if res.Error != nil {
						return nil, 0, res.Error
					}
				}
			}
			tovar.Cost = cost
		}
		tovar.Profit = tovar.Price - tovar.Cost
		if tovar.Cost != 0 {
			tovar.Margin = tovar.Profit * 100 / tovar.Cost
		}
	}

	return tovars, count, nil
}

func (r *TovarDB) GetDeletedTechCart(filter *model.Filter) ([]*model.TechCartResponse, int64, error) {
	techCarts := []*model.TechCartResponse{}
	res := r.gormDB.Model(&model.TechCart{}).Select("tech_carts.id, tech_carts.tech_cart_id, tech_carts.shop_id, shops.name as shop_name, tech_carts.name, tech_carts.image, category_tovars.name as category, tech_carts.tax, tech_carts.measure, AVG(sklad_tovars.cost) as cost, tech_carts.price, tech_carts.discount, tech_carts.deleted").Joins("inner join category_tovars on tech_carts.category = category_tovars.id inner join sklad_tovars on tech_carts.tech_cart_id = sklad_tovars.tech_cart_id inner join sklads on  sklads.id = sklad_tovars.sklad_id inner join shops on shops.id = tech_carts.shop_id").Where("tech_carts.deleted = ? and sklads.shop_id IN (?)", true, filter.AccessibleShops).Group("tech_carts.id, tech_carts.tech_cart_id, shops.name, category_tovars.name")
	if res.Error != nil {
		return nil, 0, res.Error
	}
	newRes, count, err := filter.FilterResults(res, &model.TechCart{}, utils.DefaultPageSize, "", "", "")

	if err != nil {
		return nil, 0, err
	}

	if err := newRes.Scan(&techCarts).Error; err != nil {
		return nil, 0, err
	}
	for _, tech := range techCarts {
		techIngredients := []*model.IngredientTechCart{}
		res := r.gormDB.Model(&model.IngredientTechCart{}).Select("*").Where("ingredient_tech_carts.tech_cart_id = ? and shop_id = ?", tech.TechCartID, tech.ShopID).Scan(&techIngredients)
		if res.Error != nil {
			return nil, 0, res.Error
		}

		var ingredientCost float32 = 0
		var ingredientSum float32 = 0
		for _, item := range techIngredients {
			skladIngredient := []*model.SkladIngredient{}
			res := r.gormDB.Model(&model.SkladIngredient{}).Joins("inner join sklads on sklads.id = sklad_ingredients.sklad_id").Where("sklad_ingredients.ingredient_id = ? and sklads.shop_id IN (?)", item.IngredientID, filter.AccessibleShops).Scan(&skladIngredient)
			if res.Error != nil {
				return nil, 0, res.Error
			}
			var sum float32 = 0
			var quantity float32 = 0
			for _, val := range skladIngredient {
				if val.Quantity <= 0 {
					continue
				}
				sum += (val.Cost * val.Quantity)
				quantity += val.Quantity
			}
			if quantity > 0 {
				ingredientCost = sum / quantity
			} else {
				var cost float32
				res := r.gormDB.Model(&model.ItemPostavka{}).Select("cost").Joins("inner join postavkas on postavkas.id = item_postavkas.postavka_id").Where("item_id = ? and item_postavkas.type = ? and postavkas.deleted = ?", item.IngredientID, utils.TypeIngredient, false).Order("postavkas.time desc").First(&cost)
				if res.Error != nil {
					if res.Error != gorm.ErrRecordNotFound {
						return nil, 0, res.Error
					}
				}
				if res.RowsAffected == 0 {
					deletedCost := r.gormDB.Model(&model.ItemPostavka{}).Select("cost").Joins("inner join postavkas on postavkas.id = item_postavkas.postavka_id").Where("item_id = ? and item_postavkas.type = ?", item.IngredientID, utils.TypeIngredient).Order("postavkas.time desc").First(&cost)
					if deletedCost.Error != nil {
						if deletedCost.Error != gorm.ErrRecordNotFound {
							return nil, 0, deletedCost.Error
						}
					}
					if deletedCost.RowsAffected == 0 {
						res := r.gormDB.Model(&model.SkladIngredient{}).Select("COALESCE(AVG(sklad_ingredients.cost), 0.0)").Joins("inner join sklads on sklads.id = sklad_ingredients.sklad_id").Where("sklad_ingredients.ingredient_id = ? and sklads.shop_id IN (?)", item.IngredientID, filter.AccessibleShops).Scan(&cost)
						if res.Error != nil {
							return nil, 0, res.Error
						}
					}
				}
				ingredientCost = cost
			}
			ingredientSum += ingredientCost * item.Brutto
		}
		tech.Cost = ingredientSum
		tech.Profit = tech.Price - tech.Cost
		if tech.Cost != 0 {
			tech.Margin = tech.Profit * 100 / tech.Cost
		}
	}
	return techCarts, count, nil
}

func (r *TovarDB) RecreateTovar(tovar *model.Tovar) error {
	res := r.gormDB.Model(&model.Tovar{}).Where("tovar_id = ? and shop_id = ?", tovar.TovarID, tovar.ShopID).Update("deleted", false)
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return errors.New("tovar not found")
	}
	return nil
}

func (r *TovarDB) RecreateTechCart(techCart *model.TechCart) error {
	res := r.gormDB.Model(&model.TechCart{}).Where("tech_cart_id = ? and shop_id = ?", techCart.TechCartID, techCart.ShopID).Update("deleted", false)
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return errors.New("techCart not found")
	}
	return nil
}

func (r *TovarDB) GetToAddTovar(filter *model.Filter) ([]*model.TovarMaster, int64, error) {
	tovarMaster := []*model.TovarMaster{}
	//res := r.gormDB.Model(&model.TovarMaster{}).Select("tovar_masters.id, tovar_masters.name, category_tovars.name as category, tovar_masters.image, tovar_masters.tax, tovar_masters.measure, tovar_masters.price, tovar_masters.deleted, tovar_masters.discount, tovar_masters.status").Joins("inner join category_tovars on category_tovars.id = tovar_masters.category left join tovars on tovar_masters.id = tovars.tovar_id and tovars.shop_id IN (?)", filter.AccessibleShops).Where("tovars.id is null and tovar_masters.deleted = ?", false)
	res := r.gormDB.Model(&model.TovarMaster{}).Select("tovar_masters.*").Joins("left join tovars on tovar_masters.id = tovars.tovar_id and tovars.shop_id IN (?)", filter.AccessibleShops).Where("tovars.id is null and tovar_masters.deleted = ?", false)
	if res.Error != nil {
		return nil, 0, res.Error
	}
	newRes, count, err := filter.FilterResults(res, &model.TovarMaster{}, utils.DefaultPageSize, "", "", "")

	if err != nil {
		return nil, 0, err
	}

	if err := newRes.Scan(&tovarMaster).Error; err != nil {
		return nil, 0, err
	}

	return tovarMaster, count, nil
}

func (r *TovarDB) GetToAddTechCart(filter *model.Filter) ([]*model.TechCartMaster, int64, error) {
	techCartMaster := []*model.TechCartMaster{}
	//res := r.gormDB.Model(&model.TechCartMaster{}).Select("tech_cart_masters.id, tech_cart_masters.name, category_tovars.name as category, tech_cart_masters.image, tech_cart_masters.tax, tech_cart_masters.measure, tech_cart_masters.price, tech_cart_masters.deleted, tech_cart_masters.discount, tech_cart_masters.status").Joins("inner join category_tovars on category_tovars.id = tech_cart_masters.category left join tech_carts on tech_cart_masters.id = tech_carts.tech_cart_id and tech_carts.shop_id IN (?)", filter.AccessibleShops).Where("tech_carts.id is null and tech_cart_masters.deleted = ?", false)
	res := r.gormDB.Model(&model.TechCartMaster{}).Select("tech_cart_masters.*").Joins("left join tech_carts on tech_cart_masters.id = tech_carts.tech_cart_id and tech_carts.shop_id IN (?)", filter.AccessibleShops).Where("tech_carts.id is null and tech_cart_masters.deleted = ?", false)
	if res.Error != nil {
		return nil, 0, res.Error
	}
	newRes, count, err := filter.FilterResults(res, &model.TechCartMaster{}, utils.DefaultPageSize, "", "", "")

	if err != nil {
		return nil, 0, err
	}

	if err := newRes.Scan(&techCartMaster).Error; err != nil {
		return nil, 0, err
	}

	return techCartMaster, count, nil
}

func (r *TovarDB) GetIdsOfShopsWhereTheTovarAlreadyExist(tovarID int) ([]int, error) {
	shopIds := []int{}
	err := r.gormDB.Model(&model.Tovar{}).Select("shop_id").Where("tovar_id = ?", tovarID).Scan(&shopIds).Error
	if err != nil {
		return nil, err
	}
	return shopIds, nil
}

func (r *TovarDB) GetIdsOfShopsWhereTheTechCartAlreadyExist(techCartID int) ([]int, error) {
	shopIds := []int{}
	err := r.gormDB.Model(&model.TechCart{}).Select("shop_id").Where("tech_cart_id = ?", techCartID).Scan(&shopIds).Error
	if err != nil {
		return nil, err
	}
	return shopIds, nil
}
