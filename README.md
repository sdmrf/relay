# relay

A modern CLI for managing Burp Suite installations.

relay provides deterministic, repeatable commands for installing, updating, launching, and managing Burp Suite across platforms.

## Features

- **Install**: Download and install Burp Suite with a single command
- **Launch**: Start Burp Suite with proper Java configuration
- **Update**: Check for and apply updates
- **Remove**: Clean uninstallation preserving user configuration
- **Doctor**: Diagnose system readiness and configuration issues

## Installation

### From Source

```bash
go install github.com/sdmrf/relay/cmd/relay@latest
```

### From Releases

Download the latest release from the [releases page](https://github.com/sdmrf/relay/releases).

## Quick Start

```bash
# Install Burp Suite Professional
relay install --edition professional

# Launch Burp Suite
relay launch

# Check system status
relay doctor

# Update to latest version
relay update
```

## Configuration

relay uses a YAML configuration file (`config.yaml` by default):

```yaml
product:
  name: burpsuite
  edition: professional
  version: latest

layout:
  mode: system  # or "portable"

runtime:
  java:
    strategy: auto
    min_version: 17
    jvm_args:
      - "--add-opens=java.desktop/javax.swing=ALL-UNNAMED"
      - "--add-opens=java.base/java.lang=ALL-UNNAMED"
      - "-noverify"
```

See [docs/configuration.md](docs/configuration.md) for the full configuration reference.

## Commands

| Command | Description |
|---------|-------------|
| `relay install` | Install Burp Suite |
| `relay launch` | Launch Burp Suite |
| `relay update` | Update to latest version |
| `relay remove` | Uninstall Burp Suite |
| `relay doctor` | Run diagnostic checks |
| `relay version` | Show version information |

See [docs/commands.md](docs/commands.md) for the full command reference.

## Global Flags

| Flag | Description |
|------|-------------|
| `-c, --config` | Config file path (default: `config.yaml`) |
| `--dry-run` | Preview actions without executing |
| `-v, --verbose` | Verbose output |

## Requirements

- Go 1.22+ (for building)
- Java 17+ (for running Burp Suite)

## License

MIT - See [LICENSE](LICENSE) for details.
