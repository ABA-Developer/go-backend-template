package main

import (
	"database/sql"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/kristijorgji/goseeder"

	_ "be-dashboard-nba/db/seeder"
	"be-dashboard-nba/internal/config"
	"be-dashboard-nba/internal/db"
	"be-dashboard-nba/internal/utils"
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
	log := utils.NewLogger()
	config := config.NewConfig(log)

	db, err := db.NewDatabase(config, log)
	if err != nil {
		log.Fatal().Msg("failed to connect to database")
		return nil
	}

	return db
}
