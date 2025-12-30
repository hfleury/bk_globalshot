package psql

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/hfleury/bk_globalshot/internal/model"
	"github.com/hfleury/bk_globalshot/pkg/db"
	"github.com/hfleury/bk_globalshot/pkg/repository"
	"github.com/lib/pq"
)

type PostgresUserRepository struct {
	db db.Db
}

func NewPostgresUserRepository(db db.Db) *PostgresUserRepository {
	return &PostgresUserRepository{
		db: db,
	}
}

func (r *PostgresUserRepository) WithTx(tx db.Db) repository.UserRepository {
	return &PostgresUserRepository{db: tx}
}

func (r *PostgresUserRepository) FindByEmail(ctx context.Context, email string) (*model.User, error) {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()
	query := `
        SELECT id, email, password, role, company_id 
        FROM users 
        WHERE email = $1`

	var user model.User
	var companyID sql.NullString

	row := r.db.GetDb().QueryRowContext(ctx, query, email)

	if err := row.Scan(&user.ID, &user.Email, &user.Password, &user.Role, &companyID); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	user.CompanyID = companyID.String

	return &user, nil
}

func (r *PostgresUserRepository) Create(ctx context.Context, user *model.User) error {
	query := `
        INSERT INTO users (id, email, password, role, company_id)
        VALUES ($1, $2, $3, $4, $5)`

	_, err := r.db.GetDb().ExecContext(ctx, query, user.ID, user.Email, user.Password, user.Role, user.CompanyID)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok && pqErr.Code == "23505" {
			return repository.ErrEmailAlreadyExists
		}
		return err
	}
	return nil
}

func (r *PostgresUserRepository) FindAll(ctx context.Context, limit, offset int) ([]*model.User, int64, error) {
	var total int64
	countQuery := `SELECT count(*) FROM users`
	err := r.db.GetDb().QueryRowContext(ctx, countQuery).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count users: %w", err)
	}

	query := `
		SELECT id, email, role, company_id
		FROM users
		ORDER BY email ASC
		LIMIT $1 OFFSET $2
	`
	rows, err := r.db.GetDb().QueryContext(ctx, query, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list users: %w", err)
	}
	defer rows.Close()

	var users []*model.User
	for rows.Next() {
		var u model.User
		var companyID sql.NullString
		if err := rows.Scan(&u.ID, &u.Email, &u.Role, &companyID); err != nil {
			return nil, 0, fmt.Errorf("failed to scan user: %w", err)
		}
		u.CompanyID = companyID.String
		users = append(users, &u)
	}

	return users, total, nil
}

func (r *PostgresUserRepository) FindByID(ctx context.Context, id string) (*model.User, error) {
	query := `
		SELECT id, email, role, company_id
		FROM users
		WHERE id = $1
	`
	var u model.User
	var companyID sql.NullString
	err := r.db.GetDb().QueryRowContext(ctx, query, id).Scan(&u.ID, &u.Email, &u.Role, &companyID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to find user: %w", err)
	}
	u.CompanyID = companyID.String
	return &u, nil
}

func (r *PostgresUserRepository) Update(ctx context.Context, user *model.User) error {
	query := `
		UPDATE users
		SET email = $1, role = $2, company_id = $3
		WHERE id = $4
	`
	// Note: Password update is usually handled separately for security, skipping for basic CRUD update
	// or handled if provided. For now, assuming basic details update.
	_, err := r.db.GetDb().ExecContext(ctx, query, user.Email, user.Role, user.CompanyID, user.ID)
	if err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}
	return nil
}

func (r *PostgresUserRepository) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM users WHERE id = $1`
	_, err := r.db.GetDb().ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}
	return nil
}
