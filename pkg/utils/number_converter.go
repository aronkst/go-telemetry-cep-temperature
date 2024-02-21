package utils

import (
	"strconv"
)

func StringToFloat64(value string) float64 {
	number, _ := strconv.ParseFloat(value, 64)
	return number
}
