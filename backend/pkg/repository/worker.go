package repository

import (
	"database/sql"
	"zebra/model"
	"zebra/utils"

	"gorm.io/gorm"
)

type WorkerDB struct {
	db     *sql.DB
	gormDB *gorm.DB
}

func NewWorkerDB(db *sql.DB, gormDB *gorm.DB) *WorkerDB {
	return &WorkerDB{db: db, gormDB: gormDB}
}

func (r *WorkerDB) GetWorker(id int) (*model.Worker, error) {
	worker := &model.Worker{}
	if err := r.gormDB.First(&worker, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return worker, nil
}

func (r *WorkerDB) GetAllWorkers(filter *model.Filter) ([]*model.Worker, int64, error) {
	workers := []*model.Worker{}

	res := r.gormDB.Model(workers).Where("bind_shop IN (?) and deleted != ?", filter.AccessibleShops, true).Find(&workers)
	if res.Error != nil {
		return nil, 0, res.Error
	}

	newRes, count, err := filter.FilterResults(res, model.Worker{}, utils.DefaultPageSize, "", "", "")
	if err != nil {
		return nil, 0, err
	}

	if newRes.Scan(&workers).Error != nil {
		return nil, 0, newRes.Error
	}

	return workers, count, nil
}

func (r *WorkerDB) UpdateWorker(worker *model.Worker) (*model.Worker, error) {
	if err := r.gormDB.Save(&worker).Error; err != nil {
		return nil, err
	}
	return worker, nil
}

func (r *WorkerDB) DeleteWorker(id int) error {
	if err := r.gormDB.Model(&model.Worker{}).Where("id = ?", id).Update("deleted", true).Error; err != nil {
		return err
	}
	return nil
}
