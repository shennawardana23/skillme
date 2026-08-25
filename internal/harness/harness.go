// Package harness builds a single-skill throwaway Claude Code plugin
// directory so a with_skill eval run installs exactly the skill under
// test — never the other skills in this catalog, and never its own
// evals/ directory (which would leak the grading criteria into the
// model's context).
package harness

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

// Build copies skillDir into a fresh temp directory shaped as a minimal
// plugin (.claude-plugin/plugin.json + skills/<name>/...), excluding the
// evals/ subdirectory, and returns the temp plugin root suitable for
// `claude --plugin-dir`. On any error, the partially-built temp directory
// is removed before returning — callers never need to clean up a failed
// Build themselves.
func Build(skillDir string) (dir string, err error) {
	skillName := filepath.Base(filepath.Clean(skillDir))

	root, err := newPluginRoot()
	if err != nil {
		return "", err
	}
	defer func() {
		if err != nil {
			os.RemoveAll(root)
		}
	}()

	dest := filepath.Join(root, "skills", skillName)
	if err = copyDirExcluding(skillDir, dest, "evals"); err != nil {
		return "", fmt.Errorf("copy skill into harness: %w", err)
	}
	return root, nil
}

// BuildCatalog copies every skill directory directly under skillsRootDir
// (any entry containing a SKILL.md) into a fresh temp plugin — the same
// shape Build produces, but with the entire catalog installed at once
// instead of one skill in isolation. Each skill's evals/ subdirectory is
// still excluded, for the same reason Build excludes it: none of this
// repo's grading criteria should ever reach a model's context.
//
// Use this for trigger-accuracy and coexistence testing (see
// catalog-evals/), where the question is whether Claude picks the right
// skill out of the whole catalog — Build's single-skill harness cannot
// exercise that at all, since there is nothing else present to (correctly,
// or incorrectly) not pick.
func BuildCatalog(skillsRootDir string) (dir string, err error) {
	entries, err := os.ReadDir(skillsRootDir)
	if err != nil {
		return "", fmt.Errorf("read %s: %w", skillsRootDir, err)
	}

	root, err := newPluginRoot()
	if err != nil {
		return "", err
	}
	defer func() {
		if err != nil {
			os.RemoveAll(root)
		}
	}()

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		src := filepath.Join(skillsRootDir, entry.Name())
		if _, statErr := os.Stat(filepath.Join(src, "SKILL.md")); statErr != nil {
			continue // not a skill directory (e.g. a stray file, or a dir with no SKILL.md)
		}
		dest := filepath.Join(root, "skills", entry.Name())
		if err = copyDirExcluding(src, dest, "evals"); err != nil {
			return "", fmt.Errorf("copy skill %s into harness: %w", entry.Name(), err)
		}
	}
	return root, nil
}

// newPluginRoot creates a fresh temp directory shaped as a minimal Claude
// Code plugin — just the .claude-plugin/plugin.json manifest — for Build or
// BuildCatalog to then populate with skills/. On any error, the partially
// created directory is removed before returning.
func newPluginRoot() (root string, err error) {
	root, err = os.MkdirTemp("", "smeval-harness-*")
	if err != nil {
		return "", fmt.Errorf("create harness temp dir: %w", err)
	}
	defer func() {
		if err != nil {
			os.RemoveAll(root)
		}
	}()

	manifestDir := filepath.Join(root, ".claude-plugin")
	if err = os.MkdirAll(manifestDir, 0o755); err != nil {
		return "", fmt.Errorf("create %s: %w", manifestDir, err)
	}
	manifest, merr := json.MarshalIndent(map[string]string{"name": "smeval-harness"}, "", "  ")
	if merr != nil {
		err = merr
		return "", err
	}
	if err = os.WriteFile(filepath.Join(manifestDir, "plugin.json"), manifest, 0o644); err != nil {
		return "", fmt.Errorf("write harness plugin.json: %w", err)
	}
	return root, nil
}

// copyDirExcluding recursively copies src to dst, skipping any top-level
// entry named excludeName.
func copyDirExcluding(src, dst, excludeName string) error {
	entries, err := os.ReadDir(src)
	if err != nil {
		return err
	}
	if err := os.MkdirAll(dst, 0o755); err != nil {
		return err
	}
	for _, entry := range entries {
		if entry.Name() == excludeName {
			continue
		}
		if err := copyEntry(filepath.Join(src, entry.Name()), filepath.Join(dst, entry.Name())); err != nil {
			return err
		}
	}
	return nil
}

// CopyTree recursively copies every entry under src into dst, following
// symlinked directories. Exported for reuse seeding a case's workspace with
// a fixture (see catalog-evals/README.md and evalspec.Eval.Fixture) — the
// exact same "copy everything, no exclusions" operation Build/BuildCatalog
// need internally, just applied to a fixture directory instead of a skill.
func CopyTree(src, dst string) error {
	return copyDirAll(src, dst)
}

func copyDirAll(src, dst string) error {
	entries, err := os.ReadDir(src)
	if err != nil {
		return err
	}
	if err := os.MkdirAll(dst, 0o755); err != nil {
		return err
	}
	for _, entry := range entries {
		if err := copyEntry(filepath.Join(src, entry.Name()), filepath.Join(dst, entry.Name())); err != nil {
			return err
		}
	}
	return nil
}

// copyEntry copies one directory entry, following symlinks (os.Stat rather
// than the DirEntry's own type) so a symlinked directory is recursed into
// rather than mistakenly handed to copyFile, which would fail trying to
// read a directory as a regular file.
func copyEntry(src, dst string) error {
	info, err := os.Stat(src)
	if err != nil {
		return fmt.Errorf("stat %s: %w", src, err)
	}
	if info.IsDir() {
		return copyDirAll(src, dst)
	}
	return copyFile(src, dst)
}

func copyFile(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()
	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()
	_, err = io.Copy(out, in)
	return err
}
