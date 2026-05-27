package model

import (
	"testing"
)

func TestComponentCodeGraphConstant(t *testing.T) {
	if ComponentCodeGraph != "codegraph" {
		t.Fatalf("ComponentCodeGraph = %q, want %q", ComponentCodeGraph, "codegraph")
	}
}

func TestPersonaRickSanchezConstant(t *testing.T) {
	if PersonaRickSanchez != "rick-sanchez" {
		t.Fatalf("PersonaRickSanchez = %q, want %q", PersonaRickSanchez, "rick-sanchez")
	}
}

func TestPresetFullRickConstant(t *testing.T) {
	if PresetFullRick != "full-rick" {
		t.Fatalf("PresetFullRick = %q, want %q", PresetFullRick, "full-rick")
	}
}
