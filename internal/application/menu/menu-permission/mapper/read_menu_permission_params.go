package mapper

import (
	"github.com/ABA-Developer/go-backend-template/constant"
	"github.com/ABA-Developer/go-backend-template/internal/application/menu/menu-permission/dto"
	repoParams "github.com/ABA-Developer/go-backend-template/internal/domain/contract/repository/params"
)

func ToReadMenuPermissionParams(req *dto.ReadMenuPermissionListPayload, menuID int) repoParams.ReadMenuPermissionParams {
	req.Init()

	if req.Limit <= 0 {
		req.Limit = constant.DefaultLimit
	}

	if req.Page <= 0 {
		req.Page = constant.DefaultPage
	}

	return repoParams.ReadMenuPermissionParams{
		SetSearch: req.SetSearch,
		Search:    req.Search,
		Order:     req.Order,
		Limit:     req.Limit,
		Offset:    req.Offset,
		MenuID:    menuID,
	}
}
