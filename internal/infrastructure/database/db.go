package db

import (
	"context"
	"database/sql"

	"github.com/pkg/errors"
)

func NewDatabase(ctx context.Context) (db *sql.DB, err error) {
	opt, err := newDatabaseOption(ctx)
	if err != nil {
		return
	}

	switch opt.driver {
	case "postgresql":
		db, err = NewPostgresql(opt)
	case "mysql":
		db, err = NewMySQL(opt)
	default:
		err = errors.Wrapf(errors.New("invalid datasources driver"), "db: driver=%s", opt.driver)
	}

	return
}
