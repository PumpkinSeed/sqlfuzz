package mysql

import (
	"strconv"
	"strings"
)

func length(field string, t string) []int16 {
	field = strings.ToLower(field)
	t = strings.ToLower(t)
	if !(strings.Contains(field, "(") && strings.Contains(field, ")")) {
		return nil
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
				panic(err)
			}
			result = append(result, int16(data))
		}
		return result
	}
	v, err := strconv.Atoi(str)
	if err != nil {
		panic(err)
	}
	return []int16{int16(v)}
}
