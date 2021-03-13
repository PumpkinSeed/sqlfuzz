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
				Type: "varchar(12)",
			},
			output: Field{Type: String, Length: 12},
		},
		{
			input: FieldDescriptor{
				Type: "char(100)",
			},
			output: Field{Type: String, Length: 100},
		},
		{
			input: FieldDescriptor{
				Type: "varbinary(100)",
			},
			output: Field{Type: String, Length: 100},
		},
		{
			input: FieldDescriptor{
				Type: "binary(100)",
			},
			output: Field{Type: String, Length: 100},
		},
		{
			input: FieldDescriptor{
				Type: "tinyint",
			},
			output: Field{Type: Bool, Length: -1},
		},
		{
			input: FieldDescriptor{
				Type: "smallint",
			},
			output: Field{Type: Int16, Length: -1},
		},
		{
			input: FieldDescriptor{
				Type: "mediumint",
			},
			output: Field{Type: Int16, Length: -1},
		},
		{
			input: FieldDescriptor{
				Type: "int",
			},
			output: Field{Type: Int32, Length: -1},
		},
		{
			input: FieldDescriptor{
				Type: "bigint",
			},
			output: Field{Type: Int32, Length: -1},
		},
		{
			input: FieldDescriptor{
				Type: "decimal(12, 4)",
			},
			output: Field{Type: Float, Length: 8},
		},
		{
			input: FieldDescriptor{
				Type: "float(12, 5)",
			},
			output: Field{Type: Float, Length: 7},
		},
		{
			input: FieldDescriptor{
				Type: "double(20,5)",
			},
			output: Field{Type: Float, Length: 15},
		},
		{
			input: FieldDescriptor{
				Type: "blob",
			},
			output: Field{Type: Blob, Length: -1},
		},
		{
			input: FieldDescriptor{
				Type: "tinyblob",
			},
			output: Field{Type: Blob, Length: -1},
		},
		{
			input: FieldDescriptor{
				Type: "mediumblob",
			},
			output: Field{Type: Blob, Length: -1},
		},
		{
			input: FieldDescriptor{
				Type: "longblob",
			},
			output: Field{Type: Blob, Length: -1},
		},
		{
			input: FieldDescriptor{
				Type: "text",
			},
			output: Field{Type: Text, Length: -1},
		},
		{
			input: FieldDescriptor{
				Type: "tinytext",
			},
			output: Field{Type: Text, Length: -1},
		},
		{
			input: FieldDescriptor{
				Type: "mediumtext",
			},
			output: Field{Type: Text, Length: -1},
		},
		{
			input: FieldDescriptor{
				Type: "longtext",
			},
			output: Field{Type: Text, Length: -1},
		},
		{
			input: FieldDescriptor{
				Type: "json",
			},
			output: Field{Type: Json, Length: -1},
		},
		{
			input: FieldDescriptor{
				Type: "enum(test, this, data)",
			},
			output: Field{Type: Enum, Length: -1, Enum: []string{"test", "this", "data"}},
		},
	}

	for _, scenario := range scenarios {
		output := MySQL{}.MapField(scenario.input)

		if !reflect.DeepEqual(output.Type, scenario.output.Type) {
			t.Errorf("Invalid output for %s, out: %+v scenario out: %+v", scenario.input.Field, output, scenario.output)
		}
	}
}
