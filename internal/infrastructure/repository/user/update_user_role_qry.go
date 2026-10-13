package repository

import (
	"context"
)

func (r *Query) UpdateUserRoleQry(ctx context.Context, roleID int, userID string) (err error) {
	const stmt = `
		UPDATE app_user_role
		SET
			user_id = $1,
			role_id = $2
		WHERE user_id = $1
	`

	_, err = r.db.ExecContext(ctx, stmt, userID, roleID)

	return
}
