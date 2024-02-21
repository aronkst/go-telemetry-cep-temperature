package service_test

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/aronkst/go-telemetry-cep-temperature/internal/input_server/model"
	"github.com/aronkst/go-telemetry-cep-temperature/internal/input_server/service"
	temperatureServerModel "github.com/aronkst/go-telemetry-cep-temperature/internal/temperature_server/model"
)

type MockTemperatureRepository struct {
	Temperature *temperatureServerModel.Temperature
	Err         error
}

func (m *MockTemperatureRepository) GetTemperature(*model.Zipcode, context.Context) (*temperatureServerModel.Temperature, error) {
	return m.Temperature, m.Err
}

func TestInputService_Success(t *testing.T) {
	mockTemperatureRepo := &MockTemperatureRepository{Temperature: &temperatureServerModel.Temperature{City: "Cidade", Celsius: 30.0, Fahrenheit: 86.0, Kelvin: 303.15}}

	service := service.NewInputService(mockTemperatureRepo)

	zipcode := &model.Zipcode{
		Cep: "12345678",
	}

	temperature, err := service.GetTemperatureByCep(zipcode, context.Background())
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	expected := &temperatureServerModel.Temperature{
		City:       "Cidade",
		Celsius:    30.0,
		Fahrenheit: 86.0,
		Kelvin:     303.15,
	}

	if temperature.City != expected.City {
		t.Errorf("Expected City %v, got %v", expected.City, temperature.City)
	}

	if temperature.Celsius != expected.Celsius {
		t.Errorf("Expected Celsius %v, got %v", expected.Celsius, temperature.Celsius)
	}

	if temperature.Fahrenheit != expected.Fahrenheit {
		t.Errorf("Expected Fahrenheit %v, got %v", expected.Fahrenheit, temperature.Fahrenheit)
	}

	if temperature.Kelvin != expected.Kelvin {
		t.Errorf("Expected Kelvin %v, got %v", expected.Kelvin, temperature.Kelvin)
	}
}

func TestInputService_Error(t *testing.T) {
	expectedErrorMsg := "invalid zipcode"

	mockTemperatureRepo := &MockTemperatureRepository{Err: fmt.Errorf(expectedErrorMsg)}

	service := service.NewInputService(mockTemperatureRepo)

	zipcode := &model.Zipcode{
		Cep: "0",
	}

	_, err := service.GetTemperatureByCep(zipcode, context.Background())
	if err == nil {
		t.Fatalf("Expected an error but got nil")
	}

	if !strings.Contains(err.Error(), expectedErrorMsg) {
		t.Errorf("Error message does not match expected. \nExpected to contain: %s\nGot: %s", expectedErrorMsg, err.Error())
	}
}
