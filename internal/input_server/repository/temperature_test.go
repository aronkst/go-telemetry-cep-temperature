package repository_test

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/aronkst/go-telemetry-cep-temperature/internal/input_server/model"
	"github.com/aronkst/go-telemetry-cep-temperature/internal/input_server/repository"
	temperatureServerModel "github.com/aronkst/go-telemetry-cep-temperature/internal/temperature_server/model"
)

func TestTemperatureRepository_Success(t *testing.T) {
	t.Setenv("TEST", "true")

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		responseBody := `{"city":"Cidade","temp_C":30,"temp_F":86,"temp_K":303.15}`
		w.Write([]byte(responseBody))
	}))
	defer server.Close()

	repo := repository.NewTemperatureRepository(server.URL)

	zipcode := &model.Zipcode{
		Cep: "12345678",
	}

	address, err := repo.GetTemperature(zipcode, context.Background())
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	expected := &temperatureServerModel.Temperature{
		City:       "Cidade",
		Celsius:    30.0,
		Fahrenheit: 86.0,
		Kelvin:     303.15,
	}

	if address.City != expected.City {
		t.Errorf("City mismatch: expected %v, got %v", expected.City, address.City)
	}

	if address.Celsius != expected.Celsius {
		t.Errorf("Celsius mismatch: expected %v, got %v", expected.Celsius, address.Celsius)
	}

	if address.Fahrenheit != expected.Fahrenheit {
		t.Errorf("Fahrenheit mismatch: expected %v, got %v", expected.Fahrenheit, address.Fahrenheit)
	}

	if address.Kelvin != expected.Kelvin {
		t.Errorf("Kelvin mismatch: expected %v, got %v", expected.Kelvin, address.Kelvin)
	}
}

func TestTemperatureRepository_InvalidCep(t *testing.T) {
	t.Setenv("TEST", "true")

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		responseBody := `{}`
		w.Write([]byte(responseBody))
	}))
	defer server.Close()

	repo := repository.NewTemperatureRepository(server.URL)

	zipcode := &model.Zipcode{
		Cep: "0",
	}

	_, err := repo.GetTemperature(zipcode, context.Background())
	if err == nil {
		t.Fatalf("Expected an error but got nil")
	}

	expectedErrorMsg := "invalid zipcode"
	if !strings.Contains(err.Error(), expectedErrorMsg) {
		t.Errorf("Error message does not match expected. \nExpected to contain: %s\nGot: %s", expectedErrorMsg, err.Error())
	}
}

func TestTemperatureRepository_ErrorHttp(t *testing.T) {
	t.Setenv("TEST", "true")

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "", http.StatusMovedPermanently)
	}))
	defer server.Close()

	repo := repository.NewTemperatureRepository(server.URL)

	zipcode := &model.Zipcode{
		Cep: "12345678",
	}

	_, err := repo.GetTemperature(zipcode, context.Background())
	if err == nil {
		t.Fatalf("Expected an error but got nil")
	}

	expectedErrorMsg := "error when searching for temperature by cep"
	if !strings.Contains(err.Error(), expectedErrorMsg) {
		t.Errorf("Error message does not match expected. \nExpected to contain: %s\nGot: %s", expectedErrorMsg, err.Error())
	}
}

func TestTemperatureRepository_NotFindZipcode(t *testing.T) {
	t.Setenv("TEST", "true")

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Not Found", http.StatusNotFound)
	}))
	defer server.Close()

	repo := repository.NewTemperatureRepository(server.URL)

	zipcode := &model.Zipcode{
		Cep: "12345678",
	}

	_, err := repo.GetTemperature(zipcode, context.Background())
	if err == nil {
		t.Fatalf("Expected an error but got nil")
	}

	expectedErrorMsg := "can not find zipcode"
	if !strings.Contains(err.Error(), expectedErrorMsg) {
		t.Errorf("Error message does not match expected. \nExpected to contain: %s\nGot: %s", expectedErrorMsg, err.Error())
	}
}

func TestTemperatureRepository_NotStatusOK(t *testing.T) {
	t.Setenv("TEST", "true")

	statusServerError := http.StatusInternalServerError

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Internal Server Error", statusServerError)
	}))
	defer server.Close()

	repo := repository.NewTemperatureRepository(server.URL)

	zipcode := &model.Zipcode{
		Cep: "12345678",
	}

	_, err := repo.GetTemperature(zipcode, context.Background())
	if err == nil {
		t.Fatalf("Expected an error but got nil")
	}

	expectedErrorMsg := fmt.Sprintf("temperature by cep api returned status %d", statusServerError)
	if !strings.Contains(err.Error(), expectedErrorMsg) {
		t.Errorf("Error message does not match expected. \nExpected to contain: %s\nGot: %s", expectedErrorMsg, err.Error())
	}
}

func TestTemperatureRepository_ErrorJsonDecoder(t *testing.T) {
	t.Setenv("TEST", "true")

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		responseBody := `error`
		w.Write([]byte(responseBody))
	}))
	defer server.Close()

	repo := repository.NewTemperatureRepository(server.URL)

	zipcode := &model.Zipcode{
		Cep: "12345678",
	}

	_, err := repo.GetTemperature(zipcode, context.Background())
	if err == nil {
		t.Fatalf("Expected an error but got nil")
	}

	expectedErrorMsg := "error parsing json"
	if !strings.Contains(err.Error(), expectedErrorMsg) {
		t.Errorf("Error message does not match expected. \nExpected to contain: %s\nGot: %s", expectedErrorMsg, err.Error())
	}
}
