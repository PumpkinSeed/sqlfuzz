package drivers

import (
	"fmt"

	"github.com/PumpkinSeed/sqlfuzz/drivers/mysql"
	"github.com/PumpkinSeed/sqlfuzz/drivers/postgres"
	"github.com/PumpkinSeed/sqlfuzz/drivers/types"
)

// New creates a new driver instance based on the flags
func New(f types.Flags) (types.Driver, error) {
	switch f.Driver {
	case "mysql":
		return mysql.New(f), nil
	case "postgres":
		return postgres.New(f), nil
	default:
		return nil, fmt.Errorf("driver %q not implemented", f.Driver)
	}
}

func NewTestable(f types.Flags) (types.Testable, error) {
	switch f.Driver {
	case "mysql":
		return mysql.New(f), nil
	case "postgres":
		return postgres.New(f), nil
	default:
		return nil, fmt.Errorf("testable %q not implemented", f.Driver)
	}
}
