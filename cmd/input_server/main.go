package main

import (
	"log"
	"net/http"

	"github.com/aronkst/go-telemetry-cep-temperature/internal/input_server/handler"
	"github.com/aronkst/go-telemetry-cep-temperature/internal/input_server/repository"
	"github.com/aronkst/go-telemetry-cep-temperature/internal/input_server/service"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	temperatureRepository := repository.NewTemperatureRepository("http://localhost:8080/?cep=%s")

	inputService := service.NewInputService(temperatureRepository)

	inputHandler := handler.NewInputHandler(inputService)

	router := chi.NewRouter()
	router.Use(middleware.Logger)

	router.Post("/", inputHandler.GetTemperatureByCep)

	log.Printf("server started on port 3000")

	err := http.ListenAndServe(":3000", router)
	if err != nil {
		log.Fatal("error starting server: ", err)
	}
}
