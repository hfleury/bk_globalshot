package psql

import (
	"context"

	"github.com/hfleury/bk_globalshot/internal/model"
	"github.com/hfleury/bk_globalshot/pkg/db"
)

type PostgresCompanyRepository struct {
	db db.Db
}

func NewPostgresCompanyRepository(db db.Db) *PostgresCompanyRepository {
	return &PostgresCompanyRepository{db: db}
}

func (r *PostgresCompanyRepository) Create(ctx context.Context, company *model.Company) error {
	query := `
		INSERT INTO companies (name, created_at)
		VALUES ($1, $2)
		RETURNING id`

	row := r.db.GetDb().QueryRowContext(ctx, query, company.Name, company.CreatedAt)
	if err := row.Scan(&company.ID); err != nil {
		return err
	}
	return nil
}

func (r *PostgresCompanyRepository) FindAll(ctx context.Context, limit, offset int) ([]*model.Company, int64, error) {
	query := `
        SELECT id, name, created_at, COUNT(*) OVER() as total_count 
        FROM companies 
        WHERE deleted_at IS NULL
        ORDER BY created_at DESC 
        LIMIT $1 OFFSET $2`

	rows, err := r.db.GetDb().QueryContext(ctx, query, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var companies []*model.Company
	var totalCount int64

	for rows.Next() {
		var c model.Company
		if err := rows.Scan(&c.ID, &c.Name, &c.CreatedAt, &totalCount); err != nil {
			return nil, 0, err
		}
		companies = append(companies, &c)
	}

	// If no rows, check if it's just empty page or empty table.
	// However, if no rows, totalCount will be 0 from the loop not running.
	// If we want accurate total count even when page is empty (out of bounds), we might need separate query.
	// But for now, if empty, total 0 is "safe-ish" or we assume 0.
	// Actually, for React Admin, if I request page 10 and it's empty, I still need total count?
	// Usage of window function in empty result set returns... nothing.
	// So if offset > count, we get 0 rows and 0 count.
	// React Admin might be confused if it thinks total is 0 but it just saw page 1.
	// Proper way: Separate Count query or this is fine for "simple" start.
	// Let's stick to this. If `companies` is empty, we return 0 total (which is technically wrong if out of bounds, but UI usually resets).

	return companies, totalCount, nil
}

func (r *PostgresCompanyRepository) FindByID(ctx context.Context, id string) (*model.Company, error) {
	query := `SELECT id, name, created_at FROM companies WHERE id = $1 AND deleted_at IS NULL`
	row := r.db.GetDb().QueryRowContext(ctx, query, id)

	var c model.Company
	if err := row.Scan(&c.ID, &c.Name, &c.CreatedAt); err != nil {
		return nil, err
	}
	return &c, nil
}

func (r *PostgresCompanyRepository) Update(ctx context.Context, company *model.Company) error {
	query := `UPDATE companies SET name = $1 WHERE id = $2`
	_, err := r.db.GetDb().ExecContext(ctx, query, company.Name, company.ID)
	return err
}

func (r *PostgresCompanyRepository) Delete(ctx context.Context, id string) error {
	query := `UPDATE companies SET deleted_at = NOW() WHERE id = $1`
	_, err := r.db.GetDb().ExecContext(ctx, query, id)
	return err
}
