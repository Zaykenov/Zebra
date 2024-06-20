package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math"
	"net/http"
	"os"
	"time"
	"zebra/model"
	"zebra/pkg/repository"
	"zebra/utils"
)

type CheckService struct {
	repo repository.Repository
}

func NewCheckService(repo repository.Repository) *CheckService {
	return &CheckService{repo: repo}
}

func (s *CheckService) SaveError(fullError string, request string) {
	s.repo.SaveError(fullError, request)
}

func (s *CheckService) AddCheck(check *model.ReqCheck) (*model.CheckResponse, error) {
	gormCheck, err := s.repo.CalculateCheck(check)
	if err != nil {
		return nil, err
	}
	loc := time.FixedZone("UTC+6", 5*60*60)
	now := time.Now().In(loc)

	gormCheck.Opened_at = now
	gormCheck.IdempotencyKey = check.IdempotencyKey
	res, err := s.repo.AddCheck(gormCheck)
	if err != nil {
		return nil, err
	}
	/*res, err = s.repo.RemoveFromSklad(res)
	if err != nil {
		return nil, err
	}
	*/

	return res, nil
}

func (s *CheckService) GetCheckByIdempotency(idempotency string) (*model.Check, error) {
	return s.repo.GetCheckByIdempotency(idempotency)
}

func (s *CheckService) SendToTis() error {
	checks, err := s.repo.GetUnsendTisCheck()
	if err != nil {
		return err
	}

	for _, check := range checks {

		newCheck, err := s.SendTis(check)
		if err != nil {
			return err
		}

		check = newCheck
		check.RetryCount++
		err = s.repo.UpdateToSend(check)

		if err != nil {
			return err
		}
	}
	return err
}

func (s *CheckService) SendTis(check *model.SendToTis) (*model.SendToTis, error) {
	url := os.Getenv("TIS_URL")
	body := []byte(check.Request)
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}
	req.Header.Add("Content-type", "application/json")
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Idempotency-Key", check.IdempotencyKey)
	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	var TisResponse model.TisData
	err = json.NewDecoder(resp.Body).Decode(&TisResponse)
	if err != nil {
		return nil, err
	}
	response, err := json.Marshal(TisResponse)
	if err != nil {
		return nil, err
	}
	if (resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusOK) || TisResponse.Data == nil {
		check.Status = utils.StatusError
		errorResponse, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			check.Exception = err.Error()
		} else {
			check.Exception = string(errorResponse)
		}
	} else {
		check.Status = utils.StatusSuccess
		check.Response = string(response)
		if TisResponse.Data != nil {
			if TisResponse.Data.Link != "" {
				err := s.repo.UpdateCheckLink(check.CheckID, TisResponse.Data.Link)
				if err != nil {
					return nil, err
				}
			}
		}
	}
	return check, nil
}

func (s *CheckService) UpdateCheck(check *model.ReqCheck) (*model.CheckResponse, error) {
	if check.Status == utils.StatusClosed {
		check.Closed_at = time.Now()
	}
	gormCheck, err := s.repo.CalculateCheck(check)
	if err != nil {
		return nil, err
	}
	gormCheck.IdempotencyKey = check.IdempotencyKey
	if gormCheck.Status == utils.StatusClosed {
		gormCheck.Cash = gormCheck.Sum - gormCheck.Discount - gormCheck.Card
		if gormCheck.Card < 1 {
			gormCheck.Cash = gormCheck.Cash + gormCheck.Card
			gormCheck.Card = 0
		}

		if gormCheck.Cash < 1 {
			gormCheck.Card = gormCheck.Cash + gormCheck.Card
			gormCheck.Cash = 0
		}
		gormCheck.Closed_at = time.Now()
		newCheck, err := s.repo.CloseCheck(gormCheck)
		if err != nil {
			return nil, err
		}
		res := s.repo.ConstructResponse(newCheck)

		shop, err := s.repo.GetShop(res.ShopID)
		if err != nil {
			return nil, err
		}
		revenue, err := s.repo.GetRevenueByShopID(res.ShopID)
		if err != nil {
			return nil, err
		}
		if (revenue < shop.Limit || shop.Limit == -1) && check.OFD {
			if shop.CassaType == utils.CassaTypeTis {
				tisBody, err := s.CreateTisBody(res, utils.TisTypeSale)
				if err != nil {
					return nil, err
				}
				bodyJson, err := json.Marshal(tisBody)
				if err != nil {
					return nil, err
				}
				sendToTis := &model.SendToTis{
					IdempotencyKey: res.IdempotencyKey,
					CheckID:        res.ID,
					Created_at:     time.Now().Local(),
					Request:        string(bodyJson),
					Status:         utils.StatusNew,
					RetryCount:     0,
					CassaType:      utils.CassaTypeTis,
				}
				err = s.repo.AddCheckToSend(sendToTis)
				if err != nil {
					return nil, err
				}
			} else if shop.CassaType == utils.CassaTypeWK {
				// wkBody, err := s.CreateWkBody(res, utils.TisTypeSale)
				// if err != nil {
				// 	return nil, err
				// }
				// bodyJson, err := json.Marshal(wkBody)
				// if err != nil {
				// 	return nil, err
				// }
				// sendToWk := &model.SendToTis{
				// 	IdempotencyKey: res.IdempotencyKey,
				// 	CheckID:        res.ID,
				// 	Created_at:     time.Now().Local(),
				// 	Request:        string(bodyJson),
				// 	Status:         utils.StatusNew,
				// 	RetryCount:     0,
				// 	CassaType:      utils.CassaTypeWK,
				// }
				// err = s.repo.AddCheckToSend(sendToWk)
				// if err != nil {
				// 	return nil, err
				// }
			}
		}
		res.Link = fmt.Sprintf("https://zebra-crm.kz/%s/%d", "receipt", res.ID)
		return res, nil
	}
	res, err := s.repo.UpdateCheck(gormCheck)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (s *CheckService) CreateWkBody(check *model.CheckResponse, tisType int) (*model.WKBody, error) {
	price := 0
	items := []*model.WKPositions{}
	for _, v := range check.Tovar {
		item := &model.WKPositions{
			PositionName: v.TovarName,
			Count:        int(v.Quantity),
			Price:        int(math.Floor(float64(v.Price))),
			TaxPercent:   12,
			Tax:          math.Floor(float64(v.Price * v.Quantity * 0.12)),
			TaxType:      100,
			UnitCode:     796,
		}
		price = price + item.Price
		items = append(items, item)
	}
	shop, err := s.repo.GetShop(check.ShopID)
	if err != nil {
		return nil, err
	}
	body := &model.WKBody{
		Token:               shop.TisToken,
		CashboxUniqueNumber: shop.CashboxUniqueNumber,
	}
	return body, nil
}

func (s *CheckService) CreateTisBody(check *model.CheckResponse, tisType int) (*model.TisBody, error) {
	token, err := s.repo.GetTisToken(check.ShopID)
	if err != nil {
		return nil, err
	}
	price := float32(0.0)
	items := []*model.TisItems{}
	for _, v := range check.Tovar {
		item := &model.TisItems{
			Name:     v.TovarName,
			Quantity: v.Quantity,
			Price:    (v.Price),
			Discount: (v.Discount),
			KgdCode:  796,
			CompareField: &model.TisCompareField{
				Type:  "barcode",
				Value: v.TovarName,
			},
		}
		price = price + (v.Price * v.Quantity) - v.Discount
		items = append(items, item)
	}
	for _, v := range check.TechCart {
		item := &model.TisItems{
			Name:     v.TechCartName,
			Quantity: v.Quantity,
			Price:    v.Price,
			Discount: v.Discount,
			KgdCode:  796,
			CompareField: &model.TisCompareField{
				Type:  "barcode",
				Value: v.TechCartName,
			},
		}
		price = price + (v.Price * v.Quantity) - v.Discount
		items = append(items, item)
	}
	check.Card = price - check.Cash
	if check.Card > 0 && check.Card < 1 {
		check.Cash = price
	}
	tisPayment := []*model.TisPayments{}
	if check.Card >= 1 {
		payment := &model.TisPayments{
			PaymentMethod: 1,
			Sum:           check.Card,
		}
		tisPayment = append(tisPayment, payment)
	}
	if check.Cash >= 1 {
		payment := &model.TisPayments{
			PaymentMethod: 0,
			Sum:           check.Cash,
		}
		tisPayment = append(tisPayment, payment)
	}

	if len(tisPayment) == 0 {
		payment := &model.TisPayments{
			PaymentMethod: 0,
			Sum:           1,
		}
		tisPayment = append(tisPayment, payment)
	}

	tisBody := model.TisBody{
		Token:    token,
		Type:     tisType,
		Items:    items,
		Payments: tisPayment,
	}
	return &tisBody, nil
}

func (s *CheckService) GetAllWorkerCheck(filter *model.Filter) ([]*model.CheckResponse, int64, error) {
	res, count, err := s.repo.GetAllWorkerCheck(filter)
	if err != nil {
		return nil, 0, err
	}
	pageCount := utils.CalculatePageCount(count, utils.ChecksPageSize)
	return res, pageCount, nil
}

func (s *CheckService) GetCheckByIDForPrinter(id int) (*model.CheckPrint, error) {
	check, err := s.repo.GetCheckByID(id)
	if err != nil {
		return nil, err
	}
	res := &model.CheckPrint{
		ID:              check.ID,
		UserID:          -1,
		Opened_at:       check.Opened_at,
		Closed_at:       check.Closed_at,
		WorkerID:        check.WorkerID,
		Discount:        check.Discount,
		DiscountPercent: check.DiscountPercent,
		Sum:             check.Sum,
		Cost:            check.Cost,
		Status:          check.Status,
		Payment:         check.Payment,
		Tovar:           check.Tovar,
		TechCart:        check.TechCart,
		Comment:         check.Comment,
		Feedback:        check.Feedback,
		TisCheckUrl:     check.Link,
	}

	return res, nil
}

func (s *CheckService) GetAllCheck(filter *model.Filter) (*model.GlobalCheck, int64, error) {
	res, count, err := s.repo.GetAllCheck(filter)
	if err != nil {
		return nil, 0, err
	}
	pageCount := utils.CalculatePageCount(count, utils.ChecksPageSize)
	checks := &model.GlobalCheck{
		Check: res,
	}
	for _, check := range res {
		if check.Status == utils.StatusInactive {
			continue
		}
		checks.TotalMoney += check.Card + check.Cash
		checks.TotalCard += check.Card
		checks.TotalCash += check.Cash
		checks.TotalNetCost += check.Cost
		checks.TotalDiscount += check.Discount
	}
	checks.TotalProfit = checks.TotalMoney - checks.TotalNetCost
	return checks, pageCount, nil
}

func (s *CheckService) GetCheckByID(id int) (*model.CheckResponse, error) {
	return s.repo.GetCheckByID(id)
}

func (s *CheckService) GetAllCheckView(page int) ([]*model.CheckView, int64, error) {
	checks, count, err := s.repo.GetAllCheckView(page)
	if err != nil {
		return nil, 0, err
	}
	pageCount := utils.CalculatePageCount(count, utils.ChecksPageSize)
	return checks, pageCount, nil
}

func (s *CheckService) AddTag(tag *model.Tag) error {
	return s.repo.AddTag(tag)
}

func (s *CheckService) GetAllTag(shopID int) ([]*model.Tag, error) {
	return s.repo.GetAllTag(shopID)
}

func (s *CheckService) GetTag(id int) (*model.Tag, error) {
	return s.repo.GetTag(id)
}

func (s *CheckService) UpdateTag(tag *model.Tag) error {
	return s.repo.UpdateTag(tag)
}

func (s *CheckService) DeleteTag(id int) error {
	return s.repo.DeleteTag(id)
}

func (s *CheckService) DeleteCheck(id int) error {
	return s.repo.DeleteCheck(id)
}

func (s *CheckService) DeactivateCheck(id int) error {

	items, err := s.repo.DeactivateCheck(id)
	if err != nil {
		return err
	}

	// trafficItems := []*model.AsyncJob{}
	// for _, item := range items {
	// 	if item.Time.Year() == time.Now().Year() && item.Time.Month() == time.Now().Month() && item.Time.Day() != time.Now().Day() {
	// 		trafficItems = append(trafficItems, &model.AsyncJob{
	// 			ItemID:    item.ItemID,
	// 			ItemType:  item.Type,
	// 			SkladID:   item.SkladID,
	// 			TimeStamp: item.Time,
	// 			CreatedAt: time.Now(),
	// 			Status:    utils.StatusNeedRecalculate,
	// 		})
	// 	}
	// }
	// err = s.repo.AddTrafficReportJob(trafficItems)
	// if err != nil {
	// 	return err
	// }

	check, err := s.repo.GetCheckByID(id)
	if err != nil {
		return err
	}

	if check.Opened_at.Year() == time.Now().Year() && check.Opened_at.Month() == time.Now().Month() && check.Opened_at.Day() != time.Now().Day() {
		trafficItems := []*model.AsyncJob{}
		for _, item := range check.Tovar {
			trafficItems = append(trafficItems, &model.AsyncJob{
				ItemID:    item.TovarID,
				ItemType:  utils.TypeTovar,
				SkladID:   check.SkladID,
				ShopID:    check.ShopID,
				TimeStamp: check.Opened_at,
				CreatedAt: time.Now(),
				Status:    utils.StatusNeedRecalculate,
			})
		}
		for _, techCart := range check.TechCart {
			ingredients, err := s.repo.GetIngredientsByTechCart(techCart.TechCartID)
			if err != nil {
				return err
			}
			for _, item := range ingredients {
				trafficItems = append(trafficItems, &model.AsyncJob{
					ItemID:    item.ID,
					ItemType:  utils.TypeIngredient,
					SkladID:   check.SkladID,
					ShopID:    check.ShopID,
					TimeStamp: check.Opened_at,
					CreatedAt: time.Now(),
					Status:    utils.StatusNeedRecalculate,
				})
			}
		}
		s.repo.ConcurrentRecalculationForDailyStatistics(trafficItems)
	}

	tisBody, err := s.CreateTisBody(check, utils.TisTypeRefund)
	if err != nil {
		return err
	}

	bodyJson, err := json.Marshal(tisBody)
	if err != nil {
		return err
	}

	sendToTis := &model.SendToTis{
		IdempotencyKey: check.IdempotencyKey + "_refund",
		CheckID:        check.ID,
		Created_at:     time.Now().Local(),
		Request:        string(bodyJson),
		Status:         utils.StatusNew,
		RetryCount:     0,
	}

	err = s.repo.AddCheckToSend(sendToTis)
	if err != nil {
		return err
	}
	if len(items) > 0 {
		err = s.repo.Sklad.RecalculateInventarization(items)
		if err != nil {
			return err
		}
	}

	shift, err := s.repo.Transaction.GetShiftByTime(check.Closed_at, check.ShopID)
	if err != nil {
		return err
	}

	err = s.repo.Transaction.RecalculateShift(shift.ID)
	if err != nil {
		return err
	}

	return nil
}

func (s *CheckService) IdempotencyCheck(keys *model.IdempotencyCheckArray, shopID int) error {
	res, err := s.repo.IdempotencyCheck(keys, shopID)
	if err != nil {
		s.repo.SaveError(res, "idempotencyKeyCheck")
	}
	return err

}

func (s *CheckService) SaveFailedCheck(checks []*model.FailedCheck) error {
	return s.repo.SaveFailedCheck(checks)
}

func (s *CheckService) GetStoliki(shopID int) ([]*model.Stolik, error) {
	return s.repo.GetStoliki(shopID)
}

func (s *CheckService) GetFilledStoliki(shopID int) ([]*model.Stolik, error) {
	return s.repo.GetFilledStoliki(shopID)
}
