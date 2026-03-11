package usecase

import (
	"be-dashboard-nba/constant"
	"be-dashboard-nba/internal/core/jwt"
	"be-dashboard-nba/internal/core/utils"
	"be-dashboard-nba/internal/modules/auth/domain"
	"context"
	"database/sql"

	"github.com/pkg/errors"
)

func (s *authUsecase) LoginUsecase(
	ctx context.Context,
	email, password, userAgent,
	ipAddress string,
) (data domain.Session, user domain.User, err error) {

	user, err = s.repo.ReadDetailUserByEmailQuery(ctx, email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			logAttemptArgs := domain.LoginAttemp{
				Email:     email,
				Password:  password,
				IPAddress: ipAddress,
			}
			if logErr := s.repo.CreateLoginAttemp(ctx, logAttemptArgs); logErr != nil {
				s.log.Error().Err(logErr).Msg("failed to create login attempt log (wrong email)")
			}

			err = constant.ErrWrongEmailOrPassword
			return
		}

		s.log.Error().Err(err).Str("email", email).Msg("error to read user by email")
		err = errors.WithStack(constant.ErrUnknownSource)
		return
	}

	if err = utils.CompareHashPassword(password, user.Password); err != nil {
		logAttemptArgs := domain.LoginAttemp{
			Email:     email,
			Password:  password,
			IPAddress: ipAddress,
		}
		if logErr := s.repo.CreateLoginAttemp(ctx, logAttemptArgs); logErr != nil {
			s.log.Error().Err(logErr).Msg("failed to create login attempt log (wrong password)")
		}
		err = constant.ErrWrongEmailOrPassword
		return
	}

	sessionID := utils.GenerateUUID()
	accessTokenReq := jwt.AccessTokenPayload{
		SessionID: sessionID,
		UserID:    user.ID,
	}
	accessToken, err := jwt.GenerateAccessToken(accessTokenReq)
	if err != nil {
		s.log.Error().Err(err).Msg("failed to generate access token")
		err = errors.WithStack(constant.ErrUnknownSource)
		return
	}

	refreshTokenReq := jwt.RefreshTokenPayload{
		SessionID: sessionID,
	}
	refreshToken, err := jwt.GenerateRefreshToken(refreshTokenReq)
	if err != nil {
		s.log.Error().Err(err).Msg("failed to generate refresh token")
		err = errors.WithStack(constant.ErrUnknownSource)
		return
	}

	data = domain.Session{
		ID:                    sessionID,
		UserID:                user.ID,
		AccessToken:           accessToken.Token,
		AccessTokenExpiredAt:  accessToken.ExpiresAt,
		RefreshToken:          refreshToken.Token,
		RefreshTokenExpiredAt: refreshToken.ExpiresAt,
		IPAddress:             ipAddress,
		UserAgent:             userAgent,
	}

	tx, err := s.db.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		s.log.Error().Err(err).Msg("error to begin transaction")
		err = errors.WithStack(constant.ErrUnknownSource)
		return
	}

	q := s.repo.WithTx(tx)

	defer func() {
		if err != nil {
			if !errors.Is(err, constant.ErrWrongEmailOrPassword) {
				if errRollback := tx.Rollback(); errRollback != nil {
					s.log.Error().Err(errRollback).AnErr("original_error", err).Msg("error to rollback transaction")
					err = errors.WithStack(constant.ErrUnknownSource)
					return
				}
				if !errors.Is(err, constant.ErrUnknownSource) {
					err = errors.WithStack(constant.ErrUnknownSource)
				}
			}
		}
	}()

	err = q.CreateSessionQuery(ctx, data)
	if err != nil {
		s.log.Error().Err(err).Interface("session", data).Msg("error to create session")
		err = errors.WithStack(constant.ErrUnknownSource)
		return
	}

	logRecordArgs := domain.LoginRecord{
		UserID:      user.ID,
		AccessToken: data.AccessToken,
		Status:      "login",
		IPAddress:   ipAddress,
		Type:        "web",
	}

	if logErr := q.CreateLoginRecord(ctx, logRecordArgs); logErr != nil {
		s.log.Error().Err(logErr).Msg("failed to create login record log")
	}

	if err = tx.Commit(); err != nil {
		s.log.Error().Err(err).Msg("error to commit transaction")
		err = errors.WithStack(constant.ErrUnknownSource)
	}

	return
}
