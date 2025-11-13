package repository

import (
	"be-dashboard-nba/pkg/entities"
	"context"
)

type ReadRolesParams struct {
	SetSearch bool
	Search    string
	Order     string
	Limit     int
	Offset    int
}

func (r *repository) ReadRolesQuery(
	ctx context.Context,
	args ReadRolesParams,
) (data []entities.Role, err error) {
	const stmt = `
		SELECT 
			id, name, code ,description
		FROM
			app_role
		WHERE
			(CASE WHEN $1::bool THEN(
				name ILIKE $2
			) ELSE TRUE END)
		ORDER BY
			(CASE WHEN $3 = 'name ASC' THEN name END)
			ASC,
			(CASE WHEN $3 = 'name DESC' THEN name END)
			DESC,
			(CASE WHEN $3 = 'description ASC' THEN description END)
			ASC,
			(CASE WHEN $3 = 'description DESC' THEN description END)
			DESC,
			(CASE WHEN $3 = 'id ASC' THEN id END)
			ASC,
			(CASE WHEN $3 = 'id DESC' THEN id END)
			DESC
		LIMIT $4
		OFFSET $5
	`
	rows, err := r.db.QueryContext(ctx, stmt,
		args.SetSearch,
		args.Search,
		args.Order,
		args.Limit,
		args.Offset,
	)
	if err != nil {
		return
	}
	defer rows.Close()

	for rows.Next() {
		var r entities.Role

		if err = rows.Scan(
			&r.ID,
			&r.Name,
			&r.Code,
			&r.Description,
		); err != nil {
			return
		}

		data = append(data, r)
	}

	return
}

func (r *repository) ReadRolesCount(ctx context.Context,
	args ReadRolesParams) (count int, err error) {
	const stmt = `
			SELECT COUNT(*)
			FROM
				app_role
			WHERE
				(CASE WHEN $1::bool THEN(
				name ILIKE $2
				) ELSE TRUE END)
		`

	err = r.db.QueryRowContext(ctx, stmt,
		args.SetSearch,
		args.Search,
	).Scan(&count)

	return
}

func (r *repository) ReadRoleByIDQuery(ctx context.Context, roleID int) (data entities.Role, err error) {
	const stmt = `
		SELECT 
			id, name, code, description
		FROM
			app_role
		WHERE
			id = $1
	`
	err = r.db.QueryRowContext(ctx, stmt, roleID).Scan(
		&data.ID,
		&data.Name,
		&data.Code,
		&data.Description,
	)

	return
}

func (r *repository) ReadRoleAccess(ctx context.Context, roleID int) (data []entities.RoleAccessResponse, err error) {
	const stmt = `
	SELECT 
		am.id AS menu_id,
		am.name AS menu_name,
		amp.id AS permission_id,
		amp.action_name AS permission_name,
		amp.code AS permission_code,
		CASE 
			WHEN ara.role_id IS NOT NULL THEN TRUE
			ELSE FALSE
		END AS has_access
	FROM
		app_menu am
	JOIN 
		app_menu_permission amp
		ON am.id = amp.menu_id
	LEFT JOIN
		app_role_access ara
		ON amp.id = ara.menu_permission_id
		AND ara.role_id = $1
	ORDER BY 
		am.name ASC, amp.action_name ASC
`
	rows, err := r.db.QueryContext(ctx, stmt, roleID)
	if err != nil {
		return
	}
	defer rows.Close()

	for rows.Next() {
		var r entities.RoleAccessResponse

		if err = rows.Scan(
			&r.MenuID,
			&r.MenuName,
			&r.PermissionID,
			&r.PermissionName,
			&r.PermissionCode,
			&r.HasAccess,
		); err != nil {
			return
		}
		data = append(data, r)
	}
	return
}
