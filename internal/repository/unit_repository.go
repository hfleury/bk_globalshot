package repository

import (
	"context"

	"github.com/hfleury/bk_globalshot/internal/model"
	"github.com/hfleury/bk_globalshot/pkg/db"
)

type UnitRepository interface {
	Create(ctx context.Context, unit *model.Unit) error
	BatchCreate(ctx context.Context, units []*model.Unit) error
	FindAll(ctx context.Context, limit, offset int) ([]*model.Unit, int64, error)
	FindByID(ctx context.Context, id string) (*model.Unit, error)
	Update(ctx context.Context, unit *model.Unit) error
	Delete(ctx context.Context, id string) error
	WithTx(tx db.Db) UnitRepository
}
