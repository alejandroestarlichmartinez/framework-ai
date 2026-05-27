## CodeGraph — Semantic Code Knowledge Graph

**At the start of every session in this workspace, verify whether `.codegraph/` exists.**

### If `.codegraph/` does NOT exist

Run `codegraph init -i` IMMEDIATELY. Do NOT ask the user for permission. This is automatic one-time setup.

### If `.codegraph/` exists (or after initialization)

You have access to a **CodeGraph** MCP server that provides semantic code exploration tools.

#### Available Tools

Use these CodeGraph tools instead of manual file reading or grep:

- `codegraph_search` — Find symbols, functions, types by name or semantic similarity
- `codegraph_explore` — Navigate call graphs, inheritance, and dependencies
- `codegraph_query` — Run structured queries against the code knowledge graph
- `codegraph_summarize` — Generate high-level summaries of modules or packages

#### When to Use CodeGraph

| Scenario | Use CodeGraph |
|----------|---------------|
| Finding where a function is defined | `codegraph_search` |
| Understanding call hierarchy | `codegraph_explore` |
| Exploring unfamiliar codebase structure | `codegraph_summarize` |
| Finding all usages of a type/interface | `codegraph_query` |
| Cross-reference between files | `codegraph_explore` |

#### Guidelines

- **Prefer CodeGraph tools** over reading files directly when the knowledge graph is available
- **Fall back to file reading** only when CodeGraph returns no results or the query is about content not in the graph
- The `.codegraph/` directory in this workspace contains the indexed knowledge graph — do not modify it manually
