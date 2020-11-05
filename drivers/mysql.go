package drivers

import "fmt"

type MySQL struct {
	f Flags
}

func (m MySQL) Connection() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", m.f.Username, m.f.Password, m.f.Host, m.f.Port, m.f.Database)
}

func (m MySQL) Driver() string {
	return m.f.Driver
}
