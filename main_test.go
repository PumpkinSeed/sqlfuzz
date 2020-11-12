package main

import (
	"github.com/PumpkinSeed/sqlfuzz/drivers"
	"github.com/brianvoe/gofakeit/v5"
	"log"
	"testing"
)

func TestFuzz(t *testing.T) {
	f.driver = drivers.Flags{
		Username: "fluidpay",
		Password: "fluidpay",
		Database: "fluidpay",
		Host: "10.0.0.12",
		Port: "3306",
		Driver:"mysql",
	}
	f.table = "Persons"
	f.parsed = true

	gofakeit.Seed(0)
	fields, err := describe()
	if err != nil {
		log.Fatal(err.Error())
	}

	defer db.Close()

	err = fuzz(fields)
	if err != nil {
		log.Fatal(err.Error())
	}
}
