package usecase

import (
	"context"

	appdto "github.com/ABA-Developer/go-backend-template/internal/application/menu/dto"
	"github.com/ABA-Developer/go-backend-template/internal/domain/model"
)

type MenuUseCase interface {
	CreateMenuUseCase(ctx context.Context, request appdto.CreateMenuPayload, userID string) (err error)
	UpdateMenuUseCase(ctx context.Context, request appdto.UpdateMenuPayload, userID string, MenuID int) (err error)
	ReadListMenuUseCase(ctx context.Context, request appdto.ReadMenuListPayload) (data []model.Menu, err error)
	ReadSidebarMenuUseCase(ctx context.Context, userID string) (data []model.Menu, err error)
	DeleteMenuUseCase(ctx context.Context, menuID int) (err error)
	UpdateMenuOrderUseCase(ctx context.Context, request appdto.UpdateMenuOrderPayload, userID string) (err error)
	ReadMenuDetailUseCase(ctx context.Context, menuID int) (data model.Menu, err error)
	ReadMenuParentUseCase(ctx context.Context) (data []model.Menu, err error)
}
