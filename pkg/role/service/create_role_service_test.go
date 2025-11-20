package service

import (
	rolePresenter "be-dashboard-nba/api/presenter/role"
	"be-dashboard-nba/constant"
	"context"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

func TestCreateRoleService(t *testing.T) {
	const createRoleQuery = `INSERT INTO app_role`
	getPtrString := func(s string) *string { return &s }
	tests := []struct {
		name        string
		request     rolePresenter.CreateRoleRequest
		userID      string
		wantErr     bool
		expectedErr error
		setupMock   func(mock sqlmock.Sqlmock, request rolePresenter.CreateRoleRequest, userID string)
	}{
		{
			name: "success create role",
			request: rolePresenter.CreateRoleRequest{
				Name:        "new role",
				Code:        "90900 099990 1213",
				Description: getPtrString("description role"),
			},
			userID:      "user-id",
			wantErr:     false,
			expectedErr: nil,
			setupMock: func(mock sqlmock.Sqlmock, request rolePresenter.CreateRoleRequest, userID string) {
				mock.ExpectExec(createRoleQuery).WithArgs(request.Code, request.Name, request.Description, userID).WillReturnResult(sqlmock.NewResult(1, 1))
			},
		},
		{
			name: "error sql",
			request: rolePresenter.CreateRoleRequest{
				Name:        "new role",
				Code:        "90900 099990 1213",
				Description: getPtrString("description role"),
			},
			userID:      "user-id",
			wantErr:     true,
			expectedErr: constant.ErrUnknownSource,
			setupMock: func(mock sqlmock.Sqlmock, request rolePresenter.CreateRoleRequest, userID string) {
				mock.ExpectExec(createRoleQuery).
					WithArgs(request.Code, request.Name, request.Description, userID).
					WillReturnError(errors.New("sql error"))
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			svc, mock, cleanup := newTestService(t)
			defer cleanup()

			ctx := context.Background()

			tc.setupMock(mock, tc.request, tc.userID)

			err := svc.CreateRoleService(ctx, tc.request, tc.userID)

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
