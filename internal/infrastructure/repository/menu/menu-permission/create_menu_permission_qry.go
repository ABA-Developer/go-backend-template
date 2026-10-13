package repository

import (
	repoParams "github.com/ABA-Developer/go-backend-template/internal/domain/contract/repository/params"
	"context"
)

func (r *Query) CreateMenuPermissionQry(ctx context.Context, params repoParams.CreateMenuPermissionParams) (err error) {
	const stmt = `
		INSERT INTO app_menu_permission (
			code, action_name, menu_id, created_by, created_at
		)
		VALUES($1, $2, $3, $4, NOW())
	`

	_, err = r.db.ExecContext(ctx, stmt, params.Code, params.ActionName, params.MenuID, params.CreatedBy)

	return
}
