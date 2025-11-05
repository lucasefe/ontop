# Config Module Interface Contract

**Feature**: 004-row-ui | **Module**: internal/config | **Date**: 2025-11-05

## Overview

The config module provides persistent storage and retrieval of user preferences in TOML format. This contract defines the public API for loading and saving configuration.

## Package Interface

### Package Name

```go
package config
```

### Public Types

#### Config Struct

```go
// Config represents the application configuration
type Config struct {
    UI UIConfig `toml:"ui"`
}

// UIConfig holds user interface preferences
type UIConfig struct {
    ViewMode string `toml:"view_mode"` // "column" or "row"
}
```

**Constraints**:
- `ViewMode` MUST be either "column" or "row"
- Invalid values SHOULD default to "column"
- Additional fields MAY be added in future without breaking changes

### Public Functions

#### Load

```go
// Load reads the config file from the standard location and returns
// a Config struct. If the file doesn't exist or can't be parsed,
// returns a Config with default values and logs a warning.
//
// Default location: ~/.config/ontop/ontop.toml
//
// Returns:
//   - Config with loaded values (or defaults on error)
//   - error (nil if successful, non-nil if load failed but defaults used)
func Load() (Config, error)
```

**Behavior**:
- File not found → Return defaults, error = nil
- File malformed → Return defaults, error = non-nil (logged)
- File unreadable → Return defaults, error = non-nil (logged)
- Valid file → Return parsed config, error = nil
- Never panics

**Defaults**:
```go
Config{
    UI: UIConfig{
        ViewMode: "column",
    },
}
```

#### Save

```go
// Save writes the config to the standard location. Creates the directory
// if it doesn't exist. Uses atomic write (temp file + rename).
//
// Parameters:
//   - cfg: The Config struct to persist
//
// Returns:
//   - error (nil if successful, non-nil if write failed)
func Save(cfg Config) error
```

**Behavior**:
- Creates `~/.config/ontop/` if not exists (with 0755 permissions)
- Writes to temp file first (`.ontop.toml.tmp`)
- Atomic rename to `ontop.toml` (0644 permissions)
- Cleans up temp file on error
- Returns error if write fails (non-fatal to caller)

#### GetConfigPath

```go
// GetConfigPath returns the absolute path to the config file.
// Useful for debugging and logging.
//
// Returns:
//   - string: Absolute path to config file (e.g., "/home/user/.config/ontop/ontop.toml")
func GetConfigPath() string
```

**Behavior**:
- Resolves home directory using `os.UserHomeDir()`
- Returns XDG-compliant path: `$HOME/.config/ontop/ontop.toml`
- Never panics (returns empty string if home dir unavailable)

## Error Handling

### Error Types

The module uses standard `error` interface. Common error scenarios:

| Error Scenario | Returned Error | Config Returned | Caller Action |
|----------------|----------------|-----------------|---------------|
| File not found | nil | Defaults | Use defaults silently |
| Parse error | Non-nil (toml.DecodeError) | Defaults | Log warning, use defaults |
| Write error | Non-nil (os.PathError) | N/A | Log error, continue |
| Home dir unavailable | Non-nil | Defaults | Log error, use defaults |

### Non-Fatal Errors

All config operations are designed to be non-fatal:

- **Load failure**: App continues with defaults
- **Save failure**: App continues with in-memory state (preference just not persisted)
- **No panic guarantee**: Package never panics under any error condition

## Usage Examples

### Basic Load/Save Cycle

```go
import "github.com/lucasefe/ontop/internal/config"

// Load config at startup
cfg, err := config.Load()
if err != nil {
    log.Printf("Config load warning: %v (using defaults)", err)
}

// Use config
viewMode := cfg.UI.ViewMode  // "column" or "row"

// Modify and save
cfg.UI.ViewMode = "row"
if err := config.Save(cfg); err != nil {
    log.Printf("Config save error: %v", err)
}
```

### Integration with TUI Model

```go
// In NewModel() initialization
cfg, _ := config.Load()  // Ignore error, defaults are safe
model := Model{
    viewLayout: parseViewLayout(cfg.UI.ViewMode),
    // ... other fields
}

// In view toggle handler
func (m Model) handleToggle() (Model, tea.Cmd) {
    // Toggle layout
    m.viewLayout = toggleLayout(m.viewLayout)

    // Save preference
    cfg := config.Config{
        UI: config.UIConfig{
            ViewMode: formatViewLayout(m.viewLayout),
        },
    }
    go func() {
        if err := config.Save(cfg); err != nil {
            log.Printf("Failed to save view mode preference: %v", err)
        }
    }()

    return m, nil
}
```

## Dependencies

### External Dependencies

```go
import (
    "github.com/BurntSushi/toml" // v1.4.0
)
```

### Standard Library

```go
import (
    "os"
    "path/filepath"
)
```

## File Format Contract

### TOML Schema

```toml
# OnTop Configuration
# Generated: <ISO8601 timestamp>

[ui]
view_mode = "column"  # "column" | "row"
```

**Schema Version**: 1.0 (implicit, no version field needed yet)

**Compatibility**:
- Unknown fields MUST be ignored (TOML parser default behavior)
- Missing fields MUST use defaults
- Invalid enum values MUST use defaults

### Example Files

**Default (column mode)**:
```toml
[ui]
view_mode = "column"
```

**Row mode**:
```toml
[ui]
view_mode = "row"
```

**Future extension (example)**:
```toml
[ui]
view_mode = "row"
theme = "gruvbox"  # Future field, ignored by current version
```

## Testing Contract

### Unit Tests Required

1. **Load Tests**:
   - `TestLoad_FileNotFound`: Returns defaults, no error
   - `TestLoad_ValidFile`: Parses correctly
   - `TestLoad_InvalidTOML`: Returns defaults, returns error
   - `TestLoad_InvalidViewMode`: Returns defaults ("column")

2. **Save Tests**:
   - `TestSave_Success`: Writes file correctly
   - `TestSave_CreatesDirectory`: Creates ~/.config/ontop/ if missing
   - `TestSave_AtomicWrite`: Uses temp file + rename
   - `TestSave_PermissionsError`: Returns error gracefully

3. **Path Tests**:
   - `TestGetConfigPath`: Returns XDG-compliant path

### Integration Tests Required

1. **Round-trip**: Load → Modify → Save → Load → Verify
2. **Concurrent writes**: Multiple Save() calls don't corrupt file
3. **Permission scenarios**: Read-only directory, read-only file

## Performance Contract

### Latency Requirements

| Operation | Target | Maximum | Measured Method |
|-----------|--------|---------|-----------------|
| Load() | <1ms | 5ms | Benchmark |
| Save() | <1ms | 10ms | Benchmark |

### Resource Constraints

| Resource | Limit | Justification |
|----------|-------|---------------|
| Memory | <1KB | Small struct, short-lived |
| Disk I/O | 1 read + 1 write per app session | Load once, save on toggle |
| File size | <10KB | Future extensions, human-editable |

## Security Considerations

### File Permissions

- **Config directory**: 0755 (drwxr-xr-x) - User-owned, group-readable
- **Config file**: 0644 (-rw-r--r--) - User-writable, world-readable

**Rationale**: Config contains no secrets, user preferences only. Standard permissions for non-sensitive config files.

### Attack Vectors

| Vector | Mitigation |
|--------|------------|
| Config injection | TOML parser handles escaping, no eval |
| Path traversal | Use filepath.Join(), no user-supplied paths |
| Symlink attacks | Atomic write with temp file minimizes window |
| TOCTOU races | Acceptable (config is non-critical data) |

**Note**: This is user preference config, not security-sensitive. Defense-in-depth but not paranoid.

## Backwards Compatibility

### Version 1.0 Guarantees

- Config struct will only add fields (never remove)
- ViewMode enum values "column" and "row" are permanent
- GetConfigPath() returns same location across versions
- TOML format is stable (TOML 1.0.0 spec)

### Breaking Changes Policy

Breaking changes require major version bump:
- Removing fields from Config struct
- Changing ViewMode enum values
- Moving config file location
- Changing file format (e.g., TOML → JSON)

## Version History

| Version | Date | Changes |
|---------|------|---------|
| 1.0 | 2025-11-05 | Initial contract - ViewMode persistence |

## Future Extensions (Not in Scope)

These MAY be added in future without breaking this contract:

```go
type UIConfig struct {
    ViewMode  string `toml:"view_mode"`
    Theme     string `toml:"theme"`      // Future: "gruvbox", "solarized", etc.
    ShowHelp  bool   `toml:"show_help"`  // Future: persistent help visibility
    SortMode  string `toml:"sort_mode"`  // Future: default sort mode
}

type BehaviorConfig struct {
    AutoArchive bool `toml:"auto_archive"`  // Future: auto-archive completed tasks
    ArchiveAfterDays int `toml:"archive_after_days"`
}

type Config struct {
    UI       UIConfig       `toml:"ui"`
    Behavior BehaviorConfig `toml:"behavior"`  // Future section
}
```

All future fields MUST have sensible defaults and be optional (TOML unmarshaling graceful for missing fields).
