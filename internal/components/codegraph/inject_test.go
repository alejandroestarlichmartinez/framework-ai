package codegraph

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/alejandroestarlichmartinez/framework-ai/internal/agents"
	"github.com/alejandroestarlichmartinez/framework-ai/internal/agents/claude"
	"github.com/alejandroestarlichmartinez/framework-ai/internal/agents/codex"
	"github.com/alejandroestarlichmartinez/framework-ai/internal/agents/gemini"
	"github.com/alejandroestarlichmartinez/framework-ai/internal/agents/opencode"
	"github.com/alejandroestarlichmartinez/framework-ai/internal/agents/vscode"
)

func claudeAdapter() agents.Adapter   { return claude.NewAdapter() }
func opencodeAdapter() agents.Adapter { return opencode.NewAdapter() }
func codexAdapter() agents.Adapter    { return codex.NewAdapter() }
func geminiAdapter() agents.Adapter   { return gemini.NewAdapter() }
func vscodeAdapter() agents.Adapter   { return vscode.NewAdapter() }

func TestInjectClaudeWritesMCPConfig(t *testing.T) {
	home := t.TempDir()

	result, err := Inject(home, home, claudeAdapter())
	if err != nil {
		t.Fatalf("Inject() error = %v", err)
	}
	if !result.Changed {
		t.Fatalf("Inject() changed = false")
	}

	mcpPath := filepath.Join(home, ".claude", "mcp", "codegraph.json")
	mcpContent, err := os.ReadFile(mcpPath)
	if err != nil {
		t.Fatalf("ReadFile(codegraph.json) error = %v", err)
	}

	var parsed map[string]any
	if err := json.Unmarshal(mcpContent, &parsed); err != nil {
		t.Fatalf("Unmarshal(codegraph.json) error = %v", err)
	}
	if parsed["command"] != "codegraph" {
		t.Fatalf("codegraph.json command = %v, want codegraph", parsed["command"])
	}
	args, _ := parsed["args"].([]any)
	if len(args) != 2 || args[0] != "serve" || args[1] != "--mcp" {
		t.Fatalf("codegraph.json args = %v, want [serve --mcp]", args)
	}
}

func TestInjectClaudeWritesPromptSection(t *testing.T) {
	home := t.TempDir()

	_, err := Inject(home, home, claudeAdapter())
	if err != nil {
		t.Fatalf("Inject() error = %v", err)
	}

	claudeMDPath := filepath.Join(home, ".claude", "CLAUDE.md")
	content, err := os.ReadFile(claudeMDPath)
	if err != nil {
		t.Fatalf("ReadFile(CLAUDE.md) error = %v", err)
	}

	text := string(content)
	if !strings.Contains(text, "<!-- framework-ai:codegraph -->") {
		t.Fatal("CLAUDE.md missing open marker for codegraph")
	}
	if !strings.Contains(text, "<!-- /framework-ai:codegraph -->") {
		t.Fatal("CLAUDE.md missing close marker for codegraph")
	}
	if !strings.Contains(text, "codegraph init -i") {
		t.Fatal("CLAUDE.md missing codegraph init prompt")
	}
}

func TestInjectClaudeWithCodegraphIndexWritesFullPrompt(t *testing.T) {
	home := t.TempDir()
	if err := os.MkdirAll(filepath.Join(home, ".codegraph"), 0o755); err != nil {
		t.Fatalf("MkdirAll(.codegraph) error = %v", err)
	}

	_, err := Inject(home, home, claudeAdapter())
	if err != nil {
		t.Fatalf("Inject() error = %v", err)
	}

	claudeMDPath := filepath.Join(home, ".claude", "CLAUDE.md")
	content, err := os.ReadFile(claudeMDPath)
	if err != nil {
		t.Fatalf("ReadFile(CLAUDE.md) error = %v", err)
	}

	text := string(content)
	if !strings.Contains(text, "codegraph_search") {
		t.Fatal("CLAUDE.md missing full codegraph tools content")
	}
	// The prompt is now conditional — it always includes both branches
	// (init instructions + tool usage), so the agent can verify at runtime.
	if !strings.Contains(text, "verify whether `.codegraph/` exists") {
		t.Fatal("CLAUDE.md missing conditional check instruction")
	}
	if !strings.Contains(text, "codegraph init -i") {
		t.Fatal("CLAUDE.md missing init instructions for conditional branch")
	}
}

func TestInjectClaudeIsIdempotent(t *testing.T) {
	home := t.TempDir()

	first, err := Inject(home, home, claudeAdapter())
	if err != nil {
		t.Fatalf("Inject() first error = %v", err)
	}
	if !first.Changed {
		t.Fatalf("Inject() first changed = false")
	}

	second, err := Inject(home, home, claudeAdapter())
	if err != nil {
		t.Fatalf("Inject() second error = %v", err)
	}
	if second.Changed {
		t.Fatalf("Inject() second changed = true")
	}
}

func TestInjectOpenCodeMergesCodegraphToSettings(t *testing.T) {
	home := t.TempDir()

	result, err := Inject(home, home, opencodeAdapter())
	if err != nil {
		t.Fatalf("Inject() error = %v", err)
	}
	if !result.Changed {
		t.Fatalf("Inject() changed = false")
	}

	configPath := filepath.Join(home, ".config", "opencode", "opencode.json")
	config, err := os.ReadFile(configPath)
	if err != nil {
		t.Fatalf("ReadFile(opencode.json) error = %v", err)
	}

	text := string(config)
	if !strings.Contains(text, `"codegraph"`) {
		t.Fatal("opencode.json missing codegraph server entry")
	}
	if !strings.Contains(text, `"mcp"`) {
		t.Fatal("opencode.json missing mcp key")
	}
	if strings.Contains(text, `"mcpServers"`) {
		t.Fatal("opencode.json should use 'mcp' key, not 'mcpServers'")
	}
	if !strings.Contains(text, `"type": "local"`) {
		t.Fatal("opencode.json codegraph missing type: local")
	}
	if strings.Contains(text, `"args"`) {
		t.Fatal("opencode.json must NOT have a separate args field")
	}

	agentsPath := filepath.Join(home, ".config", "opencode", "AGENTS.md")
	agentsContent, err := os.ReadFile(agentsPath)
	if err != nil {
		t.Fatalf("ReadFile(AGENTS.md) error = %v", err)
	}
	if !strings.Contains(string(agentsContent), "<!-- framework-ai:codegraph -->") {
		t.Fatal("AGENTS.md missing codegraph section marker")
	}
}

func TestInjectOpenCodeIsIdempotent(t *testing.T) {
	home := t.TempDir()

	first, err := Inject(home, home, opencodeAdapter())
	if err != nil {
		t.Fatalf("Inject() first error = %v", err)
	}
	if !first.Changed {
		t.Fatalf("Inject() first changed = false")
	}

	second, err := Inject(home, home, opencodeAdapter())
	if err != nil {
		t.Fatalf("Inject() second error = %v", err)
	}
	if second.Changed {
		t.Fatalf("Inject() second changed = true")
	}
}

func TestInjectVSCodeMergesCodegraphToMCPConfigFile(t *testing.T) {
	home := t.TempDir()
	t.Setenv("XDG_CONFIG_HOME", filepath.Join(home, ".config"))
	adapter := vscodeAdapter()

	result, err := Inject(home, home, adapter)
	if err != nil {
		t.Fatalf("Inject(vscode) error = %v", err)
	}
	if !result.Changed {
		t.Fatalf("Inject(vscode) changed = false")
	}

	mcpPath := adapter.MCPConfigPath(home, "codegraph")
	content, err := os.ReadFile(mcpPath)
	if err != nil {
		t.Fatalf("ReadFile(mcp.json) error = %v", err)
	}

	text := string(content)
	if !strings.Contains(text, `"servers"`) {
		t.Fatal("mcp.json missing servers key")
	}
	if !strings.Contains(text, `"codegraph"`) {
		t.Fatal("mcp.json missing codegraph server")
	}
	if strings.Contains(text, `"mcpServers"`) {
		t.Fatal("mcp.json should use 'servers' key, not 'mcpServers'")
	}
}

func TestInjectGeminiMergesCodegraphToSettings(t *testing.T) {
	home := t.TempDir()

	result, err := Inject(home, home, geminiAdapter())
	if err != nil {
		t.Fatalf("Inject(gemini) error = %v", err)
	}
	if !result.Changed {
		t.Fatalf("Inject(gemini) changed = false")
	}

	settingsPath := filepath.Join(home, ".gemini", "settings.json")
	content, err := os.ReadFile(settingsPath)
	if err != nil {
		t.Fatalf("ReadFile(settings.json) error = %v", err)
	}

	text := string(content)
	if !strings.Contains(text, `"mcpServers"`) {
		t.Fatal("settings.json missing mcpServers key")
	}
	if !strings.Contains(text, `"codegraph"`) {
		t.Fatal("settings.json missing codegraph entry")
	}
	if !strings.Contains(text, `"serve"`) {
		t.Fatal("settings.json missing serve arg")
	}
}

func TestInjectCodexWritesTOMLMCP(t *testing.T) {
	home := t.TempDir()

	result, err := Inject(home, home, codexAdapter())
	if err != nil {
		t.Fatalf("Inject(codex) error = %v", err)
	}
	if !result.Changed {
		t.Fatalf("Inject(codex) changed = false")
	}

	configPath := filepath.Join(home, ".codex", "config.toml")
	content, err := os.ReadFile(configPath)
	if err != nil {
		t.Fatalf("ReadFile(config.toml) error = %v", err)
	}

	text := string(content)
	if !strings.Contains(text, "[mcp_servers.codegraph]") {
		t.Fatalf("config.toml missing [mcp_servers.codegraph] block; got:\n%s", text)
	}
	if !strings.Contains(text, `command = "codegraph"`) {
		t.Fatalf("config.toml missing command field; got:\n%s", text)
	}
	if !strings.Contains(text, `"serve"`) {
		t.Fatalf("config.toml missing serve arg; got:\n%s", text)
	}
}

func TestInjectCodexIsIdempotent(t *testing.T) {
	home := t.TempDir()

	first, err := Inject(home, home, codexAdapter())
	if err != nil {
		t.Fatalf("Inject(codex) first error = %v", err)
	}
	if !first.Changed {
		t.Fatalf("Inject(codex) first changed = false")
	}

	second, err := Inject(home, home, codexAdapter())
	if err != nil {
		t.Fatalf("Inject(codex) second error = %v", err)
	}
	if second.Changed {
		t.Fatalf("Inject(codex) second changed = true (should be idempotent)")
	}

	configPath := filepath.Join(home, ".codex", "config.toml")
	content, err := os.ReadFile(configPath)
	if err != nil {
		t.Fatalf("ReadFile(config.toml) error = %v", err)
	}
	count := strings.Count(string(content), "[mcp_servers.codegraph]")
	if count != 1 {
		t.Fatalf("config.toml has %d [mcp_servers.codegraph] blocks, want exactly 1; got:\n%s", count, string(content))
	}
}

func TestHasCodegraphIndex(t *testing.T) {
	t.Run("returns true when .codegraph exists", func(t *testing.T) {
		dir := t.TempDir()
		if err := os.MkdirAll(filepath.Join(dir, ".codegraph"), 0o755); err != nil {
			t.Fatalf("MkdirAll error = %v", err)
		}
		if !hasCodegraphIndex(dir) {
			t.Fatal("hasCodegraphIndex() = false, want true")
		}
	})

	t.Run("returns false when .codegraph missing", func(t *testing.T) {
		dir := t.TempDir()
		if hasCodegraphIndex(dir) {
			t.Fatal("hasCodegraphIndex() = true, want false")
		}
	})

	t.Run("returns false for empty workspace", func(t *testing.T) {
		if hasCodegraphIndex("") {
			t.Fatal("hasCodegraphIndex(\"\") = true, want false")
		}
	})
}

func TestSelectPromptContent(t *testing.T) {
	t.Run("returns full prompt when index exists", func(t *testing.T) {
		dir := t.TempDir()
		if err := os.MkdirAll(filepath.Join(dir, ".codegraph"), 0o755); err != nil {
			t.Fatalf("MkdirAll error = %v", err)
		}
		content := selectPromptContent(dir)
		if !strings.Contains(content, "codegraph_search") {
			t.Fatal("selectPromptContent() missing full tools content")
		}
	})

	t.Run("returns init prompt when index missing", func(t *testing.T) {
		dir := t.TempDir()
		content := selectPromptContent(dir)
		if !strings.Contains(content, "codegraph init -i") {
			t.Fatal("selectPromptContent() missing init prompt")
		}
	})
}


