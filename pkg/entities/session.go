package entities

import (
	"database/sql"
	"time"
)

type Session struct {
	GUID                  string       `db:"guid"`
	UserID                int64        `db:"user_id"`
	AccessToken           string       `db:"access_token"`
	AccessTokenExpiredAt  time.Time    `db:"access_token_expired_at"`
	RefreshToken          string       `db:"refresh_token"`
	RefreshTokenExpiredAt time.Time    `db:"refresh_token_expired_at"`
	IPAddress             string       `db:"ip_address"`
	UserAgent             string       `db:"user_agent"`
	CreatedAt             time.Time    `db:"created_at"`
	UpdatedAt             sql.NullTime `db:"updated_at"`
}
