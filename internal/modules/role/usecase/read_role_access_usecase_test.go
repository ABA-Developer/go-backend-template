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

func TestReadRoleAccessUsecase(t *testing.T) {
	logger := zerolog.Nop()

	tests := []struct {
		name        string
		filter      domain.RoleAccessFilter
		wantErr     bool
		expectedErr error
		expectedRes []domain.RoleAccessResponse
		expectedCnt int
		setupMock   func(mockRepo *domain.MockRoleRepository)
	}{
		{
			name: "success read role access",
			filter: domain.RoleAccessFilter{
				RoleID: 1,
				Limit:  10,
				Page:   1,
			},
			wantErr: false,
			expectedRes: []domain.RoleAccessResponse{
				{RoleID: 1, RoleName: "Admin", MenuID: 1, PermissionName: "Read", HasAccess: true},
			},
			expectedCnt: 1,
			setupMock: func(mockRepo *domain.MockRoleRepository) {
				mockRepo.On("ReadRoleByIDQuery", mock.Anything, 1).
					Return(domain.Role{ID: 1, Name: "Admin"}, nil).Once()

				mockRepo.On("ReadRoleAccessCount", mock.Anything, mock.AnythingOfType("domain.RoleAccessFilter")).
					Return(1, nil).Once()

				mockRepo.On("ReadRoleAccessQuery", mock.Anything, mock.AnythingOfType("domain.RoleAccessFilter")).
					Return([]domain.RoleAccessResponse{
						{RoleID: 1, RoleName: "Admin", MenuID: 1, PermissionName: "Read", HasAccess: true},
					}, nil).Once()
			},
		},
		{
			name: "failed - role not found",
			filter: domain.RoleAccessFilter{
				RoleID: 99,
			},
			wantErr:     true,
			expectedErr: constant.ErrRoleIdNotFound,
			setupMock: func(mockRepo *domain.MockRoleRepository) {
				mockRepo.On("ReadRoleByIDQuery", mock.Anything, 99).
					Return(domain.Role{}, constant.ErrDataNotFound).Once()
			},
		},
		{
			name: "failed - db error on read query",
			filter: domain.RoleAccessFilter{
				RoleID: 1,
				Limit:  10,
				Page:   1,
			},
			wantErr:     true,
			expectedErr: constant.ErrUnknownSource,
			setupMock: func(mockRepo *domain.MockRoleRepository) {
				mockRepo.On("ReadRoleByIDQuery", mock.Anything, 1).
					Return(domain.Role{ID: 1, Name: "Admin"}, nil).Once()

				mockRepo.On("ReadRoleAccessCount", mock.Anything, mock.AnythingOfType("domain.RoleAccessFilter")).
					Return(1, nil).Once()

				mockRepo.On("ReadRoleAccessQuery", mock.Anything, mock.AnythingOfType("domain.RoleAccessFilter")).
					Return([]domain.RoleAccessResponse{}, context.DeadlineExceeded).Once()
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			mockRepo := new(domain.MockRoleRepository)
			tc.setupMock(mockRepo)

			usecase := NewRoleUsecase(mockRepo, &logger)
			ctx := context.Background()

			data, count, err := usecase.ReadRoleAccessUsecase(ctx, tc.filter)

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
