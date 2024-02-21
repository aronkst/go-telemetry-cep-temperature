package repository

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/aronkst/go-telemetry-cep-temperature/internal/model"
)

type WeatherByCoordinatesRepository interface {
	GetWeather(coordinates *model.Coordinates) (*model.Weather, error)
}

type weatherByCoordinatesRepository struct {
	URL string
}

func NewWeatherByCoordinatesRepository(url string) WeatherByCoordinatesRepository {
	return &weatherByCoordinatesRepository{
		URL: url,
	}
}

func (repository *weatherByCoordinatesRepository) GetWeather(coordinates *model.Coordinates) (*model.Weather, error) {
	var url string

	if os.Getenv("TEST") == "true" {
		url = repository.URL
	} else {
		url = fmt.Sprintf(repository.URL, coordinates.Latitude, coordinates.Longitude)
	}

	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("error when searching for weather forecast: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("weather api returned status %d", resp.StatusCode)
	}

	var tempWeather struct {
		CurrentWeather struct {
			Temperature float64 `json:"temperature"`
		} `json:"current_weather"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&tempWeather); err != nil {
		return nil, fmt.Errorf("error parsing json: %w", err)
	}

	if tempWeather.CurrentWeather.Temperature > 0 {
		weather := &model.Weather{Temperature: tempWeather.CurrentWeather.Temperature}

		return weather, nil
	} else {
		return nil, fmt.Errorf("temperature error")
	}
}
