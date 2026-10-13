package usecase

import (
	contractRepo "github.com/ABA-Developer/go-backend-template/internal/domain/contract/repository"
	contract "github.com/ABA-Developer/go-backend-template/internal/domain/contract/usecase"
	db "github.com/ABA-Developer/go-backend-template/internal/infrastructure/database"
	menuRepo "github.com/ABA-Developer/go-backend-template/internal/infrastructure/repository/menu"
)

type useCase struct {
	db          db.DB
	newMenuRepo func(q db.Query) contractRepo.MenuRepository
}

func NewUseCase(db db.DB) contract.MenuUseCase {
	return &useCase{
		db:          db,
		newMenuRepo: menuRepo.NewQuery,
	}
}
