package repository

import (
	"context"
	"database/sql"
	"fmt"

	"be-dashboard-nba/constant"
	"be-dashboard-nba/internal/core/db"
	"be-dashboard-nba/internal/modules/role/domain"

	"github.com/lib/pq"
	"github.com/pkg/errors"
)

func isPostgresForeignKeyViolation(err error) bool {
	if pqErr, ok := err.(*pq.Error); ok {
		return pqErr.Code == "23503"
	}
	return false
}

type roleRepository struct {
	db db.DB
}

// NewRoleRepository injects the database connection into the RoleRepository implementation.
func NewRoleRepository(db db.DB) domain.RoleRepository {
	return &roleRepository{
		db: db,
	}
}

func (r *roleRepository) CreateRoleQuery(ctx context.Context, payload domain.CreateRolePayload) (err error) {
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

func (r *roleRepository) UpdateRoleQuery(ctx context.Context, payload domain.UpdateRolePayload) (err error) {
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

func (r *roleRepository) DeleteRoleQuery(ctx context.Context, roleID int) (err error) {
	const stmt = `
		DELETE FROM app_role
		WHERE
			id = $1
	`
	_, err = r.db.ExecContext(ctx, stmt, roleID)

	return
}

func (r *roleRepository) UpdateRoleAccessTx(ctx context.Context, roleID int, payloads []domain.UpdateRoleMenuPermission) (err error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return errors.Wrap(err, "failed to begin transaction")
	}
	defer tx.Rollback()

	const insertStmt = `
		INSERT INTO app_role_access(
			role_id, menu_permission_id
		)
		VALUES ($1, $2)
		ON CONFLICT (role_id, menu_permission_id) DO NOTHING;
	`
	const deleteStmt = `
		DELETE FROM app_role_access
		WHERE 
			role_id = $1 AND menu_permission_id = $2		
	`

	for _, payload := range payloads {
		if payload.HasAccess != nil && *payload.HasAccess {
			_, err = tx.ExecContext(ctx, insertStmt, payload.RoleID, payload.MenuPermissionID)
			if err != nil {
				if isPostgresForeignKeyViolation(err) {
					return constant.ErrMenuPermissionIdNotFound
				}
				return errors.Wrap(err, "failed to insert role access")
			}
		} else {
			_, err = tx.ExecContext(ctx, deleteStmt, payload.RoleID, payload.MenuPermissionID)
			if err != nil {
				return errors.Wrap(err, "failed to delete role access")
			}
		}
	}

	if err = tx.Commit(); err != nil {
		return errors.Wrap(err, "failed to commit transaction")
	}

	return nil
}

func (r *roleRepository) ReadRolesQuery(
	ctx context.Context,
	filter domain.RoleFilter,
) (data []domain.Role, err error) {
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
			(CASE WHEN $3 = 'name ASC' THEN name END) ASC,
			(CASE WHEN $3 = 'name DESC' THEN name END) DESC,
			(CASE WHEN $3 = 'description ASC' THEN description END) ASC,
			(CASE WHEN $3 = 'description DESC' THEN description END) DESC,
			(CASE WHEN $3 = 'id ASC' THEN id END) ASC,
			(CASE WHEN $3 = 'id DESC' THEN id END) DESC
		LIMIT $4
		OFFSET $5
	`
	rows, err := r.db.QueryContext(ctx, stmt,
		filter.SetSearch,
		filter.Search,
		filter.Order,
		filter.Limit,
		filter.Offset,
	)
	if err != nil {
		return
	}
	defer rows.Close()

	for rows.Next() {
		var role domain.Role

		if err = rows.Scan(
			&role.ID,
			&role.Name,
			&role.Code,
			&role.Description,
		); err != nil {
			return
		}

		data = append(data, role)
	}

	return
}

func (r *roleRepository) ReadRolesCount(ctx context.Context, filter domain.RoleFilter) (count int, err error) {
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
		filter.SetSearch,
		filter.Search,
	).Scan(&count)
	if err != nil {
		fmt.Printf("Error in ReadRolesCount: %v\n", err)
	}
	return
}

func (r *roleRepository) ReadRoleByIDQuery(ctx context.Context, roleID int) (data domain.Role, err error) {
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

	if err != nil {
		if err == sql.ErrNoRows {
			err = constant.ErrDataNotFound
		}
		return
	}

	return
}

func (r *roleRepository) ReadRoleAccessQuery(
	ctx context.Context,
	filter domain.RoleAccessFilter,
) (data []domain.RoleAccessResponse, err error) {
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
		filter.RoleID,
		filter.SetSearch,
		filter.Search,
		filter.Order,
		filter.Limit,
		filter.Offset,
	)
	if err != nil {
		return
	}
	defer rows.Close()

	for rows.Next() {
		var res domain.RoleAccessResponse
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

func (r *roleRepository) ReadRoleAccessCount(
	ctx context.Context,
	filter domain.RoleAccessFilter,
) (count int, err error) {
	const stmt = `
		SELECT COUNT(*)
		FROM app_menu
		WHERE 
			(CASE WHEN $1::bool THEN (
				name ILIKE $2
			) ELSE TRUE END)
	`
	err = r.db.QueryRowContext(ctx, stmt,
		filter.SetSearch,
		filter.Search,
	).Scan(&count)
	return
}
