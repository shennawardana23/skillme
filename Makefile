.DEFAULT_GOAL := help

SMEVAL      := ./smeval
SKILL       ?=
CASE        ?=
INCLUDE     ?=
THRESHOLD   ?=
RUNS_DIR     = smeval-workspace/runs/$(SKILL)

# ---- build / lint (mirrors .github/workflows/skill-eval.yml's "test" job) --

.PHONY: build
build: ## Build the smeval binary
	go build -o $(SMEVAL) ./cmd/smeval

.PHONY: vet
vet: ## Run go vet
	go vet ./...

.PHONY: fmt
fmt: ## Reformat all Go source with gofmt
	gofmt -w .

.PHONY: fmt-check
fmt-check: ## Fail if any Go file is not gofmt-formatted (CI check, no writes)
	@unformatted="$$(gofmt -l .)"; \
	if [ -n "$$unformatted" ]; then \
		echo "gofmt found unformatted files:"; echo "$$unformatted"; exit 1; \
	fi

.PHONY: unit-test
unit-test: ## Run this repo's own Go unit tests (internal/*, cmd/smeval)
	go test ./...

.PHONY: ci
ci: build vet unit-test fmt-check ## Full local baseline check, same steps as CI's "test" job

# ---- skill validation (schema/spec only, no model calls, no API key) -------

.PHONY: validate
validate: build ## Validate one skill's spec + evals.json schema. Usage: make validate SKILL=<name>
	@test -n "$(SKILL)" || (echo "usage: make validate SKILL=<name>"; exit 1)
	uvx --from skills-ref agentskills validate skills/$(SKILL)
	$(SMEVAL) validate skills/$(SKILL)

.PHONY: validate-all
validate-all: build ## Validate every skill that has an evals.json (mirrors CI's validate matrix)
	@for dir in skills/*/; do \
		name=$$(basename "$$dir"); \
		[ -f "$$dir/evals/evals.json" ] || continue; \
		echo "== $$name =="; \
		uvx --from skills-ref agentskills validate "$$dir" && \
		$(SMEVAL) validate "$$dir" || exit 1; \
	done

.PHONY: security-scan
security-scan: build ## Static risk scan of one skill (no model calls). Usage: make security-scan SKILL=<name>
	@test -n "$(SKILL)" || (echo "usage: make security-scan SKILL=<name>"; exit 1)
	$(SMEVAL) security-scan skills/$(SKILL)

.PHONY: security-scan-all
security-scan-all: build ## Static risk scan of every skill (mirrors CI's security-scan matrix)
	@failed=0; \
	for dir in skills/*/; do \
		[ -f "$$dir/SKILL.md" ] || continue; \
		$(SMEVAL) security-scan "$$dir" -quiet || failed=1; \
	done; \
	exit $$failed

.PHONY: similarity-check
similarity-check: build ## Advisory-only: report skill-description pairs likely to overlap (no model calls). Optional: THRESHOLD=0.NN
	$(SMEVAL) similarity-check $(if $(THRESHOLD),-threshold $(THRESHOLD))

# ---- live evals (real `claude` CLI calls, real cost/time) ------------------

.PHONY: eval
eval: build ## Run one skill's live eval. Usage: make eval SKILL=<name> [INCLUDE=<case-substring>]
	@test -n "$(SKILL)" || (echo "usage: make eval SKILL=<name> [INCLUDE=<case-substring>]"; exit 1)
	$(SMEVAL) run skills/$(SKILL) $(if $(INCLUDE),-include $(INCLUDE))

.PHONY: eval-benchmark
eval-benchmark: build ## Run one skill's eval with a without_skill baseline for comparison. Usage: make eval-benchmark SKILL=<name>
	@test -n "$(SKILL)" || (echo "usage: make eval-benchmark SKILL=<name>"; exit 1)
	$(SMEVAL) run skills/$(SKILL) -benchmark

.PHONY: test-all
test-all: ## Run the live eval for every skill in the catalog, resuming a prior interrupted run
	scripts/test-all-skills.sh

.PHONY: test-all-fresh
test-all-fresh: ## Same as test-all, discarding any prior accumulated results first
	scripts/test-all-skills.sh -fresh

.PHONY: test-all-only
test-all-only: ## Run test-all-skills.sh for just one skill. Usage: make test-all-only SKILL=<name>
	@test -n "$(SKILL)" || (echo "usage: make test-all-only SKILL=<name>"; exit 1)
	scripts/test-all-skills.sh -only=$(SKILL)

# ---- catalog-wide trigger-accuracy / coexistence testing -------------------
# See catalog-evals/README.md — installs the WHOLE catalog at once instead
# of one skill in isolation, so these cases test whether Claude picks the
# right skill among all of them, not per-skill correctness.

.PHONY: trigger-validate
trigger-validate: build ## Validate catalog-evals/trigger-accuracy.json's schema
	$(SMEVAL) trigger-validate catalog-evals/trigger-accuracy.json

.PHONY: trigger-run
trigger-run: build ## Run catalog-wide trigger-accuracy cases against the whole catalog. Optional: INCLUDE=<case-substring>
	$(SMEVAL) trigger-run catalog-evals/trigger-accuracy.json $(if $(INCLUDE),-include $(INCLUDE))
	@echo "Inspect with: make report SKILL=catalog-triggers  /  make cat SKILL=catalog-triggers CASE=<id>"

# ---- inspecting results -----------------------------------------------------

.PHONY: list
list: ## List every skill with its evals.json case count
	scripts/list-skills.sh

.PHONY: report
report: ## Open the latest COMPLETE iteration's report.html for a skill (web UI). Usage: make report SKILL=<name>
	@test -n "$(SKILL)" || (echo "usage: make report SKILL=<name>"; exit 1)
	@dir=$$(for d in $$(ls -d $(RUNS_DIR)/iteration-* 2>/dev/null | sort -t- -k2 -nr); do \
		if [ -f "$$d/report.html" ]; then echo "$$d"; break; fi; \
	done); \
	test -n "$$dir" || (echo "no completed run found under $(RUNS_DIR) (an interrupted run leaves no report.html) - run 'make eval SKILL=$(SKILL)' first"; exit 1); \
	echo "$$dir/report.html"; \
	(command -v open >/dev/null && open "$$dir/report.html") || \
	(command -v xdg-open >/dev/null && xdg-open "$$dir/report.html") || true

.PHONY: cat
cat: ## Print a case's response + grading from the latest iteration. Usage: make cat SKILL=<name> CASE=<case-id>
	@test -n "$(SKILL)" || (echo "usage: make cat SKILL=<name> CASE=<case-id>"; exit 1)
	@test -n "$(CASE)" || (echo "usage: make cat SKILL=<name> CASE=<case-id>"; exit 1)
	@dir=$$(ls -d $(RUNS_DIR)/iteration-* 2>/dev/null | sort -t- -k2 -n | tail -1); \
	test -n "$$dir" || (echo "no runs found under $(RUNS_DIR) - run 'make eval SKILL=$(SKILL)' first"; exit 1); \
	case_dir="$$dir/$(CASE)/with_skill"; \
	test -d "$$case_dir" || (echo "no case '$(CASE)' in $$dir - check the id in evals/evals.json"; exit 1); \
	echo "--- $$case_dir/outputs/response.md ---"; cat "$$case_dir/outputs/response.md"; \
	echo; echo "--- $$case_dir/grading.json ---"; cat "$$case_dir/grading.json"

.PHONY: catalog-report
catalog-report: ## Open the latest catalog-wide trigger-run's report.html (shortcut for report SKILL=catalog-triggers)
	@$(MAKE) --no-print-directory report SKILL=catalog-triggers

.PHONY: catalog-cat
catalog-cat: ## Print one catalog-wide trigger case's response + grading. Usage: make catalog-cat CASE=<case-id>
	@$(MAKE) --no-print-directory cat SKILL=catalog-triggers CASE=$(CASE)

.PHONY: benchmark-report
benchmark-report: ## Print the latest -benchmark comparison (with_skill vs without_skill). Usage: make benchmark-report SKILL=<name>
	@test -n "$(SKILL)" || (echo "usage: make benchmark-report SKILL=<name>"; exit 1)
	@dir=$$(ls -d $(RUNS_DIR)/iteration-* 2>/dev/null | sort -t- -k2 -n | tail -1); \
	test -n "$$dir" || (echo "no runs found under $(RUNS_DIR) - run 'make eval-benchmark SKILL=$(SKILL)' first"; exit 1); \
	test -f "$$dir/benchmark.json" || (echo "$$dir has no benchmark.json - it wasn't run with -benchmark; use 'make eval-benchmark SKILL=$(SKILL)'"; exit 1); \
	(command -v jq >/dev/null && jq . "$$dir/benchmark.json") || cat "$$dir/benchmark.json"

.PHONY: feedback
feedback: ## Open the latest iteration's feedback.json for human review. Usage: make feedback SKILL=<name>
	@test -n "$(SKILL)" || (echo "usage: make feedback SKILL=<name>"; exit 1)
	@dir=$$(ls -d $(RUNS_DIR)/iteration-* 2>/dev/null | sort -t- -k2 -n | tail -1); \
	test -n "$$dir" || (echo "no runs found under $(RUNS_DIR) - run 'make eval SKILL=$(SKILL)' first"; exit 1); \
	f="$$dir/feedback.json"; \
	test -f "$$f" || (echo "no feedback.json at $$f"; exit 1); \
	echo "$$f"; \
	if [ -n "$$EDITOR" ]; then "$$EDITOR" "$$f"; else cat "$$f"; fi

# ---- housekeeping ------------------------------------------------------------

.PHONY: clean
clean: ## Remove the smeval binary and all local eval workspaces (both gitignored, safe to delete)
	rm -f $(SMEVAL)
	rm -rf smeval-workspace

.PHONY: help
help: ## Show this list of targets
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-20s\033[0m %s\n", $$1, $$2}'
