package handler

import (
	"encoding/json"
	"net/http"

	"github.com/aronkst/go-telemetry-cep-temperature/internal/input_server/model"
	"github.com/aronkst/go-telemetry-cep-temperature/internal/input_server/service"
	"go.opentelemetry.io/otel"
)

type InputHandler struct {
	inputService service.InputService
}

func NewInputHandler(inputService service.InputService) *InputHandler {
	return &InputHandler{
		inputService: inputService,
	}
}

func (h *InputHandler) GetTemperatureByCep(w http.ResponseWriter, r *http.Request) {
	tracer := otel.Tracer("Handler")

	ctx, span := tracer.Start(r.Context(), "InputHandler.GetTemperatureByCep")
	defer span.End()

	var zipcode model.Zipcode

	err := json.NewDecoder(r.Body).Decode(&zipcode)
	if err != nil {
		http.Error(w, "invalid body", http.StatusInternalServerError)
		return
	}

	temperature, err := h.inputService.GetTemperatureByCep(&zipcode, ctx)
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
