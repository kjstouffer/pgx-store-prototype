package store_test

import (
	"context"
	"db-test/store"
	"db-test/types"
	"fmt"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/suite"
)

type TxStoreSuite struct {
	suite.Suite
	ctx  context.Context
	pool *pgxpool.Pool
	db   *store.TxStore
}

// set up the suite; shared between all tests
func (s *TxStoreSuite) SetupSuite() {
	// create a new db connection
	s.ctx = context.Background()
	var err error
	s.pool, err = pgxpool.New(s.ctx, "postgres://postgres@localhost/test?sslmode=disable&timezone=UTC")
	s.Require().NoError(err)
	s.db = store.NewTxStore(s.pool)
}

func (s *TxStoreSuite) TearDownSuite() {
	s.pool.Close()
}

func (s *TxStoreSuite) SetupTest() {
	_, err := s.pool.Exec(s.ctx, "CREATE TABLE widget (id serial primary key, name text, meta text)")
	s.Require().NoError(err)
}

func (s *TxStoreSuite) TearDownTest() {
	// drop table to isolate tests
	s.pool.Exec(s.ctx, "DROP TABLE widget")
}

func (s *TxStoreSuite) TestExecBatch_Error() {
	// Tests that if one statement has an error, none will be committed
	widgets := []types.Widget{}
	numWidgets := 10
	dupeID := 1
	dupeIDwidget := types.Widget{ID: &dupeID, Name: "dupe", Metadata: "bar"}
	for i := 0; i < numWidgets; i++ {
		widgets = append(widgets, types.Widget{Name: fmt.Sprintf("foo%d", i), Metadata: "bar"})
	}
	widgets = append(widgets, dupeIDwidget)

	queries := s.db.InsertWidgetsLater(s.ctx, widgets)

	err := s.db.ExecBatch(s.ctx, queries)
	s.Require().Error(err)
	widgetsCheck, err := s.db.GetAllWidgets(s.ctx)
	s.Require().NoError(err)
	s.Require().Len(widgetsCheck, 0)
}

func (s *TxStoreSuite) TestExecBatch() {
	widgets := []types.Widget{}
	numWidgets := 10
	for i := 0; i < numWidgets; i++ {
		widgets = append(widgets, types.Widget{Name: fmt.Sprintf("foo%d", i), Metadata: "bar"})
	}

	queries := s.db.InsertWidgetsLater(s.ctx, widgets)
	err := s.db.ExecBatch(s.ctx, queries)
	s.Require().NoError(err)
	widgetsCheck, err := s.db.GetAllWidgets(s.ctx)
	s.Require().NoError(err)
	s.Require().Len(widgetsCheck, numWidgets)
}

func TestTxStoreSuite(t *testing.T) {
	suite.Run(t, new(TxStoreSuite))
}
