package handler

import (
	"encoding/json"
	"log"
	"net/http"
)

type PingResponse struct {
	Message string `json:"message"`
	Status  int    `json:"status"`
}

func PingHandler(w http.ResponseWriter, r *http.Request) {
	response := PingResponse{
		Message: "pong",
		Status:  http.StatusOK,
	}

	log.Println("Pinging a new log here")

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:5173")

	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}
