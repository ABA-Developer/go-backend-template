package db

import (
	"database/sql"

	"github.com/pkg/errors"

	"be-dashboard-nba/internal/config"
)

func NewDatabase(config *config.Config) (db *sql.DB, err error) {
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
