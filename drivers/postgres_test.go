package drivers

import (
	"encoding/json"
	"fmt"
	"reflect"
	"testing"

	_ "github.com/lib/pq"
	"github.com/volatiletech/null"
)

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

func TestPostgres_MultiDescribe(t *testing.T) {
	db, err := getPostgresConnection()
	pgDriver := Postgres{}

	if err != nil {
		t.Errorf("error getting postgres connection : %v", err.Error())
	}
	testCase, err := pgDriver.GetTestCase("multi")
	if err != nil {
		t.Errorf("error getting multi test case : %v", err.Error())
	}
	err = pgDriver.TestTable(db, "multi", "")
	if err != nil {
		t.Errorf("error initialising multi test case : %v", err.Error())
	}
	tables := testCase.TableCreationOrder
	tableFieldsMap, insertionOrder, err := Postgres{}.MultiDescribe(tables, db)
	if err != nil {
		t.Errorf("error descriving tables %v. Error : %v", tables, err)
	}
	if len(tableFieldsMap) == 0 || len(insertionOrder) != len(tableFieldsMap) || len(insertionOrder) != len(tables) {
		t.Errorf("error receiving required fields count. input len %v described fields len %v insertion order length %v", len(tables), len(tableFieldsMap), len(insertionOrder))
	}
	tableFieldMapStr, err := json.Marshal(tableFieldsMap)
	if err != nil {
		t.Error(err)
	}
	fmt.Println(string(tableFieldMapStr))
}
