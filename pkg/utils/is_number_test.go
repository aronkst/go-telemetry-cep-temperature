package utils_test

import (
	"testing"

	"github.com/aronkst/go-telemetry-cep-temperature/pkg/utils"
)

func TestIsNumber(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"123", true},
		{"-123", false},
		{"abc", false},
		{"123.5", false},
		{" ", false},
		{"1", true},
		{"0", false},
	}

	for _, test := range tests {
		if result := utils.IsNumber(test.input); result != test.expected {
			t.Errorf("IsNumber(%q) = %v; want %v", test.input, result, test.expected)
		}
	}
}
