package middleware

import (
	"github.com/ABA-Developer/go-backend-template/constant"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/limiter"
)

func RateLimiterMiddleware(app *fiber.App) {
	app.Use(limiter.New(limiter.Config{
		Max:        constant.DefaultMdwRateLimiter,
		Expiration: constant.DefaultMdwRateLimiterDuration,
	}))
}
