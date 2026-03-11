package repository

import (
	"context"
	"database/sql"
)

type CreateRolePayload struct {
	Code        string
	Name        string
	Description sql.NullString
	CreatedBy   string
}

func (r *repository) CreateRoleQuery(ctx context.Context, payload CreateRolePayload) (err error) {
	const stmt = `
		INSERT INTO app_role(
			code, name, description, created_by, created_at
		)
		VALUES ($1, $2, $3, $4, NOW())
	`

	_, err = r.db.ExecContext(ctx, stmt,
		payload.Code,
		payload.Name,
		payload.Description,
		payload.CreatedBy,
	)

	return
}

type UpdateRolePayload struct {
	RoleID      int
	Code        string
	Name        string
	Description sql.NullString
	UpdatedBy   string
}

func (r *repository) UpdateRoleQuery(ctx context.Context, payload UpdateRolePayload) (err error) {
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
		payload.Code,
		payload.Name,
		payload.Description,
		payload.UpdatedBy,
		payload.RoleID,
	)
	return
}

func (r *repository) DeleteRoleQuery(ctx context.Context, roleID int) (err error) {
	const stmt = `
		DELETE FROM app_role
		WHERE
			id = $1
	`
	_, err = r.db.ExecContext(ctx, stmt, roleID)

	return
}

type UpdateRoleMenuPermission struct {
	MenuPermissionID int
	RoleID           int
}

func (r *repository) CreateRoleAccess(ctx context.Context, payload UpdateRoleMenuPermission) (err error) {
	const stmt = `
		INSERT INTO app_role_access(
			role_id, menu_permission_id
		)
		VALUES ($1, $2)
		ON CONFLICT (role_id, menu_permission_id) DO NOTHING;
	`
	_, err = r.db.ExecContext(ctx, stmt,
		payload.RoleID,
		payload.MenuPermissionID,
	)
	return
}

func (r *repository) DeleteRoleAccess(ctx context.Context, payload UpdateRoleMenuPermission) (err error) {
	const stmt = `
		DELETE FROM app_role_access
		WHERE 
			role_id = $1 AND menu_permission_id = $2		
	`
	_, err = r.db.ExecContext(ctx, stmt,
		payload.RoleID,
		payload.MenuPermissionID,
	)
	return
}
