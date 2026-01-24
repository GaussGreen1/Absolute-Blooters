package reddit

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	"blooters/internal/models"
)

const RedditAPIURL = "https://www.reddit.com/r/soccer/new.json?limit=50"

type RedditResponse struct {
	Kind string `json:"kind"`
	Data struct {
		Children []struct {
			Kind string `json:"kind"`
			Data struct {
				Title     string  `json:"title"`
				URL       string  `json:"url"`
				Permalink string  `json:"permalink"`
				Created   float64 `json:"created_utc"`
				FlairText string  `json:"link_flair_text"`
			} `json:"data"`
		} `json:"children"`
	} `json:"data"`
}

type Comment struct {
	Kind string `json:"kind"`
	Data struct {
		Author    string    `json:"author"`
		Body      string    `json:"body"`
		Permalink string    `json:"permalink"`
		Replies   *Comments `json:"replies"`
		ID        string    `json:"id"`
	} `json:"data"`
}

type Comments struct {
	Kind string `json:"kind"`
	Data struct {
		Children []Comment `json:"children"`
	} `json:"data"`
}

type PostAndComments struct {
	Post     interface{} `json:"-"` // We don't need the post
	Comments Comments    `json:"1"`
}

// FetchGoals fetches goals from Reddit API and parses them
func FetchGoals() ([]models.Goal, error) {
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	req, err := http.NewRequest("GET", RedditAPIURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("User-Agent", "blooters/1.0 (goal scraper)")

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch from Reddit: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("reddit API returned status %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var redditResp RedditResponse
	if err := json.Unmarshal(body, &redditResp); err != nil {
		return nil, fmt.Errorf("failed to parse JSON response: %w", err)
	}

	var goals []models.Goal
	for _, child := range redditResp.Data.Children {
		if child.Kind != "t3" || child.Data.FlairText != "Media" {
			continue
		}

		// Parse the title to extract goal information
		goal, err := ParseGoalFromTitle(child.Data.Title, child.Data.URL, child.Data.Permalink)
		if err != nil {
			// Skip posts that don't match goal format
			continue
		}

		// Fetch mirrors link
		mirrorsLink, err := getMirrorsLink(goal.RedditURL)
		if err != nil {
			// Log error but don't skip the goal
			fmt.Printf("Warning: failed to get mirrors for %s: %v\n", goal.RedditURL, err)
			mirrorsLink = ""
		}
		goal.Mirrors = mirrorsLink

		// Sleep to avoid rate limiting
		time.Sleep(200 * time.Millisecond)

		goals = append(goals, goal)
	}

	return goals, nil
}

func ParseGoalFromTitle(title, url, permalink string) (models.Goal, error) {
	goal := models.Goal{
		Url:         url,
		RedditURL:   "https://www.reddit.com" + permalink,
		Description: title,
	}

	normalized := strings.NewReplacer("–", "-", "—", "-", "  ", " ").Replace(title)

	//I hate regexes but here we go:
	// 1: home team
	// 2: home score
	// 3: away score
	// 4: away team
	scorePattern := regexp.MustCompile(`(?m)^(.+?)\s+\[?(\d+)\]?\s*-\s*\[?(\d+)\]?\s+(.+?)\s*-\s*(.+)$`)
	matches := scorePattern.FindStringSubmatch(normalized)
	if len(matches) != 6 {
		// I really hate regexes. Fallback: try without the final dash (in case title has no dash before scorer)
		scorePatternAlt := regexp.MustCompile(`(?m)^(.+?)\s+\[?(\d+)\]?\s*-\s*\[?(\d+)\]?\s+(.+)$`)
		matches = scorePatternAlt.FindStringSubmatch(normalized)
		if len(matches) != 5 {
			return goal, fmt.Errorf("could not parse teams and score")
		}
		// In this fallback, scorer part is empty; we will parse manually later
		matches = append(matches, "")
	}

	goal.HomeTeam = strings.TrimSpace(matches[1])
	goal.AwayTeam = strings.TrimSpace(matches[4])

	homeScore, _ := strconv.Atoi(matches[2])
	awayScore, _ := strconv.Atoi(matches[3])
	goal.HomeScore = homeScore
	goal.AwayScore = awayScore

	scorerPart := strings.TrimSpace(matches[5])
	if scorerPart == "" {
		// TRY to extract scorer info from remaining string
		idx := strings.LastIndex(normalized, "-")
		if idx >= 0 && idx+1 < len(normalized) {
			scorerPart = strings.TrimSpace(normalized[idx+1:])
		}
	}

	// Extract minute and scorer
	minutePattern := regexp.MustCompile(`(.+?)\s+(\d+(?:\+\d+)?)[′']?$`)
	mm := minutePattern.FindStringSubmatch(scorerPart)
	if len(mm) == 3 {
		rawScorer := strings.TrimSpace(mm[1])
		min := strings.TrimSpace(mm[2])

		// Remove penalty/own goal markers
		lower := strings.ToLower(rawScorer)
		removeMods := []string{"(pen)", "pen", "penalty", "(og)", "own goal", "og"}
		for _, m := range removeMods {
			lower = strings.TrimSpace(strings.TrimSuffix(lower, " "+m))
			lower = strings.TrimSpace(strings.TrimSuffix(lower, " ("+m+")"))
		}

		goal.Goalscorer = strings.Title(lower)
		goal.Minute = min
	} else {
		// fallback: last field is minute, everything before is the scorer name
		fields := strings.Fields(scorerPart)
		if len(fields) >= 2 {
			// Last field should be the minute (possibly with ' or ′)
			minute := fields[len(fields)-1]
			minute = strings.Trim(minute, "′'")
			goal.Minute = minute

			// Everything else is the goalscorer name
			goal.Goalscorer = strings.Join(fields[:len(fields)-1], " ")
		} else if len(fields) == 1 {
			// Only one field, assume it's the scorer
			goal.Goalscorer = fields[0]
		}
	}

	// Home or away?
	goal.Away = goal.AwayScore > homeScore

	return goal, nil
}

func getMirrorsLink(postURL string) (string, error) {
	if !strings.HasSuffix(postURL, "/") {
		postURL += "/"
	}
	commentsURL := postURL + ".json"

	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	req, err := http.NewRequest("GET", commentsURL, nil)
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("User-Agent", "blooters/1.0 (goal scraper)")

	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to fetch comments: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("comments API returned status %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read comments body: %w", err)
	}

	var data []interface{}
	if err := json.Unmarshal(body, &data); err != nil {
		return "", fmt.Errorf("failed to parse comments JSON: %w", err)
	}

	if len(data) < 2 {
		return "", fmt.Errorf("unexpected JSON structure")
	}

	commentsData, ok := data[1].(map[string]interface{})
	if !ok {
		return "", fmt.Errorf("comments data not found")
	}

	children, ok := commentsData["data"].(map[string]interface{})["children"].([]interface{})
	if !ok {
		return "", fmt.Errorf("comments children not found")
	}

	for _, child := range children {
		comment, ok := child.(map[string]interface{})
		if !ok || comment["kind"].(string) != "t1" {
			continue
		}

		commentData, ok := comment["data"].(map[string]interface{})
		if !ok {
			continue
		}

		author, ok := commentData["author"].(string)
		if !ok || author != "AutoModerator" {
			continue
		}

		body, ok := commentData["body"].(string)
		if !ok || !strings.Contains(body, "Mirrors / Alternative Angles") {
			continue
		}

		permalink, ok := commentData["permalink"].(string)
		if !ok {
			continue
		}

		return "https://www.reddit.com" + permalink, nil
	}

	return "", fmt.Errorf("mirrors comment not found")
}
