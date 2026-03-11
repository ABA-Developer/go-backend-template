package usecase

import (
	"context"
	"testing"

	"be-dashboard-nba/constant"
	"be-dashboard-nba/internal/modules/role/domain"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCreateRoleUsecase(t *testing.T) {
	logger := zerolog.Nop()
	description := "Administrator Role"

	tests := []struct {
		name        string
		payload     domain.CreateRolePayload
		wantErr     bool
		expectedErr error
		setupMock   func(mockRepo *domain.MockRoleRepository)
	}{
		{
			name: "success create role",
			payload: domain.CreateRolePayload{
				Code:        "ADMIN",
				Name:        "Admin",
				Description: &description,
				CreatedBy:   "user-123",
			},
			wantErr: false,
			setupMock: func(mockRepo *domain.MockRoleRepository) {
				mockRepo.On("CreateRoleQuery", mock.Anything, mock.AnythingOfType("domain.CreateRolePayload")).
					Return(nil).Once()
			},
		},
		{
			name: "failed create role - db error",
			payload: domain.CreateRolePayload{
				Code: "ADMIN",
				Name: "Admin",
			},
			wantErr:     true,
			expectedErr: constant.ErrUnknownSource,
			setupMock: func(mockRepo *domain.MockRoleRepository) {
				mockRepo.On("CreateRoleQuery", mock.Anything, mock.AnythingOfType("domain.CreateRolePayload")).
					Return(context.DeadlineExceeded).Once()
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			mockRepo := new(domain.MockRoleRepository)
			tc.setupMock(mockRepo)

			usecase := NewRoleUsecase(mockRepo, &logger)
			ctx := context.Background()

			err := usecase.CreateRoleUsecase(ctx, tc.payload)

			if tc.wantErr {
				assert.Error(t, err)
				if tc.expectedErr != nil {
					assert.ErrorIs(t, err, tc.expectedErr)
				}
			} else {
				assert.NoError(t, err)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}
