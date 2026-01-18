package handler

import (
	"encoding/json"
	"log"
	"net/http"

	"blooters/internal/db"
	"blooters/internal/models"
)

func GamesHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Fetching games...")

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
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:5173")

	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}
