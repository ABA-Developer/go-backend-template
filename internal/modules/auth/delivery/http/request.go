package http

import (
	"be-dashboard-nba/internal/core/jwt"
	"be-dashboard-nba/internal/core/utils"
)

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type SessionPayload struct {
	SessionID string `json:"session_id"`
	UserID    string `json:"user_id"`
	UserAgent string `json:"user_agent"`
	IPAddress string `json:"ip_address"`
}

func (request *LoginRequest) ToSessionPayload(userID string, userAgent, iPAddress string) SessionPayload {
	return SessionPayload{
		SessionID: utils.GenerateUUID(),
		UserID:    userID,
		UserAgent: userAgent,
		IPAddress: iPAddress,
	}
}

func (pld *SessionPayload) ToAccessTokenRequest() jwt.AccessTokenPayload {
	return jwt.AccessTokenPayload{
		SessionID: pld.SessionID,
		UserID:    pld.UserID,
	}
}

func (pld *SessionPayload) ToRefreshTokenRequest() jwt.RefreshTokenPayload {
	return jwt.RefreshTokenPayload{
		SessionID: pld.SessionID,
	}
}
