package usecase

import (
	"context"
	"database/sql"
	"testing"

	"github.com/ABA-Developer/go-backend-template/constant"
	"github.com/ABA-Developer/go-backend-template/internal/domain/model"
	mockrepo "github.com/ABA-Developer/go-backend-template/internal/test/mock/repository"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	testifymock "github.com/stretchr/testify/mock"
)

func TestReadUserProfileUseCase(t *testing.T) {
	tests := []struct {
		name    string
		userID  string
		setup   func(userRepo *mockrepo.UserRepositoryMock)
		wantErr error
	}{
		{
			name:   "positive/returns-data",
			userID: "u1",
			setup: func(userRepo *mockrepo.UserRepositoryMock) {
				userRepo.On("ReadUserProfileQry", testifymock.Anything, "u1").Return(model.User{ID: "u1"}, nil).Once()
			},
		},
		{
			name:   "negative/not-found",
			userID: "missing",
			setup: func(userRepo *mockrepo.UserRepositoryMock) {
				userRepo.On("ReadUserProfileQry", testifymock.Anything, "missing").Return(model.User{}, sql.ErrNoRows).Once()
			},
			wantErr: constant.ErrUserIdNotFound,
		},
		{
			name:   "negative/repo-error",
			userID: "u1",
			setup: func(userRepo *mockrepo.UserRepositoryMock) {
				userRepo.On("ReadUserProfileQry", testifymock.Anything, "u1").Return(model.User{}, errors.New("db down")).Once()
			},
			wantErr: constant.ErrUnknownSource,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			userRepo := mockrepo.NewUserRepository()
			uc, _, cleanup := newTestUseCaseForDeleteUser(t, userRepo)
			defer cleanup()

			tt.setup(userRepo)

			_, err := uc.ReadUserProfileUseCase(context.Background(), tt.userID)
			if tt.wantErr != nil {
				assert.Error(t, err)
				assert.ErrorIs(t, err, tt.wantErr)
				return
			}
			assert.NoError(t, err)
		})
	}
}
