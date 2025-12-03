package auth

import (
	"context"
	"time"

	"github.com/gofiber/fiber/v2/log"

	"be-dashboard-nba/constant"
	"be-dashboard-nba/pkg/auth/repository"
)

func (a *Auth) ValidateSession(ctx context.Context) (err error) {
	r := repository.NewRepo(a.db)

	session, err := r.ReadDetailSessionQuery(ctx, a.claims.SessionID)
	if err != nil {
		log.WithContext(ctx).Error(err, "error read session by id : "+a.claims.SessionID)
		return
	}

	if time.Now().After(session.AccessTokenExpiredAt) {
		err = constant.ErrTokenExpired
		return
	}

	return
}
