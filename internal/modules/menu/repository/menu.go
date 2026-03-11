package repository

import (
	"be-dashboard-nba/constant"
	"be-dashboard-nba/internal/modules/menu/domain"
	"context"
	"database/sql"

	"github.com/pkg/errors"
)

func (r *menuRepository) CreateMenuQuery(
	ctx context.Context,
	params domain.MenuCreatePayload,
	userID string,
) (err error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return errors.Wrap(err, "failed to begin transaction")
	}
	defer tx.Rollback()

	var nextSort int

	if params.ParentID != nil {
		const stmtParent = `SELECT COALESCE(MAX(sort), -1) + 1 FROM app_menu WHERE parent_id = $1`
		err = tx.QueryRowContext(ctx, stmtParent, *params.ParentID).Scan(&nextSort)
	} else {
		const stmtGroup = `SELECT COALESCE(MAX(sort), -1) + 1 FROM app_menu WHERE parent_id IS NULL AND "group" = $1`
		err = tx.QueryRowContext(ctx, stmtGroup, params.Group).Scan(&nextSort)
	}

	if err != nil {
		return errors.Wrap(err, "failed to get next sort value")
	}

	params.Sort = nextSort

	const insertStmt = `
		INSERT INTO app_menu (
			parent_id, name, description, url, sort, "group", icon,
			active, display, created_by, created_at
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, NOW())
	`

	_, err = tx.ExecContext(ctx, insertStmt,
		params.ParentID,
		params.Name,
		params.Description,
		params.URL,
		params.Sort,
		params.Group,
		params.Icon,
		params.Active,
		params.Display,
		userID,
	)
	if err != nil {
		return errors.Wrap(err, "failed to insert new menu")
	}

	// 4. Commit Transaksi
	if err = tx.Commit(); err != nil {
		return errors.Wrap(err, "failed to commit transaction")
	}

	return nil
}

func (r *menuRepository) DeleteMenuQuery(ctx context.Context, menuID int) error {
	const stmt = `
		DELETE FROM app_menu
		WHERE id = $1
	`

	// Eksekusi query delete langsung (tanpa Tx)
	res, err := r.db.ExecContext(ctx, stmt, menuID)
	if err != nil {
		return errors.Wrap(err, "failed to execute delete menu query")
	}

	// Cek berapa banyak baris yang berhasil dihapus oleh database
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return errors.Wrap(err, "failed to check deleted rows")
	}

	// Jika 0, berarti menu dengan ID tersebut memang tidak ada di database
	if rowsAffected == 0 {
		return constant.ErrDataNotFound
	}

	return nil
}

func (r *menuRepository) ReadMenuByIDQuery(ctx context.Context, menuID int) (data domain.Menu, err error) {
	const stmt = `
		SELECT
			id, parent_id, name, description, url, sort, "group", icon,
			active, display
		FROM 
			app_menu
		WHERE
			id = $1
	`

	err = r.db.QueryRowContext(ctx, stmt, menuID).Scan(
		&data.ID, &data.ParentID, &data.Name, &data.Description, &data.URL, &data.Sort,
		&data.Group, &data.Icon, &data.Active, &data.Display,
	)

	if err != nil {
		// ✨ TERJEMAHKAN ERROR SQL DI SINI ✨
		if errors.Is(err, sql.ErrNoRows) {
			return data, constant.ErrDataNotFound
		}
		// Bungkus error lain agar log-nya detail
		return data, errors.Wrap(err, "failed to scan read menu by id")
	}

	return data, nil
}

func (r *menuRepository) ReadParentMenuQuery(
	ctx context.Context,
) (data []domain.MenuParent, err error) {
	const stmt = `
		SELECT
			id, name, "group"
		FROM
			app_menu
		WHERE
			parent_id IS NULL
		ORDER BY
			sort ASC;
	`
	rows, err := r.db.QueryContext(ctx, stmt)
	if err != nil {
		return
	}
	defer rows.Close()

	for rows.Next() {
		var m domain.MenuParent
		if err = rows.Scan(
			&m.ID, &m.Name, &m.Group,
		); err != nil {
			return
		}
		data = append(data, m)
	}
	return
}

func (r *menuRepository) ReadListMenuQuery(
	ctx context.Context,
	params domain.MenuFilter,
) (data []domain.Menu, err error) {
	const stmt = `
		WITH RECURSIVE menu_with_parents AS (
			SELECT
				id, parent_id, name, description, url, sort, "group", icon,
				active, display, created_by, created_at, updated_by, updated_at
			FROM
				app_menu
			WHERE
				(CASE WHEN $1::bool THEN
					name ILIKE $2
				ELSE TRUE 
				END)

			UNION

			SELECT
				m.id, m.parent_id, m.name, m.description, m.url, m.sort, m."group", m.icon,
				m.active, m.display, m.created_by, m.created_at, m.updated_by, m.updated_at
			FROM
				app_menu m
			JOIN
				menu_with_parents mwp ON m.id = mwp.parent_id 
		)
		SELECT DISTINCT * FROM menu_with_parents
		ORDER BY
			sort ASC;
	`

	rows, err := r.db.QueryContext(ctx, stmt,
		params.SetSearch,
		params.Search,
	)
	if err != nil {
		return
	}
	defer rows.Close()

	for rows.Next() {
		var m domain.Menu
		if err = rows.Scan(
			&m.ID, &m.ParentID, &m.Name, &m.Description, &m.URL, &m.Sort,
			&m.Group, &m.Icon, &m.Active, &m.Display,
			&m.CreatedBy, &m.CreatedAt, &m.UpdatedBy, &m.UpdatedAt,
		); err != nil {
			return
		}
		data = append(data, m)
	}

	return
}

func (r *menuRepository) ReadSidebarMenuQuery(
	ctx context.Context,
	userID string,
) (data []domain.Menu, err error) {

	const stmt = `
		SELECT DISTINCT
			m.id, m.parent_id, m.name, m.url, m.sort,
			m."group", m.icon, m.active, m.display
		FROM
			app_menu m
		JOIN
			app_menu_permission amp ON m.id = amp.menu_id
		JOIN
			app_role_access ara ON amp.id = ara.menu_permission_id
		JOIN
			app_user_role aur ON ara.role_id = aur.role_id
		WHERE
			aur.user_id = $1
			AND amp.code = 'R' 
			AND m.display = true   
			AND m.active = true  
		ORDER BY
    		m."group" ASC,
   		 	m.sort ASC;
	`

	rows, err := r.db.QueryContext(ctx, stmt, userID)
	if err != nil {
		return
	}
	defer rows.Close()

	for rows.Next() {
		var m domain.Menu
		if err = rows.Scan(
			&m.ID, &m.ParentID, &m.Name, &m.URL, &m.Sort,
			&m.Group, &m.Icon, &m.Active, &m.Display,
		); err != nil {
			return
		}
		data = append(data, m)
	}
	return
}

func (r *menuRepository) ReadNextSortForParent(ctx context.Context, parentID int32) (int, error) {
	var nextSort int
	const stmt = `
		SELECT COALESCE(MAX(sort), -1) + 1 
		FROM app_menu
		WHERE parent_id = $1
	`

	err := r.db.QueryRowContext(ctx, stmt, parentID).Scan(&nextSort)
	return nextSort, err
}

func (r *menuRepository) ReadSortForGroup(ctx context.Context, group string) (int, error) {
	var nextSort int
	const stmt = `
		SELECT COALESCE(MAX(sort), -1) + 1 
		FROM app_menu
		WHERE parent_id IS NULL AND "group" = $1
	`
	err := r.db.QueryRowContext(ctx, stmt, group).Scan(&nextSort)
	return nextSort, err
}

func (r *menuRepository) ReadNextSortForParentAndGroup(ctx context.Context, parentID int32, group string) (int, error) {
	const stmt = `SELECT COALESCE(MAX(sort), -1) + 1 FROM app_menu WHERE parent_id = $1 AND "group" = $2`
	var nextSort int
	err := r.db.QueryRowContext(ctx, stmt, parentID, group).Scan(&nextSort)
	return nextSort, err
}

func (r *menuRepository) CountMenuChildren(ctx context.Context, menuID int) (int, error) {
	const stmt = `SELECT COUNT(*) FROM app_menu WHERE parent_id = $1`
	var count int
	err := r.db.QueryRowContext(ctx, stmt, menuID).Scan(&count)
	return count, err
}

func (r *menuRepository) ReadCountMenuQuery(
	ctx context.Context,
	params domain.MenuFilter,
) (count int64, err error) {
	const stmt = `
		SELECT
			COUNT(*)
		FROM
			app_menu
		WHERE
			(CASE WHEN $1::bool THEN(
				name ILIKE $2
				OR description ILIKE $2
				OR url ILIKE $2
			) ELSE TRUE END)
	`
	err = r.db.QueryRowContext(ctx, stmt,
		params.SetSearch,
		params.Search,
	).Scan(&count)
	return
}

func (r *menuRepository) UpdateMenuQuery(
	ctx context.Context,
	params domain.MenuUpdatePayload,
	updateChildrenGroup bool,
) error {

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return errors.Wrap(err, "failed to begin transaction")
	}
	defer tx.Rollback()

	// 1. Ambil Parent dan Group lama langsung menggunakan tx untuk perbandingan Sort
	var oldParentID sql.NullInt32
	var oldGroup string
	err = tx.QueryRowContext(ctx, `SELECT parent_id, "group" FROM app_menu WHERE id = $1`, params.ID).Scan(&oldParentID, &oldGroup)
	if err != nil {
		return errors.Wrap(err, "failed to read existing menu for update tx")
	}

	// Deteksi perubahan
	var parentChanged bool
	if !oldParentID.Valid && params.ParentID == nil {
		parentChanged = false
	} else if oldParentID.Valid && params.ParentID != nil {
		parentChanged = oldParentID.Int32 != *params.ParentID
	} else {
		parentChanged = true
	}

	groupChanged := oldGroup != params.Group

	// 2. Kalkulasi ulang Sort di dalam transaksi (mencegah Race Condition)
	if parentChanged {
		if params.ParentID != nil {
			const stmt = `SELECT COALESCE(MAX(sort), -1) + 1 FROM app_menu WHERE parent_id = $1`
			err = tx.QueryRowContext(ctx, stmt, *params.ParentID).Scan(&params.Sort)
		} else {
			const stmt = `SELECT COALESCE(MAX(sort), -1) + 1 FROM app_menu WHERE parent_id IS NULL AND "group" = $1`
			err = tx.QueryRowContext(ctx, stmt, params.Group).Scan(&params.Sort)
		}
		if err != nil {
			return errors.Wrap(err, "failed to calculate new sort for parent change")
		}
	} else if groupChanged && params.ParentID == nil {
		const stmt = `SELECT COALESCE(MAX(sort), -1) + 1 FROM app_menu WHERE parent_id IS NULL AND "group" = $1`
		err = tx.QueryRowContext(ctx, stmt, params.Group).Scan(&params.Sort)
		if err != nil {
			return errors.Wrap(err, "failed to calculate new sort for group change")
		}
	}

	// 3. Update Menu Utama menggunakan tx
	const updateStmt = `
		UPDATE app_menu
		SET
			name = $2, 
			parent_id = $3,
			description = $4,
			url = $5,
			sort = $6,
			"group" = $7, 
			icon = $8,      
			active = $9,     
			display = $10, 
			updated_by = $11,
			updated_at = NOW()
		WHERE
			id = $1
	`
	_, err = tx.ExecContext(ctx, updateStmt,
		params.ID, params.Name, params.ParentID, params.Description, params.URL,
		params.Sort, params.Group, params.Icon, params.Active, params.Display, params.UpdatedBy,
	)
	if err != nil {
		return errors.Wrap(err, "failed to update menu")
	}

	// 4. Update Children Group jika flag bernilai true
	if updateChildrenGroup {
		const childStmt = `UPDATE app_menu SET "group" = $1, updated_at = NOW() WHERE parent_id = $2`
		_, err = tx.ExecContext(ctx, childStmt, params.Group, params.ID)
		if err != nil {
			return errors.Wrap(err, "failed to update children group")
		}
	}

	// 5. Commit Transaksi
	if err = tx.Commit(); err != nil {
		return errors.Wrap(err, "failed to commit transaction")
	}

	return nil
}

func (r *menuRepository) UpdateMenuOrderQuery(
	ctx context.Context,
	payloads []domain.MenuUpdateSortItemPayload,
) error {
	// 1. Mulai Transaksi
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return errors.Wrap(err, "failed to begin transaction")
	}
	defer tx.Rollback()

	// Siapkan statement query di luar loop agar efisien
	const updateSortStmt = `
		UPDATE app_menu
		SET 
			sort = $1,
			"group" = $2,
			parent_id = $3,
			updated_by = $4,
			updated_at = NOW()
		WHERE 
			id = $5
	`
	const updateChildrenStmt = `
		UPDATE app_menu 
		SET "group" = $1, updated_at = NOW() 
		WHERE parent_id = $2
	`

	for _, payload := range payloads {
		// Update detail menu utama
		_, err = tx.ExecContext(ctx, updateSortStmt,
			payload.Sort,
			payload.Group,
			payload.ParentID,
			payload.UpdatedBy,
			payload.ID,
		)
		if err != nil {
			return errors.Wrapf(err, "failed to update sort for menu id %d", payload.ID)
		}

		// Update kolom 'group' untuk semua anak-anaknya agar selaras
		_, err = tx.ExecContext(ctx, updateChildrenStmt, payload.Group, payload.ID)
		if err != nil {
			return errors.Wrapf(err, "failed to update children group for menu id %d", payload.ID)
		}
	}

	// 3. Commit Transaksi jika SEMUA menu berhasil di-update tanpa ada yang tertinggal
	if err = tx.Commit(); err != nil {
		return errors.Wrap(err, "failed to commit transaction")
	}

	return nil
}
