package main

import (
	"blooters/internal/db"
	"blooters/internal/server"
	"log"
)

func main() {
	// initialize database
	if err := db.Init(); err != nil {
		log.Fatalf("failed to initialize database: %v", err)
	}
	defer func() {
		if err := db.Close(); err != nil {
			log.Printf("error closing db: %v", err)
		}
	}()

	srv := server.NewServer()

	if err := srv.Start(":8080"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
