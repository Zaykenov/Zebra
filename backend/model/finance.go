package model

import (
	"bytes"
	"errors"
	"io/ioutil"

	"github.com/gin-gonic/gin"
)

type Schet struct {
	ID           int         `json:"id"`
	Name         string      `json:"name"`
	Currency     string      `json:"currency"`
	Type         string      `json:"type"`
	StartBalance float32     `json:"start_balance"`
	Deleted      bool        `json:"deleted" gorm:"default:false"`
	Postavka     []*Postavka `json:"postavka" gorm:"ForeignKey:SchetID"`
}

type ReqSchet struct {
	ID           int         `json:"id"`
	Name         string      `json:"name"`
	Currency     string      `json:"currency"`
	Type         string      `json:"type"`
	StartBalance float32     `json:"start_balance"`
	Deleted      bool        `json:"deleted" gorm:"default:false"`
	ShopIDs      []int       `json:"shops"`
	Postavka     []*Postavka `json:"postavka" gorm:"ForeignKey:SchetID"`
}

func (p *ReqSchet) ParseRequest(c *gin.Context) error {
	body, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		return err
	}
	c.Request.Body = &CustomReadCloser{reader: bytes.NewReader(body)}
	if err := c.BindJSON(&p); err != nil {
		return errors.New("bad request | " + err.Error())
	}
	c.Request.Body = &CustomReadCloser{reader: bytes.NewReader(body)}
	if p.Name == "" || p.Currency == "" || p.Type == "" {
		return errors.New("bad request | fill fields properly")
	}

	return nil
}
