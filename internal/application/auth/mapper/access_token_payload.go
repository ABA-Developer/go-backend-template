package mapper

import (
	"github.com/ABA-Developer/go-backend-template/internal/application/auth/dto"
	"github.com/ABA-Developer/go-backend-template/internal/application/jwt"
)

func ToAccessTokenPayload(request dto.SessionPayload) jwt.AccessTokenPayload {
	return jwt.AccessTokenPayload{SessionID: request.SessionID, UserID: request.UserID}
}
