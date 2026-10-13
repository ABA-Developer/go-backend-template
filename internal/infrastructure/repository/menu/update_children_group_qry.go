package repository

import (
	"context"
)

func (r *Query) UpdateChildrenGroupQry(ctx context.Context, parentID int, newGroup string) error {
	const stmt = `UPDATE app_menu SET "group" = $1, updated_at = NOW() WHERE parent_id = $2`
	_, err := r.db.ExecContext(ctx, stmt, newGroup, parentID)
	return err
}
