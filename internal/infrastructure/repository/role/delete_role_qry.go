package repository

import (
	"context"
)

func (r *Query) DeleteRoleQry(ctx context.Context, roleID int) (err error) {
	const stmt = `
		DELETE FROM app_role
		WHERE
			id = $1
	`
	_, err = r.db.ExecContext(ctx, stmt, roleID)

	return
}
