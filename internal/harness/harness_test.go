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
