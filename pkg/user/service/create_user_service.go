package service

import (
	userPresenter "be-dashboard-nba/api/presenter/user"
	"be-dashboard-nba/constant"
	"be-dashboard-nba/internal/utils"
	roleRepo "be-dashboard-nba/pkg/role/repository"
	userRepo "be-dashboard-nba/pkg/user/repository"
	"context"
	"database/sql"

	"github.com/pkg/errors"
)

func (s *service) CreateUserService(ctx context.Context, req userPresenter.CreateUserRequest, userID string) (err error) {
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

	ur := userRepo.NewRepository(tx)
	rr := roleRepo.NewRepository(tx)

	hashedPassword, err := utils.GenerateHashPassword(req.Password)
	if err != nil {
		s.log.Error().Err(err).Msg("failed hash password")
		err = errors.WithStack(constant.ErrUnknownSource)
		return
	}

	params := req.ToParams(userID, hashedPassword)
	roleID := params.RoleID

	_, err = rr.ReadRoleByIDQuery(ctx, roleID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			s.log.Warn().Int("Id", roleID).Msg("role detail not found")
			err = constant.ErrRoleIdNotFound
			return
		}
		s.log.Error().Err(err).Int("id", roleID).Msg("error reading role detail query")
		err = errors.WithStack(constant.ErrUnknownSource)
		return
	}

	newUserID, err := ur.CreateUserQuery(ctx, params)
	if err != nil {
		s.log.Error().Err(err).Interface("request_payload", req).Msg("error to create user")
		err = errors.WithStack(constant.ErrUnknownSource)
		return
	}

	err = ur.CreateUserRoleQuery(ctx, params.RoleID, newUserID)
	if err != nil {
		s.log.Error().Err(err).Msg("error to create user role")
		err = errors.WithStack(constant.ErrUnknownSource)
		return
	}

	if err = tx.Commit(); err != nil {
		s.log.Error().Err(err).Msg("error to commit transaction")
		err = errors.WithStack(constant.ErrUnknownSource)
	}

	return
}
