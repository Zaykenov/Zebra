package service

import (
	"zebra/model"
	"zebra/pkg/repository"
	"zebra/utils"
)

type FinanceService struct {
	repo repository.Finance
}

func NewFinanceService(repo repository.Finance) *FinanceService {
	return &FinanceService{repo: repo}
}

func (s *FinanceService) AddSchet(schet *model.ReqSchet) (*model.Schet, error) {
	return s.repo.AddSchet(schet)
}

func (s *FinanceService) GetAllSchet(filter *model.Filter) ([]*model.Schet, int64, error) {
	res, count, err := s.repo.GetAllSchet(filter)

	if err != nil {
		return nil, 0, err
	}

	pageCount := utils.CalculatePageCount(count, utils.DefaultPageSize)

	return res, pageCount, nil

}

func (s *FinanceService) GetSchet(id int) (*model.Schet, error) {
	return s.repo.GetSchetByID(id)
}

func (s *FinanceService) UpdateSchet(schet *model.Schet) error {
	return s.repo.UpdateSchet(schet)
}

func (s *FinanceService) DeleteSchet(id int) error {
	return s.repo.DeleteSchet(id)
}
