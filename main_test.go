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
	f.Num = 10
	f.Workers = 2
	f.Seed = 1

	gofakeit.Seed(f.Seed)
	driver := drivers.New(f.Driver)
	testable := drivers.NewTestable(f.Driver)
	db := connector.Connection(driver, f)
	defer db.Close()
	if _, err := db.Exec(fmt.Sprintf("DROP TABLE IF EXISTS %s", f.Table)); err != nil {
		t.Fatal(err)
	}
	if err := testable.TestTable(db, "single", f.Table); err != nil {
		t.Fatal(err)
	}
	fields, err := driver.Describe(f.Table, db)
	if err != nil {
		t.Fatal(err.Error())
	}
	err = fuzzer.Run(fields, f)
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

func TestFuzzPostgres(t *testing.T) {
	f := flags.Flags{}
	f.Driver = drivers.Flags{
		Username: "test",
		Password: "test",
		Database: "test",
		Host:     "localhost",
		Port:     "5432",
		Driver:   "postgres",
	}
	f.Table = "Persons"
	f.Parsed = true
	f.Num = 10
	f.Workers = 2
	f.Seed = 1

	gofakeit.Seed(f.Seed)
	driver := drivers.New(f.Driver)
	testable := drivers.NewTestable(f.Driver)
	db := connector.Connection(driver, f)
	defer db.Close()
	if _, err := db.Exec(fmt.Sprintf("DROP TABLE IF EXISTS %s", f.Table)); err != nil {
		t.Fatal(err)
	}
	if err := testable.TestTable(db, "single", f.Table); err != nil {
		t.Fatal(err)
	}
	fields, err := driver.Describe(f.Table, db)
	if err != nil {
		t.Fatal(err.Error())
	}
	err = fuzzer.Run(fields, f)
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
		columns, err := res.Columns()
		if err != nil {
			t.Error(err)
		}
		var valueHolders []interface{}
		for j := 0; j < len(columns); j++ {
			var tempInterface interface{}
			valueHolders = append(valueHolders, &tempInterface)
		}
		if err := res.Scan(valueHolders...); err != nil {
			t.Error(err)
			continue
		}
		for index, val := range valueHolders {
			// If given column index is valid and value is nil, return error
			if val == nil {
				t.Error(fmt.Sprintf("Invalid value received for column  %s", columns[index]))
			}
		}
		i++
	}
	if i == 0 {
		t.Error("the table should not be empty")
	}
}

func TestMysqlMultiInsert(t *testing.T) {
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
	f.Num = 10
	f.Workers = 2
	f.Seed = 1

	gofakeit.Seed(f.Seed)
	driver := drivers.New(f.Driver)
	testable := drivers.NewTestable(f.Driver)
	test, err := testable.GetTestCase("multi")
	if err != nil {
		t.Error(fmt.Sprintf("postgres : error fetching test case for multi. %v", err.Error()))
	}
	db := connector.Connection(driver, f)
	defer db.Close()
	if _, err := db.Exec(fmt.Sprintf("DROP TABLE IF EXISTS %s", f.Table)); err != nil {
		t.Fatal(err)
	}
	if err := testable.TestTable(db, "multi", f.Table); err != nil {
		t.Fatal(err)
	}
	tables := test.TableCreationOrder
	tableFieldMap, insertionOrder, err := driver.MultiDescribe(tables, db)
	if err != nil {
		t.Errorf("Error describing tables %v. Error %v", tables, err)
	}
	err = fuzzer.RunMulti(tableFieldMap, insertionOrder, f)
	if err != nil {
		t.Errorf("error during multi insert %v", err.Error())
	}
}

func TestPostgresMultiInsert(t *testing.T) {
	f := flags.Flags{}
	f.Driver = drivers.Flags{
		Username: "test",
		Password: "test",
		Database: "test",
		Host:     "localhost",
		Port:     "5432",
		Driver:   "postgres",
	}
	f.Table = "Persons"
	f.Parsed = true
	f.Num = 10
	f.Workers = 2
	f.Seed = 1

	gofakeit.Seed(f.Seed)
	driver := drivers.New(f.Driver)
	testable := drivers.NewTestable(f.Driver)
	test, err := testable.GetTestCase("multi")
	if err != nil {
		t.Error(fmt.Sprintf("postgres : error fetching test case for multi. %v", err.Error()))
	}
	db := connector.Connection(driver, f)
	defer db.Close()
	if _, err := db.Exec(fmt.Sprintf("DROP TABLE IF EXISTS %s", f.Table)); err != nil {
		t.Fatal(err)
	}
	if err := testable.TestTable(db, "multi", f.Table); err != nil {
		t.Fatal(err)
	}
	tables := test.TableCreationOrder
	tableFieldMap, insertionOrder, err := driver.MultiDescribe(tables, db)
	if err != nil {
		t.Errorf("Error describing tables %v. Error %v", tables, err)
	}
	err = fuzzer.RunMulti(tableFieldMap, insertionOrder, f)
	if err != nil {
		t.Errorf("error during multi insert %v", err.Error())
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
