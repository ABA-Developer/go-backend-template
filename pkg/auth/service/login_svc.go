package service

import (
	"context"
	"database/sql"

	"github.com/gofiber/fiber/v2/log"
	"github.com/pkg/errors"

	"be-dashboard-nba/api/handlers/auth/payload"
	"be-dashboard-nba/constant"
	"be-dashboard-nba/internal/utils"
	"be-dashboard-nba/pkg/auth/helper"
	"be-dashboard-nba/pkg/auth/repository"
	"be-dashboard-nba/pkg/entities"
)

func (s *Service) LoginService(
	ctx context.Context,
	request payload.LoginRequest,
	userAgent, iPAddress string,
) (data entities.Session, user entities.User, err error) {
	tx, err := s.db.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		log.WithContext(ctx).Error(err, "error to begin transaction")
		err = errors.WithStack(constant.ErrUnknownSource)
		return
	}

	defer func() {
		if err != nil {
			if errRollback := tx.Rollback(); errRollback != nil {
				log.WithContext(ctx).Error(errRollback, "error to rollback transaction", "original_error", err)
				err = errors.WithStack(constant.ErrUnknownSource)
				return
			}
		}
	}()

	q := repository.NewQuery(tx)

	user, err = q.ReadDetailUserByEmailQuery(ctx, request.Email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			err = errors.New(constant.ErrAccountNotFound.Error())
			return
		}

		log.WithContext(ctx).Error(err, "error to read user by email", "email", request.Email)
		return
	}

	if err = utils.CompareHashPassword(request.Password, user.Password); err != nil {
		err = errors.New(constant.ErrPasswordIncorrect.Error())
		return
	}

	sessionPayload := request.ToSessionPayload(user.ID, user.Role, userAgent, iPAddress)

	data, err = helper.GenerateSessionModel(ctx, sessionPayload)
	if err != nil {
		log.WithContext(ctx).Error(err, "error to generate session model", "payload", sessionPayload)
		return
	}

	err = q.CreateSessionQuery(ctx, repository.SessionParams{
		GUID:                  data.GUID,
		UserID:                data.UserID,
		AccessToken:           data.AccessToken,
		AccessTokenExpiredAt:  data.AccessTokenExpiredAt,
		RefreshToken:          data.RefreshToken,
		RefreshTokenExpiredAt: data.RefreshTokenExpiredAt,
		IPAddress:             data.IPAddress,
		UserAgent:             data.UserAgent,
	})
	if err != nil {
		log.WithContext(ctx).Error(err, "error to create session", "session", data)
		return
	}

	if err = tx.Commit(); err != nil {
		log.WithContext(ctx).Error(err, "error to commit transaction")
		err = errors.WithStack(constant.ErrUnknownSource)
	}

	return
}
