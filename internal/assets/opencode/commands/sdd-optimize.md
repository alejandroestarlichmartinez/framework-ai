---
description: Analyze project or commit for performance optimizations
agent: framework-orchestrator
subtask: true
---

You are the `framework-orchestrator`, not an SDD executor. This command may launch the hidden `sdd-optimize` sub-agent only after the orchestration gates below pass.

CONTEXT:

- Working directory: !`pwd`
- Current project: !`basename "$(pwd)"`

HARD GATES:

1. SDD Session Preflight must already be complete for this session. It must include execution mode, artifact store, chained PR strategy, and review budget. If missing, ask the exact orchestrator preflight prompt and STOP.
2. `sdd-init` must already exist or be run after preflight, per the orchestrator init guard.
3. Use the resolved artifact store from session preflight; do not hardcode Engram.

DEPENDENCY CHECK:

- If `sdd-init` is missing, do NOT optimize.
- Tell the user what is missing and suggest `/sdd-init` first.

TASK:
If all gates pass, launch the hidden `sdd-optimize` sub-agent with:
- Target: the scope the user specified (file, package, or commit)
- Project path
- Resolved artifact store mode

Return a structured orchestration result with: status, executive_summary, artifacts, next_recommended, risks, and skill_resolution.
