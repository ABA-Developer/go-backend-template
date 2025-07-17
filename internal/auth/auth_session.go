package auth

import (
	"context"
	"time"

	"github.com/gofiber/fiber/v2/log"

	"be-dashboard-nba/constant"
	"be-dashboard-nba/pkg/auth/repository"
)

func (a *Auth) ValidateSession(ctx context.Context) (err error) {
	q := repository.NewQuery(a.db)

	session, err := q.ReadDetailSessionQuery(ctx, a.claims.GUID)
	if err != nil {
		log.WithContext(ctx).Error(err, "error read session by guid : "+a.claims.GUID)
		return
	}

	if time.Now().After(session.AccessTokenExpiredAt) {
		err = constant.ErrTokenExpired
		return
	}

	return
}
