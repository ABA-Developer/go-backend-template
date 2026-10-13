package mapper

import (
	appdto "github.com/ABA-Developer/go-backend-template/internal/application/user/dto"
	"github.com/ABA-Developer/go-backend-template/internal/presentation/user/presenter"
)

func ToCreateUserPayload(req presenter.CreateUserRequest) appdto.CreateUserPayload {
	return appdto.CreateUserPayload(req)
}
