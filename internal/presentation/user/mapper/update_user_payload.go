package mapper

import (
	appdto "github.com/ABA-Developer/go-backend-template/internal/application/user/dto"
	"github.com/ABA-Developer/go-backend-template/internal/presentation/user/presenter"
)

func ToUpdateUserPayload(req presenter.UpdateUserRequest) appdto.UpdateUserPayload {
	return appdto.UpdateUserPayload(req)
}
