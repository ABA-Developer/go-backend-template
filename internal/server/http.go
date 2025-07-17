package server

import (
	"database/sql"
	"fmt"

	"github.com/gofiber/fiber/v2"

	"be-dashboard-nba/api/middleware"
	"be-dashboard-nba/api/routes"
	"be-dashboard-nba/internal/config"
	"be-dashboard-nba/internal/utils"
	"be-dashboard-nba/internal/validator"
	"be-dashboard-nba/pkg/user"
)

func StartHTTPServer(cfg *config.Config, database *sql.DB) {
	log := utils.NewLogger()

	app := fiber.New(fiber.Config{
		DisableStartupMessage: true,
		ErrorHandler:          errorHandler(),
	})

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

	// Setup dependencies
	userRepo := user.NewRepo(database)
	userService := user.NewService(userRepo)
	validate := validator.NewValidator()

	api := app.Group("/api").Group("/v1")
	routes.UserRouter(api, userService)
	routes.AuthRouter(api, database, validate)

	log.Printf("starting http server %v:%v", cfg.Host, cfg.Port)

	if err := app.Listen(fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)); err != nil {
		log.Fatal().Msg(fmt.Sprintf("starting http server: %v", err))
	}
}
