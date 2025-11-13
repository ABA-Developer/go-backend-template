package service

import (
	"be-dashboard-nba/api/presenter"
	rolePresenter "be-dashboard-nba/api/presenter/role"
	"be-dashboard-nba/constant"
	"be-dashboard-nba/pkg/entities"
	"be-dashboard-nba/pkg/role/repository"
	"context"
	"math"

	"github.com/pkg/errors"
)

func (s *service) ReadRolesService(
	ctx context.Context,
	args rolePresenter.ReadRolesRequest,
) (data entities.RolePaginationResponse, err error) {
	r := repository.NewRepository(s.db)

	params := args.ToParams()

	totalItems, err := r.ReadRolesCount(ctx, params)
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

	rolesData, err := r.ReadRolesQuery(ctx, params)
	if err != nil {
		s.log.Error().Err(err).Msg("error query read list menu permission")
		err = errors.WithStack(constant.ErrUnknownSource)
		return
	}

	data = entities.RolePaginationResponse{
		Data:       rolesData,
		Pagination: pagination,
	}
	return
}
