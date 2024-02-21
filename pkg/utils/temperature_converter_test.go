package utils_test

import (
	"testing"

	"github.com/aronkst/go-telemetry-cep-temperature/pkg/utils"
)

func TestCelsiusToFahrenheit(t *testing.T) {
	tests := []struct {
		celsius float64
		want    float64
	}{
		{0, 32},
		{100, 212},
		{-40, -40},
		{25, 77},
		{-273.15, -459.67},
	}

	for _, test := range tests {
		got := utils.CelsiusToFahrenheit(test.celsius)
		if got != test.want {
			t.Errorf("CelsiusToFahrenheit(%v) = %v; want %v", test.celsius, got, test.want)
		}
	}
}

func TestCelsiusToKelvin(t *testing.T) {
	tests := []struct {
		celsius float64
		want    float64
	}{
		{-273.15, 0},
		{0, 273.15},
		{100, 373.15},
		{-40, 233.15},
		{25, 298.15},
	}

	for _, test := range tests {
		got := utils.CelsiusToKelvin(test.celsius)
		if got != test.want {
			t.Errorf("CelsiusToKelvin(%v) = %v; want %v", test.celsius, got, test.want)
		}
	}
}
