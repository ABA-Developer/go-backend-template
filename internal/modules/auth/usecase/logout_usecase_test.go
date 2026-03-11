package usecase

import (
	"be-dashboard-nba/constant"
	"be-dashboard-nba/internal/core/jwt"
	"be-dashboard-nba/internal/modules/auth/domain"
	"context"
	"errors"
	"testing"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestLogoutUsecase(t *testing.T) {
	mockRepo := new(domain.MockAuthRepository)

	logger := zerolog.New(nil)
	u := NewAuthUsecase(mockRepo, &logger, nil)

	claims := &jwt.AccessTokenPayload{
		SessionID: "session-123",
		UserID:    "user-123",
	}
	ipAddress := "127.0.0.1"

	t.Run("Success - Logout", func(t *testing.T) {
		mockRepo.On("DeleteSessionQuery", mock.Anything, claims.SessionID).
			Return(nil).Once()

		recordArgs := domain.LoginRecord{
			UserID:    claims.UserID,
			Action:    "logout",
			IPAddress: ipAddress,
		}
		mockRepo.On("CreateLoginRecord", mock.Anything, recordArgs).
			Return(nil).Once()

		err := u.LogoutUsecase(context.Background(), claims, ipAddress)

		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Error - Session deletion failed", func(t *testing.T) {
		mockRepo.On("DeleteSessionQuery", mock.Anything, claims.SessionID).
			Return(errors.New("db error")).Once()

		err := u.LogoutUsecase(context.Background(), claims, ipAddress)

		assert.Error(t, err)
		assert.Equal(t, constant.ErrUnknownSource, errors.Unwrap(err))
		mockRepo.AssertExpectations(t)
	})

	t.Run("Success - Record logging failed but logout succeeds", func(t *testing.T) {
		mockRepo.On("DeleteSessionQuery", mock.Anything, claims.SessionID).
			Return(nil).Once()

		recordArgs := domain.LoginRecord{
			UserID:    claims.UserID,
			Action:    "logout",
			IPAddress: ipAddress,
		}
		mockRepo.On("CreateLoginRecord", mock.Anything, recordArgs).
			Return(errors.New("db error")).Once()

		err := u.LogoutUsecase(context.Background(), claims, ipAddress)

		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})
}
