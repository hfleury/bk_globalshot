package psql

import (
	"context"
	"database/sql"
	"time"

	"github.com/hfleury/bk_globalshot/internal/model"
	"github.com/hfleury/bk_globalshot/pkg/db"
)

type PostgresUserRepository struct {
	db db.Db
}

func NewPostgresUserRepository(db db.Db) *PostgresUserRepository {
	return &PostgresUserRepository{
		db: db,
	}
}

func (r *PostgresUserRepository) FindByEmail(ctx context.Context, email string) (*model.User, error) {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()
	query := `
        SELECT id, email, password 
        FROM users 
        WHERE email = $1`

	var user model.User
	row := r.db.GetDb().QueryRowContext(ctx, query, email)

	if err := row.Scan(&user.ID, &user.Email, &user.Password); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &user, nil
}

func (r *PostgresUserRepository) Create(ctx context.Context, user *model.User) error {
	query := `
        INSERT INTO users (id, email, password)
        VALUES ($1, $2, $3)`

	_, err := r.db.GetDb().ExecContext(ctx, query, user.ID, user.Email, user.Password)
	return err
}
