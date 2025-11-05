# Quickstart: Implementing Hierarchical Subtask Display

**Feature**: 002-subtask-display-edit
**For**: Developers implementing this feature
**Prerequisites**: Familiarity with Go, Bubble Tea framework, and OnTop codebase

---

## Architecture Overview

### Data Flow

```
Storage (Flat)  →  Hierarchy Service  →  Display (Hierarchical)
                                      ↓
                                   TUI/CLI
```

1. **Storage Layer** (`internal/storage/sqlite.go`): Tasks stored flat with ParentID field
2. **Service Layer** (`internal/service/hierarchy.go`): Builds hierarchical display structure
3. **Presentation Layer** (`internal/tui/`, `internal/cli/`): Renders with indentation

**Key Insight**: Database stays flat, hierarchy is computed on-the-fly during rendering.

---

## Implementation Order

Follow this order to minimize integration issues:

### Phase 1: Foundation (Day 1)
1. ✅ Create `internal/service/hierarchy.go` with all functions
2. ✅ Write unit tests for hierarchy service
3. ✅ Ensure all tests pass before proceeding

### Phase 2: TUI Integration (Day 2)
4. ✅ Modify `GetTasksByColumn()` to use `BuildFlatHierarchy()`
5. ✅ Update `renderTaskCard()` to accept `isSubtask` parameter
6. ✅ Update `handleDetailKeys()` to handle 'n' key
7. ✅ Update `initCreateForm()` to accept optional `parentID`
8. ✅ Update all callers of `initCreateForm()` (pass `nil` or `&id`)
9. ✅ Test TUI manually: create parent, create subtask, verify hierarchy

### Phase 3: CLI Integration (Day 3)
10. ✅ Modify `ListCommand()` to use `BuildFlatHierarchy()`
11. ✅ Modify `ShowCommand()` to display subtasks
12. ✅ Test CLI manually: verify indentation matches TUI

### Phase 4: Polish & Testing (Day 4)
13. ✅ Write integration tests (TUI and CLI)
14. ✅ Handle edge cases (orphaned subtasks, deletion protection)
15. ✅ Update documentation (README, help text)

---

## Key Functions to Understand

### 1. BuildFlatHierarchy (Most Complex)

**Location**: `internal/service/hierarchy.go`

**What it does**: Converts `[]*models.Task` (flat) → `[]*HierarchicalTask` (display order)

**Algorithm Walkthrough**:
```go
func BuildFlatHierarchy(tasks []*models.Task, sortMode tui.SortMode) []*HierarchicalTask {
    // Step 1: Separate parents from subtasks
    parents := []*models.Task{}
    subtasks := []*models.Task{}
    for _, task := range tasks {
        if task.ParentID == nil {
            parents = append(parents, task)
        } else {
            subtasks = append(subtasks, task)
        }
    }

    // Step 2: Sort parents by sortMode
    sortTasks(parents, sortMode)  // helper function

    // Step 3: Build result by interleaving parents and their subtasks
    result := []*HierarchicalTask{}
    for _, parent := range parents {
        // Add parent
        result = append(result, &HierarchicalTask{
            Task: parent,
            IsSubtask: false,
            Indentation: 0,
        })

        // Find and add subtasks for this parent
        parentSubtasks := findSubtasksForParent(subtasks, parent.ID)
        sortTasks(parentSubtasks, sortMode)  // sort subtasks too
        for _, st := range parentSubtasks {
            result = append(result, &HierarchicalTask{
                Task: st,
                IsSubtask: true,
                Indentation: 2,
                ParentTitle: parent.Title,
            })
        }
    }

    // Step 4: Handle orphaned subtasks (parent not in list)
    for _, st := range subtasks {
        if !hasParentInList(st, parents) {
            result = append(result, &HierarchicalTask{
                Task: st,
                IsSubtask: false,  // treat as top-level
                Indentation: 0,
            })
        }
    }

    return result
}
```

**Testing Strategy**:
- Test with only parents → should return sorted parents
- Test with only subtasks → should treat all as orphans
- Test with mixed → should group correctly
- Test with orphans → should handle gracefully

---

### 2. GetTasksByColumn (Critical Integration Point)

**Location**: `internal/tui/model.go`

**What it does**: Filter tasks by column, then apply hierarchical ordering

**Before**:
```go
func (m *Model) GetTasksByColumn(column string) []*models.Task {
    filtered := filterByColumn(m.tasks, column)
    return sortTasks(filtered, m.sortMode)
}
```

**After**:
```go
func (m *Model) GetTasksByColumn(column string) []*models.Task {
    filtered := filterByColumn(m.tasks, column)
    hierarchical := hierarchy.BuildFlatHierarchy(filtered, m.sortMode)

    // Extract Task pointers (unwrap HierarchicalTask)
    result := make([]*models.Task, len(hierarchical))
    for i, ht := range hierarchical {
        result[i] = ht.Task
    }
    return result
}
```

**Why this works**:
- Callers still get `[]*models.Task` slice
- Selection index still maps directly to slice position
- Navigation code (j/k keys) requires zero changes

**Testing**:
- Create parent in inbox, subtask in inbox → should appear together
- Create parent in inbox, subtask in progress → should appear separately
- Move task between columns → hierarchy should rebuild correctly

---

### 3. initCreateForm (Breaking Change, Handle Carefully)

**Location**: `internal/tui/form.go`

**Old Signature**:
```go
func (m *Model) initCreateForm()
```

**New Signature**:
```go
func (m *Model) initCreateForm(parentID *string)
```

**All Callers Must Be Updated**:
```go
// From kanban 'n' key (app.go:166):
m.initCreateForm(nil)  // ADD nil

// From detail 'n' key (app.go:NEW):
m.initCreateForm(&m.detailTask.ID)  // NEW: pass parent ID
```

**How to Find All Callers**:
```bash
grep -r "initCreateForm()" internal/tui/
```

**Common Mistake**: Forgetting to update a caller → compile error

**Testing**:
- Press 'n' in kanban → parent ID field should be empty
- Press 'n' in detail → parent ID field should be pre-filled
- Edit parent ID field → should still work (not disabled)

---

## Common Gotchas

### Gotcha 1: Selection Index Off-By-One

**Problem**: After adding hierarchy, navigation feels "jumpy"

**Cause**: Selection index points to wrong task after hierarchy rebuild

**Solution**: Always reset `selectedTask = 0` after task list mutation (create, move, delete)

**Example**:
```go
// In saveForm() after creating task:
m.viewMode = ViewModeKanban
m.selectedTask = 0  // IMPORTANT: reset to top
return m, m.loadTasks
```

---

### Gotcha 2: Subtasks in Different Columns

**Problem**: Subtask doesn't appear under parent in kanban

**Cause**: Parent in inbox, subtask in done → they're in different columns

**Expected Behavior**: Subtask appears in done column (its actual column), not under parent

**Not a Bug**: This is correct per spec. Detail view shows all subtasks regardless of column.

**Testing**:
```bash
# Create parent in inbox
ontop add --title "Parent" --column inbox

# Create subtask in done
ontop add --title "Subtask" --column done --parent <PARENT_ID>

# Verify: subtask appears in DONE column, not INBOX
# Verify: pressing Enter on parent shows subtask in detail view
```

---

### Gotcha 3: Orphaned Subtasks

**Problem**: Deleted parent, now subtask has invalid ParentID

**Current Behavior**: Deletion blocked if subtasks exist

**Implementation**:
```go
// In storage.DeleteTask():
func DeleteTask(db *sql.DB, id string) error {
    // Check for subtasks BEFORE soft-delete
    var count int
    err := db.QueryRow("SELECT COUNT(*) FROM tasks WHERE parent_id = ? AND deleted_at IS NULL", id).Scan(&count)
    if err != nil {
        return err
    }

    if count > 0 {
        return fmt.Errorf("cannot delete task with %d subtasks (delete subtasks first)", count)
    }

    // Proceed with soft-delete
    // ...
}
```

**User Experience**: User must delete/move subtasks before deleting parent

---

### Gotcha 4: Multi-Level Nesting

**Problem**: User tries to create subtask of a subtask

**Spec Requirement**: Maximum 1 level of nesting (parent → subtask only)

**Implementation**:
```go
// In form.go, saveForm():
if parentIDStr != "" {
    parent, err := storage.GetTask(m.db, parentIDStr)
    if err != nil {
        m.formErr = fmt.Errorf("parent task not found: %s", parentIDStr)
        return m, nil
    }

    // NEW: Check if parent is itself a subtask
    if parent.ParentID != nil {
        m.formErr = fmt.Errorf("cannot create subtask of subtask (max 1 level nesting)")
        return m, nil
    }

    parentID = &parentIDStr
}
```

**Error Message**: "Cannot create subtask of subtask (max 1 level nesting)"

---

## Testing Approach

### Unit Tests

**File**: `internal/service/hierarchy_test.go`

```go
func TestBuildFlatHierarchy_BasicCase(t *testing.T) {
    parent := &models.Task{ID: "P1", Title: "Parent", ParentID: nil}
    sub1 := &models.Task{ID: "S1", Title: "Sub1", ParentID: &"P1"}
    sub2 := &models.Task{ID: "S2", Title: "Sub2", ParentID: &"P1"}

    input := []*models.Task{sub1, parent, sub2}  // intentionally unordered
    result := BuildFlatHierarchy(input, tui.SortByPriority)

    assert.Len(t, result, 3)
    assert.Equal(t, "P1", result[0].Task.ID)  // parent first
    assert.Equal(t, "S1", result[1].Task.ID)  // subtask second
    assert.Equal(t, "S2", result[2].Task.ID)  // subtask third
    assert.True(t, result[1].IsSubtask)
    assert.Equal(t, 2, result[1].Indentation)
}
```

**Run Tests**:
```bash
go test ./internal/service/... -v
```

---

### Integration Tests

**File**: `tests/integration/hierarchy_test.go`

```go
func TestTUIHierarchyDisplay(t *testing.T) {
    // Setup: Create test database
    db := setupTestDB(t)
    defer cleanupTestDB(t, db)

    // Create parent task
    parent := createTask(t, db, &models.Task{
        Title: "Parent",
        Column: models.ColumnInbox,
    })

    // Create 2 subtasks
    sub1 := createTask(t, db, &models.Task{
        Title: "Subtask 1",
        ParentID: &parent.ID,
        Column: models.ColumnInbox,
    })

    // Launch TUI model
    model := tui.NewModel(db)
    tasks := model.GetTasksByColumn(models.ColumnInbox)

    // Verify hierarchy
    assert.Len(t, tasks, 2)  // parent + 1 subtask
    assert.Equal(t, parent.ID, tasks[0].ID)
    assert.Equal(t, sub1.ID, tasks[1].ID)

    // Verify subtask follows parent immediately
    assert.Equal(t, &parent.ID, tasks[1].ParentID)
}
```

**Run Integration Tests**:
```bash
go test ./tests/integration/... -v
```

---

### Manual Testing Checklist

**TUI Mode** (`./ontop` with no args):

1. ✅ Create parent task ('n' key in kanban)
2. ✅ Press Enter on parent → detail view opens
3. ✅ Press 'n' in detail view → form opens with parent ID pre-filled
4. ✅ Create subtask → verify it appears indented below parent
5. ✅ Navigate with j/k → cursor moves through parent and subtask
6. ✅ Move subtask to different column → verify it disappears from parent's column
7. ✅ View parent details → verify subtask still listed in "Subtasks" section
8. ✅ Try to delete parent → verify blocked with error message

**CLI Mode**:

1. ✅ `ontop list` → verify subtasks indented with dash prefix
2. ✅ `ontop show <parent-id>` → verify subtasks listed at bottom
3. ✅ `ontop add --title "Sub" --parent <id>` → verify subtask created
4. ✅ Output piping: `ontop list | grep "P1"` → verify scriptable

---

## Debugging Tips

### Problem: Subtasks Not Appearing

**Check**:
1. Run `ontop list --json` (if implemented) to verify ParentID is set
2. Check database: `sqlite3 ~/.ontop/tasks.db "SELECT id, title, parent_id FROM tasks"`
3. Add debug logging in `BuildFlatHierarchy()`:
   ```go
   fmt.Printf("DEBUG: Found %d parents, %d subtasks\n", len(parents), len(subtasks))
   ```

### Problem: Selection Jumps to Wrong Task

**Check**:
1. Verify `selectedTask` is reset to 0 after task mutations
2. Print task list in `renderColumn()`:
   ```go
   for i, task := range tasks {
       fmt.Printf("DEBUG: Index %d = %s (IsSubtask: %v)\n", i, task.ID, task.ParentID != nil)
   }
   ```

### Problem: Parent ID Not Pre-filling

**Check**:
1. Verify `handleDetailKeys()` passes `&m.detailTask.ID` (not just `m.detailTask.ID`)
2. Add logging in `initCreateForm()`:
   ```go
   if parentID != nil {
       fmt.Printf("DEBUG: Pre-filling parent ID: %s\n", *parentID)
   }
   ```

---

## Performance Considerations

### Expected Performance

- 100 tasks: <1ms for `BuildFlatHierarchy()`
- 1000 tasks: <10ms for `BuildFlatHierarchy()`
- 10,000 tasks: <100ms (unlikely for personal use)

### Optimization (Not Needed for MVP)

If performance becomes an issue:
1. Cache hierarchy structure in model (invalidate on mutation)
2. Use parent→children map instead of linear scan
3. Implement lazy subtask loading (fetch only when parent expanded)

**For Personal Task Manager**: None of these optimizations needed.

---

## Documentation Updates

After implementation, update:

1. **README.md**: Add subtask feature to feature list
2. **CLAUDE.md**: Update AI context with hierarchy patterns
3. **CLI help text**: Document `--parent` flag for `add` command
4. **TUI help**: Mention 'n' key in detail view creates subtask

---

## Summary

**Simplest Path to MVP**:
1. Start with `hierarchy.go` (no dependencies)
2. Integrate with TUI `GetTasksByColumn()` (single point of integration)
3. Update rendering (visual change only)
4. Add 'n' key handler (user interaction)
5. Integrate with CLI (parallel to TUI)

**Timeline**: 3-4 days for experienced Go developer familiar with codebase

**Risk Areas**: Selection index math, form parameter changes (update all callers)

**Success Criteria**: Can create parent, create subtask from detail view, see hierarchy in kanban and CLI list
