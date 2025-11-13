package repository

import (
	"context"
	"database/sql"
)

type CreateMenuParams struct {
	ParentID    sql.NullInt32
	Name        string
	Description sql.NullString
	URL         sql.NullString
	Sort        int
	Group       string
	Icon        sql.NullString
	Active      bool
	Display     bool
	CreatedBy   string
}

func (r *repository) CreateMenuQuery(
	ctx context.Context,
	params CreateMenuParams,
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

type UpdateMenuParams struct {
	ID          int
	ParentID    sql.NullInt32
	Name        string
	Description sql.NullString
	URL         sql.NullString
	Group       string
	Icon        sql.NullString
	Active      bool
	Display     bool
	UpdatedBy   sql.NullString
}

func (r *repository) UpdateMenuQuery(
	ctx context.Context,
	params UpdateMenuParams,
) (err error) {
	const stmt = `
		UPDATE app_menu
		SET
			name = $3, 
			sort = $6, 
			parent_id = COALESCE($2, parent_id),
			description = COALESCE($4, description),
			url = COALESCE($5, url),
			"group" = COALESCE($6, "group"), 
			icon = COALESCE($7, icon),      
			active = COALESCE($8, active),     
			display = COALESCE($9, display), 

			updated_by = $10, -- $11 menjadi $10
			updated_at = NOW()
		WHERE
			id = $1;
	`
	_, err = r.db.ExecContext(ctx, stmt,
		params.ID, params.ParentID, params.Name, params.Description, params.URL,
		params.Group, params.Icon, params.Active, params.Display, params.UpdatedBy,
	)
	return
}

func (r *repository) DeleteMenuQuery(ctx context.Context, menuID int) (err error) {
	const stmt = `
	DELETE FROM app_menu
	WHERE 
		id = $1
`

	_, err = r.db.ExecContext(ctx, stmt, menuID)

	return
}

type UpdateMenuSortParams struct {
	ID        int
	Sort      int
	UpdatedBy string
}

func (r *repository) UpdateMenuSortQuery(
	ctx context.Context,
	params UpdateMenuSortParams,
) (err error) {
	const stmt = `
		UPDATE app_menu
		SET 
			sort = $1,
			updated_by = $2,
			updated_at = NOW()
		WHERE 
			id = $3
	`
	// Gunakan ExecContext karena ini adalah UPDATE
	_, err = r.db.ExecContext(ctx, stmt,
		params.Sort,
		params.UpdatedBy,
		params.ID,
	)
	return
}
