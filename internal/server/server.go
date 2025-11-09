package server

import (
	"blooters/internal/handler"
	"blooters/internal/middleware"
	"log"
	"net/http"
)

type Server struct {
	mux http.Handler
}

func NewServer() *Server {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /api/ping", handler.PingHandler)
	mux.HandleFunc("GET /api/games", handler.GamesHandler)

	//Logging middleware:
	handler := middleware.LoggingMiddleware(mux)

	return &Server{
		mux: handler,
	}
}

func (s *Server) Start(addr string) error {
	log.Printf("Server starting on %s", addr)
	return http.ListenAndServe(addr, s.mux)
}
