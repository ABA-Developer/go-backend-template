package role

import (
	"be-dashboard-nba/api/presenter"
	rolePresenter "be-dashboard-nba/api/presenter/role"
	"be-dashboard-nba/constant"
	"be-dashboard-nba/usecase/entities"
	"be-dashboard-nba/repository/role"
	"context"
	"database/sql"
	"math"

	"github.com/pkg/errors"
)

func (s *service) ReadRoleAccessService(ctx context.Context, args rolePresenter.ReadRoleAccessesRequest, roleID int) (data entities.RoleAccessPaginationResponse, err error) {
	r := repository.NewRepository(s.db)
	params := args.ToParams(roleID)

	totalItems, err := r.ReadRoleAccessCount(ctx, params)
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

	_, err = r.ReadRoleByIDQuery(ctx, roleID)
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

	roleAccessData, err := r.ReadRoleAccessQuery(ctx, params)
	if err != nil {
		s.log.Error().Err(err).Int("id", roleID).Msg("error reading role detail query")
		err = errors.WithStack(constant.ErrUnknownSource)
		return
	}

	data = entities.RoleAccessPaginationResponse{
		Data:       roleAccessData,
		Pagination: pagination,
	}
	return
}
