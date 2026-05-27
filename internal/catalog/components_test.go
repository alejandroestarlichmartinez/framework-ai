package catalog

import (
	"testing"

	"github.com/alejandroestarlichmartinez/framework-ai/internal/model"
)

func TestMVPComponentsIncludesCodeGraph(t *testing.T) {
	components := MVPComponents()
	var found bool
	for _, c := range components {
		if c.ID == model.ComponentCodeGraph {
			found = true
			if c.Name != "CodeGraph" {
				t.Fatalf("CodeGraph name = %q, want %q", c.Name, "CodeGraph")
			}
			if c.Description != "Semantic code knowledge graph" {
				t.Fatalf("CodeGraph description = %q, want %q", c.Description, "Semantic code knowledge graph")
			}
			break
		}
	}
	if !found {
		t.Fatalf("MVPComponents() missing CodeGraph component")
	}
}
