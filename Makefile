tidy:
	go mod tidy

run:
	go run cmd/main.go

up:
	docker compose up -d

down:
	docker compose down