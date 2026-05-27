package gga

import (
	"github.com/alejandroestarlichmartinez/framework-ai/internal/installcmd"
	"github.com/alejandroestarlichmartinez/framework-ai/internal/model"
	"github.com/alejandroestarlichmartinez/framework-ai/internal/system"
)

func InstallCommand(profile system.PlatformProfile) ([][]string, error) {
	return installcmd.NewResolver().ResolveComponentInstall(profile, model.ComponentGGA)
}

func ShouldInstall(enabled bool) bool {
	return enabled
}
