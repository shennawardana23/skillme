// Package evalspec loads and validates a skill's evals/evals.json — the
// declarative case format this repository's eval runner (smeval) grades
// against, modeled on Anthropic's documented Agent Skills evaluation
// methodology (https://agentskills.io/skill-creation/evaluating-skills).
package evalspec

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// FileContainsCheck asserts that a file within the case's workspace
// contains a given substring.
type FileContainsCheck struct {
	Path     string `json:"path"`
	Contains string `json:"contains"`
}

// Check is a deterministic, machine-gradable condition attached to an
// Assertion. Exactly the fields that are set are evaluated; all set fields
// must pass for the assertion to pass. The output_* fields check the
// engine's final response text; the files_* fields check the case's real
// workspace directory (see engine.Options.WorkDir) — use these for cases
// that ask the agent to write files rather than answer inline.
type Check struct {
	ContainsAll  []string            `json:"contains_all,omitempty"`
	ContainsAny  []string            `json:"contains_any,omitempty"`
	NotContains  []string            `json:"not_contains,omitempty"`
	MatchesAny   []string            `json:"matches_any,omitempty"` // Go regexp; at least one must match
	NotMatches   []string            `json:"not_matches,omitempty"` // Go regexp; none may match
	FilesExist   []string            `json:"files_exist,omitempty"` // paths relative to the case workspace
	FileContains []FileContainsCheck `json:"file_contains,omitempty"`
}

// Assertion is one verifiable statement about the output, paired with a
// human-readable description and the deterministic Check that decides
// pass/fail — mirrors the "text" + machine check split in the reference
// methodology's grading.json, but computes the check locally instead of
// asking an LLM judge, keeping grading free and reproducible.
type Assertion struct {
	Text  string `json:"text"`
	Check Check  `json:"check"`
}

// Eval is a single test case: a prompt, a human-readable description of
// success, and the assertions that verify it.
//
// A case with no Fixture is greenfield-feasible: smeval's isolated,
// initially-empty workspace is all it needs, exactly like every case in
// this catalog before this field existed. A case with Fixture set is
// context-dependent: it needs pre-existing files (a small multi-file
// codebase, a config, a partial git history) seeded into that workspace
// before the prompt runs, because the scenario it tests can't be fully
// specified in prompt text alone — e.g. "trace a bug whose symptom and root
// cause are in different files" requires the model to actually go read
// more than one file, which a prompt that hands over all the code inline
// cannot exercise. See catalog-evals/README.md for the real gap this
// closes and debug/evals/evals.json's "traces-bug-across-files" case for a
// worked example.
type Eval struct {
	ID             string      `json:"id"`
	Prompt         string      `json:"prompt"`
	ExpectedOutput string      `json:"expected_output"`
	Assertions     []Assertion `json:"assertions"`
	TimeoutSeconds int         `json:"timeout_seconds,omitempty"` // 0 = use runner default
	Fixture        string      `json:"fixture,omitempty"`         // relative to the suite's evals/fixtures/ directory; empty = greenfield-feasible
}

// IsContextDependent reports whether e needs a seeded fixture rather than
// being testable in a completely fresh, empty container.
func (e Eval) IsContextDependent() bool { return e.Fixture != "" }

// Suite is the evals/evals.json document for one skill.
type Suite struct {
	SkillName string `json:"skill_name"`
	Evals     []Eval `json:"evals"`
}

// Load reads and validates evals/evals.json at path.
func Load(path string) (*Suite, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read %s: %w", path, err)
	}
	var s Suite
	if err := json.Unmarshal(data, &s); err != nil {
		return nil, fmt.Errorf("parse %s: %w", path, err)
	}
	if s.SkillName == "" {
		return nil, fmt.Errorf("%s: skill_name is required", path)
	}
	if len(s.Evals) == 0 {
		return nil, fmt.Errorf("%s: evals must contain at least one case", path)
	}
	seen := make(map[string]bool, len(s.Evals))
	for i, e := range s.Evals {
		if e.ID == "" {
			return nil, fmt.Errorf("%s: evals[%d].id is required", path, i)
		}
		if seen[e.ID] {
			return nil, fmt.Errorf("%s: duplicate eval id %q", path, e.ID)
		}
		seen[e.ID] = true
		if e.Prompt == "" {
			return nil, fmt.Errorf("%s: eval %q: prompt is required", path, e.ID)
		}
		if len(e.Assertions) == 0 {
			return nil, fmt.Errorf("%s: eval %q: at least one assertion is required", path, e.ID)
		}
		for j, a := range e.Assertions {
			if a.Text == "" {
				return nil, fmt.Errorf("%s: eval %q: assertions[%d].text is required", path, e.ID, j)
			}
			c := a.Check
			if len(c.ContainsAll) == 0 && len(c.ContainsAny) == 0 && len(c.NotContains) == 0 &&
				len(c.MatchesAny) == 0 && len(c.NotMatches) == 0 && len(c.FilesExist) == 0 && len(c.FileContains) == 0 {
				return nil, fmt.Errorf("%s: eval %q: assertions[%d] (%q) has no check conditions", path, e.ID, j, a.Text)
			}
		}
	}
	return &s, nil
}

// SkillDir returns the skill directory containing an evals/evals.json path,
// e.g. skills/go-service-idioms/evals/evals.json -> skills/go-service-idioms.
func SkillDir(evalsJSONPath string) string {
	return filepath.Dir(filepath.Dir(evalsJSONPath))
}

// minSubstringLen is the length below which a plain-alphabetic
// contains_all/contains_any string is flagged as too short to reliably
// identify a specific fact — a 1-3 letter word routinely turns up inside
// unrelated words (e.g. "id" inside "avoid", "did", "provide"), so a case
// built on one can pass against a response that never actually states what
// the assertion claims to check. A short string containing anything other
// than letters (a format verb like "%w", a bound parameter like "$1", a
// version number) is exempt — punctuation and digits already make a short
// token distinctive, and this exact "%w" case is used correctly, and
// intentionally, in this repo's own grading_test.go.
const minSubstringLen = 4

// Lint returns non-fatal quality warnings about a suite that Load's hard
// validation deliberately does not reject — a suite with these issues still
// parses and still runs, but is weaker evidence than it looks. Unlike Load,
// Lint never returns an error: every finding here is a suggestion for a
// human to strengthen the eval, not a schema violation.
//
// Deliberately does NOT flag a case for having only one assertion. This
// catalog's own standard (skill-catalog-authoring's "2-3 cases per skill")
// is about case count, never assertion count per case, and a single
// well-chosen contains_any with several specific phrasings — e.g.
// jwt-tenant-scoped-authorization's "ConstantTimeCompare"/"constant-time"/
// "constant time" check — is a solid, specific assertion on its own. An
// earlier version of this lint flagged exactly that kind of case as "weak"
// on assertion count alone; auditing the real catalog against it showed the
// count-based signal was noise, not evidence, so it was removed rather than
// tuned further. Substring specificity (below) is the signal that actually
// correlates with a weak assertion.
func Lint(s *Suite) []string {
	var warnings []string
	for _, e := range s.Evals {
		for _, a := range e.Assertions {
			for _, short := range shortSubstrings(a.Check.ContainsAll, a.Check.ContainsAny) {
				warnings = append(warnings, fmt.Sprintf(
					"eval %q: assertion %q: substring %q is shorter than %d characters and may match unrelated text — use a longer, more specific substring or matches_any with a word-boundary regex (e.g. `\\bword\\b`)",
					e.ID, a.Text, short, minSubstringLen))
			}
		}
	}
	return warnings
}

// shortSubstrings returns every string across the given contains_all/
// contains_any lists that is shorter than minSubstringLen, made up entirely
// of ASCII letters, and not all-uppercase. A short token that also contains
// punctuation or a digit (a format verb, a symbol, a version number) is
// precise enough on its own; so is a short all-uppercase token (an acronym
// like "XSS", "SRI", "ROI") — acronyms are conventionally capitalized
// specifically so they read as one distinct token, not a coincidental
// substring, unlike a short plain word ("id", "key", "log") in normal
// sentence case. Only the latter is flagged.
func shortSubstrings(lists ...[]string) []string {
	var short []string
	for _, list := range lists {
		for _, raw := range list {
			t := strings.TrimSpace(raw)
			if len(t) > 0 && len(t) < minSubstringLen && isAllLetters(t) && !isAllUpper(t) {
				short = append(short, raw)
			}
		}
	}
	return short
}

// isAllLetters reports whether s consists entirely of ASCII letters — used
// to exempt short-but-precise tokens (punctuation or digits included, like
// "%w" or "$1") from the short-substring lint.
func isAllLetters(s string) bool {
	for _, r := range s {
		if !((r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z')) {
			return false
		}
	}
	return true
}

// isAllUpper reports whether s (already confirmed all-letters by
// isAllLetters) is entirely uppercase — the conventional shape of an
// acronym, which this lint exempts from the short-substring warning.
func isAllUpper(s string) bool {
	for _, r := range s {
		if r >= 'a' && r <= 'z' {
			return false
		}
	}
	return true
}
