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

func TestUpdateRoleUsecase(t *testing.T) {
	logger := zerolog.Nop()
	description := "Administrator updated"

	tests := []struct {
		name        string
		payload     domain.UpdateRolePayload
		wantErr     bool
		expectedErr error
		setupMock   func(mockRepo *domain.MockRoleRepository)
	}{
		{
			name: "success update role",
			payload: domain.UpdateRolePayload{
				RoleID:      1,
				Code:        "ADMIN",
				Name:        "Admin Edited",
				Description: &description,
			},
			wantErr: false,
			setupMock: func(mockRepo *domain.MockRoleRepository) {
				mockRepo.On("ReadRoleByIDQuery", mock.Anything, 1).
					Return(domain.Role{ID: 1, Name: "Admin"}, nil).Once()

				mockRepo.On("UpdateRoleQuery", mock.Anything, mock.AnythingOfType("domain.UpdateRolePayload")).
					Return(nil).Once()
			},
		},
		{
			name: "failed - role not found",
			payload: domain.UpdateRolePayload{
				RoleID: 99,
				Code:   "NOTFOUND",
			},
			wantErr:     true,
			expectedErr: constant.ErrRoleIdNotFound,
			setupMock: func(mockRepo *domain.MockRoleRepository) {
				mockRepo.On("ReadRoleByIDQuery", mock.Anything, 99).
					Return(domain.Role{}, constant.ErrDataNotFound).Once()
			},
		},
		{
			name: "failed - update query error",
			payload: domain.UpdateRolePayload{
				RoleID:      1,
				Code:        "ADMIN",
				Name:        "Admin Edited",
				Description: &description,
			},
			wantErr:     true,
			expectedErr: constant.ErrUnknownSource,
			setupMock: func(mockRepo *domain.MockRoleRepository) {
				mockRepo.On("ReadRoleByIDQuery", mock.Anything, 1).
					Return(domain.Role{ID: 1, Name: "Admin"}, nil).Once()

				mockRepo.On("UpdateRoleQuery", mock.Anything, mock.AnythingOfType("domain.UpdateRolePayload")).
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

			err := usecase.UpdateRoleUsecase(ctx, tc.payload)

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
