package repository

import (
	"context"
)

func (r *Query) DeleteMenuQry(ctx context.Context, menuID int) (err error) {
	const stmt = `
	DELETE FROM app_menu
	WHERE
		id = $1
`

	_, err = r.db.ExecContext(ctx, stmt, menuID)

	return
}
