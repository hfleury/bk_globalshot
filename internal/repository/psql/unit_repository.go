package psql

import (
	"context"
	"database/sql"

	"github.com/hfleury/bk_globalshot/internal/model"
	"github.com/hfleury/bk_globalshot/internal/repository"
	"github.com/hfleury/bk_globalshot/pkg/db"
)

type unitRepository struct {
	db db.Db
}

func NewUnitRepository(db db.Db) repository.UnitRepository {
	return &unitRepository{db: db}
}

func (r *unitRepository) WithTx(tx db.Db) repository.UnitRepository {
	return &unitRepository{db: tx}
}

func (r *unitRepository) Create(ctx context.Context, unit *model.Unit) error {
	query := `
		INSERT INTO units (id, name, type, site_id, client_id, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`
	_, err := r.db.GetDb().ExecContext(ctx, query, unit.ID, unit.Name, unit.Type, unit.SiteID, unit.ClientID, unit.CreatedAt, unit.UpdatedAt)
	return err
}

func (r *unitRepository) FindAll(ctx context.Context, limit, offset int) ([]*model.Unit, int64, error) {
	// TODO: Add filtering by SiteID if needed
	var total int64
	countQuery := `SELECT count(*) FROM units`
	err := r.db.GetDb().QueryRowContext(ctx, countQuery).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	query := `
		SELECT id, name, type, site_id, client_id, created_at, updated_at
		FROM units
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`
	rows, err := r.db.GetDb().QueryContext(ctx, query, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	units := make([]*model.Unit, 0)
	for rows.Next() {
		var u model.Unit
		if err := rows.Scan(&u.ID, &u.Name, &u.Type, &u.SiteID, &u.ClientID, &u.CreatedAt, &u.UpdatedAt); err != nil {
			return nil, 0, err
		}
		units = append(units, &u)
	}

	return units, total, nil
}

func (r *unitRepository) FindByID(ctx context.Context, id string) (*model.Unit, error) {
	query := `
		SELECT id, name, type, site_id, client_id, created_at, updated_at
		FROM units
		WHERE id = $1
	`
	var u model.Unit
	err := r.db.GetDb().QueryRowContext(ctx, query, id).Scan(&u.ID, &u.Name, &u.Type, &u.SiteID, &u.ClientID, &u.CreatedAt, &u.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // Not found
		}
		return nil, err
	}
	return &u, nil
}

func (r *unitRepository) Update(ctx context.Context, unit *model.Unit) error {
	query := `
		UPDATE units
		SET name = $1, type = $2, site_id = $3, client_id = $4, updated_at = $5
		WHERE id = $6
	`
	_, err := r.db.GetDb().ExecContext(ctx, query, unit.Name, unit.Type, unit.SiteID, unit.ClientID, unit.UpdatedAt, unit.ID)
	return err
}

func (r *unitRepository) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM units WHERE id = $1`
	_, err := r.db.GetDb().ExecContext(ctx, query, id)
	return err
}
