package main

import (
	"fmt"
	"os"

	"be-dashboard-nba/internal/config"
	"be-dashboard-nba/internal/db"
	"be-dashboard-nba/internal/server"
	"be-dashboard-nba/internal/utils"
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
	defer db.Close()

	server.StartHTTPServer(cfg, db)
}
