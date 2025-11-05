# Research: README Visual Enhancement

**Feature**: README Visual Enhancement with Screenshots
**Date**: 2025-11-05
**Purpose**: Research best practices for technical documentation, screenshot capture, and markdown image optimization

## Research Areas

### 1. README Documentation Best Practices

**Decision**: Use visual-first approach with screenshots prominently placed near feature descriptions

**Rationale**:
- GitHub users scan documentation visually before reading detailed text
- Screenshots provide immediate proof of capabilities without requiring installation
- Visual examples reduce "does it do X?" questions in issues
- Industry standard: successful CLI/TUI projects (lazygit, btop, lazydocker) use prominent screenshots

**Alternatives considered**:
- Text-only documentation: Rejected - requires users to imagine interface, increases barrier to evaluation
- Video demos/GIFs: Considered but deferred - larger file sizes, more complex to update, static screenshots sufficient for MVP
- External image hosting (imgur, etc.): Rejected - external dependencies, link rot risk, self-contained repo preferred

**Best practices identified**:
- Place hero screenshot (main interface) near top of README after brief description
- Include alt text for accessibility (screen readers)
- Use descriptive filenames (tui-kanban-board.png not screenshot1.png)
- Group related screenshots in dedicated sections
- Keep total image payload under 2MB for fast GitHub rendering

**References**:
- lazygit README structure: prominent TUI screenshot with demo GIF
- GitHub's own documentation guidelines emphasize visual examples
- Accessibility guidelines require descriptive alt text

---

### 2. Screenshot Capture Techniques for Terminal Applications

**Decision**: Use native OS screenshot tools (macOS: Command+Shift+4, Linux: screenshot utils) capturing actual terminal windows

**Rationale**:
- Native tools provide highest quality PNG output
- Captures actual rendered terminal with correct colors and fonts
- No external dependencies or tools to install
- Reproducible by any contributor with standard OS tools

**Alternatives considered**:
- asciinema/termtosvg: Considered for animated recordings - deferred to future, static screenshots sufficient for MVP
- Terminal-to-image libraries (carbon.now.sh): Rejected - adds external dependency, not actual application rendering
- Manual mockups: Rejected - violates FR-005 requirement for actual application screenshots

**Best practices identified**:
- Use terminal size of at least 120x30 for readability in screenshots
- Ensure Gruvbox theme is active (dark theme provides better contrast)
- Capture during stable state (not mid-animation or with flickering cursor)
- Use PNG format for lossless text rendering (not JPEG which causes artifacting)
- Optimize PNGs after capture using built-in compression (pngcrush, optipng if available)

**Terminal size recommendations**:
- Kanban board: 140x35 to show all three columns comfortably
- Task form: 100x40 to show full form with all fields visible
- Detail view: 120x40 to show task details and subtask list

---

### 3. Sample Data Generation Strategy

**Decision**: Create shell script that uses OnTop CLI to populate database with predefined sample tasks

**Rationale**:
- Uses existing CLI commands (no new code needed)
- Reproducible - script can be run again after UI changes
- Custom database path isolates sample data from user's production tasks
- Script serves as documentation for sample data structure

**Alternatives considered**:
- Manual task creation: Rejected - not reproducible, time-consuming for future updates
- SQL insertion: Rejected - bypasses application logic, fragile to schema changes
- JSON fixture import: Rejected - would require new import feature, violates simplicity principle

**Sample data requirements (from spec FR-010 through FR-013)**:
- 5+ parent tasks distributed across columns:
  - 2 in Inbox
  - 2 in In Progress
  - 1 in Done
- 3+ subtasks demonstrating hierarchy:
  - At least one parent with 2+ subtasks
  - Subtasks in different columns than parents
- Priority diversity: tasks with P1, P2, P3, P4, P5
- Progress diversity: 0%, 25%, 50%, 75%, 100%
- Tags: development, bug, feature, documentation
- Multiline descriptions: at least 2 tasks with 3+ line descriptions

**Script structure**:
```bash
#!/bin/bash
# Use custom database to avoid polluting user data
export DB_PATH="./sample-tasks.db"

# Create parent tasks
./ontop add "Implement user authentication" --priority 1 --description "Add OAuth2 and session management
Support GitHub and Google providers
Include remember-me functionality" --tags "feature,security"

./ontop add "Fix memory leak in task sync" --priority 2 --progress 75 --tags "bug"

# Create subtasks
PARENT_ID=$(./ontop list --column inbox | grep "authentication" | awk '{print $1}')
./ontop add "Research OAuth2 libraries" --parent $PARENT_ID --priority 2 --progress 100

# ... more sample data ...
```

---

### 4. Image Optimization and File Size Management

**Decision**: Use PNG format with basic compression, target <500KB per screenshot, <2MB total

**Rationale**:
- PNG provides lossless compression suitable for text rendering
- GitHub handles PNG efficiently in markdown rendering
- Modern PNG compression (level 9) provides good size reduction without quality loss
- 2MB total is reasonable for 4-5 screenshots while keeping README load time under 3 seconds

**Alternatives considered**:
- WebP format: Rejected - not universally supported in all markdown viewers
- JPEG: Rejected - lossy compression causes text artifacting
- SVG for terminal output: Rejected - complex to generate, not actual screenshots

**Optimization workflow**:
1. Capture screenshot with native tools (produces uncompressed PNG)
2. Optional: Run through optipng or pngcrush if available (not required)
3. Verify file size < 500KB
4. If oversized, reduce terminal window size slightly and recapture

**File size benchmarks**:
- Terminal screenshot 120x30: typically 150-300KB
- Terminal screenshot 140x35: typically 200-400KB
- With Gruvbox dark theme: slightly smaller due to dark background compression

---

### 5. CLI Output Documentation Format

**Decision**: Use markdown code blocks with bash syntax highlighting for commands and plain text blocks for output

**Rationale**:
- Standard markdown approach familiar to developers
- GitHub renders bash syntax highlighting automatically
- Clear visual distinction between command and output
- Copyable - users can copy commands directly

**Format example**:
````markdown
```bash
$ ontop list
```

```
Inbox (2 tasks)
  [abc123] P1: Implement user authentication (0%)
    - [def456] P2: Research OAuth2 libraries (100%)
  [ghi789] P3: Update documentation (25%)

In Progress (1 task)
  [jkl012] P2: Fix memory leak in task sync (75%)
```
````

**Alternatives considered**:
- Screenshot of terminal output: Rejected - not searchable, not copyable, accessibility issues
- Inline text without code blocks: Rejected - poor readability, no syntax distinction
- Table format: Rejected - doesn't match actual CLI output format, misleading

**Best practices identified**:
- Show actual command with $ prompt to indicate it's runnable
- Include both successful outputs and demonstrate error handling if relevant
- Use realistic task IDs (not abc123 in actual examples)
- Show at least 3-5 tasks to demonstrate hierarchy and variety

---

### 6. README Structure and Organization

**Decision**: Restructure README with visual-first sections: Hero screenshot → Features showcase → Installation → Usage (TUI/CLI) → Development

**Rationale**:
- Visual elements grab attention before walls of text
- Progressive disclosure: overview → specifics → contributor info
- Matches successful open-source project patterns (README analysis of top Go TUI projects)

**Section flow**:
1. **Title + Description** (existing, 1-2 sentences)
2. **Hero Screenshot** (NEW - kanban board with sample tasks)
3. **Features** (existing list, enhanced with inline screenshots)
4. **Installation** (existing, minimal changes)
5. **Quick Start** (NEW - visual walkthrough)
6. **Usage - TUI Mode** (existing, enhanced with screenshots + keyboard reference)
7. **Usage - CLI Mode** (existing, enhanced with command output examples)
8. **Development** (existing, move to bottom)

**Features Showcase approach**:
- Each major feature gets 2-3 sentences + visual proof (screenshot or code block)
- Hierarchical subtasks: screenshot showing indentation
- Multiline descriptions: screenshot of form with textarea
- Dual-mode: side-by-side CLI + TUI examples

---

### 7. Accessibility Considerations

**Decision**: Include descriptive alt text for all images following WCAG guidelines

**Rationale**:
- Screen reader users need text description of visual content
- Alt text improves SEO and provides fallback if images fail to load
- Demonstrates inclusive design practices
- Required by FR-008

**Alt text guidelines**:
- Describe what the image shows, not that it's a screenshot
- Include key information visible in image (columns, tasks, hierarchy)
- Keep under 125 characters when possible
- Format: `![Alt text describing content](path/to/image.png)`

**Examples**:
- Good: `![OnTop TUI showing kanban board with three columns (Inbox, In Progress, Done) containing tasks with subtasks indented below parent tasks]`
- Bad: `![Screenshot]` or `![OnTop]`

**Markdown syntax**:
```markdown
![OnTop TUI kanban board with hierarchical tasks](docs/images/tui-kanban-board.png)
```

---

## Research Summary

All research areas completed with clear decisions. No NEEDS CLARIFICATION items remain - all technical approaches use existing tools and established best practices.

**Key takeaways**:
1. Visual-first documentation with screenshots near feature descriptions
2. Native OS tools for screenshot capture (no external dependencies)
3. Shell script using OnTop CLI for reproducible sample data
4. PNG format with basic compression, <2MB total
5. Markdown code blocks for CLI examples with clear command/output distinction
6. Restructured README flow: visual → features → usage → development
7. Descriptive alt text for accessibility

**Dependencies identified**:
- None - all work uses standard tools (bash, OS screenshot utilities, markdown)

**Next phase**: Generate data-model.md (though minimal for documentation feature) and contracts (also minimal - sample data structure documentation)
