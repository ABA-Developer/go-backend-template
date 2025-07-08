package config

import (
	"be-dashboard-nba/internal/env"

	"github.com/joho/godotenv"
	"github.com/rs/zerolog"
)

type Config struct {
	Addr string
	Db   DbConfig
}

func NewConfig(log zerolog.Logger) *Config {
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}
	log.Info().Msg("Loading config...")
	// Db Config load
	DbConfig := new(DbConfig)
	DbConfig.Addr = env.GetString("DB_ADDR", "postgres://auth_user:auth_user@localhost:5432/auth_user?sslmode=disable")
	DbConfig.MigratorDriver = env.GetString("DB_MIGRATOR_DRIVER", "postgres")
	DbConfig.Username = env.GetString("DB_USERNAME", "auth_user")
	DbConfig.Password = env.GetString("DB_PASSWORD", "auth_user")
	DbConfig.Name = env.GetString("DB_NAME", "auth_user")
	DbConfig.Host = env.GetString("DB_HOST", "localhost")
	DbConfig.Port = env.GetString("DB_PORT", "5432")
	DbConfig.SSLMode = env.GetString("DB_SSLMODE", "disable")
	DbConfig.MaxOpenConn = env.GetInt("DB_MAX_OPEN_CONNS", 30)
	DbConfig.MaxIdleConn = env.GetInt("DB_MAX_IDLE_CONNS", 30)
	DbConfig.MaxIdleTime = env.GetInt("DB_MAX_IDLE_TIME", 15)
	DbConfig.MaxLifetime = env.GetInt("DB_MAX_LIFETIME", 60)
	DbConfig.MaxConnWaitTime = env.GetInt("DB_MAX_CONN_WAIT_TIME", 1)
	DbConfig.MaxConnLifetime = env.GetInt("DB_MAX_CONN_LIFETIME", 3600)
	DbConfig.MaxConnIdleTime = env.GetInt("DB_MAX_CONN_IDLE_TIME", 60)

	// Application config load
	config := new(Config)
	config.Addr = env.GetString("ADDR", ":3000")
	config.Db = *DbConfig

	log.Info().Msgf("Config loaded: %+v", config)
	return config
}
