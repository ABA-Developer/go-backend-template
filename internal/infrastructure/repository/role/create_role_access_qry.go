package repository

import (
	repoParams "github.com/ABA-Developer/go-backend-template/internal/domain/contract/repository/params"
	"context"
)

func (r *Query) CreateRoleAccessQry(ctx context.Context, payload repoParams.UpdateRoleMenuPermission) (err error) {
	const stmt = `
		INSERT INTO app_role_access(
			role_id, menu_permission_id
		)
		VALUES ($1, $2)
		ON CONFLICT (role_id, menu_permission_id) DO NOTHING;
	`
	_, err = r.db.ExecContext(ctx, stmt,
		payload.RoleID,
		payload.MenuPermissionID,
	)
	return
}
