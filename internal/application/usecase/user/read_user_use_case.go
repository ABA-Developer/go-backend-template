package user

import (
	"be-dashboard-nba/constant"
	"be-dashboard-nba/internal/application/mapper"
	"be-dashboard-nba/internal/domain/model"
	"be-dashboard-nba/internal/presentation/presenter"
	userPresenter "be-dashboard-nba/internal/presentation/presenter/user"
	"context"
	"fmt"
	"math"

	"github.com/pkg/errors"
)

func (s *useCase) ReadUsersUseCase(ctx context.Context, req userPresenter.ReadUserRequest) (data model.UserPaginationResponse, err error) {
	r := s.newUserRepo(s.db)
	params := mapper.ToReadListUserParams(&req)
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

	data = model.UserPaginationResponse{
		Data:       userData,
		Pagination: pagination,
	}

	return

}
