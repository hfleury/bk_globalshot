package db

import (
	"context"
	"database/sql"
)

//go:generate mockgen -source=db.go -destination=../../mock/db/mock_db.go -package=mock_db
type DbTx interface {
	ExecContext(ctx context.Context, query string, arg ...interface{}) (sql.Result, error)
	QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error)
	QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row
}

type SqlDb interface {
	DbTx
	BeginTx(ctx context.Context, opts *sql.TxOptions) (*sql.Tx, error)
}

type Db interface {
	GetDb() DbTx
	BegrinTransaction(ctx context.Context) (DbTx, error)
	Commit(ctx context.Context, tx DbTx) error
	Rollback(ctx context.Context, tx DbTx) error
}
