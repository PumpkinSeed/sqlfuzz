package drivers

import (
	"fmt"
	"strings"
)

type MySQL struct {
	f Flags
}

func (m MySQL) Connection() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", m.f.Username, m.f.Password, m.f.Host, m.f.Port, m.f.Database)
}

func (m MySQL) Driver() string {
	return m.f.Driver
}

func (m MySQL) Insert(fields []string, table string) string {
	var template = "INSERT INTO %s(%s) VALUES(%s)"
	return fmt.Sprintf(template, table, strings.Join(fields, ", "), questionMarks(len(fields)))
}

func (m MySQL) MapField(field string) (Type, []string) {
	if strings.Contains(field, "varbinary") || strings.Contains(field, "char") {
		return String, nil
	}
	if strings.Contains(field, "tinyint") {
		return Bool, nil
	}
	if strings.Contains(field, "int") {
		return Uint, nil
	}
	if strings.Contains(field, "json") {
		return Json, nil
	}
	if strings.Contains(field, "datetime") {
		return Time, nil
	}
	if strings.Contains(field, "enum") {
		f := strings.Replace(field, "enum(", "", -1)
		f = strings.Replace(f, ")", "", -1)
		f = strings.Replace(f, "'", "", -1)
		return Enum, strings.Split(f, ",")
	}

	return Unknown, nil
}

func questionMarks(n int) string {
	var q []string
	for i := 0; i < n; i++ {
		q = append(q, "?")
	}

	return strings.Join(q, ", ")
}
