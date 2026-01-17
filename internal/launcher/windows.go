package launcher

import (
	"fmt"
	"os"
	"path/filepath"
	"text/template"

	"github.com/sdmrf/relay/internal/plan"
)

const psTemplate = `Start-Process "{{.JavaPath}}" -ArgumentList '{{range .JVMArgs}}{{.}} {{end}}-jar "{{.JarPath}}"' -NoNewWindow
`

// PowerShellLauncher generates PowerShell scripts for Windows.
type PowerShellLauncher struct {
	BinDir   string
	JavaPath string
}

func (p PowerShellLauncher) Path() string {
	return filepath.Join(p.BinDir, "burpsuite.ps1")
}

func (p PowerShellLauncher) Generate(lp plan.LaunchPlan) error {
	t, err := template.New("ps").Parse(psTemplate)
	if err != nil {
		return fmt.Errorf("parse template: %w", err)
	}

	f, err := os.Create(p.Path())
	if err != nil {
		return fmt.Errorf("create launcher: %w", err)
	}
	defer f.Close()

	data := struct {
		JavaPath string
		JVMArgs  []string
		JarPath  string
	}{
		JavaPath: p.JavaPath,
		JVMArgs:  lp.JVMArgs,
		JarPath:  filepath.Join(lp.Paths.InstallDir, "burpsuite.jar"),
	}

	if err := t.Execute(f, data); err != nil {
		return fmt.Errorf("execute template: %w", err)
	}

	return nil
}
