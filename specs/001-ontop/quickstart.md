# Quickstart Guide: OnTop Development

**Date**: 2025-11-04
**Branch**: 001-ontop

## Overview

This guide helps you get started developing OnTop, a dual-mode (TUI + CLI) task manager written in Go.

---

## Prerequisites

- **Go 1.21+**: [Install Go](https://golang.org/dl/)
- **Make**: Should be pre-installed on macOS/Linux. Windows users can use WSL or skip Makefile and run go commands directly.
- **Terminal**: Any terminal with ANSI color support

---

## Initial Setup

### 1. Initialize Go Module

```bash
# From repository root
go mod init github.com/yourusername/ontop

# Add dependencies
go get github.com/charmbracelet/bubbletea@latest
go get github.com/charmbracelet/bubbles@latest
go get github.com/charmbracelet/lipgloss@latest
go get modernc.org/sqlite@latest
go get github.com/stretchr/testify@latest
```

### 2. Create Project Structure

```bash
# Create directories
mkdir -p cmd/ontop
mkdir -p internal/{models,storage,tui,cli,service}
mkdir -p tests/{integration,testdata}

# Create main entry point
touch cmd/ontop/main.go

# Create initial files
touch internal/models/task.go
touch internal/storage/sqlite.go
touch internal/tui/app.go
touch internal/cli/commands.go
touch internal/service/task_service.go
```

### 3. Create Makefile

```makefile
# Makefile
.PHONY: build test run install clean

build:
	go build -o ontop cmd/ontop/main.go

test:
	go test ./... -v

run:
	go run cmd/ontop/main.go

install:
	go install cmd/ontop/main.go

clean:
	rm -f ontop
	rm -f ~/.ontop/tasks.db

help:
	@echo "Available targets:"
	@echo "  build   - Build the ontop binary"
	@echo "  test    - Run all tests"
	@echo "  run     - Run ontop directly (TUI mode)"
	@echo "  install - Install ontop to \$$GOPATH/bin"
	@echo "  clean   - Remove binary and test database"
```

---

## Development Workflow

### Phase 1: Database Layer

**Goal**: Get SQLite working with basic CRUD operations.

**Files to create**:
1. `internal/models/task.go` - Task struct
2. `internal/storage/sqlite.go` - Database connection and queries
3. `internal/storage/migrations.go` - Schema creation
4. `internal/service/id_generator.go` - Timestamp-based ID generation

**Implementation order**:
```go
// 1. Define Task model (internal/models/task.go)
type Task struct {
    ID          string
    Description string
    Priority    int
    Column      string
    Progress    int
    ParentID    *string
    Archived    bool
    Tags        []string
    CreatedAt   time.Time
    UpdatedAt   time.Time
    CompletedAt *time.Time
    DeletedAt   *time.Time
}

// 2. Implement database connection (internal/storage/sqlite.go)
func NewDB(path string) (*sql.DB, error)
func CreateSchema(db *sql.DB) error

// 3. Basic CRUD operations
func CreateTask(db *sql.DB, task *Task) error
func GetTask(db *sql.DB, id string) (*Task, error)
func ListTasks(db *sql.DB, filters map[string]interface{}) ([]Task, error)
func UpdateTask(db *sql.DB, task *Task) error
func DeleteTask(db *sql.DB, id string) error
```

**Test**:
```bash
# Run tests
go test ./internal/storage/... -v

# Or create integration test
go test ./tests/integration/storage_test.go -v
```

---

### Phase 2: CLI Commands

**Goal**: Get basic CLI working before building TUI.

**Files to create**:
1. `cmd/ontop/main.go` - CLI entry point and routing
2. `internal/cli/add.go` - Add command
3. `internal/cli/list.go` - List command
4. `internal/cli/show.go` - Show command

**Implementation order**:
```go
// 1. Main CLI router (cmd/ontop/main.go)
func main() {
    if len(os.Args) == 1 {
        // Launch TUI (implement later)
        fmt.Println("TUI mode not implemented yet")
        return
    }

    command := os.Args[1]
    switch command {
    case "add":
        cli.AddCommand(os.Args[2:])
    case "list":
        cli.ListCommand(os.Args[2:])
    case "show":
        cli.ShowCommand(os.Args[2:])
    default:
        fmt.Fprintf(os.Stderr, "Unknown command: %s\n", command)
        os.Exit(2)
    }
}

// 2. Implement each command
// internal/cli/add.go
func AddCommand(args []string) {
    // Parse flags, create task, print result
}
```

**Test**:
```bash
# Build and test CLI
make build

./ontop add "Test task" --priority 1
./ontop list
./ontop show 20251104-143000-00001
```

---

### Phase 3: TUI - Kanban View

**Goal**: Build the main kanban board interface.

**Files to create**:
1. `internal/tui/app.go` - Bubbletea app model
2. `internal/tui/kanban.go` - Kanban board view
3. `internal/tui/keys.go` - Keybinding definitions

**Implementation order**:
```go
// 1. Define app model (internal/tui/app.go)
type Model struct {
    db          *sql.DB
    tasks       []models.Task
    columns     [3][]models.Task  // inbox, in_progress, done
    activeCol   int
    activeRow   int
    width       int
    height      int
}

func (m Model) Init() tea.Cmd {
    return loadTasks(m.db)
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {
    case tea.KeyMsg:
        return m.handleKeypress(msg)
    case tea.WindowSizeMsg:
        m.width = msg.Width
        m.height = msg.Height
    case tasksLoadedMsg:
        m.tasks = msg.tasks
        m.organizeIntoColumns()
    }
    return m, nil
}

func (m Model) View() string {
    return m.renderKanban()
}

// 2. Implement navigation (internal/tui/kanban.go)
func (m Model) handleKeypress(msg tea.KeyMsg) (Model, tea.Cmd) {
    switch msg.String() {
    case "j":  // Move down
        m.activeRow++
    case "k":  // Move up
        m.activeRow--
    case "h":  // Move left
        m.activeCol--
    case "l":  // Move right
        m.activeCol++
    case "q":  // Quit
        return m, tea.Quit
    }
    return m, nil
}
```

**Test**:
```bash
# Run TUI
./ontop

# Or during development
go run cmd/ontop/main.go
```

---

### Phase 4: TUI - Detail View

**Goal**: Add task detail expansion view.

**Files to create**:
1. `internal/tui/detail.go` - Task detail view

**Implementation**:
```go
type Model struct {
    // ... existing fields
    viewMode    string  // "kanban" or "detail"
    selectedTask *models.Task
    subtasks     []models.Task
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {
    case tea.KeyMsg:
        if m.viewMode == "kanban" {
            switch msg.String() {
            case "enter":
                return m.showDetailView()
            }
        } else if m.viewMode == "detail" {
            switch msg.String() {
            case "esc":
                return m.showKanbanView()
            }
        }
    }
}

func (m Model) View() string {
    if m.viewMode == "detail" {
        return m.renderDetail()
    }
    return m.renderKanban()
}
```

---

### Phase 5: Additional Features

Implement remaining features in priority order:

1. **Task operations in TUI**: Move, archive, delete
2. **Filtering**: By priority, tags, column
3. **Subtasks**: Create, display, progress calculation
4. **Polish**: Colors, formatting, error handling

---

## Testing Strategy

### Unit Tests

```bash
# Test individual packages
go test ./internal/models/...
go test ./internal/service/...
```

### Integration Tests

```go
// tests/integration/cli_test.go
func TestCLIAddCommand(t *testing.T) {
    // Setup test database
    db, cleanup := setupTestDB(t)
    defer cleanup()

    // Run CLI command
    output := runCLI("add", "Test task", "--priority", "1")

    // Assert task was created
    assert.Contains(t, output, "Created task:")
}
```

### Manual TUI Testing

Checklist:
- [ ] j/k navigation works in columns
- [ ] h/l navigation switches columns
- [ ] Enter expands task to detail view
- [ ] Esc closes detail view
- [ ] Task cards display correctly
- [ ] Colors render properly

---

## Common Development Tasks

### Run TUI in Development

```bash
go run cmd/ontop/main.go
```

### Build and Install

```bash
make build
make install

# Then from anywhere:
ontop
```

### Reset Database

```bash
rm ~/.ontop/tasks.db
./ontop add "First task"
```

### Add Sample Data

```bash
./ontop add "Fix bug" --priority 1 --tags bug
./ontop add "Write docs" --priority 3 --tags docs
./ontop add "Deploy" --priority 2 --tags deploy
./ontop move <id> in_progress
```

---

## Troubleshooting

### "Cannot find package"

```bash
go mod tidy
go get -u all
```

### "Database locked"

SQLite database is locked by another process:
```bash
# Close all ontop instances
ps aux | grep ontop

# Or delete and recreate database
rm ~/.ontop/tasks.db
```

### TUI Not Rendering Properly

Check terminal support:
```bash
echo $TERM
# Should be xterm-256color or similar

# Test colors
tput colors
# Should output 256
```

---

## Next Steps

Once basic functionality works:

1. **Add more CLI commands**: update, filter, archive
2. **Enhance TUI**: Better styling, help screen, search
3. **Add tests**: Increase coverage
4. **Documentation**: README, usage examples
5. **Distribution**: Cross-compile for Linux/macOS/Windows

---

## Quick Reference

### Project Structure

```
ontop/
├── cmd/ontop/main.go          # Entry point
├── internal/
│   ├── models/task.go          # Data structures
│   ├── storage/sqlite.go       # Database layer
│   ├── tui/                    # Bubbletea TUI
│   ├── cli/                    # CLI commands
│   └── service/                # Business logic
├── tests/                      # Integration tests
├── Makefile                    # Build automation
└── go.mod                      # Dependencies
```

### Essential Commands

```bash
make build      # Build binary
make test       # Run tests
make run        # Launch TUI
make install    # Install to $GOPATH/bin
make clean      # Clean build artifacts
```

### Key Files

- `internal/models/task.go` - Task entity definition
- `internal/storage/sqlite.go` - Database operations
- `internal/tui/app.go` - Bubbletea app model
- `internal/cli/commands.go` - CLI command routing
- `contracts/cli-commands.md` - CLI command specifications
- `data-model.md` - Database schema and entities

---

## Getting Help

- **Bubbletea docs**: https://github.com/charmbracelet/bubbletea
- **Bubbles components**: https://github.com/charmbracelet/bubbles
- **SQLite Go driver**: https://gitlab.com/cznic/sqlite
- **Go testing**: https://golang.org/pkg/testing/

Start with Phase 1 (Database Layer) and work through each phase sequentially. Each phase builds on the previous one and delivers working functionality you can test immediately.
