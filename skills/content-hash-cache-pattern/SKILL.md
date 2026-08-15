---
name: content-hash-cache-pattern
description: Cache expensive file-processing results (PDF parsing, OCR, text extraction, image analysis) keyed by a SHA-256 hash of file content rather than file path, so renamed files still hit the cache and changed content auto-invalidates it. Use when building a file-processing pipeline that reprocesses the same files across runs, or when adding a --cache/--no-cache option without touching the existing pure processing function.
license: Apache-2.0
metadata:
  version: "0.1.0"
---

# Content-Hash File Cache Pattern

Cache expensive file-processing results using a SHA-256 hash of the file's *content* as the cache key. Unlike path-based caching, this survives file moves and renames, and auto-invalidates the moment content changes — no separate invalidation logic needed.

## When to Activate

- Building a file-processing pipeline (PDF parsing, OCR, text extraction, image analysis) where the same files reappear across runs
- Processing cost is high enough that reprocessing unchanged files is wasteful
- Adding a `--cache`/`--no-cache` flag to a CLI tool
- Retrofitting caching onto an existing pure function without modifying that function

## Core Pattern (Go)

### 1. Content hash as the cache key

Hash the file in chunks — never load the whole file into memory:

```go
package cache

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"os"
)

func ComputeFileHash(path string) (string, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", fmt.Errorf("open %s: %w", path, err)
	}
	defer f.Close()

	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		return "", fmt.Errorf("hash %s: %w", path, err)
	}
	return hex.EncodeToString(h.Sum(nil)), nil
}
```

`io.Copy` streams the file into the hasher in fixed-size buffers, so this scales to large files without loading them fully into memory.

### 2. Cache entry type

```go
type CacheEntry struct {
	FileHash   string   `json:"file_hash"`
	SourcePath string   `json:"source_path"`
	Document   Document `json:"document"` // the cached result
}
```

### 3. File-based storage: `{hash}.json`

O(1) lookup by hash, no separate index file:

```go
package cache

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

func WriteEntry(cacheDir string, entry CacheEntry) error {
	if err := os.MkdirAll(cacheDir, 0o755); err != nil {
		return fmt.Errorf("create cache dir: %w", err)
	}
	data, err := json.Marshal(entry)
	if err != nil {
		return fmt.Errorf("marshal cache entry: %w", err)
	}
	path := filepath.Join(cacheDir, entry.FileHash+".json")
	return os.WriteFile(path, data, 0o644)
}

// ReadEntry returns (entry, true) on a cache hit, (zero, false) on a
// miss — including a miss on corrupted or unreadable cache files, so
// corruption degrades to a re-process rather than a crash.
func ReadEntry(cacheDir, fileHash string) (CacheEntry, bool) {
	path := filepath.Join(cacheDir, fileHash+".json")
	data, err := os.ReadFile(path)
	if err != nil {
		return CacheEntry{}, false
	}
	var entry CacheEntry
	if err := json.Unmarshal(data, &entry); err != nil {
		return CacheEntry{}, false
	}
	return entry, true
}
```

### 4. Service-layer wrapper (single responsibility)

Keep the processing function pure. Caching is a separate layer that wraps it:

```go
// ExtractText is pure: no cache awareness, easy to test and reuse.
func ExtractText(path string) (Document, error) {
	// ... real extraction logic ...
	return Document{}, nil
}

type ExtractOptions struct {
	CacheEnabled bool
	CacheDir     string
}

// ExtractWithCache is the service layer: cache check -> extraction -> cache write.
func ExtractWithCache(path string, opts ExtractOptions) (Document, error) {
	if !opts.CacheEnabled {
		return ExtractText(path)
	}

	hash, err := ComputeFileHash(path)
	if err != nil {
		return Document{}, err
	}

	if entry, ok := ReadEntry(opts.CacheDir, hash); ok {
		return entry.Document, nil
	}

	doc, err := ExtractText(path)
	if err != nil {
		return Document{}, err
	}
	entry := CacheEntry{FileHash: hash, SourcePath: path, Document: doc}
	if err := WriteEntry(opts.CacheDir, entry); err != nil {
		// A cache write failure should not fail the caller — the result is
		// still correct, just not persisted for next time.
		return doc, nil
	}
	return doc, nil
}
```

## Key Design Decisions

| Decision | Rationale |
|---|---|
| SHA-256 content hash, not path | Path-independent; auto-invalidates the instant content changes |
| `{hash}.json` file naming | O(1) lookup, no index file to keep in sync |
| Service layer wraps a pure function | Extraction logic stays testable and reusable; caching is orthogonal |
| Corrupted/missing cache file → cache miss | Graceful degradation: reprocess rather than crash |
| Cache write failure → log, don't fail the call | The correct result still reaches the caller even if persistence fails |
| Lazy `MkdirAll` on first write | No setup step required before first use |

## Best Practices

- Hash content, not paths — paths change identity, content doesn't.
- Stream large files into the hasher (`io.Copy`) instead of reading them fully into memory first.
- Keep the processing function pure — it should have zero knowledge that caching exists.
- Log cache hit/miss with a truncated hash prefix for debugging without flooding logs.
- Treat any cache-read error (missing file, bad JSON, wrong shape) as a miss, never as a fatal error.

## Gotchas

- Two files with identical content but different intended outputs (e.g. same PDF processed with two different extraction configs) will collide on the same cache key if the config isn't part of what's hashed — if extraction behavior depends on parameters beyond file bytes, fold those parameters into the key (e.g. hash the file content concatenated with a serialized config) or this pattern silently returns the wrong result for the second config.
- A cache entry for a very large result (e.g. full-page OCR output for a 500-page PDF) written as a single JSON file works but doesn't scale indefinitely — past some size, prefer streaming the cached artifact to disk directly rather than round-tripping it through JSON.
- This pattern is wrong for anything that must always be fresh (a live price feed, a real-time API response) — the entire point of the pattern is durable reuse of a result tied to unchanging content, which doesn't apply to data that's expected to change independent of any file.

## Anti-Patterns to Avoid

```go
// BAD: path-based caching — breaks the moment a file is renamed or moved
cache := map[string]Document{"/path/to/file.pdf": result}

// BAD: caching concern baked into the pure function — now it has two
// responsibilities and can't be tested or reused without a cache
func ExtractText(path string, cacheEnabled bool, cacheDir string) (Document, error) {
    if cacheEnabled { /* ... */ }
    // ...
}
```

## Real-world grounding

Git's object model uses exactly this idea at the core of the entire system: every blob, tree, and commit is addressed by the SHA hash of its own content, not by a file path. That's precisely why renaming a file with unchanged content doesn't create a new object in Git's store — the object is identified by what it contains, not where it lives. This pattern applies the same principle to a processing cache instead of a version-control object store.

## Verification

- [ ] The cache key is a hash of file content, never the file path
- [ ] Large files are hashed via streaming (`io.Copy`), not loaded fully into memory
- [ ] The processing function has no cache-related parameters or branches
- [ ] A missing or corrupt cache file degrades to a cache miss, not a crash
- [ ] Parameters that affect output (beyond file bytes) are folded into the cache key if they vary
