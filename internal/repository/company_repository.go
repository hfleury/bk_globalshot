package repository

import (
	"context"

	"github.com/hfleury/bk_globalshot/internal/model"
	"github.com/hfleury/bk_globalshot/pkg/db"
)

type CompanyRepository interface {
	Create(ctx context.Context, company *model.Company) error
	FindAll(ctx context.Context, limit, offset int) ([]*model.Company, int64, error)
	FindByID(ctx context.Context, id string) (*model.Company, error)
	Update(ctx context.Context, company *model.Company) error
	Delete(ctx context.Context, id string) error
	WithTx(tx db.Db) CompanyRepository
}
