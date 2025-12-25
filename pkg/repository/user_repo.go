package repository

import (
	"context"

	"github.com/hfleury/bk_globalshot/internal/model"
	"github.com/hfleury/bk_globalshot/pkg/db"
)

//go:generate mockgen -source=user_repo.go -destination=../../mock/repository/mock_user_repo.go -package=mock_repository

type UserRepository interface {
	FindByEmail(ctx context.Context, email string) (*model.User, error)
	Create(ctx context.Context, user *model.User) error
	WithTx(tx db.Db) UserRepository
}
