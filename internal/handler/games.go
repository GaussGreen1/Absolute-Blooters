package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/sirupsen/logrus"
)

type Goal struct {
	Description string    `json:"description"`
	Teams       [2]string `json:"teams"`
	Goalscorer  string    `json:"goalscorer"`
	Minute      string    `json:"minute"`
	Url         string    `json:"url"`
	Score       [2]int    `json:"score"`
	Allegiance  bool      `json:"allegiance"`
}

type Game struct {
	Teams     [2]string `json:"teams"`
	Score     [2]int    `json:"score"`
	Goals     []Goal    `json:"goals"`
	Timestamp time.Time `json:"timestamp"`
}

type GamesResponse struct {
	Games  []Game `json:"games"`
	Status int    `json:"status"`
}

func GamesHandler(w http.ResponseWriter, r *http.Request) {
	merinoGoal := Goal{
		Description: "Arsenal [1]-0 Chelsea - Mikel Merino 90+7'",
		Teams:       [2]string{"Arsenal", "Chelsea"},
		Goalscorer:  "Mikel Merino",
		Minute:      "90+7'",
		Url:         "https://www.youtube.com/watch?v=wnG96bon5IQ",
		Score:       [2]int{1, 0},
		Allegiance:  false,
	}

	trossardGoal := Goal{
		Description: "Arsenal [4]-2 Sunderland - Leandro Trossard 23' (Great Goal)",
		Teams:       [2]string{"Arsenal", "Sunderland"},
		Goalscorer:  "Leandro Trossard",
		Minute:      "23'",
		Url:         "https://www.youtube.com/watch?v=BRZbmT4yMWI",
		Score:       [2]int{4, 2},
		Allegiance:  false,
	}

	log.Println("Gaaaames")
	logrus.Info("Gaaaames logrus")

	firstGame := Game{
		Teams: [2]string{"Arsenal", "Chelsea"},
		Score: [2]int{1, 0},
		Goals: []Goal{
			merinoGoal,
		},
		Timestamp: time.Now(),
	}

	secondGame := Game{
		Teams: [2]string{"Arsenal", "Sunderland"},
		Score: [2]int{4, 2},
		Goals: []Goal{
			trossardGoal,
		},
		Timestamp: time.Now(),
	}

	response := GamesResponse{
		Games:  []Game{firstGame, secondGame},
		Status: http.StatusOK,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}
