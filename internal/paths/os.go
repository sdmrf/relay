package paths

import "runtime"

type OS string

const (
	Linux   OS = "linux"
	Darwin  OS = "darwin"
	Windows OS = "windows"
)

func CurrentOS() OS {
	switch runtime.GOOS {
	case "linux":
		return Linux
	case "darwin":
		return Darwin
	case "windows":
		return Windows
	default:
		return OS(runtime.GOOS)
	}
}
