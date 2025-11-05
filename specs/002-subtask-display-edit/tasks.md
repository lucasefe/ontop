# Tasks: Hierarchical Subtask Display and Management

**Input**: Design documents from `/specs/002-subtask-display-edit/`
**Prerequisites**: plan.md, spec.md, research.md, data-model.md, contracts/, quickstart.md

**Tests**: Manual testing checklist included in Polish phase. No automated test tasks requested in specification.

**Organization**: Tasks are grouped by user story to enable independent implementation and testing of each story.

## Format: `[ID] [P?] [Story] Description`

- **[P]**: Can run in parallel (different files, no dependencies)
- **[Story]**: Which user story this task belongs to (e.g., US1, US2, US3, US4)
- Include exact file paths in descriptions

## Path Conventions

OnTop uses Go single project structure at repository root:
- `internal/` - Internal packages (tui, cli, models, service, storage)
- `cmd/ontop/` - Main application entry point
- `tests/` - Integration and unit tests
- Paths shown below follow this structure from plan.md

---

## Phase 1: Setup (Shared Infrastructure)

**Purpose**: No special setup needed - existing project structure supports this feature

- [ ] T001 Review existing project structure in /Users/lucasefe/Code/mine/todoer/ to understand codebase
- [ ] T002 Verify Go 1.25.3+ is installed and `make build` succeeds
- [ ] T003 [P] Review research.md decisions for implementation approach
- [ ] T004 [P] Review data-model.md for HierarchicalTask structure
- [ ] T005 [P] Review contracts/ for API changes to existing functions

**Checkpoint**: Development environment ready, design documents understood

---

## Phase 2: Foundational (Blocking Prerequisites)

**Purpose**: Core hierarchy service that ALL user stories depend on

**⚠️ CRITICAL**: No user story work can begin until this phase is complete

- [ ] T006 Create internal/service/hierarchy.go with package declaration and imports
- [ ] T007 Define HierarchicalTask struct in internal/service/hierarchy.go
- [ ] T008 Implement BuildFlatHierarchy() function in internal/service/hierarchy.go
- [ ] T009 Implement GetSubtasksForParent() function in internal/service/hierarchy.go
- [ ] T010 Implement IsOrphanedSubtask() function in internal/service/hierarchy.go
- [ ] T011 Implement FormatWithIndentation() function in internal/service/hierarchy.go
- [ ] T012 Create internal/service/hierarchy_test.go with unit tests for BuildFlatHierarchy
- [ ] T013 Add unit test TestBuildFlatHierarchy_EmptyInput in internal/service/hierarchy_test.go
- [ ] T014 [P] Add unit test TestBuildFlatHierarchy_OnlyParents in internal/service/hierarchy_test.go
- [ ] T015 [P] Add unit test TestBuildFlatHierarchy_MixedParentsAndSubtasks in internal/service/hierarchy_test.go
- [ ] T016 [P] Add unit test TestBuildFlatHierarchy_OrphanedSubtask in internal/service/hierarchy_test.go
- [ ] T017 Run `make test` and ensure all hierarchy service tests pass before proceeding

**Checkpoint**: Foundation ready - hierarchy service fully functional and tested, user story implementation can now begin

---

## Phase 3: User Story 1 - View Subtasks in Kanban Column (Priority: P1) 🎯 MVP

**Goal**: Display parent tasks with their subtasks visually indented in the TUI kanban board

**Independent Test**: Create parent task with 2-3 subtasks in same column, launch TUI, verify subtasks appear indented below parent with dash prefix

### Implementation for User Story 1

- [ ] T018 [US1] Modify GetTasksByColumn() in internal/tui/model.go to call hierarchy.BuildFlatHierarchy()
- [ ] T019 [US1] Update GetTasksByColumn() in internal/tui/model.go to unwrap HierarchicalTask and return []*models.Task
- [ ] T020 [US1] Add isSubtask parameter to renderTaskCard() signature in internal/tui/kanban.go
- [ ] T021 [US1] Update renderTaskCard() implementation in internal/tui/kanban.go to add dash prefix when isSubtask=true
- [ ] T022 [US1] Adjust title truncation logic in renderTaskCard() in internal/tui/kanban.go to account for dash prefix width
- [ ] T023 [US1] Update renderColumn() in internal/tui/kanban.go to pass isSubtask parameter (check task.ParentID != nil)
- [ ] T024 [US1] Build project with `make build` and verify no compilation errors
- [ ] T025 [US1] Manual test: Create parent task in inbox, create 2 subtasks in inbox, verify hierarchy displays correctly
- [ ] T026 [US1] Manual test: Create parent in inbox, subtask in progress, verify subtask appears in progress column (not under parent)
- [ ] T027 [US1] Manual test: Navigate through tasks with j/k keys, verify selection works correctly with hierarchy

**Checkpoint**: User Story 1 complete and independently testable - kanban shows hierarchical subtask display

---

## Phase 4: User Story 2 - Create Subtasks from Detail View (Priority: P1) 🎯 MVP

**Goal**: Enable quick subtask creation by pressing 'n' in detail view with parent ID pre-filled

**Independent Test**: View any task's details, press 'n', verify form opens with parent ID pre-populated, create subtask, verify it appears in parent's subtask list

### Implementation for User Story 2

- [ ] T028 [US2] Modify initCreateForm() signature in internal/tui/form.go to accept parentID *string parameter
- [ ] T029 [US2] Update initCreateForm() implementation in internal/tui/form.go to pre-fill field 5 if parentID != nil
- [ ] T030 [US2] Update initCreateForm() call in handleKanbanKeys() in internal/tui/app.go to pass nil for parent ID
- [ ] T031 [US2] Add 'n' key handler in handleDetailKeys() in internal/tui/app.go to create subtask
- [ ] T032 [US2] In handleDetailKeys() 'n' handler in internal/tui/app.go, call initCreateForm(&m.detailTask.ID)
- [ ] T033 [US2] Add validation in saveForm() in internal/tui/form.go to prevent multi-level nesting (check if parent is itself a subtask)
- [ ] T034 [US2] Add formErr message "cannot create subtask of subtask (max 1 level nesting)" in internal/tui/form.go
- [ ] T035 [US2] Build project with `make build` and verify no compilation errors
- [ ] T036 [US2] Manual test: View task detail, press 'n', verify parent ID field is pre-filled
- [ ] T037 [US2] Manual test: Create subtask from detail view, submit form, verify subtask appears in detail view subtasks list
- [ ] T038 [US2] Manual test: Try to create subtask of a subtask, verify error message appears
- [ ] T039 [US2] Manual test: Cancel form after pressing 'n' in detail, verify returns to detail view without creating task

**Checkpoint**: User Story 2 complete and independently testable - subtask creation streamlined from detail view

---

## Phase 5: User Story 3 - View Subtasks in CLI List Mode (Priority: P2)

**Goal**: Display hierarchical task structure in CLI list and show commands with indentation

**Independent Test**: Run `ontop list` command and verify subtasks appear indented with dash prefix below their parents

### Implementation for User Story 3

- [ ] T040 [US3] Add import for internal/service/hierarchy package in internal/cli/list.go
- [ ] T041 [US3] Modify ListCommand() in internal/cli/list.go to call hierarchy.BuildFlatHierarchy()
- [ ] T042 [US3] Update output loop in ListCommand() in internal/cli/list.go to use HierarchicalTask
- [ ] T043 [US3] Add indentation logic in ListCommand() in internal/cli/list.go using strings.Repeat(" ", ht.Indentation)
- [ ] T044 [US3] Add dash prefix for subtasks in ListCommand() output in internal/cli/list.go
- [ ] T045 [US3] Add import for internal/service/hierarchy package in internal/cli/show.go
- [ ] T046 [US3] Add subtasks section to ShowCommand() in internal/cli/show.go after main task display
- [ ] T047 [US3] Call hierarchy.GetSubtasksForParent() in ShowCommand() in internal/cli/show.go
- [ ] T048 [US3] Format and display subtasks with indentation in ShowCommand() in internal/cli/show.go
- [ ] T049 [US3] Add formatColumnShort() helper function in internal/cli/show.go for column display
- [ ] T050 [US3] Build project with `make build` and verify no compilation errors
- [ ] T051 [US3] Manual test: Run `./ontop list`, verify subtasks indented below parents with dash prefix
- [ ] T052 [US3] Manual test: Run `./ontop show <parent-id>`, verify subtasks section displays with all subtasks
- [ ] T053 [US3] Manual test: Pipe CLI output `./ontop list | grep "-"`, verify scriptable (no ANSI issues)

**Checkpoint**: User Story 3 complete and independently testable - CLI displays hierarchy consistently with TUI

---

## Phase 6: User Story 4 - Navigate Subtask Hierarchy in Kanban (Priority: P2)

**Goal**: Enable keyboard navigation through hierarchical task list with j/k keys

**Independent Test**: Navigate through column with parents and subtasks using j/k keys, verify cursor moves sequentially through hierarchy

### Implementation for User Story 4

- [ ] T054 [US4] Verify handleKanbanKeys() j/k navigation in internal/tui/app.go works with hierarchical list (should already work due to US1 changes)
- [ ] T055 [US4] Manual test: Create parent with 3 subtasks in same column
- [ ] T056 [US4] Manual test: Navigate with 'j' key, verify cursor moves: parent → subtask1 → subtask2 → subtask3 → next parent
- [ ] T057 [US4] Manual test: Navigate with 'k' key, verify cursor moves: subtask → previous subtask → parent
- [ ] T058 [US4] Manual test: Select subtask with Enter key, verify detail view opens for subtask
- [ ] T059 [US4] Manual test: From subtask detail, press 'e' to edit, verify edit form opens for subtask
- [ ] T060 [US4] Manual test: From subtask in kanban, press 'm' to move, verify move prompt works for subtask

**Checkpoint**: User Story 4 complete and independently testable - full keyboard navigation through hierarchy

---

## Phase 7: Edge Cases & Data Integrity

**Purpose**: Handle edge cases and prevent data corruption

- [ ] T061 Add deletion protection in storage.DeleteTask() in internal/storage/sqlite.go
- [ ] T062 Query for subtasks COUNT before soft-delete in DeleteTask() in internal/storage/sqlite.go
- [ ] T063 Return error "cannot delete task with N subtasks" if count > 0 in internal/storage/sqlite.go
- [ ] T064 Build project with `make build` and verify no compilation errors
- [ ] T065 Manual test: Create parent with 2 subtasks, try to delete parent, verify deletion blocked with error
- [ ] T066 Manual test: Delete both subtasks first, then delete parent, verify succeeds
- [ ] T067 Manual test: Create subtask with non-existent parent ID via CLI, verify hierarchy handles gracefully (treats as orphan)

**Checkpoint**: Edge cases handled - data integrity protected

---

## Phase 8: Polish & Cross-Cutting Concerns

**Purpose**: Final validation and documentation

### Manual Testing Checklist

- [ ] T068 Execute full manual test suite from quickstart.md testing section
- [ ] T069 Test with 50+ tasks to verify performance is acceptable (<100ms rendering)
- [ ] T070 Test archive functionality: verify archived subtasks don't appear under active parents
- [ ] T071 Test sorting modes (priority, date): verify hierarchy grouping preserved across all sort modes
- [ ] T072 Test task moves: move subtask between columns, verify hierarchy updates correctly
- [ ] T073 Test multi-column scenario: parent in inbox, subtasks in inbox/progress/done, verify correct display

### Documentation & Cleanup

- [ ] T074 [P] Update README.md with subtask feature description and usage examples
- [ ] T075 [P] Update internal/cli/add.go help text to mention --parent flag for creating subtasks
- [ ] T076 [P] Update TUI help text (if displayed) to mention 'n' key in detail view creates subtask
- [ ] T077 [P] Review code for TODO comments, clean up debug logging
- [ ] T078 Run `go fmt ./...` to format all Go code
- [ ] T079 Run `make lint` if available to check code quality
- [ ] T080 Final build: `make build` and verify clean compilation
- [ ] T081 Run full test suite: `make test` and verify all tests pass

**Checkpoint**: Feature complete, tested, and documented - ready for use

---

## Dependencies & Execution Order

### Phase Dependencies

- **Setup (Phase 1)**: No dependencies - can start immediately
- **Foundational (Phase 2)**: Depends on Setup completion - BLOCKS all user stories
- **User Stories (Phases 3-6)**: All depend on Foundational phase (T017) completion
  - User Story 1 (Phase 3): Can start after Foundational - No dependencies on other stories
  - User Story 2 (Phase 4): Can start after Foundational - Independent but best done after US1 for testing
  - User Story 3 (Phase 5): Can start after Foundational - Independent, reuses hierarchy service
  - User Story 4 (Phase 6): Depends on User Story 1 (Phase 3) - validates existing navigation works with hierarchy
- **Edge Cases (Phase 7)**: Can run in parallel with user stories, but best after US1-US2
- **Polish (Phase 8)**: Depends on all user stories being complete

### User Story Dependencies

- **User Story 1 (P1)**: Foundation only → Hierarchical display in TUI kanban
- **User Story 2 (P1)**: Foundation only → Subtask creation from detail view (independent)
- **User Story 3 (P2)**: Foundation only → CLI hierarchy display (independent, reuses service)
- **User Story 4 (P2)**: User Story 1 → Navigation validation (verifies US1 works correctly)

### Within Each User Story

- Build and test after implementation tasks
- Manual tests after build verification
- Each story checkpoint before moving to next priority

### Parallel Opportunities

#### During Foundational Phase (Phase 2):
- T013, T014, T015, T016 (unit tests) can run in parallel after T012 complete

#### After Foundational Complete:
- User Story 1 (Phase 3): T018-T023 are sequential (modifying same files)
- User Story 2 (Phase 4): Can start in parallel with US1 if separate developer
- User Story 3 (Phase 5): Can start in parallel with US1-US2 if separate developer
- User Story 4 (Phase 6): Should wait for US1 completion

#### Polish Phase:
- T074, T075, T076, T077 (documentation) can all run in parallel

---

## Parallel Example: Foundational Phase

```bash
# After T012 (hierarchy_test.go created), launch unit tests together:
Task T013: "Add unit test TestBuildFlatHierarchy_EmptyInput"
Task T014: "Add unit test TestBuildFlatHierarchy_OnlyParents"
Task T015: "Add unit test TestBuildFlatHierarchy_MixedParentsAndSubtasks"
Task T016: "Add unit test TestBuildFlatHierarchy_OrphanedSubtask"
```

## Parallel Example: Polish Phase

```bash
# Launch all documentation tasks together:
Task T074: "Update README.md with subtask feature"
Task T075: "Update internal/cli/add.go help text"
Task T076: "Update TUI help text"
Task T077: "Review code for cleanup"
```

---

## Implementation Strategy

### MVP First (User Stories 1 & 2 Only - Both P1)

1. Complete Phase 1: Setup (T001-T005)
2. Complete Phase 2: Foundational (T006-T017) - CRITICAL
3. Complete Phase 3: User Story 1 (T018-T027) - Hierarchical display
4. Complete Phase 4: User Story 2 (T028-T039) - Subtask creation
5. **STOP and VALIDATE**: Test both P1 stories independently
6. Complete Phase 7: Edge Cases (T061-T067) - Data protection
7. Complete Phase 8: Polish (T068-T081)
8. Deploy/demo MVP

**MVP delivers**: Visual hierarchy in kanban + streamlined subtask creation

### Incremental Delivery

1. Complete Setup + Foundational → Hierarchy service ready
2. Add User Story 1 → Test independently → Working hierarchical kanban
3. Add User Story 2 → Test independently → Subtask creation workflow
4. Add User Story 3 → Test independently → CLI parity with TUI
5. Add User Story 4 → Test independently → Navigation validation
6. Each story adds value without breaking previous stories

### Single Developer Strategy

**Recommended order** (one person):
1. Phases 1-2 (Setup + Foundation) - Days 1
2. Phase 3 (US1) - Day 2 morning
3. Phase 4 (US2) - Day 2 afternoon
4. Phase 5 (US3) - Day 3 morning
5. Phase 6 (US4) - Day 3 afternoon
6. Phase 7 (Edge Cases) - Day 4 morning
7. Phase 8 (Polish) - Day 4 afternoon

**Timeline**: 4 days for experienced Go developer familiar with codebase

### Parallel Team Strategy

With multiple developers:

1. Team completes Setup + Foundational together (Days 1)
2. Once Foundational is done:
   - Developer A: User Story 1 (Phase 3)
   - Developer B: User Story 2 (Phase 4)
   - Developer C: User Story 3 (Phase 5)
3. Developer A completes US4 (Phase 6) after US1 done
4. Team handles Edge Cases + Polish together

**Timeline**: 2-3 days with 3 developers

---

## Task Counts

- **Phase 1 (Setup)**: 5 tasks
- **Phase 2 (Foundational)**: 12 tasks (including unit tests)
- **Phase 3 (US1 - P1)**: 10 tasks
- **Phase 4 (US2 - P1)**: 12 tasks
- **Phase 5 (US3 - P2)**: 14 tasks
- **Phase 6 (US4 - P2)**: 7 tasks
- **Phase 7 (Edge Cases)**: 7 tasks
- **Phase 8 (Polish)**: 14 tasks

**Total**: 81 tasks

**Critical Path**: T001-T017 (Foundational) → T018-T027 (US1) → T054-T060 (US4) = ~29 sequential tasks

**MVP Tasks**: T001-T039 + T061-T067 + T068-T081 = ~65 tasks for production-ready MVP

---

## Notes

- [P] tasks = different files, no dependencies, can run in parallel
- [Story] label maps task to specific user story for traceability
- Each user story is independently completable and testable
- Manual testing approach chosen (no automated test tasks per spec)
- Commit after each logical group of tasks
- Stop at any checkpoint to validate story independently
- Build and test after each implementation phase before manual testing
- QuickStart.md contains detailed gotchas and debugging tips
