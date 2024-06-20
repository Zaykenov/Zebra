package repository

import (
	"database/sql"
	"zebra/model"
	"zebra/utils"

	"gorm.io/gorm"
)

type DealerDB struct {
	db     *sql.DB
	gormDB *gorm.DB
}

func NewDealerDB(db *sql.DB, gormDB *gorm.DB) *DealerDB {
	return &DealerDB{db: db, gormDB: gormDB}
}

func (r *DealerDB) AddDealer(dealer *model.Dealer) error {
	if err := r.gormDB.Create(dealer).Error; err != nil {
		return err
	}
	return nil
}

func (r *DealerDB) GetAllDealer(filter *model.Filter) ([]*model.Dealer, int64, error) {
	dealers := []*model.Dealer{}
	res := r.gormDB.Limit(10).Table("dealers").Select("*")
	if res.Error != nil {
		return nil, 0, res.Error
	}

	newRes, count, err := filter.FilterResults(res, model.Dealer{}, utils.DefaultPageSize, "", "", "")

	if err != nil {
		return nil, 0, err
	}

	if newRes.Scan(&dealers).Error != nil {
		return nil, 0, newRes.Error
	}

	return dealers, count, nil
}

func (r *DealerDB) GetDealerByID(id int) (*model.Dealer, error) {
	dealer := &model.Dealer{}
	res := r.gormDB.Table("dealers").Select("*").Where("id = ?", id).Scan(&dealer)
	if res.Error != nil {
		return nil, res.Error
	}
	return dealer, nil
}

func (r *DealerDB) UpdateDealer(dealer *model.Dealer) error {
	if err := r.gormDB.Save(dealer).Error; err != nil {
		return err
	}
	return nil
}

func (r *DealerDB) DeleteDealer(id int) error {
	if err := r.gormDB.Table("dealers").Where("id = ?", id).Update("deleted", true).Error; err != nil {
		return err
	}
	return nil
}

func (r *DealerDB) GetAnyDealer() (*model.Dealer, error) {
	dealer := &model.Dealer{}
	res := r.gormDB.Table("dealers").Select("*").Scan(&dealer)
	if res.Error != nil {
		return nil, res.Error
	}
	return dealer, nil
}
