---
description: Audit project compliance against AGENTS.md and coding conventions
---

If the native `sdd-inspect` sub-agent is available, delegate this command to it.
Otherwise, read the skill file at `~/.claude/skills/sdd-inspect/SKILL.md` FIRST, then follow its instructions exactly inline.

CONTEXT:
- Working directory: !`pwd`
- Current project: !`basename "$(pwd)"`
- Artifact store mode: engram

TASK:
Audit the project for compliance. Read AGENTS.md, project conventions, and git history. Then:

ENGRAM PERSISTENCE (artifact store mode: engram):
CRITICAL: mem_search returns 300-char PREVIEWS, not full content. You MUST call mem_get_observation(id) for EVERY artifact.
Save report:
  mem_save(title: "sdd-inspect/{timestamp-or-sha}", topic_key: "sdd-inspect/{timestamp-or-sha}", type: "architecture", project: "{project}", capture_prompt: false, content: "{inspect report}")
  Set capture_prompt: false when the Engram tool schema supports it; if an older schema rejects or does not expose the field, omit it rather than failing.

Then:
1. Check AGENTS.md compliance
2. Check Go code quality (if applicable)
3. Check git hygiene
4. Check SDD compliance (if .atl/ or openspec/ exists)
5. Check file hygiene (secrets, merge conflicts, etc.)

Return a structured inspect report with: status, executive_summary, detailed_report, artifacts, and next_recommended.
