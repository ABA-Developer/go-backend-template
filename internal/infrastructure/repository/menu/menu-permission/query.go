package repository

import (
	db "github.com/ABA-Developer/go-backend-template/internal/infrastructure/database"

	contractRepo "github.com/ABA-Developer/go-backend-template/internal/domain/contract/repository"
)

type Query struct {
	db db.Query
}

func NewQuery(db db.Query) contractRepo.MenuPermissionRepository {
	return &Query{
		db: db,
	}
}
