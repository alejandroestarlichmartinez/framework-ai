package planner

import (
	"testing"

	"github.com/alejandroestarlichmartinez/framework-ai/internal/model"
)

func TestMVPGraphIncludesCodeGraph(t *testing.T) {
	g := MVPGraph()
	if !g.Has(model.ComponentCodeGraph) {
		t.Fatal("MVPGraph() missing CodeGraph")
	}
	deps := g.DependenciesOf(model.ComponentCodeGraph)
	if len(deps) != 0 {
		t.Fatalf("CodeGraph dependencies = %v, want nil", deps)
	}
}
