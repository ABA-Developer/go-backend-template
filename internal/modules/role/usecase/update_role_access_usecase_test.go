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

func TestUpdateRoleAccessUsecase(t *testing.T) {
	logger := zerolog.Nop()

	hasAccessTrue := true

	tests := []struct {
		name        string
		roleID      int
		payloads    []domain.UpdateRoleMenuPermission
		wantErr     bool
		expectedErr error
		setupMock   func(mockRepo *domain.MockRoleRepository)
	}{
		{
			name:   "success update role access",
			roleID: 1,
			payloads: []domain.UpdateRoleMenuPermission{
				{MenuPermissionID: 10, RoleID: 1, HasAccess: &hasAccessTrue},
			},
			wantErr: false,
			setupMock: func(mockRepo *domain.MockRoleRepository) {
				mockRepo.On("ReadRoleByIDQuery", mock.Anything, 1).
					Return(domain.Role{ID: 1, Name: "Admin"}, nil).Once()

				mockRepo.On("UpdateRoleAccessTx", mock.Anything, 1, mock.AnythingOfType("[]domain.UpdateRoleMenuPermission")).
					Return(nil).Once()
			},
		},
		{
			name:   "failed - role not found",
			roleID: 99,
			payloads: []domain.UpdateRoleMenuPermission{
				{MenuPermissionID: 10, RoleID: 99, HasAccess: &hasAccessTrue},
			},
			wantErr:     true,
			expectedErr: constant.ErrRoleIdNotFound,
			setupMock: func(mockRepo *domain.MockRoleRepository) {
				mockRepo.On("ReadRoleByIDQuery", mock.Anything, 99).
					Return(domain.Role{}, constant.ErrDataNotFound).Once()
			},
		},
		{
			name:   "failed - update role access via tx",
			roleID: 1,
			payloads: []domain.UpdateRoleMenuPermission{
				{MenuPermissionID: 10, RoleID: 1, HasAccess: &hasAccessTrue},
			},
			wantErr:     true,
			expectedErr: constant.ErrUnknownSource,
			setupMock: func(mockRepo *domain.MockRoleRepository) {
				mockRepo.On("ReadRoleByIDQuery", mock.Anything, 1).
					Return(domain.Role{ID: 1, Name: "Admin"}, nil).Once()

				mockRepo.On("UpdateRoleAccessTx", mock.Anything, 1, mock.AnythingOfType("[]domain.UpdateRoleMenuPermission")).
					Return(context.DeadlineExceeded).Once()
			},
		},
		{
			name:   "failed - menu permission not found (foreign key error)",
			roleID: 1,
			payloads: []domain.UpdateRoleMenuPermission{
				{MenuPermissionID: 99, RoleID: 1, HasAccess: &hasAccessTrue},
			},
			wantErr:     true,
			expectedErr: constant.ErrMenuPermissionIdNotFound,
			setupMock: func(mockRepo *domain.MockRoleRepository) {
				mockRepo.On("ReadRoleByIDQuery", mock.Anything, 1).
					Return(domain.Role{ID: 1, Name: "Admin"}, nil).Once()

				mockRepo.On("UpdateRoleAccessTx", mock.Anything, 1, mock.AnythingOfType("[]domain.UpdateRoleMenuPermission")).
					Return(constant.ErrMenuPermissionIdNotFound).Once()
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			mockRepo := new(domain.MockRoleRepository)
			tc.setupMock(mockRepo)

			usecase := NewRoleUsecase(mockRepo, &logger)
			ctx := context.Background()

			err := usecase.UpdateRoleAccessUsecase(ctx, tc.roleID, tc.payloads)

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
