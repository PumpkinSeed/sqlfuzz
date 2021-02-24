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

type Flags struct {
	Username string
	Password string
	Database string
	Host string
	Port string
	Driver string
}

type Field struct {
	Type Type
	Length int16
	Enum []string
}

type Driver interface {
	Connection() string
	Driver() string
	Insert(fields []string, table string) string
	MapField(string) Field
}

func New(f Flags) Driver {
	switch f.Driver {
	case "mysql":
		return MySQL{f:f}
	default:
		log.Fatal("Driver not implemented")
		return nil
	}
}
