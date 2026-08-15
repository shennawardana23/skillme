package main

import (
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
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
