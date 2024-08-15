package store

import (
	"db-test/types"

	"github.com/jackc/pgx/v5"
)

type widgets []widget

type widget struct {
	ID       int    `db:"id"`
	Name     string `db:"name"`
	Metadata string `db:"meta"`
}

func (w widget) toWidget() types.Widget {
	return types.Widget{
		ID:       w.ID,
		Name:     w.Name,
		Metadata: w.Metadata,
	}
}

func (w widgets) toWidgets() []types.Widget {
	var ws []types.Widget
	for _, w := range w {
		ws = append(ws, w.toWidget())
	}
	return ws
}

func fromWidgets(otherWidgets []types.Widget) widgets {
	var ws widgets
	for _, w := range otherWidgets {
		ws = append(ws, fromWidget(w))
	}
	return ws
}

func fromWidget(w types.Widget) widget {
	return widget{
		ID:       w.ID,
		Name:     w.Name,
		Metadata: w.Metadata,
	}
}

// returns a set of insert queries meant to be used in a batch
func (ws widgets) GetInsertQueries() []Query {
	var queries []Query
	for _, w := range ws {
		// idea: use copy from? single insert query?
		queries = append(queries, Query{
			SQL:  "INSERT INTO widget (name, meta) VALUES ($1, $2)",
			Args: []any{w.Name, w.Metadata},
		})
	}
	return queries
}

type Query struct {
	SQL  string
	Args []any
}
