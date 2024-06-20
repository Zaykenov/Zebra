package service

import (
	"zebra/model"
	"zebra/pkg/repository"
	"zebra/utils"
)

type MasterService struct {
	repo repository.Repository
}

func NewMasterService(repo repository.Repository) *MasterService {
	return &MasterService{repo: repo}
}

func (s *MasterService) GetShopsMaster() ([]*model.Shop, error) {
	shops, err := s.repo.GetAllPureShops()
	if err != nil {
		return nil, err
	}
	return shops, nil
}

func (s *MasterService) AddTovarMaster(tovar *model.TovarMaster) (*model.TovarMaster, error) {
	newTovar, err := s.repo.AddTovarMaster(tovar)
	if err != nil {
		return nil, err
	}

	return newTovar, nil
}

func (s *MasterService) AddIngredientMaster(ingredient *model.IngredientMaster) (*model.IngredientMaster, error) {
	newIngredient, err := s.repo.AddIngredientMaster(ingredient)
	if err != nil {
		return nil, err
	}

	return newIngredient, nil
}

func (s *MasterService) NormaliseTovars() error {
	return s.repo.NormaliseTovars()
}

func (s *MasterService) NormaliseIngredients() error {
	return s.repo.NormaliseIngredients()
}

func (s *MasterService) NormaliseTechCarts() error {
	return s.repo.NormaliseTechCarts()
}

func (s *MasterService) GetAllTovarMaster(filter *model.Filter) ([]*model.TovarMasterResponse, int64, error) {
	res, count, err := s.repo.GetAllTovarMaster(filter)
	if err != nil {
		return nil, 0, err
	}
	pageCount := utils.CalculatePageCount(count, utils.DefaultPageSize)
	return res, pageCount, nil
}

func (s *MasterService) GetTovarMaster(id int) (*model.TovarMasterResponse, error) {
	return s.repo.GetTovarMaster(id)
}

func (s *MasterService) AddTechCartMaster(techCart *model.TechCartMaster) (*model.TechCartMaster, error) {
	return s.repo.AddTechCartMaster(techCart)
}

func (s *MasterService) AddOrUpdateNaborMaster(nabor *model.NaborMaster) (*model.NaborMaster, error) {
	var err error
	if nabor.ID != 0 {
		_, err = s.repo.UpdateNaborMaster(nabor)
	} else {
		nabor, err = s.repo.AddNaborMaster(nabor)
	}
	if err != nil {
		return nil, err
	}
	return nabor, nil
}

func (s *MasterService) NormaliseNabors() error {
	return s.repo.NormaliseNabors()
}

func (s *MasterService) GetAllIngredientMaster(filter *model.Filter) ([]*model.IngredientMasterResponse, int64, error) {
	res, count, err := s.repo.GetAllIngredientMaster(filter)
	if err != nil {
		return nil, 0, err
	}
	pageCount := utils.CalculatePageCount(count, utils.DefaultPageSize)
	return res, pageCount, nil
}

func (s *MasterService) GetAllTechCartMaster(filter *model.Filter) ([]*model.TechCartMasterResponse, int64, error) {
	res, count, err := s.repo.GetAllTechCartMaster(filter)
	if err != nil {
		return nil, 0, err
	}
	pageCount := utils.CalculatePageCount(count, utils.DefaultPageSize)
	return res, pageCount, nil
}

func (s *MasterService) GetIngredientMaster(id int) (*model.IngredientMasterResponse, error) {
	return s.repo.GetIngredientMaster(id)
}

func (s *MasterService) GetTechCartMaster(id int) (*model.TechCartMasterResponse, error) {
	return s.repo.GetTechCartMaster(id)
}

func (s *MasterService) GetAllNaborMaster(filter *model.Filter) ([]*model.NaborMasterOutput, int64, error) {
	res, count, err := s.repo.GetAllNaborMaster(filter)
	if err != nil {
		return nil, 0, err
	}
	pageCount := utils.CalculatePageCount(count, utils.DefaultPageSize)
	return res, pageCount, nil
}

func (s *MasterService) GetNaborMaster(id int) (*model.NaborMasterOutput, error) {
	return s.repo.GetNaborMaster(id)
}

func (s *MasterService) UpdateTovarMaster(tovar *model.ReqTovarMaster) error {
	tovarMaster := &model.TovarMaster{
		ID:       tovar.ID,
		Name:     tovar.Name,
		Price:    tovar.Price,
		Measure:  tovar.Measure,
		Image:    tovar.Image,
		Deleted:  false,
		Discount: tovar.Discount,
		Tax:      tovar.Tax,
		Category: tovar.Category,
		Status:   utils.MenuStatusApproved,
	}
	return s.repo.UpdateTovarMaster(tovarMaster)
}

func (s *MasterService) UpdateIngredientMaster(ingredient *model.ReqIngredientMaster) error {
	ingredientMaster := &model.IngredientMaster{
		ID:       ingredient.ID,
		Name:     ingredient.Name,
		Category: ingredient.Category,
		Measure:  ingredient.Measure,
		Image:    ingredient.Image,
		Deleted:  false,
		Status:   utils.MenuStatusApproved,
	}
	return s.repo.UpdateIngredientMaster(ingredientMaster)
}
func (s *MasterService) UpdateTechCartMaster(techCart *model.ReqTechCartMaster) error {
	techCartMaster := &model.TechCartMaster{
		ID:          techCart.ID,
		Category:    techCart.Category,
		Name:        techCart.Name,
		Image:       techCart.Image,
		Tax:         techCart.Tax,
		Measure:     techCart.Measure,
		Price:       techCart.Price,
		Discount:    techCart.Discount,
		Deleted:     false,
		Status:      utils.MenuStatusApproved,
		Ingredients: techCart.Ingredients,
		Nabor:       techCart.Nabor,
	}
	return s.repo.UpdateTechCartMaster(techCartMaster)
}

func (s *MasterService) UpdateTechCartsMaster(techCarts []*model.ReqTechCartMaster) error {
	for _, val := range techCarts {
		techCart := &model.TechCartMaster{
			ID:          val.ID,
			Name:        val.Name,
			Category:    val.Category,
			Image:       val.Image,
			Tax:         val.Tax,
			Measure:     val.Measure,
			Price:       val.Price,
			Status:      utils.MenuStatusApproved,
			Ingredients: val.Ingredients,
			Nabor:       val.Nabor,
		}
		err := s.repo.UpdateTechCartMaster(techCart)
		if err != nil {
			return err
		}
		err = s.repo.UpdateTechCartsMaster(techCart)
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *MasterService) DeleteTovarMaster(id int) error {
	return s.repo.DeleteTovarMaster(id)
}
func (s *MasterService) DeleteIngredientMaster(id int) error {
	return s.repo.DeleteIngredientMaster(id)
}
func (s *MasterService) DeleteTechCartMaster(id int) error {
	return s.repo.DeleteTechCartMaster(id)
}

func (s *MasterService) CreateNaborMaster(nabors []*model.NaborMaster) ([]*model.NaborMaster, error) {
	result := []*model.NaborMaster{}
	for _, nabor := range nabors {
		if nabor.ID != 0 {
			res, err := s.repo.UpdateNaborMaster(nabor)
			if err != nil {
				return nil, err
			}
			result = append(result, res)
			continue
		}
		res, err := s.repo.CreateNaborMaster(nabor)
		if err != nil {
			return nil, err
		}
		result = append(result, res)
	}
	return result, nil
}

func (s *MasterService) GetAllTovarMasterIds() ([]int, error) {
	return s.repo.GetAllTovarMasterIds()
}

func (s *MasterService) GetAllTechCartsMasterIds() ([]int, error) {
	return s.repo.GetAllTechCartsMasterIds()
}

func (s *MasterService) ConfirmTovarMaster(id int) (*model.TovarMaster, error) {
	return s.repo.ConfirmTovarMaster(id)
}

func (s *MasterService) RejectTovarMaster(id int) error {
	return s.repo.RejectTovarMaster(id)
}

func (s *MasterService) ConfirmTechCartMaster(id int) (*model.TechCartMaster, error) {
	return s.repo.ConfirmTechCartMaster(id)
}

func (s *MasterService) RejectTechCartMaster(id int) error {
	return s.repo.RejectTechCartMaster(id)
}

func (s *MasterService) ConfirmIngredientMaster(id int) (*model.IngredientMaster, error) {
	return s.repo.ConfirmIngredientMaster(id)
}

func (s *MasterService) RejectIngredientMaster(id int) error {
	return s.repo.RejectIngredientMaster(id)
}

func (s *MasterService) ConfirmNaborMaster(id int) (*model.NaborMaster, error) {
	return s.repo.ConfirmNaborMaster(id)
}

func (s *MasterService) RejectNaborMaster(id int) error {
	return s.repo.RejectNaborMaster(id)
}

func (s *MasterService) AddNaborMaster(nabor *model.NaborMaster) (*model.NaborMaster, error) {
	return s.repo.CreateNaborMaster(nabor)
}

func (s *MasterService) UpdateNaborMaster(nabor *model.NaborMaster) (*model.NaborMaster, error) {
	return s.repo.UpdateNaborMaster(nabor)
}

func (s *MasterService) DeleteNaborMaster(id int) error {
	return s.repo.DeleteNaborMaster(id)
}
