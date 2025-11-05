# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

**OnTop** is a dual-mode task management application written in Go. It supports both:
- **TUI Mode**: Interactive Kanban board using Bubble Tea framework (default when run without arguments)
- **CLI Mode**: Traditional command-line interface for scripting and automation

The binary is named `ontop` (not `todoer` - that was the original working name).

## Development Commands

### Building and Running

```bash
# Build the ontop binary
make build

# Run in TUI mode (interactive)
make run
# or just
./ontop

# Run CLI commands
./ontop add "Task description"
./ontop list
./ontop show <task-id>
./ontop move <task-id> "in_progress"
./ontop update <task-id> --priority 3

# Install to $GOPATH/bin
make install

# Clean build artifacts and database
make clean
```

### Testing

```bash
# Run all tests
make test

# Run tests with verbose output
go test ./... -v

# Run specific package tests
go test ./internal/storage -v
go test ./internal/cli -v
```

### Database

```bash
# Default database location
~/.ontop/tasks.db

# Use custom database path
./ontop --db-path /path/to/custom.db list

# Remove test database
rm -f ~/.ontop/tasks.db
```

## Architecture Overview

### Layered Architecture

The codebase follows a clean layered architecture:

```
cmd/ontop/main.go       → Entry point and routing
    ↓
internal/tui/           → Bubble Tea TUI implementation (if no args)
internal/cli/           → CLI command handlers (if args provided)
    ↓
internal/service/       → Business logic (ID generation, progress calculation)
    ↓
internal/storage/       → SQLite database operations
    ↓
internal/models/        → Data structures (Task, Column)
```

### Key Design Decisions

1. **Dual Interface Pattern**: `main.go` routes to either TUI or CLI based on presence of arguments
2. **State Machine TUI**: TUI uses view modes (Kanban, Detail, Move, Create, Edit) with Bubble Tea's Update/View pattern
3. **Stateless CLI**: Each CLI command is independent, no session state
4. **SQLite Storage**: Single database file with soft deletes (`deleted_at` timestamp)
5. **Timestamp-based IDs**: Format `YYYYMMDD-HHMMSS-nnnnn` for sortable, unique task IDs
6. **Fixed Workflow**: Three columns only: `inbox`, `in_progress`, `done`

### TUI Implementation (`internal/tui/`)

The TUI implements Bubble Tea's `tea.Model` interface:

- **Model struct**: Holds all state (tasks, current view mode, form inputs, selected indices)
- **Init()**: Returns initial `loadTasks` command
- **Update(msg)**: Routes messages based on view mode and message type
- **View()**: Renders current state using Lip Gloss styling

**View Modes**:
- `ViewModeKanban` (0): Three-column Kanban board
- `ViewModeDetail` (1): Single task details
- `ViewModeMove` (2): Column selection for moving tasks
- `ViewModeCreate` (3): Form for creating new task
- `ViewModeEdit` (4): Form for editing existing task

**Key Files**:
- `model.go`: Core Model struct and state management
- `kanban.go`: Kanban view rendering and navigation
- `detail.go`: Task detail view
- `move.go`: Column selection UI
- `app.go`: Update handlers and view routing
- `keys.go`: Keyboard input definitions

### CLI Implementation (`internal/cli/`)

Each command is a separate function that:
1. Parses flags using standard `flag` package
2. Validates input
3. Calls storage layer directly
4. Outputs results (plain text or JSON with `-json` flag)

**Commands**: `add.go`, `list.go`, `show.go`, `move.go`, `update.go`

### Storage Layer (`internal/storage/`)

**Database**: SQLite with WAL mode enabled

**Key Operations**:
- `CreateTask`, `GetTask`, `ListTasks`, `UpdateTask`, `DeleteTask` (soft delete)
- Flexible filtering via map: `{"priority": 3, "column": "inbox", "archived": false}`
- Schema auto-created on startup via `InitSchema()`

**Schema Features**:
- Single `tasks` table with 7 indexes
- Timestamps stored as RFC3339 strings
- Tags stored as JSON array
- Foreign key constraint for parent-child task hierarchy
- Constraints: priority 1-5, progress 0-100, valid column names

### Data Model (`internal/models/`)

```go
type Task struct {
    ID          string     // Unique timestamp-based ID
    Description string     // Task text
    Priority    int        // 1-5 (1=highest)
    Column      string     // inbox, in_progress, done
    Progress    int        // 0-100 percentage
    ParentID    *string    // For hierarchical tasks
    Archived    bool       // Soft delete flag
    Tags        []string   // Labels
    CreatedAt   time.Time
    UpdatedAt   time.Time
    CompletedAt *time.Time // Set when moved to "done"
    DeletedAt   *time.Time // Soft delete timestamp
}
```

## Common Development Patterns

### Adding a New CLI Command

1. Create new file in `internal/cli/` (e.g., `archive.go`)
2. Implement function signature: `func ArchiveCommand(db *sql.DB, args []string) error`
3. Parse flags and validate input
4. Call storage layer methods
5. Format and output results
6. Add case in `main.go` switch statement

### Adding a New TUI View Mode

1. Define new view mode constant in `internal/tui/model.go`
2. Add state fields to `Model` struct if needed
3. Implement keyboard handler function (e.g., `handleArchiveKeys`)
4. Add rendering logic in `View()` method
5. Add transitions from other view modes in `Update()`

### Working with the Database

- Always use parameterized queries (already done throughout)
- Use `scanTask` helper to convert rows to Task structs
- Filters are passed as `map[string]interface{}` to `ListTasks`
- Remember to update `updated_at` on all modifications
- Set `completed_at` when moving to "done" column, clear it when moving out

### ID Generation

Use `service.GenerateID()` for new tasks - don't create IDs manually. Format is `YYYYMMDD-HHMMSS-nnnnn` using UTC timestamps.

## Dependencies

- **Bubble Tea** (v1.2.4): TUI framework
- **Bubbles** (v0.20.0): TUI components
- **Lip Gloss** (v1.0.0): Terminal styling
- **SQLite** (modernc.org/sqlite v1.40.0): Pure Go SQLite driver
- **testify** (v1.11.1): Testing utilities

## Known Issues

1. **Scanner Duplication**: `scanTask` helper is duplicated between `service/task_service.go` and `storage/sqlite.go` to avoid circular dependencies. Consider extracting to shared utility package.

2. **Single-Level Hierarchy**: Progress calculation only updates immediate parent, not full tree traversal.

3. **No Concurrent Access**: No file locking for multiple instances accessing same database simultaneously.

<!-- MANUAL ADDITIONS START -->
<!-- MANUAL ADDITIONS END -->
