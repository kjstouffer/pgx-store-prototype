package store

import (
	"db-test/types"

	"github.com/jackc/pgx/v5"
)

type widget struct {
	ID       int    `db:"id"`
	Name     string `db:"name"`
	Metadata string `db:"metadata"`
}

func (w widget) toWidget() types.Widget {
	return types.Widget{
		ID:       w.ID,
		Name:     w.Name,
		Metadata: w.Metadata,
	}
}

func (w widget) Insert() pgx.QueuedQuery
