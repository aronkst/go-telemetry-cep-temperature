package service

import (
	"context"

	"github.com/aronkst/go-telemetry-cep-temperature/internal/input_server/model"
	"github.com/aronkst/go-telemetry-cep-temperature/internal/input_server/repository"
	temperatureServerModel "github.com/aronkst/go-telemetry-cep-temperature/internal/temperature_server/model"
	"go.opentelemetry.io/otel"
)

type InputService interface {
	GetTemperatureByCep(*model.Zipcode, context.Context) (*temperatureServerModel.Temperature, error)
}

type inputService struct {
	temperatureRepository repository.TemperatureRepository
}

func NewInputService(
	temperatureRepository repository.TemperatureRepository,
) InputService {
	return &inputService{
		temperatureRepository: temperatureRepository,
	}
}

func (s *inputService) GetTemperatureByCep(zipcode *model.Zipcode, ctx context.Context) (*temperatureServerModel.Temperature, error) {
	tracer := otel.Tracer("Service")

	ctx, span := tracer.Start(ctx, "InputService.GetTemperatureByCep")
	defer span.End()

	temperature, err := s.temperatureRepository.GetTemperature(zipcode, ctx)
	if err != nil {
		return nil, err
	}

	return temperature, nil
}
