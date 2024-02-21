package utils

import "math"

func CelsiusToFahrenheit(celsius float64) float64 {
	result := (celsius * 1.8) + 32
	return math.Round(result*100) / 100
}

func CelsiusToKelvin(celsius float64) float64 {
	result := celsius + 273.15
	return math.Round(result*100) / 100
}
