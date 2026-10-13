package main

import (
	"log"
	"os"
	"time"

	"github.com/ABA-Developer/go-backend-template/internal/infrastructure/api"
	db "github.com/ABA-Developer/go-backend-template/internal/infrastructure/database"
	"github.com/ABA-Developer/go-backend-template/internal/infrastructure/logger"
	runtime "github.com/ABA-Developer/go-backend-template/internal/infrastructure/runtime"
	"github.com/ABA-Developer/go-backend-template/internal/infrastructure/runtime/container"
	"github.com/ABA-Developer/go-backend-template/internal/infrastructure/validator"

	_ "github.com/ABA-Developer/go-backend-template/docs"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

// @title           Base APP Go
// @version         1.0
// @description     Template API for Golang project
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

	logger.NewLogger()

	ctx, stop := runtime.NewRuntimeContext()
	defer stop()

	log := logger.WithContext(ctx)

	cfg := runtime.NewConfig()

	dbConn, err := db.NewDatabase(ctx)
	if err != nil {
		log.Fatal().Msgf("Failed to initialize database, %v", err)
	}

	validatorEngine := validator.NewValidator()
	c := container.NewContainer(dbConn, log.Raw(), validatorEngine)

	if err := api.RunFiberServer(ctx, c, cfg); err != nil {
		log.Fatal().Msgf("Failed to run server, %v", err)
	}
}

func setDefaultTimezone() {
	loc, err := time.LoadLocation("Asia/Jakarta")
	if err != nil {
		loc = time.Now().Location()
	}

	time.Local = loc
}
