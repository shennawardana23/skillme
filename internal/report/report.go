// Package report writes the workspace layout, timing.json, grading.json,
// benchmark.json, and an HTML summary for a smeval run, following the
// directory convention documented at
// https://agentskills.io/skill-creation/evaluating-skills:
//
//	<skill>-workspace/iteration-N/<eval-id>/with_skill/{outputs/,timing.json,grading.json}
//	<skill>-workspace/iteration-N/<eval-id>/without_skill/... (benchmark mode only)
//	<skill>-workspace/iteration-N/benchmark.json              (benchmark mode only)
//	<skill>-workspace/iteration-N/feedback.json               (stub — a human fills this in)
//	<skill>-workspace/iteration-N/report.html
package report

import (
	"embed"
	"encoding/json"
	"fmt"
	"html/template"
	"os"
	"path/filepath"
	"time"

	"github.com/shennawardana23/skillme/internal/engine"
	"github.com/shennawardana23/skillme/internal/grading"
)

//go:embed assets/report.html.tmpl
var templateFS embed.FS

// Timing is the timing.json shape: tokens and duration for one run.
type Timing struct {
	TotalTokens int64 `json:"total_tokens"`
	DurationMS  int64 `json:"duration_ms"`
}

// RunOutcome is what one configuration (with_skill or without_skill) of one
// eval case produced.
type RunOutcome struct {
	EvalID        string
	Configuration string // "with_skill" or "without_skill"
	EngineErr     string // non-empty if the engine failed both attempts (case is ERROR, not FAIL)
	Result        *engine.Result
	Grading       *grading.Grading
}

// Dir returns this outcome's directory under the iteration directory.
func (o RunOutcome) Dir(iterationDir string) string {
	return filepath.Join(iterationDir, o.EvalID, o.Configuration)
}

// Write persists outputs/, timing.json, and grading.json for one outcome.
func (o RunOutcome) Write(iterationDir string) error {
	dir := o.Dir(iterationDir)
	if err := os.MkdirAll(filepath.Join(dir, "outputs"), 0o755); err != nil {
		return fmt.Errorf("create outputs dir: %w", err)
	}

	if o.EngineErr != "" {
		return os.WriteFile(filepath.Join(dir, "outputs", "engine-error.txt"), []byte(o.EngineErr), 0o644)
	}

	if err := os.WriteFile(filepath.Join(dir, "outputs", "response.md"), []byte(o.Result.FinalMessage), 0o644); err != nil {
		return fmt.Errorf("write response: %w", err)
	}

	timing := Timing{TotalTokens: o.Result.TotalTokens, DurationMS: o.Result.DurationMS}
	if err := writeJSON(filepath.Join(dir, "timing.json"), timing); err != nil {
		return err
	}
	return writeJSON(filepath.Join(dir, "grading.json"), o.Grading)
}

// BenchmarkSide aggregates the runs in one configuration (with_skill or
// without_skill) across all cases in an iteration.
type BenchmarkSide struct {
	PassRateMean    float64 `json:"pass_rate_mean"`
	TimeSecondsMean float64 `json:"time_seconds_mean"`
	TokensMean      float64 `json:"tokens_mean"`
}

// Benchmark is benchmark.json — only written when the run was invoked with
// --benchmark, since a delta requires both configurations to have actually
// run; a with_skill-only run must never emit a fabricated delta.
type Benchmark struct {
	WithSkill    BenchmarkSide `json:"with_skill"`
	WithoutSkill BenchmarkSide `json:"without_skill"`
	Delta        BenchmarkSide `json:"delta"`
}

// Aggregate computes a BenchmarkSide's means from the runs in one configuration.
func Aggregate(outcomes []RunOutcome) BenchmarkSide {
	var side BenchmarkSide
	n := 0
	for _, o := range outcomes {
		if o.EngineErr != "" || o.Result == nil || o.Grading == nil {
			continue
		}
		side.PassRateMean += o.Grading.Summary.PassRate
		side.TimeSecondsMean += float64(o.Result.DurationMS) / 1000.0
		side.TokensMean += float64(o.Result.TotalTokens)
		n++
	}
	if n > 0 {
		side.PassRateMean /= float64(n)
		side.TimeSecondsMean /= float64(n)
		side.TokensMean /= float64(n)
	}
	return side
}

// WriteBenchmark computes and writes benchmark.json from with_skill and
// without_skill outcome sets.
func WriteBenchmark(iterationDir string, withSkill, withoutSkill []RunOutcome) error {
	ws := Aggregate(withSkill)
	wo := Aggregate(withoutSkill)
	b := Benchmark{
		WithSkill:    ws,
		WithoutSkill: wo,
		Delta: BenchmarkSide{
			PassRateMean:    ws.PassRateMean - wo.PassRateMean,
			TimeSecondsMean: ws.TimeSecondsMean - wo.TimeSecondsMean,
			TokensMean:      ws.TokensMean - wo.TokensMean,
		},
	}
	return writeJSON(filepath.Join(iterationDir, "benchmark.json"), b)
}

// WriteFeedbackStub scaffolds feedback.json — a place for a human (or a
// separate qualitative-review pass) to record what deterministic assertions
// structurally cannot check: prose/report quality, whether an output is
// technically correct but misses the point, "does this feel right." The
// spec (https://agentskills.io/skill-creation/evaluating-skills#reviewing-results-with-a-human)
// treats this as a required pillar of the loop, not an optional extra —
// assertion grading only checks what someone thought to write an assertion
// for. Every eval ID gets an empty-string entry ready to fill in; empty
// means "reviewed, looked fine," not "not yet reviewed," so leave an entry
// empty only once you've actually looked at that case's output. Never
// overwrites an existing feedback.json — a fresh iteration directory starts
// empty by construction (NextIterationDir), so there's nothing to preserve
// on a normal run, but this stays safe if it's ever called twice.
func WriteFeedbackStub(iterationDir string, evalIDs []string) error {
	path := filepath.Join(iterationDir, "feedback.json")
	if _, err := os.Stat(path); err == nil {
		return nil
	}
	feedback := make(map[string]string, len(evalIDs))
	for _, id := range evalIDs {
		feedback[id] = ""
	}
	return writeJSON(path, feedback)
}

// htmlData is what report.html.tmpl renders.
type htmlData struct {
	SkillName   string
	GeneratedAt string
	Benchmark   bool
	Cases       []htmlCase
	TotalCases  int
	PassedCases int
}

type htmlCase struct {
	ID           string
	ModelUsed    string
	FallbackUsed bool
	Error        string
	Passed       int
	Total        int
	PassRate     float64
	Assertions   []grading.AssertionResult
	WithoutSkill *htmlCaseSide // nil unless benchmark mode
}

type htmlCaseSide struct {
	Passed   int
	Total    int
	PassRate float64
}

// WriteHTML renders the iteration's report.html from the with_skill (and,
// in benchmark mode, without_skill) outcomes.
func WriteHTML(iterationDir, skillName string, withSkill, withoutSkill []RunOutcome, benchmarkMode bool) error {
	withoutByID := make(map[string]RunOutcome, len(withoutSkill))
	for _, o := range withoutSkill {
		withoutByID[o.EvalID] = o
	}

	data := htmlData{
		SkillName:   skillName,
		GeneratedAt: time.Now().UTC().Format(time.RFC3339),
		Benchmark:   benchmarkMode,
	}
	for _, o := range withSkill {
		hc := htmlCase{ID: o.EvalID}
		if o.EngineErr != "" {
			hc.Error = o.EngineErr
		} else {
			hc.ModelUsed = o.Result.ModelUsed
			hc.FallbackUsed = o.Result.FallbackUsed
			hc.Passed = o.Grading.Summary.Passed
			hc.Total = o.Grading.Summary.Total
			hc.PassRate = o.Grading.Summary.PassRate
			hc.Assertions = o.Grading.AssertionResults
			if hc.Passed == hc.Total && hc.Total > 0 {
				data.PassedCases++
			}
		}
		if benchmarkMode {
			if wo, ok := withoutByID[o.EvalID]; ok && wo.Grading != nil {
				hc.WithoutSkill = &htmlCaseSide{
					Passed:   wo.Grading.Summary.Passed,
					Total:    wo.Grading.Summary.Total,
					PassRate: wo.Grading.Summary.PassRate,
				}
			}
		}
		data.Cases = append(data.Cases, hc)
		data.TotalCases++
	}

	tmpl, err := template.ParseFS(templateFS, "assets/report.html.tmpl")
	if err != nil {
		return fmt.Errorf("parse report template: %w", err)
	}

	// Render to a temp file and rename into place so a mid-render template
	// error never leaves a truncated report.html behind.
	final := filepath.Join(iterationDir, "report.html")
	tmpFile := final + ".tmp"
	f, err := os.Create(tmpFile)
	if err != nil {
		return fmt.Errorf("create %s: %w", tmpFile, err)
	}
	if err := tmpl.Execute(f, data); err != nil {
		f.Close()
		os.Remove(tmpFile)
		return fmt.Errorf("render report.html: %w", err)
	}
	if err := f.Close(); err != nil {
		os.Remove(tmpFile)
		return fmt.Errorf("close %s: %w", tmpFile, err)
	}
	return os.Rename(tmpFile, final)
}

func writeJSON(path string, v any) error {
	data, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal %s: %w", path, err)
	}
	if err := os.WriteFile(path, data, 0o644); err != nil {
		return fmt.Errorf("write %s: %w", path, err)
	}
	return nil
}

// NextIterationDir picks the next unused iteration-N directory under
// outputDir, starting from 1.
func NextIterationDir(outputDir string) (string, error) {
	if err := os.MkdirAll(outputDir, 0o755); err != nil {
		return "", fmt.Errorf("create %s: %w", outputDir, err)
	}
	n := 1
	for {
		dir := filepath.Join(outputDir, fmt.Sprintf("iteration-%d", n))
		if _, err := os.Stat(dir); os.IsNotExist(err) {
			if err := os.MkdirAll(dir, 0o755); err != nil {
				return "", fmt.Errorf("create %s: %w", dir, err)
			}
			return dir, nil
		}
		n++
	}
}
