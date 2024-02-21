package handler_test

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/aronkst/go-telemetry-cep-temperature/internal/input_server/handler"
	"github.com/aronkst/go-telemetry-cep-temperature/internal/input_server/model"
	temperatureServerModel "github.com/aronkst/go-telemetry-cep-temperature/internal/temperature_server/model"
)

type MockInputService struct {
	Temperature *temperatureServerModel.Temperature
	Err         error
}

func (m *MockInputService) GetTemperatureByCep(*model.Zipcode) (*temperatureServerModel.Temperature, error) {
	return m.Temperature, m.Err
}

func TestGetTemperatureByCep_ValidCEP(t *testing.T) {
	mockService := &MockInputService{
		Temperature: &temperatureServerModel.Temperature{City: "Cidade", Celsius: 30, Fahrenheit: 86, Kelvin: 303.15},
		Err:         nil,
	}

	handler := handler.NewInputHandler(mockService)

	body := bytes.NewBufferString(`{"cep": "12345678"}`)
	req, err := http.NewRequest("POST", "/", body)
	if err != nil {
		t.Fatal(err)
	}

	responseRecorder := httptest.NewRecorder()
	handler.GetTemperatureByCep(responseRecorder, req)

	if status := responseRecorder.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	expected := `{"city":"Cidade","temp_C":30,"temp_F":86,"temp_K":303.15}`
	if strings.Trim(responseRecorder.Body.String(), "\n") != expected {
		t.Errorf("handler returned unexpected body: got %v want %v", responseRecorder.Body.String(), expected)
	}
}

func TestGetTemperatureByCep_InvalidBody(t *testing.T) {
	mockService := &MockInputService{
		Temperature: nil,
		Err:         nil,
	}

	handler := handler.NewInputHandler(mockService)

	body := bytes.NewBufferString(`error`)
	req, err := http.NewRequest("POST", "/", body)
	if err != nil {
		t.Fatal(err)
	}

	responseRecorder := httptest.NewRecorder()
	handler.GetTemperatureByCep(responseRecorder, req)

	if status := responseRecorder.Code; status != http.StatusInternalServerError {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusInternalServerError)
	}

	expected := "invalid body"
	if strings.Trim(responseRecorder.Body.String(), "\n") != expected {
		t.Errorf("handler returned unexpected body: got %v want %v", responseRecorder.Body.String(), expected)
	}
}

func TestGetTemperatureByCep_InvalidCEP(t *testing.T) {
	mockService := &MockInputService{
		Temperature: nil,
		Err:         fmt.Errorf("invalid zipcode"),
	}

	handler := handler.NewInputHandler(mockService)

	body := bytes.NewBufferString(`{"cep": "12345-678"}`)
	req, err := http.NewRequest("POST", "/", body)
	if err != nil {
		t.Fatal(err)
	}

	responseRecorder := httptest.NewRecorder()
	handler.GetTemperatureByCep(responseRecorder, req)

	if status := responseRecorder.Code; status != http.StatusUnprocessableEntity {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusUnprocessableEntity)
	}

	expected := "invalid zipcode"
	if strings.Trim(responseRecorder.Body.String(), "\n") != expected {
		t.Errorf("handler returned unexpected body: got %v want %v", responseRecorder.Body.String(), expected)
	}
}

func TestGetTemperatureByCep_NotFound(t *testing.T) {
	mockService := &MockInputService{
		Temperature: nil,
		Err:         fmt.Errorf("can not find zipcode"),
	}

	handler := handler.NewInputHandler(mockService)

	body := bytes.NewBufferString(`{"cep": "99999999"}`)
	req, err := http.NewRequest("POST", "/", body)
	if err != nil {
		t.Fatal(err)
	}

	responseRecorder := httptest.NewRecorder()
	handler.GetTemperatureByCep(responseRecorder, req)

	if status := responseRecorder.Code; status != http.StatusNotFound {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusNotFound)
	}

	expected := "can not find zipcode"
	if strings.Trim(responseRecorder.Body.String(), "\n") != expected {
		t.Errorf("handler returned unexpected body: got %v want %v", responseRecorder.Body.String(), expected)
	}
}

func TestGetTemperatureByCep_InternalServerError(t *testing.T) {
	mockService := &MockInputService{
		Temperature: nil,
		Err:         fmt.Errorf("internal server error"),
	}

	handler := handler.NewInputHandler(mockService)

	body := bytes.NewBufferString(`{"cep": "00000000"}`)
	req, err := http.NewRequest("POST", "/", body)
	if err != nil {
		t.Fatal(err)
	}

	responseRecorder := httptest.NewRecorder()
	handler.GetTemperatureByCep(responseRecorder, req)

	if status := responseRecorder.Code; status != http.StatusInternalServerError {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusInternalServerError)
	}

	expected := "internal server error"
	if strings.Trim(responseRecorder.Body.String(), "\n") != expected {
		t.Errorf("handler returned unexpected body: got %v want %v", responseRecorder.Body.String(), expected)
	}
}
