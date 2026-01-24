package models

import "time"

type Goal struct {
	ID          int    `json:"id"`
	GameID      int    `json:"game_id"`
	Description string `json:"description"`
	HomeTeam    string `json:"home_team"`
	AwayTeam    string `json:"away_team"`
	Goalscorer  string `json:"goalscorer"`
	Minute      string `json:"minute"`
	Url         string `json:"url"`
	RedditURL   string `json:"reddit_url"`
	Mirrors     string `json:"mirrors"`
	HomeScore   int    `json:"home_score"`
	AwayScore   int    `json:"away_score"`
	Away        bool   `json:"away"` // true if goalscorer plays for away team
}

type Game struct {
	ID        int       `json:"id"`
	HomeTeam  string    `json:"home_team"`
	AwayTeam  string    `json:"away_team"`
	HomeScore int       `json:"home_score"`
	AwayScore int       `json:"away_score"`
	Goals     []Goal    `json:"goals"`
	Timestamp time.Time `json:"timestamp"`
}

type GamesResponse struct {
	Games  []Game `json:"games"`
	Status int    `json:"status"`
}
