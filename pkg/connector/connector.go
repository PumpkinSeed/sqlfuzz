package connector

import (
	"database/sql"
	"log"

	"github.com/PumpkinSeed/sqlfuzz/drivers"
)

var (
	db *sql.DB
)

func Connection(d drivers.Driver) *sql.DB {
	if db == nil {
		connect(d)
	}

	return db
}

func connect(d drivers.Driver) {
	var err error
	db, err = sql.Open(d.Driver(), d.Connection())
	if err != nil {
		log.Fatal(err)
	}
}
