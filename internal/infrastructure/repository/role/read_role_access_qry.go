package repository

import (
	repoParams "github.com/ABA-Developer/go-backend-template/internal/domain/contract/repository/params"
	"github.com/ABA-Developer/go-backend-template/internal/domain/model"
	"context"
)

func (r *Query) ReadRoleAccessQry(
	ctx context.Context,
	args repoParams.ReadRoleAccessParams,
) (data []model.RoleAccessResponse, err error) {
	const stmt = `
		WITH filtered_menu AS (
			SELECT
				id, name
			FROM app_menu
			WHERE
				(CASE WHEN $2::bool THEN (
					name ILIKE $3
				) ELSE TRUE END)
			ORDER BY
				(CASE WHEN $4 = 'name ASC' THEN name END) ASC,
				(CASE WHEN $4 = 'name DESC' THEN name END) DESC,
				(CASE WHEN $4 = 'id ASC' THEN id END) ASC,
				(CASE WHEN $4 = 'id DESC' THEN id END) DESC,
				name ASC
			LIMIT $5
			OFFSET $6
		)
		SELECT
			ar.id AS role_id,
			ar.name AS role_name,
			fm.id AS menu_id,
			fm.name AS menu_name,
			amp.id AS permission_id,
			amp.action_name AS permission_name,
			amp.code AS permission_code,
			CASE
				WHEN ara.role_id IS NOT NULL THEN TRUE
				ELSE FALSE
			END AS has_access
		FROM filtered_menu fm
		JOIN app_menu_permission amp
			ON fm.id = amp.menu_id
		JOIN app_role ar
			ON ar.id = $1
		LEFT JOIN app_role_access ara
			ON amp.id = ara.menu_permission_id
			AND ara.role_id = $1
		ORDER BY fm.name ASC, amp.action_name ASC
	`

	rows, err := r.db.QueryContext(ctx, stmt,
		args.RoleID,
		args.SetSearch,
		args.Search,
		args.Order,
		args.Limit,
		args.Offset,
	)
	if err != nil {
		return
	}
	defer func() {
		_ = rows.Close()
	}()

	for rows.Next() {
		var res model.RoleAccessResponse
		if err = rows.Scan(
			&res.RoleID,
			&res.RoleName,
			&res.MenuID,
			&res.MenuName,
			&res.PermissionID,
			&res.PermissionName,
			&res.PermissionCode,
			&res.HasAccess,
		); err != nil {
			return
		}
		data = append(data, res)
	}

	return
}
