package repository

import "context"

type CreateMenuPermissionPayload struct {
	Code       string
	ActionName string
	MenuID     int
	CreatedBy  string
}

func (r *repository) CreateMenuPermissionQuery(ctx context.Context, payload CreateMenuPermissionPayload) (err error) {
	const stmt = `
		INSERT INTO app_menu_permission (
			code, action_name, menu_id, created_by, created_at
		)
		VALUES($1, $2, $3, $4, NOW())
	`

	_, err = r.db.ExecContext(ctx, stmt, payload.Code, payload.ActionName, payload.MenuID, payload.CreatedBy)

	return
}

type UpdateMenuPermissionPayload struct {
	MenuPermissionID int
	Code             string
	ActionName       string
	UpdatedBy        string
}

func (r *repository) UpdateMenuPermissionQuery(ctx context.Context, payload UpdateMenuPermissionPayload) (err error) {
	const stmt = `
		UPDATE app_menu_permission
		SET
			code = $1,
			action_name = $2,
			updated_by = $3,
			updated_at = NOW()
		WHERE
			id = $4
	`
	_, err = r.db.ExecContext(ctx, stmt,
		payload.Code,
		payload.ActionName,
		payload.UpdatedBy,
		payload.MenuPermissionID,
	)

	return
}

func (r *repository) DeleteMenuPermissionQuery(ctx context.Context, MenuPermissionID int) (err error) {
	const stmt = `
		DELETE FROM app_menu_permission
		WHERE
			id = $1
	`
	_, err = r.db.ExecContext(ctx, stmt, MenuPermissionID)

	return
}
