# Tasks: Row-Based UI Mode

**Input**: Design documents from `/specs/004-row-ui/`
**Prerequisites**: plan.md, spec.md, research.md, data-model.md, contracts/

**Tests**: No test tasks included (not requested in specification)

**Organization**: Tasks are grouped by user story to enable independent implementation and testing of each story.

## Format: `[ID] [P?] [Story] Description`

- **[P]**: Can run in parallel (different files, no dependencies)
- **[Story]**: Which user story this task belongs to (e.g., US1, US2, US3)
- Include exact file paths in descriptions

## Path Conventions

This is a single Go project with structure:
- Repository root: `/Users/lucasefe/Code/mine/ontop/`
- Source: `internal/`, `cmd/`
- Tests: `tests/` (unit and integration)

---

## Phase 1: Setup (Shared Infrastructure)

**Purpose**: Add TOML dependency and prepare config infrastructure

- [x] T001 Add github.com/BurntSushi/toml v1.4.0 dependency to go.mod
- [x] T002 Create internal/config/ directory for config management package

---

## Phase 2: Foundational (Blocking Prerequisites)

**Purpose**: Core config infrastructure and model extensions that all user stories depend on

**⚠️ CRITICAL**: No user story work can begin until this phase is complete

- [x] T003 [P] Implement Config struct and TOML schema in internal/config/config.go
- [x] T004 [P] Implement Load() function with defaults and error handling in internal/config/config.go
- [x] T005 [P] Implement Save() function with atomic write in internal/config/config.go
- [x] T006 [P] Implement GetConfigPath() function in internal/config/config.go
- [x] T007 Add ViewLayout enum type to internal/tui/model.go
- [x] T008 Add viewLayout field to Model struct in internal/tui/model.go
- [x] T009 Add rowScrollOffset map field to Model struct in internal/tui/model.go
- [x] T010 Update NewModel() to load config and initialize viewLayout in internal/tui/model.go

**Checkpoint**: Foundation ready - user story implementation can now begin in parallel

---

## Phase 3: User Story 1 - Switch to Row-Based View (Priority: P1) 🎯 MVP

**Goal**: Users can toggle between column and row layouts with 'v' key, see tasks displayed in horizontal rows

**Independent Test**: Launch TUI, press 'v' to toggle, verify tasks display in three horizontal rows (inbox, in_progress, done) with clear separation, press 'v' again to return to column mode

### Implementation for User Story 1

- [x] T011 [P] [US1] Create internal/tui/row.go with package declaration and imports
- [x] T012 [P] [US1] Extract renderTaskCard() helper function to share styling between column and row layouts in internal/tui/kanban.go
- [x] T013 [US1] Implement renderKanbanRows() function with basic row structure in internal/tui/row.go
- [x] T014 [US1] Implement row header rendering (workflow stage names) in internal/tui/row.go
- [x] T015 [US1] Implement horizontal task card layout within rows in internal/tui/row.go
- [x] T016 [US1] Add Lip Gloss styling for row containers (border, padding) in internal/tui/row.go
- [x] T017 [US1] Add Lip Gloss styling for selected row highlighting in internal/tui/row.go
- [x] T018 [US1] Add 'v' key binding to KeyMap in internal/tui/keys.go
- [x] T019 [US1] Implement handleToggleView() function in internal/tui/app.go
- [x] T020 [US1] Add toggle handler to Update() switch statement in internal/tui/app.go
- [x] T021 [US1] Update View() to route to renderKanbanRows() when viewLayout is LayoutRow in internal/tui/app.go
- [x] T022 [US1] Implement config save on toggle (async) in handleToggleView() in internal/tui/app.go
- [x] T023 [US1] Reset rowScrollOffset map when toggling layouts in handleToggleView() in internal/tui/app.go
- [x] T024 [US1] Update context-sensitive help text to show 'v' key in internal/tui/keys.go
- [x] T025 [US1] Handle empty rows with placeholder message in renderKanbanRows() in internal/tui/row.go

**Checkpoint**: At this point, User Story 1 should be fully functional - users can toggle view mode and see row-based layout

---

## Phase 4: User Story 2 - Navigate Between Rows (Priority: P2)

**Goal**: Users can navigate between rows with up/down arrows and between tasks within a row with left/right arrows

**Independent Test**: Switch to row mode, use up/down to move between rows (verify row highlighting), use left/right to move between tasks within a row (verify task selection highlighting)

### Implementation for User Story 2

- [x] T026 [US2] Update handleKanbanKeys() to check viewLayout for context-sensitive navigation in internal/tui/app.go
- [x] T027 [US2] Implement up/down arrow handling for row mode (switch between rows) in handleKanbanKeys() in internal/tui/app.go
- [x] T028 [US2] Implement left/right arrow handling for row mode (navigate within row) in handleKanbanKeys() in internal/tui/app.go
- [ ] T029 [US2] Add viewport calculation logic (visible task count based on terminal width) in internal/tui/row.go
- [ ] T030 [US2] Implement horizontal scrolling logic (update rowScrollOffset on left/right) in handleKanbanKeys() in internal/tui/app.go
- [ ] T031 [US2] Add scroll indicators ([◀] [▶]) rendering in renderKanbanRows() in internal/tui/row.go
- [ ] T032 [US2] Implement bounds checking for scroll offset (clamp to valid range) in handleKanbanKeys() in internal/tui/app.go
- [ ] T033 [US2] Update selected task highlighting in row mode in renderKanbanRows() in internal/tui/row.go
- [ ] T034 [US2] Implement task selection preservation when toggling layouts in handleToggleView() in internal/tui/app.go
- [ ] T035 [US2] Add getCurrentTaskID() helper method to Model in internal/tui/model.go
- [ ] T036 [US2] Add findTaskIndexByID() helper method to Model in internal/tui/model.go
- [ ] T037 [US2] Update help text to show context-sensitive navigation keys (column vs row mode) in internal/tui/keys.go

**Checkpoint**: At this point, User Stories 1 AND 2 should both work - users can toggle view and navigate in row mode

---

## Phase 5: User Story 3 - Move Tasks Between Rows (Priority: P3)

**Goal**: Users can move tasks between workflow stage rows using 'm' command, with row-based move UI

**Independent Test**: In row mode, select a task, press 'm', use up/down to select destination row, press Enter to confirm move, verify task appears in new row with correct workflow status

### Implementation for User Story 3

- [ ] T038 [US3] Update renderMove() to check viewLayout in internal/tui/move.go
- [ ] T039 [US3] Implement row-based move UI (display rows instead of columns) in renderMove() in internal/tui/move.go
- [ ] T040 [US3] Add up/down arrow handling for move mode in row layout in handleMoveKeys() in internal/tui/app.go
- [ ] T041 [US3] Update move destination highlighting for row mode in renderMove() in internal/tui/move.go
- [ ] T042 [US3] Ensure current row is indicated and unavailable during move in renderMove() in internal/tui/move.go
- [ ] T043 [US3] Verify task move operations work identically in row mode (workflow status updates) in handleMoveKeys() in internal/tui/app.go
- [ ] T044 [US3] Preserve scroll offset after task move in row mode in handleMoveKeys() in internal/tui/app.go
- [ ] T045 [US3] Update help text for move mode to reflect row navigation in internal/tui/keys.go

**Checkpoint**: All user stories should now be independently functional - full row mode feature complete

---

## Phase 6: Polish & Cross-Cutting Concerns

**Purpose**: Edge cases, error handling, documentation, and refinements

- [ ] T046 [P] Handle narrow terminal widths gracefully (minimum width warning or fallback) in renderKanbanRows() in internal/tui/row.go
- [ ] T047 [P] Handle wide terminal widths (adjust visible task count) in renderKanbanRows() in internal/tui/row.go
- [ ] T048 [P] Add error logging for config save failures in handleToggleView() in internal/tui/app.go
- [ ] T049 [P] Verify detail view preserves viewLayout context when returning to Kanban in internal/tui/detail.go
- [ ] T050 Update README.md to document 'v' keybinding for view toggle
- [ ] T051 Update README.md to add row mode navigation explanation (context-sensitive arrows)
- [ ] T052 Update README.md to document config file location (~/.config/ontop/ontop.toml)
- [ ] T053 Update README.md to document scroll indicators ([◀] [▶]) in row mode
- [ ] T054 Update README.md with quickstart section for view mode toggle (adapt from specs/004-row-ui/quickstart.md)
- [ ] T055 [P] Consider adding screenshots of row mode to README.md (optional)
- [ ] T056 Verify performance: measure view toggle latency (target <1 second, max acceptable 500ms)
- [ ] T057 Verify performance: measure navigation responsiveness in row mode (target <16ms per frame)
- [ ] T058 Run manual testing against all acceptance scenarios from spec.md
- [ ] T059 Verify edge cases: empty rows, terminal resize, view toggle with selected task
- [ ] T060 Run quickstart.md validation - execute all examples and scenarios

---

## Dependencies & Execution Order

### Phase Dependencies

- **Setup (Phase 1)**: No dependencies - can start immediately
- **Foundational (Phase 2)**: Depends on Setup completion - BLOCKS all user stories
- **User Stories (Phase 3-5)**: All depend on Foundational phase completion
  - User stories CAN proceed in parallel (if staffed) after Phase 2
  - Or sequentially in priority order (P1 → P2 → P3)
- **Polish (Phase 6)**: Depends on all user stories being complete

### User Story Dependencies

- **User Story 1 (P1)**: Depends only on Foundational (Phase 2) - No dependencies on other stories
- **User Story 2 (P2)**: Depends only on Foundational (Phase 2) - Builds on US1 but independently testable
- **User Story 3 (P3)**: Depends only on Foundational (Phase 2) - Requires US1 and US2 rendering/navigation to be useful

**Note**: While US2 and US3 build on US1 functionality, they are designed to be independently testable. US1 provides the visual foundation, but US2 adds navigation (testable independently), and US3 adds move operations (testable independently).

### Within Each User Story

- Tasks within a story marked [P] can run in parallel
- Tasks without [P] should be executed sequentially (dependencies exist)
- Complete all tasks in a user story before validating that story's checkpoint

### Parallel Opportunities

- **Phase 1**: T001 and T002 can run in parallel
- **Phase 2**: T003-T006 (config package) can run in parallel, T007-T009 (model extensions) can run in parallel
- **User Story 1**: T011-T012 can start in parallel, T016-T017 (styling) can run in parallel, T024-T025 can run in parallel
- **User Story 2**: No explicit parallel markers (most tasks have sequential dependencies)
- **User Story 3**: No explicit parallel markers (tasks build on each other)
- **Phase 6**: T046-T049 (edge cases) can run in parallel, T050-T054 (documentation) can run in parallel, T055 can run anytime

**Between User Stories**: If team capacity allows, US1, US2, and US3 CAN be worked on in parallel after Phase 2 completes, though sequential execution (P1→P2→P3) is recommended for single developer.

---

## Parallel Example: Phase 2 (Foundational)

```bash
# Launch config package tasks in parallel:
Task: "Implement Config struct and TOML schema in internal/config/config.go"
Task: "Implement Load() function with defaults and error handling in internal/config/config.go"
Task: "Implement Save() function with atomic write in internal/config/config.go"
Task: "Implement GetConfigPath() function in internal/config/config.go"

# Then launch model extension tasks in parallel:
Task: "Add ViewLayout enum type to internal/tui/model.go"
Task: "Add viewLayout field to Model struct in internal/tui/model.go"
Task: "Add rowScrollOffset map field to Model struct in internal/tui/model.go"
```

---

## Parallel Example: User Story 1

```bash
# Launch initial file creation in parallel:
Task: "Create internal/tui/row.go with package declaration and imports"
Task: "Extract renderTaskCard() helper function to share styling between column and row layouts in internal/tui/kanban.go"

# Later, launch styling tasks in parallel:
Task: "Add Lip Gloss styling for row containers (border, padding) in internal/tui/row.go"
Task: "Add Lip Gloss styling for selected row highlighting in internal/tui/row.go"

# Near end, launch final touches in parallel:
Task: "Update context-sensitive help text to show 'v' key in internal/tui/keys.go"
Task: "Handle empty rows with placeholder message in renderKanbanRows() in internal/tui/row.go"
```

---

## Parallel Example: Phase 6 (Polish)

```bash
# Launch edge case handling in parallel:
Task: "Handle narrow terminal widths gracefully (minimum width warning or fallback) in renderKanbanRows() in internal/tui/row.go"
Task: "Handle wide terminal widths (adjust visible task count) in renderKanbanRows() in internal/tui/row.go"
Task: "Add error logging for config save failures in handleToggleView() in internal/tui/app.go"
Task: "Verify detail view preserves viewLayout context when returning to Kanban in internal/tui/detail.go"

# Launch documentation tasks in parallel:
Task: "Update README.md to document 'v' keybinding for view toggle"
Task: "Update README.md to add row mode navigation explanation (context-sensitive arrows)"
Task: "Update README.md to document config file location (~/.config/ontop/ontop.toml)"
Task: "Update README.md to document scroll indicators ([◀] [▶]) in row mode"
```

---

## Implementation Strategy

### MVP First (User Story 1 Only)

1. Complete Phase 1: Setup (T001-T002)
2. Complete Phase 2: Foundational (T003-T010) - CRITICAL
3. Complete Phase 3: User Story 1 (T011-T025)
4. **STOP and VALIDATE**: Test User Story 1 independently
   - Launch TUI
   - Press 'v' to toggle to row mode
   - Verify three rows displayed (inbox, in_progress, done)
   - Verify tasks show horizontally within rows
   - Press 'v' again to return to column mode
   - Restart app and verify preference persisted
5. Deploy/demo if ready (MVP delivered!)

**MVP Scope**: 25 tasks (T001-T025)
**Estimated LOC**: ~250-300 lines (config package + row rendering + toggle logic)
**Value**: Alternative visualization mode that may work better in certain terminals

### Incremental Delivery

1. **Foundation**: Complete Phase 1 + Phase 2 → Config infrastructure ready (10 tasks)
2. **MVP**: Add User Story 1 → Test independently → Deploy (15 tasks, total 25)
   - Users can toggle to row mode and see alternative layout
3. **Enhanced**: Add User Story 2 → Test independently → Deploy (12 tasks, total 37)
   - Users can navigate in row mode with context-sensitive arrows
4. **Complete**: Add User Story 3 → Test independently → Deploy (8 tasks, total 45)
   - Users can move tasks between rows (feature-complete)
5. **Polished**: Complete Phase 6 → Final validation → Deploy (15 tasks, total 60)
   - Edge cases handled, documentation complete, performance verified

Each increment adds measurable value without breaking previous functionality.

### Parallel Team Strategy

With multiple developers (unlikely for personal project, but documented for completeness):

1. Team completes Phase 1 + Phase 2 together (10 tasks)
2. Once Foundational is done:
   - Developer A: User Story 1 (T011-T025) - 15 tasks
   - Developer B: User Story 2 (T026-T037) - 12 tasks
   - Developer C: User Story 3 (T038-T045) - 8 tasks
3. Stories complete and integrate independently
4. Team reconvenes for Phase 6 (Polish) - 15 tasks in parallel

**Note**: Sequential execution recommended for solo developer following priority order.

---

## Task Count Summary

- **Phase 1 (Setup)**: 2 tasks
- **Phase 2 (Foundational)**: 8 tasks
- **Phase 3 (US1 - P1)**: 15 tasks
- **Phase 4 (US2 - P2)**: 12 tasks
- **Phase 5 (US3 - P3)**: 8 tasks
- **Phase 6 (Polish)**: 15 tasks

**Total**: 60 tasks

### Tasks by Category

- Setup/Infrastructure: 10 tasks (Phases 1-2)
- User Story Implementation: 35 tasks (Phases 3-5)
- Polish/Documentation: 15 tasks (Phase 6)

### Parallelization Opportunities

- **Explicit [P] markers**: 16 tasks can run in parallel with other tasks
- **Phase-level parallelism**: User Stories 1-3 can start together after Phase 2
- **Recommended approach**: Sequential execution (P1→P2→P3) for single developer

---

## Independent Test Criteria

### User Story 1 (P1) - Switch to Row-Based View

**Test Steps**:
1. Build and run: `make build && ./ontop`
2. Press 'v' key
3. Verify display switches to three horizontal rows (inbox, in_progress, done)
4. Verify tasks within each row are displayed horizontally with visual separation
5. Press 'v' again
6. Verify display switches back to original column layout
7. Press 'v' to row mode again
8. Exit and restart app
9. Verify app starts in row mode (preference persisted)

**Success Criteria**: All acceptance scenarios pass, view toggle completes in <1 second

### User Story 2 (P2) - Navigate Between Rows

**Test Steps**:
1. Build and run: `make build && ./ontop`
2. Press 'v' to switch to row mode
3. Press 'j' or down arrow
4. Verify cursor moves to next row and row is visually highlighted
5. Press 'k' or up arrow
6. Verify cursor moves to previous row
7. Press 'l' or right arrow
8. Verify cursor moves to next task within row
9. Press 'h' or left arrow
10. Verify cursor moves to previous task within row
11. If row has many tasks, verify scroll indicators ([◀] [▶]) appear
12. Navigate left/right and verify viewport scrolls

**Success Criteria**: All acceptance scenarios pass, navigation feels natural (context-sensitive)

### User Story 3 (P3) - Move Tasks Between Rows

**Test Steps**:
1. Build and run: `make build && ./ontop`
2. Press 'v' to switch to row mode
3. Select a task in inbox row
4. Press 'm' to initiate move
5. Verify move UI shows rows (not columns)
6. Use 'j'/'k' to select destination row (in_progress)
7. Press Enter to confirm
8. Verify task appears in in_progress row
9. Verify task no longer in inbox row
10. Move task from in_progress to done
11. Verify completed_at timestamp is set
12. Move task back to inbox
13. Verify completed_at is cleared

**Success Criteria**: All acceptance scenarios pass, move operations identical to column mode behavior

---

## Notes

- [P] tasks = different files or independent work, no dependencies
- [Story] label maps task to specific user story for traceability
- Each user story should be independently completable and testable
- Commit after each task or logical group
- Stop at any checkpoint to validate story independently
- All file paths are relative to repository root: `/Users/lucasefe/Code/mine/ontop/`
- Go convention: packages are lowercase (config, tui)
- No test tasks included (not requested in specification)
- Constitution principles validated: Simplicity First, Vim-Style Interaction, Dual-Mode Architecture, Documentation Currency

---

## Validation Checklist

Before considering this feature complete:

- [ ] All 60 tasks completed
- [ ] User Story 1 independently tested and passing
- [ ] User Story 2 independently tested and passing
- [ ] User Story 3 independently tested and passing
- [ ] All edge cases handled (narrow terminals, empty rows, terminal resize)
- [ ] Performance targets met (sub-second toggle, responsive navigation)
- [ ] Documentation updated in README.md
- [ ] Config file persistence verified (~/.config/ontop/ontop.toml)
- [ ] All acceptance scenarios from spec.md verified
- [ ] quickstart.md examples validated
- [ ] No regressions in column mode (existing functionality unchanged)
- [ ] Code follows project conventions and Go best practices
