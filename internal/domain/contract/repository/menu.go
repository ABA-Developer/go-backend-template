package repository

import (
	"context"

	repoParams "github.com/ABA-Developer/go-backend-template/internal/domain/contract/repository/params"
	"github.com/ABA-Developer/go-backend-template/internal/domain/model"
)

type MenuRepository interface {
	ReadSidebarMenuQry(ctx context.Context, userID string) (data []model.Menu, err error)
	ReadListMenuQry(ctx context.Context, params repoParams.ReadListMenuParams) (data []model.Menu, err error)
	ReadCountMenuQry(ctx context.Context, params repoParams.ReadListMenuParams) (count int64, err error)
	ReadParentMenuQry(ctx context.Context) (data []model.Menu, err error)
	CreateMenuQry(ctx context.Context, params repoParams.CreateMenuParams) (err error)
	UpdateMenuQry(ctx context.Context, params repoParams.UpdateMenuParams) (err error)
	DeleteMenuQry(ctx context.Context, menuID int) (err error)
	ReadMenuByIDQry(ctx context.Context, menuID int) (data model.Menu, err error)
	ReadSortForGroupQry(ctx context.Context, group string) (int, error)
	ReadNextSortForParentQry(ctx context.Context, parentID int32) (int, error)
	UpdateMenuSortQry(ctx context.Context, params repoParams.UpdateMenuSortParams) (err error)
	ReadNextSortForParentAndGroupQry(ctx context.Context, parentID int32, group string) (int, error)
	CountMenuChildrenQry(ctx context.Context, menuID int) (int, error)
	UpdateChildrenGroupQry(ctx context.Context, parentID int, newGroup string) error
}
