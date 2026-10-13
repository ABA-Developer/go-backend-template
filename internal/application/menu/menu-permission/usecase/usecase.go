package usecase

import (
	contractRepo "github.com/ABA-Developer/go-backend-template/internal/domain/contract/repository"
	contract "github.com/ABA-Developer/go-backend-template/internal/domain/contract/usecase"
	db "github.com/ABA-Developer/go-backend-template/internal/infrastructure/database"
	menuRepo "github.com/ABA-Developer/go-backend-template/internal/infrastructure/repository/menu"
	menuPermissionRepo "github.com/ABA-Developer/go-backend-template/internal/infrastructure/repository/menu/menu-permission"
)

type useCase struct {
	db                    db.DB
	newMenuRepo           func(q db.Query) contractRepo.MenuRepository
	newMenuPermissionRepo func(q db.Query) contractRepo.MenuPermissionRepository
}

func NewUseCase(db db.DB) contract.MenuPermissionUseCase {
	return &useCase{
		db:                    db,
		newMenuRepo:           menuRepo.NewQuery,
		newMenuPermissionRepo: menuPermissionRepo.NewQuery,
	}
}
