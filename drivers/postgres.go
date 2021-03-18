package drivers

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"strings"
)

/*
Below columns are not yet implemented and removed from table. Add these once implemented.
col20 cidr
col24 inet
col30 line
col31 macaddr
col32 money
*/
const CreateTable = `CREATE TABLE IF NOT EXISTS %s (
   col1 bigint not null,
   col2 int8 not null,
   col3 bit(10) not null,
   col4 bit not null,
   col5 bit varying not null,
   col6 bit varying(10) not null,
   col64 bit varying(1024) not null,
   col7 varbit not null,
   col8 varbit(10) not null,
   col9 boolean not null,
   col10 bool not null,
   col11 bytea not null,
   col12 character not null,
   col13 character(10) not null,
   col14 char not null,
   col15 char(10) not null,
   col16 character varying not null,
   col17 character varying(10) not null,
   col18 varchar not null,
   col19 varchar(10) not null,
   col21 date not null,
   col22 double precision not null,
   col23 float8 not null,
   col25 integer not null,
   col26 int not null,
   col27 int4 not null,
   col28 json not null,
   col29 jsonb not null,
   col33 numeric(5,2) not null,
   col34 decimal(5,2) not null,
   col35 real not null,
   col36 float4 not null,
   col37 smallint not null,
   col38 int2 not null,
   col39 smallserial not null,
   col40 serial2 not null,
   col41 serial not null,
   col42 serial4 not null,
   col43 text not null,
   col44 time not null,
   col45 timetz not null,
   col46 timestamp not null,
   col47 uuid not null,
   col48 xml not null
);
`

const (
	PSQLDescribeTemplate   = "select column_name, data_type, character_maximum_length, column_default, is_nullable,numeric_precision,numeric_scale from INFORMATION_SCHEMA.COLUMNS where table_name = '%s'"
	PSQLConnectionTemplate = "host=%s port=%s user=%s password=%s dbname=%s sslmode=disable"
	PSQLInsertTemplate     = `INSERT INTO %s("%s") VALUES(%s)`
	PSQLShowTablesQuery    = "SELECT tablename FROM pg_catalog.pg_tables WHERE schemaname != 'pg_catalog' AND schemaname != 'information_schema';"
)

type Postgres struct {
	f Flags
}

func (p Postgres) ShowTables(db *sql.DB) ([]string, error) {
	rows, err := db.Query(PSQLShowTablesQuery)
	if err != nil {
		return nil, err
	}
	var tables []string
	for rows.Next() {
		var table string
		rows.Scan(&table)
		tables = append(tables, table)
	}
	return tables, nil
}

func (p Postgres) Connection() string {
	return fmt.Sprintf(PSQLConnectionTemplate,
		p.f.Host, p.f.Port, p.f.Username, p.f.Password, p.f.Database)
}

func (p Postgres) Driver() string {
	return p.f.Driver
}

func (p Postgres) Insert(fields []string, table string) string {
	return fmt.Sprintf(PSQLInsertTemplate, table, strings.Join(fields, `","`), pgValPlaceholder(len(fields)))
}

func (p Postgres) MapField(descriptor FieldDescriptor) Field {
	field := Field{Type: Unknown, Length: -1}
	switch descriptor.Type {
	case "bigint":
		return Field{Type: Int32, Length: -1}
	case "bit", "bit varying", "bytea":
		return Field{Type: BinaryString, Length: int16(descriptor.Length.Int)}
	case "character", "character varying":
		if descriptor.Length.Valid && descriptor.Length.Int > 0 {
			return Field{Type: String, Length: int16(descriptor.Length.Int)}
		}
		return Field{Type: String, Length: -1}
	case "date":
		return Field{Type: Time, Length: -1}
	case "double precision", "numeric", "real":
		return Field{Type: Float, Length: -1}
	case "integer":
		return Field{Type: Int32, Length: -1}
	case "json", "jsonb":
		return Field{Type: Json, Length: -1}
	case "smallint":
		return Field{Type: Int16, Length: -1}
	case "text":
		return Field{Type: Text, Length: -1}
	case "time without time zone", "time with time zone", "timestamp without time zone":
		return Field{Type: Time, Length: -1}
	case "xml":
		return Field{Type: XML, Length: -1}
	case "uuid":
		return Field{Type: UUID, Length: -1}
	case "boolean":
		return Field{Type: Bool, Length: -1}
	default:
		log.Printf("Field not identified. Name %s Length %d", descriptor.Field, descriptor.Length.Int)
	}
	return field
}

func (p Postgres) DescribeFields(table string, db *sql.DB) ([]FieldDescriptor, error) {
	results, err := db.Query(fmt.Sprintf(PSQLDescribeTemplate, strings.ToLower(table)))
	if err != nil {
		return nil, err
	}
	return parsePostgresFields(results)
}

// TestTable only for test purposes
func (p Postgres) TestTable(db *sql.DB, table string) error {
	res, err := db.ExecContext(context.Background(), fmt.Sprintf(CreateTable, table))
	if err != nil {
		return err
	}

	_, err = res.RowsAffected()
	if err != nil {
		return err
	}
	return nil
}

func parsePostgresFields(rows *sql.Rows) ([]FieldDescriptor, error) {
	var tableFields []FieldDescriptor
	for rows.Next() {
		var field FieldDescriptor
		err := rows.Scan(&field.Field, &field.Type, &field.Length, &field.Default, &field.Null, &field.Precision, &field.Scale)
		field.HasDefaultValue = field.Default.Valid && len(field.Default.String) > 0
		if err != nil {
			return nil, err
		}
		tableFields = append(tableFields, field)
	}
	return tableFields, nil
}

func pgValPlaceholder(fieldLen int) string {
	var q []string
	for i := 1; i <= fieldLen; i++ {
		q = append(q, fmt.Sprintf("$%d", i))
	}
	return strings.Join(q, ",")
}
