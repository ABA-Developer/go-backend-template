package menu

import (
	menuPresenter "be-dashboard-nba/api/presenter/menu"
	"be-dashboard-nba/constant"
	"context"
	"database/sql"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

func TestCreateMenuService(t *testing.T) {
	getPtrString := func(s string) *string { return &s }
	getPtrBool := func(b bool) *bool { return &b }
	getPtrInt := func(i int) *int { return &i }

	tests := []struct {
		name        string
		request     menuPresenter.CreateMenuRequest
		userID      string
		wantErr     bool
		expectedErr error
		setupMock   func(mock sqlmock.Sqlmock, request menuPresenter.CreateMenuRequest, userID string)
	}{
		{
			name: "success with parent ID",
			request: menuPresenter.CreateMenuRequest{
				Name:     "Dashboard",
				Group:    "main",
				ParentID: getPtrInt(1),
				URL:      getPtrString("/dashboard"),
				Active:   getPtrBool(true),
				Display:  getPtrBool(true),
			},
			userID:      "user-123",
			wantErr:     false,
			expectedErr: nil,
			setupMock: func(mock sqlmock.Sqlmock, request menuPresenter.CreateMenuRequest, userID string) {
				mock.ExpectBegin()

				parentID := *request.ParentID
				mock.ExpectQuery(`SELECT COALESCE\(MAX\(sort\), -1\) \+ 1 FROM app_menu WHERE parent_id = \$1`).
					WithArgs(int32(parentID)).
					WillReturnRows(sqlmock.NewRows([]string{"sort"}).AddRow(5))

				params := request.ToParams(userID)
				mock.ExpectExec(`INSERT INTO app_menu`).
					WithArgs(
						sql.NullInt32{Int32: int32(parentID), Valid: true},
						params.Name,
						params.Description,
						params.URL,
						5,
						params.Group,
						params.Icon,
						params.Active,
						params.Display,
						params.CreatedBy,
					).WillReturnResult(sqlmock.NewResult(1, 1))

				mock.ExpectCommit()
			},
		},
		{
			name: "success without parent ID",
			request: menuPresenter.CreateMenuRequest{
				Name:    "Home",
				Group:   "main",
				URL:     getPtrString("/home"),
				Active:  getPtrBool(true),
				Display: getPtrBool(true),
			},
			userID:      "user-123",
			wantErr:     false,
			expectedErr: nil,
			setupMock: func(mock sqlmock.Sqlmock, request menuPresenter.CreateMenuRequest, userID string) {
				mock.ExpectBegin()

				mock.ExpectQuery(`SELECT COALESCE\(MAX\(sort\), -1\) \+ 1 FROM app_menu WHERE parent_id IS NULL AND "group" = \$1`).
					WithArgs(request.Group).
					WillReturnRows(sqlmock.NewRows([]string{"sort"}).AddRow(2))

				// Mock insert menu
				params := request.ToParams(userID)
				mock.ExpectExec(`INSERT INTO app_menu`).
					WithArgs(
						sql.NullInt32{Valid: false}, // No parent
						params.Name,
						params.Description,
						params.URL,
						2,
						params.Group,
						params.Icon,
						params.Active,
						params.Display,
						params.CreatedBy,
					).WillReturnResult(sqlmock.NewResult(1, 1))

				mock.ExpectCommit()
			},
		},
		{
			name: "success with description and icon",
			request: menuPresenter.CreateMenuRequest{
				Name:        "Settings",
				Group:       "main",
				Description: getPtrString("Application settings"),
				URL:         getPtrString("/settings"),
				Icon:        getPtrString("settings-icon"),
				Active:      getPtrBool(true),
				Display:     getPtrBool(true),
			},
			userID:      "user-123",
			wantErr:     false,
			expectedErr: nil,
			setupMock: func(mock sqlmock.Sqlmock, request menuPresenter.CreateMenuRequest, userID string) {
				mock.ExpectBegin()

				mock.ExpectQuery(`SELECT COALESCE\(MAX\(sort\), -1\) \+ 1 FROM app_menu WHERE parent_id IS NULL AND "group" = \$1`).
					WithArgs(request.Group).
					WillReturnRows(sqlmock.NewRows([]string{"sort"}).AddRow(3))

				params := request.ToParams(userID)
				mock.ExpectExec(`INSERT INTO app_menu`).
					WithArgs(
						sql.NullInt32{Valid: false},
						params.Name,
						sql.NullString{String: "Application settings", Valid: true},
						sql.NullString{String: "/settings", Valid: true},
						3,
						params.Group,
						sql.NullString{String: "settings-icon", Valid: true},
						params.Active,
						params.Display,
						params.CreatedBy,
					).WillReturnResult(sqlmock.NewResult(1, 1))

				mock.ExpectCommit()
			},
		},
		{
			name: "success with URL without leading slash",
			request: menuPresenter.CreateMenuRequest{
				Name:    "Profile",
				Group:   "main",
				URL:     getPtrString("profile"), // No leading slash
				Active:  getPtrBool(true),
				Display: getPtrBool(true),
			},
			userID:      "user-123",
			wantErr:     false,
			expectedErr: nil,
			setupMock: func(mock sqlmock.Sqlmock, request menuPresenter.CreateMenuRequest, userID string) {
				mock.ExpectBegin()

				mock.ExpectQuery(`SELECT COALESCE\(MAX\(sort\), -1\) \+ 1 FROM app_menu WHERE parent_id IS NULL AND "group" = \$1`).
					WithArgs(request.Group).
					WillReturnRows(sqlmock.NewRows([]string{"sort"}).AddRow(4))

				params := request.ToParams(userID)
				mock.ExpectExec(`INSERT INTO app_menu`).
					WithArgs(
						sql.NullInt32{Valid: false},
						params.Name,
						params.Description,
						sql.NullString{String: "/profile", Valid: true},
						4,
						params.Group,
						params.Icon,
						params.Active,
						params.Display,
						params.CreatedBy,
					).WillReturnResult(sqlmock.NewResult(1, 1))

				mock.ExpectCommit()
			},
		},
		{
			name: "error beginning transaction",
			request: menuPresenter.CreateMenuRequest{
				Name:  "Error Tx",
				Group: "main",
			},
			userID:      "user-123",
			wantErr:     true,
			expectedErr: constant.ErrUnknownSource,
			setupMock: func(mock sqlmock.Sqlmock, request menuPresenter.CreateMenuRequest, userID string) {
				mock.ExpectBegin().WillReturnError(errors.New("begin tx error"))
			},
		},
		{
			name: "error reading sort for parent",
			request: menuPresenter.CreateMenuRequest{
				Name:     "Error Sort Parent",
				Group:    "main",
				ParentID: getPtrInt(1),
			},
			userID:      "user-123",
			wantErr:     true,
			expectedErr: constant.ErrUnknownSource,
			setupMock: func(mock sqlmock.Sqlmock, request menuPresenter.CreateMenuRequest, userID string) {
				mock.ExpectBegin()

				parentID := *request.ParentID
				mock.ExpectQuery(`SELECT COALESCE\(MAX\(sort\), -1\) \+ 1 FROM app_menu WHERE parent_id = \$1`).
					WithArgs(int32(parentID)).
					WillReturnError(errors.New("query error"))

				mock.ExpectRollback()
			},
		},
		{
			name: "error reading sort for group",
			request: menuPresenter.CreateMenuRequest{
				Name:  "Error Sort Group",
				Group: "main",
			},
			userID:      "user-123",
			wantErr:     true,
			expectedErr: constant.ErrUnknownSource,
			setupMock: func(mock sqlmock.Sqlmock, request menuPresenter.CreateMenuRequest, userID string) {
				mock.ExpectBegin()

				mock.ExpectQuery(`SELECT COALESCE\(MAX\(sort\), -1\) \+ 1 FROM app_menu WHERE parent_id IS NULL AND "group" = \$1`).
					WithArgs(request.Group).
					WillReturnError(errors.New("query error"))

				mock.ExpectRollback()
			},
		},
		{
			name: "error creating menu query",
			request: menuPresenter.CreateMenuRequest{
				Name:  "Error Insert",
				Group: "main",
			},
			userID:      "user-123",
			wantErr:     true,
			expectedErr: constant.ErrUnknownSource,
			setupMock: func(mock sqlmock.Sqlmock, request menuPresenter.CreateMenuRequest, userID string) {
				mock.ExpectBegin()

				mock.ExpectQuery(`SELECT COALESCE\(MAX\(sort\), -1\) \+ 1 FROM app_menu WHERE parent_id IS NULL AND "group" = \$1`).
					WithArgs(request.Group).
					WillReturnRows(sqlmock.NewRows([]string{"sort"}).AddRow(3))

				mock.ExpectExec(`INSERT INTO app_menu`).
					WillReturnError(errors.New("insert error"))

				mock.ExpectRollback()
			},
		},
		{
			name: "error committing transaction",
			request: menuPresenter.CreateMenuRequest{
				Name:  "Error Commit",
				Group: "main",
			},
			userID:      "user-123",
			wantErr:     true,
			expectedErr: constant.ErrUnknownSource,
			setupMock: func(mock sqlmock.Sqlmock, request menuPresenter.CreateMenuRequest, userID string) {
				mock.ExpectBegin()

				mock.ExpectQuery(`SELECT COALESCE\(MAX\(sort\), -1\) \+ 1 FROM app_menu WHERE parent_id IS NULL AND "group" = \$1`).
					WithArgs(request.Group).
					WillReturnRows(sqlmock.NewRows([]string{"sort"}).AddRow(10))

				mock.ExpectExec(`INSERT INTO app_menu`).
					WillReturnResult(sqlmock.NewResult(1, 1))

				mock.ExpectCommit().WillReturnError(errors.New("commit error"))
			},
		},
		{
			name: "error rolling back transaction",
			request: menuPresenter.CreateMenuRequest{
				Name:  "Error Rollback",
				Group: "main",
			},
			userID:      "user-123",
			wantErr:     true,
			expectedErr: constant.ErrUnknownSource,
			setupMock: func(mock sqlmock.Sqlmock, request menuPresenter.CreateMenuRequest, userID string) {
				mock.ExpectBegin()

				mock.ExpectQuery(`SELECT COALESCE\(MAX\(sort\), -1\) \+ 1 FROM app_menu WHERE parent_id IS NULL AND "group" = \$1`).
					WithArgs(request.Group).
					WillReturnError(errors.New("query error"))

				// Rollback also fails
				mock.ExpectRollback().WillReturnError(errors.New("rollback failed"))
			},
		},
		{
			name: "default values when optional fields are nil",
			request: menuPresenter.CreateMenuRequest{
				Name:  "Default Values",
				Group: "main",
			},
			userID:      "user-123",
			wantErr:     false,
			expectedErr: nil,
			setupMock: func(mock sqlmock.Sqlmock, request menuPresenter.CreateMenuRequest, userID string) {
				mock.ExpectBegin()

				mock.ExpectQuery(`SELECT COALESCE\(MAX\(sort\), -1\) \+ 1 FROM app_menu WHERE parent_id IS NULL AND "group" = \$1`).
					WithArgs(request.Group).
					WillReturnRows(sqlmock.NewRows([]string{"sort"}).AddRow(1))

				params := request.ToParams(userID)
				mock.ExpectExec(`INSERT INTO app_menu`).
					WithArgs(
						sql.NullInt32{Valid: false},
						params.Name,
						sql.NullString{Valid: false},
						sql.NullString{Valid: false},
						1,
						params.Group,
						sql.NullString{Valid: false},
						params.Active,
						params.Display,
						params.CreatedBy,
					).WillReturnResult(sqlmock.NewResult(1, 1))

				mock.ExpectCommit()
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			svc, mock, cleanup := newTestService(t)
			defer cleanup()

			ctx := context.Background()

			tc.setupMock(mock, tc.request, tc.userID)

			err := svc.CreateMenuService(ctx, tc.request, tc.userID)

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
