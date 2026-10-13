package repository

import (
	repoParams "github.com/ABA-Developer/go-backend-template/internal/domain/contract/repository/params"
	"context"
)

func (r *Query) UpdateSessionQry(
	ctx context.Context,
	args repoParams.SessionParams,
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
			id = $1
	`

	_, err = r.db.ExecContext(ctx, statement,
		args.ID,
		args.AccessToken,
		args.AccessTokenExpiredAt,
		args.RefreshToken,
		args.RefreshTokenExpiredAt,
	)

	return
}
