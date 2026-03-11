package menu_permission

import (
	menuPermissionPresenter "be-dashboard-nba/api/presenter/menu_permission"
	"be-dashboard-nba/constant"
	"be-dashboard-nba/repository/menu_permission"
	"context"
	"database/sql"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

func TestCreateMenuPermissionService(t *testing.T) {
	tests := []struct {
		name        string
		payload     menuPermissionPresenter.CreateMenuPermissionRequest
		userID      string
		menuID      int
		wantErr     bool
		expectedErr error
		setupMock   func(mock sqlmock.Sqlmock, payload menuPermissionPresenter.CreateMenuPermissionRequest, userID string, menuID int)
	}{
		{
			name: "success create menu permission",
			payload: menuPermissionPresenter.CreateMenuPermissionRequest{
				Code:       "R",
				ActionName: "read",
			},
			userID:      "user-123",
			menuID:      1,
			wantErr:     false,
			expectedErr: nil,
			setupMock: func(mock sqlmock.Sqlmock, payload menuPermissionPresenter.CreateMenuPermissionRequest, userID string, menuID int) {
				mock.ExpectBegin()

				// Mock menu exists check
				mock.ExpectQuery(`SELECT .* FROM app_menu WHERE id = \$1`).
					WithArgs(menuID).
					WillReturnRows(sqlmock.NewRows([]string{
						"id", "parent_id", "name", "description", "url", "sort", "group", "icon", "active", "display",
					}).AddRow(
						menuID, nil, "Dashboard", nil, "/dashboard", 0, "main", nil, true, true,
					))

				// Mock insert menu permission
				expectedParams := payload.ToParams(userID, menuID)
				mock.ExpectExec(`INSERT INTO app_menu_permission`).
					WithArgs(
						expectedParams.Code,
						expectedParams.ActionName,
						expectedParams.MenuID,
						expectedParams.CreatedBy,
					).WillReturnResult(sqlmock.NewResult(1, 1))

				mock.ExpectCommit()
			},
		},
		{
			name: "error beginning transaction",
			payload: menuPermissionPresenter.CreateMenuPermissionRequest{
				Code:       "R",
				ActionName: "read",
			},
			userID:      "user-123",
			menuID:      1,
			wantErr:     true,
			expectedErr: constant.ErrUnknownSource,
			setupMock: func(mock sqlmock.Sqlmock, payload menuPermissionPresenter.CreateMenuPermissionRequest, userID string, menuID int) {
				mock.ExpectBegin().WillReturnError(errors.New("begin tx error"))
			},
		},
		{
			name: "menu not found",
			payload: menuPermissionPresenter.CreateMenuPermissionRequest{
				Code:       "R",
				ActionName: "read",
			},
			userID:      "user-123",
			menuID:      999,
			wantErr:     true,
			expectedErr: constant.ErrMenuIdNotFound,
			setupMock: func(mock sqlmock.Sqlmock, payload menuPermissionPresenter.CreateMenuPermissionRequest, userID string, menuID int) {
				mock.ExpectBegin()

				mock.ExpectQuery(`SELECT .* FROM app_menu WHERE id = \$1`).
					WithArgs(menuID).
					WillReturnError(sql.ErrNoRows)

				mock.ExpectRollback()
			},
		},
		{
			name: "error reading menu by ID",
			payload: menuPermissionPresenter.CreateMenuPermissionRequest{
				Code:       "R",
				ActionName: "read",
			},
			userID:      "user-123",
			menuID:      1,
			wantErr:     true,
			expectedErr: constant.ErrUnknownSource,
			setupMock: func(mock sqlmock.Sqlmock, payload menuPermissionPresenter.CreateMenuPermissionRequest, userID string, menuID int) {
				mock.ExpectBegin()

				mock.ExpectQuery(`SELECT .* FROM app_menu WHERE id = \$1`).
					WithArgs(menuID).
					WillReturnError(errors.New("select error"))

				mock.ExpectRollback()
			},
		},
		{
			name: "error creating menu permission",
			payload: menuPermissionPresenter.CreateMenuPermissionRequest{
				Code:       "R",
				ActionName: "read",
			},
			userID:      "user-123",
			menuID:      1,
			wantErr:     true,
			expectedErr: constant.ErrUnknownSource,
			setupMock: func(mock sqlmock.Sqlmock, payload menuPermissionPresenter.CreateMenuPermissionRequest, userID string, menuID int) {
				mock.ExpectBegin()

				mock.ExpectQuery(`SELECT .* FROM app_menu WHERE id = \$1`).
					WithArgs(menuID).
					WillReturnRows(sqlmock.NewRows([]string{
						"id", "parent_id", "name", "description", "url", "sort", "group", "icon", "active", "display",
					}).AddRow(
						menuID, nil, "Dashboard", nil, "/dashboard", 0, "main", nil, true, true,
					))

				expectedParams := payload.ToParams(userID, menuID)
				mock.ExpectExec(`INSERT INTO app_menu_permission`).
					WithArgs(
						expectedParams.Code,
						expectedParams.ActionName,
						expectedParams.MenuID,
						expectedParams.CreatedBy,
					).WillReturnError(errors.New("insert error"))

				mock.ExpectRollback()
			},
		},
		{
			name: "error committing transaction",
			payload: menuPermissionPresenter.CreateMenuPermissionRequest{
				Code:       "R",
				ActionName: "read",
			},
			userID:      "user-123",
			menuID:      1,
			wantErr:     true,
			expectedErr: constant.ErrUnknownSource,
			setupMock: func(mock sqlmock.Sqlmock, payload menuPermissionPresenter.CreateMenuPermissionRequest, userID string, menuID int) {
				mock.ExpectBegin()

				mock.ExpectQuery(`SELECT .* FROM app_menu WHERE id = \$1`).
					WithArgs(menuID).
					WillReturnRows(sqlmock.NewRows([]string{
						"id", "parent_id", "name", "description", "url", "sort", "group", "icon", "active", "display",
					}).AddRow(
						menuID, nil, "Dashboard", nil, "/dashboard", 0, "main", nil, true, true,
					))

				expectedParams := payload.ToParams(userID, menuID)
				mock.ExpectExec(`INSERT INTO app_menu_permission`).
					WithArgs(
						expectedParams.Code,
						expectedParams.ActionName,
						expectedParams.MenuID,
						expectedParams.CreatedBy,
					).WillReturnResult(sqlmock.NewResult(1, 1))

				mock.ExpectCommit().WillReturnError(errors.New("commit error"))
			},
		},
		{
			name: "error rolling back transaction",
			payload: menuPermissionPresenter.CreateMenuPermissionRequest{
				Code:       "R",
				ActionName: "read",
			},
			userID:      "user-123",
			menuID:      1,
			wantErr:     true,
			expectedErr: constant.ErrUnknownSource,
			setupMock: func(mock sqlmock.Sqlmock, payload menuPermissionPresenter.CreateMenuPermissionRequest, userID string, menuID int) {
				mock.ExpectBegin()

				mock.ExpectQuery(`SELECT .* FROM app_menu WHERE id = \$1`).
					WithArgs(menuID).
					WillReturnError(errors.New("select error"))

				// Mock rollback also fails
				mock.ExpectRollback().WillReturnError(errors.New("rollback failed"))
			},
		},
		{
			name: "create with different code and action name",
			payload: menuPermissionPresenter.CreateMenuPermissionRequest{
				Code:       "U",
				ActionName: "update",
			},
			userID:      "user-456",
			menuID:      2,
			wantErr:     false,
			expectedErr: nil,
			setupMock: func(mock sqlmock.Sqlmock, payload menuPermissionPresenter.CreateMenuPermissionRequest, userID string, menuID int) {
				mock.ExpectBegin()

				// Mock menu exists check for different menu
				mock.ExpectQuery(`SELECT .* FROM app_menu WHERE id = \$1`).
					WithArgs(menuID).
					WillReturnRows(sqlmock.NewRows([]string{
						"id", "parent_id", "name", "description", "url", "sort", "group", "icon", "active", "display",
					}).AddRow(
						menuID, nil, "Settings", nil, "/settings", 1, "main", nil, true, true,
					))

				// Mock insert with different parameters
				expectedParams := payload.ToParams(userID, menuID)
				mock.ExpectExec(`INSERT INTO app_menu_permission`).
					WithArgs(
						expectedParams.Code,
						expectedParams.ActionName,
						expectedParams.MenuID,
						expectedParams.CreatedBy,
					).WillReturnResult(sqlmock.NewResult(2, 1))

				mock.ExpectCommit()
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			svc, mock, cleanup := newTestService(t)
			defer cleanup()

			ctx := context.Background()

			// Setup mock expectations
			tc.setupMock(mock, tc.payload, tc.userID, tc.menuID)

			// Execute service
			err := svc.CreateMenuPermissionService(ctx, tc.payload, tc.userID, tc.menuID)

			// Assertions
			if tc.wantErr {
				assert.Error(t, err)
				if tc.expectedErr != nil {
					assert.ErrorIs(t, err, tc.expectedErr)
				}
			} else {
				assert.NoError(t, err)
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

// Test helper untuk memverifikasi ToParams function
func TestCreateMenuPermissionRequest_ToParams(t *testing.T) {
	req := menuPermissionPresenter.CreateMenuPermissionRequest{
		Code:       "C",
		ActionName: "Create",
	}

	userID := "test-user"
	menuID := 5

	params := req.ToParams(userID, menuID)

	expected := repository.CreateMenuPermissionPayload{
		Code:       "C",
		ActionName: "Create",
		MenuID:     5,
		CreatedBy:  "test-user",
	}

	assert.Equal(t, expected, params)
}
