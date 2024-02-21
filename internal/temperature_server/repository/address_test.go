package repository_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/aronkst/go-telemetry-cep-temperature/internal/temperature_server/model"
	"github.com/aronkst/go-telemetry-cep-temperature/internal/temperature_server/repository"
)

func TestAddressRepository_Success(t *testing.T) {
	t.Setenv("TEST", "true")

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		responseBody := `{"cep":"12345-678","logradouro":"Rua Exemplo","complemento":"","bairro":"Bairro","localidade":"Cidade","uf":"Estado"}`
		w.Write([]byte(responseBody))
	}))
	defer server.Close()

	repo := repository.NewAddressRepository(server.URL)

	address, err := repo.GetAddress("12345678")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	expected := &model.Address{
		PostalCode: "12345-678",
		Street:     "Rua Exemplo",
		Complement: "",
		District:   "Bairro",
		City:       "Cidade",
		State:      "Estado",
	}

	if address.PostalCode != expected.PostalCode {
		t.Errorf("PostalCode mismatch: expected %v, got %v", expected.PostalCode, address.PostalCode)
	}

	if address.Street != expected.Street {
		t.Errorf("Street mismatch: expected %v, got %v", expected.Street, address.Street)
	}

	if address.Complement != expected.Complement {
		t.Errorf("Complement mismatch: expected %v, got %v", expected.Complement, address.Complement)
	}

	if address.District != expected.District {
		t.Errorf("District mismatch: expected %v, got %v", expected.District, address.District)
	}

	if address.City != expected.City {
		t.Errorf("City mismatch: expected %v, got %v", expected.City, address.City)
	}

	if address.State != expected.State {
		t.Errorf("State mismatch: expected %v, got %v", expected.State, address.State)
	}
}

func TestAddressRepository_InvalidCep(t *testing.T) {
	t.Setenv("TEST", "true")

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		responseBody := `{}`
		w.Write([]byte(responseBody))
	}))
	defer server.Close()

	repo := repository.NewAddressRepository(server.URL)

	cep := "0"

	_, err := repo.GetAddress(cep)
	if err == nil {
		t.Fatalf("Expected an error but got nil")
	}

	expectedErrorMsg := "invalid zipcode"
	if !strings.Contains(err.Error(), expectedErrorMsg) {
		t.Errorf("Error message does not match expected. \nExpected to contain: %s\nGot: %s", expectedErrorMsg, err.Error())
	}
}

func TestAddressRepository_NotFindZipcode(t *testing.T) {
	t.Setenv("TEST", "true")

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		responseBody := `{"erro":true}`
		w.Write([]byte(responseBody))
	}))
	defer server.Close()

	repo := repository.NewAddressRepository(server.URL)

	cep := "99999999"

	_, err := repo.GetAddress(cep)
	if err == nil {
		t.Fatalf("Expected an error but got nil")
	}

	expectedErrorMsg := "can not find zipcode"
	if !strings.Contains(err.Error(), expectedErrorMsg) {
		t.Errorf("Error message does not match expected. \nExpected to contain: %s\nGot: %s", expectedErrorMsg, err.Error())
	}
}

func TestAddressRepository_ErrorHttp(t *testing.T) {
	t.Setenv("TEST", "true")

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "", http.StatusMovedPermanently)
	}))
	defer server.Close()

	repo := repository.NewAddressRepository(server.URL)

	cep := "12345678"

	_, err := repo.GetAddress(cep)
	if err == nil {
		t.Fatalf("Expected an error but got nil")
	}

	expectedErrorMsg := fmt.Sprintf("error when searching for zipcode %s information", cep)
	if !strings.Contains(err.Error(), expectedErrorMsg) {
		t.Errorf("Error message does not match expected. \nExpected to contain: %s\nGot: %s", expectedErrorMsg, err.Error())
	}
}

func TestAddressRepository_NotStatusOK(t *testing.T) {
	t.Setenv("TEST", "true")

	statusServerError := http.StatusInternalServerError

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Internal Server Error", statusServerError)
	}))
	defer server.Close()

	repo := repository.NewAddressRepository(server.URL)

	cep := "12345678"

	_, err := repo.GetAddress(cep)
	if err == nil {
		t.Fatalf("Expected an error but got nil")
	}

	expectedErrorMsg := fmt.Sprintf("ViaCEP api returned status %d for zipcode %s", statusServerError, cep)
	if !strings.Contains(err.Error(), expectedErrorMsg) {
		t.Errorf("Error message does not match expected. \nExpected to contain: %s\nGot: %s", expectedErrorMsg, err.Error())
	}
}

func TestAddressRepository_ErrorJsonDecoder(t *testing.T) {
	t.Setenv("TEST", "true")

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		responseBody := `error`
		w.Write([]byte(responseBody))
	}))
	defer server.Close()

	repo := repository.NewAddressRepository(server.URL)

	cep := "12345678"

	_, err := repo.GetAddress(cep)
	if err == nil {
		t.Fatalf("Expected an error but got nil")
	}

	expectedErrorMsg := fmt.Sprintf("error when decoding ViaCEP api response to zipcode %s", cep)
	if !strings.Contains(err.Error(), expectedErrorMsg) {
		t.Errorf("Error message does not match expected. \nExpected to contain: %s\nGot: %s", expectedErrorMsg, err.Error())
	}
}
