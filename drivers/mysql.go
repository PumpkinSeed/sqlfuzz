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

func (m MySQL) MapField(field string) Field {
	field = strings.ToLower(field)
	// String types
	if strings.Contains(field, "char") {
		l := length(field, "char")
		if l == nil || len(l) < 1 {
			return Field{Type: Unknown, Length: -1}
		}
		return Field{Type: String, Length: l[0]}
	}
	if strings.Contains(field, "varchar") {
		l := length(field, "varchar")
		if l == nil || len(l) < 1 {
			return Field{Type: Unknown, Length: -1}
		}
		return Field{Type: String, Length: l[0]}
	}
	if strings.Contains(field, "binary") {
		l := length(field, "binary")
		if l == nil || len(l) < 1 {
			return Field{Type: Unknown, Length: -1}
		}
		return Field{Type: String, Length: l[0]}
	}
	if strings.Contains(field, "varbinary") {
		l := length(field, "varbinary")
		if l == nil || len(l) < 1 {
			return Field{Type: Unknown, Length: -1}
		}
		return Field{Type: String, Length: l[0]}
	}

	// Numeric types
	if strings.Contains(field, "tinyint") {
		return Field{Type: Bool, Length: -1}
	}
	if strings.Contains(field, "smallint") {
		return Field{Type: Int16, Length: -1}
	}
	if strings.Contains(field, "mediumint") {
		return Field{Type: Int16, Length: -1}
	}
	if strings.Contains(field, "int") || strings.Contains(field, "bigint"){
		return Field{Type: Int32, Length: -1}
	}

	// Float types
	if strings.Contains(field, "decimal") {
		l := length(field, "decimal")
		if l == nil || len(l) < 2 || l[0] < l[1] {
			return Field{Type: Unknown, Length: -1}
		}
		return Field{Type: Float, Length: l[0] - l[1]}
	}
	if strings.Contains(field, "float") {
		l := length(field, "float")
		if l == nil || len(l) < 2 || l[0] < l[1] {
			return Field{Type: Unknown, Length: -1}
		}
		return Field{Type: Float, Length: l[0] - l[1]}
	}
	if strings.Contains(field, "double") {
		l := length(field, "double")
		if l == nil || len(l) < 2 || l[0] < l[1] {
			return Field{Type: Unknown, Length: -1}
		}
		return Field{Type: Float, Length: l[0] - l[1]}
	}

	// Blob
	if strings.Contains(field, "blob") {
		return Field{Type: Blob, Length: -1}
	}

	// Text
	if strings.Contains(field, "text") {
		return Field{Type: Text, Length: -1}
	}

	// Json
	if strings.Contains(field, "json") {
		return Field{Type: Json, Length: -1}
	}

	if strings.Contains(field, "datetime") {
		return Field{Type: Time, Length: -1}
	}

	// Enum
	if strings.Contains(field, "enum") {
		f := strings.Replace(field, "enum(", "", -1)
		f = strings.Replace(f, ")", "", -1)
		f = strings.Replace(f, "'", "", -1)
		return Field{Type: Enum, Length: -1, Enum: strings.Split(f, ",")}
	}

	return Field{Type: Unknown, Length: -1}
}

func questionMarks(n int) string {
	var q []string
	for i := 0; i < n; i++ {
		q = append(q, "?")
	}

	return strings.Join(q, ", ")
}
