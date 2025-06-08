package db

import (
	"context"
	"database/sql"
	"log"

	_ "github.com/jackc/pgx/v5/stdlib"
)

type PsqlDb struct {
	db *sql.DB
}

func NewPsqlDb(dsn string) (*PsqlDb, error) {
	log.Printf("Connecting to DB with DSN: %s", dsn)
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	log.Println("Connected to PostegreSQL")
	return &PsqlDb{db: db}, nil
}

func (p *PsqlDb) GetDb() DbTx {
	return p.db
}

func (p *PsqlDb) BegrinTransaction(ctx context.Context) (DbTx, error) {
	return p.db.BeginTx(ctx, nil)
}

func (p *PsqlDb) Commit(ctx context.Context, tx DbTx) error {
	return tx.(*sql.Tx).Commit()
}

func (p *PsqlDb) Rollback(ctx context.Context, tx DbTx) error {
	return tx.(*sql.Tx).Rollback()
}
