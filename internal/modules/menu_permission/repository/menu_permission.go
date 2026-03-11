package repository

import (
	"be-dashboard-nba/constant"
	"be-dashboard-nba/internal/modules/menu_permission/domain"
	"context"
	"database/sql"

	"github.com/pkg/errors"
)

func (r *menuPermissionRepository) CreateMenuPermissionQuery(
	ctx context.Context,
	payload domain.MenuPermissionCreatePayload,
) error {

	const stmt = `
		INSERT INTO app_menu_permission (
			code, action_name, menu_id, created_by, created_at
		)
		VALUES ($1, $2, $3, $4, NOW())
	`

	_, err := r.db.ExecContext(ctx, stmt,
		payload.Code,
		payload.ActionName,
		payload.MenuID,
		payload.CreatedBy,
	)

	if err != nil {
		return errors.Wrap(err, "failed to execute insert menu permission query")
	}

	return nil
}

func (r *menuPermissionRepository) DeleteMenuPermissionQuery(ctx context.Context, id int) error {
	const stmt = `
		DELETE FROM app_menu_permission
		WHERE id = $1
	`

	res, err := r.db.ExecContext(ctx, stmt, id)
	if err != nil {
		return errors.Wrap(err, "failed to execute delete menu permission query")
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return errors.Wrap(err, "failed to get rows affected on delete")
	}

	if rowsAffected == 0 {
		return constant.ErrDataNotFound
	}

	return nil
}

func (r *menuPermissionRepository) ReadMenuPermissionByIDQuery(ctx context.Context, id int) (domain.MenuPermissionDetail, error) {
	var data domain.MenuPermissionDetail

	const stmt = `
		SELECT 
			id, menu_id, code, action_name
		FROM 
			app_menu_permission
		WHERE 
			id = $1
	`

	err := r.db.QueryRowContext(ctx, stmt, id).Scan(
		&data.ID,
		&data.MenuID,
		&data.Code,
		&data.ActionName,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return data, constant.ErrDataNotFound
		}

		return data, errors.Wrap(err, "failed to execute read menu permission detail query")
	}

	return data, nil
}

func (r *menuPermissionRepository) ReadMenuPermissionCountQuery(ctx context.Context,
	args domain.MenuPermissionFilter) (count int, err error) {
	const stmt = `
			SELECT COUNT(*)
			FROM
				app_menu_permission
			WHERE
				(CASE WHEN $1::bool THEN(
				action_name ILIKE $2
				) ELSE TRUE END)
		 		AND menu_id = $3
		`

	err = r.db.QueryRowContext(ctx, stmt,
		args.SetSearch,
		args.Search,
		args.MenuID,
	).Scan(&count)

	if err != nil {
		return 0, err
	}

	return
}

func (r *menuPermissionRepository) ReadMenuPermissionListQuery(
	ctx context.Context,
	args domain.MenuPermissionFilter,
) (data []domain.MenuPermissionDetail, err error) {

	const stmt = `
	SELECT
		id, menu_id, code, action_name
	FROM
		app_menu_permission
	WHERE
		(CASE WHEN $1::bool THEN(
			action_name ILIKE $2
		) ELSE TRUE END)
		 AND menu_id = $6
	ORDER BY
		(CASE WHEN $3 = 'action_name ASC' THEN action_name END)
		ASC,
		(CASE WHEN $3 = 'action_name DESC' THEN action_name END)
		DESC,
		(CASE WHEN $3 = 'code ASC' THEN code END)
		ASC,
		(CASE WHEN $3 = 'code DESC' THEN code END)
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
		args.MenuID,
	)
	if err != nil {
		return
	}
	defer rows.Close()

	for rows.Next() {
		var mp domain.MenuPermissionDetail

		if err = rows.Scan(
			&mp.ID,
			&mp.MenuID,
			&mp.Code,
			&mp.ActionName,
		); err != nil {
			return
		}
		data = append(data, mp)
	}

	return
}

func (r *menuPermissionRepository) UpdateMenuPermissionQuery(
	ctx context.Context,
	payload domain.MenuPermissionUpdatePayload,
) error {

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

	// Eksekusi tanpa transaksi
	res, err := r.db.ExecContext(ctx, stmt,
		payload.Code,
		payload.ActionName,
		payload.UpdatedBy,
		payload.ID,
	)

	if err != nil {
		return errors.Wrap(err, "failed to execute update menu permission query")
	}

	// Memastikan benar-benar ada data yang di-update
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return errors.Wrap(err, "failed to get rows affected on update")
	}

	if rowsAffected == 0 {
		return constant.ErrDataNotFound
	}

	return nil
}
