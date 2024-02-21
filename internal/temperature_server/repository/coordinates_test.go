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

func TestCoordinatesRepository_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		responseBody := `[{"lat":"123","lon":"321"}]`
		w.Write([]byte(responseBody))
	}))
	defer server.Close()

	repo := repository.NewCoordinatesRepository(server.URL)

	address := &model.Address{
		PostalCode: "12345-678",
		Street:     "Rua Exemplo",
		Complement: "",
		District:   "Bairro",
		City:       "Cidade",
		State:      "Estado",
	}

	coordinates, err := repo.GetCoordinates(address)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	expected := &model.Coordinates{
		Latitude:  "123",
		Longitude: "321",
	}

	if coordinates.Latitude != expected.Latitude {
		t.Errorf("Latitude mismatch: expected %v, got %v", expected.Latitude, coordinates.Latitude)
	}

	if coordinates.Longitude != expected.Longitude {
		t.Errorf("Longitude mismatch: expected %v, got %v", expected.Longitude, coordinates.Longitude)
	}
}

func TestCoordinatesRepository_ErrorHttp(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "", http.StatusMovedPermanently)
	}))
	defer server.Close()

	repo := repository.NewCoordinatesRepository(server.URL)

	address := &model.Address{
		PostalCode: "12345-678",
		Street:     "Rua Exemplo",
		Complement: "",
		District:   "Bairro",
		City:       "Cidade",
		State:      "Estado",
	}

	_, err := repo.GetCoordinates(address)
	if err == nil {
		t.Fatalf("Expected an error but got nil")
	}

	expectedErrorMsg := "error when searching for coordinates for the address"
	if !strings.Contains(err.Error(), expectedErrorMsg) {
		t.Errorf("Error message does not match expected. \nExpected to contain: %s\nGot: %s", expectedErrorMsg, err.Error())
	}
}

func TestCoordinatesRepository_NotStatusOK(t *testing.T) {
	statusServerError := http.StatusInternalServerError

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Internal Server Error", statusServerError)
	}))
	defer server.Close()

	repo := repository.NewCoordinatesRepository(server.URL)

	address := &model.Address{
		PostalCode: "12345-678",
		Street:     "Rua Exemplo",
		Complement: "",
		District:   "Bairro",
		City:       "Cidade",
		State:      "Estado",
	}

	_, err := repo.GetCoordinates(address)
	if err == nil {
		t.Fatalf("Expected an error but got nil")
	}

	expectedErrorMsg := fmt.Sprintf("coordinates api returned status %d", statusServerError)
	if !strings.Contains(err.Error(), expectedErrorMsg) {
		t.Errorf("Error message does not match expected. \nExpected to contain: %s\nGot: %s", expectedErrorMsg, err.Error())
	}
}

func TestCoordinatesRepository_ErrorJsonDecoder(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		responseBody := `error`
		w.Write([]byte(responseBody))
	}))
	defer server.Close()

	repo := repository.NewCoordinatesRepository(server.URL)

	address := &model.Address{
		PostalCode: "12345-678",
		Street:     "Rua Exemplo",
		Complement: "",
		District:   "Bairro",
		City:       "Cidade",
		State:      "Estado",
	}

	_, err := repo.GetCoordinates(address)
	if err == nil {
		t.Fatalf("Expected an error but got nil")
	}

	expectedErrorMsg := "error decoding coordinates api response"
	if !strings.Contains(err.Error(), expectedErrorMsg) {
		t.Errorf("Error message does not match expected. \nExpected to contain: %s\nGot: %s", expectedErrorMsg, err.Error())
	}
}

func TestCoordinatesRepository_JsonBlank(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		responseBody := `[]`
		w.Write([]byte(responseBody))
	}))
	defer server.Close()

	repo := repository.NewCoordinatesRepository(server.URL)

	address := &model.Address{
		PostalCode: "12345-678",
		Street:     "Rua Exemplo",
		Complement: "",
		District:   "Bairro",
		City:       "Cidade",
		State:      "Estado",
	}

	_, err := repo.GetCoordinates(address)
	if err == nil {
		t.Fatalf("Expected an error but got nil")
	}

	expectedErrorMsg := "no coordinates found for the address"
	if !strings.Contains(err.Error(), expectedErrorMsg) {
		t.Errorf("Error message does not match expected. \nExpected to contain: %s\nGot: %s", expectedErrorMsg, err.Error())
	}
}
