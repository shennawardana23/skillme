// Package engine invokes the local `claude` CLI headlessly to run one eval
// prompt and reports what happened.
//
// Provider/model fallback happens in two layers, and the split is
// deliberate:
//
//  1. Native, in-process: --fallback-model is passed straight through to
//     the claude CLI, which already retries within the same invocation when
//     the primary model is overloaded or unavailable (verified empirically:
//     an unrecognized --model name still returns is_error:false, exit 0,
//     when --fallback-model is set — the CLI silently recovers). Do not
//     reimplement this; it already works.
//  2. Outer, cross-process: engaged only when the entire invocation fails
//     in a way the native layer cannot see or recover from — the process
//     exits non-zero, the context deadline is exceeded, or stdout is not
//     parseable JSON. On any of these, Run makes exactly one further
//     attempt in a fresh process using FallbackModel as the new primary
//     model (no further chaining). A run that completes and returns
//     is_error:false is never retried for content reasons — assertion
//     failures are a grading outcome, not an engine failure, and must not
//     trigger a fallback.
package engine

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

// ClaudeBin is the executable invoked for every attempt. Tests override
// this to point at a fake script so the outer-fallback branch can be proven
// deterministically without depending on live API failures.
var ClaudeBin = "claude"

// Options configures a single graded invocation (which may internally make
// up to two process attempts — primary, then fallback).
type Options struct {
	Prompt         string        // the eval's prompt, passed as the CLI's positional argument
	PrimaryModel   string        // e.g. "sonnet"
	FallbackModel  string        // e.g. "opus"; empty disables both fallback layers
	PluginDir      string        // when non-empty, the skill-under-test's harness plugin dir (with_skill run); empty = baseline (without_skill)
	WorkDir        string        // working directory for the invocation; the agent may read/write files here
	PermissionMode string        // default "bypassPermissions" if empty
	Timeout        time.Duration // per attempt; default 180s if zero
}

// Result is what a graded invocation produced, for both grading and the
// timing.json/report the workspace records.
type Result struct {
	FinalMessage string
	DurationMS   int64
	TotalTokens  int64
	CostUSD      float64
	ModelUsed    string
	FallbackUsed bool // true only when the OUTER (cross-process) layer engaged
	Attempts     int
}

// claudeJSON is the subset of `claude --output-format json`'s result object
// this package reads. Field names and shapes were captured empirically from
// a live invocation, not guessed from documentation.
type claudeJSON struct {
	IsError        bool    `json:"is_error"`
	Result         string  `json:"result"`
	DurationMS     int64   `json:"duration_ms"`
	TotalCostUSD   float64 `json:"total_cost_usd"`
	APIErrorStatus *int    `json:"api_error_status"`
	TerminalReason string  `json:"terminal_reason"`
	Usage          struct {
		InputTokens            int64 `json:"input_tokens"`
		OutputTokens           int64 `json:"output_tokens"`
		CacheCreationInputToks int64 `json:"cache_creation_input_tokens"`
		CacheReadInputTokens   int64 `json:"cache_read_input_tokens"`
	} `json:"usage"`
}

// Run executes the eval prompt and returns the graded result. A non-nil
// error means both the primary and (if configured) fallback attempts
// failed at the engine level — the case is an ERROR, not a graded FAIL.
func Run(ctx context.Context, opts Options) (*Result, error) {
	if opts.PermissionMode == "" {
		opts.PermissionMode = "bypassPermissions"
	}
	if opts.Timeout == 0 {
		opts.Timeout = 180 * time.Second
	}

	res, err := attempt(ctx, opts, opts.PrimaryModel)
	if err == nil {
		res.ModelUsed = opts.PrimaryModel
		res.Attempts = 1
		return res, nil
	}
	if opts.FallbackModel == "" {
		return nil, fmt.Errorf("primary model %q failed and no fallback model configured: %w", opts.PrimaryModel, err)
	}

	fallbackRes, fallbackErr := attempt(ctx, opts, opts.FallbackModel)
	if fallbackErr != nil {
		return nil, fmt.Errorf("primary model %q failed (%v); fallback model %q also failed: %w", opts.PrimaryModel, err, opts.FallbackModel, fallbackErr)
	}
	fallbackRes.ModelUsed = opts.FallbackModel
	fallbackRes.FallbackUsed = true
	fallbackRes.Attempts = 2
	return fallbackRes, nil
}

// attempt runs one process invocation with the given model and classifies
// the outcome. The returned error is non-nil exactly when this attempt
// should count as an engine-level failure (triggering the caller's
// fallback), never for a well-formed is_error:false response regardless of
// what the assertions will later say about its content.
func attempt(ctx context.Context, opts Options, model string) (*Result, error) {
	attemptCtx, cancel := context.WithTimeout(ctx, opts.Timeout)
	defer cancel()

	args := []string{
		"-p", opts.Prompt,
		"--output-format", "json",
		"--permission-mode", opts.PermissionMode,
		"--model", model,
		// Exclude "user" (the tester's own global ~/.claude/CLAUDE.md and
		// settings) so a case only ever reflects the skill under test (or,
		// for a without_skill baseline, the total absence of one) — not
		// whatever the machine running the eval happens to have configured
		// globally. Without this, a benchmark run on a machine whose global
		// CLAUDE.md already documents the same fact a skill teaches shows a
		// false "no difference" result, because both conditions silently
		// inherit that fact from outside the harness.
		"--setting-sources", "project,local",
	}
	if opts.FallbackModel != "" && opts.FallbackModel != model {
		args = append(args, "--fallback-model", opts.FallbackModel)
	}
	if opts.PluginDir != "" {
		args = append(args, "--plugin-dir", opts.PluginDir)
	}

	cmd := exec.CommandContext(attemptCtx, ClaudeBin, args...)
	cmd.Dir = opts.WorkDir
	// cmd.Dir changes the OS-level working directory but Go does not touch
	// the inherited PWD env var to match — it stays whatever the smeval
	// process's own PWD was (almost always this repo's root, since that's
	// where `smeval run` gets invoked from). Node.js tools sometimes trust
	// process.env.PWD over process.cwd() for path resolution, which let a
	// file-writing case ("scaffold a new skill under skills/<name>/")
	// silently escape the isolated per-case workspace and write real files
	// into this actual repo instead — confirmed by an untracked directory
	// appearing in `git status` after a live run, matching the case's
	// prompt exactly, while the intended workspace stayed empty. Overriding
	// PWD here closes that specific escape route.
	absWorkDir, err := filepath.Abs(opts.WorkDir)
	if err != nil {
		absWorkDir = opts.WorkDir
	}
	env := os.Environ()
	for i, kv := range env {
		if strings.HasPrefix(kv, "PWD=") {
			env[i] = "PWD=" + absWorkDir
		}
	}
	cmd.Env = env
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	runErr := cmd.Run()
	if attemptCtx.Err() != nil {
		return nil, fmt.Errorf("timed out after %s: %w", opts.Timeout, attemptCtx.Err())
	}
	if runErr != nil {
		return nil, fmt.Errorf("process exited with error: %w (stderr: %s)", runErr, truncate(stderr.String(), 500))
	}

	var cj claudeJSON
	if err := json.Unmarshal(stdout.Bytes(), &cj); err != nil {
		return nil, fmt.Errorf("could not parse JSON output: %w (stdout: %s)", err, truncate(stdout.String(), 500))
	}
	if cj.IsError || cj.APIErrorStatus != nil || cj.TerminalReason == "api_error" {
		return nil, fmt.Errorf("engine reported an error (is_error=%v, api_error_status=%v, terminal_reason=%q): %s",
			cj.IsError, deref(cj.APIErrorStatus), cj.TerminalReason, cj.Result)
	}

	total := cj.Usage.InputTokens + cj.Usage.OutputTokens + cj.Usage.CacheCreationInputToks + cj.Usage.CacheReadInputTokens
	return &Result{
		FinalMessage: cj.Result,
		DurationMS:   cj.DurationMS,
		TotalTokens:  total,
		CostUSD:      cj.TotalCostUSD,
	}, nil
}

func truncate(s string, n int) string {
	if len(s) <= n {
		return s
	}
	return s[:n] + "...(truncated)"
}

func deref(p *int) any {
	if p == nil {
		return nil
	}
	return *p
}
