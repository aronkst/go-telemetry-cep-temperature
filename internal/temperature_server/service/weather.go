package service

import (
	"context"

	"github.com/aronkst/go-telemetry-cep-temperature/internal/temperature_server/model"
	"github.com/aronkst/go-telemetry-cep-temperature/internal/temperature_server/repository"
	"github.com/aronkst/go-telemetry-cep-temperature/pkg/utils"
	"go.opentelemetry.io/otel"
)

type WeatherService interface {
	GetWeatherByCEP(string, context.Context, context.Context) (*model.Temperature, error)
}

type weatherService struct {
	addressRepository              repository.AddressRepository
	coordinatesRepository          repository.CoordinatesRepository
	weatherByAddressRepository     repository.WeatherByAddressRepository
	weatherByCoordinatesRepository repository.WeatherByCoordinatesRepository
}

func NewWeatherService(
	addressRepository repository.AddressRepository,
	coordinatesRepository repository.CoordinatesRepository,
	weatherByAddressRepository repository.WeatherByAddressRepository,
	weatherByCoordinatesRepository repository.WeatherByCoordinatesRepository,
) WeatherService {
	return &weatherService{
		addressRepository:              addressRepository,
		coordinatesRepository:          coordinatesRepository,
		weatherByAddressRepository:     weatherByAddressRepository,
		weatherByCoordinatesRepository: weatherByCoordinatesRepository,
	}
}

func (s *weatherService) GetWeatherByCEP(cep string, ctx context.Context, ctxDistributed context.Context) (*model.Temperature, error) {
	tracer := otel.Tracer("WeatherService")

	ctx, span := tracer.Start(ctx, "WeatherService.GetWeatherByCEP")
	defer span.End()

	ctxDistributed, spanDistributed := tracer.Start(ctxDistributed, "WeatherService.GetWeatherByCEP")
	defer spanDistributed.End()

	address, err := s.addressRepository.GetAddress(cep, ctx, ctxDistributed)
	if err != nil {
		return nil, err
	}

	var weather *model.Weather

	coordinates, err := s.coordinatesRepository.GetCoordinates(address, ctx, ctxDistributed)
	if err == nil {
		weather, err = s.weatherByCoordinatesRepository.GetWeather(coordinates, ctx, ctxDistributed)
		if err != nil {
			return nil, err
		}
	} else {
		weather, err = s.weatherByAddressRepository.GetWeather(address, ctx, ctxDistributed)
		if err != nil {
			return nil, err
		}
	}

	temperature := &model.Temperature{
		City:       address.City,
		Celsius:    weather.Temperature,
		Fahrenheit: utils.CelsiusToFahrenheit(weather.Temperature),
		Kelvin:     utils.CelsiusToKelvin(weather.Temperature),
	}

	return temperature, nil
}
