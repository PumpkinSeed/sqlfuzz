package drivers

import (
	"reflect"
	"testing"
)

func TestMapField(t *testing.T) {
	var scenarios = []struct {
		input  FieldDescriptor
		output Field
	}{
		{
			input: FieldDescriptor{
				Field: "varchar(12)",
			},
			output: Field{Type: String, Length: 12},
		},
		{
			input: FieldDescriptor{
				Field: "char(100)",
			},
			output: Field{Type: String, Length: 100},
		},
		{
			input: FieldDescriptor{
				Field: "varbinary(100)",
			},
			output: Field{Type: String, Length: 100},
		},
		{
			input: FieldDescriptor{
				Field: "binary(100)",
			},
			output: Field{Type: String, Length: 100},
		},
		{
			input: FieldDescriptor{
				Field: "tinyint",
			},
			output: Field{Type: Bool, Length: -1},
		},
		{
			input: FieldDescriptor{
				Field: "smallint",
			},
			output: Field{Type: Int16, Length: -1},
		},
		{
			input: FieldDescriptor{
				Field: "mediumint",
			},
			output: Field{Type: Int16, Length: -1},
		},
		{
			input: FieldDescriptor{
				Field: "int",
			},
			output: Field{Type: Int32, Length: -1},
		},
		{
			input: FieldDescriptor{
				Field: "bigint",
			},
			output: Field{Type: Int32, Length: -1},
		},
		{
			input: FieldDescriptor{
				Field: "decimal(12, 4)",
			},
			output: Field{Type: Float, Length: 8},
		},
		{
			input: FieldDescriptor{
				Field: "float(12, 5)",
			},
			output: Field{Type: Float, Length: 7},
		},
		{
			input: FieldDescriptor{
				Field: "double(20,5)",
			},
			output: Field{Type: Float, Length: 15},
		},
		{
			input: FieldDescriptor{
				Field: "blob",
			},
			output: Field{Type: Blob, Length: -1},
		},
		{
			input: FieldDescriptor{
				Field: "tinyblob",
			},
			output: Field{Type: Blob, Length: -1},
		},
		{
			input: FieldDescriptor{
				Field: "mediumblob",
			},
			output: Field{Type: Blob, Length: -1},
		},
		{
			input: FieldDescriptor{
				Field: "longblob",
			},
			output: Field{Type: Blob, Length: -1},
		},
		{
			input: FieldDescriptor{
				Field: "text",
			},
			output: Field{Type: Text, Length: -1},
		},
		{
			input: FieldDescriptor{
				Field: "tinytext",
			},
			output: Field{Type: Text, Length: -1},
		},
		{
			input: FieldDescriptor{
				Field: "mediumtext",
			},
			output: Field{Type: Text, Length: -1},
		},
		{
			input: FieldDescriptor{
				Field: "longtext",
			},
			output: Field{Type: Text, Length: -1},
		},
		{
			input: FieldDescriptor{
				Field: "json",
			},
			output: Field{Type: Json, Length: -1},
		},
		{
			input: FieldDescriptor{
				Field: "enum(test, this, data)",
			},
			output: Field{Type: Enum, Length: -1, Enum: []string{"test", "this", "data"}},
		},
	}

	for _, scenario := range scenarios {
		output := MySQL{}.MapField(scenario.input)

		if !reflect.DeepEqual(output, scenario.output) {
			t.Errorf("Invalid output for %s, out: %+v", scenario.input.Field, output)
		}
	}
}
