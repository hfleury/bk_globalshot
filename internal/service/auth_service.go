package service

import (
	"context"

	"github.com/hfleury/bk_globalshot/pkg/config"
	"github.com/hfleury/bk_globalshot/pkg/repository"
	"github.com/hfleury/bk_globalshot/pkg/token"
	"golang.org/x/crypto/bcrypt"
)

//go:generate mockgen -source=auth_service.go -destination=../../mock/services/mock_auth_service.go -package=mock_services
type AuthService interface {
	Login(ctx context.Context, email, password string) (string, bool, error)
}

type authService struct {
	repo     repository.UserRepository
	maker    token.Maker
	cfgToken *config.ConfigToken
}

func NewAuthService(repo repository.UserRepository, maker token.Maker) AuthService {
	return &authService{repo: repo, maker: maker}
}

func (s *authService) Login(ctx context.Context, email, password string) (string, bool, error) {
	user, err := s.repo.FindByEmail(ctx, email)
	if err != nil || user == nil {
		return "", false, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return "", false, nil
	}

	token, err := s.maker.CreateToken(user.Email, s.cfgToken.TokenExpiry)
	if err != nil {
		return "", false, err
	}

	return token, true, nil
}
