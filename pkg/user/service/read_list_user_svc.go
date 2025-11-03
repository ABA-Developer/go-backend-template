package service

import (
	"context"

	"github.com/gofiber/fiber/v2/log"

	"be-dashboard-nba/api/handlers/user/payload"
	"be-dashboard-nba/pkg/entities"
	"be-dashboard-nba/pkg/user/repository"
)

func (s *Service) ReadListUserService(
	ctx context.Context,
	request payload.ReadUserListRequest,
) (data []entities.User, totalData int64, err error) {
	q := repository.NewQuery(s.db)

	data, err = q.ReadListUserQuery(ctx, request.ToParams())
	if err != nil {
		log.WithContext(ctx).Error(err, "error read user list query", "request", request)
		return
	}

	totalData, err = q.GetCountUserQuery(ctx, request.ToParams())
	if err != nil {
		log.WithContext(ctx).Error(err, "error get count user query", "request", request)
		return
	}

	return
}
