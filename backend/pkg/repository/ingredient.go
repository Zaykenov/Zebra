package repository

import (
	"database/sql"
	"errors"
	"zebra/model"
	"zebra/utils"

	"gorm.io/gorm"
)

type IngredientDB struct {
	db     *sql.DB
	gormDB *gorm.DB
}

func NewIngredientDB(db *sql.DB, gormDB *gorm.DB) *IngredientDB {
	return &IngredientDB{db: db, gormDB: gormDB}
}

func (r *IngredientDB) AddIngredients(ingredients []*model.Ingredient) error {
	err := r.gormDB.Transaction(func(tx *gorm.DB) error {
		err := tx.Create(ingredients).Error
		if err != nil {
			return err
		}
		for _, ingredient := range ingredients {
			sklad := &model.Sklad{}
			err = r.gormDB.Model(&model.Sklad{}).Where("shop_id = ? and deleted = ?", ingredient.ShopID, false).First(sklad).Error
			if err != nil {
				return err
			}
			skladIngredient := &model.SkladIngredient{
				IngredientID: ingredient.IngredientID,
				SkladID:      sklad.ID,
				Quantity:     0,
				Cost:         ingredient.Cost,
			}
			err = r.gormDB.Create(skladIngredient).Error
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

func (r *IngredientDB) AddIngredient(ingredient *model.Ingredient) error {
	err := r.gormDB.Transaction(func(tx *gorm.DB) error {
		err := r.gormDB.Create(ingredient).Error
		if err != nil {
			return err
		}

		sklads := []*model.Sklad{}
		err = r.gormDB.Model(&model.Sklad{}).Where("deleted = ?", false).Find(&sklads).Error
		if err != nil {
			return err
		}
		for _, sklad := range sklads {
			skladIngredient := &model.SkladIngredient{
				SkladID:      sklad.ID,
				IngredientID: ingredient.IngredientID,
				Quantity:     0,
				Cost:         0,
			}
			err = r.gormDB.Create(skladIngredient).Error
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

func (r *IngredientDB) GetIngredientsForNabor(id int) ([]*model.Ingredient, error) {
	ingredients := []*model.Ingredient{}
	err := r.gormDB.Table("ingredients").Select("ingredients.id, ingredients.ingredient_id, ingredients.shop_id, ingredients.name, ingredients.measure, ingredients.image, ingredients.category, sklad_ingredients.cost").Joins("inner join ingredient_nabors on ingredients.ingredient_id = ingredient_nabors.ingredient_id inner join sklad_ingredients on sklad_ingredients.ingredient_id = ingredients.ingredient_id").Where("ingredient_nabors.nabor_id = ?", id).Find(&ingredients).Error
	if err != nil {
		return nil, err
	}

	return ingredients, nil
}

func (r *IngredientDB) GetAllIngredient(filter *model.Filter) ([]*model.IngredientOutput, int64, error) {
	/*ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()
	ingredients := []*model.IngredientOutput{}
	stmt := `select ingredient.id, ingredient.name, category_ingredient.name, ingredient.image, ingredient.measure, ingredient.cost from ingredient inner join category_ingredient on ingredient.category = category_ingredient.id where ingredient.deleted=$1`

	row, err := r.db.QueryContext(ctx, stmt, false)
	if err != nil {
		return nil, err
	}

	for row.Next() {
		ingredient := &model.IngredientOutput{}
		err := row.Scan(
			&ingredient.ID,
			&ingredient.Name,
			&ingredient.Category,
			&ingredient.Image,
			&ingredient.Measure,
			&ingredient.Cost,
		)
		if err != nil {
			return nil, err
		}
		ingredients = append(ingredients, ingredient)
	}
	*/
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
	ingredients := []*model.IngredientOutput{}
	res := r.gormDB
	if filter.Measure == "" {
		res = r.gormDB.Table("ingredients").Select("ingredients.id, ingredients.ingredient_id, ingredients.shop_id, shops.name as shop_name, ingredients.name, category_ingredients.name as category, ingredients.image, ingredients.measure, AVG(sklad_ingredients.cost) as cost").Joins("inner join category_ingredients on ingredients.category = category_ingredients.id inner join sklad_ingredients on ingredients.ingredient_id = sklad_ingredients.ingredient_id inner join sklads on  sklads.id = sklad_ingredients.sklad_id inner join shops on shops.id = ingredients.shop_id").Where("ingredients.deleted = ? and sklads.shop_id IN (?)", false, filter.AccessibleShops).Group("ingredients.id, ingredients.ingredient_id, category_ingredients.name, shops.name")
		if res.Error != nil {
			return nil, 0, res.Error
		}
	} else {
		res = r.gormDB.Table("ingredients").Select("ingredients.id, ingredients.ingredient_id, ingredients.shop_id, shops.name as shop_name, ingredients.name, category_ingredients.name as category, ingredients.image, ingredients.measure, AVG(sklad_ingredients.cost) as cost").Joins("inner join category_ingredients on ingredients.category = category_ingredients.id inner join sklad_ingredients on ingredients.ingredient_id = sklad_ingredients.ingredient_id inner join sklads on  sklads.id = sklad_ingredients.sklad_id inner join shops on shops.id = ingredients.shop_id").Where("ingredients.deleted = ? and sklads.shop_id IN (?) and ingredients.measure = ?", false, filter.AccessibleShops, filter.Measure).Group("ingredients.id, ingredients.ingredient_id, category_ingredients.name, shops.name")

		if res.Error != nil {
			return nil, 0, res.Error
		}
	}

	newRes, count, err := filter.FilterResults(res, &model.Ingredient{}, utils.DefaultPageSize, "", "", "")

	if err != nil {
		return nil, 0, err
	}

	if newRes.Scan(&ingredients).Error != nil {
		return nil, 0, newRes.Error
	}

	for _, ingredient := range ingredients {
		skladIngredient := &model.SkladIngredient{}
		res := r.gormDB.Model(&model.SkladIngredient{}).Joins("inner join sklads on sklads.id = sklad_ingredients.sklad_id").Where("sklad_ingredients.ingredient_id = ? and sklads.shop_id = ?", ingredient.IngredientID, ingredient.ShopID).Scan(&skladIngredient)
		if res.Error != nil {
			return nil, 0, res.Error
		}
		ingredient.Cost = skladIngredient.Cost
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
		// 	ingredient.Cost = sum / quantity
		// } else {
		// 	var cost float32
		// 	res := r.gormDB.Model(&model.ItemPostavka{}).Select("cost").Joins("inner join postavkas on postavkas.id = item_postavkas.postavka_id").Where("item_id = ? and item_postavkas.type = ? and postavkas.sklad_id = ?", ingredient.IngredientID, utils.TypeIngredient).Order("postavkas.time desc").First(&cost)
		// 	if res.Error != nil {
		// 		if res.Error != gorm.ErrRecordNotFound {
		// 			return nil, 0, res.Error

		// 		}
		// 	}
		// 	if res.RowsAffected == 0 {
		// 		deletedCost := r.gormDB.Model(&model.ItemPostavka{}).Select("cost").Joins("inner join postavkas on postavkas.id = item_postavkas.postavka_id").Where("item_id = ? and item_postavkas.type = ?", ingredient.IngredientID, utils.TypeIngredient).Order("postavkas.time desc").First(&cost)
		// 		if deletedCost.Error != nil {
		// 			if deletedCost.Error != gorm.ErrRecordNotFound {
		// 				return nil, 0, deletedCost.Error
		// 			}
		// 		}
		// 		if deletedCost.RowsAffected == 0 {
		// 			res := r.gormDB.Model(&model.SkladIngredient{}).Select("COALESCE(AVG(sklad_ingredients.cost), 0.0)").Joins("inner join sklads on sklads.id = sklad_ingredients.sklad_id").Where("sklad_ingredients.ingredient_id = ? and sklads.shop_id IN (?)", ingredient.IngredientID, filter.AccessibleShops).Scan(&cost)
		// 			if res.Error != nil {
		// 				return nil, 0, res.Error
		// 			}
		// 		}
		// 	}
		// 	ingredient.Cost = cost
		// }
	}

	return ingredients, count, nil

}

func (r *IngredientDB) GetIngredient(id int, filter *model.Filter) (*model.IngredientOutput, error) {
	/*ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()
	ingredient := &model.IngredientOutput{}
	stmt := `select ingredient.id, ingredient.name, category_ingredient.name, ingredient.image, ingredient.measure, ingredient.cost from ingredient inner join category_ingredient on ingredient.category = category_ingredient.id where ingredient.id=$1`

	row := r.db.QueryRowContext(ctx, stmt, id)

	if row.Err() != nil {
		return nil, row.Err()
	}

	err := row.Scan(
		&ingredient.ID,
		&ingredient.Name,
		&ingredient.Category,
		&ingredient.Image,
		&ingredient.Measure,
		&ingredient.Cost,
	)
	if err != nil {
		return nil, err
	}


	return ingredient, nil*/
	ingredient := &model.IngredientOutput{}
	err := r.gormDB.Table("ingredients").Select("ingredients.id, ingredients.ingredient_id, ingredients.shop_id, ingredients.name, category_ingredients.name as category, ingredients.category as category_id, ingredients.image, ingredients.measure, AVG(sklad_ingredients.cost) as cost").Joins("inner join category_ingredients on ingredients.category = category_ingredients.id inner join sklad_ingredients on ingredients.ingredient_id = sklad_ingredients.ingredient_id inner join sklads on  sklads.id = sklad_ingredients.sklad_id").Where("ingredients.id = ? and sklads.shop_id IN (?)", id, filter.AccessibleShops).Group("ingredients.id, category_ingredients.name").Scan(&ingredient).Error
	if err != nil {
		return nil, err
	}
	skladIngredient := &model.SkladIngredient{}
	res := r.gormDB.Model(&model.SkladIngredient{}).Joins("inner join sklads on sklads.id = sklad_ingredients.sklad_id").Where("sklad_ingredients.ingredient_id = ? and sklads.shop_id = ?", ingredient.ID, ingredient.ShopID).Scan(&skladIngredient)
	if res.Error != nil {
		return nil, res.Error
	}
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
	// 	ingredient.Cost = sum / quantity
	// } else {
	// 	var cost float32
	// 	res := r.gormDB.Model(&model.ItemPostavka{}).Select("cost").Joins("inner join postavkas on postavkas.id = item_postavkas.postavka_id").Where("item_id = ? and item_postavkas.type = ?", ingredient.IngredientID, utils.TypeIngredient).Order("postavkas.time desc").First(&cost)
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
	return ingredient, nil
}

func (r *IngredientDB) UpdateIngredient(ingredient *model.ReqIngredient, shops []int) error {
	/*ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()
	stmt := `update ingredient set name = $2, category = $3, measure = $4, cost = $5, image = $6 where id=$1`

	row := r.db.QueryRowContext(ctx, stmt,
		ingredient.ID,
		ingredient.Name,
		ingredient.Category,
		ingredient.Measure,
		ingredient.Cost,
		ingredient.Image,
	)

	if row.Err() != nil {
		return row.Err()
	}

	return nil*/
	oldItem := &model.Ingredient{}
	err := r.gormDB.Model(&model.Ingredient{}).Where("id = ?", ingredient.ID).First(&oldItem).Error
	if err != nil {
		return err
	}
	for _, shop := range shops {
		res := r.gormDB.Model(model.Ingredient{}).Where("ingredient_id = ? and shop_id = ?", oldItem.IngredientID, shop).Updates(model.Ingredient{Name: ingredient.Name, Category: ingredient.Category, Image: ingredient.Image, Measure: ingredient.Measure, Deleted: false, IsVisible: ingredient.IsVisible})
		if res.Error != nil {
			return res.Error
		}
		if res.RowsAffected == 0 {
			err := r.gormDB.Create(&model.Ingredient{IngredientID: oldItem.IngredientID, ShopID: shop, Name: ingredient.Name, Category: ingredient.Category, Image: ingredient.Image, Measure: ingredient.Measure, IsVisible: ingredient.IsVisible}).Error
			if err != nil {
				return err
			}
			sklad := &model.Sklad{}
			err = r.gormDB.Model(&model.Sklad{}).Where("shop_id = ? and deleted = ?", shop, false).First(sklad).Error
			if err != nil {
				return err
			}
			skladIngredient := &model.SkladIngredient{
				IngredientID: oldItem.IngredientID,
				SkladID:      sklad.ID,
				Quantity:     0,
				Cost:         ingredient.Cost,
			}
			err = r.gormDB.Create(skladIngredient).Error
			if err != nil {
				return err
			}
		}
	}
	return nil
}
func (r *IngredientDB) DeleteIngredient(id int) error {
	/*ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()
	stmt := `update ingredient set deleted = $2 where id=$1`
	row := r.db.QueryRowContext(ctx, stmt, id, true)

	if row.Err() != nil {
		return row.Err()
	}

	return nil*/
	if err := r.gormDB.Table("ingredients").Where("id = ?", id).Update("deleted", true).Error; err != nil {
		return err
	}
	err := r.gormDB.Model(&model.IngredientTechCart{}).Where("ingredient_tech_carts.ingredient_id = ?", id).Delete(&model.IngredientTechCart{}).Error
	if err != nil { //???
		return err
	}
	return nil
}

func (r *IngredientDB) RestoreIngredient(id int) error {
	/*ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()
	stmt := `update ingredient set deleted = $2 where id=$1`
	row := r.db.QueryRowContext(ctx, stmt, id, false)

	if row.Err() != nil {
		return row.Err()
	}

	return nil*/
	if err := r.gormDB.Table("ingredients").Where("id = ?", id).Update("deleted", false).Error; err != nil {
		return err
	}
	return nil
}

func (r *IngredientDB) AddCategoryIngredient(category *model.CategoryIngredient) error {
	/*ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	stmt := `insert into category_ingredient (name, image)
	values ($1, $2)`
	err := r.db.QueryRowContext(ctx, stmt,
		category.Name,
		category.Image,
	)

	if err != nil {
		return err.Err()
	}

	return nil*/
	err := r.gormDB.Create(&category).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *IngredientDB) GetAllCategoryIngredient(filter *model.Filter) ([]*model.CategoryIngredient, int64, error) {
	/*ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()
	categories := []*model.Category{}
	stmt := `select * from category_ingredient`

	row, err := r.db.QueryContext(ctx, stmt)
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
	categories := []*model.CategoryIngredient{}
	res := r.gormDB.Table("category_ingredients").Find(&categories)
	if res.Error != nil {
		return nil, 0, res.Error
	}
	if len(categories) == 0 {
		categories = []*model.CategoryIngredient{}
	}

	newRes, count, err := filter.FilterResults(res, model.CategoryIngredient{}, utils.DefaultPageSize, "", "", "")

	if err != nil {
		return nil, 0, err
	}

	if newRes.Scan(&categories).Error != nil {
		return nil, 0, newRes.Error
	}

	return categories, count, nil

}

func (r *IngredientDB) GetCategoryIngredient(id int) (*model.CategoryIngredient, error) {
	/*ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()
	category := &model.Category{}
	stmt := `select * from category_ingredient where id=$1`

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
	category := &model.CategoryIngredient{}
	err := r.gormDB.Table("category_ingredients").Where("id = ?", id).Scan(&category).Error
	if err != nil {
		return nil, err
	}
	return category, nil
}

func (r *IngredientDB) UpdateCategoryIngredient(category *model.CategoryIngredient) error {
	/*ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()
	stmt := `update category_ingredient set name = $2 where id=$1`

	row := r.db.QueryRowContext(ctx, stmt,
		category.ID,
		category.Name,
	)

	if row.Err() != nil {
		return row.Err()
	}

	return nil*/
	err := r.gormDB.Model(&category).Updates(category).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *IngredientDB) DeleteCategoryIngredient(id int) error {
	/*ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()
	stmt := `delete from category_ingredient where id=$1`
	row := r.db.QueryRowContext(ctx, stmt, id)

	if row.Err() != nil {
		return row.Err()
	}

	return nil*/
	if err := r.gormDB.Table("category_ingredients").Where("id = ?", id).Update("deleted", true).Error; err != nil {
		return err
	}
	return nil
}

func (r *IngredientDB) AddNabor(nabor *model.Nabor) (*model.Nabor, error) {
	/*ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	stmt := `insert into nabor (name, min, max)
	values ($1, $2, $3) returning id`
	row := tx.QueryRowContext(ctx, stmt,
		nabor.Name,
		nabor.Min,
		nabor.Max,
	)

	if row.Err() != nil {
		if err = tx.Rollback(); err != nil {
			return nil, err
		}
		return nil, row.Err()
	}
	row.Scan(&nabor.ID)

	for _, ingredient := range nabor.Ingredients {
		stmt = `insert into nabor_ingredient (nabor_id, ingredient_id, brutto, price)
		values ($1, $2, $3, $4)`
		row := tx.QueryRowContext(ctx, stmt,
			nabor.ID,
			ingredient.IngredientID,
			ingredient.Brutto,
			ingredient.Price,
		)

		if row.Err() != nil {
			if err = tx.Rollback(); err != nil {
				return nil, err
			}
			return nil, row.Err()
		}
	}
	err = tx.Commit()
	if err != nil {
		return nil, err
	}

	return nabor, nil*/
	err := r.gormDB.Transaction(func(tx *gorm.DB) error {
		ingredients := nabor.Ingredients
		nabor.Ingredients = nil
		for _, shop := range nabor.Shops {
			nabor.ID = 0
			nabor.ShopID = shop
			if err := tx.Create(&nabor).Error; err != nil {
				return err
			}
			for _, ingredient := range ingredients {
				ingredient.NaborID = nabor.NaborID
				ingredient.ShopID = shop
				ingredient.ID = 0
				if err := tx.Create(&ingredient).Error; err != nil {
					return err
				}
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return nabor, nil
}

func (r *IngredientDB) GetAllNabor(filter *model.Filter) ([]*model.NaborOutput, int64, error) {
	nabors := []*model.NaborOutput{}
	res := r.gormDB.Model(&model.Nabor{}).Select("nabors.id, nabors.nabor_id, nabors.shop_id, shops.name as shop_name, nabors.name, nabors.min, nabors.max, nabors.replaces, nabors.deleted").Joins("inner join shops on shops.id = nabors.shop_id").Where("nabors.deleted = ?", false)
	if res.Error != nil {
		return nil, 0, res.Error
	}
	newRes, count, err := filter.FilterResults(res, model.Nabor{}, utils.DefaultPageSize, "", "", "")
	if err != nil {
		return nil, 0, err
	}
	err = newRes.Scan(&nabors).Error
	if err != nil {
		return nil, 0, err
	}
	for _, nabor := range nabors {
		ingredients := []*model.IngredientOutput{}
		err = r.gormDB.Model(&model.IngredientNabor{}).Select("ingredient_nabors.id, ingredient_nabors.nabor_id, ingredient_nabors.ingredient_id, ingredient_nabors.brutto, ingredient_nabors.price, ingredient_nabors.shop_id, ingredients.name, ingredients.image, ingredients.measure").Joins("inner join ingredients on (ingredients.ingredient_id = ingredient_nabors.ingredient_id and ingredients.shop_id = ingredient_nabors.shop_id)").Where("nabor_id = ? and ingredient_nabors.shop_id = ?", nabor.NaborID, nabor.ShopID).Scan(&ingredients).Error
		if err != nil {
			return nil, 0, err
		}
		nabor.NaborIngredient = ingredients
	}
	return nabors, count, nil
}

func (r *IngredientDB) GetNabor(id int) (*model.NaborOutput, error) {
	nabor := &model.NaborOutput{}
	res := r.gormDB.Model(&model.Nabor{}).Select("nabors.id, nabors.nabor_id, nabors.shop_id, shops.name as shop_name, nabors.name, nabors.min, nabors.max, nabors.replaces, nabors.deleted").Joins("inner join shops on shops.id = nabors.shop_id").Where("nabors.id = ?", id).Scan(&nabor)
	if res.Error != nil {
		return nil, res.Error
	}
	err := r.gormDB.Model(&model.IngredientNabor{}).Select("ingredient_nabors.id, ingredient_nabors.nabor_id, ingredient_nabors.ingredient_id, ingredient_nabors.brutto, ingredient_nabors.price, ingredient_nabors.shop_id, ingredients.name, ingredients.image, ingredients.measure").Joins("inner join ingredients on (ingredients.ingredient_id = ingredient_nabors.ingredient_id and ingredients.shop_id = ingredient_nabors.shop_id)").Where("nabor_id = ? and ingredient_nabors.shop_id = ?", nabor.NaborID, nabor.ShopID).Scan(&nabor.NaborIngredient).Error
	if err != nil {
		return nil, err
	}

	return nabor, nil
}

func (r *IngredientDB) UpdateNabor(nabor *model.Nabor) error {
	/*ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()
	stmt := `update nabor set name = $2, min = $3, max = $4 where id=$1`
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	row := tx.QueryRowContext(ctx, stmt,
		nabor.ID,
		nabor.Name,
		nabor.Min,
		nabor.Max,
	)

	if row.Err() != nil {
		if err = tx.Rollback(); err != nil {
			return err
		}
		return row.Err()
	}
	delete := `delete from ingredient_nabors where nabor_id=$1`
	row = tx.QueryRowContext(ctx, delete,
		nabor.ID,
	)
	if row.Err() != nil {
		if err = tx.Rollback(); err != nil {
			return err
		}
		return row.Err()
	}
	for _, ingredient := range nabor.Ingredients {
		query := `insert into nabor_ingredient (nabor_id, ingredient_id, brutto, price)
		values ($1, $2, $3, $4)`
		row := tx.QueryRowContext(ctx, query,
			nabor.ID,
			ingredient.IngredientID,
			ingredient.Brutto,
			ingredient.Price,
		)
		if row.Err() != nil {
			if err = tx.Rollback(); err != nil {
				return err
			}
			return row.Err()
		}
	}
	err = tx.Commit()
	if err != nil {
		return err
	}
	return nil*/

	err := r.gormDB.Transaction(func(tx *gorm.DB) error {
		ingredientNabors := nabor.Ingredients
		nabor.Ingredients = nil
		for _, shop := range nabor.Shops {
			nabor.ID = 0
			nabor.ShopID = shop
			oldNabor := &model.Nabor{}
			res := tx.Model(oldNabor).Where("nabor_id = ? and shop_id = ?", nabor.NaborID, shop).Scan(&oldNabor)
			if res.Error != nil {
				return res.Error
			}
			err := tx.Model(&model.IngredientNabor{}).Where("nabor_id = ? and shop_id = ?", nabor.NaborID, nabor.ShopID).Delete(model.IngredientNabor{}).Error
			if err != nil {
				return err
			}

			if res.RowsAffected == 0 {
				err := tx.Model(nabor).Create(nabor).Error
				if err != nil {
					return err
				}
			} else {
				nabor.ID = oldNabor.ID
				res = tx.Model(nabor).Where("nabor_id = ? and shop_id = ?", nabor.NaborID, shop).Updates(nabor)
				if res.Error != nil {
					return res.Error
				}
			}
			for _, ingredient := range ingredientNabors {
				ingredient.ShopID = shop
				ingredient.NaborID = nabor.NaborID
				ingredient.ID = 0
				err := tx.Model(ingredient).Create(ingredient).Error
				if err != nil {
					return err
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

func (r *IngredientDB) DeleteNabor(id int) error {
	nabor := &model.Nabor{}
	err := r.gormDB.Model(&model.Nabor{}).Where("id = ?", id).Scan(&nabor).Error
	if err != nil {
		return err
	}
	err = r.gormDB.Model(&model.IngredientNabor{}).Where("nabor_id = ? and shop_id = ?", nabor.NaborID, nabor.ShopID).Delete(&model.IngredientNabor{}).Error
	if err != nil {
		return err
	}
	err = r.gormDB.Model(&model.Nabor{}).Where("nabor_id = ? and shop_id = ?", nabor.NaborID, nabor.ShopID).Delete(&model.Nabor{}).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *IngredientDB) GetTechCartByIngredientID(id int) ([]*model.TechCart, error) {
	techs := []*model.TechCart{}
	err := r.gormDB.Model(&model.TechCart{}).Joins("inner join ingredient_tech_carts on ingredient_tech_carts.tech_cart_id = tech_carts.id").Where("ingredient_tech_carts.ingredient_id = ?", id).Scan(&techs).Error
	if err != nil {
		return nil, err
	}
	return techs, nil
}

func (r *IngredientDB) GetPureIngredientByShopID(id, shopID int) (*model.Ingredient, error) {
	ingredient := &model.Ingredient{}
	err := r.gormDB.Model(ingredient).Where("ingredient_id = ?", id).Scan(&ingredient).Error
	if err != nil {
		return nil, err
	}
	return ingredient, nil
}

func (r *IngredientDB) GetIngredientsByTechCart(techCartID int) ([]*model.Ingredient, error) {
	ingredients := []*model.Ingredient{}
	res := r.gormDB.Model(&model.IngredientTechCart{}).Select("ingredient_id as id").Where("tech_cart_id = ?", techCartID).Scan(&ingredients)
	if res.Error != nil {
		return nil, res.Error
	}
	return ingredients, nil
}

func (r *IngredientDB) GetIngredientsIngredientIDByIngredientsID(id int) (int, error) {
	ingredients := &model.Ingredient{}
	res := r.gormDB.Model(&model.Ingredient{}).Where("id = ?", id).Scan(&ingredients)
	if res.Error != nil {
		return 0, res.Error
	}
	if ingredients.IngredientID == 0 {
		return id, nil
	}
	return ingredients.IngredientID, nil
}

func (r *IngredientDB) GetDeletedIngredient(filter *model.Filter) ([]*model.IngredientOutput, int64, error) {
	ingredients := []*model.IngredientOutput{}

	res := r.gormDB
	if filter.Measure == "" {
		res = r.gormDB.Table("ingredients").Select("ingredients.id, ingredients.ingredient_id, ingredients.shop_id, shops.name as shop_name, ingredients.name, category_ingredients.name as category, ingredients.image, ingredients.measure, AVG(sklad_ingredients.cost) as cost").Joins("inner join category_ingredients on ingredients.category = category_ingredients.id inner join sklad_ingredients on ingredients.ingredient_id = sklad_ingredients.ingredient_id inner join sklads on  sklads.id = sklad_ingredients.sklad_id inner join shops on shops.id = ingredients.shop_id").Where("ingredients.deleted = ? and sklads.shop_id IN (?)", true, filter.AccessibleShops).Group("ingredients.id, ingredients.ingredient_id, category_ingredients.name, shops.name")
		if res.Error != nil {
			return nil, 0, res.Error
		}
	} else {
		res = r.gormDB.Table("ingredients").Select("ingredients.id, ingredients.ingredient_id, ingredients.shop_id, shops.name as shop_name, ingredients.name, category_ingredients.name as category, ingredients.image, ingredients.measure, AVG(sklad_ingredients.cost) as cost").Joins("inner join category_ingredients on ingredients.category = category_ingredients.id inner join sklad_ingredients on ingredients.ingredient_id = sklad_ingredients.ingredient_id inner join sklads on  sklads.id = sklad_ingredients.sklad_id inner join shops on shops.id = ingredients.shop_id").Where("ingredients.deleted = ? and sklads.shop_id IN (?) and ingredients.measure = ?", true, filter.AccessibleShops, filter.Measure).Group("ingredients.id, ingredients.ingredient_id, category_ingredients.name, shops.name")

		if res.Error != nil {
			return nil, 0, res.Error
		}
	}

	newRes, count, err := filter.FilterResults(res, &model.Ingredient{}, utils.DefaultPageSize, "", "", "")

	if err != nil {
		return nil, 0, err
	}

	if newRes.Scan(&ingredients).Error != nil {
		return nil, 0, newRes.Error
	}

	for _, ingredient := range ingredients {
		skladIngredient := []*model.SkladIngredient{}
		res := r.gormDB.Model(&model.SkladIngredient{}).Joins("inner join sklads on sklads.id = sklad_ingredients.sklad_id").Where("sklad_ingredients.ingredient_id = ? and sklads.shop_id IN (?)", ingredient.IngredientID, filter.AccessibleShops).Scan(&skladIngredient)
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
			ingredient.Cost = sum / quantity
		} else {
			var cost float32
			res := r.gormDB.Model(&model.ItemPostavka{}).Select("cost").Joins("inner join postavkas on postavkas.id = item_postavkas.postavka_id").Where("item_id = ? and item_postavkas.type = ?", ingredient.IngredientID, utils.TypeIngredient).Order("postavkas.time desc").First(&cost)
			if res.Error != nil {
				if res.Error != gorm.ErrRecordNotFound {
					return nil, 0, res.Error

				}
			}
			if res.RowsAffected == 0 {
				deletedCost := r.gormDB.Model(&model.ItemPostavka{}).Select("cost").Joins("inner join postavkas on postavkas.id = item_postavkas.postavka_id").Where("item_id = ? and item_postavkas.type = ?", ingredient.IngredientID, utils.TypeIngredient).Order("postavkas.time desc").First(&cost)
				if deletedCost.Error != nil {
					if deletedCost.Error != gorm.ErrRecordNotFound {
						return nil, 0, deletedCost.Error
					}
				}
				if deletedCost.RowsAffected == 0 {
					res := r.gormDB.Model(&model.SkladIngredient{}).Select("AVG(sklad_ingredients.cost)").Joins("inner join sklads on sklads.id = sklad_ingredients.sklad_id").Where("sklad_ingredients.ingredient_id = ? and sklads.shop_id IN (?)", ingredient.IngredientID, filter.AccessibleShops).Scan(&cost)
					if res.Error != nil {
						return nil, 0, res.Error
					}
				}
			}
			ingredient.Cost = cost
		}
	}

	return ingredients, count, nil
}

func (r *IngredientDB) RecreateIngredient(ingredient *model.Ingredient) error {
	res := r.gormDB.Model(&model.Ingredient{}).Where("ingredient_id = ? and shop_id = ?", ingredient.IngredientID, ingredient.ShopID).Update("deleted", false)
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return errors.New("ingredient not found")
	}
	return nil
}

func (r *IngredientDB) GetToAddIngredient(filter *model.Filter) ([]*model.IngredientMaster, int64, error) {
	ingredientMaster := []*model.IngredientMaster{}
	res := r.gormDB.Model(&model.IngredientMaster{}).Select("ingredient_masters.*").Joins("left join ingredients on ingredient_masters.id = ingredients.ingredient_id and ingredients.shop_id IN (?)", filter.AccessibleShops).Where("ingredients.id is null and ingredient_masters.deleted = ?", false)
	if res.Error != nil {
		return nil, 0, res.Error
	}
	newRes, count, err := filter.FilterResults(res, &model.IngredientMaster{}, utils.DefaultPageSize, "", "", "")

	if err != nil {
		return nil, 0, err
	}

	if err := newRes.Scan(&ingredientMaster).Error; err != nil {
		return nil, 0, err
	}

	return ingredientMaster, count, nil
}

func (r *IngredientDB) GetIdsOfShopsWhereTheIngredientAlreadyExist(ingredientID int) ([]int, error) {
	shopIds := []int{}
	res := r.gormDB.Model(&model.Ingredient{}).Select("shop_id").Where("ingredient_id = ?", ingredientID).Scan(&shopIds)
	if res.Error != nil {
		return nil, res.Error
	}
	return shopIds, nil
}

func (r *IngredientDB) GetIdsOfShopsWhereTheNaborAlreadyExist(naborID int) ([]int, error) {
	shopIds := []int{}
	res := r.gormDB.Model(&model.Nabor{}).Select("shop_id").Where("nabor_id = ?", naborID).Scan(&shopIds)
	if res.Error != nil {
		return nil, res.Error
	}
	return shopIds, nil
}
