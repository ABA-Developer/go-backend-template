package middleware

import (
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func LoggerMiddleware(app *fiber.App) {
	app.Use(logger.New(logger.Config{
		Format: "[${time}] | ${ip} | ${status} | ${latency} | ${method} ${path} ${queryParams}\n",
		Next: func(c *fiber.Ctx) bool {
			return os.Getenv("APP_ENV") == "production"
		},
	}))
}
