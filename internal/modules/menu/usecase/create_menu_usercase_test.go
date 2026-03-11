package usecase_test

import (
	"context"
	"testing"

	"be-dashboard-nba/constant"
	"be-dashboard-nba/internal/modules/menu/domain"
	"be-dashboard-nba/internal/modules/menu/usecase" // Sesuaikan path

	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
)

func TestCreateMenuUsecase(t *testing.T) {
	ctx := context.Background()
	userID := "user-123"

	valInt32 := int32(1)
	valString := "/dashboard"

	defaultPayload := domain.MenuCreatePayload{
		Name:     "Dashboard",
		Group:    "main",
		ParentID: &valInt32,
		URL:      &valString,
		Active:   true,
		Display:  true,
	}

	nopLogger := zerolog.Nop()

	tests := []struct {
		name        string
		payload     domain.MenuCreatePayload
		userID      string
		mockSetup   func(mockRepo *domain.MockMenuRepository)
		wantErr     bool
		expectedErr error
	}{
		{
			name:    "success - create menu",
			payload: defaultPayload,
			userID:  userID,
			mockSetup: func(mockRepo *domain.MockMenuRepository) {
				mockRepo.On("CreateMenuQuery", ctx, defaultPayload, userID).Return(nil)
			},
			wantErr:     false,
			expectedErr: nil,
		},
		{
			name:    "error - repository failed (simulating transaction/db error)",
			payload: defaultPayload,
			userID:  userID,
			mockSetup: func(mockRepo *domain.MockMenuRepository) {
				mockRepo.On("CreateMenuQuery", ctx, defaultPayload, userID).Return(errors.New("db transaction failed"))
			},
			wantErr:     true,
			expectedErr: constant.ErrUnknownSource,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			mockRepo := new(domain.MockMenuRepository)
			tc.mockSetup(mockRepo)

			// Masukkan mock repo dan logger kosong
			uc := usecase.NewMenuUsecase(mockRepo, &nopLogger)

			err := uc.CreateMenuUsecase(ctx, tc.payload, tc.userID)

			if tc.wantErr {
				assert.Error(t, err)
				assert.ErrorIs(t, err, tc.expectedErr)
			} else {
				assert.NoError(t, err)
			}

			// Validasi bahwa mock dipanggil sesuai ekspektasi
			mockRepo.AssertExpectations(t)
		})
	}
}
