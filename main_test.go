package main

import (
	"testing"

	"github.com/PumpkinSeed/sqlfuzz/drivers"
	"github.com/PumpkinSeed/sqlfuzz/pkg/connector"
	"github.com/PumpkinSeed/sqlfuzz/pkg/flags"
	"github.com/PumpkinSeed/sqlfuzz/pkg/fuzzer"
	"github.com/brianvoe/gofakeit/v5"
)

func TestFuzz(t *testing.T) {
	f := flags.Flags{}
	f.Driver = drivers.Flags{
		Username: "mysql",
		Password: "mysql",
		Database: "mysql",
		Host:     "127.0.0.1",
		Port:     "3306",
		Driver:   "mysql",
	}
	f.Table = "Persons"
	f.Parsed = true

	gofakeit.Seed(0)
	driver := drivers.New(flags.Get().Driver)
	db := connector.Connection(driver)
	describeQuery := driver.Describe(f.Table)
	results, err := db.Query(describeQuery)
	if err != nil {
		t.Fatal(err.Error())
	}
	fields, err := driver.ParseFields(results)
	if err != nil {
		t.Fatal(err.Error())
	}

	defer db.Close()

	err = fuzzer.Run(db, fields, f)
	if err != nil {
		t.Fatal(err)
	}
}
