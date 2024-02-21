package service_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aronkst/go-telemetry-cep-temperature/internal/temperature_server/model"
	"github.com/aronkst/go-telemetry-cep-temperature/internal/temperature_server/service"
)

type MockAddressRepository struct {
	Address *model.Address
	Err     error
}

func (m *MockAddressRepository) GetAddress(cep string) (*model.Address, error) {
	return m.Address, m.Err
}

type MockCoordinatesRepository struct {
	Coordinates *model.Coordinates
	Err         error
}

func (m *MockCoordinatesRepository) GetCoordinates(address *model.Address) (*model.Coordinates, error) {
	return m.Coordinates, m.Err
}

type MockWeatherByAddressRepository struct {
	Weather *model.Weather
	Err     error
}

func (m *MockWeatherByAddressRepository) GetWeather(address *model.Address) (*model.Weather, error) {
	return m.Weather, m.Err
}

type MockWeatherByCoordinatesRepository struct {
	Weather *model.Weather
	Err     error
}

func (m *MockWeatherByCoordinatesRepository) GetWeather(coordinates *model.Coordinates) (*model.Weather, error) {
	return m.Weather, m.Err
}

func TestWeatherService_Success(t *testing.T) {
	mockAddressRepo := &MockAddressRepository{Address: &model.Address{PostalCode: "12345-678", Street: "Rua Exemplo", Complement: "", District: "Bairro", City: "Cidade", State: "Estado"}}
	mockCoordinatesRepo := &MockCoordinatesRepository{Coordinates: &model.Coordinates{Latitude: "123", Longitude: "321"}}
	mockWeatherByAddressRepo := &MockWeatherByAddressRepository{Weather: &model.Weather{Temperature: 30}}
	mockWeatherByCoordinatesRepo := &MockWeatherByCoordinatesRepository{Weather: &model.Weather{Temperature: 30}}

	service := service.NewWeatherService(mockAddressRepo, mockCoordinatesRepo, mockWeatherByAddressRepo, mockWeatherByCoordinatesRepo)

	temperature, err := service.GetWeatherByCEP("12345678")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	expected := &model.Temperature{
		City:       "Cidade",
		Celsius:    30.0,
		Fahrenheit: 86.0,
		Kelvin:     303.15,
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

func TestWeatherService_AddressNotFound(t *testing.T) {
	expectedErrorMsg := "invalid zipcode"

	mockAddressRepo := &MockAddressRepository{Err: fmt.Errorf(expectedErrorMsg)}
	mockCoordinatesRepo := &MockCoordinatesRepository{}
	mockWeatherByAddressRepo := &MockWeatherByAddressRepository{}
	mockWeatherByCoordinatesRepo := &MockWeatherByCoordinatesRepository{}

	service := service.NewWeatherService(mockAddressRepo, mockCoordinatesRepo, mockWeatherByAddressRepo, mockWeatherByCoordinatesRepo)

	_, err := service.GetWeatherByCEP("12345678")
	if err == nil {
		t.Fatalf("Expected an error but got nil")
	}

	if !strings.Contains(err.Error(), expectedErrorMsg) {
		t.Errorf("Error message does not match expected. \nExpected to contain: %s\nGot: %s", expectedErrorMsg, err.Error())
	}
}

func TestWeatherService_ErrorWhenCoordinates(t *testing.T) {
	expectedErrorMsg := "can not find zipcode"

	mockAddressRepo := &MockAddressRepository{Address: &model.Address{PostalCode: "12345-678", Street: "Rua Exemplo", Complement: "", District: "Bairro", City: "Cidade", State: "Estado"}}
	mockCoordinatesRepo := &MockCoordinatesRepository{Coordinates: &model.Coordinates{Latitude: "123", Longitude: "321"}}
	mockWeatherByAddressRepo := &MockWeatherByAddressRepository{}
	mockWeatherByCoordinatesRepo := &MockWeatherByCoordinatesRepository{Err: fmt.Errorf(expectedErrorMsg)}

	service := service.NewWeatherService(mockAddressRepo, mockCoordinatesRepo, mockWeatherByAddressRepo, mockWeatherByCoordinatesRepo)

	_, err := service.GetWeatherByCEP("12345678")
	if err == nil {
		t.Fatalf("Expected an error but got nil")
	}

	if !strings.Contains(err.Error(), expectedErrorMsg) {
		t.Errorf("Error message does not match expected. \nExpected to contain: %s\nGot: %s", expectedErrorMsg, err.Error())
	}
}

func TestWeatherService_ErrorWhenNotCoordinates(t *testing.T) {
	expectedErrorMsg := "can not find zipcode"

	mockAddressRepo := &MockAddressRepository{Address: &model.Address{PostalCode: "12345-678", Street: "Rua Exemplo", Complement: "", District: "Bairro", City: "Cidade", State: "Estado"}}
	mockCoordinatesRepo := &MockCoordinatesRepository{Err: fmt.Errorf("coordinates not found")}
	mockWeatherByAddressRepo := &MockWeatherByAddressRepository{Err: fmt.Errorf(expectedErrorMsg)}
	mockWeatherByCoordinatesRepo := &MockWeatherByCoordinatesRepository{}

	service := service.NewWeatherService(mockAddressRepo, mockCoordinatesRepo, mockWeatherByAddressRepo, mockWeatherByCoordinatesRepo)

	_, err := service.GetWeatherByCEP("12345678")
	if err == nil {
		t.Fatalf("Expected an error but got nil")
	}

	if !strings.Contains(err.Error(), expectedErrorMsg) {
		t.Errorf("Error message does not match expected. \nExpected to contain: %s\nGot: %s", expectedErrorMsg, err.Error())
	}
}
