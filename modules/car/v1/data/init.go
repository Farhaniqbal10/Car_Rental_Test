package data

import (
	"car_rental_test/pkg/errors"
	"context"
	"log"
	"slices"

	"github.com/jmoiron/sqlx"
)

type (
	Data struct {
		db   *sqlx.DB
		stmt map[string]*sqlx.Stmt
	}

	statement struct {
		key   string
		query string
	}
)

func New(db *sqlx.DB) *Data {
	d := &Data{
		db: db,
	}

	d.initStmt()

	return d
}

func (d *Data) initStmt() {
	var (
		err     error
		stmtMap = make(map[string]*sqlx.Stmt)
	)

	stmts := slices.Concat(
		BookingStmts,
		CarsStmts,
	)

	for _, v := range stmts {
		stmtMap[v.key], err = d.db.Preparex(v.query)
		if err != nil {
			log.Fatalf("[DB] Failed to initialize statement key %v, err : %v", v.key, err)
		}
	}

	d.stmt = stmtMap
}

func (d *Data) BeginTx(ctx context.Context) (*sqlx.Tx, error) {
	tx, err := d.db.BeginTxx(ctx, nil)
	if err != nil {
		return nil, errors.Wrap(err, "[DATA][BeginTx]")
	}

	return tx, nil
}

func (d *Data) CommitTx(ctx context.Context, tx *sqlx.Tx) error {
	err := tx.Commit()
	if err != nil {
		return errors.Wrap(err, "[DATA][CommitTx]")
	}

	return nil
}

func (d *Data) RollbackTx(ctx context.Context, tx *sqlx.Tx) error {
	err := tx.Rollback()
	if err != nil {
		return errors.Wrap(err, "[DATA][RollbackTx]")
	}

	return nil
}
