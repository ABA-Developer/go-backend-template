package repository

import (
	"context"
)

func (r *Query) DeleteSessionQry(
	ctx context.Context,
	id string,
) (err error) {

	const statement = `
		DELETE FROM sessions
		WHERE
			id = $1
	`

	_, err = r.db.ExecContext(ctx, statement, id)

	return
}
