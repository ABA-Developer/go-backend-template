package db

import (
	"database/sql"

	config "be-dashboard-nba/internal/infrastructure/runtime"
)

func openSQL(driver, source string, opt *config.Config) (db *sql.DB, err error) {
	db, err = sql.Open(driver, source)
	if err != nil {
		return
	}

	db.SetMaxOpenConns(opt.DB.MaxOpenConn)
	db.SetMaxIdleConns(opt.DB.MaxIdleConn)
	db.SetConnMaxIdleTime(opt.DB.MaxConnIdleTime)
	db.SetConnMaxLifetime(opt.DB.MaxConnLifetime)

	return
}
