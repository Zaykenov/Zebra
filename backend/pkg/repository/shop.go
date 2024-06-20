package repository

import (
	"database/sql"
	"time"
	"zebra/model"
	"zebra/utils"

	"gorm.io/gorm"
)

type ShopDB struct {
	db     *sql.DB
	gormDB *gorm.DB
}

func NewShopDB(db *sql.DB, gormDB *gorm.DB) *ShopDB {
	return &ShopDB{db: db, gormDB: gormDB}
}

func (r *ShopDB) CreateShop(shop *model.Shop, products *model.ProductsShop) (*model.Shop, error) {
	err := r.gormDB.Transaction(func(tx *gorm.DB) error {
		err := tx.Create(shop).Error
		if err != nil {
			return err
		}
		tovarsMap := make(map[int]bool, 100)
		for _, v := range products.Tovars {
			tovarsMap[v] = true
		}
		techCartMap := make(map[int]bool, 100)
		for _, v := range products.TechCarts {
			techCartMap[v] = true
		}

		tovarsMaster := []*model.TovarMaster{}
		err = tx.Model(&model.TovarMaster{}).Where("id in (?)", products.Tovars).Scan(&tovarsMaster).Error
		if err != nil {
			return err
		}
		tovars := make([]*model.Tovar, 0, 100)

		for _, tovarMaster := range tovarsMaster {
			tovar := &model.Tovar{
				TovarID:   tovarMaster.ID,
				ShopID:    shop.ID,
				Name:      tovarMaster.Name,
				Category:  tovarMaster.Category,
				Image:     tovarMaster.Image,
				Tax:       tovarMaster.Tax,
				Measure:   tovarMaster.Measure,
				Price:     tovarMaster.Price,
				Deleted:   false,
				IsVisible: false,
			}
			if _, ok := tovarsMap[tovarMaster.ID]; ok {
				tovar.IsVisible = true
			}
			tovars = append(tovars, tovar)
		}
		err = tx.Create(&tovars).Error
		if err != nil {
			return err
		}

		ingredientMasters := []*model.IngredientMaster{}
		err = tx.Model(&model.IngredientMaster{}).Where("deleted = ?", false).Scan(&ingredientMasters).Error
		if err != nil {
			return err
		}
		ingredients := make([]*model.Ingredient, 0, 100)
		for _, ingredientMaster := range ingredientMasters {
			ingredient := &model.Ingredient{
				IngredientID: ingredientMaster.ID,
				ShopID:       shop.ID,
				Name:         ingredientMaster.Name,
				Category:     ingredientMaster.Category,
				Image:        ingredientMaster.Image,
				Measure:      ingredientMaster.Measure,
				Deleted:      false,
				IsVisible:    false,
			}
			ingredients = append(ingredients, ingredient)
		}
		err = tx.Create(&ingredients).Error
		if err != nil {
			return err
		}
		naborMasters := []*model.NaborMaster{}
		err = tx.Model(&model.NaborMaster{}).Where("deleted = ?", false).Scan(&naborMasters).Error
		if err != nil {
			return err
		}
		for _, naborMaster := range naborMasters {
			nabor := &model.Nabor{
				NaborID:  naborMaster.ID,
				ShopID:   shop.ID,
				Name:     naborMaster.Name,
				Min:      naborMaster.Min,
				Max:      naborMaster.Max,
				Replaces: naborMaster.Replaces,
				Deleted:  false,
			}
			err = tx.Create(nabor).Error
			if err != nil {
				return err
			}
			naborIngredientsMaster := []*model.IngredientNaborMaster{}
			err = tx.Model(&model.IngredientNaborMaster{}).Debug().Where("nabor_id = ?", naborMaster.ID).Scan(&naborIngredientsMaster).Error
			if err != nil {
				return err
			}
			naborIngredients := []*model.IngredientNabor{}
			for _, naborIngredientMaster := range naborIngredientsMaster {
				naborIngredient := &model.IngredientNabor{
					NaborID:      nabor.NaborID,
					IngredientID: naborIngredientMaster.IngredientID,
					ShopID:       shop.ID,
					Brutto:       naborIngredientMaster.Brutto,
					Price:        naborIngredientMaster.Price,
				}
				naborIngredients = append(naborIngredients, naborIngredient)
			}
			if len(naborIngredients) > 0 {
				err = tx.Create(&naborIngredients).Error
				if err != nil {
					return err
				}
			}
		}

		techCartsMaster := []*model.TechCartMaster{}
		err = tx.Model(&model.TechCartMaster{}).Where("deleted = ?", false).Scan(&techCartsMaster).Error
		if err != nil {
			return err
		}

		for _, techCartMaster := range techCartsMaster {
			ingredientMasters := []*model.IngredientTechCartMaster{}
			err = tx.Model(&model.IngredientTechCartMaster{}).Where("tech_cart_id = ?", techCartMaster.ID).Scan(&ingredientMasters).Error
			if err != nil {
				return err
			}

			techCart := &model.TechCart{
				TechCartID: techCartMaster.ID,
				ShopID:     shop.ID,
				Name:       techCartMaster.Name,
				Category:   techCartMaster.Category,
				Image:      techCartMaster.Image,
				Tax:        techCartMaster.Tax,
				Measure:    techCartMaster.Measure,
				Price:      techCartMaster.Price,
				Discount:   techCartMaster.Discount,
				Deleted:    false,
				IsVisible:  false,
			}
			if _, ok := techCartMap[techCartMaster.ID]; ok {
				techCart.IsVisible = true
			}
			err = tx.Create(techCart).Error
			if err != nil {
				return err
			}
			ingredients := []*model.IngredientTechCart{}
			for _, ingredientMaster := range ingredientMasters {
				ingredient := &model.IngredientTechCart{
					IngredientID: ingredientMaster.IngredientID,
					TechCartID:   techCartMaster.ID,
					ShopID:       shop.ID,
					Brutto:       ingredientMaster.Brutto,
				}
				ingredients = append(ingredients, ingredient)
			}
			if len(ingredients) > 0 {
				err = tx.Create(&ingredients).Error
				if err != nil {
					return err
				}
			}
			techCartNaborMaster := []*model.NaborTechCartMaster{}
			err = tx.Model(&model.NaborTechCartMaster{}).Where("tech_cart_id = ?", techCartMaster.ID).Scan(&techCartNaborMaster).Error
			if err != nil {
				return err
			}
			techCartNabors := []*model.NaborTechCart{}
			for _, techCartNaborMaster := range techCartNaborMaster {
				techCartNabor := &model.NaborTechCart{
					NaborID:    techCartNaborMaster.NaborID,
					TechCartID: techCartMaster.ID,
					ShopID:     shop.ID,
				}
				techCartNabors = append(techCartNabors, techCartNabor)
			}
			if len(techCartNabors) > 0 {
				err = tx.Create(&techCartNabors).Error
				if err != nil {
					return err
				}
			}
		}

		return nil
	})
	if err != nil {
		return nil, err
	}
	return shop, nil
}
func (r *ShopDB) GetAllShop(filter *model.Filter) ([]*model.Shop, int64, error) {
	var shops []*model.Shop
	res := r.gormDB.Find(&shops).Where("id IN (?)", filter.AccessibleShops)
	newRes, count, err := filter.FilterResults(res, &shops, utils.DefaultPageSize, "", "", "")
	if err != nil {
		return nil, 0, err
	}

	err = newRes.Scan(&shops).Error
	if err != nil {
		return nil, 0, err
	}
	return shops, count, nil
}
func (r *ShopDB) GetShop(id int) (*model.Shop, error) {
	shop := &model.Shop{}
	err := r.gormDB.Where("id = ?", id).First(&shop).Error
	if err != nil {
		return nil, err
	}
	return shop, nil
}
func (r *ShopDB) UpdateShop(shop *model.Shop) error {
	err := r.gormDB.Model(&model.Stolik{}).Where("shop_id = ?", shop.ID).Delete(&model.Stolik{}).Error
	if err != nil {
		return err
	}
	stoliki := shop.Stoliki
	err = r.gormDB.Create(&stoliki).Error
	if err != nil {
		return err
	}
	err = r.gormDB.Model(&model.Shop{}).Select("*").Where("id = ?", shop.ID).Updates(shop).Error
	if err != nil {
		return err
	}
	if shop.Blocked {
		shift := &model.Shift{}
		err := r.gormDB.Model(&model.Shift{}).Where("shop_id = ?", shop.ID).Last(&model.Shift{}).Scan(&shift).Error
		if err != nil {
			return err
		}
		if !shift.IsClosed {
			shift.IsClosed = true
			shift.ClosedAt = time.Now()
			err = r.gormDB.Model(&model.Shift{}).Where("id = ?", shift.ID).Updates(shift).Error
			if err != nil {
				return err
			}
		}
	}
	return nil
}
func (r *ShopDB) DeleteShop(id int) error {
	err := r.gormDB.Delete(&model.Shop{}, id).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *ShopDB) GetShopBySchetID(id int) (*model.Shop, error) {
	shop := &model.Shop{}
	err := r.gormDB.Where("cash_schet = ? OR card_schet = ?", id, id).First(&shop).Error
	if err != nil {
		return nil, err
	}
	return shop, nil
}

func (r *ShopDB) GetShopByCashSchetID(id int) (*model.Shop, error) {
	shop := &model.Shop{}
	res := r.gormDB.Model(&model.Shop{}).Where("cash_schet = ?", id).Scan(shop)
	if res.Error != nil {
		return nil, res.Error
	}
	if res.RowsAffected == 0 {
		return nil, nil
	}

	return shop, nil
}

func (r *ShopDB) GetAllShopsWithouParam() ([]*model.Shop, error) {
	var shops []*model.Shop

	err := r.gormDB.Find(&shops).Error
	if err != nil {
		return nil, err
	}

	return shops, nil
}

func (r *ShopDB) GetAllPureShops() ([]*model.Shop, error) {
	var shops []*model.Shop
	err := r.gormDB.Model(&model.Shop{}).Find(&shops).Error
	if err != nil {
		return nil, err
	}
	return shops, nil
}

func (r *ShopDB) GetRevenueByShopID(id int) (float32, error) {
	var revenue float32
	now := time.Now()
	err := r.gormDB.Model(&model.Check{}).Where("shop_id = ? and closed_at::date = ?::date", id, now).Select("COALESCE(SUM(card+cash), 0)").Scan(&revenue).Error
	if err != nil {
		if err != gorm.ErrRecordNotFound {
			return 0, err
		}
	}
	return revenue, nil
}
