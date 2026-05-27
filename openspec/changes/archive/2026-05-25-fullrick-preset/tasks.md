# Tasks: FullRick Preset and Rick Sanchez Persona

## Task Conventions

- **TDD Required:** Every task starts with a failing test, then writes the minimal code to pass.
- **Timebox:** Each task is designed for ~30 minutes. If a task exceeds 45 minutes, stop and flag it.
- **Commit per Task:** Each task is a single, reviewable commit.
- **No Skipping Tests:** Do not write production code without a corresponding test (even if the test is trivial).

---

## Phase 1: Types + Constants

### Task 1.1 — Add PersonaRickSanchez and PresetFullRick constants

**Test First:**
- Add `TestPersonaConstants` in `internal/model/types_test.go` (create if missing) asserting:
  - `PersonaRickSanchez == "rick-sanchez"`
  - `PresetFullRick == "full-rick"`
- Run `go test ./internal/model` — expect failure (constants don't exist yet).

**Production:**
- Add `PersonaRickSanchez PersonaID = "rick-sanchez"` to `types.go`.
- Add `PresetFullRick PresetID = "full-rick"` to `types.go`.

**Verify:**
- `go test ./internal/model` passes.
- `go build ./...` succeeds.

**Commit:** `feat(types): add PersonaRickSanchez and PresetFullRick constants`

---

## Phase 2: Asset Files (All Agents)

### Task 2.1 — Create generic Rick Sanchez persona asset

**Test First:**
- Add `TestRickSanchezAssetExists` in `internal/assets/assets_test.go` asserting `assets.MustRead("generic/persona-rick-sanchez.md")` does not panic and contains `"## Personality"`.
- Run test — expect failure (file missing).

**Production:**
- Create `internal/assets/generic/persona-rick-sanchez.md`.
- Structure: `## Rules`, `## Personality`, `## Persona Scope`, `## Language`, `## Tone`, `## Philosophy`, `## Expertise`, `## Behavior`, `## Contextual Skill Loading (MANDATORY)`.
- Personality: cynical genius scientist, multiverse traveler, tired of bad code but cares deep down.
- Language: neutral Spanish (no voseo, no regionalisms) when user writes in Spanish.
- All functional guardrails from Gentleman preserved.

**Verify:**
- Test passes.
- File is readable via `assets.MustRead`.

**Commit:** `feat(assets): add generic Rick Sanchez persona asset`

### Task 2.2 — Create Claude-specific Rick Sanchez persona asset

**Test First:**
- Add `TestRickSanchezClaudeAssetExists` asserting `assets.MustRead("claude/persona-rick-sanchez.md")` does not panic.
- Run test — expect failure.

**Production:**
- Create `internal/assets/claude/persona-rick-sanchez.md`.
- Adapt generic content for Claude Code context (same structure, minor agent-specific phrasing if needed; can be identical to generic if no divergence required).

**Verify:** Test passes.

**Commit:** `feat(assets): add Claude Rick Sanchez persona asset`

### Task 2.3 — Create OpenCode/Kilocode Rick Sanchez persona asset

**Test First:**
- Add `TestRickSanchezOpencodeAssetExists` asserting `assets.MustRead("opencode/persona-rick-sanchez.md")` does not panic.
- Run test — expect failure.

**Production:**
- Create `internal/assets/opencode/persona-rick-sanchez.md`.
- Same structure as generic, adapted for OpenCode/Kilocode shared context.

**Verify:** Test passes.

**Commit:** `feat(assets): add OpenCode Rick Sanchez persona asset`

### Task 2.4 — Create Kimi Rick Sanchez persona asset

**Test First:**
- Add `TestRickSanchezKimiAssetExists` asserting `assets.MustRead("kimi/persona-rick-sanchez.md")` does not panic.
- Run test — expect failure.

**Production:**
- Create `internal/assets/kimi/persona-rick-sanchez.md`.
- Same structure, adapted for Kimi Jinja module context.

**Verify:** Test passes.

**Commit:** `feat(assets): add Kimi Rick Sanchez persona asset`

### Task 2.5 — Create Kiro Rick Sanchez persona asset

**Test First:**
- Add `TestRickSanchezKiroAssetExists` asserting `assets.MustRead("kiro/persona-rick-sanchez.md")` does not panic.
- Run test — expect failure.

**Production:**
- Create `internal/assets/kiro/persona-rick-sanchez.md`.
- Same structure, adapted for Kiro steering file context.

**Verify:** Test passes.

**Commit:** `feat(assets): add Kiro Rick Sanchez persona asset`

---

## Phase 3: Injection Logic

### Task 3.1 — Introduce isFullPersona() helper with tests

**Test First:**
- In `internal/components/persona/inject_test.go`, add `TestIsFullPersona`:
  - `isFullPersona(PersonaGentleman)` → `true`
  - `isFullPersona(PersonaRickSanchez)` → `true`
  - `isFullPersona(PersonaNeutral)` → `false`
  - `isFullPersona(PersonaCustom)` → `false`
- Run test — expect failure (function doesn't exist).

**Production:**
- Add `isFullPersona()` to `inject.go`.

**Verify:**
- `go test ./internal/components/persona -run TestIsFullPersona` passes.

**Commit:** `refactor(persona): introduce isFullPersona helper`

### Task 3.2 — Refactor output-style gating to use isFullPersona()

**Test First:**
- Add `TestInjectClaudeRickSanchezWritesOutputStyle`:
  - `Inject(home, claudeAdapter(), PersonaRickSanchez)`
  - Assert `output-styles/gentleman.md` exists.
  - Assert `settings.json` has `"outputStyle": "Gentleman"`.
- Run test — expect failure (gating still hardcoded to Gentleman only).

**Production:**
- Replace `persona == model.PersonaGentleman` with `isFullPersona(persona)` in output-style write block (line ~346).
- Replace `persona != model.PersonaGentleman` with `!isFullPersona(persona)` in cleanup block (line ~374).
- Replace `persona == model.PersonaGentleman` in Kimi Jinja module block (line ~302).

**Verify:**
- New Rick output-style test passes.
- All existing output-style tests still pass (Gentleman, Neutral, switch tests).

**Commit:** `refactor(persona): gate output styles on isFullPersona`

### Task 3.3 — Refactor agent overlay gating to use isFullPersona()

**Test First:**
- Add `TestInjectOpenCodeRickSanchezCreatesAgentOverlay`:
  - `Inject(home, opencodeAdapter(), PersonaRickSanchez)`
  - Assert `opencode.json` contains `"gentleman"` agent key.
  - Assert it does NOT contain `"sdd-orchestrator"`.
- Run test — expect failure.

**Production:**
- Replace `persona == model.PersonaGentleman` with `isFullPersona(persona)` in agent overlay creation block (line ~322).
- Replace the else branch (remove agent key) to use `!isFullPersona(persona)` instead of `persona != model.PersonaGentleman` (line ~329-341).

**Verify:**
- New Rick agent overlay test passes.
- All existing agent overlay tests still pass.

**Commit:** `refactor(persona): gate agent overlay on isFullPersona`

### Task 3.4 — Refactor preserveManagedSections gating

**Test First:**
- Add `TestPreserveManagedSectionsRickSanchez`:
  - Call `preserveManagedSections(existing, newRickContent, PersonaRickSanchez)`
  - Assert it returns `("", false)` because Rick is a full persona and does not need section preservation.
- Run test — expect failure.

**Production:**
- Replace `persona == model.PersonaGentleman` with `isFullPersona(persona)` in `preserveManagedSections()` (line ~523).

**Verify:**
- New test passes.
- Existing preservation tests (Neutral, VSCode) still pass.

**Commit:** `refactor(persona): gate preserveManagedSections on isFullPersona`

### Task 3.5 — Extend personaContent() dispatch for Rick

**Test First:**
- Add `TestPersonaContentRickSanchezDispatch`:
  - For each agent (Claude, OpenCode, Kilocode, Kimi, Kiro, Gemini, Cursor, VSCode, Codex, Windsurf, Antigravity, Qwen, OpenClaw, Pi):
    - Call `personaContent(agent, PersonaRickSanchez)`
    - Assert result is non-empty.
    - Assert result contains `"## Personality"`.
  - For Claude specifically, assert it does NOT equal the generic asset (verifies per-agent dispatch).
- Run test — expect failure (dispatch missing Rick cases).

**Production:**
- Extend `personaContent()` switch to handle `PersonaRickSanchez` with the same per-agent structure as Gentleman.

**Verify:**
- Dispatch test passes.
- Existing `personaContent` tests (Gentleman, Neutral) still pass.

**Commit:** `feat(persona): add Rick Sanchez dispatch to personaContent`

### Task 3.6 — Add Rick injection idempotency test

**Test First:**
- Add `TestInjectClaudeRickSanchezIdempotent`:
  - First inject → `Changed=true`
  - Second inject → `Changed=false`
- Run test — expect failure if idempotency is broken.

**Production:**
- No code changes expected — idempotency should be inherited from existing file-compare logic.
- If test fails, fix the root cause (likely in how content is compared).

**Verify:** Test passes.

**Commit:** `test(persona): verify Rick Sanchez injection idempotency`

---

## Phase 4: Preset + Validation

### Task 4.1 — Update CLI persona validation

**Test First:**
- In `internal/cli/validate_test.go`, add `TestNormalizePersonaAcceptsRickSanchez`:
  - `normalizePersona("rick-sanchez")` → `PersonaRickSanchez`, no error.
- Run test — expect failure.

**Production:**
- Update `normalizePersona()` to accept `PersonaRickSanchez`.
- Update signature to accept preset for default derivation (or derive default before calling `normalizePersona`).

**Verify:**
- New test passes.
- Existing persona validation tests still pass.

**Commit:** `feat(cli): accept rick-sanchez persona in validation`

### Task 4.2 — Update CLI preset validation

**Test First:**
- Add `TestNormalizePresetAcceptsFullRick`:
  - `normalizePreset("full-rick")` → `PresetFullRick`, no error.
- Add `TestComponentsForPresetFullRick`:
  - `componentsForPreset(PresetFullRick)` returns same list as `PresetFullGentleman`.
- Run tests — expect failure.

**Production:**
- Update `normalizePreset()` to accept `PresetFullRick`.
- Update `componentsForPreset()` to handle `PresetFullRick` (same list as FullGentleman).

**Verify:** Tests pass.

**Commit:** `feat(cli): accept full-rick preset in validation`

### Task 4.3 — Derive default persona from preset

**Test First:**
- Add `TestNormalizeInstallFlags_FullRickDefaultsToRickSanchez`:
  - Input: `InstallFlags{Preset: "full-rick"}`
  - Assert: `Selection.Persona == PersonaRickSanchez`.
- Add `TestNormalizeInstallFlags_FullGentlemanDefaultsToGentleman`:
  - Same for `full-gentleman` → `PersonaGentleman`.
- Add `TestNormalizeInstallFlags_PresetWithExplicitPersonaOverridesDefault`:
  - Input: `InstallFlags{Preset: "full-rick", Persona: "gentleman"}`
  - Assert: `Selection.Persona == PersonaGentleman`.
- Run tests — expect failure.

**Production:**
- Modify `NormalizeInstallFlags` to derive default persona from preset when `--persona` is omitted.
- Pass preset into persona normalization logic.

**Verify:**
- All new tests pass.
- Existing preset/persona integration tests still pass.

**Commit:** `feat(cli): derive default persona from selected preset`

---

## Phase 5: TUI Integration

### Task 5.1 — Add Rick to persona picker

**Test First:**
- In `internal/tui/screens/persona_test.go`, update `TestPersonaOptions`:
  - Assert `PersonaOptions()` includes `PersonaRickSanchez`.
  - Assert order is: Gentleman, RickSanchez, Neutral, Custom.
- Add `TestRenderPersonaIncludesRickDescription`:
  - Render with Rick selected.
  - Assert output contains `"cynical genius"`.
- Run tests — expect failure.

**Production:**
- Update `PersonaOptions()` to include `PersonaRickSanchez`.
- Update `personaDescriptions` with Rick description.

**Verify:** Tests pass.

**Commit:** `feat(tui): add Rick Sanchez to persona picker`

### Task 5.2 — Add FullRick to preset picker

**Test First:**
- In `internal/tui/screens/preset_test.go`, update `TestPresetOptions`:
  - Assert `PresetOptions()` includes `PresetFullRick`.
  - Assert order is: FullGentleman, FullRick, EcosystemOnly, Minimal, Custom.
- Add `TestRenderPresetIncludesFullRickDescription`:
  - Assert output contains `"Rick Sanchez"`.
- Run tests — expect failure.

**Production:**
- Update `PresetOptions()` to include `PresetFullRick`.
- Update `presetDescriptions` with FullRick description.

**Verify:** Tests pass.

**Commit:** `feat(tui): add FullRick preset to preset picker`

### Task 5.3 — Wire preset-to-persona default in TUI model

**Test First:**
- In `internal/tui/model_test.go`, add `TestPresetFullRickSetsPersonaToRick`:
  - Simulate preset screen confirm with `PresetFullRick`.
  - Assert `m.Selection.Persona == PersonaRickSanchez`.
- Add `TestPresetFullRickDoesNotOverrideExplicitPersona`:
  - Set `m.Selection.Persona = PersonaNeutral` before preset confirm.
  - Confirm `PresetFullRick`.
  - Assert persona stays `PersonaNeutral`.
- Run tests — expect failure.

**Production:**
- In `internal/tui/model.go`, `ScreenPreset` confirm handler:
  - After setting preset and components, if preset is `PresetFullRick` AND current persona is the default (`PersonaGentleman`), update to `PersonaRickSanchez`.

**Verify:** Tests pass.

**Commit:** `feat(tui): auto-select Rick persona when FullRick preset chosen`

### Task 5.4 — Update TUI componentsForPreset

**Test First:**
- In `internal/tui/model_test.go`, add `TestTUIComponentsForPresetFullRick`:
  - `componentsForPreset(PresetFullRick)` returns same list as `PresetFullGentleman`.
- Run test — expect failure.

**Production:**
- Update `componentsForPreset()` in `internal/tui/model.go` to handle `PresetFullRick`.

**Verify:** Test passes.

**Commit:** `feat(tui): handle FullRick preset in component resolution`

---

## Phase 6: Tests + Golden Files

### Task 6.1 — Create golden file infrastructure

**Test First:**
- Add `TestRickGoldenFiles` in `inject_test.go`:
  - Inject Rick into Claude, OpenCode, Kimi temp dirs.
  - Compare against `testdata/claude-rick-sanchez.golden.md`.
  - Expect failure (golden files don't exist yet).

**Production:**
- Create `internal/components/persona/testdata/` directory.
- Run tests with `-update` flag to generate initial golden files.
- Verify golden files contain expected markers (`<!-- framework-ai:persona -->`, `## Personality`).

**Verify:**
- `go test ./internal/components/persona -run TestRickGoldenFiles` passes.
- Golden files are committed.

**Commit:** `test(persona): add golden files for Rick Sanchez injection`

### Task 6.2 — Add Rick switch-to-neutral cleanup test

**Test First:**
- Add `TestInjectClaudeRickToNeutralCleansOutputStyle`:
  1. Inject Rick → assert output-style exists and settings has `"Gentleman"`.
  2. Inject Neutral → assert output-style removed and settings key removed.
- Run test — expect failure if cleanup gating is incorrect.

**Production:**
- No changes expected if Task 3.2 was done correctly.
- Fix any remaining hardcoded checks if test fails.

**Verify:** Test passes.

**Commit:** `test(persona): verify Rick-to-Neutral cleanup`

### Task 6.3 — Add comprehensive Rick injection coverage

**Test First:**
- Add `TestInjectOpenCodeRickSanchez`:
  - Assert `AGENTS.md` has Rick content and markers.
  - Assert `opencode.json` has `agent.gentleman`.
- Add `TestInjectKimiRickSanchez`:
  - Assert `persona.md` module has Rick content.
  - Assert `output-style.md` module exists.
- Add `TestInjectKiroRickSanchez`:
  - Assert steering file has Rick content and frontmatter.
- Run tests — expect failure.

**Production:**
- No production code changes expected — all infrastructure should be in place.
- Fix any failing tests.

**Verify:** All tests pass.

**Commit:** `test(persona): comprehensive Rick Sanchez injection coverage`

### Task 6.4 — Final integration test suite

**Test First:**
- Add `TestNormalizeInstallFlags_InvalidRickSanchezPreset`:
  - `normalizePreset("rick-sanchez")` should fail (it's a persona, not a preset).
- Add `TestNormalizeInstallFlags_InvalidFullRickPersona`:
  - `normalizePersona("full-rick")` should fail (it's a preset, not a persona).
- Run tests — verify they fail appropriately.

**Production:**
- No changes needed — validation already rejects unknown values.
- These tests document the boundary.

**Verify:** Tests pass.

**Commit:** `test(cli): boundary tests for preset/persona confusion`

### Task 6.5 — Full test suite verification

**No new code. Pure verification.**

- Run `go test ./...` — assert all tests pass.
- Run `go build ./...` — assert no compilation errors.
- Run `go vet ./...` — assert no warnings.
- Run `gofmt -l .` — assert no unformatted files.

**If failures:** Fix in place, commit as `fix: address test suite failures`.

**Commit:** `test: verify full test suite passes`

---

## Summary Commit Sequence

1. `feat(types): add PersonaRickSanchez and PresetFullRick constants`
2. `feat(assets): add generic Rick Sanchez persona asset`
3. `feat(assets): add Claude Rick Sanchez persona asset`
4. `feat(assets): add OpenCode Rick Sanchez persona asset`
5. `feat(assets): add Kimi Rick Sanchez persona asset`
6. `feat(assets): add Kiro Rick Sanchez persona asset`
7. `refactor(persona): introduce isFullPersona helper`
8. `refactor(persona): gate output styles on isFullPersona`
9. `refactor(persona): gate agent overlay on isFullPersona`
10. `refactor(persona): gate preserveManagedSections on isFullPersona`
11. `feat(persona): add Rick Sanchez dispatch to personaContent`
12. `test(persona): verify Rick Sanchez injection idempotency`
13. `feat(cli): accept rick-sanchez persona in validation`
14. `feat(cli): accept full-rick preset in validation`
15. `feat(cli): derive default persona from selected preset`
16. `feat(tui): add Rick Sanchez to persona picker`
17. `feat(tui): add FullRick preset to preset picker`
18. `feat(tui): auto-select Rick persona when FullRick preset chosen`
19. `feat(tui): handle FullRick preset in component resolution`
20. `test(persona): add golden files for Rick Sanchez injection`
21. `test(persona): verify Rick-to-Neutral cleanup`
22. `test(persona): comprehensive Rick Sanchez injection coverage`
23. `test(cli): boundary tests for preset/persona confusion`
24. `test: verify full test suite passes`

---

## Rollback Instructions

If any task introduces regressions:

1. Revert the commit(s) for that task using `git revert <hash>`.
2. Run `go test ./...` to confirm baseline.
3. Resume from the previous stable task.

For full rollback of the entire feature:

```bash
git revert --no-commit HEAD~24..HEAD  # adjust range as needed
rm internal/assets/*/persona-rick-sanchez.md
go test ./...
```
