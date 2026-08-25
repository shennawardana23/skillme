package harness

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"
)

func TestBuild_CleansUpOnError(t *testing.T) {
	before, err := os.ReadDir(os.TempDir())
	if err != nil {
		t.Fatalf("read temp dir: %v", err)
	}
	countBefore := countHarnessDirs(before)

	// A skill directory that does not exist makes copyDirExcluding fail
	// after MkdirTemp has already created the harness root — Build must
	// remove that root before returning the error.
	_, err = Build(filepath.Join(t.TempDir(), "does-not-exist"))
	if err == nil {
		t.Fatal("Build should fail for a nonexistent skill directory")
	}

	after, err := os.ReadDir(os.TempDir())
	if err != nil {
		t.Fatalf("read temp dir: %v", err)
	}
	countAfter := countHarnessDirs(after)

	if countAfter > countBefore {
		t.Fatalf("Build leaked a smeval-harness-* temp directory on error: before=%d after=%d", countBefore, countAfter)
	}
}

func countHarnessDirs(entries []os.DirEntry) int {
	n := 0
	for _, e := range entries {
		if e.IsDir() && len(e.Name()) >= len("smeval-harness-") && e.Name()[:len("smeval-harness-")] == "smeval-harness-" {
			n++
		}
	}
	return n
}

// TestBuildCatalog_IncludesMultipleSkillsExcludingEvals proves BuildCatalog
// installs every skill under the given root — the property a single-skill
// Build harness cannot offer at all, since coexistence/trigger-accuracy
// testing needs more than one skill actually present to choose among — and
// still excludes each one's evals/ directory.
func TestBuildCatalog_IncludesMultipleSkillsExcludingEvals(t *testing.T) {
	root := t.TempDir()
	for _, name := range []string{"skill-a", "skill-b"} {
		dir := filepath.Join(root, name)
		if err := os.MkdirAll(filepath.Join(dir, "evals"), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(filepath.Join(dir, "SKILL.md"), []byte("---\nname: "+name+"\n---\nbody"), 0o644); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(filepath.Join(dir, "evals", "evals.json"), []byte(`{"skill_name":"`+name+`"}`), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	// A stray non-skill entry (no SKILL.md) must be skipped, not copied or
	// treated as an error.
	if err := os.MkdirAll(filepath.Join(root, "not-a-skill"), 0o755); err != nil {
		t.Fatal(err)
	}

	harnessDir, err := BuildCatalog(root)
	if err != nil {
		t.Fatalf("BuildCatalog failed: %v", err)
	}
	defer os.RemoveAll(harnessDir)

	for _, name := range []string{"skill-a", "skill-b"} {
		if _, err := os.Stat(filepath.Join(harnessDir, "skills", name, "SKILL.md")); err != nil {
			t.Fatalf("expected %s/SKILL.md to be copied into the catalog harness: %v", name, err)
		}
		if _, err := os.Stat(filepath.Join(harnessDir, "skills", name, "evals")); !os.IsNotExist(err) {
			t.Fatalf("expected %s's evals/ to be excluded from the catalog harness, got err=%v", name, err)
		}
	}
	if _, err := os.Stat(filepath.Join(harnessDir, "skills", "not-a-skill")); !os.IsNotExist(err) {
		t.Fatal("expected the non-skill directory (no SKILL.md) to be skipped entirely")
	}
}

func TestBuild_FollowsSymlinkedDirectory(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("symlink creation requires elevated privileges on windows")
	}

	skillDir := t.TempDir()
	if err := os.WriteFile(filepath.Join(skillDir, "SKILL.md"), []byte("---\nname: x\n---\nbody"), 0o644); err != nil {
		t.Fatal(err)
	}
	realRefs := filepath.Join(t.TempDir(), "real-references")
	if err := os.MkdirAll(realRefs, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(realRefs, "notes.md"), []byte("shared notes"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink(realRefs, filepath.Join(skillDir, "references")); err != nil {
		t.Fatalf("create symlink: %v", err)
	}

	root, err := Build(skillDir)
	if err != nil {
		t.Fatalf("Build failed on a skill dir with a symlinked directory: %v", err)
	}
	defer os.RemoveAll(root)

	copied := filepath.Join(root, "skills", filepath.Base(skillDir), "references", "notes.md")
	data, err := os.ReadFile(copied)
	if err != nil {
		t.Fatalf("expected symlinked directory contents to be copied through: %v", err)
	}
	if string(data) != "shared notes" {
		t.Fatalf("copied content = %q, want %q", data, "shared notes")
	}
}
