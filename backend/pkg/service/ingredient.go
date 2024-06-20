package service

import (
	"zebra/model"
	"zebra/pkg/repository"
	"zebra/utils"
)

type IngredientService struct {
	repo repository.Ingredient
}

func NewIngredientService(repo repository.Ingredient) *IngredientService {
	return &IngredientService{repo: repo}
}

func (s *IngredientService) AddIngredient(ingredient *model.Ingredient) error {
	return s.repo.AddIngredient(ingredient)
}

func (s *IngredientService) GetAllIngredient(filter *model.Filter) ([]*model.IngredientOutput, int64, error) {

	res, count, err := s.repo.GetAllIngredient(filter)

	if err != nil {
		return nil, 0, err
	}

	pageCount := utils.CalculatePageCount(count, utils.DefaultPageSize)

	return res, pageCount, nil
}

func (s *IngredientService) GetIngredient(id int, filter *model.Filter) (*model.IngredientOutput, error) {
	return s.repo.GetIngredient(id, filter)
}

func (s *IngredientService) UpdateIngredient(ingredient *model.ReqIngredient, shops []int) error {
	return s.repo.UpdateIngredient(ingredient, shops)
}

func (s *IngredientService) DeleteIngredient(id int) error {
	return s.repo.DeleteIngredient(id)
}

func (s *IngredientService) AddCategoryIngredient(category *model.CategoryIngredient) error {
	return s.repo.AddCategoryIngredient(category)
}

func (s *IngredientService) GetAllCategoryIngredient(filter *model.Filter) ([]*model.CategoryIngredient, int64, error) {

	res, count, err := s.repo.GetAllCategoryIngredient(filter)

	if err != nil {
		return nil, 0, err
	}

	pageCount := utils.CalculatePageCount(count, utils.DefaultPageSize)

	return res, pageCount, nil

}

func (s *IngredientService) GetCategoryIngredient(id int) (*model.CategoryIngredient, error) {
	return s.repo.GetCategoryIngredient(id)
}

func (s *IngredientService) UpdateCategoryIngredient(category *model.CategoryIngredient) error {
	return s.repo.UpdateCategoryIngredient(category)
}

func (s *IngredientService) DeleteCategoryIngredient(id int) error {
	return s.repo.DeleteCategoryIngredient(id)
}

func (s *IngredientService) AddNabor(nabor *model.Nabor) (*model.Nabor, error) {
	return s.repo.AddNabor(nabor)
}

func (s *IngredientService) GetAllNabor(filter *model.Filter) ([]*model.NaborOutput, int64, error) {
	res, count, err := s.repo.GetAllNabor(filter)
	if err != nil {
		return nil, 0, err
	}
	pageCount := utils.CalculatePageCount(count, utils.DefaultPageSize)
	return res, pageCount, nil
}

func (s *IngredientService) GetNabor(id int) (*model.NaborOutput, error) {
	return s.repo.GetNabor(id)
}

func (s *IngredientService) UpdateNabor(nabor *model.Nabor) error {
	return s.repo.UpdateNabor(nabor)
}

func (s *IngredientService) DeleteNabor(id int) error {
	return s.repo.DeleteNabor(id)
}

func (s *IngredientService) AddNabors(nabors []*model.Nabor) ([]*model.Nabor, error) {
	var err error
	for i := 0; i < len(nabors); i++ {
		if nabors[i].ID != 0 {
			err = s.repo.UpdateNabor(nabors[i])
		} else {
			nabors[i], err = s.repo.AddNabor(nabors[i])
		}
		if err != nil {
			return nil, err
		}
	}
	return nabors, nil
}

func (s *IngredientService) GetTechCartByIngredientID(id int) ([]*model.TechCart, error) {
	return s.repo.GetTechCartByIngredientID(id)
}

func (s *IngredientService) AddIngredients(ingredients []*model.Ingredient) error {
	return s.repo.AddIngredients(ingredients)
}

func (s *IngredientService) GetIngredientsIngredientIDByIngredientsID(id int) (int, error) {
	return s.repo.GetIngredientsIngredientIDByIngredientsID(id)
}

func (s *IngredientService) GetDeletedIngredient(filter *model.Filter) ([]*model.IngredientOutput, int64, error) {
	res, count, err := s.repo.GetDeletedIngredient(filter)
	if err != nil {
		return nil, 0, err
	}
	pageCount := utils.CalculatePageCount(count, utils.DefaultPageSize)
	return res, pageCount, nil
}

func (s *IngredientService) RecreateIngredient(ingredient *model.Ingredient) error {
	return s.repo.RecreateIngredient(ingredient)
}

func (s *IngredientService) GetToAddIngredient(filter *model.Filter) ([]*model.IngredientMaster, int64, error) {
	res, count, err := s.repo.GetToAddIngredient(filter)
	if err != nil {
		return nil, 0, err
	}
	pageCount := utils.CalculatePageCount(count, utils.DefaultPageSize)
	return res, pageCount, nil
}

func (s *IngredientService) GetIdsOfShopsWhereTheIngredientAlreadyExist(id int) ([]int, error) {
	return s.repo.GetIdsOfShopsWhereTheIngredientAlreadyExist(id)
}

func (s *IngredientService) GetIdsOfShopsWhereTheNaborAlreadyExist(id int) ([]int, error) {
	return s.repo.GetIdsOfShopsWhereTheNaborAlreadyExist(id)
}
