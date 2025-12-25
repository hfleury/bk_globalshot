package service

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/hfleury/bk_globalshot/internal/model"
	"github.com/hfleury/bk_globalshot/internal/repository"
	"github.com/hfleury/bk_globalshot/pkg/db"
	pkgRepository "github.com/hfleury/bk_globalshot/pkg/repository"
	"golang.org/x/crypto/bcrypt"
)

//go:generate mockgen -source=company_service.go -destination=../../mock/services/mock_company_service.go -package=mock_services
type CompanyService interface {
	CreateCompany(ctx context.Context, name, email, password string) (*model.Company, error)
	GetAllCompanies(ctx context.Context, limit, offset int) ([]*model.Company, int64, error)
	GetCompanyByID(ctx context.Context, id string) (*model.Company, error)
	UpdateCompany(ctx context.Context, id string, name string) (*model.Company, error)
	DeleteCompany(ctx context.Context, id string) error
}

type companyService struct {
	db       db.Db
	repo     repository.CompanyRepository
	userRepo pkgRepository.UserRepository
}

func NewCompanyService(db db.Db, repo repository.CompanyRepository, userRepo pkgRepository.UserRepository) CompanyService {
	return &companyService{
		db:       db,
		repo:     repo,
		userRepo: userRepo,
	}
}

func (s *companyService) CreateCompany(ctx context.Context, name, email, password string) (*model.Company, error) {
	tx, err := s.db.BegrinTransaction(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}

	txAdapter := db.NewTxAdapter(tx.(*sql.Tx))

	txCompanyRepo := s.repo.WithTx(txAdapter)
	txUserRepo := s.userRepo.WithTx(txAdapter)

	company := &model.Company{
		Name:      name,
		CreatedAt: time.Now(),
	}

	if err := txCompanyRepo.Create(ctx, company); err != nil {
		s.db.Rollback(ctx, tx)
		return nil, fmt.Errorf("failed to create company: %w", err)
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		s.db.Rollback(ctx, tx)
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	user := &model.User{
		ID:        uuid.New().String(),
		Email:     email,
		Password:  string(hashedPassword),
		Role:      string(model.RoleCompany),
		CompanyID: company.ID,
	}

	if err := txUserRepo.Create(ctx, user); err != nil {
		s.db.Rollback(ctx, tx)
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	if err := s.db.Commit(ctx, tx); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return company, nil
}

func (s *companyService) GetAllCompanies(ctx context.Context, limit, offset int) ([]*model.Company, int64, error) {
	return s.repo.FindAll(ctx, limit, offset)
}

func (s *companyService) GetCompanyByID(ctx context.Context, id string) (*model.Company, error) {
	return s.repo.FindByID(ctx, id)
}

func (s *companyService) UpdateCompany(ctx context.Context, id string, name string) (*model.Company, error) {
	company, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if company == nil {
		return nil, nil // Or explicit error not found
	}

	company.Name = name

	if err := s.repo.Update(ctx, company); err != nil {
		return nil, err
	}
	return company, nil
}

func (s *companyService) DeleteCompany(ctx context.Context, id string) error {
	company, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return err
	}
	if company == nil {
		return nil // Or error not found
	}
	return s.repo.Delete(ctx, id)
}
