package routes

import (
	"github.com/gofiber/fiber/v2"

	userRoute "be-dashboard-nba/api/handlers/user/app"
	"be-dashboard-nba/internal/db"
	"be-dashboard-nba/internal/validator"
)

func UserRouter(app fiber.Router, db db.DB, validate *validator.Validator) {
	group := app.Group("/user")

	userRoute.UserRouter(group, db, validate)
}
