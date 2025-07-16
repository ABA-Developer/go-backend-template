package main

import (
	"be-dashboard-nba/api/routes"
	"be-dashboard-nba/internal/config"
	"be-dashboard-nba/internal/db"
	"be-dashboard-nba/internal/utils"
	"be-dashboard-nba/pkg/auth"
	"be-dashboard-nba/pkg/user"
	"fmt"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func main() {
	// Register logger
	log := utils.NewLogger()
	log.Info().Msg("Starting application...")

	// Register config
	cfg := config.NewConfig()

	// Register DB
	db := db.NewPostgresDB(*cfg, log)
	if db == nil {
		log.Panic().Msg("Database connection failed")
		os.Exit(0)
	}

	// Register repository
	userRepo := user.NewRepo(db)
	authRepo := auth.NewRepo(db)

	// Register service
	userService := user.NewService(userRepo)
	authService := auth.NewService(authRepo)

	// Register fiber
	app := fiber.New(
		fiber.Config{
			DisableStartupMessage: true,
		},
	)
	app.Use(cors.New())
	app.Use(logger.New())
	app.Get("/", func(ctx *fiber.Ctx) error {
		return ctx.JSON(fiber.Map{
			"message": cfg.Name + " is Running",
		})
	})

	// API group
	api := app.Group("/api").Group("/v1")

	// Register router
	routes.UserRouter(api, userService)
	routes.AuthRouter(api, authService)

	defer db.Close()

	log.Info().Msg(fmt.Sprintf("starting http server %v:%v", cfg.Host, cfg.Port))

	if err := app.Listen(fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)); err != nil {
		log.Fatal().Msg(fmt.Sprintf("starting http server: %v", err))
	}
}
