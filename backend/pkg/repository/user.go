package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"time"
	"zebra/model"
	"zebra/utils"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type UserDB struct {
	db     *sql.DB
	gormDB *gorm.DB
}

func NewUserDB(db *sql.DB, gormDB *gorm.DB) *UserDB {
	return &UserDB{db: db, gormDB: gormDB}
}

func (r *UserDB) AddUser(user *model.User) (int, error) {
	if err := r.gormDB.Create(&user).Error; err != nil {
		return -1, err
	}
	return user.ID, nil
}
func (r *UserDB) GetUser(id int) (*model.User, error) {
	user := &model.User{}
	if err := r.gormDB.First(&user, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return user, nil
}

func (r *UserDB) GetCurrentOrders(id int) ([]*model.Check, error) {
	checksGorm := []*model.Check{}

	err := r.gormDB.Transaction(func(tx *gorm.DB) error {
		res := tx.Model(checksGorm).Select("*").Where("user_id = ?", id).Scan(&checksGorm)
		if res.Error != nil {
			return res.Error
		}

		for i := 0; i < len(checksGorm); i++ {
			feedback := &model.Feedback{}
			tovarsGorm := []*model.CheckTovar{}
			techCartsGorm := []*model.CheckTechCart{}
			res := tx.Table("check_tovars").Select("*").Where("check_id = ?", checksGorm[i].ID).Scan(&tovarsGorm)
			if res.Error != nil {
				return res.Error
			}
			res = tx.Table("check_tech_carts").Select("*").Where("check_id = ?", checksGorm[i].ID).Scan(&techCartsGorm)
			if res.Error != nil {
				return res.Error
			}
			res = tx.Model(feedback).Select("*").Where("check_id = ?", checksGorm[i].ID).Scan(&feedback)
			if res.Error != nil {
				return res.Error
			}
			checksGorm[i].Tovar = tovarsGorm
			checksGorm[i].TechCart = techCartsGorm
			checksGorm[i].Feedback = feedback
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return checksGorm, err
}

func (r *UserDB) AddFeedback(feedback *model.Feedback) error {
	if err := r.gormDB.Create(&feedback).Error; err != nil {
		return err
	}
	return nil
}

func (r *UserDB) CheckCode(code string) (*model.UserQR, error) {
	userQR := &model.UserQR{}
	if err := r.gormDB.First(userQR, "code = ?", code).Scan(&userQR).Error; err != nil {
		return nil, err
	}
	now := time.Now().Local().Unix()
	if userQR.ExpireTime <= now {
		if err := r.gormDB.Delete(userQR).Error; err != nil {
			return nil, err
		}
		return nil, nil
	} else {
		return nil, errors.New("code exists")
	}
	return userQR, nil
}

func (r *UserDB) SetCode(id int, code string) (*model.UserQR, error) {
	userQR := &model.UserQR{
		UserID:     id,
		Code:       code,
		ExpireTime: time.Now().Local().Add(time.Minute * utils.QRExpire).Unix(),
	}

	err := r.gormDB.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "user_id"}},                        // key colume
		DoUpdates: clause.AssignmentColumns([]string{"code", "expire_time"}), // column needed to be updated
	}).Create(&userQR).Error
	if err != nil {
		return nil, err
	}
	return userQR, nil
}
func (r *UserDB) GetUserCode(id int) (*model.UserQR, error) {
	userQR := &model.UserQR{}
	if err := r.gormDB.First(userQR, "user_id = ?", id).Scan(&userQR).Error; err != nil {
		return nil, err
	}
	return userQR, nil
}

func (r *UserDB) GetUserByCode(code string) (*model.User, error) {
	userQR := &model.UserQR{}
	if err := r.gormDB.First(userQR, "code = ?", code).Scan(&userQR).Error; err != nil {
		return nil, err
	}
	if userQR.ExpireTime <= time.Now().Local().Unix() {
		return nil, errors.New("code expired")
	}
	user := &model.User{}
	if err := r.gormDB.First(user, "id = ?", userQR.UserID).Scan(&user).Error; err != nil {
		return nil, err
	}
	return user, nil
}

func (r *UserDB) CleanExpiredCode() error {
	now := time.Now().Unix()
	if err := r.gormDB.Model(&model.UserQR{}).Where("expire_time <= ?", now).Delete(&model.UserQR{}).Error; err != nil {
		return err
	}
	return nil
}

func (r *UserDB) GetWorkerByUsername(username string) (*model.Worker, error) {
	worker := &model.Worker{}
	if err := r.gormDB.First(worker, "username = ?", username).Scan(&worker).Error; err != nil {
		return nil, err
	}
	return worker, nil
}

func (r *UserDB) CreateWorker(worker *model.Worker) (*model.Worker, error) {
	existWorker := &model.Worker{}
	if err := r.gormDB.First(existWorker, "username = ?", worker.Username).Scan(&existWorker).Error; err != nil {
		if err.Error() != "record not found" {
			return nil, err
		}
	}
	if existWorker.Username != "" {
		return nil, fmt.Errorf("the '%s' username already exists", worker.Username)
	}
	if err := r.gormDB.Create(&worker).Error; err != nil {
		return nil, err
	}
	return worker, nil
}

func (r *UserDB) UpdateWorker(worker *model.Worker) error {
	if err := r.gormDB.Save(&worker).Error; err != nil {
		return err
	}
	return nil
}
