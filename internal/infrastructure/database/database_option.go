package db

import (
	"context"
	"os"
	"strconv"

	"github.com/pkg/errors"

	"be-dashboard-nba/internal/infrastructure/runtime/env"
)

type databaseOption struct {
	driver   string
	host     string
	port     int
	username string
	password string
	schema   string
	sslmode  string
	ctx      context.Context
	*connectionOption
}

func newDatabaseOption(ctx context.Context) (*databaseOption, error) {
	driver := os.Getenv("DB_DRIVER")
	host := os.Getenv("DB_HOST")
	portStr := os.Getenv("DB_PORT")
	username := os.Getenv("DB_USERNAME")
	password := os.Getenv("DB_PASSWORD")
	schema := os.Getenv("DB_SCHEMA")
	sslmode := os.Getenv("DB_SSLMODE")

	if portStr == "" {
		portStr = "0"
	}

	port, err := strconv.Atoi(portStr)
	if err != nil {
		return nil, errors.Wrapf(err, "error parse int on port db env : %s", portStr)
	}

	if host == "" {
		return nil, errors.Wrapf(errors.New("invalid data source host or port"), "db: host=%s port=%d", host, port)
	}

	conn := defaultConnectionOption()
	conn.maxIdle = env.GetInt("DB_MAX_IDLE_CONN", conn.maxIdle)
	conn.maxOpen = env.GetInt("DB_MAX_OPEN_CONN", conn.maxOpen)
	conn.maxLifetime = env.GetDuration("DB_MAX_LIFETIME_CONN", conn.maxLifetime)
	conn.maxConnIdleTime = env.GetDuration("DB_MAX_CONN_IDLE_TIME", conn.maxConnIdleTime)
	conn.keepAliveInterval = env.GetDuration("DB_KEEP_ALIVE_INTERVAL_CONN", conn.keepAliveInterval)

	return &databaseOption{
		driver:           driver,
		host:             host,
		port:             port,
		username:         username,
		password:         password,
		schema:           schema,
		sslmode:          sslmode,
		ctx:              ctx,
		connectionOption: conn,
	}, nil
}
