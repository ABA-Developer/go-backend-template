package service

import (
	"be-dashboard-nba/api/presenter"
	menuPermissionPresenter "be-dashboard-nba/api/presenter/menu_permission"
	"be-dashboard-nba/constant"
	"be-dashboard-nba/pkg/entities"
	menuRepository "be-dashboard-nba/pkg/menu/repository"
	menuPermissionRepository "be-dashboard-nba/pkg/menu_permission/repository"
	"context"
	"database/sql"
	"math"

	"github.com/pkg/errors"
)

func (s *service) ReadMenuPermissionListParams(
	ctx context.Context,
	args menuPermissionPresenter.ReadMenuPermissionListRequest,
	MenuID int,
) (data entities.MenuPermissionPaginationResponse, err error) {

	mpr := menuPermissionRepository.NewRepository(s.db)
	mr := menuRepository.NewRepository(s.db)

	_, err = mr.ReadMenuByIDQuery(ctx, MenuID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			s.log.Warn().Int("id", MenuID).Msg("menu detail not found")
			err = constant.ErrMenuIdNotFound
			return
		}
		s.log.Error().Err(err).Int("id", MenuID).Msg("error reading menu detail query")
		err = errors.WithStack(constant.ErrUnknownSource)
		return
	}

	params := args.ToParams(MenuID)

	totalItems, err := mpr.ReadMenuPermissionCount(ctx, params)
	if err != nil {
		s.log.Error().Err(err).Msg("error query get menu permission count")
		err = errors.WithStack(constant.ErrUnknownSource)
		return
	}

	totalPages := int(math.Ceil(float64(totalItems) / float64(args.Limit)))
	hasNext := args.Page < totalPages
	hasPrev := args.Page > 1

	pagination := presenter.Pagination{
		Page:       args.Page,
		PageSize:   args.Limit,
		TotalItems: totalItems,
		TotalPages: totalPages,
		HasNext:    hasNext,
		HasPrev:    hasPrev,
	}

	menuPermissionData, err := mpr.ReadMenuPermissionListQuery(ctx, params)
	if err != nil {
		s.log.Error().Err(err).Msg("error query read list menu permission")
		err = errors.WithStack(constant.ErrUnknownSource)
		return
	}

	data = entities.MenuPermissionPaginationResponse{
		Data:       menuPermissionData,
		Pagination: pagination,
	}

	return
}
