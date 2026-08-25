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
	"regexp"
	"strings"
	"time"

	"github.com/shennawardana23/skillme/internal/engine"
	"github.com/shennawardana23/skillme/internal/evalspec"
	"github.com/shennawardana23/skillme/internal/grading"
	"github.com/shennawardana23/skillme/internal/harness"
	"github.com/shennawardana23/skillme/internal/report"
	"github.com/shennawardana23/skillme/internal/riskscan"
	"github.com/shennawardana23/skillme/internal/similarity"
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
	case "trigger-validate":
		err = runTriggerValidate(os.Args[2:])
	case "trigger-run":
		err = runTriggerRun(os.Args[2:])
	case "security-scan":
		err = runSecurityScan(os.Args[2:])
	case "similarity-check":
		err = runSimilarityCheck(os.Args[2:])
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
  smeval trigger-validate <trigger-cases.json>
  smeval trigger-run <trigger-cases.json> [flags]
  smeval security-scan <skill-dir> [-quiet]
  smeval similarity-check [-skills-dir skills] [-threshold 0.3]

Flags for run:
  -primary-model string   Model for the primary attempt (default "sonnet")
  -fallback-model string  Model for the outer fallback attempt; also passed
                          to claude's native --fallback-model (default "opus")
  -timeout duration       Per-attempt timeout (default 3m0s)
  -benchmark              Also run each case without the skill installed and write benchmark.json
  -output-dir string      Workspace root (default "<skill-dir>-workspace", a sibling directory)
  -include string         Only run eval IDs containing this substring

trigger-run installs the ENTIRE catalog (every skill under -skills-dir) into
one harness instead of a single skill in isolation, so each case tests
whether Claude picks the right skill's guidance out of the whole catalog —
triggering accuracy and coexistence, not per-skill correctness. See
catalog-evals/trigger-accuracy.json for the case format (same schema as a
skill's evals.json) and cmd/smeval's package doc for why this exists.

Flags for trigger-run:
  -primary-model, -fallback-model, -timeout, -output-dir, -include
                          same meaning as for run
  -skills-dir string      Catalog root to install in full (default "skills")

security-scan statically scans a skill directory against Anthropic's
enterprise Risk Tier indicators (code execution, network access, hardcoded
credentials, suspicious paths, MCP references) — no model calls, free, fast.
Exits non-zero only on a High-concern finding; Medium/Informational findings
print but never fail the scan. See internal/riskscan's package doc for which
two Risk Tier indicators are deliberately NOT automated here and why.

Flags for security-scan:
  -quiet                  Only print output if there's at least one finding

similarity-check reports skill-description pairs likely to overlap in
coverage, using TF-IDF cosine similarity over name+description — no model
calls, free, fast, plain Go. It is advisory only and always exits 0: a high
score is a nudge to go read both descriptions and decide (per CONTRIBUTING's
"prefer extending an existing skill" rule), not a fail. See
internal/similarity's package doc for the free-vs-embeddings-API trade-off
this deliberately takes.

Flags for similarity-check:
  -skills-dir string      Catalog root to scan (default "skills")
  -threshold float        Minimum score to report, 0-1 (default 0.3)`)
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

// seedFixture copies a context-dependent eval's fixture files into workDir
// before the case runs, so its prompt can describe a scenario ("trace a bug
// across files") without needing to inline the entire codebase as prompt
// text — see evalspec.Eval's doc for why that distinction matters. A no-op
// for a greenfield case (the overwhelming majority). Errors clearly rather
// than silently running greenfield when a context-dependent case is routed
// through a command that doesn't support fixtures (fixturesRoot == "", the
// case for trigger-run today) or when the referenced fixture directory
// doesn't exist — a silently-empty workspace would make the case fail for
// a confusing reason days later instead of failing loudly right here.
func seedFixture(fixturesRoot string, ev evalspec.Eval, workDir string) error {
	if !ev.IsContextDependent() {
		return nil
	}
	if strings.Contains(ev.Fixture, "..") {
		return fmt.Errorf("eval %q: fixture %q must not contain \"..\"", ev.ID, ev.Fixture)
	}
	if fixturesRoot == "" {
		return fmt.Errorf("eval %q requires fixture %q, but this command doesn't support fixtures", ev.ID, ev.Fixture)
	}
	src := filepath.Join(fixturesRoot, ev.Fixture)
	if _, err := os.Stat(src); err != nil {
		return fmt.Errorf("eval %q: fixture directory %s: %w", ev.ID, src, err)
	}
	return harness.CopyTree(src, workDir)
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
	for _, w := range evalspec.Lint(suite) {
		fmt.Printf("  ⚠ %s\n", w)
	}
	for _, ev := range suite.Evals {
		if !ev.IsContextDependent() {
			fmt.Printf("  · %s: greenfield\n", ev.ID)
			continue
		}
		fixtureDir := filepath.Join(skillDir, "evals", "fixtures", ev.Fixture)
		if _, err := os.Stat(fixtureDir); err != nil {
			return fmt.Errorf("eval %q: fixture %q not found at %s", ev.ID, ev.Fixture, fixtureDir)
		}
		fmt.Printf("  · %s: context-dependent (fixture: %s)\n", ev.ID, ev.Fixture)
	}
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
		// One shared, already-gitignored root (smeval-workspace/) instead of a
		// "<skill>-workspace" sibling scattered next to every single skill —
		// 125 skills each leaving their own workspace directory inside
		// skills/ turns that listing into 50% real content, 50% run output.
		*outputDir = filepath.Join("smeval-workspace", "runs", filepath.Base(filepath.Clean(skillDir)))
	}

	harnessDir, err := harness.Build(skillDir)
	if err != nil {
		return fmt.Errorf("build harness plugin: %w", err)
	}
	defer os.RemoveAll(harnessDir)

	return runSuite(suite, harnessDir, runOptions{
		primaryModel:  *primaryModel,
		fallbackModel: *fallbackModel,
		timeout:       *timeout,
		benchmark:     *benchmark,
		outputDir:     *outputDir,
		include:       *include,
		loadErrLabel:  evalsPath(skillDir),
		fixturesRoot:  filepath.Join(skillDir, "evals", "fixtures"),
	})
}

// runOptions parameterizes runSuite across its two callers: runRun (one
// skill in isolation, benchmarking supported) and runTriggerRun (the whole
// catalog installed at once, benchmarking not applicable — there is no
// single skill to strip out for a without_skill baseline).
type runOptions struct {
	primaryModel  string
	fallbackModel string
	timeout       time.Duration
	benchmark     bool
	outputDir     string
	include       string
	loadErrLabel  string // path named in the "-include matched nothing" error
	fixturesRoot  string // dir containing named fixture subdirs for context-dependent cases; empty = fixtures unsupported (e.g. trigger-run)
}

// runSuite executes every (optionally -include-filtered) case in suite
// against harnessDir, writes the iteration's workspace and report, and
// returns errCasesFailed if any case didn't fully pass. Shared by runRun
// and runTriggerRun — the only difference between them is how harnessDir
// was built (harness.Build for one skill vs. harness.BuildCatalog for the
// whole catalog) and whether benchmarking applies.
func runSuite(suite *evalspec.Suite, harnessDir string, opts runOptions) error {
	iterationDir, err := report.NextIterationDir(opts.outputDir)
	if err != nil {
		return err
	}

	ctx := context.Background()
	var withSkill, withoutSkill []report.RunOutcome
	failures := 0
	selected := 0

	for _, ev := range suite.Evals {
		if opts.include != "" && !strings.Contains(ev.ID, opts.include) {
			continue
		}
		selected++

		perCaseTimeout := opts.timeout
		if ev.TimeoutSeconds > 0 {
			perCaseTimeout = time.Duration(ev.TimeoutSeconds) * time.Second
		}

		fmt.Printf("⏳ %s: %s\n", ev.ID, truncate(ev.ExpectedOutput, 70))

		wsWorkDir := filepath.Join(iterationDir, ev.ID, "with_skill", "workspace")
		if err := os.MkdirAll(wsWorkDir, 0o755); err != nil {
			return fmt.Errorf("create workspace for %s: %w", ev.ID, err)
		}
		if err := seedFixture(opts.fixturesRoot, ev, wsWorkDir); err != nil {
			return err
		}
		isolateWorkspace(wsWorkDir)
		wsOutcome := runOne(ctx, ev, engine.Options{
			Prompt:        ev.Prompt,
			PrimaryModel:  opts.primaryModel,
			FallbackModel: opts.fallbackModel,
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

		if opts.benchmark {
			woWorkDir := filepath.Join(iterationDir, ev.ID, "without_skill", "workspace")
			if err := os.MkdirAll(woWorkDir, 0o755); err != nil {
				return fmt.Errorf("create baseline workspace for %s: %w", ev.ID, err)
			}
			if err := seedFixture(opts.fixturesRoot, ev, woWorkDir); err != nil {
				return err
			}
			isolateWorkspace(woWorkDir)
			woOutcome := runOne(ctx, ev, engine.Options{
				Prompt:        ev.Prompt,
				PrimaryModel:  opts.primaryModel,
				FallbackModel: opts.fallbackModel,
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
		return fmt.Errorf("-include %q matched none of the %d case(s) in %s", opts.include, len(suite.Evals), opts.loadErrLabel)
	}

	if opts.benchmark {
		if err := report.WriteBenchmark(iterationDir, withSkill, withoutSkill); err != nil {
			return err
		}
	}
	if err := report.WriteHTML(iterationDir, suite.SkillName, withSkill, withoutSkill, opts.benchmark); err != nil {
		return err
	}
	evalIDs := make([]string, len(withSkill))
	for i, o := range withSkill {
		evalIDs[i] = o.EvalID
	}
	if err := report.WriteFeedbackStub(iterationDir, evalIDs); err != nil {
		return err
	}

	fmt.Printf("\n📋 Results: %d/%d cases fully passed — report: %s\n",
		len(withSkill)-failures, len(withSkill), filepath.Join(iterationDir, "report.html"))
	fmt.Printf("   Human review: open each case's outputs/response.md, then fill in %s\n",
		filepath.Join(iterationDir, "feedback.json"))
	if failures > 0 {
		// Returned rather than os.Exit'd here so the harnessDir defer in the
		// caller still runs — os.Exit skips all deferred cleanup.
		return errCasesFailed
	}
	return nil
}

// runTriggerValidate checks a trigger-cases file's schema — the same
// evals.json shape a skill uses, but with no skills/<name> directory
// backing it, so (unlike runValidate) there is no directory-name match to
// enforce.
func runTriggerValidate(args []string) error {
	fs := flag.NewFlagSet("trigger-validate", flag.ExitOnError)
	fs.Parse(args)
	if fs.NArg() != 1 {
		return fmt.Errorf("usage: smeval trigger-validate <trigger-cases.json>")
	}
	path := fs.Arg(0)
	suite, err := evalspec.Load(path)
	if err != nil {
		return err
	}
	fmt.Printf("✓ %s is valid (%d case(s))\n", path, len(suite.Evals))
	for _, w := range evalspec.Lint(suite) {
		fmt.Printf("  ⚠ %s\n", w)
	}
	return nil
}

// runTriggerRun installs the whole catalog under -skills-dir into one
// harness (harness.BuildCatalog, not harness.Build) and grades each case
// against it — see the package doc and usage()'s trigger-run section for
// why this exists: it is the only command in smeval that can measure
// triggering accuracy or cross-skill coexistence, since every other command
// installs exactly one skill in isolation.
func runTriggerRun(args []string) error {
	fs := flag.NewFlagSet("trigger-run", flag.ExitOnError)
	primaryModel := fs.String("primary-model", "sonnet", "")
	fallbackModel := fs.String("fallback-model", "opus", "")
	timeout := fs.Duration("timeout", 3*time.Minute, "")
	outputDir := fs.String("output-dir", "", "")
	include := fs.String("include", "", "")
	skillsDir := fs.String("skills-dir", "skills", "")
	fs.Parse(reorderArgs(args, map[string]bool{}))
	if fs.NArg() != 1 {
		return fmt.Errorf("usage: smeval trigger-run <trigger-cases.json> [flags]")
	}
	casesPath := fs.Arg(0)

	suite, err := evalspec.Load(casesPath)
	if err != nil {
		return err
	}

	if *outputDir == "" {
		*outputDir = filepath.Join("smeval-workspace", "runs", "catalog-triggers")
	}

	harnessDir, err := harness.BuildCatalog(*skillsDir)
	if err != nil {
		return fmt.Errorf("build catalog harness: %w", err)
	}
	defer os.RemoveAll(harnessDir)

	return runSuite(suite, harnessDir, runOptions{
		primaryModel:  *primaryModel,
		fallbackModel: *fallbackModel,
		timeout:       *timeout,
		benchmark:     false, // no single skill to strip out for a without_skill baseline
		outputDir:     *outputDir,
		include:       *include,
		loadErrLabel:  casesPath,
	})
}

// runSecurityScan statically scans one skill directory and prints every
// finding, most-severe first. Exits non-zero (via errCasesFailed, so main
// doesn't print a redundant "smeval: ..." line) only when at least one
// High-concern finding exists — Medium and Informational findings are
// printed for human review but never fail the scan on their own. See
// internal/riskscan's package doc for the two Risk Tier indicators this
// deliberately does not turn into an automated check.
func runSecurityScan(args []string) error {
	fs := flag.NewFlagSet("security-scan", flag.ExitOnError)
	quiet := fs.Bool("quiet", false, "")
	fs.Parse(reorderArgs(args, map[string]bool{"quiet": true}))
	if fs.NArg() != 1 {
		return fmt.Errorf("usage: smeval security-scan <skill-dir> [-quiet]")
	}
	skillDir := fs.Arg(0)

	findings, err := riskscan.Scan(skillDir)
	if err != nil {
		return err
	}

	if len(findings) == 0 {
		if !*quiet {
			fmt.Printf("✓ %s: no findings\n", skillDir)
		}
		return nil
	}

	fmt.Printf("%s: %d finding(s)\n", skillDir, len(findings))
	for _, f := range findings {
		loc := f.File
		if f.Line > 0 {
			loc = fmt.Sprintf("%s:%d", f.File, f.Line)
		}
		icon := "ℹ️ "
		switch f.Concern {
		case riskscan.High:
			icon = "🔴"
		case riskscan.Medium:
			icon = "🟡"
		}
		fmt.Printf("  %s [%s] %s (%s): %s\n", icon, f.Concern, f.Indicator, loc, truncate(f.Detail, 100))
	}

	if riskscan.HasFailing(findings) {
		return errCasesFailed
	}
	return nil
}

// nameAndDescription extracts a SKILL.md's name and description frontmatter
// fields with a plain regex, not a YAML parser — every skill in this
// catalog has exactly one `description:` line as a single physical line
// (verified against all 137 skills before relying on this), so a full
// parser would be unneeded weight for what this command needs.
var (
	namePattern        = regexp.MustCompile(`(?m)^name:\s*(.+)$`)
	descriptionPattern = regexp.MustCompile(`(?m)^description:\s*(.+)$`)
)

func nameAndDescription(skillMDPath string) (name, description string, err error) {
	data, err := os.ReadFile(skillMDPath)
	if err != nil {
		return "", "", err
	}
	content := string(data)
	if m := namePattern.FindStringSubmatch(content); m != nil {
		name = strings.TrimSpace(m[1])
	}
	if m := descriptionPattern.FindStringSubmatch(content); m != nil {
		description = strings.Trim(strings.TrimSpace(m[1]), `"`)
	}
	return name, description, nil
}

// runSimilarityCheck is advisory-only and always returns nil (never fails a
// build) — see internal/similarity's package doc and usage()'s
// similarity-check section for why. A human decides what a high score
// means; this command's job is only to surface candidates worth their
// attention, per CONTRIBUTING's "search the catalog first" pre-flight step.
func runSimilarityCheck(args []string) error {
	fs := flag.NewFlagSet("similarity-check", flag.ExitOnError)
	skillsDir := fs.String("skills-dir", "skills", "")
	threshold := fs.Float64("threshold", 0.3, "")
	fs.Parse(reorderArgs(args, map[string]bool{}))

	entries, err := os.ReadDir(*skillsDir)
	if err != nil {
		return fmt.Errorf("read %s: %w", *skillsDir, err)
	}

	docs := make(map[string]string)
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		skillMD := filepath.Join(*skillsDir, entry.Name(), "SKILL.md")
		name, description, err := nameAndDescription(skillMD)
		if err != nil {
			continue // not a skill directory (no SKILL.md) — skip, don't fail
		}
		if name == "" {
			name = entry.Name()
		}
		docs[name] = name + " " + description
	}

	corpus := similarity.BuildCorpus(docs)
	pairs := corpus.TopPairs(*threshold)

	if len(pairs) == 0 {
		fmt.Printf("No pairs scored >= %.2f across %d skills.\n", *threshold, len(docs))
		return nil
	}
	fmt.Printf("%d pair(s) scored >= %.2f across %d skills (highest first):\n", len(pairs), *threshold, len(docs))
	for _, p := range pairs {
		fmt.Printf("  %.3f  %s  <->  %s\n", p.Score, p.A, p.B)
	}
	fmt.Println("\nAdvisory only — read both descriptions and decide per CONTRIBUTING.md's")
	fmt.Println("\"prefer extending an existing skill\" rule; this never fails a build.")
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
