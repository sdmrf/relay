package config

import (
	"strings"
	"testing"
)

func TestValidate(t *testing.T) {
	tests := []struct {
		name    string
		config  Config
		wantErr string
	}{
		{
			name:   "valid default config",
			config: Default(),
		},
		{
			name: "missing product name",
			config: Config{
				Product: ProductConfig{Name: "", Version: "latest"},
				Layout:  LayoutConfig{Mode: SystemLayout},
				Runtime: RuntimeConfig{Java: JavaConfig{Strategy: JavaStrategyAuto, MinVersion: 17}},
				Logging: LoggingConfig{Level: LogLevelInfo},
			},
			wantErr: "product.name is required",
		},
		{
			name: "missing product version",
			config: Config{
				Product: ProductConfig{Name: "burpsuite", Version: ""},
				Layout:  LayoutConfig{Mode: SystemLayout},
				Runtime: RuntimeConfig{Java: JavaConfig{Strategy: JavaStrategyAuto, MinVersion: 17}},
				Logging: LoggingConfig{Level: LogLevelInfo},
			},
			wantErr: "product.version is required",
		},
		{
			name: "invalid layout mode",
			config: Config{
				Product: ProductConfig{Name: "burpsuite", Version: "latest"},
				Layout:  LayoutConfig{Mode: LayoutMode("invalid")},
				Runtime: RuntimeConfig{Java: JavaConfig{Strategy: JavaStrategyAuto, MinVersion: 17}},
				Logging: LoggingConfig{Level: LogLevelInfo},
			},
			wantErr: "invalid layout.mode",
		},
		{
			name: "invalid java strategy",
			config: Config{
				Product: ProductConfig{Name: "burpsuite", Version: "latest"},
				Layout:  LayoutConfig{Mode: SystemLayout},
				Runtime: RuntimeConfig{Java: JavaConfig{Strategy: JavaStrategy("invalid"), MinVersion: 17}},
				Logging: LoggingConfig{Level: LogLevelInfo},
			},
			wantErr: "invalid runtime.java.strategy",
		},
		{
			name: "java min version too low",
			config: Config{
				Product: ProductConfig{Name: "burpsuite", Version: "latest"},
				Layout:  LayoutConfig{Mode: SystemLayout},
				Runtime: RuntimeConfig{Java: JavaConfig{Strategy: JavaStrategyAuto, MinVersion: 7}},
				Logging: LoggingConfig{Level: LogLevelInfo},
			},
			wantErr: "min_version must be >= 8",
		},
		{
			name: "invalid logging level",
			config: Config{
				Product: ProductConfig{Name: "burpsuite", Version: "latest"},
				Layout:  LayoutConfig{Mode: SystemLayout},
				Runtime: RuntimeConfig{Java: JavaConfig{Strategy: JavaStrategyAuto, MinVersion: 17}},
				Logging: LoggingConfig{Level: LogLevel("invalid")},
			},
			wantErr: "invalid logging.level",
		},
		{
			name: "portable layout valid",
			config: Config{
				Product: ProductConfig{Name: "burpsuite", Version: "latest"},
				Layout:  LayoutConfig{Mode: PortableLayout},
				Runtime: RuntimeConfig{Java: JavaConfig{Strategy: JavaStrategySystem, MinVersion: 17}},
				Logging: LoggingConfig{Level: LogLevelDebug},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.config.Validate()

			if tt.wantErr == "" {
				if err != nil {
					t.Errorf("Validate() error = %v, want nil", err)
				}
				return
			}

			if err == nil {
				t.Errorf("Validate() error = nil, want error containing %q", tt.wantErr)
				return
			}

			if !strings.Contains(err.Error(), tt.wantErr) {
				t.Errorf("Validate() error = %v, want error containing %q", err, tt.wantErr)
			}
		})
	}
}

func TestDefaultConfig(t *testing.T) {
	cfg := Default()

	// Verify default values
	if cfg.Product.Name != "burpsuite" {
		t.Errorf("Default().Product.Name = %v, want burpsuite", cfg.Product.Name)
	}
	if cfg.Product.Edition != "professional" {
		t.Errorf("Default().Product.Edition = %v, want professional", cfg.Product.Edition)
	}
	if cfg.Product.Version != "latest" {
		t.Errorf("Default().Product.Version = %v, want latest", cfg.Product.Version)
	}
	if cfg.Layout.Mode != SystemLayout {
		t.Errorf("Default().Layout.Mode = %v, want system", cfg.Layout.Mode)
	}
	if cfg.Runtime.Java.MinVersion != 17 {
		t.Errorf("Default().Runtime.Java.MinVersion = %v, want 17", cfg.Runtime.Java.MinVersion)
	}
	if cfg.Runtime.Java.Strategy != JavaStrategyAuto {
		t.Errorf("Default().Runtime.Java.Strategy = %v, want auto", cfg.Runtime.Java.Strategy)
	}
	if cfg.Logging.Level != LogLevelInfo {
		t.Errorf("Default().Logging.Level = %v, want info", cfg.Logging.Level)
	}

	// Default config should be valid
	if err := cfg.Validate(); err != nil {
		t.Errorf("Default().Validate() error = %v", err)
	}
}

func TestDefaultJVMArgs(t *testing.T) {
	cfg := Default()

	if len(cfg.Runtime.Java.JVMArgs) == 0 {
		t.Error("Default().Runtime.Java.JVMArgs should not be empty")
	}

	// Check for required JVM args for Burp Suite
	hasSwingOpen := false
	hasLangOpen := false
	for _, arg := range cfg.Runtime.Java.JVMArgs {
		if strings.Contains(arg, "javax.swing") {
			hasSwingOpen = true
		}
		if strings.Contains(arg, "java.lang") {
			hasLangOpen = true
		}
	}

	if !hasSwingOpen {
		t.Error("Default JVM args should include swing module open")
	}
	if !hasLangOpen {
		t.Error("Default JVM args should include lang module open")
	}
}
