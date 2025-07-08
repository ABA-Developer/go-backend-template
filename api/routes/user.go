package routes

import (
	"be-dashboard-nba/api/handlers"
	"be-dashboard-nba/api/middleware"
	"be-dashboard-nba/pkg/user"

	"github.com/gofiber/fiber/v2"
)

func UserRouter(app fiber.Router, service user.Service) {
	protectedUser := app.Group("/user", middleware.JWTProtected())
	protectedUser.Get("/", handlers.GetAllUser(service))       // Get all user
	protectedUser.Post("/", handlers.CreateUser(service))      // Create user
	protectedUser.Get("/:id", handlers.GetUserById(service))   // Get user by Id
	protectedUser.Put("/:id", handlers.UpdateUser(service))    // Update user
	protectedUser.Delete("/:id", handlers.DeleteUser(service)) // Delete user
}
