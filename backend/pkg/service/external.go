package service

import (
	"encoding/json"
	"zebra/model"
	"zebra/pkg/repository"
)

type ExternalService struct {
	repo repository.External
}

func NewExternalService(repo repository.External) *ExternalService {
	return &ExternalService{repo: repo}
}

func (s *ExternalService) SaveCheck(check *model.ReqTisResponse) error {
	json, err := json.Marshal(check)
	if err != nil {
		return err
	}
	newCheck := &model.TisResponse{
		ID:      0,
		CheckID: check.CheckID,
		Data:    json,
	}
	return s.repo.SaveCheck(newCheck)

}
