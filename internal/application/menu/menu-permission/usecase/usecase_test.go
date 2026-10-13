package usecase

import (
	"testing"

	"github.com/ABA-Developer/go-backend-template/internal/application/utils"
	contractRepo "github.com/ABA-Developer/go-backend-template/internal/domain/contract/repository"
	db "github.com/ABA-Developer/go-backend-template/internal/infrastructure/database"
	mockrepo "github.com/ABA-Developer/go-backend-template/internal/test/mock/repository"

	"github.com/DATA-DOG/go-sqlmock"
)

func newTestUseCaseWithRepos(
	t *testing.T,
	menuRepo *mockrepo.MenuRepositoryMock,
	menuPermissionRepo *mockrepo.MenuPermissionRepositoryMock,
) (*useCase, sqlmock.Sqlmock, func()) {
	t.Helper()

	mockDB, mock, cleanup := utils.NewMockDB(t)
	uc := &useCase{
		db:                    mockDB,
		newMenuRepo:           func(db.Query) contractRepo.MenuRepository { return menuRepo },
		newMenuPermissionRepo: func(db.Query) contractRepo.MenuPermissionRepository { return menuPermissionRepo },
	}
	return uc, mock, cleanup
}

func newTestUseCaseForDelete(t *testing.T, menuPermissionRepo *mockrepo.MenuPermissionRepositoryMock) (*useCase, sqlmock.Sqlmock, func()) {
	t.Helper()

	mockDB, mock, cleanup := utils.NewMockDB(t)
	uc := &useCase{
		db: mockDB,
		newMenuPermissionRepo: func(db.Query) contractRepo.MenuPermissionRepository {
			return menuPermissionRepo
		},
	}
	return uc, mock, cleanup
}

func newTestUseCaseForUpdate(t *testing.T, menuPermissionRepo *mockrepo.MenuPermissionRepositoryMock) (*useCase, sqlmock.Sqlmock, func()) {
	t.Helper()
	return newTestUseCaseForDelete(t, menuPermissionRepo)
}
