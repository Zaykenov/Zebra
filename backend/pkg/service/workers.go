package service

import (
	"zebra/model"
	"zebra/pkg/repository"
	"zebra/utils"
)

type WorkerService struct {
	repo repository.Worker
}

func NewWorkerService(repo repository.Worker) *WorkerService {
	return &WorkerService{repo: repo}
}

func (s *WorkerService) GetAllWorkers(filter *model.Filter) ([]*model.Worker, int64, error) {
	res, count, err := s.repo.GetAllWorkers(filter)

	if err != nil {
		return nil, 0, err
	}

	pageCount := utils.CalculatePageCount(count, utils.DefaultPageSize)

	return res, pageCount, nil

}

func (s *WorkerService) GetWorker(id int) (*model.Worker, error) {
	return s.repo.GetWorker(id)
}

func (s *WorkerService) UpdateWorker(worker *model.Worker) (*model.Worker, error) {
	res, err := s.repo.GetWorker(worker.ID)
	if err != nil {
		return nil, err
	}
	worker.Token = res.Token
	return s.repo.UpdateWorker(worker)
}

func (s *WorkerService) DeleteWorker(id int) error {
	return s.repo.DeleteWorker(id)
}
