package main

import (
	"blooters/internal/db"
	"blooters/internal/reddit"
	"blooters/internal/server"
	"log"
	"time"
)

func main() {
	if err := db.Init(); err != nil {
		log.Fatalf("failed to initialize database: %v", err)
	}
	defer func() {
		if err := db.Close(); err != nil {
			log.Printf("error closing db: %v", err)
		}
	}()

	srv := server.NewServer()

	go func() {
		if err := srv.Start(":8080"); err != nil {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Periodically (10s) fetch goals from Reddit and store in database
	ticker := time.NewTicker(10 * time.Second)
	go func() {
		for range ticker.C {
			log.Println("GOALS FETCHED HERE (ticker triggered)")
			goals, err := reddit.FetchGoals()
			if err != nil {
				log.Printf("Error fetching goals: %s\n", err)
				continue
			}

			if err := db.StoreGoals(goals); err != nil {
				log.Printf("Error storing goals: %s\n", err)
			} else {
				log.Printf("StoreGoals end")
			}
		}
	}()

	// Keep the program running
	select {}
}
