package mapper

import (
	appdto "github.com/ABA-Developer/go-backend-template/internal/application/menu/dto"
	"github.com/ABA-Developer/go-backend-template/internal/presentation/menu/presenter"
)

func ToCreateMenuPayload(req presenter.CreateMenuRequest) appdto.CreateMenuPayload {
	return appdto.CreateMenuPayload(req)
}
