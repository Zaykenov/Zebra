package service

import (
	"errors"
	"math"
	"time"
	"zebra/model"
	"zebra/pkg/repository"
	"zebra/utils"

	"gorm.io/gorm"
)

type TransactionService struct {
	repo repository.Repository
}

func NewTransactionService(repo repository.Repository) *TransactionService {
	return &TransactionService{repo: repo}
}

func (r *TransactionService) OpenShift(id int, transaction *model.Transaction) error {
	shop, err := r.repo.GetShopBySchetID(transaction.SchetID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return errors.New("shop not found")
		}
		return err
	}
	res, err := r.repo.GetLastShift(shop.ID)
	if err != nil {
		if err.Error() != errors.New("record not found").Error() {
			return err
		}
	}
	if res != nil {
		if !res.IsClosed {
			return errors.New("close shift")
		}
	}
	shift := &model.Shift{
		IsClosed:  false,
		CreatedAt: transaction.Time,
		BeginSum:  transaction.Sum,
		SchetID:   transaction.SchetID,
		ShopID:    shop.ID,
	}
	shiftWithID, err := r.repo.OpenShift(id, shift)
	if err != nil {
		return err
	}
	transaction.Status = utils.TransactionNeutralStatus
	transaction.ShiftID = shiftWithID.ID
	_, err = r.repo.CreateTransaction(transaction)
	if err != nil {
		return err
	}
	return nil
}

func (r *TransactionService) CloseShift(id int, transaction *model.Transaction) error {
	if transaction.ShiftID != 0 {
		res, err := r.repo.GetShiftByID2(transaction.ShiftID)
		if err != nil {
			return err
		}

		if res.IsClosed {
			return errors.New("shift is closed")
		}

		if transaction.Time.Before(res.CreatedAt) {
			return errors.New("time is incorrect; should be in range of shift time")
		}

		res.IsClosed = true
		res.ClosedAt = transaction.Time
		res.EndSumFact = transaction.Sum
		res.EndSumPlan = res.BeginSum + res.Cash - res.Expense - res.Collection
		res.Difference = float32(math.Abs(float64(res.EndSumFact - res.EndSumPlan)))

		transaction.Status = utils.TransactionNeutralStatus
		transaction.SchetID = res.SchetID
		_, err = r.repo.CreateTransaction(transaction)
		if err != nil {
			return err
		}
		err = r.repo.UpdateShift(res)
		if err != nil {
			return err
		}
	} else {
		shop, err := r.repo.GetShopBySchetID(transaction.SchetID)
		if err != nil {
			return errors.New("shop not found")
		}
		res, err := r.repo.GetLastShift(shop.ID)
		if err != nil {
			if err.Error() == errors.New("record not found").Error() {
				return errors.New("open shift")
			}
			return err
		}
		if res.IsClosed {
			return errors.New("open shift")
		}

		openCheckCount, err := r.repo.CountOpenCheck(id, transaction.Time)
		if err != nil {
			return err
		}

		if openCheckCount > 0 {
			return errors.New("close opened checks")
		}

		res.IsClosed = true
		res.ClosedAt = transaction.Time
		res.EndSumFact = transaction.Sum
		res.EndSumPlan = res.BeginSum + res.Cash - res.Expense - res.Collection
		res.Difference = float32(math.Abs(float64(res.EndSumFact - res.EndSumPlan)))
		err = r.repo.UpdateShift(res)
		if err != nil {
			return err
		}
		transaction.Status = utils.TransactionNeutralStatus
		transaction.ShiftID = res.ID
		_, err = r.repo.CreateTransaction(transaction)
		if err != nil {
			return err
		}
	}
	return nil
}

func (r *TransactionService) Collection(role string, id int, transaction *model.Transaction) error {
	if role == utils.WorkerRole {
		shop, err := r.repo.GetShopBySchetID(transaction.SchetID)
		if err != nil {
			return errors.New("shop not found")
		}
		res, err := r.repo.GetLastShift(shop.ID)
		if err != nil {
			if err.Error() == errors.New("record not found").Error() {
				return errors.New("open shift")
			}
			return err
		}
		if res.IsClosed {
			return errors.New("open shift")
		}
		res.Collection += transaction.Sum
		err = r.repo.UpdateShift(res)
		if err != nil {
			return err
		}
		transaction.ShiftID = res.ID

	}
	if transaction.ShiftID != 0 {
		res, err := r.repo.GetShiftByID2(transaction.ShiftID)
		if err != nil {
			return err
		}
		if res.IsClosed {
			if transaction.Time.After(res.ClosedAt) {
				return errors.New("time is incorrect; should be past time")
			}
		} else {
			if transaction.Time.After(time.Now()) {
				return errors.New("time is incorrect; should be past time")
			}
		}
		if transaction.Time.Before(res.CreatedAt) {
			return errors.New("time is incorrect; should be in range of shift time")
		}
		transaction.SchetID = res.SchetID
	} else {
		shop, err := r.repo.GetShopByCashSchetID(transaction.SchetID)
		if err != nil {
			if err != gorm.ErrRecordNotFound {
				return err
			}
		}
		if shop != nil {
			shift, err := r.repo.GetShiftByTime(transaction.Time, shop.ID)
			if err != nil {
				if err.Error() == errors.New("record not found").Error() {
					return errors.New("shift by time not found")
				}
				return err
			}
			transaction.ShiftID = shift.ID
		}

	}
	_, err := r.repo.CreateTransaction(transaction)
	if err != nil {
		return err
	}
	return nil
}

func (r *TransactionService) Income(role string, id int, transaction *model.Transaction) error {
	if role == utils.WorkerRole {
		shop, err := r.repo.GetShopBySchetID(transaction.SchetID)
		if err != nil {
			return errors.New("shop not found")
		}
		res, err := r.repo.GetLastShift(shop.ID)
		if err != nil {
			if err.Error() == errors.New("record not found").Error() {
				return errors.New("open shift")
			}
			return err
		}
		if res.IsClosed {
			return errors.New("open shift")
		}
		res.Income += transaction.Sum
		err = r.repo.UpdateShift(res)
		if err != nil {
			return err
		}
		transaction.ShiftID = res.ID

	}
	if transaction.ShiftID != 0 {
		res, err := r.repo.GetShiftByID2(transaction.ShiftID)
		if err != nil {
			return err
		}
		if res.IsClosed {
			if transaction.Time.After(res.ClosedAt) {
				return errors.New("time is incorrect; should be past time")
			}
		} else {
			if transaction.Time.After(time.Now()) {
				return errors.New("time is incorrect; should be past time")
			}
		}
		if transaction.Time.Before(res.CreatedAt) {
			return errors.New("time is incorrect; should be in range of shift time")
		}
		transaction.SchetID = res.SchetID
	} else {
		shop, err := r.repo.GetShopByCashSchetID(transaction.SchetID)
		if err != nil {
			if err != gorm.ErrRecordNotFound {
				return err
			}
		}
		if shop != nil {
			shift, err := r.repo.GetShiftByTime(transaction.Time, shop.ID)
			if err != nil {
				if err.Error() == errors.New("record not found").Error() {
					return errors.New("shift by time not found")
				}
				return err
			}
			transaction.ShiftID = shift.ID
		}

	}
	_, err := r.repo.CreateTransaction(transaction)
	if err != nil {
		return err
	}
	return nil
}

func (r *TransactionService) PostavkaTransaction(role string, transaction *model.Transaction, postavkaID, shopID int) error {
	res, err := r.repo.GetShiftByTime(transaction.Time, shopID)
	if err != nil {
		if err.Error() == errors.New("record not found").Error() {
			return errors.New("open shift")
		}
		return err
	}
	/*if res.IsClosed {
		return errors.New("open shift")
	}*/
	transaction.ShiftID = res.ID

	if transaction.SchetID == 0 {
		transaction.Status = utils.TransactionNeutralStatus
	} else {
		transaction.Status = utils.TransactionNegativeStatus
	}

	transaction, err = r.repo.CreateTransaction(transaction)
	if err != nil {
		return err
	}
	transactionPostavka := &model.TransactionPostavka{
		TransactionID: transaction.ID,
		PostavkaID:    postavkaID,
	}
	err = r.repo.CreatePostavkaTransaction(transactionPostavka)
	if err != nil {
		return err
	}
	return nil
}

func (r *TransactionService) GetAllTransaction(filter *model.Filter) ([]*model.TransactionResponse, int64, error) {
	res, count, err := r.repo.GetAllTransaction(filter)
	if err != nil {
		return nil, 0, err
	}
	pageCount := utils.CalculatePageCount(count, utils.DefaultPageSize)
	return res, pageCount, nil
}

func (r *TransactionService) GetAllShifts(filter *model.Filter) ([]*model.Shift, int64, error) {
	res, count, err := r.repo.GetAllShifts(filter)
	if err != nil {
		return nil, 0, err
	}
	pageCount := utils.CalculatePageCount(count, utils.DefaultPageSize)
	return res, pageCount, nil
}

func (r *TransactionService) GetShiftByID(id int) (*model.ShiftTransaction, error) {
	return r.repo.GetShiftByID(id)
}

func (r *TransactionService) GetShiftByShopId(id int, workerID int) (*model.CurrentShift, error) {
	return r.repo.GetShiftByShopId(id, workerID)
}

func (r *TransactionService) GetLastShift(shopID int) (*model.Shift, error) {
	return r.repo.GetLastShift(shopID)
}

func (r *TransactionService) GetTransactionByID(id int) (*model.TransactionResponse, error) {
	return r.repo.GetTransactionByID(id)
}

func (r *TransactionService) UpdateTransaction(transaction *model.Transaction) error {
	return r.repo.UpdateTransaction(transaction)
}

func (r *TransactionService) DeleteTransaction(id int) error {
	return r.repo.DeleteTransaction(id)
}

func (r *TransactionService) CheckForBlockedShop(shiftID int) (bool, error) {
	return r.repo.CheckForBlockedShop(shiftID)
}
