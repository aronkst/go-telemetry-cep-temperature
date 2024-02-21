package utils_test

import (
	"testing"

	"github.com/aronkst/go-telemetry-cep-temperature/pkg/utils"
)

func TestStringToFloat64(t *testing.T) {
	tests := []struct {
		input string
		want  float64
	}{
		{"123.456", 123.456},
		{"-123.456", -123.456},
		{"0", 0},
		{"1e2", 100},
		{"abc", 0},
		{"", 0},
		{"123-456", 0},
	}

	for _, test := range tests {
		got := utils.StringToFloat64(test.input)
		if got != test.want {
			t.Errorf("StringToFloat64(%q) = %v; want %v", test.input, got, test.want)
		}
	}
}
