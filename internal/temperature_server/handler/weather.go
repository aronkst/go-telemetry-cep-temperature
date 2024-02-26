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
	tracer := otel.Tracer("WeatherHandler")

	ctx := r.Context()
	ctxDistributed := otel.GetTextMapPropagator().Extract(ctx, propagation.HeaderCarrier(r.Header))

	ctx, spanRoute := tracer.Start(ctx, "GET /")
	defer spanRoute.End()

	ctx, span := tracer.Start(ctx, "WeatherHandler.GetWeatherByCEP")
	defer span.End()

	ctxDistributed, spanDistributedRoute := tracer.Start(ctxDistributed, "GET /")
	defer spanDistributedRoute.End()

	ctxDistributed, spanDistributed := tracer.Start(ctxDistributed, "WeatherHandler.GetWeatherByCEP")
	defer spanDistributed.End()

	cep := r.URL.Query().Get("cep")

	temperature, err := h.weatherService.GetWeatherByCEP(cep, ctx, ctxDistributed)
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
