package db

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"

	config "be-dashboard-nba/internal/infrastructure/runtime"
)

func NewPostgresql(config *config.Config) (db *sql.DB, err error) {
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		config.DB.Host,
		config.DB.Port,
		config.DB.Username,
		config.DB.Password,
		config.DB.Name,
		config.DB.SSLMode,
	)

	db, err = openSQL("postgres", dsn, config)
	if err != nil {
		return
	}

	log.Printf("Successfully connected to postgresql %s:%s schema: %s", config.DB.Host, config.DB.Port, config.DB.Name)

	go keepAlive(db, config.DB.Driver, config.DB.Name, config.DB.KeepAliveInterval)

	return
}
