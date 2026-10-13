package mapper

import (
	appdto "github.com/ABA-Developer/go-backend-template/internal/application/role/dto"
	"github.com/ABA-Developer/go-backend-template/internal/presentation/role/presenter"
)

func ToUpdateRolePayload(req presenter.UpdateRoleRequest) appdto.UpdateRolePayload {
	return appdto.UpdateRolePayload(req)
}
