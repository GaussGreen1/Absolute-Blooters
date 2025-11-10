package handler

import (
	"encoding/json"
	"log"
	"net/http"

	"blooters/internal/db"
	"blooters/internal/models"
	"blooters/internal/reddit"
)

func GamesHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Fetching games...")

	// Fetch goals from r/soccer
	redditGoals, err := reddit.FetchGoals()
	if err != nil {
		log.Printf("Warning: failed to fetch from Reddit: %v\n", err)
		// Fall back to database if Reddit fetch fails
	} else if len(redditGoals) > 0 {
		// Store fetched goals in database
		if err := db.StoreGoals(redditGoals); err != nil {
			log.Printf("Warning: failed to store goals: %v\n", err)
		}
	}

	games, err := db.GetGames()
	if err != nil {
		http.Error(w, "Failed to load games", http.StatusInternalServerError)
		return
	}

	response := models.GamesResponse{
		Games:  games,
		Status: http.StatusOK,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}
