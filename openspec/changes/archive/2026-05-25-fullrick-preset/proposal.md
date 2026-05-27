# Proposal: FullRick Preset and Rick Sanchez Persona

## Intent

Add `PersonaRickSanchez` and `PresetFullRick` to framework-ai. The persona is a cynical-genius senior architect who teaches architecture, patterns, and testing in Rick Sanchez style — neutral Spanish, scientific metaphors, multiverse references, informal language — while remaining fully functional. The preset delivers the complete Gentleman AI ecosystem with Rick as the default persona.

## Scope

### In Scope
- `PersonaRickSanchez` ID and persona asset files for all 13 supported agents (generic + per-agent variants: claude, opencode, kilocode, gemini, cursor, vscode, codex, windsurf, antigravity, kimi, qwen, kiro, openclaw, pi)
- `PresetFullRick` with identical components to `PresetFullGentleman`
- Update `personaContent()` dispatch, CLI validation, TUI pickers, and preset handling
- Refactor output-style and agent-overlay logic so Rick receives the same features as Gentleman
- Tests and golden files for all agents

### Out of Scope
- New components or skills
- Changes to existing Gentleman/Neutral personas
- LLM-specific prompt engineering beyond asset files

## Capabilities

### New Capabilities
- `rick-sanchez-persona`: Teaching persona with agent-specific variants
- `fullrick-preset`: Full ecosystem preset defaulting to Rick persona

### Modified Capabilities
- `persona-injection`: Extend dispatch and feature-gating to support Rick as a full persona

## Approach

Add `PersonaRickSanchez` to `types.go`. Create 5 embedded asset markdowns following the existing Gentleman structure, adapted to Rick's tone and constraints (neutral Spanish, no regionalisms, functional architecture guidance).

In `inject.go`, refactor hardcoded `PersonaGentleman` checks for output styles and OpenCode agent overlay to treat both Gentleman and Rick as "full-featured" personas (via an `isFullPersona()` helper). Update `personaContent()` to dispatch Rick assets per agent.

Add `PresetFullRick` with the same component list as `PresetFullGentleman`. Update validation, TUI persona/preset screens, and `model.go` default selection.

## Affected Areas

| Area | Impact | Description |
|------|--------|-------------|
| `internal/model/types.go` | Modified | Add `PersonaRickSanchez` and `PresetFullRick` constants |
| `internal/assets/*/persona-rick-sanchez.md` | New | 5 persona asset files |
| `internal/components/persona/inject.go` | Modified | Dispatch + output-style/agent-overlay generalization |
| `internal/catalog/components.go` | Modified | Add preset description |
| `internal/cli/validate.go` | Modified | Add preset/persona validation |
| `internal/tui/screens/persona.go` | Modified | Add Rick option |
| `internal/tui/screens/preset.go` | Modified | Add FullRick preset |
| `internal/tui/model.go` | Modified | Handle preset default persona selection |

## Risks

| Risk | Likelihood | Mitigation |
|------|------------|----------|
| Output-style logic misses Rick | Low | Refactor to `isFullPersona()` helper, add unit test |
| Asset tone too parody, not functional | Med | Review assets against "still teaches architecture" constraint |

## Rollback Plan

1. Revert `types.go`, `validate.go`, TUI files, and `inject.go`
2. Delete the 5 new asset files
3. Run `go test ./...` to confirm baseline

## Dependencies

None.

## Success Criteria

- [x] `go test ./...` passes
- [x] TUI shows Rick persona and FullRick preset options
- [x] `framework-ai install --preset full-rick --persona rick-sanchez` injects correct assets
- [x] Rick persona receives output styles and agent overlays like Gentleman
