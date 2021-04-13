package drivers

import (
	"database/sql"
	"fmt"
	"log"
)

func getPostgresConnection() (*sql.DB, error) {
	connectionString := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable ", "127.0.0.1", "5432", "test", "test", "test")
	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		return nil, err
	}
	return db, nil
}

func createTable() {
	db, err := getPostgresConnection()
	if err != nil {
		return
	}
	db.Query(fmt.Sprintf(CreateTable, "test_table"))
}

func ExamplePostgres_ShowTables() {
	createTable()
	db, err := getPostgresConnection()
	if err != nil {
		return
	}
	driver := Postgres{}
	tables, err := driver.ShowTables(db)
	if err != nil {
		log.Printf("Error showing tables : %s", err.Error())
		return
	}
	for _, table := range tables {
		fmt.Println(table)
	}
}

func ExamplePostgres_DescribeFields() {
	createTable()
	db, err := getPostgresConnection()
	if err != nil {
		return
	}
	driver := Postgres{}
	fields, err := driver.Describe("pg_data_types", db)
	if err != nil {
		log.Printf("Error describing table : %s", err.Error())
		return
	}
	for _, field := range fields {
		fmt.Println(fmt.Sprintf("Column name %s Field Type %s", field.Field, field.Type))
	}
}
