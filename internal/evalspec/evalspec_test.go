package evalspec

import (
	"strings"
	"testing"
)

func TestLint_NoWarningsForAWellFormedSuite(t *testing.T) {
	s := &Suite{
		SkillName: "x",
		Evals: []Eval{
			{
				ID:     "case-one",
				Prompt: "p",
				Assertions: []Assertion{
					{Text: "wraps with %w", Check: Check{ContainsAll: []string{"fmt.Errorf", "%w"}}},
					{Text: "does not panic", Check: Check{NotContains: []string{"panic("}}},
				},
			},
		},
	}
	if got := Lint(s); len(got) != 0 {
		t.Fatalf("Lint = %v, want no warnings for a suite with 2+ assertions and no short substrings", got)
	}
}

// TestLint_DoesNotFlagASingleWellChosenAssertion proves assertion count
// alone is never a warning trigger — a real, previously-flagged case
// (jwt-tenant-scoped-authorization's constant-time-compare check) has
// exactly one assertion with several specific phrasings and no short
// substrings, and should get zero warnings despite having "only 1
// assertion."
func TestLint_DoesNotFlagASingleWellChosenAssertion(t *testing.T) {
	s := &Suite{
		SkillName: "x",
		Evals: []Eval{
			{
				ID:     "lonely-but-specific-case",
				Prompt: "p",
				Assertions: []Assertion{
					{Text: "uses constant-time comparison", Check: Check{ContainsAny: []string{
						"ConstantTimeCompare", "constant-time", "constant time",
					}}},
				},
			},
		},
	}
	if got := Lint(s); len(got) != 0 {
		t.Fatalf("Lint = %v, want no warnings — a single specific assertion is not weak just for being alone", got)
	}
}

func TestLint_FlagsShortSubstrings(t *testing.T) {
	s := &Suite{
		SkillName: "x",
		Evals: []Eval{
			{
				ID:     "short-substring-case",
				Prompt: "p",
				Assertions: []Assertion{
					{Text: "uses an id", Check: Check{ContainsAll: []string{"id"}, ContainsAny: []string{"ok", "fmt.Errorf"}}},
				},
			},
		},
	}
	got := Lint(s)
	if len(got) != 2 {
		t.Fatalf("Lint = %v, want exactly 2 warnings (one per short substring: %q, %q)", got, "id", "ok")
	}
	for _, w := range got {
		if !strings.Contains(w, "short-substring-case") {
			t.Fatalf("warning %q does not identify the case", w)
		}
	}
}

func TestLint_DoesNotFlagShortAcronyms(t *testing.T) {
	s := &Suite{
		SkillName: "x",
		Evals: []Eval{
			{
				ID:     "acronym-case",
				Prompt: "p",
				Assertions: []Assertion{
					{Text: "flags the XSS risk", Check: Check{ContainsAll: []string{"XSS"}}},
					{Text: "flags the missing SRI attribute", Check: Check{ContainsAny: []string{"SRI", "integrity"}}},
				},
			},
		},
	}
	if got := Lint(s); len(got) != 0 {
		t.Fatalf("Lint = %v, want no warnings — all-uppercase acronyms are precise, unlike a short plain word", got)
	}
}

func TestLint_DoesNotFlagLongSubstringsOrRegexChecks(t *testing.T) {
	s := &Suite{
		SkillName: "x",
		Evals: []Eval{
			{
				ID:     "regex-case",
				Prompt: "p",
				Assertions: []Assertion{
					{Text: "uses a bound parameter", Check: Check{MatchesAny: []string{`hotel_id\s*=\s*\$\d`}}},
					{Text: "wraps with %w", Check: Check{ContainsAll: []string{"fmt.Errorf"}}},
				},
			},
		},
	}
	if got := Lint(s); len(got) != 0 {
		t.Fatalf("Lint = %v, want no warnings — matches_any regexes are never length-checked and the other substring is long", got)
	}
}
