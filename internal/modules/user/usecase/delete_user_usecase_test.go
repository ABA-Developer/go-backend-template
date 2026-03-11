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

func TestDeleteUserUsecase(t *testing.T) {
	logger := zerolog.Nop()

	tests := []struct {
		name        string
		userID      string
		deletedBy   string
		wantErr     bool
		expectedErr error
		setupMock   func(mockRepo *domain.MockUserRepository)
	}{
		{
			name:      "success delete user",
			userID:    "user-123",
			deletedBy: "admin-456",
			wantErr:   false,
			setupMock: func(mockRepo *domain.MockUserRepository) {
				// Mock user exists and is not self
				mockRepo.On("ReadUserByIDQuery", mock.Anything, "user-123").
					Return(domain.UserDetailRow{ID: "user-123", Name: "User"}, nil).Once()

				mockRepo.On("DeleteUserQuery", mock.Anything, "user-123").
					Return(nil).Once()
			},
		},
		{
			name:        "failed - user not found",
			userID:      "user-123",
			deletedBy:   "admin-456",
			wantErr:     true,
			expectedErr: constant.ErrUserIdNotFound,
			setupMock: func(mockRepo *domain.MockUserRepository) {
				mockRepo.On("ReadUserByIDQuery", mock.Anything, "user-123").
					Return(domain.UserDetailRow{}, constant.ErrDataNotFound).Once()
			},
		},
		{
			name:        "failed - cannot delete self",
			userID:      "admin-456",
			deletedBy:   "admin-456",
			wantErr:     true,
			expectedErr: constant.ErrForbiddenSelfDelete,
			setupMock: func(mockRepo *domain.MockUserRepository) {
				mockRepo.On("ReadUserByIDQuery", mock.Anything, "admin-456").
					Return(domain.UserDetailRow{ID: "admin-456", Name: "Admin"}, nil).Once()
			},
		},
		{
			name:        "failed - db error on delete",
			userID:      "user-123",
			deletedBy:   "admin-456",
			wantErr:     true,
			expectedErr: constant.ErrUnknownSource,
			setupMock: func(mockRepo *domain.MockUserRepository) {
				mockRepo.On("ReadUserByIDQuery", mock.Anything, "user-123").
					Return(domain.UserDetailRow{ID: "user-123", Name: "User"}, nil).Once()

				mockRepo.On("DeleteUserQuery", mock.Anything, "user-123").
					Return(context.DeadlineExceeded).Once()
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			mockRepo := new(domain.MockUserRepository)
			tc.setupMock(mockRepo)

			usecase := NewUserUsecase(mockRepo, &logger)
			ctx := context.Background()

			err := usecase.DeleteUserUsecase(ctx, tc.userID, tc.deletedBy)

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
