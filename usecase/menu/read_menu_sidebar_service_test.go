package menu

import (
	"be-dashboard-nba/usecase/entities"
	"context"
	"database/sql"
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestReadSidebarMenuService(t *testing.T) {
	tests := []struct {
		name        string
		userID      string
		expected    []entities.Menu
		expectedErr error
		setupMock   func(mock sqlmock.Sqlmock, userID string)
	}{
		{
			name:   "success - single menu",
			userID: "user1",
			expected: []entities.Menu{
				{
					ID:       1,
					ParentID: sql.NullInt32{},
					Name:     "Dashboard",
					URL:      sql.NullString{String: "/dashboard", Valid: true},
					Sort:     0,
					Group:    "main",
					Icon:     sql.NullString{},
					Active:   true,
					Display:  true,
				},
			},
			expectedErr: nil,
			setupMock: func(mock sqlmock.Sqlmock, userID string) {
				mock.ExpectQuery(`SELECT DISTINCT\s+m.id, m.parent_id, m.name, m.url, m.sort,`).
					WithArgs(userID).
					WillReturnRows(sqlmock.NewRows([]string{
						"id", "parent_id", "name", "url", "sort", "group", "icon", "active", "display",
					}).AddRow(
						int64(1),
						nil,
						"Dashboard",
						"/dashboard",
						int32(0),
						"main",
						nil,
						true,
						true,
					))
			},
		},
		{
			name:   "success - multiple menus with hierarchy",
			userID: "user1",
			expected: []entities.Menu{
				{
					ID:       1,
					ParentID: sql.NullInt32{},
					Name:     "Dashboard",
					URL:      sql.NullString{String: "/dashboard", Valid: true},
					Sort:     0,
					Group:    "main",
					Icon:     sql.NullString{},
					Active:   true,
					Display:  true,
				},
				{
					ID:       2,
					ParentID: sql.NullInt32{Int32: 1, Valid: true},
					Name:     "Submenu",
					URL:      sql.NullString{String: "/submenu", Valid: true},
					Sort:     1,
					Group:    "main",
					Icon:     sql.NullString{String: "icon", Valid: true},
					Active:   true,
					Display:  true,
				},
				{
					ID:       3,
					ParentID: sql.NullInt32{},
					Name:     "Settings",
					URL:      sql.NullString{String: "/settings", Valid: true},
					Sort:     0,
					Group:    "settings",
					Icon:     sql.NullString{String: "gear", Valid: true},
					Active:   true,
					Display:  true,
				},
			},
			expectedErr: nil,
			setupMock: func(mock sqlmock.Sqlmock, userID string) {
				mock.ExpectQuery(`SELECT DISTINCT\s+m.id, m.parent_id, m.name, m.url, m.sort,`).
					WithArgs(userID).
					WillReturnRows(sqlmock.NewRows([]string{
						"id", "parent_id", "name", "url", "sort", "group", "icon", "active", "display",
					}).
						AddRow(int64(1), nil, "Dashboard", "/dashboard", int32(0), "main", nil, true, true).
						AddRow(int64(2), int32(1), "Submenu", "/submenu", int32(1), "main", "icon", true, true).
						AddRow(int64(3), nil, "Settings", "/settings", int32(0), "settings", "gear", true, true))
			},
		},
		{
			name:        "success - empty result",
			userID:      "user1",
			expected:    nil,
			expectedErr: nil,
			setupMock: func(mock sqlmock.Sqlmock, userID string) {
				mock.ExpectQuery(`SELECT DISTINCT\s+m.id, m.parent_id, m.name, m.url, m.sort,`).
					WithArgs(userID).
					WillReturnRows(sqlmock.NewRows([]string{
						"id", "parent_id", "name", "url", "sort", "group", "icon", "active", "display",
					}))
			},
		},
		{
			name:        "error - query execution fails",
			userID:      "user1",
			expected:    nil,
			expectedErr: sql.ErrConnDone,
			setupMock: func(mock sqlmock.Sqlmock, userID string) {
				mock.ExpectQuery(`SELECT DISTINCT\s+m.id, m.parent_id, m.name, m.url, m.sort,`).
					WithArgs(userID).
					WillReturnError(sql.ErrConnDone)
			},
		},
		{
			name:        "error - database connection error",
			userID:      "user1",
			expected:    nil,
			expectedErr: errors.New("database connection failed"),
			setupMock: func(mock sqlmock.Sqlmock, userID string) {
				mock.ExpectQuery(`SELECT DISTINCT\s+m.id, m.parent_id, m.name, m.url, m.sort,`).
					WithArgs(userID).
					WillReturnError(errors.New("database connection failed"))
			},
		},
		{
			name:   "success - menu with null values",
			userID: "user1",
			expected: []entities.Menu{
				{
					ID:       1,
					ParentID: sql.NullInt32{Int32: 0, Valid: false},
					Name:     "Menu with Nulls",
					URL:      sql.NullString{String: "", Valid: false},
					Sort:     0,
					Group:    "main",
					Icon:     sql.NullString{String: "", Valid: false},
					Active:   true,
					Display:  true,
				},
			},
			expectedErr: nil,
			setupMock: func(mock sqlmock.Sqlmock, userID string) {
				mock.ExpectQuery(`SELECT DISTINCT\s+m.id, m.parent_id, m.name, m.url, m.sort,`).
					WithArgs(userID).
					WillReturnRows(sqlmock.NewRows([]string{
						"id", "parent_id", "name", "url", "sort", "group", "icon", "active", "display",
					}).AddRow(
						int64(1),
						nil,
						"Menu with Nulls",
						nil,
						int32(0),
						"main",
						nil,
						true,
						true,
					))
			},
		},
		{
			name:   "success - different user ID",
			userID: "user2",
			expected: []entities.Menu{
				{
					ID:       1,
					ParentID: sql.NullInt32{},
					Name:     "User2 Menu",
					URL:      sql.NullString{String: "/user2", Valid: true},
					Sort:     0,
					Group:    "main",
					Icon:     sql.NullString{},
					Active:   true,
					Display:  true,
				},
			},
			expectedErr: nil,
			setupMock: func(mock sqlmock.Sqlmock, userID string) {
				mock.ExpectQuery(`SELECT DISTINCT\s+m.id, m.parent_id, m.name, m.url, m.sort,`).
					WithArgs(userID).
					WillReturnRows(sqlmock.NewRows([]string{
						"id", "parent_id", "name", "url", "sort", "group", "icon", "active", "display",
					}).AddRow(
						int64(1),
						nil,
						"User2 Menu",
						"/user2",
						int32(0),
						"main",
						nil,
						true,
						true,
					))
			},
		},
		{
			name:        "error - scan error (type mismatch)",
			userID:      "user1",
			expected:    nil,
			expectedErr: nil, // Error akan terjadi selama scan, bukan return error dari query
			setupMock: func(mock sqlmock.Sqlmock, userID string) {
				// Simulate type mismatch by returning wrong type for ID
				mock.ExpectQuery(`SELECT DISTINCT\s+m.id, m.parent_id, m.name, m.url, m.sort,`).
					WithArgs(userID).
					WillReturnRows(sqlmock.NewRows([]string{
						"id", "parent_id", "name", "url", "sort", "group", "icon", "active", "display",
					}).AddRow(
						"invalid_id", // Wrong type - should be int64
						nil,
						"Dashboard",
						"/dashboard",
						int32(0),
						"main",
						nil,
						true,
						true,
					))
			},
		},
		{
			name:   "success - menus ordered by group and sort",
			userID: "user1",
			expected: []entities.Menu{
				{
					ID:       1,
					ParentID: sql.NullInt32{},
					Name:     "First Group - First Item",
					URL:      sql.NullString{String: "/first", Valid: true},
					Sort:     0,
					Group:    "alpha",
					Icon:     sql.NullString{},
					Active:   true,
					Display:  true,
				},
				{
					ID:       2,
					ParentID: sql.NullInt32{},
					Name:     "First Group - Second Item",
					URL:      sql.NullString{String: "/second", Valid: true},
					Sort:     1,
					Group:    "alpha",
					Icon:     sql.NullString{},
					Active:   true,
					Display:  true,
				},
				{
					ID:       3,
					ParentID: sql.NullInt32{},
					Name:     "Second Group - First Item",
					URL:      sql.NullString{String: "/third", Valid: true},
					Sort:     0,
					Group:    "beta",
					Icon:     sql.NullString{},
					Active:   true,
					Display:  true,
				},
			},
			expectedErr: nil,
			setupMock: func(mock sqlmock.Sqlmock, userID string) {
				mock.ExpectQuery(`SELECT DISTINCT\s+m.id, m.parent_id, m.name, m.url, m.sort,`).
					WithArgs(userID).
					WillReturnRows(sqlmock.NewRows([]string{
						"id", "parent_id", "name", "url", "sort", "group", "icon", "active", "display",
					}).
						AddRow(int64(1), nil, "First Group - First Item", "/first", int32(0), "alpha", nil, true, true).
						AddRow(int64(2), nil, "First Group - Second Item", "/second", int32(1), "alpha", nil, true, true).
						AddRow(int64(3), nil, "Second Group - First Item", "/third", int32(0), "beta", nil, true, true))
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			svc, mock, cleanup := newTestService(t)
			defer cleanup()

			ctx := context.Background()

			// Setup mock expectations
			tc.setupMock(mock, tc.userID)

			// Execute service
			data, err := svc.ReadSidebarMenuService(ctx, tc.userID)

			// Assertions
			if tc.expectedErr != nil {
				assert.Error(t, err)
				if tc.expectedErr != sql.ErrConnDone { // Special case for connection done error
					assert.Contains(t, err.Error(), tc.expectedErr.Error())
				}
				assert.Nil(t, data)
			} else {
				if tc.name == "error - scan error (type mismatch)" {
					// For scan errors, we expect an error during scanning
					assert.Error(t, err)
					assert.Nil(t, data)
				} else {
					assert.NoError(t, err)
					assert.Equal(t, tc.expected, data)
					if len(tc.expected) == 0 {
						assert.Empty(t, data)
					} else {
						assert.Len(t, data, len(tc.expected))
					}
				}
			}

			// Verify all expectations were met
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}
