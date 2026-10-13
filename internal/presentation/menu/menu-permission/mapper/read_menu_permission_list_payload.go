package mapper

import (
	appdto "github.com/ABA-Developer/go-backend-template/internal/application/menu/menu-permission/dto"
	"github.com/ABA-Developer/go-backend-template/internal/presentation/menu/menu-permission/presenter"
)

func ToReadMenuPermissionListPayload(req presenter.ReadMenuPermissionListRequest) appdto.ReadMenuPermissionListPayload {
	return appdto.ReadMenuPermissionListPayload(req)
}
