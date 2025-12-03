package service

import (
	menuPermissionPresenter "be-dashboard-nba/api/presenter/menu_permission"
	"be-dashboard-nba/constant"
	"be-dashboard-nba/pkg/menu_permission/repository"
	"context"
	"database/sql"

	"github.com/pkg/errors"
)

func (s *service) UpdateMenuPermissionService(ctx context.Context, payload menuPermissionPresenter.UpdateMenuPermissionRequest, userID string, menuPermissionID int) (err error) {
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

	r := repository.NewRepository(s.db)

	_, err = r.ReadMenuPermissionByIdQuery(ctx, menuPermissionID)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			s.log.Warn().Int("id", menuPermissionID).Msg("menu permission detail not found for update")
			err = constant.ErrMenuPermissionIdNotFound
			return
		}
		s.log.Error().Err(err).Int("id", menuPermissionID).Msg("error reading menu permission detail query")
		err = errors.WithStack(constant.ErrUnknownSource)
		return
	}

	err = r.UpdateMenuPermissionQuery(ctx, payload.ToParams(userID, menuPermissionID))
	if err != nil {
		s.log.Error().Err(err).Interface("request_payload", payload).Msg("error to update menu permission")
		err = errors.WithStack(constant.ErrUnknownSource)
		return
	}

	if err = tx.Commit(); err != nil {
		s.log.Error().Err(err).Msg("error to commit transaction")
		err = errors.WithStack(constant.ErrUnknownSource)
	}

	return
}
