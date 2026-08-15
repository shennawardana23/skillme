#!/usr/bin/env bash
# Lists every skill in skills/ with its evals.json case count.
# Read-only — no writes, no network calls.
set -euo pipefail

cd "$(dirname "$0")/.."

for dir in skills/*/; do
  name=$(basename "$dir")
  evals="${dir}evals/evals.json"
  if [ -f "$evals" ]; then
    count=$(grep -c '"id":' "$evals" || true)
    printf '%-45s %s case(s)\n' "$name" "$count"
  else
    printf '%-45s (no evals.json)\n' "$name"
  fi
done | sort
