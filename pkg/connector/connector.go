package connector

import (
	"database/sql"
	"log"

	"github.com/PumpkinSeed/sqlfuzz/drivers/types"
	"github.com/PumpkinSeed/sqlfuzz/pkg/flags"
	_ "github.com/lib/pq"
)

// Connection building a singleton connection to the database for give driver
func Connection(d types.Driver, f flags.Flags) *sql.DB {
	db, err := connect(d)
	if err != nil {
		log.Fatal(err)
		return nil
	}
	db.SetConnMaxLifetime(f.ConnMaxLifetimeInSec)
	db.SetMaxIdleConns(f.MaxIdleConns)
	db.SetMaxOpenConns(f.MaxOpenConns)
	return db
}

// connect doing the direct connection open to the SQL database
func connect(d types.Driver) (*sql.DB, error) {
	return sql.Open(d.Driver(), d.Connection())
}
