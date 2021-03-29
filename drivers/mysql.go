package drivers

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
)

const (
	MySQLDescribeTableQuery = "SHOW TABLES;"
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

func (MySQL) DescribeFields(table string, db *sql.DB) ([]FieldDescriptor, error) {
	describeQuery := fmt.Sprintf("DESCRIBE %s;", table)
	results, err := db.Query(describeQuery)
	if err != nil {
		return nil, err
	}
	return parseMySQLFields(results)
}

func (MySQL) MultiDescribe(tables []string, db *sql.DB) (map[string][]FieldDescriptor, []string, error) {
	return nil, nil, errors.New("error : Not yet implemented")
}

// TestTable only for test purposes
func (m MySQL) TestTable(db *sql.DB, table string) error {
	query := `CREATE TABLE %s (
		id INT(6) UNSIGNED,
		firstname VARCHAR(30),
		lastname VARCHAR(30),
		email VARCHAR(50),
		reg_date TIMESTAMP
	)`

	res, err := db.ExecContext(context.Background(), fmt.Sprintf(query, table))
	if err != nil {
		return err
	}

	_, err = res.RowsAffected()
	if err != nil {
		return err
	}
	return nil
}

func parseMySQLFields(results *sql.Rows) ([]FieldDescriptor, error) {
	var fields []FieldDescriptor
	for results.Next() {
		var d FieldDescriptor
		err := results.Scan(&d.Field, &d.Type, &d.Null, &d.Key, &d.Default, &d.Extra)
		if err != nil {
			return nil, err
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
