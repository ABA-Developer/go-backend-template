package config

import (
	"github.com/joho/godotenv"
	"github.com/rs/zerolog"

	"be-dashboard-nba/constant"
	"be-dashboard-nba/internal/env"
)

type Swagger struct {
	Host      string
	BasePath  string
	IsEnabled bool
}

type Config struct {
	Name string
	Host string
	Port int
	DB   DBConfig
    Swagger Swagger
}

func NewConfig(log *zerolog.Logger) *Config {
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}

	// DB Config load
	DBConfig := new(DBConfig)
	DBConfig.MigratorDriver = env.MustGetEnv("DB_MIGRATOR_DRIVER")
	DBConfig.Username = env.MustGetEnv("DB_USERNAME")
	DBConfig.Password = env.MustGetEnv("DB_PASSWORD")
	DBConfig.Name = env.MustGetEnv("DB_NAME")
	DBConfig.Host = env.MustGetEnv("DB_HOST")
	DBConfig.Port = env.MustGetEnv("DB_PORT")
	DBConfig.SSLMode = env.GetString("DB_SSLMODE", constant.DefaultDBSSLMode)
	DBConfig.MaxOpenConn = env.GetInt("DB_MAX_OPEN_CONNS", constant.DefaultDBMaxOpenConns)
	DBConfig.MaxIdleConn = env.GetInt("DB_MAX_IDLE_CONNS", constant.DefaultDBMaxIdleConns)
	DBConfig.MaxIdleTime = env.GetInt("DB_MAX_IDLE_TIME", int(constant.DefaultDBMaxIdleTime))
	DBConfig.MaxLifetime = env.GetInt("DB_MAX_LIFETIME", int(constant.DefaultDBMaxLifetime))
	DBConfig.MaxConnWaitTime = env.GetInt("DB_MAX_CONN_WAIT_TIME", int(constant.DefaultDBMaxConnWaitTime))
	DBConfig.MaxConnLifetime = env.GetInt("DB_MAX_CONN_LIFETIME", int(constant.DefaultDBMaxConnLifetime))
	DBConfig.MaxConnIdleTime = env.GetInt("DB_MAX_CONN_IDLE_TIME", int(constant.DefaultDBMaxConnIdleTime))
	DBConfig.KeepAliveInterval = env.GetDuration("DB_KEEP_ALIVE_INTERVAL_CONN", constant.DefaultDBKeepAliveInterval)

	// swagger
	swagger := Swagger{
		Host:      env.GetString("SWAGGER_HOST", ""),
		BasePath:  env.GetString("SWAGGER_BASE_PATH", ""),
		IsEnabled: env.GetBool("SWAGGER_ENABLED", false),
	}

	// Application config load
	cfg := new(Config)
	cfg.Name = env.GetString("APP_NAME", constant.DefaultAppName)
	cfg.Host = env.GetString("APP_HOST", constant.DefaultAppHost)
	cfg.Port = env.GetInt("APP_PORT", constant.DefaultAppPort)
	cfg.DB = *DBConfig
	cfg.Swagger = swagger

	log.Info().Msgf("Config loaded: %+v", cfg)

	return cfg
}
