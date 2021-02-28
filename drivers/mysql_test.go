package drivers

import (
	"reflect"
	"testing"
)

func TestMapField(t *testing.T) {
	var scenarios = []struct {
		input  string
		output Field
	}{
		{
			input:  "varchar(12)",
			output: Field{Type: String, Length: 12},
		},
		{
			input:  "char(100)",
			output: Field{Type: String, Length: 100},
		},
		{
			input:  "varbinary(100)",
			output: Field{Type: String, Length: 100},
		},
		{
			input:  "binary(100)",
			output: Field{Type: String, Length: 100},
		},
		{
			input:  "tinyint",
			output: Field{Type: Bool, Length: -1},
		},
		{
			input:  "smallint",
			output: Field{Type: Int16, Length: -1},
		},
		{
			input:  "mediumint",
			output: Field{Type: Int16, Length: -1},
		},
		{
			input:  "int",
			output: Field{Type: Int32, Length: -1},
		},
		{
			input:  "bigint",
			output: Field{Type: Int32, Length: -1},
		},
		{
			input:  "decimal(12, 4)",
			output: Field{Type: Float, Length: 8},
		},
		{
			input:  "float(12, 5)",
			output: Field{Type: Float, Length: 7},
		},
		{
			input:  "double(20,5)",
			output: Field{Type: Float, Length: 15},
		},
		{
			input:  "blob",
			output: Field{Type: Blob, Length: -1},
		},
		{
			input:  "tinyblob",
			output: Field{Type: Blob, Length: -1},
		},
		{
			input:  "mediumblob",
			output: Field{Type: Blob, Length: -1},
		},
		{
			input:  "longblob",
			output: Field{Type: Blob, Length: -1},
		},
		{
			input:  "text",
			output: Field{Type: Text, Length: -1},
		},
		{
			input:  "tinytext",
			output: Field{Type: Text, Length: -1},
		},
		{
			input:  "mediumtext",
			output: Field{Type: Text, Length: -1},
		},
		{
			input:  "longtext",
			output: Field{Type: Text, Length: -1},
		},
		{
			input:  "json",
			output: Field{Type: Json, Length: -1},
		},
		{
			input:  "enum(test, this, data)",
			output: Field{Type: Enum, Length: -1, Enum: []string{"test", "this", "data"}},
		},
	}

	for _, scenario := range scenarios {
		output := MySQL{}.MapField(scenario.input)

		if !reflect.DeepEqual(output, scenario.output) {
			t.Errorf("Invalid output for %s, out: %+v", scenario.input, output)
		}
	}
}
