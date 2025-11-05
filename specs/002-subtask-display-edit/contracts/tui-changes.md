# Internal API Contract: TUI Changes

**Purpose**: Define changes to existing TUI components for hierarchical display and subtask creation
**Affected Files**: `internal/tui/model.go`, `internal/tui/app.go`, `internal/tui/kanban.go`, `internal/tui/detail.go`, `internal/tui/form.go`

---

## 1. Model Changes (`internal/tui/model.go`)

### Modified Function: GetTasksByColumn

**Before**:
```go
func (m *Model) GetTasksByColumn(column string) []*models.Task
```

**After**:
```go
func (m *Model) GetTasksByColumn(column string) []*models.Task
```

**Changes**:
- Return tasks in hierarchical display order (parents, then their subtasks)
- Internally calls `hierarchy.BuildFlatHierarchy(filtered, m.sortMode)`
- External signature unchanged for compatibility

**Implementation**:
```go
func (m *Model) GetTasksByColumn(column string) []*models.Task {
    // Filter by column (existing logic)
    var filtered []*models.Task
    for _, task := range m.tasks {
        if task.Column == column {
            filtered = append(filtered, task)
        }
    }

    // Build hierarchical display order (NEW)
    hierarchical := hierarchy.BuildFlatHierarchy(filtered, m.sortMode)

    // Extract just the Task pointers (flatten wrapper)
    result := make([]*models.Task, len(hierarchical))
    for i, ht := range hierarchical {
        result[i] = ht.Task
    }

    return result
}
```

**Behavior Change**:
- Before: Returns tasks sorted by sortMode globally
- After: Returns tasks grouped hierarchically (parents sorted, subtasks grouped under parents)

**Backward Compatibility**: Callers still receive `[]*models.Task`, selection index math unchanged

---

## 2. App (Key Handlers) Changes (`internal/tui/app.go`)

### Modified Function: handleDetailKeys

**Purpose**: Add 'n' key handler to create subtask from detail view

**Before**:
```go
func (m Model) handleDetailKeys(msg tea.KeyMsg, keys KeyMap) (tea.Model, tea.Cmd) {
    // ... existing handlers (back, edit, archive, delete)
}
```

**After**:
```go
func (m Model) handleDetailKeys(msg tea.KeyMsg, keys KeyMap) (tea.Model, tea.Cmd) {
    // ... existing handlers ...

    // NEW: Create subtask (press 'n' in detail view)
    if key.Matches(msg, keys.New) {
        if m.detailTask != nil {
            m.viewMode = ViewModeCreate
            // Pass parent ID to pre-fill form
            m.initCreateForm(&m.detailTask.ID)  // NEW: pass parent ID
        }
        return m, nil
    }

    return m, nil
}
```

**Key Behavior**:
- Pressing 'n' while viewing task details opens create form
- Parent ID field pre-populated with current task's ID
- User can still edit/clear parent ID if desired

**Edge Case**: If `detailTask` is nil (shouldn't happen), no-op (safe)

---

## 3. Kanban Rendering (`internal/tui/kanban.go`)

### Modified Function: renderTaskCard

**Purpose**: Add visual indicator for subtasks (indentation + dash prefix)

**Before**:
```go
func (m Model) renderTaskCard(task *models.Task, maxWidth int) string {
    priorityStr := priorityStyle.Render(fmt.Sprintf("P%d", task.Priority))
    title := task.Title
    // ... truncation logic ...
    return fmt.Sprintf("%s %s %s", priorityStr, title, createdStr)
}
```

**After**:
```go
func (m Model) renderTaskCard(task *models.Task, maxWidth int, isSubtask bool) string {
    priorityStr := priorityStyle.Render(fmt.Sprintf("P%d", task.Priority))

    // NEW: Add dash prefix for subtasks
    prefix := ""
    if isSubtask {
        prefix = "- "
    }

    title := task.Title
    // ... truncation logic (adjust maxWidth for prefix) ...

    return fmt.Sprintf("%s%s %s %s", prefix, priorityStr, title, createdStr)
}
```

**New Parameter**: `isSubtask bool` - indicates if task should be rendered with dash prefix

**Caller Changes**:
```go
// In renderColumn function:
for i, task := range tasks {
    isSubtask := task.ParentID != nil
    cardContent := m.renderTaskCard(task, columnWidth-4, isSubtask)  // pass isSubtask
    // ... render with selection style ...
}
```

**Visual Change**:
```
Before:                After:
P1 Parent task         P1 Parent task
P2 Subtask task        - P2 Subtask task
```

---

## 4. Detail View (`internal/tui/detail.go`)

### No API Changes Required

**Existing Behavior**:
- Detail view already shows subtasks via `m.detailSubtasks`
- Subtasks rendered in "Subtasks (N)" section
- No changes needed (already meets spec requirements)

**Enhancement** (optional, not required for MVP):
- Could use `hierarchy.FormatWithIndentation()` for consistent formatting
- Current manual formatting is acceptable

---

## 5. Form (`internal/tui/form.go`)

### Modified Function: initCreateForm

**Purpose**: Support optional pre-population of parent ID field

**Before**:
```go
func (m *Model) initCreateForm()
```

**After**:
```go
func (m *Model) initCreateForm(parentID *string)
```

**Implementation**:
```go
func (m *Model) initCreateForm(parentID *string) {
    inputWidth := m.width - 20
    m.formInputs = make([]textinput.Model, 6)

    // ... initialize fields 0-4 (unchanged) ...

    // Parent task ID (field 5)
    m.formInputs[5] = textinput.New()
    m.formInputs[5].Placeholder = "Parent task ID (optional)"
    m.formInputs[5].Width = inputWidth

    // NEW: Pre-fill if parentID provided
    if parentID != nil {
        m.formInputs[5].SetValue(*parentID)
    }

    m.formFocusIndex = 0
    m.formTask = nil
}
```

**Callers**:
1. From kanban 'n' key: `m.initCreateForm(nil)` - no parent
2. From detail 'n' key: `m.initCreateForm(&m.detailTask.ID)` - pre-fill parent
3. From edit: No changes (edit uses `initEditForm`, not affected)

**Backward Compatibility**: All existing callers must be updated to pass `nil` explicitly

---

## 6. Modified Function: initEditForm

**Purpose**: Support optional pre-population of parent ID field (same pattern as initCreateForm)

**Signature Change**:
```go
// Before:
func (m *Model) initEditForm(task *models.Task)

// After: (NO CHANGE - already receives task with ParentID)
func (m *Model) initEditForm(task *models.Task)
```

**No changes needed**: Edit form already populates ParentID from `task.ParentID`

---

## Testing Contract

### Unit Tests

**Location**: `internal/tui/model_test.go`

1. `TestGetTasksByColumn_HierarchicalOrder` - verify parents and subtasks grouped
2. `TestGetTasksByColumn_PreservesSortMode` - verify sort mode respected within hierarchy

**Location**: `internal/tui/kanban_test.go`

3. `TestRenderTaskCard_SubtaskPrefix` - verify dash appears for subtasks
4. `TestRenderTaskCard_NoPrefix` - verify no dash for parents

### Integration Tests

**Location**: `tests/integration/tui_hierarchy_test.go`

1. Test creating subtask from detail view (full flow)
2. Test navigating through hierarchical kanban (j/k keys)
3. Test parent ID pre-populated when creating from detail

---

## Migration Notes

### Breaking Changes
- `initCreateForm()` now requires parameter (was zero parameters)
- `renderTaskCard()` now requires `isSubtask` parameter

### Non-Breaking Changes
- `GetTasksByColumn()` returns different order but same type
- Selection index math unchanged (flat index still works)

### Callers to Update
1. `internal/tui/app.go::handleKanbanKeys` - pass `nil` to `initCreateForm()`
2. `internal/tui/app.go::handleDetailKeys` - pass `&m.detailTask.ID` to `initCreateForm()`
3. `internal/tui/kanban.go::renderColumn` - pass `isSubtask` to `renderTaskCard()`

---

## Performance Impact

- `GetTasksByColumn()` now calls `BuildFlatHierarchy()`: adds O(n log n) overhead
- With 100 tasks, ~1ms additional latency (imperceptible)
- Rendering logic unchanged (still linear in task count)

**No performance concerns** for personal task manager use case.
