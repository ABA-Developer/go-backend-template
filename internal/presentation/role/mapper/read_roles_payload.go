package mapper

import (
	appdto "github.com/ABA-Developer/go-backend-template/internal/application/role/dto"
	"github.com/ABA-Developer/go-backend-template/internal/presentation/role/presenter"
)

func ToReadRolesPayload(req presenter.ReadRolesRequest) appdto.ReadRolesPayload {
	return appdto.ReadRolesPayload(req)
}
