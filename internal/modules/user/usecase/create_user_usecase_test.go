package usecase_test

import (
	"context"
	"os"
	"testing"

	"be-dashboard-nba/constant"
	"be-dashboard-nba/internal/modules/user/domain"
	"be-dashboard-nba/internal/modules/user/usecase"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCreateUserUsecase(t *testing.T) {
	os.Setenv("AUTH_BCRYPT_COST", "10")
	logger := zerolog.Nop()

	phone := "08123456789"

	tests := []struct {
		name        string
		request     domain.User
		wantErr     bool
		expectedErr error
		setupMock   func(mockRepo *domain.MockUserRepository)
	}{
		{
			name: "success create user",
			request: domain.User{
				Name:      "john_doe",
				FullName:  "John Doe",
				Email:     "john.doe@example.com",
				Password:  "password123",
				RoleID:    1,
				Phone:     &phone,
				Active:    true,
				CreatedBy: "admin-123",
			},
			wantErr:     false,
			expectedErr: nil,
			setupMock: func(mockRepo *domain.MockUserRepository) {
				mockRepo.On("CreateUserWithRoleTx", mock.Anything, mock.MatchedBy(func(u domain.User) bool {
					return u.Name == "john_doe" && len(u.Password) > 10 // checking if hashed
				})).Return(nil).Once()
			},
		},
		{
			name: "failed - role id not found / email duplicate",
			request: domain.User{
				Name:     "john_doe",
				Password: "password123",
			},
			wantErr:     true,
			expectedErr: constant.ErrRoleIdNotFound,
			setupMock: func(mockRepo *domain.MockUserRepository) {
				mockRepo.On("CreateUserWithRoleTx", mock.Anything, mock.AnythingOfType("domain.User")).
					Return(constant.ErrRoleIdNotFound).Once()
			},
		},
		{
			name: "failed - unknown error from repository",
			request: domain.User{
				Name:     "john_doe",
				Password: "password123",
			},
			wantErr:     true,
			expectedErr: constant.ErrUnknownSource,
			setupMock: func(mockRepo *domain.MockUserRepository) {
				mockRepo.On("CreateUserWithRoleTx", mock.Anything, mock.AnythingOfType("domain.User")).
					Return(context.DeadlineExceeded).Once()
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			mockRepo := new(domain.MockUserRepository)

			tc.setupMock(mockRepo)

			uc := usecase.NewUserUsecase(mockRepo, &logger)
			ctx := context.Background()

			err := uc.CreateUserUsecase(ctx, tc.request)

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
