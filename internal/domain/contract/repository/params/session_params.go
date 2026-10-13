package params

import "time"

type SessionParams struct {
	ID                    string
	UserID                string
	AccessToken           string
	AccessTokenExpiredAt  time.Time
	RefreshToken          string
	RefreshTokenExpiredAt time.Time
	IPAddress             string
	UserAgent             string
}
