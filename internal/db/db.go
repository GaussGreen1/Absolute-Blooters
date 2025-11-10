package db

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"

	"blooters/internal/models"
)

var DB *sql.DB

// Init connects to Postgres using environment variables. It returns an error on failure.
func Init() error {
	host := getEnv("DB_HOST", "localhost")
	port := getEnv("DB_PORT", "5432")
	user := getEnv("DB_USER", "postgres")
	pass := getEnv("DB_PASSWORD", "postgres")
	name := getEnv("DB_NAME", "blooters")

	dsn := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s", user, pass, host, port, name)

	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return err
	}

	db.SetConnMaxLifetime(5 * time.Minute)
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(2)

	// ping with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := db.PingContext(ctx); err != nil {
		_ = db.Close()
		return err
	}

	DB = db
	return nil
}

func Close() error {
	if DB != nil {
		return DB.Close()
	}
	return nil
}

func getEnv(k, def string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return def
}

func GetGames() ([]models.Game, error) {
	if DB == nil {
		return nil, fmt.Errorf("database not initialized")
	}

	rows, err := DB.Query("SELECT id, home_team, away_team, home_score, away_score, timestamp FROM games ORDER BY id")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var games []models.Game
	for rows.Next() {
		var g models.Game
		var homeTeam, awayTeam string
		var homeScore, awayScore int
		var ts time.Time
		if err := rows.Scan(&g.ID, &homeTeam, &awayTeam, &homeScore, &awayScore, &ts); err != nil {
			return nil, err
		}
		g.HomeTeam = homeTeam
		g.AwayTeam = awayTeam
		g.HomeScore = homeScore
		g.AwayScore = awayScore
		g.Timestamp = ts

		g.Goals, err = loadGoalsForGame(g.ID, g.HomeTeam, g.AwayTeam)
		if err != nil {
			return nil, err
		}

		games = append(games, g)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return games, nil
}

func loadGoalsForGame(gameID int, homeTeam, awayTeam string) ([]models.Goal, error) {
	q := `SELECT id, description, goalscorer, minute, url, away, home_score, away_score FROM goals WHERE game_id=$1 ORDER BY id`
	rows, err := DB.Query(q, gameID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var goals []models.Goal
	for rows.Next() {
		var gl models.Goal
		var hs, as int
		if err := rows.Scan(&gl.ID, &gl.Description, &gl.Goalscorer, &gl.Minute, &gl.Url, &gl.Away, &hs, &as); err != nil {
			return nil, err
		}
		gl.GameID = gameID
		gl.HomeTeam = homeTeam
		gl.AwayTeam = awayTeam
		gl.HomeScore = hs
		gl.AwayScore = as
		goals = append(goals, gl)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return goals, nil
}

// StoreGoals stores goals from r/soccer into the database, creating games as needed
func StoreGoals(goals []models.Goal) error {
	if DB == nil {
		return fmt.Errorf("database not initialized")
	}

	// Group goals by game (home_team, away_team)
	gameMap := make(map[string]models.Game)
	for _, goal := range goals {
		key := goal.HomeTeam + " vs " + goal.AwayTeam
		if game, exists := gameMap[key]; exists {
			game.Goals = append(game.Goals, goal)
			gameMap[key] = game
		} else {
			gameMap[key] = models.Game{
				HomeTeam:  goal.HomeTeam,
				AwayTeam:  goal.AwayTeam,
				HomeScore: goal.HomeScore,
				AwayScore: goal.AwayScore,
				Goals:     []models.Goal{goal},
				Timestamp: time.Now(),
			}
		}
	}

	// Store each game and its goals
	for _, game := range gameMap {
		// Check if game already exists
		var existingGameID int
		err := DB.QueryRow(
			"SELECT id FROM games WHERE home_team=$1 AND away_team=$2",
			game.HomeTeam, game.AwayTeam,
		).Scan(&existingGameID)

		var gameID int
		if err == sql.ErrNoRows {
			// Insert new game
			err := DB.QueryRow(
				"INSERT INTO games (home_team, away_team, home_score, away_score, timestamp) VALUES ($1, $2, $3, $4, $5) RETURNING id",
				game.HomeTeam, game.AwayTeam, game.HomeScore, game.AwayScore, game.Timestamp,
			).Scan(&gameID)
			if err != nil {
				return fmt.Errorf("failed to insert game: %w", err)
			}
		} else if err != nil {
			return fmt.Errorf("failed to query game: %w", err)
		} else {
			gameID = existingGameID
		}

		// Insert goals for this game
		for _, goal := range game.Goals {
			_, err := DB.Exec(
				"INSERT INTO goals (game_id, description, goalscorer, minute, url, away, home_score, away_score) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)",
				gameID, goal.Description, goal.Goalscorer, goal.Minute, goal.Url, goal.Away, goal.HomeScore, goal.AwayScore,
			)
			if err != nil {
				// Log but continue - might be duplicate
				fmt.Printf("Warning: failed to insert goal: %v\n", err)
			}
		}
	}

	return nil
}
