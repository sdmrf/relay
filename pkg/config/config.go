package config

import "time"

type LayoutMode string

const (
	SystemLayout   LayoutMode = "system"
	PortableLayout LayoutMode = "portable"
)

type JavaStrategy string

const (
	JavaStrategyAuto    JavaStrategy = "auto"    // Prefer bundled, fallback to system
	JavaStrategySystem  JavaStrategy = "system"  // Only use system Java
	JavaStrategyBundled JavaStrategy = "bundled" // Only use bundled JRE
)

type LogLevel string

const (
	LogLevelInfo  LogLevel = "info"
	LogLevelDebug LogLevel = "debug"
	LogLevelTrace LogLevel = "trace"
)

type Config struct {
	Product ProductConfig `yaml:"product"`
	Layout  LayoutConfig  `yaml:"layout"`
	Paths   PathsConfig   `yaml:"paths"`
	Runtime RuntimeConfig `yaml:"runtime"`
	Network NetworkConfig `yaml:"network"`
	Logging LoggingConfig `yaml:"logging"`
}

type ProductConfig struct {
	Name    string `yaml:"name"`
	Edition string `yaml:"edition"`
	Version string `yaml:"version"`
}

type LayoutConfig struct {
	Mode LayoutMode `yaml:"mode"`
}

type PathsConfig struct {
	Install string `yaml:"install"`
	Data    string `yaml:"data"`
	Bin     string `yaml:"bin"`
}

type RuntimeConfig struct {
	Java JavaConfig `yaml:"java"`
}

type JavaConfig struct {
	Strategy   JavaStrategy `yaml:"strategy"`
	MinVersion int          `yaml:"min_version"`
	JVMArgs    []string     `yaml:"jvm_args"`
}

type NetworkConfig struct {
	Timeout time.Duration `yaml:"timeout"`
	Retries int           `yaml:"retries"`
}

type LoggingConfig struct {
	Level LogLevel `yaml:"level"`
}
