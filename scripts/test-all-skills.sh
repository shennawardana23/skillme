#!/usr/bin/env bash
# Runs the live smeval eval for every skill in skills/ and reports a
# pass/fail summary. Real model calls via the local `claude` CLI - this
# is slow and makes real API calls, budget real time for it.
#
# Resumable: a skill already recorded in the results file is skipped on a
# re-run, so an interrupted run can just be re-invoked. Use -fresh to
# discard prior results and test everything again.
#
# Usage:
#   scripts/test-all-skills.sh                # test every skill, resuming
#   scripts/test-all-skills.sh -fresh         # discard prior results, test everything
#   scripts/test-all-skills.sh -only=<name>   # test just one skill
set -uo pipefail

cd "$(dirname "$0")/.."

RESULTS_DIR="smeval-workspace/test-all-results"
RESULTS_FILE="$RESULTS_DIR/summary.txt"
FRESH=false
ONLY=""

for arg in "$@"; do
  case "$arg" in
    -fresh) FRESH=true ;;
    -only=*) ONLY="${arg#-only=}" ;;
    -h|--help)
      sed -n '2,16p' "$0" | sed 's/^# \{0,1\}//'
      exit 0
      ;;
  esac
done

mkdir -p "$RESULTS_DIR"
if [ "$FRESH" = true ]; then
  rm -f "$RESULTS_FILE"
fi
touch "$RESULTS_FILE"

echo "Building smeval..."
go build -o smeval ./cmd/smeval

total=0
passed_full=0
partial=0
errored=0
flagged=()

for dir in skills/*/; do
  name=$(basename "$dir")
  [ -f "${dir}evals/evals.json" ] || continue
  if [ -n "$ONLY" ] && [ "$name" != "$ONLY" ]; then
    continue
  fi
  if grep -q "^${name} :" "$RESULTS_FILE" 2>/dev/null; then
    continue
  fi

  total=$((total + 1))
  echo ""
  echo "=== $name ==="
  summary=$(./smeval run "$dir" 2>&1 | tee /dev/stderr)
  line=$(echo "$summary" | grep -o '[0-9]*/[0-9]* cases fully passed' | head -1)

  if [ -z "$line" ]; then
    echo "${name} : ERROR :: $(echo "$summary" | tr '\n' ' ')" >> "$RESULTS_FILE"
    errored=$((errored + 1))
    flagged+=("$name (ERROR)")
  else
    echo "${name} : ${line}" >> "$RESULTS_FILE"
    pass_n="${line%%/*}"
    tot_n="${line#*/}"
    tot_n="${tot_n%% *}"
    if [ "$pass_n" = "$tot_n" ]; then
      passed_full=$((passed_full + 1))
    else
      partial=$((partial + 1))
      flagged+=("$name (${line})")
    fi
  fi
done

echo ""
echo "==================================================="
echo "Ran ${total} skill(s) this invocation (results accumulate in ${RESULTS_FILE})"
echo "  fully passed : ${passed_full}"
echo "  partial fail : ${partial}"
echo "  errored      : ${errored}"
echo "==================================================="

if [ "${#flagged[@]}" -gt 0 ]; then
  echo ""
  echo "Needs attention:"
  for f in "${flagged[@]}"; do
    echo "  - $f"
  done
fi

echo ""
echo "Full accumulated results: ${RESULTS_FILE}"
echo "Per-case evidence: smeval-workspace/runs/<name>/iteration-1/<case-id>/with_skill/{outputs/response.md,grading.json}"
