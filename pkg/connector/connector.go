package connector

import (
	"database/sql"

	"github.com/PumpkinSeed/sqlfuzz/drivers/types"
	"github.com/PumpkinSeed/sqlfuzz/pkg/flags"
	_ "github.com/lib/pq"
)

// Connection builds a connection to the database for the given driver
func Connection(d types.Driver, f flags.Flags) (*sql.DB, error) {
	db, err := connect(d)
	if err != nil {
		return nil, err
	}
	db.SetConnMaxLifetime(f.ConnMaxLifetimeInSec)
	db.SetMaxIdleConns(f.MaxIdleConns)
	db.SetMaxOpenConns(f.MaxOpenConns)
	return db, nil
}

// connect doing the direct connection open to the SQL database
func connect(d types.Driver) (*sql.DB, error) {
	return sql.Open(d.Driver(), d.Connection())
}
