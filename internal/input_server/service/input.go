package service

import (
	"github.com/aronkst/go-telemetry-cep-temperature/internal/input_server/model"
	"github.com/aronkst/go-telemetry-cep-temperature/internal/input_server/repository"
	temperatureServerModel "github.com/aronkst/go-telemetry-cep-temperature/internal/temperature_server/model"
)

type InputService interface {
	GetTemperatureByCep(*model.Zipcode) (*temperatureServerModel.Temperature, error)
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

func (s *inputService) GetTemperatureByCep(zipcode *model.Zipcode) (*temperatureServerModel.Temperature, error) {
	temperature, err := s.temperatureRepository.GetTemperature(zipcode)
	if err != nil {
		return nil, err
	}

	return temperature, nil
}
