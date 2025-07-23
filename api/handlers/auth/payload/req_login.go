package payload

import "be-dashboard-nba/internal/utils"

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

func (request *LoginRequest) ToSessionPayload(userID int64, role, userAgent, iPAddress string) (
	params SessionPayload,
) {
	params = SessionPayload{
		SessionGUID: utils.GenerateUUID(),
		UserID:      userID,
		Role:        role,
		UserAgent:   userAgent,
		IPAddress:   iPAddress,
	}

	return
}
