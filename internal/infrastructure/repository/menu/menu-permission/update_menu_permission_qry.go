package repository

import (
	repoParams "github.com/ABA-Developer/go-backend-template/internal/domain/contract/repository/params"
	"context"
)

func (r *Query) UpdateMenuPermissionQry(ctx context.Context, params repoParams.UpdateMenuPermissionParams) (err error) {
	const stmt = `
		UPDATE app_menu_permission
		SET
			code = $1,
			action_name = $2,
			updated_by = $3,
			updated_at = NOW()
		WHERE
			id = $4
	`
	_, err = r.db.ExecContext(ctx, stmt,
		params.Code,
		params.ActionName,
		params.UpdatedBy,
		params.MenuPermissionID,
	)

	return
}
