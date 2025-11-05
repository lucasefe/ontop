# CLI Command Contracts: OnTop

**Date**: 2025-11-04
**Branch**: 001-ontop

## Overview

This document specifies the command-line interface contracts for OnTop. All commands follow Unix conventions: exit code 0 for success, non-zero for errors, text output to stdout, errors to stderr.

---

## Global Flags

Available for all commands:

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--db-path` | string | `~/.ontop/tasks.db` | Path to SQLite database file |
| `--json` | bool | false | Output in JSON format (where applicable) |
| `--help` | bool | false | Show help for command |

---

## Commands

### 1. Launch TUI (Default)

**Command**: `ontop` (no arguments)

**Description**: Launches the interactive TUI mode with kanban board view.

**Usage**:
```bash
ontop
```

**Behavior**:
- Opens full-screen TUI
- Displays three columns: Inbox, In Progress, Done
- Enables vim-style navigation
- Shows active (non-archived, non-deleted) tasks only

**Exit Codes**:
- 0: User quit normally
- 1: Error initializing TUI or database

---

### 2. Add Task

**Command**: `ontop add <description>`

**Description**: Creates a new task with the given description.

**Usage**:
```bash
ontop add "Fix login bug" --priority 1 --tags bug,urgent
ontop add "Write documentation" -p 3 -t docs
```

**Flags**:
| Flag | Alias | Type | Default | Description |
|------|-------|------|---------|-------------|
| `--priority` | `-p` | int | 3 | Priority (1-5, where 1=highest) |
| `--tags` | `-t` | string | "" | Comma-separated tags |
| `--parent` | - | string | "" | Parent task ID (for subtasks) |

**Output** (text):
```
Created task: 20251104-143000-00001
Description: Fix login bug
Priority: 1
Column: inbox
Tags: bug, urgent
```

**Output** (JSON with `--json`):
```json
{
  "id": "20251104-143000-00001",
  "description": "Fix login bug",
  "priority": 1,
  "column": "inbox",
  "progress": 0,
  "parent_id": null,
  "archived": false,
  "tags": ["bug", "urgent"],
  "created_at": "2025-11-04T14:30:00Z",
  "updated_at": "2025-11-04T14:30:00Z"
}
```

**Exit Codes**:
- 0: Task created successfully
- 1: Database error
- 2: Invalid arguments (priority out of range, parent not found, etc.)

---

### 3. List Tasks

**Command**: `ontop list`

**Description**: Lists all active tasks.

**Usage**:
```bash
ontop list
ontop list --column inbox
ontop list --priority 1
ontop list --tags urgent
ontop list --json
```

**Flags**:
| Flag | Type | Description |
|------|------|-------------|
| `--column` | string | Filter by column (inbox, in_progress, done) |
| `--priority` | int | Filter by priority (1-5) |
| `--tags` | string | Filter by comma-separated tags (AND logic) |
| `--archived` | bool | Show archived tasks instead of active |

**Output** (text):
```
ID                      Description                 Pri  Column       Progress  Tags
20251104-143000-00001   Fix login bug               1    inbox        0%        bug, urgent
20251104-143100-00002   Write documentation         3    in_progress  50%       docs
20251104-143200-00003   Deploy to production        2    done         100%      deploy

3 tasks found
```

**Output** (JSON with `--json`):
```json
{
  "tasks": [
    {
      "id": "20251104-143000-00001",
      "description": "Fix login bug",
      "priority": 1,
      "column": "inbox",
      "progress": 0,
      "tags": ["bug", "urgent"]
    }
  ],
  "count": 1
}
```

**Exit Codes**:
- 0: Success (even if 0 tasks found)
- 1: Database error
- 2: Invalid filter arguments

---

### 4. Show Task

**Command**: `ontop show <id>`

**Description**: Displays full details of a task including all subtasks.

**Usage**:
```bash
ontop show 20251104-143000-00001
ontop show 20251104-143000-00001 --json
```

**Output** (text):
```
Task: 20251104-143000-00001
Description: Build authentication system
Priority: 2
Column: in_progress
Progress: 50%
Tags: backend, security
Created: 2025-11-04 14:30:00
Updated: 2025-11-04 16:00:00
Completed: -
Archived: No

Subtasks (3):
  ✓ 20251104-143100-00002  Design user schema         (done, 100%)
  → 20251104-143200-00003  Implement JWT tokens      (in_progress, 50%)
    20251104-143300-00004  Write API tests           (inbox, 0%)
```

**Output** (JSON with `--json`):
```json
{
  "task": {
    "id": "20251104-143000-00001",
    "description": "Build authentication system",
    "priority": 2,
    "column": "in_progress",
    "progress": 50,
    "parent_id": null,
    "archived": false,
    "tags": ["backend", "security"],
    "created_at": "2025-11-04T14:30:00Z",
    "updated_at": "2025-11-04T16:00:00Z",
    "completed_at": null,
    "deleted_at": null
  },
  "subtasks": [
    {
      "id": "20251104-143100-00002",
      "description": "Design user schema",
      "column": "done",
      "progress": 100
    }
  ]
}
```

**Exit Codes**:
- 0: Task found and displayed
- 1: Database error
- 2: Task not found

---

### 5. Move Task

**Command**: `ontop move <id> <column>`

**Description**: Moves a task to a different column.

**Usage**:
```bash
ontop move 20251104-143000-00001 in_progress
ontop move 20251104-143000-00001 done
```

**Arguments**:
- `<id>`: Task ID
- `<column>`: Target column (inbox, in_progress, done)

**Output** (text):
```
Task 20251104-143000-00001 moved to in_progress
```

**Output** (JSON with `--json`):
```json
{
  "id": "20251104-143000-00001",
  "old_column": "inbox",
  "new_column": "in_progress",
  "completed_at": null
}
```

**Behavior**:
- If moving to 'done': Sets `completed_at` timestamp
- If moving away from 'done': Clears `completed_at` timestamp
- Updates `updated_at` timestamp

**Exit Codes**:
- 0: Task moved successfully
- 1: Database error
- 2: Invalid arguments (task not found, invalid column)

---

### 6. Archive Task

**Command**: `ontop archive <id>`

**Description**: Archives a task (removes from active views but keeps data).

**Usage**:
```bash
ontop archive 20251104-143000-00001
```

**Output** (text):
```
Task 20251104-143000-00001 archived
```

**Exit Codes**:
- 0: Task archived successfully
- 1: Database error
- 2: Task not found

---

### 7. Unarchive Task

**Command**: `ontop unarchive <id>`

**Description**: Restores an archived task to active status.

**Usage**:
```bash
ontop unarchive 20251104-143000-00001
```

**Output** (text):
```
Task 20251104-143000-00001 unarchived
```

**Exit Codes**:
- 0: Task unarchived successfully
- 1: Database error
- 2: Task not found or not archived

---

### 8. Delete Task

**Command**: `ontop delete <id>`

**Description**: Soft-deletes a task (sets deleted_at timestamp).

**Usage**:
```bash
ontop delete 20251104-143000-00001
ontop delete 20251104-143000-00001 --hard
```

**Flags**:
| Flag | Description |
|------|-------------|
| `--hard` | Permanently delete (remove from database) |

**Output** (text):
```
Task 20251104-143000-00001 deleted
```

**Warning**: Hard delete with subtasks:
```
Warning: Task has 3 subtasks. They will also be deleted.
Confirm? (y/N):
```

**Exit Codes**:
- 0: Task deleted successfully
- 1: Database error
- 2: Task not found or user canceled hard delete

---

### 9. Update Task

**Command**: `ontop update <id>`

**Description**: Updates task attributes.

**Usage**:
```bash
ontop update 20251104-143000-00001 --description "Fixed login bug"
ontop update 20251104-143000-00001 --priority 2
ontop update 20251104-143000-00001 --tags bug,fixed
ontop update 20251104-143000-00001 --progress 75
```

**Flags**:
| Flag | Type | Description |
|------|------|-------------|
| `--description` | string | New description |
| `--priority` | int | New priority (1-5) |
| `--tags` | string | New comma-separated tags (replaces existing) |
| `--progress` | int | New progress (0-100, only for tasks without subtasks) |

**Output** (text):
```
Task 20251104-143000-00001 updated
```

**Exit Codes**:
- 0: Task updated successfully
- 1: Database error
- 2: Invalid arguments (task not found, invalid values, trying to set progress on task with subtasks)

---

### 10. Filter Tasks

**Command**: `ontop filter`

**Description**: Advanced filtering with multiple criteria (alias for `list` with filters).

**Usage**:
```bash
ontop filter --priority 1 --column inbox
ontop filter --tags urgent,bug
ontop filter --archived
```

**Note**: This is essentially an alias for `ontop list` with the same filtering options. Included for discoverability.

---

## Error Handling

### Common Error Messages

**Database Errors**:
```
Error: Cannot connect to database at ~/.ontop/tasks.db
Error: Database file is corrupted
```

**Validation Errors**:
```
Error: Invalid priority. Must be between 1 and 5
Error: Invalid column. Must be one of: inbox, in_progress, done
Error: Task not found: 20251104-143000-99999
Error: Parent task not found: 20251104-143000-99999
```

**Permission Errors**:
```
Error: Cannot write to ~/.ontop/tasks.db (permission denied)
```

---

## JSON Output Schema

When `--json` flag is used, all commands output valid JSON with this structure:

**Success Response**:
```json
{
  "success": true,
  "data": { /* command-specific data */ }
}
```

**Error Response**:
```json
{
  "success": false,
  "error": {
    "code": "TASK_NOT_FOUND",
    "message": "Task not found: 20251104-143000-99999"
  }
}
```

---

## Help Text

Each command provides help via `--help`:

```bash
$ ontop add --help
Usage: ontop add <description> [flags]

Creates a new task with the given description.

Flags:
  -p, --priority int     Task priority (1-5, where 1=highest) (default 3)
  -t, --tags string      Comma-separated tags
      --parent string    Parent task ID (for creating subtasks)
      --db-path string   Path to database file (default "~/.ontop/tasks.db")
      --json             Output in JSON format
  -h, --help             Show this help message

Examples:
  ontop add "Fix login bug" --priority 1 --tags bug,urgent
  ontop add "Subtask" --parent 20251104-143000-00001
```

---

## Summary

- **10 commands**: Default TUI, add, list, show, move, archive, unarchive, delete, update, filter
- **Consistent interface**: All commands follow Unix conventions
- **Flexible output**: Text (human-readable) or JSON (machine-readable)
- **Global flags**: --db-path, --json, --help
- **Clear exit codes**: 0=success, 1=error, 2=validation/usage error
- **Comprehensive help**: Every command has `--help` documentation

This CLI design supports both interactive human use and scripted automation per the Dual-Mode Architecture constitutional principle.
