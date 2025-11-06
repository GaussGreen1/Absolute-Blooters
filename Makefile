.PHONY: build run clean

# Build the executable for Windows
build:
    go build -o bin/api.exe ./cmd/api

# Build and then run the executable
run: build
    bin/api.exe

# Clean up build artifacts
clean:
    rm -rf bin/