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

var (
	f flags
	db *sql.DB
)

type flags struct {
	driver drivers.Flags

	table string
	parsed bool
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
	fields, err := describe()
	if err != nil {
		log.Fatal(err.Error())
	}

	defer db.Close()

}

func parseFlags() {
	flag.StringVar(&f.driver.Username, "u", "", "Username for the database connection")
	flag.StringVar(&f.driver.Password, "p", "", "Password for the database connection")
	flag.StringVar(&f.driver.Database, "d", "", "Database of the database connection")
	flag.StringVar(&f.driver.Host, "h", "localhost", "Host for the database connection")
	flag.StringVar(&f.driver.Port, "P", "3306", "Port for the database connection")
	flag.StringVar(&f.driver.Driver, "D", "mysql", "Driver for the database connection (mysql, postgres, etc.)")
	flag.StringVar(&f.table, "t", "", "Table for fuzzing")
	flag.Parse()

	f.parsed = true
}

func connect(d drivers.Driver) {
	var err error
	db, err = sql.Open(d.Driver(), d.Connection())
	if err != nil {
		log.Fatal(err)
	}
}

func describe() ([]fieldDescriptor, error) {
	results, err := connection().Query(fmt.Sprintf("DESCRIBE %s;", flagsOut().table))
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

func fuzz(fields []fieldDescriptor, table string) error {

	return nil
}

func exec(fields []fieldDescriptor, table string) error {
	driver := drivers.New(flagsOut().driver)

	var f []string
	var values []interface{}
	for _, field := range fields {
		f = append(f, field.Field)


		values = append(values, genField(driver, field.Type))
	}
	driver.Insert(f, flagsOut().table)
	ins, err := connection().Prepare(driver.Insert(f, flagsOut().table))
	if err != nil {
		log.Fatal(err)
	}


	ins.Exec(values...)

	return nil
}

func genField(driver drivers.Driver, t string) interface{} {
	typ, options := driver.MapField(t)
	switch typ {
	case drivers.String:
	case drivers.Uint:
	case drivers.Enum:
		fmt.Println(options)
	}

	return nil
}

func flagsOut() flags {
	if !f.parsed {
		parseFlags()
	}

	return f
}

func connection() *sql.DB {
	if db == nil {
		connect(drivers.New(flagsOut().driver))
	}

	return db
}