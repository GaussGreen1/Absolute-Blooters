
.PHONY: build run clean db-up db-down


build:
    go build -o bin/api.exe ./cmd/api

run: build
    bin/api.exe

db-up:
    docker compose up -d

db-down:
    docker compose down

clean:
    rm -rf bin/