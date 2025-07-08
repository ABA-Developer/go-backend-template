package routes

import (
	"be-dashboard-nba/api/handlers"
	"be-dashboard-nba/pkg/auth"

	"github.com/gofiber/fiber/v2"
)

func AuthRouter(app fiber.Router, service auth.Service) {
	// app.Post("/auth/register", handlers.RegisterUser(service))
	app.Post("/auth/login", handlers.LoginUser(service))
	// app.Post("/auth/logout", handlers.LogoutUser(service))
}
