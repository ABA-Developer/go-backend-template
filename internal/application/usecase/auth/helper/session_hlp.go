package helper

import (
	"context"

	"github.com/gofiber/fiber/v2/log"

	"be-dashboard-nba/internal/application/jwt"
	"be-dashboard-nba/internal/domain/model"
	authPresenter "be-dashboard-nba/internal/presentation/presenter/auth"
)

func GenerateSessionModel(
	ctx context.Context,
	request authPresenter.SessionPayload,
) (data model.Session, err error) {
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
