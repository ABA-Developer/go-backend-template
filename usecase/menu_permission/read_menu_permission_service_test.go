package menu_permission

import (
	"be-dashboard-nba/api/presenter"
	menuPermissionPresenter "be-dashboard-nba/api/presenter/menu_permission"
	"be-dashboard-nba/constant"
	"be-dashboard-nba/internal/utils"
	"be-dashboard-nba/usecase/entities"
	"context"
	"database/sql"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestReadMenuPermissionService(t *testing.T) {
	tests := []struct {
		name          string
		search        string
		page          int
		limit         int
		menuID        int
		expectedCount int
		expectedData  []entities.MenuPermission
		wantErr       bool
		expectedErr   error
		setupMock     func(mock sqlmock.Sqlmock, menuID int, search string, page int, limit int, expectedCount int, expectedData []entities.MenuPermission)
	}{
		{
			name:          "success without search",
			search:        "",
			page:          constant.DefaultPage,
			limit:         constant.DefaultLimit,
			menuID:        1,
			expectedCount: 1,
			expectedData: []entities.MenuPermission{
				{
					ID:         1,
					MenuID:     1,
					Code:       "R",
					ActionName: "read",
				},
			},
			wantErr:     false,
			expectedErr: nil,
			setupMock: func(mock sqlmock.Sqlmock, menuID int, search string, page int, limit int, expectedCount int, expectedData []entities.MenuPermission) {
				mockMenuExists(mock, menuID)

				expectedSearch := search
				if search != "" {
					expectedSearch = "%" + search + "%"
				}
				expectedSetSearch := search != ""
				expectedOffset := (page * limit) - limit

				mock.ExpectQuery(`SELECT COUNT`).
					WithArgs(expectedSetSearch, expectedSearch, menuID).
					WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(expectedCount))

				rows := sqlmock.NewRows([]string{"id", "menu_id", "code", "action_name"})
				for _, data := range expectedData {
					rows.AddRow(data.ID, data.MenuID, data.Code, data.ActionName)
				}

				mock.ExpectQuery(`SELECT\s+id,\s*menu_id,\s*code,\s*action_name`).
					WithArgs(
						expectedSetSearch,
						expectedSearch,
						"created_at DESC",
						limit,
						expectedOffset,
						menuID,
					).
					WillReturnRows(rows)
			},
		},
		{
			name:          "success with search",
			search:        "read",
			page:          constant.DefaultPage,
			limit:         constant.DefaultLimit,
			menuID:        1,
			expectedCount: 2,
			expectedData: []entities.MenuPermission{
				{
					ID:         1,
					MenuID:     1,
					Code:       "R",
					ActionName: "read",
				},
				{
					ID:         2,
					MenuID:     1,
					Code:       "U",
					ActionName: "update",
				},
			},
			wantErr:     false,
			expectedErr: nil,
			setupMock: func(mock sqlmock.Sqlmock, menuID int, search string, page int, limit int, expectedCount int, expectedData []entities.MenuPermission) {
				mockMenuExists(mock, menuID)

				expectedSearch := search
				if search != "" {
					expectedSearch = "%" + search + "%"
				}
				expectedSetSearch := search != ""
				expectedOffset := (page * limit) - limit

				mock.ExpectQuery(`SELECT COUNT`).
					WithArgs(expectedSetSearch, expectedSearch, menuID).
					WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(expectedCount))

				rows := sqlmock.NewRows([]string{"id", "menu_id", "code", "action_name"})
				for _, data := range expectedData {
					rows.AddRow(data.ID, data.MenuID, data.Code, data.ActionName)
				}

				mock.ExpectQuery(`SELECT\s+id,\s*menu_id,\s*code,\s*action_name`).
					WithArgs(
						expectedSetSearch,
						expectedSearch,
						"created_at DESC",
						limit,
						expectedOffset,
						menuID,
					).
					WillReturnRows(rows)
			},
		},
		{
			name:          "pagination page 2",
			search:        "",
			page:          2,
			limit:         5,
			menuID:        1,
			expectedCount: 8,
			expectedData: []entities.MenuPermission{
				{
					ID:         6,
					MenuID:     1,
					Code:       "D",
					ActionName: "delete",
				},
				{
					ID:         7,
					MenuID:     1,
					Code:       "E",
					ActionName: "export",
				},
			},
			wantErr:     false,
			expectedErr: nil,
			setupMock: func(mock sqlmock.Sqlmock, menuID int, search string, page int, limit int, expectedCount int, expectedData []entities.MenuPermission) {
				mockMenuExists(mock, menuID)

				expectedSearch := search
				if search != "" {
					expectedSearch = "%" + search + "%"
				}
				expectedSetSearch := search != ""
				expectedOffset := (page * limit) - limit

				mock.ExpectQuery(`SELECT COUNT`).
					WithArgs(expectedSetSearch, expectedSearch, menuID).
					WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(expectedCount))

				rows := sqlmock.NewRows([]string{"id", "menu_id", "code", "action_name"})
				for _, data := range expectedData {
					rows.AddRow(data.ID, data.MenuID, data.Code, data.ActionName)
				}

				mock.ExpectQuery(`SELECT\s+id,\s*menu_id,\s*code,\s*action_name`).
					WithArgs(
						expectedSetSearch,
						expectedSearch,
						"created_at DESC",
						limit,
						expectedOffset,
						menuID,
					).
					WillReturnRows(rows)
			},
		},
		{
			name:          "empty result",
			search:        "nonexistent",
			page:          constant.DefaultPage,
			limit:         constant.DefaultLimit,
			menuID:        1,
			expectedCount: 0,
			expectedData:  nil,
			wantErr:       false,
			expectedErr:   nil,
			setupMock: func(mock sqlmock.Sqlmock, menuID int, search string, page int, limit int, expectedCount int, expectedData []entities.MenuPermission) {
				mockMenuExists(mock, menuID)

				expectedSearch := search
				if search != "" {
					expectedSearch = "%" + search + "%"
				}
				expectedSetSearch := search != ""
				expectedOffset := (page * limit) - limit

				mock.ExpectQuery(`SELECT COUNT`).
					WithArgs(expectedSetSearch, expectedSearch, menuID).
					WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(expectedCount))

				rows := sqlmock.NewRows([]string{"id", "menu_id", "code", "action_name"})

				mock.ExpectQuery(`SELECT\s+id,\s*menu_id,\s*code,\s*action_name`).
					WithArgs(
						expectedSetSearch,
						expectedSearch,
						"created_at DESC",
						limit,
						expectedOffset,
						menuID,
					).
					WillReturnRows(rows)
			},
		},
		{
			name:          "menu not found",
			search:        "",
			page:          constant.DefaultPage,
			limit:         constant.DefaultLimit,
			menuID:        999,
			expectedCount: 0,
			expectedData:  nil,
			wantErr:       true,
			expectedErr:   constant.ErrMenuIdNotFound,
			setupMock: func(mock sqlmock.Sqlmock, menuID int, search string, page int, limit int, expectedCount int, expectedData []entities.MenuPermission) {
				mock.ExpectQuery(`(?s)SELECT.*FROM.*app_menu.*WHERE.*id`).
					WithArgs(menuID).
					WillReturnError(sql.ErrNoRows)

			},
		},
		{
			name:          "database error on menu check",
			search:        "",
			page:          constant.DefaultPage,
			limit:         constant.DefaultLimit,
			menuID:        1,
			expectedCount: 0,
			expectedData:  nil,
			wantErr:       true,
			expectedErr:   constant.ErrUnknownSource,
			setupMock: func(mock sqlmock.Sqlmock, menuID int, search string, page int, limit int, expectedCount int, expectedData []entities.MenuPermission) {
				mock.ExpectQuery(`(?s)SELECT.*FROM.*app_menu.*WHERE.*id`).
					WithArgs(menuID).
					WillReturnError(sql.ErrConnDone)

			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			svc, mock, cleanup := newTestService(t)
			defer cleanup()

			ctx := context.Background()

			args := menuPermissionPresenter.ReadMenuPermissionListRequest{
				PaginationPayload: utils.PaginationPayload{
					Search: tc.search,
					Page:   tc.page,
					Limit:  tc.limit,
				},
			}

			tc.setupMock(mock, tc.menuID, tc.search, tc.page, tc.limit, tc.expectedCount, tc.expectedData)

			resp, err := svc.ReadMenuPermissionService(ctx, args, tc.menuID)

			if tc.wantErr {
				assert.Error(t, err)
				if tc.expectedErr != nil {
					assert.ErrorIs(t, err, tc.expectedErr)
				}
				assert.Equal(t, entities.MenuPermissionPaginationResponse{}, resp)
			} else {
				assert.NoError(t, err)

				expectedPagination := presenter.Pagination{
					Page:       tc.page,
					PageSize:   tc.limit,
					TotalItems: tc.expectedCount,
					TotalPages: (tc.expectedCount + tc.limit - 1) / tc.limit,
					HasNext:    tc.page < ((tc.expectedCount + tc.limit - 1) / tc.limit),
					HasPrev:    tc.page > 1,
				}

				assert.Equal(t, expectedPagination, resp.Pagination)
				assert.Equal(t, tc.expectedCount, resp.Pagination.TotalItems)
				assert.Equal(t, len(tc.expectedData), len(resp.Data))

				if tc.expectedData == nil {
					assert.Nil(t, resp.Data)
				} else {
					assert.Equal(t, tc.expectedData, resp.Data)
				}
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func mockMenuExists(mock sqlmock.Sqlmock, menuID int) {
	mock.ExpectQuery(`(?s)SELECT.*FROM.*app_menu.*WHERE.*id`).
		WithArgs(menuID).
		WillReturnRows(sqlmock.NewRows([]string{
			"id", "parent_id", "name", "description", "url", "sort", "group", "icon", "active", "display",
		}).AddRow(
			menuID, nil, "Dashboard", nil, "/dashboard", 0, "main", nil, true, true,
		))
}
