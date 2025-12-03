package service

import (
	"be-dashboard-nba/api/presenter"
	userPresenter "be-dashboard-nba/api/presenter/user"
	"be-dashboard-nba/constant"
	"be-dashboard-nba/pkg/entities"
	"be-dashboard-nba/pkg/user/repository"
	"context"
	"fmt"
	"math"

	"github.com/pkg/errors"
)

func (s *service) ReadUsersService(ctx context.Context, req userPresenter.ReadUserRequest) (data entities.UserPaginationResponse, err error) {
	r := repository.NewRepository(s.db)
	params := req.ToParams()
	totalItems, err := r.ReadCountUserQuery(ctx, params)
	if err != nil {
		s.log.Error().Err(err).Msg("error query get user count")
		err = errors.WithStack(constant.ErrUnknownSource)
		return
	}
	fmt.Printf("Page: %d", req.Page)
	totalPages := int(math.Ceil(float64(totalItems) / float64(req.Limit)))
	fmt.Printf("total pages: %d", totalPages)
	hasNext := req.Page < totalPages
	hasPrev := req.Page > 1

	pagination := presenter.Pagination{
		Page:       req.Page,
		PageSize:   req.Limit,
		TotalItems: totalItems,
		TotalPages: totalPages,
		HasNext:    hasNext,
		HasPrev:    hasPrev,
	}

	userData, err := r.ReadUsersQuery(ctx, params)
	if err != nil {
		s.log.Error().Err(err).Msg("error query read roles")
		err = errors.WithStack(constant.ErrUnknownSource)
		return
	}

	data = entities.UserPaginationResponse{
		Data:       userData,
		Pagination: pagination,
	}

	return

}
