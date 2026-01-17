package config

import "time"

type Config struct {
	Product  ProductConfig  `yaml:"product"`
	Layout   LayoutConfig   `yaml:"layout"`
	Paths    PathsConfig    `yaml:"paths"`
	Runtime  RuntimeConfig  `yaml:"runtime"`
	Network  NetworkConfig  `yaml:"network"`
	Logging  LoggingConfig  `yaml:"logging"`
}

type ProductConfig struct {
	Name    string `yaml:"name"`
	Edition string `yaml:"edition"`
	Version string `yaml:"version"`
}

type LayoutConfig struct {
	Mode string `yaml:"mode"`
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
	Strategy   string   `yaml:"strategy"`
	MinVersion int      `yaml:"min_version"`
	JVMArgs    []string `yaml:"jvm_args"`
}

type NetworkConfig struct {
	Timeout time.Duration `yaml:"timeout"`
	Retries int           `yaml:"retries"`
}

type LoggingConfig struct {
	Level string `yaml:"level"`
}
