package engine

import (
	"context"
	"os"
	"path/filepath"
	"runtime"
	"testing"
	"time"
)

// writeFakeClaude writes a stub "claude" binary (a shell script) that
// fails for one model name and succeeds for another, so the outer
// (cross-process) fallback path in Run can be exercised deterministically
// — without depending on a real model actually being unavailable.
func writeFakeClaude(t *testing.T, failModel, okModel string) string {
	t.Helper()
	if runtime.GOOS == "windows" {
		t.Skip("fake claude stub is a POSIX shell script")
	}
	dir := t.TempDir()
	path := filepath.Join(dir, "fake-claude")
	script := `#!/bin/sh
model=""
prev=""
for arg in "$@"; do
  if [ "$prev" = "--model" ]; then
    model="$arg"
  fi
  prev="$arg"
done
if [ "$model" = "` + failModel + `" ]; then
  echo '{"is_error":true,"result":"boom","api_error_status":404,"terminal_reason":"api_error"}'
  exit 1
fi
if [ "$model" = "` + okModel + `" ]; then
  echo '{"is_error":false,"result":"OK from ` + okModel + `","duration_ms":42,"total_cost_usd":0.01,"usage":{"input_tokens":10,"output_tokens":5,"cache_creation_input_tokens":0,"cache_read_input_tokens":0}}'
  exit 0
fi
echo '{"is_error":true,"result":"unexpected model"}'
exit 1
`
	if err := os.WriteFile(path, []byte(script), 0o755); err != nil {
		t.Fatalf("write fake claude: %v", err)
	}
	return path
}

func TestRun_PrimarySucceeds_NoFallback(t *testing.T) {
	old := ClaudeBin
	ClaudeBin = writeFakeClaude(t, "never-used", "sonnet")
	defer func() { ClaudeBin = old }()

	res, err := Run(context.Background(), Options{
		Prompt:        "hi",
		PrimaryModel:  "sonnet",
		FallbackModel: "opus",
		Timeout:       5 * time.Second,
	})
	if err != nil {
		t.Fatalf("Run returned error: %v", err)
	}
	if res.FallbackUsed {
		t.Fatal("FallbackUsed should be false when the primary attempt succeeds")
	}
	if res.Attempts != 1 {
		t.Fatalf("Attempts = %d, want 1", res.Attempts)
	}
	if res.ModelUsed != "sonnet" {
		t.Fatalf("ModelUsed = %q, want %q", res.ModelUsed, "sonnet")
	}
}

// TestRun_PrimaryFails_OuterFallbackEngages proves the cross-process
// fallback layer: the primary model's process attempt fails at the engine
// level (non-zero exit, is_error:true), and Run retries in a fresh process
// against FallbackModel, succeeding there.
func TestRun_PrimaryFails_OuterFallbackEngages(t *testing.T) {
	old := ClaudeBin
	ClaudeBin = writeFakeClaude(t, "broken-model", "opus")
	defer func() { ClaudeBin = old }()

	res, err := Run(context.Background(), Options{
		Prompt:        "hi",
		PrimaryModel:  "broken-model",
		FallbackModel: "opus",
		Timeout:       5 * time.Second,
	})
	if err != nil {
		t.Fatalf("Run returned error: %v", err)
	}
	if !res.FallbackUsed {
		t.Fatal("FallbackUsed should be true when the primary attempt fails at the engine level")
	}
	if res.Attempts != 2 {
		t.Fatalf("Attempts = %d, want 2", res.Attempts)
	}
	if res.ModelUsed != "opus" {
		t.Fatalf("ModelUsed = %q, want %q", res.ModelUsed, "opus")
	}
	if res.FinalMessage != "OK from opus" {
		t.Fatalf("FinalMessage = %q, want the fallback model's response", res.FinalMessage)
	}
}

func TestRun_BothFail_ReturnsError(t *testing.T) {
	old := ClaudeBin
	ClaudeBin = writeFakeClaude(t, "broken-model", "also-broken-but-unused")
	defer func() { ClaudeBin = old }()

	_, err := Run(context.Background(), Options{
		Prompt:        "hi",
		PrimaryModel:  "broken-model",
		FallbackModel: "another-broken-model",
		Timeout:       5 * time.Second,
	})
	if err == nil {
		t.Fatal("Run should return an error when both primary and fallback attempts fail")
	}
}

func TestRun_NoFallbackConfigured_PrimaryFailure_ReturnsError(t *testing.T) {
	old := ClaudeBin
	ClaudeBin = writeFakeClaude(t, "broken-model", "sonnet")
	defer func() { ClaudeBin = old }()

	_, err := Run(context.Background(), Options{
		Prompt:       "hi",
		PrimaryModel: "broken-model",
		Timeout:      5 * time.Second,
	})
	if err == nil {
		t.Fatal("Run should return an error when the primary fails and no fallback model is configured")
	}
}
