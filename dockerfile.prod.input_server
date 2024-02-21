FROM golang:1.22.0 AS builder

WORKDIR /app

RUN mkdir -p cmd/server
RUN mkdir -p internal/handler
RUN mkdir -p internal/model
RUN mkdir -p internal/repository
RUN mkdir -p internal/service
RUN mkdir -p pkg/utils

COPY go.mod ./
COPY go.sum ./
COPY cmd/server/main.go ./cmd/server

COPY internal/handler/weather_test.go ./internal/handler
COPY internal/handler/weather.go ./internal/handler

COPY internal/model/address.go ./internal/model
COPY internal/model/coordinates.go ./internal/model
COPY internal/model/temperature.go ./internal/model
COPY internal/model/weather.go ./internal/model

COPY internal/repository/address_test.go ./internal/repository
COPY internal/repository/address.go ./internal/repository
COPY internal/repository/coordinates_test.go ./internal/repository
COPY internal/repository/coordinates.go ./internal/repository
COPY internal/repository/weather_by_address_test.go ./internal/repository
COPY internal/repository/weather_by_address.go ./internal/repository
COPY internal/repository/weather_by_coordinates_test.go ./internal/repository
COPY internal/repository/weather_by_coordinates.go ./internal/repository

COPY internal/service/weather_test.go ./internal/service
COPY internal/service/weather.go ./internal/service

COPY pkg/utils/clean_string_test.go ./pkg/utils
COPY pkg/utils/clean_string.go ./pkg/utils
COPY pkg/utils/is_number_test.go ./pkg/utils
COPY pkg/utils/is_number.go ./pkg/utils
COPY pkg/utils/number_converter_test.go ./pkg/utils
COPY pkg/utils/number_converter.go ./pkg/utils
COPY pkg/utils/temperature_converter_test.go ./pkg/utils
COPY pkg/utils/temperature_converter.go ./pkg/utils

RUN go mod download

RUN CGO_ENABLED=0 GOOS=linux go build -o go-cep-temperature cmd/server/main.go

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/go-cep-temperature ./

CMD ["./go-cep-temperature"]
