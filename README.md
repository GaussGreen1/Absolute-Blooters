(Minimal instructions)

Run Postgres locally with Docker Compose (this automatically seeds schema and sample data):

1. Start Postgres:

	docker compose up -d

2. Build and run the API (locally):

	# ensure you have Go installed and the env file available
	go mod tidy
	go build ./...
	# run the binary created in cmd/api (or run via `go run cmd/api/main.go`)

Environment variables used (in `.env`): DB_HOST, DB_PORT, DB_USER, DB_PASSWORD, DB_NAME

Endpoints:
- GET /api/ping
- GET /api/games

