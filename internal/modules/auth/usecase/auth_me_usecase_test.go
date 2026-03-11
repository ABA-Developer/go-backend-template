package usecase

import (
	"be-dashboard-nba/constant"
	"be-dashboard-nba/internal/modules/auth/domain"
	"context"
	"errors"
	"testing"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestAuthMeUsecase(t *testing.T) {
	mockRepo := new(domain.MockAuthRepository)

	logger := zerolog.New(nil)
	u := NewAuthUsecase(mockRepo, &logger, nil)

	userID := "user-123"

	t.Run("Success - AuthMe detail retrieved", func(t *testing.T) {
		expectedUser := domain.User{
			ID:    userID,
			Email: "test@example.com",
		}

		mockRepo.On("ReadDetailUserByIdQuery", mock.Anything, userID).
			Return(expectedUser, nil).Once()

		result, err := u.AuthMeUsecase(context.Background(), userID)

		assert.NoError(t, err)
		assert.Equal(t, expectedUser, result)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Error - User Not Found", func(t *testing.T) {
		mockRepo.On("ReadDetailUserByIdQuery", mock.Anything, userID).
			Return(domain.User{}, errors.New("sql: no rows in result set")).Once()

		result, err := u.AuthMeUsecase(context.Background(), userID)

		assert.Error(t, err)
		assert.Equal(t, constant.ErrUserIdNotFound, err)
		assert.Empty(t, result)
		mockRepo.AssertExpectations(t)
	})
}
