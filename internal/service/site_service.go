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

type SiteService interface {
	CreateSite(ctx context.Context, name, address, companyID string) (*model.Site, error)
	GetAllSites(ctx context.Context, limit, offset int) ([]*model.Site, int64, error)
	GetSiteByID(ctx context.Context, id string) (*model.Site, error)
	UpdateSite(ctx context.Context, id, name, address string) (*model.Site, error)
	DeleteSite(ctx context.Context, id string) error
}

type siteService struct {
	db   db.Db
	repo repository.SiteRepository
}

func NewSiteService(db db.Db, repo repository.SiteRepository) SiteService {
	return &siteService{
		db:   db,
		repo: repo,
	}
}

func (s *siteService) CreateSite(ctx context.Context, name, address, companyID string) (*model.Site, error) {
	site := &model.Site{
		ID:        uuid.New().String(),
		Name:      name,
		Address:   address,
		CompanyID: companyID,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := s.repo.Create(ctx, site); err != nil {
		return nil, fmt.Errorf("failed to create site: %w", err)
	}
	return site, nil
}

func (s *siteService) GetAllSites(ctx context.Context, limit, offset int) ([]*model.Site, int64, error) {
	user := ctx.Value("user").(*model.User)
	if user == nil {
		return nil, 0, fmt.Errorf("user not found in context")
	}

	if user.Role == string(model.RoleCustomer) {
		return s.repo.FindAllByCustomerID(ctx, limit, offset, user.ID)
	}

	// For Company and Admin, we use CompanyID.
	// Admin might see all? Assuming Admin belongs to a specific company or we default to user.CompanyID
	return s.repo.FindAllByCompanyID(ctx, limit, offset, user.CompanyID)
}

func (s *siteService) GetSiteByID(ctx context.Context, id string) (*model.Site, error) {
	return s.repo.FindByID(ctx, id)
}

func (s *siteService) UpdateSite(ctx context.Context, id, name, address string) (*model.Site, error) {
	site, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if site == nil {
		return nil, nil // Not found
	}

	site.Name = name
	site.Address = address
	site.UpdatedAt = time.Now()

	if err := s.repo.Update(ctx, site); err != nil {
		return nil, err
	}
	return site, nil
}

func (s *siteService) DeleteSite(ctx context.Context, id string) error {
	site, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return err
	}
	if site == nil {
		return nil
	}
	return s.repo.Delete(ctx, id)
}
