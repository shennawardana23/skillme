---
name: prompt-optimizer
description: Analyze a draft prompt, diagnose gaps, match it to relevant skills already in this catalog, and output a ready-to-paste improved prompt — advisory only, never executes the task itself. Use when the user says "optimize this prompt", "improve my prompt", "how should I prompt for X", or pastes a draft prompt asking for feedback. Do not use when the user wants the task executed directly or says "just do it".
license: Apache-2.0
metadata:
  version: "0.1.0"
---

# Prompt Optimizer

Turn a vague or underspecified prompt into one that's specific enough to get a good result on the first try — without doing the task itself.

## When to Use

- "Optimize this prompt", "improve my prompt", "rewrite this prompt"
- "How should I ask for X", "what's the best way to prompt for..."
- A draft prompt is pasted with a request for feedback or enhancement

**Do not use when**: the user wants the task done directly ("just do it") — tell them this skill only produces an optimized prompt, and to make a normal task request instead if they want execution.

## Role: Advisory Only

Never write code, create files, run commands, or take any implementation action from within this skill. The only output is a diagnosis plus an optimized prompt the user can paste elsewhere to actually run the task.

## The Pipeline

Run these phases in order.

### Phase 0: Project Detection

Check whether a rules file (CLAUDE.md or equivalent) exists in the working directory and read it for conventions. Detect the tech stack from manifest files present: `go.mod` → Go, `package.json` → Node/TypeScript, `pyproject.toml`/`requirements.txt` → Python, `Cargo.toml` → Rust, and so on. If no manifest is found (the prompt is abstract or for a brand-new project), skip detection and flag "tech stack unknown" in Phase 4 rather than guessing one.

### Phase 1: Intent Detection

Classify the task into one or more categories: new feature, bug fix, refactor, research, testing, review, documentation, infrastructure, or design. The category shapes which phase-5 workflow steps matter most.

### Phase 2: Scope Assessment

Estimate size from the prompt and, if a project was detected, from the affected codebase: trivial (single file, <50 lines), low (single module), medium (multiple components, same domain), high (cross-domain, 5+ files), epic (multi-session, architectural). If no project context exists, estimate from the description alone and mark the estimate as uncertain.

### Phase 3: Existing-Skill Matching

Before recommending generic workflow steps, check whether a skill already in this repository's `skills/` directory matches the task — grep or list `skills/*/SKILL.md` names and descriptions rather than assuming none exists. If a matching skill is found, recommend invoking it by name instead of describing the process from scratch. If none matches, recommend the generic lifecycle instead: plan → implement → test → review → verify → commit, weighted toward whichever step the intent category most needs (e.g. research-classified tasks lean on the planning/investigation step, testing-classified tasks lean on the test step).

### Phase 4: Missing Context Detection

Check the prompt for whether each of these is present or needs to come from the user (mark auto-detected items as already answered by Phase 0):

- Tech stack
- Target scope (files, directories, modules)
- Acceptance criteria (how to know it's done)
- Error handling and edge cases
- Security requirements (auth, input validation, secrets)
- Testing expectations
- Performance constraints
- UI/UX requirements, if frontend
- Database changes, if data layer
- Existing patterns to follow
- Explicit scope boundaries (what NOT to do)

If 3 or more critical items are missing, ask up to 3 clarifying questions before producing the optimized prompt, then fold the answers in.

### Phase 5: Workflow Recommendation

State where the task sits in the lifecycle: research → plan → implement → test → review → verify → commit. For medium-or-larger scope, recommend writing a short plan before implementing. For epic/cross-session scope, recommend an explicit written plan document checked in before work starts, broken into phases with a verification step between each. Don't recommend a specific model version — model availability and naming change independently of this skill; instead give the size-based planning guidance above, which stays valid regardless of which model executes it.

## Output Format

Respond in the same language as the user's input.

### Section 1: Prompt Diagnosis
- **Strengths**: what the original prompt does well
- **Issues**: a table of problem / impact / suggested fix
- **Needs Clarification**: numbered questions, or "auto-detected: X" where Phase 0 already answered it

### Section 2: Recommended Skills and Workflow Steps
A table of skill-or-step / purpose, drawing on Phase 3's matches.

### Section 3: Optimized Prompt — Full Version
A single fenced code block: complete, self-contained, ready to paste. Include the task description with context, detected or specified tech stack, the workflow steps at the right stages, acceptance criteria, verification steps, and explicit scope boundaries (what NOT to do).

### Section 4: Optimized Prompt — Quick Version
A compact one- or two-line version for an experienced user who just wants the gist.

### Section 5: Enhancement Rationale
A table of what was added and why it matters.

### Footer
A one-line invitation to adjust, or to make a normal task request instead if execution (not optimization) is wanted.

## Gotchas

- Hardcoding a table of specific skill names is a stale-data trap the moment the catalog changes — always check the live `skills/` directory in Phase 3 rather than recommending a skill from memory that may no longer exist or may have been renamed.
- A prompt that looks complete because it's long is not the same as a prompt with acceptance criteria and scope boundaries — length and specificity are different axes; check Phase 4's list explicitly rather than judging completeness by word count.
- Recommending a plan-first workflow for a trivial one-file change adds overhead without benefit — match the workflow weight to the Phase 2 scope estimate, don't apply the epic-scope process uniformly.

## Real-world grounding

The core premise here — that being explicit about acceptance criteria, constraints, and what NOT to do measurably improves the first-pass quality of an AI coding assistant's output — is consistent with Anthropic's own published prompt-engineering guidance, which repeatedly emphasizes specificity and providing examples over open-ended requests. It's also just the well-established software-engineering principle that clarifying requirements before implementation is cheaper than discovering the gap during code review or after shipping — the same principle behind this catalog's `interview-me` skill, applied here specifically to the shape of the prompt itself rather than to a live back-and-forth conversation.

## Verification

- [ ] No code, file, or command was executed — only the diagnosis and optimized prompt were produced
- [ ] Phase 3's skill recommendations were checked against the actual `skills/` directory, not recalled from memory
- [ ] The optimized prompt includes acceptance criteria and explicit scope boundaries
- [ ] The workflow weight (plan-first vs. direct) matches the Phase 2 scope estimate
