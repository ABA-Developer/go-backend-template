package role

import (
	"be-dashboard-nba/constant"
	"be-dashboard-nba/internal/application/mapper"
	"be-dashboard-nba/internal/domain/model"
	"be-dashboard-nba/internal/presentation/presenter"
	rolePresenter "be-dashboard-nba/internal/presentation/presenter/role"
	"context"
	"math"

	"github.com/pkg/errors"
)

func (s *useCase) ReadRolesUseCase(
	ctx context.Context,
	args rolePresenter.ReadRolesRequest,
) (data model.RolePaginationResponse, err error) {
	r := s.newRoleRepo(s.db)

	params := mapper.ToReadRolesParams(&args)

	totalItems, err := r.ReadRolesCount(ctx, params)
	if err != nil {
		s.log.Error().Err(err).Msg("error query get role count")
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
		s.log.Error().Err(err).Msg("error query read roles")
		err = errors.WithStack(constant.ErrUnknownSource)
		return
	}

	data = model.RolePaginationResponse{
		Data:       rolesData,
		Pagination: pagination,
	}
	return
}
