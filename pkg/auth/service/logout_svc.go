package service

import (
	"context"

	"be-dashboard-nba/internal/jwt"
	"be-dashboard-nba/pkg/auth/repository"

	"github.com/gofiber/fiber/v2/log"
)

func (s *Service) LogoutService(
	ctx context.Context,
	claims *jwt.AccessTokenPayload,
) (err error) {
	q := repository.NewQuery(s.db)

	err = q.DeleteSessionQuery(ctx, claims.GUID)
	if err != nil {
		log.WithContext(ctx).Error(err, "error delete session", "session guid", claims.GUID)
		return
	}

	return
}
