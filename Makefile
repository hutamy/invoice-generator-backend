tidy:
	go mod tidy

run:
	go run cmd/main.go

up:
	docker compose -f docker-compose-pg-only.yaml up -d

down:
	docker compose down

swagger:
	swag init --generalInfo cmd/main.go --output docs