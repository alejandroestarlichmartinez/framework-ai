package codegraph

import (
	"reflect"
	"testing"

	"github.com/alejandroestarlichmartinez/framework-ai/internal/system"
)

func TestInstallCommandReturnsCurlScript(t *testing.T) {
	profile := system.PlatformProfile{OS: "linux", PackageManager: "apt"}
	got, err := InstallCommand(profile)
	if err != nil {
		t.Fatalf("InstallCommand() error = %v", err)
	}
	want := [][]string{{"sh", "-c", "curl -fsSL https://raw.githubusercontent.com/colbymchenry/codegraph/main/install.sh | sh"}}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("InstallCommand() = %v, want %v", got, want)
	}
}
