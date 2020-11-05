package drivers

import "log"

type Flags struct {
	Username string
	Password string
	Database string
	Host string
	Port string
	Driver string
}

type Driver interface {
	Connection() string
	Driver() string
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
