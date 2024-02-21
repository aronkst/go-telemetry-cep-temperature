package repository

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/aronkst/go-telemetry-cep-temperature/internal/model"
	"github.com/aronkst/go-telemetry-cep-temperature/pkg/utils"
)

type WeatherByAddressRepository interface {
	GetWeather(address *model.Address) (*model.Weather, error)
}

type weatherByAddressRepository struct {
	URL string
}

func NewWeatherByAddressRepository(url string) WeatherByAddressRepository {
	return &weatherByAddressRepository{
		URL: url,
	}
}

func (repository *weatherByAddressRepository) GetWeather(address *model.Address) (*model.Weather, error) {
	var url string

	if os.Getenv("TEST") == "true" {
		url = repository.URL
	} else {
		url = fmt.Sprintf(repository.URL, utils.CleanString(address.City), utils.CleanString(address.State))
	}

	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("error when searching for weather forecast: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("weather api returned status %d: %w", resp.StatusCode, err)
	}

	var tempWeather struct {
		CurrentCondition []struct {
			TempC string `json:"temp_C"`
		} `json:"current_condition"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&tempWeather); err != nil {
		return nil, fmt.Errorf("error parsing json: %w", err)
	}

	if len(tempWeather.CurrentCondition) > 0 {
		weather := &model.Weather{Temperature: utils.StringToFloat64(tempWeather.CurrentCondition[0].TempC)}

		return weather, nil
	} else {
		return nil, fmt.Errorf("temperature error")
	}
}
