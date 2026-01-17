# Command Reference

Complete reference for all relay commands.

## Global Flags

These flags are available for all commands:

| Flag | Short | Description | Default |
|------|-------|-------------|---------|
| `--config` | `-c` | Config file path | `config.yaml` |
| `--dry-run` | | Preview actions without executing | `false` |
| `--verbose` | `-v` | Verbose output | `false` |
| `--help` | `-h` | Help for the command | |

## Commands

### relay install

Download and install Burp Suite.

```bash
relay install [flags]
```

**Flags:**

| Flag | Description | Default |
|------|-------------|---------|
| `--edition` | Edition to install (professional, community) | from config |
| `--version` | Version to install | `latest` |

**Examples:**

```bash
# Install using config defaults
relay install

# Install Community Edition
relay install --edition community

# Preview installation without executing
relay install --dry-run

# Install with verbose output
relay install -v
```

**What it does:**

1. Creates installation directories
2. Creates data directories
3. Creates binary directory
4. Creates cache directory
5. Downloads Burp Suite JAR from PortSwigger CDN

---

### relay launch

Start Burp Suite.

```bash
relay launch [flags]
```

**Examples:**

```bash
# Launch Burp Suite
relay launch

# Preview launch without executing
relay launch --dry-run

# Launch with verbose output
relay launch -v
```

**What it does:**

1. Validates Java installation meets minimum version
2. Generates platform-specific launcher script
3. Executes the launcher to start Burp Suite

---

### relay update

Check for and apply updates.

```bash
relay update [flags]
```

**Flags:**

| Flag | Short | Description | Default |
|------|-------|-------------|---------|
| `--force` | `-f` | Force update even if at latest version | `false` |

**Examples:**

```bash
# Update to latest version
relay update

# Force update (re-download even if current)
relay update --force

# Preview update without executing
relay update --dry-run
```

**What it does:**

1. Reads current installed version from marker file
2. Compares with target version
3. Downloads new version if update needed
4. Updates version marker file

---

### relay remove

Uninstall Burp Suite.

```bash
relay remove [flags]
```

**Examples:**

```bash
# Remove Burp Suite
relay remove

# Preview removal without executing
relay remove --dry-run

# Remove with verbose output
relay remove -v
```

**What it does:**

1. Removes installation directory
2. Removes data directory
3. Removes cache directory
4. **Preserves** configuration directory

---

### relay doctor

Run diagnostic checks to verify system readiness.

```bash
relay doctor [flags]
```

**Examples:**

```bash
# Run diagnostics
relay doctor

# Run with verbose output (shows all details)
relay doctor -v
```

**Checks performed:**

1. **Java** - Verifies Java is installed and meets minimum version
2. **Config** - Validates configuration file
3. **Paths** - Checks directories exist and are writable
4. **Product** - Verifies Burp Suite JAR is present
5. **Network** - Tests connectivity to PortSwigger

**Output format:**

```
relay doctor

[✓] Java: Java 21 found
    /usr/bin/java
[✓] Config: Loaded from config.yaml
[✓] Install directory: Directory exists and writable
    /Applications/relay
[!] Data directory: Directory not writable
    /Library/Application Support/relay
[✓] Product: Burp Suite JAR present (612 MB)
    /Applications/relay/burpsuite.jar
[✓] Network: portswigger.net reachable

All critical checks passed with warnings.
```

**Status icons:**

| Icon | Meaning |
|------|---------|
| `✓` | Check passed |
| `!` | Warning (non-critical) |
| `✗` | Check failed |

---

### relay version

Print version information.

```bash
relay version
```

**Output:**

```
relay v1.0.0
  commit:  abc1234
  built:   2024-01-15T10:30:00Z
```

---

## Exit Codes

| Code | Meaning |
|------|---------|
| 0 | Success |
| 1 | Error occurred |

## Common Workflows

### Fresh Installation

```bash
# Check system readiness
relay doctor

# Install Burp Suite
relay install

# Launch Burp Suite
relay launch
```

### Update Existing Installation

```bash
# Check for updates and apply
relay update

# Or force re-download
relay update --force
```

### Clean Reinstall

```bash
# Remove existing installation
relay remove

# Fresh install
relay install
```

### Troubleshooting

```bash
# Check what's wrong
relay doctor -v

# Preview what install would do
relay install --dry-run

# Verbose install for debugging
relay install -v
```
