package service

import (
	authPresenter "be-dashboard-nba/api/presenter/auth"
	"be-dashboard-nba/constant"
	"be-dashboard-nba/internal/db"
	"be-dashboard-nba/internal/jwt"
	"be-dashboard-nba/pkg/entities"
	"context"

	// "github.com/gofiber/fiber/v2/log" // Dihapus, kita akan gunakan s.log

	"github.com/rs/zerolog"
)

type Service interface {
	LoginService(ctx context.Context, request authPresenter.LoginRequest, userAgent, iPAddress string) (data entities.Session, user entities.User, err error)
	LogoutService(ctx context.Context, claims *jwt.AccessTokenPayload, iPAddress string) (err error)
	AuthMeService(ctx context.Context, id string) (data entities.User, err error)
	CheckPermissionService(ctx context.Context, menuURL constant.MenuKey, userID string, permissionCode []string) (hasAccess bool, err error)
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
