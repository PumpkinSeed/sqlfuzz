package drivers

import (
	"database/sql"
	"log"

	"github.com/volatiletech/null"
)

type Type int16

const (
	String Type = iota
	Int16
	Int32
	Float
	Blob // base64
	Text
	Enum
	Bool
	Json
	Time
	Year
	XML
	UUID
	BinaryString
	Unknown
)

// Flags needed by the driver
type Flags struct {
	Username string
	Password string
	Database string
	Host     string
	Port     string
	Driver   string
}

// Field is the possible field definition
type Field struct {
	Type   Type
	Length int16
	Enum   []string
}

type FKDescriptor struct {
	ConstraintName    string
	TableName         string
	ColumnName        string
	ForeignTableName  string
	ForeignColumnName string
}

//FieldDescriptor represents a field described by the table in the SQL database
type FieldDescriptor struct {
	Field                string
	Type                 string
	Null                 string
	Key                  string
	Length               null.Int
	Default              null.String
	Extra                string
	Precision            null.Int
	Scale                null.Int
	HasDefaultValue      bool
	ForeignKeyDescriptor *FKDescriptor
}

// TestCase has a map of table to its create table query and table creation order
type TestCase struct {
	TableToCreateQueryMap map[string]string
	TableCreationOrder    []string
}

// Driver is the interface should satisfied by a certain driver
type Driver interface {
	ShowTables(db *sql.DB) ([]string, error)
	Connection() string
	Driver() string
	Insert(fields []string, table string) string
	MapField(descriptor FieldDescriptor) Field
	Describe(table string, db *sql.DB) ([]FieldDescriptor, error)
	MultiDescribe(tables []string, db *sql.DB) (map[string][]FieldDescriptor, []string, error)
	GetLatestColumnValue(table, column string, db *sql.DB) (interface{}, error)
	GetTestCase(name string) (TestCase, error)
}

type Testable interface {
	TestTable(conn *sql.DB, testCase, table string) error
}

// New creates a new driver instance based on the flags
func New(f Flags) Driver {
	switch f.Driver {
	case "mysql":
		return MySQL{f: f}
	case "postgres":
		return Postgres{f: f}
	default:
		log.Fatal("Driver not implemented")
		return nil
	}
}

func NewTestable(f Flags) Testable {
	switch f.Driver {
	case "mysql":
		return MySQL{f: f}
	case "postgres":
		return Postgres{f: f}
	default:
		log.Fatal("Testable not implemented")
		return nil
	}
}
