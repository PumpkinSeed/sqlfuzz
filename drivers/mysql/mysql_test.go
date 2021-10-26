package mysql

import (
	"database/sql"
	"reflect"
	"testing"

	"github.com/PumpkinSeed/sqlfuzz/drivers/types"
	_ "github.com/go-sql-driver/mysql"
)

func TestDescribe(t *testing.T) {
	// Describe(table string, db *sql.DB)
	db, err := sql.Open("mysql", "test:test@tcp(localhost:3306)/test")
	if err != nil {
		t.Error(err)
	}

	m := MySQL{}
	descriptors, err := m.Describe("t_product", db)
	if err != nil {
		t.Error(err)
	}

	if descriptors[0].Field != "id" {
		t.Error("First should be id")
	}
	if descriptors[1].Field != "name" {
		t.Error("Second should be name")
	}
	if descriptors[len(descriptors)-1].Field != "currency_id" {
		t.Error("Last should be currency_id")
	}

	t.Log(descriptors)
}

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
