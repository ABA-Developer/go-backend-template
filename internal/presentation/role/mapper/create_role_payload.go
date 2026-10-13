package mapper

import (
	appdto "github.com/ABA-Developer/go-backend-template/internal/application/role/dto"
	"github.com/ABA-Developer/go-backend-template/internal/presentation/role/presenter"
)

func ToCreateRolePayload(req presenter.CreateRoleRequest) appdto.CreateRolePayload {
	return appdto.CreateRolePayload(req)
}
