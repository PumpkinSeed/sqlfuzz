package drivers

import (
	"log"

	"github.com/PumpkinSeed/sqlfuzz/drivers/mysql"
	"github.com/PumpkinSeed/sqlfuzz/drivers/postgres"
	"github.com/PumpkinSeed/sqlfuzz/drivers/types"
)

// New creates a new driver instance based on the flags
func New(f types.Flags) types.Driver {
	switch f.Driver {
	case "mysql":
		return mysql.New(f)
	case "postgres":
		return postgres.New(f)
	default:
		log.Fatal("Driver not implemented")
		return nil
	}
}

func NewTestable(f types.Flags) types.Testable {
	switch f.Driver {
	case "mysql":
		return mysql.New(f)
	case "postgres":
		return postgres.New(f)
	default:
		log.Fatal("Testable not implemented")
		return nil
	}
}
