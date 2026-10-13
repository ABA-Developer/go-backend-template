package usecase

import (
	"context"

	appLogger "github.com/ABA-Developer/go-backend-template/internal/infrastructure/logger"
)

func log(ctx context.Context) *appLogger.Logger {
	return appLogger.WithContext(ctx)
}
