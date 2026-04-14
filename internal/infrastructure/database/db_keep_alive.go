package db

import (
	"database/sql"
	"log"
	"time"
)

func keepAlive(db *sql.DB, driver, schema string, interval time.Duration) {
	for {
		err := db.Ping()
		if err != nil {
			log.Printf("ERROR db.keepAlive driver=%s schema=%s \n%s \n\n", driver, schema, err)
		}

		time.Sleep(interval)
	}
}
