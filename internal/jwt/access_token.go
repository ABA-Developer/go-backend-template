package jwt

import (
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type AccessTokenPayload struct {
	GUID   string
	UserID int64
	Role   string
}

// Generate access token with signing JWT method HS256.
func GenerateAccessToken(request AccessTokenPayload) (response TokenPayload, err error) {
	expiredDuration, err := time.ParseDuration(os.Getenv("AUTH_ACCESS_TOKEN_EXPIRES"))
	if err != nil {
		return
	}

	expiresAt := time.Now().Add(expiredDuration)

	claims := &jwt.MapClaims{
		"jti":  request.GUID,
		"exp":  expiresAt.Unix(),
		"uri":  request.UserID,
		"role": request.Role,
	}

	token, err := GenerateJWT(claims, os.Getenv("AUTH_ACCESS_TOKEN_SECRET_KEY"))
	if err != nil {
		return
	}

	response = TokenPayload{
		Token:     token,
		ExpiresAt: expiresAt,
	}

	return
}

// Parse JWT and return access token payload.
func ClaimsAccessToken(token string) (response AccessTokenPayload, err error) {
	claims, err := ClaimsJWT(token, os.Getenv("AUTH_ACCESS_TOKEN_SECRET_KEY"))
	if err != nil {
		return
	}

	response = AccessTokenPayload{
		GUID:   claims["jti"].(string),
		UserID: claims["uri"].(int64),
		Role:   claims["role"].(string),
	}

	return
}
