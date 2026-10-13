package usecase

import (
	"context"

	"github.com/ABA-Developer/go-backend-template/constant"
	"github.com/ABA-Developer/go-backend-template/internal/application/auth/dto"
	"github.com/ABA-Developer/go-backend-template/internal/application/jwt"
	"github.com/ABA-Developer/go-backend-template/internal/domain/model"
)

type AuthUseCase interface {
	LoginUseCase(ctx context.Context, request dto.LoginPayload, userAgent, iPAddress string) (data model.Session, user model.User, err error)
	LogoutUseCase(ctx context.Context, claims *jwt.AccessTokenPayload, iPAddress string) (err error)
	AuthMeUseCase(ctx context.Context, id string) (data model.User, err error)
	CheckPermissionUseCase(ctx context.Context, menuURL constant.MenuKey, userID string, permissionCode []string) (hasAccess bool, err error)
}
