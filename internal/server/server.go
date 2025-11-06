package server

import (
	"blooters/internal/handler"
	"log"
	"net/http"
)

// Server represents the HTTP server
type Server struct {
	mux *http.ServeMux
}

// NewServer creates a new HTTP server
func NewServer() *Server {
	mux := http.NewServeMux()

	// Register routes using the enhanced pattern matching in Go 1.22+
	// This allows us to specify the HTTP method directly in the pattern
	mux.HandleFunc("GET /api/ping", handler.PingHandler)
	mux.HandleFunc("GET /api/games", handler.GamesHandler)

	return &Server{
		mux: mux,
	}
}

// Start starts the HTTP server
func (s *Server) Start(addr string) error {
	log.Printf("Server starting on %s", addr)
	return http.ListenAndServe(addr, s.mux)
}
