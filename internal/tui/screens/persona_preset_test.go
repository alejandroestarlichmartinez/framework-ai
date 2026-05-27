package screens

import (
	"strings"
	"testing"

	"github.com/alejandroestarlichmartinez/framework-ai/internal/model"
)

func TestRenderPersonaClarifiesCustomKeepsExistingPersona(t *testing.T) {
	out := RenderPersona(model.PersonaCustom, 2)

	if !strings.Contains(out, "custom") {
		t.Fatalf("RenderPersona missing custom option; output:\n%s", out)
	}
	if !strings.Contains(out, "Keep your existing persona unmanaged") {
		t.Fatalf("RenderPersona missing custom persona clarification; output:\n%s", out)
	}
	if strings.Contains(out, "Bring your own persona instructions") {
		t.Fatalf("RenderPersona still shows old custom persona wording; output:\n%s", out)
	}
}

func TestRenderPresetClarifiesCustomManualSelection(t *testing.T) {
	out := RenderPreset(model.PresetCustom, 4)

	if !strings.Contains(out, "Choose components and skills manually") {
		t.Fatalf("RenderPreset missing custom preset clarification; output:\n%s", out)
	}
	if strings.Contains(out, "Pick individual components yourself") {
		t.Fatalf("RenderPreset still shows old custom preset wording; output:\n%s", out)
	}
}

func TestPersonaOptionsIncludesRickSanchez(t *testing.T) {
	options := PersonaOptions()
	var found bool
	for _, p := range options {
		if p == model.PersonaRickSanchez {
			found = true
			break
		}
	}
	if !found {
		t.Fatal("PersonaOptions() missing PersonaRickSanchez")
	}
}

func TestPresetOptionsIncludesFullRick(t *testing.T) {
	options := PresetOptions()
	var found bool
	for _, p := range options {
		if p == model.PresetFullRick {
			found = true
			break
		}
	}
	if !found {
		t.Fatal("PresetOptions() missing PresetFullRick")
	}
}

func TestRenderPersonaIncludesRickDescription(t *testing.T) {
	out := RenderPersona(model.PersonaRickSanchez, 1)
	if !strings.Contains(out, "cynical genius") {
		t.Fatalf("RenderPersona missing Rick description; output:\n%s", out)
	}
}

func TestRenderPresetIncludesFullRickDescription(t *testing.T) {
	out := RenderPreset(model.PresetFullRick, 1)
	if !strings.Contains(out, "Rick Sanchez") {
		t.Fatalf("RenderPreset missing FullRick description; output:\n%s", out)
	}
}
