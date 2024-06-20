package service

import (
	"zebra/model"
	"zebra/pkg/repository"
	"zebra/utils"
)

type DealerService struct {
	repo repository.Dealer
}

func NewDealerService(repo repository.Dealer) *DealerService {
	return &DealerService{repo: repo}
}

func (s *DealerService) AddDealer(dealer *model.Dealer) error {
	return s.repo.AddDealer(dealer)
}

func (s *DealerService) GetAllDealer(filter *model.Filter) ([]*model.Dealer, int64, error) {
	res, count, err := s.repo.GetAllDealer(filter)
	if err != nil {
		return nil, 0, err
	}

	pageCount := utils.CalculatePageCount(count, utils.DefaultPageSize)

	return res, pageCount, nil

}

func (s *DealerService) GetDealer(id int) (*model.Dealer, error) {
	return s.repo.GetDealerByID(id)
}

func (s *DealerService) UpdateDealer(dealer *model.Dealer) error {
	return s.repo.UpdateDealer(dealer)
}

func (s *DealerService) DeleteDealer(id int) error {
	return s.repo.DeleteDealer(id)
}
