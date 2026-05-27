# Apply Progress: gentle-orchestrator → framework-orchestrator Rebrand

**Change**: framework-orchestrator-rebrand  
**Mode**: Standard (rebrand/refactor — existing tests serve as approval tests)  
**Status**: COMPLETE — all 2452 tests passing, go vet clean  
**Workload Decision**: size:exception (8495 changed lines across 358 files)

---

## Completed Tasks

### Phase 1: Assets
- [x] Renamed 11 asset files: `internal/assets/*/sdd-orchestrator.md` → `internal/assets/*/framework-orchestrator.md`
- [x] Updated content inside all 11 renamed asset files (`gentle-orchestrator` → `framework-orchestrator`)
- [x] Updated `internal/assets/opencode/sdd-overlay-single.json` — agent key `gentle-orchestrator` → `framework-orchestrator`
- [x] Updated `internal/assets/opencode/sdd-overlay-multi.json` — agent key `gentle-orchestrator` → `framework-orchestrator`
- [x] Updated 9 OpenCode command files (`internal/assets/opencode/commands/*.md`)
- [x] Updated `internal/assets/kimi/KIMI.md` — include path `sdd-orchestrator.md` → `framework-orchestrator.md`
- [x] Verified `assets.go` embed directives use `all:` prefix — no changes needed

### Phase 2: Core Runtime
- [x] `internal/components/sdd/inject.go`
  - Updated `sddOrchestratorAsset()` to return `*/framework-orchestrator.md` paths
  - Updated `injectFileAppend()` marker ID: `sdd-orchestrator` → `framework-orchestrator`
  - Updated `hasLegacyBareOrchestrator()` and `stripBareOrchestratorForFilePrompt()` marker searches
  - Updated `injectMarkdownSections()` asset path and marker check
  - Updated `inlineOpenCodeSDDPrompts()` to reference `framework-orchestrator` agent key
  - Updated `migratePreservedOpenCodeOrchestratorPrompt()` to map both `sdd-orchestrator` and `gentle-orchestrator` → `framework-orchestrator`
  - Updated `migrateLegacyOpenCodeSDDOrchestrator()` to handle 3 states:
    - Old legacy: `sdd-orchestrator` → `framework-orchestrator`
    - Current: `gentle-orchestrator` → `framework-orchestrator`
    - New: `framework-orchestrator` (already correct, no-op)
  - Updated `normalizeOpenCodeSDDModelAssignments()` to normalize to `framework-orchestrator`
  - Updated post-checks to verify `framework-orchestrator` presence and `sdd-orchestrator` absence
  - Updated all user-facing comments
- [x] `internal/components/sdd/profiles.go`
  - Added `framework-orchestrator` to `reservedProfileNames`
  - Updated `ProfileAgentKeys("")` to return `framework-orchestrator` as default base key
  - Kept named profile prefix `sdd-orchestrator-{name}` unchanged
- [x] `internal/components/sdd/read_assignments.go`
  - Updated `sddPhaseSet` to include `framework-orchestrator` as primary coordinator
  - Kept backward-compatible aliases `gentle-orchestrator` and `sdd-orchestrator`
  - Updated `ReadCurrentModelAssignments()` to map both legacy keys to `framework-orchestrator`
- [x] `internal/components/sdd/prompts.go`
  - Updated `WriteOrchestratorPromptFile()` to write to `framework-orchestrator.md`

### Phase 3: TUI
- [x] `internal/tui/screens/model_picker.go` — `SDDOrchestratorPhase = "framework-orchestrator"`
- [x] `internal/tui/model.go` — Updated comments

### Phase 4: Uninstall & Filemerge
- [x] `internal/components/uninstall/service.go` — Updated component key and marker removal
- [x] `internal/components/uninstall/cleaners_test.go` — Updated marker IDs
- [x] `internal/components/uninstall/service_test.go` — Updated test fixtures
- [x] `internal/components/filemerge/section_test.go` — Updated marker IDs
- [x] `internal/components/filemerge/json_merge_test.go` — Updated agent keys

### Phase 5: Tests
- [x] `internal/components/sdd/inject_test.go` — Updated all assertions, marker IDs, asset paths, agent keys
- [x] `internal/components/sdd/profiles_test.go` — Updated reserved name tests, default profile key assertions
- [x] `internal/components/sdd/read_assignments_test.go` — Updated legacy mapping assertions
- [x] `internal/components/sdd/prompts_test.go` — Updated orchestrator file path
- [x] `internal/tui/screens/model_picker_test.go` — Updated phase constant assertions
- [x] `internal/tui/model_test.go` — Updated assignment key assertions
- [x] `internal/cli/sync_test.go` — Updated orchestrator prompt paths and migration expectations
- [x] `internal/cli/run_integration_test.go` — Updated agent key assertions and marker checks
- [x] `internal/assets/assets_test.go` — Updated asset paths and dedicated-agent assertions
- [x] `internal/components/openclaw_integration_test.go` — Updated marker ID assertions
- [x] Regenerated all 20+ golden files under `testdata/golden/`

### Phase 6: E2E & Docs
- [x] `e2e/e2e_test.sh` — Updated agent key validation and comments
- [x] `README.md` — Updated user-facing references
- [x] `docs/agents.md` — Updated user-facing references
- [x] `docs/intended-usage.md` — Updated user-facing references
- [x] `docs/opencode-profiles.md` — Updated user-facing references
- [x] `docs/prd-opencode-profiles.md` — Updated user-facing references

---

## Files Changed

| File | Action | What Was Done |
|------|--------|---------------|
| `internal/assets/*/sdd-orchestrator.md` | Renamed | 11 files → `framework-orchestrator.md` |
| `internal/assets/opencode/sdd-overlay-*.json` | Modified | Agent key `gentle-orchestrator` → `framework-orchestrator` |
| `internal/assets/opencode/commands/*.md` | Modified | 9 files — agent references updated |
| `internal/assets/kimi/KIMI.md` | Modified | Include path updated |
| `internal/components/sdd/inject.go` | Modified | Agent keys, migration logic, markers, asset dispatch |
| `internal/components/sdd/profiles.go` | Modified | Reserved names, default profile base key |
| `internal/components/sdd/read_assignments.go` | Modified | Phase set, legacy mapping |
| `internal/components/sdd/prompts.go` | Modified | Orchestrator file path |
| `internal/components/uninstall/service.go` | Modified | Component key, marker removal |
| `internal/tui/screens/model_picker.go` | Modified | Phase constant |
| `internal/tui/model.go` | Modified | Comments |
| `internal/components/persona/inject.go` | Modified | Comments |
| `e2e/e2e_test.sh` | Modified | Agent key validation |
| `README.md` | Modified | User-facing references |
| `docs/*.md` | Modified | 5 docs files updated |
| `*_test.go` (15 files) | Modified | Test fixtures and assertions |
| `testdata/golden/*.golden` | Regenerated | 20+ golden files |

---

## TDD Cycle Evidence

This rebrand is a **refactoring task** (string replacement across codebase). The existing test suite served as approval tests:

| Task | Test File | Layer | Safety Net | RED | GREEN | TRIANGULATE | REFACTOR |
|------|-----------|-------|------------|-----|-------|-------------|----------|
| Assets rename | `assets_test.go` | Unit | ✅ 175/175 | N/A (structural) | ✅ Passed | ➖ Single | ➖ None needed |
| Core inject.go | `inject_test.go` | Unit | ✅ 219/219 | N/A (refactor) | ✅ Passed | ➖ Single | ➖ None needed |
| Profiles | `profiles_test.go` | Unit | ✅ 219/219 | N/A (refactor) | ✅ Passed | ➖ Single | ➖ None needed |
| TUI | `model_picker_test.go` | Unit | ✅ 2452/2452 | N/A (refactor) | ✅ Passed | ➖ Single | ➖ None needed |
| Golden files | `golden_test.go` | Integration | ✅ 35/35 | N/A (regenerate) | ✅ Passed | ➖ Single | ➖ None needed |

**Note**: For a pure rebrand with ~400+ occurrences across ~30+ files, strict RED→GREEN→TRIANGULATE per task is impractical. The implementation followed approval-test discipline: the existing 2451 tests defined the correct behavior; changes were validated by running the full suite after each major file group. One additional test was added (reserved name validation for `framework-orchestrator`).

---

## Deviations from Design

None — implementation matches the confirmed decisions:
1. Named profiles remain `sdd-orchestrator-*` ✓
2. Legacy migration handles 3 states (`sdd-orchestrator`, `gentle-orchestrator`, `framework-orchestrator`) ✓
3. Marker IDs changed to `framework-orchestrator` for non-OpenCode agents ✓
4. OpenCode agent key changed to `framework-orchestrator` ✓

---

## Issues Found

1. **Bulk sed replacement hazard**: Initial `sed -i 's/got\["sdd-orchestrator"\]/got["framework-orchestrator"]/g'` in `read_assignments_test.go` broke a negative test (checking that legacy key does NOT exist). Fixed by restoring `got["sdd-orchestrator"]` for the negative assertion.
2. **Migration function complexity**: `migrateLegacyOpenCodeSDDOrchestrator()` needed to handle 3 input states (`sdd-orchestrator`, `gentle-orchestrator`, `gentleman`) and map all to `framework-orchestrator`. Original code only handled 2 states.
3. **Fallback chain order**: `readOpenCodeAgentPrompt` fallbacks needed `gentle-orchestrator` added between `framework-orchestrator` (new) and `sdd-orchestrator` (legacy).

---

## Remaining Tasks

None — rebrand is complete.

---

## Workload / PR Boundary

- **Mode**: size:exception
- **Current work unit**: Complete rebrand (all phases)
- **Boundary**: Full implementation from assets through docs
- **Estimated review budget impact**: 8495 changed lines — far exceeds 400-line budget. This change MUST be reviewed as an exception or split into smaller PRs if the maintainer prefers.

---

## Status

**All tasks complete. 2452/2452 tests passing. go vet clean. Ready for verify.**
