package usecase

import (
	"context"

	appdto "github.com/ABA-Developer/go-backend-template/internal/application/user/dto"
	"github.com/ABA-Developer/go-backend-template/internal/domain/model"
)

type UserUseCase interface {
	ReadDetailUserUseCase(ctx context.Context, id string) (data model.User, err error)
	UpdateUserUseCase(ctx context.Context, request appdto.UpdateUserPayload, updatedBy string, userID string) (err error)
	ReadUsersUseCase(ctx context.Context, req appdto.ReadUserPayload) (data model.UserPaginationResponse, err error)
	CreateUserUseCase(ctx context.Context, req appdto.CreateUserPayload, userID string) (err error)
	DeleteUserUseCase(ctx context.Context, userID string, deletedBy string) (err error)
	ReadUserProfileUseCase(ctx context.Context, userID string) (data model.User, err error)
}
