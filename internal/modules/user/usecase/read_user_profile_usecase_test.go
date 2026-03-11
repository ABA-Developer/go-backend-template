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

func TestReadUserProfileUsecase(t *testing.T) {
	logger := zerolog.Nop()

	tests := []struct {
		name        string
		userID      string
		wantErr     bool
		expectedErr error
		expectedRes domain.UserProfileRow
		setupMock   func(mockRepo *domain.MockUserRepository)
	}{
		{
			name:    "success read user profile",
			userID:  "user-123",
			wantErr: false,
			expectedRes: domain.UserProfileRow{
				ID:       "user-123",
				Name:     "john_doe",
				FullName: "John Doe",
			},
			setupMock: func(mockRepo *domain.MockUserRepository) {
				mockRepo.On("ReadUserProfileQuery", mock.Anything, "user-123").
					Return(domain.UserProfileRow{
						ID:       "user-123",
						Name:     "john_doe",
						FullName: "John Doe",
					}, nil).Once()
			},
		},
		{
			name:        "failed - user not found",
			userID:      "user-123",
			wantErr:     true,
			expectedErr: constant.ErrUserIdNotFound,
			setupMock: func(mockRepo *domain.MockUserRepository) {
				mockRepo.On("ReadUserProfileQuery", mock.Anything, "user-123").
					Return(domain.UserProfileRow{}, constant.ErrDataNotFound).Once()
			},
		},
		{
			name:        "failed - db error",
			userID:      "user-123",
			wantErr:     true,
			expectedErr: constant.ErrUnknownSource,
			setupMock: func(mockRepo *domain.MockUserRepository) {
				mockRepo.On("ReadUserProfileQuery", mock.Anything, "user-123").
					Return(domain.UserProfileRow{}, context.DeadlineExceeded).Once()
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			mockRepo := new(domain.MockUserRepository)
			tc.setupMock(mockRepo)

			usecase := NewUserUsecase(mockRepo, &logger)
			ctx := context.Background()

			res, err := usecase.ReadUserProfileUsecase(ctx, tc.userID)

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
