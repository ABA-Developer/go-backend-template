package repository

import (
	"github.com/ABA-Developer/go-backend-template/internal/domain/model"
	"context"
)

func (r *Query) ReadDetailSessionQry(
	ctx context.Context,
	id string,
) (data model.Session, err error) {
	const statement = `
		SELECT
			id, user_id,
			access_token, access_token_expired_at,
			refresh_token, refresh_token_expired_at
		FROM
			sessions
		WHERE
			id = $1
	`

	err = r.db.QueryRowContext(ctx, statement, id).Scan(
		&data.ID,
		&data.UserID,
		&data.AccessToken,
		&data.AccessTokenExpiredAt,
		&data.RefreshToken,
		&data.RefreshTokenExpiredAt,
	)

	return
}
