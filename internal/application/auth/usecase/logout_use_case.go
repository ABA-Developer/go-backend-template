package usecase

import (
	"github.com/ABA-Developer/go-backend-template/constant"
	"github.com/ABA-Developer/go-backend-template/internal/application/jwt"
	repoParams "github.com/ABA-Developer/go-backend-template/internal/domain/contract/repository/params"
	"context"

	"github.com/pkg/errors"
)

func (s *useCase) LogoutUseCase(
	ctx context.Context,
	claims *jwt.AccessTokenPayload,
	iPAddress string,
) (err error) {
	q := s.newAuthRepo(s.db)
	logRepo := s.newAuthRepo(s.db)

	err = q.DeleteSessionQry(ctx, claims.SessionID)
	if err != nil {
		log(ctx).Error().Err(err).Str("session_id", claims.SessionID).Msg("error delete session")
		err = errors.WithStack(constant.ErrUnknownSource)
		return
	}

	logRecordArgs := repoParams.LoginRecord{
		UserID:      claims.UserID,
		AccessToken: claims.SessionID,
		Status:      "logout",
		IPAddress:   iPAddress,
		Type:        "web",
	}
	if logErr := logRepo.CreateLoginRecordQry(ctx, logRecordArgs); logErr != nil {
		log(ctx).Error().Err(logErr).Msg("failed to create logout record log")
	}

	return nil
}
