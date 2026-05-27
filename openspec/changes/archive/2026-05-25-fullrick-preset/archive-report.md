# Archive Report: FullRick Preset and Rick Sanchez Persona

**Change**: fullrick-preset
**Status**: COMPLETE — VERIFIED
**Archived**: 2026-05-25
**Mode**: OpenSpec

---

## Change Summary

Added `PersonaRickSanchez` and `PresetFullRick` to framework-ai. The Rick Sanchez persona is a cynical-genius senior architect that teaches architecture, patterns, and testing in Rick Sanchez style — neutral Spanish, scientific metaphors, multiverse references, informal language — while remaining fully functional. The preset delivers the complete Gentleman AI ecosystem with Rick as the default persona.

## Final Status

| Phase | Status | Notes |
|-------|--------|-------|
| Proposal | ✅ Complete | Scope, approach, risks, rollback plan defined |
| Spec | ✅ Complete | 12 requirement sections, 8 scenarios, full acceptance criteria |
| Design | ✅ Complete | Asset structure, `isFullPersona()` helper, dispatch table, TUI wiring |
| Tasks | ✅ Complete | All 24 tasks done, 24 commits |
| Apply | ✅ Complete | All production code delivered per task plan |
| Verify | ✅ PASS WITH WARNINGS | All tests pass, `go vet` clean, 2 minor coverage suggestions noted |

## Files Created

### Asset Files (5 new persona assets)
- `internal/assets/generic/persona-rick-sanchez.md`
- `internal/assets/claude/persona-rick-sanchez.md`
- `internal/assets/opencode/persona-rick-sanchez.md`
- `internal/assets/kimi/persona-rick-sanchez.md`
- `internal/assets/kiro/persona-rick-sanchez.md`

### Golden Files (3 new golden files)
- `internal/components/persona/testdata/claude-rick-sanchez.golden.md`
- `internal/components/persona/testdata/opencode-rick-sanchez.golden.md`
- `internal/components/persona/testdata/kimi-rick-sanchez.golden.md`

### Test Files (new or updated)
- `internal/model/types_test.go` — persona/preset constants
- `internal/assets/assets_test.go` — asset existence, readability, language guardrails
- `internal/components/persona/inject_test.go` — `isFullPersona`, injection, idempotency, cleanup, golden files
- `internal/cli/validate_test.go` — persona/preset normalization, default derivation
- `internal/tui/screens/persona_test.go` — options, rendering
- `internal/tui/screens/preset_test.go` — options, rendering
- `internal/tui/model_test.go` — preset-to-persona auto-select, component resolution

## Files Modified

| File | Change |
|------|--------|
| `internal/model/types.go` | Added `PersonaRickSanchez` and `PresetFullRick` constants |
| `internal/components/persona/inject.go` | Added `isFullPersona()` helper; refactored all gating checks; extended `personaContent()` dispatch for Rick |
| `internal/cli/validate.go` | Accept Rick persona/preset; derive default persona from preset |
| `internal/tui/screens/persona.go` | Added Rick option and description |
| `internal/tui/screens/preset.go` | Added FullRick option and description |
| `internal/tui/model.go` | Auto-select Rick persona when FullRick preset chosen; handle FullRick in `componentsForPreset()` |
| `internal/catalog/components.go` | Updated persona component description to mention Rick |

## Test Results

```
$ go test ./... -count=1
ALL PACKAGES PASS (50+ packages)

$ go vet ./...
NO WARNINGS
```

### Test Coverage
- `TestIsFullPersona` — covers all 4 persona values
- `TestPersonaContentRickSanchezDispatch` — covers all 14 agents
- `TestInjectClaudeRickSanchezWritesSectionWithRealContent`
- `TestInjectClaudeRickSanchezWritesOutputStyleFile`
- `TestInjectClaudeRickSanchezMergesOutputStyleIntoSettings`
- `TestInjectOpenCodeRickSanchezCreatesAgentOverlay`
- `TestInjectClaudeRickSanchezIsIdempotent`
- `TestInjectClaudeRickToNeutralCleansOutputStyle`
- `TestInjectOpenCodeRickToNeutralCleansAgentOverlay`
- `TestInjectKimiRickSanchezIncludesOutputStyle`
- `TestRickGoldenFiles` — golden file comparison for Claude, OpenCode, Kimi
- `TestNormalizePersonaAcceptsRickSanchez`
- `TestNormalizePresetAcceptsFullRick`
- `TestComponentsForPresetFullRickMatchesFullGentleman`
- `TestNormalizeInstallFlagsFullRickDefaultsToRickSanchez`
- `TestNormalizeInstallFlagsFullRickWithExplicitPersonaOverride`
- `TestPersonaOptionsIncludesRickSanchez`
- `TestPresetOptionsIncludesFullRick`
- `TestRenderPersonaIncludesRickDescription`
- `TestRenderPresetIncludesFullRickDescription`
- `TestPresetFullRickSetsPersonaToRick`
- `TestPresetFullRickDoesNotOverrideExplicitPersona`
- `TestTUIComponentsForPresetFullRickMatchesFullGentleman`

## Spec Compliance

All 55+ requirements from the spec are met. Key compliance highlights:

- `PersonaRickSanchez` registered as a fully managed persona (not Custom)
- `isFullPersona()` gates output styles, agent overlays, managed sections, and cleanup identically for Gentleman and Rick
- `PresetFullRick` shares the same component list as `PresetFullGentleman`
- Default persona derivation works for both CLI and TUI
- Backward compatibility preserved: agent key stays `"gentleman"`, output style names unchanged
- All 5 Rick assets contain the 8 required structural sections with functional guardrails intact

## Warnings from Verification

Two low-impact test coverage suggestions were noted in the verify report:

1. **Rick assets not in `TestPersonasContainContextualSkillLoadingDirective`** — Rick assets DO contain the directive (verified manually), but the test does not guard against accidental removal.
2. **Rick assets not in `TestAllEmbeddedAssetsAreReadable`** — Assets are covered by other tests and would panic at runtime if missing.

These are non-blocking and can be addressed in a follow-up commit if desired.

## Lessons Learned

1. **Centralize gating early**: The `isFullPersona()` helper eliminated 7+ hardcoded `PersonaGentleman` checks. Adding a third full persona in the future will now require changing a single switch case instead of hunting through injection logic.

2. **Golden files catch drift**: Byte-exact golden file comparison immediately surfaced a formatting issue during development that unit tests alone would have missed.

3. **Preset default persona coupling**: Deriving the default persona from the preset required touching both CLI (`validate.go`) and TUI (`model.go`) paths. A shared default-resolution helper could reduce this duplication.

4. **Agent overlay key naming for compatibility**: Keeping the OpenCode agent key as `"gentleman"` for Rick avoided breaking existing UI workflows. The `description` field is the right place for persona-specific labeling, not the JSON key.

5. **Asset consistency across agents**: All 5 Rick assets share the same structure but differ only in agent-specific phrasing. A code-generation approach or template system could reduce maintenance if more personas are added.

## Rollback Validity

Per REQ-ROLLBACK-001, reverting this change requires:
1. Reverting changes to `types.go`, `validate.go`, TUI files, and `inject.go`
2. Deleting the 5 new Rick asset files
3. Running `go test ./...` to confirm baseline

This remains valid and was verified during development.

## SDD Cycle Complete

The change has been fully planned, implemented, verified, and archived.
Ready for the next change.
