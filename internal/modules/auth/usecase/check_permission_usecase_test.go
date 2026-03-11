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

func TestCheckPermissionUsecase(t *testing.T) {
	mockRepo := new(domain.MockAuthRepository)

	logger := zerolog.New(nil)
	u := NewAuthUsecase(mockRepo, &logger, nil)

	menuURL := constant.MenuKey("dashboards")
	userID := "user-123"
	codes := []string{"read", "create"}

	t.Run("Success - Has Permission", func(t *testing.T) {
		mockRepo.On("CheckPermissionQuery", mock.Anything, menuURL, userID, codes).
			Return(true, nil).Once()

		hasAccess, err := u.CheckPermissionUsecase(context.Background(), menuURL, userID, codes)

		assert.NoError(t, err)
		assert.True(t, hasAccess)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Success - No Permission", func(t *testing.T) {
		mockRepo.On("CheckPermissionQuery", mock.Anything, menuURL, userID, codes).
			Return(false, nil).Once()

		hasAccess, err := u.CheckPermissionUsecase(context.Background(), menuURL, userID, codes)

		assert.NoError(t, err)
		assert.False(t, hasAccess)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Error - Database error", func(t *testing.T) {
		mockRepo.On("CheckPermissionQuery", mock.Anything, menuURL, userID, codes).
			Return(false, errors.New("db disconnect")).Once()

		hasAccess, err := u.CheckPermissionUsecase(context.Background(), menuURL, userID, codes)

		assert.Error(t, err)
		assert.False(t, hasAccess)
		mockRepo.AssertExpectations(t)
	})
}
