# Feature Specification: OnTop - Dual-Mode Task Manager

**Feature Branch**: `001-ontop`
**Created**: 2025-11-04
**Status**: Draft
**Input**: User description: "I want to build a TUI app + CLI. Dual mode. It should be like a crud for tasks. i want to have sub tasks inside a given task. It should have priority. Progress (based off of the completeness of the task itself, or the inside subtasks. Task belongs to a column. Columns are like trello. Inbox, in-progress, done. We can archive tasks, so tey do not show. We can delete them, to m ake them go away. We need to be able to filter them."

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Quick Task Capture via CLI (Priority: P1)

As a user working in the terminal, I need to quickly capture tasks without leaving my current workflow, so I can immediately return to what I'm doing without context switching.

**Why this priority**: This is the most fundamental and frequently used capability. Users need a frictionless way to capture tasks as they think of them. This represents the minimum viable product - a working task capture and storage system.

**Independent Test**: Can be fully tested by running CLI commands to create tasks and verifying they are persisted. Delivers immediate value as a simple task inbox.

**Acceptance Scenarios**:

1. **Given** I am at the terminal, **When** I run the command to add a new task with description and priority, **Then** the task is created in the Inbox column with a unique timestamp-based ID and confirmation is displayed
2. **Given** I have an existing task, **When** I run the command to add a subtask to it, **Then** the subtask is created under the parent task and confirmation is displayed
3. **Given** I want to see my tasks, **When** I run the list command, **Then** all active tasks are displayed with their IDs and key attributes (description, priority, column, progress)
4. **Given** I know a task ID, **When** I run the command to open that task by ID, **Then** the full task details are displayed in the terminal (all attributes, subtasks, timestamps)

---

### User Story 2 - Interactive Task Management via TUI (Priority: P2)

As a user planning my day, I need to view and organize tasks in a visual kanban-style interface, so I can see the big picture and reorganize priorities using familiar vim-style navigation.

**Why this priority**: Once tasks are captured, users need to organize and prioritize them. The TUI provides the visual overview needed for effective task management and planning sessions.

**Independent Test**: Can be tested by launching the TUI, navigating with j/k keys, moving tasks between columns, and verifying changes persist. Works independently with tasks created in User Story 1.

**Acceptance Scenarios**:

1. **Given** I launch the TUI, **When** the interface loads, **Then** I see three vertical columns (Inbox, In Progress, Done) displayed side-by-side with brief task cards in each column
2. **Given** I am viewing the TUI, **When** I press j/k keys, **Then** the selection moves down/up within the current column
3. **Given** I am viewing the TUI, **When** I press h/l keys, **Then** the selection moves left/right between columns
4. **Given** each task card is displayed, **When** viewing, **Then** I see the task description (truncated if long), priority indicator, and progress percentage on the card
5. **Given** I have selected a task card, **When** I press the expand/enter key, **Then** the card expands to a full-detail view showing all task attributes (description, priority, tags, progress, created date, updated date, completed date, subtasks) and the column view is hidden
6. **Given** I am in the full-detail view, **When** I press the close/escape key, **Then** I return to the column view with the same task still selected
7. **Given** I am viewing a task in the TUI, **When** I press the move command key, **Then** I can select a new column and the task moves to that column
8. **Given** I select a task with subtasks in the detail view, **When** viewing, **Then** I see all subtasks listed with their completion status and the parent's calculated progress

---

### User Story 3 - Task Filtering and Search (Priority: P3)

As a user with many tasks, I need to filter tasks by various criteria (priority, column, tags), so I can focus on relevant work without visual clutter.

**Why this priority**: As the task list grows, filtering becomes important for productivity. However, the system is still valuable without filtering for users with smaller task lists.

**Independent Test**: Can be tested by creating tasks with different attributes and running filter commands in both CLI and TUI modes. Filtering does not affect core task management functionality.

**Acceptance Scenarios**:

1. **Given** I have tasks with different priorities, **When** I apply a priority filter in CLI, **Then** only tasks matching that priority are displayed
2. **Given** I am in the TUI, **When** I activate the filter mode and select criteria, **Then** the view updates to show only matching tasks
3. **Given** I have applied filters, **When** I clear the filter, **Then** all active tasks are visible again

---

### User Story 4 - Archive and Delete Operations (Priority: P4)

As a user maintaining my task list over time, I need to archive completed tasks and delete unwanted tasks, so my active view stays clean and focused.

**Why this priority**: Important for long-term maintenance but not critical for initial usage. Users can work effectively with all tasks visible for a period before needing cleanup.

**Independent Test**: Can be tested by archiving and deleting tasks, verifying they disappear from active views but archived tasks remain queryable. Does not affect other functionality.

**Acceptance Scenarios**:

1. **Given** I have a completed task, **When** I archive it via CLI or TUI, **Then** it no longer appears in default views but remains accessible via archived view
2. **Given** I have a task I want to remove permanently, **When** I delete it, **Then** it is removed from the system entirely and cannot be recovered
3. **Given** I want to review old tasks, **When** I view archived tasks, **Then** I see all previously archived tasks with their original metadata

---

### Edge Cases

- What happens when a parent task is deleted but has subtasks?
- How does the system handle moving a parent task to Done when subtasks are incomplete?
- What happens if a user tries to create a subtask under a subtask (nested subtasks)?
- How does progress calculation work when subtasks are added or removed dynamically?
- What happens when the terminal window is too small to display the TUI properly?
- How are tasks sorted within a column when multiple tasks have the same priority?
- What happens when filtering results in zero matches?

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: System MUST assign each task a unique ID based on creation timestamp that is sortable and looks like an identifier
- **FR-002**: System MUST allow creating tasks with description, priority, and optional tags via CLI
- **FR-003**: System MUST allow creating subtasks under existing tasks
- **FR-004**: System MUST organize tasks into three columns: Inbox, In Progress, Done
- **FR-005**: System MUST allow viewing/opening a task by its ID in CLI mode to display all task details
- **FR-006**: System MUST allow moving tasks between columns via both CLI and TUI modes
- **FR-007**: System MUST calculate task progress automatically based on subtask completion percentage
- **FR-008**: System MUST support manual progress updates for tasks without subtasks
- **FR-009**: System MUST provide a full-screen TUI mode displaying three vertical columns side-by-side (Inbox, In Progress, Done)
- **FR-010**: System MUST display tasks as compact cards within columns showing description, priority, and progress
- **FR-011**: System MUST support vim-style navigation keybindings: j/k for moving down/up within a column, h/l for moving left/right between columns
- **FR-012**: System MUST allow expanding any task card to a full-detail view that hides the column layout and displays all task attributes
- **FR-013**: System MUST allow closing the detail view to return to the column layout with the same task selected
- **FR-014**: System MUST provide CLI mode for non-interactive scripted operations
- **FR-015**: System MUST persist all task data between sessions
- **FR-016**: System MUST allow archiving tasks to remove them from active views
- **FR-017**: System MUST allow permanent deletion of tasks
- **FR-018**: System MUST support filtering tasks by priority, column, and tags in both modes
- **FR-019**: System MUST display task attributes including ID, description, priority, tags, column, progress percentage, created date, updated date, completed date, deleted date
- **FR-020**: System MUST track creation, last update, completion, and deletion timestamps for all tasks
- **FR-021**: System MUST automatically set the completion timestamp when a task is moved to the Done column
- **FR-022**: System MUST automatically set the deletion timestamp when a task is permanently deleted
- **FR-023**: System MUST prevent archived tasks from appearing in default views
- **FR-024**: System MUST allow viewing archived tasks when explicitly requested

### Key Entities

- **Task**: Represents a work item with unique timestamp-based ID, description, priority level, set of tags, current column (Inbox/In Progress/Done), progress percentage (0-100), archived status (boolean), creation timestamp, last update timestamp, completion timestamp (when moved to Done column), and deletion timestamp (when deleted). Can contain zero or more subtasks. Progress is either manual (for leaf tasks) or calculated from subtask completion. The ID is sortable and looks like an identifier (not just a raw timestamp).

- **Subtask**: A child task belonging to exactly one parent task. Has the same attributes as a task but exists within the context of its parent. Contributes to parent's progress calculation.

- **Column**: A workflow stage (Inbox, In Progress, Done) that contains tasks. Represents the current state of a task in the workflow.

- **Tag**: A label attached to tasks for categorization and filtering. Multiple tags can be applied to a single task.

- **Priority**: A numeric level from 1-5 indicating task importance, where 1 is the highest priority and 5 is the lowest. Users can assign any value in this range when creating or updating tasks.

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: Users can capture a new task in under 5 seconds using the CLI without leaving their terminal workflow
- **SC-002**: Users can reorganize 10 tasks across columns in under 30 seconds using the TUI with vim-style navigation
- **SC-003**: The system displays task progress accurately within 1 second of subtask status changes
- **SC-004**: Users can locate a specific task among 100+ tasks in under 10 seconds using filters
- **SC-005**: All task data persists correctly across application restarts with zero data loss
- **SC-006**: The TUI responds to keyboard input with no perceptible lag (< 100ms) on standard terminal emulators
- **SC-007**: Users report that vim-style navigation feels natural and intuitive (subjective validation with vim users)

## Assumptions

- The primary user is comfortable with command-line interfaces and vim keybindings
- Tasks are managed by a single user (no multi-user collaboration required)
- The system runs on a local machine (no cloud sync or multi-device support initially)
- Terminal supports standard ANSI color codes and cursor positioning
- Users have permission to write to the local filesystem for data persistence
- Average task list size is 50-200 tasks with archived tasks potentially reaching 1000+
- Most tasks have 0-5 subtasks, deeply nested subtasks are rare
- Priority levels use numeric scale 1-5 where 1 is highest priority and 5 is lowest
- Filters use AND logic when multiple criteria are specified (e.g., "High priority AND In Progress")
- Archived tasks are retained indefinitely unless explicitly deleted
- Task descriptions are plain text without rich formatting
- Tags are simple strings without hierarchy or special characters

## Out of Scope

The following are explicitly not included in this feature:

- Multi-user collaboration or task sharing
- Cloud synchronization or backup
- Recurring tasks or scheduled reminders
- Time tracking or task duration estimates
- Task dependencies or blocking relationships
- Rich text formatting in descriptions
- File attachments or links to external resources
- Email or notification integration
- Mobile app or web interface
- Import/export from other task management tools (may be added later)
- Task templates or automation rules
- Custom columns beyond Inbox/In Progress/Done
- Undo/redo functionality
- Task history or audit log
