package usecase

import (
	"context"
	"database/sql"
	"testing"

	"github.com/ABA-Developer/go-backend-template/constant"
	appdto "github.com/ABA-Developer/go-backend-template/internal/application/role/dto"
	"github.com/ABA-Developer/go-backend-template/internal/application/utils"
	"github.com/ABA-Developer/go-backend-template/internal/domain/model"
	mockrepo "github.com/ABA-Developer/go-backend-template/internal/test/mock/repository"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	testifymock "github.com/stretchr/testify/mock"
)

func TestReadRoleAccessUseCase(t *testing.T) {
	tests := []struct {
		name    string
		roleID  int
		req     appdto.ReadRoleAccessesPayload
		setup   func(roleRepo *mockrepo.RoleRepositoryMock)
		wantErr error
	}{
		{
			name:   "positive/returns-pagination",
			roleID: 10,
			req:    appdto.ReadRoleAccessesPayload{PaginationPayload: utils.PaginationPayload{Limit: 2, Page: 1}},
			setup: func(roleRepo *mockrepo.RoleRepositoryMock) {
				roleRepo.On("ReadRoleAccessCountQry", testifymock.Anything, testifymock.Anything).Return(1, nil).Once()
				roleRepo.On("ReadRoleByIDQry", testifymock.Anything, 10).Return(model.Role{ID: 10}, nil).Once()
				roleRepo.On("ReadRoleAccessQry", testifymock.Anything, testifymock.Anything).Return([]model.RoleAccessResponse{{RoleID: 10}}, nil).Once()
			},
		},
		{
			name:   "negative/count-error",
			roleID: 10,
			req:    appdto.ReadRoleAccessesPayload{PaginationPayload: utils.PaginationPayload{Limit: 2, Page: 1}},
			setup: func(roleRepo *mockrepo.RoleRepositoryMock) {
				roleRepo.On("ReadRoleAccessCountQry", testifymock.Anything, testifymock.Anything).Return(0, errors.New("db down")).Once()
			},
			wantErr: constant.ErrUnknownSource,
		},
		{
			name:   "negative/role-not-found",
			roleID: 10,
			req:    appdto.ReadRoleAccessesPayload{PaginationPayload: utils.PaginationPayload{Limit: 2, Page: 1}},
			setup: func(roleRepo *mockrepo.RoleRepositoryMock) {
				roleRepo.On("ReadRoleAccessCountQry", testifymock.Anything, testifymock.Anything).Return(0, nil).Once()
				roleRepo.On("ReadRoleByIDQry", testifymock.Anything, 10).Return(model.Role{}, sql.ErrNoRows).Once()
			},
			wantErr: constant.ErrRoleIdNotFound,
		},
		{
			name:   "negative/list-error",
			roleID: 10,
			req:    appdto.ReadRoleAccessesPayload{PaginationPayload: utils.PaginationPayload{Limit: 2, Page: 1}},
			setup: func(roleRepo *mockrepo.RoleRepositoryMock) {
				roleRepo.On("ReadRoleAccessCountQry", testifymock.Anything, testifymock.Anything).Return(1, nil).Once()
				roleRepo.On("ReadRoleByIDQry", testifymock.Anything, 10).Return(model.Role{ID: 10}, nil).Once()
				roleRepo.On("ReadRoleAccessQry", testifymock.Anything, testifymock.Anything).Return(nil, errors.New("db down")).Once()
			},
			wantErr: constant.ErrUnknownSource,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			roleRepo := mockrepo.NewRoleRepository()
			uc, _, cleanup := newTestUseCaseForCreateRole(t, roleRepo)
			defer cleanup()

			tt.setup(roleRepo)

			_, err := uc.ReadRoleAccessUseCase(context.Background(), tt.req, tt.roleID)
			if tt.wantErr != nil {
				assert.Error(t, err)
				assert.ErrorIs(t, err, tt.wantErr)
				return
			}
			assert.NoError(t, err)
		})
	}
}
