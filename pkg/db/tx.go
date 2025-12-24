package db

import (
	"context"
	"database/sql"
	"errors"
)

// TxAdapter adapts an *sql.Tx to the Db interface.
// This allows repositories to run on a transaction without knowing it.
type TxAdapter struct {
	tx *sql.Tx
}

func NewTxAdapter(tx *sql.Tx) *TxAdapter {
	return &TxAdapter{tx: tx}
}

func (t *TxAdapter) GetDb() DbTx {
	return t.tx
}

func (t *TxAdapter) BegrinTransaction(ctx context.Context) (DbTx, error) {
	return nil, errors.New("cannot start a transaction within a transaction")
}

func (t *TxAdapter) Commit(ctx context.Context, tx DbTx) error {
	return errors.New("nested commits are not supported, use the parent transaction manager")
}

func (t *TxAdapter) Rollback(ctx context.Context, tx DbTx) error {
	return errors.New("nested rollbacks are not supported, use the parent transaction manager")
}

func (t *TxAdapter) PingContext(ctx context.Context) error {
	return nil // Transaction is active, so connection is alive
}
