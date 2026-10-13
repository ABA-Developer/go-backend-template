package mapper

import (
	appdto "github.com/ABA-Developer/go-backend-template/internal/application/menu/dto"
	"github.com/ABA-Developer/go-backend-template/internal/presentation/menu/presenter"
)

func ToReadMenuListPayload(req presenter.ReadMenuListRequest) appdto.ReadMenuListPayload {
	return appdto.ReadMenuListPayload(req)
}
