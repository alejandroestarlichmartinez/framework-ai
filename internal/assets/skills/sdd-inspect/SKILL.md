---
name: sdd-inspect
description: "Audit project compliance against AGENTS.md, conventions, and git hygiene. Trigger: sdd inspect, /sdd-inspect, audit, check compliance."
disable-model-invocation: true
user-invocable: false
license: MIT
metadata:
  author: alejandroestarlichmartinez
  version: "2.0"
  delegate_only: true
---

## Purpose

You are a sub-agent responsible for PRE-COMMIT AUDIT and DIAGNOSTICS. You analyze code against the project's AGENTS.md rules, Go quality standards, git hygiene, and SDD compliance. You do NOT fix issues — you report them and, if failures are found, generate a correction proposal.

## What You Receive

From the orchestrator:
- **Mode**: `commit` or `all`
- **Project path**
- **Artifact store mode** (`engram | openspec | hybrid | none`)

## Execution and Persistence Contract

Read and follow `skills/_shared/persistence-contract.md` for mode resolution rules.

- If mode is `engram`:

  **Save your artifact(s)**:
  ```
  mem_save(
    title: "sdd-inspect/{timestamp-or-sha}",
    topic_key: "sdd-inspect/{timestamp-or-sha}",
    type: "architecture",
    project: "{project}",
    content: "{your full inspect report markdown}"
  )
  ```
  If a correction proposal is generated, ALSO save:
  ```
  mem_save(
    title: "sdd-inspect/{timestamp-or-sha}/proposal",
    topic_key: "sdd-inspect/{timestamp-or-sha}/proposal",
    type: "architecture",
    project: "{project}",
    content: "{your full correction proposal markdown}"
  )
  ```
  `topic_key` enables upserts — saving again updates, not duplicates.

  (See `skills/_shared/engram-convention.md` for full naming conventions.)
- If mode is `openspec`: Read and follow `skills/_shared/openspec-convention.md`. Save report to `openspec/inspect-report.md` and proposal to `openspec/inspect-proposal.md`.
- If mode is `hybrid`: Follow BOTH conventions — persist to Engram AND write `inspect-report.md` / `inspect-proposal.md` to filesystem.
- If mode is `none`: Return the report inline only. Never write files.

## What to Do

### Step 1: Load Skill Registry

**Do this FIRST, before any other work.**

1. Try engram first: `mem_search(query: "skill-registry", project: "{project}")` → if found, `mem_get_observation(id)` for the full registry
2. If engram not available or not found: read `.atl/skill-registry.md` from the project root
3. If neither exists: proceed without skills (not an error)

From the registry, identify and read any skills whose triggers match your task. Also read any project convention files listed in the registry.

### Step 2: Load AGENTS.md and Project Conventions

Search for and read the project's rule files:

```
Project root:
├── AGENTS.md           ← Primary rules source
├── agents.md           ← Fallback
├── .agent/rules/       ← Rule fragments
├── persona.md          ← Persona definitions
└── CLAUDE.md           ← Claude-specific rules
```

Parse the rules into categorized checks:

| Section | Example Rules |
|---------|---------------|
| General | Conventional commits, no AI attribution, verify-before-stating, stop-and-wait |
| Go | Doc comments on exports, error handling, no naked returns |
| TypeScript | Strict mode, explicit return types, no `any` |
| Python | Type hints, PEP 8, docstrings |

**If no AGENTS.md (or equivalent) is found, flag this as CRITICAL** and suggest creating one.

### Step 3: Git Analysis

Run git commands to understand the change scope:

**If mode == `commit`:**
```bash
git diff --staged                     # Staged changes
git diff                              # Unstaged changes
git log --oneline -10                 # Recent commits for context
git status --short                    # File state overview
```

**If mode == `all`:**
```bash
git log --all --oneline --format="%h %s" -50   # Recent commit history
git ls-files                                     # Full tracked file tree
git status --short                               # Untracked/modified files
```

Capture:
- Files changed (names and approximate line counts)
- Commit messages (for conventional commit compliance)
- Author lines (for Co-Authored-By detection)

### Step 4: AGENTS.md Compliance Checks

For each rule found in Step 2, analyze the code from Step 3:

| Check | Method | Severity |
|-------|--------|----------|
| Conventional commits format | Regex against commit messages | CRITICAL if violated |
| No AI attribution (Co-Authored-By) | Search for `Co-Authored-By`, `Generated-by`, `AI-authored` | CRITICAL if found |
| Tool preferences respected (rg vs grep, bat vs cat) | Check changed files for discouraged tool usage | WARNING |
| Verify-before-stating pattern | Check for unverified claims in comments/docs | SUGGESTION |
| Stop-and-wait pattern | Check for TODO/FIXME without context | WARNING |

Flag each finding with:
- **File path** (if applicable)
- **Line or commit SHA**
- **Rule violated**
- **Suggested fix**

### Step 5: Go Code Quality (if Go project detected)

Detect Go by presence of `go.mod`. If found, run these checks on changed files (mode=`commit`) or all `.go` files (mode=`all`):

| Check | Method | Severity |
|-------|--------|----------|
| Exported functions have doc comments | AST or regex: `^func [A-Z]` without preceding `// ` | WARNING |
| No ignored errors (`_ = err`) | Search for `_ = ` assignments involving `error` | WARNING |
| No naked returns | Search for bare `return` in named return functions | WARNING |
| `go vet` equivalent | Run `go vet ./...` (mode=`all`) or `go vet {changed packages}` | CRITICAL if fails |
| `gofmt` compliance | Run `gofmt -d` on changed files | WARNING if non-empty diff |
| Correct import paths | Verify imports resolve (no broken paths) | CRITICAL if broken |

Capture:
- Command output
- Affected file paths
- Line numbers (if available)

### Step 6: SDD Compliance (if `.atl/` or `openspec/` exists)

Check SDD infrastructure health:

| Check | Method | Severity |
|-------|--------|----------|
| Skill registry present | Does `.atl/skill-registry.md` exist? | WARNING if missing |
| Skills have valid frontmatter | Parse YAML frontmatter from all `*/SKILL.md` | WARNING if invalid |
| Topic keys follow convention | Verify `sdd/{change}/{artifact}` pattern in engram refs | WARNING if violated |
| No broken references | Check that referenced files in specs/proposals exist | CRITICAL if broken |
| OpenSpec structure valid | If `openspec/` exists, verify `config.yaml`, `specs/`, `changes/` | WARNING if malformed |

### Step 7: File Hygiene

Scan all relevant files for common hygiene issues:

| Check | Method | Severity |
|-------|--------|----------|
| No secrets committed | Search for `.env`, `credentials.json`, `*.pem`, `SECRET`, `API_KEY` | CRITICAL if found |
| No merge conflict markers | Search for `<<<<<<<`, `=======`, `>>>>>>>` | CRITICAL if found |
| No trailing whitespace | `git diff --check` or regex `[ \t]+$` | WARNING |
| No TODO without context | Search `TODO` / `FIXME` without following explanation | WARNING |
| No binaries committed | Check for common binary extensions or `git diff --numstat` showing `-` | WARNING |

### Step 8: Generate Report

Produce a markdown report with these sections:

```markdown
# Inspect Report: {project}

**Mode**: {commit|all}
**Timestamp**: {ISO date}
**Scope**: {files analyzed}

---

## AGENTS.md Compliance

| Check | Status | Details |
|-------|--------|---------|
| Conventional commits | ✅/❌/⚠️ | {details} |
| No AI attribution | ✅/❌/⚠️ | {details} |
| Tool preferences | ✅/❌/⚠️ | {details} |

## Go Quality

| Check | Status | Details |
|-------|--------|---------|
| Exported docs | ✅/❌/⚠️ | {details} |
| Error handling | ✅/❌/⚠️ | {details} |
| go vet | ✅/❌/⚠️ | {details} |
| gofmt | ✅/❌/⚠️ | {details} |

## Git Hygiene

| Check | Status | Details |
|-------|--------|---------|
| Commit format | ✅/❌/⚠️ | {details} |

## SDD Compliance

| Check | Status | Details |
|-------|--------|---------|
| Skill registry | ✅/❌/⚠️ | {details} |
| Valid frontmatter | ✅/❌/⚠️ | {details} |

## File Hygiene

| Check | Status | Details |
|-------|--------|---------|
| Secrets | ✅/❌/⚠️ | {details} |
| Merge conflicts | ✅/❌/⚠️ | {details} |
| Trailing whitespace | ✅/❌/⚠️ | {details} |
| TODO context | ✅/❌/⚠️ | {details} |
| Binaries | ✅/❌/⚠️ | {details} |

---

## Summary Score

**{N}/{Total} checks passed**

- ❌ CRITICAL: {N}
- ⚠️ WARNING: {N}
- 💡 SUGGESTION: {N}
```

### Step 9: Generate Correction Proposal (if failures found)

If any ❌ or ⚠️ findings exist, create a mini-proposal:

```markdown
# Correction Proposal: Compliance Fixes

## Issues Found

### CRITICAL
1. **{issue title}** — {file/location}
   - {description}
   - **Suggested task**: {concrete fix}

### WARNING
1. **{issue title}** — {file/location}
   - {description}
   - **Suggested task**: {concrete fix}

### SUGGESTION
1. **{issue title}** — {file/location}
   - {description}
   - **Suggested task**: {concrete fix}

## Recommended Actions

- Run `/sdd-new fix-compliance` to create a full change for these fixes
- Or apply fixes manually and re-run `/sdd-inspect`
```

Save this proposal per the persistence contract (Engram topic_key: `sdd-inspect/{timestamp-or-sha}/proposal`, or `openspec/inspect-proposal.md`).

### Step 10: Persist & Return

**This step is MANDATORY — do NOT skip it.**

Persist the report (and proposal, if generated) according to the resolved `artifact_store.mode`:

- **engram**: Use `engram-convention.md` — artifact type `inspect-report` and optionally `inspect-proposal`
- **openspec**: Write to `openspec/inspect-report.md` (and `openspec/inspect-proposal.md` if applicable)
- **none**: Return the full report inline, do NOT write any files

Return to the orchestrator:

```markdown
## Inspect Report

**Project**: {project}
**Mode**: {commit|all}
**Scope**: {N} files analyzed

---

### Summary Score
**{N}/{Total} checks passed**

- ❌ CRITICAL: {N}
- ⚠️ WARNING: {N}
- 💡 SUGGESTION: {N}

---

### AGENTS.md Compliance
{summary or "All checks passed"}

### Go Quality
{summary or "All checks passed"}

### Git Hygiene
{summary or "All checks passed"}

### SDD Compliance
{summary or "All checks passed"}

### File Hygiene
{summary or "All checks passed"}

---

### Correction Proposal
{If failures found: link/location of proposal + suggestion to run `/sdd-new fix-compliance`}
{If no failures: "No issues found. Project is compliant."}

---

### Verdict
{PASS / PASS WITH WARNINGS / FAIL}

{One-line summary}
```

## Rules

- **Do NOT fix issues — only report.** The orchestrator decides what to do.
- **Be objective** — state what IS, not what should be.
- **CRITICAL** = blocks commit or indicates a serious violation.
- **WARNING** = should fix but won't block.
- **SUGGESTION** = improvements, not blockers.
- If AGENTS.md (or equivalent) is not found, report as **CRITICAL** and suggest creating one.
- If mode is `commit`, focus checks on changed files and recent commits.
- If mode is `all`, scan the entire project including history and untracked files.
- Return a structured envelope with: `status`, `executive_summary`, `detailed_report` (optional), `artifacts`, `next_recommended`, and `risks`
