package role

import (
	"be-dashboard-nba/api/presenter"
	rolePresenter "be-dashboard-nba/api/presenter/role"
	"be-dashboard-nba/constant"
	"be-dashboard-nba/internal/utils"
	"be-dashboard-nba/usecase/entities"
	"context"
	"database/sql"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestReadRolesService(t *testing.T) {
	tests := []struct {
		name          string
		search        string
		page          int
		limit         int
		order         string
		expectedCount int
		expectedData  []entities.Role
		wantErr       bool
		expectedErr   error
		setupMock     func(mock sqlmock.Sqlmock, search string, page int, limit int, order string, expectedCount int, expectedData []entities.Role)
	}{
		{
			name:          "success without search",
			search:        "",
			page:          constant.DefaultPage,
			limit:         constant.DefaultLimit,
			order:         "id DESC",
			expectedCount: 2,
			expectedData: []entities.Role{
				{
					ID:   1,
					Name: "Admin",
					Code: "ADMIN",
					Description: sql.NullString{
						String: "Administrator Role",
						Valid:  true,
					},
				},
				{
					ID:   2,
					Name: "User",
					Code: "USER",
					Description: sql.NullString{
						String: "User Role",
						Valid:  true,
					},
				},
			},
			wantErr:     false,
			expectedErr: nil,
			setupMock: func(mock sqlmock.Sqlmock, search string, page int, limit int, order string, expectedCount int, expectedData []entities.Role) {
				expectedSearch := search
				if search != "" {
					expectedSearch = "%" + search + "%"
				}
				expectedSetSearch := search != ""
				expectedOffset := (page * limit) - limit

				// Mock count query - use more flexible regex
				mock.ExpectQuery(`SELECT COUNT\(\*\) FROM app_role`).
					WithArgs(expectedSetSearch, expectedSearch).
					WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(expectedCount))

				// Mock data query - use more flexible regex
				rows := sqlmock.NewRows([]string{"id", "name", "code", "description"})
				for _, data := range expectedData {
					if data.Description.Valid {
						rows.AddRow(data.ID, data.Name, data.Code, data.Description.String)
					} else {
						rows.AddRow(data.ID, data.Name, data.Code, nil)
					}
				}

				// Use more flexible regex pattern that matches the actual query structure
				mock.ExpectQuery(`SELECT id, name, code ,description FROM app_role`).
					WithArgs(
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
			search:        "admin",
			page:          constant.DefaultPage,
			limit:         constant.DefaultLimit,
			order:         "id DESC",
			expectedCount: 1,
			expectedData: []entities.Role{
				{
					ID:   1,
					Name: "Admin",
					Code: "ADMIN",
					Description: sql.NullString{
						String: "Administrator Role",
						Valid:  true,
					},
				},
			},
			wantErr:     false,
			expectedErr: nil,
			setupMock: func(mock sqlmock.Sqlmock, search string, page int, limit int, order string, expectedCount int, expectedData []entities.Role) {
				expectedSearch := "%" + search + "%"
				expectedSetSearch := true
				expectedOffset := (page * limit) - limit

				// Mock count query
				mock.ExpectQuery(`SELECT COUNT\(\*\) FROM app_role`).
					WithArgs(expectedSetSearch, expectedSearch).
					WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(expectedCount))

				// Mock data query
				rows := sqlmock.NewRows([]string{"id", "name", "code", "description"})
				for _, data := range expectedData {
					if data.Description.Valid {
						rows.AddRow(data.ID, data.Name, data.Code, data.Description.String)
					} else {
						rows.AddRow(data.ID, data.Name, data.Code, nil)
					}
				}

				mock.ExpectQuery(`SELECT id, name, code ,description FROM app_role`).
					WithArgs(
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
			search:        "",
			page:          2,
			limit:         5,
			order:         "id DESC",
			expectedCount: 12,
			expectedData: []entities.Role{
				{
					ID:   6,
					Name: "Manager",
					Code: "MANAGER",
					Description: sql.NullString{
						String: "Manager Role",
						Valid:  true,
					},
				},
				{
					ID:   7,
					Name: "Supervisor",
					Code: "SUPERVISOR",
					Description: sql.NullString{
						String: "Supervisor Role",
						Valid:  true,
					},
				},
			},
			wantErr:     false,
			expectedErr: nil,
			setupMock: func(mock sqlmock.Sqlmock, search string, page int, limit int, order string, expectedCount int, expectedData []entities.Role) {
				expectedSearch := search
				if search != "" {
					expectedSearch = "%" + search + "%"
				}
				expectedSetSearch := search != ""
				expectedOffset := (page * limit) - limit

				// Mock count query
				mock.ExpectQuery(`SELECT COUNT\(\*\) FROM app_role`).
					WithArgs(expectedSetSearch, expectedSearch).
					WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(expectedCount))

				// Mock data query
				rows := sqlmock.NewRows([]string{"id", "name", "code", "description"})
				for _, data := range expectedData {
					if data.Description.Valid {
						rows.AddRow(data.ID, data.Name, data.Code, data.Description.String)
					} else {
						rows.AddRow(data.ID, data.Name, data.Code, nil)
					}
				}

				mock.ExpectQuery(`SELECT id, name, code ,description FROM app_role`).
					WithArgs(
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
			search:        "nonexistent",
			page:          constant.DefaultPage,
			limit:         constant.DefaultLimit,
			order:         "id DESC",
			expectedCount: 0,
			expectedData:  []entities.Role{},
			wantErr:       false,
			expectedErr:   nil,
			setupMock: func(mock sqlmock.Sqlmock, search string, page int, limit int, order string, expectedCount int, expectedData []entities.Role) {
				expectedSearch := "%" + search + "%"
				expectedSetSearch := true
				expectedOffset := (page * limit) - limit

				// Mock count query
				mock.ExpectQuery(`SELECT COUNT\(\*\) FROM app_role`).
					WithArgs(expectedSetSearch, expectedSearch).
					WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(expectedCount))

				// Mock data query - return empty rows
				rows := sqlmock.NewRows([]string{"id", "name", "code", "description"})

				mock.ExpectQuery(`SELECT id, name, code ,description FROM app_role`).
					WithArgs(
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
			name:          "database error on count query",
			search:        "",
			page:          constant.DefaultPage,
			limit:         constant.DefaultLimit,
			order:         "id DESC",
			expectedCount: 0,
			expectedData:  nil,
			wantErr:       true,
			expectedErr:   constant.ErrUnknownSource,
			setupMock: func(mock sqlmock.Sqlmock, search string, page int, limit int, order string, expectedCount int, expectedData []entities.Role) {
				expectedSearch := search
				expectedSetSearch := search != ""

				// Mock count query to return error
				mock.ExpectQuery(`SELECT COUNT\(\*\) FROM app_role`).
					WithArgs(expectedSetSearch, expectedSearch).
					WillReturnError(sql.ErrConnDone)
			},
		},
		{
			name:          "database error on data query",
			search:        "",
			page:          constant.DefaultPage,
			limit:         constant.DefaultLimit,
			order:         "id DESC",
			expectedCount: 5,
			expectedData:  nil,
			wantErr:       true,
			expectedErr:   constant.ErrUnknownSource,
			setupMock: func(mock sqlmock.Sqlmock, search string, page int, limit int, order string, expectedCount int, expectedData []entities.Role) {
				expectedSearch := search
				expectedSetSearch := search != ""
				expectedOffset := (page * limit) - limit

				// Mock count query success
				mock.ExpectQuery(`SELECT COUNT\(\*\) FROM app_role`).
					WithArgs(expectedSetSearch, expectedSearch).
					WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(expectedCount))

				// Mock data query to return error
				mock.ExpectQuery(`SELECT id, name, code ,description FROM app_role`).
					WithArgs(
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

			sortField := "id"
			direction := "DESC"

			args := rolePresenter.ReadRolesRequest{
				PaginationPayload: utils.PaginationPayload{
					Search:    tc.search,
					Page:      tc.page,
					Limit:     tc.limit,
					Sort:      sortField,
					Direction: direction,
				},
			}

			tc.setupMock(mock, tc.search, tc.page, tc.limit, tc.order, tc.expectedCount, tc.expectedData)

			resp, err := svc.ReadRolesService(ctx, args)

			if tc.wantErr {
				assert.Error(t, err)
				if tc.expectedErr != nil {
					assert.ErrorIs(t, err, tc.expectedErr)
				}
				assert.Equal(t, entities.RolePaginationResponse{}, resp)
			} else {
				assert.NoError(t, err)

				totalPages := 0
				if tc.expectedCount > 0 {
					totalPages = (tc.expectedCount + tc.limit - 1) / tc.limit
				}
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

				if tc.expectedData != nil {
					assert.Equal(t, len(tc.expectedData), len(resp.Data))
					for i, expectedRole := range tc.expectedData {
						assert.Equal(t, expectedRole.ID, resp.Data[i].ID)
						assert.Equal(t, expectedRole.Name, resp.Data[i].Name)
						assert.Equal(t, expectedRole.Code, resp.Data[i].Code)
						assert.Equal(t, expectedRole.Description, resp.Data[i].Description)
					}
				} else {
					assert.Nil(t, resp.Data)
				}
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}
