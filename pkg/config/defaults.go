package config

import "time"

func Default() Config {
	return Config{
		Product: ProductConfig{
			Name:    "burpsuite",
			Edition: "professional",
			Version: "latest",
		},
		Layout: LayoutConfig{
			Mode: SystemLayout,
		},
		Paths: PathsConfig{
			Install: "auto",
			Data:    "auto",
			Bin:     "auto",
		},
		Runtime: RuntimeConfig{
			Java: JavaConfig{
				Strategy:   JavaStrategyAuto,
				MinVersion: 17,
				JVMArgs: []string{
					"--add-opens=java.desktop/javax.swing=ALL-UNNAMED",
					"--add-opens=java.base/java.lang=ALL-UNNAMED",
					"-noverify",
				},
			},
		},
		Network: NetworkConfig{
			Timeout: 30 * time.Second,
			Retries: 3,
		},
		Logging: LoggingConfig{
			Level: LogLevelInfo,
		},
	}
}
