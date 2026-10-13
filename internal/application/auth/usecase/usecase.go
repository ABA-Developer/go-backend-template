package usecase

import (
	contractRepo "github.com/ABA-Developer/go-backend-template/internal/domain/contract/repository"
	contract "github.com/ABA-Developer/go-backend-template/internal/domain/contract/usecase"
	db "github.com/ABA-Developer/go-backend-template/internal/infrastructure/database"
	authRepo "github.com/ABA-Developer/go-backend-template/internal/infrastructure/repository/auth"
)

type useCase struct {
	db          db.DB
	newAuthRepo func(q db.Query) contractRepo.AuthRepository
}

func NewUseCase(db db.DB) contract.AuthUseCase {
	return &useCase{
		db:          db,
		newAuthRepo: authRepo.NewQuery,
	}
}
