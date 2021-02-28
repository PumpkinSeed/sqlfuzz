package drivers

import "log"

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

// Driver is the interface should satisfied by a certain driver
type Driver interface {
	Connection() string
	Driver() string
	Insert(fields []string, table string) string
	MapField(string) Field
}

// New creates a new driver instance based on the flags
func New(f Flags) Driver {
	switch f.Driver {
	case "mysql":
		return MySQL{f: f}
	default:
		log.Fatal("Driver not implemented")
		return nil
	}
}
