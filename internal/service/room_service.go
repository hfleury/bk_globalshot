package service

import (
	"context"
	"time"

	"github.com/hfleury/bk_globalshot/internal/model"
	"github.com/hfleury/bk_globalshot/internal/repository"
)

//go:generate mockgen -source=room_service.go -destination=../../mock/services/mock_room_service.go -package=mock_services
type RoomService interface {
	CreateRoom(ctx context.Context, name, unitID string) (*model.Room, error)
	GetAllRooms(ctx context.Context, limit, offset int) ([]*model.Room, int64, error)
	GetRoomByID(ctx context.Context, id string) (*model.Room, error)
	UpdateRoom(ctx context.Context, id string, name string, unitID string) (*model.Room, error)
	DeleteRoom(ctx context.Context, id string) error
}

type roomService struct {
	repo repository.RoomRepository
}

func NewRoomService(repo repository.RoomRepository) RoomService {
	return &roomService{
		repo: repo,
	}
}

func (s *roomService) CreateRoom(ctx context.Context, name, unitID string) (*model.Room, error) {
	room := &model.Room{
		Name:      name,
		UnitID:    unitID,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	return room, s.repo.Create(ctx, room)
}

func (s *roomService) GetAllRooms(ctx context.Context, limit, offset int) ([]*model.Room, int64, error) {
	return s.repo.FindAll(ctx, limit, offset)
}

func (s *roomService) GetRoomByID(ctx context.Context, id string) (*model.Room, error) {
	return s.repo.FindByID(ctx, id)
}

func (s *roomService) UpdateRoom(ctx context.Context, id string, name string, unitID string) (*model.Room, error) {
	room, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if room == nil {
		return nil, nil
	}

	room.Name = name
	room.UnitID = unitID
	room.UpdatedAt = time.Now()

	if err := s.repo.Update(ctx, room); err != nil {
		return nil, err
	}
	return room, nil
}

func (s *roomService) DeleteRoom(ctx context.Context, id string) error {
	room, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return err
	}
	if room == nil {
		return nil
	}
	return s.repo.Delete(ctx, id)
}
