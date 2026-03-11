package jwt

import (
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type RefreshTokenPayload struct {
	SessionID string
}

// Generate refresh token with signing JWT method HS256.
func GenerateRefreshToken(request RefreshTokenPayload) (response TokenPayload, err error) {
	expiredDuration, err := time.ParseDuration(os.Getenv("AUTH_REFRESH_TOKEN_EXPIRES"))
	if err != nil {
		return
	}

	expiresAt := time.Now().Add(expiredDuration)

	claims := &jwt.MapClaims{
		"jti": request.SessionID,
		"exp": expiresAt.Unix(),
	}

	token, err := GenerateJWT(claims, os.Getenv("AUTH_REFRESH_TOKEN_SECRET_KEY"))
	if err != nil {
		return
	}

	response = TokenPayload{
		Token:     token,
		ExpiresAt: expiresAt,
	}

	return
}

// Parse JWT and return refresh token payload.
func ClaimsRefreshToken(token string) (response RefreshTokenPayload, err error) {
	claims, err := ClaimsJWT(token, os.Getenv("AUTH_REFRESH_TOKEN_SECRET_KEY"))
	if err != nil {
		return
	}

	response = RefreshTokenPayload{
		SessionID: claims["jti"].(string),
	}

	return
}
