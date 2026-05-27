package codegraph

// Setup is a no-op for CodeGraph. The CodeGraph binary is installed globally
// and does not require per-agent setup like Engram does.
func Setup() error {
	return nil
}
