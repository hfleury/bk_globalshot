package service

import (
	"context"
	"fmt"

	"github.com/hfleury/bk_globalshot/pkg/config"
	"github.com/hfleury/bk_globalshot/pkg/repository"
	"github.com/hfleury/bk_globalshot/pkg/token"
	"golang.org/x/crypto/bcrypt"
)

//go:generate mockgen -source=auth_service.go -destination=../../mock/services/mock_auth_service.go -package=mock_services
type AuthService interface {
	Login(ctx context.Context, email, password string) (string, string, bool, error)
}

type authService struct {
	repo     repository.UserRepository
	maker    token.Maker
	cfgToken *config.ConfigToken
}

func NewAuthService(repo repository.UserRepository, maker token.Maker, cfgToken *config.ConfigToken) AuthService {
	return &authService{repo: repo, maker: maker, cfgToken: cfgToken}
}

func (s *authService) Login(ctx context.Context, email, password string) (string, string, bool, error) {
	fmt.Printf("DEBUG: Login attempt for email: %s\n", email)
	user, err := s.repo.FindByEmail(ctx, email)
	if err != nil {
		fmt.Printf("DEBUG: FindByEmail error: %v\n", err)
		return "", "", false, err
	}
	if user == nil {
		fmt.Println("DEBUG: User not found")
		return "", "", false, nil
	}

	fmt.Printf("DEBUG: User found: %s, Hash: %s\n", user.Email, user.Password)

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		fmt.Printf("DEBUG: Password mismatch error: %v\n", err)
		return "", "", false, nil
	}

	token, err := s.maker.CreateToken(user.ID, user.Email, user.Role, user.CompanyID, s.cfgToken.TokenExpiry)
	if err != nil {
		fmt.Printf("DEBUG: Token creation error: %v\n", err)
		return "", "", false, err
	}

	return token, user.Role, true, nil
}
