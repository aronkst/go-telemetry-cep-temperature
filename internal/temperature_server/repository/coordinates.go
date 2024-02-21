package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/aronkst/go-telemetry-cep-temperature/internal/temperature_server/model"
	"go.opentelemetry.io/otel"
)

type CoordinatesRepository interface {
	GetCoordinates(*model.Address, context.Context) (*model.Coordinates, error)
}

type coordinatesRepository struct {
	URL        string
	BaseMethod string
}

func NewCoordinatesRepository(url string) CoordinatesRepository {
	return &coordinatesRepository{
		URL: url,
	}
}

func (r *coordinatesRepository) GetCoordinates(address *model.Address, ctx context.Context) (*model.Coordinates, error) {
	tracer := otel.Tracer("Repository")

	_, span := tracer.Start(ctx, "CoordinatesRepository.GetCoordinates")
	defer span.End()

	baseURL := r.URL

	params := url.Values{}
	params.Add("city", address.City)
	params.Add("state", address.State)
	params.Add("country", "Brasil")
	params.Add("format", "json")

	req, err := http.NewRequest("GET", baseURL+"?"+params.Encode(), nil)
	if err != nil {
		return nil, fmt.Errorf("error when creating request: %w", err)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error when searching for coordinates for the address: %w", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("coordinates api returned status %d: %w", resp.StatusCode, err)
	}

	var results []struct {
		Lat string `json:"lat"`
		Lon string `json:"lon"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&results); err != nil {
		return nil, fmt.Errorf("error decoding coordinates api response: %w", err)
	}

	if len(results) == 0 {
		return nil, fmt.Errorf("no coordinates found for the address")
	}

	coordinates := &model.Coordinates{
		Latitude:  results[0].Lat,
		Longitude: results[0].Lon,
	}

	return coordinates, nil
}
