prod-start:
	docker compose -f dev.docker-compose.yml up -d

prod-stop:
	docker compose -f dev.docker-compose.yml stop

prod-down:
	docker compose -f dev.docker-compose.yml down

dev-run-service-a:
	docker compose exec dev-go-telemetry-cep-temperature go run cmd/input_server/main.go

dev-run-service-b:
	docker compose exec dev-go-telemetry-cep-temperature go run cmd/temperature_server/main.go

dev-run-tests:
	docker compose exec dev-go-telemetry-cep-temperature go test ./... -v

prod-start:
	docker compose -f prod.docker-compose.yml up -d

prod-stop:
	docker compose -f prod.docker-compose.yml stop

prod-down:
	docker compose -f prod.docker-compose.yml down
