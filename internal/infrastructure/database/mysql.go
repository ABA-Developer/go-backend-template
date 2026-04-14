package db

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"

	config "be-dashboard-nba/internal/infrastructure/runtime"
)

func NewMySQL(config *config.Config) (db *sql.DB, err error) {
	dsn := fmt.Sprintf("%s:%s@(%s:%s)/%s",
		config.DB.Username,
		config.DB.Password,
		config.DB.Host,
		config.DB.Port,
		config.DB.Name,
	)

	db, err = openSQL("mysql", dsn, config)
	if err != nil {
		return
	}

	log.Printf("Successfully connected to MySQL %s:%s database: %s", config.DB.Host, config.DB.Port, config.DB.Name)

	go keepAlive(db, config.DB.Driver, config.DB.Name, config.DB.KeepAliveInterval)

	return
}
