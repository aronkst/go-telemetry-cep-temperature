package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/aronkst/go-telemetry-cep-temperature/internal/input_server/model"
	temperatureServerModel "github.com/aronkst/go-telemetry-cep-temperature/internal/temperature_server/model"
	"github.com/aronkst/go-telemetry-cep-temperature/pkg/utils"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
)

type TemperatureRepository interface {
	GetTemperature(*model.Zipcode, context.Context, context.Context) (*temperatureServerModel.Temperature, error)
}

type temperatureRepository struct {
	URL string
}

func NewTemperatureRepository(url string) TemperatureRepository {
	return &temperatureRepository{
		URL: url,
	}
}

func (r *temperatureRepository) GetTemperature(zipcode *model.Zipcode, ctx context.Context, ctxDistributed context.Context) (*temperatureServerModel.Temperature, error) {
	tracer := otel.Tracer("TemperatureRepository")

	_, span := tracer.Start(ctx, "TemperatureRepository.GetTemperature")
	defer span.End()

	ctxDistributed, spanDistributed := tracer.Start(ctxDistributed, "TemperatureRepository.GetTemperature")
	defer spanDistributed.End()

	cep := zipcode.Cep
	if cep == "" || len(cep) != 8 || !utils.IsNumber(cep) {
		return nil, fmt.Errorf("invalid zipcode")
	}

	var url string

	if os.Getenv("TEST") == "true" {
		url = r.URL
	} else {
		url = fmt.Sprintf(r.URL, cep)
	}

	client := &http.Client{}
	req, _ := http.NewRequest("GET", url, nil)

	otel.GetTextMapPropagator().Inject(ctxDistributed, propagation.HeaderCarrier(req.Header))

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error when searching for temperature by cep: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, fmt.Errorf("can not find zipcode")
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("temperature by cep api returned status %d: %w", resp.StatusCode, err)
	}

	var temperature temperatureServerModel.Temperature
	if err := json.NewDecoder(resp.Body).Decode(&temperature); err != nil {
		return nil, fmt.Errorf("error parsing json: %w", err)
	}

	return &temperature, nil
}
