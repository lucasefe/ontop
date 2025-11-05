# Research: Row-Based UI Mode

**Feature**: 004-row-ui | **Date**: 2025-11-05

## Overview

This document captures research findings for implementing the row-based UI mode feature, focusing on technical decisions around TOML configuration management, Bubble Tea view rendering patterns, and horizontal scrolling approaches.

## Research Questions

### 1. TOML Configuration Library Selection

**Question**: Which Go TOML library should we use for config persistence at `~/.config/ontop/ontop.toml`?

**Investigation**:
- Evaluated three main options:
  1. `github.com/BurntSushi/toml` - Original, stable, widely used
  2. `github.com/pelletier/go-toml` - Feature-rich with v2 offering better API
  3. Standard library (none exists for TOML)

**Decision**: Use `github.com/BurntSushi/toml` v1.x

**Rationale**:
- Most widely adopted (32k+ stars, used by many Go projects)
- Stable API (v1.x) with minimal breaking changes over years
- Simple unmarshal/marshal API matching encoding/json patterns
- Zero transitive dependencies (aligns with Simplicity First principle)
- Already proven in production across Go ecosystem
- Example usage:
  ```go
  import "github.com/BurntSushi/toml"

  type Config struct {
      ViewMode string `toml:"view_mode"`
  }

  // Read
  var cfg Config
  _, err := toml.DecodeFile(path, &cfg)

  // Write
  file, _ := os.Create(path)
  encoder := toml.NewEncoder(file)
  encoder.Encode(cfg)
  ```

**Alternatives Considered**:
- `go-toml/v2`: More features but heavier API, unnecessary for simple config
- Manual parsing: Too error-prone, reinventing wheel

**Impact**:
- Add `github.com/BurntSushi/toml` v1.4.0 to go.mod
- Config struct in `internal/config/config.go`
- XDG directory creation handled by `os.MkdirAll`

---

### 2. Row-Based Rendering Approach with Bubble Tea

**Question**: How should we structure the row-based view rendering to minimize duplication with existing column kanban.go?

**Investigation**:
- Reviewed existing `internal/tui/kanban.go` column rendering approach
- Examined Bubble Tea view composition patterns
- Considered three architectural options:
  1. Single View() function with if/else for layout
  2. Separate renderKanban() and renderRows() functions
  3. Strategy pattern with interface for layout renderers

**Decision**: Separate rendering functions with shared task card styling

**Rationale**:
- **Option 2 chosen** (separate render functions):
  - Clean separation of concerns without over-engineering
  - Kanban view stays untouched (no risk of regression)
  - Allows independent evolution of each layout
  - Easier to test each rendering path in isolation
  - Aligns with Simplicity First (no interfaces/abstractions needed)

- Architecture:
  ```go
  // In app.go View()
  func (m Model) View() string {
      switch m.viewMode {
      case ViewModeKanban:
          if m.viewLayout == LayoutColumn {
              return m.renderKanbanColumns()
          }
          return m.renderKanbanRows()
      // ... other modes
      }
  }

  // kanban.go (existing, unchanged)
  func (m Model) renderKanbanColumns() string { /* existing code */ }

  // row.go (new)
  func (m Model) renderKanbanRows() string { /* new implementation */ }
  ```

- Shared styling:
  - Extract card rendering to `renderTaskCard(task, isSelected, width) string`
  - Used by both column and row layouts
  - Maintains visual consistency

**Alternatives Considered**:
- Single function: Would make kanban.go harder to read, mixing two layouts
- Strategy pattern: Over-engineered for two simple layout variations

**Impact**:
- Create `internal/tui/row.go` for row rendering
- Extract `renderTaskCard()` helper (may create `styles.go` or keep in kanban.go)
- Add `ViewLayout` type and field to Model

---

### 3. Horizontal Scrolling for Overflow Tasks

**Question**: How should we handle rows with more tasks than fit horizontally on screen?

**Investigation**:
- Terminal width constraint: 80-200 characters (from success criteria)
- Task card minimum width: ~25 characters (for readable description)
- Maximum tasks visible: 3-8 depending on terminal width
- Bubble Tea doesn't provide horizontal scrolling widget out-of-box

**Decision**: Viewport windowing with left/right navigation

**Rationale**:
- Implement simple windowing logic in row rendering:
  ```
  [< task1 | task2 | task3 | task4 >]

  Visible window shifts as user navigates left/right:
  Window 1: task1, task2, task3 (task4 off-screen right)
  Window 2: task2, task3, task4 (task1 off-screen left)
  ```

- Technical approach:
  1. Track `rowScrollOffset[rowIndex]` in Model
  2. Calculate visible task slice based on terminal width and offset
  3. Show indicators: `[◀]` left, `[▶]` right when more tasks exist
  4. Left/right arrows adjust offset and wrap selection
  5. Offset resets when switching rows or leaving row mode

- Indicators:
  - Left edge: `◀` (U+25C0) when tasks exist left of viewport
  - Right edge: `▶` (U+25B6) when tasks exist right of viewport
  - Dim color (using existing Gruvbox theme gray)

**Alternatives Considered**:
- Marquee scrolling: Too complex, poor UX for task management
- Wrapping to multiple physical rows: Confusing (blurs row boundaries)
- No scrolling (error on overflow): Too restrictive for real usage

**Impact**:
- Add `rowScrollOffset map[int]int` to Model struct
- Implement viewport calculation in renderKanbanRows()
- Update left/right navigation handlers to manage offset
- Document in quickstart.md: scroll indicators and navigation

---

### 4. View Mode State Persistence

**Question**: When should config be written - on every toggle, or on app exit?

**Investigation**:
- Config operations: file I/O (potential failure points)
- Toggle frequency: Likely low (users settle on preference)
- Bubble Tea lifecycle: No guaranteed cleanup hook
- Terminal signals: SIGINT, SIGTERM may skip cleanup

**Decision**: Write config immediately on view mode toggle

**Rationale**:
- **Immediate write on toggle** (not on exit):
  - Terminal apps can be killed ungracefully (Ctrl+C, SIGKILL)
  - No reliable cleanup hook in Bubble Tea for all exit paths
  - Config write is fast (<1ms), doesn't impact UX
  - User sees preference persisted immediately
  - Simpler code (no deferred writes or exit hooks)

- Error handling:
  - Log error if write fails (non-fatal)
  - Continue with view toggle regardless
  - User can still use app, just preference not saved

**Alternatives Considered**:
- Write on exit: Risk losing preference on ungraceful shutdown
- Debounced write: Added complexity for minimal benefit

**Impact**:
- Call `config.Save()` in view toggle handler (app.go)
- Non-blocking, log errors only
- No exit hooks needed

---

### 5. Navigation State Preservation Between View Modes

**Question**: Should we preserve cursor position when toggling between column and row views?

**Investigation**:
- User expectation: Likely wants same task selected after toggle
- State mapping: Column index + task index → same in both views
- Edge cases: Task might not be in same visual position

**Decision**: Preserve selected task ID, reset to first position if not found

**Rationale**:
- Track selected task by ID (not index):
  ```go
  // Before toggle
  selectedID := m.getCurrentTaskID()

  // After toggle (in layout rendering)
  newIndex := m.findTaskIndex(selectedID, currentColumn)
  if newIndex == -1 {
      newIndex = 0  // Fallback to first task
  }
  ```

- Benefits:
  - Maintains user context when switching views
  - Graceful degradation (first task if ID not found)
  - Works across all view modes (column/row/detail)

- Reset scroll offset when toggling to avoid confusion

**Alternatives Considered**:
- Always reset to first task: Jarring UX, loses context
- Preserve exact visual position: Complex, not always meaningful

**Impact**:
- Add `getCurrentTaskID()` helper method
- Modify view toggle handler to preserve selection
- Document behavior in quickstart.md

---

### 6. Keyboard Navigation Consistency

**Question**: Should up/down arrows behave differently in row mode vs column mode?

**Investigation**:
- Column mode: up/down navigate tasks within column
- Row mode proposal: up/down switch between rows
- Vim consistency: j/k usually mean "move in primary dimension"
- User mental model: Rows are primary dimension in row mode

**Decision**: Context-sensitive navigation - arrows adapt to layout

**Rationale**:
- **In column mode** (existing):
  - up/down: Navigate tasks within current column
  - left/right: Switch between columns

- **In row mode** (new):
  - up/down: Switch between rows (inbox → in_progress → done)
  - left/right: Navigate tasks within current row

- Consistency principle:
  - "Arrow keys navigate the primary visual dimension"
  - Primary dimension is vertical in column mode
  - Primary dimension is horizontal in row mode
  - Maintains vim h/j/k/l muscle memory mapping

- Help text:
  - Update context-sensitive help (already shows different keys per mode)
  - Show current layout and navigation in help footer

**Alternatives Considered**:
- Fixed behavior: Confusing in row mode (up/down would scroll within row)
- Additional modifier keys: Violates Simplicity First

**Impact**:
- Navigation logic already isolated per view mode
- Update key handlers in app.go to check layout
- Update help text in keys.go
- Document in quickstart.md

---

## Technology Summary

| Component | Technology | Version | Rationale |
|-----------|------------|---------|-----------|
| Config Format | TOML | 1.0.0 | Simple, human-readable, standard |
| TOML Parser | BurntSushi/toml | v1.4.0 | Most stable, zero dependencies |
| TUI Framework | Bubble Tea | v1.2.4 | Existing (no change) |
| Styling | Lip Gloss | v1.0.0 | Existing (no change) |
| Config Location | XDG Base Dir | ~/.config/ontop/ | Cross-platform standard |

## Best Practices Applied

### Bubble Tea Rendering

1. **Stateless View Functions**: All render functions pure (no side effects)
2. **Update Separation**: State changes only in Update(), never in View()
3. **String Building**: Use strings.Builder for efficient concatenation
4. **Lip Gloss Composition**: Reuse styled components, avoid inline styling

### Configuration Management

1. **Atomic Writes**: Write to temp file, then rename (os.Rename)
2. **Directory Creation**: Use os.MkdirAll with 0755 permissions
3. **Error Handling**: Non-fatal config errors, log and continue
4. **Defaults**: Always provide fallback values (view_mode defaults to "column")

### State Management

1. **Single Source of Truth**: ViewLayout field determines rendering path
2. **Immutable Toggles**: Toggle produces new state, no mutation
3. **Reset on Mode Change**: Clear scroll offsets when switching layouts
4. **ID-Based Selection**: Track by task ID, not array indices

## Dependencies Added

```go
// Add to go.mod
require (
    github.com/BurntSushi/toml v1.4.0
)
```

## Integration Points

### With Existing Code

- `internal/tui/model.go`: Add ViewLayout field and scroll offset map
- `internal/tui/app.go`: Add toggle handler, route to correct renderer
- `internal/tui/kanban.go`: Extract card rendering helper (minor refactor)
- `internal/tui/keys.go`: Add 'v' key binding, update help text

### New Code

- `internal/config/config.go`: Config struct, Load/Save functions
- `internal/tui/row.go`: Row-based rendering, viewport windowing

## Risk Assessment

| Risk | Likelihood | Impact | Mitigation |
|------|------------|--------|------------|
| Config write failure | Low | Low | Non-fatal, log error, app continues |
| TOML library breaking change | Very Low | Low | v1.x stable for years, lock version |
| Rendering performance | Low | Medium | Profile if needed, optimize string building |
| Navigation confusion | Medium | Medium | Clear help text, context indicators |
| Scroll offset bugs | Medium | Low | Reset on view change, comprehensive testing |

## Open Questions

None. All research questions resolved.

## References

- BurntSushi/toml: https://github.com/BurntSushi/toml
- Bubble Tea docs: https://github.com/charmbracelet/bubbletea
- XDG Base Directory: https://specifications.freedesktop.org/basedir-spec/latest/
- TOML spec: https://toml.io/en/v1.0.0
