package mapper

import (
	"be-dashboard-nba/internal/application/dto"
	menuPermissionPresenter "be-dashboard-nba/internal/presentation/presenter/menu/menu-permission"
)

func ToUpdateMenuPermissionParams(req *menuPermissionPresenter.UpdateMenuPermissionRequest, updatedBy string, menuPermissionID int) dto.UpdateMenuPermissionParams {
	return dto.UpdateMenuPermissionParams{
		Code:             req.Code,
		ActionName:       req.ActionName,
		UpdatedBy:        updatedBy,
		MenuPermissionID: menuPermissionID,
	}
}
