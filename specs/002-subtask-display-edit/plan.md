# Implementation Plan: Hierarchical Subtask Display and Management

**Branch**: `002-subtask-display-edit` | **Date**: 2025-01-04 | **Spec**: [spec.md](./spec.md)
**Input**: Feature specification from `/specs/002-subtask-display-edit/spec.md`

## Summary

This feature enhances the OnTop task manager to display subtasks hierarchically in both TUI kanban and CLI list views, with 2-character indentation and dash prefix for visual clarity. It streamlines subtask creation by allowing users to press 'n' in the detail view to create subtasks with the parent relationship automatically established. The implementation maintains dual-mode architecture (TUI + CLI) with consistent behavior across both interfaces and supports keyboard navigation through the hierarchical task structure.

## Technical Context

**Language/Version**: Go 1.25.3+
**Primary Dependencies**:
- Charmbracelet Bubble Tea (v1.2.4) - TUI framework
- Charmbracelet Bubbles (v0.20.0) - UI components
- Charmbracelet Lipgloss (v1.0.0) - Terminal styling
- modernc.org/sqlite (v1.40.0) - SQLite driver

**Storage**: SQLite (local file: `~/.ontop/tasks.db`) with WAL mode
**Testing**: Go standard `testing` package + Testify (v1.11.1)
**Target Platform**: Cross-platform CLI/TUI (Linux, macOS, Windows)
**Project Type**: Dual-mode (TUI + CLI) - single binary with interactive and scriptable modes
**Performance Goals**: Sub-100ms response for kanban rendering with 100+ tasks, instant keyboard navigation
**Constraints**: Maintain existing data model compatibility, preserve vim-style keybindings, no breaking CLI changes
**Scale/Scope**: Personal task manager (~1000 tasks expected), single-user local database

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

### ✅ I. Simplicity First
**Status**: PASS

- Feature uses existing data model (ParentID field already exists in Task struct)
- No new abstractions or external dependencies required
- Enhances existing views (kanban, detail, list) rather than adding new modes
- Implementation focuses on display logic and navigation—straightforward enhancements

**Justification**: Subtask hierarchy is a natural extension of existing parent-child relationships. The display changes are localized to rendering functions without introducing complex new patterns.

### ✅ II. Vim-Style Interaction
**Status**: PASS

- Reuses existing j/k navigation keys for hierarchical movement
- Press 'n' in detail view aligns with "new task" mental model
- No new keybindings required—extends existing behavior naturally
- Maintains keyboard-only workflow without mouse dependency

**Justification**: The 'n' key for subtask creation from detail view follows vim's principle of contextual commands. Navigation through indented lists uses familiar j/k patterns.

### ✅ III. Dual-Mode Architecture
**Status**: PASS

- TUI Mode: Hierarchical kanban display with indentation and navigation
- CLI Mode: `list` and `show` commands will display subtasks with indentation
- Both modes use same underlying storage and data model
- No interactive prompts in CLI mode (maintains scriptability)

**Justification**: Feature specification explicitly requires both TUI and CLI support. Implementation maintains clean separation—TUI gets interactive navigation, CLI gets formatted output for scripting.

**GATE RESULT**: ✅ **All checks passed** - No complexity violations. Feature aligns with all constitutional principles.

## Project Structure

### Documentation (this feature)

```text
specs/002-subtask-display-edit/
├── spec.md              # Feature specification ✓ Complete
├── plan.md              # This file (Phase 0-1 output)
├── research.md          # Phase 0: Technical decisions and patterns
├── data-model.md        # Phase 1: Data structures and relationships
├── quickstart.md        # Phase 1: Developer guide for implementation
├── contracts/           # Phase 1: Internal API contracts (if needed)
├── checklists/
│   └── requirements.md  # Spec quality checklist ✓ Complete
└── tasks.md             # Phase 2: Implementation tasks (/speckit.tasks)
```

### Source Code (repository root)

```text
/Users/lucasefe/Code/mine/todoer/
├── cmd/
│   └── ontop/
│       └── main.go                 # No changes expected
├── internal/
│   ├── cli/                        # 🔧 CLI command updates
│   │   ├── list.go                 # Add hierarchical display
│   │   └── show.go                 # Add subtask listing
│   ├── models/                     # No changes (ParentID exists)
│   │   └── task.go
│   ├── service/                    # 🔧 New helper functions
│   │   └── hierarchy.go            # NEW: Subtask grouping logic
│   ├── storage/                    # No changes expected
│   │   └── sqlite.go
│   └── tui/                        # 🔧 Major updates
│       ├── model.go                # Update GetTasksByColumn for hierarchy
│       ├── app.go                  # Add 'n' key handler in detail view
│       ├── kanban.go               # 🔧 Hierarchical rendering
│       ├── detail.go               # Add 'n' key support
│       ├── form.go                 # Add parentID pre-fill parameter
│       └── keys.go                 # No changes (reuse existing 'n' key)
├── tests/
│   ├── integration/                # NEW: Hierarchical display tests
│   └── unit/                       # NEW: Hierarchy service tests
└── Makefile                        # No changes
```

**Structure Decision**: Single project structure maintained. Changes localized to internal/tui (hierarchical display), internal/cli (formatted output), and a new internal/service/hierarchy.go for shared subtask grouping logic used by both modes.

## Complexity Tracking

> **No violations detected** - Table not needed. All constitutional gates passed without requiring complexity justifications.

---

## Phase 0: Research & Technical Decisions

### Research Areas

1. **Hierarchical Rendering in Bubble Tea**
   - How to maintain selection state when task list changes (parent + subtasks)
   - Best practice for rendering nested items in lipgloss-styled views
   - Handling dynamic list sizing with indented items

2. **Navigation Through Hierarchical Lists**
   - Cursor position tracking when list includes both parents and subtasks
   - Index mapping between flat display list and underlying data structures
   - Edge cases: navigating from last subtask to next parent

3. **CLI Formatting Patterns**
   - ANSI indentation best practices for terminal output
   - Maintaining table-like alignment with variable indentation
   - Ensuring scriptable output (no color codes when piped)

4. **Form Pre-population Pattern**
   - Passing context (parent ID) to initCreateForm without breaking existing usage
   - Optional parameters in Go for form initialization
   - Preserving form state when user navigates back without saving

### Key Decisions to Document

1. **Hierarchical List Structure**: How to represent parent+subtasks in memory for rendering
   - Option A: Flatten on-demand during render (keeps model simple)
   - Option B: Pre-compute hierarchical structure (faster rendering)
   - **Decision needed**: Trade-off between render performance and code complexity

2. **Selection Index Strategy**: How to track selected task in hierarchical view
   - Option A: Flat index (0, 1, 2... across all displayed items)
   - Option B: Two-level index (parentIndex, subtaskIndex or -1)
   - **Decision needed**: Impact on navigation code complexity

3. **Subtask Ordering in Same Column**: When parent and subtask in same column
   - Option A: Subtasks always immediately follow parent (breaks sort)
   - Option B: Respect sort order, then group subtasks (more complex)
   - **Decision needed**: User experience vs. implementation simplicity

4. **CLI Indentation Detection**: How CLI knows when to indent
   - Option A: Check ParentID != nil → indent
   - Option B: Build hierarchy map first, then indent based on structure
   - **Decision needed**: Performance vs. accuracy for complex hierarchies

---

## Phase 1: Design & Contracts

### Deliverables

1. **data-model.md**
   - Existing Task struct (no changes)
   - HierarchicalTask struct (display wrapper)
   - TaskHierarchy helper functions
   - State transitions (none new, reuse existing column moves)

2. **contracts/** (Internal APIs)
   - `hierarchy.go` function signatures
   - `model.go` updated GetTasksByColumn signature
   - `form.go` initCreateForm optional parameter
   - `detail.go` handleDetailKeys 'n' key behavior

3. **quickstart.md**
   - Architecture overview (flat storage → hierarchical display)
   - Key functions to understand (BuildHierarchy, RenderWithIndentation)
   - Testing approach (unit tests for hierarchy logic, integration for navigation)
   - Common gotchas (selection index math, subtask ordering edge cases)

---

## Phase 2: Task Breakdown (Generated by /speckit.tasks)

**Note**: Phase 2 is not executed by `/speckit.plan`. Run `/speckit.tasks` after completing Phase 0-1 to generate implementation tasks.

The task breakdown will include:
- P1 tasks for core hierarchical display (TUI kanban + CLI list)
- P1 tasks for 'n' key subtask creation from detail view
- P2 tasks for navigation through hierarchy
- P2 tasks for CLI parity (show command subtask listing)
- Test tasks for each component

---

## Implementation Notes

### Key Principles

1. **Preserve Existing Behavior**: All current functionality must continue working
2. **Vim-Style Consistency**: Hierarchical navigation should feel natural to vim users
3. **Dual-Mode Parity**: TUI and CLI must show consistent hierarchy representation
4. **Simplicity**: Favor straightforward rendering over complex optimization

### Risk Mitigation

- **Risk**: Selection index off-by-one errors in hierarchical navigation
  - **Mitigation**: Comprehensive unit tests for index math, integration tests for j/k navigation

- **Risk**: Performance degradation with many subtasks
  - **Mitigation**: Benchmark rendering with 100+ tasks, optimize only if needed (YAGNI)

- **Risk**: Breaking existing task operations (move, edit, delete)
  - **Mitigation**: Ensure all operations work on selected task regardless of hierarchy level

### Testing Strategy

1. **Unit Tests**
   - Hierarchy building logic (parent grouping, ordering)
   - Index mapping functions (flat → hierarchical, hierarchical → flat)
   - Indentation formatting (CLI and TUI)

2. **Integration Tests**
   - End-to-end navigation through hierarchical kanban
   - Subtask creation from detail view (full flow)
   - CLI list command with various task/subtask combinations

3. **Manual Testing**
   - Large task lists (50+ tasks with subtasks)
   - Edge cases (empty parents, orphaned subtasks, deep nesting)
   - Keyboard navigation feel (vim-style natural flow)

---

## Next Steps

1. ✅ **Phase 0 Complete**: Run research, fill `research.md` with technical decisions
2. ⏳ **Phase 1 Pending**: Generate `data-model.md`, `contracts/`, `quickstart.md`
3. ⏳ **Phase 2 Pending**: Run `/speckit.tasks` to generate implementation breakdown
