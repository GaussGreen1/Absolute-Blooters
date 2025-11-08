
.PHONY: build run clean db-up db-down


# Build the executable for Windows
build:
    go build -o bin/api.exe ./cmd/api

# Build and then run the executable
run: build
    bin/api.exe

# Start Postgres via Docker Compose
db-up:
    docker compose up -d

# Stop Postgres and remove containers
db-down:
    docker compose down

# Clean up build artifacts
clean:
    rm -rf bin/