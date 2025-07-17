package db

import (
	"database/sql"
	"time"

	"be-dashboard-nba/internal/config"
)

func openSQL(driver, source string, opt *config.Config) (db *sql.DB, err error) {
	db, err = sql.Open(driver, source)
	if err != nil {
		return
	}

	db.SetMaxOpenConns(opt.DB.MaxOpenConn)
	db.SetMaxIdleConns(opt.DB.MaxIdleConn)
	db.SetConnMaxIdleTime(time.Duration(opt.DB.MaxConnIdleTime))

	return
}
