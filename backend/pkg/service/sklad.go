package service

import (
	"errors"
	"sort"
	"time"
	"zebra/model"
	"zebra/pkg/repository"
	"zebra/utils"
)

type SkladService struct {
	repo repository.Repository
}

func NewSkladService(repo repository.Repository) *SkladService {
	return &SkladService{repo: repo}
}

func (s *SkladService) AddSklad(sklad *model.Sklad) error {
	return s.repo.AddSklad(sklad)
}

func (s *SkladService) GetAllSklad(filter *model.Filter) ([]*model.Sklad, int64, error) {
	res, count, err := s.repo.GetAllSklad(filter)

	if err != nil {
		return nil, 0, err
	}

	pageCount := utils.CalculatePageCount(count, utils.DefaultPageSize)

	return res, pageCount, nil

}

func (s *SkladService) GetSklad(id int) (*model.Sklad, error) {
	return s.repo.GetSklad(id)
}
func (s *SkladService) UpdateSklad(sklad *model.Sklad) error {
	return s.repo.UpdateSklad(sklad)
}
func (s *SkladService) DeleteSklad(id int) error {
	return s.repo.DeleteSklad(id)
}
func (s *SkladService) Ostatki(filter *model.Filter) ([]*model.Item, int64, error) {
	/*items, err := s.repo.Ostatki()
	if err != nil {
		return nil, err
	}
	sort.Slice(items, func(i, j int) bool {
		return items[i].Name < items[j].Name
	})
	return items, nil*/

	res, count, err := s.repo.Ostatki(filter)

	if err != nil {
		return nil, 0, err
	}

	pageCount := utils.CalculatePageCount(count, utils.DefaultPageSize)

	return res, pageCount, nil

}

func (s *SkladService) AddToSklad(postavka *model.Postavka, shopID int, id int) (*model.Postavka, error) {
	for _, item := range postavka.Items {
		postavka.Sum = postavka.Sum + item.Cost*item.Quantity
	}
	postavka, shiftID, err := s.repo.AddToSklad(postavka, shopID, id)
	if err != nil {
		return nil, err
	}
	err = s.repo.RecalculateShift(shiftID)
	if err != nil {
		return nil, err
	}
	return postavka, nil
}

func (s *SkladService) GetItems(filter *model.Filter) ([]*model.ItemOutput, error) {
	res, err := s.repo.GetItems(filter)
	if err != nil {
		return nil, err
	}
	sort.Slice(res, func(i, j int) bool {
		return res[i].Name < res[j].Name
	})
	return res, nil
}

func (s *SkladService) GetAllPostavka(filter *model.Filter) (*model.GlobalPostavka, int64, error) {
	res, count, err := s.repo.GetAllPostavka(filter)
	if err != nil {
		return nil, 0, err
	}
	pageCount := utils.CalculatePageCount(count, utils.PostavkaPageSize)
	postavkas := &model.GlobalPostavka{
		Postavka: res,
	}
	postavkas.Sum, err = s.repo.GetSumOfPostavkaForPeriod(filter)
	if err != nil {
		return nil, 0, err
	}
	return postavkas, pageCount, nil
}

func (s *SkladService) GetPostavka(id int) (*model.PostavkaOutput, error) {
	return s.repo.GetPostavka(id)
}

func (s *SkladService) UpdatePostavka(postavka *model.Postavka) error {
	for _, item := range postavka.Items {
		postavka.Sum = postavka.Sum + item.Cost*item.Quantity
	}
	return s.repo.UpdatePostavka(postavka)
}

func (s *SkladService) DeletePostavka(id int) error {
	return s.repo.DeletePostavka(id)
}

func (s *SkladService) RemoveFromSklad(spisanie *model.RemoveFromSklad) error {
	return s.repo.Sklad.RemoveFromSklad(spisanie)
}

func (s *SkladService) RequestToRemove(request *model.RemoveFromSklad) error {
	return s.repo.RequestToRemove(request)
}

func (s *SkladService) ConfirmToRemove(id int) error {
	return s.repo.ConfirmToRemove(id)
}

func (s *SkladService) UpdateSpisanie(spisanie *model.RemoveFromSklad) error {
	return s.repo.UpdateSpisanie(spisanie)
}

func (s *SkladService) DeleteSpisanie(id int) error {
	return s.repo.DeleteSpisanie(id)
}

func (s *SkladService) RejectToRemove(id int) error {
	return s.repo.RejectToRemove(id)
}

// func (s *SkladService) GetShopBySkladID(id int) (*model.Shop, error) {
// 	return s.repo.GetShopBySkladID(id)
// }

func (s *SkladService) GetRemoved(filter *model.Filter) (*model.GlobalSpisanie, int64, error) {
	res, count, err := s.repo.GetRemoved(filter)
	if err != nil {
		return nil, 0, err
	}
	pageCount := utils.CalculatePageCount(count, utils.DefaultPageSize)
	spisanies := &model.GlobalSpisanie{
		RemoveFromSklad: res,
	}
	for _, spisanie := range res {
		spisanies.Sum = spisanies.Sum + spisanie.Cost
	}
	return spisanies, pageCount, nil

}
func (s *SkladService) GetRemovedByID(id int) (*model.RemoveFromSkladResponse, error) {
	return s.repo.GetRemovedByID(id)
}
func (s *SkladService) AddTransfer(transfer *model.Transfer) error {
	transfer, err := s.repo.AddTransfer(transfer)
	if err != nil {
		return err
	}
	itemsPostavka := []model.ItemPostavka{}
	itemsSpisanie := []*model.RemoveFromSkladItem{}
	for _, item := range transfer.ItemTransfers {
		itemPostavka := model.ItemPostavka{
			ItemID:   item.ItemID,
			Type:     item.Type,
			Quantity: item.Quantity,
			Cost:     item.Sum / float32(item.Quantity),
		}
		itemsPostavka = append(itemsPostavka, itemPostavka)
		itemSpisanie := &model.RemoveFromSkladItem{
			ItemID:   item.ItemID,
			Type:     item.Type,
			Quantity: item.Quantity,
			SkladID:  transfer.FromSklad,
			Cost:     item.Sum / float32(item.Quantity),
		}
		itemsSpisanie = append(itemsSpisanie, itemSpisanie)
	}
	dealer, err := s.repo.GetAnyDealer()
	if err != nil {
		return err
	}
	postavka := &model.Postavka{
		DealerID:   dealer.ID,
		SkladID:    transfer.ToSklad,
		Time:       transfer.Time,
		Items:      itemsPostavka,
		Type:       utils.TypeTransfer,
		TransferID: transfer.ID,
		Sum:        transfer.Sum,
	}
	_, _, err = s.repo.AddToSklad(postavka, 0, transfer.Worker)
	if err != nil {
		return err
	}
	spisanie := &model.RemoveFromSklad{
		SkladID:    transfer.FromSklad,
		Time:       transfer.Time,
		Items:      itemsSpisanie,
		Type:       utils.TypeTransfer,
		TransferID: transfer.ID,
		WorkerID:   transfer.Worker,
		Cost:       0,
	}
	err = s.repo.Sklad.RemoveFromSklad(spisanie)
	if err != nil {
		return err
	}
	return nil
}

func (s *SkladService) GetTransfer(id int) (*model.TransferOutput, error) {
	return s.repo.GetTransfer(id)
}

func (s *SkladService) GetAllTransfer(filter *model.Filter) ([]*model.TransferOutput, int64, error) {
	res, count, err := s.repo.GetAllTransfer(filter)
	if err != nil {
		return nil, 0, err
	}
	pageCount := utils.CalculatePageCount(count, utils.DefaultPageSize)
	return res, pageCount, nil
}

func (s *SkladService) UpdateTransfer(transfer *model.Transfer) error {
	err := s.repo.UpdateTransfer(transfer)
	if err != nil {
		return err
	}
	postavka, err := s.repo.GetPostavkaByTransferID(transfer.ID)
	if err != nil {
		return err
	}
	postavka.SkladID = transfer.ToSklad
	postavka.Time = transfer.Time
	postavka.Sum = transfer.Sum
	itemsPostavka := []model.ItemPostavka{}
	for _, item := range transfer.ItemTransfers {
		itemPostavka := model.ItemPostavka{
			ItemID:   item.ItemID,
			Type:     item.Type,
			Quantity: item.Quantity,
			Cost:     item.Sum / float32(item.Quantity),
		}
		itemsPostavka = append(itemsPostavka, itemPostavka)
	}
	postavka.Items = itemsPostavka
	err = s.repo.UpdatePostavka(postavka)
	if err != nil {
		return err
	}
	spisanie, err := s.repo.GetSpisanieByTransferID(transfer.ID)
	if err != nil {
		return err
	}
	spisanie.SkladID = transfer.FromSklad
	spisanie.Time = transfer.Time
	spisanie.Cost = 0
	spisanie.WorkerID = transfer.Worker
	itemsSpisanie := []*model.RemoveFromSkladItem{}
	for _, item := range transfer.ItemTransfers {
		itemSpisanie := &model.RemoveFromSkladItem{
			ItemID:   item.ItemID,
			Type:     item.Type,
			Quantity: item.Quantity,
			SkladID:  transfer.FromSklad,
			Cost:     item.Sum / float32(item.Quantity),
		}
		itemsSpisanie = append(itemsSpisanie, itemSpisanie)
	}
	spisanie.Items = itemsSpisanie
	err = s.repo.UpdateSpisanie(spisanie)
	if err != nil {
		return err
	}
	return nil
}

func (s *SkladService) DeleteTransfer(id int) error {
	err := s.repo.DeleteTransfer(id)
	if err != nil {
		return err
	}
	postavka, err := s.repo.GetPostavkaByTransferID(id)
	if err != nil {
		return err
	}
	err = s.repo.DeletePostavka(postavka.ID)
	if err != nil {
		return err
	}
	spisanie, err := s.repo.GetSpisanieByTransferID(id)
	if err != nil {
		return err
	}
	err = s.repo.DeleteSpisanie(spisanie.ID)
	if err != nil {
		return err
	}
	return nil
}

func (s *SkladService) GetToCreateInventratization(inventarization *model.Inventarization) (*model.Inventarization, error) {
	ingIDMap := make(map[int]bool)
	tovarIDMap := make(map[int]bool)
	shop, err := s.repo.GetShopBySkladID(inventarization.SkladID)
	if err != nil {
		return nil, err
	}
	for _, item := range inventarization.InventarizationItems {
		item.SkladID = inventarization.SkladID
	}

	err = s.repo.CheckInventItems(inventarization.InventarizationItems, shop.ID)
	if err != nil {
		return nil, err
	}
	items := make([]*model.InventarizationItem, 0)
	for _, item := range inventarization.InventarizationItems {
		if item.Type == utils.TypeGroup {
			group, err := s.repo.GetPureInventarizationGroup(item.ItemID)
			if err != nil {
				return nil, err
			}
			for _, ing := range group.Items {
				if group.Type == utils.TypeIngredient {
					if _, ok := ingIDMap[ing.ItemID]; !ok {
						ingIDMap[ing.ItemID] = true
						newItem := &model.InventarizationItem{
							ItemID:    ing.ItemID,
							Time:      inventarization.Time,
							SkladID:   inventarization.SkladID,
							Type:      utils.TypeIngredient,
							IsVisible: false,
							GroupID:   item.ItemID,
						}
						if newItem.Time.After(time.Now()) {
							newItem.NeedToRecalculate = true
						}
						items = append(items, newItem)
					}
				} else {
					if _, ok := tovarIDMap[ing.ItemID]; !ok {
						tovarIDMap[ing.ItemID] = true
						newItem := &model.InventarizationItem{
							ItemID:    ing.ItemID,
							Time:      inventarization.Time,
							SkladID:   inventarization.SkladID,
							Type:      utils.TypeTovar,
							IsVisible: false,
							GroupID:   item.ItemID,
						}
						if newItem.Time.After(time.Now()) {
							newItem.NeedToRecalculate = true
						}
						items = append(items, newItem)
					}
				}
			}
			item.Time = inventarization.Time
			item.SkladID = inventarization.SkladID
			item.IsVisible = true
			items = append(items, item)
		}
	}
	for _, item := range inventarization.InventarizationItems {
		if item.Type == utils.TypeIngredient {
			if _, ok := ingIDMap[item.ItemID]; !ok {
				ingIDMap[item.ItemID] = true
				newItem := &model.InventarizationItem{
					ItemID:    item.ItemID,
					Time:      inventarization.Time,
					SkladID:   inventarization.SkladID,
					Type:      utils.TypeIngredient,
					IsVisible: true,
				}
				if newItem.Time.After(time.Now()) {
					newItem.NeedToRecalculate = true
				}
				items = append(items, newItem)
			}
		} else if item.Type == utils.TypeTovar {
			if _, ok := tovarIDMap[item.ItemID]; !ok {
				tovarIDMap[item.ItemID] = true
				newItem := &model.InventarizationItem{
					ItemID:    item.ItemID,
					Time:      inventarization.Time,
					SkladID:   inventarization.SkladID,
					Type:      utils.TypeTovar,
					IsVisible: true,
				}
				if newItem.Time.After(time.Now()) {
					newItem.NeedToRecalculate = true
				}
				items = append(items, newItem)
			}
		}
	}
	inventarization.InventarizationItems = items
	return s.repo.GetToCreateInventratization(inventarization)
}

func (s *SkladService) RecalculateInventarization() error {
	return s.repo.RecalculateInventarizations()
}

func (s *SkladService) RecalculateNetCost() error {
	go s.repo.RecalculateNetCost()
	return nil
}

func (s *SkladService) UpdateInventarization(inventarization *model.Inventarization) (*model.Inventarization, error) {
	ingMap := make(map[int]*model.InventarizationItem)
	tovarMap := make(map[int]*model.InventarizationItem)
	newGroup := false
	for _, item := range inventarization.InventarizationItems {
		item.InventarizationID = inventarization.ID
		item.Time = inventarization.Time
		item.SkladID = inventarization.SkladID
		item.Status = inventarization.Status //set all values to items
		if item.GroupID == 0 {               //item is visible if it is not in the group
			item.IsVisible = true
		}
		if item.Type == utils.TypeGroup {
			if item.ID == 0 { // If group's item.ID is 0 then it is new group
				newGroup = true
				continue
			}
			numOfIngredients := 0
			var allPlan float32
			for _, ing := range inventarization.InventarizationItems {
				if ing.GroupID == item.ItemID {
					numOfIngredients++
					allPlan = allPlan + ing.PlanQuantity //calculate all plan quantity of ingredients and number of ingredients
				}
			}
			difference := (item.FactQuantity - allPlan) / float32(numOfIngredients) //calculate the average difference
			for _, ing := range inventarization.InventarizationItems {
				if ing.GroupID == item.ItemID {
					ing.FactQuantity = ing.PlanQuantity + difference
					ing.Difference = ing.FactQuantity - ing.PlanQuantity
					ing.DifferenceSum = ing.Difference * ing.Cost
					inventarization.Result = inventarization.Result + ing.DifferenceSum
				}
			}
			item.Difference = 0
			item.DifferenceSum = 0
			item.FactQuantity = 0
		} else if item.GroupID == 0 {
			item.Difference = item.FactQuantity - item.PlanQuantity
			item.DifferenceSum = item.Difference * item.Cost
		}
		if item.Type == utils.TypeIngredient {
			ingMap[item.ItemID] = item
		} else if item.Type == utils.TypeTovar {
			tovarMap[item.ItemID] = item
		}
	}
	if newGroup {
		for _, item := range inventarization.InventarizationItems {
			if item.Type == utils.TypeGroup {
				group, err := s.repo.GetPureInventarizationGroup(item.ItemID)
				if err != nil {
					return nil, err
				}
				for _, ing := range group.Items {
					if group.Type == utils.TypeIngredient {
						if _, ok := ingMap[ing.ItemID]; !ok {

							newItem := &model.InventarizationItem{
								ItemID:            ing.ItemID,
								Time:              inventarization.Time,
								SkladID:           inventarization.SkladID,
								Type:              utils.TypeIngredient,
								IsVisible:         false,
								GroupID:           item.ItemID,
								InventarizationID: inventarization.ID,
							}
							if newItem.Time.After(time.Now()) {
								newItem.NeedToRecalculate = true
							}
							inventarization.InventarizationItems = append(inventarization.InventarizationItems, newItem)
						} else {
							ingMap[ing.ItemID].GroupID = item.ItemID
							ingMap[ing.ItemID].IsVisible = false
						}
					} else {
						if _, ok := tovarMap[ing.ItemID]; !ok {
							newItem := &model.InventarizationItem{
								ItemID:            ing.ItemID,
								Time:              inventarization.Time,
								SkladID:           inventarization.SkladID,
								Type:              utils.TypeTovar,
								IsVisible:         false,
								GroupID:           item.ItemID,
								InventarizationID: inventarization.ID,
							}
							if newItem.Time.After(time.Now()) {
								newItem.NeedToRecalculate = true
							}
							inventarization.InventarizationItems = append(inventarization.InventarizationItems, newItem)
						} else {
							tovarMap[ing.ItemID].GroupID = item.ItemID
							tovarMap[ing.ItemID].IsVisible = false
						}
					}
				}
			}
		}
	}
	return s.repo.UpdateInventarization(inventarization)
}

func (s *SkladService) UpdateInventarizationParams(inventarization *model.Inventarization) (*model.Inventarization, error) {

	ingMap := make(map[int]*model.InventarizationItem)
	tovarMap := make(map[int]*model.InventarizationItem)
	for _, item := range inventarization.InventarizationItems {
		item.SkladID = inventarization.SkladID
		if item.Type == utils.TypeGroup {
			item.Difference = 0
			item.DifferenceSum = 0
			item.FactQuantity = 0
			item.IsVisible = true
		}
		if item.GroupID == 0 {
			item.IsVisible = true
		}
		if item.Type == utils.TypeIngredient {
			ingMap[item.ItemID] = item
		} else if item.Type == utils.TypeTovar {
			tovarMap[item.ItemID] = item
		}
	}
	for _, item := range inventarization.InventarizationItems {
		if item.Type == utils.TypeGroup {
			group, err := s.repo.GetPureInventarizationGroup(item.ItemID)
			if err != nil {
				return nil, err
			}
			for _, ing := range group.Items {
				if group.Type == utils.TypeIngredient {
					if _, ok := ingMap[ing.ItemID]; !ok {
						newItem := &model.InventarizationItem{
							ItemID:            ing.ItemID,
							Time:              inventarization.Time,
							SkladID:           inventarization.SkladID,
							Type:              utils.TypeIngredient,
							IsVisible:         false,
							GroupID:           item.ItemID,
							InventarizationID: inventarization.ID,
						}
						if newItem.Time.After(time.Now()) {
							newItem.NeedToRecalculate = true
						}
						inventarization.InventarizationItems = append(inventarization.InventarizationItems, newItem)
					} else {
						ingMap[ing.ItemID].GroupID = item.ItemID
						ingMap[ing.ItemID].IsVisible = false
					}
				} else {
					if _, ok := tovarMap[ing.ItemID]; !ok {
						newItem := &model.InventarizationItem{
							ItemID:            ing.ItemID,
							Time:              inventarization.Time,
							SkladID:           inventarization.SkladID,
							Type:              utils.TypeTovar,
							IsVisible:         false,
							GroupID:           item.ItemID,
							InventarizationID: inventarization.ID,
						}
						if newItem.Time.After(time.Now()) {
							newItem.NeedToRecalculate = true
						}
						inventarization.InventarizationItems = append(inventarization.InventarizationItems, newItem)
					} else {
						tovarMap[ing.ItemID].GroupID = item.ItemID
						tovarMap[ing.ItemID].IsVisible = false
					}
				}
			}
		}
	}
	return s.repo.UpdateInventarizationParams(inventarization)
}

func (s *SkladService) OpenInventarization(inventarization *model.Inventarization) (*model.Inventarization, error) {
	inventarization.Status = utils.StatusOpened
	for i := range inventarization.InventarizationItems {
		inventarization.InventarizationItems[i].Status = utils.StatusOpened
	}
	return s.repo.UpdateInventarizationV2(inventarization)
}

func (s *SkladService) UpdateInventarizationV2(inventarization *model.Inventarization) (*model.Inventarization, error) {
	for i := 0; i < len(inventarization.InventarizationItems); i++ {
		item := inventarization.InventarizationItems[i]
		item.InventarizationID = inventarization.ID
		item.Time = inventarization.Time
		item.SkladID = inventarization.SkladID
		item.Status = inventarization.Status //set all values to items
		if item.Type == utils.TypeGroup {
			numOfIngredients := 0
			var allPlan float32
			for _, ing := range inventarization.InventarizationItems {
				if ing.GroupID == item.ItemID {
					numOfIngredients++
					allPlan = allPlan + ing.PlanQuantity //calculate all plan quantity of ingredients and number of ingredients
				}
			}
			difference := (item.FactQuantity - allPlan) / float32(numOfIngredients) //calculate the average difference
			for _, ing := range inventarization.InventarizationItems {
				if ing.GroupID == item.ItemID {
					ing.FactQuantity = ing.PlanQuantity + difference
					ing.Difference = ing.FactQuantity - ing.PlanQuantity
					ing.DifferenceSum = ing.Difference * ing.Cost
					inventarization.Result = inventarization.Result + ing.DifferenceSum
				}
			}
			item.Difference = 0
			item.DifferenceSum = 0
			item.FactQuantity = 0
		} else if item.GroupID == 0 {
			item.Difference = item.FactQuantity - item.PlanQuantity
			item.DifferenceSum = item.Difference * item.Cost
			item.IsVisible = true
		}
	}
	return s.repo.UpdateInventarizationV2(inventarization)
}

func (s *SkladService) GetInventarization(id int) (*model.InventarizationResponse, error) {
	return s.repo.GetInventarization(id)
}

func (s *SkladService) GetAllInventarization(filter *model.Filter) ([]*model.InventarizationResponse, int64, error) {
	res, count, err := s.repo.GetAllInventarization(filter)
	if err != nil {
		return nil, 0, err
	}
	pageCount := utils.CalculatePageCount(count, utils.DefaultPageSize)
	return res, pageCount, nil
}

func (s *SkladService) DeleteInventarization(id int) error {
	return s.repo.DeleteInventarization(id)
}

func (s *SkladService) DeleteInventarizationItem(id int) error {
	return s.repo.DeleteInventarizationItem(id)
}

func (s *SkladService) GetSkladByShopID(shopID int) (*model.Sklad, error) {
	return s.repo.GetSkladByShopID(shopID)
}

func (s *SkladService) GetInventarizationDetailsIncome(id int) ([]*model.InventarizationDetailsIncome, error) {
	return s.repo.GetInventarizationDetailsIncome(id)
}
func (s *SkladService) GetInventarizationDetailsExpence(id int) ([]*model.InventarizationDetailsExpence, error) {
	return s.repo.GetInventarizationDetailsExpence(id)
}
func (s *SkladService) GetInventarizationDetailsSpisanie(id int) ([]*model.InventarizationDetailsSpisanie, error) {
	return s.repo.GetInventarizationDetailsSpisanie(id)
}

func (s *SkladService) GetTrafficReport(filter *model.Filter) (*model.GlobalTrafficReport, int64, error) {
	res, count, err := s.repo.GetTrafficReport(filter)
	if err != nil {
		return nil, 0, err
	}
	trafficReport := &model.GlobalTrafficReport{
		TrafficReports: res,
	}
	for _, report := range res {
		trafficReport.InitialSum = trafficReport.InitialSum + report.InitialSum
		trafficReport.FinalSum = trafficReport.FinalSum + report.FinalSum
	}
	pageCount := utils.CalculatePageCount(count, utils.DefaultPageSize)
	return trafficReport, pageCount, nil
}

func (s *SkladService) DailyStatistic() error {
	shops, err := s.repo.GetAllShopsWithouParam()

	if err != nil {
		return err
	}

	for _, shop := range shops {
		s.repo.DailyStatistic(shop.ID)
	}

	return nil
}
func (s *SkladService) CreateInventarizationGroup(group *model.ReqInventarizationGroup) error {
	groupItems, err := s.CheckItems(group)
	if err != nil {
		return err
	}
	groupToAdd := &model.InventarizationGroup{
		Name:    group.Name,
		Measure: group.Measure,
		SkladID: group.SkladID,
		Items:   groupItems,
		Type:    group.Type,
	}
	return s.repo.CreateInventarizationGroup(groupToAdd)
}
func (s *SkladService) GetAllInventarizationGroup(filter *model.Filter) ([]*model.InventarizationGroupResponse, error) {
	return s.repo.GetAllInventarizationGroup(filter)
}
func (s *SkladService) GetInventarizationGroup(filter *model.Filter, id int) (*model.InventarizationGroupResponse, error) {
	return s.repo.GetInventarizationGroup(filter, id)
}
func (s *SkladService) CheckItems(group *model.ReqInventarizationGroup) ([]*model.InventarizationGroupItem, error) {
	groupItems := make([]*model.InventarizationGroupItem, 0)
	sklad, err := s.repo.GetSklad(group.SkladID)
	if err != nil {
		return nil, err
	}
	shopID := sklad.ShopID
	for _, item := range group.Items {
		groupItem := &model.InventarizationGroupItem{}
		if group.Type == utils.TypeIngredient {
			ing, err := s.repo.GetPureIngredientByShopID(item.ItemID, shopID)
			if err != nil {
				return nil, err
			}
			if ing == nil {
				return nil, errors.New("ingredient not found")
			}

			if ing.Measure != group.Measure {
				return nil, errors.New("measure is not equal")
			}
			ok := s.repo.CheckUnique(item.ItemID, group.SkladID, group.Type, group.ID)
			if ok != nil {
				return nil, ok
			}
			groupItem.ItemID = ing.IngredientID
			groupItem.SkladID = group.SkladID
		} else {
			tovar, err := s.repo.GetPureTovarByShopID(item.ItemID, shopID)
			if err != nil {
				return nil, err
			}
			if tovar == nil {
				return nil, errors.New("tovar not found")
			}

			if group.Measure != tovar.Measure {
				return nil, errors.New("measure is not equal")
			}

			ok := s.repo.CheckUnique(item.ItemID, group.SkladID, group.Type, group.ID)
			if ok != nil {
				return nil, ok
			}
			groupItem.ItemID = tovar.TovarID
			groupItem.SkladID = group.SkladID
		}
		groupItems = append(groupItems, groupItem)
	}
	return groupItems, nil
}
func (s *SkladService) UpdateInventarizationGroup(group *model.ReqInventarizationGroup) error {
	groupItems, err := s.CheckItems(group)
	if err != nil {
		return err
	}
	groupToAdd := &model.InventarizationGroup{
		ID:      group.ID,
		Name:    group.Name,
		Measure: group.Measure,
		SkladID: group.SkladID,
		Items:   groupItems,
		Type:    group.Type,
	}
	return s.repo.UpdateInventarizationGroup(groupToAdd)
}
func (s *SkladService) DeleteInventarizationGroup(id int) error {
	return s.repo.DeleteInventarizationGroup(id)
}

func (s *SkladService) RecalculateTrafficReport() error {
	items, err := s.repo.GetItemsForRecalculateTrafficReport()
	if err != nil {
		return err
	}

	for _, item := range items {
		if item.RetryCount < 5 {
			err := s.repo.RecalculateTrafficReport(item)
			item.RetryCount++
			if err != nil {
				item.Exception = err.Error()
			} else {
				item.Status = utils.StatusRecalculated
				item.Exception = utils.StatusSuccess
				item.FinishedAt = time.Now()
			}
			err = s.repo.UpdateTrafficReportJob(item)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (s *SkladService) RecalculateDailyStatistic() error {
	sklads, err := s.repo.GetAllSklads()
	if err != nil {
		return err
	}
	from := time.Date(2023, 7, 17, 0, 0, 0, 0, time.UTC)
	to := time.Date(2023, 8, 26, 0, 0, 0, 0, time.UTC)
	for from.Before(to) || from.Equal(to) {
		for _, sklad := range sklads {
			err := s.repo.RecalculateDailyStatisticByDate(from, sklad.ID)
			if err != nil {
				return err
			}
		}
		from = from.AddDate(0, 0, 1)
	}
	return nil
}
