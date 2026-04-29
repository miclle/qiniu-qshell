package sandbox

import (
	"os"

	"golang.org/x/term"
)

// IsInteractive reports whether stdin is attached to a terminal.
// Used to guard huh.Form.Run() calls so they don't hang when the CLI is
// invoked from CI, AI agents, pipes, or background contexts.
func IsInteractive() bool {
	return term.IsTerminal(int(os.Stdin.Fd()))
}
