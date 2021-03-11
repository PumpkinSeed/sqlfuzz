package main

import (
	"fmt"
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
		Username: "test",
		Password: "test",
		Database: "test",
		Host:     "mysql",
		Port:     "3306",
		Driver:   "mysql",
	}
	f.Table = "Persons"
	f.Parsed = true



	gofakeit.Seed(0)
	driver := drivers.New(f.Driver)
	db := connector.Connection(driver)

	if _, err := db.Exec(fmt.Sprintf("DROP TABLE IF EXISTS %s", f.Table)); err != nil {
		t.Fatal(err)
	}
	if err := driver.TestTable(db, f.Table); err != nil {
		t.Fatal(err)
	}
	fields, err := driver.DescribeFields(f.Table, db)
	if err != nil {
		t.Fatal(err.Error())
	}
	defer db.Close()
	err = fuzzer.Run(db, fields, f)
	if err != nil {
		t.Fatal(err)
	}
}
