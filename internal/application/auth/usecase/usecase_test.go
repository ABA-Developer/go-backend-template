package usecase

import (
	"testing"

	"github.com/ABA-Developer/go-backend-template/internal/application/utils"
	contractRepo "github.com/ABA-Developer/go-backend-template/internal/domain/contract/repository"
	db "github.com/ABA-Developer/go-backend-template/internal/infrastructure/database"
	mockrepo "github.com/ABA-Developer/go-backend-template/internal/test/mock/repository"

	"github.com/DATA-DOG/go-sqlmock"
)

func newTestUseCaseWithAuthRepo(t *testing.T, authRepo *mockrepo.AuthRepositoryMock) (*useCase, sqlmock.Sqlmock, func()) {
	t.Helper()

	mockDB, mock, cleanup := utils.NewMockDB(t)
	uc := &useCase{
		db:          mockDB,
		newAuthRepo: func(db.Query) contractRepo.AuthRepository { return authRepo },
	}
	return uc, mock, cleanup
}
