package service

import (
	"be-dashboard-nba/constant"
	"be-dashboard-nba/pkg/menu/repository"
	"context"
	"database/sql"

	"github.com/pkg/errors"
)

func (s *service) DeleteMenuService(ctx context.Context, menuID int) (err error) {
	tx, err := s.db.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		s.log.Error().Err(err).Msg("error to begin transaction")
		err = errors.WithStack(constant.ErrUnknownSource)
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

	_, err = r.ReadMenuByIDQuery(ctx, menuID)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			s.log.Warn().Int("id", menuID).Msg("menu detail not found for delete")
			err = constant.ErrMenuIdNotFound
			return
		}
		s.log.Error().Err(err).Int("id", menuID).Msg("error reading menu detail query")
		err = errors.WithStack(constant.ErrUnknownSource)
		return
	}

	err = r.DeleteMenuQuery(ctx, menuID)

	if err != nil {
		s.log.Error().Err(err).Msg("error delete menu")
		err = errors.WithStack(constant.ErrUnknownSource)
		return

	}
	if err = tx.Commit(); err != nil {
		s.log.Error().Err(err).Msg("error to commit transaction")
		err = errors.WithStack(constant.ErrUnknownSource)
	}

	return
}
