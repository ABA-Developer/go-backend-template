package db

import (
	"be-dashboard-nba/internal/config"
	"context"
	"database/sql"
	"fmt"
	"time"

	_ "github.com/lib/pq"
	"github.com/rs/zerolog"
)

func OpenPostgresDB(addr string, maxOpenConns int, maxIdleConns int, maxIdleTime int) (db *sql.DB, err error) {
	db, err = sql.Open("postgres", addr)
	if err != nil {
		return
	}

	db.SetMaxOpenConns(maxOpenConns)
	db.SetMaxIdleConns(maxIdleConns)
	db.SetConnMaxIdleTime(time.Duration(maxIdleConns))

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = db.PingContext(ctx)
	if err != nil {
		return
	}

	return
}

func NewPostgresDB(config config.Config, log zerolog.Logger) (db *sql.DB, err error) {
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		config.DB.Host,
		config.DB.Port,
		config.DB.Username,
		config.DB.Password,
		config.DB.Name,
		config.DB.SSLMode,
	)

	db, err = OpenPostgresDB(dsn, config.DB.MaxOpenConn, config.DB.MaxIdleConn, config.DB.MaxIdleTime)
	if err != nil {
		return
	}

	log.Printf("Successfully connected to postgresql %s:%s schema: %s", config.DB.Host, config.DB.Port, config.DB.Name)

	return
}
