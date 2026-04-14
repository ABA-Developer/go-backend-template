package menu_permission

import (
	"be-dashboard-nba/constant"
	"be-dashboard-nba/internal/application/mapper"
	menuPermissionPresenter "be-dashboard-nba/internal/presentation/presenter/menu/menu-permission"
	"context"
	"database/sql"

	"github.com/pkg/errors"
)

func (s *useCase) CreateMenuPermissionUseCase(ctx context.Context, payload menuPermissionPresenter.CreateMenuPermissionRequest, userID string, menuID int) (err error) {

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

	rmp := s.newMenuPermissionRepo(tx)
	rm := s.newMenuRepo(tx)

	_, err = rm.ReadMenuByIDQuery(ctx, menuID)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			s.log.Warn().Int("id", menuID).Msg("menu detail not found for update")
			err = constant.ErrMenuIdNotFound
			return
		}
		s.log.Error().Err(err).Int("id", menuID).Msg("error reading menu detail query")
		err = errors.WithStack(constant.ErrUnknownSource)
		return
	}

	err = rmp.CreateMenuPermissionQuery(ctx, mapper.ToCreateMenuPermissionParams(&payload, userID, menuID))
	if err != nil {
		s.log.Error().Err(err).Interface("request_payload", payload).Msg("error to create menu permission")
		err = errors.WithStack(constant.ErrUnknownSource)
		return
	}

	if err = tx.Commit(); err != nil {
		s.log.Error().Err(err).Msg("error to commit transaction")
		err = errors.WithStack(constant.ErrUnknownSource)
	}
	return
}
