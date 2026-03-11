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

func TestDeleteRoleUsecase(t *testing.T) {
	logger := zerolog.Nop()

	tests := []struct {
		name        string
		roleID      int
		wantErr     bool
		expectedErr error
		setupMock   func(mockRepo *domain.MockRoleRepository)
	}{
		{
			name:    "success delete role",
			roleID:  1,
			wantErr: false,
			setupMock: func(mockRepo *domain.MockRoleRepository) {
				mockRepo.On("ReadRoleByIDQuery", mock.Anything, 1).
					Return(domain.Role{ID: 1, Name: "Admin"}, nil).Once()

				mockRepo.On("DeleteRoleQuery", mock.Anything, 1).
					Return(nil).Once()
			},
		},
		{
			name:        "failed - role not found",
			roleID:      99,
			wantErr:     true,
			expectedErr: constant.ErrRoleIdNotFound,
			setupMock: func(mockRepo *domain.MockRoleRepository) {
				mockRepo.On("ReadRoleByIDQuery", mock.Anything, 99).
					Return(domain.Role{}, constant.ErrDataNotFound).Once()
			},
		},
		{
			name:        "failed - db error on read",
			roleID:      1,
			wantErr:     true,
			expectedErr: constant.ErrUnknownSource,
			setupMock: func(mockRepo *domain.MockRoleRepository) {
				mockRepo.On("ReadRoleByIDQuery", mock.Anything, 1).
					Return(domain.Role{}, context.DeadlineExceeded).Once()
			},
		},
		{
			name:        "failed - db error on delete",
			roleID:      1,
			wantErr:     true,
			expectedErr: constant.ErrUnknownSource,
			setupMock: func(mockRepo *domain.MockRoleRepository) {
				mockRepo.On("ReadRoleByIDQuery", mock.Anything, 1).
					Return(domain.Role{ID: 1, Name: "Admin"}, nil).Once()

				mockRepo.On("DeleteRoleQuery", mock.Anything, 1).
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

			err := usecase.DeleteRoleUsecase(ctx, tc.roleID)

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
