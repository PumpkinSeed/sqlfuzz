package connector

import (
	"database/sql"
	"github.com/PumpkinSeed/sqlfuzz/drivers"
	_ "github.com/lib/pq"
	"log"
	"sync"
)

var (
	dbMap = make(map[string]*sql.DB)
	mu    = sync.Mutex{}
)

// Connection building a singleton connection to the SQL database
func Connection(d drivers.Driver) *sql.DB {
	mu.Lock()
	defer mu.Unlock()
	if db, ok := dbMap[d.Driver()]; ok {
		return db
	}
	db, err := connect(d)
	if err != nil {
		log.Fatal(err)
		return nil
	}
	dbMap[d.Driver()] = db
	return db
}

func Close(d drivers.Driver) error {
	mu.Lock()
	defer mu.Unlock()
	db, ok := dbMap[d.Driver()]
	if !ok {
		return nil
	}
	delete(dbMap, d.Driver())
	err := db.Close()
	return err
}

// connect doing the direct connection open to the SQL database
func connect(d drivers.Driver) (*sql.DB, error) {
	return sql.Open(d.Driver(), d.Connection())
}
