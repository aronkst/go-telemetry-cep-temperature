package main

import (
	"log"
	"net/http"

	"github.com/aronkst/go-telemetry-cep-temperature/internal/temperature_server/handler"
	"github.com/aronkst/go-telemetry-cep-temperature/internal/temperature_server/repository"
	"github.com/aronkst/go-telemetry-cep-temperature/internal/temperature_server/service"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
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
