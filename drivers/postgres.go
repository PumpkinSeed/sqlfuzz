package drivers

import (
	"database/sql"
	"fmt"
	"log"
	"strings"
)

type Postgres struct {
	f Flags
}

func (p Postgres) Connection() string {
	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s "+
		"password=%s dbname=%s sslmode=disable",
		p.f.Host, p.f.Port, p.f.Username, p.f.Password, p.f.Database)
	return psqlInfo
}

func (p Postgres) Driver() string {
	return p.f.Driver
}

func (p Postgres) Insert(fields []string, table string) string {
	// Syntax : VALUES($1, $2, $3)
	var template = "INSERT INTO %s(\"%s\") VALUES(%s)"
	return fmt.Sprintf(template, table, strings.Join(fields, "\",\""), pgValPlaceholder(len(fields)))
}

func (p Postgres) MapField(descriptor FieldDescriptor) Field {
	field := Field{Type: Unknown, Length: -1}
	// Has default values such as auto increments. No need to insert them specifically.
	switch descriptor.Type {
	case "bigint":
		return Field{Type: Int32, Length: -1}
	case "bit", "bit varying", "bytea":
		return Field{Type: BinaryString, Length: int16(descriptor.Length.Int)}
	case "character", "character varying":
		if descriptor.Length.Valid && descriptor.Length.Int > 0 {
			return Field{Type: String, Length: int16(descriptor.Length.Int)}
		}
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
	default:
		log.Printf("Field not identified. Name %s Length %d", descriptor.Field, descriptor.Length.Int)
	}
	return field
}

/*
ARRAY
select * from INFORMATION_SCHEMA.COLUMNS where table_name = 'pg_data_types' AND data_type = 'ARRAY';
*/
func (p Postgres) Describe(table string) string {
	//return fmt.Sprintf("SELECT col_attr.attname as 'ColumnName', pg_catalog.format_type(col_attr.atttypid, col_attr.atttypmod) as 'DataType'	FROM pg_catalog.pg_attribute col_attr WHERE col_attr.attnum > 0 AND NOT col_attr.attisdropped AND col_attr.attrelid = ( SELECT cls.oid FROM pg_catalog.pg_class cls LEFT JOIN pg_catalog.pg_namespace ns ON ns.oid = cls.relnamespace WHERE cls.relname = '%s');", table)
	return fmt.Sprintf("select column_name, data_type, character_maximum_length, column_default, is_nullable,numeric_precision,numeric_scale from INFORMATION_SCHEMA.COLUMNS where table_name = '%s'", table)

}

func (p Postgres) ParseFields(rows *sql.Rows) ([]FieldDescriptor, error) {
	var tableFields []FieldDescriptor
	for rows.Next() {
		var field FieldDescriptor
		err := rows.Scan(&field.Field, &field.Type, &field.Length, &field.Default, &field.Null, &field.Precision, &field.Scale)
		if field.Default.Valid && len(field.Default.String) > 0 {
			field.AutoIncrement = true
		}
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

func DescribeTable(db *sql.DB, table string) (sql.Result, error) {
	describeTableQuery := fmt.Sprintf("\\d+ %s;", table)
	//describeTableQuery := fmt.Sprintf("Select * from %s;", table)

	results, err := db.Exec(describeTableQuery)
	if err != nil {
		return nil, err
	}
	return results, nil
}
