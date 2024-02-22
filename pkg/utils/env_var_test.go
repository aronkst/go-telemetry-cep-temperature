package utils_test

import (
	"os"
	"testing"

	"github.com/aronkst/go-telemetry-cep-temperature/pkg/utils"
)

func TestGetEnvOrDefaultSet(t *testing.T) {
	const envKey = "TEST_ENV_VAR"
	const expectedValue = "testValue"

	t.Setenv(envKey, expectedValue)

	if value := utils.GetEnvOrDefault(envKey, "defaultValue"); value != expectedValue {
		t.Errorf("Expected %s, got %s", expectedValue, value)
	}

	os.Unsetenv(envKey)
}

func TestGetEnvOrDefaultNotSet(t *testing.T) {
	const envKey = "NON_EXISTENT_ENV_VAR"
	const defaultValue = "defaultValue"

	if value := utils.GetEnvOrDefault(envKey, defaultValue); value != defaultValue {
		t.Errorf("Expected %s, got %s", defaultValue, value)
	}
}
