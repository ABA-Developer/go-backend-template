package app

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"

	"be-dashboard-nba/api/middleware"
	"be-dashboard-nba/internal/config"

	"github.com/gofiber/swagger"
)



func startHTTPServer(cfg *config.Config) *fiber.App {

	app := fiber.New(fiber.Config{
		DisableStartupMessage: true,
		ErrorHandler: errorHandler(),
	})
	app.Use(logger.New())
	// Register middlewares
	middleware.TimeoutMiddleware(app)
	middleware.CorsMiddleware(app)
	middleware.RecoverMiddleware(app)
	middleware.RateLimiterMiddleware(app)
	middleware.LoggerMiddleware(app)

	app.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"message": cfg.Name + " is Running",
		})
	})

	if cfg.Swagger.IsEnabled {
		app.Get("/docs/*", swagger.HandlerDefault)
	}


	return app
}

