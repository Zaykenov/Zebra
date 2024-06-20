package service

import (
	"errors"
	"sort"
	"zebra/model"
	"zebra/pkg/repository"
	"zebra/utils"
)

type TovarService struct {
	repo repository.Tovar
}

func NewTovarService(repo repository.Tovar) *TovarService {
	return &TovarService{repo: repo}
}

func (s *TovarService) AddTovars(tovars []*model.Tovar) error {
	return s.repo.AddTovars(tovars)
}

func (s *TovarService) AddTovar(tovar *model.Tovar) error {
	return s.repo.AddTovar(tovar)
}

func (s *TovarService) GetAllTovar(filter *model.Filter) ([]*model.TovarOutput, int64, error) {
	res, count, err := s.repo.GetAllTovar(filter)
	if err != nil {
		return nil, 0, err
	}
	pageCount := utils.CalculatePageCount(count, utils.DefaultPageSize)
	return res, pageCount, nil
}

func (s *TovarService) GetTovar(id int, filter *model.Filter) (*model.TovarOutput, error) {

	return s.repo.GetTovar(id, filter)
}

func (s *TovarService) UpdateTovar(tovar *model.ReqTovar) error {
	return s.repo.UpdateTovar(tovar)
}

func (s *TovarService) DeleteTovar(id int) error {
	return s.repo.DeleteTovar(id)
}

func (s *TovarService) AddCategoryTovar(category *model.CategoryTovar) error {
	return s.repo.AddCategoryTovar(category)
}

func (s *TovarService) GetAllCategoryTovar(filter *model.Filter) ([]*model.CategoryTovar, int64, error) {
	res, count, err := s.repo.GetAllCategoryTovar(filter)

	if err != nil {
		return nil, 0, err
	}

	pageCount := utils.CalculatePageCount(count, utils.DefaultPageSize)

	return res, pageCount, nil

}

func (s *TovarService) GetCategoryTovar(id int) (*model.CategoryTovar, error) {
	return s.repo.GetCategoryTovar(id)
}

func (s *TovarService) UpdateCategoryTovar(category *model.CategoryTovar) error {
	return s.repo.UpdateCategoryTovar(category)
}

func (s *TovarService) DeleteCategoryTovar(id int) error {
	return s.repo.DeleteCategoryTovar(id)
}

func (s *TovarService) AddTechCart(techCart *model.TechCart) error {
	return s.repo.AddTechCart(techCart)
}

func (s *TovarService) AddTechCarts(techCarts []*model.TechCart) error {
	return s.repo.AddTechCarts(techCarts)
}

func (s *TovarService) GetAllTechCart(filter *model.Filter) ([]*model.TechCartResponse, int64, error) {
	res, count, err := s.repo.GetAllTechCart(filter)
	if err != nil {
		return nil, 0, err
	}
	pageCount := utils.CalculatePageCount(count, utils.DefaultPageSize)
	return res, pageCount, nil
}

func (s *TovarService) GetTechCart(id int, filter *model.Filter) (*model.TechCartOutput, error) {
	return s.repo.GetTechCart(id, filter)
}

func (s *TovarService) UpdateTechCart(techCart *model.ReqTechCart, role string) error {
	return s.repo.UpdateTechCart(techCart, role)
}

func (s *TovarService) DeleteTechCart(id int) error {
	return s.repo.DeleteTechCart(id)
}

func (s *TovarService) GetTovarWithParams(sortParam, sklad, search string, category int) ([]*model.TovarOutput, error) {
	search = "%" + search + "%"
	res, err := s.repo.GetTovarWithParams(sortParam, sklad, search, category)
	if err != nil {
		return nil, err
	}
	switch sortParam {
	case "id.asc":
		sort.Slice(res, func(i, j int) bool {
			return res[i].ID < res[j].ID
		})
	case "id.desc":
		sort.Slice(res, func(i, j int) bool {
			return res[i].ID > res[j].ID
		})
	case "name.asc":
		sort.Slice(res, func(i, j int) bool {
			return res[i].Name < res[j].Name
		})
	case "name.desc":
		sort.Slice(res, func(i, j int) bool {
			return res[i].Name > res[j].Name
		})
	case "price.asc":
		sort.Slice(res, func(i, j int) bool {
			return res[i].Price < res[j].Price
		})
	case "price.desc":
		sort.Slice(res, func(i, j int) bool {
			return res[i].Price > res[j].Price
		})
	case "cost.asc":
		sort.Slice(res, func(i, j int) bool {
			return res[i].Cost < res[j].Cost
		})
	case "cost.desc":
		sort.Slice(res, func(i, j int) bool {
			return res[i].Cost > res[j].Cost
		})
	case "profit.asc":
		sort.Slice(res, func(i, j int) bool {
			return res[i].Profit < res[j].Profit
		})
	case "profit.desc":
		sort.Slice(res, func(i, j int) bool {
			return res[i].Profit > res[j].Profit
		})
	case "category.asc":
		sort.Slice(res, func(i, j int) bool {
			return res[i].Category < res[j].Category
		})
	case "category.desc":
		sort.Slice(res, func(i, j int) bool {
			return res[i].Category > res[j].Category
		})
	case "margin.asc":
		sort.Slice(res, func(i, j int) bool {
			return res[i].Margin < res[j].Margin
		})
	case "margin.desc":
		sort.Slice(res, func(i, j int) bool {
			return res[i].Margin > res[j].Margin
		})
	case "measure.asc":
		sort.Slice(res, func(i, j int) bool {
			return res[i].Measure < res[j].Measure
		})
	case "measure.desc":
		sort.Slice(res, func(i, j int) bool {
			return res[i].Measure > res[j].Measure
		})
	default:
		return nil, errors.New("invalid sort param")
	}

	return res, nil
}

func (s *TovarService) GetTechCartWithParams(sortParam, sklad string, category int) ([]*model.TechCartOutput, error) {
	res, err := s.repo.GetTechCartWithParams(sortParam, sklad, category)
	if err != nil {
		return nil, err
	}

	switch sortParam {
	case "id.asc":
		sort.Slice(res, func(i, j int) bool {
			return res[i].ID < res[j].ID
		})
	case "id.desc":
		sort.Slice(res, func(i, j int) bool {
			return res[i].ID > res[j].ID
		})
	case "name.asc":
		sort.Slice(res, func(i, j int) bool {
			return res[i].Name < res[j].Name
		})
	case "name.desc":
		sort.Slice(res, func(i, j int) bool {
			return res[i].Name > res[j].Name
		})
	case "price.asc":
		sort.Slice(res, func(i, j int) bool {
			return res[i].Price < res[j].Price
		})
	case "price.desc":
		sort.Slice(res, func(i, j int) bool {
			return res[i].Price > res[j].Price
		})
	case "cost.asc":
		sort.Slice(res, func(i, j int) bool {
			return res[i].Cost < res[j].Cost
		})
	case "cost.desc":
		sort.Slice(res, func(i, j int) bool {
			return res[i].Cost > res[j].Cost
		})
	case "profit.asc":
		sort.Slice(res, func(i, j int) bool {
			return res[i].Profit < res[j].Profit
		})
	case "profit.desc":
		sort.Slice(res, func(i, j int) bool {
			return res[i].Profit > res[j].Profit
		})
	case "category.asc":
		sort.Slice(res, func(i, j int) bool {
			return res[i].Category < res[j].Category
		})
	case "category.desc":
		sort.Slice(res, func(i, j int) bool {
			return res[i].Category > res[j].Category
		})
	case "margin.asc":
		sort.Slice(res, func(i, j int) bool {
			return res[i].Margin < res[j].Margin
		})
	case "margin.desc":
		sort.Slice(res, func(i, j int) bool {
			return res[i].Margin > res[j].Margin
		})
	case "measure.asc":
		sort.Slice(res, func(i, j int) bool {
			return res[i].Measure < res[j].Measure
		})
	case "measure.desc":
		sort.Slice(res, func(i, j int) bool {
			return res[i].Measure > res[j].Measure
		})
	default:
		return nil, errors.New("invalid sort param")
	}
	return res, nil
}

func (s *TovarService) GetTechCartNabor(id int) ([]*model.NaborOutput, error) {
	return s.repo.GetTechCartNabor(id)
}

func (s *TovarService) GetEverything() ([]*model.ItemOutput, error) {
	return s.repo.GetEverything()
}

func (s *TovarService) GetTovarsTovarIDByTovarsID(id int) (int, error) {
	return s.repo.GetTovarsTovarIDByTovarsID(id)
}

func (s *TovarService) GetTechCartsTechCartIDByTechCartsID(id int) (int, error) {
	return s.repo.GetTechCartsTechCartIDByTechCartsID(id)
}

func (s *TovarService) GetDeletedTovar(filter *model.Filter) ([]*model.TovarOutput, int64, error) {
	res, count, err := s.repo.GetDeletedTovar(filter)
	if err != nil {
		return nil, 0, err
	}
	pageCount := utils.CalculatePageCount(count, utils.DefaultPageSize)
	return res, pageCount, nil
}

func (s *TovarService) GetDeletedTechCart(filter *model.Filter) ([]*model.TechCartResponse, int64, error) {
	res, count, err := s.repo.GetDeletedTechCart(filter)
	if err != nil {
		return nil, 0, err
	}
	pageCount := utils.CalculatePageCount(count, utils.DefaultPageSize)
	return res, pageCount, nil
}

func (s *TovarService) RecreateTovar(tovar *model.Tovar) error {
	return s.repo.RecreateTovar(tovar)
}

func (s *TovarService) RecreateTechCart(techCart *model.TechCart) error {
	return s.repo.RecreateTechCart(techCart)
}

func (s *TovarService) GetToAddTovar(filter *model.Filter) ([]*model.TovarMaster, int64, error) {
	res, count, err := s.repo.GetToAddTovar(filter)
	if err != nil {
		return nil, 0, err
	}
	pageCount := utils.CalculatePageCount(count, utils.DefaultPageSize)
	return res, pageCount, nil
}

func (s *TovarService) GetToAddTechCart(filter *model.Filter) ([]*model.TechCartMaster, int64, error) {
	res, count, err := s.repo.GetToAddTechCart(filter)
	if err != nil {
		return nil, 0, err
	}
	pageCount := utils.CalculatePageCount(count, utils.DefaultPageSize)
	return res, pageCount, nil
}

func (s *TovarService) GetIdsOfShopsWhereTheTovarAlreadyExist(tovarID int) ([]int, error) {
	return s.repo.GetIdsOfShopsWhereTheTovarAlreadyExist(tovarID)
}

func (s *TovarService) GetIdsOfShopsWhereTheTechCartAlreadyExist(techCartID int) ([]int, error) {
	return s.repo.GetIdsOfShopsWhereTheTechCartAlreadyExist(techCartID)
}
