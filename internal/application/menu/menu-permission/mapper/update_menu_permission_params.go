package mapper

import (
	"github.com/ABA-Developer/go-backend-template/internal/application/menu/menu-permission/dto"
	repoParams "github.com/ABA-Developer/go-backend-template/internal/domain/contract/repository/params"
)

func ToUpdateMenuPermissionParams(req *dto.UpdateMenuPermissionPayload, updatedBy string, menuPermissionID int) repoParams.UpdateMenuPermissionParams {
	return repoParams.UpdateMenuPermissionParams{
		Code:             req.Code,
		ActionName:       req.ActionName,
		UpdatedBy:        updatedBy,
		MenuPermissionID: menuPermissionID,
	}
}
