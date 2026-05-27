# CodeGraph Component Specification

## Purpose

Define behavior for adding CodeGraph as a selectable framework-ai component, following the Engram/Context7 pattern.

## Requirements

### Requirement: Component Registration

The system MUST register CodeGraph as ComponentID `codegraph` with nil dependencies in the planner graph, include it in the component catalog, and add it to PresetFullGentleman.

#### Scenario: Registration

- GIVEN framework-ai initializes
- WHEN components are listed
- THEN CodeGraph appears with name "CodeGraph", description "Semantic code knowledge graph", nil deps, and is in PresetFullGentleman

### Requirement: Installation

The system MUST resolve a curl-based install command and ensure the binary is globally available.

#### Scenario: Install command

- GIVEN a platform profile is known
- WHEN the resolver queries CodeGraph
- THEN it returns `curl -fsSL https://raw.githubusercontent.com/colbymchenry/codegraph/main/install.sh | sh`

#### Scenario: Idempotent install

- GIVEN CodeGraph is already installed
- WHEN the install step runs
- THEN the binary is not re-downloaded

### Requirement: MCP Configuration Injection

The system MUST inject CodeGraph MCP server config per agent strategy using `codegraph serve --mcp`.

#### Scenario: Separate files

- GIVEN the adapter uses StrategySeparateMCPFiles
- WHEN inject runs
- THEN `~/.claude/mcp/codegraph.json` is written with command `codegraph` and args `["serve", "--mcp"]`

#### Scenario: Merge into settings

- GIVEN the adapter uses StrategyMergeIntoSettings
- WHEN inject runs
- THEN `mcpServers.codegraph` is merged into the agent settings file

#### Scenario: MCP config file

- GIVEN the adapter uses StrategyMCPConfigFile
- WHEN inject runs
- THEN the agent mcp.json is merged with a `codegraph` server entry

#### Scenario: TOML file

- GIVEN the adapter uses StrategyTOMLFile
- WHEN inject runs
- THEN `[mcp_servers.codegraph]` is upserted in `~/.codex/config.toml`

### Requirement: System Prompt Injection

The system MUST inject CodeGraph instructions into system prompts when `.codegraph/` exists; otherwise it MUST inject a prompt asking the user to initialize CodeGraph.

#### Scenario: Index present

- GIVEN `.codegraph/` exists in the workspace
- WHEN system prompt injection runs
- THEN the full tool-selection table is injected, instructing the agent to use CodeGraph tools instead of grep

#### Scenario: Index absent

- GIVEN `.codegraph/` does not exist in the workspace
- WHEN system prompt injection runs
- THEN a minimal prompt asks: "Would you like me to run `codegraph init -i` to build a code knowledge graph?"

#### Scenario: Per-strategy injection

- GIVEN the adapter uses any supported SystemPromptStrategy
- WHEN inject runs
- THEN the appropriate CodeGraph block is injected (markdown section, file replace, append, or Jinja module)

### Requirement: Per-Project Init Detection

The system MUST detect `.codegraph/` presence at injection time and select the appropriate prompt variant.

#### Scenario: Detection

- GIVEN the workspace directory
- WHEN the CodeGraph component apply step runs
- THEN the orchestrator checks for `.codegraph/` and passes the full instructions variant if present, else the minimal variant

### Requirement: Uninstall

The system MUST remove all CodeGraph-managed MCP configs and system prompt sections during uninstall.

#### Scenario: MCP cleanup

- GIVEN CodeGraph was injected
- WHEN uninstall runs
- THEN `codegraph` entries are removed from all agent config files per strategy

#### Scenario: Prompt cleanup

- GIVEN CodeGraph instructions exist in system prompts
- WHEN uninstall runs
- THEN `<!-- framework-ai:codegraph -->` sections are stripped from markdown files

#### Scenario: Full uninstall

- GIVEN a complete uninstall is requested
- WHEN the plan executes
- THEN no CodeGraph references remain in managed files

### Requirement: Sync

The system MUST refresh CodeGraph MCP and prompt configs without reinstalling the binary.

#### Scenario: Sync refresh

- GIVEN CodeGraph is selected for sync
- WHEN sync runs
- THEN MCP and prompt configs are re-injected for all target agents

#### Scenario: No binary reinstall

- GIVEN the CodeGraph binary exists
- WHEN sync runs
- THEN no download or install commands execute

#### Scenario: Idempotent sync

- GIVEN all CodeGraph configs are current
- WHEN sync runs
- THEN no files change and sync reports a no-op
