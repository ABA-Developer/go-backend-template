package repository

import (
	"context"
	"time"
)

type SessionParams struct {
	GUID                  string
	UserID                int64
	AccessToken           string
	AccessTokenExpiredAt  time.Time
	RefreshToken          string
	RefreshTokenExpiredAt time.Time
	IPAddress             string
	UserAgent             string
}

func (q *Query) CreateSessionQuery(
	ctx context.Context,
	args SessionParams,
) (err error) {
	const statement = `
		INSERT INTO sessions (
			guid, user_id,
			access_token, access_token_expired_at,
			refresh_token, refresh_token_expired_at,
			ip_address, user_agent,
			created_at
		)
		VALUES (
			$1, $2,
			$3, $4,
			$5, $6,
			$7, $8,
			(now() at time zone 'UTC')::TIMESTAMP
		)
	`

	_, err = q.db.ExecContext(ctx, statement,
		args.GUID,
		args.UserID,
		args.AccessToken,
		args.AccessTokenExpiredAt,
		args.RefreshToken,
		args.RefreshTokenExpiredAt,
		args.IPAddress,
		args.UserAgent,
	)

	return
}

func (q *Query) UpdateSessionQuery(
	ctx context.Context,
	args SessionParams,
) (err error) {
	const statement = `
		UPDATE
			sessions
		SET
			access_token = $2,
			access_token_expired_at = $3,
			refresh_token = $4,
			refresh_token_expired_at = $5,
			updated_at = (now() at time zone 'UTC')::TIMESTAMP
		WHERE
			guid = $1
	`

	_, err = q.db.ExecContext(ctx, statement,
		args.GUID,
		args.AccessToken,
		args.AccessTokenExpiredAt,
		args.RefreshToken,
		args.RefreshTokenExpiredAt,
	)

	return
}

func (q *Query) DeleteSessionQuery(
	ctx context.Context,
	guid string,
) (err error) {
	const statement = `
		DELETE FROM sessions
		WHERE
			guid = $1
	`

	_, err = q.db.ExecContext(ctx, statement, guid)

	return
}
