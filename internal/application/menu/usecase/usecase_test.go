package usecase

import (
	"testing"

	"github.com/ABA-Developer/go-backend-template/internal/application/utils"
	contractRepo "github.com/ABA-Developer/go-backend-template/internal/domain/contract/repository"
	db "github.com/ABA-Developer/go-backend-template/internal/infrastructure/database"
	mockrepo "github.com/ABA-Developer/go-backend-template/internal/test/mock/repository"

	"github.com/DATA-DOG/go-sqlmock"
)

func newTestUseCaseWithMenuRepo(t *testing.T, menuRepo *mockrepo.MenuRepositoryMock) (*useCase, sqlmock.Sqlmock, func()) {
	t.Helper()

	mockDB, mock, cleanup := utils.NewMockDB(t)
	uc := &useCase{
		db:          mockDB,
		newMenuRepo: func(db.Query) contractRepo.MenuRepository { return menuRepo },
	}
	return uc, mock, cleanup
}

func newTestUseCaseForDeleteMenu(t *testing.T, menuRepo *mockrepo.MenuRepositoryMock) (*useCase, sqlmock.Sqlmock, func()) {
	t.Helper()
	return newTestUseCaseWithMenuRepo(t, menuRepo)
}

func newTestUseCaseForReadMenu(t *testing.T, menuRepo *mockrepo.MenuRepositoryMock) (*useCase, sqlmock.Sqlmock, func()) {
	t.Helper()
	return newTestUseCaseWithMenuRepo(t, menuRepo)
}

func newTestUseCaseForUpdateMenuOrder(t *testing.T, menuRepo *mockrepo.MenuRepositoryMock) (*useCase, sqlmock.Sqlmock, func()) {
	t.Helper()
	return newTestUseCaseWithMenuRepo(t, menuRepo)
}
