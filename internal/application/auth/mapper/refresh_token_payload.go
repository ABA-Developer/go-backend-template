package mapper

import (
	"github.com/ABA-Developer/go-backend-template/internal/application/auth/dto"
	"github.com/ABA-Developer/go-backend-template/internal/application/jwt"
)

func ToRefreshTokenPayload(request dto.SessionPayload) jwt.RefreshTokenPayload {
	return jwt.RefreshTokenPayload{SessionID: request.SessionID}
}
