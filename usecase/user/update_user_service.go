package user

import (
	presenter "be-dashboard-nba/api/presenter/user"
	"be-dashboard-nba/constant"
	roleRepo "be-dashboard-nba/repository/role"
	userRepo "be-dashboard-nba/repository/user"
	"context"
	"database/sql"

	"github.com/pkg/errors"
)

func (s *service) UpdateUserService(
	ctx context.Context,
	request presenter.UpdateUserRequest,
	updatedBy string,
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
		}
	}()

	ru := userRepo.NewRepository(tx)
	rr := roleRepo.NewRepository(tx)

	existingUser, err := ru.ReadUserByIDQuery(ctx, userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			s.log.Warn().Str("id", userID).Msg("user detail not found")
			err = constant.ErrUserIdNotFound
			return
		}
		s.log.Error().Err(err).Str("id", userID).Msg("error reading user detail query")
		err = errors.WithStack(constant.ErrUnknownSource)
		return
	}

	existingRoleUser, err := rr.ReadRoleByIDQuery(ctx, request.RoleID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			s.log.Warn().Int("id", request.RoleID).Msg("role detail not found")
			err = constant.ErrRoleIdNotFound
			return
		}
		s.log.Error().Err(err).Int("id", request.RoleID).Msg("error reading role detail query")
		err = errors.WithStack(constant.ErrUnknownSource)
		return
	}

	params := request.ToParams(userID, updatedBy, existingUser)

	if existingUser.Role != existingRoleUser.Name {
		err = ru.UpdateUserRoleQuery(ctx, params.RoleID, params.ID)
		if err != nil {
			s.log.Error().Err(err).Interface("request_payload", request).Msg("error to update user role")
			err = errors.WithStack(constant.ErrUnknownSource)
			return
		}
	}

	err = ru.UpdateUserQuery(ctx, params)
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
