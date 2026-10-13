package repository

import (
	"context"
)

func (r *Query) IsUserEmailExistsQry(
	ctx context.Context,
	email string,
) (exists bool, err error) {
	const statement = `
		SELECT EXISTS (
			SELECT
				1
			FROM
				app_user
			WHERE
				email = $1
		)
	`

	err = r.db.QueryRowContext(ctx, statement, email).Scan(&exists)

	return
}
