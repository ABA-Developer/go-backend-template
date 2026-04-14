package auth

import (
	"be-dashboard-nba/constant"
	"be-dashboard-nba/internal/application/dto"
	"be-dashboard-nba/internal/application/usecase/auth/helper"
	"be-dashboard-nba/internal/application/utils"
	"be-dashboard-nba/internal/domain/model"
	authPresenter "be-dashboard-nba/internal/presentation/presenter/auth"
	"context"
	"database/sql"

	"github.com/pkg/errors"
)

func (s *useCase) LoginUseCase(
	ctx context.Context,
	request authPresenter.LoginRequest,
	userAgent,
	iPAddress string,
) (data model.Session, user model.User, err error) {

	logRepo := s.newAuthRepo(s.db)

	tx, err := s.db.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		s.log.Error().Err(err).Msg("error to begin transaction")
		err = errors.WithStack(constant.ErrUnknownSource)
		return
	}

	q := s.newAuthRepo(tx)

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

	user, err = q.ReadDetailUserByEmailQuery(ctx, request.Email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			logAttemptArgs := dto.LoginAttempParams{
				Email:     request.Email,
				IPAddress: iPAddress,
				Password:  request.Password,
			}
			if logErr := logRepo.CreateLoginAttemp(ctx, logAttemptArgs); logErr != nil {
				s.log.Error().Err(logErr).Msg("failed to create login attempt log (wrong email)")
			}
			err = constant.ErrWrongEmailOrPassword

			return
		}

		s.log.Error().Err(err).Str("email", request.Email).Msg("error to read user by email")
		err = errors.WithStack(constant.ErrUnknownSource)
		return
	}

	if err = utils.CompareHashPassword(request.Password, user.Password); err != nil {
		logAttemptArgs := dto.LoginAttempParams{
			Email:     request.Email,
			IPAddress: iPAddress,
			Password:  request.Password,
		}
		if logErr := logRepo.CreateLoginAttemp(ctx, logAttemptArgs); logErr != nil {
			s.log.Error().Err(logErr).Msg("failed to create login attempt log (wrong password)")
		}
		err = constant.ErrWrongEmailOrPassword
		return
	}

	sessionPayload := request.ToSessionPayload(user.ID, userAgent, iPAddress)

	data, err = helper.GenerateSessionModel(ctx, sessionPayload)
	if err != nil {
		s.log.Error().Err(err).Interface("payload", sessionPayload).Msg("error to generate session model")
		err = errors.WithStack(constant.ErrUnknownSource)
		return
	}

	err = q.CreateSessionQuery(ctx, dto.SessionParams{
		ID:                    data.ID,
		UserID:                data.UserID,
		AccessToken:           data.AccessToken,
		AccessTokenExpiredAt:  data.AccessTokenExpiredAt,
		RefreshToken:          data.RefreshToken,
		RefreshTokenExpiredAt: data.RefreshTokenExpiredAt,
		IPAddress:             data.IPAddress,
		UserAgent:             data.UserAgent,
	})
	if err != nil {
		s.log.Error().Err(err).Interface("session", data).Msg("error to create session")
		err = errors.WithStack(constant.ErrUnknownSource)
		return
	}

	logRecordArgs := dto.LoginRecord{
		UserID:      user.ID,
		AccessToken: data.AccessToken,
		Status:      "login",
		IPAddress:   iPAddress,
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
