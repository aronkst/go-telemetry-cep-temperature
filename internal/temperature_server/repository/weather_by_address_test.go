package repository_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/aronkst/go-telemetry-cep-temperature/internal/model"
	"github.com/aronkst/go-telemetry-cep-temperature/internal/repository"
)

func TestWeatherByAddressRepository_Success(t *testing.T) {
	t.Setenv("TEST", "true")

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		responseBody := `{"current_condition":[{"temp_C":"30"}]}`
		w.Write([]byte(responseBody))
	}))
	defer server.Close()

	repo := repository.NewWeatherByAddressRepository(server.URL)

	address := &model.Address{
		PostalCode: "12345-678",
		Street:     "Rua Exemplo",
		Complement: "",
		District:   "Bairro",
		City:       "Cidade",
		State:      "Estado",
	}

	temperature, err := repo.GetWeather(address)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	expected := &model.Weather{
		Temperature: 30.0,
	}

	if temperature.Temperature != expected.Temperature {
		t.Errorf("Latitude mismatch: expected %v, got %v", expected.Temperature, temperature.Temperature)
	}
}

func TestWeatherByAddressRepository_ErrorHttp(t *testing.T) {
	t.Setenv("TEST", "true")

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "", http.StatusMovedPermanently)
	}))
	defer server.Close()

	repo := repository.NewWeatherByAddressRepository(server.URL)

	address := &model.Address{
		PostalCode: "12345-678",
		Street:     "Rua Exemplo",
		Complement: "",
		District:   "Bairro",
		City:       "Cidade",
		State:      "Estado",
	}

	_, err := repo.GetWeather(address)
	if err == nil {
		t.Fatalf("Expected an error but got nil")
	}

	expectedErrorMsg := "error when searching for weather forecast"
	if !strings.Contains(err.Error(), expectedErrorMsg) {
		t.Errorf("Error message does not match expected. \nExpected to contain: %s\nGot: %s", expectedErrorMsg, err.Error())
	}
}

func TestWeatherByAddressRepository_NotStatusOK(t *testing.T) {
	t.Setenv("TEST", "true")

	statusServerError := http.StatusInternalServerError

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Internal Server Error", statusServerError)
	}))
	defer server.Close()

	repo := repository.NewWeatherByAddressRepository(server.URL)

	address := &model.Address{
		PostalCode: "12345-678",
		Street:     "Rua Exemplo",
		Complement: "",
		District:   "Bairro",
		City:       "Cidade",
		State:      "Estado",
	}

	_, err := repo.GetWeather(address)
	if err == nil {
		t.Fatalf("Expected an error but got nil")
	}

	expectedErrorMsg := fmt.Sprintf("weather api returned status %d", statusServerError)
	if !strings.Contains(err.Error(), expectedErrorMsg) {
		t.Errorf("Error message does not match expected. \nExpected to contain: %s\nGot: %s", expectedErrorMsg, err.Error())
	}
}

func TestWeatherByAddressRepository_ErrorJsonDecoder(t *testing.T) {
	t.Setenv("TEST", "true")

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		responseBody := `error`
		w.Write([]byte(responseBody))
	}))
	defer server.Close()

	repo := repository.NewWeatherByAddressRepository(server.URL)

	address := &model.Address{
		PostalCode: "12345-678",
		Street:     "Rua Exemplo",
		Complement: "",
		District:   "Bairro",
		City:       "Cidade",
		State:      "Estado",
	}

	_, err := repo.GetWeather(address)
	if err == nil {
		t.Fatalf("Expected an error but got nil")
	}

	expectedErrorMsg := "error parsing json"
	if !strings.Contains(err.Error(), expectedErrorMsg) {
		t.Errorf("Error message does not match expected. \nExpected to contain: %s\nGot: %s", expectedErrorMsg, err.Error())
	}
}

func TestWeatherByAddressRepository_JsonBlank(t *testing.T) {
	t.Setenv("TEST", "true")

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		responseBody := `{"current_condition":[]}`
		w.Write([]byte(responseBody))
	}))
	defer server.Close()

	repo := repository.NewWeatherByAddressRepository(server.URL)

	address := &model.Address{
		PostalCode: "12345-678",
		Street:     "Rua Exemplo",
		Complement: "",
		District:   "Bairro",
		City:       "Cidade",
		State:      "Estado",
	}

	_, err := repo.GetWeather(address)
	if err == nil {
		t.Fatalf("Expected an error but got nil")
	}

	expectedErrorMsg := "temperature error"
	if !strings.Contains(err.Error(), expectedErrorMsg) {
		t.Errorf("Error message does not match expected. \nExpected to contain: %s\nGot: %s", expectedErrorMsg, err.Error())
	}
}
