package mapper

import (
	appdto "github.com/ABA-Developer/go-backend-template/internal/application/role/dto"
	"github.com/ABA-Developer/go-backend-template/internal/presentation/role/presenter"
)

func ToUpdateRoleAccessPayload(req presenter.UpdateRoleAccessRequest) appdto.UpdateRoleAccessPayload {
	items := make([]appdto.UpdateRoleAccessPayloadItem, len(req.AccessItem))
	for index, item := range req.AccessItem {
		items[index] = appdto.UpdateRoleAccessPayloadItem{AccessID: item.AccessID, HasAccess: item.HasAccess}
	}
	return appdto.UpdateRoleAccessPayload{AccessItem: items}
}
