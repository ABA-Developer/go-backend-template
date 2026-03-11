package usecase

import (
	"be-dashboard-nba/constant"
	"be-dashboard-nba/internal/core/jwt"
	"be-dashboard-nba/internal/modules/auth/domain"
	"context"

	"github.com/pkg/errors"
)

func (s *authUsecase) LogoutUsecase(
	ctx context.Context,
	claims *jwt.AccessTokenPayload,
	iPAddress string,
) (err error) {

	err = s.repo.DeleteSessionQuery(ctx, claims.SessionID)
	if err != nil {
		s.log.Error().Err(err).Str("session_id", claims.SessionID).Msg("error delete session")
		err = errors.WithStack(constant.ErrUnknownSource)
		return
	}

	logRecordArgs := domain.LoginRecord{
		UserID:    claims.UserID,
		Action:    "logout",
		IPAddress: iPAddress,
	}

	if logErr := s.repo.CreateLoginRecord(ctx, logRecordArgs); logErr != nil {
		s.log.Error().Err(logErr).Msg("failed to create logout record log")
	}

	return nil
}
