package service

import (
	"be-dashboard-nba/api/presenter"
	userPresenter "be-dashboard-nba/api/presenter/user"
	"be-dashboard-nba/constant"
	"be-dashboard-nba/internal/utils"
	"be-dashboard-nba/pkg/entities"
	"context"
	"database/sql"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestReadUsersService(t *testing.T) {
	tests := []struct {
		name          string
		search        string
		page          int
		limit         int
		order         string
		expectedCount int
		expectedData  []entities.User
		wantErr       bool
		expectedErr   error
		setupMock     func(mock sqlmock.Sqlmock, search string, page int, limit int, order string, expectedCount int, expectedData []entities.User)
	}{
		{
			name:          "success without search",
			search:        "",
			page:          constant.DefaultPage,
			limit:         constant.DefaultLimit,
			order:         "created_at DESC",
			expectedCount: 2,
			expectedData: []entities.User{
				{
					ID:       "user-1",
					FullName: "John Doe",
					Active:   true,
					Role:     "Admin",
					RoleID:   1,
				},
				{
					ID:       "user-2",
					FullName: "Jane Smith",
					Active:   true,
					Role:     "User",
					RoleID:   2,
				},
			},
			wantErr:     false,
			expectedErr: nil,
			setupMock: func(mock sqlmock.Sqlmock, search string, page int, limit int, order string, expectedCount int, expectedData []entities.User) {
				expectedSearch := search
				if search != "" {
					expectedSearch = "%" + search + "%"
				}
				expectedSetSearch := search != ""
				expectedOffset := (page * limit) - limit

				// Mock count query
				mock.ExpectQuery(`SELECT COUNT\(\*\) FROM app_user`).
					WithArgs(expectedSetSearch, expectedSearch).
					WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(expectedCount))

				// Mock data query
				rows := sqlmock.NewRows([]string{"id", "full_name", "active", "role", "role_id"})
				for _, data := range expectedData {
					rows.AddRow(data.ID, data.FullName, data.Active, data.Role, data.RoleID)
				}

				mock.ExpectQuery(`SELECT`).
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
			search:        "john",
			page:          constant.DefaultPage,
			limit:         constant.DefaultLimit,
			order:         "created_at DESC",
			expectedCount: 1,
			expectedData: []entities.User{
				{
					ID:       "user-1",
					FullName: "John Doe",
					Active:   true,
					Role:     "Admin",
					RoleID:   1,
				},
			},
			wantErr:     false,
			expectedErr: nil,
			setupMock: func(mock sqlmock.Sqlmock, search string, page int, limit int, order string, expectedCount int, expectedData []entities.User) {
				expectedSearch := "%" + search + "%"
				expectedSetSearch := true
				expectedOffset := (page * limit) - limit

				// Mock count query
				mock.ExpectQuery(`SELECT COUNT\(\*\) FROM app_user`).
					WithArgs(expectedSetSearch, expectedSearch).
					WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(expectedCount))

				// Mock data query
				rows := sqlmock.NewRows([]string{"id", "full_name", "active", "role", "role_id"})
				for _, data := range expectedData {
					rows.AddRow(data.ID, data.FullName, data.Active, data.Role, data.RoleID)
				}

				mock.ExpectQuery(`SELECT`).
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
			order:         "created_at DESC",
			expectedCount: 12,
			expectedData: []entities.User{
				{
					ID:       "user-6",
					FullName: "Manager User",
					Active:   true,
					Role:     "Manager",
					RoleID:   1,
				},
				{
					ID:       "user-7",
					FullName: "Supervisor User",
					Active:   false,
					Role:     "Supervisor",
					RoleID:   2,
				},
			},
			wantErr:     false,
			expectedErr: nil,
			setupMock: func(mock sqlmock.Sqlmock, search string, page int, limit int, order string, expectedCount int, expectedData []entities.User) {
				expectedSearch := search
				if search != "" {
					expectedSearch = "%" + search + "%"
				}
				expectedSetSearch := search != ""
				expectedOffset := (page * limit) - limit

				// Mock count query
				mock.ExpectQuery(`SELECT COUNT\(\*\) FROM app_user`).
					WithArgs(expectedSetSearch, expectedSearch).
					WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(expectedCount))

				// Mock data query
				rows := sqlmock.NewRows([]string{"id", "full_name", "active", "role", "role_id"})
				for _, data := range expectedData {
					rows.AddRow(data.ID, data.FullName, data.Active, data.Role, data.RoleID)
				}

				mock.ExpectQuery(`SELECT`).
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
			order:         "created_at DESC",
			expectedCount: 0,
			expectedData:  []entities.User{},
			wantErr:       false,
			expectedErr:   nil,
			setupMock: func(mock sqlmock.Sqlmock, search string, page int, limit int, order string, expectedCount int, expectedData []entities.User) {
				expectedSearch := "%" + search + "%"
				expectedSetSearch := true
				expectedOffset := (page * limit) - limit

				// Mock count query
				mock.ExpectQuery(`SELECT COUNT\(\*\) FROM app_user`).
					WithArgs(expectedSetSearch, expectedSearch).
					WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(expectedCount))

				// Mock data query - return empty rows
				rows := sqlmock.NewRows([]string{"id", "full_name", "active", "role"})

				mock.ExpectQuery(`SELECT`).
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
			order:         "created_at DESC",
			expectedCount: 0,
			expectedData:  nil,
			wantErr:       true,
			expectedErr:   constant.ErrUnknownSource,
			setupMock: func(mock sqlmock.Sqlmock, search string, page int, limit int, order string, expectedCount int, expectedData []entities.User) {
				expectedSearch := search
				expectedSetSearch := search != ""

				// Mock count query to return error
				mock.ExpectQuery(`SELECT COUNT\(\*\) FROM app_user`).
					WithArgs(expectedSetSearch, expectedSearch).
					WillReturnError(sql.ErrConnDone)
			},
		},
		{
			name:          "database error on data query",
			search:        "",
			page:          constant.DefaultPage,
			limit:         constant.DefaultLimit,
			order:         "created_at DESC",
			expectedCount: 5,
			expectedData:  nil,
			wantErr:       true,
			expectedErr:   constant.ErrUnknownSource,
			setupMock: func(mock sqlmock.Sqlmock, search string, page int, limit int, order string, expectedCount int, expectedData []entities.User) {
				expectedSearch := search
				expectedSetSearch := search != ""
				expectedOffset := (page * limit) - limit

				// Mock count query success
				mock.ExpectQuery(`SELECT COUNT\(\*\) FROM app_user`).
					WithArgs(expectedSetSearch, expectedSearch).
					WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(expectedCount))

				// Mock data query to return error
				mock.ExpectQuery(`SELECT`).
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

			// Set default order to match what's used in the service
			sortField := "created_at"
			direction := "DESC"

			args := userPresenter.ReadUserRequest{
				PaginationPayload: utils.PaginationPayload{
					Search:    tc.search,
					Page:      tc.page,
					Limit:     tc.limit,
					Sort:      sortField,
					Direction: direction,
				},
			}

			tc.setupMock(mock, tc.search, tc.page, tc.limit, tc.order, tc.expectedCount, tc.expectedData)

			resp, err := svc.ReadUsersService(ctx, args)

			if tc.wantErr {
				assert.Error(t, err)
				if tc.expectedErr != nil {
					assert.ErrorIs(t, err, tc.expectedErr)
				}
				assert.Equal(t, entities.UserPaginationResponse{}, resp)
			} else {
				assert.NoError(t, err)

				// Calculate expected pagination
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

				// Compare each user individually
				if tc.expectedData != nil {
					assert.Equal(t, len(tc.expectedData), len(resp.Data))
					for i, expectedUser := range tc.expectedData {
						assert.Equal(t, expectedUser.ID, resp.Data[i].ID)
						assert.Equal(t, expectedUser.FullName, resp.Data[i].FullName)
						assert.Equal(t, expectedUser.Active, resp.Data[i].Active)
						assert.Equal(t, expectedUser.Role, resp.Data[i].Role)
					}
				} else {
					assert.Nil(t, resp.Data)
				}
			}

			// Check if all expectations were met
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}
