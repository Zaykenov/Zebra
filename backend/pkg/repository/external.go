package repository

import (
	"database/sql"
	"zebra/model"

	"gorm.io/gorm"
)

type ExternalDB struct {
	db     *sql.DB
	gormDB *gorm.DB
}

func NewExternalDB(db *sql.DB, gormDB *gorm.DB) *ExternalDB {
	return &ExternalDB{db: db, gormDB: gormDB}
}

func (r *ExternalDB) SaveCheck(check *model.TisResponse) error {
	err := r.gormDB.Model(&model.TisResponse{}).Create(check).Error
	if err != nil {
		return err
	}
	return nil
}
