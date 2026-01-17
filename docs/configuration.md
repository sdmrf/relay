# Configuration Reference

relay uses YAML configuration files to control behavior. By default, it looks for `config.yaml` in the current directory.

## Configuration File Location

```bash
# Use default config.yaml
relay install

# Use custom config file
relay install -c /path/to/config.yaml
```

## Full Configuration Reference

```yaml
# Product configuration
product:
  name: burpsuite           # Product name (currently only "burpsuite" supported)
  edition: professional     # "professional" or "community"
  version: latest           # Version string or "latest"

# Layout configuration
layout:
  mode: system              # "system" or "portable"

# Path hints (optional, usually auto-detected)
paths:
  install: auto             # Installation directory
  data: auto                # Data directory
  bin: auto                 # Binary/launcher directory

# Runtime configuration
runtime:
  java:
    strategy: auto          # "auto" or "system"
    min_version: 17         # Minimum required Java version
    jvm_args:               # JVM arguments passed to Java
      - "--add-opens=java.desktop/javax.swing=ALL-UNNAMED"
      - "--add-opens=java.base/java.lang=ALL-UNNAMED"
      - "-noverify"

# Network configuration
network:
  timeout: 30s              # Download timeout
  retries: 3                # Number of retry attempts

# Logging configuration
logging:
  level: info               # "info", "debug", or "trace"
```

## Layout Modes

### System Layout

Uses standard OS-specific paths:

| OS | Install | Data | Bin |
|---|---|---|---|
| Linux | `/usr/lib/relay` | `/var/lib/relay` | `/usr/bin` |
| macOS | `/Applications/relay` | `/Library/Application Support/relay` | `/usr/local/bin` |
| Windows | `C:\Program Files\relay` | `C:\ProgramData\relay` | `C:\Program Files\relay\bin` |

### Portable Layout

Self-contained directory structure:

```
<install_path>/
├── bin/
├── data/
├── config/
└── cache/
```

Enable portable mode:

```yaml
layout:
  mode: portable

paths:
  install: /path/to/burpsuite
```

## Product Editions

### Professional Edition

```yaml
product:
  edition: professional
```

Downloads Burp Suite Professional. Requires a valid license to use.

### Community Edition

```yaml
product:
  edition: community
```

Downloads the free Burp Suite Community Edition.

## Java Configuration

### Auto Strategy

Automatically detects Java installation:

```yaml
runtime:
  java:
    strategy: auto
    min_version: 17
```

### System Strategy

Uses the system Java installation:

```yaml
runtime:
  java:
    strategy: system
    min_version: 17
```

### Custom JVM Arguments

Add custom JVM arguments for Burp Suite:

```yaml
runtime:
  java:
    jvm_args:
      - "-Xmx4g"                    # 4GB max heap
      - "-Xms1g"                    # 1GB initial heap
      - "--add-opens=java.desktop/javax.swing=ALL-UNNAMED"
      - "--add-opens=java.base/java.lang=ALL-UNNAMED"
      - "-noverify"
```

## Environment Variables

relay respects the following environment variables:

| Variable | Description |
|----------|-------------|
| `JAVA_HOME` | Java installation directory |
| `PATH` | Used to locate Java executable |

## Default Configuration

If no config file exists, relay uses sensible defaults:

```yaml
product:
  name: burpsuite
  edition: professional
  version: latest

layout:
  mode: system

paths:
  install: auto
  data: auto
  bin: auto

runtime:
  java:
    strategy: auto
    min_version: 17
    jvm_args:
      - "--add-opens=java.desktop/javax.swing=ALL-UNNAMED"
      - "--add-opens=java.base/java.lang=ALL-UNNAMED"
      - "-noverify"

network:
  timeout: 30s
  retries: 3

logging:
  level: info
```
