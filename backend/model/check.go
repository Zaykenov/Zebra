package model

import (
	"errors"
	"time"

	"github.com/gin-gonic/gin"
)

type ModificatorExpense struct {
	Sum float32 `json:"sum"`
}

type Modificator struct {
	ID       int     `json:"id"`
	NaborID  int     `json:"nabor_id"`
	Name     string  `json:"name"`
	Brutto   float32 `json:"brutto"`
	Cost     float32 `json:"cost"`
	Quantity float32 `json:"quantity"`
}

type ModificatorReader struct {
	Modificators string `json:"modificators"`
}

type Test struct {
	ID   int    `json:"id" gorm:"primary_key"`
	Text string `json:"text"`
}

type CheckPrint struct {
	ID              int                      `json:"id" gorm:"primaryKey"`
	UserID          int                      `json:"user_id"`
	WorkerID        int                      `json:"worker_id"`
	Opened_at       time.Time                `json:"opened_at"`
	Closed_at       time.Time                `json:"closed_at"`
	Sum             float32                  `json:"sum"`
	Cost            float32                  `json:"cost"`
	Status          string                   `json:"status"`
	Payment         string                   `json:"payment"`
	Discount        float32                  `json:"discount"`
	DiscountPercent float32                  `json:"discount_percent"`
	Tovar           []*CheckTovar            `json:"tovarCheck" gorm:"ForeignKey:CheckID"`
	TechCart        []*CheckTechCartResponse `json:"techCartCheck" gorm:"ForeignKey:CheckID"`
	Comment         string                   `json:"comment"`
	Feedback        *Feedback                `json:"feedback" gorm:"ForeignKey:CheckID"`
	TisCheckUrl     string                   `json:"tisCheckUrl"`
	WorkerName      string                   `json:"workerName"`
}

type ErrorCheck struct {
	ID      int       `json:"id" gorm:"primaryKey"`
	Request string    `json:"request"`
	Error   string    `json:"error"`
	Time    time.Time `json:"time"`
}

type SendToTis struct {
	ID             int       `json:"id" gorm:"primaryKey"`
	CheckID        int       `json:"check_id"`
	IdempotencyKey string    `json:"idempotency_key"`
	Created_at     time.Time `json:"created_at"`
	Request        string    `json:"request"`
	Response       string    `json:"response"`
	Exception      string    `json:"exception"`
	Status         string    `json:"status"`
	RetryCount     int       `json:"retry_count"`
	CassaType      string    `json:"cassa_type"`
}

type FailedCheck struct {
	ID         int       `json:"id" gorm:"primaryKey"`
	Request    string    `json:"request"`
	Response   string    `json:"response"`
	Created_at time.Time `json:"created_at"`
	ShopID     int       `json:"shop_id"`
}

func (p *FailedCheck) ParseRequest(c *gin.Context) error {
	if err := c.BindJSON(&p); err != nil {
		return err
	}

	if p.Request == "" || p.Response == "" {
		return errors.New("bad request | fill fields properly")
	}

	return nil
}
