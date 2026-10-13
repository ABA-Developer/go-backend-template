package repository

import (
	"context"
)

func (r *Query) DeleteUserQry(
	ctx context.Context,
	id string,
) (err error) {
	const statement = `
		DELETE FROM app_user
		WHERE
			id = $1
	`

	_, err = r.db.ExecContext(ctx, statement, id)

	return
}
