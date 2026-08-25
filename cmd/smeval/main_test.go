package main

import (
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/shennawardana23/skillme/internal/evalspec"
)

func TestReorderArgs(t *testing.T) {
	boolFlags := map[string]bool{"benchmark": true}

	tests := []struct {
		name string
		args []string
		want []string
	}{
		{
			name: "flags before positional (already fine)",
			args: []string{"-include", "foo", "skills/x"},
			want: []string{"-include", "foo", "skills/x"},
		},
		{
			name: "flags after positional (the footgun this fixes)",
			args: []string{"skills/x", "-include", "foo", "-output-dir", "/tmp/out"},
			want: []string{"-include", "foo", "-output-dir", "/tmp/out", "skills/x"},
		},
		{
			name: "bool flag after positional does not consume the next token",
			args: []string{"skills/x", "-benchmark", "-include", "foo"},
			want: []string{"-benchmark", "-include", "foo", "skills/x"},
		},
		{
			name: "flag=value form after positional",
			args: []string{"skills/x", "-include=foo"},
			want: []string{"-include=foo", "skills/x"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := reorderArgs(tt.args, boolFlags)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("reorderArgs(%v) = %v, want %v", tt.args, got, tt.want)
			}
		})
	}
}

// TestRunRun_IncludeMatchingNothingFails proves a typo'd -include (or a
// future eval-id rename that silently orphans one somewhere) fails loudly
// instead of a false-green "0/0 cases fully passed" exit 0. Because no
// case is selected, the loop body — and therefore the engine, which would
// need a real or fake `claude` binary — never runs, so this test needs no
// engine stub at all.
func TestRunRun_IncludeMatchingNothingFails(t *testing.T) {
	skillDir := t.TempDir()
	if err := os.WriteFile(filepath.Join(skillDir, "SKILL.md"), []byte("---\nname: x\ndescription: x\n---\nbody"), 0o644); err != nil {
		t.Fatal(err)
	}
	evalsDir := filepath.Join(skillDir, "evals")
	if err := os.MkdirAll(evalsDir, 0o755); err != nil {
		t.Fatal(err)
	}
	suiteJSON := `{
		"skill_name": "` + filepath.Base(skillDir) + `",
		"evals": [{
			"id": "real-case",
			"prompt": "p",
			"expected_output": "e",
			"assertions": [{"text": "t", "check": {"contains_all": ["x"]}}]
		}]
	}`
	if err := os.WriteFile(filepath.Join(evalsDir, "evals.json"), []byte(suiteJSON), 0o644); err != nil {
		t.Fatal(err)
	}

	err := runRun([]string{skillDir, "-include", "this-id-matches-nothing", "-output-dir", t.TempDir()})
	if err == nil {
		t.Fatal("runRun should fail when -include matches zero cases, not silently succeed")
	}
	if strings.Contains(err.Error(), "matched none") == false {
		t.Fatalf("expected a clear 'matched none' error, got: %v", err)
	}
}

func TestSeedFixture_GreenfieldCaseIsNoOp(t *testing.T) {
	ev := evalspec.Eval{ID: "x"} // no Fixture set
	workDir := t.TempDir()
	if err := seedFixture("/does/not/matter", ev, workDir); err != nil {
		t.Fatalf("seedFixture on a greenfield case should be a no-op, got: %v", err)
	}
	entries, _ := os.ReadDir(workDir)
	if len(entries) != 0 {
		t.Fatalf("seedFixture should not have written anything into %s for a greenfield case", workDir)
	}
}

// TestSeedFixture_CopiesFixtureFilesIntoWorkspace proves the actual
// mechanism this whole feature is for: a context-dependent case's fixture
// files land in the case workspace exactly as authored, ready for the
// model's own Read/Grep tools to find — not inlined into the prompt.
func TestSeedFixture_CopiesFixtureFilesIntoWorkspace(t *testing.T) {
	fixturesRoot := t.TempDir()
	fixtureDir := filepath.Join(fixturesRoot, "two-file-bug")
	if err := os.MkdirAll(filepath.Join(fixtureDir, "pkg"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(fixtureDir, "main.go"), []byte("package main\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(fixtureDir, "pkg", "config.go"), []byte("package pkg\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	ev := evalspec.Eval{ID: "x", Fixture: "two-file-bug"}
	workDir := t.TempDir()
	if err := seedFixture(fixturesRoot, ev, workDir); err != nil {
		t.Fatalf("seedFixture failed: %v", err)
	}

	for _, rel := range []string{"main.go", "pkg/config.go"} {
		if _, err := os.Stat(filepath.Join(workDir, rel)); err != nil {
			t.Fatalf("expected %s to be copied into the workspace: %v", rel, err)
		}
	}
}

func TestSeedFixture_ErrorsWhenFixturesUnsupported(t *testing.T) {
	ev := evalspec.Eval{ID: "x", Fixture: "some-fixture"}
	err := seedFixture("", ev, t.TempDir()) // empty fixturesRoot == trigger-run today
	if err == nil {
		t.Fatal("seedFixture should error, not silently run greenfield, when fixturesRoot is empty but a fixture is requested")
	}
}

func TestSeedFixture_ErrorsWhenFixtureDirMissing(t *testing.T) {
	ev := evalspec.Eval{ID: "x", Fixture: "does-not-exist"}
	err := seedFixture(t.TempDir(), ev, t.TempDir())
	if err == nil {
		t.Fatal("seedFixture should error when the referenced fixture directory doesn't exist")
	}
}

func TestSeedFixture_RejectsPathTraversal(t *testing.T) {
	ev := evalspec.Eval{ID: "x", Fixture: "../../etc"}
	err := seedFixture(t.TempDir(), ev, t.TempDir())
	if err == nil || !strings.Contains(err.Error(), "..") {
		t.Fatalf("seedFixture should reject a fixture path containing '..', got: %v", err)
	}
}

// TestRunValidate_FailsOnMissingFixtureDirectory proves a context-dependent
// case with a typo'd or deleted fixture is caught by the free, schema-only
// validate path — before a live run would waste an API call discovering
// the same problem the hard way.
func TestRunValidate_FailsOnMissingFixtureDirectory(t *testing.T) {
	skillDir := t.TempDir()
	if err := os.WriteFile(filepath.Join(skillDir, "SKILL.md"), []byte("---\nname: "+filepath.Base(skillDir)+"\ndescription: x\n---\nbody"), 0o644); err != nil {
		t.Fatal(err)
	}
	evalsDir := filepath.Join(skillDir, "evals")
	if err := os.MkdirAll(evalsDir, 0o755); err != nil {
		t.Fatal(err)
	}
	suiteJSON := `{
		"skill_name": "` + filepath.Base(skillDir) + `",
		"evals": [{
			"id": "needs-a-fixture",
			"prompt": "p",
			"expected_output": "e",
			"fixture": "missing-fixture",
			"assertions": [{"text": "t", "check": {"contains_all": ["x"]}}]
		}]
	}`
	if err := os.WriteFile(filepath.Join(evalsDir, "evals.json"), []byte(suiteJSON), 0o644); err != nil {
		t.Fatal(err)
	}

	err := runValidate([]string{skillDir})
	if err == nil {
		t.Fatal("runValidate should fail when a context-dependent case's fixture directory doesn't exist")
	}
	if !strings.Contains(err.Error(), "fixture") {
		t.Fatalf("expected a fixture-related error, got: %v", err)
	}
}
