package app

import (
	"database/sql"
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog"

	"be-dashboard-nba/docs"
	"be-dashboard-nba/internal/config"
	"be-dashboard-nba/internal/db"
	"be-dashboard-nba/internal/validator"
)

type Application struct {
	Server    *fiber.App
	DB        *sql.DB
	Log       *zerolog.Logger
	Config    *config.Config
	Validator *validator.Validator
}

func NewApplication(log *zerolog.Logger) (*Application, error) {
	config := config.NewConfig(log)

	setSwaggerConfiguration(config)

	db, err := db.NewDatabase(config, log)
	if err != nil {
		log.Fatal().Msg("failed to connect to database")
		return nil, fmt.Errorf("failed to connect to database")
	}

	validator := validator.NewValidator()

	server := startHTTPServer(config)

	application := &Application{
		Server:    server,
		DB:        db,
		Log:       log,
		Config:    config,
		Validator: validator,
	}

	return application, nil
}

func setSwaggerConfiguration(config *config.Config) {
	docs.SwaggerInfo.BasePath = config.Swagger.BasePath
	docs.SwaggerInfo.Host = config.Swagger.Host
}
