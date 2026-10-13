package mapper

import (
	appdto "github.com/ABA-Developer/go-backend-template/internal/application/menu/menu-permission/dto"
	"github.com/ABA-Developer/go-backend-template/internal/presentation/menu/menu-permission/presenter"
)

func ToUpdateMenuPermissionPayload(req presenter.UpdateMenuPermissionRequest) appdto.UpdateMenuPermissionPayload {
	return appdto.UpdateMenuPermissionPayload(req)
}
