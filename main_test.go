package main

import (
	"testing"

	"github.com/PumpkinSeed/sqlfuzz/drivers"
	"github.com/PumpkinSeed/sqlfuzz/pkg/connector"
	"github.com/PumpkinSeed/sqlfuzz/pkg/descriptor"
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
		Host: "127.0.0.1",
		Port: "3306",
		Driver:"mysql",
	}
	f.Table = "Persons"
	f.Parsed = true

	gofakeit.Seed(0)
	db := connector.Connection(drivers.New(flags.Get().Driver))
	fields, err := descriptor.Describe(db, f.Table)
	if err != nil {
		t.Fatal(err)
	}

	defer db.Close()

	err = fuzzer.Run(db, fields, f)
	if err != nil {
		t.Fatal(err)
	}
}
