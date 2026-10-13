package usecase

import (
	"context"
	"testing"

	"github.com/ABA-Developer/go-backend-template/constant"
	appdto "github.com/ABA-Developer/go-backend-template/internal/application/user/dto"
	"github.com/ABA-Developer/go-backend-template/internal/application/utils"
	"github.com/ABA-Developer/go-backend-template/internal/domain/model"
	mockrepo "github.com/ABA-Developer/go-backend-template/internal/test/mock/repository"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	testifymock "github.com/stretchr/testify/mock"
)

func TestReadUsersUseCase(t *testing.T) {
	tests := []struct {
		name    string
		req     appdto.ReadUserPayload
		setup   func(userRepo *mockrepo.UserRepositoryMock)
		wantErr error
	}{
		{
			name: "positive/returns-pagination",
			req:  appdto.ReadUserPayload{PaginationPayload: utils.PaginationPayload{Limit: 2, Page: 1}},
			setup: func(userRepo *mockrepo.UserRepositoryMock) {
				userRepo.On("ReadCountUserQry", testifymock.Anything, testifymock.Anything).Return(3, nil).Once()
				userRepo.On("ReadUsersQry", testifymock.Anything, testifymock.Anything).Return([]model.User{{ID: "u1"}}, nil).Once()
			},
		},
		{
			name: "negative/count-error",
			req:  appdto.ReadUserPayload{PaginationPayload: utils.PaginationPayload{Limit: 2, Page: 1}},
			setup: func(userRepo *mockrepo.UserRepositoryMock) {
				userRepo.On("ReadCountUserQry", testifymock.Anything, testifymock.Anything).Return(0, errors.New("db down")).Once()
			},
			wantErr: constant.ErrUnknownSource,
		},
		{
			name: "negative/list-error",
			req:  appdto.ReadUserPayload{PaginationPayload: utils.PaginationPayload{Limit: 2, Page: 1}},
			setup: func(userRepo *mockrepo.UserRepositoryMock) {
				userRepo.On("ReadCountUserQry", testifymock.Anything, testifymock.Anything).Return(1, nil).Once()
				userRepo.On("ReadUsersQry", testifymock.Anything, testifymock.Anything).Return(nil, errors.New("db down")).Once()
			},
			wantErr: constant.ErrUnknownSource,
		},
		{
			name: "edge/default-limit-page",
			req:  appdto.ReadUserPayload{},
			setup: func(userRepo *mockrepo.UserRepositoryMock) {
				userRepo.On("ReadCountUserQry", testifymock.Anything, testifymock.Anything).Return(0, nil).Once()
				userRepo.On("ReadUsersQry", testifymock.Anything, testifymock.Anything).Return([]model.User{}, nil).Once()
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			userRepo := mockrepo.NewUserRepository()
			uc, _, cleanup := newTestUseCaseForDeleteUser(t, userRepo)
			defer cleanup()

			tt.setup(userRepo)

			_, err := uc.ReadUsersUseCase(context.Background(), tt.req)
			if tt.wantErr != nil {
				assert.Error(t, err)
				assert.ErrorIs(t, err, tt.wantErr)
				return
			}
			assert.NoError(t, err)
		})
	}
}
