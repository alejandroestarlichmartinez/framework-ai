package app

import (
	"fmt"
	"io"
)

func printHelp(w io.Writer, version string) {
	fmt.Fprintf(w, `framework-ai — Framework-AI: Ecosystem, Frameworks, Workflows (%s)

USAGE
  framework-ai                     Launch interactive TUI
  framework-ai <command> [flags]

COMMANDS
  install      Configure AI coding agents on this machine
  uninstall    Remove Framework AI managed files from this machine
  sync         Sync agent configs and skills to current version
  skill-registry refresh
               Refresh .atl/skill-registry.md with cache-hit fast path
  update       Check for available updates
  upgrade      Apply updates to managed tools
  restore      Restore a config backup
  version      Print version

FLAGS
  --help, -h    Show this help

Run 'framework-ai help' for this message.
Documentation: https://github.com/alejandroestarlichmartinez/framework-ai
`, version)
}
