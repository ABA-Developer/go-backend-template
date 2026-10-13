package helper

import (
	"context"

	"github.com/ABA-Developer/go-backend-template/internal/infrastructure/logger"

	"github.com/ABA-Developer/go-backend-template/internal/application/auth/dto"
	"github.com/ABA-Developer/go-backend-template/internal/application/auth/mapper"
	"github.com/ABA-Developer/go-backend-template/internal/application/jwt"
	"github.com/ABA-Developer/go-backend-template/internal/domain/model"
)

func GenerateSessionModel(
	ctx context.Context,
	request dto.SessionPayload,
) (data model.Session, err error) {
	accessToken, err := jwt.GenerateAccessToken(mapper.ToAccessTokenPayload(request))
	if err != nil {
		logger.WithContext(ctx).Errorf("failed to generate access token for user %s: %v", request.UserID, err)
		return
	}

	refreshToken, err := jwt.GenerateRefreshToken(mapper.ToRefreshTokenPayload(request))
	if err != nil {
		logger.WithContext(ctx).Errorf("failed to generate refresh token for user %s: %v", request.UserID, err)
		return
	}

	data = model.Session{
		ID:                    request.SessionID,
		UserID:                request.UserID,
		AccessToken:           accessToken.Token,
		AccessTokenExpiredAt:  accessToken.ExpiresAt,
		RefreshToken:          refreshToken.Token,
		RefreshTokenExpiredAt: refreshToken.ExpiresAt,
		IPAddress:             request.IPAddress,
		UserAgent:             request.UserAgent,
	}

	return
}
