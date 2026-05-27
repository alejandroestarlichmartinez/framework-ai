# Proposal: Add CodeGraph as a Selectable Component

## Intent

Add **CodeGraph** (code knowledge graph) as a framework-ai component, following the exact Engram/Context7/GGA pattern.

## Scope

### In Scope
- Register in types, catalog, planner graph, TUI presets
- Global binary installation via curl; per-agent MCP injection
- System prompt instruction injection
- CodeGraph per-project `.codegraph/` init detection
- Uninstall, sync, and full test coverage

### Out of Scope
- CodeGraph MCP server or binary development (external)
- RTK integration (deferred to future iteration)

## Capabilities

### New Capabilities
- `codegraph-component`: Install, MCP injection, init detection, instructions

### Modified Capabilities
- None

## Approach

Follow the **Engram pattern**:

1. **Types/Graph/Catalog**: Add constant, description, `nil` deps
2. **TUI**: Include in `PresetFullGentleman`
3. **Install**: Download binary, inject MCP, inject system prompts
4. **Sync**: Inject-only (no binary install)
5. **Uninstall**: Remove MCP configs, system prompt sections
6. **Resolver**: Add curl-based install command
7. **New package**: `internal/components/codegraph/`

## Affected Areas

| Area | Impact |
|------|--------|
| `internal/model/types.go` | Add constant |
| `internal/catalog/components.go` | Add to `mvpComponents` |
| `internal/planner/graph.go` | Add to `MVPGraph` |
| `internal/tui/model.go` | Add to preset components |
| `internal/cli/run.go` | Add apply + path cases |
| `internal/cli/validate.go` | Add to validation |
| `internal/cli/sync.go` | Add sync cases |
| `internal/installcmd/resolver.go` | Add install command |
| `internal/components/codegraph/` | New: MCP + instructions + init |
| `internal/components/uninstall/service.go` | Add uninstall ops |
| `*_test.go` / testdata | Update tests and golden files |

## Risks

| Risk | Likelihood | Mitigation |
|------|------------|------------|
| Install script URL unavailable | Low | Pin to commit; document fallback |
| CodeGraph MCP format changes | Low | Isolate in adapter strategies |
| CodeGraph init fails | Low | Skip gracefully outside projects |

## Rollback Plan

1. Revert PR (all changes additive)
2. Run `framework-ai uninstall --agent <agent>` with CodeGraph selected
3. No migrations required

## Dependencies

- CodeGraph: `curl -fsSL https://raw.githubusercontent.com/colbymchenry/codegraph/main/install.sh | sh`

## Success Criteria

- [ ] Component appears in picker and presets
- [ ] Global curl install works on Linux/macOS
- [ ] MCP config injects for all supported agents
- [ ] Instructions inject and clean up correctly
- [ ] Uninstall removes all managed config
- [ ] Sync refreshes without reinstalling binary
- [ ] `go test ./...` passes
