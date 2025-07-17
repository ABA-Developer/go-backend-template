package app

import (
	"github.com/gofiber/fiber/v2"

	"be-dashboard-nba/api/middleware"
	"be-dashboard-nba/internal/db"
	"be-dashboard-nba/internal/validator"
	"be-dashboard-nba/pkg/auth/service"
)

func AuthRouter(app fiber.Router, db db.DB, validate *validator.Validator) {
	svc := service.NewService(db)
	mdw := middleware.NewEnsureToken(db)

	app.Post("/login", loginApp(svc, validate))
	app.Post("/logout", mdw.ValidateToken(), logoutApp(svc))
}
