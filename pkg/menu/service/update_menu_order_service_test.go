package service

import (
	menuPresenter "be-dashboard-nba/api/presenter/menu"
	"be-dashboard-nba/constant"
	"context"
	"database/sql"
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestUpdateMenuOrderService(t *testing.T) {
	tests := []struct {
		name        string
		request     menuPresenter.UpdateMenuOrderRequest
		userID      string
		wantErr     bool
		expectedErr error
		setupMock   func(mock sqlmock.Sqlmock)
	}{
		{
			name: "success - single menu update",
			request: menuPresenter.UpdateMenuOrderRequest{
				SortedIDs: []int{1},
			},
			userID:      "user123",
			wantErr:     false,
			expectedErr: nil,
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec(`UPDATE app_menu`).
					WithArgs(0, "user123", 1).
					WillReturnResult(sqlmock.NewResult(0, 1))
				mock.ExpectCommit()
			},
		},
		{
			name: "success - multiple menus update",
			request: menuPresenter.UpdateMenuOrderRequest{
				SortedIDs: []int{1, 2, 3},
			},
			userID:      "user123",
			wantErr:     false,
			expectedErr: nil,
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec(`UPDATE app_menu`).
					WithArgs(0, "user123", 1).
					WillReturnResult(sqlmock.NewResult(0, 1))
				mock.ExpectExec(`UPDATE app_menu`).
					WithArgs(1, "user123", 2).
					WillReturnResult(sqlmock.NewResult(0, 1))
				mock.ExpectExec(`UPDATE app_menu`).
					WithArgs(2, "user123", 3).
					WillReturnResult(sqlmock.NewResult(0, 1))
				mock.ExpectCommit()
			},
		},
		{
			name: "error - begin transaction fails",
			request: menuPresenter.UpdateMenuOrderRequest{
				SortedIDs: []int{1, 2, 3},
			},
			userID:      "user123",
			wantErr:     true,
			expectedErr: constant.ErrUnknownSource,
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin().WillReturnError(errors.New("begin transaction error"))
			},
		},
		{
			name: "error - first update fails",
			request: menuPresenter.UpdateMenuOrderRequest{
				SortedIDs: []int{1, 2, 3},
			},
			userID:      "user123",
			wantErr:     true,
			expectedErr: constant.ErrUnknownSource,
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec(`UPDATE app_menu`).
					WithArgs(0, "user123", 1).
					WillReturnError(errors.New("update failed"))
				mock.ExpectRollback()
			},
		},
		{
			name: "error - middle update fails",
			request: menuPresenter.UpdateMenuOrderRequest{
				SortedIDs: []int{1, 2, 3, 4, 5},
			},
			userID:      "user123",
			wantErr:     true,
			expectedErr: constant.ErrUnknownSource,
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec(`UPDATE app_menu`).
					WithArgs(0, "user123", 1).
					WillReturnResult(sqlmock.NewResult(0, 1))
				mock.ExpectExec(`UPDATE app_menu`).
					WithArgs(1, "user123", 2).
					WillReturnResult(sqlmock.NewResult(0, 1))

				mock.ExpectExec(`UPDATE app_menu`).
					WithArgs(2, "user123", 3).
					WillReturnError(errors.New("update failed"))
				mock.ExpectRollback()
			},
		},
		{
			name: "error - commit fails",
			request: menuPresenter.UpdateMenuOrderRequest{
				SortedIDs: []int{1, 2, 3},
			},
			userID:      "user123",
			wantErr:     true,
			expectedErr: constant.ErrUnknownSource,
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec(`UPDATE app_menu`).
					WithArgs(0, "user123", 1).
					WillReturnResult(sqlmock.NewResult(0, 1))
				mock.ExpectExec(`UPDATE app_menu`).
					WithArgs(1, "user123", 2).
					WillReturnResult(sqlmock.NewResult(0, 1))
				mock.ExpectExec(`UPDATE app_menu`).
					WithArgs(2, "user123", 3).
					WillReturnResult(sqlmock.NewResult(0, 1))
				mock.ExpectCommit().WillReturnError(errors.New("commit error"))
			},
		},
		{
			name: "error - rollback fails during error handling",
			request: menuPresenter.UpdateMenuOrderRequest{
				SortedIDs: []int{1, 2, 3},
			},
			userID:      "user123",
			wantErr:     true,
			expectedErr: constant.ErrUnknownSource,
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec(`UPDATE app_menu`).
					WithArgs(0, "user123", 1).
					WillReturnError(errors.New("update failed"))
				mock.ExpectRollback().WillReturnError(errors.New("rollback also failed"))
			},
		},
		{
			name: "success - empty parent ID",
			request: menuPresenter.UpdateMenuOrderRequest{
				ParentID:  nil,
				SortedIDs: []int{1, 2},
			},
			userID:      "user123",
			wantErr:     false,
			expectedErr: nil,
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec(`UPDATE app_menu`).
					WithArgs(0, "user123", 1).
					WillReturnResult(sqlmock.NewResult(0, 1))
				mock.ExpectExec(`UPDATE app_menu`).
					WithArgs(1, "user123", 2).
					WillReturnResult(sqlmock.NewResult(0, 1))
				mock.ExpectCommit()
			},
		},
		{
			name: "error - update returns no rows affected",
			request: menuPresenter.UpdateMenuOrderRequest{
				SortedIDs: []int{999},
			},
			userID:      "user123",
			wantErr:     true,
			expectedErr: constant.ErrUnknownSource,
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec(`UPDATE app_menu`).
					WithArgs(0, "user123", 999).
					WillReturnResult(sqlmock.NewResult(0, 0)) // 0 rows affected
			},
		},
		{
			name: "error - database connection lost during update",
			request: menuPresenter.UpdateMenuOrderRequest{
				SortedIDs: []int{1, 2},
			},
			userID:      "user123",
			wantErr:     true,
			expectedErr: constant.ErrUnknownSource,
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec(`UPDATE app_menu`).
					WithArgs(0, "user123", 1).
					WillReturnError(sql.ErrConnDone)
				mock.ExpectRollback()
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			svc, mock, cleanup := newTestService(t)
			defer cleanup()

			ctx := context.Background()

			tc.setupMock(mock)

			err := svc.UpdateMenuOrderService(ctx, tc.request, tc.userID)

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
