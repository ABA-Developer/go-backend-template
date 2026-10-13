package repository

import (
	"context"
)

func (r *Query) ReadNextSortForParentQry(ctx context.Context, parentID int32) (int, error) {
	var nextSort int
	const stmt = `
		SELECT COALESCE(MAX(sort), -1) + 1
		FROM app_menu
		WHERE parent_id = $1
	`

	err := r.db.QueryRowContext(ctx, stmt, parentID).Scan(&nextSort)
	return nextSort, err
}
