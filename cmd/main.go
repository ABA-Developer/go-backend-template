package main

import (
	"be-dashboard-nba/api/routes"
	"be-dashboard-nba/internal/config"
	"be-dashboard-nba/internal/db"
	"be-dashboard-nba/internal/utils"
	"be-dashboard-nba/pkg/auth"
	"be-dashboard-nba/pkg/user"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func main() {
	// Register logger
	log := utils.NewLogger("log/app.log")
	log.Info().Msg("Starting application...")

	// Register config
	config := config.NewConfig(log)

	// Register DB
	db := db.NewPostgresDB(*config, log)
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
	app := fiber.New()
	app.Use(cors.New())
	app.Use(logger.New())
	app.Get("/", func(ctx *fiber.Ctx) error {
		return ctx.Send([]byte("Welcome to the go backend nba"))
	})

	// API group
	api := app.Group("/api").Group("/v1")

	// Register router
	routes.UserRouter(api, userService)
	routes.AuthRouter(api, authService)

	defer db.Close()
	log.Fatal().Msg(app.Listen(":8080").Error())
}
