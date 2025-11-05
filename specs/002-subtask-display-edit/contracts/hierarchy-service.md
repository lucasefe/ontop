# Internal API Contract: Hierarchy Service

**File**: `internal/service/hierarchy.go` (NEW)
**Purpose**: Shared logic for building and displaying hierarchical task structures
**Consumers**: TUI kanban renderer, CLI list command, CLI show command

---

## Data Structures

### HierarchicalTask

```go
// HierarchicalTask wraps a Task with display metadata
type HierarchicalTask struct {
    Task        *models.Task  // The underlying task data
    IsSubtask   bool          // True if this task has a parent
    Indentation int           // Spaces for visual indentation (0 or 2)
    ParentTitle string        // Title of parent (empty if not subtask)
}
```

**Invariants**:
- If `IsSubtask == true`, then `Task.ParentID != nil` and `Indentation == 2`
- If `IsSubtask == false`, then `Indentation == 0`
- `ParentTitle` populated only if `IsSubtask == true`

---

## Functions

### BuildFlatHierarchy

```go
func BuildFlatHierarchy(tasks []*models.Task, sortMode tui.SortMode) []*HierarchicalTask
```

**Purpose**: Converts flat task list into hierarchical display order

**Parameters**:
- `tasks`: Flat list of tasks (may include parents and subtasks)
- `sortMode`: How to sort tasks (Priority, Description, CreatedAt, UpdatedAt)

**Returns**: Tasks in hierarchical display order (parents sorted, subtasks immediately follow their parent)

**Algorithm**:
1. Separate parents (ParentID == nil) and subtasks (ParentID != nil)
2. Sort parents by `sortMode`
3. For each parent:
   - Add parent to result with Indentation=0, IsSubtask=false
   - Find subtasks where ParentID == parent.ID
   - Sort subtasks by `sortMode`
   - Add each subtask with Indentation=2, IsSubtask=true, ParentTitle=parent.Title
4. Handle orphaned subtasks (ParentID points to non-existent task):
   - Treat as top-level tasks (Indentation=0, IsSubtask=false)

**Example**:
```go
input := []*models.Task{
    {ID: "P1", Title: "Parent 1", ParentID: nil, Priority: 2},
    {ID: "S1", Title: "Subtask 1", ParentID: &"P1", Priority: 1},
    {ID: "P2", Title: "Parent 2", ParentID: nil, Priority: 1},
}

output := BuildFlatHierarchy(input, tui.SortByPriority)
// Result:
// [
//   {Task: P2, IsSubtask: false, Indentation: 0},     // Priority 1
//   {Task: P1, IsSubtask: false, Indentation: 0},     // Priority 2
//   {Task: S1, IsSubtask: true, Indentation: 2, ParentTitle: "Parent 1"},
// ]
```

**Time Complexity**: O(n log n) - dominated by sorting

**Edge Cases**:
- Empty input → returns empty slice
- All parents, no subtasks → returns sorted parents
- All subtasks, no parents → treats all as orphans (top-level)
- Subtask with deleted parent → treated as orphan

---

### GetSubtasksForParent

```go
func GetSubtasksForParent(tasks []*models.Task, parentID string) []*models.Task
```

**Purpose**: Returns all subtasks for a specific parent task

**Parameters**:
- `tasks`: Full list of tasks to search
- `parentID`: ID of parent task

**Returns**: Subtasks where ParentID == parentID (unsorted)

**Usage**: Detail view uses this to show "Subtasks (N)" section

**Example**:
```go
tasks := []*models.Task{
    {ID: "P1", ParentID: nil},
    {ID: "S1", ParentID: &"P1"},
    {ID: "S2", ParentID: &"P1"},
    {ID: "P2", ParentID: nil},
}

subtasks := GetSubtasksForParent(tasks, "P1")
// Returns: [S1, S2]
```

**Time Complexity**: O(n) - linear scan

---

### FormatWithIndentation

```go
func FormatWithIndentation(ht *HierarchicalTask) string
```

**Purpose**: Formats task title with appropriate indentation for CLI output

**Parameters**:
- `ht`: HierarchicalTask with indentation metadata

**Returns**: Formatted string with indentation and dash prefix for subtasks

**Format**:
- Parent: `"P1 Title of parent task"`
- Subtask: `"  - P2 Title of subtask"` (2 spaces + dash + space + priority + title)

**Example**:
```go
parent := &HierarchicalTask{
    Task: &models.Task{ID: "P1", Title: "Parent", Priority: 1},
    IsSubtask: false,
    Indentation: 0,
}
fmt.Println(FormatWithIndentation(parent))
// Output: "P1 Parent"

subtask := &HierarchicalTask{
    Task: &models.Task{ID: "S1", Title: "Subtask", Priority: 2, ParentID: &"P1"},
    IsSubtask: true,
    Indentation: 2,
    ParentTitle: "Parent",
}
fmt.Println(FormatWithIndentation(subtask))
// Output: "  - P2 Subtask"
```

---

### IsOrphanedSubtask

```go
func IsOrphanedSubtask(task *models.Task, allTasks []*models.Task) bool
```

**Purpose**: Checks if a subtask's parent no longer exists

**Parameters**:
- `task`: Task to check
- `allTasks`: Full task list to search for parent

**Returns**: `true` if task has ParentID but parent doesn't exist, `false` otherwise

**Usage**: Called by `BuildFlatHierarchy()` to identify orphans

**Example**:
```go
subtask := &models.Task{ID: "S1", ParentID: &"DELETED"}
allTasks := []*models.Task{
    {ID: "P1", ParentID: nil},
}

isOrphan := IsOrphanedSubtask(subtask, allTasks)
// Returns: true (parent "DELETED" not in allTasks)
```

**Time Complexity**: O(n) - linear scan for parent

---

## Testing Contract

### Unit Tests Required

**Location**: `internal/service/hierarchy_test.go`

**Test Cases**:
1. `TestBuildFlatHierarchy_EmptyInput` - handles empty slice
2. `TestBuildFlatHierarchy_OnlyParents` - no subtasks
3. `TestBuildFlatHierarchy_OnlySubtasks` - all orphans
4. `TestBuildFlatHierarchy_MixedParentsAndSubtasks` - standard case
5. `TestBuildFlatHierarchy_OrphanedSubtask` - parent doesn't exist
6. `TestBuildFlatHierarchy_SortByPriority` - parents and subtasks sorted by priority
7. `TestBuildFlatHierarchy_SortByDate` - parents and subtasks sorted by date
8. `TestGetSubtasksForParent_NoSubtasks` - returns empty
9. `TestGetSubtasksForParent_MultipleSubtasks` - returns all
10. `TestFormatWithIndentation_Parent` - no indentation
11. `TestFormatWithIndentation_Subtask` - includes dash and indent
12. `TestIsOrphanedSubtask_ValidParent` - returns false
13. `TestIsOrphanedSubtask_MissingParent` - returns true

### Integration Tests

**Location**: `tests/integration/hierarchy_test.go`

**Test Cases**:
1. End-to-end: Create parent, add subtasks, verify hierarchy in CLI list
2. End-to-end: Move subtask to different column, verify still appears under parent in detail view
3. End-to-end: Delete parent attempt blocked when subtasks exist

---

## Performance Guarantees

- `BuildFlatHierarchy()`: O(n log n) time, O(n) space
- `GetSubtasksForParent()`: O(n) time, O(k) space where k = number of subtasks
- `IsOrphanedSubtask()`: O(n) time, O(1) space
- All functions handle up to 10,000 tasks efficiently (<100ms)

---

## Error Handling

All functions are panic-free and handle edge cases gracefully:
- `nil` slices → return empty results
- Invalid task data → skip and continue processing
- Missing parent → treat subtask as orphan

**No errors returned**: These are display/formatting functions, not critical operations
