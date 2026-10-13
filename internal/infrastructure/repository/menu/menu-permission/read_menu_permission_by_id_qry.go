package repository

import (
	"github.com/ABA-Developer/go-backend-template/internal/domain/model"
	"context"
)

func (r *Query) ReadMenuPermissionByIdQry(ctx context.Context, MenuPermissionID int) (data model.MenuPermission, err error) {
	const stmt = `
		SELECT
			id, menu_id, code, action_name
		FROM
			app_menu_permission
		WHERE
			id = $1
	`
	err = r.db.QueryRowContext(ctx, stmt, MenuPermissionID).Scan(
		&data.ID,
		&data.MenuID,
		&data.Code,
		&data.ActionName,
	)

	return
}
