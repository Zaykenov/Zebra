package model

import (
	"bytes"
	"errors"
	"io/ioutil"

	"github.com/gin-gonic/gin"
)

type Shop struct {
	ID                  int       `json:"id"`
	Name                string    `json:"name"`
	Address             string    `json:"address"`
	TisToken            string    `json:"tis_token" gorm:"tis_token"`
	CashSchet           int       `json:"cash_schet" gorm:"cash_schet"`
	CardSchet           int       `json:"card_schet" gorm:"card_schet"`
	CassaType           string    `json:"cassa_type" gorm:"cassa_type"`
	CashboxUniqueNumber string    `json:"cashbox_unique_number" gorm:"cashbox_unique_number"`
	Shifts              []*Shift  `json:"shifts" gorm:"ForeignKey:ShopID"`
	Checks              []*Check  `json:"checks" gorm:"ForeignKey:ShopID"`
	Sklads              []*Sklad  `json:"sklads" gorm:"ForeignKey:ShopID"`
	Blocked             bool      `json:"block" gorm:"default:false"`
	Limit               float32   `json:"limit" gorm:"default:-1"`
	Schets              []Schet   `json:"schets" gorm:"many2many:shop_schets;"`
	ServicePercent      float32   `json:"service_percent" gorm:"default:0"`
	Stoliki             []*Stolik `json:"stoliki" gorm:"ForeignKey:ShopID"`
}

type Stolik struct {
	ID       int  `json:"id"`
	StolikID int  `json:"stolik_id"`
	ShopID   int  `json:"shop_id"`
	Empty    bool `json:"empty"`
}

type ShopSchet struct {
	ShopID  int `json:"shop_id"`
	SchetID int `json:"schet_id"`
}

type ShopGlobal struct {
	Shop         *ReqShop                 `json:"shop"`
	Workers      []*ReqWorkerRegistration `json:"workers"`
	Sklad        *Sklad                   `json:"sklad"`
	ProductsShop *ProductsShop            `json:"products_shop"`
}

type ProductsShop struct {
	Tovars    []int `json:"tovars"`
	TechCarts []int `json:"tech_carts"`
}

func (s *ShopGlobal) ParseRequest(c *gin.Context) error {
	body, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		return err
	}
	c.Request.Body = &CustomReadCloser{reader: bytes.NewReader(body)}
	if err := c.BindJSON(&s); err != nil {
		return errors.New("bad request | " + err.Error())
	}
	c.Request.Body = &CustomReadCloser{reader: bytes.NewReader(body)}
	return nil
}

type ReqShop struct {
	ID                  int     `json:"id"`
	Name                string  `json:"name"`
	Address             string  `json:"address"`
	TisToken            string  `json:"tis_token"`
	CashSchet           int     `json:"cash_schet"`
	CardSchet           int     `json:"card_schet"`
	CassaType           string  `json:"cassa_type"`
	CashboxUniqueNumber string  `json:"cashbox_unique_number"`
	Blocked             bool    `json:"blocked"`
	Limit               float32 `json:"limit"`
	ServicePercent      float32 `json:"service_percent"`
	Stoliki             int     `json:"stoliki"`
}

func (s *ReqShop) ParseRequest(c *gin.Context) error {
	body, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		return err
	}
	c.Request.Body = &CustomReadCloser{reader: bytes.NewReader(body)}
	if err := c.BindJSON(&s); err != nil {
		return errors.New("bad request | " + err.Error())
	}
	c.Request.Body = &CustomReadCloser{reader: bytes.NewReader(body)}
	return nil //??? TO ADD VALIDATION
}

type CategoryTovarTerminal struct {
	ID      int    `json:"id"`
	Name    string `json:"name"`
	Image   string `json:"image"`
	Deleted bool   `json:"deleted" gorm:"default:false"`
}

type ShopFromMaster struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Token string `json:"token"`
}
