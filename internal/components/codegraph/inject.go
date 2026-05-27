package codegraph

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/alejandroestarlichmartinez/framework-ai/internal/agents"
	"github.com/alejandroestarlichmartinez/framework-ai/internal/assets"
	"github.com/alejandroestarlichmartinez/framework-ai/internal/components/filemerge"
	"github.com/alejandroestarlichmartinez/framework-ai/internal/model"
)

type InjectionResult struct {
	Changed bool
	Files   []string
}

var defaultCodegraphServerJSON = []byte("{\n  \"command\": \"codegraph\",\n  \"args\": [\"serve\", \"--mcp\"]\n}\n")

var defaultCodegraphOverlayJSON = []byte("{\n  \"mcpServers\": {\n    \"codegraph\": {\n      \"command\": \"codegraph\",\n      \"args\": [\"serve\", \"--mcp\"]\n    }\n  }\n}\n")

var openCodeCodegraphOverlayJSON = []byte("{\n  \"mcp\": {\n    \"codegraph\": {\n      \"__replace__\": {\n        \"command\": [\"codegraph\", \"serve\", \"--mcp\"],\n        \"type\": \"local\"\n      }\n    }\n  }\n}\n")

var openClawCodegraphOverlayJSON = []byte("{\n  \"mcp\": {\n    \"servers\": {\n      \"codegraph\": {\n        \"command\": \"codegraph\",\n        \"args\": [\"serve\", \"--mcp\"]\n      }\n    }\n  }\n}\n")

var vsCodeCodegraphOverlayJSON = []byte("{\n  \"servers\": {\n    \"codegraph\": {\n      \"command\": \"codegraph\",\n      \"args\": [\"serve\", \"--mcp\"]\n    }\n  }\n}\n")

func DefaultCodegraphServerJSON() []byte {
	content := make([]byte, len(defaultCodegraphServerJSON))
	copy(content, defaultCodegraphServerJSON)
	return content
}

func DefaultCodegraphOverlayJSON() []byte {
	content := make([]byte, len(defaultCodegraphOverlayJSON))
	copy(content, defaultCodegraphOverlayJSON)
	return content
}

func OpenCodeCodegraphOverlayJSON() []byte {
	content := make([]byte, len(openCodeCodegraphOverlayJSON))
	copy(content, openCodeCodegraphOverlayJSON)
	return content
}

func OpenClawCodegraphOverlayJSON() []byte {
	content := make([]byte, len(openClawCodegraphOverlayJSON))
	copy(content, openClawCodegraphOverlayJSON)
	return content
}

func VSCodeCodegraphOverlayJSON() []byte {
	content := make([]byte, len(vsCodeCodegraphOverlayJSON))
	copy(content, vsCodeCodegraphOverlayJSON)
	return content
}

func hasCodegraphIndex(workspaceDir string) bool {
	if workspaceDir == "" {
		return false
	}
	_, err := os.Stat(filepath.Join(workspaceDir, ".codegraph"))
	return err == nil
}

func selectPromptContent(workspaceDir string) string {
	if hasCodegraphIndex(workspaceDir) {
		return assets.MustRead("claude/codegraph.md")
	}
	return assets.MustRead("claude/codegraph-init.md")
}

func Inject(homeDir string, workspaceDir string, adapter agents.Adapter) (InjectionResult, error) {
	if !adapter.SupportsMCP() {
		// Still inject system prompt even if MCP is not supported.
		return injectSystemPromptOnly(homeDir, workspaceDir, adapter)
	}

	files := make([]string, 0, 2)
	changed := false

	// 1. Write MCP server config using the adapter's strategy.
	switch adapter.MCPStrategy() {
	case model.StrategySeparateMCPFiles:
		mcpPath := adapter.MCPConfigPath(homeDir, "codegraph")
		mcpWrite, err := filemerge.WriteFileAtomic(mcpPath, DefaultCodegraphServerJSON(), 0o644)
		if err != nil {
			return InjectionResult{}, err
		}
		changed = changed || mcpWrite.Changed
		files = append(files, mcpPath)

	case model.StrategyMergeIntoSettings:
		settingsPath := adapter.SettingsPath(homeDir)
		if settingsPath == "" {
			break
		}
		var overlay []byte
		if adapter.Agent() == model.AgentOpenCode || adapter.Agent() == model.AgentKilocode {
			overlay = OpenCodeCodegraphOverlayJSON()
		} else if adapter.Agent() == model.AgentOpenClaw {
			overlay = OpenClawCodegraphOverlayJSON()
		} else {
			overlay = DefaultCodegraphOverlayJSON()
		}
		settingsWrite, err := mergeJSONFile(settingsPath, overlay)
		if err != nil {
			return InjectionResult{}, err
		}
		changed = changed || settingsWrite.Changed
		files = append(files, settingsPath)

	case model.StrategyMCPConfigFile:
		mcpPath := adapter.MCPConfigPath(homeDir, "codegraph")
		if mcpPath == "" {
			break
		}
		var overlay []byte
		if adapter.Agent() == model.AgentVSCodeCopilot {
			overlay = VSCodeCodegraphOverlayJSON()
		} else {
			overlay = DefaultCodegraphOverlayJSON()
		}
		mcpWrite, err := mergeJSONFile(mcpPath, overlay)
		if err != nil {
			return InjectionResult{}, err
		}
		changed = changed || mcpWrite.Changed
		files = append(files, mcpPath)

	case model.StrategyTOMLFile:
		configPath := adapter.MCPConfigPath(homeDir, "codegraph")
		if configPath == "" {
			break
		}
		existing, err := readFileOrEmpty(configPath)
		if err != nil {
			return InjectionResult{}, err
		}
		updated := upsertCodexCodegraphBlock(existing)
		tomlWrite, err := filemerge.WriteFileAtomic(configPath, []byte(updated), 0o644)
		if err != nil {
			return InjectionResult{}, err
		}
		changed = changed || tomlWrite.Changed
		files = append(files, configPath)
	}

	// 2. Inject CodeGraph instructions into system prompt (if supported).
	if adapter.SupportsSystemPrompt() {
		promptPath := adapter.SystemPromptFile(homeDir)
		content := selectPromptContent(workspaceDir)

		existing, err := readFileOrEmpty(promptPath)
		if err != nil {
			return InjectionResult{}, err
		}

		updated := filemerge.InjectMarkdownSection(existing, "codegraph", content)

		mdWrite, err := filemerge.WriteFileAtomic(promptPath, []byte(updated), 0o644)
		if err != nil {
			return InjectionResult{}, err
		}
		changed = changed || mdWrite.Changed
		files = append(files, promptPath)
	}

	return InjectionResult{Changed: changed, Files: files}, nil
}

func injectSystemPromptOnly(homeDir string, workspaceDir string, adapter agents.Adapter) (InjectionResult, error) {
	if !adapter.SupportsSystemPrompt() {
		return InjectionResult{}, nil
	}

	promptPath := adapter.SystemPromptFile(homeDir)
	content := selectPromptContent(workspaceDir)

	existing, err := readFileOrEmpty(promptPath)
	if err != nil {
		return InjectionResult{}, err
	}

	updated := filemerge.InjectMarkdownSection(existing, "codegraph", content)

	mdWrite, err := filemerge.WriteFileAtomic(promptPath, []byte(updated), 0o644)
	if err != nil {
		return InjectionResult{}, err
	}

	return InjectionResult{Changed: mdWrite.Changed, Files: []string{promptPath}}, nil
}

func upsertCodexCodegraphBlock(content string) string {
	block := "[mcp_servers.codegraph]\ncommand = \"codegraph\"\nargs = [\"serve\", \"--mcp\"]"
	content = strings.ReplaceAll(content, "\r\n", "\n")
	lines := strings.Split(content, "\n")

	var kept []string
	for i := 0; i < len(lines); {
		trimmed := strings.TrimSpace(lines[i])
		if trimmed == "[mcp_servers.codegraph]" {
			i++
			for i < len(lines) {
				next := strings.TrimSpace(lines[i])
				if strings.HasPrefix(next, "[") && strings.HasSuffix(next, "]") {
					break
				}
				i++
			}
			continue
		}
		kept = append(kept, lines[i])
		i++
	}

	base := strings.TrimSpace(strings.Join(kept, "\n"))
	if base == "" {
		return block + "\n"
	}

	return base + "\n\n" + block + "\n"
}

func mergeJSONFile(path string, overlay []byte) (filemerge.WriteResult, error) {
	baseJSON, err := osReadFile(path)
	if err != nil {
		return filemerge.WriteResult{}, err
	}

	merged, err := filemerge.MergeJSONObjects(baseJSON, overlay)
	if err != nil {
		return filemerge.WriteResult{}, err
	}

	return filemerge.WriteFileAtomic(path, merged, 0o644)
}

var osReadFile = func(path string) ([]byte, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, fmt.Errorf("read json file %q: %w", path, err)
	}
	return content, nil
}

func readFileOrEmpty(path string) (string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return "", nil
		}
		return "", fmt.Errorf("read file %q: %w", path, err)
	}
	return string(data), nil
}
