package flags

import (
	"flag"

	"github.com/PumpkinSeed/sqlfuzz/drivers"
)

var f Flags

// Flags represents the CLI flags
type Flags struct {
	Driver drivers.Flags

	Num     int
	Workers int
	Table   string
	Parsed  bool
}

// Get the parsed flags and parsing them if needed
func Get() Flags {
	if !f.Parsed {
		parse()
	}

	return f
}

// parse parsing the flags into the f variable
func parse() {
	if !f.Parsed {
		flag.StringVar(&f.Driver.Username, "u", "test", "Username for the database connection")
		flag.StringVar(&f.Driver.Password, "p", "test", "Password for the database connection")
		flag.StringVar(&f.Driver.Database, "d", "test", "Database of the database connection")
		flag.StringVar(&f.Driver.Host, "h", "localhost", "Host for the database connection")
		flag.StringVar(&f.Driver.Port, "P", "3306", "Port for the database connection")
		flag.StringVar(&f.Driver.Driver, "D", "mysql", "Driver for the database connection (mysql, postgres, etc.)")
		flag.StringVar(&f.Table, "t", "", "Table for fuzzing")
		flag.IntVar(&f.Num, "n", 1000, "Number of rows")
		flag.IntVar(&f.Workers, "w", 20, "Number of workers")
		flag.Parse()
	}

	f.Parsed = true
}
