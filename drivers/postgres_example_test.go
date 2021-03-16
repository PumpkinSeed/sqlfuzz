package drivers

import (
	"database/sql"
	"fmt"
	"log"
)

const CreateTable = `CREATE TABLE IF NOT EXISTS pg_data_types (
   col1 bigint,
   col2 int8,
   col3 bit(10),
   col4 bit,
   col5 bit varying,
   col6 bit varying(10),
   col61 bit varying(128),
   col62 bit varying(256),
   col63 bit varying(512),
   col64 bit varying(1024),
   col7 varbit,
   col8 varbit(10),
   col9 boolean,
   col10 bool,
   col11 bytea,
   col12 character,
   col13 character(10),
   col14 char,
   col15 char(10),
   col16 character varying,
   col17 character varying(10),
   col18 varchar,
   col19 varchar(10),
   col20 cidr,
   col21 date,
   col22 double precision,
   col23 float8,
   col24 inet,
   col25 integer,
   col26 int,
   col27 int4,
   col28 json,
   col29 jsonb,
   col30 line,
   col31 macaddr,
   col32 money,
   col33 numeric(5,2),
   col34 decimal(5,2),
   col35 real,
   col36 float4,
   col37 smallint,
   col38 int2,
   col39 smallserial,
   col40 serial2,
   col41 serial,
   col42 serial4,
   col43 text,
   col44 time,
   col45 timetz,
   col46 timestamp,
   col47 uuid,
   col48 xml
);
`

func getConnection() (*sql.DB, error) {
	connectionString := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable ", "127.0.0.1", "5432", "postgres", "password", "fuzzpostgres")
	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		return nil, err
	}
	return db, nil
}

func createTable() {
	db, err := getConnection()
	if err != nil {
		return
	}
	db.Query(CreateTable)
}

func ExamplePostgres_ShowTables() {
	createTable()
	db, err := getConnection()
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
	db, err := getConnection()
	if err != nil {
		return
	}
	driver := Postgres{}
	fields, err := driver.DescribeFields("pg_data_types", db)
	if err != nil {
		log.Printf("Error describing table : %s", err.Error())
		return
	}
	for _, field := range fields {
		fmt.Println(fmt.Sprintf("Column name %s Field Type %s", field.Field, field.Type))
	}
}
