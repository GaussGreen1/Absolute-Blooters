package main

import (
	"blooters/internal/db"
	"blooters/internal/metrics"
	"blooters/internal/reddit"
	"blooters/internal/server"
	"log"
	"net/http"
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
				metrics.GoalsFetchCount.WithLabelValues("error").Inc()
				continue
			}
			metrics.GoalsFetchCount.WithLabelValues("success").Inc()

			if err := db.StoreGoals(goals); err != nil {
				log.Printf("Error storing goals: %s\n", err)
				metrics.GoalsStoreCount.WithLabelValues("error").Inc()
			} else {
				log.Printf("StoreGoals end")
				metrics.GoalsStoreCount.WithLabelValues("success").Inc()
			}

			// Populate mirrors for goals that don't have them
			if err := reddit.PopulateMirrors(); err != nil {
				log.Printf("Error populating mirrors: %s\n", err)
				metrics.MirrorsPopulateCount.WithLabelValues("error").Inc()
			} else {
				metrics.MirrorsPopulateCount.WithLabelValues("success").Inc()
			}

			// Call the Ping API to keep the server active:
			resp, err := http.Get("https://absolute-blooters.onrender.com/api/ping")
			if err != nil {
				log.Fatal(err)
			}
			resp.Body.Close() // close immediately
		}
	}()

	// Periodically (5h) remove old goals
	tickerLimit := time.NewTicker(5 * time.Hour)
	go func() {
		for range tickerLimit.C {
			if err := db.RemoveOldGoals(); err != nil {
				log.Printf("Error removing old goals: %s\n", err)
				metrics.RemoveOldGoalsCount.WithLabelValues("error").Inc()
			} else {
				log.Printf("RemoveOldGoals end")
				metrics.RemoveOldGoalsCount.WithLabelValues("success").Inc()
			}

		}
	}()

	// Keep the program running
	select {}
}
