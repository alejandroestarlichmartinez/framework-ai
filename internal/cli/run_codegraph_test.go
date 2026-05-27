package cli

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/alejandroestarlichmartinez/framework-ai/internal/model"
	"github.com/alejandroestarlichmartinez/framework-ai/internal/system"
)

func TestRunInstallCodeGraphSkipsInstallWhenOnPath(t *testing.T) {
	origLookPath := cmdLookPath
	cmdLookPath = func(file string) (string, error) {
		if file == "codegraph" {
			return "/usr/local/bin/codegraph", nil
		}
		return origLookPath(file)
	}
	t.Cleanup(func() { cmdLookPath = origLookPath })

	origHomeDir := osUserHomeDir
	homeDir := t.TempDir()
	osUserHomeDir = func() (string, error) { return homeDir, nil }
	t.Cleanup(func() { osUserHomeDir = origHomeDir })

	detection := system.DetectionResult{
		Configs: []system.ConfigState{{Agent: "claude-code", Exists: true}},
	}

	result, err := RunInstall([]string{"--agent", "claude-code", "--component", "codegraph"}, detection)
	if err != nil {
		t.Fatalf("RunInstall() error = %v", err)
	}

	if len(result.Plan.Apply) == 0 {
		t.Fatal("expected at least one apply step")
	}

	// Verify CodeGraph MCP config was written.
	mcpPath := filepath.Join(homeDir, ".claude", "mcp", "codegraph.json")
	if _, err := os.Stat(mcpPath); os.IsNotExist(err) {
		t.Fatalf("codegraph.json should exist at %q", mcpPath)
	}
}

func TestRunInstallCodeGraphInstallsWhenMissing(t *testing.T) {
	origLookPath := cmdLookPath
	cmdLookPath = func(file string) (string, error) {
		if file == "codegraph" {
			return "", os.ErrNotExist
		}
		return origLookPath(file)
	}
	t.Cleanup(func() { cmdLookPath = origLookPath })

	var ranCommands [][]string
	origRunCommand := runCommand
	runCommand = func(name string, args ...string) error {
		ranCommands = append(ranCommands, append([]string{name}, args...))
		return nil
	}
	t.Cleanup(func() { runCommand = origRunCommand })

	detection := system.DetectionResult{
		Configs: []system.ConfigState{{Agent: "claude-code", Exists: true}},
	}

	_, err := RunInstall([]string{"--agent", "claude-code", "--component", "codegraph", "--dry-run"}, detection)
	if err != nil {
		t.Fatalf("RunInstall() error = %v", err)
	}

	var foundInstall bool
	for _, cmd := range ranCommands {
		if len(cmd) > 2 && cmd[0] == "sh" && cmd[1] == "-c" && strings.Contains(cmd[2], "codegraph") {
			foundInstall = true
			break
		}
	}
	if foundInstall {
		t.Fatal("dry-run should not execute install commands")
	}
}

func TestComponentPathsWithWorkspaceIncludesCodeGraphPaths(t *testing.T) {
	homeDir := t.TempDir()
	workspaceDir := t.TempDir()

	adapters := resolveAdapters([]model.AgentID{model.AgentClaudeCode})
	paths := componentPathsWithWorkspace(homeDir, workspaceDir, model.Selection{}, adapters, model.ComponentCodeGraph)

	var foundMCP, foundPrompt bool
	for _, p := range paths {
		if strings.Contains(p, "codegraph.json") {
			foundMCP = true
		}
		if strings.Contains(p, "CLAUDE.md") {
			foundPrompt = true
		}
	}
	if !foundMCP {
		t.Fatal("componentPathsWithWorkspace should include codegraph MCP path")
	}
	if !foundPrompt {
		t.Fatal("componentPathsWithWorkspace should include codegraph prompt path")
	}
}
