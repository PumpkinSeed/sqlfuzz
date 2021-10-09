package drivers

import (
	"reflect"
	"testing"

	"github.com/PumpkinSeed/sqlfuzz/drivers/types"
)

func TestMapField(t *testing.T) {
	var scenarios = []struct {
		input  types.FieldDescriptor
		output types.Field
	}{
		{
			input: types.FieldDescriptor{
				Type: "varchar(12)",
			},
			output: types.Field{Type: types.String, Length: 12},
		},
		{
			input: types.FieldDescriptor{
				Type: "char(100)",
			},
			output: types.Field{Type: types.String, Length: 100},
		},
		{
			input: types.FieldDescriptor{
				Type: "varbinary(100)",
			},
			output: types.Field{Type: types.String, Length: 100},
		},
		{
			input: types.FieldDescriptor{
				Type: "binary(100)",
			},
			output: types.Field{Type: types.String, Length: 100},
		},
		{
			input: types.FieldDescriptor{
				Type: "tinyint",
			},
			output: types.Field{Type: types.Bool, Length: -1},
		},
		{
			input: types.FieldDescriptor{
				Type: "smallint",
			},
			output: types.Field{Type: types.Int16, Length: -1},
		},
		{
			input: types.FieldDescriptor{
				Type: "mediumint",
			},
			output: types.Field{Type: types.Int16, Length: -1},
		},
		{
			input: types.FieldDescriptor{
				Type: "int",
			},
			output: types.Field{Type: types.Int32, Length: -1},
		},
		{
			input: types.FieldDescriptor{
				Type: "bigint",
			},
			output: types.Field{Type: types.Int32, Length: -1},
		},
		{
			input: types.FieldDescriptor{
				Type: "decimal(12, 4)",
			},
			output: types.Field{Type: types.Float, Length: 8},
		},
		{
			input: types.FieldDescriptor{
				Type: "float(12, 5)",
			},
			output: types.Field{Type: types.Float, Length: 7},
		},
		{
			input: types.FieldDescriptor{
				Type: "double(20,5)",
			},
			output: types.Field{Type: types.Float, Length: 15},
		},
		{
			input: types.FieldDescriptor{
				Type: "blob",
			},
			output: types.Field{Type: types.Blob, Length: -1},
		},
		{
			input: types.FieldDescriptor{
				Type: "tinyblob",
			},
			output: types.Field{Type: types.Blob, Length: -1},
		},
		{
			input: types.FieldDescriptor{
				Type: "mediumblob",
			},
			output: types.Field{Type: types.Blob, Length: -1},
		},
		{
			input: types.FieldDescriptor{
				Type: "longblob",
			},
			output: types.Field{Type: types.Blob, Length: -1},
		},
		{
			input: types.FieldDescriptor{
				Type: "text",
			},
			output: types.Field{Type: types.Text, Length: -1},
		},
		{
			input: types.FieldDescriptor{
				Type: "tinytext",
			},
			output: types.Field{Type: types.Text, Length: -1},
		},
		{
			input: types.FieldDescriptor{
				Type: "mediumtext",
			},
			output: types.Field{Type: types.Text, Length: -1},
		},
		{
			input: types.FieldDescriptor{
				Type: "longtext",
			},
			output: types.Field{Type: types.Text, Length: -1},
		},
		{
			input: types.FieldDescriptor{
				Type: "json",
			},
			output: types.Field{Type: types.Json, Length: -1},
		},
		{
			input: types.FieldDescriptor{
				Type: "enum(test, this, data)",
			},
			output: types.Field{Type: types.Enum, Length: -1, Enum: []string{"test", "this", "data"}},
		},
	}

	for _, scenario := range scenarios {
		output := MySQL{}.MapField(scenario.input)

		if !reflect.DeepEqual(output.Type, scenario.output.Type) {
			t.Errorf("Invalid output for %s, out: %+v scenario out: %+v", scenario.input.Field, output, scenario.output)
		}
	}
}
