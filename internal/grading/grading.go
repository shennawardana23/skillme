// Package grading evaluates evalspec.Assertion checks against an engine
// response and produces the grading.json shape documented at
// https://agentskills.io/skill-creation/evaluating-skills (assertion_results
// + summary), computed deterministically rather than by an LLM judge.
package grading

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/shennawardana23/skillme/internal/evalspec"
)

// AssertionResult is one graded assertion, with concrete evidence for a
// PASS or FAIL rather than a bare boolean — mirrors the reference
// methodology's principle that a PASS must cite what was actually found.
type AssertionResult struct {
	Text     string `json:"text"`
	Passed   bool   `json:"passed"`
	Evidence string `json:"evidence"`
}

// Summary aggregates a case's assertion results.
type Summary struct {
	Passed   int     `json:"passed"`
	Failed   int     `json:"failed"`
	Total    int     `json:"total"`
	PassRate float64 `json:"pass_rate"`
}

// Grading is the full grading.json document for one run (one case, one
// configuration — with_skill or without_skill).
type Grading struct {
	AssertionResults []AssertionResult `json:"assertion_results"`
	Summary          Summary           `json:"summary"`
}

// Grade evaluates every assertion in an eval against the engine's final
// output text and, for files_exist/file_contains checks, the case's real
// workspace directory (workDir may be empty when a case has no file
// checks).
func Grade(assertions []evalspec.Assertion, output, workDir string) Grading {
	results := make([]AssertionResult, 0, len(assertions))
	passed := 0
	for _, a := range assertions {
		ok, evidence := evaluateCheck(a.Check, output, workDir)
		results = append(results, AssertionResult{Text: a.Text, Passed: ok, Evidence: evidence})
		if ok {
			passed++
		}
	}
	total := len(results)
	rate := 0.0
	if total > 0 {
		rate = float64(passed) / float64(total)
	}
	return Grading{
		AssertionResults: results,
		Summary: Summary{
			Passed:   passed,
			Failed:   total - passed,
			Total:    total,
			PassRate: rate,
		},
	}
}

// evaluateCheck runs every set condition on a Check and requires all of
// them to hold; it returns human-readable evidence either way so a FAIL is
// as legible as a PASS.
//
// contains_all/contains_any/not_contains match case-insensitively. Live
// grading runs showed models phrase the same correct answer with different
// capitalization ("Payment gateway" vs. an eval expecting "payment") —
// case-sensitive matching produced false FAILs on genuinely correct
// responses. matches_any/not_matches stay case-sensitive by default since
// they're regexps; authors add an inline (?i) flag when a pattern should
// ignore case.
func evaluateCheck(c evalspec.Check, output, workDir string) (bool, string) {
	var reasons []string
	lowerOutput := strings.ToLower(output)

	if len(c.ContainsAll) > 0 {
		var missing []string
		for _, s := range c.ContainsAll {
			if !strings.Contains(lowerOutput, strings.ToLower(s)) {
				missing = append(missing, s)
			}
		}
		if len(missing) > 0 {
			return false, fmt.Sprintf("missing required substring(s): %s", strings.Join(missing, ", "))
		}
		reasons = append(reasons, fmt.Sprintf("contains all of: %s", strings.Join(c.ContainsAll, ", ")))
	}

	if len(c.ContainsAny) > 0 {
		found := ""
		for _, s := range c.ContainsAny {
			if strings.Contains(lowerOutput, strings.ToLower(s)) {
				found = s
				break
			}
		}
		if found == "" {
			return false, fmt.Sprintf("none of the required substrings found: %s", strings.Join(c.ContainsAny, ", "))
		}
		reasons = append(reasons, fmt.Sprintf("found %q", found))
	}

	if len(c.NotContains) > 0 {
		for _, s := range c.NotContains {
			if strings.Contains(lowerOutput, strings.ToLower(s)) {
				return false, fmt.Sprintf("found disallowed substring: %q", s)
			}
		}
		reasons = append(reasons, "none of the disallowed substrings present")
	}

	if len(c.MatchesAny) > 0 {
		matched := ""
		for _, pat := range c.MatchesAny {
			re, err := regexp.Compile(pat)
			if err != nil {
				return false, fmt.Sprintf("invalid regexp %q: %v", pat, err)
			}
			if re.MatchString(output) {
				matched = pat
				break
			}
		}
		if matched == "" {
			return false, fmt.Sprintf("none of the required patterns matched: %s", strings.Join(c.MatchesAny, ", "))
		}
		reasons = append(reasons, fmt.Sprintf("matched pattern %q", matched))
	}

	if len(c.NotMatches) > 0 {
		for _, pat := range c.NotMatches {
			re, err := regexp.Compile(pat)
			if err != nil {
				return false, fmt.Sprintf("invalid regexp %q: %v", pat, err)
			}
			if re.MatchString(output) {
				return false, fmt.Sprintf("disallowed pattern matched: %q", pat)
			}
		}
		reasons = append(reasons, "none of the disallowed patterns matched")
	}

	if len(c.FilesExist) > 0 {
		var missing []string
		for _, rel := range c.FilesExist {
			if _, err := os.Stat(filepath.Join(workDir, rel)); err != nil {
				missing = append(missing, rel)
			}
		}
		if len(missing) > 0 {
			return false, fmt.Sprintf("missing expected file(s) in workspace: %s", strings.Join(missing, ", "))
		}
		reasons = append(reasons, fmt.Sprintf("file(s) exist: %s", strings.Join(c.FilesExist, ", ")))
	}

	if len(c.FileContains) > 0 {
		for _, fc := range c.FileContains {
			data, err := os.ReadFile(filepath.Join(workDir, fc.Path))
			if err != nil {
				return false, fmt.Sprintf("could not read %s: %v", fc.Path, err)
			}
			if !strings.Contains(string(data), fc.Contains) {
				return false, fmt.Sprintf("%s does not contain %q", fc.Path, fc.Contains)
			}
		}
		reasons = append(reasons, "file content check(s) passed")
	}

	return true, strings.Join(reasons, "; ")
}
