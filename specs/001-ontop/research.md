# Research: OnTop - Dual-Mode Task Manager

**Date**: 2025-11-04
**Branch**: 001-ontop

## Overview

This document captures technical research and decisions for implementing a dual-mode (TUI + CLI) task manager in Go with SQLite persistence.

---

## 1. TUI Framework: Bubbletea

**Decision**: Use Charm's Bubbletea framework for TUI implementation

**Rationale**:
- Elm architecture (Model-Update-View) provides clean separation of concerns
- Built-in support for keyboard input handling
- Active maintenance and strong community
- Excellent documentation and examples
- Integrates well with Bubbles (reusable components) and Lipgloss (styling)
- No C dependencies, pure Go

**Alternatives Considered**:
- **tview**: More widget-focused, less control over rendering, heavier abstraction
- **termui**: Dashboard/graph focused, not ideal for interactive keyboard-driven UIs
- **tcell** (lower-level): More control but requires building everything from scratch

**Best Practices**:
- Use `tea.Model` interface for each view (kanban, detail)
- Handle `tea.KeyMsg` for vim-style keybindings
- Use `tea.Cmd` for async operations (database queries)
- Separate business logic from UI rendering
- Use Bubbles components for common patterns (lists, viewports)

**References**:
- https://github.com/charmbracelet/bubbletea
- https://github.com/charmbracelet/bubbles
- https://github.com/charmbracelet/lipgloss

---

## 2. SQLite Driver Selection

**Decision**: Use `modernc.org/sqlite` (pure Go SQLite implementation)

**Rationale**:
- Pure Go, no CGo required
- Cross-platform without C compiler dependencies
- Simpler builds and distribution
- Performance sufficient for personal task management workload
- Active maintenance

**Alternatives Considered**:
- **mattn/go-sqlite3**: Most popular, requires CGo, complicates cross-compilation
- **crawshaw/sqlite**: Good performance but less actively maintained
- **zombiezen.com/go/sqlite**: Another pure Go option, less mature

**Best Practices**:
- Use prepared statements for all queries
- Enable WAL mode for better concurrency: `PRAGMA journal_mode=WAL`
- Use transactions for multi-statement operations
- Close connections properly with `defer db.Close()`
- Use connection pooling (default with `database/sql`)

**Schema Design Notes**:
- Single `tasks` table with self-referencing foreign key for subtasks
- Indexes on: `column`, `priority`, `created_at`, `parent_id`
- Use TEXT for timestamp-based IDs (sortable, human-readable)
- Store tags as JSON array for simplicity (filter with JSON functions)

**References**:
- https://gitlab.com/cznic/sqlite
- https://www.sqlite.org/wal.html

---

## 3. Timestamp-Based ID Generation

**Decision**: Use format `YYYYMMDD-HHMMSS-nnnnn` (e.g., `20251104-143022-00001`)

**Rationale**:
- Sortable lexicographically (matches chronological order)
- Human-readable (can see creation date/time at a glance)
- Compact (20 characters)
- Collision-resistant with microsecond precision + counter
- No external dependencies (use standard `time` package)

**Implementation Approach**:
```go
// Pseudo-code
func GenerateID() string {
    now := time.Now()
    dateTime := now.Format("20060102-150405")
    micro := now.Nanosecond() / 1000 % 100000
    return fmt.Sprintf("%s-%05d", dateTime, micro)
}
```

**Alternatives Considered**:
- **UUID**: Not sortable, less human-readable, longer
- **Snowflake**: Requires distributed ID generation setup, overkill for single-user
- **Unix timestamp**: Less human-readable, no built-in collision resistance

**References**:
- Go `time` package documentation

---

## 4. CLI Framework

**Decision**: Use standard library `flag` package + custom command routing

**Rationale**:
- No external dependencies needed for simple CLI
- Stdlib `flag` sufficient for our command structure
- Custom routing keeps implementation simple and transparent
- Easy to add subcommands as needed

**Command Structure**:
```
ontop                    # Launch TUI (default)
ontop add "task" -p 1    # Add task
ontop list               # List tasks
ontop show <id>          # Show task details
ontop move <id> <column> # Move task
ontop filter -p 1        # Filter by priority
```

**Alternatives Considered**:
- **cobra**: Full-featured but adds complexity and dependency for simple CLI
- **cli (urfave)**: Good middle ground, but still unnecessary for our needs

**Best Practices**:
- Use subcommands pattern: `ontop <command> [args]`
- Provide `--help` for each command
- Output structured text (easy to parse) or JSON with `--json` flag
- Use exit codes: 0 (success), 1 (error), 2 (usage error)

---

## 5. Testing Strategy

**Decision**: Integration tests with real SQLite + unit tests for business logic

**Rationale**:
- TUI testing is complex; focus on CLI and data layer
- Use in-memory SQLite (`:memory:`) for fast integration tests
- Test business logic separately from UI
- Manual TUI testing for user experience validation

**Testing Approach**:
- **Unit tests**: ID generation, progress calculation, filtering logic
- **Integration tests**: CLI commands end-to-end, database operations
- **Manual testing**: TUI navigation, visual layout, keybindings

**Tools**:
- Standard `testing` package
- `testify/assert` for readable assertions
- Table-driven tests for CLI command variations

**References**:
- Go testing best practices
- https://github.com/stretchr/testify

---

## 6. Build and Distribution

**Decision**: Use `Makefile` + `go build` for local development, Go releases for distribution

**Rationale**:
- `make build`: Quick local builds
- `make test`: Run all tests
- `make install`: Install to $GOPATH/bin
- Native Go cross-compilation for multiple platforms

**Makefile Targets**:
```makefile
build:    # Build binary
test:     # Run tests
install:  # Install to system
clean:    # Remove build artifacts
```

**Cross-Compilation**:
```bash
# Linux
GOOS=linux GOARCH=amd64 go build -o ontop-linux

# macOS
GOOS=darwin GOARCH=amd64 go build -o ontop-mac

# Windows
GOOS=windows GOARCH=amd64 go build -o ontop.exe
```

**References**:
- Go build documentation
- https://golang.org/doc/install/source#environment

---

## 7. Data Storage Location

**Decision**: Store SQLite database at `$HOME/.ontop/tasks.db`

**Rationale**:
- Standard location for user data in Unix systems
- Works across Linux/macOS
- For Windows, fallback to `%USERPROFILE%\.ontop\tasks.db`
- Easy to backup/restore
- Doesn't clutter working directory

**Implementation**:
```go
// Use os.UserHomeDir() to get home directory
// Create .ontop directory if it doesn't exist
// Default path: ~/.ontop/tasks.db
// Allow override with --db-path flag for testing
```

**References**:
- XDG Base Directory Specification (optional future enhancement)

---

## 8. Vim-Style Keybinding Implementation

**Decision**: Use Bubbletea's key message handling with custom key map

**Rationale**:
- Bubbletea provides `tea.KeyMsg` with key detection
- Define key bindings in central location
- Support both single-key (j/k/h/l) and modified keys (Enter/Esc)

**Key Mapping**:
- `j`: Move down within column
- `k`: Move up within column
- `h`: Move left to previous column
- `l`: Move right to next column
- `Enter`: Expand task to detail view
- `Esc`: Close detail view, return to kanban
- `m`: Move task (prompt for column)
- `a`: Archive task
- `d`: Delete task
- `f`: Filter tasks
- `q`: Quit application

**Implementation Pattern**:
```go
type keyMap struct {
    Down  key.Binding
    Up    key.Binding
    Left  key.Binding
    Right key.Binding
    // ... more keys
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {
    case tea.KeyMsg:
        switch msg.String() {
        case "j":
            return m.moveDown()
        case "k":
            return m.moveUp()
        // ... handle other keys
        }
    }
}
```

**References**:
- Bubbletea key handling examples

---

## Summary of Decisions

| Area | Decision | Key Benefit |
|------|----------|-------------|
| TUI Framework | Bubbletea + Bubbles + Lipgloss | Clean architecture, pure Go |
| SQLite Driver | modernc.org/sqlite | No CGo, easy cross-compilation |
| ID Format | YYYYMMDD-HHMMSS-nnnnn | Sortable, human-readable |
| CLI | Standard library flag + routing | Simple, no dependencies |
| Testing | Integration + unit tests | Fast, reliable, focused |
| Build | Makefile + go build | Standard Go tooling |
| Data Location | ~/.ontop/tasks.db | Standard, portable |
| Keybindings | Custom key map | Vim-style, intuitive |

All technical decisions align with the **Simplicity First** constitution principle.
