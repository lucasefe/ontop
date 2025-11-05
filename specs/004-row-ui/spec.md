# Feature Specification: Row-Based UI Mode

**Feature Branch**: `004-row-ui`
**Created**: 2025-11-05
**Status**: Draft
**Input**: User description: "I want to experiment with a new ui. I don't think it works to have columns in the terminal. I want to introduce a new UI mode where we see the columns as rows, and when we move the tasks, we move between rows."

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Switch to Row-Based View (Priority: P1)

Users want to view their tasks in a row-based layout instead of the current column layout, where each workflow stage (inbox, in_progress, done) is displayed as a separate row with tasks listed horizontally within each row.

**Why this priority**: This is the foundation of the entire feature - without the ability to switch to and view the row-based layout, no other functionality can be tested or used. It's the minimum viable product that demonstrates the alternative UI approach.

**Independent Test**: Can be fully tested by launching the TUI, pressing the view toggle key, and verifying that tasks are displayed in three horizontal rows (one per workflow stage) instead of three vertical columns. Delivers immediate value by providing an alternative visualization that may work better in certain terminal environments.

**Acceptance Scenarios**:

1. **Given** the TUI is running in column mode, **When** user presses the view toggle key (e.g., 'v' for view), **Then** the display switches to row-based layout with three rows labeled "inbox", "in_progress", and "done"
2. **Given** the TUI is in row-based mode, **When** user views their tasks, **Then** tasks within each row are displayed horizontally with clear visual separation between workflow stages
3. **Given** the TUI is in row-based mode, **When** user presses the view toggle key again, **Then** the display switches back to the original column-based layout

---

### User Story 2 - Navigate Between Rows (Priority: P2)

Users want to navigate between workflow stage rows using keyboard controls, similar to how they currently navigate between columns, so they can select tasks in different workflow stages.

**Why this priority**: Once users can see the row-based layout, they need basic navigation to select tasks. This is essential for any task interaction but builds on the P1 visualization foundation.

**Independent Test**: Can be tested by switching to row mode, using up/down arrow keys to move between rows (inbox → in_progress → done), and verifying that the current row is visually highlighted. Delivers value by enabling users to focus on different workflow stages.

**Acceptance Scenarios**:

1. **Given** the TUI is in row-based mode with cursor on inbox row, **When** user presses down arrow, **Then** cursor moves to in_progress row and that row is visually highlighted
2. **Given** the TUI is in row-based mode with cursor on done row, **When** user presses up arrow, **Then** cursor moves to in_progress row
3. **Given** the TUI is in row-based mode with cursor on a row, **When** user presses left/right arrows, **Then** cursor moves between tasks within that row
4. **Given** the TUI is in row-based mode, **When** user navigates to a row, **Then** the currently selected task in that row is clearly indicated

---

### User Story 3 - Move Tasks Between Rows (Priority: P3)

Users want to move tasks between workflow stage rows using keyboard commands, similar to the current move functionality in column mode, so they can update task workflow status in row-based view.

**Why this priority**: This completes the core workflow functionality in row mode. While users could switch back to column mode to move tasks, providing this in row mode makes it a fully functional alternative view.

**Independent Test**: Can be tested by selecting a task in inbox row, pressing the move command key, choosing a destination row (in_progress or done), and verifying the task appears in the new row. Delivers value by making row mode feature-complete for basic task management.

**Acceptance Scenarios**:

1. **Given** a task is selected in inbox row, **When** user initiates move command and selects in_progress row, **Then** task is moved to in_progress row and cursor follows
2. **Given** a task is selected in in_progress row, **When** user initiates move command and selects done row, **Then** task is moved to done row and completed_at timestamp is set
3. **Given** a task is selected in done row, **When** user initiates move command and selects inbox row, **Then** task is moved back to inbox and completed_at is cleared
4. **Given** user is moving a task, **When** rows are displayed as move destinations, **Then** current row is clearly indicated and unavailable as destination

---

### Edge Cases

- What happens when a row has more tasks than can fit horizontally on screen? (Should support horizontal scrolling or pagination)
- What happens when terminal width is very narrow? (Should gracefully handle or show minimum width warning)
- What happens when user switches between column and row mode? (Should preserve cursor position/selected task when possible)
- What happens when a row is empty? (Should display empty state message like "No tasks in inbox")
- What happens when viewing task details in row mode? (Should maintain row mode context when returning from detail view)

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: System MUST provide a toggle mechanism to switch between column-based and row-based layout modes
- **FR-002**: System MUST display workflow stages (inbox, in_progress, done) as horizontal rows when in row-based mode
- **FR-003**: System MUST display tasks within each row horizontally with clear visual separation between tasks
- **FR-004**: Users MUST be able to navigate between rows using up/down navigation keys
- **FR-005**: Users MUST be able to navigate between tasks within a row using left/right navigation keys
- **FR-006**: System MUST visually indicate which row is currently selected
- **FR-007**: System MUST visually indicate which task within a row is currently selected
- **FR-008**: Users MUST be able to initiate move operations on tasks while in row-based mode
- **FR-009**: System MUST display available destination rows when user initiates task move
- **FR-010**: System MUST move tasks between rows and update workflow status accordingly
- **FR-011**: System MUST persist the user's preferred view mode between sessions in a configuration file located at ~/.config/ontop/ontop.toml
- **FR-012**: System MUST handle horizontal scrolling when row contains more tasks than fit on screen width
- **FR-013**: System MUST maintain task data integrity when switching between view modes

### Key Entities

This feature primarily affects the UI presentation layer and does not introduce new data entities. Existing entities (Task, Column) remain unchanged.

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: Users can switch between column and row views in under 1 second with a single keypress
- **SC-002**: Users can navigate to any workflow stage row and select any task within 5 keypresses
- **SC-003**: Users can move tasks between workflow stages in row mode with the same number of actions as column mode
- **SC-004**: Row-based layout displays correctly on terminal widths from 80 to 200 characters
- **SC-005**: Task selection and movement success rate remains at 100% regardless of view mode
- **SC-006**: No data loss or corruption occurs when switching between view modes during active session

## Assumptions *(optional)*

- The existing three-column workflow (inbox, in_progress, done) remains unchanged
- Row-based mode is an alternative visualization of the same underlying data structure
- Keyboard navigation patterns are consistent between column and row modes (arrows for navigation, same keys for commands)
- Terminal environments support the minimum width requirements for readable row-based display
- Users may frequently switch between modes during a session to find their preferred view
- The default view mode on first launch will be column-based; subsequent launches will use the persisted preference
- Performance characteristics (responsiveness, rendering speed) are equivalent between view modes
- Configuration file at ~/.config/ontop/ontop.toml will be created automatically if it doesn't exist
- TOML format is suitable for storing user preferences (view mode, and potentially other settings in future)
- XDG Base Directory specification is followed for config file location (~/.config/)

## Out of Scope *(optional)*

- Customizing the number of workflow stages (remains fixed at three)
- Mixed mode where some stages are columns and others are rows
- Responsive layout that automatically switches based on terminal dimensions
- Custom keybindings for view mode toggle
- Animation or transitions when switching between view modes
- Drag-and-drop task movement (remains keyboard-only)
- Concurrent display of both column and row views
