package model

import (
	"bytes"
	"errors"
	"io/ioutil"

	"github.com/gin-gonic/gin"
)

type Ingredient struct {
	ID           int     `json:"id"`
	IngredientID int     `json:"ingredient_id"`
	ShopID       int     `json:"shop_id"`
	Name         string  `json:"name"`
	Category     int     `json:"category"`
	Image        string  `json:"image"`
	Measure      string  `json:"measure"`
	Cost         float32 `json:"cost" gorm:"-"`
	Deleted      bool    `json:"deleted" gorm:"default:false"`
	IsVisible    bool    `json:"is_visible" gorm:"default:true"`
}

type ReqIngredient struct {
	ID           int     `json:"id"`
	IngredientID int     `json:"ingredient_id"`
	ShopID       []int   `json:"shop_id"`
	Name         string  `json:"name"`
	Category     int     `json:"category"`
	Image        string  `json:"image"`
	Measure      string  `json:"measure"`
	Cost         float32 `json:"cost"`
	IsVisible    bool    `json:"is_visible" gorm:"default:true"`
}

func (p *ReqIngredient) ParseRequest(c *gin.Context) error {

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

type CategoryIngredient struct {
	ID      int    `json:"id"`
	Name    string `json:"name"`
	Image   string `json:"image"`
	ShopID  int    `json:"shop_id"`
	Deleted bool   `json:"deleted" gorm:"default:false"`
}

func (p *CategoryIngredient) ParseRequest(c *gin.Context) error {
	body, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		return err
	}
	c.Request.Body = &CustomReadCloser{reader: bytes.NewReader(body)}
	if err := c.BindJSON(&p); err != nil {
		return errors.New("bad request | " + err.Error())
	}
	c.Request.Body = &CustomReadCloser{reader: bytes.NewReader(body)}

	if p.Name == "" || p.ShopID <= 0 {
		return errors.New("bad request | fill fields properly")
	}

	return nil
}

type IngredientOutput struct {
	ID           int     `json:"id"`
	IngredientID int     `json:"ingredient_id"`
	ShopID       int     `json:"shop_id"`
	ShopName     string  `json:"shop_name"`
	Name         string  `json:"name"`
	Category     string  `json:"category"`
	CategoryID   int     `json:"category_id"`
	NaborID      int     `json:"nabor_id"`
	Image        string  `json:"image"`
	Measure      string  `json:"measure"`
	Cost         float32 `json:"cost"`
	Deleted      bool    `json:"deleted"`
	Brutto       float32 `json:"brutto"`
	Price        float32 `json:"price"`
}

func (p *IngredientOutput) ParseRequest(c *gin.Context) error {
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

type Nabor struct {
	ID          int                `json:"id"`
	NaborID     int                `json:"nabor_id"`
	ShopID      int                `json:"shop_id"`
	Shops       []int              `json:"shops" gorm:"-"`
	Name        string             `json:"name"`
	Min         int                `json:"min"`
	Max         int                `json:"max"`
	Ingredients []*IngredientNabor `json:"ingredient_nabor" gorm:"-"`
	Replaces    IntArray           `json:"replaces" gorm:"type:integer[]"`
	Deleted     bool               `json:"deleted" gorm:"default:false"`
}

func (p *Nabor) ParseRequest(c *gin.Context) error {

	body, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		return err
	}
	c.Request.Body = &CustomReadCloser{reader: bytes.NewReader(body)}
	if err := c.BindJSON(&p); err != nil {
		return errors.New("bad request | " + err.Error())
	}
	c.Request.Body = &CustomReadCloser{reader: bytes.NewReader(body)}

	if p.Name == "" || p.Min < 0 || p.Max < 0 {
		return errors.New("bad request | fill fields properly")
	}

	return nil
}

type IngredientNabor struct {
	ID           int     `json:"id" gorm:"primary_key"`
	NaborID      int     `json:"nabor_id"`
	IngredientID int     `json:"ingredient_id" gorm:"ForeignKey:IngredientID"`
	ShopID       int     `json:"shop_id"`
	Name         string  `json:"name" gorm:"-"`
	Measure      string  `json:"measure" gorm:"-"`
	Image        string  `json:"image" gorm:"-"`
	Brutto       float32 `json:"brutto"`
	Price        float32 `json:"price"`
}
