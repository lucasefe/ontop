# TUI Rendering Contract

**Feature**: 004-row-ui | **Module**: internal/tui | **Date**: 2025-11-05

## Overview

This contract defines the rendering interface and behavior for the new row-based layout mode, ensuring consistency with the existing column-based rendering while maintaining code separation.

## Type Definitions

### ViewLayout Enum

```go
// ViewLayout represents the visual organization of the kanban board
type ViewLayout int

const (
    LayoutColumn ViewLayout = iota  // Vertical columns (existing)
    LayoutRow                       // Horizontal rows (new)
)
```

**Contract**:
- Type-safe enum (compiler-enforced)
- Values are immutable constants
- Used to determine rendering path in View()

### Model Extensions

```go
// Model represents the Bubbletea application state
type Model struct {
    // ... existing fields (unchanged) ...

    // NEW: Layout mode for Kanban view
    viewLayout ViewLayout

    // NEW: Horizontal scroll offsets per row (row mode only)
    // Key: row index (0=inbox, 1=in_progress, 2=done)
    // Value: leftmost visible task index
    rowScrollOffset map[int]int
}
```

## Public Functions

### Rendering Functions

#### renderKanbanColumns (Existing)

```go
// renderKanbanColumns renders the traditional vertical column layout.
// UNCHANGED from current implementation.
func (m Model) renderKanbanColumns() string
```

**Contract**:
- Returns fully-formatted terminal string
- Uses existing column navigation logic
- No behavior changes (stability guaranteed)

#### renderKanbanRows (New)

```go
// renderKanbanRows renders the horizontal row layout where each
// workflow stage (inbox, in_progress, done) is displayed as a
// horizontal row with tasks listed left-to-right.
//
// Returns:
//   - string: Formatted terminal output with rows, scroll indicators,
//             and help text
func (m Model) renderKanbanRows() string
```

**Behavior**:
- Renders three horizontal rows (one per workflow stage)
- Each row shows workflow stage name + tasks laid out horizontally
- Implements viewport windowing for horizontal scrolling
- Shows `[◀]` and `[▶]` indicators when tasks overflow viewport
- Highlights currently selected row and task
- Returns string compatible with Bubble Tea's View() contract

**Layout Structure**:
```
╔═══════════════════════════════════════════╗
║ INBOX                              (3)    ║
║ [◀] │ Task 1 │ Task 2 │ Task 3 │ [▶]     ║
╚═══════════════════════════════════════════╝
╔═══════════════════════════════════════════╗
║ IN PROGRESS                        (2)  * ║  ← Selected row
║     │ Task 4 │ Task 5 │                   ║
╚═══════════════════════════════════════════╝
╔═══════════════════════════════════════════╗
║ DONE                               (1)    ║
║     │ Task 6 │                            ║
╚═══════════════════════════════════════════╝

* = selected row indicator
```

#### renderTaskCard (Extracted Helper)

```go
// renderTaskCard renders a single task as a styled card.
// Used by both column and row layouts for consistency.
//
// Parameters:
//   - task: The task to render
//   - isSelected: Whether this task is currently selected
//   - width: Maximum width for the card (characters)
//
// Returns:
//   - string: Styled task card (Lip Gloss formatted)
func renderTaskCard(task *models.Task, isSelected bool, width int) string
```

**Contract**:
- Visual consistency across both layouts
- Handles text truncation for width constraint
- Applies selection styling (border, color) via Lip Gloss
- Shows task priority, description, progress if >0
- Returns single-line or multi-line string depending on layout needs

### View Routing

```go
// View implements the Bubble Tea View() contract.
// Routes to appropriate renderer based on viewMode and viewLayout.
func (m Model) View() string {
    switch m.viewMode {
    case ViewModeKanban:
        if m.viewLayout == LayoutRow {
            return m.renderKanbanRows()
        }
        return m.renderKanbanColumns()
    case ViewModeDetail:
        return m.renderDetail()  // Unchanged
    // ... other modes unchanged
    }
}
```

## Update Handlers

### View Toggle Handler

```go
// handleToggleView toggles between column and row layouts.
// Triggered by 'v' key in Kanban mode.
//
// Behavior:
//   - Toggles m.viewLayout (column ↔ row)
//   - Resets m.rowScrollOffset to empty map
//   - Preserves selected task by ID (best effort)
//   - Saves preference to config file (async)
//   - Returns updated model (no command)
func (m Model) handleToggleView() Model
```

**Side Effects**:
- Writes config file asynchronously (non-blocking)
- Logs error if config save fails (non-fatal)

### Navigation Handler Updates

```go
// handleKanbanKeys processes keyboard input in Kanban mode.
// Navigation behavior depends on m.viewLayout:
//
// Column Mode:
//   - up/down: Navigate tasks within column
//   - left/right: Switch columns
//
// Row Mode:
//   - up/down: Switch rows (workflow stages)
//   - left/right: Navigate tasks within row (with scrolling)
func (m Model) handleKanbanKeys(msg tea.KeyMsg) (Model, tea.Cmd)
```

**Contract**:
- Navigation is context-sensitive (adapts to layout)
- No mode confusion (help text always shows correct keys)
- Preserves task selection when possible

## Rendering Specifications

### Row Layout Dimensions

| Element | Size | Behavior |
|---------|------|----------|
| Row header width | 20 chars | Fixed, right-padded |
| Task card min width | 20 chars | Minimum readable width |
| Task card max width | 40 chars | Prevents excessive width |
| Row padding | 2 chars | Left/right margins |
| Task separator | ` │ ` | Vertical bar with spaces |
| Scroll indicator | `[◀]` `[▶]` | 3 chars each |

### Viewport Calculation

```go
// Pseudocode for viewport windowing
availableWidth := terminalWidth - headerWidth - padding - indicatorSpace
taskCardWidth := 25  // Target width per card
visibleCount := max(1, availableWidth / (taskCardWidth + separatorWidth))

startIndex := rowScrollOffset[currentRow]
endIndex := min(startIndex + visibleCount, len(tasks))

visibleTasks := tasks[startIndex:endIndex]
showLeftIndicator := startIndex > 0
showRightIndicator := endIndex < len(tasks)
```

### Styling (Lip Gloss)

**Row Container**:
```go
rowStyle := lipgloss.NewStyle().
    Border(lipgloss.RoundedBorder()).
    BorderForeground(lipgloss.Color("#665c54")).  // Gruvbox gray
    Padding(0, 1).
    Width(terminalWidth - 4)
```

**Selected Row**:
```go
selectedRowStyle := rowStyle.Copy().
    BorderForeground(lipgloss.Color("#b8bb26")).  // Gruvbox green
    Bold(true)
```

**Task Card** (same for both layouts):
```go
taskStyle := lipgloss.NewStyle().
    Border(lipgloss.NormalBorder()).
    BorderForeground(lipgloss.Color("#928374")).  // Gruvbox gray
    Padding(0, 1).
    Width(cardWidth)

selectedTaskStyle := taskStyle.Copy().
    BorderForeground(lipgloss.Color("#fabd2f")).  // Gruvbox yellow
    Bold(true)
```

**Scroll Indicators**:
```go
indicatorStyle := lipgloss.NewStyle().
    Foreground(lipgloss.Color("#928374")).  // Gruvbox gray (dim)
    Faint(true)
```

## State Management

### Scroll Offset Lifecycle

```
Initialization:
    rowScrollOffset = make(map[int]int)  // Empty map, all offsets implicitly 0

Left Arrow (Row Mode):
    offset := m.rowScrollOffset[m.currentColumn]  // Default 0 if not set
    m.rowScrollOffset[m.currentColumn] = max(0, offset - 1)
    m.selectedTask = max(0, m.selectedTask - 1)

Right Arrow (Row Mode):
    tasks := m.GetTasksByColumn(currentColumnName)
    maxOffset := max(0, len(tasks) - visibleCount)
    offset := m.rowScrollOffset[m.currentColumn]
    m.rowScrollOffset[m.currentColumn] = min(maxOffset, offset + 1)
    m.selectedTask = min(len(tasks)-1, m.selectedTask + 1)

Layout Toggle:
    m.rowScrollOffset = make(map[int]int)  // Reset all offsets

Row Change (Up/Down):
    // Preserve scroll offsets for each row independently
    // Only update currentColumn, don't modify rowScrollOffset
```

### Task Selection Preservation

```go
// When toggling layout, preserve selected task if possible
func (m Model) preserveTaskSelection() Model {
    selectedID := m.getCurrentTaskID()
    // After layout change and re-render
    newIndex := m.findTaskIndexByID(selectedID)
    if newIndex != -1 {
        m.selectedTask = newIndex
    } else {
        m.selectedTask = 0  // Fallback to first task
    }
    return m
}
```

## Performance Requirements

### Rendering Performance

| Metric | Target | Maximum | Measurement |
|--------|--------|---------|-------------|
| View() latency | <16ms | 50ms | Frame time at 60fps |
| Layout toggle | <100ms | 500ms | User-perceptible delay |
| Scroll operation | <16ms | 50ms | Smooth navigation |

**Optimization Strategies**:
- Use `strings.Builder` for string concatenation
- Cache Lip Gloss styles (don't recreate each frame)
- Limit rendered tasks to visible viewport only

### Memory Usage

| Component | Size | Justification |
|-----------|------|---------------|
| rowScrollOffset map | ~100 bytes | 3 rows × ~32 bytes per entry |
| Rendered strings | <10KB | Terminal output buffer |
| Style cache | <1KB | Reused Lip Gloss styles |

## Testing Contract

### Unit Tests (internal/tui/)

1. **renderKanbanRows()**:
   - Empty rows (no tasks)
   - Single task per row
   - Overflow tasks (triggers scroll indicators)
   - Narrow terminal (minimum width handling)
   - Wide terminal (maximum visible tasks)

2. **Viewport windowing**:
   - Calculate visible task count correctly
   - Show left indicator when offset > 0
   - Show right indicator when more tasks exist
   - Bounds checking (offset never negative, never exceeds max)

3. **Layout toggle**:
   - Column → Row → Column preserves data
   - Scroll offsets reset on toggle
   - Task selection preserved (when possible)
   - Config save triggered (mock verification)

### Visual Regression Tests

Use snapshot testing or manual verification:

1. **Layout consistency**:
   - Task cards look identical in both layouts
   - Borders and colors match theme
   - Selection highlighting works in both modes

2. **Responsive behavior**:
   - Terminal resize updates viewport correctly
   - Narrow terminals don't break layout
   - Wide terminals show more tasks

## Compatibility Guarantees

### Backwards Compatibility

- Column mode rendering unchanged (no visual differences)
- Existing key bindings work identically in column mode
- No changes to task data, storage, or CLI

### Forward Compatibility

- Additional ViewLayout values MAY be added (e.g., LayoutGrid)
- renderKanbanRows() MAY add features without breaking contract
- Model fields MAY be extended (additive changes only)

## Error Handling

### Rendering Errors

| Error Scenario | Handling |
|----------------|----------|
| Terminal too narrow | Show warning message, fallback to minimal layout |
| No tasks in row | Show "No tasks" placeholder |
| Invalid scroll offset | Clamp to valid range [0, maxOffset] |
| Style rendering failure | Fallback to unstyled text (graceful degradation) |

**Contract**: View() MUST NEVER panic. Always return valid string, even if degraded.

## Integration Points

### With Config Module

```go
// On startup
cfg, _ := config.Load()
model.viewLayout = parseViewMode(cfg.UI.ViewMode)  // "column" → LayoutColumn

// On toggle
model.viewLayout = toggleLayout(model.viewLayout)
cfg.UI.ViewMode = formatViewMode(model.viewLayout)  // LayoutRow → "row"
config.Save(cfg)  // Async, errors logged only
```

### With Bubble Tea

```go
// Implements tea.Model interface
func (m Model) Init() tea.Cmd { /* unchanged */ }
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) { /* route based on layout */ }
func (m Model) View() string { /* route to renderKanbanColumns or renderKanbanRows */ }
```

### With Storage/Models

- No changes to storage layer
- No changes to Task or Column models
- Read-only access to m.tasks (no mutations in View())

## Documentation Requirements

### Code Comments

All new/modified functions MUST have:
- Purpose description
- Parameter documentation
- Return value description
- Behavior notes (especially context-sensitive logic)

### Help Text

```go
// keys.go: Update help text based on layout
func (k KeyMap) ShortHelp() []key.Binding {
    if currentLayout == LayoutRow {
        return []key.Binding{
            k.Up,    // "↑/k: previous row"
            k.Down,  // "↓/j: next row"
            k.Left,  // "←/h: previous task"
            k.Right, // "→/l: next task"
            k.Toggle, // "v: column view"
        }
    }
    // ... column mode help (existing)
}
```

## Version History

| Version | Date | Changes |
|---------|------|---------|
| 1.0 | 2025-11-05 | Initial contract - Row layout rendering |

## Future Enhancements (Not in Scope)

Potential future additions that don't break this contract:

- **LayoutGrid**: 2D grid view (rows × columns)
- **Custom row ordering**: User-defined row sequence
- **Row collapsing**: Hide/show individual rows
- **Multi-select**: Select multiple tasks across rows
- **Drag-and-drop**: Mouse-based task movement (if mouse support added)

All enhancements MUST maintain backwards compatibility with LayoutColumn and LayoutRow.
