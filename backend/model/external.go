package model

import (
	"bytes"
	"errors"
	"io/ioutil"
	"time"

	"github.com/gin-gonic/gin"
)

type TisData struct {
	Data *ReqTisResponse `json:"data"`
}

type ReqTisResponse struct {
	ID            int       `json:"id"`
	CreatedAt     time.Time `json:"created_at"`
	TicketType    string    `json:"ticket_type"`
	ReceiptNumber string    `json:"receipt_number"`
	TotalSum      float32   `json:"total_sum"`
	Change        float32   `json:"change"`
	Link          string    `json:"link"`
	Items         []struct {
		Name     string  `json:"name"`
		Price    float32 `json:"price"`
		Quantity int     `json:"quantity"`
		Discount float32 `json:"discount"`
		Sum      float32 `json:"sum"`
	} `json:"items"`
	Payments []struct {
		PaymentMethod int     `json:"payment_method"`
		Sum           float32 `json:"sum"`
	} `json:"payments"`
	CheckID int `json:"check_id"`
}

func (p *ReqTisResponse) ParseRequest(c *gin.Context) error {
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

type TisResponse struct {
	ID      int    `json:"id"`
	CheckID int    `json:"check_id"`
	Data    []byte `json:"data"`
}

type TisBody struct {
	Token    string         `json:"token"`
	Type     int            `json:"type"`
	Items    []*TisItems    `json:"items"`
	Payments []*TisPayments `json:"payments"`
}

type TisItems struct {
	Name         string           `json:"name"`
	Price        float32          `json:"price"`
	Quantity     float32          `json:"quantity"`
	Discount     float32          `json:"discount"`
	KgdCode      int              `json:"kgd_code"`
	CompareField *TisCompareField `json:"compare_field"`
}

type TisCompareField struct {
	Type  string `json:"type"`
	Value string `json:"value"`
}

type TisPayments struct {
	PaymentMethod int     `json:"payment_method"`
	Sum           float32 `json:"sum"`
}

type WKBody struct {
	Token               string       `json:"Token"`
	CashboxUniqueNumber string       `json:"CashboxUniqueNumber"`
	OperationType       int          `json:"OperationType"`
	Positions           *WKPositions `json:"Positions"`
	TicketModifiers     *WKModifiers `json:"TicketModifiers"`
	Payments            *WKPayments  `json:"Payments"`
	Change              int          `json:"Change"`
	RoundType           int          `json:"RoundType"`
	ExternalCheckNumber int          `json:"ExternalCheckNumber"`
}

type WKPayments struct {
	Sum         int `json:"Sum"`
	PaymentType int `json:"PaymentType"`
}

type WKPositions struct {
	Count        int     `json:"Count"`
	Price        int     `json:"Price"`
	TaxPercent   int     `json:"TaxPercent"`
	Tax          float64 `json:"Tax"`
	TaxType      int     `json:"TaxType"`
	PositionName string  `json:"PositionName"`
	UnitCode     int     `json:"UnitCode"`
}

type WKModifiers struct {
	Sum        int    `json:"Sum"`
	Text       string `json:"Text"`
	Type       int    `json:"Type"`
	TaxPercent int    `json:"TaxPercent"`
	Tax        int    `json:"Tax"`
	TaxType    int    `json:"TaxType"`
}
