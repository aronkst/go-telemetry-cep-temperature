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
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
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
	collectorENV := utils.GetEnvOrDefault("COLLECTOR_URL", "collector")
	collectorURL := fmt.Sprintf("%s:4317", collectorENV)

	ctx := context.Background()

	exporter, err := otlptracegrpc.New(ctx,
		otlptracegrpc.WithInsecure(),
		otlptracegrpc.WithEndpoint(collectorURL),
	)
	if err != nil {
		log.Fatalf("failed to create OTLP gRPC trace exporter: %v", err)
	}

	resource := resource.NewWithAttributes(
		semconv.SchemaURL,
		semconv.ServiceNameKey.String("Service B"),
	)

	tracerProvider := trace.NewTracerProvider(
		trace.WithBatcher(exporter),
		trace.WithResource(resource),
	)

	otel.SetTracerProvider(tracerProvider)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))

	return func() {
		if err := tracerProvider.Shutdown(ctx); err != nil {
			log.Fatalf("error closing trace provider: %v", err)
		}
	}
}
