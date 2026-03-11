package usecase

import (
	"be-dashboard-nba/constant"
	"context"
)

func (s *authUsecase) CheckPermissionUsecase(ctx context.Context, menuURL constant.MenuKey, userID string, permissionCode []string) (bool, error) {
	hasAccess, err := s.repo.CheckPermissionQuery(ctx, menuURL, userID, permissionCode)
	return hasAccess, err
}
