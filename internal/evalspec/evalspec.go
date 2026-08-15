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
type Eval struct {
	ID             string      `json:"id"`
	Prompt         string      `json:"prompt"`
	ExpectedOutput string      `json:"expected_output"`
	Assertions     []Assertion `json:"assertions"`
	TimeoutSeconds int         `json:"timeout_seconds,omitempty"` // 0 = use runner default
}

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
