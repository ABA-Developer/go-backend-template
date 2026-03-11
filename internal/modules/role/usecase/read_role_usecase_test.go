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

func TestReadRoleUsecase(t *testing.T) {
	logger := zerolog.Nop()

	tests := []struct {
		name        string
		filter      domain.RoleFilter
		wantErr     bool
		expectedErr error
		expectedRes []domain.Role
		expectedCnt int
		setupMock   func(mockRepo *domain.MockRoleRepository)
	}{
		{
			name: "success read roles",
			filter: domain.RoleFilter{
				Limit: 10,
				Page:  1,
			},
			wantErr: false,
			expectedRes: []domain.Role{
				{ID: 1, Name: "Admin"},
			},
			expectedCnt: 1,
			setupMock: func(mockRepo *domain.MockRoleRepository) {
				mockRepo.On("ReadRolesCount", mock.Anything, mock.AnythingOfType("domain.RoleFilter")).
					Return(1, nil).Once()

				mockRepo.On("ReadRolesQuery", mock.Anything, mock.AnythingOfType("domain.RoleFilter")).
					Return([]domain.Role{
						{ID: 1, Name: "Admin"},
					}, nil).Once()
			},
		},
		{
			name: "failed - count db error",
			filter: domain.RoleFilter{
				Limit: 10,
				Page:  1,
			},
			wantErr:     true,
			expectedErr: constant.ErrUnknownSource,
			setupMock: func(mockRepo *domain.MockRoleRepository) {
				mockRepo.On("ReadRolesCount", mock.Anything, mock.AnythingOfType("domain.RoleFilter")).
					Return(0, context.DeadlineExceeded).Once()
			},
		},
		{
			name: "failed - data db error",
			filter: domain.RoleFilter{
				Limit: 10,
				Page:  1,
			},
			wantErr:     true,
			expectedErr: constant.ErrUnknownSource,
			setupMock: func(mockRepo *domain.MockRoleRepository) {
				mockRepo.On("ReadRolesCount", mock.Anything, mock.AnythingOfType("domain.RoleFilter")).
					Return(1, nil).Once()

				mockRepo.On("ReadRolesQuery", mock.Anything, mock.AnythingOfType("domain.RoleFilter")).
					Return([]domain.Role{}, context.DeadlineExceeded).Once()
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			mockRepo := new(domain.MockRoleRepository)
			tc.setupMock(mockRepo)

			usecase := NewRoleUsecase(mockRepo, &logger)
			ctx := context.Background()

			data, count, err := usecase.ReadRoleUsecase(ctx, tc.filter)

			if tc.wantErr {
				assert.Error(t, err)
				if tc.expectedErr != nil {
					assert.ErrorIs(t, err, tc.expectedErr)
				}
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expectedRes, data)
				assert.Equal(t, tc.expectedCnt, count)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}
