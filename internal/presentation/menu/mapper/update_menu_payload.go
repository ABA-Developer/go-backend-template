package mapper

import (
	appdto "github.com/ABA-Developer/go-backend-template/internal/application/menu/dto"
	"github.com/ABA-Developer/go-backend-template/internal/presentation/menu/presenter"
)

func ToUpdateMenuPayload(req presenter.UpdateMenuRequest) appdto.UpdateMenuPayload {
	return appdto.UpdateMenuPayload(req)
}
