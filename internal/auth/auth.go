package auth

import (
	"github.com/gofiber/fiber/v2"

	"be-dashboard-nba/constant"
	"be-dashboard-nba/internal/db"
	"be-dashboard-nba/internal/jwt"
)

type Auth struct {
	db     db.DB
	claims *jwt.AccessTokenPayload
}

func NewAuth(db db.DB) *Auth {
	return &Auth{
		db: db,
	}
}

func GetAuth(c *fiber.Ctx) (*Auth, error) {
	a, ok := c.Locals("auth").(Auth)
	if !ok {
		return nil, constant.ErrTokenUnauthorized
	}

	return &a, nil
}

func (a *Auth) GetClaims() *jwt.AccessTokenPayload {
	return a.claims
}

func (a *Auth) SetClaims(claims *jwt.AccessTokenPayload) {
	a.claims = claims
}
