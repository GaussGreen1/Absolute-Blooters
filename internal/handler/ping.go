package handler

import (
	"encoding/json"
	"log"
	"net/http"
)

// PingResponse represents the response structure for the ping endpoint
type PingResponse struct {
	Message string `json:"message"`
	Status  int    `json:"status"`
}

// PingHandler handles the /api/ping endpoint
func PingHandler(w http.ResponseWriter, r *http.Request) {
	// Create the response
	response := PingResponse{
		Message: "pong",
		Status:  http.StatusOK,
	}

	log.Println("Pinging a new log here")

	// Set content type to JSON
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	// Encode and send the response
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}
