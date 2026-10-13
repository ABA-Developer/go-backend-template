package usecase

import (
	contractRepo "github.com/ABA-Developer/go-backend-template/internal/domain/contract/repository"
	contract "github.com/ABA-Developer/go-backend-template/internal/domain/contract/usecase"
	db "github.com/ABA-Developer/go-backend-template/internal/infrastructure/database"
	roleRepo "github.com/ABA-Developer/go-backend-template/internal/infrastructure/repository/role"
	userRepo "github.com/ABA-Developer/go-backend-template/internal/infrastructure/repository/user"
)

type useCase struct {
	db db.DB

	newUserRepo func(q db.Query) contractRepo.UserRepository
	newRoleRepo func(q db.Query) contractRepo.RoleRepository
}

func NewUseCase(db db.DB) contract.UserUseCase {
	return &useCase{
		db:          db,
		newUserRepo: userRepo.NewQuery,
		newRoleRepo: roleRepo.NewQuery,
	}
}
