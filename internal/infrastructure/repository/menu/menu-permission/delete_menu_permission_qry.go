package repository

import (
	"context"
)

func (r *Query) DeleteMenuPermissionQry(ctx context.Context, MenuPermissionID int) (err error) {
	const stmt = `
		DELETE FROM app_menu_permission
		WHERE
			id = $1
	`
	_, err = r.db.ExecContext(ctx, stmt, MenuPermissionID)

	return
}
