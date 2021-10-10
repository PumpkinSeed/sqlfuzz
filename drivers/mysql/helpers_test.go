package mysql

import (
	"reflect"
	"testing"
)

func TestLength(t *testing.T) {
	var scenarios = []struct {
		input  string
		t      string
		output []int16
	}{
		{
			"DECIMAL(12, 4)",
			"decimal",
			[]int16{12, 4},
		},
		{
			"FLOAT(12)",
			"float",
			[]int16{12},
		},
		{
			"DOUBLE(23)",
			"DOUBLE",
			[]int16{23},
		},
	}

	for _, scenario := range scenarios {
		out := length(scenario.input, scenario.t)
		if !reflect.DeepEqual(scenario.output, out) {
			t.Errorf("Output doesn't match with the scenario: %v, out: %v", scenario.output, out)
		}
	}
}
