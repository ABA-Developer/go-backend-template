package repository

import (
	repoParams "github.com/ABA-Developer/go-backend-template/internal/domain/contract/repository/params"
	"context"
)

func (r *Query) UpdateRoleQry(ctx context.Context, params repoParams.UpdateRoleParams) (err error) {
	const stmt = `
		UPDATE app_role
		SET
			code = $1,
			name = $2,
			description = COALESCE($3, description),
			updated_by = $4,
			updated_at = NOW()
		WHERE
			id = $5
	`
	_, err = r.db.ExecContext(ctx, stmt,
		params.Code,
		params.Name,
		params.Description,
		params.UpdatedBy,
		params.RoleID,
	)
	return
}
