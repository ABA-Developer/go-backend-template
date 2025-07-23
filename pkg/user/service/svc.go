package service

import (
	"be-dashboard-nba/internal/db"
)

type Service struct {
	db db.DB
}

func NewService(db db.DB) *Service {
	return &Service{
		db: db,
	}
}
