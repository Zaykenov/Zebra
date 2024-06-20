package service

import (
	"log"
	"zebra/model"
	"zebra/pkg/repository"
	"zebra/utils"
)

type MobileService struct {
	repo repository.Mobile
}

func NewMobileService(repo repository.Mobile) *MobileService {
	return &MobileService{repo: repo}
}

func (s *MobileService) Registrate(req *model.ReqRegistrate) (string, error) {
	registered, err := s.repo.CheckForRegister(req.Email)
	if err != nil {
		return utils.NotOk, err
	}
	if registered {
		return utils.AlreadyRegistered, nil
	}
	code, err := s.repo.SendCode(req.Email)
	if err != nil {
		return utils.NotOk, err
	}
	err = s.repo.Registrate(code, req)
	if err != nil {
		return utils.NotOk, err
	}
	return utils.Ok, nil
}

func (s *MobileService) VerifyEmail(req *model.ReqVerifyEmail) (int, string, error) {
	client, err := s.repo.GetClientByDeviceID(req.DeviceID)
	if err != nil {
		return 0, utils.NotOk, err
	}
	if client.LastCode != req.Code {
		return 0, utils.IncorrectCode, nil
	}
	return client.ID, utils.Ok, nil
}

func (s *MobileService) SendCode(email string) error {
	code, err := s.repo.SendCode(email)
	if err != nil {
		return err
	}
	err = s.repo.UpdateCode(email, code)
	if err != nil {
		return err
	}
	return nil
}

func (s *MobileService) GetAllMobileUsers(filter *model.Filter) ([]*model.MobileUser, int64, error) {
	log.Print("GetAllMobileUsers")
	res, count, err := s.repo.GetAllMobileUsers(filter)
	if err != nil {
		return nil, 0, err
	}
	pageCount := utils.CalculatePageCount(count, utils.DefaultPageSize)
	return res, pageCount, nil
}
func (s *MobileService) GetMobileUser(id string, shopIDs []int) (*model.MobileUser, error) {
	return s.repo.GetMobileUser(id, shopIDs)
}
func (s *MobileService) GetAllFeedbacks(filter *model.Filter) ([]*model.MobileUserFeedbackResponse, int64, error) {
	res, count, err := s.repo.GetAllFeedbacks(filter)
	if err != nil {
		return nil, 0, err
	}
	pageCount := utils.CalculatePageCount(count, utils.DefaultPageSize)
	return res, pageCount, nil
}
func (s *MobileService) GetFeedback(id int) (*model.MobileUserFeedbackResponse, error) {
	return s.repo.GetFeedback(id)
}
