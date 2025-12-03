package service

import (
	"be-dashboard-nba/api/presenter"
	rolePresenter "be-dashboard-nba/api/presenter/role"
	"be-dashboard-nba/constant"
	"be-dashboard-nba/internal/utils"
	"be-dashboard-nba/pkg/entities"
	"context"
	"database/sql"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestReadRoleAccessService(t *testing.T) {
	tests := []struct {
		name          string
		roleID        int
		search        string
		page          int
		limit         int
		order         string
		expectedCount int
		expectedData  []entities.RoleAccessResponse
		wantErr       bool
		expectedErr   error
		setupMock     func(mock sqlmock.Sqlmock, roleID int, search string, page int, limit int, order string, expectedCount int, expectedData []entities.RoleAccessResponse)
	}{
		{
			name:          "success without search",
			roleID:        1,
			search:        "",
			page:          constant.DefaultPage,
			limit:         constant.DefaultLimit,
			order:         "name ASC",
			expectedCount: 3,
			expectedData: []entities.RoleAccessResponse{
				{
					RoleID:         1,
					RoleName:       "Admin",
					MenuID:         1,
					MenuName:       "Dashboard",
					PermissionID:   1,
					PermissionName: "Read",
					PermissionCode: "R",
					HasAccess:      true,
				},
				{
					RoleID:         1,
					RoleName:       "Admin",
					MenuID:         1,
					MenuName:       "Dashboard",
					PermissionID:   2,
					PermissionName: "Write",
					PermissionCode: "W",
					HasAccess:      true,
				},
				{
					RoleID:         1,
					RoleName:       "Admin",
					MenuID:         2,
					MenuName:       "Users",
					PermissionID:   3,
					PermissionName: "Read",
					PermissionCode: "R",
					HasAccess:      true,
				},
			},
			wantErr:     false,
			expectedErr: nil,
			setupMock: func(mock sqlmock.Sqlmock, roleID int, search string, page int, limit int, order string, expectedCount int, expectedData []entities.RoleAccessResponse) {
				expectedSearch := search
				if search != "" {
					expectedSearch = "%" + search + "%"
				}
				expectedSetSearch := search != ""
				expectedOffset := (page * limit) - limit

				// 1. Mock count query (dipanggil pertama di service)
				mock.ExpectQuery(`SELECT COUNT\(\*\) FROM app_menu`).
					WithArgs(expectedSetSearch, expectedSearch).
					WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(expectedCount))

				// 2. Mock role existence check (dipanggil kedua di service)
				mock.ExpectQuery(`SELECT id, name, code, description FROM app_role WHERE id =`).
					WithArgs(roleID).
					WillReturnRows(sqlmock.NewRows([]string{"id", "name", "code", "description"}).
						AddRow(1, "Admin", "ADMIN", "Administrator Role"))

				// 3. Mock data query (dipanggil terakhir di service)
				rows := sqlmock.NewRows([]string{
					"role_id", "role_name", "menu_id", "menu_name",
					"permission_id", "permission_name", "permission_code", "has_access",
				})
				for _, data := range expectedData {
					rows.AddRow(
						data.RoleID,
						data.RoleName,
						data.MenuID,
						data.MenuName,
						data.PermissionID,
						data.PermissionName,
						data.PermissionCode,
						data.HasAccess,
					)
				}

				mock.ExpectQuery(`WITH filtered_menu AS`).
					WithArgs(
						roleID,
						expectedSetSearch,
						expectedSearch,
						order,
						limit,
						expectedOffset,
					).
					WillReturnRows(rows)
			},
		},
		{
			name:          "success with search",
			roleID:        1,
			search:        "dashboard",
			page:          constant.DefaultPage,
			limit:         constant.DefaultLimit,
			order:         "name ASC",
			expectedCount: 1,
			expectedData: []entities.RoleAccessResponse{
				{
					RoleID:         1,
					RoleName:       "Admin",
					MenuID:         1,
					MenuName:       "Dashboard",
					PermissionID:   1,
					PermissionName: "Read",
					PermissionCode: "R",
					HasAccess:      true,
				},
				{
					RoleID:         1,
					RoleName:       "Admin",
					MenuID:         1,
					MenuName:       "Dashboard",
					PermissionID:   2,
					PermissionName: "Write",
					PermissionCode: "W",
					HasAccess:      true,
				},
			},
			wantErr:     false,
			expectedErr: nil,
			setupMock: func(mock sqlmock.Sqlmock, roleID int, search string, page int, limit int, order string, expectedCount int, expectedData []entities.RoleAccessResponse) {
				expectedSearch := "%" + search + "%"
				expectedSetSearch := true
				expectedOffset := (page * limit) - limit

				// 1. Count query
				mock.ExpectQuery(`SELECT COUNT\(\*\) FROM app_menu`).
					WithArgs(expectedSetSearch, expectedSearch).
					WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(expectedCount))

				// 2. Role existence check
				mock.ExpectQuery(`SELECT id, name, code, description FROM app_role WHERE id =`).
					WithArgs(roleID).
					WillReturnRows(sqlmock.NewRows([]string{"id", "name", "code", "description"}).
						AddRow(1, "Admin", "ADMIN", "Administrator Role"))

				// 3. Data query
				rows := sqlmock.NewRows([]string{
					"role_id", "role_name", "menu_id", "menu_name",
					"permission_id", "permission_name", "permission_code", "has_access",
				})
				for _, data := range expectedData {
					rows.AddRow(
						data.RoleID,
						data.RoleName,
						data.MenuID,
						data.MenuName,
						data.PermissionID,
						data.PermissionName,
						data.PermissionCode,
						data.HasAccess,
					)
				}

				mock.ExpectQuery(`WITH filtered_menu AS`).
					WithArgs(
						roleID,
						expectedSetSearch,
						expectedSearch,
						order,
						limit,
						expectedOffset,
					).
					WillReturnRows(rows)
			},
		},
		{
			name:          "pagination page 2",
			roleID:        1,
			search:        "",
			page:          2,
			limit:         5,
			order:         "name ASC",
			expectedCount: 12,
			expectedData: []entities.RoleAccessResponse{
				{
					RoleID:         1,
					RoleName:       "Admin",
					MenuID:         6,
					MenuName:       "Reports",
					PermissionID:   11,
					PermissionName: "Read",
					PermissionCode: "R",
					HasAccess:      true,
				},
				{
					RoleID:         1,
					RoleName:       "Admin",
					MenuID:         6,
					MenuName:       "Reports",
					PermissionID:   12,
					PermissionName: "Export",
					PermissionCode: "E",
					HasAccess:      false,
				},
			},
			wantErr:     false,
			expectedErr: nil,
			setupMock: func(mock sqlmock.Sqlmock, roleID int, search string, page int, limit int, order string, expectedCount int, expectedData []entities.RoleAccessResponse) {
				expectedSearch := search
				if search != "" {
					expectedSearch = "%" + search + "%"
				}
				expectedSetSearch := search != ""
				expectedOffset := (page * limit) - limit

				// 1. Count query
				mock.ExpectQuery(`SELECT COUNT\(\*\) FROM app_menu`).
					WithArgs(expectedSetSearch, expectedSearch).
					WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(expectedCount))

				// 2. Role existence check
				mock.ExpectQuery(`SELECT id, name, code, description FROM app_role WHERE id =`).
					WithArgs(roleID).
					WillReturnRows(sqlmock.NewRows([]string{"id", "name", "code", "description"}).
						AddRow(1, "Admin", "ADMIN", "Administrator Role"))

				// 3. Data query
				rows := sqlmock.NewRows([]string{
					"role_id", "role_name", "menu_id", "menu_name",
					"permission_id", "permission_name", "permission_code", "has_access",
				})
				for _, data := range expectedData {
					rows.AddRow(
						data.RoleID,
						data.RoleName,
						data.MenuID,
						data.MenuName,
						data.PermissionID,
						data.PermissionName,
						data.PermissionCode,
						data.HasAccess,
					)
				}

				mock.ExpectQuery(`WITH filtered_menu AS`).
					WithArgs(
						roleID,
						expectedSetSearch,
						expectedSearch,
						order,
						limit,
						expectedOffset,
					).
					WillReturnRows(rows)
			},
		},
		{
			name:          "empty result",
			roleID:        1,
			search:        "nonexistent",
			page:          constant.DefaultPage,
			limit:         constant.DefaultLimit,
			order:         "name ASC",
			expectedCount: 0,
			expectedData:  nil,
			wantErr:       false,
			expectedErr:   nil,
			setupMock: func(mock sqlmock.Sqlmock, roleID int, search string, page int, limit int, order string, expectedCount int, expectedData []entities.RoleAccessResponse) {
				expectedSearch := "%" + search + "%"
				expectedSetSearch := true
				expectedOffset := (page * limit) - limit

				// 1. Count query
				mock.ExpectQuery(`SELECT COUNT\(\*\) FROM app_menu`).
					WithArgs(expectedSetSearch, expectedSearch).
					WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(expectedCount))

				// 2. Role existence check
				mock.ExpectQuery(`SELECT id, name, code, description FROM app_role WHERE id =`).
					WithArgs(roleID).
					WillReturnRows(sqlmock.NewRows([]string{"id", "name", "code", "description"}).
						AddRow(1, "Admin", "ADMIN", "Administrator Role"))

				// 3. Data query - empty rows
				rows := sqlmock.NewRows([]string{
					"role_id", "role_name", "menu_id", "menu_name",
					"permission_id", "permission_name", "permission_code", "has_access",
				})

				mock.ExpectQuery(`WITH filtered_menu AS`).
					WithArgs(
						roleID,
						expectedSetSearch,
						expectedSearch,
						order,
						limit,
						expectedOffset,
					).
					WillReturnRows(rows)
			},
		},
		{
			name:          "role not found",
			roleID:        999,
			search:        "",
			page:          constant.DefaultPage,
			limit:         constant.DefaultLimit,
			order:         "name ASC",
			expectedCount: 0,
			expectedData:  nil,
			wantErr:       true,
			expectedErr:   constant.ErrRoleIdNotFound,
			setupMock: func(mock sqlmock.Sqlmock, roleID int, search string, page int, limit int, order string, expectedCount int, expectedData []entities.RoleAccessResponse) {
				// 1. Count query
				mock.ExpectQuery(`SELECT COUNT\(\*\) FROM app_menu`).
					WithArgs(false, "").
					WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(5))

				// 2. Role existence check to return no rows
				mock.ExpectQuery(`SELECT id, name, code, description FROM app_role WHERE id =`).
					WithArgs(roleID).
					WillReturnError(sql.ErrNoRows)
			},
		},
		{
			name:          "database error on role check",
			roleID:        1,
			search:        "",
			page:          constant.DefaultPage,
			limit:         constant.DefaultLimit,
			order:         "name ASC",
			expectedCount: 0,
			expectedData:  nil,
			wantErr:       true,
			expectedErr:   constant.ErrUnknownSource,
			setupMock: func(mock sqlmock.Sqlmock, roleID int, search string, page int, limit int, order string, expectedCount int, expectedData []entities.RoleAccessResponse) {
				// 1. Count query
				mock.ExpectQuery(`SELECT COUNT\(\*\) FROM app_menu`).
					WithArgs(false, "").
					WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(5))

				// 2. Role existence check to return error
				mock.ExpectQuery(`SELECT id, name, code, description FROM app_role WHERE id =`).
					WithArgs(roleID).
					WillReturnError(sql.ErrConnDone)
			},
		},
		{
			name:          "database error on count query",
			roleID:        1,
			search:        "",
			page:          constant.DefaultPage,
			limit:         constant.DefaultLimit,
			order:         "name ASC",
			expectedCount: 0,
			expectedData:  nil,
			wantErr:       true,
			expectedErr:   constant.ErrUnknownSource,
			setupMock: func(mock sqlmock.Sqlmock, roleID int, search string, page int, limit int, order string, expectedCount int, expectedData []entities.RoleAccessResponse) {
				expectedSearch := search
				expectedSetSearch := search != ""

				// 1. Count query to return error
				mock.ExpectQuery(`SELECT COUNT\(\*\) FROM app_menu`).
					WithArgs(expectedSetSearch, expectedSearch).
					WillReturnError(sql.ErrConnDone)
			},
		},
		{
			name:          "database error on data query",
			roleID:        1,
			search:        "",
			page:          constant.DefaultPage,
			limit:         constant.DefaultLimit,
			order:         "name ASC",
			expectedCount: 5,
			expectedData:  nil,
			wantErr:       true,
			expectedErr:   constant.ErrUnknownSource,
			setupMock: func(mock sqlmock.Sqlmock, roleID int, search string, page int, limit int, order string, expectedCount int, expectedData []entities.RoleAccessResponse) {
				expectedSearch := search
				expectedSetSearch := search != ""
				expectedOffset := (page * limit) - limit

				// 1. Count query success
				mock.ExpectQuery(`SELECT COUNT\(\*\) FROM app_menu`).
					WithArgs(expectedSetSearch, expectedSearch).
					WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(expectedCount))

				// 2. Role existence check
				mock.ExpectQuery(`SELECT id, name, code, description FROM app_role WHERE id =`).
					WithArgs(roleID).
					WillReturnRows(sqlmock.NewRows([]string{"id", "name", "code", "description"}).
						AddRow(1, "Admin", "ADMIN", "Administrator Role"))

				// 3. Data query to return error
				mock.ExpectQuery(`WITH filtered_menu AS`).
					WithArgs(
						roleID,
						expectedSetSearch,
						expectedSearch,
						order,
						limit,
						expectedOffset,
					).
					WillReturnError(sql.ErrConnDone)
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			svc, mock, cleanup := newTestService(t)
			defer cleanup()

			ctx := context.Background()

			args := rolePresenter.ReadRoleAccessesRequest{
				PaginationPayload: utils.PaginationPayload{
					Search:    tc.search,
					Page:      tc.page,
					Limit:     tc.limit,
					Sort:      "name",
					Direction: "ASC",
				},
			}

			tc.setupMock(mock, tc.roleID, tc.search, tc.page, tc.limit, tc.order, tc.expectedCount, tc.expectedData)

			resp, err := svc.ReadRoleAccessService(ctx, args, tc.roleID)

			if tc.wantErr {
				assert.Error(t, err)
				if tc.expectedErr != nil {
					assert.ErrorIs(t, err, tc.expectedErr)
				}
				assert.Equal(t, entities.RoleAccessPaginationResponse{}, resp)
			} else {
				assert.NoError(t, err)

				// Calculate expected pagination - sesuai dengan perhitungan di service
				totalPages := 0
				if tc.expectedCount > 0 {
					totalPages = (tc.expectedCount + tc.limit - 1) / tc.limit
				}
				// Untuk empty result, totalPages akan 0 sesuai dengan math.Ceil(0/limit) = 0
				hasNext := tc.page < totalPages
				hasPrev := tc.page > 1

				expectedPagination := presenter.Pagination{
					Page:       tc.page,
					PageSize:   tc.limit,
					TotalItems: tc.expectedCount,
					TotalPages: totalPages,
					HasNext:    hasNext,
					HasPrev:    hasPrev,
				}

				assert.Equal(t, expectedPagination, resp.Pagination)
				assert.Equal(t, tc.expectedCount, resp.Pagination.TotalItems)
				assert.Equal(t, len(tc.expectedData), len(resp.Data))

				// Compare each role access individually
				if len(tc.expectedData) > 0 {
					for i, expected := range tc.expectedData {
						assert.Equal(t, expected.RoleID, resp.Data[i].RoleID)
						assert.Equal(t, expected.RoleName, resp.Data[i].RoleName)
						assert.Equal(t, expected.MenuID, resp.Data[i].MenuID)
						assert.Equal(t, expected.MenuName, resp.Data[i].MenuName)
						assert.Equal(t, expected.PermissionID, resp.Data[i].PermissionID)
						assert.Equal(t, expected.PermissionName, resp.Data[i].PermissionName)
						assert.Equal(t, expected.PermissionCode, resp.Data[i].PermissionCode)
						assert.Equal(t, expected.HasAccess, resp.Data[i].HasAccess)
					}
				} else {
					// Untuk empty result, pastikan data adalah empty slice, bukan nil
					assert.Nil(t, resp.Data)
				}
			}

			// Check if all expectations were met
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}
