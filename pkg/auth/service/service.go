package service

import (
	authPresenter "be-dashboard-nba/api/presenter/auth"
	"be-dashboard-nba/constant"
	"be-dashboard-nba/internal/db"
	"be-dashboard-nba/internal/jwt"
	"be-dashboard-nba/internal/utils"
	"be-dashboard-nba/pkg/auth/helper"
	"be-dashboard-nba/pkg/auth/repository"
	"be-dashboard-nba/pkg/entities"
	"context"
	"database/sql"

	// "github.com/gofiber/fiber/v2/log" // Dihapus, kita akan gunakan s.log
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
)

type Service interface {
	LoginService(ctx context.Context, request authPresenter.LoginRequest, userAgent, iPAddress string) (data entities.Session, user entities.User, err error)
	LogoutService(ctx context.Context, claims *jwt.AccessTokenPayload, iPAddress string) (err error)
	AuthMeService(ctx context.Context, id string) (data entities.User, err error)
}

type service struct {
	db  db.DB
	log *zerolog.Logger
}

func NewService(db db.DB, log *zerolog.Logger) Service {
	return &service{
		db:  db,
		log: log,
	}
}

func (s *service) LoginService(
	ctx context.Context,
	request authPresenter.LoginRequest,
	userAgent,
	iPAddress string,
) (data entities.Session, user entities.User, err error) {

	logRepo := repository.NewRepo(s.db)

	tx, err := s.db.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		s.log.Error().Err(err).Msg("error to begin transaction")
		err = errors.WithStack(constant.ErrUnknownSource)
		return
	}

	q := repository.NewRepo(tx)

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
			logAttemptArgs := repository.LoginAttempParams{
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
		logAttemptArgs := repository.LoginAttempParams{
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

	err = q.CreateSessionQuery(ctx, repository.SessionParams{
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

	logRecordArgs := repository.LoginRecord{
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

func (s *service) AuthMeService(ctx context.Context, id string) (data entities.User, err error) {
	r := repository.NewRepo(s.db)
	data, err = r.ReadDetailUserByIdQuery(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			s.log.Warn().Str("id", id).Msg("user detail not found")
			err = constant.ErrUserIdNotFound
			return
		}
		s.log.Error().Err(err).Str("id", id).Msg("error reading user detail query")
		err = errors.WithStack(constant.ErrUnknownSource)
		return
	}

	return
}
