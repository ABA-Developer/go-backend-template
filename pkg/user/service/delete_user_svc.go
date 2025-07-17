package service

import (
	"context"

	"github.com/gofiber/fiber/v2/log"

	"be-dashboard-nba/pkg/user/repository"
)

func (s *Service) DeleteUserService(
	ctx context.Context,
	id int64,
) (err error) {
	q := repository.NewQuery(s.db)

	err = q.DeleteUserQuery(ctx, id)
	if err != nil {
		log.WithContext(ctx).Error(err, "error delete user query", "id", id)
		return
	}

	return
}
