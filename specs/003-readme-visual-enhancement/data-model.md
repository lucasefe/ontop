# Data Model: Sample Task Dataset

**Feature**: README Visual Enhancement with Screenshots
**Date**: 2025-11-05
**Purpose**: Define the structure and content of sample tasks used for screenshot generation

## Overview

This feature doesn't modify application data models. Instead, it defines a **sample dataset** of tasks that will be created for demonstration purposes in screenshots and CLI examples.

## Sample Dataset Structure

### Parent Tasks (5 total)

The sample dataset includes 5 parent tasks distributed across columns to demonstrate the kanban workflow:

#### Inbox Column (2 tasks)

**Task 1: User Authentication Feature**
- **Title**: "Implement user authentication"
- **Description**: Multi-line text:
  ```
  Add OAuth2 and session management
  Support GitHub and Google providers
  Include remember-me functionality
  ```
- **Priority**: 1 (highest)
- **Progress**: 0%
- **Tags**: ["feature", "security"]
- **Parent ID**: null
- **Purpose**: Demonstrates high-priority feature planning with detailed multi-line description

**Task 2: Documentation Update**
- **Title**: "Update API documentation"
- **Description**: "Review and update all API endpoint documentation to reflect recent changes"
- **Priority**: 3 (medium)
- **Progress**: 25%
- **Tags**: ["documentation"]
- **Parent ID**: null
- **Purpose**: Demonstrates medium-priority task in progress

#### In Progress Column (2 tasks)

**Task 3: Memory Leak Fix**
- **Title**: "Fix memory leak in task sync"
- **Description**: "Investigate and resolve memory growth issue in background sync process"
- **Priority**: 2 (high)
- **Progress**: 75%
- **Tags**: ["bug", "performance"]
- **Parent ID**: null
- **Purpose**: Demonstrates high-priority bug fix nearing completion

**Task 4: Search Feature Enhancement**
- **Title**: "Add fuzzy search to task titles"
- **Description**: Multi-line text:
  ```
  Implement fuzzy matching for task search
  Consider using fzf-like algorithm
  Test with large task datasets
  ```
- **Priority**: 4 (low)
- **Progress**: 50%
- **Tags**: ["feature", "enhancement"]
- **Parent ID**: null
- **Purpose**: Demonstrates lower-priority enhancement with subtasks (see below)

#### Done Column (1 task)

**Task 5: Database Migration**
- **Title**: "Migrate to WAL mode for SQLite"
- **Description**: "Enable Write-Ahead Logging for better concurrent access"
- **Priority**: 2 (high)
- **Progress**: 100%
- **Tags**: ["infrastructure", "database"]
- **Parent ID**: null
- **Purpose**: Demonstrates completed high-priority infrastructure work

---

### Subtasks (3+ total)

Subtasks demonstrate parent-child task relationships and hierarchy:

**Subtask 1: OAuth Library Research** (child of Task 1)
- **Title**: "Research OAuth2 libraries"
- **Description**: "Evaluate golang.org/x/oauth2 and alternatives"
- **Priority**: 2 (high)
- **Progress**: 100%
- **Tags**: ["research"]
- **Parent ID**: Task 1 ID
- **Column**: Done (different from parent's Inbox)
- **Purpose**: Shows completed research subtask for planned feature

**Subtask 2: GitHub Provider Implementation** (child of Task 1)
- **Title**: "Implement GitHub OAuth provider"
- **Description**: "Create provider for GitHub authentication flow"
- **Priority**: 2 (high)
- **Progress**: 0%
- **Tags**: ["development"]
- **Parent ID**: Task 1 ID
- **Column**: Inbox (same as parent)
- **Purpose**: Shows pending subtask that will start after research

**Subtask 3: Search Algorithm Selection** (child of Task 4)
- **Title**: "Evaluate fuzzy search algorithms"
- **Description**: "Compare Levenshtein distance vs. trigram matching"
- **Priority**: 4 (low)
- **Progress**: 100%
- **Tags**: ["research"]
- **Parent ID**: Task 4 ID
- **Column**: Done
- **Purpose**: Shows research completed for in-progress feature

**Subtask 4: Search UI Integration** (child of Task 4)
- **Title**: "Add search input to TUI"
- **Description**: "Integrate search field into kanban view header"
- **Priority**: 4 (low)
- **Progress**: 25%
- **Tags**: ["development", "ui"]
- **Parent ID**: Task 4 ID
- **Column**: In Progress (same as parent)
- **Purpose**: Shows active development subtask

---

## Sample Data Characteristics

### Priority Distribution
- **P1 (highest)**: 1 task (User Authentication)
- **P2 (high)**: 4 tasks/subtasks (Memory Leak, OAuth Research, GitHub Provider, Database Migration)
- **P3 (medium)**: 1 task (Documentation Update)
- **P4 (low)**: 3 tasks/subtasks (Search Enhancement + 2 subtasks)
- **P5 (lowest)**: 0 tasks

**Rationale**: Weighted toward higher priorities to demonstrate realistic project prioritization

### Progress Distribution
- **0%**: 2 tasks (not started)
- **25%**: 2 tasks (early stage)
- **50%**: 1 task (midpoint)
- **75%**: 1 task (nearly complete)
- **100%**: 3 tasks (completed)

**Rationale**: Shows full range of progress states from planning to completion

### Column Distribution
- **Inbox**: 2 parent tasks + 1 subtask (3 total)
- **In Progress**: 2 parent tasks + 1 subtask (3 total)
- **Done**: 1 parent task + 2 subtasks (3 total)

**Rationale**: Balanced distribution shows active workflow with visible progress

### Tag Distribution
- **feature**: 3 tasks (new functionality)
- **bug**: 1 task (defect fixes)
- **documentation**: 2 tasks (docs work)
- **security**: 1 task (security-related)
- **performance**: 1 task (optimization)
- **research**: 3 tasks (investigation work)
- **development**: 2 tasks (implementation)
- **infrastructure**: 1 task (system-level)

**Rationale**: Diverse tags demonstrate realistic project categorization

### Hierarchical Relationships
- **Parent tasks**: 5 (3 with subtasks, 2 standalone)
- **Subtasks**: 4 total
- **Parent-child relationships**:
  - Task 1 (User Authentication) → 2 subtasks
  - Task 4 (Search Enhancement) → 2 subtasks
  - Tasks 2, 3, 5 → no subtasks

**Rationale**: Demonstrates both simple tasks and complex features broken into subtasks

---

## Visual Representation Goals

This sample dataset is designed to showcase:

1. **Hierarchical display**: Subtasks visually indented below parents with dash prefix
2. **Cross-column relationships**: Subtasks can be in different columns than parents
3. **Priority-based sorting**: Tasks displayed by priority within columns
4. **Progress indicators**: Visual progress bars and percentages
5. **Tag variety**: Multiple tags showing different project dimensions
6. **Multiline descriptions**: Textarea functionality with newlines preserved
7. **Workflow stages**: Clear progression from Inbox → In Progress → Done

## Sample Data Generation

The sample dataset will be created using the `scripts/generate-sample-data.sh` script, which uses OnTop CLI commands to populate a separate database (not the user's production database).

**Script approach**:
```bash
# Use isolated database
export DB_PATH="./sample-tasks.db"

# Create parent tasks
ontop add "Implement user authentication" --priority 1 --description "..." --tags "feature,security"
TASK1_ID=$(ontop list --column inbox | grep "authentication" | awk '{print $1}')

# Create subtasks with parent relationship
ontop add "Research OAuth2 libraries" --parent $TASK1_ID --priority 2 --tags "research"

# Move completed subtasks to Done column
SUBTASK1_ID=$(ontop list | grep "Research OAuth2" | awk '{print $1}')
ontop move $SUBTASK1_ID done
ontop update $SUBTASK1_ID --progress 100
```

---

## No Application Data Model Changes

**Important**: This feature does NOT modify the application's Task data model. The existing Task structure (as defined in `internal/models/task.go`) already supports all required fields:

- ID, Title, Description
- Priority (1-5), Progress (0-100)
- Column (inbox, in_progress, done)
- ParentID (for subtask relationships)
- Tags (array of strings)
- Archived, CreatedAt, UpdatedAt, CompletedAt, DeletedAt

The sample dataset simply uses these existing fields to create representative tasks for documentation purposes.
