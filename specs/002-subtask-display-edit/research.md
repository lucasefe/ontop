# Research: Hierarchical Subtask Display Technical Decisions

**Feature**: 002-subtask-display-edit
**Date**: 2025-01-04
**Status**: Complete

## Overview

This document captures technical research and decisions for implementing hierarchical subtask display in the OnTop task manager. The feature must work consistently across both TUI (Bubble Tea) and CLI modes while maintaining simplicity and vim-style interaction patterns.

---

## Decision 1: Hierarchical List Structure

### Problem
How should we represent parent tasks with their subtasks in memory for rendering? The display needs to show subtasks indented below parents, but the underlying storage is flat (all tasks in a single slice).

### Options Evaluated

**Option A: Flatten on-demand during render**
- Keep existing `GetTasksByColumn()` returning `[]*models.Task`
- During rendering, check each task's ParentID and compute indentation
- Pros: Minimal changes to existing code, simple data flow
- Cons: Must scan full list multiple times to find parent-child relationships

**Option B: Pre-compute hierarchical structure**
- Create new `GetHierarchicalTasksByColumn()` returning custom struct
- Build parent → children map once, return nested structure
- Pros: Faster rendering (O(n) build, O(1) lookups), clearer intent
- Cons: New data structure, more complex selection index tracking

### Decision: **Option A - Flatten on-demand during render**

**Rationale**:
1. **Simplicity First** (Constitution principle): Existing code already works with flat `[]*models.Task` slices
2. **Performance adequate**: With ~1000 tasks expected (personal use), O(n²) parent lookups are sub-millisecond
3. **Easier selection tracking**: Flat index directly maps to rendered list position
4. **Less code churn**: No need to refactor `GetSelectedTask()`, navigation logic, or other consumers

**Implementation approach**:
- Add new `internal/service/hierarchy.go` with helper function: `BuildFlatHierarchy(tasks []*models.Task) []*models.Task`
- Returns tasks in display order: parent, then its subtasks (if in same column), then next parent
- TUI kanban calls this before rendering
- CLI list uses same function for consistent output

**Trade-off accepted**: Slightly less efficient than pre-computed hierarchy, but maintains code simplicity and reduces risk of regression bugs.

---

## Decision 2: Selection Index Strategy

### Problem
How do we track which task is selected when the displayed list includes both parent tasks and indented subtasks? Pressing j/k needs to move through the hierarchical list naturally.

### Options Evaluated

**Option A: Flat index (0, 1, 2... across all displayed items)**
- `selectedTask` remains a single integer
- Index corresponds to position in flattened hierarchical list
- Pros: Simple, minimal code changes, works with existing navigation
- Cons: Index becomes stale if task list changes (parent moves, subtask added)

**Option B: Two-level index (parentIndex, subtaskIndex or -1)**
- Track selection as `{parentIdx: 2, subtaskIdx: 1}` = "3rd parent, 2nd subtask"
- Pros: Survives list re-ordering, clearer semantic meaning
- Cons: Complex navigation code, harder to debug, breaks existing patterns

### Decision: **Option A - Flat index across displayed items**

**Rationale**:
1. **Existing code compatibility**: All current navigation uses single int index
2. **Selection staleness acceptable**: Task list only changes on explicit user actions (create, move, delete), not during browsing
3. **Simpler edge case handling**: No need to handle "parent deleted, what happens to subtask index?"
4. **Vim-like behavior**: Users expect j/k to move through visible list sequentially

**Implementation approach**:
- `GetTasksByColumn()` now returns tasks in hierarchical display order
- Navigation code (j/k handlers) unchanged—operates on flat list
- When list refreshes (after create/move/delete), reset `selectedTask = 0` (consistent with current behavior)

**Trade-off accepted**: Selection may jump to top of column after modifying tasks. This matches current behavior (spec doesn't require preserving selection across mutations).

---

## Decision 3: Subtask Ordering in Same Column

### Problem
When parent and subtask are both in the same column (e.g., both in "In Progress"), should subtasks appear immediately below parent (breaking sort order) or respect the column's sort mode (e.g., priority)?

### Options Evaluated

**Option A: Subtasks always immediately follow parent**
- Regardless of sort mode, subtasks appear directly under parent
- Pros: Clear visual grouping, intuitive hierarchy
- Cons: Breaks expected sort behavior (P1 parent with P3 subtask breaks priority ordering)

**Option B: Respect sort order, then group subtasks**
- Sort all tasks by current mode (priority/date/etc.), then for each parent, immediately append its subtasks
- Pros: Honors user's chosen sort mode
- Cons: More complex, parent-subtask grouping may be visually scattered

**Option C: Sort parents only, subtasks unsorted below parent**
- Parents sorted by current mode, subtasks appear in creation order below parent
- Pros: Balance between hierarchy and sort order
- Cons: Subtasks don't benefit from sorting (can't prioritize within parent)

### Decision: **Option A - Subtasks immediately follow parent**

**Rationale**:
1. **User expectation from spec**: "show right below the parent task each subtask" is explicit requirement
2. **Hierarchy over sort**: Visual parent-child relationship is primary goal of this feature
3. **Simplicity**: No need to handle sort edge cases or explain inconsistent behavior
4. **Subtasks can be prioritized independently**: Users can manually order by editing subtask priorities if needed

**Implementation approach**:
- In `BuildFlatHierarchy()`:
  1. Get all parent tasks (ParentID == nil) and sort by current sort mode
  2. For each parent, immediately append its subtasks (ParentID == parent.ID) in their own sort order
  3. Continue to next parent
- Result: Parents are sorted, each parent's subtasks grouped immediately below it

**Trade-off accepted**: Sort mode won't produce global priority ordering when subtasks are present. We accept this because hierarchy visualization is the feature's primary value.

**Future consideration**: If users request "global sort ignoring hierarchy," add a toggle (e.g., "s" cycles through "priority, priority-global, date, date-global").

---

## Decision 4: CLI Indentation Detection

### Problem
How does the CLI `list` and `show` commands know which tasks to indent? Should it check ParentID directly or build a hierarchy structure first?

### Options Evaluated

**Option A: Check ParentID != nil → indent**
- Simple conditional: if `task.ParentID != nil`, print with indentation
- Pros: Fast, simple, works for single-level hierarchy
- Cons: Doesn't handle orphaned subtasks (parent deleted), no validation

**Option B: Build hierarchy map first, then indent based on structure**
- Create parent→children map, only indent tasks that appear in map
- Pros: Validates parent exists, handles edge cases, consistent with TUI
- Cons: More code, slower (two passes: build map + render)

### Decision: **Option B - Build hierarchy map, then indent**

**Rationale**:
1. **Code reuse**: Same `BuildFlatHierarchy()` function used by TUI and CLI
2. **Consistent behavior**: TUI and CLI handle orphaned subtasks identically
3. **Future-proofing**: If we add multi-level nesting, hierarchy map already in place
4. **Performance acceptable**: CLI is interactive (human speed), ~1ms extra latency imperceptible

**Implementation approach**:
- `internal/cli/list.go`:
  ```go
  tasks, _ := storage.ListTasks(db, filters)
  hierarchical := hierarchy.BuildFlatHierarchy(tasks)
  for _, task := range hierarchical {
      indent := hierarchy.GetIndentation(task) // returns 0 or 2
      fmt.Printf("%s- %s\n", strings.Repeat(" ", indent), task.Title)
  }
  ```
- `hierarchy.GetIndentation(task)` returns 2 if ParentID != nil, else 0

**Trade-off accepted**: Slightly more code than simple ParentID check, but ensures TUI/CLI consistency and handles edge cases gracefully.

---

## Decision 5: Form Pre-population Pattern (Optional ParentID)

### Problem
How do we pass parent ID to `initCreateForm()` when user presses 'n' in detail view, without breaking existing calls from kanban view (where parent ID is unknown)?

### Options Evaluated

**Option A: Add optional parameter to initCreateForm(parentID ...string)**
- Variadic parameter allows 0 or 1 parent IDs
- Pros: Backward compatible, single function
- Cons: Variadic params can be confusing in Go, no type safety for "exactly 0 or 1"

**Option B: Create separate function initCreateSubtaskForm(parentID string)**
- Two functions: `initCreateForm()` and `initCreateSubtaskForm(parentID)`
- Pros: Clear intent, type-safe, explicit behavior
- Cons: Code duplication (both initialize same 6 fields)

**Option C: Add optional parameter with pointer *string (nil = no parent)**
- `initCreateForm(parentID *string)` where nil means top-level task
- Pros: Idiomatic Go pattern for optional values, type-safe
- Cons: All callers must pass argument (even if nil)

### Decision: **Option C - Optional pointer parameter**

**Rationale**:
1. **Idiomatic Go**: Pointer for optional string is standard Go pattern
2. **Type safety**: Compiler ensures all callers provide argument (prevents forgetting)
3. **Explicit nil handling**: Forces caller to think about parent vs top-level
4. **Minimal duplication**: Single function, conditional logic only for pre-filling field

**Implementation approach**:
- Signature: `func (m *Model) initCreateForm(parentID *string)`
- From kanban 'n' key: `m.initCreateForm(nil)`
- From detail 'n' key: `m.initCreateForm(&m.detailTask.ID)`
- Inside function:
  ```go
  if parentID != nil {
      m.formInputs[5].SetValue(*parentID)
  }
  ```

**Trade-off accepted**: All existing calls must be updated to pass `nil`, but this is a small one-time change for better long-term clarity.

---

## Decision 6: Subtask Display in Other Columns

### Problem
Spec says: "When the parent is in Inbox column, all its subtasks display in their respective columns." How do we visually indicate parent-child relationship when they're in different columns?

### Options Evaluated

**Option A: Show subtasks only in their actual column (don't duplicate)**
- Subtask in "In Progress" appears only in that column, even if parent is in "Inbox"
- Pros: No duplication, clear separation, each column shows its actual tasks
- Cons: Hierarchy not visible when tasks split across columns

**Option B: Show subtasks in parent's column with "(→ Progress)" indicator**
- Subtask appears indented under parent in parent's column with badge showing actual column
- Pros: Hierarchy always visible, parent-child link clear
- Cons: Confusing UX (task appears in column it's not actually in)

**Option C: Show subtle visual link (color/style) between parent and subtasks across columns**
- Matching color codes or symbols connect related tasks across columns
- Pros: Hierarchy visible without duplicating tasks
- Cons: Complex rendering, limited terminal color palette, cluttered view

### Decision: **Option A - Subtasks only in their actual column**

**Rationale**:
1. **Spec alignment**: Spec explicitly states "subtasks display in their respective columns"
2. **Column integrity**: Each column should show only tasks actually in that column (matches current behavior)
3. **User can view parent details**: Pressing Enter on parent shows all subtasks regardless of column
4. **Simplicity**: No special rendering logic for cross-column relationships

**Implementation approach**:
- Hierarchical grouping only applies when `parent.Column == subtask.Column`
- If subtask is in different column than parent, it appears as top-level task in its own column
- Detail view shows all subtasks with column indicators (already in spec: "Progress", "Done", etc.)

**Trade-off accepted**: Users won't see full hierarchy in kanban view when parent/subtasks split across columns. This is acceptable because:
- Users can press Enter on parent to see all subtasks (detail view)
- Tasks moving across columns is expected workflow (Inbox → In Progress → Done)
- Forcing visual linkage across columns would clutter the interface

---

## Research Artifacts

### Bubble Tea Hierarchy Pattern Research

**Finding**: Bubble Tea doesn't have built-in hierarchy support. Standard pattern is to flatten data before rendering and use styling (indentation, colors) to indicate structure.

**Best practice**:
- Build flat display list in `Update()` or `View()`, not during model mutations
- Use lipgloss `Padding()` or `Margin()` for indentation (not manual spaces)
- Example from Bubble Tea examples:
  ```go
  indent := lipgloss.NewStyle().PaddingLeft(2)
  if isSubtask {
      return indent.Render(taskContent)
  }
  return taskContent
  ```

**Decision impact**: Confirms Option A (flatten on-demand) is standard Bubble Tea pattern. We'll use lipgloss padding for indentation in TUI, plain spaces in CLI.

### Go Optional Parameter Patterns

**Finding**: Go doesn't have optional parameters natively. Three common patterns:
1. Variadic parameters (`func foo(args ...string)`)
2. Functional options (`func foo(opts ...Option)`)
3. Pointer parameters (`func foo(opt *string)`)

**Best practice for simple cases**: Pointer parameters (Option 3) is most idiomatic for single optional value.

**Decision impact**: Confirms Option C (pointer parameter) for `initCreateForm(parentID *string)`.

---

## Open Questions (Resolved)

### Q1: Should we support multi-level nesting (grandchild subtasks)?
**Answer**: No. Spec assumes 1 level (parent → subtask). Implementation will ignore ParentID on tasks that themselves are subtasks (prevents accidental deep nesting). Document this limitation in quickstart.md.

### Q2: What happens when parent is deleted?
**Answer**: Spec mentions this edge case ("orphaned subtasks"). We'll implement option: "block deletion if subtasks exist, require user to delete or move subtasks first." This prevents data loss and is simpler than auto-promoting subtasks.

### Q3: Should archived subtasks appear under their parent?
**Answer**: No. Archive view (`showArchived=true`) and active view (`showArchived=false`) are mutually exclusive. If parent is active and subtask is archived, subtask won't appear (consistent with current filter behavior). Detail view shows all subtasks regardless of archive status (with visual indicator).

---

## Summary

All research complete. Key decisions:
1. ✅ Flatten hierarchy on-demand during render (simplicity)
2. ✅ Use flat index for selection (compatibility)
3. ✅ Subtasks immediately follow parent (hierarchy over sort)
4. ✅ Build hierarchy map in CLI (consistency)
5. ✅ Pointer parameter for optional parent ID (idiomatic Go)
6. ✅ Subtasks display in their actual column only (spec alignment)

No NEEDS CLARIFICATION items remain. Ready for Phase 1 (Design & Contracts).
