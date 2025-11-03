package service

import (
	"context"

	"github.com/gofiber/fiber/v2/log"

	"be-dashboard-nba/pkg/entities"
	"be-dashboard-nba/pkg/user/repository"
)

func (s *Service) ReadDetailUserService(
	ctx context.Context,
	id int64,
) (data entities.User, err error) {
	q := repository.NewQuery(s.db)

	data, err = q.ReadDetailUserQuery(ctx, id)
	if err != nil {
		log.WithContext(ctx).Error(err, "error read user detail query", "id", id)
		return
	}

	return
}
