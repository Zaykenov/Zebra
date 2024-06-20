package model

import (
	"bytes"
	"errors"
	"io"
	"io/ioutil"
	"math"
	"strconv"
	"strings"
	"time"
	"zebra/utils"

	"github.com/gin-gonic/gin"
)

type CustomReadCloser struct {
	reader io.Reader
}

func (crc *CustomReadCloser) Read(p []byte) (n int, err error) {
	return crc.reader.Read(p)
}

func (crc *CustomReadCloser) Close() error {
	// You can handle any cleanup or close operations if needed
	return nil
}

// ReqID struct
type ReqID struct {
	ID int `json:"id"`
}

func (p *ReqID) ParseRequest(c *gin.Context) error {
	body, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		return err
	}
	c.Request.Body = &CustomReadCloser{reader: bytes.NewReader(body)}
	if err := c.BindJSON(&p); err != nil {
		return errors.New("bad request | " + err.Error())
	}
	c.Request.Body = &CustomReadCloser{reader: bytes.NewReader(body)}

	if p.ID == 0 {
		return errors.New("bad request | id is required")
	}

	return nil
}

type ReqInventarizationDetail struct {
	ID   int    `json:"id"`
	Type string `json:"type"`
}

func (p *ReqInventarizationDetail) ParseRequest(c *gin.Context) error {
	body, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		return err
	}
	c.Request.Body = &CustomReadCloser{reader: bytes.NewReader(body)}
	if err := c.BindJSON(&p); err != nil {
		return errors.New("bad request | " + err.Error())
	}
	c.Request.Body = &CustomReadCloser{reader: bytes.NewReader(body)}

	if p.ID == 0 {
		return errors.New("bad request | id is required")
	}

	if p.Type == "" {
		return errors.New("bad request | type is required")
	}

	return nil
}

type ReqInventarizationGroup struct {
	ID      int                            `json:"id"`
	Name    string                         `json:"name"`
	SkladID int                            `json:"sklad_id"`
	Type    string                         `json:"type"`
	Measure string                         `json:"measure"`
	Items   []*ReqInventarizationGroupItem `json:"items"`
}

type ReqInventarizationGroupItem struct {
	ItemID int `json:"item_id"`
}

func (p *ReqInventarizationGroup) ParseRequest(c *gin.Context) error {
	body, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		return err
	}
	c.Request.Body = &CustomReadCloser{reader: bytes.NewReader(body)}
	if err := c.BindJSON(&p); err != nil {
		return errors.New("bad request | " + err.Error())
	}
	c.Request.Body = &CustomReadCloser{reader: bytes.NewReader(body)}
	if p.Name == "" || p.Type == "" || p.Measure == "" {
		return errors.New("bad request | name is required")
	}

	if len(p.Items) == 0 {
		return errors.New("bad request | items is required")
	}

	return nil
}

type ReqCheck struct {
	ID              int                 `json:"id"`
	IdempotencyKey  string              `json:"idempotency_key"`
	Version         int                 `json:"version"`
	WorkerID        int                 `json:"worker_id"`
	ShopID          int                 `json:"shop_id"`
	SkladID         int                 `json:"sklad_id"`
	Opened_at       time.Time           `json:"opened_at"`
	Closed_at       time.Time           `json:"closed_at"`
	Cash            float32             `json:"cash"`
	Card            float32             `json:"card"`
	Sum             float32             `json:"sum"`
	Cost            float32             `json:"cost"`
	Status          string              `json:"status"`
	Payment         string              `json:"payment"`
	Discount        float32             `json:"discount"`
	DiscountSum     float32             `json:"discount_sum"`
	DiscountPercent float32             `json:"discount_percent"`
	Tovars          []*ReqTovarCheck    `json:"tovarCheck"`
	TechCarts       []*ReqTechCartCheck `json:"techCartCheck"`
	Comment         string              `json:"comment"`
	MobileUserID    string              `json:"mobile_user_id"`
	Stolik          int                 `json:"stolik"`
	Service         bool                `json:"service"`
	OFD             bool                `json:"ofd"`
}

func (p *ReqCheck) ParseRequest(c *gin.Context) error {
	body, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		return err
	}
	c.Request.Body = &CustomReadCloser{reader: bytes.NewReader(body)}
	if err := c.BindJSON(&p); err != nil {
		return errors.New("bad request | " + err.Error())
	}
	c.Request.Body = &CustomReadCloser{reader: bytes.NewReader(body)}
	key := c.GetHeader("Idempotency-Key")
	keyParts := strings.Split(key, "_")
	if len(keyParts) != 2 {
		return errors.New("wrong idempotency key")
	}
	intVersion, err := strconv.Atoi(keyParts[1])
	if err != nil {
		intVersion = 1
	}
	p.IdempotencyKey = keyParts[0]
	p.Version = intVersion
	if len(p.Tovars) <= 0 && len(p.TechCarts) <= 0 && p.Status == utils.StatusClosed {
		return errors.New("bad request | fill fields properly")
	}

	if p.Opened_at.IsZero() {
		p.Opened_at = time.Now().Local()

	}

	p.Card = float32(math.Floor(float64(p.Card)))
	p.Cash = float32(math.Floor(float64(p.Cash)))
	p.Sum = float32(math.Floor(float64(p.Sum)))

	return nil
}

type ReqTovarCheck struct {
	TovarID       int     `json:"tovar_id"`
	TovarName     string  `json:"name"`
	Quantity      float32 `json:"quantity"`
	Cost          float32 `json:"cost"`
	Price         float32 `json:"price"`
	Modifications string  `json:"modifications"`
	Comments      string  `json:"comments"`
}

type ReqTechCartCheck struct {
	TechCartID   int            `json:"tech_cart_id"`
	TechCartName string         `json:"name"`
	Quantity     float32        `json:"quantity"`
	Cost         float32        `json:"cost"`
	Price        float32        `json:"price"`
	Modificators []*Modificator `json:"modificators"`
	Comments     string         `json:"comments"`
}

type Tag struct {
	ID     int    `json:"id"`
	Text   string `json:"text"`
	ShopID int    `json:"shop_id"`
}

func (p *Tag) ParseRequest(c *gin.Context) error {
	body, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		return err
	}
	c.Request.Body = &CustomReadCloser{reader: bytes.NewReader(body)}
	if err := c.BindJSON(&p); err != nil {
		return errors.New("bad request | " + err.Error())
	}
	c.Request.Body = &CustomReadCloser{reader: bytes.NewReader(body)}

	if p.Text == "" {
		return errors.New("bad request | fill fields properly")
	}

	return nil
}

type IdempotencyCheck struct {
	IdempotencyKey string    `json:"idempotency_key"`
	Time           time.Time `json:"time"`
}

type IdempotencyCheckArray struct {
	Keys       []*IdempotencyCheck `json:"keys"`
	ShopID     int                 `json:"shop_id"`
	Date       time.Time           `json:"date"`
	TimeIsLast bool                `json:"time_is_last"`
}

func (i *IdempotencyCheckArray) ParseRequest(c *gin.Context) error {
	body, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		return err
	}
	c.Request.Body = &CustomReadCloser{reader: bytes.NewReader(body)}
	if err := c.BindJSON(&i); err != nil {
		return errors.New("bad request | " + err.Error())
	}
	c.Request.Body = &CustomReadCloser{reader: bytes.NewReader(body)}

	if i.Date.IsZero() {
		i.TimeIsLast = true
	} else {
		i.TimeIsLast = false
	}
	return nil
}
