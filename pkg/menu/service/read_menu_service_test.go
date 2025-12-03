package service

import (
	menuPresenter "be-dashboard-nba/api/presenter/menu"
	"be-dashboard-nba/pkg/entities"
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestReadListMenuService(t *testing.T) {
	// Fixed time for consistent testing
	fixedTime := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
	fixedUpdatedTime := time.Date(2025, 1, 2, 0, 0, 0, 0, time.UTC)

	tests := []struct {
		name      string
		request   menuPresenter.ReadMenuListRequest
		expected  []entities.Menu
		wantErr   bool
		setupMock func(mock sqlmock.Sqlmock, request menuPresenter.ReadMenuListRequest)
	}{
		{
			name: "success without search",
			request: menuPresenter.ReadMenuListRequest{
				SetSearch: false,
				Search:    "",
			},
			expected: []entities.Menu{
				{
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
					CreatedBy:   "user1",
					CreatedAt:   fixedTime,
					UpdatedBy:   sql.NullString{String: "user2", Valid: true},
					UpdatedAt:   sql.NullTime{Time: fixedUpdatedTime, Valid: true},
				},
			},
			wantErr: false,
			setupMock: func(mock sqlmock.Sqlmock, request menuPresenter.ReadMenuListRequest) {
				params := request.ToParams()
				mock.ExpectQuery(`WITH RECURSIVE menu_with_parents AS`).
					WithArgs(params.SetSearch, params.Search).
					WillReturnRows(sqlmock.NewRows([]string{
						"id", "parent_id", "name", "description", "url", "sort", "group", "icon",
						"active", "display", "created_by", "created_at", "updated_by", "updated_at",
					}).AddRow(
						1,
						nil,
						"Dashboard",
						nil,
						"/dashboard",
						0,
						"main",
						nil,
						true,
						true,
						"user1",
						fixedTime,
						"user2",
						sql.NullTime{Time: fixedUpdatedTime, Valid: true},
					))
			},
		},
		{
			name: "success with search",
			request: menuPresenter.ReadMenuListRequest{
				SetSearch: true,
				Search:    "Dash",
			},
			expected: []entities.Menu{
				{
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
					CreatedBy:   "user1",
					CreatedAt:   fixedTime,
					UpdatedBy:   sql.NullString{String: "user2", Valid: true},
					UpdatedAt:   sql.NullTime{Time: fixedUpdatedTime, Valid: true},
				},
			},
			wantErr: false,
			setupMock: func(mock sqlmock.Sqlmock, request menuPresenter.ReadMenuListRequest) {
				params := request.ToParams()
				mock.ExpectQuery(`WITH RECURSIVE menu_with_parents AS`).
					WithArgs(params.SetSearch, params.Search).
					WillReturnRows(sqlmock.NewRows([]string{
						"id", "parent_id", "name", "description", "url", "sort", "group", "icon",
						"active", "display", "created_by", "created_at", "updated_by", "updated_at",
					}).AddRow(
						1,
						nil,
						"Dashboard",
						nil,
						"/dashboard",
						0,
						"main",
						nil,
						true,
						true,
						"user1",
						fixedTime,
						"user2",
						sql.NullTime{Time: fixedUpdatedTime, Valid: true},
					))
			},
		},
		{
			name: "empty result",
			request: menuPresenter.ReadMenuListRequest{
				SetSearch: false,
				Search:    "",
			},
			expected: nil,
			wantErr:  false,
			setupMock: func(mock sqlmock.Sqlmock, request menuPresenter.ReadMenuListRequest) {
				params := request.ToParams()
				mock.ExpectQuery(`WITH RECURSIVE menu_with_parents AS`).
					WithArgs(params.SetSearch, params.Search).
					WillReturnRows(sqlmock.NewRows([]string{
						"id", "parent_id", "name", "description", "url", "sort", "group", "icon",
						"active", "display", "created_by", "created_at", "updated_by", "updated_at",
					}))
			},
		},
		{
			name: "query error",
			request: menuPresenter.ReadMenuListRequest{
				SetSearch: false,
				Search:    "",
			},
			expected: nil,
			wantErr:  true,
			setupMock: func(mock sqlmock.Sqlmock, request menuPresenter.ReadMenuListRequest) {
				params := request.ToParams()
				mock.ExpectQuery(`WITH RECURSIVE menu_with_parents AS`).
					WithArgs(params.SetSearch, params.Search).
					WillReturnError(sql.ErrConnDone)
			},
		},
		{
			name: "multiple menus with hierarchy",
			request: menuPresenter.ReadMenuListRequest{
				SetSearch: false,
				Search:    "",
			},
			expected: []entities.Menu{
				{
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
					CreatedBy:   "user1",
					CreatedAt:   fixedTime,
					UpdatedBy:   sql.NullString{},
					UpdatedAt:   sql.NullTime{},
				},
				{
					ID:          2,
					ParentID:    sql.NullInt32{Int32: 1, Valid: true},
					Name:        "Submenu",
					Description: sql.NullString{String: "A submenu", Valid: true},
					URL:         sql.NullString{String: "/submenu", Valid: true},
					Sort:        0,
					Group:       "main",
					Icon:        sql.NullString{String: "icon", Valid: true},
					Active:      true,
					Display:     true,
					CreatedBy:   "user1",
					CreatedAt:   fixedTime,
					UpdatedBy:   sql.NullString{String: "user2", Valid: true},
					UpdatedAt:   sql.NullTime{Time: fixedUpdatedTime, Valid: true},
				},
			},
			wantErr: false,
			setupMock: func(mock sqlmock.Sqlmock, request menuPresenter.ReadMenuListRequest) {
				params := request.ToParams()
				mock.ExpectQuery(`WITH RECURSIVE menu_with_parents AS`).
					WithArgs(params.SetSearch, params.Search).
					WillReturnRows(sqlmock.NewRows([]string{
						"id", "parent_id", "name", "description", "url", "sort", "group", "icon",
						"active", "display", "created_by", "created_at", "updated_by", "updated_at",
					}).
						AddRow(
							1,
							nil,
							"Dashboard",
							nil,
							"/dashboard",
							0,
							"main",
							nil,
							true,
							true,
							"user1",
							fixedTime,
							nil,
							nil,
						).
						AddRow(
							2,
							1,
							"Submenu",
							"A submenu",
							"/submenu",
							0,
							"main",
							"icon",
							true,
							true,
							"user1",
							fixedTime,
							"user2",
							sql.NullTime{Time: fixedUpdatedTime, Valid: true},
						))
			},
		},
		{
			name: "search with formatted term",
			request: menuPresenter.ReadMenuListRequest{
				SetSearch: true,
				Search:    "searchterm",
			},
			expected: []entities.Menu{
				{
					ID:          3,
					ParentID:    sql.NullInt32{},
					Name:        "Search Result",
					Description: sql.NullString{},
					URL:         sql.NullString{String: "/search", Valid: true},
					Sort:        0,
					Group:       "main",
					Icon:        sql.NullString{},
					Active:      true,
					Display:     true,
					CreatedBy:   "user1",
					CreatedAt:   fixedTime,
					UpdatedBy:   sql.NullString{},
					UpdatedAt:   sql.NullTime{},
				},
			},
			wantErr: false,
			setupMock: func(mock sqlmock.Sqlmock, request menuPresenter.ReadMenuListRequest) {
				params := request.ToParams()
				// Verify that search term is properly formatted with %%
				assert.Equal(t, true, params.SetSearch)
				assert.Equal(t, "%searchterm%", params.Search)

				mock.ExpectQuery(`WITH RECURSIVE menu_with_parents AS`).
					WithArgs(params.SetSearch, params.Search).
					WillReturnRows(sqlmock.NewRows([]string{
						"id", "parent_id", "name", "description", "url", "sort", "group", "icon",
						"active", "display", "created_by", "created_at", "updated_by", "updated_at",
					}).AddRow(
						3,
						nil,
						"Search Result",
						nil,
						"/search",
						0,
						"main",
						nil,
						true,
						true,
						"user1",
						fixedTime,
						nil,
						nil,
					))
			},
		},
		{
			name: "inactive menu",
			request: menuPresenter.ReadMenuListRequest{
				SetSearch: false,
				Search:    "",
			},
			expected: []entities.Menu{
				{
					ID:          4,
					ParentID:    sql.NullInt32{},
					Name:        "Inactive Menu",
					Description: sql.NullString{},
					URL:         sql.NullString{String: "/inactive", Valid: true},
					Sort:        0,
					Group:       "main",
					Icon:        sql.NullString{},
					Active:      false,
					Display:     false,
					CreatedBy:   "user1",
					CreatedAt:   fixedTime,
					UpdatedBy:   sql.NullString{},
					UpdatedAt:   sql.NullTime{},
				},
			},
			wantErr: false,
			setupMock: func(mock sqlmock.Sqlmock, request menuPresenter.ReadMenuListRequest) {
				params := request.ToParams()
				mock.ExpectQuery(`WITH RECURSIVE menu_with_parents AS`).
					WithArgs(params.SetSearch, params.Search).
					WillReturnRows(sqlmock.NewRows([]string{
						"id", "parent_id", "name", "description", "url", "sort", "group", "icon",
						"active", "display", "created_by", "created_at", "updated_by", "updated_at",
					}).AddRow(
						4,
						nil,
						"Inactive Menu",
						nil,
						"/inactive",
						0,
						"main",
						nil,
						false,
						false,
						"user1",
						fixedTime,
						nil,
						nil,
					))
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			svc, mock, cleanup := newTestService(t)
			defer cleanup()

			ctx := context.Background()

			// Setup mock expectations
			tc.setupMock(mock, tc.request)

			// Execute service
			data, err := svc.ReadListMenuService(ctx, tc.request)

			// Assertions
			if tc.wantErr {
				assert.Error(t, err)
				assert.Nil(t, data)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expected, data)
				assert.Len(t, data, len(tc.expected))
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

// Test untuk memverifikasi ToParams function
func TestReadMenuListRequest_ToParams(t *testing.T) {
	tests := []struct {
		name     string
		request  menuPresenter.ReadMenuListRequest
		expected struct {
			SetSearch bool
			Search    string
		}
	}{
		{
			name: "without search",
			request: menuPresenter.ReadMenuListRequest{
				SetSearch: false,
				Search:    "",
			},
			expected: struct {
				SetSearch bool
				Search    string
			}{
				SetSearch: false,
				Search:    "",
			},
		},
		{
			name: "with search",
			request: menuPresenter.ReadMenuListRequest{
				SetSearch: true,
				Search:    "test",
			},
			expected: struct {
				SetSearch bool
				Search    string
			}{
				SetSearch: true,
				Search:    "%test%",
			},
		},
		{
			name: "search auto formatting",
			request: menuPresenter.ReadMenuListRequest{
				SetSearch: false, // Should be set to true by ToParams
				Search:    "auto",
			},
			expected: struct {
				SetSearch bool
				Search    string
			}{
				SetSearch: true,
				Search:    "%auto%",
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			params := tc.request.ToParams()

			assert.Equal(t, tc.expected.SetSearch, params.SetSearch)
			assert.Equal(t, tc.expected.Search, params.Search)
		})
	}
}
