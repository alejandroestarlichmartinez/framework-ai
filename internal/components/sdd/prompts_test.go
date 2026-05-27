package sdd

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// TestSharedPromptDir verifies the expected directory path is returned.
func TestSharedPromptDir(t *testing.T) {
	want := filepath.FromSlash("/home/testuser/.config/opencode/prompts/sdd")
	got := SharedPromptDir(filepath.FromSlash("/home/testuser"))
	if got != want {
		t.Fatalf("SharedPromptDir(%q) = %q, want %q", "/home/testuser", got, want)
	}
}

// TestWriteSharedPromptFilesCreates12Files verifies that WriteSharedPromptFiles
// creates exactly the 12 expected prompt files under {homeDir}/.config/opencode/prompts/sdd/.
func TestWriteSharedPromptFilesCreates12Files(t *testing.T) {
	home := t.TempDir()

	changed, err := WriteSharedPromptFiles(home)
	if err != nil {
		t.Fatalf("WriteSharedPromptFiles() error = %v", err)
	}
	if !changed {
		t.Fatal("WriteSharedPromptFiles() first call changed = false, want true")
	}

	expectedFiles := []string{
		"sdd-init.md",
		"sdd-explore.md",
		"sdd-propose.md",
		"sdd-spec.md",
		"sdd-design.md",
		"sdd-tasks.md",
		"sdd-apply.md",
		"sdd-verify.md",
		"sdd-archive.md",
		"sdd-onboard.md",
		"sdd-inspect.md",
		"sdd-optimize.md",
	}

	promptDir := SharedPromptDir(home)
	for _, fileName := range expectedFiles {
		path := filepath.Join(promptDir, fileName)
		info, statErr := os.Stat(path)
		if statErr != nil {
			t.Errorf("prompt file %q not found: %v", path, statErr)
			continue
		}
		if info.Size() == 0 {
			t.Errorf("prompt file %q is empty", path)
		}
	}
}

// TestWriteSharedPromptFilesIdempotent verifies that calling WriteSharedPromptFiles
// twice returns changed=false on the second call.
func TestWriteSharedPromptFilesIdempotent(t *testing.T) {
	home := t.TempDir()

	first, err := WriteSharedPromptFiles(home)
	if err != nil {
		t.Fatalf("WriteSharedPromptFiles() first error = %v", err)
	}
	if !first {
		t.Fatal("WriteSharedPromptFiles() first call changed = false, want true")
	}

	second, err := WriteSharedPromptFiles(home)
	if err != nil {
		t.Fatalf("WriteSharedPromptFiles() second error = %v", err)
	}
	if second {
		t.Fatal("WriteSharedPromptFiles() second call changed = true, want false (idempotent)")
	}
}

// TestWriteSharedPromptFilesContent verifies each prompt file contains the
// executor-scoped sub-agent prompt content for the correct phase.
func TestWriteSharedPromptFilesContent(t *testing.T) {
	home := t.TempDir()

	if _, err := WriteSharedPromptFiles(home); err != nil {
		t.Fatalf("WriteSharedPromptFiles() error = %v", err)
	}

	promptDir := SharedPromptDir(home)

	phases := []struct {
		file  string
		phase string
	}{
		{"sdd-init.md", "init"},
		{"sdd-explore.md", "explore"},
		{"sdd-propose.md", "propose"},
		{"sdd-spec.md", "spec"},
		{"sdd-design.md", "design"},
		{"sdd-tasks.md", "tasks"},
		{"sdd-apply.md", "apply"},
		{"sdd-verify.md", "verify"},
		{"sdd-archive.md", "archive"},
		{"sdd-onboard.md", "onboard"},
		{"sdd-inspect.md", "inspect"},
		{"sdd-optimize.md", "optimize"},
	}

	for _, tc := range phases {
		path := filepath.Join(promptDir, tc.file)
		data, readErr := os.ReadFile(path)
		if readErr != nil {
			t.Errorf("ReadFile(%q) error = %v", path, readErr)
			continue
		}

		content := string(data)

		// Each file must contain the phase name (executor-scoped prompt).
		if !strings.Contains(content, tc.phase) {
			t.Errorf("prompt file %q missing phase %q in content", tc.file, tc.phase)
		}

		// Each file must contain the key executor-scoped markers.
		for _, marker := range []string{"not the orchestrator", "Do NOT delegate", "Do NOT launch sub-agents"} {
			if !strings.Contains(content, marker) {
				t.Errorf("prompt file %q missing required marker %q", tc.file, marker)
			}
		}
	}
}

// TestInjectOpenCodeMultiModeSubagentPromptsUseFilePaths verifies that after
// injection in multi-mode, each sub-agent's prompt field in opencode.json
// contains a {file:...} reference (not an inline string).
func TestInjectOpenCodeMultiModeSubagentPromptsUseFilePaths(t *testing.T) {
	home := t.TempDir()
	mockNoPackageManager(t)

	if _, err := Inject(home, opencodeAdapter(), "multi"); err != nil {
		t.Fatalf("Inject(multi) error = %v", err)
	}

	settingsPath := filepath.Join(home, ".config", "opencode", "opencode.json")
	content, err := os.ReadFile(settingsPath)
	if err != nil {
		t.Fatalf("ReadFile(opencode.json) error = %v", err)
	}

	promptDir := SharedPromptDir(home)

	text := strings.ReplaceAll(string(content), `\\`, `/`)
	for _, phase := range []string{"sdd-init", "sdd-explore", "sdd-propose", "sdd-spec", "sdd-design", "sdd-tasks", "sdd-apply", "sdd-verify", "sdd-archive", "sdd-onboard", "sdd-inspect", "sdd-optimize"} {
		expectedRef := "{file:" + filepath.Join(promptDir, phase+".md") + "}"
		expectedRef = strings.ReplaceAll(expectedRef, `\`, `/`)
		if !strings.Contains(text, expectedRef) {
			t.Errorf("opencode.json sub-agent %q missing {file:...} reference %q", phase, expectedRef)
		}
	}
}

// TestInjectOpenCodeMultiModeOrchestratorPromptUsesFileRef verifies that
// the orchestrator prompt uses a {file:...} reference (not inline text) after injection.
func TestInjectOpenCodeMultiModeOrchestratorPromptUsesFileRef(t *testing.T) {
	home := t.TempDir()
	mockNoPackageManager(t)

	if _, err := Inject(home, opencodeAdapter(), "multi"); err != nil {
		t.Fatalf("Inject(multi) error = %v", err)
	}

	settingsPath := filepath.Join(home, ".config", "opencode", "opencode.json")
	content, err := os.ReadFile(settingsPath)
	if err != nil {
		t.Fatalf("ReadFile(opencode.json) error = %v", err)
	}

	text := string(content)

	// The orchestrator should reference its prompt via {file:...} pointing to agents/.
	orchestratorFile := filepath.Join(OrchestratorPromptDir(home), "framework-orchestrator.md")
	expectedRef := "{file:" + orchestratorFile + "}"
	if !strings.Contains(text, expectedRef) {
		t.Fatalf("orchestrator prompt should be a file reference %q in opencode.json", expectedRef)
	}

	// Verify the referenced file actually exists and contains the orchestrator content.
	fileContent, err := os.ReadFile(orchestratorFile)
	if err != nil {
		t.Fatalf("ReadFile(%q) error = %v", orchestratorFile, err)
	}
	if !strings.Contains(string(fileContent), "Gentle AI") {
		t.Fatal("orchestrator prompt file missing expected content")
	}
}

// TestInjectOpenCodeMultiModeIdempotentWithPromptFiles verifies that the second
// Inject call returns changed=false when prompt files are already on disk.
func TestInjectOpenCodeMultiModeIdempotentWithPromptFiles(t *testing.T) {
	home := t.TempDir()
	mockNoPackageManager(t)

	first, err := Inject(home, opencodeAdapter(), "multi")
	if err != nil {
		t.Fatalf("Inject(multi) first error = %v", err)
	}
	if !first.Changed {
		t.Fatal("Inject(multi) first changed = false")
	}

	second, err := Inject(home, opencodeAdapter(), "multi")
	if err != nil {
		t.Fatalf("Inject(multi) second error = %v", err)
	}
	if second.Changed {
		t.Fatal("Inject(multi) second changed = true — should be idempotent with prompt files")
	}
}
