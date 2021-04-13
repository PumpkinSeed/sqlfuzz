package drivers

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"
)

const (
	MySQLDescribeTableQuery = "SHOW TABLES;"
	mysqlFKQuery            = "SELECT CONSTRAINT_NAME,TABLE_NAME,COLUMN_NAME,REFERENCED_TABLE_NAME,REFERENCED_COLUMN_NAME from INFORMATION_SCHEMA.KEY_COLUMN_USAGE where REFERENCED_TABLE_NAME <> 'NULL' and REFERENCED_COLUMN_NAME <> 'NULL' and TABLE_NAME = '%s'"
)

var (
	mySQLNameToTestCase = map[string]TestCase{
		"single": {
			TableToCreateQueryMap: map[string]string{"": `CREATE TABLE %s (
		id INT(6) UNSIGNED,
		firstname VARCHAR(30),
		lastname VARCHAR(30),
		email VARCHAR(50),
		reg_date TIMESTAMP
		)`},
			TableCreationOrder: nil,
		},
		"multi": {
			TableToCreateQueryMap: map[string]string{
				"t_currency": "CREATE TABLE IF NOT EXISTS t_currency ( id int not null,shortcut char (3) not null,PRIMARY KEY (id));",
				"t_location": "CREATE TABLE IF NOT EXISTS t_location	( id int not null,location_name  text not null,PRIMARY KEY (id));",
				"t_product": "CREATE TABLE IF NOT EXISTS t_product( id int not null,name text not null,currency_id int ,PRIMARY KEY (id), FOREIGN KEY (currency_id) REFERENCES t_currency(id));",
				"t_product_desc": "CREATE TABLE IF NOT EXISTS t_product_desc (id int not null,product_id  int ,	description text not null,	PRIMARY KEY (id), FOREIGN KEY (product_id) REFERENCES t_currency(id) );",
				"t_product_stock": "CREATE TABLE IF NOT EXISTS t_product_stock(product_id  int ,	location_id int ,amount numeric not null, FOREIGN KEY (product_id) REFERENCES t_currency(id),FOREIGN KEY(location_id) REFERENCES t_location(id));",
			},
			TableCreationOrder: []string{"t_currency", "t_location", "t_product", "t_product_desc", "t_product_stock"},
		},
	}
)

// MySQL implementation of the Driver
type MySQL struct {
	f Flags
}

func (m MySQL) ShowTables(db *sql.DB) ([]string, error) {
	results, err := db.Query(MySQLDescribeTableQuery)
	if err != nil {
		return nil, err
	}
	defer results.Close()
	var tables []string
	for results.Next() {
		var table string
		if err := results.Scan(&table); err != nil {
			return nil, err
		}
		tables = append(tables, table)
	}

	return tables, nil
}

// Connection returns the specific connection string
func (m MySQL) Connection() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", m.f.Username, m.f.Password, m.f.Host, m.f.Port, m.f.Database)
}

// Driver returns the name of the driver
func (m MySQL) Driver() string {
	return m.f.Driver
}

// Insert inserts the data into
func (m MySQL) Insert(fields []string, table string) string {
	var template = "INSERT INTO %s(`%s`) VALUES(%s)"
	return fmt.Sprintf(template, table, strings.Join(fields, "`,`"), questionMarks(len(fields)))
}

// MapField returns the actual fields
func (m MySQL) MapField(descriptor FieldDescriptor) Field {
	field := strings.ToLower(descriptor.Type)
	// String types
	if strings.HasPrefix(field, "varchar") {
		l := length(field, "varchar")
		if l == nil || len(l) < 1 {
			return Field{Type: Unknown, Length: -1}
		}
		return Field{Type: String, Length: l[0]}
	}
	if strings.HasPrefix(field, "char") {
		l := length(field, "char")
		if l == nil || len(l) < 1 {
			return Field{Type: Unknown, Length: -1}
		}
		return Field{Type: String, Length: l[0]}
	}
	if strings.HasPrefix(field, "varbinary") {
		l := length(field, "varbinary")
		if l == nil || len(l) < 1 {
			return Field{Type: Unknown, Length: -1}
		}
		return Field{Type: String, Length: l[0]}
	}
	if strings.HasPrefix(field, "binary") {
		l := length(field, "binary")
		if l == nil || len(l) < 1 {
			return Field{Type: Unknown, Length: -1}
		}
		return Field{Type: String, Length: l[0]}
	}

	// Numeric types
	if strings.HasPrefix(field, "tinyint") {
		return Field{Type: Bool, Length: -1}
	}
	if strings.HasPrefix(field, "smallint") {
		return Field{Type: Int16, Length: -1}
	}
	if strings.HasPrefix(field, "mediumint") {
		return Field{Type: Int16, Length: -1}
	}
	if strings.HasPrefix(field, "int") || strings.HasPrefix(field, "bigint") {
		return Field{Type: Int32, Length: -1}
	}

	// Float types
	if strings.HasPrefix(field, "decimal") {
		l := length(field, "decimal")
		if l == nil || len(l) < 2 || l[0] < l[1] {
			return Field{Type: Unknown, Length: -1}
		}
		return Field{Type: Float, Length: l[0] - l[1]}
	}
	if strings.HasPrefix(field, "float") {
		l := length(field, "float")
		if l == nil || len(l) < 2 || l[0] < l[1] {
			return Field{Type: Unknown, Length: -1}
		}
		return Field{Type: Float, Length: l[0] - l[1]}
	}
	if strings.HasPrefix(field, "double") {
		l := length(field, "double")
		if l == nil || len(l) < 2 || l[0] < l[1] {
			return Field{Type: Unknown, Length: -1}
		}
		return Field{Type: Float, Length: l[0] - l[1]}
	}

	// Blob
	if strings.HasPrefix(field, "blob") || strings.HasPrefix(field, "tinyblob") ||
		strings.HasPrefix(field, "mediumblob") || strings.HasPrefix(field, "longblob") {
		return Field{Type: Blob, Length: -1}
	}

	// Text
	if strings.HasPrefix(field, "text") || strings.HasPrefix(field, "tinytext") ||
		strings.HasPrefix(field, "mediumtext") || strings.HasPrefix(field, "longtext") {
		return Field{Type: Text, Length: -1}
	}

	// Json
	if strings.HasPrefix(field, "json") {
		return Field{Type: Json, Length: -1}
	}

	// Year
	if strings.HasPrefix(field, "year") {
		return Field{Type: Year, Length: 4}
	}

	// Time
	// Date
	// Timestamp
	// Datetime
	if strings.HasPrefix(field, "datetime") || strings.HasPrefix(field, "date") ||
		strings.HasPrefix(field, "timestamp") || strings.HasPrefix(field, "time") {
		return Field{Type: Time, Length: -1}
	}

	// Enum
	if strings.HasPrefix(field, "enum") {
		f := strings.Replace(field, "enum(", "", -1)
		f = strings.Replace(f, ")", "", -1)
		f = strings.Replace(f, "'", "", -1)
		f = strings.Replace(f, " ", "", -1)
		return Field{Type: Enum, Length: -1, Enum: strings.Split(f, ",")}
	}

	return Field{Type: Unknown, Length: -1}
}

func (MySQL) Describe(table string, db *sql.DB) ([]FieldDescriptor, error) {
	describeQuery := fmt.Sprintf("DESCRIBE %s;", table)
	results, err := db.Query(describeQuery)
	if err != nil {
		return nil, err
	}
	fkRows, err := db.Query(fmt.Sprintf(mysqlFKQuery, strings.ToLower(table)))
	if err != nil {
		return nil, err
	}
	return parseMySQLFields(results, fkRows)
}

func (m MySQL) MultiDescribe(tables []string, db *sql.DB) (map[string][]FieldDescriptor, []string, error) {
	processedTables := make(map[string]bool)
	tableToDescriptorMap := make(map[string][]FieldDescriptor)
	for {
		newTableToDescriptorMap, newlyReferencedTables, err := multiDescribeHelper(tables, processedTables, db, m)
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
	insertionOrder, err := getInsertionOrder(tableToDescriptorMap)
	if err != nil {
		return nil, nil, err
	}
	return tableToDescriptorMap, insertionOrder, nil
}

func (MySQL) GetLatestColumnValue(table, column string, db *sql.DB) (interface{}, error) {
	query := fmt.Sprintf("select %v from %v order by %v desc limit 1", column, table, column)
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	var val interface{}
	for rows.Next() {
		rows.Scan(&val)
	}
	return val, nil
}

// TestTable only for test purposes
func (m MySQL) TestTable(db *sql.DB, testCase, table string) error {
	return testTable(db, testCase, table, m)
}

func (MySQL) GetTestCase(name string) (TestCase, error) {
	if val, ok := mySQLNameToTestCase[name]; ok {
		return val, nil
	}
	return TestCase{}, errors.New(fmt.Sprintf("postgres: Error getting testcase with name %v", name))
}

func parseMySQLFields(results, fkRows *sql.Rows) ([]FieldDescriptor, error) {
	var fields []FieldDescriptor
	columnToFKMap := make(map[string]FKDescriptor)
	for fkRows.Next() {
		var fk FKDescriptor
		err := fkRows.Scan(&fk.ConstraintName, &fk.TableName, &fk.ColumnName, &fk.ForeignTableName, &fk.ForeignColumnName)
		if err != nil {
			return nil, err
		}
		columnToFKMap[fk.ColumnName] = fk
	}
	for results.Next() {
		var d FieldDescriptor
		err := results.Scan(&d.Field, &d.Type, &d.Null, &d.Key, &d.Default, &d.Extra)
		if err != nil {
			return nil, err
		}
		if val, ok := columnToFKMap[d.Field]; ok {
			d.ForeignKeyDescriptor = &val
		}
		fields = append(fields, d)
	}
	return fields, nil
}

func questionMarks(n int) string {
	var q []string
	for i := 0; i < n; i++ {
		q = append(q, "?")
	}

	return strings.Join(q, ",")
}
