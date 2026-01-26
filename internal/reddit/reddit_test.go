package reddit

import (
	"fmt"
	"testing"
)

func TestParseGoalFromTitle(t *testing.T) {
	tests := []struct {
		title    string
		url      string
		wantHome string
		wantAway string
		wantHs   int
		wantAs   int
		wantGoal string
	}{
		{
			title:    "Newcastle United 0-1 Arsenal - Dennis Bergkamp 11'",
			url:      "https://www.youtube.com/watch?v=IicmCu47pMo",
			wantHome: "Newcastle United",
			wantAway: "Arsenal",
			wantHs:   0,
			wantAs:   1,
			wantGoal: "Dennis Bergkamp",
		},
		{
			title:    "Arsenal [1]-0 Leeds United - Thierry Henry 78'",
			url:      "https://www.youtube.com/watch?v=_bNBN9XlTK0",
			wantHome: "Arsenal",
			wantAway: "Leeds United",
			wantHs:   1,
			wantAs:   0,
			wantGoal: "Thierry Henry",
		},
		{
			title:    "Cheltenham 0-[2] Notts County - Gabriel 27'",
			url:      "https://www.youtube.com/watch?v=example",
			wantHome: "Cheltenham",
			wantAway: "Notts County",
			wantHs:   0,
			wantAs:   2,
			wantGoal: "Gabriel",
		},
		{
			title:    "Cheltenham 0-[2] Notts County - Tyrese Hall 27'",
			url:      "https://www.youtube.com/watch?v=example",
			wantHome: "Cheltenham",
			wantAway: "Notts County",
			wantHs:   0,
			wantAs:   2,
			wantGoal: "Tyrese Hall",
		},
	}

	for _, tt := range tests {
		goal, err := ParseGoalFromTitle(tt.title, tt.url, "/r/soccer/comments/example")
		if err != nil {
			t.Errorf("ParseGoalFromTitle(%q) error = %v", tt.title, err)
			continue
		}

		if goal.HomeTeam != tt.wantHome {
			t.Errorf("ParseGoalFromTitle(%q) HomeTeam = %q, want %q", tt.title, goal.HomeTeam, tt.wantHome)
		}
		if goal.AwayTeam != tt.wantAway {
			t.Errorf("ParseGoalFromTitle(%q) AwayTeam = %q, want %q", tt.title, goal.AwayTeam, tt.wantAway)
		}
		if goal.HomeScore != tt.wantHs {
			t.Errorf("ParseGoalFromTitle(%q) HomeScore = %d, want %d", tt.title, goal.HomeScore, tt.wantHs)
		}
		if goal.AwayScore != tt.wantAs {
			t.Errorf("ParseGoalFromTitle(%q) AwayScore = %d, want %d", tt.title, goal.AwayScore, tt.wantAs)
		}
		if goal.Goalscorer != tt.wantGoal {
			t.Errorf("ParseGoalFromTitle(%q) Goalscorer = %q, want %q", tt.title, goal.Goalscorer, tt.wantGoal)
		}
	}
}

func ExampleParseGoalFromTitle() {
	title := "Arsenal [1]-0 Leeds United - Thierry Henry 78'"
	goal, _ := ParseGoalFromTitle(title, "https://www.youtube.com/watch?v=_bNBN9XlTK0", "/r/soccer/comments/example")
	fmt.Printf("Team: %s vs %s\n", goal.HomeTeam, goal.AwayTeam)
	fmt.Printf("Score: %d-%d\n", goal.HomeScore, goal.AwayScore)
	fmt.Printf("Goal: %s in minute %s\n", goal.Goalscorer, goal.Minute)
}
