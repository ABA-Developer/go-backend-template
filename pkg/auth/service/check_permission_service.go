package service

import (
	"be-dashboard-nba/constant"
	"be-dashboard-nba/pkg/auth/repository"
	"context"
)

func (s *service) CheckPermissionService(ctx context.Context, menuURL constant.MenuKey, userID string, permissionCode []string) (hasAccess bool, err error) {
	r := repository.NewRepo(s.db)
	hasAccess, err = r.CheckPermissionQuery(ctx, menuURL, userID, permissionCode)

	return
}
