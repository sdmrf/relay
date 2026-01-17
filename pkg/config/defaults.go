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
					// Memory settings
					"-Xmx4g",
					"-Xms1g",
					// Modern G1 garbage collector (default in JDK 17+, explicit for clarity)
					"-XX:+UseG1GC",
					"-XX:+UseStringDeduplication",
					"-XX:MaxGCPauseMillis=100",
					// Required module access for Burp Suite GUI
					"--add-opens=java.desktop/javax.swing=ALL-UNNAMED",
					"--add-opens=java.desktop/java.awt=ALL-UNNAMED",
					"--add-opens=java.desktop/java.awt.event=ALL-UNNAMED",
					"--add-opens=java.base/java.lang=ALL-UNNAMED",
					"--add-opens=java.base/java.lang.reflect=ALL-UNNAMED",
					"--add-opens=java.base/java.io=ALL-UNNAMED",
					"--add-opens=java.base/java.util=ALL-UNNAMED",
					"--add-opens=java.base/java.util.concurrent=ALL-UNNAMED",
					"--add-opens=java.base/jdk.internal.misc=ALL-UNNAMED",
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
