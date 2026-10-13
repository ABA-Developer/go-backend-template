package mapper

import (
	"github.com/ABA-Developer/go-backend-template/internal/application/menu/menu-permission/dto"
	repoParams "github.com/ABA-Developer/go-backend-template/internal/domain/contract/repository/params"
)

func ToCreateMenuPermissionParams(req *dto.CreateMenuPermissionPayload, createdBy string, menuID int) repoParams.CreateMenuPermissionParams {
	return repoParams.CreateMenuPermissionParams{
		Code:       req.Code,
		ActionName: req.ActionName,
		MenuID:     menuID,
		CreatedBy:  createdBy,
	}
}
