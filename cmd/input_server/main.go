package main

import (
	"context"
	"log"
	"net/http"

	"github.com/aronkst/go-telemetry-cep-temperature/internal/input_server/handler"
	"github.com/aronkst/go-telemetry-cep-temperature/internal/input_server/repository"
	"github.com/aronkst/go-telemetry-cep-temperature/internal/input_server/service"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/zipkin"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
)

func main() {
	cleanup := initTracer()
	defer cleanup()

	temperatureRepository := repository.NewTemperatureRepository("http://localhost:8080/?cep=%s")

	inputService := service.NewInputService(temperatureRepository)

	inputHandler := handler.NewInputHandler(inputService)

	router := chi.NewRouter()
	router.Use(middleware.Logger)

	router.Post("/", otelhttp.NewHandler(http.HandlerFunc(inputHandler.GetTemperatureByCep), "POST /").ServeHTTP)

	log.Printf("server started on port 3000")

	err := http.ListenAndServe(":3000", router)
	if err != nil {
		log.Fatal("error starting server: ", err)
	}
}

func initTracer() func() {
	exporter, err := zipkin.New("http://zipkin:9411/api/v2/spans")
	if err != nil {
		log.Fatal("failed to create Zipkin exporter: ", err)
	}

	tracerProvider := trace.NewTracerProvider(
		trace.WithBatcher(exporter),
		trace.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String("Service A"),
		)),
	)

	otel.SetTracerProvider(tracerProvider)

	return func() {
		if err := tracerProvider.Shutdown(context.Background()); err != nil {
			log.Fatal("error closing trace provider: ", err)
		}
	}
}
