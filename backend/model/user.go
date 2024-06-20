package model

import (
	"bytes"
	"database/sql/driver"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
	"zebra/utils"

	"github.com/gin-gonic/gin"
)

type User struct {
	ID       int     `json:"id"`
	DeviceID string  `json:"device_id"`
	Phone    string  `json:"phone"`
	Name     string  `json:"name"`
	Discount float32 `json:"discount"`
	Balance  int     `json:"balance"`
	UserQR   *UserQR `json:"user_qr" gorm:"ForeignKey:UserID"`
}

type UserRequest struct {
	ID       int    `json:"id"`
	DeviceID string `json:"device_id"`
	Phone    string `json:"phone"`
	Name     string `json:"name"`
}

func (u *UserRequest) ParseRequest(c *gin.Context) error {
	if err := c.BindJSON(&u); err != nil {
		return err
	}
	return nil
}

type UserResponse struct {
	ID       int     `json:"id"`
	DeviceID string  `json:"device_id"`
	Phone    string  `json:"phone"`
	Name     string  `json:"name"`
	Discount float32 `json:"discount"`
	Balance  int     `json:"balance"`
}

type Feedback struct {
	ID           int    `json:"id"`
	CheckID      int    `json:"check_id"`
	ScoreQuality int    `json:"score_quality"`
	ScoreService int    `json:"score_service"`
	Feedback     string `json:"feedback"`
	MobileUserID string `json:"mobile_user_id"`
}

type FeedbackRequest struct {
	ID           int    `json:"id"`
	CheckID      int    `json:"check_id"`
	ScoreQuality int    `json:"score_quality"`
	ScoreService int    `json:"score_service"`
	Feedback     string `json:"feedback"`
}

func (f *FeedbackRequest) ParseRequest(c *gin.Context) error {
	if err := c.BindJSON(&f); err != nil {
		return err
	}
	return nil
}

type UserQR struct {
	UserID     int    `json:"user_id" gorm:"primary_key"`
	Code       string `json:"code"`
	ExpireTime int64  `json:"expire_time"`
}

type ReqWorkerRegistration struct {
	Name     string `json:"name"`
	Phone    string `json:"phone"`
	Username string `json:"username"`
	Password string `json:"password"`
	Role     string `json:"role"`
	Shops    []int  `json:"shops"`
	BindShop int    `json:"bind_shop"`
	New      bool   `json:"new"`
}

func (w *ReqWorkerRegistration) ParseRequest(c *gin.Context) error {
	if err := c.BindJSON(&w); err != nil {
		return err
	}
	if w.Username == "" || w.Password == "" {
		return errors.New("bad request | fill fields properly")
	}
	if w.Role != utils.ManagerRole && w.Role != utils.WorkerRole {
		return errors.New("bad request | fill fields properly")
	}
	if len(w.Shops) == 0 {
		return errors.New("bad request | fill fields properly")
	}
	w.BindShop = w.Shops[0]
	return nil
}

type Worker struct {
	ID       int      `json:"id"`
	Token    string   `json:"token"`
	Name     string   `json:"name"`
	Phone    string   `json:"phone"`
	Username string   `json:"username"`
	Password string   `json:"password"`
	Role     string   `json:"role"`
	Shops    IntArray `json:"shops" gorm:"type:integer[]"`
	BindShop int      `json:"bind_shop"`
	Deleted  bool     `json:"deleted"`
}

type IntArray []int

func (a IntArray) Value() (driver.Value, error) {
	if a == nil {
		return nil, nil
	}

	// Convert the IntArray to a string representation of a PostgreSQL array
	var buf bytes.Buffer
	buf.WriteString("{")
	for i, v := range a {
		if i > 0 {
			buf.WriteString(",")
		}
		fmt.Fprintf(&buf, "%d", v)
	}
	buf.WriteString("}")

	return buf.String(), nil
}

func (a *IntArray) Scan(src interface{}) error {
	if src == nil {
		*a = []int{}
		return nil
	}

	// Convert the PostgreSQL array to a slice of integers
	var s string
	switch v := src.(type) {
	case []byte:
		s = string(v)
	case string:
		s = v
	default:
		return fmt.Errorf("unsupported Scan, storing driver.Value type %T into type %T", src, a)
	}

	if len(s) < 2 || s[0] != '{' || s[len(s)-1] != '}' {
		return fmt.Errorf("invalid array format: %s", s)
	}

	s = s[1 : len(s)-1]
	if s == "" {
		*a = []int{}
		return nil
	}

	parts := strings.Split(s, ",")
	res := make([]int, len(parts))
	for i, part := range parts {
		n, err := strconv.Atoi(part)
		if err != nil {
			return err
		}
		res[i] = n
	}

	*a = res
	return nil
}

type WorkerStat struct {
	ID       int     `json:"id"`
	Name     string  `json:"name"`
	Revenue  float32 `json:"revenue"`
	Cost     float32 `json:"cost"`
	Profit   float32 `json:"profit"`
	CheckNum int     `json:"check_num"`
	AvgCheck float32 `json:"avg_check"`
}

func (w *Worker) ParseRequest(c *gin.Context) error {
	if err := c.BindJSON(&w); err != nil {
		return err
	}
	if w.Username == "" || w.Password == "" {
		return errors.New("bad request | fill fields properly")
	}
	if len(w.Shops) == 0 {
		return errors.New("bad request | fill fields properly")
	}
	w.BindShop = w.Shops[0]
	return nil
}

type ReqLogin struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (w *ReqLogin) ParseRequest(c *gin.Context) error {
	if err := c.BindJSON(&w); err != nil {
		return err
	}
	if w.Username == "" || w.Password == "" {
		return errors.New("bad request | fill fields properly")
	}
	return nil
}

type MobileUser struct {
	ID               string                `json:"id" gorm:"primary_key"`
	Email            string                `json:"email"`
	Name             string                `json:"name" gorm:"not null"`
	BirthDate        time.Time             `json:"birth_date"`
	RegDate          time.Time             `json:"reg_date" gorm:"not null"`
	Discount         float32               `json:"discount" gorm:"not null"`
	ZebraCoinBalance float32               `json:"zebra_coin_balance" gorm:"not null"`
	RemoveDate       time.Time             `json:"remove_date"`
	Status           int                   `json:"status" gorm:"not null"`
	Feedbacks        []*MobileUserFeedback `json:"feedbacks" gorm:"foreignkey:UserID"`
	Checks           []*Check              `json:"checks" gorm:"foreignkey:MobileUserID"`
}

type MobileUserFeedback struct {
	ID           int     `json:"id" gorm:"primary_key"`
	UserID       string  `json:"user_id"`
	CheckID      int     `json:"check_id"`
	ShopID       int     `json:"shop_id"`
	WorkerID     int     `json:"worker_id"`
	ScoreQuality float32 `json:"score_quality"`
	ScoreService float32 `json:"score_service"`
	FeedbackText string  `json:"feedback_text"`
	CheckJson    string  `json:"check_json"`
}

func (m *MobileUserFeedback) TableName() string {
	return "user_feedback"
}

type MobileUserFeedbackResponse struct {
	ID           int       `json:"id" gorm:"primary_key"`
	UserID       string    `json:"user_id"`
	Username     string    `json:"username"`
	CheckID      int       `json:"check_id"`
	ShopID       int       `json:"shop_id"`
	ShopName     string    `json:"shop_name"`
	WorkerID     int       `json:"worker_id"`
	WorkerName   string    `json:"worker_name"`
	ScoreQuality float32   `json:"score_quality"`
	ScoreService float32   `json:"score_service"`
	FeedbackText string    `json:"feedback_text"`
	FeedbackDate time.Time `json:"feedback_date"`
	CheckJson    string    `json:"check_json"`
}
