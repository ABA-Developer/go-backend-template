package usecase

import (
	"context"

	appdto "github.com/ABA-Developer/go-backend-template/internal/application/role/dto"
	"github.com/ABA-Developer/go-backend-template/internal/domain/model"
)

type RoleUseCase interface {
	ReadRolesUseCase(ctx context.Context, args appdto.ReadRolesPayload) (data model.RolePaginationResponse, err error)
	ReadRoleDetailUseCase(ctx context.Context, roleID int) (data model.Role, err error)
	CreateRoleUseCase(ctx context.Context, paylaod appdto.CreateRolePayload, userID string) (err error)
	UpdateRoleUseCase(ctx context.Context, payload appdto.UpdateRolePayload, userID string, roleID int) (err error)
	DeleteRoleUseCase(ctx context.Context, roleID int) (err error)
	ReadRoleAccessUseCase(ctx context.Context, args appdto.ReadRoleAccessesPayload, roleID int) (data model.RoleAccessPaginationResponse, err error)
	UpdateRoleAccessUseCase(ctx context.Context, roleID int, request appdto.UpdateRoleAccessPayload) (err error)
}
