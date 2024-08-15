package store_test

import (
	"context"
	"db-test/store"
	"testing"

	"github.com/jackc/pgx/v5"
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
	batch := pgx.Batch{}
	batch.Queue("INSERT INTO widget (name, meta) VALUES ($1, $2)", "foo", "bar")
	batch.Queue("INSERT INTO widget (name, meta) VALUES ($1, $2)", "foo1", "baz")
	batch.Queue("INSERT INTO widget (name1, meta) VALUES ($1, $2)", "foo2", "bax")

	err := s.db.ExecBatch(s.ctx, &batch)
	s.Require().Error(err)
	widgets, err := s.db.GetAllWidgets(s.ctx)
	s.Require().NoError(err)
	s.Require().Len(widgets, 0)
}

func (s *TxStoreSuite) TestExecBatch() {
	batch := pgx.Batch{}
	batch.Queue("INSERT INTO widget (name, meta) VALUES ($1, $2)", "foo", "bar")
	batch.Queue("INSERT INTO widget (name, meta) VALUES ($1, $2)", "foo1", "baz")
	batch.Queue("INSERT INTO widget (name, meta) VALUES ($1, $2)", "foo2", "bax")

	err := s.db.ExecBatch(s.ctx, &batch)
	s.Require().NoError(err)
	widgets, err := s.db.GetAllWidgets(s.ctx)
	s.Require().NoError(err)
	s.Require().Len(widgets, 3)
}

func TestTxStoreSuite(t *testing.T) {
	suite.Run(t, new(TxStoreSuite))
}
