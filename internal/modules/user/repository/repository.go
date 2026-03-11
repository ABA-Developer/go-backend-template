package repository

import (
	"be-dashboard-nba/internal/core/db"
	"be-dashboard-nba/internal/modules/user/domain"
)

type userRepository struct {
	db db.DB
}

// NewUserRepository creates a new user repository that implements the domain.UserRepository interface
func NewUserRepository(db db.DB) domain.UserRepository {
	return &userRepository{
		db: db,
	}
}
