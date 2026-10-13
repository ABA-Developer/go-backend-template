package repository

import (
	repoParams "github.com/ABA-Developer/go-backend-template/internal/domain/contract/repository/params"
	"context"
)

func (r *Query) DeleteRoleAccessQry(ctx context.Context, payload repoParams.UpdateRoleMenuPermission) (err error) {
	const stmt = `
		DELETE FROM app_role_access
		WHERE
			role_id = $1 AND menu_permission_id = $2
	`
	_, err = r.db.ExecContext(ctx, stmt,
		payload.RoleID,
		payload.MenuPermissionID,
	)
	return
}
