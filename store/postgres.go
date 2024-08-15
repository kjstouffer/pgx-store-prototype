package store

import (
	"db-test/types"

	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type TxStore struct {
	db *pgxpool.Pool
}

func NewTxStore(db *pgxpool.Pool) *TxStore {
	return &TxStore{db}
}

func (ts *TxStore) GetWidgetByID(ctx context.Context, id int) (types.Widget, error) {
	var w widget
	err := pgx.BeginTxFunc(ctx, ts.db, pgx.TxOptions{AccessMode: pgx.ReadOnly}, func(tx pgx.Tx) error {
		rows, err := tx.Query(ctx, "select id, name, meta from widgets where id = $1", id)
		if err != nil {
			return err
		}
		defer rows.Close()
		w, err = pgx.CollectExactlyOneRow[widget](rows, pgx.RowToStructByName)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return types.Widget{}, err
	}

	return w.toWidget(), nil
}

type later interface {
}

func (ts *TxStore) ExecBatch(ctx context.Context, batch *pgx.Batch) error {
	err := pgx.BeginTxFunc(ctx, ts.db, pgx.TxOptions{AccessMode: pgx.ReadWrite}, func(tx pgx.Tx) error {
		for _, q := range batch.QueuedQueries {
			q.Exec(func(ct c) error {})
			_, err := tx.Exec(ctx, q.SQL, q.Arguments...)
			if err != nil {
				return err
			}
		}
	})
	return err
}
