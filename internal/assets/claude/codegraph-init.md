## CodeGraph — Conditional Setup

**At the start of every session in this workspace, verify whether `.codegraph/` exists.**

### If `.codegraph/` does NOT exist

Run `codegraph init -i` IMMEDIATELY. Do NOT ask the user for permission. This is automatic one-time setup, identical to how Engram saves memory automatically. The user already chose to install CodeGraph as a component.

### If `.codegraph/` exists

Use CodeGraph tools for all code exploration instead of file-reading or grep:

- `codegraph_search` — Find symbols, functions, types by name
- `codegraph_explore` — Navigate call graphs and dependencies  
- `codegraph_query` — Run structured queries
- `codegraph_summarize` — Generate module summaries

**Prefer CodeGraph tools** over reading files directly. Fall back to file reading only when CodeGraph returns no results.

The `.codegraph/` directory contains the indexed knowledge graph — do not modify it manually.
