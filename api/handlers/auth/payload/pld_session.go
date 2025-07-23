package payload

import "be-dashboard-nba/internal/jwt"

type SessionPayload struct {
	SessionGUID string `json:"session_guid"`
	UserID      int64  `json:"user_id"`
	Role        string `json:"role"`
	UserAgent   string `json:"user_agent"`
	IPAddress   string `json:"ip_address"`
}

func (request *SessionPayload) ToAccessTokenRequest() (
	params jwt.AccessTokenPayload,
) {
	params = jwt.AccessTokenPayload{
		GUID:   request.SessionGUID,
		UserID: request.UserID,
		Role:   request.Role,
	}

	return
}

func (request *SessionPayload) ToRefreshTokenRequest() (
	params jwt.RefreshTokenPayload,
) {
	params = jwt.RefreshTokenPayload{
		GUID: request.SessionGUID,
	}

	return
}
