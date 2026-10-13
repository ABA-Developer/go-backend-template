package usecase

import (
	"testing"

	"github.com/ABA-Developer/go-backend-template/internal/application/utils"
	contractRepo "github.com/ABA-Developer/go-backend-template/internal/domain/contract/repository"
	db "github.com/ABA-Developer/go-backend-template/internal/infrastructure/database"
	mockrepo "github.com/ABA-Developer/go-backend-template/internal/test/mock/repository"

	"github.com/DATA-DOG/go-sqlmock"
)

func newTestUseCaseWithUserAndRoleRepo(
	t *testing.T,
	userRepo *mockrepo.UserRepositoryMock,
	roleRepo *mockrepo.RoleRepositoryMock,
) (*useCase, sqlmock.Sqlmock, func()) {
	t.Helper()

	mockDB, mock, cleanup := utils.NewMockDB(t)
	uc := &useCase{
		db:          mockDB,
		newUserRepo: func(db.Query) contractRepo.UserRepository { return userRepo },
		newRoleRepo: func(db.Query) contractRepo.RoleRepository { return roleRepo },
	}
	return uc, mock, cleanup
}

func newTestUseCaseForDeleteUser(t *testing.T, userRepo *mockrepo.UserRepositoryMock) (*useCase, sqlmock.Sqlmock, func()) {
	t.Helper()

	mockDB, mock, cleanup := utils.NewMockDB(t)
	uc := &useCase{
		db:          mockDB,
		newUserRepo: func(db.Query) contractRepo.UserRepository { return userRepo },
	}
	return uc, mock, cleanup
}

func newTestUseCaseForUpdateUser(
	t *testing.T,
	userRepo *mockrepo.UserRepositoryMock,
	roleRepo *mockrepo.RoleRepositoryMock,
) (*useCase, sqlmock.Sqlmock, func()) {
	t.Helper()
	return newTestUseCaseWithUserAndRoleRepo(t, userRepo, roleRepo)
}
