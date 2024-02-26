package handler

import (
	"context"
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
	tracer := otel.Tracer("InputHandler")

	ctx := r.Context()
	ctxDistributed := context.Background()

	ctx, spanRoute := tracer.Start(ctx, "POST /")
	defer spanRoute.End()

	ctx, span := tracer.Start(ctx, "InputHandler.GetTemperatureByCep")
	defer span.End()

	ctxDistributed, spanDistributedStart := tracer.Start(ctxDistributed, "Distributed")
	defer spanDistributedStart.End()

	ctxDistributed, spanDistributedStartRoute := tracer.Start(ctxDistributed, "POST /")
	defer spanDistributedStartRoute.End()

	ctxDistributed, spanDistributed := tracer.Start(ctxDistributed, "InputHandler.GetTemperatureByCep")
	defer spanDistributed.End()

	var zipcode model.Zipcode

	err := json.NewDecoder(r.Body).Decode(&zipcode)
	if err != nil {
		http.Error(w, "invalid body", http.StatusInternalServerError)
		return
	}

	temperature, err := h.inputService.GetTemperatureByCep(&zipcode, ctx, ctxDistributed)
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
