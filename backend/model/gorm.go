package model

import (
	"encoding/json"
	"fmt"
	"time"
)

type CheckResponse struct {
	ID               int                      `json:"id" gorm:"primaryKey"`
	MobileUserID     string                   `json:"mobile_user_id"`
	SkladID          int                      `json:"sklad_id"`
	ShopID           int                      `json:"shop_id"`
	WorkerID         int                      `json:"worker_id"`
	Worker           string                   `json:"worker"`
	Opened_at        time.Time                `json:"opened_at"`
	Closed_at        time.Time                `json:"closed_at"`
	IdempotencyKey   string                   `json:"idempotency_key"`
	Cash             float32                  `json:"cash"`
	Card             float32                  `json:"card"`
	Sum              float32                  `json:"sum"`
	Cost             float32                  `json:"cost"`
	Status           string                   `json:"status"`
	Payment          string                   `json:"payment"`
	Discount         float32                  `json:"discount"`
	DiscountPercent  float32                  `json:"discount_percent"`
	Link             string                   `json:"link"`
	Tovar            []*CheckTovar            `json:"tovarCheck" gorm:"ForeignKey:CheckID"`
	TechCart         []*CheckTechCartResponse `json:"techCartCheck" gorm:"ForeignKey:CheckID"`
	ModificatorCheck []*CheckModificator      `json:"modificatorCheck" gorm:"ForeignKey:CheckID"`
	Comment          string                   `json:"comment"`
	Feedback         *Feedback                `json:"feedback" gorm:"ForeignKey:CheckID"`
	ServicePercent   float32                  `json:"service_percent"`
	ServiceSum       float32                  `json:"service_sum"`
	Stolik           int                      `json:"stolik"`
}

type CheckTechCartResponse struct {
	ID           int     `json:"id" gorm:"primaryKey"`
	CheckID      int     `json:"check_id"`
	TechCartID   int     `json:"tech_cart_id"`
	TechCartName string  `json:"name"`
	Quantity     float32 `json:"quantity"`
	Cost         float32 `json:"cost"`
	Price        float32 `json:"price"`
	Discount     float32 `json:"discount"`
	Comments     string  `json:"comments"`
	Modificators string  `json:"modificators"`
	Ingredients  string  `json:"ingredients"`
}

type Check struct {
	ID               int                 `json:"id" gorm:"primaryKey"`
	MobileUserID     string              `json:"mobile_user_id"`
	SkladID          int                 `json:"sklad_id"`
	ShopID           int                 `json:"shop_id"`
	WorkerID         int                 `json:"worker_id"`
	Worker           string              `json:"worker"`
	Opened_at        time.Time           `json:"opened_at"`
	Closed_at        time.Time           `json:"closed_at"`
	IdempotencyKey   string              `json:"idempotency_key"`
	Version          int                 `json:"version"`
	Cash             float32             `json:"cash"`
	Card             float32             `json:"card"`
	Sum              float32             `json:"sum"`
	Cost             float32             `json:"cost"`
	Status           string              `json:"status"`
	Payment          string              `json:"payment"`
	Discount         float32             `json:"discount"`
	DiscountPercent  float32             `json:"discount_percent"`
	Link             string              `json:"link"`
	Tovar            []*CheckTovar       `json:"tovarCheck" gorm:"ForeignKey:CheckID"`
	TechCart         []*CheckTechCart    `json:"techCartCheck" gorm:"ForeignKey:CheckID"`
	ModificatorCheck []*CheckModificator `json:"modificatorCheck" gorm:"ForeignKey:CheckID"`
	Comment          string              `json:"comment"`
	Feedback         *Feedback           `json:"feedback" gorm:"ForeignKey:CheckID"`
	ServicePercent   float32             `json:"service_percent"`
	ServiceSum       float32             `json:"service_sum"`
	Stolik           int                 `json:"stolik"`
}

type CheckNeOutput struct {
	ID   int       `json:"id" gorm:"primaryKey"`
	Info CheckInfo `json:"checks_info"`
}

type CheckInfo struct {
	MobileUserID     string              `json:"mobile_user_id"`
	SkladID          int                 `json:"sklad_id"`
	ShopID           int                 `json:"shop_id"`
	WorkerID         int                 `json:"worker_id"`
	Worker           string              `json:"worker"`
	Opened_at        time.Time           `json:"opened_at"`
	Closed_at        time.Time           `json:"closed_at"`
	IdempotencyKey   string              `json:"idempotency_key"`
	Version          int                 `json:"version"`
	Cash             float32             `json:"cash"`
	Card             float32             `json:"card"`
	Sum              float32             `json:"sum"`
	Cost             float32             `json:"cost"`
	Status           string              `json:"status"`
	Payment          string              `json:"payment"`
	Discount         float32             `json:"discount"`
	DiscountPercent  float32             `json:"discount_percent"`
	Tovar            []*CheckTovar       `json:"tovarCheck" gorm:"ForeignKey:CheckID"`
	TechCart         []*CheckTechCart    `json:"techCartCheck" gorm:"ForeignKey:CheckID"`
	ModificatorCheck []*CheckModificator `json:"modificatorCheck" gorm:"ForeignKey:CheckID"`
	Comment          string              `json:"comment"`
}

type GlobalCheck struct {
	TotalMoney    float32  `json:"total_money"`
	TotalCash     float32  `json:"total_cash"`
	TotalCard     float32  `json:"total_card"`
	TotalDiscount float32  `json:"total_discount"`
	TotalNetCost  float32  `json:"total_net_cost"`
	TotalProfit   float32  `json:"total_profit"`
	Check         []*Check `json:"check"`
}

func (cis *CheckInfo) Scan(value interface{}) error {
	if value == nil {
		return nil
	}

	// Convert the value to a []byte slice
	bytes, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("failed to unmarshal CheckInfoSlice value: %v", value)
	}

	// Unmarshal the JSON-encoded value into the CheckInfoSlice struct
	if err := json.Unmarshal(bytes, cis); err != nil {
		return fmt.Errorf("failed to unmarshal CheckInfoSlice value: %v", err)
	}

	return nil
}

type CheckModificator struct {
	ID      int     `json:"id"`
	ItemID  int     `json:"item_id"`
	CheckID int     `json:"check_id"`
	Name    string  `json:"name"`
	Brutto  float32 `json:"brutto"`
	Cost    float32 `json:"cost"`
}

type CheckTovar struct {
	ID            int           `json:"id" gorm:"primaryKey"`
	CheckID       int           `json:"check_id"`
	TovarID       int           `json:"tovar_id"`
	TovarName     string        `json:"tovar_name"`
	Quantity      float32       `json:"quantity"`
	Cost          float32       `json:"cost"`
	Price         float32       `json:"price"`
	Discount      float32       `json:"discount"`
	Modifications string        `json:"modifications"`
	Comments      string        `json:"comments"`
	ExpenceTovar  *ExpenceTovar `json:"expence_tovar" gorm:"ForeignKey:CheckTovarID"`
}

type CheckTechCart struct {
	ID                int                  `json:"id" gorm:"primaryKey"`
	CheckID           int                  `json:"check_id"`
	TechCartID        int                  `json:"tech_cart_id"`
	TechCartName      string               `json:"name"`
	Quantity          float32              `json:"quantity"`
	Cost              float32              `json:"cost"`
	Price             float32              `json:"price"`
	Discount          float32              `json:"discount"`
	Comments          string               `json:"comments"`
	ExpenceIngredient []*ExpenceIngredient `json:"expence_ingredient" gorm:"ForeignKey:CheckTechCartID"`
}

type CheckView struct {
	ID    int    `json:"id" gorm:"primaryKey"`
	Check string `json:"check"`
}
