package usecase

import (
	"context"
	"database/sql"
	"testing"

	"github.com/ABA-Developer/go-backend-template/constant"
	appdto "github.com/ABA-Developer/go-backend-template/internal/application/menu/menu-permission/dto"
	"github.com/ABA-Developer/go-backend-template/internal/application/utils"
	"github.com/ABA-Developer/go-backend-template/internal/domain/model"
	mockrepo "github.com/ABA-Developer/go-backend-template/internal/test/mock/repository"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	testifymock "github.com/stretchr/testify/mock"
)

func TestReadMenuPermissionUseCase(t *testing.T) {
	tests := []struct {
		name    string
		menuID  int
		req     appdto.ReadMenuPermissionListPayload
		setup   func(menuRepo *mockrepo.MenuRepositoryMock, menuPermissionRepo *mockrepo.MenuPermissionRepositoryMock)
		wantErr error
	}{
		{
			name:   "positive/returns-pagination",
			menuID: 1,
			req: appdto.ReadMenuPermissionListPayload{
				PaginationPayload: utils.PaginationPayload{Limit: 2, Page: 1},
			},
			setup: func(menuRepo *mockrepo.MenuRepositoryMock, menuPermissionRepo *mockrepo.MenuPermissionRepositoryMock) {
				menuRepo.On("ReadMenuByIDQry", testifymock.Anything, 1).Return(model.Menu{ID: 1}, nil).Once()
				menuPermissionRepo.On("ReadMenuPermissionCountQry", testifymock.Anything, testifymock.Anything).Return(3, nil).Once()
				menuPermissionRepo.On("ReadMenuPermissionListQry", testifymock.Anything, testifymock.Anything).Return([]model.MenuPermission{{ID: 1}}, nil).Once()
			},
		},
		{
			name:   "negative/menu-not-found",
			menuID: 999,
			req:    appdto.ReadMenuPermissionListPayload{PaginationPayload: utils.PaginationPayload{Limit: 2, Page: 1}},
			setup: func(menuRepo *mockrepo.MenuRepositoryMock, _ *mockrepo.MenuPermissionRepositoryMock) {
				menuRepo.On("ReadMenuByIDQry", testifymock.Anything, 999).Return(model.Menu{}, sql.ErrNoRows).Once()
			},
			wantErr: constant.ErrMenuIdNotFound,
		},
		{
			name:   "negative/count-error",
			menuID: 1,
			req:    appdto.ReadMenuPermissionListPayload{PaginationPayload: utils.PaginationPayload{Limit: 2, Page: 1}},
			setup: func(menuRepo *mockrepo.MenuRepositoryMock, menuPermissionRepo *mockrepo.MenuPermissionRepositoryMock) {
				menuRepo.On("ReadMenuByIDQry", testifymock.Anything, 1).Return(model.Menu{ID: 1}, nil).Once()
				menuPermissionRepo.On("ReadMenuPermissionCountQry", testifymock.Anything, testifymock.Anything).Return(0, errors.New("db down")).Once()
			},
			wantErr: constant.ErrUnknownSource,
		},
		{
			name:   "negative/list-error",
			menuID: 1,
			req:    appdto.ReadMenuPermissionListPayload{PaginationPayload: utils.PaginationPayload{Limit: 2, Page: 1}},
			setup: func(menuRepo *mockrepo.MenuRepositoryMock, menuPermissionRepo *mockrepo.MenuPermissionRepositoryMock) {
				menuRepo.On("ReadMenuByIDQry", testifymock.Anything, 1).Return(model.Menu{ID: 1}, nil).Once()
				menuPermissionRepo.On("ReadMenuPermissionCountQry", testifymock.Anything, testifymock.Anything).Return(1, nil).Once()
				menuPermissionRepo.On("ReadMenuPermissionListQry", testifymock.Anything, testifymock.Anything).Return(nil, errors.New("db down")).Once()
			},
			wantErr: constant.ErrUnknownSource,
		},
		{
			name:   "edge/default-limit-page",
			menuID: 1,
			req:    appdto.ReadMenuPermissionListPayload{},
			setup: func(menuRepo *mockrepo.MenuRepositoryMock, menuPermissionRepo *mockrepo.MenuPermissionRepositoryMock) {
				menuRepo.On("ReadMenuByIDQry", testifymock.Anything, 1).Return(model.Menu{ID: 1}, nil).Once()
				menuPermissionRepo.On("ReadMenuPermissionCountQry", testifymock.Anything, testifymock.Anything).Return(0, nil).Once()
				menuPermissionRepo.On("ReadMenuPermissionListQry", testifymock.Anything, testifymock.Anything).Return([]model.MenuPermission{}, nil).Once()
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			menuRepo := mockrepo.NewMenuRepository()
			menuPermissionRepo := mockrepo.NewMenuPermissionRepository()
			uc, _, cleanup := newTestUseCaseWithRepos(t, menuRepo, menuPermissionRepo)
			defer cleanup()

			tt.setup(menuRepo, menuPermissionRepo)

			_, err := uc.ReadMenuPermissionUseCase(context.Background(), tt.req, tt.menuID)
			if tt.wantErr != nil {
				assert.Error(t, err)
				assert.ErrorIs(t, err, tt.wantErr)
				return
			}
			assert.NoError(t, err)
		})
	}
}
