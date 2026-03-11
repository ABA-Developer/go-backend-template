package repository

import (
	"be-dashboard-nba/internal/modules/auth/domain"
	"be-dashboard-nba/internal/core/db"
	"database/sql"
)

type repository struct {
	db db.Query
}

func NewAuthRepository(db db.Query) domain.AuthRepository {
	return &repository{
		db: db,
	}
}

func (r *repository) WithTx(tx *sql.Tx) domain.AuthRepository {
	return &repository{
		db: tx,
	}
}
