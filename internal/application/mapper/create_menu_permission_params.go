package mapper

import (
	"be-dashboard-nba/internal/application/dto"
	menuPermissionPresenter "be-dashboard-nba/internal/presentation/presenter/menu/menu-permission"
)

func ToCreateMenuPermissionParams(req *menuPermissionPresenter.CreateMenuPermissionRequest, createdBy string, menuID int) dto.CreateMenuPermissionParams {
	return dto.CreateMenuPermissionParams{
		Code:       req.Code,
		ActionName: req.ActionName,
		MenuID:     menuID,
		CreatedBy:  createdBy,
	}
}
