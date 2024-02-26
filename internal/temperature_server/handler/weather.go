package handler

import (
	"encoding/json"
	"net/http"

	"github.com/aronkst/go-telemetry-cep-temperature/internal/temperature_server/service"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
)

type WeatherHandler struct {
	weatherService service.WeatherService
}

func NewWeatherHandler(weatherService service.WeatherService) *WeatherHandler {
	return &WeatherHandler{
		weatherService: weatherService,
	}
}

func (h *WeatherHandler) GetWeatherByCEP(w http.ResponseWriter, r *http.Request) {
	ctx := otel.GetTextMapPropagator().Extract(r.Context(), propagation.HeaderCarrier(r.Header))

	tracer := otel.Tracer("WeatherHandler")

	_, span := tracer.Start(ctx, "WeatherHandler.GetWeatherByCEP")
	defer span.End()

	cep := r.URL.Query().Get("cep")

	temperature, err := h.weatherService.GetWeatherByCEP(cep, ctx)
	if err != nil {
		var errorStatusCode int

		switch err.Error() {
		case "invalid zipcode":
			errorStatusCode = http.StatusUnprocessableEntity
		case "can not find zipcode":
			errorStatusCode = http.StatusNotFound
		default:
			errorStatusCode = http.StatusInternalServerError
		}

		http.Error(w, err.Error(), errorStatusCode)

		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(temperature)
}
