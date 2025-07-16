package routes

import (
	"be-dashboard-nba/api/handlers"
	"be-dashboard-nba/api/middleware"
	"be-dashboard-nba/pkg/user"

	"github.com/gofiber/fiber/v2"
)

func UserRouter(app fiber.Router, service user.Service) {
	routes := app.Group("/user", middleware.JWTProtected())

	routes.Get("/", handlers.GetAllUser(service))       // Get all user
	routes.Post("/", handlers.CreateUser(service))      // Create user
	routes.Get("/:id", handlers.GetUserById(service))   // Get user by Id
	routes.Put("/:id", handlers.UpdateUser(service))    // Update user
	routes.Delete("/:id", handlers.DeleteUser(service)) // Delete user
}
