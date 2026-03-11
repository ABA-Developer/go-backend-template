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

func TestReadDetailRoleUsecase(t *testing.T) {
	logger := zerolog.Nop()

	tests := []struct {
		name        string
		roleID      int
		wantErr     bool
		expectedErr error
		expectedRes domain.Role
		setupMock   func(mockRepo *domain.MockRoleRepository)
	}{
		{
			name:    "success read role detail",
			roleID:  1,
			wantErr: false,
			expectedRes: domain.Role{
				ID:   1,
				Name: "Admin",
			},
			setupMock: func(mockRepo *domain.MockRoleRepository) {
				mockRepo.On("ReadRoleByIDQuery", mock.Anything, 1).
					Return(domain.Role{ID: 1, Name: "Admin"}, nil).Once()
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
			name:        "failed - db error",
			roleID:      1,
			wantErr:     true,
			expectedErr: constant.ErrUnknownSource,
			setupMock: func(mockRepo *domain.MockRoleRepository) {
				mockRepo.On("ReadRoleByIDQuery", mock.Anything, 1).
					Return(domain.Role{}, context.DeadlineExceeded).Once()
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			mockRepo := new(domain.MockRoleRepository)
			tc.setupMock(mockRepo)

			usecase := NewRoleUsecase(mockRepo, &logger)
			ctx := context.Background()

			res, err := usecase.ReadDetailRoleUsecase(ctx, tc.roleID)

			if tc.wantErr {
				assert.Error(t, err)
				if tc.expectedErr != nil {
					assert.ErrorIs(t, err, tc.expectedErr)
				}
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expectedRes, res)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}
