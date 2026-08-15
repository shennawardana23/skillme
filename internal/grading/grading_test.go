package grading

import (
	"testing"

	"github.com/shennawardana23/skillme/internal/evalspec"
)

func TestGrade(t *testing.T) {
	assertions := []evalspec.Assertion{
		{Text: "wraps with %w", Check: evalspec.Check{ContainsAll: []string{"%w", "fmt.Errorf"}}},
		{Text: "does not panic", Check: evalspec.Check{NotContains: []string{"panic("}}},
		{Text: "uses a bound parameter", Check: evalspec.Check{MatchesAny: []string{`hotel_id\s*=\s*\$\d`}}},
	}

	tests := []struct {
		name       string
		output     string
		wantPassed int
	}{
		{"all pass", `fmt.Errorf("load: %w", err); q := "hotel_id = $1"`, 3},
		{"missing wrap, has panic", `panic(err)`, 0},
		{"partial", `fmt.Errorf("load: %w", err)`, 2},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := Grade(assertions, tt.output, "")
			if g.Summary.Passed != tt.wantPassed {
				t.Errorf("Passed = %d, want %d (results: %+v)", g.Summary.Passed, tt.wantPassed, g.AssertionResults)
			}
			if g.Summary.Total != len(assertions) {
				t.Errorf("Total = %d, want %d", g.Summary.Total, len(assertions))
			}
		})
	}
}

// TestGrade_CaseInsensitive proves contains_all/contains_any/not_contains
// match regardless of capitalization — a live smeval run against
// test-driven-development's mock-vs-real-judgment eval false-FAILed
// because the model wrote "Payment gateway" while the assertion checked
// for lowercase "payment".
func TestGrade_CaseInsensitive(t *testing.T) {
	assertions := []evalspec.Assertion{
		{Text: "mentions payment", Check: evalspec.Check{ContainsAll: []string{"payment"}}},
		{Text: "mentions db or database", Check: evalspec.Check{ContainsAny: []string{"database"}}},
		{Text: "no todo left", Check: evalspec.Check{NotContains: []string{"todo"}}},
	}

	g := Grade(assertions, "Payment gateway: mock it. Postgres DATABASE: use the real thing.", "")
	if g.Summary.Passed != 3 {
		t.Fatalf("Passed = %d, want 3 — payment/database match case-insensitively and todo is absent (results: %+v)", g.Summary.Passed, g.AssertionResults)
	}

	g = Grade(assertions, "Payment gateway: mock it. Postgres DATABASE: use the real thing. TODO: revisit.", "")
	if g.Summary.Passed != 2 {
		t.Fatalf("Passed = %d, want 2 once TODO/todo collides case-insensitively with not_contains (results: %+v)", g.Summary.Passed, g.AssertionResults)
	}
}
