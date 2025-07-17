package middleware

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"

	"be-dashboard-nba/constant"
	"be-dashboard-nba/internal/auth"
	"be-dashboard-nba/internal/db"
	"be-dashboard-nba/internal/jwt"
)

type EnsureToken struct {
	auth *auth.Auth
}

func NewEnsureToken(db db.DB) *EnsureToken {
	ah := auth.NewAuth(db)
	return &EnsureToken{auth: ah}
}

func (et *EnsureToken) ValidateToken() fiber.Handler {
	return func(c *fiber.Ctx) error {
		tokenHeader := c.Get(constant.DefaultMdwHeaderToken)
		token, err := parseHeaderToken(tokenHeader)
		if err != nil {
			log.WithContext(c.Context()).Error(err, "error parse header token")
			return fiber.NewError(http.StatusUnauthorized, err.Error())
		}

		accessTokenClaims, err := jwt.ClaimsAccessToken(token)
		if err != nil {
			log.WithContext(c.Context()).Error(err, "error claims access token")
			return fiber.NewError(http.StatusUnauthorized, err.Error())
		}

		et.auth.SetClaims(&accessTokenClaims)

		err = et.auth.ValidateSession(c.UserContext())
		if err != nil {
			return fiber.NewError(http.StatusUnauthorized, err.Error())
		}

		c.Locals("auth", *et.auth)
		return c.Next()
	}
}

func parseHeaderToken(headerDataToken string) (string, error) {
	if !strings.Contains(headerDataToken, "Bearer") {
		return "", constant.ErrHeaderTokenNotFound
	}

	splitToken := strings.Split(headerDataToken, fmt.Sprintf("%s ", constant.DefaultMdwHeaderBearer))
	if len(splitToken) <= 1 {
		return "", constant.ErrHeaderTokenInvalid
	}

	return splitToken[1], nil
}
