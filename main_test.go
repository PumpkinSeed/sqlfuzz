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
		Host:     "localhost",
		Port:     "3306",
		Driver:   "mysql",
	}
	f.Table = "Persons"
	f.Parsed = true

	gofakeit.Seed(0)
	driver := drivers.New(f.Driver)
	testable := drivers.NewTestable(f.Driver)
	db := connector.Connection(driver)

	if _, err := db.Exec(fmt.Sprintf("DROP TABLE IF EXISTS %s", f.Table)); err != nil {
		t.Fatal(err)
	}
	if err := testable.TestTable(db, f.Table); err != nil {
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

	res, err := db.Query(fmt.Sprintf("SELECT * FROM %s", f.Table))
	if err != nil {
		t.Fatal(err)
	}
	defer res.Close()
	var i int
	for res.Next() {
		tt := testTable{}
		if err := res.Scan(&tt.id, &tt.firstname, &tt.lastname, &tt.email, &tt.reqDate); err != nil {
			t.Error(err)
			continue
		}
		if err := tt.Validate(); err != nil {
			t.Error(err)
		}
		i++
	}
	if i == 0 {
		t.Error("the table should not be empty")
	}
}

type testTable struct {
	id        int
	firstname string
	lastname  string
	email     string
	reqDate   string
}

func (tt testTable) Validate() error {
	var err error
	if tt.id == 0 {
		err = warpErr(err, "id should not be 0")
	}
	if tt.firstname == "" {
		err = warpErr(err, "firstname should not be empty string")
	}
	if tt.lastname == "" {
		err = warpErr(err, "lastname should not be empty string")
	}
	if tt.email == "" {
		err = warpErr(err, "email should not be empty string")
	}
	if tt.reqDate == "" {
		err = warpErr(err, "reqDate should not be null")
	}

	return err
}

func warpErr(err error, msg string) error {
	if err != nil {
		return fmt.Errorf("%w;"+msg, err)
	}
	return fmt.Errorf(msg)
}
