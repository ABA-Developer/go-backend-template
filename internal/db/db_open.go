package db

import (
	"context"
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

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = db.PingContext(ctx)
	if err != nil {
		return
	}

	return
}
