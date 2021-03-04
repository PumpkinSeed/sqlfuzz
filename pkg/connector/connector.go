package connector

import (
	"database/sql"
	"github.com/PumpkinSeed/sqlfuzz/drivers"
	_ "github.com/lib/pq"
	"log"
)

var (
	db *sql.DB
)

// Connection building a singleton connection to the SQL database
func Connection(d drivers.Driver) *sql.DB {
	if db == nil {
		connect(d)
	}

	return db
}

// connect doing the direct connection open to the SQL database
func connect(d drivers.Driver) {
	var err error
	db, err = sql.Open(d.Driver(), d.Connection())
	if err != nil {
		log.Fatal(err)
	}
}
