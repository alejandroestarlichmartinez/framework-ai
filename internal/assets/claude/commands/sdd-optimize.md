---
description: Analyze project or commit for performance optimizations
---

If the native `sdd-optimize` sub-agent is available, delegate this command to it.
Otherwise, read the skill file at `~/.claude/skills/sdd-optimize/SKILL.md` FIRST, then follow its instructions exactly inline.

CONTEXT:
- Working directory: !`pwd`
- Current project: !`basename "$(pwd)"`
- Artifact store mode: engram

TASK:
Analyze the project or specified commit for performance optimizations. Then:

ENGRAM PERSISTENCE (artifact store mode: engram):
CRITICAL: mem_search returns 300-char PREVIEWS, not full content. You MUST call mem_get_observation(id) for EVERY artifact.
Save report:
  mem_save(title: "sdd-optimize/{timestamp-or-sha}", topic_key: "sdd-optimize/{timestamp-or-sha}", type: "architecture", project: "{project}", capture_prompt: false, content: "{optimize report}")
  Set capture_prompt: false when the Engram tool schema supports it; if an older schema rejects or does not expose the field, omit it rather than failing.

Then:
1. Identify performance bottlenecks
2. Propose concrete optimizations
3. Estimate impact of each optimization
4. Generate a correction proposal if applicable

Return a structured optimize report with: status, executive_summary, detailed_report, artifacts, and next_recommended.
