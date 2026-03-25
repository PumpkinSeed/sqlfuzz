package mysql

import (
	"fmt"
	"strconv"
	"strings"
)

func length(field string, t string) ([]int16, error) {
	field = strings.ToLower(field)
	t = strings.ToLower(t)
	if !(strings.Contains(field, "(") && strings.Contains(field, ")")) {
		return nil, nil
	}
	str := strings.ReplaceAll(field, t, "")
	str = strings.ReplaceAll(str, "(", "")
	str = strings.ReplaceAll(str, ")", "")
	str = strings.ReplaceAll(str, " ", "")
	if strings.Contains(str, ",") {
		p := strings.Split(str, ",")
		var result []int16
		for _, v := range p {
			data, err := strconv.Atoi(v)
			if err != nil {
				return nil, fmt.Errorf("length: failed to parse %q: %w", v, err)
			}
			result = append(result, int16(data))
		}
		return result, nil
	}
	v, err := strconv.Atoi(str)
	if err != nil {
		return nil, fmt.Errorf("length: failed to parse %q: %w", str, err)
	}
	return []int16{int16(v)}, nil
}
