package model

import (
	"bytes"
	"errors"
	"io/ioutil"
	"time"

	"github.com/gin-gonic/gin"
)

type Client struct {
	ID             int       `json:"id"`
	Email          string    `json:"email"`
	ClientName     string    `json:"client_name"`
	BirthDate      string    `json:"birth_date"`
	DeviceID       string    `json:"device_id"`
	RegisteredDate time.Time `json:"registered_date"`
	LastCode       string    `json:"last_code"`
}

type ReqRegistrate struct {
	ID         int    `json:"id"`
	Email      string `json:"email"`
	ClientName string `json:"client_name"`
	BirthDate  string `json:"birth_date"`
	DeviceID   string `json:"device_id"`
}

func (w *ReqRegistrate) ParseRequest(c *gin.Context) error {
	body, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		return err
	}
	c.Request.Body = &CustomReadCloser{reader: bytes.NewReader(body)}
	if err := c.BindJSON(&w); err != nil {
		return errors.New("bad request | " + err.Error())
	}
	c.Request.Body = &CustomReadCloser{reader: bytes.NewReader(body)}
	if w.Email == "" || w.ClientName == "" {
		return errors.New("bad request | fill fields properly")
	}
	return nil
}

type ReqVerifyEmail struct {
	Code     string `json:"code"`
	DeviceID string `json:"device_id"`
}

func (w *ReqVerifyEmail) ParseRequest(c *gin.Context) error {
	body, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		return err
	}
	c.Request.Body = &CustomReadCloser{reader: bytes.NewReader(body)}
	if err := c.BindJSON(&w); err != nil {
		return errors.New("bad request | " + err.Error())
	}
	c.Request.Body = &CustomReadCloser{reader: bytes.NewReader(body)}
	if w.Code == "" {
		return errors.New("bad request | fill fields properly")
	}
	return nil
}

type ReqSignIn struct {
	Email    string `json:"email"`
	DeviceID string `json:"device_id"`
}

func (w *ReqSignIn) ParseRequest(c *gin.Context) error {
	body, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		return err
	}
	c.Request.Body = &CustomReadCloser{reader: bytes.NewReader(body)}
	if err := c.BindJSON(&w); err != nil {
		return errors.New("bad request | " + err.Error())
	}
	c.Request.Body = &CustomReadCloser{reader: bytes.NewReader(body)}
	if w.Email == "" {
		return errors.New("bad request | fill fields properly")
	}
	return nil
}
