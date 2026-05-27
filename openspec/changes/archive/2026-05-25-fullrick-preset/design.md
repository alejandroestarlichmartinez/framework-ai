# Design: FullRick Preset and Rick Sanchez Persona

## 1. Asset File Structure

### 1.1 New Files

Five embedded asset markdowns are created, following the existing Gentleman directory structure:

```
internal/assets/
├── claude/persona-rick-sanchez.md      # Claude Code CLAUDE.md content
├── opencode/persona-rick-sanchez.md    # OpenCode AGENTS.md + Kilocode shared
├── kimi/persona-rick-sanchez.md        # Kimi KIMI.md persona module
├── kiro/persona-rick-sanchez.md        # Kiro steering file content
└── generic/persona-rick-sanchez.md     # Fallback for all other agents
```

**No new directories are required.** The existing `//go:embed` directive in `internal/assets/assets.go` already covers `all:claude all:opencode all:kimi all:kiro all:generic`.

### 1.2 Asset Content Strategy

Each Rick asset mirrors the Gentleman asset structure section-for-section:

| Section | Gentleman Content | Rick Content Adaptation |
|---------|-------------------|------------------------|
| `## Rules` | Response-length contract, one-question rule, verification rule | Same functional rules, adapted tone |
| `## Personality` | Senior Architect, GDE & MVP, passionate teacher | Cynical genius scientist, multiverse traveler, tired of bad code but cares deep down |
| `## Persona Scope` | Chat reply only; artifacts in English | Identical scope guardrails |
| `## Language` | Match user language; Rioplatense Spanish | Match user language; **neutral Spanish** (no voseo, no regionalisms) |
| `## Tone` | Passionate and direct, from a place of CARING | Sardonic, impatient with stupidity, but invested in making the user smarter |
| `## Philosophy` | CONCEPTS > CODE, AI IS A TOOL, etc. | Same principles, expressed with scientific metaphor |
| `## Expertise` | Clean/Hexagonal/Screaming Architecture, testing, etc. | Identical expertise list |
| `## Behavior` | Push back without context, use analogies, correct errors | Same behaviors, Rick-style delivery |
| `## Contextual Skill Loading` | Mandatory self-check before every response | Identical requirement |

**Key constraint:** The functional rules (skill loading, response length, one question, verification) MUST remain verbatim or equivalent in meaning. Only the personality, tone, and language flavoring change.

### 1.3 Output Style Reuse

The output-style file (`claude/output-style-gentleman.md` and `kimi/output-style-gentleman.md`) is **persona-agnostic**. It defines formatting rules, not personality. Rick reuses the same output style without modification.

---

## 2. isFullPersona() Helper Design

### 2.1 Rationale

Currently `inject.go` has 7+ hardcoded `persona == model.PersonaGentleman` checks that gate:
- Output style file writing
- Settings JSON merge for outputStyle
- OpenCode/Kilocode agent overlay creation
- `preserveManagedSections` bypass
- Cleanup of residuals when switching away

Adding Rick requires touching all of these. A centralized helper eliminates drift risk.

### 2.2 Implementation

```go
// isFullPersona returns true for personas that receive the complete feature
// set: output styles, agent overlays, and full persona asset injection.
// Gentleman and Rick are "full" personas; Neutral and Custom are not.
func isFullPersona(persona model.PersonaID) bool {
    switch persona {
    case model.PersonaGentleman, model.PersonaRickSanchez:
        return true
    default:
        return false
    }
}
```

### 2.3 Replacement Map

| Line (approx) | Current Check | New Check |
|--------------|---------------|-----------|
| 302 | `persona == model.PersonaGentleman` (Kimi output-style module) | `isFullPersona(persona)` |
| 319 | `persona != model.PersonaCustom` (agent overlay guard) | Keep as-is (Custom is already excluded) |
| 322 | `persona == model.PersonaGentleman` (gentleman agent key) | `isFullPersona(persona)` |
| 330-341 | `persona == model.PersonaGentleman` else branch (remove agent key) | `isFullPersona(persona)` else branch |
| 346 | `persona == model.PersonaGentleman` (output style write) | `isFullPersona(persona)` |
| 374 | `persona != model.PersonaGentleman` (cleanup) | `!isFullPersona(persona)` |
| 523 | `persona == model.PersonaGentleman` (preserveManagedSections) | `isFullPersona(persona)` |

### 2.4 Agent Overlay Key Naming

The OpenCode/Kilocode agent overlay JSON uses the key `"gentleman"` inside the `"agent"` object:

```json
{
  "agent": {
    "gentleman": {
      "mode": "primary",
      "description": "...",
      "prompt": "{file:./AGENTS.md}"
    }
  }
}
```

**Decision:** Keep the key name `"gentleman"` for both Gentleman and Rick personas. Changing it to `"rick"` would break:
- Existing OpenCode UI agent switching workflows
- Documentation references
- User muscle memory

The `description` field inside the object MAY be updated dynamically based on persona, but the key name stays `gentleman` for backward compatibility.

---

## 3. personaContent() Dispatch Table

### 3.1 Current Structure

```go
func personaContent(agent model.AgentID, persona model.PersonaID) string {
    switch persona {
    case model.PersonaNeutral:
        return assets.MustRead("generic/persona-neutral.md")
    case model.PersonaCustom:
        return ""
    default:
        // Gentleman
        switch agent {
        case model.AgentClaudeCode:
            return assets.MustRead("claude/persona-gentleman.md")
        case model.AgentOpenCode, model.AgentKilocode:
            return assets.MustRead("opencode/persona-gentleman.md")
        case model.AgentKimi:
            return assets.MustRead("kimi/persona-gentleman.md")
        case model.AgentKiroIDE:
            return assets.MustRead("kiro/persona-gentleman.md")
        default:
            return assets.MustRead("generic/persona-gentleman.md")
        }
    }
}
```

### 3.2 New Structure

```go
func personaContent(agent model.AgentID, persona model.PersonaID) string {
    switch persona {
    case model.PersonaNeutral:
        return assets.MustRead("generic/persona-neutral.md")
    case model.PersonaCustom:
        return ""
    case model.PersonaRickSanchez:
        switch agent {
        case model.AgentClaudeCode:
            return assets.MustRead("claude/persona-rick-sanchez.md")
        case model.AgentOpenCode, model.AgentKilocode:
            return assets.MustRead("opencode/persona-rick-sanchez.md")
        case model.AgentKimi:
            return assets.MustRead("kimi/persona-rick-sanchez.md")
        case model.AgentKiroIDE:
            return assets.MustRead("kiro/persona-rick-sanchez.md")
        default:
            return assets.MustRead("generic/persona-rick-sanchez.md")
        }
    default:
        // Gentleman
        switch agent {
        case model.AgentClaudeCode:
            return assets.MustRead("claude/persona-gentleman.md")
        case model.AgentOpenCode, model.AgentKilocode:
            return assets.MustRead("opencode/persona-gentleman.md")
        case model.AgentKimi:
            return assets.MustRead("kimi/persona-gentleman.md")
        case model.AgentKiroIDE:
            return assets.MustRead("kiro/persona-gentleman.md")
        default:
            return assets.MustRead("generic/persona-gentleman.md")
        }
    }
}
```

### 3.3 Alternative: Table-Driven Dispatch (Rejected)

A map-based dispatch was considered:

```go
var personaAssets = map[model.PersonaID]map[model.AgentID]string{...}
```

**Rejected** because:
- The existing switch is explicit and compile-time safe
- A map would require initialization logic and nil checks
- The agent fallback (`default` case) is cleaner in switch form

---

## 4. Preset Component List

### 4.1 Preset Definition

`PresetFullRick` is added to `internal/model/types.go`:

```go
const (
    PresetFullGentleman PresetID = "full-gentleman"
    PresetFullRick      PresetID = "full-rick"
    PresetEcosystemOnly PresetID = "ecosystem-only"
    PresetMinimal       PresetID = "minimal"
    PresetCustom        PresetID = "custom"
)
```

### 4.2 Component Resolution

Both `internal/cli/validate.go` and `internal/tui/model.go` have `componentsForPreset()` helpers. Both are updated:

```go
func componentsForPreset(preset model.PresetID) []model.ComponentID {
    switch preset {
    case model.PresetMinimal:
        return []model.ComponentID{model.ComponentEngram}
    case model.PresetEcosystemOnly:
        return []model.ComponentID{model.ComponentEngram, model.ComponentSDD, model.ComponentSkills, model.ComponentContext7, model.ComponentGGA}
    case model.PresetCustom:
        return nil
    case model.PresetFullRick:
        fallthrough
    default:
        return []model.ComponentID{
            model.ComponentEngram,
            model.ComponentSDD,
            model.ComponentSkills,
            model.ComponentContext7,
            model.ComponentPersona,
            model.ComponentPermission,
            model.ComponentGGA,
            model.ComponentClaudeTheme,
            model.ComponentOpenCodeGentleLogo,
            model.ComponentCodeGraph,
        }
    }
}
```

**Note:** `PresetFullGentleman` and `PresetFullRick` share the same component list. Using `fallthrough` or collapsing them into the `default` case keeps the code DRY.

### 4.3 Default Persona Resolution

In `internal/cli/validate.go`, `NormalizeInstallFlags` is updated to derive the default persona from the preset when no explicit `--persona` flag is provided:

```go
func normalizePersona(value string, preset model.PresetID) (model.PersonaID, error) {
    if strings.TrimSpace(value) == "" {
        if preset == model.PresetFullRick {
            return model.PersonaRickSanchez, nil
        }
        return model.PersonaGentleman, nil
    }
    // ... existing switch
}
```

**TUI equivalent:** In `internal/tui/model.go`, when the preset screen confirms `PresetFullRick`, `m.Selection.Persona` is set to `PersonaRickSanchez` unless already explicitly overridden.

---

## 5. TUI Wiring

### 5.1 Persona Screen (`internal/tui/screens/persona.go`)

```go
func PersonaOptions() []model.PersonaID {
    return []model.PersonaID{
        model.PersonaGentleman,
        model.PersonaRickSanchez,
        model.PersonaNeutral,
        model.PersonaCustom,
    }
}

var personaDescriptions = map[model.PersonaID]string{
    model.PersonaGentleman:    "Managed Gentleman persona with teaching-first guidance",
    model.PersonaRickSanchez:  "Managed Rick Sanchez persona — cynical genius who teaches architecture with scientific flair",
    model.PersonaNeutral:      "Managed neutral persona with the same guidance and less regional tone",
    model.PersonaCustom:       "Keep your existing persona unmanaged; framework-ai does not inject a persona",
}
```

### 5.2 Preset Screen (`internal/tui/screens/preset.go`)

```go
func PresetOptions() []model.PresetID {
    return []model.PresetID{
        model.PresetFullGentleman,
        model.PresetFullRick,
        model.PresetEcosystemOnly,
        model.PresetMinimal,
        model.PresetCustom,
    }
}

var presetDescriptions = map[model.PresetID]string{
    model.PresetFullGentleman: "Everything: memory, SDD, skills, docs, persona & security",
    model.PresetFullRick:      "Everything: memory, SDD, skills, docs, persona & security — with Rick Sanchez as default persona",
    model.PresetEcosystemOnly: "Core tools only: memory, SDD, skills & docs (no persona/security)",
    model.PresetMinimal:       "Just Engram persistent memory",
    model.PresetCustom:        "Choose components and skills manually; keep existing persona/settings unmanaged",
}
```

### 5.3 Model State Transitions

In `internal/tui/model.go`, the preset confirmation handler (`ScreenPreset` case in `confirmSelection`) is updated:

```go
case ScreenPreset:
    options := screens.PresetOptions()
    if m.Cursor < len(options) {
        m.Selection.Preset = options[m.Cursor]
        m.Selection.Components = componentsForPreset(options[m.Cursor])
        
        // Auto-set default persona for FullRick
        if m.Selection.Preset == model.PresetFullRick && m.Selection.Persona == model.PersonaGentleman {
            m.Selection.Persona = model.PersonaRickSanchez
        }
        // ... rest of flow
    }
```

**Backward compatibility note:** We only auto-switch to Rick if the current persona is the global default (`PersonaGentleman`). If the user explicitly selected Neutral or Custom on the persona screen before choosing the preset, their choice is preserved.

---

## 6. Test Strategy

### 6.1 Golden Files Pattern

Golden files store the exact expected output of persona injection for deterministic verification. They protect against accidental drift in asset content.

#### Directory Structure

```
internal/components/persona/testdata/
├── claude-rick-sanchez.golden.md
├── opencode-rick-sanchez.golden.md
├── kimi-rick-sanchez.golden.md
└── README.md
```

#### Generation

Golden files are generated by running the real `Inject()` function against a temp directory and copying the result:

```go
//go:generate go test ./... -run TestGenerateGoldenFiles -update
```

Or via a dedicated test that writes when `-update` is passed:

```go
var updateGolden = flag.Bool("update", false, "update golden files")

func TestRickGoldenFiles(t *testing.T) {
    cases := []struct {
        name    string
        agent   string
        path    string
    }{
        {"claude", "claude-code", "CLAUDE.md"},
        {"opencode", "opencode", "AGENTS.md"},
        {"kimi", "kimi", ".kimi/persona.md"},
    }
    
    for _, tc := range cases {
        t.Run(tc.name, func(t *testing.T) {
            home := t.TempDir()
            adapter, _ := agents.NewAdapter(tc.agent)
            Inject(home, adapter, model.PersonaRickSanchez)
            
            got, _ := os.ReadFile(filepath.Join(home, tc.path))
            goldenPath := filepath.Join("testdata", tc.name+"-rick-sanchez.golden.md")
            
            if *updateGolden {
                os.WriteFile(goldenPath, got, 0o644)
                return
            }
            
            want, _ := os.ReadFile(goldenPath)
            if !bytes.Equal(got, want) {
                t.Fatalf("golden mismatch. run with -update to regenerate.\n--- want ---\n%s\n--- got ---\n%s", want, got)
            }
        })
    }
}
```

### 6.2 Unit Test Coverage Matrix

| Function / Behavior | Test Name | Personas Covered |
|---------------------|-----------|------------------|
| `isFullPersona()` | `TestIsFullPersona` | Gentleman, Rick, Neutral, Custom |
| `personaContent()` dispatch | `TestPersonaContentRickSanchez` | Rick × 14 agents |
| `Inject()` marker injection | `TestInjectClaudeRickSanchez` | Rick |
| `Inject()` output style | `TestInjectClaudeRickSanchezWritesOutputStyle` | Rick |
| `Inject()` agent overlay | `TestInjectOpenCodeRickSanchezCreatesAgent` | Rick |
| `Inject()` idempotency | `TestInjectClaudeRickSanchezIdempotent` | Rick |
| `Inject()` switch cleanup | `TestInjectClaudeRickToNeutralCleansOutputStyle` | Rick → Neutral |
| `normalizePersona()` | `TestNormalizePersonaAcceptsRickSanchez` | Rick |
| `normalizePreset()` | `TestNormalizePresetAcceptsFullRick` | FullRick |
| `componentsForPreset()` | `TestComponentsForPresetFullRick` | FullRick |
| TUI persona options | `TestPersonaOptionsIncludesRickSanchez` | Rick |
| TUI preset options | `TestPresetOptionsIncludesFullRick` | FullRick |

### 6.3 Test Data Helpers

The existing `assertGentlemanLanguageGuardrails` helper in `inject_test.go` is generalized or duplicated for Rick:

```go
func assertRickLanguageGuardrails(t *testing.T, text string) {
    t.Helper()
    required := []string{
        "Match the user's current language in your REPLY ONLY",
        "Do not switch languages unless the user does",
    }
    banned := []string{
        "Rioplatense",   // Rick uses neutral Spanish
        "voseo",
        `Say "déjame verificar"`,
    }
    // ... same logic as gentleman variant
}
```

### 6.4 Integration Test

An integration-level test runs `NormalizeInstallFlags` with `--preset full-rick` and asserts the resulting `Selection`:

```go
func TestNormalizeInstallFlags_FullRickDefaultsToRickSanchez(t *testing.T) {
    input := InstallFlags{Preset: "full-rick"}
    result, err := NormalizeInstallFlags(input, system.DetectionResult{})
    if err != nil {
        t.Fatalf("unexpected error: %v", err)
    }
    if result.Selection.Persona != model.PersonaRickSanchez {
        t.Fatalf("want persona %q, got %q", model.PersonaRickSanchez, result.Selection.Persona)
    }
    if result.Selection.Preset != model.PresetFullRick {
        t.Fatalf("want preset %q, got %q", model.PresetFullRick, result.Selection.Preset)
    }
}
```

---

## 7. Risks and Mitigations

| Risk | Mitigation |
|------|-----------|
| `isFullPersona()` misses a hardcoded check | Audit all `PersonaGentleman` references in `inject.go` via grep; add exhaustive unit test |
| Asset tone becomes parody, loses teaching value | Review each asset against "still teaches architecture" constraint; keep expertise/behavior sections functional |
| Golden file drift not caught in CI | Golden tests fail on any byte diff; CI runs them on every PR |
| TUI cursor positions shift when adding Rick | Update TUI tests that assert cursor indices; use named constants |

---

## 8. Affected Files Summary

| File | Change |
|------|--------|
| `internal/model/types.go` | Add `PersonaRickSanchez`, `PresetFullRick` |
| `internal/assets/claude/persona-rick-sanchez.md` | New |
| `internal/assets/opencode/persona-rick-sanchez.md` | New |
| `internal/assets/kimi/persona-rick-sanchez.md` | New |
| `internal/assets/kiro/persona-rick-sanchez.md` | New |
| `internal/assets/generic/persona-rick-sanchez.md` | New |
| `internal/components/persona/inject.go` | Add `isFullPersona()`, update `personaContent()`, refactor gating checks |
| `internal/components/persona/inject_test.go` | Add Rick tests, golden file tests, `isFullPersona` tests |
| `internal/components/persona/testdata/*` | New golden files |
| `internal/cli/validate.go` | Accept Rick persona/preset, default persona from preset |
| `internal/cli/validate_test.go` | Add Rick validation tests |
| `internal/tui/screens/persona.go` | Add Rick option |
| `internal/tui/screens/preset.go` | Add FullRick option |
| `internal/tui/screens/persona_test.go` | Update for Rick |
| `internal/tui/screens/preset_test.go` | Update for FullRick |
| `internal/tui/model.go` | Handle preset default persona selection |
| `internal/tui/model_test.go` | Update preset flow tests |
