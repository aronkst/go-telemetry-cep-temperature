package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/aronkst/go-telemetry-cep-temperature/internal/temperature_server/handler"
	"github.com/aronkst/go-telemetry-cep-temperature/internal/temperature_server/repository"
	"github.com/aronkst/go-telemetry-cep-temperature/internal/temperature_server/service"
	"github.com/aronkst/go-telemetry-cep-temperature/pkg/utils"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/zipkin"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	cleanup := initTracer()
	defer cleanup()

	addressRepository := repository.NewAddressRepository("https://viacep.com.br/ws/%s/json/")
	coordinatesRepository := repository.NewCoordinatesRepository("https://nominatim.openstreetmap.org/search")
	weatherByAddressRepository := repository.NewWeatherByAddressRepository("https://wttr.in/%s,%s,Brazil?format=j1")
	weatherByCoordinatesRepository := repository.NewWeatherByCoordinatesRepository("https://api.open-meteo.com/v1/forecast?latitude=%s&longitude=%s&current_weather=true")

	weatherService := service.NewWeatherService(addressRepository, coordinatesRepository, weatherByAddressRepository, weatherByCoordinatesRepository)

	weatherHandler := handler.NewWeatherHandler(weatherService)

	router := chi.NewRouter()
	router.Use(middleware.Logger)

	router.Get("/", weatherHandler.GetWeatherByCEP)

	log.Printf("server started on port 8080")

	err := http.ListenAndServe(":8080", router)
	if err != nil {
		log.Fatal("error starting server: ", err)
	}
}

func initTracer() func() {
	zipkinENV := utils.GetEnvOrDefault("ZIPKIN_URL", "zipkin")
	zipkinURL := fmt.Sprintf("http://%s:9411/api/v2/spans", zipkinENV)

	exporter, err := zipkin.New(zipkinURL)
	if err != nil {
		log.Fatal("failed to create Zipkin exporter: ", err)
	}

	tracerProvider := trace.NewTracerProvider(
		trace.WithBatcher(exporter),
		trace.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String("Service B"),
		)),
	)

	otel.SetTracerProvider(tracerProvider)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))

	return func() {
		if err := tracerProvider.Shutdown(context.Background()); err != nil {
			log.Fatal("error closing trace provider: ", err)
		}
	}
}
