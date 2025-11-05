# Implementation Plan: OnTop - Dual-Mode Task Manager

**Branch**: `001-ontop` | **Date**: 2025-11-04 | **Spec**: [spec.md](./spec.md)
**Input**: Feature specification from `/specs/001-ontop/spec.md`

**Note**: This template is filled in by the `/speckit.plan` command. See `.specify/templates/commands/plan.md` for the execution workflow.

## Summary

Build a personal task management tool with dual-mode operation: an interactive TUI with vim-style keybindings for visual task organization across kanban columns (Inbox, In Progress, Done), and a CLI mode for scriptable automation. Tasks support subtasks, priorities (1-5), tags, filtering, and progress tracking. All data persists locally using SQLite with timestamp-based unique IDs.

## Technical Context

**Language/Version**: Go 1.21+
**Primary Dependencies**:
- Bubbletea (TUI framework)
- Bubbles (TUI components)
- Lipgloss (styling)
- SQLite driver (modernc.org/sqlite or mattn/go-sqlite3)

**Storage**: SQLite (local file-based database)
**Testing**: Go testing package + testify for assertions
**Target Platform**: Cross-platform CLI/TUI (Linux, macOS, Windows)
**Project Type**: Single standalone binary
**Performance Goals**:
- TUI keyboard response < 100ms
- CLI commands complete < 500ms for typical datasets (200 tasks)
- Database queries < 50ms for filtered lists

**Constraints**:
- Single-user local operation
- No network dependencies
- Terminal must support ANSI colors
- Offline-first (no cloud sync)

**Scale/Scope**:
- Personal use (single user)
- 50-200 active tasks typical
- 1000+ archived tasks supported
- Subtasks: 0-5 per task typical

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

### Principle I: Simplicity First

✅ **Pass**:
- Single binary, no microservices
- Direct SQLite access, no ORM complexity
- Standard Go project layout
- Minimal dependencies (Bubbletea ecosystem + SQLite)

### Principle II: Vim-Style Interaction

✅ **Pass**:
- j/k for vertical navigation
- h/l for horizontal navigation between columns
- Enter to expand, Esc to close detail views
- Modal interaction patterns

### Principle III: Dual-Mode Architecture

✅ **Pass**:
- TUI mode: Interactive full-screen interface
- CLI mode: Scriptable commands with text output
- Shared data layer (SQLite)
- No mode conflicts

**Constitution Compliance**: All principles satisfied. No violations to justify.

## Project Structure

### Documentation (this feature)

```text
specs/001-ontop/
├── spec.md              # Feature specification (complete)
├── plan.md              # This file (/speckit.plan command output)
├── research.md          # Phase 0 output (/speckit.plan command)
├── data-model.md        # Phase 1 output (/speckit.plan command)
├── quickstart.md        # Phase 1 output (/speckit.plan command)
├── contracts/           # Phase 1 output (/speckit.plan command)
│   └── cli-commands.md  # CLI command specifications
└── checklists/
    └── requirements.md  # Specification quality checklist (complete)
```

### Source Code (repository root)

```text
ontop/
├── cmd/
│   └── ontop/
│       └── main.go           # Application entry point
├── internal/
│   ├── models/
│   │   ├── task.go           # Task entity
│   │   └── column.go         # Column enum/constants
│   ├── storage/
│   │   ├── sqlite.go         # SQLite implementation
│   │   ├── migrations.go     # Schema migrations
│   │   └── queries.go        # SQL queries
│   ├── tui/
│   │   ├── app.go            # Bubbletea app model
│   │   ├── kanban.go         # Kanban board view
│   │   ├── detail.go         # Task detail view
│   │   └── keys.go           # Keybinding definitions
│   ├── cli/
│   │   ├── commands.go       # CLI command definitions
│   │   ├── add.go            # Add task command
│   │   ├── list.go           # List tasks command
│   │   ├── show.go           # Show task by ID command
│   │   ├── move.go           # Move task command
│   │   └── filter.go         # Filter tasks command
│   └── service/
│       ├── task_service.go   # Task business logic
│       └── id_generator.go   # Timestamp-based ID generation
├── tests/
│   ├── integration/
│   │   ├── cli_test.go       # CLI integration tests
│   │   └── storage_test.go   # Storage integration tests
│   └── testdata/
│       └── test.db           # Test database fixture
├── go.mod
├── go.sum
├── Makefile                  # Build automation
└── README.md                 # Project documentation
```

**Structure Decision**: Single project layout chosen. This is a standalone CLI/TUI application with no separate frontend/backend. All code lives under `internal/` following Go conventions. The `cmd/ontop/` directory contains only the main entry point. Tests are organized by type (integration tests, with potential for unit tests per package).

## Complexity Tracking

> No constitution violations. This section is empty.
