package db

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
)

type PsqlDb struct {
	db *sql.DB
}

func NewPsqlDb(dsn string) (*PsqlDb, error) {
	var db *sql.DB
	var err, pingErr error

	for i := 0; i < 10; i++ {
		db, err = sql.Open("pgx", dsn)
		if err == nil && db != nil {
			pingErr = db.Ping()
			if pingErr == nil {
				return &PsqlDb{db: db}, nil
			}
		}

		log.Printf("waiting for DB connection... (%v)", pingErr)
		time.Sleep(5 * time.Second)
	}

	return nil, fmt.Errorf("failed to connect to DB after retries")
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

func (p *PsqlDb) PingContext(ctx context.Context) error {
	return p.db.Ping()
}
