package repository

import (
	repoParams "github.com/ABA-Developer/go-backend-template/internal/domain/contract/repository/params"
	"context"
)

func (r *Query) UpdateMenuSortQry(
	ctx context.Context,
	params repoParams.UpdateMenuSortParams,
) (err error) {
	const stmt = `
		UPDATE app_menu
		SET
			sort = $1,
			"group" = $2,
			parent_id = $3,
			updated_by = $4,
			updated_at = NOW()
		WHERE
			id = $5
	`
	_, err = r.db.ExecContext(ctx, stmt,
		params.Sort,
		params.Group,
		params.ParentID,
		params.UpdatedBy,
		params.ID,
	)
	return
}
