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

var count int64 = 0

func OpenPostgresDB(addr string, maxOpenConns int, maxIdleConns int, maxIdleTime int) (*sql.DB, error) {
	db, err := sql.Open("postgres", addr)
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(maxOpenConns)
	db.SetMaxIdleConns(maxIdleConns)
	db.SetConnMaxIdleTime(time.Duration(maxIdleConns))

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = db.PingContext(ctx)
	if err != nil {
		return nil, err
	}

	return db, nil
}

func NewPostgresDB(config config.Config, log zerolog.Logger) *sql.DB {
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		config.DB.Host,
		config.DB.Port,
		config.DB.Username,
		config.DB.Password,
		config.DB.Name,
		config.DB.SSLMode,
	)

	for {
		db, err := OpenPostgresDB(dsn, config.DB.MaxOpenConn, config.DB.MaxIdleConn, config.DB.MaxIdleTime)
		if err != nil {
			log.Info().Msg("PostgreSQL is not ready yet")
			count++
		} else {
			log.Info().Msg("Connected to PostgreSQL database")
			log.Info().Msg("Database connection pool established!")
			return db
		}

		if count > 10 {
			log.Info().Msg(err.Error())
			return nil
		}

		log.Info().Msg("Waiting for two seconds for reconnecting...")
		time.Sleep(2 * time.Second)

		continue
	}
}
