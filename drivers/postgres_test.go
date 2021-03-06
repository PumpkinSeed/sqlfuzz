package drivers

import (
	_ "github.com/lib/pq"
	"github.com/volatiletech/null"
	"reflect"
	"testing"
)

/*
Added and commented for the purpose of testing as the values will vary in each local setup.
*/
//func TestPostgres_DescribeFields(t *testing.T) {
//
//	connectionString := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable ", "127.0.0.1", "5432", "postgres", "password", "fuzzpostgres")
//	db, err := sql.Open("postgres", connectionString)
//	if err != nil {
//		t.Errorf("Error opening postgres connectoin : %s", err.Error())
//		return
//	}
//	driver := Postgres{}
//	tables, err := driver.ShowTables(db)
//	if err != nil {
//		t.Errorf("Error describing table : %s", err.Error())
//		return
//	}
//	for _, table := range tables {
//		fmt.Println(table)
//	}
//}

func TestPostgres_MapField(t *testing.T) {
	var scenarios = []struct {
		input  FieldDescriptor
		output Field
	}{
		{
			input: FieldDescriptor{
				Type: "bigint",
			},
			output: Field{Type: Int32, Length: -1},
		},
		{
			input:  FieldDescriptor{Type: "bit varying", Length: null.IntFrom(10)},
			output: Field{Type: BinaryString, Length: 10},
		},
		{
			input:  FieldDescriptor{Type: "character varying", Length: null.IntFrom(21)},
			output: Field{Type: String, Length: int16(21)},
		},
		{
			input:  FieldDescriptor{Type: "json"},
			output: Field{Type: Json, Length: -1},
		},
	}

	for _, scenario := range scenarios {
		output := Postgres{}.MapField(scenario.input)
		if !reflect.DeepEqual(output, scenario.output) {
			t.Errorf("Invalid output for %s, out: %+v", scenario.input.Field, output)
		}
	}
}
