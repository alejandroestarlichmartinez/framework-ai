# Design: Add CodeGraph as a Selectable Component

## Technical Approach

Follow the **Engram component pattern** exactly. CodeGraph is a binary-plus-MCP component like Engram, not a config-only component like Context7. The implementation adds a new `internal/components/codegraph/` package, wires it into the existing catalog/planner/TUI/CLI/uninstall pipelines, and reuses all existing adapter strategies for MCP injection and system prompt section management.

## Architecture Decisions

| Decision | Options | Tradeoffs | Choice |
|----------|---------|-----------|--------|
| Install method | a) Run official curl \| sh script<br>b) Download binary manually (like Engram Linux) | Script is official, handles all platforms, auto-updates; manual download requires per-OS asset URL logic | **Run the script** — it's the supported installer |
| System prompt markers | a) `<!-- framework-ai:codegraph -->`<br>b) Reuse engram-protocol marker | Separate marker allows independent uninstall and avoids conflating memory vs. knowledge-graph instructions | **Separate `codegraph` marker** |
| Per-project init detection | a) Check at `componentApplyStep.Run()` time<br>b) Check at TUI selection time | Run-time check captures the current workspace state accurately; selection-time may stale if user switches directories | **Check at apply-step run time** |
| Codex TOML support | a) Upsert `[mcp_servers.codegraph]`<br>b) Skip Codex (like Context7) | CodeGraph MCP works with Codex; skipping would leave a gap | **Upsert TOML block** |

## Data Flow

```
Install Flow:
  CLI run.go ──→ componentApplyStep.Run()
                      │
                      ├──→ codegraph.InstallCommand() → curl | sh
                      ├──→ codegraph.Inject() per adapter
                      │         ├──→ MCP config (strategy-specific)
                      │         └──→ System prompt section (marker-based)
                      │               └──→ Detect .codegraph/ → pick variant
                      └──→ backup/verify

Sync Flow:
  CLI sync.go ──→ componentSyncStep.Run()
                      │
                      └──→ codegraph.Inject() only (no binary reinstall)

Uninstall Flow:
  uninstall/service.go ──→ componentOperations()
                                └──→ codegraphTargets() + codegraphOperations()
                                      ├──→ Remove MCP entries per strategy
                                      └──→ Strip <!-- framework-ai:codegraph --> sections
```

## File Changes

| File | Action | Description |
|------|--------|-------------|
| `internal/model/types.go` | Modify | Add `ComponentCodeGraph` constant |
| `internal/catalog/components.go` | Modify | Add CodeGraph to `mvpComponents` |
| `internal/planner/graph.go` | Modify | Add `ComponentCodeGraph` to `MVPGraph()` with `nil` deps |
| `internal/tui/model.go` | Modify | Add CodeGraph to `PresetFullGentleman` via `componentsForPreset()` |
| `internal/cli/run.go` | Modify | Add `ComponentCodeGraph` case in `componentApplyStep.Run()` |
| `internal/cli/validate.go` | Modify | Add CodeGraph to `componentsForPreset()` |
| `internal/cli/sync.go` | Modify | Add `ComponentCodeGraph` case in `componentSyncStep.Run()` |
| `internal/installcmd/resolver.go` | Modify | Add `resolveCodeGraphInstall()` → curl script CommandSequence |
| `internal/components/codegraph/install.go` | Create | Wrap `installcmd.NewResolver().ResolveComponentInstall()` |
| `internal/components/codegraph/inject.go` | Create | MCP injection + system prompt injection with `.codegraph/` detection |
| `internal/components/codegraph/setup.go` | Create | No-op setup stub (CodeGraph has no per-agent setup like Engram) |
| `internal/components/uninstall/service.go` | Modify | Add `codegraphTargets()` and `codegraphOperations()` |
| `internal/assets/claude/codegraph.md` | Create | Full instructions when `.codegraph/` exists |
| `internal/assets/claude/codegraph-init.md` | Create | Minimal prompt when `.codegraph/` absent |
| `*_test.go` | Modify | Update tests for new component count, golden files, uninstall coverage |

## Interfaces / Contracts

```go
// internal/components/codegraph/inject.go
func Inject(homeDir string, adapter agents.Adapter) (InjectionResult, error)

// internal/components/codegraph/install.go
func InstallCommand(profile system.PlatformProfile) ([][]string, error)

// internal/installcmd/resolver.go
func resolveCodeGraphInstall(profile system.PlatformProfile) CommandSequence {
    return CommandSequence{{"sh", "-c",
        "curl -fsSL https://raw.githubusercontent.com/colbymchenry/codegraph/main/install.sh | sh"}}
}
```

**MCP config shapes per strategy:**
- `StrategySeparateMCPFiles`: `{"command": "codegraph", "args": ["serve", "--mcp"]}`
- `StrategyMergeIntoSettings`: `mcpServers.codegraph` or `mcp.codegraph` (OpenCode)
- `StrategyMCPConfigFile`: `mcpServers.codegraph` or `servers.codegraph` (VS Code)
- `StrategyTOMLFile`: `[mcp_servers.codegraph]` in `~/.codex/config.toml`

**System prompt markers:**
- Open marker: `<!-- framework-ai:codegraph -->`
- Close marker: `<!-- /framework-ai:codegraph -->`

**Per-project detection:**
```go
func hasCodegraphIndex(workspaceDir string) bool {
    if workspaceDir == "" { return false }
    _, err := os.Stat(filepath.Join(workspaceDir, ".codegraph"))
    return err == nil
}
```

## Testing Strategy

| Layer | What to Test | Approach |
|-------|-------------|----------|
| Unit | `codegraph.Inject()` for all 4 MCP strategies | Table-driven tests with fake adapters |
| Unit | `.codegraph/` detection logic | Tempdir creation + `hasCodegraphIndex()` assertions |
| Unit | `resolveCodeGraphInstall()` | Assert exact CommandSequence output |
| Unit | Uninstall `codegraphOperations()` | Verify correct JSON/TOML paths and markdown strip |
| Integration | `componentApplyStep.Run()` dispatch | Mock adapter + assert Inject called |
| E2E | Full install with CodeGraph selected | Docker E2E run (optional, existing pattern) |

## Migration / Rollout

No migration required. All changes are additive. Users who previously installed without CodeGraph will see it appear in the component picker and presets. Existing installs are unaffected until the user explicitly selects CodeGraph or runs a full-gentleman preset install.

## Open Questions

- [ ] Should CodeGraph be included in `PresetEcosystemOnly` or only `PresetFullGentleman`? (Proposal says full-gentleman only.)
- [ ] Does the CodeGraph install script require `sudo` on Linux? If so, we may need platform-specific preflight checks.