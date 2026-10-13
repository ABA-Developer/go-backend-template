package repository

import (
	"context"
)

func (r *Query) CreateUserRoleQry(ctx context.Context, roleID int, userID string) (err error) {
	const stmt = `
		INSERT INTO app_user_role(
			user_id, role_id
		)
		VALUES (
			$1, $2
		)
	`

	_, err = r.db.ExecContext(ctx, stmt, userID, roleID)

	return
}
