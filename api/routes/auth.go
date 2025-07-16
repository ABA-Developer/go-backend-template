package routes

import (
	"be-dashboard-nba/api/handlers"
	"be-dashboard-nba/pkg/auth"

	"github.com/gofiber/fiber/v2"
)

func AuthRouter(app fiber.Router, service auth.Service) {
	routes := app.Group("/auth")

	routes.Post("/login", handlers.LoginUser(service))
	// app.Post("/logout", handlers.LogoutUser(service))
}
