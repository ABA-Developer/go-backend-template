package api

import (
	"github.com/ABA-Developer/go-backend-template/internal/infrastructure/api/routes"
	"github.com/ABA-Developer/go-backend-template/internal/infrastructure/runtime/container"

	"github.com/gofiber/fiber/v2"
)

func routesAPI(app *fiber.App, c *container.Container) {
	api := app.Group("/api").Group("/v1")
	routes.AddAuthRoutes(api, c)
	routes.AddUserRoutes(api, c)
	routes.AddMenuRoutes(api, c)
	routes.AddRoleRoutes(api, c)
}
