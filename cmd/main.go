package main

import (
	"fmt"
	"os"

	"github.com/gofiber/fiber/v2"

	"be-dashboard-nba/api/middleware"
	"be-dashboard-nba/api/routes"
	"be-dashboard-nba/internal/config"
	"be-dashboard-nba/internal/db"
	"be-dashboard-nba/internal/utils"
	"be-dashboard-nba/pkg/auth"
	"be-dashboard-nba/pkg/user"
)

func main() {
	// Register logger
	log := utils.NewLogger()
	log.Info().Msg("Starting application...")

	// Register config
	cfg := config.NewConfig()

	// Register DB
	db, err := db.NewDatabase(cfg)
	if err != nil {
		log.Panic().Msg(fmt.Sprintf("Database connection failed, %v", err))
		os.Exit(0)
		return
	}

	// Register fiber
	app := fiber.New(
		fiber.Config{
			DisableStartupMessage: true,
		},
	)

	// Register middlewares
	middleware.TimeoutMiddleware(app)
	middleware.CorsMiddleware(app)
	middleware.RecoverMiddleware(app)
	middleware.RateLimiterMiddleware(app)
	middleware.LoggerMiddleware(app)

	app.Get("/", func(ctx *fiber.Ctx) error {
		return ctx.JSON(fiber.Map{
			"message": cfg.Name + " is Running",
		})
	})

	// Register repository
	userRepo := user.NewRepo(db)
	authRepo := auth.NewRepo(db)

	// Register service
	userService := user.NewService(userRepo)
	authService := auth.NewService(authRepo)

	// API group
	api := app.Group("/api").Group("/v1")

	// Register router
	routes.UserRouter(api, userService)
	routes.AuthRouter(api, authService)

	defer db.Close()

	log.Printf("starting http server %v:%v", cfg.Host, cfg.Port)

	if err := app.Listen(fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)); err != nil {
		log.Fatal().Msg(fmt.Sprintf("starting http server: %v", err))
	}
}
