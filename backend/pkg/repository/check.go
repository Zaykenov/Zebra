package repository

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"math"
	"time"
	"zebra/model"
	"zebra/utils"

	"gorm.io/gorm"
)

type CheckDB struct {
	db     *sql.DB
	gormDB *gorm.DB
}

func NewCheckDB(db *sql.DB, gormDB *gorm.DB) *CheckDB {
	return &CheckDB{db: db, gormDB: gormDB}
}

func (r *CheckDB) UpdateToSend(check *model.SendToTis) error {
	err := r.gormDB.Save(check).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *CheckDB) DeleteToSend(id int) error {
	if err := r.gormDB.Delete(&model.SendToTis{}, id).Error; err != nil {
		return err
	}
	return nil
}

func (r *CheckDB) GetUnsendTisCheck() ([]*model.SendToTis, error) {
	checks := []*model.SendToTis{}
	if err := r.gormDB.Model(&model.SendToTis{}).Where("retry_count < ? and status != ? and cassa_type = ?", 5, utils.StatusSuccess, utils.CassaTypeTis).Limit(1000).Scan(&checks).Error; err != nil {
		return nil, err
	}
	return checks, nil
}

func (r *CheckDB) AddCheckToSend(check *model.SendToTis) error {
	if err := r.gormDB.Create(check).Error; err != nil {
		return err
	}
	return nil
}

func (r *CheckDB) GetUnsendCheck() ([]*model.ErrorCheck, error) {
	checks := []*model.ErrorCheck{}
	if err := r.gormDB.Model(&model.Check{}).Scan(&checks).Error; err != nil {
		return nil, err
	}
	return checks, nil
}

func (r *CheckDB) GetModificatorsCost(modif []*model.Modificator) (float32, error) {
	ingredients := &model.Ingredient{}
	var cost float32
	for _, modificator := range modif {
		if err := r.gormDB.Model(&ingredients).Where("ingredient_id = ?", modificator.ID).Scan(ingredients).Error; err != nil {
			return -1, err
		}
		cost = cost + modificator.Brutto*ingredients.Cost
	}
	return cost, nil

}

func (r *CheckDB) AddCheck(check *model.Check) (*model.CheckResponse, error) {
	// checkTechCartResponses := []*model.CheckTechCartResponse{}
	// for i := 0; i < len(check.TechCart); i++ {
	// 	item := check.TechCart[i]
	// 	techCartsResponse := &model.CheckTechCartResponse{
	// 		ID:           item.ID,
	// 		CheckID:      item.CheckID,
	// 		TechCartID:   item.TechCartID,
	// 		TechCartName: item.TechCartName,
	// 		Quantity:     item.Quantity,
	// 		Cost:         item.Cost,
	// 		Price:        item.Price,
	// 		Discount:     item.Discount,
	// 		Comments:     item.Comments,
	// 	}
	// 	ingredients := []*model.IngredientNumOutput{}

	// 	modificators := []*model.IngredientOutput{}

	// 	ingredient, err := json.Marshal(ingredients)
	// 	if err != nil {
	// 		return nil, err
	// 	}
	// 	modificator, err := json.Marshal(modificators)
	// 	if err != nil {
	// 		return nil, err
	// 	}
	// 	techCartsResponse.Ingredients = string(ingredient)
	// 	techCartsResponse.Modificators = string(modificator)
	// 	checkTechCartResponses = append(checkTechCartResponses, techCartsResponse)

	// 	for _, ingredient := range check.TechCart[i].ExpenceIngredient {
	// 		if ingredient.Type == utils.TypeIngredient {
	// 			ing := &model.IngredientNumOutput{}
	// 			res := r.gormDB.Table("ingredients").Select("ingredients.id,  ingredients.name, ingredients.measure").Where("ingredient.id", ingredient.IngredientID).Scan(&ing)
	// 			if res.Error != nil {
	// 				return nil, res.Error
	// 			}
	// 			ing.Brutto = ingredient.Quantity
	// 			ingredients = append(ingredients, ing)
	// 		} else if ingredient.Type == utils.TypeModificator {
	// 			mod := &model.IngredientOutput{}
	// 			res := r.gormDB.Table("ingredients").Select("ingredients.id,  ingredients.name, ingredients.measure").Where("ingredient.id", ingredient.IngredientID).Scan(&mod)
	// 			if res.Error != nil {
	// 				return nil, res.Error
	// 			}
	// 			mod.Brutto = ingredient.Quantity
	// 			modificators = append(modificators, mod)
	// 		}
	// 	}
	// }

	if check.Stolik != 0 {
		stolik := &model.Stolik{}
		res := r.gormDB.Model(&model.Stolik{}).Where("shop_id = ? and stolik_id = ?", check.ShopID, check.Stolik).Scan(&stolik)
		if res.Error != nil {
			return nil, res.Error
		}
		stolik.Empty = false
		res = r.gormDB.Model(&model.Stolik{}).Select("*").Where("shop_id = ? and stolik_id = ?", check.ShopID, check.Stolik).Updates(stolik)
		if res.Error != nil {
			return nil, res.Error
		}
	}

	checkResponse := &model.CheckResponse{
		ID:               check.ID,
		MobileUserID:     check.MobileUserID,
		SkladID:          check.SkladID,
		ShopID:           check.ShopID,
		WorkerID:         check.WorkerID,
		Worker:           check.Worker,
		Opened_at:        check.Opened_at,
		Closed_at:        check.Closed_at,
		IdempotencyKey:   check.IdempotencyKey,
		Cash:             check.Cash,
		Card:             check.Card,
		Sum:              check.Sum,
		Cost:             check.Cost,
		Status:           check.Status,
		Payment:          check.Payment,
		Discount:         check.Discount,
		DiscountPercent:  check.DiscountPercent,
		Tovar:            check.Tovar,
		ModificatorCheck: check.ModificatorCheck,
		Comment:          check.Comment,
		ServicePercent:   check.ServicePercent,
		ServiceSum:       check.ServiceSum,
		Stolik:           check.Stolik,
		//TechCart:         checkTechCart.TechCart,
	}
	err := r.gormDB.Transaction(func(tx *gorm.DB) error {

		for _, tovar := range check.Tovar {
			tovar.ExpenceTovar.Status = utils.StatusOpened
		}
		for _, techCart := range check.TechCart {
			for _, ingredients := range techCart.ExpenceIngredient {
				ingredients.Status = utils.StatusOpened
			}
		}
		if err := tx.Debug().Create(check).Error; err != nil {
			return err
		}
		return nil
	},
	)
	if err != nil {
		return nil, err
	}
	checkTechCart, err := r.GetCheckByID(check.ID)
	if err != nil {
		return nil, err
	}
	checkResponse.TechCart = checkTechCart.TechCart

	return checkResponse, nil
}
func (repo *CheckDB) ConstructResponse(r *model.Check) *model.CheckResponse {
	newTechCarts := []*model.CheckTechCartResponse{}
	for _, techCart := range r.TechCart {
		newTechCart := &model.CheckTechCartResponse{
			ID:           techCart.ID,
			CheckID:      techCart.CheckID,
			TechCartID:   techCart.TechCartID,
			TechCartName: techCart.TechCartName,
			Quantity:     techCart.Quantity,
			Cost:         techCart.Cost,
			Price:        techCart.Price,
			Discount:     techCart.Discount,
			Comments:     techCart.Comments,
		}
		modificators := []*model.IngredientOutput{}
		ingredients := []*model.IngredientOutput{}
		for _, ingredient := range techCart.ExpenceIngredient {
			ing := model.IngredientOutput{}
			repo.gormDB.Table("ingredient_masters").Select("ingredient_masters.name").Where("ingredient_masters.id", ingredient.IngredientID).Scan(&ing)
			if ingredient.Type == utils.TypeModificator {
				modificators = append(modificators, &model.IngredientOutput{
					ID:           ingredient.ID,
					Name:         fmt.Sprintf("%s x%v", ing.Name, ingredient.Quantity),
					IngredientID: ingredient.IngredientID,
					Brutto:       ingredient.Quantity,
				})
			} else {
				ingredients = append(ingredients, &model.IngredientOutput{
					ID:           ingredient.ID,
					Name:         fmt.Sprintf("%s x%v", ing.Name, ingredient.Quantity),
					IngredientID: ingredient.IngredientID,
					Brutto:       ingredient.Quantity,
				})
			}
		}
		modif, err := json.Marshal(modificators)
		if err != nil {
			fmt.Println(err)
		}

		ingr, err := json.Marshal(ingredients)
		if err != nil {
			fmt.Println(err)
		}
		newTechCart.Modificators = string(modif)
		newTechCart.Ingredients = string(ingr)
		newTechCarts = append(newTechCarts, newTechCart)
	}
	return &model.CheckResponse{
		ID:               r.ID,
		MobileUserID:     r.MobileUserID,
		SkladID:          r.SkladID,
		ShopID:           r.ShopID,
		WorkerID:         r.WorkerID,
		Worker:           r.Worker,
		Opened_at:        r.Opened_at,
		Closed_at:        r.Closed_at,
		IdempotencyKey:   r.IdempotencyKey,
		Cash:             r.Cash,
		Card:             r.Card,
		Sum:              r.Sum,
		Cost:             r.Cost,
		Status:           r.Status,
		Payment:          r.Payment,
		Discount:         r.Discount,
		DiscountPercent:  r.DiscountPercent,
		Link:             r.Link,
		Tovar:            r.Tovar,
		TechCart:         newTechCarts,
		ModificatorCheck: r.ModificatorCheck,
		Comment:          r.Comment,
		Feedback:         r.Feedback,
	}
}

func (r *CheckDB) UpdateCheck(check *model.Check) (*model.CheckResponse, error) {
	err := r.gormDB.Transaction(func(tx *gorm.DB) error {
		oldCheck := &model.Check{}
		if err := tx.Model(&check).Where("id = ?", check.ID).Scan(oldCheck).Error; err != nil {
			return err
		}
		if oldCheck.Version > check.Version {
			return errors.New("old version")
		}
		if err := tx.Model(&check.Tovar).Where("check_id = ?", check.ID).Delete(&model.CheckTovar{}).Error; err != nil {
			return err
		}
		if err := tx.Model(&check.TechCart).Where("check_id = ?", check.ID).Delete(&model.CheckTechCart{}).Error; err != nil {
			return err
		}
		check.Opened_at = oldCheck.Opened_at
		for _, tovar := range check.Tovar {
			tovar.ExpenceTovar = nil
		}
		for _, techCart := range check.TechCart {
			techCart.ExpenceIngredient = nil
		}
		if err := tx.Save(check).Where("id = ?", check.ID).Error; err != nil {
			return err
		}
		return nil
	},
	)

	if err != nil {
		return nil, err
	}

	if check.Stolik != 0 {
		stolik := &model.Stolik{}
		res := r.gormDB.Model(&model.Stolik{}).Where("shop_id = ? and stolik_id = ?", check.ShopID, check.Stolik).Scan(&stolik)
		if res.Error != nil {
			return nil, res.Error
		}
		stolik.Empty = true
		res = r.gormDB.Debug().Model(&model.Stolik{}).Select("*").Where("shop_id = ? and stolik_id = ?", check.ShopID, check.Stolik).Updates(stolik)
		if res.Error != nil {
			return nil, res.Error
		}
	}

	checkResponse := &model.CheckResponse{
		ID:               check.ID,
		MobileUserID:     check.MobileUserID,
		SkladID:          check.SkladID,
		ShopID:           check.ShopID,
		WorkerID:         check.WorkerID,
		Worker:           check.Worker,
		Opened_at:        check.Opened_at,
		Closed_at:        check.Closed_at,
		IdempotencyKey:   check.IdempotencyKey,
		Cash:             check.Cash,
		Card:             check.Card,
		Sum:              check.Sum,
		Cost:             check.Cost,
		Status:           check.Status,
		Payment:          check.Payment,
		Discount:         check.Discount,
		DiscountPercent:  check.DiscountPercent,
		Tovar:            check.Tovar,
		ModificatorCheck: check.ModificatorCheck,
		Comment:          check.Comment,
		ServicePercent:   check.ServicePercent,
		ServiceSum:       check.ServiceSum,
	}
	return checkResponse, nil
}

func (r *CheckDB) DeleteCheck(id int) error {
	err := r.gormDB.Transaction(func(tx *gorm.DB) error {
		check := &model.Check{}
		err := tx.Select("*").Where("id = ?", id).First(&check).Error
		if err != nil {
			return err
		}
		tovars := []*model.CheckTovar{}

		err = tx.Model(tovars).Where("check_id = ?", id).Scan(&tovars).Error
		if err != nil {
			return err
		}
		for _, tovar := range tovars {
			err := tx.Delete(tovar).Error
			if err != nil {
				return err
			}
		}

		tech_carts := []*model.CheckTechCart{}

		err = tx.Model(tech_carts).Where("check_id = ?", id).Scan(&tech_carts).Error
		if err != nil {
			return err
		}
		for _, tech_cart := range tech_carts {
			err := tx.Delete(tech_cart).Error
			if err != nil {
				return err
			}
		}

		modificators := []*model.CheckModificator{}
		err = tx.Model(modificators).Where("check_id = ?", id).Scan(&modificators).Error
		if err != nil {
			return err
		}
		for _, modificator := range modificators {
			err := tx.Delete(modificator).Error
			if err != nil {
				return err
			}
		}

		err = tx.Delete(check).Error

		if err != nil {
			return err
		}

		return nil
	},
	)
	if err != nil {
		return err
	}
	return nil
}

func (r *CheckDB) CloseCheck(check *model.Check) (*model.Check, error) {
	err := r.gormDB.Transaction(func(tx *gorm.DB) error {
		if check.ID != 0 {
			oldCheck := &model.Check{}
			if err := tx.Model(&check).Where("id = ?", check.ID).Scan(oldCheck).Error; err != nil {
				return err
			}
			if oldCheck.Version > check.Version {
				return errors.New("old version")
			}
			check.Opened_at = oldCheck.Opened_at
		}
		if err := tx.Model(&check.Tovar).Where("check_id = ?", check.ID).Delete(&model.CheckTovar{}).Error; err != nil {
			return err
		}
		if err := tx.Model(&check.TechCart).Where("check_id = ?", check.ID).Delete(&model.CheckTechCart{}).Error; err != nil {
			return err
		}
		if err := tx.Model(&check.ModificatorCheck).Where("check_id = ?", check.ID).Delete(&model.CheckModificator{}).Error; err != nil {
			return err
		}
		for _, tovar := range check.Tovar {
			tovar.ExpenceTovar.Status = utils.StatusClosed
		}
		for _, techCart := range check.TechCart {
			for _, ingredients := range techCart.ExpenceIngredient {
				ingredients.Status = utils.StatusClosed
			}
		}
		if err := tx.Save(check).Where("id = ?", check.ID).Error; err != nil {
			return err
		}

		for _, tovar := range check.Tovar {
			res := tx.Table("sklad_tovars").Where("tovar_id = ? and sklad_id = ?", tovar.TovarID, check.SkladID).Update("quantity", gorm.Expr("quantity - ?", tovar.ExpenceTovar.Quantity))
			if res.Error != nil {
				return res.Error
			}
		}
		for i := 0; i < len(check.TechCart); i++ {
			for _, ingredient := range check.TechCart[i].ExpenceIngredient {
				res := tx.Table("sklad_ingredients").Where("ingredient_id = ? and sklad_id = ?", ingredient.IngredientID, check.SkladID).Update("quantity", gorm.Expr("quantity - ?", ingredient.Quantity))
				if res.Error != nil {
					return res.Error
				}
			}
		}
		return nil
	},
	)
	if err != nil {
		return nil, err
	}
	return check, nil
}

func (r *CheckDB) GetCheckByID(id int) (*model.CheckResponse, error) {
	checkGorm := &model.CheckResponse{}
	tovarsGorm := []*model.CheckTovar{}
	techCartsGorm := []*model.CheckTechCart{}
	techCartsResponses := []*model.CheckTechCartResponse{}
	err := r.gormDB.Transaction(func(tx *gorm.DB) error {
		res := tx.Table("checks").Select("*").Where("id = ?", id).Scan(&checkGorm)
		if res.Error != nil {
			return res.Error
		}
		res = tx.Table("check_tovars").Select("*").Where("check_id = ?", id).Scan(&tovarsGorm)
		if res.Error != nil {
			return res.Error
		}
		res = tx.Table("check_tech_carts").Select("*").Where("check_id = ?", id).Scan(&techCartsGorm)
		if res.Error != nil {
			return res.Error
		}
		for _, item := range techCartsGorm {

			techCartsResponse := &model.CheckTechCartResponse{
				ID:           item.ID,
				CheckID:      item.CheckID,
				TechCartID:   item.TechCartID,
				TechCartName: item.TechCartName,
				Quantity:     item.Quantity,
				Cost:         item.Cost,
				Price:        item.Price,
				Discount:     item.Discount,
				Comments:     item.Comments,
			}

			ingredients := []*model.IngredientNumOutput{}
			res = tx.Table("expence_ingredients").Select("ingredients.ingredient_id, expence_ingredients.quantity as brutto, ingredients.name, ingredients.measure, expence_ingredients.cost").Joins("INNER JOIN ingredients on ingredients.ingredient_id = expence_ingredients.ingredient_id").Where("check_tech_cart_id = ? and expence_ingredients.type = ? and ingredients.shop_id = ?", item.ID, utils.TypeIngredient, checkGorm.ShopID).Scan(&ingredients)
			if res.Error != nil {
				return res.Error
			}

			modificators := []*model.IngredientOutput{}
			res = tx.Table("expence_ingredients").Select("ingredients.ingredient_id, expence_ingredients.quantity as brutto, ingredients.name, ingredients.measure, expence_ingredients.cost, expence_ingredients.price").Joins("INNER JOIN ingredients on ingredients.ingredient_id = expence_ingredients.ingredient_id").Where("check_tech_cart_id = ? and expence_ingredients.type = ? and ingredients.shop_id = ?", item.ID, utils.TypeModificator, checkGorm.ShopID).Scan(&modificators)
			if res.Error != nil {
				return res.Error
			}
			ingredient, err := json.Marshal(ingredients)
			if err != nil {
				return err
			}
			modificator, err := json.Marshal(modificators)
			if err != nil {
				return err
			}
			techCartsResponse.Ingredients = string(ingredient)
			techCartsResponse.Modificators = string(modificator)
			techCartsResponses = append(techCartsResponses, techCartsResponse)
		}
		// for i := 0; i < len(checkGorm.TechCart); i++ {
		// 	ingredients := []*model.IngredientNumOutput{}
		// 	res = tx.Table("ingredient_tech_carts").Select("ingredients.id, ingredient_tech_carts.brutto, ingredients.name, ingredients.measure, sklad_ingredients.cost").Joins("inner join ingredients on ingredient_tech_carts.ingredient_id = ingredients.id inner join sklad_ingredients on sklad_ingredients.ingredient_id = ingredients.id").Where("ingredient_tech_carts.tech_cart_id = ? and sklad_ingredients.sklad_id = ?", checkGorm.TechCart[i].TechCartID, checkGorm.SkladID).Scan(&ingredients)
		// 	if res.Error != nil {
		// 		return res.Error
		// 	}
		// 	ingredient, err := json.Marshal(ingredients)
		// 	if err != nil {
		// 		return err
		// 	}
		// 	checkGorm.TechCart[i].Ingredients = string(ingredient)
		// }
		checkGorm.Tovar = tovarsGorm
		checkGorm.TechCart = techCartsResponses
		return nil
	})
	if err != nil {
		return nil, err
	}
	return checkGorm, nil
}

//SELECT ingredients.id, ingredient_tech_carts.brutto, ingredients.name, ingredients.measure, ingredients.cost
//FROM public.ingredient_tech_carts inner join ingredients on ingredient_tech_carts.ingredient_id = ingredients.id where ingredient_tech_carts.tech_cart_id = 15;

func (r *CheckDB) GetAllCheck(filter *model.Filter) ([]*model.Check, int64, error) {
	checksGorm := []*model.Check{}
	checks := []*model.CheckNeOutput{}

	res := r.gormDB.Table("checks").Select("checks.id, json_build_object('user_id', checks.user_id, 'sklad_id', checks.sklad_id, 'shop_id', checks.shop_id, 'worker_id', checks.worker_id, 'worker', workers.name, 'opened_at', checks.opened_at, 'closed_at', checks.closed_at, 'idempotency_key', checks.idempotency_key, 'version', checks.version, 'cash', checks.cash, 'card', checks.card, 'sum', checks.sum, 'cost', checks.cost, 'status', checks.status, 'payment', checks.payment, 'discount', checks.discount, 'discount_percent', checks.discount_percent, 'comment', checks.comment, 'tovarCheck', (SELECT json_agg(check_tovars) from check_tovars where check_tovars.check_id = checks.id), 'techCartCheck',  (SELECT json_agg(check_tech_carts) from check_tech_carts where check_tech_carts.check_id = checks.id), 'modificatorCheck', (SELECT json_agg(check_modificators) from check_modificators where check_modificators.check_id = checks.id)) AS checks_info").Joins("left join check_tovars on check_tovars.check_id = checks.id left join check_tech_carts on check_tech_carts.check_id = checks.id left join check_modificators on check_modificators.check_id = checks.id left join workers on workers.id = checks.worker_id").Group("checks.id, workers.name")

	if res.Error != nil {
		return nil, 0, res.Error
	}

	newRes, countRes, err := filter.FilterResults(res, checksGorm, utils.ChecksPageSize, "opened_at", fmt.Sprintf("check_tovars.tovar_name ilike '%%%s%%' OR check_tech_carts.tech_cart_name ilike '%%%s%%' OR check_modificators.name ilike '%%%s%%'", filter.Search, filter.Search, filter.Search), "")
	if err != nil {
		return nil, 0, err
	}
	rows, err := newRes.Rows()
	if err != nil {
		return nil, 0, err
	}
	var result model.CheckInfo

	for rows.Next() {
		type IDStruct struct {
			ID int64 `json:"id"`
		}
		var id IDStruct
		err = rows.Scan(&id.ID, &result)

		if err != nil {
			return nil, 0, err
		}
		checks = append(checks, &model.CheckNeOutput{
			ID:   int(id.ID),
			Info: result,
		})
	}

	for _, check := range checks {
		if check.Info.Tovar == nil {
			check.Info.Tovar = []*model.CheckTovar{}
		}
		if check.Info.TechCart == nil {
			check.Info.TechCart = []*model.CheckTechCart{}
		}
		if check.Info.ModificatorCheck == nil {
			check.Info.ModificatorCheck = []*model.CheckModificator{}
		}
		checksGorm = append(checksGorm, &model.Check{
			ID:               check.ID,
			MobileUserID:     check.Info.MobileUserID,
			SkladID:          check.Info.SkladID,
			ShopID:           check.Info.ShopID,
			WorkerID:         check.Info.WorkerID,
			Worker:           check.Info.Worker,
			Opened_at:        check.Info.Opened_at,
			Closed_at:        check.Info.Closed_at,
			IdempotencyKey:   check.Info.IdempotencyKey,
			Version:          check.Info.Version,
			Cash:             check.Info.Cash,
			Card:             check.Info.Card,
			Sum:              check.Info.Sum,
			Cost:             check.Info.Cost,
			Status:           check.Info.Status,
			Payment:          check.Info.Payment,
			Discount:         check.Info.Discount,
			DiscountPercent:  check.Info.DiscountPercent,
			Tovar:            check.Info.Tovar,
			TechCart:         check.Info.TechCart,
			ModificatorCheck: check.Info.ModificatorCheck,
			Comment:          check.Info.Comment,
		})
	}

	return checksGorm, countRes, err

	//Without search
	/*
			checksGorm := []*model.Check{}

			res := r.gormDB.Table("checks").Select("*")
			if res.Error != nil {
				return nil, 0, res.Error
			}

			newRes, countRes, err := filter.FilterResults(res, checksGorm, utils.ChecksPageSize, "opened_at", "")
			if err != nil {
				return nil, 0, err
			}

			if newRes.Scan(&checksGorm).Error != nil {
			return nil, 0, newRes.Error
		}
		for i := 0; i < len(checksGorm); i++ {
			tovarsGorm := []*model.CheckTovar{}
			techCartsGorm := []*model.CheckTechCart{}
			res := r.gormDB.Table("check_tovars").Select("*").Where("check_id = ?", checksGorm[i].ID).Scan(&tovarsGorm)
			if res.Error != nil {
				return nil, 0, res.Error
			}
			res = r.gormDB.Table("check_tech_carts").Select("*").Where("check_id = ?", checksGorm[i].ID).Scan(&techCartsGorm)
			if res.Error != nil {
				return nil, 0, res.Error
			}
			checksGorm[i].TechCart = techCartsGorm
			checksGorm[i].Tovar = tovarsGorm
		}

		if err != nil {
			return nil, 0, err
		}

			return checksGorm, countRes, err

	*/
}

func (r *CheckDB) GetTisCheck(id int) (*model.ReqTisResponse, error) {
	readCheck := &model.SendToTis{}
	err := r.gormDB.Model(&model.SendToTis{}).Where("check_id = ?", id).Scan(readCheck).Error
	if err != nil {
		return nil, err
	}
	res := &model.ReqTisResponse{}
	err = json.Unmarshal([]byte(readCheck.Response), &res)
	if err != nil {
		return nil, nil
	}

	return res, err
}

func (r *CheckDB) GetAllWorkerCheck(filter *model.Filter) ([]*model.CheckResponse, int64, error) {
	checkTovars := []*model.CheckTovar{}
	checkTechCarts := []*model.CheckTechCart{}
	checkModificators := []*model.CheckModificator{}
	checkOutput := []*model.CheckResponse{}
	shift := &model.Shift{}
	err := r.gormDB.Model(&model.Shift{}).Where("shop_id = ?", filter.BindShop).Last(&model.Shift{}).Scan(shift).Error
	if err != nil {
		return nil, 0, err
	}
	from := shift.CreatedAt
	res := r.gormDB.Model(&model.Check{}).Where("status != ? and shop_id IN (?) and opened_at::date >= ?::date", utils.StatusInactive, filter.AccessibleShops, from)
	if res.Error != nil {
		return nil, 0, res.Error
	}
	newRes, count, err := filter.FilterResults(res, &model.Check{}, utils.ChecksPageSize, "", "", "id desc")

	if err != nil {
		return nil, 0, err
	}

	if newRes.Scan(&checkOutput).Error != nil {
		return nil, 0, newRes.Error
	}

	if err := r.gormDB.Table("checks").Select("*").Joins("inner join check_tovars on check_tovars.check_id = checks.id").Where("checks.status != ? and checks.shop_id IN (?) and opened_at::date >= ?::date", utils.StatusInactive, filter.AccessibleShops, from).Order("check_tovars.check_id desc").Scan(&checkTovars).Error; err != nil {
		return nil, 0, err
	}

	if err := r.gormDB.Table("checks").Select("*").Joins("inner join check_tech_carts on check_tech_carts.check_id = checks.id").Where("checks.status != ? and checks.shop_id IN (?) and opened_at::date >= ?::date", utils.StatusInactive, filter.AccessibleShops, from).Order("check_tech_carts.check_id desc").Scan(&checkTechCarts).Error; err != nil {
		return nil, 0, err
	}

	if err := r.gormDB.Table("checks").Select("*").Joins("inner join check_modificators on check_modificators.check_id = checks.id").Where("checks.status != ? and checks.shop_id IN (?) and opened_at::date >= ?::date", utils.StatusInactive, filter.AccessibleShops, from).Order("check_modificators.check_id desc").Scan(&checkModificators).Error; err != nil {
		return nil, 0, err
	}

	for _, check := range checkOutput {
		techCartSaved := false
		for _, checkTechCart := range checkTechCarts {
			if checkTechCart.CheckID == check.ID {
				ingredients := []*model.IngredientNumOutput{}
				res = r.gormDB.Table("expence_ingredients").Select("ingredients.ingredient_id, expence_ingredients.quantity as brutto, ingredients.name, ingredients.measure, expence_ingredients.cost").Joins("INNER JOIN ingredients on ingredients.ingredient_id = expence_ingredients.ingredient_id").Where("check_tech_cart_id = ? and expence_ingredients.type = ? and ingredients.shop_id = ?", checkTechCart.ID, utils.TypeIngredient, filter.BindShop).Scan(&ingredients)
				if res.Error != nil {
					return nil, 0, res.Error
				}
				ingredient, err := json.Marshal(ingredients)
				if err != nil {
					return nil, 0, err
				}
				modificators := []*model.IngredientOutput{}
				res = r.gormDB.Table("expence_ingredients").Select("ingredients.ingredient_id, expence_ingredients.quantity as brutto, ingredients.name, ingredients.measure, expence_ingredients.cost, expence_ingredients.price").Joins("INNER JOIN ingredients on ingredients.ingredient_id = expence_ingredients.ingredient_id").Where("check_tech_cart_id = ? and expence_ingredients.type = ? and ingredients.shop_id = ?", checkTechCart.ID, utils.TypeModificator, filter.BindShop).Scan(&modificators)
				if res.Error != nil {
					return nil, 0, res.Error
				}
				modificator, err := json.Marshal(modificators)
				if err != nil {
					return nil, 0, err
				}
				check.TechCart = append(check.TechCart, &model.CheckTechCartResponse{
					ID:           checkTechCart.ID,
					CheckID:      checkTechCart.CheckID,
					TechCartID:   checkTechCart.TechCartID,
					TechCartName: checkTechCart.TechCartName,
					Quantity:     checkTechCart.Quantity,
					Cost:         checkTechCart.Cost,
					Price:        checkTechCart.Price,
					Discount:     checkTechCart.Discount,
					Comments:     checkTechCart.Comments,
					Modificators: string(modificator),
					Ingredients:  string(ingredient),
				})
				techCartSaved = true
			} else {
				if techCartSaved {
					break
				}
			}
		}
		tovarSaved := false
		for _, checkTovar := range checkTovars {
			if checkTovar.CheckID == check.ID {
				check.Tovar = append(check.Tovar, checkTovar)
				tovarSaved = true
			} else {
				if tovarSaved {
					break
				}
			}
		}
		modificatorSaved := false
		for _, checkModificator := range checkModificators {
			if checkModificator.CheckID == check.ID {
				check.ModificatorCheck = append(check.ModificatorCheck, checkModificator)
				modificatorSaved = true
			} else {
				if modificatorSaved {
					break
				}
			}
		}
	}

	for _, check := range checkOutput {
		if check.Tovar == nil {
			check.Tovar = []*model.CheckTovar{}
		}
		if check.TechCart == nil {
			check.TechCart = []*model.CheckTechCartResponse{}
		}
		if check.ModificatorCheck == nil {
			check.ModificatorCheck = []*model.CheckModificator{}
		}
	}

	if err != nil {
		return nil, 0, err
	}

	return checkOutput, count, nil
}

func (r *CheckDB) GetAllCheckView(page int) ([]*model.CheckView, int64, error) {
	checksGorm := []*model.CheckView{}
	var count int64
	var res *gorm.DB
	err := r.gormDB.Transaction(func(tx *gorm.DB) error {
		if page != 0 {
			offset := (page - 1) * utils.ChecksPageSize
			res := tx.Table("check_views").Select("*").Offset(offset).Limit(utils.ChecksPageSize).Order("id desc").Scan(&checksGorm)
			if res.Error != nil {
				return res.Error
			}
		} else {
			res := tx.Table("check_views").Select("*").Order("id desc").Scan(&checksGorm)
			if res.Error != nil {
				return res.Error
			}
		}

		res = tx.Table("check_views").Count(&count)
		if res.Error != nil {
			return res.Error
		}
		return nil
	})
	if err != nil {
		return nil, 0, err
	}
	return checksGorm, count, err
}

func (r *CheckDB) RemoveFromSklad(check *model.Check) (*model.Check, error) {
	err := r.gormDB.Transaction(func(tx *gorm.DB) error {
		for _, tovar := range check.Tovar {
			res := tx.Table("sklad_tovar").Where("id = ?", tovar.TovarID).Update("quantity", gorm.Expr("quantity - ?", tovar.Quantity))
			if res.Error != nil {
				return res.Error
			}
		}
		for _, techCart := range check.TechCart {
			ingredients := []*model.IngredientOutput{}
			res := tx.Table("tech_cart_ingredient").Select("tech_cart_ingredient.ingredient_id, tech_cart_ingredient.brutto, ingredient.name, ingredient.measure, ingredient.cost, ingredient.price").Joins("inner join ingredient on ingredient.id = tech_cart_ingredient.ingredient_id").Where("tech_cart_id = ?", techCart.TechCartID).Scan(&ingredients)
			if res.Error != nil {
				return res.Error
			}
			for _, ingredient := range ingredients {
				res := tx.Table("sklad_ingredient").Where("ingredient_id = ?", ingredient.ID).Update("quantity", gorm.Expr("quantity - ?", ingredient.Brutto*techCart.Quantity))
				if res.Error != nil {
					return res.Error
				}
			}
			//techCart.Ingredients = ingredients
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return check, nil
}

func (r *CheckDB) GetTisToken(shopID int) (string, error) {
	shop := &model.Shop{}
	err := r.gormDB.Transaction(func(tx *gorm.DB) error {
		res := tx.Table("shops").Where("id = ?", shopID).First(&shop)
		if res.Error != nil {
			return res.Error
		}
		return nil
	})
	if err != nil {
		return "", err
	}
	return shop.TisToken, nil
}

func (r *CheckDB) CalculateCheck(check *model.ReqCheck) (*model.Check, error) {
	var cost float32 = 0
	var sum float32 = 0
	var discount float32 = 0

	newCheck := &model.Check{
		ID:              check.ID,
		MobileUserID:    check.MobileUserID,
		WorkerID:        check.WorkerID,
		IdempotencyKey:  check.IdempotencyKey,
		Version:         check.Version,
		Opened_at:       check.Opened_at,
		Closed_at:       check.Closed_at,
		Cash:            check.Cash,
		Card:            check.Card,
		Status:          check.Status,
		Payment:         check.Payment,
		Comment:         check.Comment,
		DiscountPercent: check.DiscountPercent,
		ShopID:          check.ShopID,
		SkladID:         check.SkladID,
		Stolik:          check.Stolik,
	}
	err := r.gormDB.Transaction(func(tx *gorm.DB) error {
		newTovars := []*model.CheckTovar{}
		newTechCarts := []*model.CheckTechCart{}

		for _, tovar := range check.Tovars {
			tempTovar := &model.Tovar{}
			err := tx.Table("tovars").Select("id, tovar_id, shop_id, name , price, category, discount").Where("tovar_id = ? and shop_id = ?", tovar.TovarID, newCheck.ShopID).Scan(tempTovar).Error
			if err != nil {
				return err
			}
			skladTovar := &model.SkladTovar{}
			err = tx.Table("sklad_tovars").Select("cost").Where("tovar_id = ? and sklad_id = ?", tovar.TovarID, newCheck.SkladID).Scan(skladTovar).Error
			if err != nil {
				return err
			}
			newTovar := &model.CheckTovar{
				TovarID:   tempTovar.TovarID,
				TovarName: tempTovar.Name,
				Cost:      skladTovar.Cost,
				Price:     tempTovar.Price,
			}
			if tempTovar.Discount {
				discount = discount + tempTovar.Price*check.DiscountPercent*tovar.Quantity
				newTovar.Discount = tempTovar.Price * check.DiscountPercent * tovar.Quantity
			}
			newTovar.Quantity = tovar.Quantity
			newTovar.CheckID = check.ID
			newTovar.Comments = tovar.Comments
			cost = cost + newTovar.Cost*tovar.Quantity
			sum = sum + newTovar.Price*tovar.Quantity
			expence := &model.ExpenceTovar{
				TovarID:  tempTovar.TovarID,
				Quantity: tovar.Quantity,
				SkladID:  newCheck.SkladID,
				Cost:     newTovar.Cost * tovar.Quantity,
				Time:     check.Closed_at,
			}
			newTovar.ExpenceTovar = expence
			newTovars = append(newTovars, newTovar)
		}

		for _, techCart := range check.TechCarts {
			log.Print("tech cart is ", techCart)
			tempTechCart := &model.TechCart{}
			err := tx.Table("tech_carts").Select("id, tech_cart_id, shop_id, name, price, category, discount").Where("tech_cart_id = ? and shop_id = ?", techCart.TechCartID, newCheck.ShopID).Scan(tempTechCart).Error
			if err != nil {
				return err
			}
			newTechCart := &model.CheckTechCart{
				TechCartID:   tempTechCart.TechCartID,
				TechCartName: tempTechCart.Name,
				Price:        tempTechCart.Price,
			}

			if tempTechCart.Discount {
				newTechCart.Discount = tempTechCart.Price * check.DiscountPercent * techCart.Quantity
				discount = discount + newTechCart.Discount
			}

			newTechCart.Quantity = techCart.Quantity
			newTechCart.CheckID = check.ID
			newTechCart.Comments = techCart.Comments

			var techCartCost float32
			ingredientsToRemove := []*model.ExpenceIngredient{}
			ingredients := []*model.IngredientOutput{}

			res := tx.Table("ingredient_tech_carts").Select("ingredients.id, ingredients.ingredient_id, ingredients.shop_id, ingredient_tech_carts.brutto, ingredients.name, ingredients.measure, sklad_ingredients.cost").Joins("inner join ingredients on ingredients.ingredient_id = ingredient_tech_carts.ingredient_id inner join sklad_ingredients on sklad_ingredients.ingredient_id = ingredients.ingredient_id").Where("tech_cart_id = ? and sklad_ingredients.sklad_id = ? and ingredient_tech_carts.shop_id = ? and ingredients.shop_id = ?", techCart.TechCartID, newCheck.SkladID, newCheck.ShopID, newCheck.ShopID).Scan(&ingredients)
			if res.Error != nil {
				return res.Error
			}
			replaceIngredientsGlobal := make(map[int]bool)
			for _, modificator := range techCart.Modificators {
				replaceIngredients := make(map[int]bool)
				log.Print("Modificator is ", modificator)
				nabor := &model.Nabor{}
				res := tx.Table("nabors").Select("*").Where("nabor_id = ? and shop_id = ?", modificator.NaborID, newCheck.ShopID).Scan(&nabor)
				if res.Error != nil {
					return res.Error
				}
				for _, replace := range nabor.Replaces {
					replaceIngredients[replace] = true
					replaceIngredientsGlobal[replace] = true
				}

				modif := &model.IngredientOutput{}
				res = tx.Table("ingredient_nabors").Select("ingredients.id, ingredients.ingredient_id, ingredients.shop_id, ingredient_nabors.price, ingredient_nabors.brutto, sklad_ingredients.cost, ingredients.name, ingredients.measure").Joins("inner join ingredients on ingredient_nabors.ingredient_id = ingredients.ingredient_id inner join sklad_ingredients on sklad_ingredients.ingredient_id = ingredients.ingredient_id").Where("ingredient_nabors.nabor_id = ? and ingredient_nabors.ingredient_id = ? and sklad_ingredients.sklad_id = ? and ingredients.shop_id = ? and ingredient_nabors.shop_id = ?", modificator.NaborID, modificator.ID, newCheck.SkladID, newCheck.ShopID, newCheck.ShopID).Scan(&modif)
				if res.Error != nil {
					return res.Error
				}
				for _, ingredient := range ingredients {
					if replaceIngredients[ingredient.IngredientID] {
						modif.Brutto = ingredient.Brutto
					}
				}
				techCartCost = techCartCost + modif.Brutto*modif.Cost*techCart.Quantity*modificator.Quantity
				newTechCart.Price = newTechCart.Price + modif.Price*modificator.Quantity
				ingredientToRemove := &model.ExpenceIngredient{
					SkladID:      newCheck.SkladID,
					IngredientID: modif.IngredientID,
					Quantity:     modif.Brutto * techCart.Quantity * modificator.Quantity,
					Cost:         modif.Cost * modif.Brutto * techCart.Quantity * modificator.Quantity,
					Time:         check.Closed_at,
					Type:         utils.TypeModificator,
					Price:        modif.Price * techCart.Quantity,
				}
				ingredientsToRemove = append(ingredientsToRemove, ingredientToRemove)
			}

			for _, ingredient := range ingredients {
				if !replaceIngredientsGlobal[ingredient.IngredientID] {
					techCartCost = techCartCost + ingredient.Cost*ingredient.Brutto*techCart.Quantity
					ingredientsToRemove = append(ingredientsToRemove, &model.ExpenceIngredient{
						SkladID:      newCheck.SkladID,
						IngredientID: ingredient.IngredientID,
						Quantity:     ingredient.Brutto * techCart.Quantity,
						Cost:         ingredient.Cost * ingredient.Brutto * techCart.Quantity,
						Time:         check.Closed_at,
						Type:         utils.TypeIngredient,
					})
				}
			}
			cost = cost + techCartCost
			sum = sum + newTechCart.Price*techCart.Quantity
			newTechCart.Cost = techCartCost
			newTechCart.ExpenceIngredient = ingredientsToRemove
			newTechCarts = append(newTechCarts, newTechCart)
		}
		newCheck.Tovar = newTovars
		newCheck.TechCart = newTechCarts
		newCheck.Cost = cost
		newCheck.Sum = float32(math.Floor(float64(sum)))
		newCheck.Discount = float32(math.Round(float64(discount)))
		return nil
	})
	if err != nil {
		return nil, err
	}
	if check.DiscountSum != 0 {
		newCheck.Discount += check.DiscountSum
	}
	if check.Service {
		var servicePercent float32
		res := r.gormDB.Table("shops").Select("service_percent").Where("id = ?", newCheck.ShopID).Scan(&servicePercent)
		if res.Error != nil {
			return nil, res.Error
		}
		newCheck.ServicePercent = servicePercent
		newCheck.ServiceSum = newCheck.Sum * servicePercent / 100
		newCheck.Sum = newCheck.Sum + newCheck.ServiceSum
	}
	if check.Stolik != 0 {
		stolik := &model.Stolik{}
		res := r.gormDB.Model(&model.Stolik{}).Where("shop_id = ? and stolik_id = ?", check.ShopID, check.Stolik).Scan(&stolik)
		if res.Error != nil {
			return nil, res.Error
		}
		stolik.Empty = true
		res = r.gormDB.Debug().Model(&model.Stolik{}).Select("*").Where("shop_id = ? and stolik_id = ?", check.ShopID, check.Stolik).Updates(stolik)
		if res.Error != nil {
			return nil, res.Error
		}
	}
	return newCheck, nil
}

func (r *CheckDB) AddTag(tag *model.Tag) error {
	res := r.gormDB.Create(tag)

	if res.Error != nil {
		return res.Error
	}

	return nil
}

func (r *CheckDB) GetTag(id int) (*model.Tag, error) {
	tag := &model.Tag{}
	res := r.gormDB.First(tag, id)

	if res.Error != nil {
		return nil, res.Error
	}

	return tag, nil
}

func (r *CheckDB) GetAllTag(shopID int) ([]*model.Tag, error) {
	tags := []*model.Tag{}
	res := r.gormDB.Where("shop_id = ?", shopID).Order("id asc").Find(&tags)

	if res.Error != nil {
		return nil, res.Error
	}

	return tags, nil
}

func (r *CheckDB) UpdateTag(tag *model.Tag) error {
	res := r.gormDB.Save(tag)

	if res.Error != nil {
		return res.Error
	}

	return nil
}

func (r *CheckDB) DeleteTag(id int) error {
	res := r.gormDB.Delete(&model.Tag{}, id)

	if res.Error != nil {
		return res.Error
	}

	return nil
}

func (r *CheckDB) GetCheckByIdempotency(idempotencyKey string) (*model.Check, error) {
	check := &model.Check{}
	res := r.gormDB.Where("idempotency_key = ?", idempotencyKey).First(check)

	if res.Error != nil {
		return nil, res.Error
	}

	return check, nil
}

func (r *CheckDB) SaveCheck(check *model.TisResponse) error {
	err := r.gormDB.Model(&model.TisResponse{}).Create(check).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *CheckDB) SaveError(fullError string, request string) error {
	errorCheck := &model.ErrorCheck{
		Error:   fullError,
		Request: request,
		Time:    time.Now(),
	}
	err := r.gormDB.Model(&model.ErrorCheck{}).Create(errorCheck).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *CheckDB) DeactivateCheck(id int) ([]*model.InventarizationItem, error) {
	invItems := []*model.InventarizationItem{}

	err := r.gormDB.Transaction(func(tx *gorm.DB) error {
		check := &model.Check{}
		res := tx.Preload("Tovar").Preload("TechCart").First(&check, id)
		if res.Error != nil {
			return res.Error
		}
		itemsMap := make(map[int]bool)
		for _, tovar := range check.Tovar {
			expences := []*model.ExpenceTovar{}
			res = tx.Where("check_tovar_id = ?", tovar.ID).Find(&expences)
			if res.Error != nil {
				return res.Error
			}
			for _, expence := range expences {
				invItem := &model.InventarizationItem{}
				err := tx.Model(invItem).Where("sklad_id = ? and item_id = ? and type = ? and time > ? and before_time < ?", expence.SkladID, expence.TovarID, utils.TypeTovar, expence.Time, expence.Time).Scan(invItem).Error
				if err != nil {
					if err.Error() != "record not found" {
						return err
					}
				}
				if invItem.ID != 0 && invItem.ItemID != 0 && !itemsMap[invItem.ID] {
					itemsMap[invItem.ID] = true
					invItems = append(invItems, invItem)
				}
				res := tx.Model(invItem).Where("sklad_id = ? and item_id = ? and type = ? and before_time > ?", expence.SkladID, expence.TovarID, utils.TypeTovar, expence.Time).Scan(invItem)
				if err != nil {
					if err.Error() != "record not found" {
						return err
					}
				}
				if res.RowsAffected == 0 {
					res := tx.Table("sklad_tovars").Where("tovar_id = ? and sklad_id = ?", tovar.TovarID, check.SkladID).Update("quantity", gorm.Expr("quantity + ?", expence.Quantity))
					if res.Error != nil {
						return res.Error
					}
				}
				err = tx.Model(&model.ExpenceTovar{}).Where("id = ?", expence.ID).Delete(&model.ExpenceTovar{}).Error
				if err != nil {
					return err
				}

			}
		}
		for _, techCart := range check.TechCart {
			expences := []*model.ExpenceIngredient{}
			res = tx.Where("check_tech_cart_id = ?", techCart.ID).Find(&expences)
			if res.Error != nil {
				return res.Error
			}
			for _, expence := range expences {
				invItem := &model.InventarizationItem{}
				err := tx.Model(invItem).Where("sklad_id = ? and item_id = ? and type = ? and time > ? and before_time < ?", expence.SkladID, expence.IngredientID, utils.TypeIngredient, expence.Time, expence.Time).Scan(invItem).Error
				if err != nil {
					if err.Error() != "record not found" {
						return err
					}
				}
				if invItem.ID != 0 && invItem.ItemID != 0 && !itemsMap[invItem.ID] {
					itemsMap[invItem.ID] = true
					invItems = append(invItems, invItem)
				}
				res := tx.Model(invItem).Where("sklad_id = ? and item_id = ? and type = ? and before_time > ?", expence.SkladID, expence.IngredientID, utils.TypeIngredient, expence.Time).Scan(invItem)
				if err != nil {
					if err.Error() != "record not found" {
						return err
					}
				}
				if res.RowsAffected == 0 {
					res := tx.Table("sklad_ingredients").Where("ingredient_id = ? and sklad_id = ?", expence.IngredientID, check.SkladID).Update("quantity", gorm.Expr("quantity + ?", expence.Quantity))
					if res.Error != nil {
						return res.Error
					}
				}
				err = tx.Model(&model.ExpenceIngredient{}).Where("id = ?", expence.ID).Delete(&model.ExpenceIngredient{}).Error
				if err != nil {
					return err
				}
			}
		}
		res = tx.Model(&model.Check{}).Where("id = ?", id).Update("status", utils.StatusInactive)
		if res.Error != nil {
			return res.Error
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return invItems, nil
}

func (r *CheckDB) UpdateCheckLink(id int, link string) error {
	res := r.gormDB.Model(&model.Check{}).Where("id = ?", id).Update("link", link)
	if res.Error != nil {
		return res.Error
	}
	return nil
}

func (r *CheckDB) IdempotencyCheck(key *model.IdempotencyCheckArray, shopID int) (string, error) {
	keys := key.Keys
	keysMap := make(map[string]bool, len(keys))
	for _, key := range keys {
		keysMap[key.IdempotencyKey] = true
	}
	shift := &model.Shift{}
	if key.TimeIsLast {
		err := r.gormDB.Model(&model.Shift{}).Where("shop_id = ?", shopID).Last(&model.Shift{}).Debug().Scan(shift).Error
		if err != nil {
			return "", err
		}
	} else {
		err := r.gormDB.Model(&model.Shift{}).Where("shop_id = ? and created_at::date = ?::date", shopID, key.Date).Last(&model.Shift{}).Debug().Scan(shift).Error
		if err != nil {
			return "", err
		}
	}
	sklad := &model.Sklad{}
	err := r.gormDB.Model(&model.Sklad{}).Where("shop_id = ?", shopID).First(&model.Sklad{}).Debug().Scan(sklad).Error
	if err != nil {
		return "", err
	}
	checks := []*model.Check{}
	if key.TimeIsLast {
		err := r.gormDB.Model(&model.Check{}).Where("closed_at > ? and sklad_id = ?", shift.CreatedAt, sklad.ID).Debug().Find(&checks).Error
		if err != nil {
			return "", err
		}
	} else {
		err := r.gormDB.Model(&model.Check{}).Where("closed_at > ? and closed_at < ? and sklad_id = ?", shift.CreatedAt, shift.ClosedAt, sklad.ID).Debug().Find(&checks).Error
		if err != nil {
			return "", err
		}
	}
	inKeyButNotInCheck := ""
	countInCheck := 0
	inCheckButNotInKey := ""
	countInKey := 0
	res := ""
	doesNotExist := false
	for _, check := range checks {
		idKey := check.IdempotencyKey + "_1"
		if _, exists := keysMap[idKey]; !exists {
			inCheckButNotInKey = inCheckButNotInKey + idKey + " "
			countInCheck++
			doesNotExist = true
		} else {
			keysMap[idKey] = false
		}
	}
	for key, value := range keysMap {
		if value {
			inKeyButNotInCheck = inKeyButNotInCheck + key + " "
			countInKey++
			doesNotExist = true
		}
	}
	res = fmt.Sprintf("There are %d id in keys but not in checks: %s\nThere are %d id in checks but not in keys: %s", countInKey, inKeyButNotInCheck, countInCheck, inCheckButNotInKey)
	if doesNotExist {
		return res, errors.New("idempotency key does not exist")
	}

	return "", nil
}

func (r *CheckDB) SaveFailedCheck(checks []*model.FailedCheck) error {
	res := r.gormDB.Debug().Model([]*model.FailedCheck{}).Create(checks)
	if res.Error != nil {
		return res.Error
	}
	return nil
}

func (r *CheckDB) GetStoliki(shopID int) ([]*model.Stolik, error) {
	stoliki := []*model.Stolik{}
	res := r.gormDB.Model(&model.Stolik{}).Where("shop_id = ?", shopID).Scan(&stoliki)
	if res.Error != nil {
		return nil, res.Error
	}
	return stoliki, nil
}

func (r *CheckDB) GetFilledStoliki(shopID int) ([]*model.Stolik, error) {
	stoliki := []*model.Stolik{}
	res := r.gormDB.Model(&model.Stolik{}).Where("shop_id = ? and empty = ?", shopID, false).Scan(&stoliki)
	if res.Error != nil {
		return nil, res.Error
	}
	return stoliki, nil
}
