package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"be-dashboard-nba/api/app"
	"be-dashboard-nba/api/routes"
	"be-dashboard-nba/internal/utils"

	_ "be-dashboard-nba/docs"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

// @title           Base APP Go
// @version         1.0
// @description     Template API for Golang project

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:8000
// @BasePath  /api/v1
// @schemes   http https

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer {your token}" into the field below.

func main() {
	if os.Getenv("APP_ENV") == "" {
		if err := godotenv.Load(".env"); err != nil {
			log.Fatalf("ERROR loading env file: %s", err.Error())
		}
	}

	setDefaultTimezone()

	log := utils.NewLogger()

	// Create a context that listens for OS signals
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
	defer stop()

	// Initialize application dependencies
	application, err := app.NewApplication(log)
	if err != nil {
		log.Fatal().Msgf("Failed to initialize application, %v", err)
	}

	routes.Routes(application)

	go func() {
		if err := application.Server.Listen(fmt.Sprintf(":%d", application.Config.Port)); err != nil {
			log.Error().Err(err).Msg("Server failed to start or stopped with an error")
			stop()
			os.Exit(1)
		}
	}()

	log.Info().Msgf("Server is listening on port %d", application.Config.Port)

	<-ctx.Done()
	log.Warn().Msg("Shutdown signal received, starting graceful shutdown...")

	app.GracefulShutdown(application)
}

func setDefaultTimezone() {
	loc, err := time.LoadLocation("Asia/Jakarta")
	if err != nil {
		loc = time.Now().Location()
	}

	time.Local = loc
}
