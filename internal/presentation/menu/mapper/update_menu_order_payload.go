package mapper

import (
	appdto "github.com/ABA-Developer/go-backend-template/internal/application/menu/dto"
	"github.com/ABA-Developer/go-backend-template/internal/presentation/menu/presenter"
)

func ToUpdateMenuOrderPayload(req presenter.UpdateMenuOrderRequest) appdto.UpdateMenuOrderPayload {
	return appdto.UpdateMenuOrderPayload(req)
}
