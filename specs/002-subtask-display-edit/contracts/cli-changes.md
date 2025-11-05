# Internal API Contract: CLI Changes

**Purpose**: Define changes to CLI commands for hierarchical display
**Affected Files**: `internal/cli/list.go`, `internal/cli/show.go`

---

## 1. List Command (`internal/cli/list.go`)

### Modified Function: ListCommand

**Purpose**: Display tasks with hierarchical indentation in CLI output

**Before**:
```go
func ListCommand(db *sql.DB, args []string) {
    // ... parse flags ...
    tasks, err := storage.ListTasks(db, filters)
    // ... error handling ...

    for _, task := range tasks {
        fmt.Printf("[%s] P%d | %s | %s\n",
            task.ID, task.Priority, task.Column, task.Title)
    }
}
```

**After**:
```go
func ListCommand(db *sql.DB, args []string) {
    // ... parse flags (unchanged) ...
    tasks, err := storage.ListTasks(db, filters)
    // ... error handling (unchanged) ...

    // NEW: Build hierarchical display order
    hierarchical := hierarchy.BuildFlatHierarchy(tasks, getSortModeFromFlags(fs))

    for _, ht := range hierarchical {
        // NEW: Add indentation for subtasks
        indent := strings.Repeat(" ", ht.Indentation)
        prefix := ""
        if ht.IsSubtask {
            prefix = "- "
        }

        fmt.Printf("%s%s[%s] P%d | %s | %s\n",
            indent, prefix,
            ht.Task.ID, ht.Task.Priority, ht.Task.Column, ht.Task.Title)
    }
}
```

**Visual Change**:
```
Before:
[P1] P2 | inbox | Parent task
[S1] P1 | inbox | Subtask 1
[S2] P3 | inbox | Subtask 2

After:
[P1] P2 | inbox | Parent task
  - [S1] P1 | inbox | Subtask 1
  - [S2] P3 | inbox | Subtask 2
```

**Key Changes**:
1. Import `internal/service/hierarchy` package
2. Call `BuildFlatHierarchy()` to get hierarchical order
3. Add 2-space indentation for subtasks
4. Add dash prefix for subtasks

**Backward Compatibility**:
- Output format slightly different (indentation added)
- Parsing scripts may need update if they rely on fixed column positions
- **Mitigation**: Add `--format=json` flag for machine-readable output (future enhancement)

---

## 2. Show Command (`internal/cli/show.go`)

### Modified Function: ShowCommand

**Purpose**: Display subtasks when showing a parent task

**Before**:
```go
func ShowCommand(db *sql.DB, args []string) {
    // ... parse args ...
    task, err := storage.GetTask(db, taskID)
    // ... error handling ...

    fmt.Printf("ID: %s\n", task.ID)
    fmt.Printf("Title: %s\n", task.Title)
    // ... other fields ...
}
```

**After**:
```go
func ShowCommand(db *sql.DB, args []string) {
    // ... parse args (unchanged) ...
    task, err := storage.GetTask(db, taskID)
    // ... error handling (unchanged) ...

    // Display task details (unchanged)
    fmt.Printf("ID: %s\n", task.ID)
    fmt.Printf("Title: %s\n", task.Title)
    // ... other fields (unchanged) ...

    // NEW: Display subtasks if any
    allTasks, _ := storage.ListTasks(db, map[string]interface{}{})
    subtasks := hierarchy.GetSubtasksForParent(allTasks, task.ID)

    if len(subtasks) > 0 {
        fmt.Printf("\nSubtasks (%d):\n", len(subtasks))
        for _, st := range subtasks {
            fmt.Printf("  - [%s] P%d: %s", st.ID, st.Priority, st.Title)
            if st.Progress > 0 {
                fmt.Printf(" (%d%%)", st.Progress)
            }
            fmt.Printf(" [%s]\n", formatColumnShort(st.Column))
        }
    }
}
```

**Output Example**:
```
ID: P1
Title: Implement authentication
Description: Add OAuth2 and JWT support
Priority: 2
Column: in_progress
Progress: 40%
Tags: backend, security
Created: 2025-01-04 10:30:00
Updated: 2025-01-04 14:15:00

Subtasks (3):
  - [S1] P1: Add login form (50%) [in_progress]
  - [S2] P2: JWT validation (0%) [inbox]
  - [S3] P1: Write tests (100%) [done]
```

**Key Changes**:
1. After displaying main task, query for subtasks
2. Display subtasks section if any exist
3. Include progress and column for each subtask

---

## 3. Helper Functions

### New Function: getSortModeFromFlags

**Location**: `internal/cli/list.go`

**Purpose**: Convert CLI flag to TUI SortMode enum

**Signature**:
```go
func getSortModeFromFlags(fs *flag.FlagSet) tui.SortMode
```

**Implementation**:
```go
func getSortModeFromFlags(fs *flag.FlagSet) tui.SortMode {
    sortFlag := fs.Lookup("sort")
    if sortFlag == nil {
        return tui.SortByPriority // default
    }

    switch sortFlag.Value.String() {
    case "priority":
        return tui.SortByPriority
    case "title":
        return tui.SortByDescription
    case "created":
        return tui.SortByCreated
    case "updated":
        return tui.SortByUpdated
    default:
        return tui.SortByPriority
    }
}
```

**Usage**: Pass result to `BuildFlatHierarchy()` for consistent sorting

---

### New Function: formatColumnShort

**Location**: `internal/cli/show.go`

**Purpose**: Format column name for compact display

**Signature**:
```go
func formatColumnShort(column string) string
```

**Implementation**:
```go
func formatColumnShort(column string) string {
    switch column {
    case models.ColumnInbox:
        return "Inbox"
    case models.ColumnInProgress:
        return "Progress"
    case models.ColumnDone:
        return "Done"
    default:
        return column
    }
}
```

**Usage**: Display subtask column in `show` command output

---

## 4. Flag Additions

### List Command

**New Flag** (optional enhancement, not required for MVP):
```go
sortFlag := fs.String("sort", "priority", "Sort order: priority, title, created, updated")
```

**Purpose**: Allow user to control sort mode in CLI (matches TUI 's' key cycling)

**Example**:
```bash
ontop list --sort=created  # Show tasks sorted by creation date
```

**Current Status**: Not required for MVP (default to priority sort)

---

## Testing Contract

### Unit Tests

**Location**: `internal/cli/list_test.go`

1. `TestListCommand_HierarchicalDisplay` - verify indentation in output
2. `TestListCommand_SubtaskPrefixDash` - verify dash appears
3. `TestListCommand_SortMode` - verify sort flag respected

**Location**: `internal/cli/show_test.go`

4. `TestShowCommand_WithSubtasks` - verify subtasks section appears
5. `TestShowCommand_NoSubtasks` - verify no empty section

### Integration Tests

**Location**: `tests/integration/cli_hierarchy_test.go`

1. Create parent + 2 subtasks via CLI → run list → verify hierarchy in output
2. Show parent task → verify subtasks listed
3. List with different sort modes → verify hierarchy maintained

---

## Output Format Specification

### List Command Output Format

```
Standard task:
[TASK_ID] P{priority} | {column} | {title}

Subtask:
  - [TASK_ID] P{priority} | {column} | {title}
```

**Example**:
```
[01K98X44S6T] P1 | inbox | Implement feature X
  - [01K98X4521M] P2 | inbox | Write unit tests
  - [01K98X4688N] P1 | in_progress | Update documentation
[01K98X49TWG] P3 | done | Fix bug #123
```

**Invariants**:
- Subtasks always have 2-space indentation
- Subtasks always have dash prefix after indentation
- Column width fixed for alignment (use padding if needed)

### Show Command Output Format

```
{field}: {value}
{field}: {value}
...

Subtasks ({count}):
  - [{ID}] P{priority}: {title} ({progress}%) [{column}]
  - [{ID}] P{priority}: {title} ({progress}%) [{column}]
```

**Example**: See "Output Example" in Show Command section above

---

## Migration Notes

### Breaking Changes
- List output format changed (indentation added for subtasks)
- Scripts parsing list output may break if they assume fixed column positions

### Mitigation Strategies
1. Document output format change in release notes
2. Consider adding `--format=plain` flag to preserve old behavior (future)
3. Recommend using `--format=json` for scripting (future enhancement)

### Non-Breaking Changes
- Show command adds new optional section (doesn't affect existing field parsing)
- All existing flags preserved and functional

---

## Performance Impact

- `list` command: Adds `BuildFlatHierarchy()` call - O(n log n) overhead
- `show` command: Adds `GetSubtasksForParent()` call - O(n) overhead
- With 1000 tasks, both add <10ms latency (CLI is interactive, acceptable)

**No performance concerns** for personal task manager use case.

---

## Future Enhancements (Not in Scope)

1. **JSON output format**: `ontop list --format=json` for machine-readable output
2. **Depth limiting**: `ontop list --max-depth=1` to collapse deep hierarchies
3. **Tree drawing**: Use Unicode box characters (├──, └──) instead of dashes
   ```
   Parent task
   ├─ Subtask 1
   └─ Subtask 2
   ```

These are explicitly out of scope for MVP. Focus on simple 2-space indentation with dash prefix.
