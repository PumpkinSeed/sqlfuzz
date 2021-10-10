package postgres

import (
	"encoding/json"
	"reflect"
	"testing"

	"github.com/PumpkinSeed/sqlfuzz/drivers/types"
	_ "github.com/lib/pq"
	"github.com/volatiletech/null"
)

func TestPostgres_MapField(t *testing.T) {
	var scenarios = []struct {
		input  types.FieldDescriptor
		output types.Field
	}{
		{
			input: types.FieldDescriptor{
				Type: "bigint",
			},
			output: types.Field{Type: types.Int32, Length: -1},
		},
		{
			input:  types.FieldDescriptor{Type: "bit varying", Length: null.IntFrom(10)},
			output: types.Field{Type: types.BinaryString, Length: 10},
		},
		{
			input:  types.FieldDescriptor{Type: "character varying", Length: null.IntFrom(21)},
			output: types.Field{Type: types.String, Length: int16(21)},
		},
		{
			input:  types.FieldDescriptor{Type: "json"},
			output: types.Field{Type: types.Json, Length: -1},
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
		t.Errorf("error initializing multi test case : %v", err.Error())
	}
	tables := testCase.TableCreationOrder
	tableFieldsMap, insertionOrder, err := Postgres{}.MultiDescribe(tables, db)
	if err != nil {
		t.Errorf("error descriving tables %v. Error : %v", tables, err)
	}
	if len(tableFieldsMap) == 0 || len(insertionOrder) != len(tableFieldsMap) || len(insertionOrder) != len(tables) {
		t.Errorf(
			"error receiving required fields count. input len %v described fields len %v insertion order length %v",
			len(tables), len(tableFieldsMap), len(insertionOrder),
		)
	}
	tableFieldMapStr, err := json.Marshal(tableFieldsMap)
	if err != nil {
		t.Error(err)
	}
	t.Log(string(tableFieldMapStr))
}
