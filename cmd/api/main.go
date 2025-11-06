package main

import (
	"blooters/internal/server"
	"log"
)

func main() {
	srv := server.NewServer()

	if err := srv.Start(":8080"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
