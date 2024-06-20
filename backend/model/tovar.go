package model

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"

	"github.com/gin-gonic/gin"
)

type ReqTovar struct {
	ID       int     `json:"id"`
	TovarID  int     `json:"tovar_id"`
	ShopID   []int   `json:"shop_id"`
	Name     string  `json:"name"`
	Category int     `json:"category" gorm:"ForeignKey:CategoryID"`
	Image    string  `json:"image"`
	Tax      string  `json:"tax"`
	Measure  string  `json:"measure"`
	Cost     float32 `json:"cost" gorm:"-"`
	Price    float32 `json:"price" gorm:"column:price;not null"`
	Discount bool    `json:"discount" gorm:"default:false"`
}

func (p *ReqTovar) ParseRequest(c *gin.Context) error {
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

type Tovar struct {
	ID        int     `json:"id"`
	TovarID   int     `json:"tovar_id"`
	ShopID    int     `json:"shop_id"`
	Name      string  `json:"name"`
	Category  int     `json:"category" gorm:"ForeignKey:CategoryID"`
	Image     string  `json:"image"`
	Tax       string  `json:"tax"`
	Cost      float32 `json:"cost" gorm:"-"`
	Measure   string  `json:"measure"`
	Price     float32 `json:"price" gorm:"column:price;not null"`
	Deleted   bool    `json:"deleted" gorm:"default:false"`
	Discount  bool    `json:"discount" gorm:"default:false"`
	IsVisible bool    `json:"is_visible" gorm:"default:true"`
}

type TovarResponse struct {
	ID        int     `json:"id"`
	TovarID   int     `json:"tovar_id"`
	ShopID    int     `json:"shop_id"`
	Name      string  `json:"name"`
	Category  int     `json:"category" gorm:"ForeignKey:CategoryID"`
	Image     string  `json:"image"`
	Tax       string  `json:"tax"`
	Measure   string  `json:"measure"`
	Cost      float32 `json:"cost" gorm:"column:cost;not null"`
	Price     float32 `json:"price" gorm:"column:price;not null"`
	Profit    float32 `json:"profit" gorm:"column:profit;not null"`
	Margin    float32 `json:"margin" gorm:"column:margin;not null"`
	Deleted   bool    `json:"deleted" gorm:"default:false"`
	Discount  bool    `json:"discount" gorm:"default:false"`
	IsVisible bool    `json:"is_visible" gorm:"default:true"`
}

type CategoryTovar struct {
	ID      int    `json:"id"`
	Name    string `json:"name"`
	Image   string `json:"image"`
	Deleted bool   `json:"deleted" gorm:"default:false"`
}

func (p *CategoryTovar) ParseRequest(c *gin.Context) error {
	body, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		return err
	}
	c.Request.Body = &CustomReadCloser{reader: bytes.NewReader(body)}
	if err := c.BindJSON(&p); err != nil {
		return errors.New("bad request | " + err.Error())
	}
	c.Request.Body = &CustomReadCloser{reader: bytes.NewReader(body)}

	if p.Name == "" {
		return errors.New("bad request | fill fields properly")
	}

	return nil
}

type TovarOutput struct {
	ID         int     `json:"id"`
	TovarID    int     `json:"tovar_id"`
	ShopID     int     `json:"shop_id"`
	ShopName   string  `json:"shop_name"`
	Name       string  `json:"name"`
	Category   string  `json:"category" gorm:"ForeignKey:CategoryID"`
	CategoryID int     `json:"category_id" gorm:"ForeignKey:CategoryID"`
	Image      string  `json:"image"`
	Tax        string  `json:"tax"`
	Measure    string  `json:"measure"`
	Cost       float32 `json:"cost" gorm:"column:cost;not null"`
	Price      float32 `json:"price" gorm:"column:price;not null"`
	Profit     float32 `json:"profit" gorm:"column:profit;not null"`
	Margin     float32 `json:"margin" gorm:"column:margin;not null"`
	Deleted    bool    `json:"deleted" gorm:"default:false"`
	Discount   bool    `json:"discount" gorm:"default:false"`
	IsVisible  bool    `json:"is_visible" gorm:"default:true"`
}

func (p *TovarOutput) ParseRequest(c *gin.Context) error {
	body, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		return err
	}
	c.Request.Body = &CustomReadCloser{reader: bytes.NewReader(body)}
	if err := c.BindJSON(&p); err != nil {
		return errors.New("bad request | " + err.Error())
	}
	c.Request.Body = &CustomReadCloser{reader: bytes.NewReader(body)}

	return nil
}

type TechCart struct {
	ID          int                   `json:"id" gorm:"primaryKey"`
	TechCartID  int                   `json:"tech_cart_id" gorm:"ForeignKey:TechCartID;autoIncrement:false"`
	ShopID      int                   `json:"shop_id" gorm:"ForeignKey:ShopID;autoIncrement:false"`
	Name        string                `json:"name"`
	Category    int                   `json:"category" gorm:"ForeignKey:CategoryID"`
	Image       string                `json:"image"`
	Tax         string                `json:"tax"`
	Measure     string                `json:"measure"`
	Price       float32               `json:"price"`
	Deleted     bool                  `json:"deleted" gorm:"default:false"`
	Discount    bool                  `json:"discount" gorm:"default:false"`
	IsVisible   bool                  `json:"is_visible" gorm:"default:true"`
	Ingredients []*IngredientTechCart `json:"ingredient_tech_cart" gorm:"ForeignKey:TechCartID"`
	Nabor       []*NaborTechCart      `json:"nabor" gorm:"ForeignKey:TechCartID"`
}

type ReqTechCart struct {
	ID          int                   `json:"id"`
	TechCartID  int                   `json:"tech_cart_id"`
	ShopID      []int                 `json:"shop_id"`
	Name        string                `json:"name"`
	Category    int                   `json:"category" gorm:"ForeignKey:CategoryID"`
	Image       string                `json:"image"`
	Tax         string                `json:"tax"`
	Measure     string                `json:"measure"`
	Price       float32               `json:"price"`
	Deleted     bool                  `json:"deleted" gorm:"default:false"`
	Discount    bool                  `json:"discount" gorm:"default:false"`
	Ingredients []*IngredientTechCart `json:"ingredient_tech_cart" gorm:"ForeignKey:TechCartID"`
	Nabor       []*NaborTechCart      `json:"nabor" gorm:"ForeignKey:TechCartID"`
}

func (p *ReqTechCart) ParseRequest(c *gin.Context) error {
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

type NaborTechCart struct {
	TechCartID int `json:"tech_cart_id"`
	NaborID    int `json:"nabor_id"`
	ShopID     int `json:"shop_id"`
}

type IngredientTechCart struct {
	TechCartID   int     `json:"tech_cart_id"`
	IngredientID int     `json:"ingredient_id"`
	ShopID       int     `json:"shop_id"`
	Brutto       float32 `json:"brutto"`
}
type TechCartResponse struct {
	ID         int     `json:"id"`
	TechCartID int     `json:"tech_cart_id"`
	ShopID     int     `json:"shop_id"`
	ShopName   string  `json:"shop_name"`
	Name       string  `json:"name"`
	Category   string  `json:"category"`
	CategoryID int     `json:"category_id"`
	Image      string  `json:"image"`
	Tax        string  `json:"tax"`
	Measure    string  `json:"measure"`
	Cost       float32 `json:"cost"`
	Price      float32 `json:"price"`
	Profit     float32 `json:"profit"`
	Margin     float32 `json:"margin"`
	Discount   bool    `json:"discount"`
}

func (p *TechCartResponse) ParseRequest(c *gin.Context) error {
	body, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		return err
	}
	c.Request.Body = &CustomReadCloser{reader: bytes.NewReader(body)}
	if err := c.BindJSON(&p); err != nil {
		return errors.New("bad request | " + err.Error())
	}
	c.Request.Body = &CustomReadCloser{reader: bytes.NewReader(body)}

	return nil
}

type TechCartInfo struct {
	ID         int     `json:"id"`
	Name       string  `json:"name"`
	Category   string  `json:"category"`
	CategoryID int     `json:"category_id"`
	Image      string  `json:"image"`
	Tax        string  `json:"tax"`
	Measure    string  `json:"measure"`
	Cost       float32 `json:"cost"`
	Price      float32 `json:"price"`
	Profit     float32 `json:"profit"`
	Margin     float32 `json:"margin"`
	Discount   bool    `json:"discount"`
}

type TechCartOutput struct {
	ID          int                    `json:"id"`
	TechCartID  int                    `json:"tech_cart_id"`
	ShopID      int                    `json:"shop_id"`
	Shops       []int                  `json:"shops"`
	Name        string                 `json:"name"`
	Category    string                 `json:"category"`
	CategoryID  int                    `json:"category_id"`
	Image       string                 `json:"image"`
	Tax         string                 `json:"tax"`
	Measure     string                 `json:"measure"`
	Cost        float32                `json:"cost"`
	Price       float32                `json:"price"`
	Profit      float32                `json:"profit"`
	Margin      float32                `json:"margin"`
	Discount    bool                   `json:"discount"`
	Ingredients []*IngredientNumOutput `json:"ingredient_tech_cart" gorm:"-"`
	Nabors      []*NaborOutput         `json:"nabors" gorm:"-"`
}

type NaborInfo struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Min  int    `json:"min"`
	Max  int    `json:"max"`
}

type NaborOutput struct {
	ID              int                 `json:"id"`
	NaborID         int                 `json:"nabor_id"`
	ShopID          int                 `json:"shop_id"`
	ShopName        string              `json:"shop_name"`
	Name            string              `json:"name"`
	Min             int                 `json:"min"`
	Max             int                 `json:"max"`
	NaborIngredient []*IngredientOutput `json:"nabor_ingredient" gorm:"-"`
	Replaces        IntArray            `json:"replaces" gorm:"type:integer[]"`
	Deleted         bool                `json:"deleted" gorm:"default:false"`
}

type IngredientNaborOutput []*IngredientOutput

func (ls *IngredientNaborOutput) Scan(src any) error {
	var data []byte
	switch v := src.(type) {
	case string:
		data = []byte(v)
	case []byte:
		data = v
	}
	return json.Unmarshal(data, ls)
}

type IngredientNumOutput struct {
	ID           int     `json:"id"`
	IngredientID int     `json:"ingredient_id"`
	Name         string  `json:"name"`
	Measure      string  `json:"measure"`
	Image        string  `json:"image"`
	Brutto       float32 `json:"brutto"`
	Netto        float32 `json:"netto"`
	Cost         float32 `json:"cost"`
}

type IngredientNumOuputArray []*IngredientNumOutput

func (ls *IngredientNumOuputArray) Scan(src any) error {
	var data []byte
	switch v := src.(type) {
	case string:
		data = []byte(v)
	case []byte:
		data = v
	}
	return json.Unmarshal(data, ls)
}
