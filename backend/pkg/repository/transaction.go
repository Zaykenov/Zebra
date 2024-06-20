package repository

import (
	"database/sql"
	"errors"
	"time"
	"zebra/model"
	"zebra/utils"

	"gorm.io/gorm"
)

type TransactionDB struct {
	db     *sql.DB
	gormDB *gorm.DB
}

func NewTransactionDB(db *sql.DB, gormDB *gorm.DB) *TransactionDB {
	return &TransactionDB{db: db, gormDB: gormDB}
}

func (r *TransactionDB) RecalculateShift(id int) error {
	shift := &model.Shift{}
	err := r.gormDB.Model(&model.Shift{}).Where("id = ?", id).Scan(&shift).Error

	if err != nil {
		return err
	}

	payment := model.Payment{}
	time := time.Now()
	if shift.IsClosed {
		time = shift.ClosedAt
	}
	err = r.gormDB.Table("checks").Select("SUM(checks.cash) AS cash, SUM(checks.card) AS card").Where("checks.status = ? AND checks.closed_at >= ? AND checks.closed_at <= ? and checks.shop_id = ?", utils.StatusClosed, shift.CreatedAt, time, shift.ShopID).Scan(&payment).Error
	if err != nil {
		return err
	}
	shift.Cash = payment.Cash
	shift.Card = payment.Card
	shift.Expense = 0
	shift.Income = 0
	shift.EndSumPlan = 0
	shift.Collection = 0
	transactions := []*model.Transaction{}
	err = r.gormDB.Model(&model.Transaction{}).Where("shift_id = ? and deleted = ?", shift.ID, false).Order("id DESC").Scan(&transactions).Error
	if err != nil {
		return err
	}
	for _, transaction := range transactions {
		if transaction.SchetID == shift.SchetID {
			if transaction.Category == utils.Collection {
				shift.Collection += transaction.Sum
			} else if transaction.Category == utils.Income {
				shift.Income += transaction.Sum
			} else {
				if transaction.Status == utils.TransactionNegativeStatus {
					shift.Expense += transaction.Sum
				} else if transaction.Status == utils.TransactionPositiveStatus {
					shift.Cash += transaction.Sum
				}
			}
		}
	}
	shift.EndSumPlan = shift.BeginSum + shift.Cash - shift.Expense - shift.Collection + shift.Income
	if shift.IsClosed {
		shift.Difference = shift.EndSumFact - shift.EndSumPlan
	}
	// err = r.gormDB.Model(&model.Shift{}).Select("begin_sum, end_sum_fact, end_sum_plan, expense, cash,card,collection, difference, equal_cash").Where("id = ?", shift.ID).Updates(shift).Scan(shift).Error
	err = r.gormDB.Model(&model.Shift{}).Where("id = ?", shift.ID).Save(&shift).Error
	if err != nil {
		return err
	}
	return err
}

func (r *TransactionDB) CountOpenCheck(id int, time time.Time) (int, error) {
	OpenChecks := []*model.Check{}

	res := r.gormDB.Table("checks").Where("worker_id = ? and opened_at::date = ?::date and status = ?", id, time, utils.StatusOpened).Scan(&OpenChecks)
	if res.Error != nil {
		return 0, res.Error
	}

	if len(OpenChecks) > 0 {
		return 1, nil
	}

	return 0, nil
}

func (r *TransactionDB) OpenShift(id int, shift *model.Shift) (*model.Shift, error) {
	err := r.gormDB.Model(&model.Shift{}).Create(shift).Error
	if err != nil {
		return nil, err
	}
	return shift, nil

}

func (r *TransactionDB) CloseShift(id int, shift *model.Shift) error {
	return nil
}

func (r *TransactionDB) UpdateShift(shift *model.Shift) error {
	err := r.gormDB.Model(&model.Shift{}).Where("id = ?", shift.ID).Updates(shift).Error
	if err != nil {
		return err
	}
	err = r.RecalculateShift(shift.ID)
	if err != nil {
		return err
	}
	return nil
}

func (r *TransactionDB) CreateTransaction(transaction *model.Transaction) (*model.Transaction, error) {
	err := r.gormDB.Transaction(func(tx *gorm.DB) error {
		err := r.gormDB.Debug().Model(&model.Transaction{}).Create(transaction).Error
		if err != nil {
			return err
		}
		if transaction.Status == utils.TransactionNegativeStatus {
			res := tx.Table("schets").Where("id = ?", transaction.SchetID).Update("start_balance", gorm.Expr("start_balance - ?", transaction.Sum))
			if res.Error != nil {
				return res.Error
			}
		} else if transaction.Status == utils.TransactionPositiveStatus {
			res := tx.Table("schets").Where("id = ?", transaction.SchetID).Update("start_balance", gorm.Expr("start_balance + ?", transaction.Sum))
			if res.Error != nil {
				return res.Error
			}
		}
		return nil
	})

	if err != nil {
		return nil, err
	}
	if transaction.ShiftID != 0 {
		err = r.RecalculateShift(transaction.ShiftID)
		if err != nil {
			return nil, err
		}
	}
	return transaction, nil
}

func (r *TransactionDB) CreatePostavkaTransaction(transactionPostavka *model.TransactionPostavka) error {
	err := r.gormDB.Model(&model.TransactionPostavka{}).Create(transactionPostavka).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *TransactionDB) GetLastShift(shopID int) (*model.Shift, error) {
	shift := &model.Shift{}
	err := r.gormDB.Model(&model.Shift{}).Where("shop_id = ?", shopID).Last(&model.Shift{}).Scan(shift).Error
	if err != nil {
		return nil, err
	}
	payment := model.Payment{}
	err = r.gormDB.Table("checks").Select("SUM(checks.cash) AS cash, SUM(checks.card) AS card").Where("checks.status = ? AND checks.closed_at::date = ?::date and checks.shop_id = ?", utils.StatusClosed, shift.CreatedAt, shopID).Scan(&payment).Error
	if err != nil {
		return nil, err
	}
	shift.Cash = payment.Cash
	shift.Card = payment.Card
	return shift, nil
}

func (r *TransactionDB) GetAllShifts(filter *model.Filter) ([]*model.Shift, int64, error) {
	shifts := []*model.Shift{}
	if filter.Sort == "" {
		filter.Sort = "created_at DESC"
	}
	res := r.gormDB.Model(&model.Shift{}).Select("shifts.id, shifts.schet_id, shifts.shop_id, shops.name as shop, shifts.created_at, shifts.closed_at, shifts.begin_sum, shifts.end_sum_fact, shifts.end_sum_plan, shifts.expense, shifts.cash, shifts.card, shifts.collection, shifts.difference, shifts.is_closed").Joins("inner join shops on shops.id = shifts.shop_id")
	newRes, count, err := filter.FilterResults(res, shifts, utils.DefaultPageSize, "created_at", "", "")
	if err != nil {
		return nil, 0, err
	}
	err = newRes.Scan(&shifts).Error

	if err != nil {
		return nil, 0, err
	}

	for _, shift := range shifts {
		previousShift, err := r.GetPreviousShift(shift)
		if err != nil {
			return nil, 0, err
		}
		if previousShift.EndSumFact != shift.BeginSum {
			shift.EqualCash = false
			shift.DifferenceWithPrevious = shift.BeginSum - previousShift.EndSumFact
		} else {
			shift.EqualCash = true
		}
		if !shift.IsClosed {
			payment := model.Payment{}
			err := r.gormDB.Table("checks").Select("SUM(checks.cash) AS cash, SUM(checks.card) AS card").Where("checks.status = ? AND checks.closed_at::date = ?::date and checks.shop_id = ?", utils.StatusClosed, shift.CreatedAt, shift.ShopID).Scan(&payment).Error
			if err != nil {
				return nil, 0, err
			}
			shift.Cash = payment.Cash
			shift.Card = payment.Card
		}
		shift.Expense = 0
		shift.EndSumPlan = 0
		shift.Collection = 0
		transactions := []*model.Transaction{}
		err = r.gormDB.Model(&model.Transaction{}).Where("shift_id = ? and deleted = ?", shift.ID, false).Order("id DESC").Scan(&transactions).Error
		if err != nil {
			return nil, 0, err
		}
		newTransactions := []*model.Transaction{}

		for _, transaction := range transactions {
			if transaction.SchetID == shift.SchetID {
				newTransactions = append(newTransactions, transaction)
				if transaction.Category == utils.Collection {
					shift.Collection += transaction.Sum
				} else {
					if transaction.Status == utils.TransactionNegativeStatus {
						shift.Expense += transaction.Sum
					} else if transaction.Status == utils.TransactionPositiveStatus {
						shift.Cash += transaction.Sum
					}
				}
			}
		}
		shift.EndSumPlan = shift.BeginSum + shift.Cash - shift.Expense - shift.Collection
		if shift.IsClosed {
			shift.Difference = shift.EndSumFact - shift.EndSumPlan
		}
	}

	return shifts, count, nil
}

func (r *TransactionDB) GetPreviousShift(shift *model.Shift) (*model.Shift, error) {
	previousShift := &model.Shift{}
	err := r.gormDB.Model(&model.Shift{}).Where("created_at < ? and shop_id = ?", shift.CreatedAt, shift.ShopID).Order("created_at desc").First(previousShift).Error
	if err != nil {
		if err != gorm.ErrRecordNotFound {
			return nil, err
		}
	}
	return previousShift, nil
}

func (r *TransactionDB) GetShiftByTime(timeStamp time.Time, shopID int) (*model.Shift, error) {
	shift := &model.Shift{}
	if shopID == 0 {
		err := r.gormDB.Model(&model.Shift{}).Where("created_at <= ?", timeStamp).Order("created_at desc").First(shift).Error
		if err != nil {
			return nil, err
		}
	} else {
		err := r.gormDB.Model(&model.Shift{}).Where("created_at <= ? and shop_id = ?", timeStamp, shopID).Order("created_at desc").First(shift).Error
		if err != nil {
			return nil, err
		}
	}

	if timeStamp.After(shift.ClosedAt) && shift.IsClosed {
		return nil, errors.New("shoud be in range of shift")
	}
	return shift, nil
}

func (r *TransactionDB) GetShiftByID(id int) (*model.ShiftTransaction, error) {
	shift := &model.ShiftTransaction{}
	err := r.gormDB.Model(&model.Shift{}).Where("id = ?", id).Scan(shift).Error
	if err != nil {
		return nil, err
	}
	if !shift.IsClosed {
		payment := model.Payment{}
		err := r.gormDB.Table("checks").Select("SUM(checks.cash) AS cash, SUM(checks.card) AS card").Where("checks.status = ? AND checks.closed_at::date = ?::date and checks.shop_id = ?", utils.StatusClosed, shift.CreatedAt, shift.ShopID).Scan(&payment).Error
		if err != nil {
			return nil, err
		}
		shift.Cash = payment.Cash
		shift.Card = payment.Card
	}
	shift.Expense = 0
	shift.Income = 0
	shift.EndSumPlan = 0
	shift.Collection = 0
	transactions := []*model.TransactionResponse{}
	err = r.gormDB.Model(&model.Transaction{}).Select("transactions.id, transactions.shift_id, transactions.schet_id, transactions.worker_id, worker.name as worker, transactions.updated_worker_id, updatedWorker.name as updated_worker, transactions.category, transactions.status, transactions.time, transactions.updated_time, transactions.sum, transactions.comment, transactions.deleted").Joins("inner join workers worker on worker.id = transactions.worker_id left join workers updatedWorker on updatedWorker.id = transactions.updated_worker_id").Where("shift_id = ? and transactions.deleted = ?", id, false).Order("time DESC").Scan(&transactions).Error
	if err != nil {
		return nil, err
	}
	newTransactions := []*model.TransactionResponse{}
	for _, transaction := range transactions {
		if transaction.SchetID == shift.SchetID {
			newTransactions = append(newTransactions, transaction)
			if transaction.Category == utils.Collection {
				shift.Collection += transaction.Sum
			} else if transaction.Category == utils.Income {
				shift.Income += transaction.Sum
			} else {
				if transaction.Status == utils.TransactionNegativeStatus {
					shift.Expense += transaction.Sum
				} else if transaction.Status == utils.TransactionPositiveStatus {
					shift.Cash += transaction.Sum
				}
			}
		}
	}
	shift.EndSumPlan = shift.BeginSum + shift.Cash - shift.Expense - shift.Collection + shift.Income
	if shift.IsClosed {
		shift.Difference = shift.EndSumFact - shift.EndSumPlan
	}
	shift.Transactions = newTransactions

	return shift, nil
}

func (r *TransactionDB) GetShiftByID2(id int) (*model.Shift, error) {
	shift := &model.Shift{}
	err := r.gormDB.Model(&model.Shift{}).Where("id = ?", id).Scan(shift).Error
	if err != nil {
		return nil, err
	}
	return shift, nil
}

func (r *TransactionDB) GetShiftByShopId(id int, workerID int) (*model.CurrentShift, error) {
	shiftOutput := &model.CurrentShift{}
	err := r.gormDB.Model(&model.Shift{}).Select("shifts.created_at, shifts.id, shifts.is_closed, shifts.shop_id, workers.name as worker, workers.id as worker_id, shops.name as shop_name").Joins("inner join workers on workers.bind_shop = shifts.shop_id inner join shops on shops.id = shifts.shop_id").Where("shifts.shop_id = ? and workers.id = ?", id, workerID).Order("shifts.created_at desc").First(&shiftOutput).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			err = r.gormDB.Model(&model.Worker{}).Select("workers.name as worker, workers.id as worker_id").Where("workers.id = ?", workerID).First(&shiftOutput).Error
			if err != nil {
				return nil, err
			}
			err = r.gormDB.Model(&model.Shop{}).Select("shops.name as shop_name").Where("shops.id = ?", id).First(&shiftOutput).Error
			if err != nil {
				return nil, err
			}
			shiftOutput.IsClosed = true
		} else {
			return nil, err
		}
	}

	return shiftOutput, nil
}

func (r *TransactionDB) GetAllTransaction(filter *model.Filter) ([]*model.TransactionResponse, int64, error) {
	transactions := []*model.TransactionResponse{}
	var count int64
	res := r.gormDB.Model(&model.Transaction{}).Select("transactions.id,transactions.shift_id, transactions.category, transactions.status, transactions.time, transactions.sum,transactions.comment, schets.name as schet, workers.name as worker, transactions.deleted ").Joins("inner join schets on transactions.schet_id = schets.id inner join workers on transactions.worker_id = workers.id inner join shifts on shifts.id = transactions.shift_id inner join shops on shops.id = shifts.shop_id").Where("transactions.deleted = ? and shops.id IN (?)", false, filter.AccessibleShops)
	if res.Error != nil {
		return nil, 0, res.Error
	}

	newRes, count, err := filter.FilterResults(res, &model.Transaction{}, utils.DefaultPageSize, "time", "", "")
	if err != nil {
		return nil, 0, err
	}
	if filter.Search != "" {
		res.Or("workers.name like ? or schets.name like ?", filter.Search, filter.Search)
	}
	if newRes.Scan(&transactions).Error != nil {
		return nil, 0, newRes.Error
	}

	return transactions, count, nil
}
func (r *TransactionDB) GetTransactionsByID(id int) ([]*model.Transaction, error) {
	transactions := []*model.Transaction{}
	err := r.gormDB.Model(&model.Transaction{}).Where("transactions.shift_id = ? and transactions.deleted = ?", id, false).Scan(&transactions).Error

	if err != nil {
		return nil, err
	}

	return transactions, nil
}

func (r *TransactionDB) GetTransactionByID(id int) (*model.TransactionResponse, error) {
	transaction := &model.TransactionResponse{}
	err := r.gormDB.Model(&model.Transaction{}).Select("transactions.id, transactions.shift_id, transactions.schet_id,transactions.category, transactions.status, transactions.time, transactions.sum,transactions.comment, schets.name as schet, workers.name as worker , transactions.deleted").Joins("inner join schets on transactions.schet_id = schets.id inner join workers on transactions.worker_id = workers.id").Where("transactions.id = ?", id).Scan(transaction).Error

	if err != nil {
		return nil, err
	}
	return transaction, nil
}

func (r *TransactionDB) UpdateTransaction(transaction *model.Transaction) error {
	err := r.gormDB.Transaction(func(tx *gorm.DB) error {
		//transaction.Time = time.Date(transaction.Time.Year(), transaction.Time.Month(), transaction.Time.Day(), transaction.Time.Hour(), transaction.Time.Minute()+1, 0, 0, time.UTC)
		oldTransaction := &model.Transaction{}
		err := r.gormDB.Model(&model.Transaction{}).Where("id = ?", transaction.ID).Scan(oldTransaction).Error
		if err != nil {
			return err
		}
		if oldTransaction.Category != transaction.Category {
			return errors.New("category can't be changed")
		}

		shift := &model.Shift{}
		err = r.gormDB.Model(&model.Shift{}).Where("id = ?", oldTransaction.ShiftID).Scan(shift).Error
		if err != nil {
			return err
		}
		if transaction.Category != utils.OpenShift && transaction.Category != utils.CloseShift && shift.ID != 0 {
			if shift.IsClosed {
				if transaction.Time.After(shift.ClosedAt) {
					return errors.New("time is incorrect; should be past time")
				}
			} else {
				if time.Until(transaction.Time).Minutes() > 1 {
					return errors.New("time is incorrect; should be past time")
				}
			}
			if transaction.Time.Before(shift.CreatedAt) {
				return errors.New("time is incorrect; should be in range of shift time")
			}
		}
		oldSum := oldTransaction.Sum
		oldTransaction.Sum = transaction.Sum
		oldTransaction.Time = transaction.Time
		oldTransaction.UpdatedWorkerID = transaction.UpdatedWorkerID
		oldTransaction.UpdatedTime = time.Now()
		oldTransaction.Comment = transaction.Comment
		oldTransaction.SchetID = transaction.SchetID

		err = r.gormDB.Model(&model.Transaction{}).Where("id = ?", transaction.ID).Updates(oldTransaction).Error
		if err != nil {
			return err
		}

		err = r.gormDB.Model(&model.Schet{}).Where("id = ?", transaction.SchetID).Update("start_balance", gorm.Expr("start_balance + ?", transaction.Sum-oldSum)).Error

		if err != nil {
			return err
		}
		if oldTransaction.Category == utils.OpenShift {

			shift.CreatedAt = transaction.Time
			shift.BeginSum = transaction.Sum
			err = r.gormDB.Model(&model.Shift{}).Where("id = ?", oldTransaction.ShiftID).Updates(shift).Error
			if err != nil {
				return err
			}
		}
		if oldTransaction.Category == utils.CloseShift {

			shift.ClosedAt = transaction.Time
			shift.EndSumFact = transaction.Sum
			err = r.gormDB.Model(&model.Shift{}).Where("id = ?", oldTransaction.ShiftID).Updates(shift).Error
			if err != nil {
				return err
			}
		}
		err = r.RecalculateShift(oldTransaction.ShiftID)
		if err != nil {
			return err
		}
		return nil
	})
	return err
}

func (r *TransactionDB) DeleteTransaction(id int) error {
	oldTransaction := &model.Transaction{}
	err := r.gormDB.Model(&model.Transaction{}).Where("id = ?", id).Scan(oldTransaction).Error
	if err != nil {
		return err
	}
	err = r.gormDB.Model(&model.Transaction{}).Debug().Where("id = ?", id).Update("deleted", true).Error
	if err != nil {
		return err
	}
	err = r.gormDB.Model(&model.Schet{}).Where("id = ?", oldTransaction.SchetID).Update("start_balance", gorm.Expr("start_balance - ?", oldTransaction.Sum)).Error
	return err
}

func (r *TransactionDB) CheckForBlockedShop(shiftID int) (bool, error) {
	var shopID int
	err := r.gormDB.Model(&model.Shift{}).Select("shop_id").Where("id = ?", shiftID).Scan(&shopID).Error
	if err != nil {
		return false, err
	}
	shop := &model.Shop{}
	err = r.gormDB.Model(&model.Shop{}).Where("id = ?", shopID).Scan(&shop).Error
	if err != nil {
		return false, err
	}
	if shop.Blocked {
		return true, nil
	}
	return false, nil
}
