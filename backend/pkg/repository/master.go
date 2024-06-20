package repository

import (
	"database/sql"
	"zebra/model"
	"zebra/utils"

	"gorm.io/gorm"
)

type MasterDB struct {
	db     *sql.DB
	gormDB *gorm.DB
}

func NewMasterDB(db *sql.DB, gormDB *gorm.DB) *MasterDB {
	return &MasterDB{db: db, gormDB: gormDB}
}

func (m *MasterDB) AddTovarMaster(tovar *model.TovarMaster) (*model.TovarMaster, error) {
	if err := m.gormDB.Create(tovar).Error; err != nil {
		return nil, err
	}
	return tovar, nil
}

func (m *MasterDB) AddIngredientMaster(ingredient *model.IngredientMaster) (*model.IngredientMaster, error) {

	if err := m.gormDB.Create(ingredient).Error; err != nil {
		return nil, err
	}
	return ingredient, nil
}

func (m *MasterDB) NormaliseTovars() error {
	tovars := []*model.Tovar{}
	if err := m.gormDB.Find(&tovars).Error; err != nil {
		return err
	}
	for _, tovar := range tovars {
		tovarMaster := &model.TovarMaster{
			ID:       tovar.ID,
			Name:     tovar.Name,
			Category: tovar.Category,
			Image:    tovar.Image,
			Tax:      tovar.Tax,
			Measure:  tovar.Measure,
			Price:    tovar.Price,
			Discount: tovar.Discount,
			Deleted:  tovar.Deleted,
			Status:   utils.MenuStatusApproved,
		}
		if err := m.gormDB.Create(tovarMaster).Error; err != nil {
			return err
		}
		shops := []*model.Shop{}
		if err := m.gormDB.Find(&shops).Error; err != nil {
			return err
		}
		for _, shop := range shops {
			newTovar := &model.Tovar{
				ShopID:    shop.ID,
				TovarID:   tovarMaster.ID,
				Name:      tovarMaster.Name,
				Category:  tovarMaster.Category,
				Image:     tovarMaster.Image,
				Tax:       tovarMaster.Tax,
				Measure:   tovarMaster.Measure,
				Discount:  tovarMaster.Discount,
				IsVisible: true,
				Price:     tovarMaster.Price,
			}
			if err := m.gormDB.Create(newTovar).Error; err != nil {
				return err
			}
		}
	}
	return nil
}

func (m *MasterDB) NormaliseIngredients() error {
	ingredients := []*model.Ingredient{}
	if err := m.gormDB.Find(&ingredients).Error; err != nil {
		return err
	}
	for _, ingredient := range ingredients {
		ingredientMaster := &model.IngredientMaster{
			ID:       ingredient.ID,
			Name:     ingredient.Name,
			Measure:  ingredient.Measure,
			Deleted:  ingredient.Deleted,
			Image:    ingredient.Image,
			Category: ingredient.Category,
			Status:   utils.MenuStatusApproved,
		}
		if err := m.gormDB.Create(ingredientMaster).Error; err != nil {
			return err
		}
		shops := []*model.Shop{}
		if err := m.gormDB.Find(&shops).Error; err != nil {
			return err
		}
		for _, shop := range shops {
			newIngredient := &model.Ingredient{
				ShopID:       shop.ID,
				IngredientID: ingredientMaster.ID,
				Name:         ingredientMaster.Name,
				Measure:      ingredientMaster.Measure,
				Category:     ingredientMaster.Category,
				Image:        ingredientMaster.Image,
				IsVisible:    true,
			}
			if err := m.gormDB.Create(newIngredient).Error; err != nil {
				return err
			}
		}
	}
	return nil
}

func (m *MasterDB) AddTechCartMaster(techCart *model.TechCartMaster) (*model.TechCartMaster, error) {
	ingredients := techCart.Ingredients
	nabors := techCart.Nabor
	techCart.Nabor = nil
	techCart.Ingredients = nil
	if err := m.gormDB.Model(&model.TechCartMaster{}).Create(techCart).Error; err != nil {
		return nil, err
	}
	for _, ingredient := range ingredients {
		ingredient.TechCartID = techCart.ID
		if err := m.gormDB.Model(&model.IngredientTechCartMaster{}).Create(ingredient).Error; err != nil {
			return nil, err
		}
	}
	for _, nabor := range nabors {
		nabor.TechCartID = techCart.ID
		if err := m.gormDB.Model(&model.NaborTechCartMaster{}).Create(nabor).Error; err != nil {
			return nil, err
		}
	}
	return techCart, nil
}

func (m *MasterDB) NormaliseTechCarts() error {
	techCarts := []*model.TechCart{}
	if err := m.gormDB.Find(&techCarts).Error; err != nil {
		return err
	}
	for _, techCart := range techCarts {
		ingredients := []*model.IngredientTechCart{}
		if err := m.gormDB.Model(ingredients).Where("tech_cart_id = ?", techCart.ID).Scan(&ingredients).Error; err != nil {
			return err
		}
		nabors := []*model.NaborTechCart{}
		if err := m.gormDB.Model(nabors).Where("tech_cart_id = ?", techCart.ID).Scan(&nabors).Error; err != nil {
			return err
		}
		techCart.Ingredients = ingredients
		techCart.Nabor = nabors
	}
	for _, techCart := range techCarts {
		ingredientsMaster := []*model.IngredientTechCartMaster{}
		for _, ingredient := range techCart.Ingredients {
			ingredientMaster := &model.IngredientTechCartMaster{
				TechCartID:   ingredient.TechCartID,
				IngredientID: ingredient.IngredientID,
				Brutto:       ingredient.Brutto,
			}
			ingredientsMaster = append(ingredientsMaster, ingredientMaster)
		}
		nabors := []*model.NaborTechCartMaster{}
		for _, nabor := range techCart.Nabor {
			naborMaster := &model.NaborTechCartMaster{
				TechCartID: nabor.TechCartID,
				NaborID:    nabor.NaborID,
			}
			nabors = append(nabors, naborMaster)
		}
		techCartMaster := &model.TechCartMaster{
			ID:       techCart.ID,
			Name:     techCart.Name,
			Deleted:  techCart.Deleted,
			Status:   utils.MenuStatusApproved,
			Category: techCart.Category,
			Tax:      techCart.Tax,
			Price:    techCart.Price,
			Measure:  techCart.Measure,
			Discount: techCart.Discount,
		}
		if err := m.gormDB.Create(techCartMaster).Error; err != nil {
			return err
		}
		for _, ingredient := range ingredientsMaster {
			if err := m.gormDB.Create(ingredient).Error; err != nil {
				return err
			}
		}
		for _, nabor := range nabors {
			if err := m.gormDB.Create(nabor).Error; err != nil {
				return err
			}
		}
		shops := []*model.Shop{}
		if err := m.gormDB.Find(&shops).Error; err != nil {
			return err
		}
		for _, shop := range shops {
			newTechCart := &model.TechCart{
				ShopID:     shop.ID,
				TechCartID: techCartMaster.ID,
				Name:       techCart.Name,
				Deleted:    techCart.Deleted,
				Category:   techCart.Category,
				Tax:        techCart.Tax,
				Price:      techCart.Price,
				Measure:    techCart.Measure,
				Discount:   techCart.Discount,
			}
			if err := m.gormDB.Create(newTechCart).Error; err != nil {
				return err
			}
			for _, ingredient := range ingredientsMaster {
				newIngredient := &model.IngredientTechCart{
					ShopID:       shop.ID,
					TechCartID:   newTechCart.TechCartID,
					IngredientID: ingredient.IngredientID,
					Brutto:       ingredient.Brutto,
				}
				if err := m.gormDB.Create(newIngredient).Error; err != nil {
					return err
				}
			}
			for _, nabor := range nabors {
				newNabor := &model.NaborTechCart{
					ShopID:     shop.ID,
					TechCartID: newTechCart.TechCartID,
					NaborID:    nabor.NaborID,
				}
				if err := m.gormDB.Create(newNabor).Error; err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func (m *MasterDB) AddNaborMaster(nabor *model.NaborMaster) (*model.NaborMaster, error) {
	ingredients := nabor.Ingredients
	nabor.Ingredients = nil
	if err := m.gormDB.Create(nabor).Error; err != nil {
		return nil, err
	}
	for _, ingredient := range ingredients {
		ingredient.NaborID = nabor.ID
		if err := m.gormDB.Create(ingredient).Error; err != nil {
			return nil, err
		}
	}
	nabor.Ingredients = ingredients
	return nabor, nil
}

func (m *MasterDB) CreateNaborMaster(nabor *model.NaborMaster) (*model.NaborMaster, error) {
	nabor.Status = utils.MenuStatusApproved
	if err := m.gormDB.Create(nabor).Error; err != nil {
		return nil, err
	}
	shops := []*model.Shop{}
	res := m.gormDB.Model(&model.Shop{}).Scan(&shops)
	if res.Error != nil {
		return nil, res.Error
	}
	for _, shop := range shops {
		newNabor := &model.Nabor{
			ShopID:      shop.ID,
			NaborID:     nabor.ID,
			Name:        nabor.Name,
			Deleted:     nabor.Deleted,
			Min:         nabor.Min,
			Max:         nabor.Max,
			Ingredients: nil,
			Replaces:    nabor.Replaces,
		}
		if err := m.gormDB.Create(newNabor).Error; err != nil {
			return nil, err
		}

		for _, ingredient := range nabor.Ingredients {
			ingredientNabor := &model.IngredientNabor{
				NaborID:      nabor.ID,
				IngredientID: ingredient.IngredientID,
				Brutto:       ingredient.Brutto,
				Price:        ingredient.Price,
				ShopID:       shop.ID,
			}
			if err := m.gormDB.Create(ingredientNabor).Error; err != nil {
				return nil, err
			}
		}

	}

	return nabor, nil
}

func (m *MasterDB) UpdateNaborMaster(naborMaster *model.NaborMaster) (*model.NaborMaster, error) {
	naborMaster.Status = utils.MenuStatusApproved
	if err := m.gormDB.Save(naborMaster).Error; err != nil {
		return nil, err
	}

	ingredientNaborMaster := naborMaster.Ingredients
	naborMaster.Ingredients = nil
	if err := m.gormDB.Model(&model.IngredientNaborMaster{}).Where("nabor_id = ?", naborMaster.ID).Delete(&model.IngredientNaborMaster{}).Error; err != nil {
		return nil, err
	}

	for _, ingredient := range ingredientNaborMaster {
		ingredient.NaborID = naborMaster.ID
		if err := m.gormDB.Create(ingredient).Error; err != nil {
			return nil, err
		}
	}

	shops := []*model.Shop{}
	res := m.gormDB.Model(&model.Shop{}).Scan(&shops)
	if res.Error != nil {
		return nil, res.Error
	}
	for _, shop := range shops {
		nabor := &model.Nabor{}
		if err := m.gormDB.Where("shop_id = ? AND nabor_id = ?", shop.ID, naborMaster.ID).First(nabor).Error; err != nil {
			return nil, err
		}

		naborIngredients := []*model.IngredientNabor{}
		if err := m.gormDB.Model(&model.IngredientNabor{}).Where("shop_id = ? AND nabor_id = ?", nabor.ShopID, nabor.NaborID).Scan(&naborIngredients).Error; err != nil {
			return nil, err
		}

		for _, ingredient := range naborIngredients {
			if err := m.gormDB.Delete(ingredient).Error; err != nil {
				return nil, err
			}
		}

		newNabor := &model.Nabor{
			ShopID:      shop.ID,
			NaborID:     naborMaster.ID,
			Name:        naborMaster.Name,
			Deleted:     naborMaster.Deleted,
			Min:         naborMaster.Min,
			Max:         naborMaster.Max,
			Ingredients: nil,
			Replaces:    naborMaster.Replaces,
		}
		nabor = newNabor
		if err := m.gormDB.Where("shop_id = ? AND nabor_id = ?", shop.ID, naborMaster.ID).Updates(nabor).Error; err != nil {
			return nil, err
		}
		for _, ingredient := range ingredientNaborMaster {
			ingredientNabor := &model.IngredientNabor{
				NaborID:      naborMaster.ID,
				IngredientID: ingredient.IngredientID,
				Brutto:       ingredient.Brutto,
				Price:        ingredient.Price,
				ShopID:       shop.ID,
			}
			if err := m.gormDB.Create(ingredientNabor).Error; err != nil {
				return nil, err
			}
		}
	}
	return naborMaster, nil
}

func (m *MasterDB) NormaliseNabors() error {
	nabors := []*model.Nabor{}
	if err := m.gormDB.Find(&nabors).Error; err != nil {
		return err
	}
	for _, nabor := range nabors {
		ingredients := []*model.IngredientNabor{}
		if err := m.gormDB.Model(ingredients).Where("nabor_id = ?", nabor.ID).Scan(&ingredients).Error; err != nil {
			return err
		}
		nabor.Ingredients = ingredients
	}
	for _, nabor := range nabors {
		ingredientsNaborsMaster := []*model.IngredientNaborMaster{}

		for _, v := range nabor.Ingredients {
			ingredientsNaborsMaster = append(ingredientsNaborsMaster, &model.IngredientNaborMaster{
				IngredientID: v.IngredientID,
				NaborID:      v.NaborID,
				Brutto:       v.Brutto,
				Price:        v.Price,
			})
		}

		naborMaster := &model.NaborMaster{
			ID:       nabor.ID,
			Name:     nabor.Name,
			Deleted:  nabor.Deleted,
			Status:   utils.MenuStatusApproved,
			Min:      nabor.Min,
			Max:      nabor.Max,
			Replaces: nabor.Replaces,
		}
		if err := m.gormDB.Create(naborMaster).Error; err != nil {
			return err
		}
		for _, ingredientMaster := range ingredientsNaborsMaster {
			if err := m.gormDB.Create(ingredientMaster).Error; err != nil {
				return err
			}
		}
		shops := []*model.Shop{}
		if err := m.gormDB.Find(&shops).Error; err != nil {
			return err
		}
		for _, shop := range shops {
			newNabor := &model.Nabor{
				ShopID:   shop.ID,
				NaborID:  naborMaster.ID,
				Name:     nabor.Name,
				Deleted:  nabor.Deleted,
				Replaces: nabor.Replaces,
				Min:      nabor.Min,
				Max:      nabor.Max,
			}
			if err := m.gormDB.Create(newNabor).Error; err != nil {
				return err
			}
			for _, ingredient := range nabor.Ingredients {
				newIngredient := &model.IngredientNabor{
					ShopID:       shop.ID,
					NaborID:      newNabor.NaborID,
					Name:         ingredient.Name,
					IngredientID: ingredient.IngredientID,
					Brutto:       ingredient.Brutto,
					Price:        ingredient.Price,
					Measure:      ingredient.Measure,
					Image:        ingredient.Image,
				}
				if err := m.gormDB.Create(newIngredient).Error; err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func (r *MasterDB) GetAllTovarMaster(filter *model.Filter) ([]*model.TovarMasterResponse, int64, error) {
	var tovars []*model.TovarMasterResponse
	res := r.gormDB.Model(&model.TovarMaster{}).Select("tovar_masters.id, tovar_masters.name, tovar_masters.image, category_tovars.name as category,  tovar_masters.tax, tovar_masters.measure, AVG(sklad_tovars.cost) as cost, tovar_masters.price, tovar_masters.discount, tovar_masters.status").Joins("inner join category_tovars on tovar_masters.category = category_tovars.id inner join sklad_tovars on tovar_masters.id = sklad_tovars.tovar_id inner join sklads on  sklads.id = sklad_tovars.sklad_id").Where("tovar_masters.deleted = ? and sklads.shop_id IN (?)", false, filter.AccessibleShops).Group("tovar_masters.id, category_tovars.name")
	if res.Error != nil {
		return nil, 0, res.Error
	}
	newRes, count, err := filter.FilterResults(res, &model.TovarMaster{}, utils.DefaultPageSize, "", "", "")

	if err != nil {
		return nil, 0, err
	}

	if newRes.Scan(&tovars).Error != nil {
		return nil, 0, newRes.Error
	}
	for _, tovar := range tovars {
		skladTovars := []*model.SkladTovar{}
		res := r.gormDB.Model(&model.SkladTovar{}).Joins("inner join sklads on sklads.id = sklad_tovars.sklad_id").Where("sklad_tovars.tovar_id = ? and sklads.shop_id IN (?)", tovar.ID, filter.AccessibleShops).Scan(&skladTovars)
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
			res := r.gormDB.Model(&model.ItemPostavka{}).Select("cost").Joins("inner join postavkas on postavkas.id = item_postavkas.postavka_id").Where("item_id = ? and item_postavkas.type = ? and postavkas.deleted = ?", tovar.ID, utils.TypeTovar, false).Order("postavkas.time desc").First(&cost)
			if res.Error != nil {
				if res.Error != gorm.ErrRecordNotFound { //???
					return nil, 0, res.Error
				}
			}
			if res.RowsAffected == 0 {
				deletedCost := r.gormDB.Model(&model.ItemPostavka{}).Select("cost").Joins("inner join postavkas on postavkas.id = item_postavkas.postavka_id").Where("item_id = ? and item_postavkas.type = ?", tovar.ID, utils.TypeTovar).Order("postavkas.time desc").First(&cost)
				if deletedCost.Error != nil {
					if deletedCost.Error != gorm.ErrRecordNotFound {
						return nil, 0, deletedCost.Error
					}
				}
				if deletedCost.RowsAffected == 0 {
					res := r.gormDB.Model(&model.SkladTovar{}).Select("AVG(sklad_tovars.cost)").Joins("inner join sklads on sklads.id = sklad_tovars.sklad_id").Where("sklad_tovars.tovar_id = ? and sklads.shop_id IN (?)", tovar.ID, filter.AccessibleShops).Scan(&cost)
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

func (r *MasterDB) GetTovarMaster(id int) (*model.TovarMasterResponse, error) {
	var tovar model.TovarMasterResponse
	res := r.gormDB.Model(&model.TovarMaster{}).Select("tovar_masters.id, tovar_masters.name, tovar_masters.image, category_tovars.id as category_id, category_tovars.name as category, tovar_masters.tax, tovar_masters.measure, AVG(sklad_tovars.cost) as cost, tovar_masters.price, tovar_masters.discount, tovar_masters.status").Joins("inner join category_tovars on tovar_masters.category = category_tovars.id inner join sklad_tovars on tovar_masters.id = sklad_tovars.tovar_id inner join sklads on  sklads.id = sklad_tovars.sklad_id").Where("tovar_masters.deleted = ? and tovar_masters.id = ?", false, id).Group("tovar_masters.id, category_tovars.name, category_tovars.id")
	if res.Error != nil {
		return nil, res.Error
	}
	if res.Scan(&tovar).Error != nil {
		return nil, res.Error
	}
	skladTovars := []*model.SkladTovar{}
	res = r.gormDB.Model(&model.SkladTovar{}).Joins("inner join sklads on sklads.id = sklad_tovars.sklad_id").Where("sklad_tovars.tovar_id = ?", tovar.ID).Scan(&skladTovars)
	if res.Error != nil {
		return nil, res.Error
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
		res := r.gormDB.Model(&model.ItemPostavka{}).Select("cost").Joins("inner join postavkas on postavkas.id = item_postavkas.postavka_id").Where("item_id = ? and item_postavkas.type = ? and postavkas.deleted = ?", tovar.ID, utils.TypeTovar, false).Order("postavkas.time desc").First(&cost)
		if res.Error != nil {
			if res.Error != gorm.ErrRecordNotFound {
				return nil, res.Error
			}
		}
		if res.RowsAffected == 0 {
			deletedCost := r.gormDB.Model(&model.ItemPostavka{}).Select("cost").Joins("inner join postavkas on postavkas.id = item_postavkas.postavka_id").Where("item_id = ? and item_postavkas.type = ?", tovar.ID, utils.TypeTovar).Order("postavkas.time desc").First(&cost)
			if deletedCost.Error != nil {
				if deletedCost.Error != gorm.ErrRecordNotFound {
					return nil, deletedCost.Error
				}
			}
			if deletedCost.RowsAffected == 0 {
				res := r.gormDB.Model(&model.SkladTovar{}).Select("AVG(sklad_tovars.cost)").Joins("inner join sklads on sklads.id = sklad_tovars.sklad_id").Where("sklad_tovars.tovar_id = ?", tovar.ID).Scan(&cost)
				if res.Error != nil {
					return nil, res.Error
				}
			}
		}
		tovar.Cost = cost
	}
	tovar.Profit = tovar.Price - tovar.Cost
	if tovar.Cost != 0 {
		tovar.Margin = tovar.Profit * 100 / tovar.Cost
	}
	return &tovar, nil
}

func (r *MasterDB) GetAllIngredientMaster(filter *model.Filter) ([]*model.IngredientMasterResponse, int64, error) {
	ingredients := []*model.IngredientMasterResponse{}
	res := r.gormDB.Model(&model.IngredientMaster{}).Select("ingredient_masters.id, ingredient_masters.name, category_ingredients.name as category, ingredient_masters.image, ingredient_masters.measure, AVG(sklad_ingredients.cost) as cost, ingredient_masters.status").Joins("inner join category_ingredients on ingredient_masters.category = category_ingredients.id inner join sklad_ingredients on ingredient_masters.id = sklad_ingredients.ingredient_id inner join sklads on  sklads.id = sklad_ingredients.sklad_id").Where("ingredient_masters.deleted = ?", false).Group("ingredient_masters.id,  category_ingredients.name")
	if res.Error != nil {
		return nil, 0, res.Error
	}
	newRes, count, err := filter.FilterResults(res, &model.IngredientMaster{}, utils.DefaultPageSize, "", "", "")

	if err != nil {
		return nil, 0, err
	}

	if newRes.Scan(&ingredients).Error != nil {
		return nil, 0, newRes.Error
	}
	for _, ingredient := range ingredients {
		skladIngredient := []*model.SkladIngredient{}
		res := r.gormDB.Model(&model.SkladIngredient{}).Joins("inner join sklads on sklads.id = sklad_ingredients.sklad_id").Where("sklad_ingredients.ingredient_id = ?", ingredient.ID).Scan(&skladIngredient)
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
			res := r.gormDB.Model(&model.ItemPostavka{}).Select("cost").Joins("inner join postavkas on postavkas.id = item_postavkas.postavka_id").Where("item_id = ? and item_postavkas.type = ?", ingredient.ID, utils.TypeIngredient).Order("postavkas.time desc").First(&cost)
			if res.Error != nil {
				if res.Error != gorm.ErrRecordNotFound {
					return nil, 0, res.Error

				}
			}
			if res.RowsAffected == 0 {
				deletedCost := r.gormDB.Model(&model.ItemPostavka{}).Select("cost").Joins("inner join postavkas on postavkas.id = item_postavkas.postavka_id").Where("item_id = ? and item_postavkas.type = ?", ingredient.ID, utils.TypeIngredient).Order("postavkas.time desc").First(&cost)
				if deletedCost.Error != nil {
					if deletedCost.Error != gorm.ErrRecordNotFound {
						return nil, 0, deletedCost.Error
					}
				}
				if deletedCost.RowsAffected == 0 {
					res := r.gormDB.Model(&model.SkladIngredient{}).Select("COALESCE(AVG(sklad_ingredients.cost), 0.0)").Joins("inner join sklads on sklads.id = sklad_ingredients.sklad_id").Where("sklad_ingredients.ingredient_id = ? and sklads.shop_id IN (?)", ingredient.ID, filter.AccessibleShops).Scan(&cost)
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
func (r *MasterDB) GetAllTechCartMaster(filter *model.Filter) ([]*model.TechCartMasterResponse, int64, error) {
	techs := []*model.TechCartMasterResponse{}
	res := r.gormDB.Model(&model.TechCartMaster{}).Select("tech_cart_masters.id, tech_cart_masters.name, category_tovars.name as category, tech_cart_masters.image, tech_cart_masters.tax, tech_cart_masters.measure, tech_cart_masters.price, tech_cart_masters.discount").Joins("inner join category_tovars on tech_cart_masters.category = category_tovars.id").Where("tech_cart_masters.deleted = ?", false)
	if res.Error != nil {
		return nil, 0, res.Error
	}

	newRes, count, err := filter.FilterResults(res, model.TechCartMaster{}, utils.DefaultPageSize, "", "", "")
	if err != nil {
		return nil, 0, err
	}
	if newRes.Scan(&techs).Error != nil {
		return nil, 0, newRes.Error
	}

	for _, tech := range techs {
		techIngredients := []*model.IngredientTechCartMaster{}
		res := r.gormDB.Model(&model.IngredientTechCartMaster{}).Select("*").Where("ingredient_tech_cart_masters.tech_cart_id = ?", tech.ID).Scan(&techIngredients)
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
	return techs, count, nil
}
func (r *MasterDB) GetIngredientMaster(id int) (*model.IngredientMasterResponse, error) {
	ingredient := &model.IngredientMasterResponse{}
	res := r.gormDB.Model(&model.IngredientMaster{}).Select("ingredient_masters.id, ingredient_masters.name, category_ingredients.id as category_id, category_ingredients.name as category, ingredient_masters.image, ingredient_masters.measure, AVG(sklad_ingredients.cost) as cost, ingredient_masters.status").Joins("inner join category_ingredients on ingredient_masters.category = category_ingredients.id inner join sklad_ingredients on ingredient_masters.id = sklad_ingredients.ingredient_id inner join sklads on  sklads.id = sklad_ingredients.sklad_id").Where("ingredient_masters.deleted = ? and ingredient_masters.id = ?", false, id).Group("ingredient_masters.id,  category_ingredients.name, category_ingredients.id")
	if res.Error != nil {
		return nil, res.Error
	}
	if res.Scan(&ingredient).Error != nil {
		return nil, res.Error
	}
	skladIngredient := []*model.SkladIngredient{}
	res = r.gormDB.Model(&model.SkladIngredient{}).Joins("inner join sklads on sklads.id = sklad_ingredients.sklad_id").Where("sklad_ingredients.ingredient_id = ?", ingredient.ID).Scan(&skladIngredient)
	if res.Error != nil {
		return nil, res.Error
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
		res := r.gormDB.Model(&model.ItemPostavka{}).Select("cost").Joins("inner join postavkas on postavkas.id = item_postavkas.postavka_id").Where("item_postavkas.item_id = ? and item_postavkas.type = ?", ingredient.ID, utils.TypeIngredient).Order("postavkas.time desc").First(&cost)
		if res.Error != nil {
			if res.Error != gorm.ErrRecordNotFound {
				return nil, res.Error

			}
		}
		if res.RowsAffected == 0 {
			deletedCost := r.gormDB.Model(&model.ItemPostavka{}).Select("cost").Joins("inner join postavkas on postavkas.id = item_postavkas.postavka_id").Where("item_postavkas.item_id = ? and item_postavkas.type = ?", ingredient.ID, utils.TypeIngredient).Order("postavkas.time desc").First(&cost)
			if deletedCost.Error != nil {
				if deletedCost.Error != gorm.ErrRecordNotFound {
					return nil, deletedCost.Error
				}
			}
			if deletedCost.RowsAffected == 0 {
				res := r.gormDB.Model(&model.SkladIngredient{}).Select("AVG(sklad_ingredients.cost)").Joins("inner join sklads on sklads.id = sklad_ingredients.sklad_id").Where("sklad_ingredients.ingredient_id = ?", ingredient.ID).Scan(&cost)
				if res.Error != nil {
					return nil, res.Error
				}
			}
		}
		ingredient.Cost = cost
	}
	return ingredient, nil

}
func (r *MasterDB) GetTechCartMaster(id int) (*model.TechCartMasterResponse, error) {
	tech := &model.TechCartMasterResponse{}
	err := r.gormDB.Model(&model.TechCartMaster{}).Select("tech_cart_masters.id, tech_cart_masters.name, category_tovars.id as category_id, category_tovars.name as category, tech_cart_masters.image, tech_cart_masters.tax, tech_cart_masters.measure, tech_cart_masters.price, tech_cart_masters.discount, tech_cart_masters.status").Joins("inner join category_tovars on tech_cart_masters.category = category_tovars.id").Where("tech_cart_masters.deleted = ? and tech_cart_masters.id = ?", false, id).Scan(&tech).Error
	if err != nil {
		return nil, err
	}

	ingredients := []*model.IngredientNumOutput{}
	ingre := r.gormDB.Model(&model.IngredientTechCartMaster{}).Select("ingredient_masters.id, ingredient_masters.name, ingredient_masters.measure, ingredient_masters.image, ingredient_tech_cart_masters.brutto, AVG(sklad_ingredients.cost) as cost").Joins("inner join ingredient_masters on ingredient_masters.id = ingredient_tech_cart_masters.ingredient_id inner join sklad_ingredients on ingredient_masters.id = sklad_ingredients.ingredient_id").Where("ingredient_tech_cart_masters.tech_cart_id = ?", tech.ID).Group("ingredient_masters.id, ingredient_tech_cart_masters.brutto").Scan(&ingredients)
	if ingre.Error != nil {
		return nil, ingre.Error
	}
	var ingredientCost float32 = 0
	var ingredientSum float32 = 0

	for _, ingredient := range ingredients {
		skladIngredient := []*model.SkladIngredient{}
		res := r.gormDB.Model(&model.SkladIngredient{}).Joins("inner join sklads on sklads.id = sklad_ingredients.sklad_id").Where("sklad_ingredients.ingredient_id = ?", ingredient.ID).Scan(&skladIngredient)
		if res.Error != nil {
			return nil, res.Error
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
			res := r.gormDB.Model(&model.ItemPostavka{}).Select("cost").Joins("inner join postavkas on postavkas.id = item_postavkas.postavka_id").Where("item_id = ? and item_postavkas.type = ? and postavkas.deleted = ?", ingredient.ID, utils.TypeIngredient, false).Order("postavkas.time desc").First(&cost)
			if res.Error != nil {
				if res.Error != gorm.ErrRecordNotFound {
					return nil, res.Error
				}
			}
			if res.RowsAffected == 0 {
				deletedCost := r.gormDB.Model(&model.ItemPostavka{}).Select("cost").Joins("inner join postavkas on postavkas.id = item_postavkas.postavka_id").Where("item_id = ? and item_postavkas.type = ?", ingredient.ID, utils.TypeIngredient).Order("postavkas.time desc").First(&cost)
				if deletedCost.Error != nil {
					if deletedCost.Error != gorm.ErrRecordNotFound {
						return nil, deletedCost.Error
					}
				}
				if deletedCost.RowsAffected == 0 {
					res := r.gormDB.Model(&model.SkladIngredient{}).Select("AVG(sklad_ingredients.cost)").Joins("inner join sklads on sklads.id = sklad_ingredients.sklad_id").Where("sklad_ingredients.ingredient_id = ?", ingredient.ID).Scan(&cost)
					if res.Error != nil {
						return nil, res.Error
					}
				}
			}
			ingredient.Cost = cost
		}
		ingredientSum += ingredientCost * ingredient.Brutto
		ingredient.Netto = ingredient.Brutto
	}
	tech.Ingredients = ingredients
	tech.Cost = ingredientSum

	nabors := []*model.NaborInfo{}
	naborOutput := []*model.NaborOutput{}
	nab := r.gormDB.Model(&model.NaborTechCartMaster{}).Select("nabor_tech_cart_masters.nabor_id as id, nabor_masters.name as name, nabor_masters.min as min, nabor_masters.max as max").Joins("inner join nabor_masters on nabor_masters.id = nabor_tech_cart_masters.nabor_id").Where("nabor_tech_cart_masters.tech_cart_id = ?", tech.ID).Scan(&nabors)
	if nab.Error != nil {
		return nil, nab.Error
	}
	for _, nabor := range nabors {
		naborIngredients := []*model.IngredientOutput{}

		nabingre := r.gormDB.Model(&model.IngredientNaborMaster{}).Select("ingredient_masters.id, ingredient_masters.name, ingredient_masters.category as category_id, category_ingredients.name as category,  ingredient_nabor_masters.brutto, ingredient_nabor_masters.price, ingredient_masters.measure, AVG(sklad_ingredients.cost)").Joins("inner join ingredient_masters on ingredient_masters.id = ingredient_nabor_masters.ingredient_id inner join category_ingredients on category_ingredients.id = ingredient_masters.category inner join sklad_ingredients on sklad_ingredients.ingredient_id = ingredient_masters.id").Where("nabor_id = ?", nabor.ID).Group("ingredient_masters.id, category_ingredients.name, ingredient_nabor_masters.id").Scan(&naborIngredients)
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

func (r *MasterDB) GetAllNaborMaster(filter *model.Filter) ([]*model.NaborMasterOutput, int64, error) {
	nabors := []*model.NaborMasterOutput{}
	res := r.gormDB.Model(&model.NaborMaster{}).Where("deleted = ?", false)
	if res.Error != nil {
		return nil, 0, res.Error
	}
	newRes, count, err := filter.FilterResults(res, model.NaborMaster{}, utils.DefaultPageSize, "", "", "")
	if err != nil {
		return nil, 0, err
	}
	err = newRes.Scan(&nabors).Error
	if err != nil {
		return nil, 0, err
	}
	for _, nabor := range nabors {
		err = r.gormDB.Model(&model.IngredientNaborMaster{}).Select("ingredient_nabor_masters.id, ingredient_nabor_masters.nabor_id,ingredient_nabor_masters.ingredient_id, ingredient_nabor_masters.brutto, ingredient_nabor_masters.price, ingredient_masters.name, ingredient_masters.image, ingredient_masters.measure").Joins("inner join ingredient_masters on ingredient_nabor_masters.ingredient_id = ingredient_masters.id").Where("nabor_id = ?", nabor.ID).Scan(&nabor.Ingredients).Error
		if err != nil {
			return nil, 0, err
		}
	}
	return nabors, count, nil
}

func (r *MasterDB) GetNaborMaster(id int) (*model.NaborMasterOutput, error) {
	nabor := &model.NaborMasterOutput{}
	err := r.gormDB.Model(&model.NaborMaster{}).Where("id = ?", id).Scan(&nabor).Error
	if err != nil {
		return nil, err
	}
	err = r.gormDB.Model(&model.IngredientNaborMaster{}).Select("ingredient_nabor_masters.id, ingredient_nabor_masters.nabor_id,ingredient_nabor_masters.ingredient_id, ingredient_nabor_masters.brutto, ingredient_nabor_masters.price, ingredient_masters.name, ingredient_masters.image, ingredient_masters.measure").Joins("inner join ingredient_masters on ingredient_nabor_masters.ingredient_id = ingredient_masters.id").Where("nabor_id = ?", nabor.ID).Scan(&nabor.Ingredients).Error
	if err != nil {
		return nil, err
	}

	return nabor, nil
}

func (r *MasterDB) UpdateTovarMaster(tovar *model.TovarMaster) error {
	return r.gormDB.Model(tovar).Where("id = ?", tovar.ID).Updates(tovar).Error
}

func (r *MasterDB) UpdateIngredientMaster(ingredient *model.IngredientMaster) error {
	return r.gormDB.Model(ingredient).Where("id = ?", ingredient.ID).Updates(ingredient).Error
}
func (r *MasterDB) UpdateTechCartMaster(techCart *model.TechCartMaster) error {
	err := r.gormDB.Transaction(func(tx *gorm.DB) error {
		ingredients := techCart.Ingredients
		nabors := techCart.Nabor
		techCart.Ingredients = nil
		techCart.Nabor = nil
		err := r.gormDB.Debug().Select("*").Model(techCart).Where("id = ?", techCart.ID).Updates(techCart).Error
		if err != nil {
			return err
		}
		err = r.gormDB.Model(&model.IngredientTechCartMaster{}).Where("tech_cart_id = ?", techCart.ID).Delete(&model.IngredientTechCartMaster{}).Error
		if err != nil {
			return err
		}
		techCart.Ingredients = ingredients
		for _, ingredient := range techCart.Ingredients {
			ingredient.TechCartID = techCart.ID
			err = r.gormDB.Create(ingredient).Error
			if err != nil {
				return err
			}
		}
		err = r.gormDB.Model(&model.NaborTechCartMaster{}).Where("tech_cart_id = ?", techCart.ID).Delete(&model.NaborTechCartMaster{}).Error
		if err != nil {
			return err
		}
		techCart.Nabor = nabors
		for _, nabor := range techCart.Nabor {
			nabor.TechCartID = techCart.ID
			err = r.gormDB.Create(nabor).Error
			if err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return err
	}
	err = r.UpdateTechCartsMaster(techCart)
	if err != nil {
		return err
	}
	return nil
}

func (r *MasterDB) UpdateTechCartsMaster(techCartMaster *model.TechCartMaster) error {
	shops := []*model.Shop{}
	res := r.gormDB.Model(&model.Shop{}).Scan(&shops)
	if res.Error != nil {
		return res.Error
	}
	for _, shop := range shops {
		techCart := &model.TechCart{}
		err := r.gormDB.Model(&model.TechCart{}).Where("tech_cart_id = ? and shop_id = ?", techCartMaster.ID, shop.ID).Scan(&techCart).Error
		if err != nil {
			return err
		}
		techCart.Name = techCartMaster.Name
		techCart.Category = techCartMaster.Category
		techCart.Image = techCartMaster.Image
		techCart.Tax = techCartMaster.Tax
		techCart.Measure = techCartMaster.Measure
		techCart.Price = techCartMaster.Price
		techCart.Discount = techCartMaster.Discount
		techCart.ShopID = shop.ID
		techCart.Ingredients = nil
		techCart.Nabor = nil
		err = r.gormDB.Model(&model.TechCart{}).Select("*").Omit("deleted, is_visible").Where("tech_cart_id = ? and shop_id = ?", techCartMaster.ID, shop.ID).Updates(&techCart).Error
		if err != nil {
			return err
		}

		err = r.gormDB.Model(&model.IngredientTechCart{}).Where("tech_cart_id = ? and shop_id = ?", techCartMaster.ID, shop.ID).Delete(&model.IngredientTechCart{}).Error
		if err != nil {
			return err
		}

		for _, ingredient := range techCartMaster.Ingredients {
			ingredientTechCart := &model.IngredientTechCart{}
			ingredientTechCart.TechCartID = ingredient.TechCartID
			ingredientTechCart.ShopID = shop.ID
			ingredientTechCart.IngredientID = ingredient.IngredientID
			ingredientTechCart.Brutto = ingredient.Brutto
			err = r.gormDB.Create(ingredientTechCart).Error
			if err != nil {
				return err
			}
		}

		err = r.gormDB.Model(&model.NaborTechCart{}).Where("tech_cart_id = ? and shop_id = ?", techCartMaster.ID, shop.ID).Delete(&model.NaborTechCart{}).Error
		if err != nil {
			return err
		}

		for _, nabor := range techCartMaster.Nabor {
			naborTechCart := &model.NaborTechCart{}
			naborTechCart.TechCartID = nabor.TechCartID
			naborTechCart.ShopID = shop.ID
			naborTechCart.NaborID = nabor.NaborID
			err = r.gormDB.Create(naborTechCart).Error
			if err != nil {
				return err
			}
		}

	}
	return nil
}

func (r *MasterDB) DeleteTovarMaster(id int) error {
	res := r.gormDB.Model(&model.TovarMaster{}).Where("id = ?", id).Update("deleted", true)
	if res.Error != nil {
		return res.Error
	}
	return nil
}
func (r *MasterDB) DeleteIngredientMaster(id int) error {
	res := r.gormDB.Model(&model.IngredientMaster{}).Where("id = ?", id).Update("deleted", true)
	if res.Error != nil {
		return res.Error
	}
	return nil
}
func (r *MasterDB) DeleteTechCartMaster(id int) error {
	res := r.gormDB.Model(&model.TechCartMaster{}).Where("id = ?", id).Delete(&model.TechCartMaster{})
	if res.Error != nil {
		return res.Error
	}
	res = r.gormDB.Model(&model.IngredientTechCartMaster{}).Where("tech_cart_id = ?", id).Delete(&model.IngredientTechCartMaster{})
	if res.Error != nil {
		return res.Error
	}
	res = r.gormDB.Model(&model.NaborTechCartMaster{}).Where("tech_cart_id = ?", id).Delete(&model.NaborTechCartMaster{})
	if res.Error != nil {
		return res.Error
	}
	res = r.gormDB.Model(&model.TechCart{}).Where("tech_cart_id = ?", id).Delete(&model.TechCart{})
	if res.Error != nil {
		return res.Error
	}
	res = r.gormDB.Model(&model.IngredientTechCart{}).Where("tech_cart_id = ?", id).Delete(&model.IngredientTechCart{})
	if res.Error != nil {
		return res.Error
	}
	res = r.gormDB.Model(&model.NaborTechCart{}).Where("tech_cart_id = ?", id).Delete(&model.NaborTechCart{})
	if res.Error != nil {
		return res.Error
	}
	return nil
}

func (r *MasterDB) DeleteNaborMaster(id int) error {
	res := r.gormDB.Model(&model.NaborMaster{}).Where("id = ?", id).Delete(&model.NaborMaster{})
	if res.Error != nil {
		return res.Error
	}
	res = r.gormDB.Model(&model.IngredientNaborMaster{}).Where("nabor_id = ?", id).Delete(&model.IngredientNaborMaster{})
	if res.Error != nil {
		return res.Error
	}
	res = r.gormDB.Model(&model.NaborTechCartMaster{}).Where("nabor_id = ?", id).Delete(&model.NaborTechCartMaster{})
	if res.Error != nil {
		return res.Error
	}
	res = r.gormDB.Model(&model.Nabor{}).Where("nabor_id = ?", id).Delete(&model.Nabor{})
	if res.Error != nil {
		return res.Error
	}
	res = r.gormDB.Model(&model.IngredientNabor{}).Where("nabor_id = ?", id).Delete(&model.IngredientNabor{})
	if res.Error != nil {
		return res.Error
	}
	res = r.gormDB.Model(&model.NaborTechCart{}).Where("nabor_id = ?", id).Delete(&model.NaborTechCart{})
	if res.Error != nil {
		return res.Error
	}
	return nil
}

func (r *MasterDB) GetAllTovarMasterIds() ([]int, error) {
	ids := []int{}
	res := r.gormDB.Model(&model.TovarMaster{}).Where("deleted = ?", false).Pluck("id", &ids)
	if res.Error != nil {
		return nil, res.Error
	}
	return ids, nil
}

func (r *MasterDB) GetAllTechCartsMasterIds() ([]int, error) {
	ids := []int{}
	res := r.gormDB.Model(&model.TechCartMaster{}).Where("deleted = ?", false).Pluck("id", &ids)
	if res.Error != nil {
		return nil, res.Error
	}
	return ids, nil
}

func (r *MasterDB) ConfirmTovarMaster(id int) (*model.TovarMaster, error) {
	tovar := &model.TovarMaster{}
	res := r.gormDB.Model(&model.TovarMaster{}).Where("id = ?", id).Scan(&tovar)
	if res.Error != nil {
		return nil, res.Error
	}
	tovar.Status = utils.MenuStatusApproved
	res = r.gormDB.Model(&model.TovarMaster{}).Where("id = ?", id).Updates(&tovar)
	if res.Error != nil {
		return nil, res.Error
	}
	return tovar, nil
}

func (r *MasterDB) RejectTovarMaster(id int) error {
	res := r.gormDB.Model(&model.TovarMaster{}).Where("id = ?", id).Update("status", utils.MenuStatusRejected)
	if res.Error != nil {
		return res.Error
	}
	return nil
}

func (r *MasterDB) ConfirmTechCartMaster(id int) (*model.TechCartMaster, error) {
	tech := &model.TechCartMaster{}
	res := r.gormDB.Model(&model.TechCartMaster{}).Where("id = ?", id).Scan(&tech)
	if res.Error != nil {
		return nil, res.Error
	}
	tech.Status = utils.MenuStatusApproved
	res = r.gormDB.Model(&model.TechCartMaster{}).Where("id = ?", id).Updates(&tech)
	if res.Error != nil {
		return nil, res.Error
	}
	ingredients := []*model.IngredientTechCartMaster{}
	res = r.gormDB.Model(&model.IngredientTechCartMaster{}).Where("tech_cart_id = ?", id).Scan(&ingredients)
	if res.Error != nil {
		return nil, res.Error
	}
	tech.Ingredients = ingredients
	nabors := []*model.NaborTechCartMaster{}
	res = r.gormDB.Model(&model.NaborTechCartMaster{}).Where("tech_cart_id = ?", id).Scan(&nabors)
	if res.Error != nil {
		return nil, res.Error
	}
	tech.Nabor = nabors
	return tech, nil
}

func (r *MasterDB) RejectTechCartMaster(id int) error {
	res := r.gormDB.Model(&model.TechCartMaster{}).Where("id = ?", id).Update("status", utils.MenuStatusRejected)
	if res.Error != nil {
		return res.Error
	}
	return nil
}

func (r *MasterDB) ConfirmIngredientMaster(id int) (*model.IngredientMaster, error) {
	ingredient := &model.IngredientMaster{}
	res := r.gormDB.Model(&model.IngredientMaster{}).Where("id = ?", id).Scan(&ingredient)
	if res.Error != nil {
		return nil, res.Error
	}
	ingredient.Status = utils.MenuStatusApproved
	res = r.gormDB.Model(&model.IngredientMaster{}).Where("id = ?", id).Updates(&ingredient)
	if res.Error != nil {
		return nil, res.Error
	}
	return ingredient, nil
}

func (r *MasterDB) RejectIngredientMaster(id int) error {
	res := r.gormDB.Model(&model.IngredientMaster{}).Where("id = ?", id).Update("status", utils.MenuStatusRejected)
	if res.Error != nil {
		return res.Error
	}
	return nil
}

func (r *MasterDB) ConfirmNaborMaster(id int) (*model.NaborMaster, error) {
	nabor := &model.NaborMaster{}
	res := r.gormDB.Model(&model.NaborMaster{}).Where("id = ?", id).Scan(&nabor)
	if res.Error != nil {
		return nil, res.Error
	}
	nabor.Status = utils.MenuStatusApproved
	res = r.gormDB.Model(&model.NaborMaster{}).Where("id = ?", id).Updates(&nabor)
	if res.Error != nil {
		return nil, res.Error
	}
	ingredients := []*model.IngredientNaborMaster{}
	res = r.gormDB.Model(&model.IngredientNaborMaster{}).Where("nabor_id = ?", id).Scan(&ingredients)
	if res.Error != nil {
		return nil, res.Error
	}
	nabor.Ingredients = ingredients
	return nabor, nil
}

func (r *MasterDB) RejectNaborMaster(id int) error {
	res := r.gormDB.Model(&model.NaborMaster{}).Where("id = ?", id).Update("status", utils.MenuStatusRejected)
	if res.Error != nil {
		return res.Error
	}
	return nil
}
