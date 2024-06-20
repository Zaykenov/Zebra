package repository

import (
	"context"
	"database/sql"

	"github.com/bxcodec/faker/v3"
)

type SeederDB struct {
	db *sql.DB
}

func NewSeederDB(db *sql.DB) *SeederDB {
	return &SeederDB{db: db}
}

func (r *SeederDB) SeedCategory() error {
	for i := 0; i < 15; i++ {
		stmt := `INSERT INTO category_tovar(name)
	values ($1)`
		row := r.db.QueryRowContext(context.TODO(), stmt, faker.FirstName())
		if row.Err() != nil {
			return row.Err()
		}
		stmt = `INSERT INTO category_ingredient(name)
	values ($1)`
		row = r.db.QueryRowContext(context.TODO(), stmt, faker.FirstName())
		if row.Err() != nil {
			return row.Err()
		}
	}
	return nil
}
