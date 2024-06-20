package service

import (
	"zebra/model"
	"zebra/pkg/repository"
	"zebra/utils"
)

type ShopService struct {
	repo repository.Shop
}

func NewShopService(repo repository.Shop) *ShopService {
	return &ShopService{repo: repo}
}

func (s *ShopService) GetShopBySchetID(id int) (*model.Shop, error) {
	shop, err := s.repo.GetShopBySchetID(id)
	if err != nil {
		return nil, err
	}
	return shop, nil
}

func (s *ShopService) CreateShop(shop *model.ReqShop, products *model.ProductsShop) (*model.Shop, error) {
	stoliki := []*model.Stolik{}
	for i := 0; i < shop.Stoliki; i++ {
		stoliki = append(stoliki, &model.Stolik{
			StolikID: i + 1,
			ShopID:   shop.ID,
			Empty:    true,
		})
	}
	newShop := &model.Shop{
		Name:                shop.Name,
		Address:             shop.Address,
		TisToken:            shop.TisToken,
		CashSchet:           shop.CashSchet,
		CardSchet:           shop.CardSchet,
		CassaType:           "tis",
		CashboxUniqueNumber: "",
		Limit:               shop.Limit,
		ServicePercent:      shop.ServicePercent,
		Stoliki:             stoliki,
	}
	res, err := s.repo.CreateShop(newShop, products)
	if err != nil {
		return nil, err
	}
	return res, nil
}
func (s *ShopService) GetAllShops(filter *model.Filter) ([]*model.Shop, int64, error) {
	shops, count, err := s.repo.GetAllShop(filter)
	if err != nil {
		return nil, 0, err
	}
	totalPages := utils.CalculatePageCount(count, utils.DefaultPageSize)
	return shops, totalPages, nil
}
func (s *ShopService) GetShopByID(id int) (*model.Shop, error) {
	shop, err := s.repo.GetShop(id)
	if err != nil {
		return nil, err
	}
	return shop, nil
}
func (s *ShopService) UpdateShop(shop *model.ReqShop) error {
	stoliki := []*model.Stolik{}
	for i := 0; i < shop.Stoliki; i++ {
		stoliki = append(stoliki, &model.Stolik{
			StolikID: i + 1,
			ShopID:   shop.ID,
			Empty:    true,
		})
	}
	newShop := &model.Shop{
		ID:                  shop.ID,
		Name:                shop.Name,
		Address:             shop.Address,
		TisToken:            shop.TisToken,
		CashSchet:           shop.CashSchet,
		CardSchet:           shop.CardSchet,
		CassaType:           shop.CassaType,
		CashboxUniqueNumber: shop.CashboxUniqueNumber,
		Blocked:             shop.Blocked,
		Limit:               shop.Limit,
		ServicePercent:      shop.ServicePercent,
		Stoliki:             stoliki,
	}
	err := s.repo.UpdateShop(newShop)
	if err != nil {
		return err
	}
	return nil
}

func (s *ShopService) DeleteShop(id int) error {
	err := s.repo.DeleteShop(id)
	if err != nil {
		return err
	}
	return nil
}

func (r *ShopService) GetShopByCashSchetID(id int) (*model.Shop, error) {
	return r.repo.GetShopByCashSchetID(id)
}
