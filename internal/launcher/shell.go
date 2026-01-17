package launcher

import (
	"fmt"
	"os"
	"path/filepath"
	"text/template"

	"github.com/sdmrf/relay/internal/plan"
)

const shellTemplate = `#!/bin/sh
exec "{{.JavaPath}}" {{range .JVMArgs}}{{.}} {{end}}-jar "{{.JarPath}}" "$@" &
`

// ShellLauncher generates shell scripts for Linux/macOS.
type ShellLauncher struct {
	BinDir   string
	JavaPath string
}

func (s ShellLauncher) Path() string {
	return filepath.Join(s.BinDir, "burpsuite")
}

func (s ShellLauncher) Generate(p plan.LaunchPlan) error {
	t, err := template.New("launcher").Parse(shellTemplate)
	if err != nil {
		return fmt.Errorf("parse template: %w", err)
	}

	f, err := os.Create(s.Path())
	if err != nil {
		return fmt.Errorf("create launcher: %w", err)
	}
	defer f.Close()

	data := struct {
		JavaPath string
		JVMArgs  []string
		JarPath  string
	}{
		JavaPath: s.JavaPath,
		JVMArgs:  p.JVMArgs,
		JarPath:  filepath.Join(p.Paths.InstallDir, "burpsuite.jar"),
	}

	if err := t.Execute(f, data); err != nil {
		return fmt.Errorf("execute template: %w", err)
	}

	if err := os.Chmod(s.Path(), 0o755); err != nil {
		return fmt.Errorf("chmod launcher: %w", err)
	}

	return nil
}
