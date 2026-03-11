package menu

import (
	"be-dashboard-nba/constant"
	"be-dashboard-nba/usecase/entities"
	"context"
	"database/sql"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

func TestReadMenuDetailService(t *testing.T) {
	tests := []struct {
		name        string
		menuID      int
		expected    entities.Menu
		wantErr     bool
		expectedErr error
		setupMock   func(mock sqlmock.Sqlmock, menuID int)
	}{
		{
			name:   "success read menu detail",
			menuID: 1,
			expected: entities.Menu{
				ID:          1,
				ParentID:    sql.NullInt32{},
				Name:        "Dashboard",
				Description: sql.NullString{},
				URL:         sql.NullString{String: "/dashboard", Valid: true},
				Sort:        0,
				Group:       "main",
				Icon:        sql.NullString{},
				Active:      true,
				Display:     true,
			},
			wantErr:     false,
			expectedErr: nil,
			setupMock: func(mock sqlmock.Sqlmock, menuID int) {
				mock.ExpectQuery(`SELECT .* FROM app_menu WHERE id = \$1`).
					WithArgs(menuID).
					WillReturnRows(sqlmock.NewRows([]string{
						"id", "parent_id", "name", "description", "url", "sort", "group", "icon", "active", "display",
					}).AddRow(
						menuID, nil, "Dashboard", nil, "/dashboard", 0, "main", nil, true, true,
					))
			},
		},
		{
			name:   "success with parent ID",
			menuID: 2,
			expected: entities.Menu{
				ID:          2,
				ParentID:    sql.NullInt32{Int32: 1, Valid: true},
				Name:        "Submenu",
				Description: sql.NullString{},
				URL:         sql.NullString{String: "/submenu", Valid: true},
				Sort:        1,
				Group:       "main",
				Icon:        sql.NullString{},
				Active:      true,
				Display:     true,
			},
			wantErr:     false,
			expectedErr: nil,
			setupMock: func(mock sqlmock.Sqlmock, menuID int) {
				mock.ExpectQuery(`SELECT .* FROM app_menu WHERE id = \$1`).
					WithArgs(menuID).
					WillReturnRows(sqlmock.NewRows([]string{
						"id", "parent_id", "name", "description", "url", "sort", "group", "icon", "active", "display",
					}).AddRow(
						menuID, 1, "Submenu", nil, "/submenu", 1, "main", nil, true, true,
					))
			},
		},
		{
			name:   "success with description and icon",
			menuID: 3,
			expected: entities.Menu{
				ID:          3,
				ParentID:    sql.NullInt32{},
				Name:        "Settings",
				Description: sql.NullString{String: "Application settings", Valid: true},
				URL:         sql.NullString{String: "/settings", Valid: true},
				Sort:        2,
				Group:       "main",
				Icon:        sql.NullString{String: "settings-icon", Valid: true},
				Active:      true,
				Display:     true,
			},
			wantErr:     false,
			expectedErr: nil,
			setupMock: func(mock sqlmock.Sqlmock, menuID int) {
				mock.ExpectQuery(`SELECT .* FROM app_menu WHERE id = \$1`).
					WithArgs(menuID).
					WillReturnRows(sqlmock.NewRows([]string{
						"id", "parent_id", "name", "description", "url", "sort", "group", "icon", "active", "display",
					}).AddRow(
						menuID, nil, "Settings", "Application settings", "/settings", 2, "main", "settings-icon", true, true,
					))
			},
		},
		{
			name:   "success with all fields null except required",
			menuID: 4,
			expected: entities.Menu{
				ID:          4,
				ParentID:    sql.NullInt32{},
				Name:        "Minimal Menu",
				Description: sql.NullString{},
				URL:         sql.NullString{},
				Sort:        3,
				Group:       "main",
				Icon:        sql.NullString{},
				Active:      false,
				Display:     false,
			},
			wantErr:     false,
			expectedErr: nil,
			setupMock: func(mock sqlmock.Sqlmock, menuID int) {
				mock.ExpectQuery(`SELECT .* FROM app_menu WHERE id = \$1`).
					WithArgs(menuID).
					WillReturnRows(sqlmock.NewRows([]string{
						"id", "parent_id", "name", "description", "url", "sort", "group", "icon", "active", "display",
					}).AddRow(
						menuID, nil, "Minimal Menu", nil, nil, 3, "main", nil, false, false,
					))
			},
		},
		{
			name:        "menu not found",
			menuID:      999,
			expected:    entities.Menu{},
			wantErr:     true,
			expectedErr: constant.ErrMenuIdNotFound,
			setupMock: func(mock sqlmock.Sqlmock, menuID int) {
				mock.ExpectQuery(`SELECT .* FROM app_menu WHERE id = \$1`).
					WithArgs(menuID).
					WillReturnError(sql.ErrNoRows)
			},
		},
		{
			name:        "database query error",
			menuID:      1,
			expected:    entities.Menu{},
			wantErr:     true,
			expectedErr: constant.ErrUnknownSource,
			setupMock: func(mock sqlmock.Sqlmock, menuID int) {
				mock.ExpectQuery(`SELECT .* FROM app_menu WHERE id = \$1`).
					WithArgs(menuID).
					WillReturnError(errors.New("database connection failed"))
			},
		},
		{
			name:        "scan error - missing columns",
			menuID:      1,
			expected:    entities.Menu{},
			wantErr:     true,
			expectedErr: constant.ErrUnknownSource,
			setupMock: func(mock sqlmock.Sqlmock, menuID int) {
				// Return rows with missing columns to cause scan error
				mock.ExpectQuery(`SELECT .* FROM app_menu WHERE id = \$1`).
					WithArgs(menuID).
					WillReturnRows(sqlmock.NewRows([]string{
						"id", "name", // Missing other required columns
					}).AddRow(
						menuID, "Dashboard",
					))
			},
		},
		{
			name:        "zero menu ID",
			menuID:      0,
			expected:    entities.Menu{},
			wantErr:     true,
			expectedErr: constant.ErrMenuIdNotFound,
			setupMock: func(mock sqlmock.Sqlmock, menuID int) {
				mock.ExpectQuery(`SELECT .* FROM app_menu WHERE id = \$1`).
					WithArgs(menuID).
					WillReturnError(sql.ErrNoRows)
			},
		},
		{
			name:        "negative menu ID",
			menuID:      -1,
			expected:    entities.Menu{},
			wantErr:     true,
			expectedErr: constant.ErrMenuIdNotFound,
			setupMock: func(mock sqlmock.Sqlmock, menuID int) {
				mock.ExpectQuery(`SELECT .* FROM app_menu WHERE id = \$1`).
					WithArgs(menuID).
					WillReturnError(sql.ErrNoRows)
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			svc, mock, cleanup := newTestService(t)
			defer cleanup()

			ctx := context.Background()

			// Setup mock expectations
			tc.setupMock(mock, tc.menuID)

			// Execute service
			data, err := svc.ReadMenuDetailService(ctx, tc.menuID)

			// Assertions
			if tc.wantErr {
				assert.Error(t, err)
				if tc.expectedErr != nil {
					assert.ErrorIs(t, err, tc.expectedErr)
				}
				assert.Equal(t, tc.expected, data) // Should be empty on error
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expected, data)
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}
