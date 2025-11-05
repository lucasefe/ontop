# Feature Specification: Hierarchical Subtask Display and Management

**Feature Branch**: `002-subtask-display-edit`
**Created**: 2025-01-04
**Status**: Draft
**Input**: User description: "we need to find a good way to add sub tasks to parent tasks. See them in the task detail, and edit them in the form. I think in the clumn view we should show right below the parent task each subtask, 2 chars indentation, with a dash if it looks right. same in cli mode when listing the tasks, or a specific one. When in detail view, also show them. For editing, in the detail view when pressing n we should be able to add a sub task."

## User Scenarios & Testing *(mandatory)*

### User Story 1 - View Subtasks in Kanban Column (Priority: P1)

Users need to see parent-child task relationships directly in the kanban board without navigating to detail views. When viewing tasks in a column, subtasks should appear immediately below their parent task with visual indentation to indicate hierarchy.

**Why this priority**: This is the primary visual interface where users spend most time managing tasks. Without hierarchical display in the kanban view, users must memorize or repeatedly check task details to understand relationships.

**Independent Test**: Can be fully tested by creating a parent task with 2-3 subtasks and verifying they appear indented below the parent in the kanban column view. Delivers immediate value by showing task hierarchy at a glance.

**Acceptance Scenarios**:

1. **Given** a parent task exists in a column with 3 subtasks, **When** user views the kanban board, **Then** all 3 subtasks appear immediately below the parent with 2-character indentation
2. **Given** a parent task with subtasks, **When** the parent is in "Inbox" column, **Then** all its subtasks display in their respective columns (may be different from parent)
3. **Given** multiple parent tasks with subtasks in the same column, **When** user views the column, **Then** each parent's subtasks group together directly below that parent
4. **Given** a subtask exists, **When** displayed in the kanban column, **Then** it shows with a dash prefix (e.g., "- Subtask title") to visually indicate hierarchy

---

### User Story 2 - Create Subtasks from Detail View (Priority: P1)

Users need a quick way to add subtasks to an existing parent task without manually copying task IDs. When viewing a task's details, pressing 'n' should create a new subtask with the parent relationship automatically established.

**Why this priority**: This is the most common workflow for breaking down complex tasks into smaller pieces. Without this, users must remember task IDs and manually enter them, which is error-prone and slow.

**Independent Test**: Can be fully tested by viewing any task's details, pressing 'n', filling the form, and verifying the new task appears as a subtask of the viewed task. Delivers immediate value by streamlining subtask creation.

**Acceptance Scenarios**:

1. **Given** user is viewing a task's detail view, **When** user presses 'n', **Then** the task creation form opens with the parent task ID pre-filled
2. **Given** user creates a subtask from detail view, **When** the form is submitted, **Then** the new subtask appears in the parent's subtask list
3. **Given** user is viewing a subtask's detail view (not a parent), **When** user presses 'n', **Then** a new subtask is created as a sibling (sharing the same parent)
4. **Given** user creates a subtask via detail view, **When** user cancels the form, **Then** no subtask is created and user returns to the detail view

---

### User Story 3 - View Subtasks in CLI List Mode (Priority: P2)

Users working in CLI mode need to see task hierarchy when listing tasks. Subtasks should appear indented below their parent tasks with the same visual pattern as the TUI kanban view.

**Why this priority**: CLI users need feature parity with TUI users. While less visually rich, the CLI is often preferred for scripting and quick queries.

**Independent Test**: Can be fully tested by running `list` command and verifying parent tasks show their subtasks indented with dash prefix. Delivers value by providing hierarchy visibility in CLI workflows.

**Acceptance Scenarios**:

1. **Given** user runs `list` command, **When** results are displayed, **Then** subtasks appear 2 characters indented below their parent with dash prefix
2. **Given** user runs `list` with column filter (e.g., `--column inbox`), **When** results display, **Then** only parent tasks in that column show, along with all their subtasks regardless of subtask column
3. **Given** user runs `show <task-id>` for a parent task, **When** details are displayed, **Then** all subtasks are listed with their key attributes

---

### User Story 4 - Navigate Subtask Hierarchy in Kanban (Priority: P2)

Users need to navigate and select both parent tasks and subtasks using keyboard controls in the kanban view. Arrow keys should move through the hierarchical list, allowing selection of any task or subtask.

**Why this priority**: After displaying hierarchy (P1), users need to interact with it. Without proper navigation, users cannot select subtasks to view details or perform actions.

**Independent Test**: Can be fully tested by navigating through a column containing parent tasks with subtasks using j/k keys and verifying cursor moves through both parents and subtasks. Delivers value by making subtasks fully actionable.

**Acceptance Scenarios**:

1. **Given** user's cursor is on a parent task with subtasks, **When** user presses 'j' (down), **Then** cursor moves to the first subtask below the parent
2. **Given** user's cursor is on a subtask, **When** user presses 'j' (down), **Then** cursor moves to the next subtask or to the next parent task if no more subtasks exist
3. **Given** user's cursor is on a subtask, **When** user presses 'k' (up), **Then** cursor moves to the previous subtask or parent task
4. **Given** user selects a subtask in kanban, **When** user presses Enter, **Then** the detail view for that subtask opens

---

### Edge Cases

- **Orphaned subtasks**: What happens when a parent task is deleted but subtasks exist? Subtasks should either become top-level tasks or prevent parent deletion.
- **Deep nesting**: What happens if a subtask itself has subtasks (grandchild tasks)? Assume maximum 1 level of nesting (parent → subtask only) unless multi-level hierarchy is explicitly supported.
- **Empty parents**: How are parent tasks with no subtasks displayed? They appear as regular tasks without indentation markers.
- **Subtask moves**: When a subtask moves to a different column than its parent, does it still appear under the parent in the parent's column? No - subtasks appear in their own column but remain visually grouped when in the same column as parent.
- **Filtering**: When filtering tasks (by tag, priority, etc.), should subtasks be included if their parent matches the filter? Yes - if parent matches, all its subtasks display regardless of subtask attributes.
- **Sorting**: When sorting tasks in a column (by priority, date, etc.), should parent-subtask groupings remain intact? Yes - parents sort first, then their subtasks appear immediately below before the next parent.

## Requirements *(mandatory)*

### Functional Requirements

#### Display Requirements

- **FR-001**: System MUST display subtasks immediately below their parent task in kanban column view
- **FR-002**: System MUST indent subtasks by 2 characters relative to parent tasks in both TUI and CLI views
- **FR-003**: System MUST prefix subtasks with a dash character (e.g., "- Subtask title") to visually indicate hierarchy
- **FR-004**: System MUST display subtasks in their actual column, not forced into parent's column
- **FR-005**: System MUST maintain parent-subtask grouping when both are in the same column
- **FR-006**: System MUST show all subtasks for a parent in the detail view, regardless of subtask column

#### Navigation Requirements

- **FR-007**: System MUST allow keyboard navigation (j/k keys) to move cursor through both parent tasks and subtasks in hierarchical order
- **FR-008**: System MUST allow selecting subtasks in kanban view to view their details or perform actions
- **FR-009**: System MUST support all standard task operations (move, edit, delete, archive) on subtasks

#### Creation Requirements

- **FR-010**: System MUST open task creation form with parent ID pre-filled when user presses 'n' in a task's detail view
- **FR-011**: System MUST establish parent-child relationship automatically when creating subtask from detail view
- **FR-012**: System MUST allow creating subtasks from subtasks (creating sibling subtasks with shared parent)
- **FR-013**: System MUST allow manual parent ID entry when creating tasks from kanban view (existing behavior)

#### CLI Requirements

- **FR-014**: CLI `list` command MUST display subtasks indented below their parent with dash prefix
- **FR-015**: CLI `show <task-id>` command MUST display all subtasks when showing a parent task
- **FR-016**: CLI list output MUST maintain parent-subtask grouping and hierarchy

#### Data Integrity Requirements

- **FR-017**: System MUST prevent orphaned subtasks when deleting parent tasks (either block deletion or promote subtasks to top-level)
- **FR-018**: System MUST maintain referential integrity between parent and subtask IDs

### Key Entities

- **Task**: Represents both parent tasks and subtasks; contains optional ParentID field to establish hierarchy
- **Subtask**: A task with a non-null ParentID field, creating a parent-child relationship

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: Users can visually identify parent-child task relationships in the kanban view without opening detail views
- **SC-002**: Users can create subtasks from detail view in under 10 seconds (vs. current manual ID entry workflow)
- **SC-003**: Subtask hierarchy displays consistently across TUI kanban, TUI detail view, and CLI list/show commands
- **SC-004**: Users can navigate to and select any subtask in the kanban view using keyboard controls
- **SC-005**: 100% of subtasks display with 2-character indentation and dash prefix in all views
- **SC-006**: Creating a subtask from detail view automatically establishes the parent relationship without manual ID entry

## Assumptions *(optional)*

- Tasks already support a ParentID field in the data model (based on user's description referencing existing parent/subtask functionality)
- Subtask hierarchy is limited to 1 level (parent → subtask only, no grandchild subtasks)
- When parent and subtask are in different columns, they appear in their respective columns but subtasks don't visually group under parent
- Subtask creation from detail view defaults to placing subtask in the same column as parent (user can change this in the form)
- Deleting a parent task with subtasks will block deletion and require user to handle subtasks first (conservative approach to prevent data loss)

## Open Questions

- Should subtasks inherit any attributes from their parent (column, tags, priority)? **Assumption**: No inheritance - subtasks are independent except for ParentID relationship.
- How many levels of nesting should be supported (parent → subtask → sub-subtask)? **Assumption**: Maximum 1 level for simplicity.
- Should moving a parent task automatically move all its subtasks to the same column? **Assumption**: No - subtasks remain in their assigned columns and must be moved independently.
