package usecase

import (
	contractRepo "github.com/ABA-Developer/go-backend-template/internal/domain/contract/repository"
	contract "github.com/ABA-Developer/go-backend-template/internal/domain/contract/usecase"
	db "github.com/ABA-Developer/go-backend-template/internal/infrastructure/database"
	menuPermissionRepo "github.com/ABA-Developer/go-backend-template/internal/infrastructure/repository/menu/menu-permission"
	roleRepo "github.com/ABA-Developer/go-backend-template/internal/infrastructure/repository/role"
)

type useCase struct {
	db db.DB

	newRoleRepo           func(q db.Query) contractRepo.RoleRepository
	newMenuPermissionRepo func(q db.Query) contractRepo.MenuPermissionRepository
}

func NewUseCase(db db.DB) contract.RoleUseCase {
	return &useCase{
		db:                    db,
		newRoleRepo:           roleRepo.NewQuery,
		newMenuPermissionRepo: menuPermissionRepo.NewQuery,
	}
}
