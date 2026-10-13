package usecase

import (
	"github.com/ABA-Developer/go-backend-template/constant"
	"context"
)

func (s *useCase) CheckPermissionUseCase(ctx context.Context, menuURL constant.MenuKey, userID string, permissionCode []string) (hasAccess bool, err error) {
	r := s.newAuthRepo(s.db)
	hasAccess, err = r.CheckPermissionQry(ctx, menuURL, userID, permissionCode)

	return
}
