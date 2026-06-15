package routes

import (
	"be-dashboard-nba/internal/infrastructure/runtime/container"

	"github.com/gofiber/fiber/v2"
)

func Routes(app *fiber.App, c *container.Container) {
	api := app.Group("/api").Group("/v1")
	AuthRouter(api, c)
	UserRouter(api, c)
	MenuRouter(api, c)
	RoleRouter(api, c)
}
