package model

import (
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"
	"sync"
	"time"
	"zebra/utils"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

type Filter struct {
	Sort            string    `json:"sort" bson:"sort"`
	Search          string    `json:"search" bson:"search"`
	Page            int       `json:"page" bson:"page"`
	From            time.Time `json:"from" bson:"from"`
	To              time.Time `json:"to" bson:"to"`
	Shop            []int     `json:"shop"`
	BindShop        int       `json:"bind_shop"`
	Category        int       `json:"category"`
	SkladID         []int     `json:"sklad_id"`
	WorkerID        int       `json:"worker_id"`
	SchetID         int       `json:"schet_id"`
	PaymentType     string    `json:"payment_type"` //Nalichka, Kartoi
	DealerID        int       `json:"dealer_id"`
	Type            string    `json:"type"` // tovar, techcart, ingredient
	Status          string    `json:"status"`
	AccessibleShops []int     `json:"accessible_shops"`
	Measure         string    `json:"measure"`
	Role            string    `json:"role"`
}

func (f *Filter) ParseRequest(c *gin.Context) error {
	category, err := strconv.Atoi(c.Query("category"))
	if err != nil {
		category = 0
	}
	measure := c.Query("measure")

	page, err := strconv.Atoi(c.Query("page"))
	if err != nil {
		page = 0
	}
	shopsInFilter := make([]int, 0)
	if values, ok := c.GetQueryArray("shop"); ok {
		for _, v := range values {
			shop, err := strconv.Atoi(v)
			if err != nil {
				return errors.New("shop must be integer")
			}
			shopsInFilter = append(shopsInFilter, shop)
		}
	}

	skladsInFilter := make([]int, 0)
	if values, ok := c.GetQueryArray("sklad_id"); ok {
		for _, v := range values {
			sklad, err := strconv.Atoi(v)
			if err != nil {
				return errors.New("sklad_id must be integer")
			}
			skladsInFilter = append(skladsInFilter, sklad)
		}
	}

	workerID, err := strconv.Atoi(c.Query("worker_id"))
	if err != nil {
		workerID = 0
	}

	schetID, err := strconv.Atoi(c.Query("schet_id"))
	if err != nil {
		schetID = 0
	}

	dealerID, err := strconv.Atoi(c.Query("dealer_id"))
	if err != nil {
		dealerID = 0
	}

	status := c.Query("status")
	if status != utils.StatusClosed && status != utils.StatusOpened {
		status = ""
	}

	f.PaymentType = c.Query("payment_type")
	if f.PaymentType != utils.PaymentCard && f.PaymentType != utils.PaymentCash {
		f.PaymentType = ""
	}

	f.Type = c.Query("type")
	if f.Type != utils.TypeTovar && f.Type != utils.TypeTechCart && f.Type != utils.TypeIngredient {
		f.Type = ""
	}

	now := time.Now()
	layout := "2006-01-02"

	from, err := time.Parse(layout, c.Query("from"))
	if err != nil {
		from = now.AddDate(0, -1, 0)
	}

	to, err := time.Parse(layout, c.Query("to"))
	if err != nil {
		to = now.AddDate(0, 0, 2)
	}

	sort := c.Query("sort")
	if sort != "" {
		splitArr := strings.Split(sort, "_")
		log.Print(len(splitArr))
		if len(splitArr) == 2 {
			f.Sort = splitArr[0] + " " + splitArr[1]
		}
	}

	if f.Sort == "" {
		f.Sort = "id DESC"
	}
	shops, err := getAvailableShopsModel(c)
	if err != nil {
		return err
	}
	role, err := getUserRole(c)
	if err != nil {
		return err
	}
	bindShop, err := getBindShop(c)
	if err != nil {
		return err
	}
	f.Measure = measure
	f.BindShop = bindShop
	f.Role = role
	f.AccessibleShops = shops
	f.Status = status
	f.DealerID = dealerID
	f.SchetID = schetID
	f.WorkerID = workerID
	f.SkladID = skladsInFilter
	f.Search = c.Query("search")
	f.Page = page
	f.From = from
	f.To = to
	f.Shop = shopsInFilter
	f.Category = category
	return nil
}

func (f *Filter) FilterResults(db *gorm.DB, model interface{}, pageSize int, timeField string, search string, sort string) (*gorm.DB, int64, error) {
	s, err := schema.Parse(&model, &sync.Map{}, schema.NamingStrategy{})
	if err != nil {
		return nil, 0, err
	}
	tableName := s.Table

	searchQuery := search

	shopIdExists := false

	dbCols := s.FieldsByDBName
	for _, col := range dbCols {
		if col.Name == "ShopID" {
			shopIdExists = true
		}
	}

	if search == "" {
		for _, col := range dbCols {
			if f.Search != "" {
				if col.DataType == schema.String {
					if col.Name != "category" && col.Name != "IdempotencyKey" {
						if searchQuery != "" {
							searchQuery += " OR "
						}
						searchQuery = searchQuery + fmt.Sprintf("%s.%s ilike '%%%s%%' ", tableName, col.Name, f.Search)
					}
				}
			}
		}
	}
	if len(f.Shop) != 0 {
		db = db.Where(tableName+".shop_id IN (?)", f.Shop)
	}
	if searchQuery != "" {
		db = db.Where(searchQuery)
	}

	if len(f.AccessibleShops) != 0 && shopIdExists {
		db = db.Where(tableName+".shop_id IN (?)", f.AccessibleShops)
	}

	if f.Status != "" {
		db = db.Where("status = ?", f.Status)
	}

	if f.Category != 0 {
		if tableName != "daily_statistics" && tableName != "postavkas" && tableName != "remove_from_sklads" {
			db = db.Where("category = ?", f.Category)
		}
	}

	if len(f.SkladID) != 0 {
		if tableName != "transfers" {
			skladIDExists := false
			for _, col := range dbCols {
				if col.Name == "SkladID" {
					db = db.Where(tableName+".sklad_id IN (?)", f.SkladID)
					skladIDExists = true
				}
			}
			if !skladIDExists {
				db = db.Where("sklad_id IN (?)", f.SkladID)
			}

		} else {
			db = db.Where("skladFrom.id IN (?) or skladTo.id IN (?)", f.SkladID, f.SkladID)
		}
	}

	if f.WorkerID != 0 {
		db = db.Where("worker_id = ?", f.WorkerID)
	}

	if f.SchetID != 0 {
		db = db.Where(tableName+".schet_id = ?", f.SchetID)
	}

	if f.DealerID != 0 {
		db = db.Where("dealer_id = ?", f.DealerID)
	}

	if f.PaymentType != "" {
		db = db.Where("payment = ?", f.PaymentType)
	}

	// if f.Type != "" {
	// 	db = db.Where("type = ?", f.Type)
	// }

	//opened_at::date >= ?::date and opened_at::date <= ?::date
	if !f.From.IsZero() && !f.To.IsZero() && timeField != "" {
		db = db.Where(tableName+"."+timeField+"::date BETWEEN ?::date AND ?::date", f.From, f.To)
	}
	var count int64
	newDB := *db
	newnewDB := &newDB
	newnewDB = newnewDB.Count(&count)
	if f.Sort != "" {
		if sort != "" {
			db = db.Order(sort)
		} else {
			db = db.Order(tableName + "." + f.Sort)
		}
		// db = db.Order(tableName + "." + f.Sort)
	}

	if f.Page != 0 {
		offset := (f.Page - 1) * pageSize
		db = db.Offset(offset).Limit(pageSize)
	}
	return db, count, nil
}

var accessShopsCtx = "accessShops"
var roleCtx = "role"
var bindShopCtx = "bindShop"

func getAvailableShopsModel(c *gin.Context) ([]int, error) {
	shops, ok := c.Get(accessShopsCtx)
	if !ok {
		return nil, errors.New("user id not found")
	}

	idInt, ok := shops.([]int)
	if !ok {
		return nil, errors.New("user id is of invalid type")
	}

	return idInt, nil
}

func getUserRole(c *gin.Context) (string, error) {
	role, ok := c.Get(roleCtx)
	if !ok {
		return "", errors.New("user id not found")
	}

	idInt, ok := role.(string)
	if !ok {
		return "", errors.New("user id is of invalid type")
	}

	return idInt, nil
}

func getBindShop(c *gin.Context) (int, error) {
	bindShopCtx, ok := c.Get(bindShopCtx)
	if !ok {
		return 0, errors.New("user id not found")
	}

	idInt, ok := bindShopCtx.(int)
	if !ok {
		return 0, errors.New("user id is of invalid type")
	}

	return idInt, nil
}
