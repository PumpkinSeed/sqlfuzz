package postgres

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
	if _, err = db.Query(fmt.Sprintf(CreateTable, "test_table")); err != nil {
		return
	}
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
