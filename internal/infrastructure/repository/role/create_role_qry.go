package repository

import (
	repoParams "github.com/ABA-Developer/go-backend-template/internal/domain/contract/repository/params"
	"context"
)

func (r *Query) CreateRoleQry(ctx context.Context, params repoParams.CreateRoleParams) (err error) {
	const stmt = `
		INSERT INTO app_role(
			code, name, description, created_by, created_at
		)
		VALUES ($1, $2, $3, $4, NOW())
	`

	_, err = r.db.ExecContext(ctx, stmt,
		params.Code,
		params.Name,
		params.Description,
		params.CreatedBy,
	)

	return
}
