package drivers

import (
	"database/sql"
	"fmt"
	"log"
	"strings"

	"github.com/PumpkinSeed/sqlfuzz/drivers/types"
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
	PSQLDescribeTemplate = `select column_name, data_type, character_maximum_length, column_default, is_nullable,numeric_precision,numeric_scale
                            from INFORMATION_SCHEMA.COLUMNS where table_name = '%s'`
	PSQLConnectionTemplate = "host=%s port=%s user=%s password=%s dbname=%s sslmode=disable"
	PSQLInsertTemplate     = `INSERT INTO %s("%s") VALUES(%s)`
	PSQLShowTablesQuery    = "SELECT tablename FROM pg_catalog.pg_tables WHERE schemaname != 'pg_catalog' AND schemaname != 'information_schema';"
	psqlForeignKeysQuery   = `
	SELECT
    tc.constraint_name, 
    tc.table_name, 
    kcu.column_name, 
    ccu.table_name AS foreign_table_name,
    ccu.column_name AS foreign_column_name 
FROM 
    information_schema.table_constraints AS tc 
    JOIN information_schema.key_column_usage AS kcu
      ON tc.constraint_name = kcu.constraint_name
      AND tc.table_schema = kcu.table_schema
    JOIN information_schema.constraint_column_usage AS ccu
      ON ccu.constraint_name = tc.constraint_name
      AND ccu.table_schema = tc.table_schema
WHERE tc.constraint_type = 'FOREIGN KEY' AND tc.table_name='%s'
	`
)

var (
	pgNameToTestCase = map[string]types.TestCase{
		"single": {
			TableToCreateQueryMap: map[string]string{DefaultTableCreateQueryKey: CreateTable},
			TableCreationOrder:    nil,
		},
		"multi": {
			TableToCreateQueryMap: map[string]string{"t_currency": `CREATE TABLE IF NOT EXISTS "t_currency"
	(
		id      int not null,
		shortcut    char (3) not null,
		PRIMARY KEY (id)
	);
	`,
				"t_location": `CREATE TABLE IF NOT EXISTS "t_location"
	(
		id      int not null,
		location_name   text not null,
		PRIMARY KEY (id)
	);
	`,
				"t_product": `CREATE TABLE IF NOT EXISTS "t_product"
	(
		id      int not null,
		name        text not null,
		currency_id int REFERENCES t_currency (id) not null,
		PRIMARY KEY (id)
	);
	`,
				"t_product_desc": `CREATE TABLE IF NOT EXISTS "t_product_desc"
	(
		id      int not null,
		product_id  int REFERENCES t_product (id) not null,
		description text not null,
		PRIMARY KEY (id)
	);
	`,
				"t_product_stock": `CREATE TABLE IF NOT EXISTS "t_product_stock"
	(
		product_id  int REFERENCES t_product (id) not null,
		location_id int REFERENCES t_location (id) not null,
		amount      numeric not null
	);
	`,
			},
			TableCreationOrder: []string{"t_currency", "t_location", "t_product", "t_product_desc", "t_product_stock"},
		},
	}
)

type Postgres struct {
	f types.Flags
}

func (p Postgres) ShowTables(db *sql.DB) ([]string, error) {
	rows, err := db.Query(PSQLShowTablesQuery)
	if err != nil {
		return nil, err
	}
	var tables []string
	for rows.Next() {
		var table string
		if err := rows.Scan(&table); err != nil {
			return nil, err
		}
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

//nolint:cyclop
func (p Postgres) MapField(descriptor types.FieldDescriptor) types.Field {
	field := types.Field{Type: types.Unknown, Length: -1}
	switch descriptor.Type {
	case "bigint":
		return types.Field{Type: types.Int32, Length: -1}
	case "bit", "bit varying", "bytea":
		return types.Field{Type: types.BinaryString, Length: int16(descriptor.Length.Int)}
	case "character", "character varying":
		if descriptor.Length.Valid && descriptor.Length.Int > 0 {
			return types.Field{Type: types.String, Length: int16(descriptor.Length.Int)}
		}
		return types.Field{Type: types.String, Length: -1}
	case "date":
		return types.Field{Type: types.Time, Length: -1}
	case "double precision", "numeric", "real":
		return types.Field{Type: types.Float, Length: -1}
	case "integer":
		return types.Field{Type: types.Int32, Length: -1}
	case "json", "jsonb":
		return types.Field{Type: types.Json, Length: -1}
	case "smallint":
		return types.Field{Type: types.Int16, Length: -1}
	case "text":
		return types.Field{Type: types.Text, Length: -1}
	case "time without time zone", "time with time zone", "timestamp without time zone":
		return types.Field{Type: types.Time, Length: -1}
	case "xml":
		return types.Field{Type: types.XML, Length: -1}
	case "uuid":
		return types.Field{Type: types.UUID, Length: -1}
	case "boolean":
		return types.Field{Type: types.Bool, Length: -1}
	default:
		log.Printf("Field not identified. Name %s Length %d", descriptor.Field, descriptor.Length.Int)
	}
	return field
}

func (p Postgres) MultiDescribe(tables []string, db *sql.DB) (tableToDescriptorMap map[string][]types.FieldDescriptor, insertionOrder []string, err error) {
	processedTables := make(map[string]struct{})
	tableToDescriptorMap = make(map[string][]types.FieldDescriptor)
	for {
		newTableToDescriptorMap, newlyReferencedTables, err := multiDescribeHelper(tables, processedTables, db, p)
		if err != nil {
			return nil, nil, err
		}
		for key, val := range newTableToDescriptorMap {
			tableToDescriptorMap[key] = val
		}
		if len(newlyReferencedTables) == 0 {
			break
		}
		tables = newlyReferencedTables
	}
	insertionOrder, err = getInsertionOrder(tableToDescriptorMap)
	if err != nil {
		return nil, nil, err
	}
	return tableToDescriptorMap, insertionOrder, nil
}

func (p Postgres) Describe(table string, db *sql.DB) ([]types.FieldDescriptor, error) {
	results, err := db.Query(fmt.Sprintf(PSQLDescribeTemplate, strings.ToLower(table)))
	if err != nil {
		return nil, err
	}
	fkResults, err := db.Query(fmt.Sprintf(psqlForeignKeysQuery, strings.ToLower(table)))
	if err != nil {
		return nil, err
	}
	return parsePostgresFields(results, fkResults)
}

func (p Postgres) GetLatestColumnValue(table, column string, db *sql.DB) (interface{}, error) {
	query := fmt.Sprintf("select %s from %s order by %s desc limit 1", column, table, column)
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	var val interface{}
	for rows.Next() {
		if err := rows.Scan(&val); err != nil {
			return nil, err
		}
	}
	return val, nil
}

// TestTable only for test purposes
func (p Postgres) TestTable(db *sql.DB, testCase, table string) error {
	return testTable(db, testCase, table, p)
}

func (Postgres) GetTestCase(name string) (types.TestCase, error) {
	if val, ok := pgNameToTestCase[name]; ok {
		return val, nil
	}
	return types.TestCase{}, fmt.Errorf("postgres: Error getting testcase with name %v", name)
}

func parsePostgresFields(rows, fkRows *sql.Rows) ([]types.FieldDescriptor, error) {
	var tableFields []types.FieldDescriptor
	columnToFKMap := make(map[string]types.FKDescriptor)
	for fkRows.Next() {
		var fk types.FKDescriptor
		err := fkRows.Scan(&fk.ConstraintName, &fk.TableName, &fk.ColumnName, &fk.ForeignTableName, &fk.ForeignColumnName)
		if err != nil {
			return nil, err
		}
		columnToFKMap[fk.ColumnName] = fk
	}
	for rows.Next() {
		var field types.FieldDescriptor
		err := rows.Scan(&field.Field, &field.Type, &field.Length, &field.Default, &field.Null, &field.Precision, &field.Scale)
		field.HasDefaultValue = field.Default.Valid && len(field.Default.String) > 0
		if err != nil {
			return nil, err
		}
		if val, ok := columnToFKMap[field.Field]; ok {
			field.ForeignKeyDescriptor = &val
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
