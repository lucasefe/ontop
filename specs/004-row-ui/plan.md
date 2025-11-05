# Implementation Plan: Row-Based UI Mode

**Branch**: `004-row-ui` | **Date**: 2025-11-05 | **Spec**: [spec.md](./spec.md)
**Input**: Feature specification from `/specs/004-row-ui/spec.md`

**Note**: This template is filled in by the `/speckit.plan` command. See `.specify/templates/commands/plan.md` for the execution workflow.

## Summary

Introduce an alternative row-based visualization mode for the TUI where workflow stages (inbox, in_progress, done) are displayed as horizontal rows instead of vertical columns. Users can toggle between column and row views with a single keypress, navigate between rows with up/down arrows, navigate within rows with left/right arrows, and move tasks between rows. The preferred view mode persists between sessions in a TOML configuration file at `~/.config/ontop/ontop.toml`.

## Technical Context

**Language/Version**: Go 1.25.3
**Primary Dependencies**: Bubble Tea v1.2.4, Bubbles v0.20.0, Lip Gloss v1.0.0, modernc.org/sqlite v1.40.0
**Storage**: SQLite database at `~/.ontop/tasks.db` + TOML config at `~/.config/ontop/ontop.toml`
**Testing**: Go testing (go test), testify v1.11.1
**Target Platform**: Terminal/CLI (cross-platform: macOS, Linux, Windows)
**Project Type**: Single project (TUI + CLI dual mode)
**Performance Goals**: Sub-second view switching, 60fps rendering, responsive navigation
**Constraints**: Terminal width 80-200 characters, no external dependencies beyond existing stack
**Scale/Scope**: Personal task management tool, <1000 tasks expected, single user

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

### I. Simplicity First ✓ PASS (Re-validated)

- **Pre-Design Assessment**: Feature adds alternative visualization without introducing new abstractions or external dependencies
- **Post-Design Validation**:
  - ✓ Single new dependency: `github.com/BurntSushi/toml` v1.4.0 (mature, zero transitive deps)
  - ✓ Config module is 3 functions (Load, Save, GetConfigPath) - minimal API
  - ✓ Row rendering isolated to single new file (row.go) - no cross-cutting changes
  - ✓ ViewLayout is simple enum (2 values) - no complex state machine
  - ✓ No new abstractions (no interfaces, strategy patterns, or over-engineering)
- **Complexity Metrics**:
  - New code: ~300-400 LOC (row.go, config.go)
  - Modified code: ~100 LOC (app.go, model.go, keys.go)
  - Total impact: <500 LOC in well-isolated modules
- **Verdict**: PASS - Design maintains simplicity, single-purpose modules

### II. Vim-Style Interaction ✓ PASS (Re-validated)

- **Pre-Design Assessment**: Maintains vim-style navigation patterns with intuitive keybindings
- **Post-Design Validation**:
  - ✓ Context-sensitive navigation preserves vim semantics:
    - Column mode: j/k moves in primary (vertical) dimension
    - Row mode: h/l moves in primary (horizontal) dimension
  - ✓ Toggle key 'v' is vim-conventional (visual/view mode)
  - ✓ All existing keybindings unchanged (a, e, m, d, n, s, q, ?)
  - ✓ Help text is context-aware (shows correct keys for current layout)
- **Muscle Memory Preservation**:
  - Navigation arrows adapt to primary dimension (users navigate "toward" tasks)
  - No conflicts with existing bindings
  - Discoverable via help (?) command
- **Verdict**: PASS - Design enhances vim-style interaction without compromising it

### III. Dual-Mode Architecture ✓ PASS (Re-validated)

- **Pre-Design Assessment**: Feature affects TUI mode only, CLI remains unchanged
- **Post-Design Validation**:
  - ✓ Zero changes to CLI module (internal/cli/)
  - ✓ Zero changes to storage layer (internal/storage/)
  - ✓ Zero changes to business logic (internal/service/)
  - ✓ Both TUI layouts consume same SQLite database
  - ✓ Config is TUI preference only (CLI unaware, unnecessary)
- **Separation Verified**:
  - TUI changes: internal/tui/, internal/config/
  - CLI: Completely untouched
  - Data model: Completely untouched
- **Verdict**: PASS - Design respects dual-mode architecture perfectly

### IV. Documentation Currency ⚠️ ATTENTION REQUIRED (Re-validated)

- **Pre-Design Assessment**: README will need updates after implementation
- **Post-Design Requirements** (for tasks.md polish phase):
  - [ ] Document 'v' keybinding for view toggle in TUI controls section
  - [ ] Add row mode navigation explanation (context-sensitive arrows)
  - [ ] Document config file location (~/.config/ontop/ontop.toml) and format
  - [ ] Update screenshots if they show column layout only
  - [ ] Add quickstart section for view mode toggle
  - [ ] Document scroll indicators ([◀] [▶]) in row mode
- **Artifacts Ready**:
  - quickstart.md contains complete user documentation
  - Can be adapted for README.md in polish phase
- **Verdict**: ATTENTION REQUIRED - Documentation debt tracked, will be addressed in tasks.md

**Overall Gate Status (Post-Design)**: ✓ PASS - All principles satisfied
- Design maintains project values
- Implementation plan respects constitution
- Documentation update tracked for execution phase
- No violations or exceptions needed

## Project Structure

### Documentation (this feature)

```text
specs/004-row-ui/
├── plan.md              # This file (/speckit.plan command output)
├── research.md          # Phase 0 output (/speckit.plan command)
├── data-model.md        # Phase 1 output (/speckit.plan command)
├── quickstart.md        # Phase 1 output (/speckit.plan command)
├── checklists/
│   └── requirements.md  # Quality validation from /speckit.specify
└── tasks.md             # Phase 2 output (/speckit.tasks command - NOT created by /speckit.plan)
```

### Source Code (repository root)

```text
# Single project structure (existing)
cmd/ontop/
└── main.go              # Entry point - no changes needed

internal/
├── models/              # Task, Column - no changes needed
├── storage/             # SQLite operations - no changes needed
├── service/             # Business logic - no changes needed
├── cli/                 # CLI commands - no changes needed
└── tui/                 # TUI implementation - PRIMARY WORK AREA
    ├── model.go         # MODIFY: Add ViewLayout field (column vs row)
    ├── app.go           # MODIFY: Add 'v' key handler, view switching logic
    ├── kanban.go        # EXISTING: Column-based rendering (keep as-is)
    ├── row.go           # NEW: Row-based rendering
    ├── keys.go          # MODIFY: Add toggle key, document navigation changes
    ├── detail.go        # MODIFY: Remember view mode context
    ├── move.go          # MODIFY: Support row-based move UI
    └── styles.go        # MAY CREATE: Shared styles for row rendering

# NEW: Config management
internal/config/
└── config.go            # NEW: TOML config read/write, view mode persistence

tests/
├── integration/
│   └── tui_test.go      # MAY ADD: View switching integration tests
└── unit/
    └── config_test.go   # NEW: Config persistence unit tests

# NEW: Config file location (runtime, not in repo)
~/.config/ontop/
└── ontop.toml           # Runtime: view_mode = "column" | "row"
```

**Structure Decision**: Single project structure maintained. All changes isolated to TUI layer (`internal/tui/`) with new config package (`internal/config/`) for TOML persistence. No changes to CLI, storage, or data models.

## Complexity Tracking

> **Fill ONLY if Constitution Check has violations that must be justified**

No violations identified. All constitution principles satisfied.
