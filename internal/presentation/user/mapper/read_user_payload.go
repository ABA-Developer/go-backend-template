package mapper

import (
	appdto "github.com/ABA-Developer/go-backend-template/internal/application/user/dto"
	"github.com/ABA-Developer/go-backend-template/internal/presentation/user/presenter"
)

func ToReadUserPayload(req presenter.ReadUserRequest) appdto.ReadUserPayload {
	return appdto.ReadUserPayload(req)
}
