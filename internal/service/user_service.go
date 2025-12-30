package service

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/hfleury/bk_globalshot/internal/model"
	"github.com/hfleury/bk_globalshot/pkg/repository"
	"golang.org/x/crypto/bcrypt"
)

type UserService interface {
	CreateUser(ctx context.Context, email, password, role string, companyID string) (*model.User, error)
	GetAllUsers(ctx context.Context, limit, offset int) ([]*model.User, int64, error)
	GetUserByID(ctx context.Context, id string) (*model.User, error)
	UpdateUser(ctx context.Context, id, email, role string, companyID string) (*model.User, error)
	DeleteUser(ctx context.Context, id string) error
}

type userService struct {
	repo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) UserService {
	return &userService{
		repo: repo,
	}
}

func (s *userService) CreateUser(ctx context.Context, email, password, role string, companyID string) (*model.User, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	user := &model.User{
		ID:        uuid.New().String(),
		Email:     email,
		Password:  string(hashedPassword),
		Role:      role,
		CompanyID: companyID,
	}

	if err := s.repo.Create(ctx, user); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *userService) GetAllUsers(ctx context.Context, limit, offset int) ([]*model.User, int64, error) {
	return s.repo.FindAll(ctx, limit, offset)
}

func (s *userService) GetUserByID(ctx context.Context, id string) (*model.User, error) {
	return s.repo.FindByID(ctx, id)
}

func (s *userService) UpdateUser(ctx context.Context, id, email, role string, companyID string) (*model.User, error) {
	user, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, nil // Not found
	}

	user.Email = email
	user.Role = role
	user.CompanyID = companyID

	if err := s.repo.Update(ctx, user); err != nil {
		return nil, err
	}
	return user, nil
}

func (s *userService) DeleteUser(ctx context.Context, id string) error {
	user, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return err
	}
	if user == nil {
		return nil
	}
	return s.repo.Delete(ctx, id)
}
