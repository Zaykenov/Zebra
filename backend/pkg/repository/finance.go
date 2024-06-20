package repository

import (
	"database/sql"
	"log"
	"zebra/model"
	"zebra/utils"

	"gorm.io/gorm"
)

type FinanceDB struct {
	db     *sql.DB
	gormDB *gorm.DB
}

func NewFinanceDB(db *sql.DB, gormDB *gorm.DB) *FinanceDB {
	return &FinanceDB{db: db, gormDB: gormDB}
}

func (r *FinanceDB) AddSchet(schet *model.ReqSchet) (*model.Schet, error) {
	/*ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	stmt := `insert into schet (name, currency, type, start_balance)
	values ($1, $2, $3, $4)`
	err := r.db.QueryRowContext(ctx, stmt,
		schet.Name,
		schet.Currency,
		schet.Type,
		schet.StartBalance,
	)

	if err.Err() != nil {
		return err.Err()
	}

	return nil*/
	schetRes := &model.Schet{
		ID:           schet.ID,
		Name:         schet.Name,
		Currency:     schet.Currency,
		Type:         schet.Type,
		Deleted:      schet.Deleted,
		StartBalance: schet.StartBalance,
	}
	err := r.gormDB.Create(schetRes).Error
	if err != nil {
		return nil, err
	}
	shops := []*model.Shop{}
	err = r.gormDB.Model(shops).Where("id IN ?", schet.ShopIDs).Scan(&shops).Error
	if err != nil {
		return nil, err
	}
	for _, shop := range shops {
		shopSchet := &model.ShopSchet{
			ShopID:  shop.ID,
			SchetID: schetRes.ID,
		}
		err = r.gormDB.Table("shop_schets").Create(shopSchet).Error
		if err != nil {
			return nil, err
		}
	}
	return schetRes, nil
}

func (r *FinanceDB) GetAllSchet(filter *model.Filter) ([]*model.Schet, int64, error) {
	/*
		ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
		defer cancel()
		schets := []*model.Schet{}
		stmt := `select * from schet where deleted = $1`

		row, err := r.db.QueryContext(ctx, stmt, false)
		if err != nil {
			return nil, err
		}

		for row.Next() {
			schet := &model.Schet{}
			err := row.Scan(
				&schet.ID,
				&schet.Name,
				&schet.Currency,
				&schet.Type,
				&schet.StartBalance,
				&schet.Deleted,
			)
			if err != nil {
				return nil, err
			}
			schets = append(schets, schet)
		}

		return schets, nil*/

	schetIDs := []int{}
	availableShops := []*model.Shop{}
	err := r.gormDB.Model(availableShops).Where("id IN (?)", filter.AccessibleShops).Find(&availableShops).Error
	if err != nil {
		return nil, 0, err
	}
	for _, shop := range availableShops {
		schetIDs = append(schetIDs, shop.CardSchet)
		schetIDs = append(schetIDs, shop.CashSchet)
	}
	shopSchets := []*model.ShopSchet{}
	err = r.gormDB.Debug().Table("shop_schets").Where("shop_id IN (?)", filter.AccessibleShops).Find(&shopSchets).Error
	if err != nil {
		return nil, 0, err
	}
	for _, shopSchet := range shopSchets {
		schetIDs = append(schetIDs, shopSchet.SchetID)
	}
	schets := []*model.Schet{}
	res := r.gormDB
	if filter.Role == utils.MasterRole {
		res = r.gormDB.Model(schets).Debug().Where("deleted = ?", false).Scan(&schets)
		if res.Error != nil {
			return nil, 0, res.Error
		}
	} else {
		res = r.gormDB.Debug().Where("deleted = ? and id IN (?)", false, schetIDs).Find(&schets)
		if res.Error != nil {
			return nil, 0, res.Error
		}
	}
	log.Print("schets", schetIDs)
	newRes, count, err := filter.FilterResults(res, model.Schet{}, utils.DefaultPageSize, "", "", "")

	if err != nil {
		return nil, 0, err
	}

	if newRes.Scan(&schets).Error != nil {
		return nil, 0, newRes.Error
	}

	return schets, count, nil

}

func (r *FinanceDB) GetSchetByID(id int) (*model.Schet, error) {
	/*ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	schet := &model.Schet{}
	stmt := `select * from schet where id = $1 and deleted = $2`

	err := r.db.QueryRowContext(ctx, stmt, id, false).Scan(
		&schet.ID,
		&schet.Name,
		&schet.Currency,
		&schet.Type,
		&schet.StartBalance,
		&schet.Deleted,
	)

	if err != nil {
		return nil, err
	}

	return schet, nil*/
	schet := &model.Schet{}
	err := r.gormDB.Where("id = ? and deleted = ?", id, false).First(schet).Error
	if err != nil {
		return nil, err
	}
	return schet, nil
}

func (r *FinanceDB) UpdateSchet(schet *model.Schet) error {
	/*ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	stmt := `update schet set name = $1, currency = $2, type = $3, start_balance = $4 where id = $5`
	_, err := r.db.ExecContext(ctx, stmt,
		schet.Name,
		schet.Currency,
		schet.Type,
		schet.StartBalance,
		schet.ID,
	)

	if err != nil {
		return err
	}

	return nil*/
	return r.gormDB.Save(schet).Error
}

func (r *FinanceDB) DeleteSchet(id int) error {
	/*ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	stmt := `update schet set deleted = $1 where id = $2`
	_, err := r.db.ExecContext(ctx, stmt, true, id)

	if err != nil {
		return err
	}

	return nil*/
	if err := r.gormDB.Table("schets").Where("id = ?", id).Update("deleted", true).Error; err != nil {
		return err
	}
	return nil
}
