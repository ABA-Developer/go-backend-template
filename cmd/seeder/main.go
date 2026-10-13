package main

import (
	"database/sql"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/kristijorgji/goseeder"

	_ "github.com/ABA-Developer/go-backend-template/database/seeders"
	db "github.com/ABA-Developer/go-backend-template/internal/infrastructure/database"
	"github.com/ABA-Developer/go-backend-template/internal/infrastructure/logger"
	config "github.com/ABA-Developer/go-backend-template/internal/infrastructure/runtime"
)

func main() {
	if os.Getenv("APP_ENV") == "" {
		if err := godotenv.Load(".env"); err != nil {
			log.Fatalf("ERROR loading env file: %s", err.Error())
		}
	}

	goseeder.WithSeeder(conProvider, func() {
		log.Printf("running seeder")
	})

	log.Println("done run seeder")
}

func conProvider() *sql.DB {
	logger.NewLogger()
	ctx, cancel := config.NewRuntimeContext()
	defer cancel()

	log := logger.WithContext(ctx)
	db, err := db.NewDatabase(ctx)
	if err != nil {
		log.Fatal().Msg("failed to connect to database")
		return nil
	}

	return db
}
