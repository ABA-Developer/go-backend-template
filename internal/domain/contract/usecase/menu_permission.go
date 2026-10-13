package usecase

import (
	"context"

	appdto "github.com/ABA-Developer/go-backend-template/internal/application/menu/menu-permission/dto"
	"github.com/ABA-Developer/go-backend-template/internal/domain/model"
)

type MenuPermissionUseCase interface {
	ReadMenuPermissionUseCase(ctx context.Context, args appdto.ReadMenuPermissionListPayload, MenuID int) (data model.MenuPermissionPaginationResponse, err error)
	CreateMenuPermissionUseCase(ctx context.Context, payload appdto.CreateMenuPermissionPayload, userID string, menuID int) (err error)
	UpdateMenuPermissionUseCase(ctx context.Context, payload appdto.UpdateMenuPermissionPayload, userID string, menuPermissionID int) (err error)
	ReadMenuPermissionDetailUseCase(ctx context.Context, MenuPermissionID int) (data model.MenuPermission, err error)
	DeleteMenuPermissionUseCase(ctx context.Context, menuPermissionID int) (err error)
}
