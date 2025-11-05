# Feature Specification: README Visual Enhancement with Screenshots

**Feature Branch**: `003-readme-visual-enhancement`
**Created**: 2025-11-05
**Status**: Draft
**Input**: User description: "i want to have an amazing readme file, ideally with an image of ontop showing a bunch of sample tasks. It should show how it cam be used from cli and from the TUI."

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Visual TUI Demonstration (Priority: P1)

New users or contributors visiting the GitHub repository need to quickly understand what OnTop looks like and how it works without installing it. A screenshot of the TUI showing the kanban board with sample tasks provides immediate visual context about the application's interface and capabilities.

**Why this priority**: The TUI is the primary interface and most distinctive feature of OnTop. Without visual proof, users cannot assess if the tool meets their needs or if it's worth installing. This is the first impression for potential users.

**Independent Test**: Can be fully tested by viewing the README on GitHub and verifying that a clear, high-quality screenshot of the TUI kanban board is visible with sample tasks showing hierarchy, columns, and keyboard shortcuts. Delivers immediate value by showing the tool in action.

**Acceptance Scenarios**:

1. **Given** a user visits the GitHub repository, **When** they view the README, **Then** they see a screenshot of the TUI showing the kanban board with at least 5 sample tasks
2. **Given** the screenshot displays the TUI, **When** user examines it, **Then** they can identify all three columns (Inbox, In Progress, Done) with tasks in each
3. **Given** the screenshot shows sample tasks, **When** user looks closely, **Then** they can see parent-child task relationships with subtasks indented below parents
4. **Given** the screenshot shows the interface, **When** user reviews it, **Then** keyboard shortcuts or help panel is visible showing how to interact with the TUI

---

### User Story 2 - CLI Usage Examples with Output (Priority: P1)

Users evaluating OnTop for scripting or automation need to see actual CLI command outputs to understand the format and structure of responses. Showing real terminal output examples demonstrates both the command syntax and the expected results.

**Why this priority**: CLI mode is essential for automation workflows. Users need to know exactly what output format to expect for parsing in scripts. This is equally critical as TUI visualization for users who prefer command-line workflows.

**Independent Test**: Can be fully tested by viewing the README CLI section and verifying that each command example includes both the command and actual formatted output showing tasks in hierarchical structure. Delivers value by showing exact CLI behavior.

**Acceptance Scenarios**:

1. **Given** user views the CLI section, **When** they read command examples, **Then** each example shows both the command and sample output
2. **Given** CLI examples include `list` command, **When** output is shown, **Then** it displays hierarchical task structure with indented subtasks using dash prefix
3. **Given** CLI examples include `show` command, **When** output is shown, **Then** it displays full task details including all subtasks
4. **Given** CLI examples include `add` command, **When** output is shown, **Then** it demonstrates both parent task and subtask creation with `--parent` flag

---

### User Story 3 - Feature Showcase Section (Priority: P2)

Users browsing the README need a prominent visual section that highlights key features with supporting screenshots or GIFs showing those features in action. This helps users quickly identify which capabilities matter to them.

**Why this priority**: After seeing the basic interface, users need to understand specific capabilities like subtask hierarchy, multiline descriptions, vim navigation, and archive system. Visual proof of features increases confidence and reduces "does it do X?" questions.

**Independent Test**: Can be fully tested by reviewing a dedicated "Features Showcase" section that combines feature descriptions with corresponding screenshots or terminal recordings demonstrating each feature. Delivers value by proving capabilities visually.

**Acceptance Scenarios**:

1. **Given** README has a Features Showcase section, **When** user views it, **Then** each major feature (subtasks, multiline descriptions, vim navigation, archives) has a visual example
2. **Given** subtask feature is showcased, **When** user views the screenshot, **Then** they can clearly see parent-child relationships with visual indentation
3. **Given** form/editing feature is showcased, **When** user views the screenshot, **Then** they can see the multiline textarea with scrollable content
4. **Given** user wants to understand dual-mode architecture, **When** they review the showcase, **Then** side-by-side or sequential examples show the same task workflow in both TUI and CLI modes

---

### User Story 4 - Quick Start Visual Guide (Priority: P3)

Users who just installed OnTop need a visual quick start guide showing the first actions to take. A screenshot or short animated GIF showing "create your first task" workflow helps users get started immediately.

**Why this priority**: After installation, users benefit from a guided first experience. While less critical than understanding capabilities (P1-P2), this reduces friction for first-time users and improves initial success rate.

**Independent Test**: Can be fully tested by following the Quick Start section and verifying it includes a visual walkthrough (screenshot or GIF) showing the task creation form with sample data being entered. Delivers value by reducing onboarding time.

**Acceptance Scenarios**:

1. **Given** user reads Quick Start section, **When** they view it, **Then** it includes a screenshot showing the task creation form
2. **Given** the form screenshot is displayed, **When** user examines it, **Then** all fields (title, description, priority, progress, tags, parent) are visible with example data
3. **Given** user follows the Quick Start, **When** they complete it, **Then** they understand how to launch TUI, create a task, and navigate the board
4. **Given** user prefers CLI, **When** they review Quick Start, **Then** an alternative CLI-based quick start workflow is provided

---

### Edge Cases

- **Screenshot staleness**: What happens when UI changes and screenshots become outdated? Document in a comment in README that screenshots should be regenerated after UI changes, with reference to the sample data script used.
- **Image size and performance**: How do we ensure screenshots don't bloat the repository? Use optimized PNG format and keep total image size under 2MB combined.
- **Accessibility**: How do screen readers handle screenshots? Include descriptive alt text for each image describing what's shown.
- **Dark/Light theme**: Which theme should screenshots use? Assume dark theme (Gruvbox) as it's the current implementation and more visually distinctive.
- **Different terminal sizes**: What if user's terminal doesn't match screenshot dimensions? Include note about recommended terminal size (80x24 minimum) and that UI adapts to terminal size.
- **GIF vs static images**: Should we use animated GIFs or static screenshots? Prefer static screenshots for simplicity and file size, but allow GIFs for demonstrating workflows if tool integration exists.

## Requirements *(mandatory)*

### Functional Requirements

#### Visual Assets

- **FR-001**: README MUST include at least one high-quality screenshot of the TUI kanban board showing all three columns with sample tasks
- **FR-002**: README MUST include a screenshot showing hierarchical subtask display with parent tasks and indented subtasks
- **FR-003**: README MUST include CLI command examples with actual formatted output showing task lists and details
- **FR-004**: README MUST include a screenshot of the task creation/edit form showing all input fields with sample data
- **FR-005**: All screenshots MUST use actual OnTop application running with sample data (not mockups or edited images)

#### Content Requirements

- **FR-006**: README MUST have a dedicated "Features Showcase" or similar section prominently displaying key capabilities with visual examples
- **FR-007**: CLI examples section MUST show both the command syntax and expected output for each command (add, list, show, move, update)
- **FR-008**: README MUST include descriptive alt text for all images for accessibility
- **FR-009**: README MUST document the sample data setup used for screenshots to enable reproduction after UI changes

#### Sample Data

- **FR-010**: Sample data MUST include at least 5 parent tasks distributed across all three columns
- **FR-011**: Sample data MUST include at least 3 subtasks demonstrating parent-child relationships
- **FR-012**: Sample data MUST demonstrate various priority levels (1-5) and progress percentages (0-100)
- **FR-013**: Sample data MUST include tasks with tags and multiline descriptions to showcase all features
- **FR-014**: Sample data MUST be reproducible via script or documented commands for future screenshot updates

#### Technical Requirements

- **FR-015**: Screenshots MUST be stored in repository under `docs/images/` or similar directory
- **FR-016**: Total size of all screenshot images MUST not exceed 2MB to avoid repository bloat
- **FR-017**: Images MUST use PNG format for screenshots (lossless, good text rendering)
- **FR-018**: Images MUST be committed to the repository (not external hosting) for README self-containment

### Key Entities

- **Screenshot Asset**: An image file showing the OnTop application in use (TUI or CLI output)
- **Sample Task Dataset**: A predefined set of tasks used consistently across all visual examples
- **CLI Output Example**: Terminal text showing command execution and formatted results
- **Feature Showcase Entry**: A combination of feature description + supporting visual evidence

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: Users can understand OnTop's primary interface and capabilities within 30 seconds of viewing the README without reading full text documentation
- **SC-002**: At least 3 distinct screenshots are visible in the README showing different aspects (kanban board, CLI output, form interface)
- **SC-003**: Every major feature listed in the Features section has at least one corresponding visual example (screenshot or code block with output)
- **SC-004**: CLI examples demonstrate at least 5 different commands with actual formatted output
- **SC-005**: Screenshot showing hierarchical subtasks clearly displays at least 2 levels of indentation (parent → subtask)
- **SC-006**: 100% of screenshots include descriptive alt text that conveys the image content to screen reader users
- **SC-007**: Users can reproduce the sample data and screenshots using documented commands or scripts
- **SC-008**: README loads on GitHub within 3 seconds despite including multiple images (total image size under 2MB)

## Assumptions *(optional)*

- OnTop can be run with a custom database path to create isolated sample data without affecting user's production tasks
- Terminal screenshots can be captured using standard tools (macOS: Command+Shift+4, Linux: screenshot tools, or terminal recording tools like `asciinema` for static export)
- The Gruvbox dark theme is the default and preferred theme for screenshots
- Sample data script or commands will be documented but do not need automated tooling (manual setup is acceptable)
- GIF animations are optional and not required for MVP - static screenshots are sufficient
- Images will be stored directly in the repository, not using external image hosting services
- README is primarily viewed on GitHub's web interface where markdown image syntax is fully supported

## Open Questions

None - all aspects have reasonable defaults based on common documentation practices and existing OnTop capabilities.
