package model

import (
	"bytes"
	"errors"
	"io/ioutil"

	"github.com/gin-gonic/gin"
)

type TovarMaster struct {
	ID       int     `json:"id"`
	Name     string  `json:"name"`
	Category int     `json:"category" gorm:"ForeignKey:CategoryID"`
	Image    string  `json:"image"`
	Tax      string  `json:"tax"`
	Measure  string  `json:"measure"`
	Price    float32 `json:"price" gorm:"column:price;not null"`
	Deleted  bool    `json:"deleted" gorm:"default:false"`
	Discount bool    `json:"discount" gorm:"default:false"`
	Status   string  `json:"status"`
}

type TovarMasterResponse struct {
	ID         int     `json:"id"`
	Name       string  `json:"name"`
	Category   string  `json:"category"`
	CategoryID int     `json:"category_id"`
	Image      string  `json:"image"`
	Tax        string  `json:"tax"`
	Measure    string  `json:"measure"`
	Cost       float32 `json:"cost" gorm:"column:cost;not null"`
	Price      float32 `json:"price" gorm:"column:price;not null"`
	Profit     float32 `json:"profit" gorm:"column:profit;not null"`
	Margin     float32 `json:"margin" gorm:"column:margin;not null"`
	Deleted    bool    `json:"deleted" gorm:"default:false"`
	Discount   bool    `json:"discount" gorm:"default:false"`
	Status     string  `json:"status"`
}

type IngredientMaster struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Category int    `json:"category"`
	Image    string `json:"image"`
	Measure  string `json:"measure"`
	Deleted  bool   `json:"deleted" gorm:"default:false"`
	Status   string `json:"status"`
}

type IngredientMasterResponse struct {
	ID         int     `json:"id"`
	Name       string  `json:"name"`
	Category   string  `json:"category"`
	CategoryID int     `json:"category_id"`
	Cost       float32 `json:"cost"`
	Image      string  `json:"image"`
	Measure    string  `json:"measure"`
	Deleted    bool    `json:"deleted" gorm:"default:false"`
	Status     string  `json:"status"`
}

type TechCartMaster struct {
	ID              int                         `json:"id" `
	Name            string                      `json:"name"`
	Category        int                         `json:"category" gorm:"ForeignKey:CategoryID"`
	Image           string                      `json:"image"`
	Tax             string                      `json:"tax"`
	Measure         string                      `json:"measure"`
	Cost            float32                     `json:"cost" gorm:"-"`
	Price           float32                     `json:"price"`
	Deleted         bool                        `json:"deleted" gorm:"default:false"`
	Discount        bool                        `json:"discount" gorm:"default:false"`
	Status          string                      `json:"status"`
	Ingredients     []*IngredientTechCartMaster `json:"ingredient_tech_cart" gorm:"ForeignKey:TechCartID"`
	Nabor           []*NaborTechCartMaster      `json:"nabor" gorm:"ForeignKey:TechCartID"`
	TechCarts       []*TechCart                 `json:"tech_carts" gorm:"ForeignKey:TechCartID"`
	ShopIngredients []*IngredientTechCart       `json:"shop_ingredients" gorm:"ForeignKey:TechCartID"`
}

type TechCartMasterResponse struct {
	ID          int                    `json:"id"`
	CategoryID  int                    `json:"category_id"`
	Name        string                 `json:"name"`
	Category    string                 `json:"category"`
	Image       string                 `json:"image"`
	Tax         string                 `json:"tax"`
	Measure     string                 `json:"measure"`
	Cost        float32                `json:"cost" gorm:"-"`
	Profit      float32                `json:"profit" gorm:"-"`
	Margin      float32                `json:"margin" gorm:"-"`
	Price       float32                `json:"price"`
	Deleted     bool                   `json:"deleted" gorm:"default:false"`
	Discount    bool                   `json:"discount" gorm:"default:false"`
	Status      string                 `json:"status"`
	Ingredients []*IngredientNumOutput `json:"ingredient_tech_cart" gorm:"-"`
	Nabors      []*NaborOutput         `json:"nabors" gorm:"-"`
}

type NaborTechCartMaster struct {
	TechCartID int `json:"tech_cart_id"`
	NaborID    int `json:"nabor_id"`
}

type IngredientTechCartMaster struct {
	TechCartID   int     `json:"tech_cart_id" gorm:"ForeignKey:TechCartID"`
	IngredientID int     `json:"ingredient_id" gorm:"ForeignKey:IngredientID"`
	Brutto       float32 `json:"brutto"`
}

type NaborMaster struct {
	ID              int                      `json:"id"`
	Name            string                   `json:"name"`
	Min             int                      `json:"min"`
	Max             int                      `json:"max"`
	Ingredients     []*IngredientNaborMaster `json:"ingredient_nabor" gorm:"ForeignKey:NaborID"`
	Replaces        IntArray                 `json:"replaces" gorm:"type:integer[]"`
	Status          string                   `json:"status"`
	Deleted         bool                     `json:"deleted" gorm:"default:false"`
	Nabors          []*Nabor                 `json:"nabors" gorm:"ForeignKey:NaborID"`
	ShopIngredients []*IngredientNabor       `json:"-" gorm:"ForeignKey:NaborID"`
}

type IngredientNaborMaster struct {
	ID           int     `json:"id" gorm:"primary_key"`
	NaborID      int     `json:"nabor_id"`
	IngredientID int     `json:"ingredient_id" gorm:"ForeignKey:IngredientID"`
	Brutto       float32 `json:"brutto"`
	Price        float32 `json:"price"`
}

type ReqTovarMaster struct {
	ID       int     `json:"id"`
	Name     string  `json:"name"`
	Category int     `json:"category" gorm:"ForeignKey:CategoryID"`
	Image    string  `json:"image"`
	Tax      string  `json:"tax"`
	Measure  string  `json:"measure"`
	Price    float32 `json:"price" gorm:"column:price;not null"`
	Discount bool    `json:"discount" gorm:"default:false"`
}

func (p *ReqTovarMaster) ParseRequest(c *gin.Context) error {
	body, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		return err
	}
	c.Request.Body = &CustomReadCloser{reader: bytes.NewReader(body)}
	if err := c.BindJSON(&p); err != nil {
		return errors.New("bad request | " + err.Error())
	}
	c.Request.Body = &CustomReadCloser{reader: bytes.NewReader(body)}

	if p.Name == "" || p.Category <= 0 || p.Measure == "" || p.Price <= 0 {
		return errors.New("bad request | fill fields properly")
	}
	return nil
}

type ReqTechCartMaster struct {
	ID          int                         `json:"id"`
	Name        string                      `json:"name"`
	Category    int                         `json:"category" gorm:"ForeignKey:CategoryID"`
	Image       string                      `json:"image"`
	Tax         string                      `json:"tax"`
	Measure     string                      `json:"measure"`
	Price       float32                     `json:"price"`
	Deleted     bool                        `json:"deleted" gorm:"default:false"`
	Discount    bool                        `json:"discount" gorm:"default:false"`
	Ingredients []*IngredientTechCartMaster `json:"ingredient_tech_cart" gorm:"ForeignKey:TechCartID"`
	Nabor       []*NaborTechCartMaster      `json:"nabor" gorm:"ForeignKey:TechCartID"`
}

func (p *ReqTechCartMaster) ParseRequest(c *gin.Context) error {
	body, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		return err
	}
	c.Request.Body = &CustomReadCloser{reader: bytes.NewReader(body)}
	if err := c.BindJSON(&p); err != nil {
		return errors.New("bad request | " + err.Error())
	}
	c.Request.Body = &CustomReadCloser{reader: bytes.NewReader(body)}

	if p.Name == "" || p.Category <= 0 || p.Measure == "" || p.Price < 0 {
		return errors.New("bad request | fill fields properly")
	}

	return nil
}

type ReqIngredientMaster struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	Category  int    `json:"category"`
	Image     string `json:"image"`
	Measure   string `json:"measure"`
	IsVisible bool   `json:"is_visible" gorm:"default:true"`
}

func (p *ReqIngredientMaster) ParseRequest(c *gin.Context) error {

	body, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		return err
	}
	c.Request.Body = &CustomReadCloser{reader: bytes.NewReader(body)}
	if err := c.BindJSON(&p); err != nil {
		return errors.New("bad request | " + err.Error())
	}
	c.Request.Body = &CustomReadCloser{reader: bytes.NewReader(body)}

	if p.Name == "" || p.Category <= 0 || p.Measure == "" {
		return errors.New("bad request | fill fields properly")
	}

	return nil
}

type NaborMasterOutput struct {
	ID          int                            `json:"id"`
	Name        string                         `json:"name"`
	Min         int                            `json:"min"`
	Max         int                            `json:"max"`
	Ingredients []*IngredientNaborMasterOutput `json:"ingredient_nabor" gorm:"ForeignKey:NaborID"`
	Replaces    IntArray                       `json:"replaces" gorm:"type:integer[]"`
	Status      string                         `json:"status"`
	Deleted     bool                           `json:"deleted" gorm:"default:false"`
}

type IngredientNaborMasterOutput struct {
	ID           int     `json:"id" gorm:"primary_key"`
	NaborID      int     `json:"nabor_id"`
	IngredientID int     `json:"ingredient_id" gorm:"ForeignKey:IngredientID"`
	Brutto       float32 `json:"brutto"`
	Price        float32 `json:"price"`
	Name         string  `json:"name"`
	Measure      string  `json:"measure"`
	Image        string  `json:"image"`
}
