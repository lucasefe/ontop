# Tasks: OnTop - Dual-Mode Task Manager

**Input**: Design documents from `/specs/001-ontop/`
**Prerequisites**: plan.md (required), spec.md (required for user stories), research.md, data-model.md, contracts/

**Tests**: Tests are OPTIONAL - not included in this task list per specification

**Organization**: Tasks are grouped by user story to enable independent implementation and testing of each story.

## Format: `[ID] [P?] [Story] Description`

- **[P]**: Can run in parallel (different files, no dependencies)
- **[Story]**: Which user story this task belongs to (e.g., US1, US2, US3, US4)
- Include exact file paths in descriptions

## Path Conventions

- **Single project**: `internal/`, `cmd/ontop/`, `tests/` at repository root
- Paths shown below use single project structure from plan.md

---

## Phase 1: Setup (Shared Infrastructure)

**Purpose**: Project initialization and basic structure

- [x] T001 Initialize Go module with `go mod init github.com/lucasefe/ontop`
- [x] T002 [P] Create directory structure: cmd/ontop, internal/{models,storage,tui,cli,service}, tests/{integration,testdata}
- [x] T003 [P] Add dependencies: bubbletea, bubbles, lipgloss, modernc.org/sqlite, testify
- [x] T004 [P] Create Makefile with build, test, run, install, clean targets
- [x] T005 [P] Create .gitignore for Go project (bin/, *.db, vendor/)

---

## Phase 2: Foundational (Blocking Prerequisites)

**Purpose**: Core infrastructure that MUST be complete before ANY user story can be implemented

**⚠️ CRITICAL**: No user story work can begin until this phase is complete

- [x] T006 Define Task model struct in internal/models/task.go with all attributes from data-model.md
- [x] T007 Define Column constants (Inbox, InProgress, Done) in internal/models/column.go
- [x] T008 Implement timestamp-based ID generator in internal/service/id_generator.go using format YYYYMMDD-HHMMSS-nnnnn
- [x] T009 Create database connection function in internal/storage/sqlite.go with WAL mode enabled
- [x] T010 Implement database schema creation in internal/storage/migrations.go with tasks table and indexes
- [x] T011 Implement CreateTask() function in internal/storage/sqlite.go
- [x] T012 Implement GetTask() function in internal/storage/sqlite.go
- [x] T013 Implement ListTasks() function with filter support in internal/storage/sqlite.go
- [x] T014 Implement UpdateTask() function in internal/storage/sqlite.go
- [x] T015 Implement DeleteTask() function (soft delete) in internal/storage/sqlite.go
- [x] T016 Implement progress calculation logic in internal/service/task_service.go for tasks with subtasks
- [x] T017 Create database helper to get home directory path (~/.ontop/tasks.db) in internal/storage/sqlite.go

**Checkpoint**: Foundation ready - user story implementation can now begin in parallel

---

## Phase 3: User Story 1 - Quick Task Capture via CLI (Priority: P1) 🎯 MVP

**Goal**: Users can capture and manage tasks via CLI commands without leaving the terminal

**Independent Test**: Run `ontop add "Test task" -p 1`, then `ontop list` to verify task appears, then `ontop show <id>` to see details

### Implementation for User Story 1

- [x] T018 [P] [US1] Create CLI command router in cmd/ontop/main.go that parses subcommands (add, list, show)
- [x] T019 [P] [US1] Implement `ontop add` command in internal/cli/add.go with flags for --priority, --tags, --parent
- [x] T020 [P] [US1] Implement `ontop list` command in internal/cli/list.go with text output formatting
- [x] T021 [P] [US1] Implement `ontop show` command in internal/cli/show.go to display full task details and subtasks
- [x] T022 [US1] Add --json flag support to add command in internal/cli/add.go
- [x] T023 [US1] Add --json flag support to list command in internal/cli/list.go
- [x] T024 [US1] Add --json flag support to show command in internal/cli/show.go
- [x] T025 [US1] Add error handling and exit codes (0=success, 1=error, 2=validation) to all CLI commands
- [x] T026 [US1] Implement help text (--help) for add, list, and show commands

**Checkpoint**: At this point, User Story 1 should be fully functional and testable independently (MVP complete!)

---

## Phase 4: User Story 2 - Interactive Task Management via TUI (Priority: P2)

**Goal**: Users can view and organize tasks in a visual kanban interface with vim-style navigation

**Independent Test**: Run `ontop` to launch TUI, use j/k to navigate tasks, h/l to switch columns, Enter to expand task, Esc to close, verify all navigation works

### Implementation for User Story 2

- [x] T027 [P] [US2] Define Bubbletea Model struct in internal/tui/app.go with state for columns, selection, view mode
- [x] T028 [P] [US2] Define keybinding map in internal/tui/keys.go for j/k/h/l/Enter/Esc/q
- [x] T029 [P] [US2] Implement Init() function in internal/tui/app.go to load tasks from database
- [x] T030 [US2] Implement Update() function in internal/tui/app.go to handle KeyMsg events
- [x] T031 [US2] Implement renderKanban() function in internal/tui/kanban.go to display three-column layout
- [x] T032 [US2] Implement task card rendering in internal/tui/kanban.go showing description, priority, progress
- [x] T033 [US2] Implement vertical navigation (j/k) within columns in internal/tui/app.go
- [x] T034 [US2] Implement horizontal navigation (h/l) between columns in internal/tui/app.go
- [x] T035 [US2] Create detail view struct and render function in internal/tui/detail.go
- [x] T036 [US2] Implement Enter key handler to switch from kanban to detail view in internal/tui/app.go
- [x] T037 [US2] Implement Esc key handler to switch from detail to kanban view in internal/tui/app.go
- [x] T038 [US2] Implement detail view rendering with all task attributes in internal/tui/detail.go
- [x] T039 [US2] Implement subtask display in detail view in internal/tui/detail.go
- [x] T040 [US2] Add task move functionality with 'm' key in internal/tui/app.go (prompt for column)
- [x] T041 [US2] Add Lipgloss styling for colors and layout in internal/tui/kanban.go
- [x] T042 [US2] Update main.go to launch TUI when no arguments provided
- [x] T043 [US2] Handle terminal resize events in internal/tui/app.go

**Checkpoint**: At this point, User Stories 1 AND 2 should both work independently

---

## Phase 5: User Story 3 - Task Filtering and Search (Priority: P3)

**Goal**: Users can filter tasks by priority, column, tags in both CLI and TUI modes

**Independent Test**: Create tasks with different priorities/tags, run `ontop list --priority 1`, verify only P1 tasks shown; in TUI press 'f' to filter, verify filtered view works

### Implementation for User Story 3

- [ ] T044 [P] [US3] Add filter flags (--priority, --column, --tags) to list command in internal/cli/list.go
- [ ] T045 [P] [US3] Implement filter query building in internal/storage/sqlite.go for priority filtering
- [ ] T046 [P] [US3] Implement filter query building in internal/storage/sqlite.go for column filtering
- [ ] T047 [US3] Implement filter query building in internal/storage/sqlite.go for tag filtering (using JSON functions)
- [ ] T048 [US3] Add filter mode state to TUI Model in internal/tui/app.go
- [ ] T049 [US3] Implement 'f' key handler to activate filter mode in internal/tui/app.go
- [ ] T050 [US3] Create filter input UI in internal/tui/filter.go using Bubbles text input component
- [ ] T051 [US3] Apply filters to task list and update kanban display in internal/tui/app.go
- [ ] T052 [US3] Implement filter clear functionality in TUI in internal/tui/app.go
- [ ] T053 [US3] Add filter indicator to TUI status bar in internal/tui/kanban.go

**Checkpoint**: All three user stories should now be independently functional

---

## Phase 6: User Story 4 - Archive and Delete Operations (Priority: P4)

**Goal**: Users can archive completed tasks and permanently delete unwanted tasks

**Independent Test**: Archive a task with `ontop archive <id>`, verify it doesn't appear in `ontop list`, verify it appears in `ontop list --archived`; delete a task and verify it's gone

### Implementation for User Story 4

- [ ] T054 [P] [US4] Implement `ontop archive` command in internal/cli/archive.go
- [ ] T055 [P] [US4] Implement `ontop unarchive` command in internal/cli/unarchive.go
- [ ] T056 [P] [US4] Implement `ontop delete` command in internal/cli/delete.go with soft delete
- [ ] T057 [US4] Add --hard flag to delete command for permanent deletion in internal/cli/delete.go
- [ ] T058 [US4] Add confirmation prompt for hard delete with subtasks in internal/cli/delete.go
- [ ] T059 [US4] Add --archived flag to list command in internal/cli/list.go
- [ ] T060 [US4] Implement archived view filter in internal/storage/sqlite.go
- [ ] T061 [US4] Add 'a' key handler for archive in TUI in internal/tui/app.go
- [ ] T062 [US4] Add 'd' key handler for delete in TUI in internal/tui/app.go
- [ ] T063 [US4] Add confirmation dialog for delete in TUI using Bubbles in internal/tui/app.go

**Checkpoint**: All user stories should now be independently functional

---

## Phase 7: Additional CLI Commands

**Purpose**: Complete remaining CLI functionality for full dual-mode support

- [x] T064 [P] Implement `ontop move` command in internal/cli/move.go
- [x] T065 [P] Implement `ontop update` command in internal/cli/update.go with flags for description, priority, tags, progress
- [ ] T066 [P] Implement `ontop filter` command in internal/cli/filter.go (alias for list with filters)
- [ ] T067 Add validation for priority range (1-5) in all CLI commands
- [ ] T068 Add validation for column values (inbox, in_progress, done) in all CLI commands
- [ ] T069 Implement global --db-path flag in cmd/ontop/main.go

---

## Phase 8: Polish & Cross-Cutting Concerns

**Purpose**: Improvements that affect multiple user stories

- [ ] T070 [P] Add comprehensive error messages to CLI commands in internal/cli/
- [ ] T071 [P] Add help screen ('?' key) to TUI in internal/tui/help.go
- [ ] T072 [P] Create README.md with installation and usage instructions
- [ ] T073 Implement completed_at timestamp setting when moving to Done column in internal/storage/sqlite.go
- [ ] T074 Implement updated_at timestamp updating on any field change in internal/storage/sqlite.go
- [ ] T075 Add status bar to TUI showing current filter, task count in internal/tui/kanban.go
- [ ] T076 Optimize database queries with prepared statements in internal/storage/sqlite.go
- [ ] T077 [P] Add integration test for CLI add/list/show workflow in tests/integration/cli_test.go
- [ ] T078 [P] Add integration test for database CRUD operations in tests/integration/storage_test.go
- [ ] T079 Add cross-compilation instructions to Makefile for Linux/macOS/Windows
- [ ] T080 Run quickstart.md validation to ensure all commands work end-to-end

---

## Dependencies & Execution Order

### Phase Dependencies

- **Setup (Phase 1)**: No dependencies - can start immediately
- **Foundational (Phase 2)**: Depends on Setup completion - BLOCKS all user stories
- **User Stories (Phase 3-6)**: All depend on Foundational phase completion
  - User stories can then proceed in parallel (if staffed) or sequentially by priority
- **Additional CLI (Phase 7)**: Can proceed once US1 complete (shares CLI infrastructure)
- **Polish (Phase 8)**: Depends on desired user stories being complete

### User Story Dependencies

- **User Story 1 (P1)**: Can start after Foundational (Phase 2) - No dependencies on other stories ✅ MVP
- **User Story 2 (P2)**: Can start after Foundational (Phase 2) - No dependencies on US1 (uses same database layer)
- **User Story 3 (P3)**: Can start after Foundational (Phase 2) - No dependencies on US1/US2 (extends both modes)
- **User Story 4 (P4)**: Can start after Foundational (Phase 2) - No dependencies on US1/US2/US3 (extends both modes)

### Within Each User Story

- Tasks marked [P] can run in parallel (different files, no dependencies)
- Foundational tasks (Phase 2) must complete before any user story tasks
- JSON flag tasks depend on base command implementation
- Error handling tasks can run after core implementation

### Parallel Opportunities

- **Setup (Phase 1)**: All tasks marked [P] can run in parallel
- **Foundational (Phase 2)**: Tasks T006-T008 can run in parallel, then T009-T017 after T009 completes
- **User Story 1**: Tasks T018-T021 can run in parallel, then T022-T026 after core commands work
- **User Story 2**: Tasks T027-T028 can run in parallel, then T029-T043 proceed with some parallelism
- **User Story 3**: Tasks T044-T047 can run in parallel (different filters)
- **User Story 4**: Tasks T054-T057 can run in parallel (different commands)
- **All User Stories**: Once Phase 2 completes, US1, US2, US3, and US4 can all start in parallel by different developers

---

## Parallel Example: User Story 1 (MVP)

```bash
# After Phase 2 (Foundational) is complete, launch these tasks together:
Task T018: Create CLI command router (cmd/ontop/main.go)
Task T019: Implement add command (internal/cli/add.go)
Task T020: Implement list command (internal/cli/list.go)
Task T021: Implement show command (internal/cli/show.go)

# All four tasks work on different files and can proceed in parallel
```

---

## Implementation Strategy

### MVP First (User Story 1 Only)

1. Complete Phase 1: Setup (5 tasks)
2. Complete Phase 2: Foundational (12 tasks - CRITICAL PATH)
3. Complete Phase 3: User Story 1 (9 tasks)
4. **STOP and VALIDATE**: Test CLI commands end-to-end
5. You now have a working MVP: CLI task manager with add/list/show!

### Incremental Delivery

1. **MVP**: Setup + Foundational + US1 (26 tasks) → CLI task manager
2. **V1**: Add US2 (17 tasks) → TUI visual interface
3. **V2**: Add US3 (10 tasks) → Filtering capabilities
4. **V3**: Add US4 (10 tasks) → Archive and delete
5. **Final**: Add Phase 7-8 (17 tasks) → Polish and complete feature set

### Parallel Team Strategy

With multiple developers after Phase 2 completion:

1. **Developer A**: User Story 1 (CLI commands) - 9 tasks
2. **Developer B**: User Story 2 (TUI interface) - 17 tasks
3. **Developer C**: User Story 3 (Filtering) - 10 tasks
4. **Developer D**: User Story 4 (Archive/Delete) - 10 tasks

All stories complete independently and integrate through shared database layer.

---

## Notes

- [P] tasks = different files, can run in parallel
- [Story] label maps task to specific user story for traceability
- Each user story is independently completable and testable
- Foundation (Phase 2) is the critical path - must complete before any user story work
- US1 alone delivers a working MVP (CLI task manager)
- Commit after each task or logical group for incremental progress
- Stop at any checkpoint to validate story independently
- Total tasks: 80 (Setup: 5, Foundational: 12, US1: 9, US2: 17, US3: 10, US4: 10, CLI: 6, Polish: 11)
