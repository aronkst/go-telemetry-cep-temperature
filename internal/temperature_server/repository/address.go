package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/aronkst/go-telemetry-cep-temperature/internal/temperature_server/model"
	"github.com/aronkst/go-telemetry-cep-temperature/pkg/utils"
	"go.opentelemetry.io/otel"
)

type AddressRepository interface {
	GetAddress(string, context.Context, context.Context) (*model.Address, error)
}

type addressRepository struct {
	URL string
}

func NewAddressRepository(url string) AddressRepository {
	return &addressRepository{
		URL: url,
	}
}

func (r *addressRepository) GetAddress(cep string, ctx context.Context, ctxDistributed context.Context) (*model.Address, error) {
	tracer := otel.Tracer("AddressRepository")

	_, span := tracer.Start(ctx, "AddressRepository.GetAddress")
	defer span.End()

	_, spanDistributed := tracer.Start(ctxDistributed, "AddressRepository.GetAddress")
	defer spanDistributed.End()

	if cep == "" || len(cep) != 8 || !utils.IsNumber(cep) {
		return nil, fmt.Errorf("invalid zipcode")
	}

	var url string

	if os.Getenv("TEST") == "true" {
		url = r.URL
	} else {
		url = fmt.Sprintf(r.URL, cep)
	}

	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("error when searching for zipcode %s information: %w", cep, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("ViaCEP api returned status %d for zipcode %s: %w", resp.StatusCode, cep, err)
	}

	var address model.Address
	if err := json.NewDecoder(resp.Body).Decode(&address); err != nil {
		return nil, fmt.Errorf("error when decoding ViaCEP api response to zipcode %s: %w", cep, err)
	}

	if address.PostalCode == "" {
		return nil, fmt.Errorf("can not find zipcode")
	}

	return &address, nil
}
