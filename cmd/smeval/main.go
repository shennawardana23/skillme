// Command smeval is this repository's own eval runner for Claude Code
// skills — no third-party dependency, built to follow Anthropic's
// documented Agent Skills evaluation methodology
// (https://agentskills.io/skill-creation/evaluating-skills): a prompt and
// assertions per case, graded with concrete evidence, results aggregated
// into a workspace of iteration-N/<case>/{with_skill,without_skill}
// directories.
package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/shennawardana23/skillme/internal/engine"
	"github.com/shennawardana23/skillme/internal/evalspec"
	"github.com/shennawardana23/skillme/internal/grading"
	"github.com/shennawardana23/skillme/internal/harness"
	"github.com/shennawardana23/skillme/internal/report"
)

// errCasesFailed signals that runRun completed normally but at least one
// case failed or errored — main exits non-zero for it without printing a
// redundant "smeval: ..." line, since the per-case output already said so.
var errCasesFailed = errors.New("one or more cases failed")

func main() {
	if len(os.Args) < 2 {
		usage()
		os.Exit(2)
	}
	var err error
	switch os.Args[1] {
	case "validate":
		err = runValidate(os.Args[2:])
	case "run":
		err = runRun(os.Args[2:])
	case "-h", "--help", "help":
		usage()
		return
	default:
		usage()
		os.Exit(2)
	}
	if err != nil {
		if !errors.Is(err, errCasesFailed) {
			fmt.Fprintln(os.Stderr, "smeval:", err)
		}
		os.Exit(1)
	}
}

func usage() {
	fmt.Fprintln(os.Stderr, `smeval — this repo's own skill eval runner (no third-party dependency)

Usage:
  smeval validate <skill-dir>
  smeval run <skill-dir> [flags]

Flags for run:
  -primary-model string   Model for the primary attempt (default "sonnet")
  -fallback-model string  Model for the outer fallback attempt; also passed
                          to claude's native --fallback-model (default "opus")
  -timeout duration       Per-attempt timeout (default 3m0s)
  -benchmark              Also run each case without the skill installed and write benchmark.json
  -output-dir string      Workspace root (default "<skill-dir>-workspace", a sibling directory)
  -include string         Only run eval IDs containing this substring`)
}

func evalsPath(skillDir string) string {
	return filepath.Join(skillDir, "evals", "evals.json")
}

// isolateWorkspace makes dir its own empty git repository. dir already
// lives inside this repo's own working tree (skills/<name>-workspace/,
// gitignored but not outside the tree), so any tool that discovers its
// "project root" by walking up to the nearest .git — rather than trusting
// its own process working directory — would otherwise walk straight past
// dir and land on this actual repo's root. That escape was proven live: a
// case whose prompt asked the model to write files "under skills/<name>/"
// wrote them into this repo's real skills/ directory instead of into dir,
// leaving the intended workspace empty. Best-effort: a missing git binary
// or a failed init just leaves dir without this extra guard, it does not
// fail the run.
func isolateWorkspace(dir string) {
	cmd := exec.Command("git", "init", "-q", dir)
	_ = cmd.Run()
}

func runValidate(args []string) error {
	fs := flag.NewFlagSet("validate", flag.ExitOnError)
	fs.Parse(args)
	if fs.NArg() != 1 {
		return fmt.Errorf("usage: smeval validate <skill-dir>")
	}
	skillDir := fs.Arg(0)
	suite, err := evalspec.Load(evalsPath(skillDir))
	if err != nil {
		return err
	}
	dirName := filepath.Base(filepath.Clean(skillDir))
	if suite.SkillName != dirName {
		return fmt.Errorf("%s: skill_name %q does not match directory name %q", evalsPath(skillDir), suite.SkillName, dirName)
	}
	fmt.Printf("✓ %s is valid (%d case(s))\n", evalsPath(skillDir), len(suite.Evals))
	return nil
}

func runRun(args []string) error {
	fs := flag.NewFlagSet("run", flag.ExitOnError)
	primaryModel := fs.String("primary-model", "sonnet", "")
	fallbackModel := fs.String("fallback-model", "opus", "")
	timeout := fs.Duration("timeout", 3*time.Minute, "")
	benchmark := fs.Bool("benchmark", false, "")
	outputDir := fs.String("output-dir", "", "")
	include := fs.String("include", "", "")
	fs.Parse(reorderArgs(args, map[string]bool{"benchmark": true}))
	if fs.NArg() != 1 {
		return fmt.Errorf("usage: smeval run <skill-dir> [flags]")
	}
	skillDir := fs.Arg(0)

	suite, err := evalspec.Load(evalsPath(skillDir))
	if err != nil {
		return err
	}

	if *outputDir == "" {
		*outputDir = strings.TrimRight(filepath.Clean(skillDir), "/") + "-workspace"
	}
	iterationDir, err := report.NextIterationDir(*outputDir)
	if err != nil {
		return err
	}

	harnessDir, err := harness.Build(skillDir)
	if err != nil {
		return fmt.Errorf("build harness plugin: %w", err)
	}
	defer os.RemoveAll(harnessDir)

	ctx := context.Background()
	var withSkill, withoutSkill []report.RunOutcome
	failures := 0
	selected := 0

	for _, ev := range suite.Evals {
		if *include != "" && !strings.Contains(ev.ID, *include) {
			continue
		}
		selected++

		perCaseTimeout := *timeout
		if ev.TimeoutSeconds > 0 {
			perCaseTimeout = time.Duration(ev.TimeoutSeconds) * time.Second
		}

		fmt.Printf("⏳ %s: %s\n", ev.ID, truncate(ev.ExpectedOutput, 70))

		wsWorkDir := filepath.Join(iterationDir, ev.ID, "with_skill", "workspace")
		if err := os.MkdirAll(wsWorkDir, 0o755); err != nil {
			return fmt.Errorf("create workspace for %s: %w", ev.ID, err)
		}
		isolateWorkspace(wsWorkDir)
		wsOutcome := runOne(ctx, ev, engine.Options{
			Prompt:        ev.Prompt,
			PrimaryModel:  *primaryModel,
			FallbackModel: *fallbackModel,
			PluginDir:     harnessDir,
			WorkDir:       wsWorkDir,
			Timeout:       perCaseTimeout,
		}, "with_skill", wsWorkDir)
		if err := wsOutcome.Write(iterationDir); err != nil {
			return fmt.Errorf("write %s with_skill outcome: %w", ev.ID, err)
		}
		withSkill = append(withSkill, wsOutcome)
		reportOne(wsOutcome)
		if wsOutcome.EngineErr != "" || wsOutcome.Grading.Summary.Passed != wsOutcome.Grading.Summary.Total {
			failures++
		}

		if *benchmark {
			woWorkDir := filepath.Join(iterationDir, ev.ID, "without_skill", "workspace")
			if err := os.MkdirAll(woWorkDir, 0o755); err != nil {
				return fmt.Errorf("create baseline workspace for %s: %w", ev.ID, err)
			}
			isolateWorkspace(woWorkDir)
			woOutcome := runOne(ctx, ev, engine.Options{
				Prompt:        ev.Prompt,
				PrimaryModel:  *primaryModel,
				FallbackModel: *fallbackModel,
				PluginDir:     "",
				WorkDir:       woWorkDir,
				Timeout:       perCaseTimeout,
			}, "without_skill", woWorkDir)
			if err := woOutcome.Write(iterationDir); err != nil {
				return fmt.Errorf("write %s without_skill outcome: %w", ev.ID, err)
			}
			withoutSkill = append(withoutSkill, woOutcome)
		}
	}

	if selected == 0 {
		return fmt.Errorf("-include %q matched none of the %d case(s) in %s", *include, len(suite.Evals), evalsPath(skillDir))
	}

	if *benchmark {
		if err := report.WriteBenchmark(iterationDir, withSkill, withoutSkill); err != nil {
			return err
		}
	}
	if err := report.WriteHTML(iterationDir, suite.SkillName, withSkill, withoutSkill, *benchmark); err != nil {
		return err
	}

	fmt.Printf("\n📋 Results: %d/%d cases fully passed — report: %s\n",
		len(withSkill)-failures, len(withSkill), filepath.Join(iterationDir, "report.html"))
	if failures > 0 {
		// Returned rather than os.Exit'd here so the harnessDir defer above
		// still runs — os.Exit skips all deferred cleanup.
		return errCasesFailed
	}
	return nil
}

func runOne(ctx context.Context, ev evalspec.Eval, opts engine.Options, configuration, workDir string) report.RunOutcome {
	res, err := engine.Run(ctx, opts)
	if err != nil {
		return report.RunOutcome{EvalID: ev.ID, Configuration: configuration, EngineErr: err.Error()}
	}
	g := grading.Grade(ev.Assertions, res.FinalMessage, workDir)
	return report.RunOutcome{EvalID: ev.ID, Configuration: configuration, Result: res, Grading: &g}
}

func reportOne(o report.RunOutcome) {
	if o.EngineErr != "" {
		fmt.Printf("   ⚠️  %s: ERROR — %s\n", o.EvalID, truncate(o.EngineErr, 120))
		return
	}
	status := "✅ PASS"
	if o.Grading.Summary.Passed != o.Grading.Summary.Total {
		status = "❌ FAIL"
	}
	fallback := ""
	if o.Result.FallbackUsed {
		fallback = " (fallback engaged)"
	}
	fmt.Printf("   %s %s: %d/%d assertions%s\n", status, o.EvalID, o.Grading.Summary.Passed, o.Grading.Summary.Total, fallback)
}

// reorderArgs lets flags appear before or after the positional skill-dir
// argument (Go's flag package otherwise stops parsing at the first
// non-flag token, which is a common footgun: `smeval run skills/x
// -include foo` would silently ignore -include). boolFlags names the
// flags that take no value, so their following token is treated as
// positional rather than consumed as the flag's value.
func reorderArgs(args []string, boolFlags map[string]bool) []string {
	var flags, positional []string
	for i := 0; i < len(args); i++ {
		a := args[i]
		if !strings.HasPrefix(a, "-") {
			positional = append(positional, a)
			continue
		}
		flags = append(flags, a)
		name := strings.TrimLeft(a, "-")
		if idx := strings.IndexByte(name, '='); idx >= 0 {
			continue // "-flag=value" is self-contained
		}
		if !boolFlags[name] && i+1 < len(args) {
			i++
			flags = append(flags, args[i])
		}
	}
	return append(flags, positional...)
}

func truncate(s string, n int) string {
	s = strings.TrimSpace(strings.ReplaceAll(s, "\n", " "))
	if len(s) <= n {
		return s
	}
	return s[:n] + "..."
}
