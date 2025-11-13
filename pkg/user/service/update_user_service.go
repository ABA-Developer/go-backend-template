package service

import (
	presenter "be-dashboard-nba/api/presenter/user"
	"be-dashboard-nba/constant"
	"be-dashboard-nba/internal/utils"
	"be-dashboard-nba/pkg/user/repository"
	"context"
	"database/sql"

	"github.com/pkg/errors"
)

func (s *service) UpdateUserService(
	ctx context.Context,
	request presenter.UpdateUserRequest,
	userID string,
) (err error) {
	tx, err := s.db.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		s.log.Error().Err(err).Msg("error to begin transaction")
		err = errors.WithStack(constant.ErrUnknownSource)
		return
	}

	defer func() {
		if err != nil {
			if errRollback := tx.Rollback(); errRollback != nil {
				s.log.Error().Err(errRollback).AnErr("original_error", err).Msg("error to rollback transaction")
				err = errors.WithStack(constant.ErrUnknownSource)
				return
			}

			if !errors.Is(err, constant.ErrUnknownSource) {
				err = errors.WithStack(constant.ErrUnknownSource)
			}
		}
	}()

	r := repository.NewRepository(tx)

	var password string

	if request.Password != "" {
		password, err = utils.GenerateHashPassword(request.Password)
		if err != nil {
			s.log.Error().Err(err).Msg("error generate hash password")
			err = errors.WithStack(constant.ErrUnknownSource)
			return
		}
	}

	err = r.UpdateUserQuery(ctx, request.ToParams(userID, password))
	if err != nil {
		s.log.Error().Err(err).Interface("request_payload", request).Msg("error to update user")
		err = errors.WithStack(constant.ErrUnknownSource)
		return
	}

	if err = tx.Commit(); err != nil {
		s.log.Error().Err(err).Msg("error to commit transaction")
		err = errors.WithStack(constant.ErrUnknownSource)
	}

	return
}
