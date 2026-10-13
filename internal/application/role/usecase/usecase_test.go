package usecase

import (
	"testing"

	"github.com/ABA-Developer/go-backend-template/internal/application/utils"
	contractRepo "github.com/ABA-Developer/go-backend-template/internal/domain/contract/repository"
	db "github.com/ABA-Developer/go-backend-template/internal/infrastructure/database"
	mockrepo "github.com/ABA-Developer/go-backend-template/internal/test/mock/repository"

	"github.com/DATA-DOG/go-sqlmock"
)

func newTestUseCaseForUpdateRole(t *testing.T, roleRepo *mockrepo.RoleRepositoryMock) (*useCase, sqlmock.Sqlmock, func()) {
	t.Helper()

	mockDB, mock, cleanup := utils.NewMockDB(t)
	uc := &useCase{
		db:          mockDB,
		newRoleRepo: func(db.Query) contractRepo.RoleRepository { return roleRepo },
	}
	return uc, mock, cleanup
}

func newTestUseCaseForCreateRole(t *testing.T, roleRepo *mockrepo.RoleRepositoryMock) (*useCase, sqlmock.Sqlmock, func()) {
	t.Helper()

	return newTestUseCaseForUpdateRole(t, roleRepo)
}

func newTestUseCaseForUpdateAccess(
	t *testing.T,
	roleRepo *mockrepo.RoleRepositoryMock,
	menuPermissionRepo *mockrepo.MenuPermissionRepositoryMock,
) (*useCase, sqlmock.Sqlmock, func()) {
	t.Helper()

	mockDB, mock, cleanup := utils.NewMockDB(t)
	uc := &useCase{
		db:          mockDB,
		newRoleRepo: func(db.Query) contractRepo.RoleRepository { return roleRepo },
		newMenuPermissionRepo: func(db.Query) contractRepo.MenuPermissionRepository {
			return menuPermissionRepo
		},
	}
	return uc, mock, cleanup
}
