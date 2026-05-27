package cli

import (
	"testing"

	"github.com/alejandroestarlichmartinez/framework-ai/internal/model"
	"github.com/alejandroestarlichmartinez/framework-ai/internal/system"
)

func TestComponentsForPresetFullGentlemanIncludesCodeGraph(t *testing.T) {
	components := componentsForPreset(model.PresetFullGentleman)
	var found bool
	for _, c := range components {
		if c == model.ComponentCodeGraph {
			found = true
			break
		}
	}
	if !found {
		t.Fatal("PresetFullGentleman missing ComponentCodeGraph")
	}
}

func TestNormalizePersonaAcceptsRickSanchez(t *testing.T) {
	got, err := normalizePersona("rick-sanchez", model.PresetFullGentleman)
	if err != nil {
		t.Fatalf("normalizePersona(rick-sanchez) error = %v", err)
	}
	if got != model.PersonaRickSanchez {
		t.Fatalf("normalizePersona(rick-sanchez) = %q, want %q", got, model.PersonaRickSanchez)
	}
}

func TestNormalizePresetAcceptsFullRick(t *testing.T) {
	got, err := normalizePreset("full-rick")
	if err != nil {
		t.Fatalf("normalizePreset(full-rick) error = %v", err)
	}
	if got != model.PresetFullRick {
		t.Fatalf("normalizePreset(full-rick) = %q, want %q", got, model.PresetFullRick)
	}
}

func TestComponentsForPresetFullRickMatchesFullGentleman(t *testing.T) {
	rick := componentsForPreset(model.PresetFullRick)
	gentleman := componentsForPreset(model.PresetFullGentleman)
	if len(rick) != len(gentleman) {
		t.Fatalf("PresetFullRick has %d components, PresetFullGentleman has %d", len(rick), len(gentleman))
	}
	for i := range gentleman {
		if gentleman[i] != rick[i] {
			t.Fatalf("component mismatch at index %d: FullGentleman=%q, FullRick=%q", i, gentleman[i], rick[i])
		}
	}
}

func TestNormalizeInstallFlagsFullRickDefaultsToRickSanchez(t *testing.T) {
	flags := InstallFlags{Preset: "full-rick"}
	input, err := NormalizeInstallFlags(flags, system.DetectionResult{})
	if err != nil {
		t.Fatalf("NormalizeInstallFlags error = %v", err)
	}
	if input.Selection.Persona != model.PersonaRickSanchez {
		t.Fatalf("want persona %q, got %q", model.PersonaRickSanchez, input.Selection.Persona)
	}
	if input.Selection.Preset != model.PresetFullRick {
		t.Fatalf("want preset %q, got %q", model.PresetFullRick, input.Selection.Preset)
	}
}

func TestNormalizeInstallFlagsFullRickWithExplicitPersonaOverride(t *testing.T) {
	flags := InstallFlags{Preset: "full-rick", Persona: "gentleman"}
	input, err := NormalizeInstallFlags(flags, system.DetectionResult{})
	if err != nil {
		t.Fatalf("NormalizeInstallFlags error = %v", err)
	}
	if input.Selection.Persona != model.PersonaGentleman {
		t.Fatalf("want persona %q, got %q", model.PersonaGentleman, input.Selection.Persona)
	}
	if input.Selection.Preset != model.PresetFullRick {
		t.Fatalf("want preset %q, got %q", model.PresetFullRick, input.Selection.Preset)
	}
}
