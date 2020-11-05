package main

import (
	"database/sql"
	"flag"
	"fmt"
	"github.com/PumpkinSeed/sqlfuzz/drivers"
	_ "github.com/go-sql-driver/mysql"
	"github.com/volatiletech/null"
	"log"
)

type flags struct {
	driver drivers.Flags

	table string
}

type fieldDescriptor struct {
	Field string
	Type string
	Null string
	Key string
	Default null.String
	Extra string
}

func main() {
	f := parse()

	d := drivers.New(f.driver)
	db, err := connect(d)
	if err != nil {
		log.Fatal(err.Error())
	}
	defer db.Close()

	fields, err := describe(db, f.table)
	if err != nil {
		log.Fatal(err.Error())
	}


}

func parse() flags {
	var f flags
	flag.StringVar(&f.driver.Username, "u", "", "Username for the database connection")
	flag.StringVar(&f.driver.Password, "p", "", "Password for the database connection")
	flag.StringVar(&f.driver.Database, "d", "", "Database of the database connection")
	flag.StringVar(&f.driver.Host, "h", "localhost", "Host for the database connection")
	flag.StringVar(&f.driver.Port, "P", "3306", "Port for the database connection")
	flag.StringVar(&f.driver.Driver, "D", "mysql", "Driver for the database connection (mysql, postgres, etc.)")
	flag.StringVar(&f.table, "t", "", "Table for fuzzing")
	flag.Parse()

	return f
}

func connect(d drivers.Driver) (*sql.DB, error) {
	return sql.Open(d.Driver(), d.Connection())
}

func describe(db *sql.DB, table string) ([]fieldDescriptor, error) {
	results, err := db.Query(fmt.Sprintf("DESCRIBE %s;", table))
	if err != nil {
		return nil, err
	}

	var fields []fieldDescriptor
	for results.Next() {
		var d fieldDescriptor

		err = results.Scan(&d.Field, &d.Type, &d.Null, &d.Key, &d.Default, &d.Extra)
		if err != nil {
			return nil, err
		}

		fields = append(fields, d)
	}

	return fields, nil
}

func fuzz(fields) error {

}