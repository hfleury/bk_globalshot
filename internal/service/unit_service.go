package service

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/hfleury/bk_globalshot/internal/model"
	"github.com/hfleury/bk_globalshot/internal/repository"
	"github.com/hfleury/bk_globalshot/pkg/db"
)

type BatchCreateUnitItem struct {
	Name     string  `json:"name"`
	Type     string  `json:"type"`
	SiteID   string  `json:"site_id"`
	ClientID *string `json:"client_id"`
}

type UnitService interface {
	CreateUnit(ctx context.Context, name string, unitType string, siteID string, clientID *string) (*model.Unit, error)
	BatchCreateUnits(ctx context.Context, items []BatchCreateUnitItem) ([]*model.Unit, error)
	GetAllUnits(ctx context.Context, limit, offset int) ([]*model.Unit, int64, error)
	GetUnitByID(ctx context.Context, id string) (*model.Unit, error)
	UpdateUnit(ctx context.Context, id, name, unitType, siteID string, clientID *string) (*model.Unit, error)
	DeleteUnit(ctx context.Context, id string) error
}

type unitService struct {
	db   db.Db
	repo repository.UnitRepository
}

func NewUnitService(db db.Db, repo repository.UnitRepository) UnitService {
	return &unitService{
		db:   db,
		repo: repo,
	}
}

func (s *unitService) CreateUnit(ctx context.Context, name string, unitType string, siteID string, clientID *string) (*model.Unit, error) {
	unit := &model.Unit{
		ID:        uuid.New().String(),
		Name:      name,
		Type:      model.UnitType(unitType),
		SiteID:    siteID,
		ClientID:  clientID,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := s.repo.Create(ctx, unit); err != nil {
		return nil, fmt.Errorf("failed to create unit: %w", err)
	}
	return unit, nil
}

func (s *unitService) BatchCreateUnits(ctx context.Context, items []BatchCreateUnitItem) ([]*model.Unit, error) {
	if len(items) == 0 {
		return []*model.Unit{}, nil
	}

	units := make([]*model.Unit, len(items))
	now := time.Now()

	for i, item := range items {
		units[i] = &model.Unit{
			ID:        uuid.New().String(),
			Name:      item.Name,
			Type:      model.UnitType(item.Type),
			SiteID:    item.SiteID,
			ClientID:  item.ClientID,
			CreatedAt: now,
			UpdatedAt: now,
		}
	}

	if err := s.repo.BatchCreate(ctx, units); err != nil {
		return nil, fmt.Errorf("failed to batch create units: %w", err)
	}

	return units, nil
}

func (s *unitService) GetAllUnits(ctx context.Context, limit, offset int) ([]*model.Unit, int64, error) {
	return s.repo.FindAll(ctx, limit, offset)
}

func (s *unitService) GetUnitByID(ctx context.Context, id string) (*model.Unit, error) {
	return s.repo.FindByID(ctx, id)
}

func (s *unitService) UpdateUnit(ctx context.Context, id, name, unitType, siteID string, clientID *string) (*model.Unit, error) {
	unit, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if unit == nil {
		return nil, nil // Not found
	}

	unit.Name = name
	unit.Type = model.UnitType(unitType)
	unit.SiteID = siteID
	unit.ClientID = clientID
	unit.UpdatedAt = time.Now()

	if err := s.repo.Update(ctx, unit); err != nil {
		return nil, err
	}
	return unit, nil
}

func (s *unitService) DeleteUnit(ctx context.Context, id string) error {
	unit, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return err
	}
	if unit == nil {
		return nil
	}
	return s.repo.Delete(ctx, id)
}
