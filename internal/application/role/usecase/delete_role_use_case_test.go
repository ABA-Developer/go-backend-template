package usecase

import (
	"context"
	"database/sql"
	"testing"

	"github.com/ABA-Developer/go-backend-template/constant"
	"github.com/ABA-Developer/go-backend-template/internal/domain/model"
	mockrepo "github.com/ABA-Developer/go-backend-template/internal/test/mock/repository"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	testifymock "github.com/stretchr/testify/mock"
)

func TestDeleteRoleUseCase(t *testing.T) {
	tests := []struct {
		name    string
		roleID  int
		setup   func(m sqlmock.Sqlmock, roleRepo *mockrepo.RoleRepositoryMock)
		wantErr error
	}{
		{
			name:   "positive/delete-role",
			roleID: 10,
			setup: func(m sqlmock.Sqlmock, roleRepo *mockrepo.RoleRepositoryMock) {
				m.ExpectBegin()
				m.ExpectCommit()

				roleRepo.On("ReadRoleByIDQry", testifymock.Anything, 10).Return(model.Role{ID: 10}, nil).Once()
				roleRepo.On("DeleteRoleQry", testifymock.Anything, 10).Return(nil).Once()
			},
		},
		{
			name:   "negative/not-found",
			roleID: 999,
			setup: func(m sqlmock.Sqlmock, roleRepo *mockrepo.RoleRepositoryMock) {
				m.ExpectBegin()
				m.ExpectRollback()

				roleRepo.On("ReadRoleByIDQry", testifymock.Anything, 999).Return(model.Role{}, sql.ErrNoRows).Once()
			},
			wantErr: constant.ErrRoleIdNotFound,
		},
		{
			name:   "negative/delete-error",
			roleID: 10,
			setup: func(m sqlmock.Sqlmock, roleRepo *mockrepo.RoleRepositoryMock) {
				m.ExpectBegin()
				m.ExpectRollback()

				roleRepo.On("ReadRoleByIDQry", testifymock.Anything, 10).Return(model.Role{ID: 10}, nil).Once()
				roleRepo.On("DeleteRoleQry", testifymock.Anything, 10).Return(errors.New("delete failed")).Once()
			},
			wantErr: constant.ErrUnknownSource,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			roleRepo := mockrepo.NewRoleRepository()
			uc, mock, cleanup := newTestUseCaseForUpdateRole(t, roleRepo)
			defer cleanup()

			tt.setup(mock, roleRepo)

			err := uc.DeleteRoleUseCase(context.Background(), tt.roleID)
			if tt.wantErr != nil {
				assert.Error(t, err)
				assert.ErrorIs(t, err, tt.wantErr)
				return
			}
			assert.NoError(t, err)
		})
	}
}
