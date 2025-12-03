package presenter

import "be-dashboard-nba/internal/utils"

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

func (request *LoginRequest) ToSessionPayload(userID string, userAgent, iPAddress string) (
	params SessionPayload,
) {
	params = SessionPayload{
		SessionID: utils.GenerateUUID(), 
		UserID:    userID,             
		UserAgent: userAgent,
		IPAddress: iPAddress,
	}
	return
}
