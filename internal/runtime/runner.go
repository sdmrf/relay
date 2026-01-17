package runtime

import (
	"context"
	"os"
	"os/exec"
)

// ExecRunner executes launcher scripts as child processes.
type ExecRunner struct{}

// Run starts the launcher and returns immediately (non-blocking).
// Burp runs as a background process - relay exits cleanly.
func (ExecRunner) Run(ctx context.Context, path string) error {
	cmd := exec.CommandContext(ctx, path)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	return cmd.Start()
}
