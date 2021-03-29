package connector

import (
	"database/sql"
	"log"
	"time"

	"github.com/PumpkinSeed/sqlfuzz/drivers"
	_ "github.com/lib/pq"
)

// Connection building a singleton connection to the database for give driver
func Connection(d drivers.Driver) *sql.DB {
	db, err := connect(d)
	if err != nil {
		log.Fatal(err)
		return nil
	}
	db.SetConnMaxLifetime(100 * time.Second)
	db.SetMaxIdleConns(1000)
	db.SetMaxOpenConns(200)
	return db
}

// connect doing the direct connection open to the SQL database
func connect(d drivers.Driver) (*sql.DB, error) {
	return sql.Open(d.Driver(), d.Connection())
}
