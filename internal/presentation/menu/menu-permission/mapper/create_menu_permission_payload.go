package mapper

import (
	appdto "github.com/ABA-Developer/go-backend-template/internal/application/menu/menu-permission/dto"
	"github.com/ABA-Developer/go-backend-template/internal/presentation/menu/menu-permission/presenter"
)

func ToCreateMenuPermissionPayload(req presenter.CreateMenuPermissionRequest) appdto.CreateMenuPermissionPayload {
	return appdto.CreateMenuPermissionPayload(req)
}
