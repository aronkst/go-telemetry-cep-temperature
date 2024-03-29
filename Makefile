dev-start:
	docker compose -f docker-compose.dev.yml up -d

dev-stop:
	docker compose -f docker-compose.dev.yml stop

dev-down:
	docker compose -f docker-compose.dev.yml down

dev-run-service-a:
	docker compose -f docker-compose.dev.yml exec dev go run cmd/input_server/main.go

dev-run-service-b:
	docker compose -f docker-compose.dev.yml exec dev go run cmd/temperature_server/main.go

dev-run-tests:
	docker compose -f docker-compose.dev.yml exec dev go test ./... -v

prod-start:
	docker compose -f docker-compose.prod.yml up -d

prod-stop:
	docker compose -f docker-compose.prod.yml stop

prod-down:
	docker compose -f docker-compose.prod.yml down
