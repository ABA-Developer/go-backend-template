package mapper

import (
	"github.com/ABA-Developer/go-backend-template/internal/application/auth/dto"
	"github.com/ABA-Developer/go-backend-template/internal/application/utils"
)

func ToSessionPayload(userID, userAgent, iPAddress string) dto.SessionPayload {
	return dto.SessionPayload{
		SessionID: utils.GenerateUUID(),
		UserID:    userID,
		UserAgent: userAgent,
		IPAddress: iPAddress,
	}
}
