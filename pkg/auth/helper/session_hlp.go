package helper

import (
	"context"

	"github.com/gofiber/fiber/v2/log"

	authPresenter "be-dashboard-nba/api/presenter/auth"
	"be-dashboard-nba/internal/jwt"
	"be-dashboard-nba/pkg/entities"
)

func GenerateSessionModel(
	ctx context.Context,
	request authPresenter.SessionPayload,
) (data entities.Session, err error) {
	accessToken, err := jwt.GenerateAccessToken(request.ToAccessTokenRequest())
	if err != nil {
		log.WithContext(ctx).Errorf("failed to generate access token for user %s: %v", request.UserID, err)
		return
	}

	refreshToken, err := jwt.GenerateRefreshToken(request.ToRefreshTokenRequest())
	if err != nil {
		log.WithContext(ctx).Errorf("failed to generate refresh token for user %s: %v", request.UserID, err)
		return
	}

	data = entities.Session{
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
