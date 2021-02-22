package main

import (
	"fmt"
	"log"
	"time"

	"github.com/PumpkinSeed/sqlfuzz/drivers"
	"github.com/PumpkinSeed/sqlfuzz/pkg/connector"
	"github.com/PumpkinSeed/sqlfuzz/pkg/descriptor"
	"github.com/PumpkinSeed/sqlfuzz/pkg/flags"
	"github.com/PumpkinSeed/sqlfuzz/pkg/fuzzer"
	"github.com/brianvoe/gofakeit/v5"
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	f := flags.Get()
	gofakeit.Seed(0)
	db := connector.Connection(drivers.New(flags.Get().Driver))
	fields, err := descriptor.Describe(db, f)
	if err != nil {
		log.Fatal(err.Error())
	}

	defer db.Close()

	t := time.Now()
	err = fuzzer.Run(db, fields, f)
	if err != nil {
		log.Fatal(err.Error())
	}
	fmt.Println("Fuzzing taken: ", time.Since(t))
}
