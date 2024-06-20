package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"math"
	"strconv"
	"sync"
	"time"
	"zebra/model"
	"zebra/utils"

	"gorm.io/gorm"
)

type SkladDB struct {
	db     *sql.DB
	gormDB *gorm.DB
}

func NewSkladDB(db *sql.DB, gormDB *gorm.DB) *SkladDB {
	return &SkladDB{db: db, gormDB: gormDB}
}

func (r *SkladDB) RecalculateCost(tx *gorm.DB, item *model.ItemPostavka, status string, skladID int, flag bool) error {
	if status == utils.Postavka || status == utils.UpdatePostavka {
		if item.Type == utils.TypeIngredient {
			skladIngredient := &model.SkladIngredient{}
			err := tx.Model(skladIngredient).Where("sklad_id = ? and ingredient_id = ?", skladID, item.ItemID).Scan(skladIngredient).Error
			if err != nil {
				return err
			}

			costBefore := skladIngredient.Cost

			if skladIngredient.Quantity < 0 {
				if skladIngredient.Quantity+item.Quantity > 0 {
					if flag {
						skladIngredient.Quantity += item.Quantity
					}
					skladIngredient.Cost = item.Cost
				} else {
					if flag {
						skladIngredient.Quantity += item.Quantity
					}
					invCost := &model.InventarizationForNetCost{}
					err := tx.Model(&model.InventarizationItem{}).Select("time, cost, fact_quantity").Where("inventarization_items.item_id = ? and inventarization_items.type = ? and inventarization_items.sklad_id = ? and inventarization_items.time <= ?", item.ItemID, item.Type, skladID, time.Now()).Order("time desc").Limit(1).Scan(&invCost).Error
					if err != nil {
						return err
					}
					var netCost []*model.NetCost
					err = tx.Model(&model.ItemPostavka{}).Select("cost, quantity").Joins("inner join postavkas on postavkas.id = item_postavkas.postavka_id").Where("item_id = ? and sklad_id = ? and postavkas.time >= ? and postavkas.deleted = ?", item.ItemID, skladID, invCost.Time, false).Scan(&netCost).Error
					if err != nil {
						return err
					}

					var sum float32 = 0
					var quantity float32 = 0
					for _, val := range netCost {
						sum += (val.Cost * val.Quantity) //actual also consider
						quantity += val.Quantity
					}
					if invCost.FactQuantity > 0 {
						sum += (invCost.Cost * invCost.FactQuantity)
						quantity += invCost.FactQuantity
					}
					if quantity != 0 {
						skladIngredient.Cost = sum / quantity
					}
				}
			} else {
				skladIngredient.Cost = ((skladIngredient.Cost * skladIngredient.Quantity) + (item.Cost * item.Quantity)) / (skladIngredient.Quantity + item.Quantity)
				if flag {
					skladIngredient.Quantity += item.Quantity
				}
			}

			var prevCost float32 = 0
			err = tx.Model(&model.SkladIngredient{}).Select("cost").Where("ingredient_id = ? and sklad_id = ?", item.ItemID, skladID).Scan(&prevCost).Error
			if err != nil {
				return err
			}
			if skladIngredient.Cost < 0 {
				skladIngredient.Cost = costBefore
			}
			err = tx.Where("sklad_id = ? and ingredient_id = ?", skladIngredient.SkladID, skladIngredient.IngredientID).Save(skladIngredient).Error
			if err != nil {
				return err
			}

			if item.Cost >= prevCost+(prevCost*0.5) && prevCost != 0 {
				item.Risky = true
				err := tx.Model(&model.ItemPostavka{}).Select("risky").Where("id = ?", item.ID).Updates(item).Error
				if err != nil {
					return err
				}
			}
		} else if item.Type == utils.TypeTovar {
			skladTovar := &model.SkladTovar{}
			err := tx.Model(skladTovar).Where("sklad_id = ? and tovar_id = ?", skladID, item.ItemID).Scan(skladTovar).Error
			if err != nil {
				return err
			}

			costBefore := skladTovar.Cost

			if skladTovar.Quantity < 0 {
				if skladTovar.Quantity+item.Quantity > 0 {
					if flag {
						skladTovar.Quantity += item.Quantity
					}
					skladTovar.Cost = item.Cost
				} else {
					if flag {
						skladTovar.Quantity += item.Quantity
					}
					invCost := &model.InventarizationForNetCost{}
					err := tx.Model(&model.InventarizationItem{}).Select("time, cost, fact_quantity").Where("inventarization_items.item_id = ? and inventarization_items.type = ? and inventarization_items.sklad_id = ? and inventarization_items.time <= ?", item.ItemID, item.Type, skladID, time.Now()).Order("time desc").Limit(1).Scan(&invCost).Error
					if err != nil {
						return err
					}
					var netCost []*model.NetCost
					err = tx.Model(&model.ItemPostavka{}).Select("cost, quantity").Joins("inner join postavkas on postavkas.id = item_postavkas.postavka_id").Where("item_id = ? and sklad_id = ? and postavkas.time >= ? and postavkas.deleted = ?", item.ItemID, skladID, invCost.Time, false).Scan(&netCost).Error
					if err != nil {
						return err
					}
					var sum float32 = 0
					var quantity float32 = 0
					for _, val := range netCost {
						sum += (val.Cost * val.Quantity)
						quantity += val.Quantity
					}
					if invCost.FactQuantity > 0 {
						sum += (invCost.Cost * invCost.FactQuantity)
						quantity += invCost.FactQuantity
					}
					if quantity != 0 {
						skladTovar.Cost = sum / quantity
					}
				}
			} else {
				skladTovar.Cost = ((skladTovar.Cost * skladTovar.Quantity) + (item.Cost * item.Quantity)) / (skladTovar.Quantity + item.Quantity)
				if flag {
					skladTovar.Quantity += item.Quantity
				}
			}
			var prevCost float32 = 0
			err = tx.Model(&model.SkladTovar{}).Select("cost").Where("tovar_id = ? and sklad_id = ?", item.ItemID, skladID).Scan(&prevCost).Error
			if err != nil {
				return err
			}
			if skladTovar.Cost < 0 {
				skladTovar.Cost = costBefore
			}
			err = tx.Where("sklad_id = ? and tovar_id = ?", skladTovar.SkladID, skladTovar.TovarID).Save(skladTovar).Error
			if err != nil {
				return err
			}

			if item.Cost >= prevCost+(prevCost*0.5) && prevCost != 0 {
				item.Risky = true
				err := tx.Model(&model.ItemPostavka{}).Select("risky").Where("id = ?", item.ID).Updates(item).Error
				if err != nil {
					return err
				}
			}
		} else {
			return errors.New("incorrect type for recalculating netCost")
		}
	} else if status == utils.DeletePostavka {
		if item.Type == utils.TypeIngredient {
			skladIngredient := &model.SkladIngredient{}
			err := tx.Model(skladIngredient).Where("sklad_id = ? and ingredient_id = ?", skladID, item.ItemID).Scan(skladIngredient).Error
			if err != nil {
				return err
			}

			costBefore := skladIngredient.Cost

			if skladIngredient.Quantity > 0 {
				if skladIngredient.Quantity-item.Quantity > 0 {
					skladIngredient.Cost = ((skladIngredient.Cost * skladIngredient.Quantity) - (item.Cost * item.Quantity)) / (skladIngredient.Quantity - item.Quantity)
					if flag {
						skladIngredient.Quantity -= item.Quantity
					}
				} else {
					if flag {
						skladIngredient.Quantity -= item.Quantity
					}
					invCost := &model.InventarizationForNetCost{}
					err := tx.Model(&model.InventarizationItem{}).Select("time, cost, fact_quantity").Where("inventarization_items.item_id = ? and inventarization_items.type = ? and inventarization_items.sklad_id = ? and inventarization_items.time <= ?", item.ItemID, item.Type, skladID, time.Now()).Order("time desc").Limit(1).Scan(&invCost).Error
					if err != nil {
						return err
					}
					var netCost []*model.NetCost
					err = tx.Model(&model.ItemPostavka{}).Select("cost, quantity").Joins("inner join postavkas on postavkas.id = item_postavkas.postavka_id").Where("item_id = ? and sklad_id = ? and postavkas.time >= ? and postavkas.deleted = ?", item.ItemID, skladID, invCost.Time, false).Scan(&netCost).Error
					if err != nil {
						return err
					}
					var sum float32 = 0
					var quantity float32 = 0
					for _, val := range netCost {
						sum += (val.Cost * val.Quantity)
						quantity += val.Quantity
					}
					if invCost.FactQuantity > 0 {
						sum += (invCost.Cost * invCost.FactQuantity)
						quantity += invCost.FactQuantity
					}
					if quantity != 0 {
						skladIngredient.Cost = sum / quantity
					}
				}
			} else {
				if flag {
					skladIngredient.Quantity -= item.Quantity
				}
				invCost := &model.InventarizationForNetCost{}
				err := tx.Model(&model.InventarizationItem{}).Select("time, cost, fact_quantity").Where("inventarization_items.item_id = ? and inventarization_items.type = ? and inventarization_items.sklad_id = ? and inventarization_items.time <= ?", item.ItemID, item.Type, skladID, time.Now()).Order("time desc").Limit(1).Scan(&invCost).Error
				if err != nil {
					return err
				}
				var netCost []*model.NetCost
				err = tx.Model(&model.ItemPostavka{}).Select("cost, quantity").Joins("inner join postavkas on postavkas.id = item_postavkas.postavka_id").Where("item_id = ? and sklad_id = ? and postavkas.time >= ? and postavkas.deleted = ?", item.ItemID, skladID, invCost.Time, false).Scan(&netCost).Error
				if err != nil {
					return err
				}
				var sum float32 = 0
				var quantity float32 = 0
				for _, val := range netCost {
					sum += (val.Cost * val.Quantity)
					quantity += val.Quantity
				}
				if invCost.FactQuantity > 0 {
					sum += (invCost.Cost * invCost.FactQuantity)
					quantity += invCost.FactQuantity
				}
				if quantity != 0 {
					skladIngredient.Cost = sum / quantity
				}
			}

			if skladIngredient.Cost < 0 {
				skladIngredient.Cost = costBefore
			}

			err = tx.Where("sklad_id = ? and ingredient_id = ?", skladIngredient.SkladID, skladIngredient.IngredientID).Save(skladIngredient).Error
			if err != nil {
				return err
			}
		} else if item.Type == utils.TypeTovar {
			skladTovar := &model.SkladTovar{}
			err := tx.Model(skladTovar).Where("sklad_id = ? and tovar_id = ?", skladID, item.ItemID).Scan(skladTovar).Error
			if err != nil {
				return err
			}

			costBefore := skladTovar.Cost

			if skladTovar.Quantity > 0 {
				if skladTovar.Quantity-item.Quantity > 0 {
					skladTovar.Cost = ((skladTovar.Cost * skladTovar.Quantity) - (item.Cost * item.Quantity)) / (skladTovar.Quantity - item.Quantity)
					if flag {
						skladTovar.Quantity -= item.Quantity
					}
				} else {
					if flag {
						skladTovar.Quantity -= item.Quantity
					}
					invCost := &model.InventarizationForNetCost{}
					err := tx.Model(&model.InventarizationItem{}).Select("time, cost, fact_quantity").Where("inventarization_items.item_id = ? and inventarization_items.type = ? and inventarization_items.sklad_id = ? and inventarization_items.time <= ?", item.ItemID, item.Type, skladID, time.Now()).Order("time desc").Limit(1).Scan(&invCost).Error
					if err != nil {
						return err
					}
					var netCost []*model.NetCost
					err = tx.Model(&model.ItemPostavka{}).Select("cost, quantity").Joins("inner join postavkas on postavkas.id = item_postavkas.postavka_id").Where("item_id = ? and sklad_id = ? and postavkas.time >= ? and postavkas.deleted = ?", item.ItemID, skladID, invCost.Time, false).Scan(&netCost).Error
					if err != nil {
						return err
					}
					var sum float32 = 0
					var quantity float32 = 0
					for _, val := range netCost {
						sum += (val.Cost * val.Quantity)
						quantity += val.Quantity
					}
					if invCost.FactQuantity > 0 {
						sum += (invCost.Cost * invCost.FactQuantity)
						quantity += invCost.FactQuantity
					}
					if quantity != 0 {
						skladTovar.Cost = sum / quantity
					}
				}
			} else {
				if flag {
					skladTovar.Quantity -= item.Quantity
				}
				invCost := &model.InventarizationForNetCost{}
				err := tx.Model(&model.InventarizationItem{}).Select("time, cost, fact_quantity").Where("inventarization_items.item_id = ? and inventarization_items.type = ? and inventarization_items.sklad_id = ? and inventarization_items.time <= ?", item.ItemID, item.Type, skladID, time.Now()).Order("time desc").Limit(1).Scan(&invCost).Error
				if err != nil {
					return err
				}
				var netCost []*model.NetCost
				err = tx.Model(&model.ItemPostavka{}).Select("cost, quantity").Joins("inner join postavkas on postavkas.id = item_postavkas.postavka_id").Where("item_id = ? and sklad_id = ? and postavkas.time >= ? and postavkas.deleted = ?", item.ItemID, skladID, invCost.Time, false).Scan(&netCost).Error
				if err != nil {
					return err
				}
				var sum float32 = 0
				var quantity float32 = 0
				for _, val := range netCost {
					sum += (val.Cost * val.Quantity)
					quantity += val.Quantity
				}
				if invCost.FactQuantity > 0 {
					sum += (invCost.Cost * invCost.FactQuantity)
					quantity += invCost.FactQuantity
				}
				if quantity != 0 {
					skladTovar.Cost = sum / quantity
				}
			}

			if skladTovar.Cost < 0 {
				skladTovar.Cost = costBefore
			}

			err = tx.Where("sklad_id = ? and tovar_id = ?", skladTovar.SkladID, skladTovar.TovarID).Save(skladTovar).Error
			if err != nil {
				return err
			}
		} else {
			return errors.New("incorrect type for recalculating netCost")
		}
	} else {
		return errors.New("incorrect status for recalculating netCost")
	}
	return nil
}

func (r *SkladDB) RemoveFromSklad(spisanie *model.RemoveFromSklad) error {
	invItems := []*model.InventarizationItem{}
	trafficItems := []*model.AsyncJob{}
	var shopID int
	if err := r.gormDB.Model(&model.Sklad{}).Select("shop_id").Where("id = ?", spisanie.SkladID).Scan(&shopID).Error; err != nil {
		return err
	}
	err := r.gormDB.Transaction(func(tx *gorm.DB) error {
		for _, item := range spisanie.Items {
			flag := true
			invItem := &model.InventarizationItem{}
			if item.Type == utils.TypeIngredient {
				ingredient := &model.Ingredient{}
				err := tx.Model(&model.Ingredient{}).Select("ingredients.ingredient_id, ingredients.name, ingredients.measure, sklad_ingredients.cost").Joins("inner join sklad_ingredients on sklad_ingredients.ingredient_id = ingredients.ingredient_id").Where("ingredients.ingredient_id = ? and sklad_ingredients.sklad_id = ?", item.ItemID, spisanie.SkladID).Scan(&ingredient).Error
				if err != nil {
					return err
				}
				skladIngredient := &model.SkladIngredient{}
				if err := tx.Where("sklad_id = ? and ingredient_id = ?", spisanie.SkladID, item.ItemID).First(skladIngredient).Error; err != nil {
					if err != gorm.ErrRecordNotFound {
						return err
					}
				}
				item.Cost = skladIngredient.Cost * item.Quantity
				spisanie.Cost += skladIngredient.Cost * item.Quantity

				var lastInv time.Time
				err = tx.Model(&model.InventarizationItem{}).Select("time").Where("sklad_id = ? and item_id = ? and type = ? and status != ?", spisanie.SkladID, item.ItemID, utils.TypeIngredient, utils.StatusNew).Order("time desc").First(&lastInv).Error
				if err != nil {
					if err != gorm.ErrRecordNotFound {
						return err
					}
				}
				if lastInv.After(spisanie.Time) && lastInv.Before(time.Now()) {
					flag = false
				}
				if flag {
					skladIngredient := &model.SkladIngredient{}
					if err := tx.Where("sklad_id = ? and ingredient_id = ?", spisanie.SkladID, item.ItemID).First(skladIngredient).Error; err != nil {
						return err
					}
					skladIngredient.Quantity -= item.Quantity

					if err := tx.Model(skladIngredient).Select("quantity").Where("sklad_id = ? and ingredient_id = ?", spisanie.SkladID, item.ItemID).Updates(skladIngredient).Error; err != nil {
						return err
					}
				}
				if spisanie.Time.Year() != time.Now().Year() || spisanie.Time.Month() != time.Now().Month() || spisanie.Time.Day() != time.Now().Day() {
					trafficItems = append(trafficItems, &model.AsyncJob{
						ItemID:    item.ItemID,
						ItemType:  item.Type,
						SkladID:   spisanie.SkladID,
						TimeStamp: spisanie.Time.UTC(),
						CreatedAt: time.Now(),
						Status:    utils.StatusNeedRecalculate,
					})
				}
				err = tx.Model(invItem).Where("sklad_id = ? and item_id = ? and type = ? and time > ? and before_time < ?", spisanie.SkladID, item.ItemID, item.Type, spisanie.Time, spisanie.Time).Scan(invItem).Error
				if err != nil {
					if err.Error() != "record not found" {
						return err
					}
				}
				if invItem.ID != 0 && invItem.ItemID != 0 {
					invItems = append(invItems, invItem)
				}
			} else if item.Type == utils.TypeTovar {
				tovar := &model.Tovar{}
				err := tx.Model(&model.Tovar{}).Select("tovars.tovar_id, tovars.name, tovars.measure, sklad_tovars.cost").Joins("inner join sklad_tovars on sklad_tovars.tovar_id = tovars.id").Where("tovars.tovar_id = ? and sklad_tovars.sklad_id = ?", item.ItemID, spisanie.SkladID).Scan(&tovar).Error
				if err != nil {
					return err
				}
				skladTovar := &model.SkladTovar{}
				if err := tx.Where("sklad_id = ? and tovar_id = ?", spisanie.SkladID, item.ItemID).First(skladTovar).Error; err != nil {
					if err != gorm.ErrRecordNotFound {
						return err
					}
				}
				item.Cost = skladTovar.Cost * item.Quantity
				spisanie.Cost += skladTovar.Cost * item.Quantity
				var lastInv time.Time
				err = tx.Model(&model.InventarizationItem{}).Select("time").Where("sklad_id = ? and item_id = ? and type = ? and status != ?", spisanie.SkladID, item.ItemID, utils.TypeTovar, utils.StatusNew).Order("time desc").First(&lastInv).Error
				if err != nil {
					if err != gorm.ErrRecordNotFound {
						return err
					}
				}
				if lastInv.After(spisanie.Time) && lastInv.Before(time.Now()) {
					flag = false
				}
				if flag {
					skladTovar := &model.SkladTovar{}
					if err := tx.Where("sklad_id = ? and tovar_id = ?", spisanie.SkladID, item.ItemID).First(skladTovar).Error; err != nil {
						return err
					}
					skladTovar.Quantity -= item.Quantity
					if err := tx.Model(skladTovar).Select("quantity").Where("sklad_id = ? and tovar_id = ?", spisanie.SkladID, item.ItemID).Updates(skladTovar).Error; err != nil {
						return err
					}
				}
				if spisanie.Time.Year() != time.Now().Year() || spisanie.Time.Month() != time.Now().Month() || spisanie.Time.Day() != time.Now().Day() {
					trafficItems = append(trafficItems, &model.AsyncJob{
						ItemID:    item.ItemID,
						ItemType:  item.Type,
						SkladID:   spisanie.SkladID,
						TimeStamp: spisanie.Time.UTC(),
						CreatedAt: time.Now(),
						Status:    utils.StatusNeedRecalculate,
					})
				}
				err = tx.Model(invItem).Where("sklad_id = ? and item_id = ? and type = ? and time > ? and before_time < ?", spisanie.SkladID, item.ItemID, item.Type, spisanie.Time, spisanie.Time).Scan(invItem).Error
				if err != nil {
					if err.Error() != "record not found" {
						return err
					}
				}
				if invItem.ID != 0 && invItem.ItemID != 0 {
					invItems = append(invItems, invItem)
				}
			} else if item.Type == utils.TypeTechCart {
				techCart := &model.TechCart{}
				if err := tx.Model(&model.TechCart{}).Where("tech_cart_id = ? and shop_id = ?", item.ItemID, shopID).First(techCart).Error; err != nil {
					return err
				}
				ingredientsTech := []*model.IngredientTechCart{}
				if err := tx.Model(&model.IngredientTechCart{}).Where("tech_cart_id = ? and shop_id = ?", item.ItemID, shopID).Scan(&ingredientsTech).Error; err != nil {
					return err
				}
				for _, ingredient := range ingredientsTech {
					removeItem := &model.RemoveFromSkladItem{
						Type:           utils.TypeIngredient,
						Quantity:       ingredient.Brutto * item.Quantity,
						Cost:           0,
						Details:        item.Details,
						SkladID:        item.SkladID,
						ItemID:         ingredient.IngredientID,
						PartOfTechCart: true,
					}
					flag = true
					ingredientNew := &model.Ingredient{}
					err := tx.Model(&model.Ingredient{}).Select("ingredients.ingredient_id, ingredients.name, ingredients.measure, sklad_ingredients.cost").Joins("inner join sklad_ingredients on sklad_ingredients.ingredient_id = ingredients.ingredient_id").Where("ingredients.ingredient_id = ? and sklad_ingredients.sklad_id = ?", ingredient.IngredientID, spisanie.SkladID).Scan(&ingredientNew).Error
					if err != nil {
						return err
					}
					skladIngredient := &model.SkladIngredient{}
					if err := tx.Where("sklad_id = ? and ingredient_id = ?", spisanie.SkladID, ingredient.IngredientID).First(skladIngredient).Error; err != nil {
						return err
					}
					item.Cost += skladIngredient.Cost * ingredient.Brutto * item.Quantity
					removeItem.Cost = skladIngredient.Cost * ingredient.Brutto * item.Quantity
					spisanie.Items = append(spisanie.Items, removeItem)
					spisanie.Cost += skladIngredient.Cost * ingredient.Brutto * item.Quantity

					var lastInv time.Time
					err = tx.Model(&model.InventarizationItem{}).Select("time").Where("sklad_id = ? and item_id = ? and type = ? and status != ?", spisanie.SkladID, item.ItemID, utils.TypeIngredient, utils.StatusNew).Order("time desc").First(&lastInv).Error
					if err != nil {
						if err != gorm.ErrRecordNotFound {
							return err
						}
					}
					if lastInv.After(spisanie.Time) && lastInv.Before(time.Now()) {
						flag = false
					}
					if flag {
						skladIngredient := &model.SkladIngredient{}
						if err := tx.Where("sklad_id = ? and ingredient_id = ?", spisanie.SkladID, ingredient.IngredientID).First(skladIngredient).Error; err != nil {
							return err
						}
						skladIngredient.Quantity -= ingredient.Brutto * item.Quantity
						if err := tx.Model(skladIngredient).Select("quantity").Where("sklad_id = ? and ingredient_id = ?", spisanie.SkladID, ingredient.IngredientID).Updates(skladIngredient).Error; err != nil {
							return err
						}
					}

					if spisanie.Time.Year() != time.Now().Year() || spisanie.Time.Month() != time.Now().Month() || spisanie.Time.Day() != time.Now().Day() {
						trafficItems = append(trafficItems, &model.AsyncJob{
							ItemID:    ingredient.IngredientID,
							ItemType:  utils.TypeIngredient,
							SkladID:   spisanie.SkladID,
							TimeStamp: spisanie.Time.UTC(),
							CreatedAt: time.Now(),
							Status:    utils.StatusNeedRecalculate,
						})
					}
					err = tx.Model(invItem).Where("sklad_id = ? and item_id = ? and type = ? and time > ? and before_time < ?", spisanie.SkladID, ingredient.IngredientID, utils.TypeIngredient, spisanie.Time, spisanie.Time).Scan(invItem).Error
					if err != nil {
						if err.Error() != "record not found" {
							return err
						}
					}
					if invItem.ID != 0 && invItem.ItemID != 0 {
						invItems = append(invItems, invItem)
					}
				}
			}

		}
		spisanie.Status = utils.StatusClosed
		if err := tx.Model(spisanie).Create(spisanie).Error; err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return err
	}
	r.ConcurrentRecalculationForDailyStatistics(trafficItems)
	if len(invItems) > 0 {
		err = r.RecalculateInventarization(invItems)
		if err != nil {
			return err
		}
	}
	return nil
}

func (r *SkladDB) RequestToRemove(request *model.RemoveFromSklad) error {
	var shopID int
	if err := r.gormDB.Model(&model.Sklad{}).Select("shop_id").Where("id = ?", request.SkladID).Scan(&shopID).Error; err != nil {
		return err
	}
	err := r.gormDB.Transaction(func(tx *gorm.DB) error {
		for _, item := range request.Items {
			if item.Type == utils.TypeIngredient {
				ingredient := &model.Ingredient{}
				err := tx.Model(&model.Ingredient{}).Select("ingredients.ingredient_id, ingredients.name, ingredients.measure, sklad_ingredients.cost").Joins("inner join sklad_ingredients on sklad_ingredients.ingredient_id = ingredients.ingredient_id").Where("ingredients.ingredient_id = ? and sklad_ingredients.sklad_id = ?", item.ItemID, request.SkladID).Scan(&ingredient).Error
				if err != nil {
					return err
				}
				skladIngredient := &model.SkladIngredient{}
				if err := tx.Where("sklad_id = ? and ingredient_id = ?", request.SkladID, item.ItemID).First(skladIngredient).Error; err != nil {
					return err
				}
				request.Cost += skladIngredient.Cost * item.Quantity
				item.Cost = skladIngredient.Cost * item.Quantity
			} else if item.Type == utils.TypeTovar {
				tovar := &model.Tovar{}
				err := tx.Model(&model.Tovar{}).Select("tovars.tovar_id, tovars.name, tovars.measure, sklad_tovars.cost").Joins("inner join sklad_tovars on sklad_tovars.tovar_id = tovars.tovar_id").Where("tovars.tovar_id = ? and sklad_tovars.sklad_id = ?", item.ItemID, request.SkladID).Scan(&tovar).Error
				if err != nil {
					return err
				}
				skladTovar := &model.SkladTovar{}
				if err := tx.Where("sklad_id = ? and tovar_id = ?", request.SkladID, item.ItemID).First(skladTovar).Error; err != nil {
					return err
				}
				request.Cost += skladTovar.Cost * item.Quantity
				item.Cost = skladTovar.Cost * item.Quantity
			} else if item.Type == utils.TypeTechCart {
				techCart := &model.TechCart{}
				if err := tx.Model(&model.TechCart{}).Where("tech_cart_id = ? and shop_id = ?", item.ItemID, shopID).First(techCart).Error; err != nil {
					return err
				}
				ingredientsTech := []*model.IngredientTechCart{}
				if err := tx.Model(&model.IngredientTechCart{}).Where("tech_cart_id = ? and shop_id = ?", item.ItemID, shopID).Scan(&ingredientsTech).Error; err != nil {
					return err
				}
				for _, ingredient := range ingredientsTech {
					removeItem := &model.RemoveFromSkladItem{
						Type:           utils.TypeIngredient,
						Quantity:       ingredient.Brutto * item.Quantity,
						Cost:           0,
						Details:        item.Details,
						SkladID:        item.SkladID,
						ItemID:         ingredient.IngredientID,
						PartOfTechCart: true,
					}
					ingredientNew := &model.Ingredient{}
					err := tx.Model(&model.Ingredient{}).Select("ingredients.ingredient_id, ingredients.name, ingredients.measure, sklad_ingredients.cost").Joins("inner join sklad_ingredients on sklad_ingredients.ingredient_id = ingredients.ingredient_id").Where("ingredients.ingredient_id = ? and sklad_ingredients.sklad_id = ?", ingredient.IngredientID, request.SkladID).Scan(&ingredientNew).Error
					if err != nil {
						return err
					}
					skladIngredient := &model.SkladIngredient{}
					if err := tx.Where("sklad_id = ? and ingredient_id = ?", request.SkladID, ingredient.IngredientID).First(skladIngredient).Error; err != nil {
						return err
					}
					item.Cost += skladIngredient.Cost * ingredient.Brutto * item.Quantity
					removeItem.Cost = skladIngredient.Cost * ingredient.Brutto * item.Quantity
					request.Items = append(request.Items, removeItem)
					request.Cost += skladIngredient.Cost * ingredient.Brutto * item.Quantity
					if err := tx.Model(skladIngredient).Where("sklad_id = ? and ingredient_id = ?", request.SkladID, ingredient.IngredientID).Updates(skladIngredient).Error; err != nil {
						return err
					}
				}
			}
		}
		request.Status = utils.StatusOpened
		if err := tx.Model(request).Create(request).Error; err != nil {
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

func (r *SkladDB) ConfirmToRemove(id int) error {
	spisanie := model.RemoveFromSklad{}
	err := r.gormDB.Model(spisanie).Select("*").Where("remove_from_sklads.id = ?", id).Scan(&spisanie).Error
	if err != nil {
		return err
	}
	var shopID int
	if err := r.gormDB.Model(&model.Sklad{}).Select("shop_id").Where("id = ?", spisanie.SkladID).Scan(&shopID).Error; err != nil {
		return err
	}
	items := []*model.RemoveFromSkladItem{}

	err = r.gormDB.Model(items).Select("*").Where("remove_from_sklad_items.remove_id = ?", spisanie.ID).Scan(&items).Error
	if err != nil {
		return err
	}

	spisanie.Items = items
	invItems := []*model.InventarizationItem{}
	trafficItems := []*model.AsyncJob{}

	err = r.gormDB.Transaction(func(tx *gorm.DB) error {
		for _, item := range spisanie.Items {
			flag := true
			invItem := &model.InventarizationItem{}
			if item.Type == utils.TypeIngredient {
				var lastInv time.Time
				err = tx.Model(&model.InventarizationItem{}).Select("time").Where("sklad_id = ? and item_id = ? and type = ? and status != ?", spisanie.SkladID, item.ItemID, utils.TypeIngredient, utils.StatusNew).Order("time desc").First(&lastInv).Error
				if err != nil {
					if err != gorm.ErrRecordNotFound {
						return err
					}
				}
				if lastInv.After(spisanie.Time) && lastInv.Before(time.Now()) {
					flag = false
				}
				if flag {
					skladIngredient := &model.SkladIngredient{}
					if err := tx.Where("sklad_id = ? and ingredient_id = ?", spisanie.SkladID, item.ItemID).First(skladIngredient).Error; err != nil {
						return err
					}
					skladIngredient.Quantity -= item.Quantity
					if err := tx.Model(skladIngredient).Select("quantity").Where("sklad_id = ? and ingredient_id = ?", spisanie.SkladID, item.ItemID).Save(skladIngredient).Error; err != nil {
						return err
					}
				}
				if spisanie.Time.Year() != time.Now().Year() || spisanie.Time.Month() != time.Now().Month() || spisanie.Time.Day() != time.Now().Day() {
					trafficItems = append(trafficItems, &model.AsyncJob{
						ItemID:    item.ItemID,
						ItemType:  item.Type,
						SkladID:   spisanie.SkladID,
						TimeStamp: spisanie.Time.UTC(),
						CreatedAt: time.Now(),
						Status:    utils.StatusNeedRecalculate,
					})
				}
				err := tx.Model(invItem).Where("sklad_id = ? and item_id = ? and type = ? and time > ? and before_time < ?", spisanie.SkladID, item.ItemID, item.Type, spisanie.Time, spisanie.Time).Scan(invItem).Error
				if err != nil {
					if err.Error() != "record not found" {
						return err
					}
				}
				if invItem.ID != 0 && invItem.ItemID != 0 {
					invItems = append(invItems, invItem)
				}
			} else if item.Type == utils.TypeTovar {

				var lastInv time.Time
				err = tx.Model(&model.InventarizationItem{}).Select("time").Where("sklad_id = ? and item_id = ? and type = ? and status != ?", spisanie.SkladID, item.ItemID, utils.TypeTovar, utils.StatusNew).Order("time desc").First(&lastInv).Error
				if err != nil {
					if err != gorm.ErrRecordNotFound {
						return err
					}
				}
				if lastInv.After(spisanie.Time) && lastInv.Before(time.Now()) {
					flag = false
				}
				if flag {
					skladTovar := &model.SkladTovar{}
					if err := tx.Where("sklad_id = ? and tovar_id = ?", spisanie.SkladID, item.ItemID).First(skladTovar).Error; err != nil {
						return err
					}
					skladTovar.Quantity -= item.Quantity

					if err := tx.Model(skladTovar).Select("quantity").Where("sklad_id = ? and tovar_id = ?", spisanie.SkladID, item.ItemID).Save(skladTovar).Error; err != nil {
						return err
					}
				}
				if spisanie.Time.Year() != time.Now().Year() || spisanie.Time.Month() != time.Now().Month() || spisanie.Time.Day() != time.Now().Day() {
					trafficItems = append(trafficItems, &model.AsyncJob{
						ItemID:    item.ItemID,
						ItemType:  item.Type,
						SkladID:   spisanie.SkladID,
						TimeStamp: spisanie.Time.UTC(),
						CreatedAt: time.Now(),
						Status:    utils.StatusNeedRecalculate,
					})
				}
				err := tx.Model(invItem).Where("sklad_id = ? and item_id = ? and type = ? and time > ? and before_time < ?", spisanie.SkladID, item.ItemID, item.Type, spisanie.Time, spisanie.Time).Scan(invItem).Error
				if err != nil {
					if err.Error() != "record not found" {
						return err
					}
				}
				if invItem.ID != 0 && invItem.ItemID != 0 {
					invItems = append(invItems, invItem)
				}
				// } else if item.Type == utils.TypeTechCart {
				// 	techCart := &model.TechCart{}
				// 	if err := tx.Model(&model.TechCart{}).Where("tech_cart_id = ? and shop_id = ?", item.ItemID, shopID).First(techCart).Error; err != nil {
				// 		return err
				// 	}
				// 	ingredientsTech := []*model.IngredientTechCart{}
				// 	if err := tx.Model(&model.IngredientTechCart{}).Where("tech_cart_id = ? and shop_id = ?", item.ItemID, shopID).Scan(&ingredientsTech).Error; err != nil {
				// 		return err
				// 	}
				// 	for _, ingredient := range ingredientsTech {
				// 		removeItem := &model.RemoveFromSkladItem{
				// 			Type:           utils.TypeIngredient,
				// 			Quantity:       ingredient.Brutto * item.Quantity,
				// 			Cost:           0,
				// 			Details:        item.Details,
				// 			SkladID:        item.SkladID,
				// 			ItemID:         ingredient.IngredientID,
				// 			PartOfTechCart: true,
				// 		}
				// 		spisanie.Items = append(spisanie.Items, removeItem)
				// 		flag = true
				// 		var lastInv time.Time
				// 		err = tx.Model(&model.InventarizationItem{}).Select("time").Where("sklad_id = ? and item_id = ? and type = ? and status != ?", spisanie.SkladID, item.ItemID, utils.TypeIngredient, utils.StatusNew).Order("time desc").First(&lastInv).Error
				// 		if err != nil {
				// 			if err != gorm.ErrRecordNotFound {
				// 				return err
				// 			}
				// 		}
				// 		if lastInv.After(spisanie.Time) && lastInv.Before(time.Now()) {
				// 			flag = false
				// 		}
				// 		if flag {
				// 			skladIngredient := &model.SkladIngredient{}
				// 			if err := tx.Where("sklad_id = ? and ingredient_id = ?", spisanie.SkladID, ingredient.IngredientID).First(skladIngredient).Error; err != nil {
				// 				return err
				// 			}
				// 			skladIngredient.Quantity -= ingredient.Brutto * item.Quantity
				// 			if err := tx.Model(skladIngredient).Select("quantity").Where("sklad_id = ? and ingredient_id = ?", spisanie.SkladID, ingredient.IngredientID).Updates(skladIngredient).Error; err != nil {
				// 				return err
				// 			}
				// 		}
				// 		if spisanie.Time.Year() == time.Now().Year() && spisanie.Time.Month() == time.Now().Month() && spisanie.Time.Day() != time.Now().Day() {
				// 			trafficItems = append(trafficItems, &model.AsyncJob{
				// 				ItemID:    ingredient.IngredientID,
				// 				ItemType:  utils.TypeIngredient,
				// 				SkladID:   spisanie.SkladID,
				// 				TimeStamp: spisanie.Time.UTC(),
				// 				CreatedAt: time.Now(),
				// 				Status:    utils.StatusNeedRecalculate,
				// 			})
				// 		}
				// 		err := tx.Model(invItem).Where("sklad_id = ? and item_id = ? and type = ? and time > ? and before_time < ?", spisanie.SkladID, item.ItemID, item.Type, spisanie.Time, spisanie.Time).Scan(invItem).Error
				// 		if err != nil {
				// 			if err.Error() != "record not found" {
				// 				return err
				// 			}
				// 		}
				// 		if invItem.ID != 0 && invItem.ItemID != 0 {
				// 			invItems = append(invItems, invItem)
				// 		}
				// 	}
			}
		}
		spisanie.Status = utils.StatusClosed
		err = tx.Save(spisanie).Error
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return err
	}
	r.ConcurrentRecalculationForDailyStatistics(trafficItems)
	if len(invItems) > 0 {
		err = r.RecalculateInventarization(invItems)
		if err != nil {
			return err
		}
	}
	return nil
}

func (r *SkladDB) RejectToRemove(id int) error {
	spisanie := model.RemoveFromSklad{}
	err := r.gormDB.Model(spisanie).Select("*").Where("remove_from_sklads.id = ?", id).Scan(&spisanie).Error
	if err != nil {
		return err
	}
	spisanie.Status = utils.StatusRejected
	err = r.gormDB.Save(spisanie).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *SkladDB) UpdateSpisanie(spisanie *model.RemoveFromSklad) error {
	id := spisanie.ID
	invItemMap := make(map[int]int)
	invItems := []*model.InventarizationItem{}
	trafficItems := []*model.AsyncJob{}
	var shopID int
	if err := r.gormDB.Model(&model.Sklad{}).Select("shop_id").Where("id = ?", spisanie.SkladID).Scan(&shopID).Error; err != nil {
		return err
	}
	err := r.gormDB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&model.RemoveFromSklad{}).Where("id = ?", id).Error; err != nil {
			return err
		}
		oldSpisanie := model.RemoveFromSklad{}
		err := tx.Where("id = ?", id).First(&oldSpisanie).Error
		if err != nil {
			return err
		}
		items := []*model.RemoveFromSkladItem{}

		err = tx.Model(items).Where("remove_id = ?", id).Scan(&items).Error
		if err != nil {
			return err
		}
		for _, item := range items {
			flag := true
			err := tx.Delete(item).Error
			if err != nil {
				return err
			}
			invItem := &model.InventarizationItem{}
			if oldSpisanie.Status != utils.StatusClosed {
				continue
			}

			if item.Type == utils.TypeIngredient {
				var lastInv time.Time
				err := tx.Model(&model.InventarizationItem{}).Select("time").Where("sklad_id = ? and item_id = ? and type = ? and status != ?", spisanie.SkladID, item.ItemID, utils.TypeIngredient, utils.StatusNew).Order("time desc").First(&lastInv).Error
				if err != nil {
					if err != gorm.ErrRecordNotFound {
						return err
					}
				}
				if lastInv.After(spisanie.Time) && lastInv.Before(time.Now()) {
					flag = false
				}
				if flag {
					skladIngredient := &model.SkladIngredient{}
					if err := tx.Where("sklad_id = ? and ingredient_id = ?", oldSpisanie.SkladID, item.ItemID).First(skladIngredient).Error; err != nil {
						return err
					}
					skladIngredient.Quantity += item.Quantity
					if err := tx.Model(skladIngredient).Select("quantity").Where("sklad_id = ? and ingredient_id = ?", oldSpisanie.SkladID, item.ItemID).Save(skladIngredient).Error; err != nil {
						return err
					}
				}
				if _, exists := invItemMap[invItem.ID]; !exists && invItem.ID != 0 && invItem.ItemID != 0 {
					invItemMap[invItem.ID] = invItem.ID
					invItems = append(invItems, invItem)
				}
			} else if item.Type == utils.TypeTovar {
				var lastInv time.Time
				err := tx.Model(&model.InventarizationItem{}).Select("time").Where("sklad_id = ? and item_id = ? and type = ? and status != ?", spisanie.SkladID, item.ItemID, utils.TypeTovar, utils.StatusNew).Order("time desc").First(&lastInv).Error
				if err != nil {
					if err != gorm.ErrRecordNotFound {
						return err
					}
				}
				if lastInv.After(spisanie.Time) && lastInv.Before(time.Now()) {
					flag = false
				}
				if flag {
					skladTovar := &model.SkladTovar{}
					if err := tx.Where("sklad_id = ? and tovar_id = ?", oldSpisanie.SkladID, item.ItemID).First(skladTovar).Error; err != nil {
						return err
					}
					skladTovar.Quantity += item.Quantity

					if err := tx.Model(skladTovar).Select("quantity").Where("sklad_id = ? and tovar_id = ?", oldSpisanie.SkladID, item.ItemID).Save(skladTovar).Error; err != nil {
						return err
					}
				}
				if _, exists := invItemMap[invItem.ID]; !exists && invItem.ID != 0 && invItem.ItemID != 0 {
					invItemMap[invItem.ID] = invItem.ID
					invItems = append(invItems, invItem)
				}
				// } else if item.Type == utils.TypeTechCart {
				// 	techCart := &model.TechCart{}
				// 	if err := tx.Model(&model.TechCart{}).Where("tech_cart_id = ? and shop_id = ?", item.ItemID, shopID).First(techCart).Error; err != nil {
				// 		return err
				// 	}
				// 	ingredientsTech := []*model.IngredientTechCart{}
				// 	if err := tx.Model(&model.IngredientTechCart{}).Where("tech_cart_id = ? and shop_id = ?", item.ItemID, shopID).Scan(&ingredientsTech).Error; err != nil {
				// 		return err
				// 	}
				// 	for _, ingredient := range ingredientsTech {
				// 		flag = true
				// 		var lastInv time.Time
				// 		err = tx.Model(&model.InventarizationItem{}).Select("time").Where("sklad_id = ? and item_id = ? and type = ? and status != ?", spisanie.SkladID, item.ItemID, utils.TypeIngredient, utils.StatusNew).Order("time desc").First(&lastInv).Error
				// 		if err != nil {
				// 			if err != gorm.ErrRecordNotFound {
				// 				return err
				// 			}
				// 		}
				// 		if lastInv.After(spisanie.Time) && lastInv.Before(time.Now()) {
				// 			flag = false
				// 		}
				// 		if flag {
				// 			skladIngredient := &model.SkladIngredient{}
				// 			if err := tx.Where("sklad_id = ? and ingredient_id = ?", spisanie.SkladID, ingredient.IngredientID).First(skladIngredient).Error; err != nil {
				// 				return err
				// 			}
				// 			skladIngredient.Quantity += ingredient.Brutto * item.Quantity
				// 			if err := tx.Model(skladIngredient).Select("quantity").Where("sklad_id = ? and ingredient_id = ?", spisanie.SkladID, ingredient.IngredientID).Updates(skladIngredient).Error; err != nil {
				// 				return err
				// 			}
				// 		}
				// 		if spisanie.Time.Year() == time.Now().Year() && spisanie.Time.Month() == time.Now().Month() && spisanie.Time.Day() != time.Now().Day() {
				// 			trafficItems = append(trafficItems, &model.AsyncJob{
				// 				ItemID:    ingredient.IngredientID,
				// 				ItemType:  utils.TypeIngredient,
				// 				SkladID:   spisanie.SkladID,
				// 				TimeStamp: spisanie.Time.UTC(),
				// 				CreatedAt: time.Now(),
				// 				Status:    utils.StatusNeedRecalculate,
				// 			})
				// 		}
				// 		err = tx.Model(invItem).Where("sklad_id = ? and item_id = ? and type = ? and time > ? and before_time < ?", spisanie.SkladID, item.ItemID, item.Type, spisanie.Time, spisanie.Time).Scan(invItem).Error
				// 		if err != nil {
				// 			if err.Error() != "record not found" {
				// 				return err
				// 			}
				// 		}
				// 		if invItem.ID != 0 && invItem.ItemID != 0 {
				// 			invItems = append(invItems, invItem)
				// 		}
				// 	}
			}

		}

		spisanie.Cost = 0

		for _, item := range spisanie.Items {
			flag := true
			invItem := &model.InventarizationItem{}
			if item.Type == utils.TypeIngredient {
				ingredient := &model.Ingredient{}
				err := tx.Model(&model.Ingredient{}).Select("ingredients.ingredient_id, ingredients.name, ingredients.measure, sklad_ingredients.cost").Joins("inner join sklad_ingredients on sklad_ingredients.ingredient_id = ingredients.ingredient_id").Where("ingredients.ingredient_id = ? and sklad_ingredients.sklad_id = ?", item.ItemID, spisanie.SkladID).Scan(&ingredient).Error
				if err != nil {
					return err
				}
				skladIngredient := &model.SkladIngredient{}
				if err := tx.Where("sklad_id = ? and ingredient_id = ?", spisanie.SkladID, item.ItemID).First(skladIngredient).Error; err != nil {
					return err
				}
				spisanie.Cost += skladIngredient.Cost * item.Quantity
				item.Cost = skladIngredient.Cost * item.Quantity

				var lastInv time.Time
				err = tx.Model(&model.InventarizationItem{}).Select("time").Where("sklad_id = ? and item_id = ? and type = ? and status != ?", spisanie.SkladID, item.ItemID, utils.TypeIngredient, utils.StatusNew).Order("time desc").First(&lastInv).Error
				if err != nil {
					if err != gorm.ErrRecordNotFound {
						return err
					}
				}
				if lastInv.After(spisanie.Time) && lastInv.Before(time.Now()) {
					flag = false
				}
				if flag {
					skladIngredient := &model.SkladIngredient{}
					if err := tx.Where("sklad_id = ? and ingredient_id = ?", spisanie.SkladID, item.ItemID).First(skladIngredient).Error; err != nil {
						return err
					}
					skladIngredient.Quantity -= item.Quantity

					if err := tx.Model(skladIngredient).Select("quantity").Where("sklad_id = ? and ingredient_id = ?", spisanie.SkladID, item.ItemID).Updates(skladIngredient).Error; err != nil {
						return err
					}
				}
				if spisanie.Time.Year() != time.Now().Year() || spisanie.Time.Month() != time.Now().Month() || spisanie.Time.Day() != time.Now().Day() {
					trafficItems = append(trafficItems, &model.AsyncJob{
						ItemID:    item.ItemID,
						ItemType:  item.Type,
						SkladID:   spisanie.SkladID,
						TimeStamp: spisanie.Time.UTC(),
						CreatedAt: time.Now(),
						Status:    utils.StatusNeedRecalculate,
					})
				}
				err = tx.Model(invItem).Where("sklad_id = ? and item_id = ? and type = ? and time > ? and before_time < ?", spisanie.SkladID, item.ItemID, item.Type, spisanie.Time, spisanie.Time).Scan(invItem).Error
				if err != nil {
					if err.Error() != "record not found" {
						return err
					}
				}
				if _, exists := invItemMap[invItem.ID]; !exists && invItem.ID != 0 && invItem.ItemID != 0 {
					invItemMap[invItem.ID] = invItem.ID
					invItems = append(invItems, invItem)
				}
			} else if item.Type == utils.TypeTovar {
				tovar := &model.Tovar{}
				err := tx.Model(&model.Tovar{}).Select("tovars.tovar_id, tovars.name, tovars.measure, sklad_tovars.cost").Joins("inner join sklad_tovars on sklad_tovars.tovar_id = tovars.tovar_id").Where("tovars.tovar_id = ? and sklad_tovars.sklad_id = ?", item.ItemID, spisanie.SkladID).Scan(&tovar).Error
				if err != nil {
					return err
				}
				skladTovar := &model.SkladTovar{}
				if err := tx.Where("sklad_id = ? and tovar_id = ?", spisanie.SkladID, item.ItemID).First(skladTovar).Error; err != nil {
					return err
				}
				spisanie.Cost += skladTovar.Cost * item.Quantity
				item.Cost = skladTovar.Cost * item.Quantity
				var lastInv time.Time
				err = tx.Model(&model.InventarizationItem{}).Select("time").Where("sklad_id = ? and item_id = ? and type = ? and status != ?", spisanie.SkladID, item.ItemID, utils.TypeTovar, utils.StatusNew).Order("time desc").First(&lastInv).Error
				if err != nil {
					if err != gorm.ErrRecordNotFound {
						return err
					}
				}
				if lastInv.After(spisanie.Time) && lastInv.Before(time.Now()) {
					flag = false
				}
				if flag {
					skladTovar := &model.SkladTovar{}
					if err := tx.Where("sklad_id = ? and tovar_id = ?", spisanie.SkladID, item.ItemID).First(skladTovar).Error; err != nil {
						return err
					}
					skladTovar.Quantity -= item.Quantity
					if err := tx.Model(skladTovar).Select("quantity").Where("sklad_id = ? and tovar_id = ?", spisanie.SkladID, item.ItemID).Updates(skladTovar).Error; err != nil {
						return err
					}
				}
				if spisanie.Time.Year() != time.Now().Year() || spisanie.Time.Month() != time.Now().Month() || spisanie.Time.Day() != time.Now().Day() {
					trafficItems = append(trafficItems, &model.AsyncJob{
						ItemID:    item.ItemID,
						ItemType:  item.Type,
						SkladID:   spisanie.SkladID,
						TimeStamp: spisanie.Time.UTC(),
						CreatedAt: time.Now(),
						Status:    utils.StatusNeedRecalculate,
					})
				}
				err = tx.Model(invItem).Where("sklad_id = ? and item_id = ? and type = ? and time > ? and before_time < ?", spisanie.SkladID, item.ItemID, item.Type, spisanie.Time, spisanie.Time).Scan(invItem).Error
				if err != nil {
					if err.Error() != "record not found" {
						return err
					}
				}
				if _, exists := invItemMap[invItem.ID]; !exists && invItem.ID != 0 && invItem.ItemID != 0 {
					invItemMap[invItem.ID] = invItem.ID
					invItems = append(invItems, invItem)
				}
			} else if item.Type == utils.TypeTechCart {
				techCart := &model.TechCart{}
				if err := tx.Model(&model.TechCart{}).Where("tech_cart_id = ? and shop_id = ?", item.ItemID, shopID).First(techCart).Error; err != nil {
					return err
				}
				ingredientsTech := []*model.IngredientTechCart{}
				if err := tx.Model(&model.IngredientTechCart{}).Where("tech_cart_id = ? and shop_id = ?", item.ItemID, shopID).Scan(&ingredientsTech).Error; err != nil {
					return err
				}
				for _, ingredient := range ingredientsTech {
					removeItem := &model.RemoveFromSkladItem{
						Type:           utils.TypeIngredient,
						Quantity:       ingredient.Brutto * item.Quantity,
						Cost:           0,
						Details:        item.Details,
						SkladID:        item.SkladID,
						ItemID:         ingredient.IngredientID,
						PartOfTechCart: true,
					}
					flag = true
					ingredientNew := &model.Ingredient{}
					err := tx.Model(&model.Ingredient{}).Select("ingredients.ingredient_id, ingredients.name, ingredients.measure, sklad_ingredients.cost").Joins("inner join sklad_ingredients on sklad_ingredients.ingredient_id = ingredients.ingredient_id").Where("ingredients.ingredient_id = ? and sklad_ingredients.sklad_id = ?", ingredient.IngredientID, spisanie.SkladID).Scan(&ingredientNew).Error
					if err != nil {
						return err
					}
					skladIngredient := &model.SkladIngredient{}
					if err := tx.Where("sklad_id = ? and ingredient_id = ?", spisanie.SkladID, ingredient.IngredientID).First(skladIngredient).Error; err != nil {
						return err
					}
					item.Cost += skladIngredient.Cost * ingredient.Brutto * item.Quantity
					removeItem.Cost = skladIngredient.Cost * ingredient.Brutto * item.Quantity
					spisanie.Items = append(spisanie.Items, removeItem)
					spisanie.Cost += skladIngredient.Cost * ingredient.Brutto * item.Quantity

					var lastInv time.Time
					err = tx.Model(&model.InventarizationItem{}).Select("time").Where("sklad_id = ? and item_id = ? and type = ? and status != ?", spisanie.SkladID, item.ItemID, utils.TypeIngredient, utils.StatusNew).Order("time desc").First(&lastInv).Error
					if err != nil {
						if err != gorm.ErrRecordNotFound {
							return err
						}
					}
					if lastInv.After(spisanie.Time) && lastInv.Before(time.Now()) {
						flag = false
					}
					if flag {
						skladIngredient := &model.SkladIngredient{}
						if err := tx.Where("sklad_id = ? and ingredient_id = ?", spisanie.SkladID, ingredient.IngredientID).First(skladIngredient).Error; err != nil {
							return err
						}
						skladIngredient.Quantity -= ingredient.Brutto * item.Quantity
						if err := tx.Model(skladIngredient).Select("quantity").Where("sklad_id = ? and ingredient_id = ?", spisanie.SkladID, ingredient.IngredientID).Updates(skladIngredient).Error; err != nil {
							return err
						}
					}

					if spisanie.Time.Year() != time.Now().Year() || spisanie.Time.Month() != time.Now().Month() || spisanie.Time.Day() != time.Now().Day() {
						trafficItems = append(trafficItems, &model.AsyncJob{
							ItemID:    ingredient.IngredientID,
							ItemType:  utils.TypeIngredient,
							SkladID:   spisanie.SkladID,
							TimeStamp: spisanie.Time.UTC(),
							CreatedAt: time.Now(),
							Status:    utils.StatusNeedRecalculate,
						})
					}
					err = tx.Model(invItem).Where("sklad_id = ? and item_id = ? and type = ? and time > ? and before_time < ?", spisanie.SkladID, ingredient.IngredientID, utils.TypeIngredient, spisanie.Time, spisanie.Time).Scan(invItem).Error
					if err != nil {
						if err.Error() != "record not found" {
							return err
						}
					}
					if invItem.ID != 0 && invItem.ItemID != 0 {
						invItems = append(invItems, invItem)
					}
				}
			}

		}
		spisanie.Status = utils.StatusClosed
		if err := tx.Model(spisanie).Debug().Save(spisanie).Error; err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return err
	}
	r.ConcurrentRecalculationForDailyStatistics(trafficItems)
	if len(invItems) > 0 {
		err = r.RecalculateInventarization(invItems)
		if err != nil {
			return err
		}
	}
	return nil
}

func (r *SkladDB) DeleteSpisanie(id int) error {
	invItems := []*model.InventarizationItem{}
	trafficItems := []*model.AsyncJob{}
	err := r.gormDB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&model.RemoveFromSklad{}).Where("id = ?", id).Update("deleted", true).Error; err != nil {
			return err
		}
		spisanie := model.RemoveFromSklad{}
		err := tx.Where("id = ?", id).First(&spisanie).Error
		if err != nil {
			return err
		}
		var shopID int
		if err := r.gormDB.Model(&model.Sklad{}).Select("shop_id").Where("id = ?", spisanie.SkladID).Scan(&shopID).Error; err != nil {
			return err
		}
		items := []*model.RemoveFromSkladItem{}

		err = tx.Model(items).Where("remove_id = ?", id).Scan(&items).Error
		if err != nil {
			return err
		}

		for _, item := range items {
			flag := true
			err := tx.Delete(item).Error
			if err != nil {
				return err
			}
			invItem := &model.InventarizationItem{}
			if item.Type == utils.TypeIngredient {
				var lastInv time.Time
				err := tx.Model(&model.InventarizationItem{}).Select("time").Where("sklad_id = ? and item_id = ? and type = ? and status != ?", spisanie.SkladID, item.ItemID, utils.TypeIngredient, utils.StatusNew).Order("time desc").First(&lastInv).Error
				if err != nil {
					if err != gorm.ErrRecordNotFound {
						return err
					}
				}
				if lastInv.After(spisanie.Time) && lastInv.Before(time.Now()) {
					flag = false
				}
				if flag {
					skladIngredient := &model.SkladIngredient{}
					if err := tx.Where("sklad_id = ? and ingredient_id = ?", spisanie.SkladID, item.ItemID).First(skladIngredient).Error; err != nil {
						return err
					}
					skladIngredient.Quantity += item.Quantity
					if err := tx.Model(skladIngredient).Select("quantity").Where("sklad_id = ? and ingredient_id = ?", spisanie.SkladID, item.ItemID).Save(skladIngredient).Error; err != nil {
						return err
					}
				}
				if spisanie.Time.Year() != time.Now().Year() || spisanie.Time.Month() != time.Now().Month() || spisanie.Time.Day() != time.Now().Day() {
					trafficItems = append(trafficItems, &model.AsyncJob{
						ItemID:    item.ItemID,
						ItemType:  item.Type,
						SkladID:   spisanie.SkladID,
						TimeStamp: spisanie.Time.UTC(),
						CreatedAt: time.Now(),
						Status:    utils.StatusNeedRecalculate,
					})
				}
				err = tx.Model(invItem).Where("sklad_id = ? and item_id = ? and type = ? and time > ? and before_time < ?", spisanie.SkladID, item.ItemID, item.Type, spisanie.Time, spisanie.Time).Scan(invItem).Error
				if err != nil {
					if err.Error() != "record not found" {
						return err
					}
				}
				if invItem.ID != 0 && invItem.ItemID != 0 {
					invItems = append(invItems, invItem)
				}
			} else if item.Type == utils.TypeTovar {
				var lastInv time.Time
				err := tx.Model(&model.InventarizationItem{}).Select("time").Where("sklad_id = ? and item_id = ? and type = ? and status != ?", spisanie.SkladID, item.ItemID, utils.TypeTovar, utils.StatusNew).Order("time desc").First(&lastInv).Error
				if err != nil {
					if err != gorm.ErrRecordNotFound {
						return err
					}
				}
				if lastInv.After(spisanie.Time) && lastInv.Before(time.Now()) {
					flag = false
				}
				if flag {
					skladTovar := &model.SkladTovar{}
					if err := tx.Where("sklad_id = ? and tovar_id = ?", spisanie.SkladID, item.ItemID).First(skladTovar).Error; err != nil {
						return err
					}
					skladTovar.Quantity += item.Quantity

					if err := tx.Model(skladTovar).Select("quantity").Where("sklad_id = ? and tovar_id = ?", spisanie.SkladID, item.ItemID).Save(skladTovar).Error; err != nil {
						return err
					}
				}
				if spisanie.Time.Year() != time.Now().Year() || spisanie.Time.Month() != time.Now().Month() || spisanie.Time.Day() != time.Now().Day() {
					trafficItems = append(trafficItems, &model.AsyncJob{
						ItemID:    item.ItemID,
						ItemType:  item.Type,
						SkladID:   spisanie.SkladID,
						TimeStamp: spisanie.Time.UTC(),
						CreatedAt: time.Now(),
						Status:    utils.StatusNeedRecalculate,
					})
				}
				err = tx.Model(invItem).Where("sklad_id = ? and item_id = ? and type = ? and time > ? and before_time < ?", spisanie.SkladID, item.ItemID, item.Type, spisanie.Time, spisanie.Time).Scan(invItem).Error
				if err != nil {
					if err.Error() != "record not found" {
						return err
					}
				}
				if invItem.ID != 0 && invItem.ItemID != 0 {
					invItems = append(invItems, invItem)
				}
				// } else if item.Type == utils.TypeTechCart {
				// 	techCart := &model.TechCart{}
				// 	if err := tx.Model(&model.TechCart{}).Where("tech_cart_id = ? and shop_id = ?", item.ItemID, shopID).First(techCart).Error; err != nil {
				// 		return err
				// 	}
				// 	ingredientsTech := []*model.IngredientTechCart{}
				// 	if err := tx.Model(&model.IngredientTechCart{}).Where("tech_cart_id = ? and shop_id = ?", item.ItemID, shopID).Scan(&ingredientsTech).Error; err != nil {
				// 		return err
				// 	}
				// 	for _, ingredient := range ingredientsTech {
				// 		flag = true
				// 		var lastInv time.Time
				// 		err := tx.Model(&model.InventarizationItem{}).Select("time").Where("sklad_id = ? and item_id = ? and type = ? and status != ?", spisanie.SkladID, item.ItemID, utils.TypeIngredient, utils.StatusNew).Order("time desc").First(&lastInv).Error
				// 		if err != nil {
				// 			if err != gorm.ErrRecordNotFound {
				// 				return err
				// 			}
				// 		}
				// 		if lastInv.After(spisanie.Time) && lastInv.Before(time.Now()) {
				// 			flag = false
				// 		}
				// 		if flag {
				// 			skladIngredient := &model.SkladIngredient{}
				// 			if err := tx.Where("sklad_id = ? and ingredient_id = ?", spisanie.SkladID, ingredient.IngredientID).First(skladIngredient).Error; err != nil {
				// 				return err
				// 			}
				// 			skladIngredient.Quantity += ingredient.Brutto * item.Quantity
				// 			if err := tx.Model(skladIngredient).Select("quantity").Where("sklad_id = ? and ingredient_id = ?", spisanie.SkladID, ingredient.IngredientID).Updates(skladIngredient).Error; err != nil {
				// 				return err
				// 			}
				// 		}
				// 		if spisanie.Time.Year() == time.Now().Year() && spisanie.Time.Month() == time.Now().Month() && spisanie.Time.Day() != time.Now().Day() {
				// 			trafficItems = append(trafficItems, &model.AsyncJob{
				// 				ItemID:    ingredient.IngredientID,
				// 				ItemType:  utils.TypeIngredient,
				// 				SkladID:   spisanie.SkladID,
				// 				TimeStamp: spisanie.Time.UTC(),
				// 				CreatedAt: time.Now(),
				// 				Status:    utils.StatusNeedRecalculate,
				// 			})
				// 		}
				// 		err = tx.Model(invItem).Where("sklad_id = ? and item_id = ? and type = ? and time > ? and before_time < ?", spisanie.SkladID, item.ItemID, item.Type, spisanie.Time, spisanie.Time).Scan(invItem).Error
				// 		if err != nil {
				// 			if err.Error() != "record not found" {
				// 				return err
				// 			}
				// 		}
				// 		if invItem.ID != 0 && invItem.ItemID != 0 {
				// 			invItems = append(invItems, invItem)
				// 		}
				// 	}
			}

		}
		return nil
	})
	if err != nil {
		return err
	}
	r.ConcurrentRecalculationForDailyStatistics(trafficItems)
	if len(invItems) > 0 {
		err = r.RecalculateInventarization(invItems)
		if err != nil {
			return err
		}
	}
	return nil
}

func (r *SkladDB) GetRemoved(filter *model.Filter) ([]*model.RemoveFromSkladResponse, int64, error) {
	res := r.gormDB
	if filter.Category != 0 {
		res = r.gormDB.Model(&model.RemoveFromSklad{}).Select("remove_from_sklads.id, json_agg(json_build_object('sklad', sklads.name,'worker_id', remove_from_sklads.worker_id,'reason', remove_from_sklads.reason, 'comment', remove_from_sklads.comment,'cost', remove_from_sklads.cost, 'time', remove_from_sklads.time,'status', remove_from_sklads.status, 'item_id', remove_from_sklad_items.item_id,'type', remove_from_sklad_items.type, 'quantity', remove_from_sklad_items.quantity,'item_cost', remove_from_sklad_items.cost,'details', remove_from_sklad_items.details,'name', CASE WHEN remove_from_sklad_items.type = 'ingredient' THEN ingredients.name When remove_from_sklad_items.type = 'techCart' Then tech_carts.name Else tovars.name END, 'measure', CASE WHEN remove_from_sklad_items.type = 'ingredient' THEN ingredients.measure When remove_from_sklad_items.type = 'techCart' Then tech_carts.measure Else tovars.measure END)) as remove_from_sklads_info").Joins("inner join remove_from_sklad_items on remove_from_sklad_items.remove_id = remove_from_sklads.id inner join sklads on remove_from_sklads.sklad_id = sklads.id left join ingredients on (ingredients.ingredient_id = remove_from_sklad_items.item_id and remove_from_sklad_items.type = 'ingredient') left join tovars on (tovars.tovar_id = remove_from_sklad_items.item_id and  remove_from_sklad_items.type = 'tovar') left join tech_carts on (tech_carts.tech_cart_id = remove_from_sklad_items.item_id and remove_from_sklad_items.type = 'techCart')").Where("remove_from_sklads.type != ? and remove_from_sklads.deleted = ? and sklads.shop_id IN (?) and (ingredients.shop_id = sklads.shop_id or tovars.shop_id = sklads.shop_id or tech_carts.shop_id = sklads.shop_id) and (ingredients.category = ? or tovars.category = ? or tech_carts.category = ?) and remove_from_sklad_items.part_of_tech_cart != ?", utils.TypeTransfer, false, filter.AccessibleShops, filter.Category, filter.Category, filter.Category, true).Group("remove_from_sklads.id")
		if res.Error != nil {
			return nil, 0, res.Error
		}
	} else {
		res = r.gormDB.Model(&model.RemoveFromSklad{}).Select("remove_from_sklads.id, json_agg(json_build_object('sklad', sklads.name,'worker_id', remove_from_sklads.worker_id,'reason', remove_from_sklads.reason, 'comment', remove_from_sklads.comment,'cost', remove_from_sklads.cost, 'time', remove_from_sklads.time,'status', remove_from_sklads.status, 'item_id', remove_from_sklad_items.item_id,'type', remove_from_sklad_items.type, 'quantity', remove_from_sklad_items.quantity,'item_cost', remove_from_sklad_items.cost,'details', remove_from_sklad_items.details,'name', CASE WHEN remove_from_sklad_items.type = 'ingredient' THEN ingredients.name When remove_from_sklad_items.type = 'techCart' Then tech_carts.name Else tovars.name END, 'measure', CASE WHEN remove_from_sklad_items.type = 'ingredient' THEN ingredients.measure When remove_from_sklad_items.type = 'techCart' Then tech_carts.measure Else tovars.measure END)) as remove_from_sklads_info").Joins("inner join remove_from_sklad_items on remove_from_sklad_items.remove_id = remove_from_sklads.id inner join sklads on remove_from_sklads.sklad_id = sklads.id left join ingredients on (ingredients.ingredient_id = remove_from_sklad_items.item_id and remove_from_sklad_items.type = 'ingredient') left join tovars on (tovars.tovar_id = remove_from_sklad_items.item_id and  remove_from_sklad_items.type = 'tovar') left join tech_carts on (tech_carts.tech_cart_id = remove_from_sklad_items.item_id and remove_from_sklad_items.type = 'techCart')").Where("remove_from_sklads.type != ? and remove_from_sklads.deleted = ? and sklads.shop_id IN (?) and (ingredients.shop_id = sklads.shop_id or tovars.shop_id = sklads.shop_id or tech_carts.shop_id = sklads.shop_id) and remove_from_sklad_items.part_of_tech_cart != ?", utils.TypeTransfer, false, filter.AccessibleShops, true).Group("remove_from_sklads.id")
		if res.Error != nil {
			return nil, 0, res.Error
		}
	}

	newRes, count, err := filter.FilterResults(res, model.RemoveFromSklad{}, utils.DefaultPageSize, "time", fmt.Sprintf("sklads.name ilike '%%%s%%' OR tovars.name ilike '%%%s%%' OR ingredients.name ilike '%%%s%%' OR tech_carts.name ilike '%%%s%%'", filter.Search, filter.Search, filter.Search, filter.Search), "remove_from_sklads.time desc")
	if err != nil {
		return nil, 0, err
	}

	var result model.ItemsRemovedFromSkladArr
	rows, err := newRes.Rows()
	if err != nil {
		return nil, 0, err
	}

	spisanieNeOuput := []*model.RemoveFromSkladOutput{}

	for rows.Next() {
		type IDStruct struct {
			ID int64 `json:"id"`
		}
		var id IDStruct
		err = rows.Scan(&id.ID, &result)
		if err != nil {
			return nil, 0, err
		}
		spisanieNeOuput = append(spisanieNeOuput, &model.RemoveFromSkladOutput{
			ID:   int(id.ID),
			Info: result,
		})
	}

	removed := []*model.RemoveFromSkladResponse{}
	for _, remove := range spisanieNeOuput {
		itemsRemoved := []*model.RemoveFromSkladItemResponse{}
		for _, item := range remove.Info {
			for _, item1 := range item.Items {
				itemRemoved := &model.RemoveFromSkladItemResponse{
					ItemID:   item1.ItemID,
					Name:     item1.Name,
					Measure:  item1.Measure,
					Type:     item1.Type,
					Quantity: item1.Quantity,
					Cost:     item1.ItemCost,
					Details:  item1.Details,
				}
				itemsRemoved = append(itemsRemoved, itemRemoved)
			}

		}
		removed = append(removed, &model.RemoveFromSkladResponse{
			ID:       remove.ID,
			Sklad:    remove.Info[0].Items[0].Sklad,
			WorkerID: remove.Info[0].Items[0].WorkerID,
			Reason:   remove.Info[0].Items[0].Reason,
			Comment:  remove.Info[0].Items[0].Comment,
			Cost:     remove.Info[0].Items[0].Cost,
			Time:     remove.Info[0].Items[0].Time,
			Status:   remove.Info[0].Items[0].Status,
			Items:    itemsRemoved,
		})
	}
	return removed, count, nil
}

func (r *SkladDB) GetRemovedByID(id int) (*model.RemoveFromSkladResponse, error) {
	removed := &model.RemoveFromSkladResponse{}

	res := r.gormDB.Model(&model.RemoveFromSklad{}).Select("remove_from_sklads.id, sklads.name as sklad, remove_from_sklads.sklad_id,remove_from_sklads.worker_id, remove_from_sklads.reason, remove_from_sklads.comment, remove_from_sklads.cost, remove_from_sklads.time, remove_from_sklads.status").Joins("inner join sklads on remove_from_sklads.sklad_id = sklads.id").Where("remove_from_sklads.deleted = ? and remove_from_sklads.id = ?", false, id).Scan(&removed)
	if res.Error != nil {
		return nil, res.Error
	}
	items := []*model.RemoveFromSkladItemResponse{}
	res = r.gormDB.Model(&model.RemoveFromSkladItem{}).Select("remove_from_sklad_items.id, remove_from_sklad_items.item_id, remove_from_sklad_items.quantity, remove_from_sklad_items.cost, remove_from_sklad_items.type, remove_from_sklad_items.details").Where("remove_from_sklad_items.remove_id = ? and remove_from_sklad_items.part_of_tech_cart != ?", id, true).Scan(&items)
	if res.Error != nil {
		return nil, res.Error
	}

	for _, item := range items {
		if item.Type == utils.TypeIngredient {
			ingredient := &model.Ingredient{}
			res := r.gormDB.Model(&model.Ingredient{}).Select("ingredients.ingredient_id, ingredients.name, ingredients.measure, sklad_ingredients.cost").Joins("inner join sklad_ingredients on sklad_ingredients.ingredient_id = ingredients.ingredient_id inner join sklads on sklads.id = sklad_ingredients.sklad_id").Where("ingredients.ingredient_id = ? and sklads.id = ? and sklads.shop_id = ingredients.shop_id", item.ItemID, removed.SkladID).Scan(&ingredient).Error
			if res != nil {
				return nil, res
			}
			// skladIngredient := &model.SkladIngredient{}
			// res = r.gormDB.Model(&model.SkladIngredient{}).Where("sklad_id = ? and ingredient_id = ?", removed.SkladID, item.ItemID).Scan(&skladIngredient).Error
			// if res != nil {
			// 	return nil, res
			// }
			item.Name = ingredient.Name
			item.Measure = ingredient.Measure
			// item.Cost = item.Cost * item.Quantity
			removed.Cost = removed.Cost + item.Cost

		}
		if item.Type == utils.TypeTovar {
			tovar := &model.Tovar{}
			err := r.gormDB.Model(&model.Tovar{}).Select("tovars.tovar_id, tovars.name, tovars.measure, sklad_tovars.cost").Joins("inner join sklad_tovars on sklad_tovars.tovar_id = tovars.tovar_id inner join sklads on sklads.id = sklad_tovars.sklad_id").Where("tovars.tovar_id = ? and sklads.id = ? and sklads.shop_id = tovars.shop_id", item.ItemID, removed.SkladID).Scan(&tovar).Error
			if err != nil {
				return nil, err
			}
			// skladTovar := &model.SkladTovar{}
			// err = r.gormDB.Model(&model.SkladTovar{}).Where("sklad_id = ? and tovar_id = ?", removed.SkladID, item.ItemID).Scan(&skladTovar).Error
			// if err != nil {
			// 	return nil, err
			// }
			item.Name = tovar.Name
			item.Measure = tovar.Measure
			// item.Cost = item.Cost * item.Quantity
			removed.Cost = removed.Cost + item.Cost

		}
		if item.Type == utils.TypeTechCart {
			techCart := &model.TechCart{}
			err := r.gormDB.Model(&model.TechCart{}).Select("tech_carts.tech_cart_id, tech_carts.shop_id, tech_carts.name, tech_carts.measure").Joins("inner join sklads on sklads.shop_id = tech_carts.shop_id").Where("tech_carts.tech_cart_id = ? and sklads.id = ? and sklads.shop_id = tech_carts.shop_id", item.ItemID, removed.SkladID).Scan(&techCart).Error
			if err != nil {
				return nil, err
			}

			// techIngredients := []*model.IngredientTechCart{}
			// res := r.gormDB.Model(&model.IngredientTechCart{}).Select("*").Where("ingredient_tech_carts.tech_cart_id = ? and shop_id = ?", techCart.TechCartID, techCart.ShopID).Scan(&techIngredients)
			// if res.Error != nil {
			// 	return nil, res.Error
			// }
			// var ingredientSum float32 = 0
			// for _, ingredient := range techIngredients {
			// 	skladIngredient := &model.SkladIngredient{}
			// 	res := r.gormDB.Model(&model.SkladIngredient{}).Where("sklad_ingredients.ingredient_id = ? and sklad_ingredients.sklad_id = ?", ingredient.IngredientID, removed.SkladID).Scan(&skladIngredient)
			// 	if res.Error != nil {
			// 		return nil, res.Error
			// 	}
			// 	ingredientSum += skladIngredient.Cost * ingredient.Brutto * item.Quantity
			// }
			item.Name = techCart.Name
			item.Measure = techCart.Measure
			// item.Cost = item.Cost * item.Quantity
			removed.Cost = removed.Cost + item.Cost
		}
	}
	removed.Items = items

	return removed, nil
}

func (r *SkladDB) AddSklad(sklad *model.Sklad) error {
	/*ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	stmt := `insert into sklad (name, address)
	values ($1, $2)`
	err := r.db.QueryRowContext(ctx, stmt,
		sklad.Name,
		sklad.Address,
	)

	if err.Err() != nil {
		return err.Err()
	}

	return nil*/
	err := r.gormDB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(sklad).Error; err != nil {
			return err
		}
		ingredients := []*model.Ingredient{}
		if err := tx.Where("shop_id = ?", sklad.ShopID).Find(&ingredients).Error; err != nil {
			return err
		}
		skladIngredients := []*model.SkladIngredient{}

		for _, ingredient := range ingredients {
			skladIngredient := &model.SkladIngredient{
				SkladID:      sklad.ID,
				IngredientID: ingredient.IngredientID,
				Quantity:     0,
				Cost:         ingredient.Cost,
			}
			skladIngredients = append(skladIngredients, skladIngredient)
		}
		err := tx.Model(&model.SkladIngredient{}).Create(skladIngredients).Error
		if err != nil {
			return err
		}
		tovars := []*model.Tovar{}
		if err := tx.Where("shop_id = ?", sklad.ShopID).Find(&tovars).Error; err != nil {
			return err
		}
		skladTovars := []*model.SkladTovar{}
		for _, tovar := range tovars {
			skladTovar := &model.SkladTovar{
				SkladID:  sklad.ID,
				TovarID:  tovar.TovarID,
				Quantity: 0,
				Cost:     tovar.Cost,
			}
			skladTovars = append(skladTovars, skladTovar)
		}
		err = tx.Model(&model.SkladTovar{}).Create(skladTovars).Error
		if err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		return err
	}

	return nil
}

func (r *SkladDB) GetAllSklad(filter *model.Filter) ([]*model.Sklad, int64, error) {
	/*ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()
	sklads := []*model.Sklad{}
	stmt := `select * from sklad where deleted = $1`

	row, err := r.db.QueryContext(ctx, stmt, false)
	if err != nil {
		return nil, err
	}

	for row.Next() {
		sklad := &model.Sklad{}
		err := row.Scan(
			&sklad.ID,
			&sklad.Name,
			&sklad.Address,
			&sklad.Deleted,
		)
		if err != nil {
			return nil, err
		}
		sklads = append(sklads, sklad)
	}

	if len(sklads) == 0 {
		sklads = []*model.Sklad{}
	}

	return sklads, nil*/
	sklads := []*model.Sklad{}
	res := r.gormDB.Where("deleted = ? and shop_id IN (?)", false, filter.AccessibleShops).Find(&sklads)
	if res.Error != nil {
		return nil, 0, res.Error
	}

	newRes, count, err := filter.FilterResults(res, model.Sklad{}, utils.DefaultPageSize, "", "", "")

	if err != nil {
		return nil, 0, err
	}

	if newRes.Scan(&sklads).Error != nil {
		return nil, 0, newRes.Error
	}

	return sklads, count, nil
}

func (r *SkladDB) GetSklad(id int) (*model.Sklad, error) {
	/*ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()
	sklad := &model.Sklad{}
	stmt := `select * from sklad where id=$1`

	row := r.db.QueryRowContext(ctx, stmt, id)

	if row.Err() != nil {
		return nil, row.Err()
	}

	err := row.Scan(
		&sklad.ID,
		&sklad.Name,
		&sklad.Address,
		&sklad.Deleted,
	)
	if err != nil {
		return nil, err
	}

	return sklad, nil*/
	sklad := &model.Sklad{}
	err := r.gormDB.Where("id = ?", id).First(sklad).Error
	if err != nil {
		return nil, err
	}
	return sklad, nil
}
func (r *SkladDB) UpdateSklad(sklad *model.Sklad) error {
	/*ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()
	stmt := `update sklad set name = $2, address = $3 where id=$1`

	row := r.db.QueryRowContext(ctx, stmt,
		sklad.ID,
		sklad.Name,
		sklad.Address,
	)

	if row.Err() != nil {
		return row.Err()
	}

	return nil*/
	err := r.gormDB.Save(sklad).Error
	if err != nil {
		return err
	}
	return nil
}
func (r *SkladDB) DeleteSklad(id int) error {
	/*ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()
	stmt := `update sklad set deleted = $2 where id=$1`
	row := r.db.QueryRowContext(ctx, stmt, id, true)

	if row.Err() != nil {
		return row.Err()
	}

	return nil*/
	if err := r.gormDB.Table("sklads").Where("id = ?", id).Update("deleted", true).Error; err != nil {
		return err
	}
	return nil
}

func (r *SkladDB) Ostatki(filter *model.Filter) ([]*model.Item, int64, error) {
	/*ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()
	items := []*model.Item{}
	stmtIngredient := `select ingredients.id, sklads.id, ingredients.measure, ingredients.name, sklads.name, category_ingredients.name, sklad_ingredients.quantity, sklad_ingredients.cost from sklad_ingredients inner join sklads on sklad_ingredients.sklad_id = sklads.id inner join ingredients on sklad_ingredients.ingredient_id = ingredients.id inner join category_ingredients on ingredients.category = category_ingredients.id`
	row, err := r.db.QueryContext(ctx, stmtIngredient)
	if err != nil {
		return nil, err
	}

	for row.Next() {
		item := &model.Item{}
		err := row.Scan(
			&item.ItemID,
			&item.SkladID,
			&item.Measurement,
			&item.Name,
			&item.SkladName,
			&item.Category,
			&item.Quantity,
			&item.Cost,
		)
		if err != nil {
			return nil, err
		}
		item.Type = utils.TypeIngredient
		item.Sum = item.Quantity * item.Cost
		items = append(items, item)
	}
	stmtTovar := `select tovars.id, sklads.id, tovars.name, sklads.name, category_tovars.name, sklad_tovars.quantity, sklad_tovars.cost from sklad_tovars inner join sklads on sklad_tovars.sklad_id = sklads.id inner join tovars on sklad_tovars.tovar_id = tovars.id inner join category_tovars on tovars.category = category_tovars.id`

	rowTovar, err := r.db.QueryContext(ctx, stmtTovar)
	if err != nil {
		return nil, err
	}

	for rowTovar.Next() {
		item := &model.Item{}
		err := rowTovar.Scan(
			&item.ItemID,
			&item.SkladID,
			&item.Name,
			&item.SkladName,
			&item.Category,
			&item.Quantity,
			&item.Cost,
		)

		if err != nil {
			return nil, err
		}
		item.Type = utils.TypeTovar
		item.Measurement = "шт."
		item.Sum = item.Quantity * item.Cost
		items = append(items, item)
	}

	if len(items) == 0 {
		items = []*model.Item{}
	}

	return items, nil*/
	// pagination
	/*ostatki := []*model.Item{}

	res := r.gormDB.Model(&model.SkladIngredient{}).Select("ingredients.id as item_id, sklads.id as sklad_id, ingredients.measure as measurement, ingredients.name, sklads.name as sklad_name, category_ingredients.name as category, sklad_ingredients.quantity, sklad_ingredients.cost").Joins("inner join sklads on sklad_ingredients.sklad_id = sklads.id inner join ingredients on sklad_ingredients.ingredient_id = ingredients.id inner join category_ingredients on ingredients.category = category_ingredients.id").Where("ingredients.deleted = ?", false).Where("sklads.deleted = ?", false)

	if res.Error != nil {
		return nil, 0, res.Error
	}

	newRes, count, err := filter.FilterResults(res, &model.Ingredient{}, utils.DefaultPageSize, "")

	if err != nil {
		return nil, 0, err
	}

	if newRes.Scan(&ostatki).Error != nil {
		return nil, 0, newRes.Error
	}

	for _, item := range ostatki {
		item.Type = utils.TypeIngredient
		item.Sum = item.Quantity * item.Cost
	}

	tempOstatki := []*model.Item{}

	res1 := r.gormDB.Model(&model.SkladTovar{}).Select("tovars.id as item_id, sklads.id as sklad_id, tovars.name, sklads.name as sklad_name, category_tovars.name as category, sklad_tovars.quantity, sklad_tovars.cost").Joins("inner join sklads on sklad_tovars.sklad_id = sklads.id inner join tovars on sklad_tovars.tovar_id = tovars.id inner join category_tovars on tovars.category = category_tovars.id").Where("tovars.deleted = ?", false).Where("sklads.deleted = ?", false)

	if res1.Error != nil {
		return nil, 0, res1.Error
	}

	NewRes, Count, Err := filter.FilterResults(res1, &model.Tovar{}, utils.DefaultPageSize, "")

	if Err != nil {
		return nil, 0, Err
	}

	if NewRes.Scan(&tempOstatki).Error != nil {
		return nil, 0, NewRes.Error
	}

	for _, val := range tempOstatki {
		val.Type = utils.TypeTovar
		val.Measurement = "шт."
		val.Sum = val.Quantity * val.Cost
		ostatki = append(ostatki, val)
	}

	return ostatki, count + Count, nil*/

	ostatki := []*model.Item{}

	if filter.Type != "" {
		if filter.Type == utils.TypeIngredient {
			res := r.gormDB
			if filter.Category != 0 {
				res = r.gormDB.Model(&model.SkladIngredient{}).Select("ingredients.ingredient_id as item_id, sklads.id as sklad_id, ingredients.measure as measurement, ingredients.name, sklads.name as sklad_name, category_ingredients.name as category, sklad_ingredients.quantity, sklad_ingredients.cost").Joins("inner join sklads on sklad_ingredients.sklad_id = sklads.id inner join ingredients on sklad_ingredients.ingredient_id = ingredients.ingredient_id inner join category_ingredients on ingredients.category = category_ingredients.id").Where("ingredients.deleted = ? and sklads.shop_id IN (?) and ingredients.shop_id = sklads.shop_id and ingredients.category = ?", false, filter.AccessibleShops, filter.Category).Where("sklads.deleted = ?", false)
				if res.Error != nil {
					return nil, 0, res.Error
				}
			} else {
				res = r.gormDB.Model(&model.SkladIngredient{}).Select("ingredients.ingredient_id as item_id, sklads.id as sklad_id, ingredients.measure as measurement, ingredients.name, sklads.name as sklad_name, category_ingredients.name as category, sklad_ingredients.quantity, sklad_ingredients.cost").Joins("inner join sklads on sklad_ingredients.sklad_id = sklads.id inner join ingredients on sklad_ingredients.ingredient_id = ingredients.ingredient_id inner join category_ingredients on ingredients.category = category_ingredients.id").Where("ingredients.deleted = ? and sklads.shop_id IN (?) and ingredients.shop_id = sklads.shop_id", false, filter.AccessibleShops).Where("sklads.deleted = ?", false)
				if res.Error != nil {
					return nil, 0, res.Error
				}
			}

			newRes, count, err := filter.FilterResults(res, &model.Ingredient{}, utils.DefaultPageSize, "", "", "")

			if err != nil {
				return nil, 0, err
			}

			if newRes.Scan(&ostatki).Error != nil {
				return nil, 0, newRes.Error
			}

			for _, item := range ostatki {
				item.Type = utils.TypeIngredient
				item.Sum = item.Quantity * item.Cost
			}
			return ostatki, count, nil
		} else if filter.Type == utils.TypeTovar {
			res1 := r.gormDB
			if filter.Category != 0 {
				res1 = r.gormDB.Model(&model.SkladTovar{}).Select("tovars.tovar_id as item_id, sklads.id as sklad_id, tovars.name, sklads.name as sklad_name, category_tovars.name as category, sklad_tovars.quantity, sklad_tovars.cost").Joins("inner join sklads on sklad_tovars.sklad_id = sklads.id inner join tovars on sklad_tovars.tovar_id = tovars.tovar_id inner join category_tovars on tovars.category = category_tovars.id").Where("tovars.deleted = ? and sklads.shop_id IN (?) and tovars.shop_id = sklads.shop_id and tovars.category = ?", false, filter.AccessibleShops, filter.Category).Where("sklads.deleted = ?", false)
				if res1.Error != nil {
					return nil, 0, res1.Error
				}
			} else {
				res1 = r.gormDB.Model(&model.SkladTovar{}).Select("tovars.tovar_id as item_id, sklads.id as sklad_id, tovars.name, sklads.name as sklad_name, category_tovars.name as category, sklad_tovars.quantity, sklad_tovars.cost").Joins("inner join sklads on sklad_tovars.sklad_id = sklads.id inner join tovars on sklad_tovars.tovar_id = tovars.tovar_id inner join category_tovars on tovars.category = category_tovars.id").Where("tovars.deleted = ? and sklads.shop_id IN (?) and tovars.shop_id = sklads.shop_id", false, filter.AccessibleShops).Where("sklads.deleted = ?", false)
				if res1.Error != nil {
					return nil, 0, res1.Error
				}
			}

			NewRes, Count, Err := filter.FilterResults(res1, &model.Tovar{}, utils.DefaultPageSize, "", "", "")
			if Err != nil {
				return nil, 0, Err
			}

			if NewRes.Scan(&ostatki).Error != nil {
				return nil, 0, NewRes.Error
			}
			for _, val := range ostatki {
				val.Type = utils.TypeTovar
				val.Measurement = "шт."
				val.Sum = val.Quantity * val.Cost
			}
			return ostatki, Count, nil
		}
	}

	res := r.gormDB
	if filter.Category != 0 {
		res = r.gormDB.Model(&model.SkladIngredient{}).Select("ingredients.ingredient_id as item_id, sklads.id as sklad_id, ingredients.measure as measurement, ingredients.name, sklads.name as sklad_name, category_ingredients.name as category, sklad_ingredients.quantity, sklad_ingredients.cost").Joins("inner join sklads on sklad_ingredients.sklad_id = sklads.id inner join ingredients on sklad_ingredients.ingredient_id = ingredients.ingredient_id inner join category_ingredients on ingredients.category = category_ingredients.id").Where("ingredients.deleted = ? and sklads.shop_id IN (?) and ingredients.shop_id = sklads.shop_id and ingredients.category = ?", false, filter.AccessibleShops, filter.Category).Where("sklads.deleted = ?", false)
		if res.Error != nil {
			return nil, 0, res.Error
		}
	} else {
		res = r.gormDB.Model(&model.SkladIngredient{}).Select("ingredients.ingredient_id as item_id, sklads.id as sklad_id, ingredients.measure as measurement, ingredients.name, sklads.name as sklad_name, category_ingredients.name as category, sklad_ingredients.quantity, sklad_ingredients.cost").Joins("inner join sklads on sklad_ingredients.sklad_id = sklads.id inner join ingredients on sklad_ingredients.ingredient_id = ingredients.ingredient_id inner join category_ingredients on ingredients.category = category_ingredients.id").Where("ingredients.deleted = ? and sklads.shop_id IN (?) and ingredients.shop_id = sklads.shop_id", false, filter.AccessibleShops).Where("sklads.deleted = ?", false)
		if res.Error != nil {
			return nil, 0, res.Error
		}
	}

	newRes, count, err := filter.FilterResults(res, &model.Ingredient{}, utils.DefaultPageSize, "", "", "")

	if err != nil {
		return nil, 0, err
	}

	if newRes.Scan(&ostatki).Error != nil {
		return nil, 0, newRes.Error
	}

	for _, item := range ostatki {
		item.Type = utils.TypeIngredient
		item.Sum = item.Quantity * item.Cost
	}

	restOstatki := utils.DefaultPageSize - len(ostatki)
	var countTovar int64
	countPage := r.gormDB
	if filter.Category != 0 {
		countPage = r.gormDB.Model(&model.SkladTovar{}).Select("tovars.tovar_id as item_id, sklads.id as sklad_id, tovars.name, sklads.name as sklad_name, category_tovars.name as category, sklad_tovars.quantity, sklad_tovars.cost").Joins("inner join sklads on sklad_tovars.sklad_id = sklads.id inner join tovars on sklad_tovars.tovar_id = tovars.tovar_id inner join category_tovars on tovars.category = category_tovars.id").Where("tovars.deleted = ? and sklads.shop_id IN (?) and tovars.shop_id = sklads.shop_id and tovars.category = ?", false, filter.AccessibleShops, filter.Category).Where("sklads.deleted = ?", false)
		if countPage.Error != nil {
			return nil, 0, countPage.Error
		}
	} else {
		countPage = r.gormDB.Model(&model.SkladTovar{}).Select("tovars.tovar_id as item_id, sklads.id as sklad_id, tovars.name, sklads.name as sklad_name, category_tovars.name as category, sklad_tovars.quantity, sklad_tovars.cost").Joins("inner join sklads on sklad_tovars.sklad_id = sklads.id inner join tovars on sklad_tovars.tovar_id = tovars.tovar_id inner join category_tovars on tovars.category = category_tovars.id").Where("tovars.deleted = ? and sklads.shop_id IN (?) and tovars.shop_id = sklads.shop_id", false, filter.AccessibleShops).Where("sklads.deleted = ?", false)
		if countPage.Error != nil {
			return nil, 0, countPage.Error
		}
	}

	countPage.Count(&countTovar)
	if restOstatki > 0 {
		tempOstatki := []*model.Item{}

		res1 := r.gormDB
		if filter.Category != 0 {
			res1 = r.gormDB.Model(&model.SkladTovar{}).Select("tovars.tovar_id as item_id, sklads.id as sklad_id, tovars.name, sklads.name as sklad_name, category_tovars.name as category, sklad_tovars.quantity, sklad_tovars.cost").Joins("inner join sklads on sklad_tovars.sklad_id = sklads.id inner join tovars on sklad_tovars.tovar_id = tovars.tovar_id inner join category_tovars on tovars.category = category_tovars.id").Where("tovars.deleted = ? and sklads.shop_id IN (?) and tovars.shop_id = sklads.shop_id and tovars.category = ?", false, filter.AccessibleShops, filter.Category).Where("sklads.deleted = ?", false)
			if res1.Error != nil {
				return nil, 0, res1.Error
			}
		} else {
			res1 = r.gormDB.Model(&model.SkladTovar{}).Select("tovars.tovar_id as item_id, sklads.id as sklad_id, tovars.name, sklads.name as sklad_name, category_tovars.name as category, sklad_tovars.quantity, sklad_tovars.cost").Joins("inner join sklads on sklad_tovars.sklad_id = sklads.id inner join tovars on sklad_tovars.tovar_id = tovars.tovar_id inner join category_tovars on tovars.category = category_tovars.id").Where("tovars.deleted = ? and sklads.shop_id IN (?) and tovars.shop_id = sklads.shop_id", false, filter.AccessibleShops).Where("sklads.deleted = ?", false)
			if res1.Error != nil {
				return nil, 0, res1.Error
			}
		}

		ingredientPages := 0
		if count%utils.DefaultPageSize == 0 {
			ingredientPages = int(count / utils.DefaultPageSize)
		} else {
			ingredientPages = int(count/utils.DefaultPageSize) + 1
		}
		oldFilterPage := filter.Page
		if count%utils.DefaultPageSize == 0 {
			filter.Page = filter.Page - ingredientPages
		} else {
			filter.Page = filter.Page - ingredientPages + 1
		}
		NewRes, Count, Err := filter.FilterResults(res1, &model.Tovar{}, restOstatki, "", "", "")
		filter.Page = oldFilterPage

		if Err != nil {
			return nil, 0, Err
		}

		if NewRes.Scan(&tempOstatki).Error != nil {
			return nil, 0, NewRes.Error
		}
		if Count != countTovar {
			countTovar = Count
		}
		for _, val := range tempOstatki {
			val.Type = utils.TypeTovar
			val.Measurement = "шт."
			val.Sum = val.Quantity * val.Cost
			ostatki = append(ostatki, val)
		}

	}
	return ostatki, count + countTovar, nil
}

func (r *SkladDB) AddToSklad(postavka *model.Postavka, shopID int, id int) (*model.Postavka, int, error) {
	invItems := []*model.InventarizationItem{}
	shift := &model.Shift{}
	trafficItems := []*model.AsyncJob{}

	err := r.gormDB.Transaction(func(tx *gorm.DB) error {
		err := tx.Debug().Select("*").Create(postavka).Error
		if err != nil {
			return err
		}

		for _, item := range postavka.Items {
			item.PostavkaID = postavka.ID
			newCost := float32(0)
			invItem := &model.InventarizationItem{}
			if item.Type == utils.TypeIngredient {
				skladIngredient := &model.SkladIngredient{}
				err := tx.Model(skladIngredient).Where("sklad_id = ? and ingredient_id = ?", postavka.SkladID, item.ItemID).Scan(skladIngredient).Error
				if err != nil {
					return err
				}
				newCost = skladIngredient.Cost

				var lastInv time.Time
				err = tx.Model(&model.InventarizationItem{}).Select("time").Where("sklad_id = ? and item_id = ? and type = ? and status != ?", postavka.SkladID, item.ItemID, utils.TypeIngredient, utils.StatusNew).Order("time desc").First(&lastInv).Error
				if err != nil {
					if err != gorm.ErrRecordNotFound {
						return err
					}
				}
				flag := true
				if lastInv.After(postavka.Time) && lastInv.Before(time.Now()) { //?
					flag = false
				}
				if postavka.Type != utils.TypeTransfer {
					err = r.RecalculateCost(tx, &item, utils.Postavka, postavka.SkladID, flag)
					if err != nil {
						return err
					}
				} else {
					if flag {
						if err := tx.Where("sklad_id = ? and ingredient_id = ?", postavka.SkladID, item.ItemID).First(skladIngredient).Error; err != nil {
							return err
						}
						skladIngredient.Quantity += item.Quantity
						if err := tx.Select("quantity").Where("sklad_id = ? and ingredient_id = ?", postavka.SkladID, item.ItemID).Updates(skladIngredient).Error; err != nil {
							return err
						}
					}
				}
				skladIngredient = &model.SkladIngredient{}
				err = tx.Model(skladIngredient).Where("sklad_id = ? and ingredient_id = ?", postavka.SkladID, item.ItemID).Scan(skladIngredient).Error
				if err != nil {
					return err
				}
				newCost = skladIngredient.Cost
			} else if item.Type == utils.TypeTovar {
				skladTovar := &model.SkladTovar{}
				err := tx.Model(skladTovar).Where("sklad_id = ? and tovar_id = ?", postavka.SkladID, item.ItemID).Scan(skladTovar).Error
				if err != nil {
					if err != gorm.ErrRecordNotFound {
						return err
					}
				}
				newCost = skladTovar.Cost

				var lastInv time.Time
				err = tx.Model(&model.InventarizationItem{}).Select("time").Where("sklad_id = ? and item_id = ? and type = ? and status != ?", postavka.SkladID, item.ItemID, utils.TypeTovar, utils.StatusNew).Order("time desc").First(&lastInv).Error
				if err != nil {
					if err != gorm.ErrRecordNotFound {
						return err
					}
				}
				flag := true
				if lastInv.After(postavka.Time) && lastInv.Before(time.Now()) {
					flag = false
				}
				if postavka.Type != utils.TypeTransfer {
					err = r.RecalculateCost(tx, &item, utils.Postavka, postavka.SkladID, flag)
					if err != nil {
						return err
					}
				} else {
					if flag {
						if err := tx.Where("sklad_id = ? and tovar_id = ?", postavka.SkladID, item.ItemID).First(skladTovar).Error; err != nil {
							return err
						}
						skladTovar.Quantity += item.Quantity
						if err := tx.Select("quantity").Where("sklad_id = ? and tovar_id = ?", postavka.SkladID, item.ItemID).Updates(skladTovar).Error; err != nil {
							return err
						}
					}
				}

				skladTovar = &model.SkladTovar{}
				err = tx.Model(skladTovar).Where("sklad_id = ? and tovar_id = ?", postavka.SkladID, item.ItemID).Scan(skladTovar).Error
				if err != nil {
					return err
				}
				newCost = skladTovar.Cost
			} else {
				return errors.New("unknown type")
			}
			if postavka.Time.Year() != time.Now().Year() || postavka.Time.Month() != time.Now().Month() || postavka.Time.Day() != time.Now().Day() {
				trafficItems = append(trafficItems, &model.AsyncJob{
					ItemID:    item.ItemID,
					ItemType:  item.Type,
					SkladID:   postavka.SkladID,
					TimeStamp: postavka.Time.UTC(),
					CreatedAt: time.Now(),
					Status:    utils.StatusNeedRecalculate,
				})
			}
			err = tx.Model(invItem).Where("sklad_id = ? and item_id = ? and type = ? and time > ? and before_time < ?", postavka.SkladID, item.ItemID, item.Type, postavka.Time, postavka.Time).Scan(invItem).Error

			if err != nil {
				if err.Error() != "record not found" {
					return err
				}
			}

			if invItem.ID != 0 && invItem.ItemID != 0 {
				invItem.Cost = newCost
				invItems = append(invItems, invItem)
			}

			if err != nil {
				if err.Error() != "record not found" {
					return err
				}
			}

		}
		if postavka.Type != utils.TypeTransfer {
			transaction := &model.Transaction{
				WorkerID: id,
				SchetID:  postavka.SchetID,
				Category: utils.Postavka,
				Time:     postavka.Time,
				Sum:      postavka.Sum,
				Comment:  "Поставка №" + strconv.Itoa(postavka.ID),
			}
			err = r.gormDB.Model(&model.Shift{}).Where("created_at <= ? and shop_id = ?", transaction.Time, shopID).Order("created_at desc").First(shift).Error
			if err != nil {
				if err == gorm.ErrRecordNotFound {
					err = errors.New("could not find the shift for given shop")
				}
				return err
			}
			if transaction.Time.After(shift.ClosedAt) && shift.IsClosed {
				return errors.New("shoud be in range of shift")
			}

			transaction.ShiftID = shift.ID

			if transaction.SchetID == 0 {
				transaction.Status = utils.TransactionNeutralStatus
			} else {
				transaction.Status = utils.TransactionNegativeStatus
			}

			err = r.gormDB.Model(&model.Transaction{}).Debug().Create(transaction).Error
			if err != nil {
				return err
			}
			if transaction.Status == utils.TransactionNegativeStatus {
				res := tx.Table("schets").Where("id = ?", transaction.SchetID).Update("start_balance", gorm.Expr("start_balance - ?", transaction.Sum))
				if res.Error != nil {
					return res.Error
				}
			} else if transaction.Status == utils.TransactionPositiveStatus {
				res := tx.Table("schets").Where("id = ?", transaction.SchetID).Update("start_balance", gorm.Expr("start_balance + ?", transaction.Sum))
				if res.Error != nil {
					return res.Error
				}
			}
			transactionPostavka := &model.TransactionPostavka{
				TransactionID: transaction.ID,
				PostavkaID:    postavka.ID,
			}
			err = tx.Model(&model.TransactionPostavka{}).Create(transactionPostavka).Error
			if err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return nil, -1, err
	}
	r.ConcurrentRecalculationForDailyStatistics(trafficItems)

	if len(invItems) > 0 {
		r.ConcurrentRecalculation(invItems)
	}
	return postavka, shift.ID, err
}

func (r *SkladDB) GetItems(filter *model.Filter) ([]*model.ItemOutput, error) {
	// ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	// defer cancel()
	// items := []*model.ItemOutput{}
	// stmtIngredient := `select ingredients.id, ingredients.name, ingredients.cost, ingredients.measure from ingredients where deleted = false`

	// row, err := r.db.QueryContext(ctx, stmtIngredient)
	// if err != nil {
	// 	return nil, err
	// }

	// for row.Next() {
	// 	item := &model.ItemOutput{}
	// 	err := row.Scan(
	// 		&item.ID,
	// 		&item.Name,
	// 		&item.Cost,
	// 		&item.Measure,
	// 	)
	// 	if err != nil {
	// 		return nil, err
	// 	}
	// 	item.Type = utils.TypeIngredient
	// 	items = append(items, item)
	// }
	// stmtTovar := `select tovars.id, tovars.name, tovars.cost, tovars.measure from tovars where deleted = false`

	// rowTovar, err := r.db.QueryContext(ctx, stmtTovar)
	// if err != nil {
	// 	return nil, err
	// }

	// for rowTovar.Next() {
	// 	item := &model.ItemOutput{}
	// 	err := rowTovar.Scan(
	// 		&item.ID,
	// 		&item.Name,
	// 		&item.Cost,
	// 		&item.Measure,
	// 	)
	// 	if err != nil {
	// 		return nil, err
	// 	}
	// 	item.Type = utils.TypeTovar
	// 	items = append(items, item)
	// }

	// if len(items) == 0 {
	// 	items = []*model.ItemOutput{}
	// }

	// return items, nil
	items := []*model.ItemOutput{}

	if filter.Type == utils.TypeIngredient || filter.Type == "" {
		ingredients := []*model.ItemOutput{}
		if filter.Measure == "" {
			err := r.gormDB.Model(&model.Sklad{}).Select("ingredients.ingredient_id as id, ingredients.name, ingredients.measure, sklad_ingredients.cost, sklads.id as sklad_id, (SELECT MAX(postavkas.time) FROM postavkas inner join item_postavkas on item_postavkas.postavka_id = postavkas.id where item_postavkas.item_id = ingredients.ingredient_id and postavkas.sklad_id = sklads.id and item_postavkas.type = 'ingredient') as last_postavka_time, (SELECT item_postavkas.cost FROM postavkas inner join item_postavkas on item_postavkas.postavka_id = postavkas.id where item_postavkas.item_id = ingredients.ingredient_id and postavkas.sklad_id = sklads.id and item_postavkas.type = 'ingredient' order by time desc limit 1) as last_postavka_cost").Joins("inner join sklad_ingredients on sklad_ingredients.sklad_id = sklads.id inner join ingredients on ingredients.ingredient_id = sklad_ingredients.ingredient_id").Where("ingredients.deleted = ? and sklads.shop_id IN (?) and ingredients.shop_id = sklads.shop_id", false, filter.AccessibleShops).Scan(&ingredients).Error
			if err != nil {
				return nil, err
			}
		} else {
			err := r.gormDB.Model(&model.Sklad{}).Select("ingredients.ingredient_id as id, ingredients.name, ingredients.measure, sklad_ingredients.cost, sklads.id as sklad_id, (SELECT MAX(postavkas.time) FROM postavkas inner join item_postavkas on item_postavkas.postavka_id = postavkas.id where item_postavkas.item_id = ingredients.ingredient_id and postavkas.sklad_id = sklads.id and item_postavkas.type = 'ingredient') as last_postavka_time, (SELECT item_postavkas.cost FROM postavkas inner join item_postavkas on item_postavkas.postavka_id = postavkas.id where item_postavkas.item_id = ingredients.ingredient_id and postavkas.sklad_id = sklads.id and item_postavkas.type = 'ingredient' order by time desc limit 1) as last_postavka_cost").Joins("inner join sklad_ingredients on sklad_ingredients.sklad_id = sklads.id inner join ingredients on ingredients.ingredient_id = sklad_ingredients.ingredient_id").Where("ingredients.deleted = ? and sklads.shop_id IN (?) and ingredients.measure = ?", false, filter.AccessibleShops, filter.Measure).Scan(&ingredients).Error
			if err != nil {
				return nil, err
			}
		}

		for _, val := range ingredients {
			items = append(items, &model.ItemOutput{
				ID:               val.ID,
				Name:             val.Name,
				Measure:          val.Measure,
				Cost:             val.Cost,
				SkladID:          val.SkladID,
				LastPostavkaTime: val.LastPostavkaTime,
				LastPostavkaCost: val.LastPostavkaCost,
				Type:             utils.TypeIngredient,
			})
		}

	}
	if filter.Type == utils.TypeTovar || filter.Type == "" {
		tovars := []*model.ItemOutput{}
		err := r.gormDB.Model(&model.Sklad{}).Select("tovars.tovar_id as id, tovars.name, tovars.measure, sklad_tovars.cost, sklads.id as sklad_id, (SELECT MAX(postavkas.time) FROM postavkas inner join item_postavkas on item_postavkas.postavka_id = postavkas.id where item_postavkas.item_id = tovars.tovar_id and postavkas.sklad_id = sklads.id and item_postavkas.type = 'tovar') as last_postavka_time, (SELECT item_postavkas.cost FROM postavkas inner join item_postavkas on item_postavkas.postavka_id = postavkas.id where item_postavkas.item_id = tovars.tovar_id and postavkas.sklad_id = sklads.id and item_postavkas.type = 'tovar' order by time desc limit 1) as last_postavka_cost").Joins("inner join sklad_tovars on sklad_tovars.sklad_id = sklads.id inner join tovars on tovars.tovar_id = sklad_tovars.tovar_id").Where("tovars.deleted = ? and sklads.shop_id IN (?) and tovars.shop_id = sklads.shop_id", false, filter.AccessibleShops).Scan(&tovars).Error
		if err != nil {
			return nil, err
		}

		for _, val := range tovars {
			items = append(items, &model.ItemOutput{
				ID:               val.ID,
				Name:             val.Name,
				Measure:          val.Measure,
				Cost:             val.Cost,
				SkladID:          val.SkladID,
				LastPostavkaTime: val.LastPostavkaTime,
				LastPostavkaCost: val.LastPostavkaCost,
				Type:             utils.TypeTovar,
			})
		}
	}
	return items, nil
}

func (r *SkladDB) GetAllPostavka(filter *model.Filter) ([]*model.PostavkaOutput, int64, error) {
	res := r.gormDB
	if filter.Category != 0 {
		res = r.gormDB.Table("postavkas").Select("postavkas.id, json_agg(json_build_object('sum', postavkas.sum, 'dealer', dealers.name, 'sklad', sklads.name, 'schet', schets.name,'time', postavkas.time, 'item_id', item_postavkas.item_id,'postavka_id', item_postavkas.postavka_id, 'type', item_postavkas.type,'quantity', item_postavkas.quantity, 'cost', item_postavkas.cost,'name', CASE WHEN item_postavkas.type = 'ingredient' THEN ingredients.name Else tovars.name END, 'measurement', CASE WHEN item_postavkas.type = 'ingredient' THEN ingredients.measure Else tovars.measure END,'category', CASE WHEN item_postavkas.type = 'ingredient' THEN category_ingredients.name Else category_tovars.name END,'risky', item_postavkas.risky, 'deleted', item_postavkas.deleted)) as postavka_items").Joins("inner join item_postavkas on item_postavkas.postavka_id = postavkas.id inner join sklads on postavkas.sklad_id = sklads.id inner join dealers on postavkas.dealer_id = dealers.id inner join schets on postavkas.schet_id = schets.id left join ingredients on (ingredients.ingredient_id = item_postavkas.item_id and item_postavkas.type = 'ingredient') left join category_ingredients on category_ingredients.id = ingredients.category left join tovars on (tovars.tovar_id = item_postavkas.item_id and  item_postavkas.type = 'tovar') left join category_tovars on category_tovars.id = tovars.category").Where("postavkas.type != ? and sklads.shop_id IN (?) and (ingredients.shop_id = sklads.shop_id or tovars.shop_id = sklads.shop_id) and (tovars.category = ? or ingredients.category = ?)", utils.TypeTransfer, filter.AccessibleShops, filter.Category, filter.Category).Group("postavkas.id")
		if res.Error != nil {
			return nil, 0, res.Error
		}
	} else {
		res = r.gormDB.Table("postavkas").Select("postavkas.id, json_agg(json_build_object('sum', postavkas.sum, 'dealer', dealers.name, 'sklad', sklads.name, 'schet', schets.name,'time', postavkas.time, 'item_id', item_postavkas.item_id,'postavka_id', item_postavkas.postavka_id, 'type', item_postavkas.type,'quantity', item_postavkas.quantity, 'cost', item_postavkas.cost,'name', CASE WHEN item_postavkas.type = 'ingredient' THEN ingredients.name Else tovars.name END, 'measurement', CASE WHEN item_postavkas.type = 'ingredient' THEN ingredients.measure Else tovars.measure END,'category', CASE WHEN item_postavkas.type = 'ingredient' THEN category_ingredients.name Else category_tovars.name END,'risky', item_postavkas.risky, 'deleted', item_postavkas.deleted)) as postavka_items").Joins("inner join item_postavkas on item_postavkas.postavka_id = postavkas.id inner join sklads on postavkas.sklad_id = sklads.id inner join dealers on postavkas.dealer_id = dealers.id inner join schets on postavkas.schet_id = schets.id left join ingredients on (ingredients.ingredient_id = item_postavkas.item_id and item_postavkas.type = 'ingredient') left join category_ingredients on category_ingredients.id = ingredients.category left join tovars on (tovars.tovar_id = item_postavkas.item_id and  item_postavkas.type = 'tovar') left join category_tovars on category_tovars.id = tovars.category").Where("postavkas.type != ? and sklads.shop_id IN (?) and (ingredients.shop_id = sklads.shop_id or tovars.shop_id = sklads.shop_id)", utils.TypeTransfer, filter.AccessibleShops).Group("postavkas.id")
		if res.Error != nil {
			return nil, 0, res.Error
		}
	}

	newRes, count, err := filter.FilterResults(res, model.Postavka{}, utils.DefaultPageSize, "time", fmt.Sprintf("dealers.name ilike '%%%s%%' OR sklads.name ilike '%%%s%%' OR tovars.name ilike '%%%s%%' OR ingredients.name ilike '%%%s%%'", filter.Search, filter.Search, filter.Search, filter.Search), "postavkas.time desc")
	if err != nil {
		return nil, 0, err
	}
	var test model.PostavkaItemsArr
	rows, err := newRes.Rows()
	if err != nil {
		return nil, 0, err
	}

	postavkaNeOutput := []*model.PostavkaAll{}

	for rows.Next() {
		type IDStruct struct {
			ID int64 `json:"id"`
		}
		var id IDStruct
		err = rows.Scan(&id.ID, &test)
		if err != nil {
			return nil, 0, err
		}
		postavkaNeOutput = append(postavkaNeOutput, &model.PostavkaAll{
			ID:            int(id.ID),
			PostavkaItems: test,
		})
	}
	postavkas := []*model.PostavkaOutput{}
	for _, postavka := range postavkaNeOutput {
		flag := false
		itemsPostavka := []*model.ItemPostavkaOutput{}
		for _, items := range postavka.PostavkaItems {
			for _, items1 := range items.PostavkaInfo {
				if items1.Risky {
					flag = true
				}
				itemPostavka := &model.ItemPostavkaOutput{
					ItemID:      items1.ItemID,
					PostavkaID:  postavka.ID,
					Name:        items1.Name,
					Category:    items1.Category,
					Type:        items1.Type,
					Measurement: items1.Measurement,
					Quantity:    items1.Quantity,
					Cost:        items1.Cost,
					Risky:       items1.Risky,
					Deleted:     items1.Deleted,
				}
				itemsPostavka = append(itemsPostavka, itemPostavka)
			}
		}
		postavkas = append(postavkas, &model.PostavkaOutput{
			ID:      postavka.ID,
			Dealer:  postavka.PostavkaItems[0].PostavkaInfo[0].Dealer,
			Sklad:   postavka.PostavkaItems[0].PostavkaInfo[0].Sklad,
			Schet:   postavka.PostavkaItems[0].PostavkaInfo[0].Schet,
			Time:    postavka.PostavkaItems[0].PostavkaInfo[0].Time,
			Sum:     postavka.PostavkaItems[0].PostavkaInfo[0].Sum,
			Risky:   flag,
			Deleted: postavka.PostavkaItems[0].PostavkaInfo[0].Deleted,
			Items:   itemsPostavka,
		})
	}
	return postavkas, count, nil
}

func (r *SkladDB) GetPostavka(id int) (*model.PostavkaOutput, error) {
	postavka := &model.PostavkaOutput{}
	err := r.gormDB.Table("postavkas").Select("postavkas.id, dealers.name as dealer, dealers.id as dealer_id, sklads.name as sklad, sklads.id as sklad_id, schets.name as schet, schets.id as schet_id, postavkas.time, postavkas.deleted, postavkas.sum").Joins("inner join sklads on postavkas.sklad_id = sklads.id inner join dealers on postavkas.dealer_id = dealers.id inner join schets on postavkas.schet_id = schets.id").Where("postavkas.id = ?", id).Scan(&postavka).Error
	if err != nil {
		return nil, err
	}
	items := []*model.ItemPostavka{}
	err = r.gormDB.Model(items).Where("item_postavkas.postavka_id = ?", postavka.ID).Scan(&items).Error
	if err != nil {
		return nil, err
	}
	shopID, err := r.GetShopIDBySkladID(postavka.SkladID)
	if err != nil {
		return nil, err
	}
	itemPostavkas := []*model.ItemPostavkaOutput{}
	for _, item := range items {
		itemPostavka := &model.ItemPostavkaOutput{}
		if item.Type == utils.TypeIngredient {
			err := r.gormDB.Model(items).Select("item_postavkas.item_id, item_postavkas.postavka_id, ingredients.name, ingredients.measure as measurement, category_ingredients.name as category, item_postavkas.type, item_postavkas.quantity, item_postavkas.cost, item_postavkas.deleted").Joins("inner join ingredients on item_postavkas.item_id = ingredients.ingredient_id inner join category_ingredients on ingredients.category = category_ingredients.id").Where("item_postavkas.id = ? and ingredients.shop_id = ?", item.ID, shopID).Scan(&itemPostavka).Error
			if err != nil {
				return nil, err
			}
		} else if item.Type == utils.TypeTovar {
			err := r.gormDB.Model(items).Select("item_postavkas.item_id, item_postavkas.postavka_id, tovars.name, category_tovars.name as category, item_postavkas.type, item_postavkas.quantity, item_postavkas.cost, item_postavkas.deleted").Joins("inner join tovars on item_postavkas.item_id = tovars.tovar_id inner join category_tovars on category_tovars.id = tovars.category").Where("item_postavkas.id = ? and tovars.shop_id = ?", item.ID, shopID).Scan(&itemPostavka).Error
			if err != nil {
				return nil, err
			}
			itemPostavka.Measurement = "шт."
		} else {
			return nil, errors.New("unknown type")
		}
		itemPostavkas = append(itemPostavkas, itemPostavka)
	}
	postavka.Items = itemPostavkas

	return postavka, nil
}

func (r *SkladDB) GetTransactionFromPostavkaID(id int) (*model.Transaction, error) {
	transactionPostavka := &model.TransactionPostavka{}
	err := r.gormDB.Model(&model.TransactionPostavka{}).Where("postavka_id = ?", id).First(transactionPostavka).Error
	if err != nil {
		return nil, err
	}
	transaction := &model.Transaction{}
	err = r.gormDB.Model(&model.Transaction{}).Where("id = ?", transactionPostavka.TransactionID).First(transaction).Error
	if err != nil {
		return nil, err
	}
	return transaction, nil
}

func (r *SkladDB) RecalculateShift(id int) error {
	shift := &model.Shift{}
	err := r.gormDB.Model(&model.Shift{}).Where("id = ?", id).Scan(&shift).Error

	if err != nil {
		return err
	}

	if !shift.IsClosed {
		payment := model.Payment{}
		time := time.Now()
		err := r.gormDB.Table("checks").Select("SUM(checks.cash) AS cash, SUM(checks.card) AS card").Where("checks.status = ? AND checks.closed_at >= ? AND checks.closed_at <= ? and checks.shop_id = ?", utils.StatusClosed, shift.CreatedAt, time, shift.ShopID).Scan(&payment).Error
		if err != nil {
			return err
		}
		shift.Cash = payment.Cash
		shift.Card = payment.Card
	}
	shift.Expense = 0
	shift.EndSumPlan = 0
	shift.Collection = 0
	transactions := []*model.Transaction{}
	err = r.gormDB.Model(&model.Transaction{}).Where("shift_id = ? and deleted = ?", shift.ID, false).Order("id DESC").Scan(&transactions).Error
	if err != nil {
		return err
	}
	for _, transaction := range transactions {
		if transaction.SchetID == shift.SchetID {
			if transaction.Category == utils.Collection {
				shift.Collection += transaction.Sum
			} else {
				if transaction.Status == utils.TransactionNegativeStatus {
					shift.Expense += transaction.Sum
				} else if transaction.Status == utils.TransactionPositiveStatus {
					shift.Cash += transaction.Sum
				}
			}
		}
	}
	shift.EndSumPlan = shift.BeginSum + shift.Cash - shift.Expense - shift.Collection
	if shift.IsClosed {
		shift.Difference = shift.EndSumFact - shift.EndSumPlan
	}
	err = r.gormDB.Model(&model.Shift{}).Where("id = ?", shift.ID).Updates(&shift).Error
	if err != nil {
		return err
	}
	return err
}

func (r *SkladDB) UpdatePostavka(postavka *model.Postavka) error {
	invItemMap := make(map[int]int)
	invItems := []*model.InventarizationItem{}
	trafficItems := []*model.AsyncJob{}

	err := r.gormDB.Transaction(func(tx *gorm.DB) error {
		oldPostavka := &model.Postavka{}
		err := tx.Model(&model.Postavka{}).Where("id = ?", postavka.ID).Scan(&oldPostavka).Error
		if err != nil {
			return err
		}
		items := []*model.ItemPostavka{}
		err = tx.Model(items).Where("postavka_id = ?", postavka.ID).Scan(&items).Error
		if err != nil {
			return err
		}
		swap := false
		if oldPostavka.SkladID != postavka.SkladID {
			swap = true
		}

		for _, item := range items {
			err := tx.Delete(item).Error
			if err != nil {
				return err
			}
			newCost := float32(0)
			invItem := &model.InventarizationItem{}

			if item.Type == utils.TypeIngredient {
				var lastInv time.Time
				err := tx.Model(&model.InventarizationItem{}).Select("time").Where("sklad_id = ? and item_id = ? and type = ? and status != ?", postavka.SkladID, item.ItemID, utils.TypeIngredient, utils.StatusNew).Order("time desc").First(&lastInv).Error
				if err != nil {
					if err != gorm.ErrRecordNotFound {
						return err
					}
				}
				flag := true
				if lastInv.After(postavka.Time) && lastInv.Before(time.Now()) {
					flag = false
				}
				if postavka.Type != utils.TypeTransfer {
					err = r.RecalculateCost(tx, item, utils.Postavka, postavka.SkladID, flag)
					if err != nil {
						return err
					}
				} else {
					if flag {
						skladIngredient := &model.SkladIngredient{}
						if err := tx.Where("sklad_id = ? and ingredient_id = ?", postavka.SkladID, item.ItemID).First(skladIngredient).Error; err != nil {
							return err
						}
						skladIngredient.Quantity -= item.Quantity
						if err := tx.Select("quantity").Where("sklad_id = ? and ingredient_id = ?", postavka.SkladID, item.ItemID).Updates(skladIngredient).Error; err != nil {
							return err
						}
					}
				}
			} else if item.Type == utils.TypeTovar {
				var lastInv time.Time
				err := tx.Model(&model.InventarizationItem{}).Select("time").Where("sklad_id = ? and item_id = ? and type = ? and status != ?", postavka.SkladID, item.ItemID, utils.TypeTovar, utils.StatusNew).Order("time desc").First(&lastInv).Error
				if err != nil {
					if err != gorm.ErrRecordNotFound {
						return err
					}
				}
				flag := true
				if lastInv.After(postavka.Time) && lastInv.Before(time.Now()) {
					flag = false
				}
				if postavka.Type != utils.TypeTransfer {
					err = r.RecalculateCost(tx, item, utils.Postavka, postavka.SkladID, flag)
					if err != nil {
						return err
					}
				} else {
					if flag {
						skladTovar := &model.SkladTovar{}
						if err := tx.Where("sklad_id = ? and tovar_id = ?", postavka.SkladID, item.ItemID).First(skladTovar).Error; err != nil {
							return err
						}
						skladTovar.Quantity -= item.Quantity
						if err := tx.Select("quantity").Where("sklad_id = ? and tovar_id = ?", postavka.SkladID, item.ItemID).Updates(skladTovar).Error; err != nil {
							return err
						}
					}
				}
			} else {
				return errors.New("wrong type")
			}

			err = tx.Model(invItem).Where("sklad_id = ? and item_id = ? and type = ? and time > ? and before_time < ?", oldPostavka.SkladID, item.ItemID, item.Type, oldPostavka.Time, oldPostavka.Time).Scan(invItem).Error
			if err != nil {
				if err.Error() != "record not found" {
					return err
				}
			}
			if _, exists := invItemMap[invItem.ID]; !exists && invItem.ID != 0 && invItem.ItemID != 0 {
				invItemMap[invItem.ID] = invItem.ID
				invItem.Cost = newCost
				invItems = append(invItems, invItem)
			}
		}
		err = tx.Save(postavka).Error
		if err != nil {
			return err
		}
		for _, item := range postavka.Items {
			newCost := float32(0)
			invItem := &model.InventarizationItem{}

			if item.Type == utils.TypeIngredient {
				skladIngredient := &model.SkladIngredient{}
				err := tx.Where("sklad_id = ? AND ingredient_id = ?", postavka.SkladID, item.ItemID).First(skladIngredient).Error
				if err != nil {
					return err
				}
				newCost = skladIngredient.Cost

				var lastInv time.Time
				err = tx.Model(&model.InventarizationItem{}).Select("time").Where("sklad_id = ? and item_id = ? and type = ? and status != ?", postavka.SkladID, item.ItemID, utils.TypeIngredient, utils.StatusNew).Order("time desc").First(&lastInv).Error
				if err != nil {
					if err != gorm.ErrRecordNotFound {
						return err
					}
				}
				flag := true
				if lastInv.After(postavka.Time) && lastInv.Before(time.Now()) {
					flag = false
				}
				if postavka.Type != utils.TypeTransfer {
					err = r.RecalculateCost(tx, &item, utils.Postavka, postavka.SkladID, flag)
					if err != nil {
						return err
					}
				} else {
					if flag {
						if err := tx.Where("sklad_id = ? and ingredient_id = ?", postavka.SkladID, item.ItemID).First(skladIngredient).Error; err != nil {
							return err
						}
						skladIngredient.Quantity += item.Quantity
						if err := tx.Select("quantity").Where("sklad_id = ? and ingredient_id = ?", postavka.SkladID, item.ItemID).Updates(skladIngredient).Error; err != nil {
							return err
						}
					}
				}
				err = tx.Where("sklad_id = ? AND ingredient_id = ?", postavka.SkladID, item.ItemID).First(skladIngredient).Error
				if err != nil {
					return err
				}
				newCost = skladIngredient.Cost

			} else if item.Type == utils.TypeTovar {
				skladTovar := &model.SkladTovar{}
				err := tx.Where("sklad_id = ? AND tovar_id = ?", postavka.SkladID, item.ItemID).First(skladTovar).Error
				if err != nil {
					return err
				}
				newCost = skladTovar.Cost
				var lastInv time.Time
				err = tx.Model(&model.InventarizationItem{}).Select("time").Where("sklad_id = ? and item_id = ? and type = ? and status != ?", postavka.SkladID, item.ItemID, utils.TypeTovar, utils.StatusNew).Order("time desc").First(&lastInv).Error
				if err != nil {
					if err != gorm.ErrRecordNotFound {
						return err
					}
				}
				flag := true
				if lastInv.After(postavka.Time) && lastInv.Before(time.Now()) {
					flag = false
				}
				if postavka.Type != utils.TypeTransfer {
					err = r.RecalculateCost(tx, &item, utils.Postavka, postavka.SkladID, flag)
					if err != nil {
						return err
					}
				} else {
					if flag {
						skladTovar := &model.SkladTovar{}
						if err := tx.Where("sklad_id = ? and tovar_id = ?", postavka.SkladID, item.ItemID).First(skladTovar).Error; err != nil {
							return err
						}
						skladTovar.Quantity -= item.Quantity
						if err := tx.Select("quantity").Where("sklad_id = ? and tovar_id = ?", postavka.SkladID, item.ItemID).Updates(skladTovar).Error; err != nil {
							return err
						}
					}
				}
				err = tx.Where("sklad_id = ? AND tovar_id = ?", postavka.SkladID, item.ItemID).First(skladTovar).Error
				if err != nil {
					return err
				}
				newCost = skladTovar.Cost

			} else {
				return errors.New("wrong type")
			}
			if postavka.Time.Year() != time.Now().Year() || postavka.Time.Month() != time.Now().Month() || postavka.Time.Day() != time.Now().Day() {
				trafficItems = append(trafficItems, &model.AsyncJob{
					ItemID:    item.ItemID,
					ItemType:  item.Type,
					SkladID:   postavka.SkladID,
					TimeStamp: postavka.Time.UTC(),
					CreatedAt: time.Now(),
					Status:    utils.StatusNeedRecalculate,
				})
			}
			err = tx.Model(invItem).Where("sklad_id = ? and item_id = ? and type = ? and time > ? and before_time < ?", postavka.SkladID, item.ItemID, item.Type, postavka.Time, postavka.Time).Scan(invItem).Error
			if err != nil {
				if err.Error() != "record not found" {
					return err
				}
			}
			if _, exists := invItemMap[invItem.ID]; !exists && invItem.ID != 0 && invItem.ItemID != 0 {
				invItemMap[invItem.ID] = invItem.ID
				invItem.Cost = newCost
				invItems = append(invItems, invItem)
			}
		}
		if postavka.Type != utils.TypeTransfer {
			transactionPostavka := &model.TransactionPostavka{}
			err = r.gormDB.Model(&model.TransactionPostavka{}).Where("postavka_id = ?", postavka.ID).First(transactionPostavka).Error
			if err != nil {
				return nil
			}
			transaction := &model.Transaction{}
			err = r.gormDB.Model(&model.Transaction{}).Where("id = ?", transactionPostavka.TransactionID).First(transaction).Error
			if err != nil {
				return nil
			}
			change := false
			if transaction.SchetID != postavka.SchetID {
				change = true
			}

			shifID := transaction.ShiftID
			err = r.gormDB.Model(&model.Schet{}).Where("id = ?", transaction.SchetID).Update("start_balance", gorm.Expr("start_balance + ?", transaction.Sum)).Error
			if err != nil {
				return nil
			}
			if oldPostavka.Time != postavka.Time {
				change = true
			}
			if change {
				err = r.gormDB.Model(&model.Shift{}).Select("id").Where("schet_id = ? and created_at <= ?", postavka.SchetID, postavka.Time).Order("created_at desc").First(&shifID).Error
				if err != nil {
					if err != gorm.ErrRecordNotFound {
						return nil
					}
				}
				if transaction.ShiftID != shifID {
					swap = false
				}
			}
			if swap {
				var shopID int
				err = r.gormDB.Model(&model.Sklad{}).Select("shop_id").Where("id = ?", postavka.SkladID).Scan(&shopID).Error
				if err != nil {
					return nil
				}
				err = r.gormDB.Model(&model.Shift{}).Select("id").Where("shop_id = ? and created_at <= ?", transaction.SchetID, postavka.Time).Order("created_at desc").First(&shifID).Error
				if err != nil {
					return nil
				}
			}
			transaction.Time = postavka.Time
			transaction.Sum = postavka.Sum
			transaction.SchetID = postavka.SchetID
			transaction.ShiftID = shifID
			err = r.UpdateTransaction(transaction)
			if err != nil {
				return nil
			}
			err = r.RecalculateShift(transaction.ShiftID)
			if err != nil {
				return nil
			}
		}
		return nil
	})

	if err != nil {
		return err
	}
	r.ConcurrentRecalculationForDailyStatistics(trafficItems)
	if len(invItems) > 0 {
		err = r.RecalculateInventarization(invItems)
		if err != nil {
			return err
		}
	}
	return nil
}

func (r *SkladDB) UpdateTransaction(transaction *model.Transaction) error {
	err := r.gormDB.Model(&model.Transaction{}).Where("id = ?", transaction.ID).Updates(transaction).Error
	if err != nil {
		return err
	}
	err = r.gormDB.Model(&model.Schet{}).Where("id = ?", transaction.SchetID).Update("start_balance", gorm.Expr("start_balance - ?", transaction.Sum)).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *SkladDB) DeletePostavka(id int) error {
	invItems := []*model.InventarizationItem{}
	trafficItems := []*model.AsyncJob{}
	err := r.gormDB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&model.Postavka{}).Where("id = ?", id).Update("deleted", true).Error; err != nil {
			return err
		}
		postavka := model.Postavka{}
		err := tx.Where("id = ?", id).First(&postavka).Error
		if err != nil {
			return err
		}
		items := []*model.ItemPostavka{}

		err = tx.Model(items).Where("postavka_id = ?", id).Scan(&items).Error
		if err != nil {
			return err
		}
		for _, item := range items {
			err := tx.Model(&model.ItemPostavka{}).Where("item_id = ? and postavka_id = ?", item.ItemID, item.PostavkaID).Update("deleted", true).Error
			if err != nil {
				return err
			}
			invItem := &model.InventarizationItem{}
			newCost := float32(0)
			if item.Type == utils.TypeIngredient {
				skladIngredient := &model.SkladIngredient{}
				err := tx.Where("sklad_id = ? AND ingredient_id = ?", postavka.SkladID, item.ItemID).First(skladIngredient).Error
				if err != nil {
					return err
				}
				newCost = skladIngredient.Cost

				var lastInv time.Time
				err = tx.Model(&model.InventarizationItem{}).Select("time").Where("sklad_id = ? and item_id = ? and type = ? and status != ?", postavka.SkladID, item.ItemID, utils.TypeIngredient, utils.StatusNew).Order("time desc").First(&lastInv).Error
				if err != nil {
					if err != gorm.ErrRecordNotFound {
						return err
					}
				}
				flag := true
				if lastInv.After(postavka.Time) && lastInv.Before(time.Now()) {
					flag = false
				}
				if postavka.Type != utils.TypeTransfer {
					err = r.RecalculateCost(tx, item, utils.Postavka, postavka.SkladID, flag)
					if err != nil {
						return err
					}
				} else {
					if flag {
						skladIngredient := &model.SkladIngredient{}
						if err := tx.Where("sklad_id = ? and ingredient_id = ?", postavka.SkladID, item.ItemID).First(skladIngredient).Error; err != nil {
							return err
						}
						skladIngredient.Quantity -= item.Quantity
						if err := tx.Select("quantity").Where("sklad_id = ? and ingredient_id = ?", postavka.SkladID, item.ItemID).Updates(skladIngredient).Error; err != nil {
							return err
						}
					}
				}
				err = tx.Where("sklad_id = ? AND ingredient_id = ?", postavka.SkladID, item.ItemID).First(skladIngredient).Error
				if err != nil {
					return err
				}
				newCost = skladIngredient.Cost
			} else if item.Type == utils.TypeTovar {
				skladTovar := &model.SkladTovar{}
				err := tx.Where("sklad_id = ? AND tovar_id = ?", postavka.SkladID, item.ItemID).First(skladTovar).Error
				if err != nil {
					return err
				}
				newCost = skladTovar.Cost

				var lastInv time.Time
				err = tx.Model(&model.InventarizationItem{}).Select("time").Where("sklad_id = ? and item_id = ? and type = ? and status != ?", postavka.SkladID, item.ItemID, utils.TypeTovar, utils.StatusNew).Order("time desc").First(&lastInv).Error
				if err != nil {
					if err != gorm.ErrRecordNotFound {
						return err
					}
				}
				flag := true
				if lastInv.After(postavka.Time) && lastInv.Before(time.Now()) {
					flag = false
				}
				if postavka.Type != utils.TypeTransfer {
					err = r.RecalculateCost(tx, item, utils.Postavka, postavka.SkladID, flag)
					if err != nil {
						return err
					}
				} else {
					if flag {
						skladTovar := &model.SkladTovar{}
						if err := tx.Where("sklad_id = ? and tovar_id = ?", postavka.SkladID, item.ItemID).First(skladTovar).Error; err != nil {
							return err
						}
						skladTovar.Quantity -= item.Quantity
						if err := tx.Select("quantity").Where("sklad_id = ? and tovar_id = ?", postavka.SkladID, item.ItemID).Updates(skladTovar).Error; err != nil {
							return err
						}
					}
				}
				err = tx.Where("sklad_id = ? AND tovar_id = ?", postavka.SkladID, item.ItemID).First(skladTovar).Error
				if err != nil {
					return err
				}
				newCost = skladTovar.Cost
			} else {
				return errors.New("wrong type")
			}
			if postavka.Time.Year() != time.Now().Year() || postavka.Time.Month() != time.Now().Month() || postavka.Time.Day() != time.Now().Day() {
				trafficItems = append(trafficItems, &model.AsyncJob{
					ItemID:    item.ItemID,
					ItemType:  item.Type,
					SkladID:   postavka.SkladID,
					TimeStamp: postavka.Time.UTC(),
					CreatedAt: time.Now(),
					Status:    utils.StatusNeedRecalculate,
				})
			}
			err = tx.Model(invItem).Where("sklad_id = ? and item_id = ? and type = ? and time > ? and before_time < ?", postavka.SkladID, item.ItemID, item.Type, postavka.Time, postavka.Time).Scan(invItem).Error
			if err != nil {
				if err.Error() != "record not found" {
					return err
				}
			}
			if invItem.ID != 0 && invItem.ItemID != 0 {
				invItem.Cost = newCost
				invItems = append(invItems, invItem)
			}
		}
		if postavka.Type != utils.TypeTransfer {
			transactionPostavka := &model.TransactionPostavka{}
			err = r.gormDB.Model(&model.TransactionPostavka{}).Where("postavka_id = ?", postavka.ID).First(transactionPostavka).Error
			if err != nil {
				return err
			}
			err = r.gormDB.Model(&model.TransactionPostavka{}).Where("transaction_id = ? and postavka_id = ?", transactionPostavka.TransactionID, transactionPostavka.PostavkaID).Delete(&model.TransactionPostavka{}).Error
			if err != nil {
				return err
			}

			transaction := &model.Transaction{}
			err = r.gormDB.Model(&model.Transaction{}).Where("id = ?", transactionPostavka.TransactionID).First(transaction).Error
			if err != nil {
				return err
			}
			err = r.gormDB.Model(&model.Schet{}).Where("id = ?", transaction.SchetID).Update("start_balance", gorm.Expr("start_balance + ?", transaction.Sum)).Error
			if err != nil {
				return err
			}
			err = r.gormDB.Model(&model.Transaction{}).Where("id = ?", transaction.ID).Delete(&model.Transaction{}).Error
			if err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return err
	}
	r.ConcurrentRecalculationForDailyStatistics(trafficItems)
	if len(invItems) > 0 {
		err = r.RecalculateInventarization(invItems)
		if err != nil {
			return err
		}
	}
	return nil
}

func (r *SkladDB) AddTransfer(transfer *model.Transfer) (*model.Transfer, error) {
	trafficItems := []*model.AsyncJob{}
	err := r.gormDB.Transaction(func(tx *gorm.DB) error {
		var cost float32
		for _, t := range transfer.ItemTransfers {
			if t.Type == utils.TypeIngredient {
				skladIngredient := &model.SkladIngredient{}

				if err := tx.Model(skladIngredient).Select("*").Where("sklad_id = ? and ingredient_id = ?", transfer.FromSklad, t.ItemID).Scan(&skladIngredient).Error; err != nil {
					return err
				}

				if skladIngredient.Quantity < t.Quantity {
					transfer.Less = true
				}

				cost += (skladIngredient.Cost * float32(t.Quantity))
				t.Sum = skladIngredient.Cost * float32(t.Quantity)

				err := tx.Model(skladIngredient).Where("sklad_id = ? and ingredient_id = ?", skladIngredient.SkladID, skladIngredient.IngredientID).Where("sklad_id = ? and ingredient_id = ?", skladIngredient.SkladID, skladIngredient.IngredientID).Updates(skladIngredient).Error
				if err != nil {
					return err
				}
			} else if t.Type == utils.TypeTovar {
				skladTovars := &model.SkladTovar{}

				if err := tx.Model(skladTovars).Where("sklad_id = ? and tovar_id = ?", transfer.FromSklad, t.ItemID).Scan(&skladTovars).Error; err != nil {
					return err
				}

				if skladTovars.Quantity < t.Quantity {
					transfer.Less = true
				}

				cost += (skladTovars.Cost * float32(t.Quantity))
				t.Sum = skladTovars.Cost * float32(t.Quantity)

				err := tx.Model(skladTovars).Where("sklad_id = ? and tovar_id = ?", skladTovars.SkladID, skladTovars.TovarID).Where("sklad_id = ? and tovar_id = ?", skladTovars.SkladID, skladTovars.TovarID).Updates(skladTovars).Error
				if err != nil {
					return err
				}
			} else {
				return errors.New("unknown type")
			}
			if transfer.Time.Year() != time.Now().Year() || transfer.Time.Month() != time.Now().Month() || transfer.Time.Day() != time.Now().Day() {
				trafficItems = append(trafficItems, &model.AsyncJob{
					ItemID:    t.ItemID,
					ItemType:  t.Type,
					SkladID:   transfer.FromSklad,
					TimeStamp: transfer.Time.UTC(),
					CreatedAt: time.Now(),
					Status:    utils.StatusNeedRecalculate,
				}, &model.AsyncJob{
					ItemID:    t.ItemID,
					ItemType:  t.Type,
					SkladID:   transfer.ToSklad,
					TimeStamp: transfer.Time.UTC(),
					CreatedAt: time.Now(),
					Status:    utils.StatusNeedRecalculate,
				})
			}
		}
		transfer.Sum = cost

		if err := tx.Create(transfer).Error; err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	r.ConcurrentRecalculationForDailyStatistics(trafficItems)
	return transfer, nil
}

func (r *SkladDB) GetTransfer(id int) (*model.TransferOutput, error) {
	transfer := &model.TransferOutput{}
	err := r.gormDB.Table("transfers").Select("transfers.id, transfers.time, transfers.from_sklad, transfers.to_sklad, skladFrom.name as from_sklad_name, skladTo.name as to_sklad_name, transfers.worker, workers.name as worker_name, transfers.sum").Joins("inner join sklads skladFrom on skladFrom.id = transfers.from_sklad inner join sklads skladTo on skladTo.id = transfers.to_sklad inner join workers on workers.id = transfers.worker").Where("transfers.id = ?", id).Scan(&transfer).Error
	if err != nil {
		return nil, err
	}
	items := []*model.ItemTransfer{}
	err = r.gormDB.Model(items).Where("item_transfers.transfer_id = ?", transfer.ID).Scan(&items).Error
	if err != nil {
		return nil, err
	}
	itemTransfers := []*model.ItemTransferOutput{}
	for _, item := range items {
		itemTransfer := &model.ItemTransferOutput{}
		if item.Type == utils.TypeIngredient {
			err := r.gormDB.Model(items).Select("item_transfers.id, ingredients.name, sklad_ingredients.cost, item_transfers.item_ID, item_transfers.transfer_id, ingredients.measure as measurement, category_ingredients.name as category, item_transfers.type, item_transfers.quantity, item_transfers.sum").Joins("inner join ingredients on item_transfers.item_id = ingredients.ingredient_id inner join category_ingredients on ingredients.category = category_ingredients.id inner join sklad_ingredients on sklad_ingredients.ingredient_id = item_transfers.item_ID").Where("item_transfers.item_id = ? and item_transfers.transfer_id = ?", item.ItemID, transfer.ID).Scan(&itemTransfer).Error
			if err != nil {
				return nil, err
			}
		} else if item.Type == utils.TypeTovar {
			err := r.gormDB.Model(items).Select("item_transfers.id, tovars.name, sklad_tovars.cost, item_transfers.item_ID, item_transfers.transfer_id, category_tovars.name as category, item_transfers.type, item_transfers.quantity, item_transfers.sum").Joins("inner join tovars on item_transfers.item_id = tovars.tovar_id inner join category_tovars on category_tovars.id = tovars.category inner join sklad_tovars on sklad_tovars.tovar_id = item_transfers.item_ID").Where("item_transfers.item_id = ? and item_transfers.transfer_id = ?", item.ItemID, transfer.ID).Scan(&itemTransfer).Error
			if err != nil {
				return nil, err
			}
			itemTransfer.Measurement = "шт."
		} else {
			return nil, errors.New("unknown type")
		}
		itemTransfers = append(itemTransfers, itemTransfer)
	}
	transfer.ItemTransfers = itemTransfers

	return transfer, nil
}

func (r *SkladDB) GetAllTransfer(filter *model.Filter) ([]*model.TransferOutput, int64, error) {
	transfer := []*model.TransferOutput{}
	res := r.gormDB.Table("transfers").Select("transfers.id, transfers.sum, skladFrom.name as from_sklad_name, skladTo.name as to_sklad_name, transfers.from_sklad, transfers.to_sklad, transfers.worker, workers.name as worker_name, transfers.time").Joins("inner join sklads skladFrom on skladFrom.id = transfers.from_sklad inner join sklads skladTo on skladTo.id = transfers.to_sklad inner join workers on workers.id = transfers.worker").Where("skladFrom.shop_id IN (?) and skladTo.shop_id IN (?) and transfers.deleted = ?", filter.AccessibleShops, filter.AccessibleShops, false)
	if res.Error != nil {
		return nil, 0, res.Error
	}
	newRes, count, err := filter.FilterResults(res, model.Transfer{}, utils.DefaultPageSize, "time", fmt.Sprintf("skladFrom.name ilike '%%%s%%' OR skladTo.name ilike '%%%s%%'", filter.Search, filter.Search), "")
	if err != nil {
		return nil, 0, err
	}
	err = newRes.Scan(&transfer).Error
	if err != nil {
		return nil, 0, err
	}

	for i := 0; i < len(transfer); i++ {
		items := []*model.ItemTransfer{}
		err := r.gormDB.Model(items).Where("item_transfers.transfer_id = ?", transfer[i].ID).Scan(&items).Error
		if err != nil {
			return nil, 0, err
		}
		itemTransfers := []*model.ItemTransferOutput{}
		for _, item := range items {
			itemTransfer := &model.ItemTransferOutput{}
			if item.Type == utils.TypeIngredient {
				err := r.gormDB.Model(items).Select("item_transfers.id, sklad_ingredients.cost, item_transfers.item_ID, item_transfers.transfer_id, ingredients.name, ingredients.measure as measurement, category_ingredients.name as category, item_transfers.type, item_transfers.quantity, item_transfers.sum").Joins("inner join ingredients on ingredients.ingredient_id = item_transfers.item_id inner join category_ingredients on category_ingredients.id = ingredients.category inner join sklad_ingredients on sklad_ingredients.ingredient_id = item_transfers.item_ID").Where("item_transfers.item_id = ? and item_transfers.transfer_id = ?", item.ItemID, transfer[i].ID).Scan(&itemTransfer).Error
				if err != nil {
					return nil, 0, err
				}
			} else if item.Type == utils.TypeTovar {
				err := r.gormDB.Model(items).Select("item_transfers.id, sklad_tovars.cost, item_transfers.item_ID, item_transfers.transfer_id, tovars.name, category_tovars.name as category, item_transfers.type, item_transfers.quantity, item_transfers.sum").Joins("inner join tovars on item_transfers.item_id = tovars.tovar_id inner join category_tovars on category_tovars.id = tovars.category inner join sklad_tovars on sklad_tovars.tovar_id = item_transfers.item_ID").Where("item_transfers.item_id = ? and item_transfers.transfer_id = ?", item.ItemID, transfer[i].ID).Scan(&itemTransfer).Error
				if err != nil {
					return nil, 0, err
				}
				itemTransfer.Measurement = "шт."
			} else {
				return nil, 0, errors.New("unknown type")
			}
			itemTransfers = append(itemTransfers, itemTransfer)
		}
		transfer[i].ItemTransfers = itemTransfers
	}
	return transfer, count, nil
}

func (r *SkladDB) UpdateTransfer(transfer *model.Transfer) error {
	id := transfer.ID
	trafficItems := []*model.AsyncJob{}
	err := r.gormDB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&model.Transfer{}).Where("id = ?", id).Error; err != nil {
			return err
		}
		updTransfer := model.Transfer{}
		err := tx.Where("id = ?", id).First(&updTransfer).Error
		if err != nil {
			return err
		}
		items := []*model.ItemTransfer{}

		err = tx.Model(items).Where("transfer_id = ?", id).Scan(&items).Error
		if err != nil {
			return err
		}

		for _, item := range items {
			err := tx.Delete(item).Error
			if err != nil {
				return err
			}
			// if item.Type == utils.TypeIngredient {
			// 	skladIngredient := &model.SkladIngredient{}
			// 	err := tx.Where("sklad_id = ? AND ingredient_id = ?", updTransfer.ToSklad, item.ItemID).First(skladIngredient).Error
			// 	if err != nil {
			// 		return err
			// 	}
			// 	skladIngredient.Quantity -= float32(item.Quantity)
			// 	err = tx.Where("sklad_id = ? AND ingredient_id = ?", updTransfer.ToSklad, item.ItemID).Save(skladIngredient).Error
			// 	if err != nil {
			// 		return err
			// 	}

			// 	err = tx.Where("sklad_id = ? AND ingredient_id = ?", updTransfer.FromSklad, item.ItemID).First(skladIngredient).Error
			// 	if err != nil {
			// 		return err
			// 	}
			// 	skladIngredient.Quantity += float32(item.Quantity)
			// 	err = tx.Where("sklad_id = ? AND ingredient_id = ?", updTransfer.FromSklad, item.ItemID).Save(skladIngredient).Error
			// 	if err != nil {
			// 		return err
			// 	}
			// } else if item.Type == utils.TypeTovar {
			// 	skladTovar := &model.SkladTovar{}
			// 	err := tx.Where("sklad_id = ? AND tovar_id = ?", updTransfer.ToSklad, item.ItemID).First(skladTovar).Error
			// 	if err != nil {
			// 		return err
			// 	}
			// 	skladTovar.Quantity -= float32(item.Quantity)
			// 	err = tx.Where("sklad_id = ? AND tovar_id = ?", updTransfer.ToSklad, item.ItemID).Save(skladTovar).Error
			// 	if err != nil {
			// 		return err
			// 	}

			// 	err = tx.Where("sklad_id = ? AND tovar_id = ?", updTransfer.FromSklad, item.ItemID).First(skladTovar).Error
			// 	if err != nil {
			// 		return err
			// 	}
			// 	skladTovar.Quantity += float32(item.Quantity)
			// 	err = tx.Where("sklad_id = ? AND tovar_id = ?", updTransfer.FromSklad, item.ItemID).Save(skladTovar).Error
			// 	if err != nil {
			// 		return err
			// 	}
			// } else {
			// 	return errors.New("wrong type")
			// }
		}

		var cost float32

		for _, t := range transfer.ItemTransfers {
			if t.Type == utils.TypeIngredient {
				skladIngredient := &model.SkladIngredient{}

				if err := tx.Model(skladIngredient).Select("*").Where("sklad_id = ? and ingredient_id = ?", transfer.FromSklad, t.ItemID).Scan(&skladIngredient).Error; err != nil {
					return err
				}

				cost += (skladIngredient.Cost * float32(t.Quantity))
				t.Sum = skladIngredient.Cost * float32(t.Quantity)

				if skladIngredient.Quantity < float32(t.Quantity) {
					transfer.Less = true
				}

				err = tx.Model(skladIngredient).Where("sklad_id = ? and ingredient_id = ?", skladIngredient.SkladID, skladIngredient.IngredientID).Save(skladIngredient).Where("sklad_id = ? and ingredient_id = ?", skladIngredient.SkladID, skladIngredient.IngredientID).Error
				if err != nil {
					return err
				}

			} else if t.Type == utils.TypeTovar {
				skladTovars := &model.SkladTovar{}

				if err := tx.Model(skladTovars).Where("sklad_id = ? and tovar_id = ?", transfer.FromSklad, t.ItemID).Scan(&skladTovars).Error; err != nil {
					return err
				}

				cost += (skladTovars.Cost * float32(t.Quantity))
				t.Sum = skladTovars.Cost * float32(t.Quantity)

				if skladTovars.Quantity < float32(t.Quantity) {
					transfer.Less = true
				}

				err = tx.Model(skladTovars).Where("sklad_id = ? and tovar_id = ?", skladTovars.SkladID, skladTovars.TovarID).Save(skladTovars).Where("sklad_id = ? and tovar_id = ?", skladTovars.SkladID, skladTovars.TovarID).Error
				if err != nil {
					return err
				}

			} else {
				return errors.New("unknown type")
			}
			t.TransferID = transfer.ID
			if err := tx.Create(t).Error; err != nil {
				return err
			}
			if transfer.Time.Year() != time.Now().Year() || transfer.Time.Month() != time.Now().Month() || transfer.Time.Day() != time.Now().Day() {
				trafficItems = append(trafficItems, &model.AsyncJob{
					ItemID:    t.ItemID,
					ItemType:  t.Type,
					SkladID:   transfer.FromSklad,
					TimeStamp: transfer.Time.UTC(),
					CreatedAt: time.Now(),
					Status:    utils.StatusNeedRecalculate,
				}, &model.AsyncJob{
					ItemID:    t.ItemID,
					ItemType:  t.Type,
					SkladID:   transfer.ToSklad,
					TimeStamp: transfer.Time.UTC(),
					CreatedAt: time.Now(),
					Status:    utils.StatusNeedRecalculate,
				})
			}
		}
		transfer.Sum = cost
		err = tx.Save(transfer).Error
		if err != nil {
			return err
		}
		return nil
	},
	)

	if err != nil {
		return err
	}
	r.ConcurrentRecalculationForDailyStatistics(trafficItems)

	return nil
}

func (r *SkladDB) DeleteTransfer(id int) error {
	trafficItems := []*model.AsyncJob{}
	err := r.gormDB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&model.Transfer{}).Where("id = ?", id).Update("deleted", true).Error; err != nil {
			return err
		}
		transfer := model.Transfer{}
		err := tx.Where("id = ?", id).First(&transfer).Error
		if err != nil {
			return err
		}
		items := []*model.ItemTransfer{}

		err = tx.Model(items).Where("transfer_id = ?", id).Scan(&items).Error
		if err != nil {
			return err
		}

		for _, item := range items {
			err := tx.Delete(item).Error
			if err != nil {
				return err
			}
			// if item.Type == utils.TypeIngredient {
			// 	skladIngredient := &model.SkladIngredient{}
			// 	err := tx.Where("sklad_id = ? AND ingredient_id = ?", transfer.ToSklad, item.ItemID).Scan(&skladIngredient).Error
			// 	if err != nil {
			// 		return err
			// 	}
			// 	skladIngredient.Quantity -= float32(item.Quantity)
			// 	err = tx.Where("sklad_id = ? AND ingredient_id = ?", transfer.ToSklad, item.ItemID).Updates(skladIngredient).Error
			// 	if err != nil {
			// 		return err
			// 	}

			// 	err = tx.Where("sklad_id = ? AND ingredient_id = ?", transfer.FromSklad, item.ItemID).Scan(&skladIngredient).Error
			// 	if err != nil {
			// 		return err
			// 	}
			// 	skladIngredient.Quantity += float32(item.Quantity)
			// 	err = tx.Where("sklad_id = ? AND ingredient_id = ?", transfer.FromSklad, item.ItemID).Updates(skladIngredient).Error
			// 	if err != nil {
			// 		return err
			// 	}
			// } else if item.Type == utils.TypeTovar {
			// 	skladTovar := &model.SkladTovar{}
			// 	err := tx.Where("sklad_id = ? AND tovar_id = ?", transfer.ToSklad, item.ItemID).Scan(&skladTovar).Error
			// 	if err != nil {
			// 		return err
			// 	}
			// 	skladTovar.Quantity -= float32(item.Quantity)
			// 	err = tx.Where("sklad_id = ? AND tovar_id = ?", transfer.ToSklad, item.ItemID).Updates(skladTovar).Error
			// 	if err != nil {
			// 		return err
			// 	}

			// 	err = tx.Where("sklad_id = ? AND tovar_id = ?", transfer.FromSklad, item.ItemID).Scan(&skladTovar).Error
			// 	if err != nil {
			// 		return err
			// 	}
			// 	skladTovar.Quantity += float32(item.Quantity)
			// 	err = tx.Where("sklad_id = ? AND tovar_id = ?", transfer.FromSklad, item.ItemID).Updates(skladTovar).Error
			// 	if err != nil {
			// 		return err
			// 	}
			// } else {
			// 	return errors.New("wrong type")
			// }
			if transfer.Time.Year() != time.Now().Year() || transfer.Time.Month() != time.Now().Month() || transfer.Time.Day() != time.Now().Day() {
				trafficItems = append(trafficItems, &model.AsyncJob{
					ItemID:    item.ItemID,
					ItemType:  item.Type,
					SkladID:   transfer.FromSklad,
					TimeStamp: transfer.Time.UTC(),
					CreatedAt: time.Now(),
					Status:    utils.StatusNeedRecalculate,
				}, &model.AsyncJob{
					ItemID:    item.ItemID,
					ItemType:  item.Type,
					SkladID:   transfer.ToSklad,
					TimeStamp: transfer.Time.UTC(),
					CreatedAt: time.Now(),
					Status:    utils.StatusNeedRecalculate,
				})
			}
		}
		err = tx.Save(transfer).Error
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return err
	}
	r.ConcurrentRecalculationForDailyStatistics(trafficItems)
	return nil
}

// func (r *SkladDB) GetOpenInventarizationResponse(inventarization *model.InventarizationResponse) (*model.InventarizationResponse, error) {
// 	invItems := []*model.InventarizationItemResponse{}
// 	for _, items := range inventarization.InventarizationItems {
// 		invItem := &model.InventarizationItemResponse{}
// 		err := r.gormDB.Model(&model.InventarizationItem{}).Where("item_id = ? AND type = ? AND status = ? AND sklad_id = ? AND after = ?", items.ItemID, items.Type, utils.StatusClosed, inventarization.SkladID, 0).Scan(invItem).Error
// 		if err != nil {
// 			if err == gorm.ErrRecordNotFound {
// 				itemToScan := &model.InventarizationItemResponse{}
// 				if items.Type == utils.TypeIngredient {
// 					err = r.gormDB.Model(&model.Ingredient{}).Select("ingredients.name as item_name, ingredients.measure, ingredients.cost").Where("ingredients.id = ?", items.ItemID).Scan(itemToScan).Error
// 					if err != nil {
// 						return nil, err
// 					}
// 					err = r.gormDB.Model(&model.Postavka{}).Select("SUM(item_postavkas.quantity) as income").Joins("LEFT JOIN item_postavkas ON item_postavkas.postavka_id = postavkas.id").Where("postavkas.time < ? AND item_postavkas.item_id = ? AND item_postavkas.type = ?", inventarization.Time, items.ItemID, utils.TypeIngredient).Scan(&itemToScan).Error
// 					if err != nil {
// 						return nil, err
// 					}
// 					err = r.gormDB.Model(&model.RemoveFromSklad{}).Select("SUM(remove_from_sklad_items.quantity) as removed, SUM(ingredients.cost * remove_from_sklad_items.quantity) as removed_sum").Joins("LEFT JOIN remove_from_sklad_items ON remove_from_sklad_items.remove_id = remove_from_sklads.id").Joins("LEFT JOIN ingredients ON ingredients.id = remove_from_sklad_items.item_id").Where("remove_from_sklads.time < ? AND remove_from_sklad_items.item_id = ? AND remove_from_sklad_items.type = ?", inventarization.Time, items.ItemID, utils.TypeIngredient).Scan(&itemToScan).Error
// 					if err != nil {
// 						return nil, err
// 					}
// 				} else if items.Type == utils.TypeTovar {
// 					err = r.gormDB.Model(&model.Tovar{}).Select("tovars.name as item_name, tovars.measure, tovars.cost").Where("tovars.id = ?", items.ItemID).Scan(itemToScan).Error
// 					if err != nil {
// 						return nil, err
// 					}
// 					err = r.gormDB.Model(&model.Postavka{}).Select("SUM(item_postavkas.quantity) as income").Joins("LEFT JOIN item_postavkas ON item_postavkas.postavka_id = postavkas.id").Where("postavkas.time < ? AND item_postavkas.item_id = ? AND item_postavkas.type = ?", inventarization.Time, items.ItemID, utils.TypeTovar).Scan(&itemToScan).Error
// 					if err != nil {
// 						return nil, err
// 					}
// 					err = r.gormDB.Model(&model.RemoveFromSklad{}).Select("SUM(remove_from_sklad_items.quantity) as removed, SUM(tovars.cost * remove_from_sklad_items.quantity) as removed_sum").Joins("LEFT JOIN remove_from_sklad_items ON remove_from_sklad_items.remove_id = remove_from_sklads.id").Joins("LEFT JOIN tovars ON tovars.id = remove_from_sklad_items.item_id").Where("remove_from_sklads.time < ? AND remove_from_sklad_items.item_id = ? AND remove_from_sklad_items.type = ?", inventarization.Time, items.ItemID, utils.TypeTovar).Scan(&itemToScan).Error
// 					if err != nil {
// 						return nil, err
// 					}
// 				} else {
// 					return nil, errors.New("type not found")
// 				}
// 				itemToScan.PlanQuantity = items.PlanQuantity
// 				itemToScan.ItemID = items.ItemID
// 				itemToScan.Time = inventarization.Time
// 				itemToScan.Type = items.Type
// 				itemToScan.Status = utils.StatusOpened
// 				itemToScan.SkladID = inventarization.SkladID
// 				itemToScan.StartQuantity = itemToScan.PlanQuantity - itemToScan.Income + itemToScan.Removed
// 				itemToScan.Expenses = 0
// 				itemToScan.Before = 0
// 				itemToScan.After = 0
// 				invItems = append(invItems, itemToScan)
// 			} else {
// 				return nil, err
// 			}
// 		} else {
// 			itemToScan := &model.InventarizationItemResponse{}
// 			if items.Type == utils.TypeIngredient {
// 				err = r.gormDB.Model(&model.Ingredient{}).Select("ingredients.name as item_name, ingredients.measure, ingredients.cost").Where("ingredients.id = ?", items.ItemID).Scan(itemToScan).Error
// 				if err != nil {
// 					return nil, err
// 				}
// 				err = r.gormDB.Model(&model.Postavka{}).Select("SUM(item_postavkas.quantity) as income").Joins("LEFT JOIN item_postavkas ON item_postavkas.postavka_id = postavkas.id").Where("postavkas.time < ? AND postavkas.time > ? AND item_postavkas.item_id = ? AND item_postavkas.type = ?", inventarization.Time, invItem.Time, items.ItemID, utils.TypeIngredient).Scan(&itemToScan).Error
// 				if err != nil {
// 					return nil, err
// 				}
// 				err = r.gormDB.Model(&model.RemoveFromSklad{}).Select("SUM(remove_from_sklad_items.quantity) as removed, SUM(ingredients.cost * remove_from_sklad_items.quantity) as removed_sum").Joins("LEFT JOIN remove_from_sklad_items ON remove_from_sklad_items.remove_id = remove_from_sklads.id").Joins("LEFT JOIN ingredients ON ingredients.id = remove_from_sklad_items.item_id").Where("remove_from_sklads.time < ? AND remove_from_sklads.time > ? AND remove_from_sklad_items.item_id = ? AND remove_from_sklad_items.type = ?", items.Time, invItem.Time, items.ItemID, utils.TypeIngredient).Scan(&itemToScan).Error
// 				if err != nil {
// 					return nil, err
// 				}
// 			} else if items.Type == utils.TypeTovar {
// 				err = r.gormDB.Model(&model.Tovar{}).Select("tovars.name as item_name, tovars.measure, tovars.cost").Where("tovars.id = ?", items.ItemID).Scan(itemToScan).Error
// 				if err != nil {
// 					return nil, err
// 				}
// 				err = r.gormDB.Model(&model.Postavka{}).Select("SUM(item_postavkas.quantity) as income").Joins("LEFT JOIN item_postavkas ON item_postavkas.postavka_id = postavkas.id").Where("postavkas.time < ? AND postavkas.time > ? AND item_postavkas.item_id = ? AND item_postavkas.type = ?", inventarization.Time, invItem.Time, items.ItemID, utils.TypeTovar).Scan(&itemToScan).Error
// 				if err != nil {
// 					return nil, err
// 				}
// 				err = r.gormDB.Model(&model.RemoveFromSklad{}).Select("SUM(remove_from_sklad_items.quantity) as removed, SUM(tovars.cost * remove_from_sklad_items.quantity) as removed_sum").Joins("LEFT JOIN remove_from_sklad_items ON remove_from_sklad_items.remove_id = remove_from_sklads.id").Joins("LEFT JOIN tovars ON tovars.id = remove_from_sklad_items.item_id").Where("remove_from_sklads.time < ? AND remove_from_sklads.time > ? AND remove_from_sklad_items.item_id = ? AND remove_from_sklad_items.type = ?", items.Time, invItem.Time, items.ItemID, utils.TypeTovar).Scan(&itemToScan).Error
// 				if err != nil {
// 					return nil, err
// 				}
// 			} else {
// 				return nil, errors.New("type not found")
// 			}
// 			itemToScan.PlanQuantity = items.PlanQuantity
// 			itemToScan.Time = inventarization.Time
// 			itemToScan.ItemID = items.ItemID
// 			itemToScan.Type = items.Type
// 			itemToScan.Status = utils.StatusOpened
// 			itemToScan.SkladID = inventarization.SkladID
// 			itemToScan.StartQuantity = invItem.FactQuantity
// 			itemToScan.Difference = itemToScan.FactQuantity - itemToScan.PlanQuantity
// 			itemToScan.Expenses = itemToScan.StartQuantity - itemToScan.PlanQuantity + itemToScan.Income - itemToScan.Removed
// 			itemToScan.Before = invItem.ID
// 			itemToScan.BeforeTime = invItem.Time
// 			invItems = append(invItems, itemToScan)
// 		}
// 	}
// 	inventarization.InventarizationItems = invItems
// 	return inventarization, nil
// }

func (r *SkladDB) GetOpenInventarization(inventarization *model.Inventarization) (*model.Inventarization, error) {
	invItems := []*model.InventarizationItem{}
	for _, item := range inventarization.InventarizationItems {
		oldItem, err := r.GetBeforeInventarizationItem(item)
		if err != nil {
			if err != gorm.ErrRecordNotFound {
				return nil, err
			}
		}
		if oldItem == nil {
			oldItem = &model.InventarizationItem{}
		}
		item.BeforeTime = oldItem.Time
		item.Status = utils.StatusNew
		itemToScan, err := r.CalculateItemToScan(item)
		if err != nil {
			return nil, err
		}
		itemToScan.Type = item.Type
		if itemToScan.Type != utils.TypeGroup {
			itemToScan.BeforeTime = oldItem.Time
			itemToScan.Time = inventarization.Time
			itemToScan.ItemID = item.ItemID
			itemToScan.Type = item.Type
			itemToScan.Status = utils.StatusNew
			itemToScan.SkladID = inventarization.SkladID
			itemToScan.StartQuantity = oldItem.FactQuantity
			itemToScan.PlanQuantity = itemToScan.StartQuantity + itemToScan.Income - itemToScan.Removed - itemToScan.Expenses
			itemToScan.Difference = itemToScan.FactQuantity - itemToScan.PlanQuantity
			itemToScan.DifferenceSum = itemToScan.Difference * itemToScan.Cost
			itemToScan.LoadingStatus = utils.StatusLoaded
			itemToScan.IsVisible = item.IsVisible
			itemToScan.NeedToRecalculate = item.NeedToRecalculate
			itemToScan.GroupID = item.GroupID
			invItems = append(invItems, itemToScan)
		} else {
			itemToScan.BeforeTime = oldItem.Time
			itemToScan.Time = inventarization.Time
			itemToScan.ItemID = item.ItemID
			itemToScan.Type = item.Type
			itemToScan.Status = utils.StatusNew
			itemToScan.SkladID = inventarization.SkladID
			itemToScan.LoadingStatus = utils.StatusLoaded
			itemToScan.IsVisible = item.IsVisible
			itemToScan.NeedToRecalculate = item.NeedToRecalculate
			invItems = append(invItems, itemToScan)
		}
	}
	inventarization.InventarizationItems = invItems
	return inventarization, nil
}

func (r *SkladDB) GetToCreateInventratization(inventarization *model.Inventarization) (*model.Inventarization, error) {
	//invItems := []*model.InventarizationItem{}
	err := r.gormDB.Transaction(func(tx *gorm.DB) error {
		if inventarization.Type == utils.TypeFull {
			items := []*model.InventarizationItem{}
			ingre := []*model.Ingredient{}
			ingredientMap := make(map[int]bool, 100)
			tovarMap := make(map[int]bool, 100)
			groups := []*model.InventarizationGroup{}
			if err := r.gormDB.Model(groups).Select("*").Where("sklad_id = ?", inventarization.SkladID).Scan(&groups).Error; err != nil {
				return err
			}

			for _, group := range groups {
				item := &model.InventarizationItem{
					ItemID:    group.ID,
					Type:      utils.TypeGroup,
					Status:    utils.StatusNew,
					SkladID:   inventarization.SkladID,
					Time:      inventarization.Time,
					IsVisible: true,
				}
				if item.Time.After(time.Now()) {
					item.NeedToRecalculate = true
				}
				items = append(items, item)
				groupItems := []*model.InventarizationGroupItem{}
				if err := r.gormDB.Model(groupItems).Select("*").Where("group_id = ?", group.ID).Scan(&groupItems).Error; err != nil {
					return err
				}
				for _, groupItem := range groupItems {
					item := &model.InventarizationItem{
						ItemID:    groupItem.ItemID,
						Type:      utils.TypeTovar,
						Status:    utils.StatusNew,
						SkladID:   inventarization.SkladID,
						Time:      inventarization.Time,
						IsVisible: false,
						GroupID:   group.ID,
					}
					if group.Type == utils.TypeTovar {
						tovarMap[groupItem.ItemID] = true
					} else {
						ingredientMap[groupItem.ItemID] = true
						item.Type = utils.TypeIngredient
					}
					if item.Time.After(time.Now()) {
						item.NeedToRecalculate = true
					}
					items = append(items, item)
				}

			}
			var shopID int
			if err := r.gormDB.Model(&model.Sklad{}).Select("shop_id").Where("id = ?", inventarization.SkladID).Scan(&shopID).Error; err != nil {
				return err
			}
			err := tx.Model(&model.Ingredient{}).Where("deleted = ? and shop_id = ?", false, shopID).Order("category asc, id desc").Scan(&ingre).Error
			if err != nil {
				return err
			}
			for _, ing := range ingre {
				if _, ok := ingredientMap[ing.IngredientID]; ok {
					continue
				}
				item := &model.InventarizationItem{
					ItemID:    ing.IngredientID,
					Type:      utils.TypeIngredient,
					Status:    utils.StatusNew,
					SkladID:   inventarization.SkladID,
					Time:      inventarization.Time,
					IsVisible: true,
				}
				if item.Time.After(time.Now()) {
					item.NeedToRecalculate = true
				}
				items = append(items, item)
			}
			tovars := []*model.Tovar{}
			err = tx.Model(&model.Tovar{}).Where("deleted = ? and shop_id = ?", false, shopID).Order("category asc, id desc").Scan(&tovars).Error
			if err != nil {
				return err
			}
			for _, tovar := range tovars {
				if _, ok := tovarMap[tovar.TovarID]; ok {
					continue
				}
				item := &model.InventarizationItem{
					ItemID:    tovar.TovarID,
					Type:      utils.TypeTovar,
					Status:    utils.StatusNew,
					SkladID:   inventarization.SkladID,
					Time:      inventarization.Time,
					IsVisible: true,
				}
				if item.Time.After(time.Now()) {
					item.NeedToRecalculate = true
				}
				items = append(items, item)
			}
			inventarization.InventarizationItems = items
		} else if inventarization.Type == utils.TypeFullPartial {
			items := []*model.InventarizationItem{}
			ingre := []*model.Ingredient{}
			ingredientMap := make(map[int]bool, 100)
			tovarMap := make(map[int]bool, 100)
			groups := []*model.InventarizationGroup{}
			if err := r.gormDB.Model(groups).Select("*").Where("sklad_id = ?", inventarization.SkladID).Scan(&groups).Error; err != nil {
				return err
			}

			for _, group := range groups {
				item := &model.InventarizationItem{
					ItemID:    group.ID,
					Type:      utils.TypeGroup,
					Status:    utils.StatusNew,
					SkladID:   inventarization.SkladID,
					Time:      inventarization.Time,
					IsVisible: true,
				}
				if item.Time.After(time.Now()) {
					item.NeedToRecalculate = true
				}
				items = append(items, item)
				groupItems := []*model.InventarizationGroupItem{}
				if err := r.gormDB.Model(groupItems).Select("*").Where("group_id = ?", group.ID).Scan(&groupItems).Error; err != nil {
					return err
				}
				for _, groupItem := range groupItems {
					item := &model.InventarizationItem{
						ItemID:    groupItem.ItemID,
						Type:      utils.TypeTovar,
						Status:    utils.StatusNew,
						SkladID:   inventarization.SkladID,
						Time:      inventarization.Time,
						IsVisible: false,
						GroupID:   group.ID,
					}
					if group.Type == utils.TypeTovar {
						tovarMap[groupItem.ItemID] = true
					} else {
						ingredientMap[groupItem.ItemID] = true
						item.Type = utils.TypeIngredient
					}
					if item.Time.After(time.Now()) {
						item.NeedToRecalculate = true
					}
					items = append(items, item)
				}

			}
			var shopID int
			if err := r.gormDB.Model(&model.Sklad{}).Select("shop_id").Where("id = ?", inventarization.SkladID).Scan(&shopID).Error; err != nil {
				return err
			}
			err := tx.Model(&model.Ingredient{}).Where("deleted = ? and shop_id = ? and category != ?", false, shopID, 103).Order("category asc, id desc").Scan(&ingre).Error
			if err != nil {
				return err
			}
			for _, ing := range ingre {
				if _, ok := ingredientMap[ing.IngredientID]; ok {
					continue
				}
				item := &model.InventarizationItem{
					ItemID:    ing.IngredientID,
					Type:      utils.TypeIngredient,
					Status:    utils.StatusNew,
					SkladID:   inventarization.SkladID,
					Time:      inventarization.Time,
					IsVisible: true,
				}
				if item.Time.After(time.Now()) {
					item.NeedToRecalculate = true
				}
				items = append(items, item)
			}
			tovars := []*model.Tovar{}
			err = tx.Model(&model.Tovar{}).Where("deleted = ? and shop_id = ?", false, shopID).Order("category asc, id desc").Scan(&tovars).Error
			if err != nil {
				return err
			}
			for _, tovar := range tovars {
				if _, ok := tovarMap[tovar.TovarID]; ok {
					continue
				}
				item := &model.InventarizationItem{
					ItemID:    tovar.TovarID,
					Type:      utils.TypeTovar,
					Status:    utils.StatusNew,
					SkladID:   inventarization.SkladID,
					Time:      inventarization.Time,
					IsVisible: true,
				}
				if item.Time.After(time.Now()) {
					item.NeedToRecalculate = true
				}
				items = append(items, item)
			}
			inventarization.InventarizationItems = items
		}

		inventarization, err := r.GetOpenInventarization(inventarization)
		if err != nil {
			return err
		}
		err = tx.Model(inventarization).Create(inventarization).Error
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return inventarization, nil
}

func (r *SkladDB) CalculateItemToScan(item *model.InventarizationItem) (*model.InventarizationItem, error) {
	itemToScan := &model.InventarizationItem{}
	log.Print(item)
	if item.Type == utils.TypeIngredient {
		if item.Status == utils.StatusNew || item.Cost == 0 {
			err := r.gormDB.Model(&model.SkladIngredient{}).Select("sklad_ingredients.cost").Where("sklad_id = ? and ingredient_id = ?", item.SkladID, item.ItemID).Scan(&itemToScan).Error
			if err != nil { //index
				return nil, err
			}
		}
		err := r.gormDB.Model(&model.Postavka{}).Select("SUM(item_postavkas.quantity) as income").Joins("LEFT JOIN item_postavkas ON item_postavkas.postavka_id = postavkas.id").Where("postavkas.time < ? AND postavkas.time > ? AND item_postavkas.item_id = ? AND item_postavkas.type = ? AND postavkas.sklad_id = ? AND item_postavkas.deleted = ?", item.Time, item.BeforeTime, item.ItemID, utils.TypeIngredient, item.SkladID, false).Scan(&itemToScan).Error
		if err != nil {
			return nil, err
		}

		err = r.gormDB.Model(&model.RemoveFromSklad{}).Select("SUM(remove_from_sklad_items.quantity) as removed").Joins("LEFT JOIN remove_from_sklad_items ON remove_from_sklad_items.remove_id = remove_from_sklads.id").Where("remove_from_sklads.time < ? AND remove_from_sklads.time > ? AND remove_from_sklad_items.item_id = ? AND remove_from_sklad_items.type = ? AND remove_from_sklads.sklad_id = ?", item.Time, item.BeforeTime, item.ItemID, utils.TypeIngredient, item.SkladID).Scan(&itemToScan).Error
		if err != nil {
			return nil, err
		}

		itemToScanExpenses := &model.InventarizationItem{}
		err = r.gormDB.Table("expence_ingredients").Select("SUM(expence_ingredients.quantity) as expenses").Where("expence_ingredients.ingredient_id = ? AND expence_ingredients.time between ? AND ? and sklad_id = ? and expence_ingredients.status = ?", item.ItemID, item.BeforeTime, item.Time, item.SkladID, utils.StatusClosed).Scan(&itemToScanExpenses).Error
		if err != nil {
			return nil, err
		}
		itemToScan.Expenses += itemToScanExpenses.Expenses

	} else if item.Type == utils.TypeTovar {
		if item.Status == utils.StatusNew {
			err := r.gormDB.Model(&model.SkladTovar{}).Select("sklad_tovars.cost").Where("sklad_id = ? and tovar_id = ?", item.SkladID, item.ItemID).Scan(&itemToScan).Error
			if err != nil {
				return nil, err
			}
		}
		err := r.gormDB.Model(&model.Postavka{}).Select("SUM(item_postavkas.quantity) as income").Joins("LEFT JOIN item_postavkas ON item_postavkas.postavka_id = postavkas.id").Where("postavkas.time < ? AND postavkas.time > ? AND item_postavkas.item_id = ? AND item_postavkas.type = ? AND postavkas.sklad_id = ? and item_postavkas.deleted = ?", item.Time, item.BeforeTime, item.ItemID, utils.TypeTovar, item.SkladID, false).Scan(&itemToScan).Error
		if err != nil {
			return nil, err
		}
		err = r.gormDB.Model(&model.RemoveFromSklad{}).Select("SUM(remove_from_sklad_items.quantity) as removed").Joins("LEFT JOIN remove_from_sklad_items ON remove_from_sklad_items.remove_id = remove_from_sklads.id").Where("remove_from_sklads.time < ? AND remove_from_sklads.time > ? AND remove_from_sklad_items.item_id = ? AND remove_from_sklad_items.type = ? AND remove_from_sklads.sklad_id = ?", item.Time, item.BeforeTime, item.ItemID, utils.TypeTovar, item.SkladID).Scan(&itemToScan).Error
		if err != nil {
			return nil, err
		}

		itemToScanExpenses := &model.InventarizationItem{}
		err = r.gormDB.Table("expence_tovars").Select("SUM(expence_tovars.quantity) as expenses").Where("expence_tovars.tovar_id = ? AND expence_tovars.time between ? AND ? and sklad_id = ?  AND expence_tovars.status = ?", item.ItemID, item.BeforeTime, item.Time, item.SkladID, utils.StatusClosed).Scan(&itemToScanExpenses).Error
		if err != nil {
			return nil, err
		}
		itemToScan.Expenses += itemToScanExpenses.Expenses
	} else if item.Type == utils.TypeGroup {

	} else {
		return nil, errors.New("type not found")
	}
	itemToScan.RemovedSum = itemToScan.RemovedSum * itemToScan.Cost

	return itemToScan, nil
}

func (r *SkladDB) GetShopBySkladID(id int) (*model.Shop, error) {
	sklad := &model.Sklad{}
	err := r.gormDB.Where("id = ?", id).First(sklad).Error
	if err != nil {
		return nil, err
	}
	shop := &model.Shop{}
	err = r.gormDB.Where("id = ?", sklad.ShopID).First(shop).Error
	if err != nil {
		return nil, err
	}
	return shop, nil
}

func (r *SkladDB) CheckInventItems(items []*model.InventarizationItem, shopID int) error {
	for _, item := range items {
		if item.Type == utils.TypeTovar {
			tovar := &model.Tovar{}
			err := r.gormDB.Where("tovar_id = ? and shop_id = ?", item.ItemID, shopID).First(tovar).Error
			if err != nil {
				if err == gorm.ErrRecordNotFound {
					return errors.New("tovar not found")
				}
				return err
			}
			if tovar.ID == 0 {
				return errors.New("tovar not found")
			}
		} else if item.Type == utils.TypeIngredient {
			ingredient := &model.Ingredient{}
			err := r.gormDB.Where("ingredient_id = ? and shop_id = ?", item.ItemID, shopID).First(ingredient).Error
			if err != nil {
				if err == gorm.ErrRecordNotFound {
					return errors.New("ingredient not found")
				}
				return err
			}
			if ingredient.ID == 0 {
				return errors.New("ingredient not found")
			}
		} else {
			inventGroup := &model.InventarizationGroup{}
			err := r.gormDB.Where("id = ?", item.ItemID).First(inventGroup).Error
			if err != nil {
				return err
			}
			if inventGroup.SkladID != item.SkladID {
				return errors.New("group not found")
			}
		}
	}
	return nil
}

func (r *SkladDB) RecalculateOneInventarization(item *model.InventarizationItem) error {
	item.LoadingStatus = utils.StatusLoading
	err := r.gormDB.Model(item).Update("loading_status", item.LoadingStatus).Error
	if err != nil {
		return err
	}
	err = r.gormDB.Transaction(func(tx *gorm.DB) error {
		oldItem, err := r.GetBeforeInventarizationItem(item)
		if err != nil {
			if err != gorm.ErrRecordNotFound {
				return err
			}
		}
		if oldItem == nil {
			oldItem = &model.InventarizationItem{}
		}
		item.BeforeTime = oldItem.Time
		itemToScan, err := r.CalculateItemToScan(item)
		if err != nil {
			return err
		}
		if item.Status != utils.StatusNew {
			itemToScan.Cost = item.Cost
		}
		itemToScan.InventarizationID = item.InventarizationID
		itemToScan.Time = item.Time
		itemToScan.ItemID = item.ItemID
		itemToScan.Type = item.Type
		itemToScan.Status = item.Status
		itemToScan.SkladID = item.SkladID
		itemToScan.StartQuantity = oldItem.FactQuantity
		itemToScan.FactQuantity = item.FactQuantity
		itemToScan.PlanQuantity = itemToScan.StartQuantity + itemToScan.Income - itemToScan.Removed - itemToScan.Expenses
		itemToScan.Difference = itemToScan.FactQuantity - itemToScan.PlanQuantity
		itemToScan.DifferenceSum = itemToScan.Difference * itemToScan.Cost
		itemToScan.BeforeTime = oldItem.Time
		itemToScan.LoadingStatus = utils.StatusLoaded
		itemToScan.NeedToRecalculate = item.NeedToRecalculate
		err = tx.Model(&model.InventarizationItem{}).Where("inventarization_items.id = ?", item.ID).Select("before_time", "cost", "difference_sum", "difference", "fact_quantity", "plan_quantity", "removed_sum", "removed", "start_quantity", "income", "expenses", "loading_status").Updates(itemToScan).Error
		if err != nil {
			return err
		}
		invent := model.Inventarization{}
		err = tx.Model(&model.Inventarization{}).Where("inventarizations.id = ?", item.InventarizationID).Scan(&invent).Error

		if err != nil {
			return err
		}

		invent.Result += itemToScan.DifferenceSum - item.DifferenceSum

		err = tx.Model(&model.Inventarization{}).Where("inventarizations.id = ?", item.InventarizationID).Update("result", invent.Result).Error
		if err != nil {
			return err
		}
		return nil
	})
	return err
}

func (r *SkladDB) RecalculateInventarization(invItems []*model.InventarizationItem) error {
	for _, item := range invItems {
		item.LoadingStatus = utils.StatusLoading
		err := r.gormDB.Model(item).Update("loading_status", item.LoadingStatus).Error
		if err != nil {
			errorCheck := &model.ErrorCheck{
				Error:   err.Error(),
				Request: "RecalculateInventarization" + fmt.Sprintf("%v", invItems),
				Time:    time.Now(),
			}
			_ = r.gormDB.Model(&model.ErrorCheck{}).Create(errorCheck).Error

			return err
		}
	}
	for _, item := range invItems {
		itemToScan := &model.InventarizationItem{}
		err := r.gormDB.Transaction(func(tx *gorm.DB) error {
			oldItem, err := r.GetBeforeInventarizationItem(item)
			if err != nil {
				if err != gorm.ErrRecordNotFound {
					errorCheck := &model.ErrorCheck{
						Error:   err.Error(),
						Request: "RecalculateInventarization" + fmt.Sprintf("%v", invItems),
						Time:    time.Now(),
					}
					_ = r.gormDB.Model(&model.ErrorCheck{}).Create(errorCheck).Error

					return err
				}
			}
			if oldItem == nil {
				oldItem = &model.InventarizationItem{}
			}
			item.BeforeTime = oldItem.Time
			itemToScan, err = r.CalculateItemToScan(item)
			if err != nil {
				errorCheck := &model.ErrorCheck{
					Error:   err.Error(),
					Request: "RecalculateInventarization" + fmt.Sprintf("%v", invItems),
					Time:    time.Now(),
				}
				_ = r.gormDB.Model(&model.ErrorCheck{}).Create(errorCheck).Error

				return err
			}
			itemToScan.InventarizationID = item.InventarizationID
			if itemToScan.Cost == 0 {
				itemToScan.Cost = item.Cost
			}
			itemToScan.Time = item.Time
			itemToScan.ItemID = item.ItemID
			itemToScan.Type = item.Type
			itemToScan.Status = item.Status
			itemToScan.SkladID = item.SkladID
			itemToScan.StartQuantity = oldItem.FactQuantity
			itemToScan.FactQuantity = item.FactQuantity
			itemToScan.PlanQuantity = itemToScan.StartQuantity + itemToScan.Income - itemToScan.Removed - itemToScan.Expenses
			itemToScan.Difference = itemToScan.FactQuantity - itemToScan.PlanQuantity
			itemToScan.DifferenceSum = itemToScan.Difference * itemToScan.Cost
			itemToScan.BeforeTime = oldItem.Time
			itemToScan.LoadingStatus = utils.StatusLoaded
			itemToScan.NeedToRecalculate = item.NeedToRecalculate
			err = tx.Model(&model.InventarizationItem{}).Where("inventarization_items.id = ?", item.ID).Select("before_time", "difference_sum", "difference", "fact_quantity", "plan_quantity", "removed_sum", "removed", "start_quantity", "income", "expenses", "loading_status", "cost").Debug().Updates(itemToScan).Error
			if err != nil {
				errorCheck := &model.ErrorCheck{
					Error:   err.Error(),
					Request: "RecalculateInventarization" + fmt.Sprintf("%v", invItems),
					Time:    time.Now(),
				}
				_ = r.gormDB.Model(&model.ErrorCheck{}).Create(errorCheck).Error

				return err
			}
			invent := model.Inventarization{}
			err = tx.Model(&model.Inventarization{}).Where("inventarizations.id = ?", item.InventarizationID).Scan(&invent).Error

			if err != nil {
				errorCheck := &model.ErrorCheck{
					Error:   err.Error(),
					Request: "RecalculateInventarization" + fmt.Sprintf("%v", invItems),
					Time:    time.Now(),
				}
				_ = r.gormDB.Model(&model.ErrorCheck{}).Create(errorCheck).Error

				return err
			}

			invent.Result += itemToScan.DifferenceSum - item.DifferenceSum

			err = tx.Model(&model.Inventarization{}).Where("inventarizations.id = ?", item.InventarizationID).Update("result", invent.Result).Error
			if err != nil {
				errorCheck := &model.ErrorCheck{
					Error:   err.Error(),
					Request: "RecalculateInventarization" + fmt.Sprintf("%v", invItems),
					Time:    time.Now(),
				}
				_ = r.gormDB.Model(&model.ErrorCheck{}).Create(errorCheck).Error

				return err
			}
			return nil

		})
		if err != nil {
			errorCheck := &model.ErrorCheck{
				Error:   err.Error(),
				Request: "RecalculateInventarization" + fmt.Sprintf("%v", invItems),
				Time:    time.Now(),
			}
			_ = r.gormDB.Model(&model.ErrorCheck{}).Create(errorCheck).Error

			return err
		}
		afterItem, err := r.GetAfterInventarizationItem(itemToScan)
		if err != nil {
			if err != gorm.ErrRecordNotFound {
				errorCheck := &model.ErrorCheck{
					Error:   err.Error(),
					Request: "RecalculateInventarization" + fmt.Sprintf("%v", invItems),
					Time:    time.Now(),
				}
				_ = r.gormDB.Model(&model.ErrorCheck{}).Create(errorCheck).Error

				return err
			}
		}
		if afterItem != nil {
			err = r.RecalculateOneInventarization(afterItem)
			if err != nil {
				errorCheck := &model.ErrorCheck{
					Error:   err.Error(),
					Request: "RecalculateInventarization" + fmt.Sprintf("%v", invItems),
					Time:    time.Now(),
				}
				_ = r.gormDB.Model(&model.ErrorCheck{}).Create(errorCheck).Error

				return err
			}
		}
	}

	return nil
}

func (r *SkladDB) RecalculateInventarizations() error {
	items := []*model.InventarizationItem{}
	err := r.gormDB.Model(&model.InventarizationItem{}).Order("time asc").Scan(&items).Error
	if err != nil {
		return err
	}
	go r.RecalculateInventarization(items)
	return nil
}

func (r *SkladDB) GetBeforeInventarizationItem(inventarizationItem *model.InventarizationItem) (*model.InventarizationItem, error) {
	item := &model.InventarizationItem{}
	res := r.gormDB.Model(&model.InventarizationItem{}).Where("item_id = ? AND type = ? AND sklad_id = ? AND time < ? AND id != ?", inventarizationItem.ItemID, inventarizationItem.Type, inventarizationItem.SkladID, inventarizationItem.Time, inventarizationItem.ID).Order("time desc").Limit(1).Scan(item)
	if res.Error != nil {
		return nil, res.Error
	}
	if res.RowsAffected == 0 {
		return nil, gorm.ErrRecordNotFound
	}
	return item, nil
}

func (r *SkladDB) GetAfterInventarizationItem(inventarizationItem *model.InventarizationItem) (*model.InventarizationItem, error) {
	item := &model.InventarizationItem{}
	res := r.gormDB.Model(&model.InventarizationItem{}).Where("item_id = ? AND type = ? AND sklad_id = ? AND time > ? and id != ?", inventarizationItem.ItemID, inventarizationItem.Type, inventarizationItem.SkladID, inventarizationItem.Time, inventarizationItem.ID).Order("time asc").Limit(1).Scan(item)
	if res.Error != nil {
		return nil, res.Error
	}
	if res.RowsAffected == 0 {
		return nil, gorm.ErrRecordNotFound
	}
	return item, nil
}

func (r *SkladDB) UpdateInventarizationParams(inventarization *model.Inventarization) (*model.Inventarization, error) {
	itemsToRecalculate := []*model.InventarizationItem{}
	itemsToRecalculateAfter := []*model.InventarizationItem{}
	trafficItems := []*model.AsyncJob{}

	err := r.gormDB.Transaction(func(tx *gorm.DB) error {
		inventarization.Result = 0
		for i := 0; i < len(inventarization.InventarizationItems); i++ {
			invItem := inventarization.InventarizationItems[i]
			oldItem := &model.InventarizationItem{}
			itemAddedToRecalculate := false
			res := tx.Model(oldItem).Where("inventarization_items.id = ?", invItem.ID).First(oldItem)
			err := res.Error
			if err != nil {
				if err != gorm.ErrRecordNotFound {
					return err
				}
			}

			if oldItem.LoadingStatus == utils.StatusLoading {
				return errors.New("inventarization is loading")
			}
			if res.RowsAffected != 0 {
				if invItem.HasChanged(oldItem) { // if nothing is changed no need to recalculate
					if !oldItem.Time.Equal(invItem.Time) {
						afterItem, err := r.GetAfterInventarizationItem(oldItem)
						if err != nil {
							if err != gorm.ErrRecordNotFound {
								return err
							}
						}
						if afterItem != nil {
							itemsToRecalculateAfter = append(itemsToRecalculateAfter, afterItem)
						}
						oldItem.Time = invItem.Time
						if oldItem.Time.After(time.Now()) {
							oldItem.NeedToRecalculate = true
						}
						newAfterItem, err := r.GetAfterInventarizationItem(oldItem)
						if err != nil {
							if err != gorm.ErrRecordNotFound {
								return err
							}
						}
						if newAfterItem != nil {
							if newAfterItem.ID != afterItem.ID {
								itemsToRecalculateAfter = append(itemsToRecalculateAfter, newAfterItem)
							}
						}
						itemsToRecalculate = append(itemsToRecalculate, oldItem)
						itemAddedToRecalculate = true
					}
					oldItem.IsVisible = invItem.IsVisible
					oldItem.GroupID = invItem.GroupID
					if oldItem.NeedToRecalculate {
						if oldItem.Time.Before(time.Now()) {
							oldItem.NeedToRecalculate = false
						}
						if !itemAddedToRecalculate {
							itemsToRecalculate = append(itemsToRecalculate, oldItem)
						}
					}
					inventarization.InventarizationItems[i] = oldItem

					err = tx.Model(oldItem).Save(oldItem).Error
					if err != nil {
						return err
					}
				}
			} else {
				if invItem.GroupID == 0 {
					invItem.IsVisible = true
				}
				invItem.LoadingStatus = utils.StatusLoading
				invItem.Status = utils.StatusNew
				err = r.gormDB.Model(invItem).Create(invItem).Error //creates new Item if new item is added
				if err != nil {
					return err
				}
				itemsToRecalculate = append(itemsToRecalculate, invItem)
			}
			inventarization.Result += oldItem.DifferenceSum
			if invItem.Time.Year() != time.Now().Year() || invItem.Time.Month() != time.Now().Month() || invItem.Time.Day() != time.Now().Day() {
				trafficItems = append(trafficItems, &model.AsyncJob{
					ItemID:    invItem.ItemID,
					ItemType:  invItem.Type,
					SkladID:   invItem.SkladID,
					TimeStamp: invItem.Time.UTC(),
					CreatedAt: time.Now(),
					Status:    utils.StatusNeedRecalculate,
				})
			}
		}
		err := tx.Model(inventarization).Select("time", "status", "result", "loading_status").Updates(inventarization).Error
		if err != nil {
			return err
		}
		return nil
	},
	)
	if err != nil {
		return nil, err
	}
	r.ConcurrentRecalculationForDailyStatistics(trafficItems)
	r.ConcurrentRecalculation(itemsToRecalculate)
	r.ConcurrentRecalculation(itemsToRecalculateAfter)
	return inventarization, nil
}

func (r *SkladDB) ConcurrentRecalculation(itemsToRecalculate []*model.InventarizationItem) {
	wait := sync.WaitGroup{}
	myChannel := make(chan *model.InventarizationItem, len(itemsToRecalculate))
	ctx, cancel := context.WithCancel(context.Background())
	wait.Add(6)
	for i := 0; i < 6; i++ {
		go r.GoRecalcReader(&wait, myChannel, ctx)
	}

	for _, item := range itemsToRecalculate {
		log.Print("recaltulate ", item)
		myChannel <- item
	}

	close(myChannel)
	wait.Wait()
	cancel()
}

func (r *SkladDB) UpdateInventarization(inventarization *model.Inventarization) (*model.Inventarization, error) {
	itemsToRecalculate := []*model.InventarizationItem{}
	trafficItems := []*model.AsyncJob{}
	err := r.gormDB.Transaction(func(tx *gorm.DB) error {
		inventarization.Result = 0
		for i := 0; i < len(inventarization.InventarizationItems); i++ {
			invItem := inventarization.InventarizationItems[i]
			oldItem := &model.InventarizationItem{}
			res := tx.Model(oldItem).Where("inventarization_items.inventarization_id = ? AND inventarization_items.item_id = ? and inventarization_items.type = ?", inventarization.ID, invItem.ItemID, invItem.Type).First(oldItem)
			err := res.Error
			if err != nil {
				if err != gorm.ErrRecordNotFound {
					return err
				}
			}
			if oldItem.LoadingStatus == utils.StatusLoading {
				return errors.New("inventarization is loading")
			}
			if res.RowsAffected != 0 {
				if invItem.HasChanged(oldItem) || invItem.FactQuantity == 0 { // if nothing is changed no need to recalculate
					factQuantityNotEqual := false
					if oldItem.FactQuantity != invItem.FactQuantity { // if fact quantity is changed then we need to recalculate
						oldItem.FactQuantity = invItem.FactQuantity
						oldItem.Difference = oldItem.FactQuantity - oldItem.PlanQuantity
						oldItem.DifferenceSum = oldItem.Difference * oldItem.Cost
						factQuantityNotEqual = true
					}
					oldItem.GroupID = invItem.GroupID
					oldStatus := oldItem.Status
					oldItem.Status = invItem.Status
					if !oldItem.Time.Equal(invItem.Time) || oldItem.NeedToRecalculate {
						afterItem, err := r.GetAfterInventarizationItem(oldItem)
						if err != nil {
							if err != gorm.ErrRecordNotFound {
								return err
							}
						}
						if afterItem != nil {
							itemsToRecalculate = append(itemsToRecalculate, afterItem)
						}
						oldItem.Time = invItem.Time
						if oldItem.NeedToRecalculate {
							if oldItem.Time.Before(time.Now()) {
								oldItem.NeedToRecalculate = false
							}
						}
						itemsToRecalculate = append(itemsToRecalculate, oldItem)
					}
					inventarization.InventarizationItems[i] = oldItem

					err = tx.Model(oldItem).Save(oldItem).Error
					if err != nil {
						return err
					}
					inventarization.Result += oldItem.DifferenceSum

					afterItem, err := r.GetAfterInventarizationItem(oldItem)

					if err != nil {
						if err != gorm.ErrRecordNotFound {
							return err
						}
					}

					if afterItem != nil {
						if oldStatus == utils.StatusNew {
							itemsToRecalculate = append(itemsToRecalculate, afterItem)
						} else if factQuantityNotEqual {
							afterItem.StartQuantity = invItem.FactQuantity
							afterItem.PlanQuantity = afterItem.StartQuantity + afterItem.Income - afterItem.Removed - afterItem.Expenses
							afterItem.Difference = afterItem.FactQuantity - afterItem.PlanQuantity
							oldDifferenceSum := afterItem.DifferenceSum
							afterItem.DifferenceSum = afterItem.Difference * afterItem.Cost
							err = tx.Model(afterItem).Save(afterItem).Error
							if err != nil {
								return err
							}
							afterInventarization := &model.Inventarization{}
							err = tx.Model(afterInventarization).Where("inventarizations.id = ?", afterItem.InventarizationID).First(afterInventarization).Error
							if err != nil {
								return err
							}
							afterInventarization.Result += afterItem.DifferenceSum - oldDifferenceSum
							err = tx.Model(afterInventarization).Save(afterInventarization).Error
							if err != nil {
								return err
							}
						}
					}
					if afterItem == nil {
						itemToScan := &model.InventarizationItem{}
						if oldItem.Status == utils.StatusClosed {
							if oldItem.Type == utils.TypeIngredient {
								err := r.gormDB.Model(&model.Postavka{}).Select("SUM(item_postavkas.quantity) as income").Joins("LEFT JOIN item_postavkas ON item_postavkas.postavka_id = postavkas.id").Where("postavkas.time < ? AND postavkas.time > ? AND item_postavkas.item_id = ? AND item_postavkas.type = ? AND postavkas.sklad_id = ?", time.Now(), oldItem.Time, oldItem.ItemID, utils.TypeIngredient, oldItem.SkladID).Scan(&itemToScan).Error
								if err != nil {
									return err
								}
								err = r.gormDB.Model(&model.RemoveFromSklad{}).Select("SUM(remove_from_sklad_items.quantity) as removed, SUM(sklad_ingredients.cost * remove_from_sklad_items.quantity) as removed_sum").Joins("LEFT JOIN remove_from_sklad_items ON remove_from_sklad_items.remove_id = remove_from_sklads.id").Joins("LEFT JOIN sklad_ingredients ON sklad_ingredients.ingredient_id = remove_from_sklad_items.item_id").Where("remove_from_sklads.time < ? AND remove_from_sklads.time > ? AND remove_from_sklad_items.item_id = ? AND remove_from_sklad_items.type = ? AND remove_from_sklads.sklad_id = ? AND sklad_ingredients.sklad_id = ?", time.Now(), oldItem.Time, oldItem.ItemID, utils.TypeIngredient, oldItem.SkladID, oldItem.SkladID).Scan(&itemToScan).Error
								if err != nil {
									return err
								}
								itemToScanExpenses := &model.InventarizationItem{}
								err = r.gormDB.Table("expence_ingredients").Select("SUM(expence_ingredients.quantity) as expenses").Where("expence_ingredients.ingredient_id = ? AND expence_ingredients.time between ? AND ? and sklad_id = ? and expence_ingredients.status = ?", oldItem.ItemID, oldItem.Time, time.Now(), oldItem.SkladID, utils.StatusClosed).Scan(&itemToScanExpenses).Error
								if err != nil {
									return err
								}
								itemToScan.Expenses += itemToScanExpenses.Expenses
								factQuantity := oldItem.FactQuantity + itemToScan.Income - itemToScan.Removed - itemToScan.Expenses

								skladIngredient := &model.SkladIngredient{}
								err = tx.Model(skladIngredient).Where("sklad_ingredients.sklad_id = ? AND sklad_ingredients.ingredient_id = ?", oldItem.SkladID, oldItem.ItemID).Update("quantity", factQuantity).First(skladIngredient).Error
								if err != nil {
									return err
								}
							} else if oldItem.Type == utils.TypeTovar {
								err := r.gormDB.Model(&model.Postavka{}).Select("SUM(item_postavkas.quantity) as income").Joins("LEFT JOIN item_postavkas ON item_postavkas.postavka_id = postavkas.id").Where("postavkas.time < ? AND postavkas.time > ? AND item_postavkas.item_id = ? AND item_postavkas.type = ? AND postavkas.sklad_id = ?", time.Now(), oldItem.Time, oldItem.ItemID, utils.TypeTovar, oldItem.SkladID).Scan(&itemToScan).Error
								if err != nil {
									return err
								}
								err = r.gormDB.Model(&model.RemoveFromSklad{}).Select("SUM(remove_from_sklad_items.quantity) as removed, SUM(sklad_tovars.cost * remove_from_sklad_items.quantity) as removed_sum").Joins("LEFT JOIN remove_from_sklad_items ON remove_from_sklad_items.remove_id = remove_from_sklads.id").Joins("LEFT JOIN sklad_tovars ON sklad_tovars.tovar_id = remove_from_sklad_items.item_id").Where("remove_from_sklads.time < ? AND remove_from_sklads.time > ? AND remove_from_sklad_items.item_id = ? AND remove_from_sklad_items.type = ? AND remove_from_sklads.sklad_id = ? AND sklad_tovars.sklad_id = ?", time.Now(), oldItem.Time, oldItem.ItemID, utils.TypeTovar, oldItem.SkladID, oldItem.SkladID).Scan(&itemToScan).Error
								if err != nil {
									return err
								}
								itemToScanExpenses := &model.InventarizationItem{}
								err = r.gormDB.Table("expence_tovars").Select("SUM(expence_tovars.quantity) as expenses").Where("expence_tovars.tovar_id = ? AND expence_tovars.time between ? AND ? and sklad_id = ?  AND expence_tovars.status = ?", oldItem.ItemID, oldItem.Time, time.Now(), oldItem.SkladID, utils.StatusClosed).Scan(&itemToScanExpenses).Error
								if err != nil {
									return err
								}
								itemToScan.Expenses += itemToScanExpenses.Expenses

								factQuantity := oldItem.FactQuantity + itemToScan.Income - itemToScan.Removed - itemToScan.Expenses

								skladTovar := &model.SkladTovar{}
								err = tx.Model(skladTovar).Where("sklad_tovars.sklad_id = ? AND sklad_tovars.tovar_id = ?", oldItem.SkladID, oldItem.ItemID).Update("quantity", factQuantity).First(skladTovar).Error
								if err != nil {
									return err
								}
							}
						}
					}
				}
			} else {
				invItem.LoadingStatus = utils.StatusLoading
				invItem.Status = utils.StatusNew
				err = r.gormDB.Model(invItem).Create(invItem).Error //creates new Item if new item is added
				if err != nil {
					return err
				}
				itemsToRecalculate = append(itemsToRecalculate, invItem)
			}
			if invItem.Time.Year() != time.Now().Year() || invItem.Time.Month() != time.Now().Month() || invItem.Time.Day() != time.Now().Day() {
				trafficItems = append(trafficItems, &model.AsyncJob{
					ItemID:    invItem.ItemID,
					ItemType:  invItem.Type,
					SkladID:   invItem.SkladID,
					TimeStamp: invItem.Time.UTC(),
					CreatedAt: time.Now(),
					Status:    utils.StatusNeedRecalculate,
				})
			}
		}
		err := tx.Model(inventarization).Select("time", "status", "result", "loading_status").Updates(inventarization).Error
		if err != nil {
			return err
		}
		return nil
	},
	)
	if err != nil {
		return nil, err
	}
	r.ConcurrentRecalculationForDailyStatistics(trafficItems)
	go r.RecalculateInventarization(itemsToRecalculate)

	return inventarization, nil
}

func (r *SkladDB) ChangeAfterItem(tx *gorm.DB, afterItem *model.InventarizationItem, oldItem *model.InventarizationItem) error {
	afterItem.BeforeTime = oldItem.Time
	afterItem.StartQuantity = oldItem.FactQuantity
	afterItem.PlanQuantity = afterItem.StartQuantity + afterItem.Income - afterItem.Removed - afterItem.Expenses
	afterItem.Difference = afterItem.FactQuantity - afterItem.PlanQuantity
	oldDifferenceSum := afterItem.DifferenceSum
	afterItem.DifferenceSum = afterItem.Difference * afterItem.Cost
	err := tx.Model(afterItem).Save(afterItem).Error
	if err != nil {
		return err
	}
	afterInventarization := &model.Inventarization{}
	err = tx.Model(afterInventarization).Where("inventarizations.id = ?", afterItem.InventarizationID).First(afterInventarization).Error
	if err != nil {
		return err
	}
	afterInventarization.Result += afterItem.DifferenceSum - oldDifferenceSum
	err = tx.Model(afterInventarization).Save(afterInventarization).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *SkladDB) UpdateInventarizationV2(inventarization *model.Inventarization) (*model.Inventarization, error) {
	itemsToRecalculate := []*model.InventarizationItem{}
	trafficItems := []*model.AsyncJob{}
	err := r.gormDB.Transaction(func(tx *gorm.DB) error {
		inventarization.Result = 0
		for i := 0; i < len(inventarization.InventarizationItems); i++ {
			invItem := inventarization.InventarizationItems[i]
			oldItem := &model.InventarizationItem{}
			res := tx.Model(oldItem).Where("inventarization_items.inventarization_id = ? AND inventarization_items.item_id = ? and inventarization_items.type = ?", inventarization.ID, invItem.ItemID, invItem.Type).First(oldItem)
			err := res.Error
			if err != nil {
				if err != gorm.ErrRecordNotFound {
					return err
				}
			}
			if oldItem.LoadingStatus == utils.StatusLoading {
				return errors.New("inventarization is loading")
			}

			if res.RowsAffected != 0 {
				if invItem.Status == utils.StatusClosed { //if they press save inventarization
					if oldItem.FactQuantity != invItem.FactQuantity || oldItem.Status == utils.StatusNew || oldItem.NeedToRecalculate { // if they changed the fact quantity
						oldItem.FactQuantity = invItem.FactQuantity                      //change fact quantity
						oldItem.Difference = oldItem.FactQuantity - oldItem.PlanQuantity //recalculate difference
						oldItem.DifferenceSum = oldItem.Difference * oldItem.Cost        //recalculate difference sum
						oldItem.Status = invItem.Status                                  //change status

						afterItem, err := r.GetAfterInventarizationItem(oldItem) //get after inventarization item
						if err != nil {
							if err != gorm.ErrRecordNotFound {
								return err
							}
						}

						if afterItem != nil { //if after inventarization item exists then change it
							err = r.ChangeAfterItem(tx, afterItem, oldItem)
							if err != nil {
								return err
							}
						} else { // if it does not exist change ostatok is sklads
							itemToScan := &model.InventarizationItem{}
							if oldItem.Status == utils.StatusClosed {
								if oldItem.Type == utils.TypeIngredient {
									err := r.gormDB.Model(&model.Postavka{}).Select("SUM(item_postavkas.quantity) as income").Joins("LEFT JOIN item_postavkas ON item_postavkas.postavka_id = postavkas.id").Where("postavkas.time < ? AND postavkas.time > ? AND item_postavkas.item_id = ? AND item_postavkas.type = ? AND postavkas.sklad_id = ?", time.Now(), oldItem.Time, oldItem.ItemID, utils.TypeIngredient, oldItem.SkladID).Scan(&itemToScan).Error
									if err != nil {
										return err
									}
									err = r.gormDB.Model(&model.RemoveFromSklad{}).Select("SUM(remove_from_sklad_items.quantity) as removed, SUM(sklad_ingredients.cost * remove_from_sklad_items.quantity) as removed_sum").Joins("LEFT JOIN remove_from_sklad_items ON remove_from_sklad_items.remove_id = remove_from_sklads.id").Joins("LEFT JOIN sklad_ingredients ON sklad_ingredients.ingredient_id = remove_from_sklad_items.item_id").Where("remove_from_sklads.time < ? AND remove_from_sklads.time > ? AND remove_from_sklad_items.item_id = ? AND remove_from_sklad_items.type = ? AND remove_from_sklads.sklad_id = ? AND sklad_ingredients.sklad_id = ?", time.Now(), oldItem.Time, oldItem.ItemID, utils.TypeIngredient, oldItem.SkladID, oldItem.SkladID).Scan(&itemToScan).Error
									if err != nil {
										return err
									}
									itemToScanExpenses := &model.InventarizationItem{}
									err = r.gormDB.Table("expence_ingredients").Select("SUM(expence_ingredients.quantity) as expenses").Where("expence_ingredients.ingredient_id = ? AND expence_ingredients.time between ? AND ? and sklad_id = ? and expence_ingredients.status = ?", oldItem.ItemID, oldItem.Time, time.Now(), oldItem.SkladID, utils.StatusClosed).Scan(&itemToScanExpenses).Error
									if err != nil {
										return err
									}
									itemToScan.Expenses += itemToScanExpenses.Expenses
									factQuantity := oldItem.FactQuantity + itemToScan.Income - itemToScan.Removed - itemToScan.Expenses

									skladIngredient := &model.SkladIngredient{}
									err = tx.Model(skladIngredient).Where("sklad_ingredients.sklad_id = ? AND sklad_ingredients.ingredient_id = ?", oldItem.SkladID, oldItem.ItemID).Update("quantity", factQuantity).First(skladIngredient).Error
									if err != nil {
										return err
									}
								} else if oldItem.Type == utils.TypeTovar {
									err := r.gormDB.Model(&model.Postavka{}).Select("SUM(item_postavkas.quantity) as income").Joins("LEFT JOIN item_postavkas ON item_postavkas.postavka_id = postavkas.id").Where("postavkas.time < ? AND postavkas.time > ? AND item_postavkas.item_id = ? AND item_postavkas.type = ? AND postavkas.sklad_id = ?", time.Now(), oldItem.Time, oldItem.ItemID, utils.TypeTovar, oldItem.SkladID).Scan(&itemToScan).Error
									if err != nil {
										return err
									}
									err = r.gormDB.Model(&model.RemoveFromSklad{}).Select("SUM(remove_from_sklad_items.quantity) as removed, SUM(sklad_tovars.cost * remove_from_sklad_items.quantity) as removed_sum").Joins("LEFT JOIN remove_from_sklad_items ON remove_from_sklad_items.remove_id = remove_from_sklads.id").Joins("LEFT JOIN sklad_tovars ON sklad_tovars.tovar_id = remove_from_sklad_items.item_id").Where("remove_from_sklads.time < ? AND remove_from_sklads.time > ? AND remove_from_sklad_items.item_id = ? AND remove_from_sklad_items.type = ? AND remove_from_sklads.sklad_id = ? AND sklad_tovars.sklad_id = ?", time.Now(), oldItem.Time, oldItem.ItemID, utils.TypeTovar, oldItem.SkladID, oldItem.SkladID).Scan(&itemToScan).Error
									if err != nil {
										return err
									}
									itemToScanExpenses := &model.InventarizationItem{}
									err = r.gormDB.Table("expence_tovars").Select("SUM(expence_tovars.quantity) as expenses").Where("expence_tovars.tovar_id = ? AND expence_tovars.time between ? AND ? and sklad_id = ?  AND expence_tovars.status = ?", oldItem.ItemID, oldItem.Time, time.Now(), oldItem.SkladID, utils.StatusClosed).Scan(&itemToScanExpenses).Error
									if err != nil {
										return err
									}
									itemToScan.Expenses += itemToScanExpenses.Expenses

									factQuantity := oldItem.FactQuantity + itemToScan.Income - itemToScan.Removed - itemToScan.Expenses

									skladTovar := &model.SkladTovar{}
									err = tx.Model(skladTovar).Where("sklad_tovars.sklad_id = ? AND sklad_tovars.tovar_id = ?", oldItem.SkladID, oldItem.ItemID).Update("quantity", factQuantity).First(skladTovar).Error
									if err != nil {
										return err
									}
								}
							}
						}
						if oldItem.NeedToRecalculate {
							if oldItem.Time.Before(time.Now()) {
								oldItem.NeedToRecalculate = false
							}
							itemsToRecalculate = append(itemsToRecalculate, oldItem)
						}
						inventarization.InventarizationItems[i] = oldItem

						err = tx.Model(oldItem).Save(oldItem).Error
						if err != nil {
							return err
						}
						if invItem.Time.Year() != time.Now().Year() || invItem.Time.Month() != time.Now().Month() || invItem.Time.Day() != time.Now().Day() {
							trafficItems = append(trafficItems, &model.AsyncJob{
								ItemID:    invItem.ItemID,
								ItemType:  invItem.Type,
								SkladID:   invItem.SkladID,
								TimeStamp: invItem.Time.UTC(),
								CreatedAt: time.Now(),
								Status:    utils.StatusNeedRecalculate,
							})
						}
					}
				} else if invItem.Status == utils.StatusOpened { //if inventarization is opened just change status
					oldItem.Status = utils.StatusOpened
					err := tx.Model(oldItem).Select("status").Updates(oldItem).Error
					if err != nil {
						return err
					}
				} else {
					return errors.New("wrong status")
				}
			}
			inventarization.Result += oldItem.DifferenceSum
		}

		err := tx.Model(inventarization).Select("time", "status", "result", "loading_status").Debug().Updates(inventarization).Error
		if err != nil {
			return err
		}

		return nil
	},
	)
	if err != nil {
		return nil, err
	}
	r.ConcurrentRecalculation(itemsToRecalculate)
	if len(trafficItems) > 0 {
		r.ConcurrentRecalculationForDailyStatistics(trafficItems)
	}
	return inventarization, nil
}

func (r *SkladDB) GoRecalcReader(wtg *sync.WaitGroup, myChannel chan *model.InventarizationItem, ctx context.Context) {
	defer wtg.Done()
	for {
		select {
		case <-ctx.Done():
			return
		case item, ok := <-myChannel:
			if !ok {
				return
			}
			if item != nil {
				err := r.RecalculateOneInventarization(item)
				if err != nil {
					log.Print(err)
				}
			}
		}
	}
}

func (r *SkladDB) FindInventarizationPlace(tx *gorm.DB, afterID int, invItem *model.InventarizationItem) error {
	return nil
}
func (r *SkladDB) GetInventarizationByID(id int) (*model.Inventarization, error) {
	inv := &model.Inventarization{}
	err := r.gormDB.Model(inv).Where("inventarizations.id = ?", id).First(inv).Error
	if err != nil {
		return nil, err
	}
	err = r.gormDB.Model(model.InventarizationItem{}).Select("*").Where("inventarization_items.inventarization_id = ?", id).Order("inventarization_items.id desc").Find(&inv.InventarizationItems).Error
	if err != nil {
		return nil, err
	}
	return inv, nil
}

func (r *SkladDB) GetInventarization(id int) (*model.InventarizationResponse, error) {
	inv := &model.InventarizationResponse{}
	err := r.gormDB.Model(model.Inventarization{}).Select("inventarizations.id, inventarizations.sklad_id, inventarizations.time, inventarizations.type, inventarizations.result, inventarizations.status, sklads.name as sklad").Joins("inner join sklads on inventarizations.sklad_id = sklads.id").Where("inventarizations.id = ?", id).First(inv).Error
	if err != nil {
		return nil, err
	}

	err = r.gormDB.Model(model.InventarizationItem{}).Select("inventarization_items.id, inventarization_items.inventarization_id, inventarization_items.time, inventarization_items.item_id, inventarization_items.sklad_id, inventarization_items.status, inventarization_items.type, inventarization_items.start_quantity, inventarization_items.expenses, inventarization_items.income, inventarization_items.removed, inventarization_items.removed_sum, inventarization_items.plan_quantity, inventarization_items.fact_quantity, inventarization_items.difference, inventarization_items.difference_sum, inventarization_items.cost, sklads.name as sklad_name, inventarization_items.before_time, inventarization_items.loading_status, inventarization_items.group_id, inventarization_items.is_visible").Joins("inner join sklads on sklads.id = inventarization_items.sklad_id").Where("inventarization_items.inventarization_id = ?", id).Order("inventarization_items.id desc").Find(&inv.InventarizationItems).Error
	if err != nil {
		return nil, err
	}
	// if inv.Status == utils.StatusOpened {
	// 	invent, err := r.GetOpenInventarizationResponse(inv)
	// 	if err != nil {
	// 		return nil, err
	// 	}
	// 	return invent, nil
	// }
	inv.LoadingStatus = utils.StatusLoaded
	sklad := &model.Sklad{}
	err = r.gormDB.Model(sklad).Select("*").Where("sklads.id = ?", inv.SkladID).First(sklad).Error
	if err != nil {
		return nil, err
	}
	for _, invItem := range inv.InventarizationItems {
		if invItem.LoadingStatus == utils.StatusLoading {
			inv.LoadingStatus = utils.StatusLoading
		}
		if invItem.Type == utils.TypeIngredient {
			ingredient := &model.Ingredient{}
			err := r.gormDB.Model(ingredient).Select("ingredients.ingredient_id, ingredients.name, ingredients.measure").Where("ingredients.ingredient_id = ? and ingredients.shop_id = ?", invItem.ItemID, sklad.ShopID).First(ingredient).Error
			if err != nil {
				return nil, err
			}
			invItem.ItemName = ingredient.Name
			invItem.Measure = ingredient.Measure
		} else if invItem.Type == utils.TypeTovar {
			tovar := &model.Tovar{}
			err := r.gormDB.Model(tovar).Select("tovars.tovar_id, tovars.name, tovars.measure").Where("tovars.tovar_id = ? and tovars.shop_id = ?", invItem.ItemID, sklad.ShopID).First(tovar).Error
			if err != nil {
				return nil, err
			}
			invItem.ItemName = tovar.Name
			invItem.Measure = tovar.Measure
		} else if invItem.Type == utils.TypeGroup {
			group, err := r.GetPureInventarizationGroup(invItem.ItemID)
			if err != nil {
				return nil, err
			}
			invItem.ItemName = group.Name
			var cost float64
			var quantity float64
			var numOfIngredients int
			var costOfIngredients float64
			for _, ingredient := range inv.InventarizationItems {
				if ingredient.GroupID == invItem.ItemID {
					costOfIngredients += math.Abs(float64(ingredient.Cost))
					cost += math.Abs(float64(ingredient.Cost)) * math.Abs(float64(ingredient.PlanQuantity))
					quantity += math.Abs(float64(ingredient.PlanQuantity))
					numOfIngredients++
					invItem.Time = ingredient.Time
					invItem.StartQuantity += ingredient.StartQuantity
					invItem.FactQuantity += ingredient.FactQuantity
					invItem.PlanQuantity += ingredient.PlanQuantity
					invItem.Difference += ingredient.Difference
					invItem.DifferenceSum += ingredient.DifferenceSum
					invItem.Measure = ingredient.Measure
					invItem.BeforeTime = ingredient.BeforeTime
					invItem.Expenses += ingredient.Expenses
					invItem.Income += ingredient.Income
					invItem.Removed += ingredient.Removed
					invItem.RemovedSum += ingredient.RemovedSum
					invItem.LoadingStatus = ingredient.LoadingStatus
				}
			}
			if quantity != 0 {
				invItem.Cost = float32(cost / quantity)
			} else {
				if numOfIngredients != 0 {
					invItem.Cost = float32(costOfIngredients / float64(numOfIngredients))
				} else {
					invItem.Cost = 0
				}
			}
		}
	}

	return inv, nil
}

func (r *SkladDB) GetAllInventarization(filter *model.Filter) ([]*model.InventarizationResponse, int64, error) {
	inv := []*model.InventarizationResponse{}
	var count int64
	res := r.gormDB.Model(model.Inventarization{}).Select("inventarizations.id, inventarizations.sklad_id, inventarizations.time, inventarizations.type, inventarizations.result, inventarizations.status, sklads.name as sklad").Joins("inner join sklads on inventarizations.sklad_id = sklads.id").Where("inventarizations.deleted = false and sklads.shop_id IN (?)", filter.AccessibleShops).Order("time desc")
	if res.Error != nil {
		return nil, 0, res.Error
	}
	newRes, count, err := filter.FilterResults(res, model.Inventarization{}, utils.DefaultPageSize, "time", "", "")
	if err != nil {
		return nil, 0, err
	}
	err = newRes.Scan(&inv).Error

	if err != nil {
		return nil, 0, err
	}

	return inv, count, nil
}

func (r *SkladDB) DeleteInventarization(id int) error {
	inv := &model.Inventarization{}
	err := r.gormDB.Model(inv).Where("inventarizations.id = ?", id).First(inv).Error
	if err != nil {
		return err
	}
	items := []*model.InventarizationItem{}
	trafficItems := []*model.AsyncJob{}
	err = r.gormDB.Model(model.InventarizationItem{}).Where("inventarization_items.inventarization_id = ?", id).Find(&items).Error
	if err != nil {
		return err
	}
	for _, item := range items {
		if item.GroupID == 0 {
			inv.InventarizationItems = append(inv.InventarizationItems, item)
		}
	}

	for _, item := range inv.InventarizationItems {
		err = r.DeleteInventarizationItem(item.ID)
		if err != nil {
			return err
		}
		if inv.Status == utils.StatusClosed {
			if item.Type == utils.TypeIngredient {
				ingredient := &model.SkladIngredient{}
				err := r.gormDB.Model(ingredient).Where("sklad_ingredients.ingredient_id = ? and sklad_ingredients.sklad_id = ?", item.ItemID, inv.SkladID).First(ingredient).Error
				if err != nil {
					return err
				}
				ingredient.Quantity = ingredient.Quantity - item.Difference
				err = r.gormDB.Model(ingredient).Where("sklad_ingredients.ingredient_id = ? and sklad_ingredients.sklad_id = ?", item.ItemID, inv.SkladID).Save(ingredient).Error
				if err != nil {
					return err
				}
			} else if item.Type == utils.TypeTovar {
				tovar := &model.SkladTovar{}
				err := r.gormDB.Model(tovar).Where("sklad_tovars.tovar_id = ? and sklad_tovars.sklad_id = ?", item.ItemID, inv.SkladID).First(tovar).Error
				if err != nil {
					return err
				}
				tovar.Quantity = tovar.Quantity - item.Difference
				err = r.gormDB.Model(tovar).Where("sklad_tovars.tovar_id = ? and sklad_tovars.sklad_id = ?", item.ItemID, inv.SkladID).Save(tovar).Error
				if err != nil {
					return err
				}
			}
		}

		if item.Time.Year() != time.Now().Year() || item.Time.Month() != time.Now().Month() || item.Time.Day() != time.Now().Day() {
			trafficItems = append(trafficItems, &model.AsyncJob{
				ItemID:    item.ItemID,
				ItemType:  item.Type,
				SkladID:   item.SkladID,
				TimeStamp: item.Time.UTC(),
				CreatedAt: time.Now(),
				Status:    utils.StatusNeedRecalculate,
			})
		}
	}
	err = r.gormDB.Model(inv).Where("inventarizations.id = ?", id).Delete(inv).Error
	if err != nil {
		return err
	}
	r.ConcurrentRecalculationForDailyStatistics(trafficItems)
	return nil
}

func (r *SkladDB) DeleteInventarizationItem(id int) error {
	item := &model.InventarizationItem{}
	trafficItems := []*model.AsyncJob{}
	err := r.gormDB.Model(item).Where("inventarization_items.id = ?", id).First(item).Error
	if err != nil {
		return err
	}
	if item.Type == utils.TypeGroup {
		inv := &model.Inventarization{}
		err := r.gormDB.Model(inv).Where("inventarizations.id = ?", item.InventarizationID).First(inv).Error
		if err != nil {
			return err
		}
		invItems := []*model.InventarizationItem{}
		err = r.gormDB.Model(model.InventarizationItem{}).Where("inventarization_items.inventarization_id = ?", inv.ID).Find(&invItems).Error
		if err != nil {
			return err
		}
		inv.InventarizationItems = invItems
		for _, invItem := range inv.InventarizationItems {
			if invItem.GroupID == item.ItemID {
				afterItem, err := r.GetAfterInventarizationItem(invItem)
				if err != nil {
					if err != gorm.ErrRecordNotFound {
						return err
					}
				}
				err = r.gormDB.Model(item).Delete(invItem).Where("inventarization_items.id = ?", invItem.ID).Error
				if err != nil {
					return err
				}
				if afterItem != nil {
					r.RecalculateOneInventarization(afterItem)
				}
			}
			if item.Time.Year() != time.Now().Year() || item.Time.Month() != time.Now().Month() || item.Time.Day() != time.Now().Day() {
				trafficItems = append(trafficItems, &model.AsyncJob{
					ItemID:    item.ItemID,
					ItemType:  item.Type,
					SkladID:   item.SkladID,
					TimeStamp: item.Time.UTC(),
					CreatedAt: time.Now(),
					Status:    utils.StatusNeedRecalculate,
				})
			}
		}
	}
	afterItem, err := r.GetAfterInventarizationItem(item)
	if err != nil {
		if err != gorm.ErrRecordNotFound {
			return err
		}
	}

	r.ConcurrentRecalculationForDailyStatistics(trafficItems)

	err = r.gormDB.Model(item).Delete(item).Where("inventarization_items.id = ?", id).Error
	if err != nil {
		return err
	}
	if afterItem != nil {
		r.RecalculateOneInventarization(afterItem)

	}
	if item.Time.Year() != time.Now().Year() || item.Time.Month() != time.Now().Month() || item.Time.Day() != time.Now().Day() {
		trafficItems = append(trafficItems, &model.AsyncJob{
			ItemID:    item.ItemID,
			ItemType:  item.Type,
			SkladID:   item.SkladID,
			TimeStamp: item.Time.UTC(),
			CreatedAt: time.Now(),
			Status:    utils.StatusNeedRecalculate,
		})
	}

	r.ConcurrentRecalculationForDailyStatistics(trafficItems)

	return nil
}

func (r *SkladDB) GetSkladByShopID(shopID int) (*model.Sklad, error) {
	sklad := &model.Sklad{}
	err := r.gormDB.Transaction(func(tx *gorm.DB) error {
		res := tx.Table("sklads").Where("shop_id = ?", shopID).First(&sklad)
		if res.Error != nil {
			return res.Error
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return sklad, nil
}

func (r *SkladDB) RecalculateNetCost() error {
	type IngredientNetCost struct {
		SkladID      int     `json:"sklad_id"`
		IngredientID int     `json:"ingredient_id"`
		Quantity     float32 `json:"quantity"`
		Sum          float32 `json:"sum"`
	}
	ingredients := []*IngredientNetCost{}
	err := r.gormDB.Model(&model.SkladIngredient{}).Select("sklad_ingredients.sklad_id, sklad_ingredients.ingredient_id, SUM(item_postavkas.quantity) as quantity, SUM(item_postavkas.cost * item_postavkas.quantity) as sum").Joins("inner join item_postavkas on item_postavkas.item_id = sklad_ingredients.ingredient_id inner join postavkas on postavkas.id = item_postavkas.postavka_id").Where("sklad_ingredients.sklad_id = postavkas.sklad_id and item_postavkas.type = ?", utils.TypeIngredient).Group("sklad_ingredients.sklad_id, sklad_ingredients.ingredient_id").Scan(&ingredients).Error
	if err != nil {
		return err
	}
	for _, ingredient := range ingredients {
		if ingredient.Quantity == 0 {
			continue
		}

		err := r.gormDB.Model(&model.SkladIngredient{}).Where("sklad_id = ? and ingredient_id = ?", ingredient.SkladID, ingredient.IngredientID).Update("cost", ingredient.Sum/ingredient.Quantity).Error
		if err != nil {
			return err
		}
		err = r.gormDB.Model(&model.InventarizationItem{}).Where("sklad_id = ? and item_id = ? and type = ?", ingredient.SkladID, ingredient.IngredientID, utils.TypeIngredient).Update("cost", ingredient.Sum/ingredient.Quantity).Error
		if err != nil {
			return err
		}
		err = r.gormDB.Model(&model.InventarizationItem{}).Where("sklad_id = ? and item_id = ? and type = ?", ingredient.SkladID, ingredient.IngredientID, utils.TypeIngredient).Update("difference_sum", gorm.Expr("difference * ?", ingredient.Sum/ingredient.Quantity)).Error
		if err != nil {
			return err
		}
	}

	type TovarNetCost struct {
		SkladID  int     `json:"sklad_id"`
		TovarID  int     `json:"tovar_id"`
		Quantity float32 `json:"quantity"`
		Sum      float32 `json:"sum"`
	}
	tovars := []*TovarNetCost{}
	err = r.gormDB.Model(&model.SkladTovar{}).Select("sklad_tovars.sklad_id, sklad_tovars.tovar_id, SUM(item_postavkas.quantity) as quantity, SUM(item_postavkas.cost * item_postavkas.quantity) as sum").Joins("inner join item_postavkas on item_postavkas.item_id = sklad_tovars.tovar_id inner join postavkas on postavkas.id = item_postavkas.postavka_id").Where("sklad_tovars.sklad_id = postavkas.sklad_id and item_postavkas.type = ?", utils.TypeTovar).Group("sklad_tovars.sklad_id, sklad_tovars.tovar_id").Scan(&tovars).Error
	if err != nil {
		return err
	}
	for _, tovar := range tovars {
		if tovar.Quantity == 0 {
			continue
		}
		err := r.gormDB.Model(&model.SkladTovar{}).Where("sklad_id = ? and tovar_id = ?", tovar.SkladID, tovar.TovarID).Update("cost", tovar.Sum/tovar.Quantity).Error
		if err != nil {
			return err
		}
		err = r.gormDB.Model(&model.InventarizationItem{}).Where("sklad_id = ? and item_id = ? and type = ?", tovar.SkladID, tovar.TovarID, utils.TypeTovar).Update("cost", tovar.Sum/tovar.Quantity).Error
		if err != nil {
			return err
		}
		err = r.gormDB.Model(&model.InventarizationItem{}).Where("sklad_id = ? and item_id = ? and type = ?", tovar.SkladID, tovar.TovarID, utils.TypeTovar).Update("difference_sum", gorm.Expr("difference * ?", tovar.Sum/tovar.Quantity)).Error
		if err != nil {
			return err
		}
	}
	return nil
}

func (r *SkladDB) GetSkladsByShopID(shopID int) ([]*model.Sklad, error) {
	sklads := []*model.Sklad{}
	res := r.gormDB.Model(&model.Sklad{}).Where("shop_id = ?", shopID).Scan(&sklads)
	if res.Error != nil {
		return nil, res.Error
	}
	return sklads, nil
}

func (r *SkladDB) DailyStatistic(shopID int) error {
	stat := []*model.DailyStatistic{}
	sklad, err := r.GetSkladByShopID(shopID)
	if err != nil {
		return err
	}
	type Statistic struct {
		PostavkaQuantity float32 `json:"postavka_quantity"`
		PostavkaCost     float32 `json:"postavka_cost"`
		RemoveFromSklad  float32 `json:"remove_from_sklad"`
		Inventarization  float32 `json:"inventarization"`
		TransferFrom     float32 `json:"transfer_from"`
		TransferTo       float32 `json:"transfer_to"`
		Sales            float32 `json:"sales"`
		CheckCost        float32 `json:"check_cost"`
		CheckPrice       float32 `json:"check_price"`
	}

	type Ingredient struct {
		ID       int     `json:"ingredient_id"`
		Name     string  `json:"ingredient_name"`
		Cost     float32 `json:"cost"`
		Quantity float32 `json:"quantity"`
	}

	now := time.Now().UTC()

	date := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)
	from := time.Date(now.Year(), now.Month(), now.Day()-1, 18, 0, 0, 0, time.UTC)
	to := time.Date(now.Year(), now.Month(), now.Day(), 18, 0, 0, 0, time.UTC)

	ingredients := []*Ingredient{}
	err = r.gormDB.Model(&model.SkladIngredient{}).Select("sklad_ingredients.ingredient_id as id, ingredients.name as name, sklad_ingredients.cost, sklad_ingredients.quantity").Joins("inner join ingredients on ingredients.ingredient_id = sklad_ingredients.ingredient_id").Where("sklad_ingredients.sklad_id = ? and ingredients.shop_id = ? and ingredients.deleted = ?", sklad.ID, shopID, false).Scan(&ingredients).Error
	if err != nil {
		return err
	}
	for _, ingredient := range ingredients {
		daily := &Statistic{}
		err := r.gormDB.Model(&model.Postavka{}).Select("CASE WHEN postavkas.type != 'transfer' THEN SUM(item_postavkas.quantity) ELSE 0 END as postavka_quantity, CASE WHEN postavkas.type = 'transfer' THEN SUM(item_postavkas.quantity) ELSE 0 END as transfer_to, CASE WHEN postavkas.type != 'transfer' THEN SUM(item_postavkas.quantity * item_postavkas.cost) ELSE 0 END as postavka_cost").Joins("inner join item_postavkas on item_postavkas.postavka_id = postavkas.id").Where("postavkas.sklad_id = ? and postavkas.time >= ? and postavkas.time < ? and item_postavkas.item_id = ? and item_postavkas.type = ? and postavkas.deleted = false", sklad.ID, from, to, ingredient.ID, utils.TypeIngredient).Group("postavkas.type").Scan(&daily).Error
		if err != nil {
			return err
		}
		err = r.gormDB.Model(&model.RemoveFromSklad{}).Select("CASE WHEN remove_from_sklads.type != 'transfer' THEN SUM(remove_from_sklad_items.quantity) ELSE 0 END as remove_from_sklad, CASE WHEN remove_from_sklads.type = 'transfer' THEN SUM(remove_from_sklad_items.quantity) ELSE 0 END as transfer_from").Joins("inner join remove_from_sklad_items on remove_from_sklad_items.remove_id = remove_from_sklads.id").Where("remove_from_sklads.sklad_id = ? and remove_from_sklads.time >= ? and remove_from_sklads.time < ? and remove_from_sklad_items.item_id = ? and remove_from_sklad_items.type = ? and remove_from_sklads.deleted = false", sklad.ID, from, to, ingredient.ID, utils.TypeIngredient).Group("remove_from_sklads.type").Scan(&daily).Error
		if err != nil {
			return err
		}
		err = r.gormDB.Model(&model.InventarizationItem{}).Select("SUM(inventarization_items.difference) as inventarization").Joins("inner join inventarizations on inventarizations.id = inventarization_items.inventarization_id").Where("inventarization_items.sklad_id = ? and inventarization_items.time >= ? and inventarization_items.time < ? and inventarization_items.item_id = ? and inventarization_items.type = ? and inventarization_items.status = ? and inventarizations.deleted = false", sklad.ID, from, to, ingredient.ID, utils.TypeIngredient, utils.StatusClosed).Scan(&daily).Error
		if err != nil {
			return err
		}
		err = r.gormDB.Model(&model.ExpenceIngredient{}).Select("SUM(expence_ingredients.quantity) as sales").Where("expence_ingredients.sklad_id = ? and expence_ingredients.time >= ? and expence_ingredients.time <= ? and expence_ingredients.ingredient_id = ? and expence_ingredients.status = ?", sklad.ID, from, to, ingredient.ID, utils.StatusClosed).Scan(&daily).Error
		if err != nil {
			return err
		}

		stat = append(stat, &model.DailyStatistic{
			Date:            date,
			ItemID:          ingredient.ID,
			ItemName:        ingredient.Name,
			SkladID:         sklad.ID,
			ShopID:          shopID,
			Type:            utils.TypeIngredient,
			Postavka:        daily.PostavkaQuantity,
			PostavkaCost:    daily.PostavkaCost,
			Inventarization: daily.Inventarization,
			Transfer:        daily.TransferTo - daily.TransferFrom,
			Sales:           daily.Sales,
			RemoveFromSklad: daily.RemoveFromSklad,
			Cost:            ingredient.Cost,
			Quantity:        ingredient.Quantity,
		})
	}

	type Tovar struct {
		ID       int     `json:"tovar_id"`
		Name     string  `json:"tovar_name"`
		Cost     float32 `json:"cost"`
		Quantity float32 `json:"quantity"`
	}

	tovars := []*Tovar{}
	err = r.gormDB.Model(&model.SkladTovar{}).Select("sklad_tovars.tovar_id as id, tovars.name as name, sklad_tovars.cost, sklad_tovars.quantity").Joins("inner join tovars on tovars.tovar_id = sklad_tovars.tovar_id").Where("sklad_tovars.sklad_id = ? and tovars.shop_id = ? and tovars.deleted = ?", sklad.ID, shopID, false).Scan(&tovars).Error
	if err != nil {
		return err
	}
	for _, tovar := range tovars {
		daily := &Statistic{}
		err := r.gormDB.Model(&model.Postavka{}).Select("CASE WHEN postavkas.type != 'transfer' THEN SUM(item_postavkas.quantity) ELSE 0 END as postavka_quantity, CASE WHEN postavkas.type = 'transfer' THEN SUM(item_postavkas.quantity) ELSE 0 END as transfer_to, CASE WHEN postavkas.type != 'transfer' THEN SUM(item_postavkas.quantity * item_postavkas.cost) ELSE 0 END as postavka_cost").Joins("inner join item_postavkas on item_postavkas.postavka_id = postavkas.id").Where("postavkas.sklad_id = ? and postavkas.time >= ? and postavkas.time < ? and item_postavkas.item_id = ? and item_postavkas.type = ? and postavkas.deleted = false", sklad.ID, from, to, tovar.ID, utils.TypeTovar).Group("postavkas.type").Scan(&daily).Error
		if err != nil {
			return err
		}
		err = r.gormDB.Model(&model.RemoveFromSklad{}).Select("CASE WHEN remove_from_sklads.type != 'transfer' THEN SUM(remove_from_sklad_items.quantity) ELSE 0 END as remove_from_sklad, CASE WHEN remove_from_sklads.type = 'transfer' THEN SUM(remove_from_sklad_items.quantity) ELSE 0 END as transfer_from").Joins("inner join remove_from_sklad_items on remove_from_sklad_items.remove_id = remove_from_sklads.id").Where("remove_from_sklads.sklad_id = ? and remove_from_sklads.time >= ? and remove_from_sklads.time < ? and remove_from_sklad_items.item_id = ? and remove_from_sklad_items.type = ? and remove_from_sklads.deleted = false", sklad.ID, from, to, tovar.ID, utils.TypeTovar).Group("remove_from_sklads.type").Scan(&daily).Error
		if err != nil {
			return err
		}
		err = r.gormDB.Model(&model.InventarizationItem{}).Select("SUM(inventarization_items.difference) as inventarization").Joins("inner join inventarizations on inventarizations.id = inventarization_items.inventarization_id").Where("inventarization_items.sklad_id = ? and inventarization_items.time >= ? and inventarization_items.time < ? and inventarization_items.item_id = ? and inventarization_items.type = ? and inventarization_items.status = ? and inventarizations.deleted = false", sklad.ID, from, to, tovar.ID, utils.TypeTovar, utils.StatusClosed).Scan(&daily).Error
		if err != nil {
			return err
		}
		err = r.gormDB.Model(&model.ExpenceTovar{}).Select("SUM(expence_tovars.quantity) as sales").Where("expence_tovars.sklad_id = ? and expence_tovars.time >= ? and expence_tovars.time < ? and expence_tovars.tovar_id = ? and expence_tovars.status = ?", sklad.ID, from, to, tovar.ID, utils.StatusClosed).Scan(&daily).Error
		if err != nil {
			return err
		}
		err = r.gormDB.Model(&model.CheckTovar{}).Select("SUM(check_tovars.cost) as check_cost, SUM(check_tovars.price) as check_price").Joins("inner join checks on checks.id = check_tovars.check_id").Where("checks.shop_id = ? and checks.opened_at >= ? and checks.closed_at < ? and check_tovars.tovar_id = ? and checks.status = ?", shopID, from, to, tovar.ID, utils.StatusClosed).Scan(&daily).Error
		if err != nil {
			return err
		}

		stat = append(stat, &model.DailyStatistic{
			Date:            date,
			ItemID:          tovar.ID,
			ItemName:        tovar.Name,
			SkladID:         sklad.ID,
			ShopID:          shopID,
			Type:            utils.TypeTovar,
			CheckCost:       daily.CheckCost,
			CheckPrice:      daily.CheckPrice,
			Postavka:        daily.PostavkaQuantity,
			PostavkaCost:    daily.PostavkaCost,
			Inventarization: daily.Inventarization,
			Transfer:        daily.TransferTo - daily.TransferFrom,
			Sales:           daily.Sales,
			RemoveFromSklad: daily.RemoveFromSklad,
			Cost:            tovar.Cost,
			Quantity:        tovar.Quantity,
		})
	}

	type TechCart struct {
		ID   int    `json:"tech_cart_id"`
		Name string `json:"tech_cart_name"`
	}

	tech_carts := []*TechCart{}
	err = r.gormDB.Model(&model.TechCart{}).Select("tech_cart_id as id, name").Where("deleted = ? and shop_id = ?", false, shopID).Scan(&tech_carts).Error
	if err != nil {
		return err
	}
	for _, tech_cart := range tech_carts {
		daily := &Statistic{}
		err = r.gormDB.Model(&model.CheckTechCart{}).Select("SUM(check_tech_carts.cost) as check_cost, SUM(check_tech_carts.price) as check_price, SUM(check_tech_carts.quantity) as sales").Joins("inner join checks on checks.id = check_tech_carts.check_id").Where("checks.shop_id = ? and checks.opened_at >= ? and checks.closed_at < ? and check_tech_carts.tech_cart_id = ? and checks.status = ?", shopID, from, to, tech_cart.ID, utils.StatusClosed).Scan(&daily).Error
		if err != nil {
			return err
		}
		stat = append(stat, &model.DailyStatistic{
			Date:       date,
			ItemID:     tech_cart.ID,
			ItemName:   tech_cart.Name,
			ShopID:     shopID,
			Type:       utils.TypeTechCart,
			CheckCost:  daily.CheckCost,
			CheckPrice: daily.CheckPrice,
			Sales:      daily.Sales,
		})
	}

	err = r.gormDB.Create(stat).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *SkladDB) RecalculateTrafficReport(item *model.AsyncJob) error {
	now := item.TimeStamp
	var from time.Time
	var to time.Time
	var date time.Time
	if now.Hour() >= 18 {
		date = time.Date(now.Year(), now.Month(), now.Day()+1, 0, 0, 0, 0, time.UTC)
		from = time.Date(now.Year(), now.Month(), now.Day(), 18, 0, 0, 0, time.UTC)
		to = time.Date(now.Year(), now.Month(), now.Day()+1, 18, 0, 0, 0, time.UTC)
	} else {
		date = time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)
		from = time.Date(now.Year(), now.Month(), now.Day()-1, 18, 0, 0, 0, time.UTC)
		to = time.Date(now.Year(), now.Month(), now.Day(), 18, 0, 0, 0, time.UTC)
	}

	var count int64
	res := r.gormDB.Model(&model.DailyStatistic{}).Where("date = ? and sklad_id = ?", date, item.SkladID).Count(&count)
	if res.Error != nil {
		if res.Error != gorm.ErrRecordNotFound {
			return res.Error
		}
	}
	if count != 0 {
		type Transfer struct {
			TransferFrom float32 `json:"transfer_from"`
			TransferTo   float32 `json:"transfer_to"`
		}
		transfer := &Transfer{}
		daily := &model.DailyStatistic{}
		res := r.gormDB.Model(&model.DailyStatistic{}).Where("date = ? and sklad_id = ? and item_id = ? and type = ?", date, item.SkladID, item.ItemID, item.ItemType).First(&daily)
		if res.Error != nil {
			return res.Error
		}
		newDaily := &model.DailyStatistic{}
		var inital, diff float32
		inital = daily.Postavka - daily.RemoveFromSklad - daily.Sales + daily.Inventarization + daily.Transfer
		if item.ItemType == utils.TypeIngredient {
			err := r.gormDB.Model(&model.Postavka{}).Select("CASE WHEN postavkas.type != 'transfer' THEN SUM(item_postavkas.quantity) ELSE 0 END as postavka, CASE WHEN postavkas.type != 'transfer' THEN SUM(item_postavkas.quantity * item_postavkas.cost) ELSE 0 END as postavka_cost").Joins("inner join item_postavkas on item_postavkas.postavka_id = postavkas.id").Where("postavkas.sklad_id = ? and postavkas.time >= ? and postavkas.time < ? and item_postavkas.item_id = ? and item_postavkas.type = ? and postavkas.deleted = false", item.SkladID, from, to, item.ItemID, utils.TypeIngredient).Group("postavkas.type").Scan(&newDaily).Error
			if err != nil {
				if err != gorm.ErrRecordNotFound {
					return err
				}
			}
			err = r.gormDB.Model(&model.RemoveFromSklad{}).Select("CASE WHEN remove_from_sklads.type != 'transfer' THEN SUM(remove_from_sklad_items.quantity) ELSE 0 END as remove_from_sklad").Joins("inner join remove_from_sklad_items on remove_from_sklad_items.remove_id = remove_from_sklads.id").Where("remove_from_sklads.sklad_id = ? and remove_from_sklads.time >= ? and remove_from_sklads.time < ? and remove_from_sklad_items.item_id = ? and remove_from_sklad_items.type = ? and remove_from_sklads.deleted = false", item.SkladID, from, to, item.ItemID, utils.TypeIngredient).Group("remove_from_sklads.type").Scan(&newDaily).Error
			if err != nil {
				if err != gorm.ErrRecordNotFound {
					return err
				}
			}
			err = r.gormDB.Model(&model.InventarizationItem{}).Select("SUM(inventarization_items.difference) as inventarization").Joins("inner join inventarizations on inventarizations.id = inventarization_items.inventarization_id").Where("inventarization_items.sklad_id = ? and inventarization_items.time >= ? and inventarization_items.time < ? and inventarization_items.item_id = ? and inventarization_items.type = ? and inventarization_items.status = ? and inventarizations.deleted = false", item.SkladID, from, to, item.ItemID, utils.TypeIngredient, utils.StatusClosed).Scan(&newDaily).Error
			if err != nil {
				if err != gorm.ErrRecordNotFound {
					return err
				}
			}
			err = r.gormDB.Model(&model.ExpenceIngredient{}).Select("SUM(expence_ingredients.quantity) as sales").Where("expence_ingredients.sklad_id = ? and expence_ingredients.time >= ? and expence_ingredients.time < ? and expence_ingredients.ingredient_id = ? and expence_ingredients.status = ?", item.SkladID, from, to, item.ItemID, utils.StatusClosed).Scan(&newDaily).Error
			if err != nil {
				if err != gorm.ErrRecordNotFound {
					return err
				}
			}
			err = r.gormDB.Model(&model.Postavka{}).Select("CASE WHEN postavkas.type = 'transfer' THEN SUM(item_postavkas.quantity) ELSE 0 END as transfer_to").Joins("inner join item_postavkas on item_postavkas.postavka_id = postavkas.id").Where("postavkas.sklad_id = ? and postavkas.time >= ? and postavkas.time < ? and item_postavkas.item_id = ? and item_postavkas.type = ? and postavkas.deleted = false", item.SkladID, from, to, item.ItemID, utils.TypeIngredient).Group("postavkas.type").Scan(&transfer).Error
			if err != nil {
				if err != gorm.ErrRecordNotFound {
					return err
				}
			}
			err = r.gormDB.Model(&model.RemoveFromSklad{}).Select("CASE WHEN remove_from_sklads.type = 'transfer' THEN SUM(remove_from_sklad_items.quantity) ELSE 0 END as transfer_from").Joins("inner join remove_from_sklad_items on remove_from_sklad_items.remove_id = remove_from_sklads.id").Where("remove_from_sklads.sklad_id = ? and remove_from_sklads.time >= ? and remove_from_sklads.time < ? and remove_from_sklad_items.item_id = ? and remove_from_sklad_items.type = ? and remove_from_sklads.deleted = false", item.SkladID, from, to, item.ItemID, utils.TypeIngredient).Group("remove_from_sklads.type").Scan(&transfer).Error
			if err != nil {
				if err != gorm.ErrRecordNotFound {
					return err
				}
			}
			newDaily.Transfer = transfer.TransferTo - transfer.TransferFrom
			diff = newDaily.Postavka - newDaily.RemoveFromSklad - newDaily.Sales + newDaily.Transfer + newDaily.Inventarization
			difference := diff - inital
			daily.Quantity += difference
		} else if item.ItemType == utils.TypeTovar {
			err := r.gormDB.Model(&model.Postavka{}).Select("CASE WHEN postavkas.type != 'transfer' THEN SUM(item_postavkas.quantity) ELSE 0 END as postavka, CASE WHEN postavkas.type != 'transfer' THEN SUM(item_postavkas.quantity * item_postavkas.cost) ELSE 0 END as postavka_cost").Joins("inner join item_postavkas on item_postavkas.postavka_id = postavkas.id").Where("postavkas.sklad_id = ? and postavkas.time >= ? and postavkas.time < ? and item_postavkas.item_id = ? and item_postavkas.type = ? and postavkas.deleted = false", item.SkladID, from, to, item.ItemID, utils.TypeTovar).Group("postavkas.type").Scan(&newDaily).Error
			if err != nil {
				if err != gorm.ErrRecordNotFound {
					return err
				}
			}
			err = r.gormDB.Model(&model.RemoveFromSklad{}).Select("CASE WHEN remove_from_sklads.type != 'transfer' THEN SUM(remove_from_sklad_items.quantity) ELSE 0 END as remove_from_sklad").Joins("inner join remove_from_sklad_items on remove_from_sklad_items.remove_id = remove_from_sklads.id").Where("remove_from_sklads.sklad_id = ? and remove_from_sklads.time >= ? and remove_from_sklads.time < ? and remove_from_sklad_items.item_id = ? and remove_from_sklad_items.type = ? and remove_from_sklads.deleted = false", item.SkladID, from, to, item.ItemID, utils.TypeTovar).Group("remove_from_sklads.type").Scan(&newDaily).Error
			if err != nil {
				if err != gorm.ErrRecordNotFound {
					return err
				}
			}
			err = r.gormDB.Model(&model.InventarizationItem{}).Select("SUM(inventarization_items.difference) as inventarization").Joins("inner join inventarizations on inventarizations.id = inventarization_items.inventarization_id").Where("inventarization_items.sklad_id = ? and inventarization_items.time >= ? and inventarization_items.time < ? and inventarization_items.item_id = ? and inventarization_items.type = ? and inventarization_items.status = ? and inventarizations.deleted = false", item.SkladID, from, to, item.ItemID, utils.TypeTovar, utils.StatusClosed).Scan(&newDaily).Error
			if err != nil {
				if err != gorm.ErrRecordNotFound {
					return err
				}
			}
			err = r.gormDB.Model(&model.ExpenceTovar{}).Select("SUM(expence_tovars.quantity) as sales").Where("expence_tovars.sklad_id = ? and expence_tovars.time >= ? and expence_tovars.time < ? and expence_tovars.tovar_id = ? and expence_tovars.status = ?", item.SkladID, from, to, item.ItemID, utils.StatusClosed).Scan(&newDaily).Error
			if err != nil {
				if err != gorm.ErrRecordNotFound {
					return err
				}
			}
			err = r.gormDB.Model(&model.CheckTovar{}).Select("SUM(check_tovars.cost) as check_cost, SUM(check_tovars.price) as check_price").Joins("inner join checks on checks.id = check_tovars.check_id").Where("checks.shop_id = ? and checks.opened_at >= ? and checks.closed_at < ? and check_tovars.tovar_id = ? and checks.status = ?", item.ShopID, from, to, item.ItemID, utils.StatusClosed).Scan(&newDaily).Error
			if err != nil {
				if err != gorm.ErrRecordNotFound {
					return err
				}
			}
			err = r.gormDB.Model(&model.Postavka{}).Select("CASE WHEN postavkas.type = 'transfer' THEN SUM(item_postavkas.quantity) ELSE 0 END as transfer_to").Joins("inner join item_postavkas on item_postavkas.postavka_id = postavkas.id").Where("postavkas.sklad_id = ? and postavkas.time >= ? and postavkas.time < ? and item_postavkas.item_id = ? and item_postavkas.type = ? and postavkas.deleted = false", item.SkladID, from, to, item.ItemID, utils.TypeTovar).Group("postavkas.type").Scan(&transfer).Error
			if err != nil {
				if err != gorm.ErrRecordNotFound {
					return err
				}
			}
			err = r.gormDB.Model(&model.RemoveFromSklad{}).Select("CASE WHEN remove_from_sklads.type = 'transfer' THEN SUM(remove_from_sklad_items.quantity) ELSE 0 END as transfer_from").Joins("inner join remove_from_sklad_items on remove_from_sklad_items.remove_id = remove_from_sklads.id").Where("remove_from_sklads.sklad_id = ? and remove_from_sklads.time >= ? and remove_from_sklads.time < ? and remove_from_sklad_items.item_id = ? and remove_from_sklad_items.type = ? and remove_from_sklads.deleted = false", item.SkladID, from, to, item.ItemID, utils.TypeTovar).Group("remove_from_sklads.type").Scan(&transfer).Error
			if err != nil {
				if err != gorm.ErrRecordNotFound {
					return err
				}
			}
			newDaily.Transfer = transfer.TransferTo - transfer.TransferFrom
			diff = newDaily.Postavka - newDaily.RemoveFromSklad - newDaily.Sales + newDaily.Transfer + newDaily.Inventarization
			difference := diff - inital
			daily.Quantity += difference
		} else if item.ItemType == utils.TypeTechCart {
			err := r.gormDB.Model(&model.CheckTechCart{}).Select("SUM(check_tech_carts.cost) as check_cost, SUM(check_tech_carts.price) as check_price, SUM(check_tech_carts.quantity) as sales").Joins("inner join checks on checks.id = check_tech_carts.check_id").Where("checks.shop_id = ? and checks.opened_at >= ? and checks.closed_at < ? and check_tech_carts.tech_cart_id = ? and checks.status = ?", item.ShopID, from, to, item.ItemID, utils.StatusClosed).Scan(&newDaily).Error
			if err != nil {
				if err != gorm.ErrRecordNotFound {
					return err
				}
			}
		}
		daily.Postavka = newDaily.Postavka
		daily.PostavkaCost = newDaily.PostavkaCost
		daily.RemoveFromSklad = newDaily.RemoveFromSklad
		daily.Inventarization = newDaily.Inventarization
		daily.Sales = newDaily.Sales
		daily.CheckCost = newDaily.CheckCost
		daily.CheckPrice = newDaily.CheckPrice
		res = r.gormDB.Model(&model.DailyStatistic{}).Select("*").Where("id = ?", daily.ID).Save(daily)
		if res.Error != nil {
			return res.Error
		}
	}
	return nil
}

func (r *SkladDB) GetTrafficReport(filter *model.Filter) ([]*model.TrafficReport, int64, error) {
	flag := false
	if filter.From == filter.To && filter.From.Day() == time.Now().Day() {
		flag = true
		filter.From = time.Date(filter.From.Year(), filter.From.Month(), filter.From.Day()-1, 0, 0, 0, 0, time.UTC)
	}
	traffic := []*model.TrafficReport{}
	var res *gorm.DB
	if filter.Type != "" {
		if filter.Type == utils.TypeTovar {
			if filter.Category != 0 {
				res = r.gormDB.Model(&model.DailyStatistic{}).Select("daily_statistics.item_id, daily_statistics.item_name, tovars.measure, daily_statistics.type, SUM(daily_statistics.postavka) as postavka, SUM(daily_statistics.postavka_cost) as postavka_cost, SUM(daily_statistics.inventarization) as inventarization, SUM(daily_statistics.transfer) as transfer, SUM(daily_statistics.sales) as sales, SUM(daily_statistics.remove_from_sklad) as remove_from_sklad").Joins("inner join sklads on sklads.id = daily_statistics.sklad_id left join tovars on (tovars.tovar_id = daily_statistics.item_id and daily_statistics.type = 'tovar' and tovars.shop_id = sklads.shop_id)").Where("daily_statistics.date between ? and ? and daily_statistics.sklad_id IN (?) and (daily_statistics.postavka > 0 or daily_statistics.inventarization != 0 or daily_statistics.sales > 0 or daily_statistics.remove_from_sklad > 0) and daily_statistics.type = ? and tovars.category = ?", filter.From, filter.To, filter.SkladID, filter.Type, filter.Category).Group("daily_statistics.item_id, daily_statistics.item_name, tovars.measure, daily_statistics.type")
				if res.Error != nil {
					return nil, 0, res.Error
				}
			} else {
				res = r.gormDB.Model(&model.DailyStatistic{}).Select("daily_statistics.item_id, daily_statistics.item_name, tovars.measure, daily_statistics.type, SUM(daily_statistics.postavka) as postavka, SUM(daily_statistics.postavka_cost) as postavka_cost, SUM(daily_statistics.inventarization) as inventarization, SUM(daily_statistics.transfer) as transfer, SUM(daily_statistics.sales) as sales, SUM(daily_statistics.remove_from_sklad) as remove_from_sklad").Joins("inner join sklads on sklads.id = daily_statistics.sklad_id left join tovars on (tovars.tovar_id = daily_statistics.item_id and daily_statistics.type = 'tovar' and tovars.shop_id = sklads.shop_id)").Where("daily_statistics.date between ? and ? and daily_statistics.sklad_id IN (?) and (daily_statistics.postavka > 0 or daily_statistics.inventarization != 0 or daily_statistics.sales > 0 or daily_statistics.remove_from_sklad > 0) and daily_statistics.type = ?", filter.From, filter.To, filter.SkladID, filter.Type).Group("daily_statistics.item_id, daily_statistics.item_name, tovars.measure, daily_statistics.type")
				if res.Error != nil {
					return nil, 0, res.Error
				}
			}
		} else if filter.Type == utils.TypeIngredient {
			if filter.Category != 0 {
				res = r.gormDB.Model(&model.DailyStatistic{}).Select("daily_statistics.item_id, daily_statistics.item_name, ingredients.measure, daily_statistics.type, SUM(daily_statistics.postavka) as postavka, SUM(daily_statistics.postavka_cost) as postavka_cost, SUM(daily_statistics.inventarization) as inventarization, SUM(daily_statistics.transfer) as transfer, SUM(daily_statistics.sales) as sales, SUM(daily_statistics.remove_from_sklad) as remove_from_sklad").Joins("inner join sklads on sklads.id = daily_statistics.sklad_id left join ingredients on (ingredients.ingredient_id = daily_statistics.item_id and daily_statistics.type = 'ingredient' and ingredients.shop_id = sklads.shop_id)").Where("daily_statistics.date between ? and ? and daily_statistics.sklad_id IN (?) and (daily_statistics.postavka > 0 or daily_statistics.inventarization != 0 or daily_statistics.sales > 0 or daily_statistics.remove_from_sklad > 0) and daily_statistics.type = ? and ingredients.category = ?", filter.From, filter.To, filter.SkladID, filter.Type, filter.Category).Group("daily_statistics.item_id, daily_statistics.item_name, ingredients.measure, daily_statistics.type")
				if res.Error != nil {
					return nil, 0, res.Error
				}
			} else {
				res = r.gormDB.Model(&model.DailyStatistic{}).Select("daily_statistics.item_id, daily_statistics.item_name, ingredients.measure, daily_statistics.type, SUM(daily_statistics.postavka) as postavka, SUM(daily_statistics.postavka_cost) as postavka_cost, SUM(daily_statistics.inventarization) as inventarization, SUM(daily_statistics.transfer) as transfer, SUM(daily_statistics.sales) as sales, SUM(daily_statistics.remove_from_sklad) as remove_from_sklad").Joins("inner join sklads on sklads.id = daily_statistics.sklad_id left join ingredients on (ingredients.ingredient_id = daily_statistics.item_id and daily_statistics.type = 'ingredient' and ingredients.shop_id = sklads.shop_id)").Where("daily_statistics.date between ? and ? and daily_statistics.sklad_id IN (?) and (daily_statistics.postavka > 0 or daily_statistics.inventarization != 0 or daily_statistics.sales > 0 or daily_statistics.remove_from_sklad > 0) and daily_statistics.type = ?", filter.From, filter.To, filter.SkladID, filter.Type).Group("daily_statistics.item_id, daily_statistics.item_name, ingredients.measure, daily_statistics.type")
				if res.Error != nil {
					return nil, 0, res.Error
				}
			}
		} else {
			return nil, 0, errors.New("invalid type")
		}
	} else if filter.Category != 0 {
		res = r.gormDB.Model(&model.DailyStatistic{}).Select("daily_statistics.item_id, daily_statistics.item_name, CASE WHEN daily_statistics.type = 'ingredient' THEN ingredients.measure ELSE tovars.measure END,  daily_statistics.type, SUM(daily_statistics.postavka) as postavka, SUM(daily_statistics.postavka_cost) as postavka_cost, SUM(daily_statistics.inventarization) as inventarization, SUM(daily_statistics.transfer) as transfer, SUM(daily_statistics.sales) as sales, SUM(daily_statistics.remove_from_sklad) as remove_from_sklad").Joins("inner join sklads on sklads.id = daily_statistics.sklad_id left join tovars on (tovars.tovar_id = daily_statistics.item_id and daily_statistics.type = 'tovar' and tovars.shop_id = sklads.shop_id) left join ingredients on (ingredients.ingredient_id = daily_statistics.item_id and daily_statistics.type = 'ingredient' and ingredients.shop_id = sklads.shop_id)").Where("daily_statistics.type != ? and daily_statistics.date >= ? and daily_statistics.date <= ? and daily_statistics.sklad_id IN (?) and (daily_statistics.postavka > 0 or daily_statistics.inventarization != 0 or daily_statistics.sales > 0 or daily_statistics.remove_from_sklad > 0) and (tovars.category = ? or ingredients.category = ?)", utils.TypeTechCart, filter.From, filter.To, filter.SkladID, filter.Category, filter.Category).Group("daily_statistics.item_id, daily_statistics.item_name, daily_statistics.type, CASE WHEN daily_statistics.type = 'ingredient' THEN ingredients.measure ELSE tovars.measure END")
		if res.Error != nil {
			return nil, 0, res.Error
		}
	} else {
		res = r.gormDB.Model(&model.DailyStatistic{}).Select("daily_statistics.item_id, daily_statistics.item_name, CASE WHEN daily_statistics.type = 'ingredient' THEN ingredients.measure ELSE tovars.measure END,  daily_statistics.type, SUM(daily_statistics.postavka) as postavka, SUM(daily_statistics.postavka_cost) as postavka_cost, SUM(daily_statistics.inventarization) as inventarization, SUM(daily_statistics.transfer) as transfer, SUM(daily_statistics.sales) as sales, SUM(daily_statistics.remove_from_sklad) as remove_from_sklad").Joins("inner join sklads on sklads.id = daily_statistics.sklad_id left join tovars on (tovars.tovar_id = daily_statistics.item_id and daily_statistics.type = 'tovar' and tovars.shop_id = sklads.shop_id) left join ingredients on (ingredients.ingredient_id = daily_statistics.item_id and daily_statistics.type = 'ingredient' and ingredients.shop_id = sklads.shop_id)").Where("daily_statistics.type != ? and daily_statistics.date >= ? and daily_statistics.date <= ? and daily_statistics.sklad_id IN (?) and (daily_statistics.postavka > 0 or daily_statistics.inventarization != 0 or daily_statistics.sales > 0 or daily_statistics.remove_from_sklad > 0)", utils.TypeTechCart, filter.From, filter.To, filter.SkladID).Group("daily_statistics.item_id, daily_statistics.item_name, daily_statistics.type, CASE WHEN daily_statistics.type = 'ingredient' THEN ingredients.measure ELSE tovars.measure END")
		if res.Error != nil {
			return nil, 0, res.Error
		}
	}

	newRes, count, err := filter.FilterResults(res, model.DailyStatistic{}, utils.DefaultPageSize, "", fmt.Sprintf("daily_statistics.item_name ilike '%%%s%%'", filter.Search), "daily_statistics.item_id DESC")
	if err != nil {
		return nil, 0, err
	}
	err = newRes.Scan(&traffic).Error
	if err != nil {
		return nil, 0, err
	}
	if flag {
		filter.From = filter.To
	}
	for _, val := range traffic {
		if flag {
			val.Inventarization = 0
			val.Transfer = 0
			val.Sales = 0
			val.RemoveFromSklad = 0
			val.Postavka = 0
			val.PostavkaCost = 0
		}
		val, err = r.GetPreviousDayStatistic(val, filter)
		if err != nil {
			return nil, 0, err
		}
		if filter.To.Day() == time.Now().Day() {
			val, err = r.GetTodayDayStatistic(val, filter)
			if err != nil {
				return nil, 0, err
			}
		}
		var invPlus float32 = 0
		var invMinus float32 = 0
		if val.Inventarization > 0 {
			invPlus = val.Inventarization
		} else {
			invMinus = float32(math.Abs(float64(val.Inventarization)))
		}
		var transferPlus float32 = 0
		var transferMinus float32 = 0
		if val.Transfer > 0 {
			transferPlus = val.Transfer
		} else {
			transferMinus = float32(math.Abs(float64(val.Transfer)))
		}
		val.Income = val.Postavka + invPlus + transferPlus
		val.Consumption = val.RemoveFromSklad + invMinus + val.Sales + transferMinus
		val.FinalOstatki = val.InitialOstatki + val.Income - val.Consumption
		if val.InitialOstatki != 0 || val.Postavka != 0 {
			if val.InitialOstatki < 0 {
				if val.Postavka > 0 {
					val.FinalNetCost = val.PostavkaCost / val.Postavka
				} else {
					val.FinalNetCost = 0
				}
			} else {
				if val.FinalOstatki < 0 {
					val.FinalNetCost = 0
				} else {
					val.FinalNetCost = (val.InitialSum + val.PostavkaCost) / (val.InitialOstatki + val.Postavka)
				}
			}
		} else {
			val.FinalNetCost = 0
		}
		if val.FinalOstatki < 0 {
			val.FinalSum = 0
		}
		val.FinalSum = val.FinalNetCost * val.FinalOstatki
		if val.FinalSum <= 0 {
			val.FinalSum = 0
		}
	}

	if err != nil {
		return nil, 0, err
	}

	return traffic, count, nil
}

func (r *SkladDB) GetTodayDayStatistic(traffic *model.TrafficReport, filter *model.Filter) (*model.TrafficReport, error) {
	from := time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day()-1, 18, 0, 0, 0, time.UTC)
	to := time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), 18, 0, 0, 0, time.UTC)
	type Statistic struct {
		PostavkaQuantity float32 `json:"postavka_quantity"`
		PostavkaCost     float32 `json:"postavka_cost"`
		RemoveFromSklad  float32 `json:"remove_from_sklad"`
		Inventarization  float32 `json:"inventarization"`
		TransferFrom     float32 `json:"transfer_from"`
		TransferTo       float32 `json:"transfer_to"`
		Sales            float32 `json:"sales"`
		CheckCost        float32 `json:"check_cost"`
		CheckPrice       float32 `json:"check_price"`
	}
	TodayStat := &Statistic{}
	if traffic.Type == utils.TypeIngredient {
		err := r.gormDB.Model(&model.Postavka{}).Select("CASE WHEN postavkas.type != 'transfer' THEN SUM(item_postavkas.quantity) ELSE 0 END as postavka_quantity, CASE WHEN postavkas.type = 'transfer' THEN SUM(item_postavkas.quantity) ELSE 0 END as transfer_to, CASE WHEN postavkas.type != 'transfer' THEN SUM(item_postavkas.quantity * item_postavkas.cost) ELSE 0 END as postavka_cost").Joins("inner join item_postavkas on item_postavkas.postavka_id = postavkas.id").Where("postavkas.sklad_id IN (?) and postavkas.time >= ? and postavkas.time < ? and item_postavkas.item_id = ? and item_postavkas.type = ? and postavkas.deleted = false", filter.SkladID, from, to, traffic.ItemID, utils.TypeIngredient).Group("postavkas.type").Scan(&TodayStat).Error
		if err != nil {
			if err != gorm.ErrRecordNotFound {
				return nil, err
			}
		}
		err = r.gormDB.Model(&model.RemoveFromSklad{}).Select("CASE WHEN remove_from_sklads.type != 'transfer' THEN SUM(remove_from_sklad_items.quantity) ELSE 0 END as remove_from_sklad, CASE WHEN remove_from_sklads.type = 'transfer' THEN SUM(remove_from_sklad_items.quantity) ELSE 0 END as transfer_from").Joins("inner join remove_from_sklad_items on remove_from_sklad_items.remove_id = remove_from_sklads.id").Where("remove_from_sklads.sklad_id IN (?) and remove_from_sklads.time >= ? and remove_from_sklads.time < ? and remove_from_sklad_items.item_id = ? and remove_from_sklad_items.type = ? and remove_from_sklads.deleted = false", filter.SkladID, from, to, traffic.ItemID, utils.TypeIngredient).Group("remove_from_sklads.type").Scan(&TodayStat).Error
		if err != nil {
			if err != gorm.ErrRecordNotFound {
				return nil, err
			}
		}
		err = r.gormDB.Model(&model.InventarizationItem{}).Select("SUM(inventarization_items.difference) as inventarization").Joins("inner join inventarizations on inventarizations.id = inventarization_items.inventarization_id").Where("inventarization_items.sklad_id IN (?) and inventarization_items.time >= ? and inventarization_items.time < ? and inventarization_items.item_id = ? and inventarization_items.type = ? and inventarization_items.status = ? and inventarizations.deleted = false", filter.SkladID, from, to, traffic.ItemID, utils.TypeIngredient, utils.StatusClosed).Scan(&TodayStat).Error
		if err != nil {
			if err != gorm.ErrRecordNotFound {
				return nil, err
			}
		}
		err = r.gormDB.Model(&model.ExpenceIngredient{}).Select("SUM(expence_ingredients.quantity) as sales").Where("expence_ingredients.sklad_id IN (?) and expence_ingredients.time >= ? and expence_ingredients.time < ? and expence_ingredients.ingredient_id = ? and expence_ingredients.status = ?", filter.SkladID, from, to, traffic.ItemID, utils.StatusClosed).Scan(&TodayStat).Error
		if err != nil {
			if err != gorm.ErrRecordNotFound {
				return nil, err
			}
		}
	} else if traffic.Type == utils.TypeTovar {
		err := r.gormDB.Model(&model.Postavka{}).Select("CASE WHEN postavkas.type != 'transfer' THEN SUM(item_postavkas.quantity) ELSE 0 END as postavka_quantity, CASE WHEN postavkas.type = 'transfer' THEN SUM(item_postavkas.quantity) ELSE 0 END as transfer_to, CASE WHEN postavkas.type != 'transfer' THEN SUM(item_postavkas.quantity * item_postavkas.cost) ELSE 0 END as postavka_cost").Joins("inner join item_postavkas on item_postavkas.postavka_id = postavkas.id").Where("postavkas.sklad_id IN (?) and postavkas.time >= ? and postavkas.time < ? and item_postavkas.item_id = ? and item_postavkas.type = ? and postavkas.deleted = false", filter.SkladID, from, to, traffic.ItemID, utils.TypeTovar).Group("postavkas.type").Scan(&TodayStat).Error
		if err != nil {
			if err != gorm.ErrRecordNotFound {
				return nil, err
			}
		}
		err = r.gormDB.Model(&model.RemoveFromSklad{}).Select("CASE WHEN remove_from_sklads.type != 'transfer' THEN SUM(remove_from_sklad_items.quantity) ELSE 0 END as remove_from_sklad, CASE WHEN remove_from_sklads.type = 'transfer' THEN SUM(remove_from_sklad_items.quantity) ELSE 0 END as transfer_from").Joins("inner join remove_from_sklad_items on remove_from_sklad_items.remove_id = remove_from_sklads.id").Where("remove_from_sklads.sklad_id IN (?) and remove_from_sklads.time >= ? and remove_from_sklads.time < ? and remove_from_sklad_items.item_id = ? and remove_from_sklad_items.type = ? and remove_from_sklads.deleted = false", filter.SkladID, from, to, traffic.ItemID, utils.TypeTovar).Group("remove_from_sklads.type").Scan(&TodayStat).Error
		if err != nil {
			if err != gorm.ErrRecordNotFound {
				return nil, err
			}
		}
		err = r.gormDB.Model(&model.InventarizationItem{}).Select("SUM(inventarization_items.difference) as inventarization").Joins("inner join inventarizations on inventarizations.id = inventarization_items.inventarization_id").Where("inventarization_items.sklad_id IN (?) and inventarization_items.time >= ? and inventarization_items.time < ? and inventarization_items.item_id = ? and inventarization_items.type = ? and inventarization_items.status = ? and inventarizations.deleted = false", filter.SkladID, from, to, traffic.ItemID, utils.TypeTovar, utils.StatusClosed).Scan(&TodayStat).Error
		if err != nil {
			if err != gorm.ErrRecordNotFound {
				return nil, err
			}
		}
		err = r.gormDB.Model(&model.ExpenceTovar{}).Select("SUM(expence_tovars.quantity) as sales").Where("expence_tovars.sklad_id IN (?) and expence_tovars.time >= ? and expence_tovars.time < ? and expence_tovars.tovar_id = ? and expence_tovars.status = ?", filter.SkladID, from, to, traffic.ItemID, utils.StatusClosed).Scan(&TodayStat).Error
		if err != nil {
			if err != gorm.ErrRecordNotFound {
				return nil, err
			}
		}
	} else {
		return nil, errors.New("unknown type")
	}
	traffic.Postavka += TodayStat.PostavkaQuantity
	traffic.PostavkaCost += TodayStat.PostavkaCost
	traffic.Transfer += TodayStat.TransferTo - TodayStat.TransferFrom
	traffic.RemoveFromSklad += TodayStat.RemoveFromSklad
	traffic.Inventarization += TodayStat.Inventarization
	traffic.Sales += TodayStat.Sales
	return traffic, nil
}

func (r *SkladDB) GetPreviousDayStatistic(traffic *model.TrafficReport, filter *model.Filter) (*model.TrafficReport, error) {
	PrevDayStat := &model.DailyStatistic{}
	res := r.gormDB.Model(&model.DailyStatistic{}).Select("cost, quantity").Where("sklad_id IN (?) and item_id = ? and type = ? and date < ?", filter.SkladID, traffic.ItemID, traffic.Type, filter.From).Order("date desc").First(&PrevDayStat)
	if res.Error != nil {
		if res.Error != gorm.ErrRecordNotFound {
			return nil, res.Error
		}
	}
	if res.RowsAffected == 0 {
		traffic.InitialOstatki = 0
		traffic.InitialNetCost = 0
		traffic.InitialSum = traffic.InitialNetCost * traffic.InitialOstatki
	} else {
		traffic.InitialOstatki = PrevDayStat.Quantity
		traffic.InitialNetCost = PrevDayStat.Cost
		traffic.InitialSum = traffic.InitialNetCost * traffic.InitialOstatki
		if traffic.InitialSum < 0 {
			traffic.InitialSum = 0
		}
	}
	return traffic, nil
}

func (r *SkladDB) GetInventarizationDetailsIncome(id int) ([]*model.InventarizationDetailsIncome, error) {
	inventarizationDetailsIncome := []*model.InventarizationDetailsIncome{}
	item := &model.InventarizationItem{}
	err := r.gormDB.Model(&model.InventarizationItem{}).Where("id = ?", id).Scan(&item).Error
	if err != nil {
		return nil, err
	}
	sklad := &model.Sklad{}
	err = r.gormDB.Model(&model.Sklad{}).Where("id = ?", item.SkladID).Scan(&sklad).Error
	if err != nil {
		return nil, err
	}
	if item.Type == utils.TypeGroup {
		items := []*model.InventarizationItem{}
		err = r.gormDB.Model(&model.InventarizationItem{}).Where("inventarization_id = ? and group_id = ?", item.InventarizationID, item.ItemID).Scan(&items).Error
		if err != nil {
			return nil, err
		}
		ids := []int{}
		for _, val := range items {
			ids = append(ids, val.ItemID)
		}

		err := r.gormDB.Model(&model.InventarizationItem{}).Select("postavkas.time, item_postavkas.quantity, item_postavkas.cost, (item_postavkas.quantity * item_postavkas.cost) as sum, ingredients.measure as measurement, ingredients.name").Joins("inner join item_postavkas on item_postavkas.item_id = inventarization_items.item_id inner join postavkas on postavkas.id = item_postavkas.postavka_id left join ingredients on ingredients.ingredient_id = inventarization_items.item_id and ingredients.shop_id = ?", sklad.ShopID).Where("postavkas.time between inventarization_items.before_time and inventarization_items.time and inventarization_items.type = ? and inventarization_items.item_id IN (?) and inventarization_items.inventarization_id = ? and postavkas.sklad_id = inventarization_items.sklad_id and item_postavkas.deleted = ? and item_postavkas.type = ?", utils.TypeIngredient, ids, item.InventarizationID, false, utils.TypeIngredient).Order("postavkas.time asc").Scan(&inventarizationDetailsIncome).Error
		if err != nil {
			return nil, err
		}
	} else {
		err := r.gormDB.Model(&model.InventarizationItem{}).Select("postavkas.time, item_postavkas.quantity, item_postavkas.cost, (item_postavkas.quantity * item_postavkas.cost) as sum, CASE WHEN inventarization_items.type = 'ingredient' THEN ingredients.measure Else tovars.measure END as measurement").Joins("inner join item_postavkas on item_postavkas.item_id = inventarization_items.item_id inner join postavkas on postavkas.id = item_postavkas.postavka_id left join tovars on tovars.tovar_id = inventarization_items.item_id and tovars.shop_id = ? left join ingredients on ingredients.ingredient_id = inventarization_items.item_id and ingredients.shop_id = ?", sklad.ShopID, sklad.ShopID).Where("postavkas.time between inventarization_items.before_time and inventarization_items.time and inventarization_items.type = item_postavkas.type and inventarization_items.id = ? and postavkas.sklad_id = inventarization_items.sklad_id and item_postavkas.deleted = ?", id, false).Order("postavkas.time asc").Scan(&inventarizationDetailsIncome).Error
		if err != nil {
			return nil, err
		}
	}

	return inventarizationDetailsIncome, nil
}
func (r *SkladDB) GetInventarizationDetailsExpence(id int) ([]*model.InventarizationDetailsExpence, error) {
	inventarizationDetailsExpence := []*model.InventarizationDetailsExpence{}
	item := &model.InventarizationItem{}
	err := r.gormDB.Model(&model.InventarizationItem{}).Where("id = ?", id).Scan(&item).Error
	if err != nil {
		return nil, err
	}
	if item.Type == utils.TypeIngredient {
		err = r.gormDB.Model(&model.ExpenceIngredient{}).Select("checks.closed_at as time, check_tech_carts.tech_cart_name as name, expence_ingredients.quantity, ingredients.measure as measurement").Joins("inner join check_tech_carts on check_tech_carts.id = expence_ingredients.check_tech_cart_id inner join checks on checks.id = check_tech_carts.check_id inner join ingredients on ingredients.ingredient_id = expence_ingredients.ingredient_id").Where("checks.closed_at between ? and ? and expence_ingredients.ingredient_id = ? and expence_ingredients.sklad_id = ? and expence_ingredients.status = ? and ingredients.shop_id = checks.shop_id", item.BeforeTime, item.Time, item.ItemID, item.SkladID, utils.StatusClosed).Order("checks.closed_at asc").Scan(&inventarizationDetailsExpence).Error
		if err != nil {
			return nil, err
		}
	} else if item.Type == utils.TypeTovar {
		err = r.gormDB.Model(&model.ExpenceTovar{}).Select("checks.closed_at as time, check_tovars.tovar_name as name, expence_tovars.quantity, tovars.measure as measurement").Joins("inner join check_tovars on check_tovars.id = expence_tovars.check_tovar_id inner join checks on checks.id = check_tovars.check_id inner join tovars on tovars.tovar_id = expence_tovars.tovar_id").Where("checks.closed_at between ? and ? and expence_tovars.tovar_id = ? and expence_tovars.sklad_id = ? and expence_tovars.status = ? and tovars.shop_id = checks.shop_id", item.BeforeTime, item.Time, item.ItemID, item.SkladID, utils.StatusClosed).Order("checks.closed_at asc").Scan(&inventarizationDetailsExpence).Error
		if err != nil {
			return nil, err
		}
	} else if item.Type == utils.TypeGroup {
		items := []*model.InventarizationItem{}
		err = r.gormDB.Model(&model.InventarizationItem{}).Where("inventarization_id = ? and group_id = ?", item.InventarizationID, item.ItemID).Scan(&items).Error
		if err != nil {
			return nil, err
		}
		ids := []int{}
		for _, val := range items {
			ids = append(ids, val.ItemID)
		}
		err = r.gormDB.Model(&model.ExpenceIngredient{}).Select("checks.closed_at as time, check_tech_carts.tech_cart_name as name, expence_ingredients.quantity, ingredients.measure as measurement").Joins("inner join check_tech_carts on check_tech_carts.id = expence_ingredients.check_tech_cart_id inner join checks on checks.id = check_tech_carts.check_id inner join ingredients on ingredients.ingredient_id = expence_ingredients.ingredient_id").Where("checks.closed_at between ? and ? and expence_ingredients.ingredient_id IN (?) and expence_ingredients.sklad_id = ? and expence_ingredients.status = ? and ingredients.shop_id = checks.shop_id", item.BeforeTime, item.Time, ids, item.SkladID, utils.StatusClosed).Order("checks.closed_at asc").Scan(&inventarizationDetailsExpence).Error
		if err != nil {
			return nil, err
		}
	}
	return inventarizationDetailsExpence, nil
}

func (r *SkladDB) GetInventarizationDetailsSpisanie(id int) ([]*model.InventarizationDetailsSpisanie, error) {
	inventarizationDetailsSpisanie := []*model.InventarizationDetailsSpisanie{}
	item := &model.InventarizationItem{}
	err := r.gormDB.Model(&model.InventarizationItem{}).Where("id = ?", id).Scan(&item).Error
	if err != nil {
		return nil, err
	}
	sklad := &model.Sklad{}
	err = r.gormDB.Model(&model.Sklad{}).Where("id = ?", item.SkladID).Scan(&sklad).Error
	if err != nil {
		return nil, err
	}
	if item.Type == utils.TypeGroup {
		items := []*model.InventarizationItem{}
		err = r.gormDB.Model(&model.InventarizationItem{}).Where("inventarization_id = ? and group_id = ?", item.InventarizationID, item.ItemID).Scan(&items).Error
		if err != nil {
			return nil, err
		}
		ids := []int{}
		for _, val := range items {
			ids = append(ids, val.ItemID)
		}
		err := r.gormDB.Model(&model.InventarizationItem{}).Select("remove_from_sklads.time, remove_from_sklad_items.quantity, remove_from_sklad_items.cost, (remove_from_sklad_items.quantity * remove_from_sklad_items.cost) as sum, ingredients.measure as measurement, ingredients.name").Joins("inner join remove_from_sklad_items on remove_from_sklad_items.item_id = inventarization_items.item_id inner join remove_from_sklads on remove_from_sklads.id = remove_from_sklad_items.remove_id left join ingredients on ingredients.ingredient_id = inventarization_items.item_id and ingredients.shop_id = ?", sklad.ShopID).Where("remove_from_sklads.time between inventarization_items.before_time and inventarization_items.time and inventarization_items.type = ? and inventarization_items.item_id IN (?) and remove_from_sklads.sklad_id = inventarization_items.sklad_id and inventarization_items.inventarization_id = ?", utils.TypeIngredient, ids, item.InventarizationID).Order("remove_from_sklads.time asc").Scan(&inventarizationDetailsSpisanie).Error
		if err != nil {
			return nil, err
		}
	} else {
		err := r.gormDB.Model(&model.InventarizationItem{}).Select("remove_from_sklads.time, remove_from_sklad_items.quantity, remove_from_sklad_items.cost, (remove_from_sklad_items.quantity * remove_from_sklad_items.cost) as sum, CASE WHEN inventarization_items.type = 'ingredient' THEN ingredients.measure Else tovars.measure END as measurement").Joins("inner join remove_from_sklad_items on remove_from_sklad_items.item_id = inventarization_items.item_id inner join remove_from_sklads on remove_from_sklads.id = remove_from_sklad_items.remove_id  left join tovars on tovars.tovar_id = inventarization_items.item_id and tovars.shop_id = ? left join ingredients on ingredients.ingredient_id = inventarization_items.item_id and ingredients.shop_id = ?", sklad.ShopID, sklad.ShopID).Where("remove_from_sklads.time between inventarization_items.before_time and inventarization_items.time and inventarization_items.type = remove_from_sklad_items.type and inventarization_items.id = ? and remove_from_sklads.sklad_id = inventarization_items.sklad_id", id).Order("remove_from_sklads.time asc").Scan(&inventarizationDetailsSpisanie).Error
		if err != nil {
			return nil, err
		}
	}

	return inventarizationDetailsSpisanie, nil
}

func (r *SkladDB) CheckUnique(itemID, skladID int, itemType string, groupID int) error {
	var count int64
	err := r.gormDB.Model(&model.InventarizationGroupItem{}).Joins("inner join inventarization_groups on inventarization_groups.id = inventarization_group_items.group_id").Where("inventarization_group_items.item_id = ? and inventarization_group_items.sklad_id = ? and inventarization_groups.type = ? and inventarization_groups.id != ?", itemID, skladID, itemType, groupID).Count(&count).Error
	if err != nil {
		return err
	}

	if count > 0 {
		return errors.New("not unique")
	}
	return nil
}

func (r *SkladDB) CreateInventarizationGroup(groupToAdd *model.InventarizationGroup) error {
	err := r.gormDB.Create(groupToAdd).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *SkladDB) GetAllInventarizationGroup(filter *model.Filter) ([]*model.InventarizationGroupResponse, error) {
	inventarizationGroups := []*model.InventarizationGroupResponse{}
	if len(filter.SkladID) == 0 {
		err := r.gormDB.Model(&model.InventarizationGroup{}).Select("inventarization_groups.id, inventarization_groups.name,inventarization_groups.type, inventarization_groups.measure, sklads.name as sklad_name").Joins("inner join sklads on sklads.id = inventarization_groups.sklad_id inner join shops on shops.id = sklads.shop_id").Where("shops.id IN (?)", filter.AccessibleShops).Scan(&inventarizationGroups).Error
		if err != nil {
			return nil, err
		}
	} else {
		err := r.gormDB.Model(&model.InventarizationGroup{}).Select("inventarization_groups.id, inventarization_groups.name,inventarization_groups.type, inventarization_groups.measure, sklads.name as sklad_name").Joins("inner join sklads on sklads.id = inventarization_groups.sklad_id").Where("sklads.id IN (?)", filter.SkladID).Scan(&inventarizationGroups).Error
		if err != nil {
			return nil, err
		}
	}

	return inventarizationGroups, nil
}
func (r *SkladDB) GetInventarizationGroup(filter *model.Filter, id int) (*model.InventarizationGroupResponse, error) {
	inventarizationGroup := &model.InventarizationGroupResponse{}
	err := r.gormDB.Model(&model.InventarizationGroup{}).Select("inventarization_groups.id, inventarization_groups.name, inventarization_groups.measure, inventarization_groups.type, inventarization_groups.sklad_id, sklads.name as sklad_name").Joins("inner join sklads on sklads.id = inventarization_groups.sklad_id").Where("inventarization_groups.id = ?", id).Scan(inventarizationGroup).Error
	if err != nil {
		return nil, err
	}
	sklad := &model.Sklad{}
	err = r.gormDB.Model(&model.Sklad{}).Where("id = ?", inventarizationGroup.SkladID).Scan(sklad).Error
	if err != nil {
		return nil, err
	}
	items := []*model.InventarizationGroupItemResponse{}
	err = r.gormDB.Model(&model.InventarizationGroupItem{}).Select("inventarization_group_items.id,inventarization_group_items.item_id, inventarization_group_items.group_id, inventarization_group_items.sklad_id, CASE WHEN inventarization_groups.type = 'ingredient' THEN ingredients.name ELSE tovars.name END").Joins("inner join inventarization_groups on inventarization_groups.id = inventarization_group_items.group_id left join tovars on tovars.tovar_id = inventarization_group_items.item_id and tovars.shop_id = ? left join ingredients on ingredients.ingredient_id = inventarization_group_items.item_id and ingredients.shop_id = ?", sklad.ShopID, sklad.ShopID).Where("inventarization_group_items.group_id = ?", id).Scan(&items).Error
	if err != nil {
		return nil, err
	}
	inventarizationGroup.Items = items

	return inventarizationGroup, nil
}
func (r *SkladDB) UpdateInventarizationGroup(group *model.InventarizationGroup) error {
	err := r.gormDB.Model(&model.InventarizationGroupItem{}).Where("group_id = ?", group.ID).Delete(&model.InventarizationGroupItem{}).Error
	if err != nil {
		return err
	}
	err = r.gormDB.Model(&model.InventarizationGroup{}).Where("id = ?", group.ID).Save(group).Error
	if err != nil {
		return err
	}
	for _, val := range group.Items {
		val.GroupID = group.ID
		err = r.gormDB.Model(&model.InventarizationGroupItem{}).Create(val).Error
		if err != nil {
			return err
		}
	}

	return nil
}
func (r *SkladDB) DeleteInventarizationGroup(id int) error {
	err := r.gormDB.Delete(&model.InventarizationGroup{}, id).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *SkladDB) GetPureInventarizationGroup(id int) (*model.InventarizationGroup, error) {
	group := &model.InventarizationGroup{}
	err := r.gormDB.Model(&model.InventarizationGroup{}).Where("id = ?", id).Scan(group).Error
	if err != nil {
		return nil, err
	}
	items := []*model.InventarizationGroupItem{}
	err = r.gormDB.Model(&model.InventarizationGroupItem{}).Where("group_id = ?", id).Scan(&items).Error
	if err != nil {
		return nil, err
	}
	group.Items = items
	return group, nil
}

func (r *SkladDB) AddTrafficReportJob(items []*model.AsyncJob) error {
	for _, val := range items {
		var count int64
		res := r.gormDB.Model(&model.AsyncJob{}).Where("item_id = ? and item_type = ? and status = ? and retry_count < ? and time_stamp::date = ?::date and sklad_id = ?", val.ItemID, val.ItemType, utils.StatusNeedRecalculate, utils.RetryCount, val.TimeStamp, val.SkladID).Count(&count)
		if res.Error != nil {
			return res.Error
		}
		if count != 0 {
			continue
		}
		err := r.gormDB.Model(&model.AsyncJob{}).Create(val).Error
		if err != nil {
			return err
		}
	}
	return nil
}

func (r *SkladDB) GetItemsForRecalculateTrafficReport() ([]*model.AsyncJob, error) {
	items := []*model.AsyncJob{}
	res := r.gormDB.Model(&model.AsyncJob{}).Where("status = ?", utils.StatusNeedRecalculate).Scan(&items)
	if res.Error != nil {
		return nil, res.Error
	}
	return items, nil
}

func (r *SkladDB) UpdateTrafficReportJob(val *model.AsyncJob) error {
	res := r.gormDB.Model(&model.AsyncJob{}).Where("id = ?", val.ID).Updates(val)
	if res.Error != nil {
		return res.Error
	}
	return nil
}

func (r *SkladDB) GetShopIDBySkladID(id int) (int, error) {
	var shopID int
	res := r.gormDB.Model(&model.Sklad{}).Select("sklads.shop_id").Where("sklads.id = ?", id).Scan(&shopID)
	if res.Error != nil {
		return 0, res.Error
	}
	return shopID, nil
}

// Daily statitstics update
// 1. if daily time is 24 get last invent FROM 24 00 - 25 00. -> quantity
// 2. Sales of teech carts
func (r *SkladDB) GetPostavkaByTransferID(id int) (*model.Postavka, error) {
	postavka := &model.Postavka{}
	res := r.gormDB.Model(&model.Postavka{}).Where("transfer_id = ?", id).Scan(postavka)
	if res.Error != nil {
		return nil, res.Error
	}
	return postavka, nil
}

func (r *SkladDB) GetSpisanieByTransferID(id int) (*model.RemoveFromSklad, error) {
	spisanie := &model.RemoveFromSklad{}
	res := r.gormDB.Model(&model.RemoveFromSklad{}).Where("transfer_id = ?", id).Scan(spisanie)
	if res.Error != nil {
		return nil, res.Error
	}
	return spisanie, nil
}

func (r *SkladDB) RecalculateDailyStatisticByDate(date time.Time, skladID int) error {
	from := time.Date(date.Year(), date.Month(), date.Day()-1, 18, 0, 0, 0, time.UTC)
	to := time.Date(date.Year(), date.Month(), date.Day(), 18, 0, 0, 0, time.UTC)
	shopID, err := r.GetShopIDBySkladID(skladID)
	if err != nil {
		return err
	}
	daily := []*model.DailyStatistic{}
	res := r.gormDB.Model(&model.DailyStatistic{}).Where("sklad_id = ? and date = ?", skladID, date).Scan(&daily)
	if res.Error != nil {
		if res.Error != gorm.ErrRecordNotFound {
			return res.Error
		}
	}
	flag := false
	if len(daily) == 0 {
		flag = true
		oldDaily := []*model.DailyStatistic{}
		err = r.gormDB.Model(&model.DailyStatistic{}).Where("sklad_id = ? and date = ?", skladID, date.AddDate(0, 0, -1)).Scan(&oldDaily).Error
		if err != nil {
			if err != gorm.ErrRecordNotFound {
				return err
			}
		}
		daily = []*model.DailyStatistic{}
		for _, val := range oldDaily {
			newDaily := &model.DailyStatistic{
				Date:     date,
				SkladID:  skladID,
				ShopID:   shopID,
				Type:     val.Type,
				ItemID:   val.ItemID,
				ItemName: val.ItemName,
				Cost:     val.Cost,
				Quantity: val.Quantity,
			}
			daily = append(daily, newDaily)
		}
		if len(daily) == 0 {
			return nil
		}
		err = r.gormDB.Create(daily).Error
		if err != nil {
			return err
		}
	}
	for _, val := range daily {
		newDaily := &model.DailyStatistic{}
		if val.Type == utils.TypeIngredient {
			err := r.gormDB.Model(&model.Postavka{}).Select("SUM(item_postavkas.quantity) as postavka, SUM(item_postavkas.quantity * item_postavkas.cost) as postavka_cost").Joins("inner join item_postavkas on item_postavkas.postavka_id = postavkas.id").Where("postavkas.sklad_id = ? and postavkas.time >= ? and postavkas.time < ? and item_postavkas.item_id = ? and item_postavkas.type = ? and postavkas.deleted = false", skladID, from, to, val.ItemID, utils.TypeIngredient).Scan(&newDaily).Error
			if err != nil {
				if err != gorm.ErrRecordNotFound {
					return err
				}
			}
			err = r.gormDB.Model(&model.RemoveFromSklad{}).Select("SUM(remove_from_sklad_items.quantity) as remove_from_sklad").Joins("inner join remove_from_sklad_items on remove_from_sklad_items.remove_id = remove_from_sklads.id").Where("remove_from_sklads.sklad_id = ? and remove_from_sklads.time >= ? and remove_from_sklads.time < ? and remove_from_sklad_items.item_id = ? and remove_from_sklad_items.type = ? and remove_from_sklads.deleted = false", skladID, from, to, val.ItemID, utils.TypeIngredient).Scan(&newDaily).Error
			if err != nil {
				if err != gorm.ErrRecordNotFound {
					return err
				}
			}
			err = r.gormDB.Model(&model.InventarizationItem{}).Select("SUM(inventarization_items.difference) as inventarization").Joins("inner join inventarizations on inventarizations.id = inventarization_items.inventarization_id").Where("inventarization_items.sklad_id = ? and inventarization_items.time >= ? and inventarization_items.time < ? and inventarization_items.item_id = ? and inventarization_items.type = ? and inventarization_items.status = ? and inventarizations.deleted = false", skladID, from, to, val.ItemID, utils.TypeIngredient, utils.StatusClosed).Scan(&newDaily).Error
			if err != nil {
				if err != gorm.ErrRecordNotFound {
					return err
				}
			}
			err = r.gormDB.Model(&model.ExpenceIngredient{}).Select("SUM(expence_ingredients.quantity) as sales").Where("expence_ingredients.sklad_id = ? and expence_ingredients.time >= ? and expence_ingredients.time < ? and expence_ingredients.ingredient_id = ? and expence_ingredients.status = ?", skladID, from, to, val.ItemID, utils.StatusClosed).Scan(&newDaily).Error
			if err != nil {
				if err != gorm.ErrRecordNotFound {
					return err
				}
			}
			if (val.Quantity+val.Postavka) != 0 && flag {
				val.Cost = ((val.Cost * val.Quantity) + val.PostavkaCost) / (val.Quantity + val.Postavka)
			}
			err = r.RecalculateQuantity(to, val, skladID)
			if err != nil {
				if err != gorm.ErrRecordNotFound {
					return err
				}
			}
		} else if val.Type == utils.TypeTovar {
			err = r.gormDB.Model(&model.Postavka{}).Select("SUM(item_postavkas.quantity) as postavka, SUM(item_postavkas.quantity * item_postavkas.cost) as postavka_cost").Joins("inner join item_postavkas on item_postavkas.postavka_id = postavkas.id").Where("postavkas.sklad_id = ? and postavkas.time >= ? and postavkas.time < ? and item_postavkas.item_id = ? and item_postavkas.type = ? and postavkas.deleted = false", skladID, from, to, val.ItemID, utils.TypeTovar).Scan(&newDaily).Error
			if err != nil {
				if err != gorm.ErrRecordNotFound {
					return err
				}
			}
			err = r.gormDB.Model(&model.RemoveFromSklad{}).Select("SUM(remove_from_sklad_items.quantity) as remove_from_sklad").Joins("inner join remove_from_sklad_items on remove_from_sklad_items.remove_id = remove_from_sklads.id").Where("remove_from_sklads.sklad_id = ? and remove_from_sklads.time >= ? and remove_from_sklads.time < ? and remove_from_sklad_items.item_id = ? and remove_from_sklad_items.type = ? and remove_from_sklads.deleted = false", skladID, from, to, val.ItemID, utils.TypeTovar).Scan(&newDaily).Error
			if err != nil {
				if err != gorm.ErrRecordNotFound {
					return err
				}
			}
			err = r.gormDB.Model(&model.InventarizationItem{}).Select("SUM(inventarization_items.difference) as inventarization").Joins("inner join inventarizations on inventarizations.id = inventarization_items.inventarization_id").Where("inventarization_items.sklad_id = ? and inventarization_items.time >= ? and inventarization_items.time < ? and inventarization_items.item_id = ? and inventarization_items.type = ? and inventarization_items.status = ? and inventarizations.deleted = false", skladID, from, to, val.ItemID, utils.TypeTovar, utils.StatusClosed).Scan(&newDaily).Error
			if err != nil {
				if err != gorm.ErrRecordNotFound {
					return err
				}
			}
			err = r.gormDB.Model(&model.ExpenceTovar{}).Select("SUM(expence_tovars.quantity) as sales").Where("expence_tovars.sklad_id = ? and expence_tovars.time >= ? and expence_tovars.time < ? and expence_tovars.tovar_id = ? and expence_tovars.status = ?", skladID, from, to, val.ItemID, utils.StatusClosed).Scan(&newDaily).Error
			if err != nil {
				if err != gorm.ErrRecordNotFound {
					return err
				}
			}
			err = r.gormDB.Model(&model.CheckTovar{}).Select("SUM(check_tovars.cost) as check_cost, SUM(check_tovars.price) as check_price").Joins("inner join checks on checks.id = check_tovars.check_id").Where("checks.shop_id = ? and checks.opened_at >= ? and checks.closed_at < ? and check_tovars.tovar_id = ? and checks.status = ?", shopID, from, to, val.ItemID, utils.StatusClosed).Scan(&newDaily).Error
			if err != nil {
				if err != gorm.ErrRecordNotFound {
					return err
				}
			}
			if (val.Quantity+val.Postavka) != 0 && flag {
				val.Cost = ((val.Cost * val.Quantity) + val.PostavkaCost) / (val.Quantity + val.Postavka)
			}
			err = r.RecalculateQuantity(to, val, skladID)
			if err != nil {
				if err != gorm.ErrRecordNotFound {
					return err
				}
			}
		}
		val.Postavka = newDaily.Postavka
		val.PostavkaCost = newDaily.PostavkaCost
		val.RemoveFromSklad = newDaily.RemoveFromSklad
		val.Inventarization = newDaily.Inventarization
		val.Sales = newDaily.Sales
		val.CheckCost = newDaily.CheckCost
		val.CheckPrice = newDaily.CheckPrice
		err := r.gormDB.Model(&model.DailyStatistic{}).Select("*").Where("id = ?", val.ID).Updates(val).Error
		if err != nil {
			return err
		}
	}
	res = r.gormDB.Model(&model.DailyStatistic{}).Where("shop_id = ? and date = ? and type = ?", shopID, date, utils.TypeTechCart).Scan(&daily)
	if res.Error != nil {
		if res.Error != gorm.ErrRecordNotFound {
			return res.Error
		}
	}
	for _, val := range daily {
		newDaily := &model.DailyStatistic{}
		if val.Type == utils.TypeTechCart {
			err := r.gormDB.Model(&model.CheckTechCart{}).Select("SUM(check_tech_carts.quantity) as sales, SUM(check_tech_carts.cost) as check_cost, SUM(check_tech_carts.price) as check_price").Joins("inner join checks on checks.id = check_tech_carts.check_id").Where("checks.shop_id = ? and checks.opened_at >= ? and checks.closed_at < ? and check_tech_carts.tech_cart_id = ? and checks.status = ?", shopID, from, to, val.ItemID, utils.StatusClosed).Scan(&newDaily).Error
			if err != nil {
				if err != gorm.ErrRecordNotFound {
					return err
				}
			}
		}
		val.CheckCost = newDaily.CheckCost
		val.CheckPrice = newDaily.CheckPrice
		err := r.gormDB.Model(&model.DailyStatistic{}).Select("*").Where("id = ?", val.ID).Updates(val).Error
		if err != nil {
			return err
		}
	}
	return nil
}

func (r *SkladDB) RecalculateQuantity(date time.Time, item *model.DailyStatistic, skladID int) error {
	type Quantity struct {
		Quantity float32
		LastTime time.Time
	}
	var quantity Quantity
	var postavka Quantity
	var spisanie Quantity
	var sales Quantity
	if item.Type == utils.TypeIngredient {
		err := r.gormDB.Model(&model.InventarizationItem{}).Select("inventarization_items.fact_quantity as quantity, inventarization_items.time as last_time").Joins("inner join inventarizations on inventarizations.id = inventarization_items.inventarization_id").Where("inventarization_items.sklad_id = ? and inventarization_items.time < ? and inventarization_items.item_id = ? and inventarization_items.type = ? and inventarization_items.status = ? and inventarizations.deleted = false", skladID, date, item.ItemID, utils.TypeIngredient, utils.StatusClosed).Order("inventarization_items.time desc").First(&quantity).Error
		if err != nil {
			if err != gorm.ErrRecordNotFound {
				return err
			}
		}
		err = r.gormDB.Model(&model.Postavka{}).Select("SUM(item_postavkas.quantity) as quantity").Joins("inner join item_postavkas on item_postavkas.postavka_id = postavkas.id").Where("postavkas.sklad_id = ? and postavkas.time >= ? and postavkas.time < ? and item_postavkas.item_id = ? and item_postavkas.type = ? and postavkas.deleted = false", skladID, quantity.LastTime, date, item.ItemID, utils.TypeIngredient).Scan(&postavka).Error
		if err != nil {
			if err != gorm.ErrRecordNotFound {
				return err
			}
		}
		err = r.gormDB.Model(&model.RemoveFromSklad{}).Select("SUM(remove_from_sklad_items.quantity) as quantity").Joins("inner join remove_from_sklad_items on remove_from_sklad_items.remove_id = remove_from_sklads.id").Where("remove_from_sklads.sklad_id = ? and remove_from_sklads.time >= ? and remove_from_sklads.time < ? and remove_from_sklad_items.item_id = ? and remove_from_sklad_items.type = ? and remove_from_sklads.deleted = false", skladID, quantity.LastTime, date, item.ItemID, utils.TypeIngredient).Scan(&spisanie).Error
		if err != nil {
			if err != gorm.ErrRecordNotFound {
				return err
			}
		}
		err = r.gormDB.Model(&model.ExpenceIngredient{}).Select("SUM(expence_ingredients.quantity) as quantity").Where("expence_ingredients.sklad_id = ? and expence_ingredients.time >= ? and expence_ingredients.time < ? and expence_ingredients.ingredient_id = ? and expence_ingredients.status = ?", skladID, quantity.LastTime, date, item.ItemID, utils.StatusClosed).Scan(&sales).Error
		if err != nil {
			if err != gorm.ErrRecordNotFound {
				return err
			}
		}
	} else if item.Type == utils.TypeTovar {
		err := r.gormDB.Model(&model.InventarizationItem{}).Select("inventarization_items.fact_quantity as quantity, inventarization_items.time as last_time").Joins("inner join inventarizations on inventarizations.id = inventarization_items.inventarization_id").Where("inventarization_items.sklad_id = ? and inventarization_items.time < ? and inventarization_items.item_id = ? and inventarization_items.type = ? and inventarization_items.status = ? and inventarizations.deleted = false", skladID, date, item.ItemID, utils.TypeTovar, utils.StatusClosed).Order("inventarization_items.time desc").First(&quantity).Error
		if err != nil {
			if err != gorm.ErrRecordNotFound {
				return err
			}
		}
		err = r.gormDB.Model(&model.Postavka{}).Select("SUM(item_postavkas.quantity) as quantity").Joins("inner join item_postavkas on item_postavkas.postavka_id = postavkas.id").Where("postavkas.sklad_id = ? and postavkas.time >= ? and postavkas.time < ? and item_postavkas.item_id = ? and item_postavkas.type = ? and postavkas.deleted = false", skladID, quantity.LastTime, date, item.ItemID, utils.TypeTovar).Scan(&postavka).Error
		if err != nil {
			if err != gorm.ErrRecordNotFound {
				return err
			}
		}
		err = r.gormDB.Model(&model.RemoveFromSklad{}).Select("SUM(remove_from_sklad_items.quantity) as quantity").Joins("inner join remove_from_sklad_items on remove_from_sklad_items.remove_id = remove_from_sklads.id").Where("remove_from_sklads.sklad_id = ? and remove_from_sklads.time >= ? and remove_from_sklads.time < ? and remove_from_sklad_items.item_id = ? and remove_from_sklad_items.type = ? and remove_from_sklads.deleted = false", skladID, quantity.LastTime, date, item.ItemID, utils.TypeTovar).Scan(&spisanie).Error
		if err != nil {
			if err != gorm.ErrRecordNotFound {
				return err
			}
		}
		err = r.gormDB.Model(&model.ExpenceTovar{}).Select("SUM(expence_tovars.quantity) as quantity").Where("expence_tovars.sklad_id = ? and expence_tovars.time >= ? and expence_tovars.time < ? and expence_tovars.tovar_id = ? and expence_tovars.status = ?", skladID, quantity.LastTime, date, item.ItemID, utils.StatusClosed).Scan(&sales).Error
		if err != nil {
			if err != gorm.ErrRecordNotFound {
				return err
			}
		}
	}
	item.Quantity = quantity.Quantity + postavka.Quantity - spisanie.Quantity - sales.Quantity
	return nil
}

func (r *SkladDB) GetAllSklads() ([]*model.Sklad, error) {
	sklads := []*model.Sklad{}
	err := r.gormDB.Model(sklads).Where("deleted = false").Scan(&sklads).Error
	if err != nil {
		return nil, err
	}
	return sklads, nil
}

func (r *SkladDB) CreateIngredientsOfTechCartForSpisanie(skladID int, from, to time.Time) error {
	type TechCart struct {
		ID       int
		RemoveID int
		Quantity float32
	}
	techCarts := []*TechCart{}
	res := r.gormDB.Model(&model.RemoveFromSklad{}).Select("remove_from_sklad_items.item_id as id, remove_from_sklads.id as remove_id, remove_from_sklad_items.quantity").Joins("inner join remove_from_sklad_items on remove_from_sklad_items.remove_id = remove_from_sklads.id").Where("remove_from_sklads.sklad_id = ? and remove_from_sklads.time >= ? and remove_from_sklads.time < ? and remove_from_sklad_items.type = ? and remove_from_sklads.deleted = false", skladID, from, to, utils.TypeTechCart).Scan(&techCarts)
	if res.Error != nil {
		return res.Error
	}
	for _, val := range techCarts {
		ingredients := []*model.IngredientTechCart{}
		res := r.gormDB.Model(&model.IngredientTechCart{}).Where("tech_cart_id = ?", val.ID).Scan(&ingredients)
		if res.Error != nil {
			return res.Error
		}
		for _, ingredient := range ingredients {
			item := &model.RemoveFromSkladItem{
				RemoveID:       val.RemoveID,
				ItemID:         ingredient.IngredientID,
				Type:           utils.TypeIngredient,
				Quantity:       ingredient.Brutto * val.Quantity,
				PartOfTechCart: true,
			}
			res := r.gormDB.Model(&model.RemoveFromSkladItem{}).Create(item)
			if res.Error != nil {
				return res.Error
			}
		}
	}
	return nil
}

func (r *SkladDB) ConcurrentRecalculationForDailyStatistics(itemsToRecalculate []*model.AsyncJob) {
	wait := sync.WaitGroup{}
	myChannel := make(chan *model.AsyncJob, len(itemsToRecalculate))
	ctx, cancel := context.WithCancel(context.Background())
	wait.Add(6)
	for i := 0; i < 6; i++ {
		go r.GoRecalcReaderForDailyStatistics(&wait, myChannel, ctx)
	}

	for _, item := range itemsToRecalculate {
		myChannel <- item
	}

	close(myChannel)
	wait.Wait()
	cancel()
}

func (r *SkladDB) GoRecalcReaderForDailyStatistics(wtg *sync.WaitGroup, myChannel chan *model.AsyncJob, ctx context.Context) {
	defer wtg.Done()
	for {
		select {
		case <-ctx.Done():
			return
		case item, ok := <-myChannel:
			if !ok {
				return
			}
			if item != nil {
				err := r.RecalculateTrafficReport(item)
				if err != nil {
					log.Print(err)
				}
			}
		}
	}
}

func (r *SkladDB) GetSumOfPostavkaForPeriod(filter *model.Filter) (float32, error) {
	var sum float32
	sklads := []int{}
	if len(filter.SkladID) > 0 {
		sklads = append(sklads, filter.SkladID...)
	} else {
		res := r.gormDB.Model(&model.Sklad{}).Select("id").Where("shop_id IN (?) and deleted = false", filter.AccessibleShops).Scan(&sklads)
		if res.Error != nil {
			return 0, res.Error
		}
	}
	filter.To = filter.To.AddDate(0, 0, 1)
	res := r.gormDB.Model(&model.Postavka{}).Select("COALESCE(SUM(item_postavkas.quantity * item_postavkas.cost), 0.0)").Joins("inner join item_postavkas on item_postavkas.postavka_id = postavkas.id").Where("postavkas.sklad_id IN (?) and postavkas.time >= ? and postavkas.time <= ? and postavkas.deleted = false", sklads, filter.From, filter.To).Scan(&sum)
	if res.Error != nil {
		return 0, res.Error
	}
	return sum, nil
}
