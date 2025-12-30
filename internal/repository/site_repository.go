package repository

import (
	"context"

	"github.com/hfleury/bk_globalshot/internal/model"
	"github.com/hfleury/bk_globalshot/pkg/db"
)

type SiteRepository interface {
	Create(ctx context.Context, site *model.Site) error
	FindAllByCompanyID(ctx context.Context, limit, offset int, companyID string) ([]*model.Site, int64, error)
	FindAllByCustomerID(ctx context.Context, limit, offset int, customerID string) ([]*model.Site, int64, error)
	FindByID(ctx context.Context, id string) (*model.Site, error)
	Update(ctx context.Context, site *model.Site) error
	Delete(ctx context.Context, id string) error
	WithTx(tx db.Db) SiteRepository
}
