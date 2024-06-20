package model

import (
	"bytes"
	"errors"
	"io/ioutil"
	"time"
	"zebra/utils"

	"github.com/gin-gonic/gin"
)

type Shift struct {
	ID                     int       `json:"id"`
	SchetID                int       `json:"schet_id"`
	ShopID                 int       `json:"shop_id"`
	Shop                   string    `json:"shop"`
	CreatedAt              time.Time `json:"created_at"`
	ClosedAt               time.Time `json:"closed_at"`
	BeginSum               float32   `json:"begin_sum"`
	EndSumFact             float32   `json:"end_sum_fact"`
	EndSumPlan             float32   `json:"end_sum_plan"`
	Expense                float32   `json:"expense"`
	Income                 float32   `json:"income"`
	Cash                   float32   `json:"cash"`
	Card                   float32   `json:"card"`
	Collection             float32   `json:"collection"`
	Difference             float32   `json:"difference"`
	IsClosed               bool      `json:"is_closed"`
	EqualCash              bool      `json:"is_equal" gorm:"default:true"`
	DifferenceWithPrevious float32   `json:"difference_with_previous" gorm:"-"`
}

type CurrentShift struct {
	ID        int       `json:"id"`
	ShopID    int       `json:"shop_id"`
	IsClosed  bool      `json:"is_closed"`
	CreatedAt time.Time `json:"created_at"`
	WorkerID  int       `json:"worker_id"`
	Worker    string    `json:"worker"`
	ShopName  string    `json:"shop_name"`
}

func (p *Shift) ParseRequest(c *gin.Context) error {
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

type ShiftTransaction struct {
	ID           int                    `json:"id"`
	SchetID      int                    `json:"schet_id"`
	ShopID       int                    `json:"shop_id"`
	CreatedAt    time.Time              `json:"created_at"`
	ClosedAt     time.Time              `json:"closed_at"`
	BeginSum     float32                `json:"begin_sum"`
	EndSumFact   float32                `json:"end_sum_fact"`
	EndSumPlan   float32                `json:"end_sum_plan"`
	Expense      float32                `json:"expense"`
	Income       float32                `json:"income"`
	Cash         float32                `json:"cash"`
	Card         float32                `json:"card"`
	Collection   float32                `json:"collection"`
	Difference   float32                `json:"difference"`
	Transactions []*TransactionResponse `json:"transactions" gorm:"-"`
	IsClosed     bool                   `json:"is_closed"`
}
type Transaction struct {
	ID                  int                 `json:"id"`
	ShiftID             int                 `json:"shift_id"`
	SchetID             int                 `json:"schet_id"`
	WorkerID            int                 `json:"worker_id"`
	UpdatedWorkerID     int                 `json:"updated_worker_id"`
	Category            string              `json:"category"`
	Status              string              `json:"status"`
	Time                time.Time           `json:"time"`
	UpdatedTime         time.Time           `json:"updated_time"`
	Sum                 float32             `json:"sum"`
	Comment             string              `json:"comment"`
	TransactionPostavka TransactionPostavka `json:"transaction_postavka" gorm:"foreignKey:TransactionID;references:ID;"`
	TransactionTransfer TransactionTransfer `json:"transaction_transfer" gorm:"foreignKey:TransactionID;references:ID;"`
	Deleted             bool                `json:"deleted"`
}

func (p *Transaction) ParseRequest(c *gin.Context) error {
	body, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		return err
	}
	c.Request.Body = &CustomReadCloser{reader: bytes.NewReader(body)}
	if err := c.BindJSON(&p); err != nil {
		return errors.New("bad request | " + err.Error())
	}
	c.Request.Body = &CustomReadCloser{reader: bytes.NewReader(body)}
	if p.Time.IsZero() {
		p.Time = time.Now()
	}

	if p.Category != utils.OpenShift && p.Category != utils.CloseShift && p.Category != utils.Collection && p.Category != utils.Postavka && p.Category != utils.Income {
		return errors.New("bad request | category is required")
	}

	return nil
}

type TransactionResponse struct {
	ID                  int                 `json:"id"`
	ShiftID             int                 `json:"shift_id"`
	SchetID             int                 `json:"schet_id"`
	WorkerID            int                 `json:"worker_id"`
	Worker              string              `json:"worker"`
	UpdatedWorkerID     int                 `json:"updated_worker_id"`
	UpdatedWorker       string              `json:"updated_worker"`
	Category            string              `json:"category"`
	Status              string              `json:"status"`
	Time                time.Time           `json:"time"`
	UpdatedTime         time.Time           `json:"updated_time"`
	Sum                 float32             `json:"sum"`
	Comment             string              `json:"comment"`
	TransactionPostavka TransactionPostavka `json:"transaction_postavka" gorm:"foreignKey:TransactionID;references:ID;"`
	TransactionTransfer TransactionTransfer `json:"transaction_transfer" gorm:"foreignKey:TransactionID;references:ID;"`
	Deleted             bool                `json:"deleted"`
}

type TransactionPostavka struct {
	TransactionID int `json:"transaction_id" gorm:"ForeignKey:TransactionID"`
	PostavkaID    int `json:"postavka_id" gorm:"ForeignKey:PostavkaID"`
}

type TransactionTransfer struct {
	TransactionID int `json:"transaction_id" gorm:"ForeignKey:TransactionID"`
	TransferID    int `json:"transfer_id" gorm:"ForeignKey:TransferID"`
}
