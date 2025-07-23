package routes

import (
	"github.com/gofiber/fiber/v2"

	authRoute "be-dashboard-nba/api/handlers/auth/app"
	"be-dashboard-nba/internal/db"
	"be-dashboard-nba/internal/validator"
)

func AuthRouter(app fiber.Router, db db.DB, validate *validator.Validator) {
	group := app.Group("/auth")

	authRoute.AuthRouter(group, db, validate)
}
