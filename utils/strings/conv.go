package strings

import (
	"strconv"
)

func StringToInt(intStr string) int {
	v, err := strconv.Atoi(intStr)
	if err != nil {
		return 0
	}
	return v
}
