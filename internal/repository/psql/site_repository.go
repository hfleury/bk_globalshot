package psql

import (
	"context"
	"database/sql"

	"github.com/hfleury/bk_globalshot/internal/model"
	"github.com/hfleury/bk_globalshot/internal/repository"
	"github.com/hfleury/bk_globalshot/pkg/db"
)

type siteRepository struct {
	db db.Db
}

func NewSiteRepository(db db.Db) repository.SiteRepository {
	return &siteRepository{db: db}
}

func (r *siteRepository) WithTx(tx db.Db) repository.SiteRepository {
	return &siteRepository{db: tx}
}

func (r *siteRepository) Create(ctx context.Context, site *model.Site) error {
	query := `
		INSERT INTO construction_sites (id, name, address, company_id, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6)
	`
	_, err := r.db.GetDb().ExecContext(ctx, query, site.ID, site.Name, site.Address, site.CompanyID, site.CreatedAt, site.UpdatedAt)
	return err
}

func (r *siteRepository) FindAllByCompanyID(ctx context.Context, limit, offset int, companyID string) ([]*model.Site, int64, error) {
	var total int64
	countQuery := `SELECT count(*) FROM construction_sites WHERE company_id = $1`
	err := r.db.GetDb().QueryRowContext(ctx, countQuery, companyID).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	query := `
		SELECT id, name, address, company_id, created_at, updated_at
		FROM construction_sites
		WHERE company_id = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`
	rows, err := r.db.GetDb().QueryContext(ctx, query, companyID, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var sites []*model.Site
	for rows.Next() {
		var s model.Site
		if err := rows.Scan(&s.ID, &s.Name, &s.Address, &s.CompanyID, &s.CreatedAt, &s.UpdatedAt); err != nil {
			return nil, 0, err
		}
		sites = append(sites, &s)
	}

	return sites, total, nil
}

func (r *siteRepository) FindAllByCustomerID(ctx context.Context, limit, offset int, customerID string) ([]*model.Site, int64, error) {
	// Count distinct sites connected to units assigned to this customer
	countQuery := `
		SELECT count(DISTINCT s.id)
		FROM construction_sites s
		JOIN units u ON u.site_id = s.id
		WHERE u.client_id = $1
	`
	var total int64
	err := r.db.GetDb().QueryRowContext(ctx, countQuery, customerID).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	query := `
		SELECT DISTINCT s.id, s.name, s.address, s.company_id, s.created_at, s.updated_at
		FROM construction_sites s
		JOIN units u ON u.site_id = s.id
		WHERE u.client_id = $1
		ORDER BY s.created_at DESC
		LIMIT $2 OFFSET $3
	`
	rows, err := r.db.GetDb().QueryContext(ctx, query, customerID, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var sites []*model.Site
	for rows.Next() {
		var s model.Site
		if err := rows.Scan(&s.ID, &s.Name, &s.Address, &s.CompanyID, &s.CreatedAt, &s.UpdatedAt); err != nil {
			return nil, 0, err
		}
		sites = append(sites, &s)
	}

	return sites, total, nil
}

func (r *siteRepository) FindByID(ctx context.Context, id string) (*model.Site, error) {
	query := `
		SELECT id, name, address, company_id, created_at, updated_at
		FROM construction_sites
		WHERE id = $1
	`
	var s model.Site
	err := r.db.GetDb().QueryRowContext(ctx, query, id).Scan(&s.ID, &s.Name, &s.Address, &s.CompanyID, &s.CreatedAt, &s.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // Not found
		}
		return nil, err
	}
	return &s, nil
}

func (r *siteRepository) Update(ctx context.Context, site *model.Site) error {
	query := `
		UPDATE construction_sites
		SET name = $1, address = $2, updated_at = $3
		WHERE id = $4
	`
	_, err := r.db.GetDb().ExecContext(ctx, query, site.Name, site.Address, site.UpdatedAt, site.ID)
	return err
}

func (r *siteRepository) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM construction_sites WHERE id = $1`
	_, err := r.db.GetDb().ExecContext(ctx, query, id)
	return err
}
