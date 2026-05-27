package opencodeplugin

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/alejandroestarlichmartinez/framework-ai/internal/components/filemerge"
	"github.com/alejandroestarlichmartinez/framework-ai/internal/model"
)

type Definition struct {
	ID          model.OpenCodeCommunityPluginID
	Name        string
	PackageName string
	RepoURL     string
	Owner       string
	Repo        string
	Description string
}

type Result struct {
	Changed bool
	Files   []string
}

var definitions = []Definition{
	{
		ID:          model.OpenCodePluginSubAgentStatusline,
		Name:        "Sub-agent Statusline",
		PackageName: "opencode-subagent-statusline",
		RepoURL:     "https://github.com/Joaquinvesapa/sub-agent-statusline",
		Owner:       "Joaquinvesapa",
		Repo:        "sub-agent-statusline",
		Description: "OpenCode sidebar/statusline for sub-agent activity",
	},
	{
		ID:          model.OpenCodePluginSDDEngramManage,
		Name:        "SDD Engram Manager",
		PackageName: "opencode-sdd-engram-manage",
		RepoURL:     "https://github.com/j0k3r-dev-rgl/sdd-engram-plugin",
		Owner:       "j0k3r-dev-rgl",
		Repo:        "sdd-engram-plugin",
		Description: "OpenCode TUI for SDD profiles and Engram memories",
	},
}

const gentleLogoPluginFile = "gentle-logo.tsx"

const gentleLogoPluginSource = `// @ts-nocheck
/** @jsxImportSource @opentui/solid */
import type { TuiPlugin, TuiThemeCurrent } from "@opencode-ai/plugin/tui"
import { useTerminalDimensions } from "@opentui/solid"
import { createMemo } from "solid-js"

const id = "gentle-logo"

const roseArt = [
  "⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⡀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀",
  "⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⣿⣆⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⣠⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀",
  "⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⣿⣿⣦⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⢀⣼⣿⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀",
  "⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⣿⣿⣿⣷⡄⠀⠀⠀⠀⠀⠀⠀⠀⢀⣴⣿⣿⣿⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀",
  "⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⣿⣿⣿⣿⣿⣆⠀⠀⠀⠀⠀⢀⣴⣿⣿⣿⣿⣿⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀",
  "⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⢠⣿⣿⣿⣿⣿⣿⣧⠀⠀⠀⣴⣿⣿⣿⣿⣿⣿⡇⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀",
  "⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⢸⣿⣿⣿⣿⣿⣿⣿⣧⣀⣾⣿⣿⣿⣿⣿⣿⣿⠁⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀",
  "⠀⠀⠀⠀⢹⣿⣶⣦⣤⣀⡀⠀⠀⠀⠀⠀⣼⣿⣿⣿⡿⠿⠟⠛⠛⠿⠿⣿⣿⣿⣿⣿⣿⡟⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀",
  "⠀⠀⠀⠀⠀⢿⣿⣿⣿⣿⣿⣿⣶⣶⣤⣤⡿⠟⠉⢴⣶⣿⣿⣿⣿⣿⣷⣦⣍⠻⣿⣿⣿⡇⠀⠀⠀⠀⠀⣀⣀⣠⣤⣶⡶",
  "⠀⠀⠀⠀⠀⠈⣿⣿⣿⣿⣿⣿⣿⡿⠟⣋⣀⣙⡻⢶⣝⢿⣿⣿⣿⣿⣿⣿⣿⣿⣌⠻⣿⣷⣶⣶⣿⣿⣿⣿⣿⣿⣿⠏⠀",
  "⠀⠀⠀⠀⠀⠀⠘⣿⣿⣿⣿⣿⠏⣴⣿⡿⠿⢿⣿⣦⡙⢦⣽⣿⣿⣿⣿⣿⣿⣿⣿⡧⠹⣿⣿⣿⣿⣿⣿⣿⣿⡿⠁⠀⠀",
  "⠀⠀⠀⠀⠀⠀⠀⢹⣿⣿⣿⣿⡆⢉⣥⣶⣾⣶⣌⠻⣿⣎⠻⣿⣿⣿⡿⠟⣋⣭⣴⣶⡄⢹⣿⣿⣿⣿⣿⣿⡟⠀⠀⠀⠀",
  "⠀⠀⠀⠀⠀⠀⠀⠀⢿⣿⣿⣿⢃⣿⣿⡿⠿⠿⠿⣧⡙⢿⣷⣶⣶⣶⣶⣿⠿⠟⠋⣩⣴⡌⣿⣿⣿⣿⣿⡟⠀⠀⠀⠀⠀",
  "⠀⠀⠀⠀⠀⠀⠀⢀⣸⣿⣿⡟⢸⠟⣡⣶⣾⣿⣿⣶⣌⠲⣬⣉⠉⣉⣥⣴⣾⣿⣷⣦⡙⣧⢹⣿⣿⣿⠟⠀⠀⠀⠀⠀⠀",
  "⠀⠀⠀⢀⣠⣴⣾⣿⣿⣿⣿⡇⡎⣼⣿⣿⣿⣿⣿⣿⠉⢢⢹⡿⢰⣿⣿⣿⣿⣿⣿⠉⣳⠈⢸⣿⣿⡋⠀⠀⠀⠀⠀⠀⠀",
  "⠠⣴⣾⣿⣿⣿⣿⣿⣿⣿⣿⢁⡇⣿⣿⣿⣿⣿⣿⣿⣿⣿⢸⡇⣾⣿⣿⣿⣿⣿⣿⣿⣿⠀⢸⣿⣿⣿⣷⣶⣤⣄⣀⣀⠀",
  "⠀⠀⠉⠻⢿⣿⣿⣿⣿⣿⣿⢸⡇⢿⣿⣿⣿⣿⣿⣿⣿⠇⣼⣧⠸⣿⣿⣿⣿⣿⣿⣿⡿⢠⢸⣿⣿⣿⣿⣿⣿⣿⣿⠟⠁",
  "⠀⠀⠀⠀⠀⠈⠛⢿⣿⣿⣿⢸⣿⣌⠻⢿⣿⣿⣿⡿⢋⣼⣿⣿⣧⡙⠿⣿⣿⣿⡿⠟⣡⣿⢸⣿⣿⣿⣿⣿⡿⠋⠀⠀⠀",
  "⠀⠀⠀⠀⠀⠀⣠⣾⣿⣿⣿⣾⣿⣿⣿⣶⣤⣤⣤⣶⣿⠋⣿⣿⢻⣿⣷⣶⣤⣴⣶⣿⣿⣿⢸⣿⣿⣿⡿⠋⠀⠀⠀⠀⠀",
  "⠀⠀⠀⠀⣠⣾⣿⣿⣿⣿⣿⢹⣿⣷⣬⣛⣛⠛⣛⣩⣽⠀⣿⣿⢀⣷⣬⣙⡛⠛⣛⣫⣴⣿⢸⣿⣿⡟⠀⠀⠀⠀⠀⠀⠀",
  "⠀⠀⠠⣾⣿⣿⣿⣿⣿⣿⠟⢸⣿⣿⣿⣿⣿⣿⣿⣿⣿⡇⣿⣿⢸⣿⣿⣿⣿⣿⣿⣿⣿⣿⠘⢿⣿⣷⡀⠀⠀⠀⠀⠀⠀",
  "⠀⠀⠀⠈⠙⠻⢿⣿⣿⢃⣾⠘⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣌⣡⣾⣿⣿⣿⣿⣿⣿⣿⣿⣿⠀⣦⡙⣿⣿⣷⣤⣀⠀⠀⠀",
  "⠀⠀⠀⠀⠀⠀⠀⠈⣿⢸⣿⡇⣿⣿⣿⣿⡿⢿⣿⡿⢻⣟⢹⡟⢻⣟⠻⣿⣿⣿⣿⣿⣿⣿⢸⣿⡇⣿⣿⣿⠿⠟⠁⠀⠀",
  "⠀⠀⠀⠀⠀⠀⠀⢰⣿⣦⡙⠇⢸⣿⣿⡟⡰⠁⠈⠁⠀⠁⠀⠀⠀⠁⠀⠉⠀⠙⣌⢻⣿⣿⠘⣋⣴⠉⠁⠀⠀⠀⠀⠀⠀",
  "⠀⠀⠀⠀⠀⠀⢀⣿⣿⣿⣿⣷⡌⣿⣿⢰⡇⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠘⡌⣿⡇⣼⣿⣿⡆⠀⠀⠀⠀⠀⠀⠀",
  "⠀⠀⠀⠀⠀⠀⠘⠛⠛⠛⠛⠻⣷⠹⣿⠸⣧⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⢠⡇⣿⢡⣿⣿⣿⣷⠀⠀⠀⠀⠀⠀⠀",
  "⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⢀⣿⣇⢻⣧⡙⠿⠶⠴⢦⡀⠶⣶⣶⡶⠆⢠⣤⠴⢏⣴⢃⡎⠀⠈⠉⠉⠀⠀⠀⠀⠀⠀⠀",
  "⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠘⠟⠛⠀⠻⣿⣷⣶⣾⣿⣿⡇⢹⠏⣴⣶⣶⣶⣶⡿⠃⠚⠇⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀",
  "⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠈⠻⣿⣿⣿⣿⣿⣤⣾⣿⣿⣿⣿⠟⠁⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀",
  "⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠈⠙⠛⠛⠻⠿⠿⠛⠛⠉⠁⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀",
]

const compactArt = ["✦ Framework AI ✦"]

const Logo = (props: { theme: TuiThemeCurrent }) => {
  const dim = useTerminalDimensions()
  const lines = createMemo(() => {
    const term = dim()
    return term.height >= roseArt.length + 6 && term.width >= 64 ? roseArt : compactArt
  })

  return (
    <box flexDirection="column" alignItems="center">
      {lines().map((line) => (
        <text fg={props.theme.accent}>{line}</text>
      ))}
    </box>
  )
}

const tui: TuiPlugin = async (api) => {
  api.slots.register({
    id,
    order: 100,
    slots: {
      home_logo(ctx) {
        return <Logo theme={ctx.theme.current} />
      },
    },
  })
}

const plugin = { id: "gentle-logo", tui }
export default plugin
`

func Definitions() []Definition {
	out := make([]Definition, len(definitions))
	copy(out, definitions)
	return out
}

func DefinitionFor(id model.OpenCodeCommunityPluginID) (Definition, bool) {
	for _, def := range definitions {
		if def.ID == id {
			return def, true
		}
	}
	return Definition{}, false
}

func Install(homeDir string, id model.OpenCodeCommunityPluginID) (Result, error) {
	if id == model.OpenCodePluginGentleLogo {
		return installGentleLogo(homeDir)
	}

	def, ok := DefinitionFor(id)
	if !ok {
		return Result{}, fmt.Errorf("unknown OpenCode community plugin %q", id)
	}

	opencodeDir := filepath.Join(homeDir, ".config", "opencode")
	if err := os.MkdirAll(opencodeDir, 0o755); err != nil {
		return Result{}, fmt.Errorf("create OpenCode config dir: %w", err)
	}

	tuiPath := filepath.Join(opencodeDir, "tui.json")
	written, err := ensureTUIPlugin(tuiPath, def.PackageName)
	if err != nil {
		return Result{}, err
	}

	return Result{Changed: written, Files: []string{tuiPath}}, nil
}

func installGentleLogo(homeDir string) (Result, error) {
	opencodeDir := filepath.Join(homeDir, ".config", "opencode")
	pluginDir := filepath.Join(opencodeDir, "tui-plugins")
	pluginPath := filepath.Join(pluginDir, gentleLogoPluginFile)
	tuiPath := filepath.Join(opencodeDir, "tui.json")

	// Remove competing logo plugins so ours takes precedence.
	removedPaths, err := removeCompetingLogos(tuiPath, pluginDir, pluginPath)
	if err != nil {
		// Non-fatal: log and continue.
		fmt.Fprintf(os.Stderr, "warning: could not clean competing logos: %v\n", err)
	}

	pluginWrite, err := filemerge.WriteFileAtomic(pluginPath, []byte(gentleLogoPluginSource), 0o644)
	if err != nil {
		return Result{}, fmt.Errorf("write Gentle Logo TUI plugin: %w", err)
	}
	tuiChanged, err := ensureTUIPlugin(tuiPath, pluginPath)
	if err != nil {
		return Result{}, err
	}

	files := []string{pluginPath, tuiPath}
	files = append(files, removedPaths...)
	return Result{
		Changed: pluginWrite.Changed || tuiChanged || len(removedPaths) > 0,
		Files:   files,
	}, nil
}

// removeCompetingLogos strips other logo plugins from tui.json and deletes
// their files so the Gentle Logo is the only home_logo renderer.
func removeCompetingLogos(tuiPath, pluginDir, ourPluginPath string) ([]string, error) {
	data, err := os.ReadFile(tuiPath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}

	var root map[string]any
	if err := json.Unmarshal(data, &root); err != nil {
		return nil, fmt.Errorf("parse tui.json: %w", err)
	}

	plugins := stringSlice(root["plugin"])
	if len(plugins) == 0 {
		return nil, nil
	}

	var removed []string
	var kept []string
	for _, p := range plugins {
		if isCompetingLogo(p, pluginDir, ourPluginPath) {
			removed = append(removed, p)
			// Delete the file if it's inside our pluginDir.
			if filepath.IsAbs(p) && strings.HasPrefix(p, pluginDir) {
				_ = os.Remove(p)
			}
		} else {
			kept = append(kept, p)
		}
	}

	if len(removed) == 0 {
		return nil, nil
	}

	root["plugin"] = kept
	out, err := json.MarshalIndent(root, "", "  ")
	if err != nil {
		return nil, err
	}
	out = append(out, '\n')
	if _, err := filemerge.WriteFileAtomic(tuiPath, out, 0o644); err != nil {
		return nil, err
	}

	return removed, nil
}

// isCompetingLogo identifies other logo plugins in the tui-plugins directory.
func isCompetingLogo(pluginPath, pluginDir, ourPluginPath string) bool {
	if pluginPath == ourPluginPath {
		return false
	}
	// Only consider local .tsx files inside the tui-plugins directory.
	if !filepath.IsAbs(pluginPath) {
		return false
	}
	if !strings.HasPrefix(pluginPath, pluginDir) {
		return false
	}
	if !strings.HasSuffix(pluginPath, ".tsx") {
		return false
	}
	// Any other .tsx in tui-plugins is treated as a competing custom plugin.
	return true
}

func ensureTUIPlugin(path, pkg string) (bool, error) {
	root := map[string]any{"$schema": "https://opencode.ai/tui.json"}
	if data, err := os.ReadFile(path); err == nil && len(bytes.TrimSpace(data)) > 0 {
		if err := json.Unmarshal(data, &root); err != nil {
			return false, fmt.Errorf("parse OpenCode TUI config %q: %w", path, err)
		}
	} else if err != nil && !os.IsNotExist(err) {
		return false, fmt.Errorf("read OpenCode TUI config %q: %w", path, err)
	}

	plugins := stringSlice(root["plugin"])
	for _, existing := range plugins {
		if existing == pkg {
			return false, nil
		}
	}
	plugins = append(plugins, pkg)
	root["plugin"] = plugins

	out, err := json.MarshalIndent(root, "", "  ")
	if err != nil {
		return false, err
	}
	out = append(out, '\n')
	wr, err := filemerge.WriteFileAtomic(path, out, 0o644)
	if err != nil {
		return false, err
	}
	return wr.Changed, nil
}

func stringSlice(value any) []string {
	items, ok := value.([]any)
	if !ok {
		return nil
	}
	out := make([]string, 0, len(items))
	for _, item := range items {
		if s, ok := item.(string); ok && strings.TrimSpace(s) != "" {
			out = append(out, s)
		}
	}
	return out
}
