FROM golang:1.25.2 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 go build -o /app/bin/api ./cmd/api

FROM alpine:latest

WORKDIR /app

# Copy compiled binary
COPY --from=builder /app/bin/api /app/api

# Expose app port
EXPOSE 8080

CMD ["./api"]