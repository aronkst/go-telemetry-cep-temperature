package utils

import "strconv"

func IsNumber(value string) bool {
	num, err := strconv.Atoi(value)
	return err == nil && num > 0
}
