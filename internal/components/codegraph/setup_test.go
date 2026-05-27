package codegraph

import "testing"

func TestSetupReturnsNil(t *testing.T) {
	if err := Setup(); err != nil {
		t.Fatalf("Setup() error = %v, want nil", err)
	}
}
