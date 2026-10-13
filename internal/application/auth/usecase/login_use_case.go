package usecase

import (
	"github.com/ABA-Developer/go-backend-template/constant"
	"github.com/ABA-Developer/go-backend-template/internal/application/auth/dto"
	"github.com/ABA-Developer/go-backend-template/internal/application/auth/mapper"
	"github.com/ABA-Developer/go-backend-template/internal/application/auth/usecase/helper"
	"github.com/ABA-Developer/go-backend-template/internal/application/utils"
	repoParams "github.com/ABA-Developer/go-backend-template/internal/domain/contract/repository/params"
	"github.com/ABA-Developer/go-backend-template/internal/domain/model"
	"context"
	"database/sql"

	"github.com/pkg/errors"
)

func (s *useCase) LoginUseCase(
	ctx context.Context,
	request dto.LoginPayload,
	userAgent,
	iPAddress string,
) (data model.Session, user model.User, err error) {

	logRepo := s.newAuthRepo(s.db)

	tx, err := s.db.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		log(ctx).Error().Err(err).Msg("error to begin transaction")
		err = errors.WithStack(constant.ErrUnknownSource)
		return
	}

	q := s.newAuthRepo(tx)

	defer func() {
		if err != nil {
			if !errors.Is(err, constant.ErrWrongEmailOrPassword) {
				if errRollback := tx.Rollback(); errRollback != nil {
					log(ctx).Error().Err(errRollback).AnErr("original_error", err).Msg("error to rollback transaction")
					err = errors.WithStack(constant.ErrUnknownSource)
					return
				}
				if !errors.Is(err, constant.ErrUnknownSource) {
					err = errors.WithStack(constant.ErrUnknownSource)
				}
			}
		}
	}()

	user, err = q.ReadDetailUserByEmailQry(ctx, request.Email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			logAttemptArgs := repoParams.LoginAttempParams{
				Email:     request.Email,
				IPAddress: iPAddress,
				Password:  request.Password,
			}
			if logErr := logRepo.CreateLoginAttempQry(ctx, logAttemptArgs); logErr != nil {
				log(ctx).Error().Err(logErr).Msg("failed to create login attempt log (wrong email)")
			}
			err = constant.ErrWrongEmailOrPassword

			return
		}

		log(ctx).Error().Err(err).Str("email", request.Email).Msg("error to read user by email")
		err = errors.WithStack(constant.ErrUnknownSource)
		return
	}

	if err = utils.CompareHashPassword(request.Password, user.Password); err != nil {
		logAttemptArgs := repoParams.LoginAttempParams{
			Email:     request.Email,
			IPAddress: iPAddress,
			Password:  request.Password,
		}
		if logErr := logRepo.CreateLoginAttempQry(ctx, logAttemptArgs); logErr != nil {
			log(ctx).Error().Err(logErr).Msg("failed to create login attempt log (wrong password)")
		}
		err = constant.ErrWrongEmailOrPassword
		return
	}

	sessionPayload := mapper.ToSessionPayload(user.ID, userAgent, iPAddress)

	data, err = helper.GenerateSessionModel(ctx, sessionPayload)
	if err != nil {
		log(ctx).Error().Err(err).Interface("payload", sessionPayload).Msg("error to generate session model")
		err = errors.WithStack(constant.ErrUnknownSource)
		return
	}

	err = q.CreateSessionQry(ctx, repoParams.SessionParams{
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
		log(ctx).Error().Err(err).Interface("session", data).Msg("error to create session")
		err = errors.WithStack(constant.ErrUnknownSource)
		return
	}

	logRecordArgs := repoParams.LoginRecord{
		UserID:      user.ID,
		AccessToken: data.AccessToken,
		Status:      "login",
		IPAddress:   iPAddress,
		Type:        "web",
	}

	if logErr := q.CreateLoginRecordQry(ctx, logRecordArgs); logErr != nil {
		log(ctx).Error().Err(logErr).Msg("failed to create login record log")
	}

	if err = tx.Commit(); err != nil {
		log(ctx).Error().Err(err).Msg("error to commit transaction")
		err = errors.WithStack(constant.ErrUnknownSource)
	}

	return
}
