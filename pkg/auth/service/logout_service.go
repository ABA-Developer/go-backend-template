package service

import (
	"be-dashboard-nba/constant"
	"be-dashboard-nba/internal/jwt"
	"be-dashboard-nba/pkg/auth/repository"
	"context"

	"github.com/pkg/errors"
)

func (s *service) LogoutService(
	ctx context.Context,
	claims *jwt.AccessTokenPayload,
	iPAddress string,
) (err error) {
	q := repository.NewRepo(s.db)
	logRepo := repository.NewRepo(s.db)

	err = q.DeleteSessionQuery(ctx, claims.SessionID)
	if err != nil {
		s.log.Error().Err(err).Str("session_id", claims.SessionID).Msg("error delete session")
		err = errors.WithStack(constant.ErrUnknownSource)
		return
	}

	logRecordArgs := repository.LoginRecord{
		UserID:      claims.UserID,
		AccessToken: claims.SessionID,
		Status:      "logout",
		IPAddress:   iPAddress,
		Type:        "web",
	}
	if logErr := logRepo.CreateLoginRecord(ctx, logRecordArgs); logErr != nil {
		s.log.Error().Err(logErr).Msg("failed to create logout record log")
	}

	return nil
}
