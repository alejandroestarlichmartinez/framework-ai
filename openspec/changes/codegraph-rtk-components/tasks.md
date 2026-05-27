# Tasks: Add CodeGraph as a Selectable Component

## Review Workload Forecast

| Field | Value |
|-------|-------|
| Estimated changed lines | 450–550 |
| 400-line budget risk | High |
| Chained PRs recommended | Yes |
| Suggested split | PR 1: Infrastructure + Core → PR 2: CLI + Uninstall → PR 3: Tests + Golden |
| Delivery strategy | ask-on-risk |
| Chain strategy | stacked-to-main |

Decision needed before apply: Yes
Chained PRs recommended: Yes
Chain strategy: stacked-to-main
400-line budget risk: High

### Suggested Work Units

| Unit | Goal | Likely PR | Notes |
|------|------|-----------|-------|
| 1 | Types, catalog, graph, TUI, resolver, core package skeleton | PR 1 | Base: main; includes RED tests |
| 2 | CLI wiring (run, sync, validate), uninstall | PR 2 | Base: main; depends on PR 1 |
| 3 | Full test suite, golden files, integration | PR 3 | Base: main; depends on PR 2 |

---

## Phase 1: Infrastructure

- [x] 1.1 **RED** Add `ComponentCodeGraph` constant to `internal/model/types.go`; write failing assertion in `types_test.go`
- [x] 1.2 **GREEN** Make constant pass; update any golden files that enumerate components
- [x] 1.3 **RED** Register CodeGraph in `internal/catalog/components.go` (`mvpComponents`); write catalog test
- [x] 1.4 **GREEN** Make catalog pass; add description "Semantic code knowledge graph"
- [x] 1.5 **RED** Add `ComponentCodeGraph` with `nil` deps to `internal/planner/graph.go` (`MVPGraph`); write planner test
- [x] 1.6 **GREEN** Make planner pass; verify topological order unchanged
- [x] 1.7 Add CodeGraph to `PresetFullGentleman` in `internal/tui/model.go`; update TUI preset test
- [x] 1.8 Add `resolveCodeGraphInstall()` to `internal/installcmd/resolver.go` returning curl script; test exact `CommandSequence`

## Phase 2: Core Component

- [x] 2.1 Create `internal/components/codegraph/install.go` wrapping resolver; test returns correct command sequence
- [x] 2.2 Create `internal/assets/claude/codegraph.md` (full instructions for `.codegraph/` present)
- [x] 2.3 Create `internal/assets/claude/codegraph-init.md` (minimal prompt when `.codegraph/` absent)
- [x] 2.4 **RED** Create `internal/components/codegraph/inject.go` with `Inject()` signature and `hasCodegraphIndex()`; write failing table-driven tests for all 4 MCP strategies
- [x] 2.5 **GREEN** Implement MCP injection per strategy (Separate, Merge, MCPConfig, TOML); verify all tests pass
- [x] 2.6 **RED** Add system prompt injection with `.codegraph/` detection; test full vs minimal prompt selection
- [x] 2.7 **GREEN** Implement marker-based injection (`<!-- framework-ai:codegraph -->`); verify strip/upsert logic
- [x] 2.8 Create `internal/components/codegraph/setup.go` as no-op stub; test returns nil error

## Phase 3: CLI Integration

- [x] 3.1 Wire `ComponentCodeGraph` case into `internal/cli/run.go` `componentApplyStep.Run()`; dispatch to `codegraph.InstallCommand()` then `codegraph.Inject()`
- [x] 3.2 Wire `ComponentCodeGraph` case into `internal/cli/sync.go` `componentSyncStep.Run()`; inject only (no binary reinstall)
- [x] 3.3 Update `internal/cli/validate.go` `componentsForPreset()` to include CodeGraph in validation
- [x] 3.4 Add sync idempotency check: skip file write when config already current; test no-op scenario

## Phase 4: Uninstall

- [x] 4.1 Add `codegraphTargets()` to `internal/components/uninstall/service.go` listing managed paths per strategy
- [x] 4.2 Add `codegraphOperations()` removing MCP entries and stripping `<!-- framework-ai:codegraph -->` sections; test JSON/TOML removal
- [x] 4.3 Verify uninstall leaves other component markers intact; test engram/context7 markers untouched

## Phase 5: Testing

- [x] 5.1 Unit test `codegraph.Inject()` with fake adapters for all 4 strategies (spec: Scenario Separate files, Merge, MCP config, TOML)
- [x] 5.2 Unit test `.codegraph/` detection: temp dir with/without index (spec: Scenario Detection, Index present, Index absent)
- [x] 5.3 Unit test `componentApplyStep.Run()` dispatch: mock adapter, assert Inject called once (spec: Scenario Install command)
- [x] 5.4 Unit test sync no-op: binary exists, configs current, no download executes (spec: Scenario No binary reinstall, Idempotent sync)
- [x] 5.5 Integration test uninstall end-to-end: inject then uninstall, assert zero CodeGraph references remain (spec: Scenario Full uninstall)
- [x] 5.6 Update all golden files and `go test ./...`; verify `go vet ./...` clean
