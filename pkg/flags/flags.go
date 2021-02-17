package flags

import (
	"flag"

	"github.com/PumpkinSeed/sqlfuzz/drivers"
)

var f Flags

type Flags struct {
	Driver drivers.Flags

	Num     int
	Workers int
	Table   string
	Parsed  bool
}

func Get() Flags {
	if !f.Parsed {
		parseFlags()
	}

	return f
}

func parseFlags() {
	if !f.Parsed {
		flag.StringVar(&f.Driver.Username, "u", "", "Username for the database connection")
		flag.StringVar(&f.Driver.Password, "p", "", "Password for the database connection")
		flag.StringVar(&f.Driver.Database, "d", "", "Database of the database connection")
		flag.StringVar(&f.Driver.Host, "h", "", "Host for the database connection")
		flag.StringVar(&f.Driver.Port, "P", "3306", "Port for the database connection")
		flag.StringVar(&f.Driver.Driver, "D", "mysql", "Driver for the database connection (mysql, postgres, etc.)")
		flag.StringVar(&f.Table, "t", "", "Table for fuzzing")
		flag.IntVar(&f.Num, "n", 1000, "Number of rows")
		flag.IntVar(&f.Workers, "w", 20, "Number of workers")
		flag.Parse()
	}

	f.Parsed = true
}
