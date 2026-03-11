package usecase

import (
	"context"
	"testing"

	"be-dashboard-nba/constant"
	"be-dashboard-nba/internal/modules/user/domain"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestUpdateUserUsecase(t *testing.T) {
	logger := zerolog.Nop()

	activeTrue := true

	tests := []struct {
		name        string
		payload     domain.UpdateUserPayload
		wantErr     bool
		expectedErr error
		setupMock   func(mockRepo *domain.MockUserRepository)
	}{
		{
			name: "success update user",
			payload: domain.UpdateUserPayload{
				ID:     "user-123",
				Name:   "john_doe",
				RoleID: 1,
				Active: &activeTrue,
			},
			wantErr: false,
			setupMock: func(mockRepo *domain.MockUserRepository) {
				mockRepo.On("ReadUserByIDQuery", mock.Anything, "user-123").
					Return(domain.UserDetailRow{ID: "user-123", Active: false}, nil).Once()

				// Expectation untuk WithTx dan ExpectBegin/Commit dihapus
				// Langsung test pemanggilan UpdateUserWithRoleTx
				mockRepo.On("UpdateUserWithRoleTx", mock.Anything, mock.MatchedBy(func(u domain.User) bool {
					return u.ID == "user-123" && u.Active == true
				})).Return(nil).Once()
			},
		},
		{
			name: "failed - user not found",
			payload: domain.UpdateUserPayload{
				ID: "user-123",
			},
			wantErr:     true,
			expectedErr: constant.ErrUserIdNotFound,
			setupMock: func(mockRepo *domain.MockUserRepository) {
				mockRepo.On("ReadUserByIDQuery", mock.Anything, "user-123").
					Return(domain.UserDetailRow{}, constant.ErrDataNotFound).Once()
			},
		},
		{
			name: "failed - role not found constraint mapping",
			payload: domain.UpdateUserPayload{
				ID: "user-123",
			},
			wantErr:     true,
			expectedErr: constant.ErrRoleIdNotFound,
			setupMock: func(mockRepo *domain.MockUserRepository) {
				mockRepo.On("ReadUserByIDQuery", mock.Anything, "user-123").
					Return(domain.UserDetailRow{ID: "user-123"}, nil).Once()

				// Expectation rollback dihapus, langsung return error constraint
				mockRepo.On("UpdateUserWithRoleTx", mock.Anything, mock.AnythingOfType("domain.User")).
					Return(constant.ErrRoleIdNotFound).Once()
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			mockRepo := new(domain.MockUserRepository)

			// Setup sqlmock dihapus sepenuhnya dari test runner ini

			tc.setupMock(mockRepo)

			// PERHATIAN:
			// Pastikan constructor NewUserUsecase sudah tidak menerima parameter *sql.DB.
			// Jika sebelumnya NewUserUsecase(mockRepo, &logger, dbMock), ubah menjadi:
			usecase := NewUserUsecase(mockRepo, &logger)
			ctx := context.Background()

			err := usecase.UpdateUserUsecase(ctx, tc.payload)

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
