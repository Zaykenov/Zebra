package service

import (
	"sort"
	"zebra/model"
	"zebra/pkg/repository"
)

type TerminalService struct {
	repo repository.Repository
}

func NewTerminalService(repo repository.Repository) *TerminalService {
	return &TerminalService{repo: repo}
}

func (s *TerminalService) GetAllProducts(filter *model.Filter) (*model.Terminal, error) {
	products, err := s.repo.GetAllProducts(filter)
	if err != nil {
		return nil, err
	}
	categories, _, err := s.repo.GetAllCategoryTovar(filter)
	if err != nil {
		return nil, err
	}
	terminal := &model.Terminal{}
	j := 0
	categoryProducts := []*model.CategoryProduct{}
	for j := 0; j < len(products); j++ {
		if products[j].Category == 13 {
			terminal.MainDisplay = append(terminal.MainDisplay, products[j])
		}
	}
	j = 0
	for i := 0; i < len(categories); i++ {
		CategoryProduct := &model.CategoryProduct{
			ID:       categories[i].ID,
			Image:    categories[i].Image,
			Category: categories[i].Name,
		}
		if j >= len(products) {
			break
		}
		for products[j].Category == categories[i].ID {
			CategoryProduct.Products = append(CategoryProduct.Products, products[j])
			j++
			if j >= len(products) {
				break
			}
		}
		if len(CategoryProduct.Products) == 0 {
			CategoryProduct.Products = []*model.Product{}
		}
		categoryProducts = append(categoryProducts, CategoryProduct)
	}
	terminal.Categories = categoryProducts
	for _, categories := range terminal.Categories {
		sort.Slice(categories.Products, func(i, j int) bool {
			return categories.Products[i].Name < categories.Products[j].Name
		})
	}
	if len(terminal.MainDisplay) == 0 {
		terminal.MainDisplay = []*model.Product{}
	}
	sort.Slice(terminal.MainDisplay, func(i, j int) bool {
		return terminal.MainDisplay[i].Name < terminal.MainDisplay[j].Name
	})
	return terminal, nil
}

func (s *TerminalService) TerminalStart(filter *model.Filter) (*model.Terminal, error) {
	return nil, nil
}
