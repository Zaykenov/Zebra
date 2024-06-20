package model

import (
	"bytes"
	"errors"
	"io/ioutil"

	"github.com/gin-gonic/gin"
)

type Dealer struct {
	ID       int         `json:"id"`
	Name     string      `json:"name"`
	Address  string      `json:"address"`
	Phone    string      `json:"phone"`
	Comment  string      `json:"comment"`
	Deleted  bool        `json:"deleted" gorm:"default:false"`
	Postavka []*Postavka `json:"postavka" gorm:"ForeignKey:DealerID"`
}

func (p *Dealer) ParseRequest(c *gin.Context) error {
	body, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		return err
	}
	c.Request.Body = &CustomReadCloser{reader: bytes.NewReader(body)}
	if err := c.BindJSON(&p); err != nil {
		return errors.New("bad request | " + err.Error())
	}
	c.Request.Body = &CustomReadCloser{reader: bytes.NewReader(body)}

	if p.Name == "" || p.Address == "" || p.Phone == "" {
		return errors.New("bad request | fill fields properly")
	}

	return nil
}
