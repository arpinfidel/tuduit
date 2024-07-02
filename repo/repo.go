package repo

import (
	"context"

	"github.com/arpinfidel/tuduit/pkg/db"
	"github.com/arpinfidel/tuduit/pkg/trace"
	"github.com/jmoiron/sqlx"
)

type DBConnection struct {
	db *db.DB
}

func NewDBConnection(db *db.DB) *DBConnection {
	return &DBConnection{db: db}
}

func (r *DBConnection) StartTx(ctx context.Context, f func(ctx context.Context, tx *sqlx.Tx, data any) (any, error)) (wf func(ctx context.Context, data any) (any, error), commit, rollback func() error, err error) {
	defer trace.Default(&ctx, &err)()

	dbTx, err := r.db.GetMaster().BeginTxx(ctx, nil)
	if err != nil {
		return nil, nil, nil, err
	}

	return func(ctx context.Context, data any) (any, error) {
		return f(ctx, dbTx, data)
	}, dbTx.Commit, dbTx.Rollback, nil
}
