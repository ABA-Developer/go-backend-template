package service

import (
	"be-dashboard-nba/constant"
	"be-dashboard-nba/pkg/role/repository"
	"context"
	"database/sql"

	"github.com/pkg/errors"
)

func (s *service) DeleteRoleService(ctx context.Context, roleID int) (err error) {

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

	r := repository.NewRepository(tx)

	_, err = r.ReadRoleByIDQuery(ctx, roleID)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			s.log.Warn().Int("id", roleID).Msg("menu detail not found for update")
			err = constant.ErrRoleIdNotFound
			return
		}
		s.log.Error().Err(err).Int("id", roleID).Msg("error reading menu detail query")
		err = errors.WithStack(constant.ErrUnknownSource)
		return
	}

	err = r.DeleteRoleQuery(ctx, roleID)
	if err != nil {
		s.log.Error().Err(err).Msg("error to delete role")
		err = errors.WithStack(constant.ErrUnknownSource)
		return
	}

	if err = tx.Commit(); err != nil {
		s.log.Error().Err(err).Msg("error to commit transaction")
		err = errors.WithStack(constant.ErrUnknownSource)
	}
	return
}
