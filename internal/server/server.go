package server

import (
	"blooters/internal/handler"
	"blooters/internal/middleware"
	"log"
	"net/http"
)

// Server represents the HTTP server
type Server struct {
	mux http.Handler
}

// NewServer creates a new HTTP server
func NewServer() *Server {
	mux := http.NewServeMux()

	// Register routes using the enhanced pattern matching in Go 1.22+
	// This allows us to specify the HTTP method directly in the pattern
	mux.HandleFunc("GET /api/ping", handler.PingHandler)
	mux.HandleFunc("GET /api/games", handler.GamesHandler)

	// Wrap the mux with logging middleware
	handler := middleware.LoggingMiddleware(mux)

	return &Server{
		mux: handler,
	}
}

// Start starts the HTTP server
func (s *Server) Start(addr string) error {
	log.Printf("Server starting on %s", addr)
	return http.ListenAndServe(addr, s.mux)
}
