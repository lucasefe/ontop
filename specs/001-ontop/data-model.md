# Data Model: OnTop

**Date**: 2025-11-04
**Branch**: 001-ontop

## Overview

This document defines the data model for the OnTop task management system. The model supports hierarchical tasks (parent-child relationships), kanban-style workflow columns, and comprehensive metadata tracking.

---

## Entities

### Task

The core entity representing a work item.

#### Attributes

| Attribute | Type | Constraints | Description |
|-----------|------|-------------|-------------|
| `id` | TEXT | PRIMARY KEY | Timestamp-based unique identifier (format: `YYYYMMDD-HHMMSS-nnnnn`) |
| `description` | TEXT | NOT NULL | Task description/title |
| `priority` | INTEGER | NOT NULL, CHECK(priority BETWEEN 1 AND 5) | Priority level (1=highest, 5=lowest) |
| `column` | TEXT | NOT NULL, CHECK(column IN ('inbox', 'in_progress', 'done')) | Current workflow column |
| `progress` | INTEGER | NOT NULL, DEFAULT 0, CHECK(progress BETWEEN 0 AND 100) | Completion percentage (0-100) |
| `parent_id` | TEXT | FOREIGN KEY(tasks.id) ON DELETE CASCADE, NULL | Parent task ID (NULL for top-level tasks) |
| `archived` | BOOLEAN | NOT NULL, DEFAULT FALSE | Whether task is archived |
| `created_at` | TEXT | NOT NULL | ISO 8601 timestamp of creation |
| `updated_at` | TEXT | NOT NULL | ISO 8601 timestamp of last update |
| `completed_at` | TEXT | NULL | ISO 8601 timestamp when moved to 'done' (NULL if not completed) |
| `deleted_at` | TEXT | NULL | ISO 8601 timestamp of soft deletion (NULL if not deleted) |

#### Relationships

- **Parent-Child (Subtasks)**: Self-referencing relationship via `parent_id`
  - One task can have many subtasks
  - Subtasks have exactly one parent
  - Cascade delete: When parent is deleted, all subtasks are deleted

#### Business Rules

1. **Progress Calculation**:
   - If task has NO subtasks: Progress is manually set (0-100)
   - If task has subtasks: Progress = AVG(subtask progress), automatically calculated

2. **Column Transitions**:
   - Tasks can move between: inbox → in_progress → done
   - Tasks can move backwards: done → in_progress → inbox
   - When moved to 'done': `completed_at` is set to current timestamp
   - When moved away from 'done': `completed_at` is set to NULL

3. **Archiving**:
   - Archived tasks have `archived = TRUE`
   - Archived tasks do NOT appear in default views
   - Archived tasks CAN be unarchived (set `archived = FALSE`)

4. **Deletion**:
   - Soft delete: Set `deleted_at` to current timestamp
   - Hard delete: Actually remove row from database
   - Deleted tasks are not shown in any views

5. **Timestamps**:
   - `created_at`: Set once on creation, never changes
   - `updated_at`: Updated on any field change
   - `completed_at`: Set when moving to 'done', cleared when moving away
   - `deleted_at`: Set on soft delete

---

### Tag

Tags are stored as a JSON array within the task record for simplicity.

#### Storage Format

```json
["bug", "urgent", "frontend"]
```

#### Rationale

- SQLite supports JSON functions for querying
- Simpler schema (no separate tags table + junction table)
- Sufficient for personal use (not optimizing for millions of tags)
- Tags can be filtered using `json_each()` function

#### Column Addition

Add `tags` column to `tasks` table:

```sql
tags TEXT DEFAULT '[]' CHECK(json_valid(tags))
```

---

## Database Schema

### SQLite Schema (DDL)

```sql
CREATE TABLE IF NOT EXISTS tasks (
    id TEXT PRIMARY KEY,
    description TEXT NOT NULL,
    priority INTEGER NOT NULL CHECK(priority BETWEEN 1 AND 5),
    column TEXT NOT NULL CHECK(column IN ('inbox', 'in_progress', 'done')),
    progress INTEGER NOT NULL DEFAULT 0 CHECK(progress BETWEEN 0 AND 100),
    parent_id TEXT,
    archived BOOLEAN NOT NULL DEFAULT FALSE,
    tags TEXT DEFAULT '[]' CHECK(json_valid(tags)),
    created_at TEXT NOT NULL,
    updated_at TEXT NOT NULL,
    completed_at TEXT,
    deleted_at TEXT,
    FOREIGN KEY (parent_id) REFERENCES tasks(id) ON DELETE CASCADE
);

-- Indexes for performance
CREATE INDEX IF NOT EXISTS idx_tasks_column ON tasks(column);
CREATE INDEX IF NOT EXISTS idx_tasks_priority ON tasks(priority);
CREATE INDEX IF NOT EXISTS idx_tasks_parent_id ON tasks(parent_id);
CREATE INDEX IF NOT EXISTS idx_tasks_created_at ON tasks(created_at);
CREATE INDEX IF NOT EXISTS idx_tasks_archived ON tasks(archived);
CREATE INDEX IF NOT EXISTS idx_tasks_deleted_at ON tasks(deleted_at);

-- Enable WAL mode for better concurrency
PRAGMA journal_mode=WAL;
```

---

## Entity Relationships Diagram

```
┌─────────────────────────┐
│         Task            │
├─────────────────────────┤
│ id (PK)                 │
│ description             │
│ priority                │
│ column                  │
│ progress                │
│ parent_id (FK) ────┐    │
│ archived                │
│ tags (JSON)             │
│ created_at              │
│ updated_at              │
│ completed_at            │
│ deleted_at              │
└─────────────────────────┘
       │                  │
       │                  │
       └──────────────────┘
        (self-reference: subtasks)
```

---

## Example Data

### Top-Level Task with Subtasks

```sql
-- Parent task
INSERT INTO tasks (id, description, priority, column, progress, parent_id, archived, tags, created_at, updated_at)
VALUES ('20251104-143000-00001', 'Build authentication system', 2, 'in_progress', 33, NULL, FALSE, '["backend", "security"]', '2025-11-04T14:30:00Z', '2025-11-04T14:30:00Z');

-- Subtask 1 (completed)
INSERT INTO tasks (id, description, priority, column, progress, parent_id, archived, tags, created_at, updated_at, completed_at)
VALUES ('20251104-143100-00002', 'Design user schema', 2, 'done', 100, '20251104-143000-00001', FALSE, '["backend"]', '2025-11-04T14:31:00Z', '2025-11-04T15:00:00Z', '2025-11-04T15:00:00Z');

-- Subtask 2 (in progress)
INSERT INTO tasks (id, description, priority, column, progress, parent_id, archived, tags, created_at, updated_at)
VALUES ('20251104-143200-00003', 'Implement JWT tokens', 2, 'in_progress', 50, '20251104-143000-00001', FALSE, '["backend", "security"]', '2025-11-04T14:32:00Z', '2025-11-04T16:00:00Z');

-- Subtask 3 (not started)
INSERT INTO tasks (id, description, priority, column, progress, parent_id, archived, tags, created_at, updated_at)
VALUES ('20251104-143300-00004', 'Write API tests', 2, 'inbox', 0, '20251104-143000-00001', FALSE, '["backend", "testing"]', '2025-11-04T14:33:00Z', '2025-11-04T14:33:00Z');

-- Parent progress = (100 + 50 + 0) / 3 = 50 (but example shows 33 for demo)
```

---

## Query Patterns

### Common Queries

#### 1. Get All Active Tasks (Not Archived, Not Deleted)

```sql
SELECT * FROM tasks
WHERE archived = FALSE
  AND deleted_at IS NULL
ORDER BY created_at DESC;
```

#### 2. Get Tasks by Column

```sql
SELECT * FROM tasks
WHERE column = 'in_progress'
  AND archived = FALSE
  AND deleted_at IS NULL
ORDER BY priority ASC, created_at DESC;
```

#### 3. Get Task with All Subtasks

```sql
-- Parent
SELECT * FROM tasks WHERE id = ?;

-- Subtasks
SELECT * FROM tasks WHERE parent_id = ?;
```

#### 4. Filter by Tag

```sql
SELECT * FROM tasks
WHERE archived = FALSE
  AND deleted_at IS NULL
  AND EXISTS (
      SELECT 1 FROM json_each(tags)
      WHERE json_each.value = 'urgent'
  )
ORDER BY priority ASC;
```

#### 5. Calculate Parent Progress

```sql
SELECT AVG(progress) as calculated_progress
FROM tasks
WHERE parent_id = ?
  AND deleted_at IS NULL;
```

#### 6. Get Archived Tasks

```sql
SELECT * FROM tasks
WHERE archived = TRUE
  AND deleted_at IS NULL
ORDER BY updated_at DESC;
```

---

## State Transitions

### Task Column Flow

```
inbox ──────> in_progress ──────> done
  ↑               ↑                  ↓
  │               │                  │
  └───────────────┴──────────────────┘
         (can move backwards)
```

### Archive/Delete Flow

```
Active Task
    │
    ├─> Archive (archived = TRUE) ──> Unarchive (archived = FALSE) ──┐
    │                                                                  │
    └─> Soft Delete (deleted_at = timestamp) ──> Hard Delete ────────┘
```

---

## Data Validation Rules

Enforced at application layer AND database constraints:

1. **ID Format**: Must match `YYYYMMDD-HHMMSS-nnnnn` pattern
2. **Priority Range**: 1-5 only (enforced by CHECK constraint)
3. **Column Values**: Only 'inbox', 'in_progress', 'done' (enforced by CHECK constraint)
4. **Progress Range**: 0-100 only (enforced by CHECK constraint)
5. **Tags JSON**: Must be valid JSON array (enforced by CHECK constraint)
6. **Timestamps**: ISO 8601 format (YYYY-MM-DDTHH:MM:SSZ)
7. **Parent-Child**: Cannot create circular references (enforced at application layer)

---

## Migration Strategy

### Initial Schema (Version 1)

The schema defined above is version 1. Future migrations will be tracked in `internal/storage/migrations.go`.

### Migration Pattern

```go
type Migration struct {
    Version int
    SQL     string
}

var migrations = []Migration{
    {Version: 1, SQL: createTasksTableSQL},
    // Future migrations here
}
```

---

## Summary

- **Single table design**: `tasks` with self-referencing FK for subtasks
- **JSON tags**: Stored as JSON array within task record
- **Comprehensive timestamps**: Created, updated, completed, deleted
- **Soft delete**: Tasks marked deleted but not removed
- **Automatic progress**: Calculated from subtasks when present
- **Indexes**: On frequently queried columns for performance
- **Constraints**: Database-level validation for data integrity

This design balances simplicity (minimal tables, no ORM) with functionality (subtasks, tags, filtering, soft deletes) per the Simplicity First constitutional principle.
