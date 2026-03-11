package usecase

import (
	"be-dashboard-nba/constant"
	"be-dashboard-nba/internal/modules/auth/domain"
	"be-dashboard-nba/internal/core/db"
	"be-dashboard-nba/internal/core/jwt"
	"context"

	"github.com/rs/zerolog"
)

type AuthUsecase interface {
	LoginService(ctx context.Context, email, password, userAgent, iPAddress string) (data domain.Session, user domain.User, err error)
	LogoutService(ctx context.Context, claims *jwt.AccessTokenPayload, iPAddress string) (err error)
	AuthMeService(ctx context.Context, id string) (data domain.User, err error)
	CheckPermissionService(ctx context.Context, menuURL constant.MenuKey, userID string, permissionCode []string) (hasAccess bool, err error)
}

type authUsecase struct {
	repo domain.AuthRepository
	log  *zerolog.Logger
	db   db.DB
}

func NewAuthUsecase(repo domain.AuthRepository, log *zerolog.Logger, db db.DB) domain.AuthUsecase {
	return &authUsecase{
		repo: repo,
		log:  log,
		db:   db,
	}
}
