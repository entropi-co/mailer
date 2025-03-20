package snowflakes

import "strconv"

func ValueFromString(str string) (uint64, error) {
	return strconv.ParseUint(str, 10, 64)
}
