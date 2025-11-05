# Data Model: Hierarchical Subtask Display

**Feature**: 002-subtask-display-edit
**Date**: 2025-01-04

## Overview

This feature requires NO changes to the existing data model. The Task struct already supports parent-child relationships via the `ParentID` field. This document describes the existing model and new display-only helper structures.

---

## Existing Data Model (No Changes)

### Task Entity

**Location**: `internal/models/task.go`

```go
type Task struct {
    ID          string     `json:"id"`
    Title       string     `json:"title"`
    Description string     `json:"description"`
    Priority    int        `json:"priority"`        // 1-5 (1=highest)
    Column      string     `json:"column"`          // inbox, in_progress, done
    Progress    int        `json:"progress"`        // 0-100
    ParentID    *string    `json:"parent_id"`       // 🔑 NULL = parent task, non-NULL = subtask
    Archived    bool       `json:"archived"`
    Tags        []string   `json:"tags"`
    CreatedAt   time.Time  `json:"created_at"`
    UpdatedAt   time.Time  `json:"updated_at"`
    CompletedAt *time.Time `json:"completed_at"`
    DeletedAt   *time.Time `json:"deleted_at"`
}
```

**Key Field**: `ParentID *string`
- `nil`: This is a parent task (or standalone task with no subtasks)
- `"01K98X44S6T..."`: This is a subtask of the task with given ID

**Constraints**:
- ParentID foreign key references tasks(id) (enforced in SQLite schema)
- Parent tasks can have multiple subtasks
- **New constraint** (this feature): Subtasks cannot have subtasks (max 1 level nesting)

### Column Constants

**Location**: `internal/models/column.go`

```go
const (
    ColumnInbox      = "inbox"
    ColumnInProgress = "in_progress"
    ColumnDone       = "done"
)
```

---

## New Display Structures (Display Layer Only)

### HierarchicalTask (Helper Struct)

**Location**: `internal/service/hierarchy.go` (NEW FILE)

**Purpose**: Wraps Task with display metadata for rendering hierarchical views.

```go
// HierarchicalTask wraps a Task with display context for hierarchical rendering
type HierarchicalTask struct {
    Task        *models.Task  // The actual task data
    IsSubtask   bool          // True if ParentID != nil
    Indentation int           // Number of spaces for indentation (0 or 2)
    ParentTitle string        // Title of parent task (empty if not subtask)
}
```

**Usage**: Returned by `BuildFlatHierarchy()` for rendering in TUI kanban and CLI list.

**Example**:
```go
[
    {Task: parentTask1, IsSubtask: false, Indentation: 0},
    {Task: subtask1, IsSubtask: true, Indentation: 2, ParentTitle: "Parent Task 1"},
    {Task: subtask2, IsSubtask: true, Indentation: 2, ParentTitle: "Parent Task 1"},
    {Task: parentTask2, IsSubtask: false, Indentation: 0},
]
```

---

## Hierarchy Helper Functions

### BuildFlatHierarchy

**Signature**:
```go
func BuildFlatHierarchy(tasks []*models.Task, sortMode SortMode) []*HierarchicalTask
```

**Purpose**: Converts flat task list into hierarchical display order.

**Algorithm**:
1. Separate tasks into parents (ParentID == nil) and subtasks (ParentID != nil)
2. Sort parents according to `sortMode` (priority, date, etc.)
3. For each parent:
   a. Add parent to result list
   b. Find all subtasks where ParentID == parent.ID
   c. Sort subtasks by same `sortMode`
   d. Add subtasks immediately after parent
4. Return flattened list with HierarchicalTask wrappers

**Complexity**: O(n log n) for sorting + O(n) for grouping = O(n log n)

**Example Input**:
```go
tasks = [
    {ID: "P1", ParentID: nil, Priority: 2, Column: "inbox"},
    {ID: "S1", ParentID: &"P1", Priority: 1, Column: "inbox"},
    {ID: "P2", ParentID: nil, Priority: 1, Column: "inbox"},
    {ID: "S2", ParentID: &"P1", Priority: 3, Column: "inbox"},
]
sortMode = SortByPriority
```

**Example Output**:
```go
[
    {Task: P2, IsSubtask: false, Indentation: 0},              // Priority 1 parent
    {Task: P1, IsSubtask: false, Indentation: 0},              // Priority 2 parent
    {Task: S1, IsSubtask: true, Indentation: 2, ParentTitle: "P1"}, // P1's subtasks, sorted by priority
    {Task: S2, IsSubtask: true, Indentation: 2, ParentTitle: "P1"},
]
```

### GetSubtasksForParent

**Signature**:
```go
func GetSubtasksForParent(tasks []*models.Task, parentID string) []*models.Task
```

**Purpose**: Returns all subtasks for a given parent task (used in detail view).

**Algorithm**:
1. Filter tasks where ParentID == parentID
2. Return filtered list (no sorting—detail view handles sort)

**Usage**: Detail view calls this to display "Subtasks (3)" section.

### IsOrphanedSubtask

**Signature**:
```go
func IsOrphanedSubtask(task *models.Task, allTasks []*models.Task) bool
```

**Purpose**: Checks if a subtask's parent no longer exists (orphaned).

**Algorithm**:
1. If task.ParentID == nil, return false (not a subtask)
2. Search allTasks for task with ID == *task.ParentID
3. Return true if not found (orphaned), false otherwise

**Usage**: Used during hierarchy building to handle edge cases (orphaned subtasks treated as top-level tasks).

---

## State Transitions (No Changes)

Task state transitions remain unchanged:
- **Create**: Task can be created with ParentID set (subtask) or nil (parent/standalone)
- **Move**: Task moves between columns (inbox → in_progress → done) independently of parent
- **Archive**: Task can be archived independently of parent/subtasks
- **Delete**: Task deletion blocked if subtasks exist (prevents orphans)

**New validation** (this feature):
- Before deleting task with subtasks, return error: "Cannot delete task with subtasks. Delete or move subtasks first."

---

## Database Schema (No Changes)

Existing schema in `internal/storage/migrations.go`:

```sql
CREATE TABLE IF NOT EXISTS tasks (
    id TEXT PRIMARY KEY,
    title TEXT NOT NULL DEFAULT '',
    description TEXT NOT NULL DEFAULT '',
    priority INTEGER NOT NULL CHECK (priority >= 1 AND priority <= 5),
    column TEXT NOT NULL CHECK (column IN ('inbox', 'in_progress', 'done')),
    progress INTEGER NOT NULL DEFAULT 0 CHECK (progress >= 0 AND progress <= 100),
    parent_id TEXT,  -- 🔑 NULL = parent, non-NULL = subtask
    archived INTEGER NOT NULL DEFAULT 0,
    tags TEXT NOT NULL DEFAULT '[]',
    created_at TEXT NOT NULL,
    updated_at TEXT NOT NULL,
    completed_at TEXT,
    deleted_at TEXT,
    FOREIGN KEY (parent_id) REFERENCES tasks(id)
);

CREATE INDEX IF NOT EXISTS idx_tasks_parent_id ON tasks(parent_id);
```

**Key points**:
- `parent_id` column already exists with foreign key constraint
- Index on `parent_id` ensures efficient subtask lookups
- No schema migration needed

---

## Validation Rules

### Existing Rules (Enforced)
1. Priority must be 1-5
2. Column must be one of: inbox, in_progress, done
3. Progress must be 0-100
4. ParentID must reference valid task ID (or be NULL)

### New Rules (This Feature)
1. **No multi-level nesting**: If task has ParentID, it cannot itself be a parent
   - Enforced in: `initCreateForm()` and CLI `add` command
   - Check: If creating subtask, validate that parent is not itself a subtask

2. **Deletion protection**: Cannot delete task if it has subtasks
   - Enforced in: `storage.DeleteTask()` (add check before soft-delete)
   - Check: Query for tasks where ParentID == task.ID, return error if found

3. **Orphan prevention**: Subtask ParentID must point to existing, non-deleted task
   - Enforced in: Existing foreign key constraint + app-level check
   - On parent deletion attempt: Block if subtasks exist

---

## Example Data Scenarios

### Scenario 1: Simple Parent-Subtask Hierarchy
```go
Parent: {ID: "P1", Title: "Implement auth", ParentID: nil, Column: "in_progress"}
Subtasks:
  - {ID: "S1", Title: "Add login form", ParentID: &"P1", Column: "in_progress"}
  - {ID: "S2", Title: "Add JWT validation", ParentID: &"P1", Column: "done"}
  - {ID: "S3", Title: "Write tests", ParentID: &"P1", Column: "inbox"}
```

**Display in Kanban**:
```
IN PROGRESS                  DONE
P1 Implement auth            S2 Add JWT validation
  - S1 Add login form

INBOX
S3 Write tests
```

*Note: S2 and S3 appear in their actual columns, not grouped under P1*

### Scenario 2: Multiple Parents in Same Column
```go
Parent1: {ID: "P1", Title: "Feature A", Priority: 1, Column: "inbox"}
Subtask1: {ID: "S1", Title: "Subtask 1", Priority: 2, ParentID: &"P1", Column: "inbox"}
Parent2: {ID: "P2", Title: "Feature B", Priority: 2, Column: "inbox"}
Subtask2: {ID: "S2", Title: "Subtask 2", Priority: 1, ParentID: &"P2", Column: "inbox"}
```

**Display Order** (Priority sort):
```
INBOX
P1 Feature A (Priority 1)
  - S1 Subtask 1 (Priority 2)
P2 Feature B (Priority 2)
  - S2 Subtask 2 (Priority 1)
```

*Note: P2's subtask S2 has higher priority (1) than P2 itself (2), but still appears under parent*

### Scenario 3: Orphaned Subtask (Edge Case)
```go
Subtask: {ID: "S1", Title: "Orphaned task", ParentID: &"DELETED_PARENT"}
```

**Handling**:
- `BuildFlatHierarchy()` detects orphan (parent not in list)
- Treats orphan as top-level task (IsSubtask: false, Indentation: 0)
- Optional: Display warning in CLI ("Warning: Task S1 has invalid parent")

---

## Summary

- ✅ **No database schema changes required**
- ✅ **No changes to Task struct**
- ✅ **New display-only helper structures** (HierarchicalTask)
- ✅ **New validation rules** (no multi-level nesting, deletion protection)
- ✅ **Existing data model fully supports feature** (ParentID field)

All data model requirements met. Ready for contract definitions and implementation.
