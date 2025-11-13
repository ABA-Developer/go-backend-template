package service

import (
	presenter "be-dashboard-nba/api/presenter/user"
	"be-dashboard-nba/internal/db"
	"be-dashboard-nba/pkg/entities"
	"context"

	"github.com/rs/zerolog"
)

type Service interface {
	ReadDetailUserService(ctx context.Context, id string) (data entities.User, err error)
	UpdateUserService(ctx context.Context, request presenter.UpdateUserRequest, userID string) (err error)
}

type service struct {
	db  db.DB
	log *zerolog.Logger
}

func NewService(db db.DB, log *zerolog.Logger) Service {
	return &service{
		db:  db,
		log: log,
	}
}
