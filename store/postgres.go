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
		rows, err := tx.Query(ctx, "select id, name, meta from widget where id = $1", id)
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

func (ts *TxStore) GetAllWidgets(ctx context.Context) ([]types.Widget, error) {
	var w widgets
	err := pgx.BeginTxFunc(ctx, ts.db, pgx.TxOptions{AccessMode: pgx.ReadOnly}, func(tx pgx.Tx) error {
		rows, err := tx.Query(ctx, "select id, name, meta from widget")
		if err != nil {
			return err
		}
		defer rows.Close()
		w, err = pgx.CollectRows[widget](rows, pgx.RowToStructByName)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return w.toWidgets(), nil
}

func (ts *TxStore) ExecBatch(ctx context.Context, qs []Query) error {
	batch := &pgx.Batch{}
	for _, q := range qs {
		batch.Queue(q.SQL, q.Args...)
	}
	err := pgx.BeginTxFunc(ctx, ts.db, pgx.TxOptions{AccessMode: pgx.ReadWrite}, func(tx pgx.Tx) error {
		br := tx.SendBatch(ctx, batch)
		return br.Close()
	})
	return err
}

// InsterWidgetLater returns a queued query to batch with other queries. Meant to be used with `TxStore.ExecBatch`.
func (ts *TxStore) InsertWidgetsLater(ctx context.Context, widgets []types.Widget) []Query {
	return fromWidgets(widgets).GetInsertQueries()
}
