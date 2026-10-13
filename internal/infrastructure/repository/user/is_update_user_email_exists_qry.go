package repository

import (
	"context"
)

func (r *Query) IsUpdateUserEmailExistsQry(
	ctx context.Context,
	email, id string,
) (exists bool, err error) {
	statement := `
		SELECT EXISTS (
			SELECT
				1
			FROM
				app_user
			WHERE
				email = $1
				AND id != $2
		)
	`

	err = r.db.QueryRowContext(ctx, statement, email, id).Scan(&exists)

	return
}
