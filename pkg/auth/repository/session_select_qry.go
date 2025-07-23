package repository

import (
	"context"

	"be-dashboard-nba/pkg/entities"
)

func (q *Query) ReadDetailSessionQuery(
	ctx context.Context,
	guid string,
) (data entities.Session, err error) {
	const statement = `
		SELECT
			guid, user_id,
			access_token, access_token_expired_at,
			refresh_token, refresh_token_expired_at
		FROM
			sessions
		WHERE
			guid = $1
	`

	err = q.db.QueryRowContext(ctx, statement, guid).Scan(
		&data.GUID,
		&data.UserID,
		&data.AccessToken,
		&data.AccessTokenExpiredAt,
		&data.RefreshToken,
		&data.RefreshTokenExpiredAt,
	)

	return
}
