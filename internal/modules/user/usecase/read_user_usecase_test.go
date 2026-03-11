package usecase

import (
	"context"
	"testing"

	"be-dashboard-nba/api/presenter"
	"be-dashboard-nba/constant"
	"be-dashboard-nba/internal/modules/user/domain"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestReadUsersUsecase(t *testing.T) {
	logger := zerolog.Nop()

	tests := []struct {
		name        string
		filter      domain.UserFilter
		wantErr     bool
		expectedErr error
		expectedRes domain.UserPaginationResponse
		setupMock   func(mockRepo *domain.MockUserRepository)
	}{
		{
			name: "success read users",
			filter: domain.UserFilter{
				Limit: 10,
				Page:  1,
			},
			wantErr: false,
			expectedRes: domain.UserPaginationResponse{
				Data: []domain.UserListRow{
					{ID: "1", FullName: "admin"},
				},
				Pagination: presenter.Pagination{
					Page:       1,
					PageSize:   10,
					TotalPages: 1,
					TotalItems: 1,
					HasNext:    false,
					HasPrev:    false,
				},
			},
			setupMock: func(mockRepo *domain.MockUserRepository) {
				mockRepo.On("ReadCountUserQuery", mock.Anything, mock.AnythingOfType("domain.UserFilter")).
					Return(1, nil).Once()

				mockRepo.On("ReadUsersQuery", mock.Anything, mock.AnythingOfType("domain.UserFilter")).
					Return([]domain.UserListRow{
						{ID: "1", FullName: "admin"},
					}, nil).Once()
			},
		},
		{
			name: "failed - count db error",
			filter: domain.UserFilter{
				Limit: 10,
				Page:  1,
			},
			wantErr:     true,
			expectedErr: constant.ErrUnknownSource,
			setupMock: func(mockRepo *domain.MockUserRepository) {
				mockRepo.On("ReadCountUserQuery", mock.Anything, mock.AnythingOfType("domain.UserFilter")).
					Return(0, context.DeadlineExceeded).Once()
			},
		},
		{
			name: "failed - data db error",
			filter: domain.UserFilter{
				Limit: 10,
				Page:  1,
			},
			wantErr:     true,
			expectedErr: constant.ErrUnknownSource,
			setupMock: func(mockRepo *domain.MockUserRepository) {
				mockRepo.On("ReadCountUserQuery", mock.Anything, mock.AnythingOfType("domain.UserFilter")).
					Return(1, nil).Once()

				mockRepo.On("ReadUsersQuery", mock.Anything, mock.AnythingOfType("domain.UserFilter")).
					Return([]domain.UserListRow{}, context.DeadlineExceeded).Once()
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			mockRepo := new(domain.MockUserRepository)
			tc.setupMock(mockRepo)

			usecase := NewUserUsecase(mockRepo, &logger)
			ctx := context.Background()

			res, err := usecase.ReadUsersUsecase(ctx, tc.filter)

			if tc.wantErr {
				assert.Error(t, err)
				if tc.expectedErr != nil {
					assert.ErrorIs(t, err, tc.expectedErr)
				}
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expectedRes.Data, res.Data)
				assert.Equal(t, tc.expectedRes.Pagination, res.Pagination)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}
