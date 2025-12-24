package repository

import (
	"context"

	"github.com/hfleury/bk_globalshot/internal/model"
)

type RoomRepository interface {
	Create(ctx context.Context, room *model.Room) error
	FindAll(ctx context.Context, limit, offset int) ([]*model.Room, int64, error)
	FindByID(ctx context.Context, id string) (*model.Room, error)
	Update(ctx context.Context, room *model.Room) error
	Delete(ctx context.Context, id string) error
	// Additional filters can be added to FindAll later
}
