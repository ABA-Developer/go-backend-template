package usecase

import (
	"context"
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

func TestReadRolesUseCase(t *testing.T) {
	tests := []struct {
		name    string
		req     appdto.ReadRolesPayload
		setup   func(roleRepo *mockrepo.RoleRepositoryMock)
		wantErr error
	}{
		{
			name: "positive/returns-pagination",
			req:  appdto.ReadRolesPayload{PaginationPayload: utils.PaginationPayload{Limit: 2, Page: 1}},
			setup: func(roleRepo *mockrepo.RoleRepositoryMock) {
				roleRepo.On("ReadRolesCountQry", testifymock.Anything, testifymock.Anything).Return(3, nil).Once()
				roleRepo.On("ReadRolesQry", testifymock.Anything, testifymock.Anything).Return([]model.Role{{ID: 1}}, nil).Once()
			},
		},
		{
			name: "negative/count-error",
			req:  appdto.ReadRolesPayload{PaginationPayload: utils.PaginationPayload{Limit: 2, Page: 1}},
			setup: func(roleRepo *mockrepo.RoleRepositoryMock) {
				roleRepo.On("ReadRolesCountQry", testifymock.Anything, testifymock.Anything).Return(0, errors.New("db down")).Once()
			},
			wantErr: constant.ErrUnknownSource,
		},
		{
			name: "negative/list-error",
			req:  appdto.ReadRolesPayload{PaginationPayload: utils.PaginationPayload{Limit: 2, Page: 1}},
			setup: func(roleRepo *mockrepo.RoleRepositoryMock) {
				roleRepo.On("ReadRolesCountQry", testifymock.Anything, testifymock.Anything).Return(1, nil).Once()
				roleRepo.On("ReadRolesQry", testifymock.Anything, testifymock.Anything).Return(nil, errors.New("db down")).Once()
			},
			wantErr: constant.ErrUnknownSource,
		},
		{
			name: "edge/default-limit-page",
			req:  appdto.ReadRolesPayload{},
			setup: func(roleRepo *mockrepo.RoleRepositoryMock) {
				roleRepo.On("ReadRolesCountQry", testifymock.Anything, testifymock.Anything).Return(int64(0), nil).Once()
				roleRepo.On("ReadRolesQry", testifymock.Anything, testifymock.Anything).Return([]model.Role{}, nil).Once()
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			roleRepo := mockrepo.NewRoleRepository()
			uc, _, cleanup := newTestUseCaseForCreateRole(t, roleRepo)
			defer cleanup()

			tt.setup(roleRepo)

			_, err := uc.ReadRolesUseCase(context.Background(), tt.req)
			if tt.wantErr != nil {
				assert.Error(t, err)
				assert.ErrorIs(t, err, tt.wantErr)
				return
			}
			assert.NoError(t, err)
		})
	}
}
