# Spec: FullRick Preset and Rick Sanchez Persona

## Overview

This spec defines the requirements for adding `PersonaRickSanchez` and `PresetFullRick` to framework-ai. The Rick Sanchez persona is a cynical-genius senior architect that teaches architecture, patterns, and testing in Rick Sanchez style — neutral Spanish, scientific metaphors, multiverse references, informal language — while remaining fully functional. The preset delivers the complete Gentleman AI ecosystem with Rick as the default persona.

## References

- Proposal: `openspec/changes/fullrick-preset/proposal.md`
- Persona injection system: `internal/components/persona/inject.go`
- Reference persona (Gentleman): `internal/assets/generic/persona-gentleman.md`
- Reference neutral persona: `internal/assets/generic/persona-neutral.md`
- TUI persona picker: `internal/tui/screens/persona.go`
- TUI preset picker: `internal/tui/screens/preset.go`
- Types: `internal/model/types.go`

---

## 1. Persona Registration

### 1.1 Type Constants

**REQ-PERSONA-001:** The `PersonaID` type MUST support a new constant `PersonaRickSanchez` with value `"rick-sanchez"`.

**REQ-PERSONA-002:** The `PersonaRickSanchez` constant MUST be added to `internal/model/types.go` alongside existing `PersonaGentleman`, `PersonaNeutral`, and `PersonaCustom`.

### 1.2 Persona Catalog

**REQ-PERSONA-003:** The persona system MUST recognize `rick-sanchez` as a valid, managed persona ID in all validation and dispatch paths.

**REQ-PERSONA-004:** `rick-sanchez` MUST NOT be treated as `PersonaCustom` — it receives full injection, output styles, and agent overlays.

---

## 2. Persona Content Dispatch

### 2.1 Asset File Coverage

**REQ-DISPATCH-001:** The persona injection system MUST provide Rick Sanchez persona assets for all 14 supported agents. The dispatch MUST follow the same per-agent fallback strategy as Gentleman:

| Agent | Asset Path |
|-------|-----------|
| Claude Code | `claude/persona-rick-sanchez.md` |
| OpenCode | `opencode/persona-rick-sanchez.md` |
| Kilocode | `opencode/persona-rick-sanchez.md` (same as OpenCode) |
| Kimi | `kimi/persona-rick-sanchez.md` |
| Kiro IDE | `kiro/persona-rick-sanchez.md` |
| All others (Gemini, Cursor, VS Code, Codex, Windsurf, Antigravity, Qwen, OpenClaw, Pi) | `generic/persona-rick-sanchez.md` |

**REQ-DISPATCH-002:** Each asset file MUST be embedded via `//go:embed` and accessible through `assets.MustRead()`.

**REQ-DISPATCH-003:** The `personaContent()` function in `internal/components/persona/inject.go` MUST be extended to dispatch Rick assets per agent, using the same switch structure as Gentleman.

### 2.2 Asset Content Requirements

**REQ-ASSET-001:** All Rick Sanchez persona assets MUST include the same structural sections as Gentleman: `## Rules`, `## Personality`, `## Persona Scope`, `## Language`, `## Tone`, `## Philosophy`, `## Expertise`, `## Behavior`, and `## Contextual Skill Loading (MANDATORY)`.

**REQ-ASSET-002:** The Rick persona MUST maintain functional teaching capability — it MUST still guide users on architecture, patterns, testing, and SDD workflows.

**REQ-ASSET-003:** The Rick persona's language rules MUST specify neutral Spanish (no regionalisms, no voseo) when responding in Spanish, matching the user's input language.

**REQ-ASSET-004:** The Rick persona MUST NOT inject Rioplatense slang, voseo, or stylistic emphasis (CAPS, exclamations, rhetorical questions) into generated code, UI strings, or task artifacts. The persona styles chat replies only.

**REQ-ASSET-005:** Rick persona assets MUST use scientific metaphors, multiverse references, and informal language in the personality and tone sections while preserving all functional guardrails from Gentleman (e.g., skill loading self-check, response-length contract, one-question-at-a-time rule).

---

## 3. Output Style and Agent Overlay Gating

### 3.1 Full-Persona Helper

**REQ-GATING-001:** An `isFullPersona(persona model.PersonaID) bool` helper MUST be introduced in `internal/components/persona/inject.go`.

**REQ-GATING-002:** `isFullPersona()` MUST return `true` for `PersonaGentleman` and `PersonaRickSanchez`, and `false` for all other personas (including `PersonaNeutral` and `PersonaCustom`).

**REQ-GATING-003:** All hardcoded `persona == model.PersonaGentleman` checks in `inject.go` that gate output-style writing, agent overlay creation, or managed-section preservation logic MUST be refactored to use `isFullPersona()`.

**REQ-GATING-004:** Specifically, the following behaviors gated by `isFullPersona()` MUST apply identically to Rick and Gentleman:
- Writing the output-style file (`output-styles/gentleman.md` for Claude, `output-style.md` module for Kimi).
- Merging `"outputStyle": "Gentleman"` into agent settings JSON.
- Creating the OpenCode/Kilocode `agent.gentleman` tab-switchable agent definition.
- Skipping managed-section preservation in `preserveManagedSections()` (full personas replace the entire file safely).

**REQ-GATING-005:** The cleanup logic for switching AWAY from a full persona MUST also use `isFullPersona()`: when switching from Rick to Neutral, the system MUST remove output-style artifacts and agent overlays just as it does when switching from Gentleman to Neutral.

### 3.2 Output Style Compatibility

**REQ-STYLE-001:** Rick persona MUST reuse the existing `output-style-gentleman.md` content and settings overlay. The output style is a formatting layer independent of persona identity.

**REQ-STYLE-002:** The output style file path and settings key MUST remain `gentleman.md` and `"Gentleman"` respectively for backward compatibility with Claude Code's output style system.

---

## 4. Preset Registration

### 4.1 Preset Constant

**REQ-PRESET-001:** The `PresetID` type MUST support a new constant `PresetFullRick` with value `"full-rick"`.

**REQ-PRESET-002:** `PresetFullRick` MUST be added to `internal/model/types.go` alongside existing presets.

### 4.2 Preset Components

**REQ-PRESET-003:** `PresetFullRick` MUST have the exact same component list as `PresetFullGentleman`:
- `ComponentEngram`
- `ComponentSDD`
- `ComponentSkills`
- `ComponentContext7`
- `ComponentPersona`
- `ComponentPermission`
- `ComponentGGA`
- `ComponentClaudeTheme`
- `ComponentOpenCodeGentleLogo`
- `ComponentCodeGraph`

### 4.3 Preset Default Persona

**REQ-PRESET-004:** When `PresetFullRick` is selected, the default persona MUST be `PersonaRickSanchez` instead of `PersonaGentleman`.

**REQ-PRESET-005:** When the user explicitly overrides the persona via CLI (`--persona`) or TUI selection, the override MUST take precedence over the preset default.

---

## 5. TUI Integration

### 5.1 Persona Picker

**REQ-TUI-PERSONA-001:** `screens.PersonaOptions()` MUST include `model.PersonaRickSanchez` as a selectable option.

**REQ-TUI-PERSONA-002:** `screens.personaDescriptions` MUST include a description for Rick: `"Managed Rick Sanchez persona — cynical genius who teaches architecture with scientific flair"`.

**REQ-TUI-PERSONA-003:** The persona picker screen MUST render the Rick option with the same radio-button styling as existing personas.

### 5.2 Preset Picker

**REQ-TUI-PRESET-001:** `screens.PresetOptions()` MUST include `model.PresetFullRick` as a selectable option.

**REQ-TUI-PRESET-002:** `screens.presetDescriptions` MUST include a description for FullRick: `"Everything: memory, SDD, skills, docs, persona & security — with Rick Sanchez as default persona"`.

**REQ-TUI-PRESET-003:** The preset picker MUST render the FullRick option with the same radio-button styling as existing presets.

### 5.3 Model Default Selection

**REQ-TUI-MODEL-001:** In `internal/tui/model.go`, the `NewModel()` function's default `Selection.Persona` MUST remain `PersonaGentleman` (the global default when no preset is chosen yet).

**REQ-TUI-MODEL-002:** When the user selects `PresetFullRick` on the preset screen, `m.Selection.Persona` MUST automatically update to `PersonaRickSanchez` unless the user had already explicitly chosen a different persona.

**REQ-TUI-MODEL-003:** The `componentsForPreset()` helper in `internal/tui/model.go` MUST handle `PresetFullRick` identically to `PresetFullGentleman`.

---

## 6. CLI Validation

### 6.1 Persona Validation

**REQ-CLI-PERSONA-001:** `normalizePersona()` in `internal/cli/validate.go` MUST accept `"rick-sanchez"` as a valid persona flag value.

**REQ-CLI-PERSONA-002:** An invalid persona value MUST continue to return the error: `unsupported persona %q`.

### 6.2 Preset Validation

**REQ-CLI-PRESET-001:** `normalizePreset()` in `internal/cli/validate.go` MUST accept `"full-rick"` as a valid preset flag value.

**REQ-CLI-PRESET-002:** `componentsForPreset()` in `internal/cli/validate.go` MUST handle `PresetFullRick` identically to `PresetFullGentleman`.

**REQ-CLI-PRESET-003:** An invalid preset value MUST continue to return the error: `unsupported preset %q`.

### 6.3 Default Flag Behavior

**REQ-CLI-DEFAULT-001:** When `--preset full-rick` is provided without `--persona`, the normalized selection MUST default to `PersonaRickSanchez`.

**REQ-CLI-DEFAULT-002:** When `--preset full-gentleman` is provided without `--persona`, the normalized selection MUST continue to default to `PersonaGentleman`.

---

## 7. Catalog Update

**REQ-CATALOG-001:** `internal/catalog/components.go` MUST NOT require changes for the persona/preset itself, but the persona component description MAY be updated to mention Rick: `"Gentleman, Rick Sanchez, neutral or custom behavior"`.

---

## 8. Tests and Golden Files

### 8.1 Unit Tests for Injection

**REQ-TEST-INJECT-001:** A test MUST verify that `Inject(claudeAdapter(), PersonaRickSanchez)` writes the persona marker section with real Rick content to `CLAUDE.md`.

**REQ-TEST-INJECT-002:** A test MUST verify that `Inject(kimiAdapter(), PersonaRickSanchez)` writes the persona module and output-style module correctly.

**REQ-TEST-INJECT-003:** A test MUST verify that `Inject(opencodeAdapter(), PersonaRickSanchez)` creates the `agent.gentleman` overlay in `opencode.json` (the agent key name remains `gentleman` for compatibility).

**REQ-TEST-INJECT-004:** A test MUST verify that switching from Rick to Neutral cleans output-style artifacts and agent overlays.

**REQ-TEST-INJECT-005:** A test MUST verify idempotency: running `Inject()` twice with Rick on the same adapter MUST report `Changed=false` on the second run.

### 8.2 Unit Tests for isFullPersona

**REQ-TEST-FULL-001:** A test MUST verify that `isFullPersona(PersonaGentleman)` returns `true`.

**REQ-TEST-FULL-002:** A test MUST verify that `isFullPersona(PersonaRickSanchez)` returns `true`.

**REQ-TEST-FULL-003:** A test MUST verify that `isFullPersona(PersonaNeutral)` returns `false`.

**REQ-TEST-FULL-004:** A test MUST verify that `isFullPersona(PersonaCustom)` returns `false`.

### 8.3 Golden Files

**REQ-TEST-GOLDEN-001:** Golden files MUST be created under `internal/components/persona/testdata/` containing the expected injected persona content for at least 3 representative agents: Claude, OpenCode, and Kimi.

**REQ-TEST-GOLDEN-002:** Each golden file MUST be named using the pattern: `<agent>-rick-sanchez.golden.md`.

**REQ-TEST-GOLDEN-003:** Golden file tests MUST compare the exact bytes written by `Inject()` against the golden file, failing on any drift.

**REQ-TEST-GOLDEN-004:** Golden files MUST be regenerated via `go test ./... -update` or a documented `go generate` command when asset content changes.

### 8.4 TUI Tests

**REQ-TEST-TUI-001:** Tests for `screens.PersonaOptions()` MUST assert that Rick is present and in the expected position.

**REQ-TEST-TUI-002:** Tests for `screens.PresetOptions()` MUST assert that FullRick is present and in the expected position.

**REQ-TEST-TUI-003:** Tests for persona picker rendering MUST assert that the Rick description appears in the rendered output.

### 8.5 CLI Validation Tests

**REQ-TEST-CLI-001:** Tests for `normalizePersona()` MUST accept `"rick-sanchez"`.

**REQ-TEST-CLI-002:** Tests for `normalizePreset()` MUST accept `"full-rick"`.

**REQ-TEST-CLI-003:** Tests for `componentsForPreset(PresetFullRick)` MUST return the same component list as `PresetFullGentleman`.

---

## 9. Scenarios

### Scenario 1: Fresh Install with FullRick Preset

```gherkin
Given the user runs `framework-ai install --preset full-rick`
When the installation completes
Then the default persona is Rick Sanchez
And all FullGentleman components are installed
And the output style overlay is applied to supported agents
And the agent overlay is applied to OpenCode/Kilocode
```

### Scenario 2: TUI Selection of Rick Persona

```gherkin
Given the user is on the persona picker screen
When the user navigates to and selects "rick-sanchez"
Then the selection persona is set to PersonaRickSanchez
And the next screen is the preset picker
```

### Scenario 3: TUI Selection of FullRick Preset

```gherkin
Given the user is on the preset picker screen
When the user navigates to and selects "full-rick"
Then the selection preset is set to PresetFullRick
And the selection persona is updated to PersonaRickSanchez
And the component list matches PresetFullGentleman
```

### Scenario 4: Switching from Rick to Neutral

```gherkin
Given an existing installation with Rick persona
When the user runs `framework-ai install --persona neutral`
Then the persona content is replaced with neutral
And output-style files are removed
And agent.gentleman is removed from OpenCode/Kilocode settings
And managed sections (SDD, engram) are preserved
```

### Scenario 5: Rick Persona Injection Idempotency

```gherkin
Given an existing installation with Rick persona
When the user runs `framework-ai install --preset full-rick` again
Then no files are modified
And Inject reports Changed=false
```

### Scenario 6: Explicit Persona Override with FullRick Preset

```gherkin
Given the user runs `framework-ai install --preset full-rick --persona gentleman`
When the installation completes
Then the preset is PresetFullRick
And the persona is PersonaGentleman
And all components are installed
```

---

## 10. Backward Compatibility

**REQ-BWC-001:** Existing installations using `PersonaGentleman` or `PresetFullGentleman` MUST continue to work identically after this change.

**REQ-BWC-002:** The `agent.gentleman` key in OpenCode/Kilocode settings MUST remain named `gentleman` even when Rick persona is active, to avoid breaking existing agent-switching workflows.

**REQ-BWC-003:** Output style files and settings keys MUST retain their existing names (`gentleman.md`, `"Gentleman"`) for compatibility with Claude Code's output style system.

---

## 11. Rollback

**REQ-ROLLBACK-001:** Reverting this change MUST require only:
1. Reverting changes to `types.go`, `validate.go`, TUI files, `inject.go`
2. Deleting the 5 new Rick asset files
3. Running `go test ./...` to confirm baseline

---

## 12. Acceptance Criteria

- [x] `go test ./...` passes
- [x] TUI shows Rick persona and FullRick preset options
- [x] `framework-ai install --preset full-rick --persona rick-sanchez` injects correct assets
- [x] Rick persona receives output styles and agent overlays like Gentleman
- [x] Switching from Rick to Neutral cleans all Rick-specific artifacts
- [x] Golden files exist for at least Claude, OpenCode, and Kimi
- [x] `isFullPersona()` helper has unit test coverage for all 4 persona values
