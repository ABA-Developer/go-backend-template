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

func TestAuthMeUseCase(t *testing.T) {
	tests := []struct {
		name    string
		userID  string
		setup   func(authRepo *mockrepo.AuthRepositoryMock)
		wantErr error
	}{
		{
			name:   "positive/returns-user",
			userID: "u1",
			setup: func(authRepo *mockrepo.AuthRepositoryMock) {
				authRepo.On("ReadDetailUserByIdQry", testifymock.Anything, "u1").Return(model.User{ID: "u1"}, nil).Once()
			},
		},
		{
			name:   "negative/not-found",
			userID: "missing",
			setup: func(authRepo *mockrepo.AuthRepositoryMock) {
				authRepo.On("ReadDetailUserByIdQry", testifymock.Anything, "missing").Return(model.User{}, sql.ErrNoRows).Once()
			},
			wantErr: constant.ErrUserIdNotFound,
		},
		{
			name:   "negative/repo-error",
			userID: "u1",
			setup: func(authRepo *mockrepo.AuthRepositoryMock) {
				authRepo.On("ReadDetailUserByIdQry", testifymock.Anything, "u1").Return(model.User{}, errors.New("db down")).Once()
			},
			wantErr: constant.ErrUnknownSource,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			authRepo := mockrepo.NewAuthRepository()
			uc, _, cleanup := newTestUseCaseWithAuthRepo(t, authRepo)
			defer cleanup()

			tt.setup(authRepo)

			_, err := uc.AuthMeUseCase(context.Background(), tt.userID)
			if tt.wantErr != nil {
				assert.Error(t, err)
				assert.ErrorIs(t, err, tt.wantErr)
				return
			}
			assert.NoError(t, err)
		})
	}
}
