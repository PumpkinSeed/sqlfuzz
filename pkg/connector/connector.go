package connector

import (
	"database/sql"
	"github.com/PumpkinSeed/sqlfuzz/drivers"
	_ "github.com/lib/pq"
	"log"
)

// Connection building a singleton connection to the database for give driver
func Connection(d drivers.Driver) *sql.DB {
	db, err := connect(d)
	if err != nil {
		log.Fatal(err)
		return nil
	}

	return db
}

// connect doing the direct connection open to the SQL database
func connect(d drivers.Driver) (*sql.DB, error) {
	return sql.Open(d.Driver(), d.Connection())
}
