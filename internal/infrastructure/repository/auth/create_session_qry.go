package repository

import (
	repoParams "github.com/ABA-Developer/go-backend-template/internal/domain/contract/repository/params"
	"context"
)

func (r *Query) CreateSessionQry(
	ctx context.Context,
	args repoParams.SessionParams,
) (err error) {
	const statement = `
		INSERT INTO sessions (
			id, user_id,
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

	_, err = r.db.ExecContext(ctx, statement,
		args.ID,
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
