package psql

import (
	"context"
	"errors"
	"testing"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/golang/mock/gomock"
	"github.com/hfleury/bk_globalshot/internal/model"
	mock_repository "github.com/hfleury/bk_globalshot/mock/repository"
	"github.com/stretchr/testify/assert"
)

func TestFindByEmail(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockedUserRepo := mock_repository.NewMockUserRepository(ctrl)

	t.Run("Success - Find user by email", func(t *testing.T) {
		email := gofakeit.Email()
		user := &model.User{
			ID:       gofakeit.UUID(),
			Email:    email,
			Password: gofakeit.Password(true, true, true, true, true, 26),
		}

		mockedUserRepo.EXPECT().FindByEmail(gomock.Any(), email).Return(user, nil)
		result, err := mockedUserRepo.FindByEmail(context.Background(), email)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, email, user.Email)
	})

	t.Run("No user found", func(t *testing.T) {
		email := "notfound@example.com"

		mockedUserRepo.EXPECT().
			FindByEmail(gomock.Any(), email).
			Return(nil, nil)

		result, err := mockedUserRepo.FindByEmail(context.Background(), email)

		assert.NoError(t, err)
		assert.Nil(t, result)
	})

	t.Run("Database error", func(t *testing.T) {
		email := "error@example.com"
		mockErr := errors.New("database error")

		mockedUserRepo.EXPECT().
			FindByEmail(gomock.Any(), email).
			Return(nil, mockErr)

		result, err := mockedUserRepo.FindByEmail(context.Background(), email)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, mockErr, err)
	})

}

func TestCreate(t *testing.T) {

}
