package middleware

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"

	"be-dashboard-nba/constant"
	"be-dashboard-nba/internal/auth"
	"be-dashboard-nba/internal/db"
	"be-dashboard-nba/internal/jwt"
	"be-dashboard-nba/internal/permissions"
	"be-dashboard-nba/pkg/auth/service"
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
			return fiber.NewError(http.StatusUnauthorized, "Token tidak valid")
		}

		accessTokenClaims, err := jwt.ClaimsAccessToken(token)
		if err != nil {
			log.WithContext(c.Context()).Error(err, "error claims access token")
			return fiber.NewError(http.StatusUnauthorized, "Token tidak valid atau kedaluwarsa")
		}

		et.auth.SetClaims(&accessTokenClaims)

		err = et.auth.ValidateSession(c.UserContext())
		if err != nil {

			if errors.Is(err, sql.ErrNoRows) {
				log.WithContext(c.Context()).Warn(err, "Invalid session (token valid, but session not found)")
				return fiber.NewError(http.StatusUnauthorized, "Sesi tidak valid atau telah berakhir")
			}

			log.WithContext(c.Context()).Error(err, "Failed to validate session (database error)")
			return fiber.NewError(http.StatusInternalServerError, "Gagal memvalidasi sesi")
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

func Authorize(
	svc service.Service,
	menuURL constant.MenuKey,
	permissionCode constant.PermissionCode,
) fiber.Handler {
	return func(c *fiber.Ctx) (err error) {
		ah, err := auth.GetAuth(c)
		if err != nil {
			log.WithContext(c.UserContext()).Error(err, "error get auth handler")
			return fiber.NewError(http.StatusUnauthorized, "Unauthorized")
		}

		userID := ah.GetClaims().UserID
		codesToCheck := permissions.GetInheritedPermissions(permissionCode)
		fmt.Println(codesToCheck)

		hasAccess, err := svc.CheckPermissionService(c.UserContext(), menuURL, userID, codesToCheck)
		if err != nil {
			log.WithContext(c.UserContext()).Error(err, "Failed to check permissions")
			return fiber.NewError(http.StatusInternalServerError, "Failed to check permissions")

		}
		if !hasAccess {
			return fiber.NewError(http.StatusForbidden, "Access Forbidden")
		}
		return c.Next()
	}
}
