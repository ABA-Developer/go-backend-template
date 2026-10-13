package mapper

import (
	"github.com/ABA-Developer/go-backend-template/internal/application/auth/dto"
	"github.com/ABA-Developer/go-backend-template/internal/presentation/auth/presenter"
)

func ToLoginPayload(req presenter.LoginRequest) dto.LoginPayload {
	return dto.LoginPayload(req)
}
