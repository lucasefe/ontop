# Tasks: README Visual Enhancement with Screenshots

**Input**: Design documents from `/specs/003-readme-visual-enhancement/`
**Prerequisites**: plan.md (required), spec.md (required for user stories), research.md, data-model.md, contracts/

**Tests**: This is a documentation feature - no code tests required.

**Organization**: Tasks are grouped by user story to enable independent implementation and testing of each story.

## Format: `[ID] [P?] [Story] Description`

- **[P]**: Can run in parallel (different files, no dependencies)
- **[Story]**: Which user story this task belongs to (e.g., US1, US2, US3, US4)
- Include exact file paths in descriptions

## Path Conventions

This is a documentation-only feature:
- **Documentation**: `README.md` at repository root
- **Images**: `docs/images/` directory
- **Scripts**: `scripts/` directory for sample data generation

---

## Phase 1: Setup (Shared Infrastructure)

**Purpose**: Create directory structure and foundational assets needed for all user stories

- [x] T001 [P] Create `docs/images/` directory for screenshot storage
- [x] T002 [P] Create `scripts/` directory for sample data generation script
- [x] T003 Create `scripts/generate-sample-data.sh` with sample task definitions per data-model.md

---

## Phase 2: Foundational (Blocking Prerequisites)

**Purpose**: Generate sample data that MUST be complete before ANY screenshots or CLI examples can be captured

**⚠️ CRITICAL**: No user story work can begin until sample data exists

- [ ] T004 Implement sample data generation script in `scripts/generate-sample-data.sh` creating 5 parent tasks distributed across columns (2 Inbox, 2 In Progress, 1 Done)
- [ ] T005 Add 4 subtask creation commands to sample data script with proper parent relationships per data-model.md
- [ ] T006 Add task metadata to sample script (priorities 1-5, progress 0-100%, tags, multiline descriptions per sample-data-contract.md)
- [ ] T007 Make sample data script executable (`chmod +x scripts/generate-sample-data.sh`)
- [ ] T008 Test sample data script execution and verify output with `./ontop --db-path ./sample-tasks.db list`
- [ ] T009 Verify hierarchical task display shows subtasks indented below parents with dash prefix

**Checkpoint**: Sample database created - user story implementation can now begin in parallel

---

## Phase 3: User Story 1 - Visual TUI Demonstration (Priority: P1) 🎯 MVP

**Goal**: Provide immediate visual proof of OnTop's TUI interface showing kanban board with hierarchical tasks

**Independent Test**: View README on GitHub and verify clear screenshot of kanban board with 5+ sample tasks showing columns, hierarchy, and navigation

### Implementation for User Story 1

- [ ] T010 [US1] Launch OnTop TUI with sample database (`./ontop --db-path ./sample-tasks.db`)
- [ ] T011 [US1] Resize terminal to 140x35 characters for optimal kanban board display
- [ ] T012 [US1] Navigate TUI to position cursor on task in Inbox column with help panel visible
- [ ] T013 [US1] Capture kanban board screenshot showing all three columns with hierarchical tasks and save as `docs/images/tui-kanban-board.png`
- [ ] T014 [US1] Verify screenshot file size is under 500KB (optimize with `optipng` if needed)
- [ ] T015 [US1] Add hero screenshot to `README.md` after project description with descriptive alt text per research.md accessibility guidelines
- [ ] T016 [US1] Verify screenshot displays correctly when viewing README (push branch and check GitHub rendering)

**Checkpoint**: At this point, User Story 1 should be fully functional and testable independently - README shows TUI interface visually

---

## Phase 4: User Story 2 - CLI Usage Examples with Output (Priority: P1) 🎯 MVP

**Goal**: Demonstrate CLI command syntax and formatted output for automation workflows

**Independent Test**: View README CLI section and verify each command example shows both syntax and hierarchical output format

### Implementation for User Story 2

- [ ] T017 [P] [US2] Execute `./ontop --db-path ./sample-tasks.db list` and capture formatted output to temp file
- [ ] T018 [P] [US2] Execute `./ontop --db-path ./sample-tasks.db show <task-id>` for parent task and capture output showing subtask list
- [ ] T019 [P] [US2] Execute add command examples (parent and subtask with `--parent` flag) and capture output
- [ ] T020 [US2] Add CLI Mode section to `README.md` with command + output pairs using markdown code blocks per research.md format
- [ ] T021 [US2] Format CLI examples with bash syntax highlighting for commands and plain text blocks for output
- [ ] T022 [US2] Verify CLI output examples show hierarchical indentation and dash prefix for subtasks
- [ ] T023 [US2] Add at least 5 different CLI command examples (list, show, add parent, add subtask, move) with actual output

**Checkpoint**: At this point, User Stories 1 AND 2 should both work independently - README shows both TUI and CLI modes

---

## Phase 5: User Story 3 - Feature Showcase Section (Priority: P2)

**Goal**: Highlight specific capabilities (subtasks, multiline descriptions, vim navigation, archives) with visual proof

**Independent Test**: Review Features Showcase section and verify each major feature has corresponding screenshot or example

### Implementation for User Story 3

- [ ] T024 [P] [US3] Navigate TUI to parent task with visible subtasks and capture focused screenshot as `docs/images/tui-subtasks-hierarchy.png`
- [ ] T025 [P] [US3] Open task creation form (press `n`) or edit form (press `e`), fill all fields including multiline description, capture as `docs/images/tui-task-form.png`
- [ ] T026 [US3] Create "Features Showcase" section in `README.md` after Features list with subsections for each major feature
- [ ] T027 [US3] Add Hierarchical Subtasks showcase subsection with description and `tui-subtasks-hierarchy.png` screenshot with alt text
- [ ] T028 [US3] Add Rich Task Forms showcase subsection with description and `tui-task-form.png` screenshot with alt text
- [ ] T029 [US3] Add Vim Navigation showcase subsection with keyboard shortcut table or inline examples
- [ ] T030 [US3] Add Dual-Mode Architecture showcase with side-by-side TUI and CLI examples showing same task workflow
- [ ] T031 [US3] Verify all showcase screenshots have descriptive alt text per accessibility requirements
- [ ] T032 [US3] Verify total image directory size is under 2MB (`du -sh docs/images/`)

**Checkpoint**: At this point, User Stories 1, 2, AND 3 should all work independently - README has comprehensive visual feature documentation

---

## Phase 6: User Story 4 - Quick Start Visual Guide (Priority: P3)

**Goal**: Provide visual walkthrough for first-time users to get started immediately

**Independent Test**: Follow Quick Start section and verify it includes visual elements showing task creation workflow

### Implementation for User Story 4

- [ ] T033 [US4] Create "Quick Start" section in `README.md` after Installation section
- [ ] T034 [US4] Add step-by-step quick start instructions: 1) Install, 2) Launch TUI, 3) Create first task, 4) Navigate with vim keys
- [ ] T035 [US4] Reference `tui-task-form.png` screenshot in Quick Start to show form example
- [ ] T036 [US4] Add CLI-based quick start alternative showing command sequence for CLI users
- [ ] T037 [US4] Add note about recommended terminal size (minimum 80x24) and UI adaptability
- [ ] T038 [US4] Verify Quick Start walkthrough is complete and user can successfully create first task

**Checkpoint**: All user stories should now be independently functional - README provides complete visual onboarding

---

## Phase 7: Polish & Cross-Cutting Concerns

**Purpose**: Improvements that affect multiple user stories and final verification

- [ ] T039 [P] Add TUI Keyboard Shortcuts reference section to `README.md` with complete keybinding table per quickstart.md
- [ ] T040 [P] Restructure README with visual-first flow: Title → Hero Screenshot → Features → Showcase → Installation → Quick Start → Usage → Development
- [ ] T041 [P] Add HTML comment in `README.md` documenting screenshot regeneration process and referencing `scripts/generate-sample-data.sh`
- [ ] T042 Verify all images committed to git (`git add docs/images/*.png && git status`)
- [ ] T043 Review all image alt text for WCAG compliance (descriptive, under 125 characters, explains content not "screenshot")
- [ ] T044 Test README rendering on GitHub by pushing branch and viewing on web interface
- [ ] T045 Verify README loads within 3 seconds despite multiple images
- [ ] T046 Verify sample data reproducibility by deleting database and running `scripts/generate-sample-data.sh` again
- [ ] T047 Run constitution compliance check - verify Documentation Currency principle satisfied (README synchronized with features)
- [ ] T048 Final review: All 4 user stories independently testable and acceptance criteria met

---

## Dependencies & Execution Order

### Phase Dependencies

- **Setup (Phase 1)**: No dependencies - can start immediately
- **Foundational (Phase 2)**: Depends on Setup completion - BLOCKS all user stories
- **User Stories (Phase 3-6)**: All depend on Foundational phase completion
  - User stories can then proceed in parallel (if multiple contributors)
  - Or sequentially in priority order (P1 → P1 → P2 → P3)
- **Polish (Phase 7)**: Depends on all desired user stories being complete

### User Story Dependencies

- **User Story 1 (P1 - TUI Screenshots)**: Can start after Foundational (Phase 2) - No dependencies on other stories
- **User Story 2 (P1 - CLI Examples)**: Can start after Foundational (Phase 2) - No dependencies on other stories
- **User Story 3 (P2 - Feature Showcase)**: Can start after Foundational (Phase 2) - May reuse screenshots from US1 but independently testable
- **User Story 4 (P3 - Quick Start)**: Can start after Foundational (Phase 2) - May reference US1/US3 screenshots but independently testable

### Within Each User Story

- **US1**: Sample data → TUI launch → Screenshot capture → README addition → Verification
- **US2**: Sample data → CLI execution → Output capture → README formatting → Verification
- **US3**: Sample data → Additional screenshots → Showcase section creation → Verification
- **US4**: Existing screenshots → Quick Start writing → Verification

### Parallel Opportunities

- All Setup tasks (T001-T003) can run in parallel
- Within Foundational phase: Tasks T004-T006 can run in sequence, T007 must follow
- Once Foundational completes, User Stories 1 and 2 can run completely in parallel
- User Story 3 can run in parallel with US1 and US2 after Foundational
- User Story 4 can start after any screenshots are available from US1 or US3
- Within Phase 7: Tasks T039-T041 can run in parallel, T042-T048 must be sequential

---

## Parallel Example: User Stories 1 and 2 (Both P1)

```bash
# After Phase 2 (Foundational) completes, launch both user stories in parallel:

# Terminal 1: User Story 1 (TUI Screenshots)
Task: "Launch OnTop TUI with sample database"
Task: "Capture kanban board screenshot"
Task: "Add hero screenshot to README"

# Terminal 2: User Story 2 (CLI Examples)
Task: "Execute CLI list command and capture output"
Task: "Execute CLI show command and capture output"
Task: "Add CLI Mode section to README with examples"
```

---

## Implementation Strategy

### MVP First (User Stories 1 AND 2 Only)

Both User Story 1 and User Story 2 are marked P1 (highest priority) because they demonstrate the dual-mode architecture - OnTop's core value proposition.

1. Complete Phase 1: Setup (directory creation)
2. Complete Phase 2: Foundational (sample data generation - CRITICAL)
3. Complete Phase 3: User Story 1 (TUI screenshot)
4. Complete Phase 4: User Story 2 (CLI examples)
5. **STOP and VALIDATE**: README now shows both TUI and CLI modes
6. Deploy/demo if ready - users can evaluate OnTop visually

### Incremental Delivery

1. Complete Setup + Foundational → Sample data ready
2. Add User Story 1 → Test independently → Commit (TUI visualization live)
3. Add User Story 2 → Test independently → Commit (CLI examples live)
4. Add User Story 3 → Test independently → Commit (Feature showcase live)
5. Add User Story 4 → Test independently → Commit (Quick start guide live)
6. Each story adds value without breaking previous stories

### Parallel Team Strategy

With multiple contributors (or parallel work sessions):

1. Complete Setup + Foundational together (required foundation)
2. Once Foundational is done:
   - Contributor A: User Story 1 (TUI screenshots)
   - Contributor B: User Story 2 (CLI examples)
   - Contributor C: User Story 3 (Feature showcase - waits for US1 screenshots)
3. Stories complete and integrate independently
4. Contributor D: User Story 4 (Quick Start - waits for US1 or US3 screenshots)

### Single-Person Sequential Strategy

For solo implementation (recommended flow):

1. Setup + Foundational (30 minutes)
2. User Story 1 - TUI Screenshots (20 minutes)
3. User Story 2 - CLI Examples (20 minutes)
4. **MVP CHECKPOINT** - Validate dual-mode documentation
5. User Story 3 - Feature Showcase (30 minutes)
6. User Story 4 - Quick Start (20 minutes)
7. Polish phase (30 minutes)

**Total time**: 2.5-3 hours as estimated in quickstart.md

---

## Notes

- **[P] tasks**: Different files or independent actions, no blocking dependencies
- **[Story] label**: Maps task to specific user story for traceability
- **Each user story**: Independently completable and testable - can stop after any phase
- **No code changes**: This is documentation-only - no application code modified
- **Sample data critical**: Phase 2 (Foundational) must complete before any screenshots/examples
- **Commit strategy**: Commit after each user story phase completion for incremental delivery
- **Stop at any checkpoint**: Validate that user story independently delivers value
- **Constitution alignment**: Implements Documentation Currency principle (Principle IV)

---

## Task Summary

**Total Tasks**: 48
- Phase 1 (Setup): 3 tasks
- Phase 2 (Foundational): 6 tasks
- Phase 3 (User Story 1 - P1): 7 tasks
- Phase 4 (User Story 2 - P1): 7 tasks
- Phase 5 (User Story 3 - P2): 9 tasks
- Phase 6 (User Story 4 - P3): 6 tasks
- Phase 7 (Polish): 10 tasks

**Parallel Opportunities**: 15 tasks marked with [P]

**Independent Test Criteria**:
- US1: README shows TUI screenshot with kanban board
- US2: README shows CLI examples with hierarchical output
- US3: README has Feature Showcase with visual proof
- US4: README has Quick Start guide with visual walkthrough

**Suggested MVP Scope**: Phase 1 + Phase 2 + Phase 3 (US1) + Phase 4 (US2) = 23 tasks

This MVP delivers the core value: Visual documentation of both TUI and CLI modes demonstrating OnTop's dual-mode architecture.
