package repository

import (
	repoParams "github.com/ABA-Developer/go-backend-template/internal/domain/contract/repository/params"
	"context"
)

func (r *Query) CreateMenuQry(
	ctx context.Context,
	params repoParams.CreateMenuParams,
) (err error) {
	const stmt = `
		INSERT INTO app_menu (
			parent_id, name, description, url, sort, "group", icon,
			active, display, created_by, created_at
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, NOW())

	`
	_, err = r.db.ExecContext(ctx, stmt,
		params.ParentID, params.Name, params.Description, params.URL, params.Sort,
		params.Group, params.Icon, params.Active, params.Display, params.CreatedBy,
	)
	return
}
