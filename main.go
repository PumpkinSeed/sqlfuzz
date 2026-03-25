package main

import (
	"log"
	"time"

	"github.com/PumpkinSeed/sqlfuzz/drivers"
	"github.com/PumpkinSeed/sqlfuzz/pkg/connector"
	"github.com/PumpkinSeed/sqlfuzz/pkg/flags"
	"github.com/PumpkinSeed/sqlfuzz/pkg/fuzzer"
	"github.com/brianvoe/gofakeit/v5"
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	f := flags.Get()
	gofakeit.Seed(int64(f.Seed))
	driver, err := drivers.New(f.Driver)
	if err != nil {
		log.Fatal(err)
	}
	db, err := connector.Connection(driver, f)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	var tables []string
	if f.Table == "" {
		tables, err = driver.ShowTables(db)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		tables = []string{f.Table}
	}
	for _, table := range tables {
		f.Table = table
		fields, err := driver.Describe(f.Table, db)
		if err != nil {
			log.Fatal(err)
		}
		t := time.Now()
		if err := fuzzer.Run(fields, f); err != nil {
			log.Fatal(err)
		}
		log.Printf("Fuzzing %s table taken: %v \n", table, time.Since(t))
	}
}
