# Build stage
FROM golang:1.25.2 AS builder

WORKDIR /app

# Copy go mod files first for better caching
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest
COPY . .

# Build the application with static linking
RUN CGO_ENABLED=0 go build -o /app/bin/api ./cmd/api

# Run stage — small image
FROM alpine:latest

WORKDIR /app

# Copy compiled binary
COPY --from=builder /app/bin/api /app/api

# Expose app port
EXPOSE 8080

CMD ["./api"]