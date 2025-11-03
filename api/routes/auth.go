package routes

import (
	"be-dashboard-nba/api/app"
	"be-dashboard-nba/api/handlers"
	"be-dashboard-nba/api/middleware"
	"be-dashboard-nba/pkg/auth/service"

	"github.com/gofiber/fiber/v2"
)

func AuthRouter(http fiber.Router, application *app.Application) {
	svc := service.NewService(application.DB, application.Log)
	mdw := middleware.NewEnsureToken(application.DB)

	routes := http.Group("/auth")

	routes.Post("/login", handlers.Login(svc, application.Validator))
	routes.Post("/logout", mdw.ValidateToken(), handlers.Logout(svc))
	routes.Get("/me", mdw.ValidateToken(), handlers.AuthMe(svc))
}
