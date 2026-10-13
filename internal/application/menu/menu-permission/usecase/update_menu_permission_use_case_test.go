package usecase

import (
	"context"
	"database/sql"
	"testing"

	"github.com/ABA-Developer/go-backend-template/constant"
	"github.com/ABA-Developer/go-backend-template/internal/application/menu/menu-permission/dto"
	repoParams "github.com/ABA-Developer/go-backend-template/internal/domain/contract/repository/params"
	"github.com/ABA-Developer/go-backend-template/internal/domain/model"
	mockrepo "github.com/ABA-Developer/go-backend-template/internal/test/mock/repository"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	testifymock "github.com/stretchr/testify/mock"
)

func TestUpdateMenuPermissionUseCase(t *testing.T) {
	tests := []struct {
		name             string
		menuPermissionID int
		req              dto.UpdateMenuPermissionPayload
		setup            func(m sqlmock.Sqlmock, menuPermissionRepo *mockrepo.MenuPermissionRepositoryMock)
		wantErr          error
	}{
		{
			name:             "positive/update-menu-permission",
			menuPermissionID: 1,
			req: dto.UpdateMenuPermissionPayload{
				Code:       "U",
				ActionName: "update",
			},
			setup: func(m sqlmock.Sqlmock, menuPermissionRepo *mockrepo.MenuPermissionRepositoryMock) {
				m.ExpectBegin()
				m.ExpectCommit()

				menuPermissionRepo.On("ReadMenuPermissionByIdQry", testifymock.Anything, 1).Return(model.MenuPermission{ID: 1}, nil).Once()
				menuPermissionRepo.
					On("UpdateMenuPermissionQry", testifymock.Anything, testifymock.Anything).
					Run(func(args testifymock.Arguments) {
						params := args.Get(1).(repoParams.UpdateMenuPermissionParams)
						assert.Equal(t, 1, params.MenuPermissionID)
						assert.Equal(t, "U", params.Code)
						assert.Equal(t, "update", params.ActionName)
						assert.Equal(t, "user-1", params.UpdatedBy)
					}).
					Return(nil).
					Once()
			},
		},
		{
			name:             "negative/not-found",
			menuPermissionID: 999,
			setup: func(m sqlmock.Sqlmock, menuPermissionRepo *mockrepo.MenuPermissionRepositoryMock) {
				m.ExpectBegin()
				m.ExpectRollback()

				menuPermissionRepo.On("ReadMenuPermissionByIdQry", testifymock.Anything, 999).Return(model.MenuPermission{}, sql.ErrNoRows).Once()
			},
			wantErr: constant.ErrMenuPermissionIdNotFound,
		},
		{
			name:             "negative/update-error",
			menuPermissionID: 1,
			setup: func(m sqlmock.Sqlmock, menuPermissionRepo *mockrepo.MenuPermissionRepositoryMock) {
				m.ExpectBegin()
				m.ExpectRollback()

				menuPermissionRepo.On("ReadMenuPermissionByIdQry", testifymock.Anything, 1).Return(model.MenuPermission{ID: 1}, nil).Once()
				menuPermissionRepo.On("UpdateMenuPermissionQry", testifymock.Anything, testifymock.Anything).Return(errors.New("update failed")).Once()
			},
			wantErr: constant.ErrUnknownSource,
		},
		{
			name:             "edge/rollback-error",
			menuPermissionID: 1,
			setup: func(m sqlmock.Sqlmock, menuPermissionRepo *mockrepo.MenuPermissionRepositoryMock) {
				m.ExpectBegin()
				m.ExpectRollback().WillReturnError(sql.ErrConnDone)

				menuPermissionRepo.On("ReadMenuPermissionByIdQry", testifymock.Anything, 1).Return(model.MenuPermission{}, errors.New("read failed")).Once()
			},
			wantErr: constant.ErrUnknownSource,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			menuPermissionRepo := mockrepo.NewMenuPermissionRepository()
			uc, mock, cleanup := newTestUseCaseForUpdate(t, menuPermissionRepo)
			defer cleanup()

			tt.setup(mock, menuPermissionRepo)

			err := uc.UpdateMenuPermissionUseCase(context.Background(), tt.req, "user-1", tt.menuPermissionID)
			if tt.wantErr != nil {
				assert.Error(t, err)
				assert.ErrorIs(t, err, tt.wantErr)
				return
			}
			assert.NoError(t, err)
		})
	}
}
