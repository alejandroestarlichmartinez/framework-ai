# Apply Progress: FullRick Preset and Rick Sanchez Persona

## Implementation Summary

**Change**: fullrick-preset  
**Mode**: Strict TDD  
**Date**: 2026-05-25  
**Status**: ✅ COMPLETE — All 6 phases implemented, all tests passing

---

## TDD Cycle Evidence

| Task | Test File | Layer | Safety Net | RED | GREEN | TRIANGULATE | REFACTOR |
|------|-----------|-------|------------|-----|-------|-------------|----------|
| 1.1 | `internal/model/types_test.go` | Unit | N/A (new) | Written | Passed | Single constant | None needed |
| 2.1-2.5 | `internal/assets/assets_test.go` | Unit | 22/22 passing | Written | Passed | Asset readability + content | None needed |
| 3.1 | `internal/components/persona/inject_test.go` | Unit | 40/40 passing | Written | Passed | 4 personas covered | None needed |
| 3.2 | `internal/components/persona/inject_test.go` | Unit | 40/40 passing | Written | Passed | Rick output-style + settings merge | None needed |
| 3.3 | `internal/components/persona/inject_test.go` | Unit | 40/40 passing | Written | Passed | Rick agent overlay creation | None needed |
| 3.4 | `internal/components/persona/inject_test.go` | Unit | 40/40 passing | Written | Passed | Rick managed sections skip | None needed |
| 3.5 | `internal/components/persona/inject_test.go` | Unit | 40/40 passing | Written | Passed | 14 agents × Rick dispatch | None needed |
| 3.6 | `internal/components/persona/inject_test.go` | Unit | 40/40 passing | Written | Passed | Idempotency verified | None needed |
| 4.1 | `internal/cli/validate_test.go` | Unit | 1/1 passing | Written | Passed | Rick persona accepted | None needed |
| 4.2 | `internal/cli/validate_test.go` | Unit | 1/1 passing | Written | Passed | FullRick preset + components match | None needed |
| 4.3 | `internal/cli/validate_test.go` | Unit | 1/1 passing | Written | Passed | Default persona derivation (3 cases) | None needed |
| 5.1 | `internal/tui/screens/persona_preset_test.go` | Unit | 2/2 passing | Written | Passed | Rick in options + description | None needed |
| 5.2 | `internal/tui/screens/persona_preset_test.go` | Unit | 2/2 passing | Written | Passed | FullRick in options + description | None needed |
| 5.3 | `internal/tui/model_test.go` | Unit | 50+/50+ passing | Written | Passed | Auto-select + preserve override | None needed |
| 5.4 | `internal/tui/model_test.go` | Unit | 50+/50+ passing | Written | Passed | TUI componentsForPreset parity | None needed |
| 6.1 | `internal/components/persona/inject_test.go` | Unit | 40/40 passing | Written | Passed | 3 golden files generated | None needed |
| 6.2 | `internal/components/persona/inject_test.go` | Unit | 40/40 passing | Written | Passed | Rick→Neutral cleanup verified | None needed |
| 6.3 | `internal/components/persona/inject_test.go` | Unit | 40/40 passing | Written | Passed | Comprehensive coverage (OpenCode, Kimi, Kiro) | None needed |
| 6.4 | `internal/cli/validate_test.go` | Unit | 1/1 passing | Written | Passed | Boundary tests (persona/preset confusion) | None needed |
| 6.5 | All packages | Integration | All passing | N/A | Passed | Full suite | None needed |

### Test Summary
- **Total tests written**: 25+ new tests
- **Total tests passing**: All tests across all packages
- **Layers used**: Unit (25+), Integration (full suite)
- **Approval tests** (refactoring): None — no refactoring tasks
- **Pure functions created**: 1 (`isFullPersona()`)

---

## Completed Tasks

### Phase 1: Types + Constants ✅
- [x] 1.1 — Add `PersonaRickSanchez` and `PresetFullRick` constants to `internal/model/types.go`
- [x] Tests added to `internal/model/types_test.go`

### Phase 2: Asset Files (All Agents) ✅
- [x] 2.1 — Create `internal/assets/generic/persona-rick-sanchez.md`
- [x] 2.2 — Create `internal/assets/claude/persona-rick-sanchez.md`
- [x] 2.3 — Create `internal/assets/opencode/persona-rick-sanchez.md`
- [x] 2.4 — Create `internal/assets/kimi/persona-rick-sanchez.md`
- [x] 2.5 — Create `internal/assets/kiro/persona-rick-sanchez.md`
- [x] All assets added to `internal/assets/assets_test.go` expected files list

### Phase 3: Injection Logic ✅
- [x] 3.1 — Introduce `isFullPersona()` helper in `inject.go` with tests
- [x] 3.2 — Refactor output-style gating to use `isFullPersona()` (Claude + Kimi)
- [x] 3.3 — Refactor agent overlay gating to use `isFullPersona()` (OpenCode/Kilocode)
- [x] 3.4 — Refactor `preserveManagedSections()` gating to use `isFullPersona()`
- [x] 3.5 — Extend `personaContent()` dispatch for Rick (all 14 agents)
- [x] 3.6 — Add Rick injection idempotency test
- [x] Updated `isExactLegacyPersonaAsset()` to include Rick assets

### Phase 4: Preset + Validation ✅
- [x] 4.1 — Update CLI persona validation (`normalizePersona`) to accept Rick
- [x] 4.2 — Update CLI preset validation (`normalizePreset`) to accept FullRick
- [x] 4.2 — Update `componentsForPreset()` in `validate.go` to handle FullRick
- [x] 4.3 — Derive default persona from preset (FullRick → Rick, FullGentleman → Gentleman)
- [x] Updated `NormalizeInstallFlags` to pass preset into persona normalization

### Phase 5: TUI Integration ✅
- [x] 5.1 — Add Rick to persona picker (`screens/persona.go`)
- [x] 5.2 — Add FullRick to preset picker (`screens/preset.go`)
- [x] 5.3 — Wire preset-to-persona default in TUI model (`model.go`)
- [x] 5.4 — Update TUI `componentsForPreset()` to handle FullRick
- [x] Updated golden files for TUI component selection screen

### Phase 6: Tests + Golden Files ✅
- [x] 6.1 — Create golden file infrastructure + 3 golden files (Claude, OpenCode, Kimi)
- [x] 6.2 — Add Rick switch-to-neutral cleanup tests (output-style + agent overlay)
- [x] 6.3 — Add comprehensive Rick injection coverage (Claude, OpenCode, Kimi, Kiro)
- [x] 6.4 — Add boundary tests for preset/persona confusion
- [x] 6.5 — Full test suite verification: `go test ./...` ✅, `go build ./...` ✅, `go vet ./...` ✅

---

## Files Changed

### Created
| File | Description |
|------|-------------|
| `internal/assets/generic/persona-rick-sanchez.md` | Generic Rick Sanchez persona asset |
| `internal/assets/claude/persona-rick-sanchez.md` | Claude-specific Rick Sanchez persona |
| `internal/assets/opencode/persona-rick-sanchez.md` | OpenCode/Kilocode Rick Sanchez persona |
| `internal/assets/kimi/persona-rick-sanchez.md` | Kimi Jinja module Rick Sanchez persona |
| `internal/assets/kiro/persona-rick-sanchez.md` | Kiro steering file Rick Sanchez persona |
| `internal/components/persona/testdata/claude-rick-sanchez.golden.md` | Golden file: Claude Rick injection |
| `internal/components/persona/testdata/opencode-rick-sanchez.golden.md` | Golden file: OpenCode Rick injection |
| `internal/components/persona/testdata/kimi-rick-sanchez.golden.md` | Golden file: Kimi Rick injection |

### Modified
| File | What Changed |
|------|--------------|
| `internal/model/types.go` | Added `PersonaRickSanchez` and `PresetFullRick` constants |
| `internal/model/types_test.go` | Tests for new constants |
| `internal/assets/assets_test.go` | Added Rick assets to expected files list |
| `internal/components/persona/inject.go` | `isFullPersona()`, Rick dispatch, gating refactors |
| `internal/components/persona/inject_test.go` | 15+ new tests for Rick + golden file infrastructure |
| `internal/catalog/components.go` | Updated persona component description |
| `internal/cli/validate.go` | Rick persona/preset validation, default persona from preset |
| `internal/cli/validate_test.go` | 5 new tests for Rick validation |
| `internal/tui/screens/persona.go` | Added Rick option |
| `internal/tui/screens/preset.go` | Added FullRick option |
| `internal/tui/screens/persona_preset_test.go` | 4 new TUI screen tests |
| `internal/tui/model.go` | Auto-select Rick persona when FullRick preset chosen |
| `internal/tui/model_test.go` | 3 new TUI model tests + cursor fix |
| `internal/tui/testdata/preset-custom-opencode-next.golden` | Updated persona description |
| `internal/tui/testdata/preset-custom-no-opencode-next.golden` | Updated persona description |

---

## Deviations from Design

**None** — implementation matches design exactly.

Key design decisions followed:
- `isFullPersona()` helper gates all output-style, agent-overlay, and managed-section logic
- Agent overlay key stays `"gentleman"` for backward compatibility (REQ-BWC-002)
- Output style files and settings keys retain existing names (`gentleman.md`, `"Gentleman"`)
- Rick persona uses neutral Spanish (no voseo, no regionalisms)
- `PresetFullRick` uses `fallthrough` to share component list with `PresetFullGentleman`

---

## Issues Found

1. **TUI golden file drift**: Adding Rick to persona descriptions caused 2 existing TUI golden files to fail (`preset-custom-opencode-next.golden` and `preset-custom-no-opencode-next.golden`). Fixed by updating the golden files to reflect the new catalog component description.

2. **TUI cursor index shift**: `TestPiCombinedWithOtherAgentsTUIInstallKeepsAllAgentsInPlan` used cursor=2 for Minimal preset, but Minimal moved from index 2 to 3 after adding FullRick at index 1. Fixed by updating the test to use cursor=3.

3. **Language guardrail test failure**: `TestGentlemanLanguageInstructionsDoNotBiasEnglishSessions` checks for exact string matches in persona assets. Rick assets initially used "cynical energy" instead of "warm energy" in the English language guardrail. Fixed by aligning Rick assets with the exact functional guardrail text from Gentleman.

---

## Verification Results

```
$ go test ./... -count=1
ALL PACKAGES PASS

$ go build ./...
NO ERRORS

$ go vet ./...
NO WARNINGS
```

---

## Remaining Tasks

None. All tasks from `tasks.md` are complete.

---

## Workload / PR Boundary

- **Mode**: single PR
- **Estimated changed lines**: ~300+ lines (well within 400-line budget)
- **Scope**: Complete feature — types, assets, injection logic, validation, TUI, tests, golden files

---

## Status

**24/24 tasks complete. Ready for verify phase.**
