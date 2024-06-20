package repository

import (
	"database/sql"
	"fmt"
	"math/rand"
	"net/smtp"
	"strconv"
	"time"
	"zebra/model"
	"zebra/utils"

	"gorm.io/gorm"
)

type MobileDB struct {
	db     *sql.DB
	gormDB *gorm.DB
}

func NewMobileDB(db *sql.DB, gormDB *gorm.DB) *MobileDB {
	return &MobileDB{db: db, gormDB: gormDB}
}

func (r *MobileDB) CheckForRegister(email string) (bool, error) {
	var count int64
	res := r.gormDB.Model(model.Client{}).Where("email = ?", email).Count(&count)
	if res.Error != nil {
		return true, res.Error
	}
	if count > 0 {
		return true, nil
	}
	return false, nil
}

func (r *MobileDB) SendCode(email string) (string, error) {
	code := GenerateVerificationCode()
	from := "your-email@example.com"
	password := "your-email-password"
	smtpHost := "smtp.example.com"
	smtpPort := 587

	auth := smtp.PlainAuth("", from, password, smtpHost)

	message := fmt.Sprintf("Subject: User Registration Verification\n\nYour verification code is: %s", code)

	err := smtp.SendMail(smtpHost+":"+strconv.Itoa(smtpPort), auth, from, []string{email}, []byte(message))
	if err != nil {
		return utils.NotOk, err
	}

	return code, nil
}

func GenerateVerificationCode() string {
	rand.Seed(time.Now().UnixNano())
	code := rand.Intn(9999)
	return fmt.Sprintf("%04d", code)
}

func (r *MobileDB) Registrate(code string, req *model.ReqRegistrate) error {
	client := &model.Client{
		ID:             req.ID,
		Email:          req.Email,
		ClientName:     req.ClientName,
		BirthDate:      req.BirthDate,
		DeviceID:       req.DeviceID,
		RegisteredDate: time.Now(),
		LastCode:       code,
	}
	if err := r.gormDB.Create(&client).Error; err != nil {
		return err
	}
	return nil
}

func (r *MobileDB) GetClientByDeviceID(deviceID string) (*model.Client, error) {
	client := &model.Client{}
	res := r.gormDB.Model(&model.Client{}).Where("device_id = ?", deviceID).First(&client)
	if res.Error != nil {
		return nil, res.Error
	}
	return client, nil
}

func (r *MobileDB) UpdateCode(email, code string) error {
	res := r.gormDB.Model(&model.Client{}).Where("email = ?", email).Update("last_code", code)
	if res.Error != nil {
		return res.Error
	}
	return nil
}

func (r *MobileDB) GetAllMobileUsers(filter *model.Filter) ([]*model.MobileUser, int64, error) {
	users := []*model.MobileUser{}
	res := r.gormDB.Table("mobile_api.mobile_users").Order("id DESC")
	newRes, count, err := filter.FilterResults(res, model.MobileUser{}, utils.DefaultPageSize, "", "", "")
	if err != nil {
		return nil, 0, err
	}
	err = newRes.Scan(&users).Error
	if err != nil {
		return nil, 0, err
	}
	return users, count, nil
}
func (r *MobileDB) GetMobileUser(id string, shopIDs []int) (*model.MobileUser, error) {
	user := &model.MobileUser{}
	err := r.gormDB.Table("mobile_api.mobile_users").Where("id = ?", id).Scan(user).Error
	if err != nil {
		return nil, err
	}
	feedbacks := []*model.MobileUserFeedback{}
	checks := []*model.Check{}
	err = r.gormDB.Table("mobile_api.user_feedback").Where("user_id = ?", id).Order("\"Id\" desc").Limit(10).Scan(&feedbacks).Error
	if err != nil {
		return nil, err
	}
	err = r.gormDB.Table("public.checks").Where("mobile_user_id = ?", id).Order("id desc").Limit(10).Scan(&checks).Error
	if err != nil {
		return nil, err
	}
	user.Feedbacks = feedbacks
	user.Checks = checks
	return user, nil
}
func (r *MobileDB) GetAllFeedbacks(filter *model.Filter) ([]*model.MobileUserFeedbackResponse, int64, error) {
	feedbacks := []*model.MobileUserFeedbackResponse{}
	res := r.gormDB.Table("mobile_api.user_feedback").Select("mobile_api.user_feedback.\"Id\" as id,mobile_api.user_feedback.user_id, mobile_api.user_feedback.check_id, mobile_api.user_feedback.shop_id, mobile_api.user_feedback.worker_id, mobile_api.user_feedback.score_quality, mobile_api.user_feedback.score_service, mobile_api.user_feedback.feedback_text, mobile_api.user_feedback.check_json, mobile_api.user_feedback.feedback_date, public.workers.name as worker_name, public.shops.name as shop_name, mobile_api.mobile_users.name as username").Joins("inner join mobile_api.mobile_users on mobile_api.mobile_users.id = mobile_api.user_feedback.user_id inner join public.workers on public.workers.id = mobile_api.user_feedback.worker_id inner join public.shops on public.shops.id = mobile_api.user_feedback.shop_id")
	newRes, count, err := filter.FilterResults(res, model.MobileUserFeedback{}, utils.DefaultPageSize, "", "", "\"Id\" desc")
	if err != nil {
		return nil, 0, err
	}
	err = newRes.Scan(&feedbacks).Error
	if err != nil {
		return nil, 0, err
	}
	return feedbacks, count, nil
}
func (r *MobileDB) GetFeedback(id int) (*model.MobileUserFeedbackResponse, error) {
	feedback := &model.MobileUserFeedbackResponse{}
	err := r.gormDB.Table("mobile_api.user_feedback").Select("mobile_api.user_feedback.\"Id\" as id,mobile_api.user_feedback.user_id, mobile_api.user_feedback.check_id, mobile_api.user_feedback.shop_id, mobile_api.user_feedback.worker_id, mobile_api.user_feedback.score_quality, mobile_api.user_feedback.score_service, mobile_api.user_feedback.feedback_text, mobile_api.user_feedback.check_json, mobile_api.user_feedback.feedback_date, public.workers.name as worker_name, public.shops.name as shop_name, mobile_api.mobile_users.name as username").Joins("inner join mobile_api.mobile_users on mobile_api.mobile_users.id = mobile_api.user_feedback.user_id inner join public.workers on public.workers.id = mobile_api.user_feedback.worker_id inner join public.shops on public.shops.id = mobile_api.user_feedback.shop_id").Where("mobile_api.user_feedback.\"Id\" = ?", id).Scan(&feedback).Error
	if err != nil {
		return nil, err
	}
	return feedback, nil
}
