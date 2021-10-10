package mysql

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/PumpkinSeed/sqlfuzz/drivers/types"
	"github.com/PumpkinSeed/sqlfuzz/drivers/utils"
)

const (
	MySQLDescribeTableQuery = "SHOW TABLES;"
	mysqlFKQuery            = `SELECT CONSTRAINT_NAME,TABLE_NAME,COLUMN_NAME,REFERENCED_TABLE_NAME,REFERENCED_COLUMN_NAME 
							   from INFORMATION_SCHEMA.KEY_COLUMN_USAGE 
                               where REFERENCED_TABLE_NAME <> 'NULL' and REFERENCED_COLUMN_NAME <> 'NULL' and TABLE_NAME = '%s'`
)

var (
	mySQLNameToTestCase = map[string]types.TestCase{
		"single": {
			TableToCreateQueryMap: map[string]string{utils.DefaultTableCreateQueryKey: `CREATE TABLE %s (
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
				"t_product": `CREATE TABLE IF NOT EXISTS t_product( id int not null,name text not null,currency_id int,
                              PRIMARY KEY (id), FOREIGN KEY (currency_id) REFERENCES t_currency(id));`,
				"t_product_desc": `CREATE TABLE IF NOT EXISTS t_product_desc (id int not null,product_id  int ,	description text not null,
                                   PRIMARY KEY (id), FOREIGN KEY (product_id) REFERENCES t_currency(id) );`,
				"t_product_stock": `CREATE TABLE IF NOT EXISTS 
									t_product_stock(product_id  int, location_id int, amount numeric not null, 
								    FOREIGN KEY (product_id) REFERENCES t_currency(id),FOREIGN KEY(location_id) REFERENCES t_location(id));`,
			},
			TableCreationOrder: []string{"t_currency", "t_location", "t_product", "t_product_desc", "t_product_stock"},
		},
	}
)

// MySQL implementation of the Driver
type MySQL struct {
	f types.Flags
}

func New(f types.Flags) MySQL {
	return MySQL{f: f}
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
//nolint:gocognit,cyclop
func (m MySQL) MapField(descriptor types.FieldDescriptor) types.Field {
	field := strings.ToLower(descriptor.Type)
	// String types
	if strings.HasPrefix(field, "varchar") {
		l := length(field, "varchar")
		if l == nil || len(l) < 1 {
			return types.Field{Type: types.Unknown, Length: -1}
		}
		return types.Field{Type: types.String, Length: l[0]}
	}
	if strings.HasPrefix(field, "char") {
		l := length(field, "char")
		if l == nil || len(l) < 1 {
			return types.Field{Type: types.Unknown, Length: -1}
		}
		return types.Field{Type: types.String, Length: l[0]}
	}
	if strings.HasPrefix(field, "varbinary") {
		l := length(field, "varbinary")
		if l == nil || len(l) < 1 {
			return types.Field{Type: types.Unknown, Length: -1}
		}
		return types.Field{Type: types.String, Length: l[0]}
	}
	if strings.HasPrefix(field, "binary") {
		l := length(field, "binary")
		if l == nil || len(l) < 1 {
			return types.Field{Type: types.Unknown, Length: -1}
		}
		return types.Field{Type: types.String, Length: l[0]}
	}

	// Numeric types
	if strings.HasPrefix(field, "tinyint") {
		return types.Field{Type: types.Bool, Length: -1}
	}
	if strings.HasPrefix(field, "smallint") {
		return types.Field{Type: types.Int16, Length: -1}
	}
	if strings.HasPrefix(field, "mediumint") {
		return types.Field{Type: types.Int16, Length: -1}
	}
	if strings.HasPrefix(field, "int") || strings.HasPrefix(field, "bigint") {
		return types.Field{Type: types.Int32, Length: -1}
	}

	// Float types
	if strings.HasPrefix(field, "decimal") {
		l := length(field, "decimal")
		if l == nil || len(l) < 2 || l[0] < l[1] {
			return types.Field{Type: types.Unknown, Length: -1}
		}
		return types.Field{Type: types.Float, Length: l[0] - l[1]}
	}
	if strings.HasPrefix(field, "float") {
		l := length(field, "float")
		if l == nil || len(l) < 2 || l[0] < l[1] {
			return types.Field{Type: types.Unknown, Length: -1}
		}
		return types.Field{Type: types.Float, Length: l[0] - l[1]}
	}
	if strings.HasPrefix(field, "double") {
		l := length(field, "double")
		if l == nil || len(l) < 2 || l[0] < l[1] {
			return types.Field{Type: types.Unknown, Length: -1}
		}
		return types.Field{Type: types.Float, Length: l[0] - l[1]}
	}

	// Blob
	if strings.HasPrefix(field, "blob") || strings.HasPrefix(field, "tinyblob") ||
		strings.HasPrefix(field, "mediumblob") || strings.HasPrefix(field, "longblob") {
		return types.Field{Type: types.Blob, Length: -1}
	}

	// Text
	if strings.HasPrefix(field, "text") || strings.HasPrefix(field, "tinytext") ||
		strings.HasPrefix(field, "mediumtext") || strings.HasPrefix(field, "longtext") {
		return types.Field{Type: types.Text, Length: -1}
	}

	// Json
	if strings.HasPrefix(field, "json") {
		return types.Field{Type: types.Json, Length: -1}
	}

	// Year
	if strings.HasPrefix(field, "year") {
		return types.Field{Type: types.Year, Length: 4}
	}

	// Time
	// Date
	// Timestamp
	// Datetime
	if strings.HasPrefix(field, "datetime") || strings.HasPrefix(field, "date") ||
		strings.HasPrefix(field, "timestamp") || strings.HasPrefix(field, "time") {
		return types.Field{Type: types.Time, Length: -1}
	}

	// Enum
	if strings.HasPrefix(field, "enum") {
		f := strings.ReplaceAll(field, "enum(", "")
		f = strings.ReplaceAll(f, ")", "")
		f = strings.ReplaceAll(f, "'", "")
		f = strings.ReplaceAll(f, " ", "")
		return types.Field{Type: types.Enum, Length: -1, Enum: strings.Split(f, ",")}
	}

	return types.Field{Type: types.Unknown, Length: -1}
}

func (m MySQL) MultiDescribe(tables []string, db *sql.DB) (tableToDescriptorMap map[string][]types.FieldDescriptor, insertionOrder []string, err error) {
	processedTables := make(map[string]struct{})
	tableToDescriptorMap = make(map[string][]types.FieldDescriptor)
	for {
		newTableToDescriptorMap, newlyReferencedTables, err := utils.MultiDescribeHelper(tables, processedTables, db, m)
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
	insertionOrder, err = utils.GetInsertionOrder(tableToDescriptorMap)
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
		if err := rows.Scan(&val); err != nil {
			return nil, err
		}
	}
	return val, nil
}

// TestTable only for test purposes
func (m MySQL) TestTable(db *sql.DB, testCase, table string) error {
	return utils.TestTable(db, testCase, table, m)
}

func (MySQL) GetTestCase(name string) (types.TestCase, error) {
	if val, ok := mySQLNameToTestCase[name]; ok {
		return val, nil
	}
	return types.TestCase{}, fmt.Errorf("postgres: Error getting testcase with name %v", name)
}

func parseMySQLFields(results, fkRows *sql.Rows) ([]types.FieldDescriptor, error) {
	var fields []types.FieldDescriptor
	columnToFKMap := make(map[string]types.FKDescriptor)
	for fkRows.Next() {
		var fk types.FKDescriptor
		err := fkRows.Scan(&fk.ConstraintName, &fk.TableName, &fk.ColumnName, &fk.ForeignTableName, &fk.ForeignColumnName)
		if err != nil {
			return nil, err
		}
		columnToFKMap[fk.ColumnName] = fk
	}
	for results.Next() {
		var d types.FieldDescriptor
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
