package mapper

import (
	appdto "github.com/ABA-Developer/go-backend-template/internal/application/role/dto"
	"github.com/ABA-Developer/go-backend-template/internal/presentation/role/presenter"
)

func ToReadRoleAccessesPayload(req presenter.ReadRoleAccessesRequest) appdto.ReadRoleAccessesPayload {
	return appdto.ReadRoleAccessesPayload(req)
}
