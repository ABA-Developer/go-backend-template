package db

import (
	"database/sql"

	"github.com/pkg/errors"
	"github.com/rs/zerolog"

	"be-dashboard-nba/internal/config"
)

func NewDatabase(config *config.Config,  log *zerolog.Logger) (db *sql.DB, err error) {
	log.Info().Msg("Connecting to database...")
	switch config.DB.MigratorDriver {
	case "postgresql":
		db, err = NewPostgresql(config)
	case "mysql":
		db, err = NewMySQL(config)
	default:
		err = errors.Wrapf(errors.New("invalid datasources driver"), "db: driver=%s", config.DB.MigratorDriver)
	}

	return
}
