# Sample Data Contract

**Feature**: README Visual Enhancement with Screenshots
**Purpose**: Define the requirements and contract for sample task data generation
**Date**: 2025-11-05

## Overview

This contract defines the requirements for the sample dataset that must be created to support README screenshots and CLI examples. While this is not a traditional API contract, it serves as a specification for what the sample data generation script must produce.

## Sample Data Requirements Contract

### Input

**Script Invocation**:
```bash
./scripts/generate-sample-data.sh [--db-path PATH]
```

**Parameters**:
- `--db-path`: Optional custom database path (default: `./sample-tasks.db`)

**Preconditions**:
- OnTop binary (`./ontop`) must exist and be executable
- Database path directory must be writable
- No existing data at target database path (or will be overwritten)

---

### Output Guarantees

The script MUST create a database containing exactly:

#### Parent Tasks (5 required)

| Column | Count | Priority Range | Progress Range | Tags Required |
|--------|-------|----------------|----------------|---------------|
| Inbox | 2 | P1-P3 | 0%-50% | At least 3 unique tags total |
| In Progress | 2 | P2-P4 | 50%-75% | At least 3 unique tags total |
| Done | 1 | P2 | 100% | At least 1 tag |

#### Subtasks (4 minimum required)

| Parent Task | Subtasks | Column Distribution | Progress Range |
|-------------|----------|---------------------|----------------|
| Task in Inbox | 2 | 1 in Inbox, 1 in Done | 0% and 100% |
| Task in In Progress | 2 | 1 in In Progress, 1 in Done | 25% and 100% |

#### Data Characteristics

**MUST include**:
- At least 2 tasks with multi-line descriptions (3+ lines each)
- At least 5 unique tag values across all tasks
- Priority distribution covering P1, P2, P3, P4 (minimum one task per priority)
- Progress distribution: 0%, 25%, 50%, 75%, 100% (at least one task at each value)
- At least one parent task with 2+ subtasks
- At least one subtask in a different column than its parent

**MUST NOT include**:
- More than 10 total tasks (to keep screenshots uncluttered)
- Subtasks of subtasks (max 1 level hierarchy)
- Invalid data (empty titles, priority <1 or >5, progress <0 or >100)

---

### Postconditions

After successful execution:

1. **Database exists** at specified path
2. **Queryable via CLI**:
   ```bash
   ./ontop --db-path ./sample-tasks.db list
   ```
   Returns hierarchical task list with correct indentation

3. **Column distribution verified**:
   ```bash
   ./ontop --db-path ./sample-tasks.db list --column inbox
   ./ontop --db-path ./sample-tasks.db list --column in_progress
   ./ontop --db-path ./sample-tasks.db list --column done
   ```
   Each command returns expected number of tasks

4. **Subtask relationships verified**:
   ```bash
   ./ontop --db-path ./sample-tasks.db show <parent-task-id>
   ```
   Shows all subtasks for parent

5. **Total task count**:
   ```bash
   ./ontop --db-path ./sample-tasks.db list | grep -c "^\["
   ```
   Returns 9 (5 parents + 4 subtasks)

---

## Screenshot Requirements Contract

### TUI Screenshots

Each screenshot MUST meet these criteria:

#### 1. Kanban Board Screenshot
**File**: `docs/images/tui-kanban-board.png`

**Requirements**:
- Shows all three columns (Inbox, In Progress, Done) side by side
- All 5 parent tasks visible with correct priority badges
- At least 2 subtasks visible, indented below parents with dash prefix
- Help text or keyboard shortcuts visible at bottom
- Terminal size: minimum 140x35 characters
- File size: < 500KB
- Format: PNG with transparency removed (solid background)
- Alt text: "OnTop TUI showing kanban board with three columns (Inbox, In Progress, Done) containing parent tasks and indented subtasks"

#### 2. Hierarchical Subtasks Screenshot
**File**: `docs/images/tui-subtasks-hierarchy.png`

**Requirements**:
- Focuses on a column with parent task + 2 subtasks
- Clear visual indentation (2 characters) and dash prefix visible
- Priority indicators visible for parent and subtasks
- Progress percentages visible
- Terminal size: minimum 100x30 characters
- File size: < 500KB
- Format: PNG
- Alt text: "OnTop TUI showing parent task with two subtasks indented below, demonstrating hierarchical task structure"

#### 3. Task Form Screenshot
**File**: `docs/images/tui-task-form.png`

**Requirements**:
- Shows complete task creation or edit form
- All fields visible: Title, Description (textarea), Priority, Progress, Tags, Parent ID
- Description field shows multiline text with line numbers
- Help text visible showing "Ctrl+S to save" and "Tab to next field"
- Terminal size: minimum 100x40 characters
- File size: < 500KB
- Format: PNG
- Alt text: "OnTop task creation form showing all input fields including multiline description textarea with scrollbar"

### CLI Output Examples

#### 1. List Command Output
**Required in README**:

```bash
$ ./ontop list
```

**Output MUST show**:
- All columns labeled (Inbox, In Progress, Done)
- Parent tasks with full information (ID, priority, title, progress)
- Subtasks indented 2 characters below parents with dash prefix
- At least 5 parent tasks + 3 subtasks visible
- Hierarchical grouping maintained (subtasks immediately follow parent)

#### 2. Show Command Output
**Required in README**:

```bash
$ ./ontop show <task-id>
```

**Output MUST show**:
- Full task details (ID, Title, Description, Priority, Column, Progress)
- Tags listed
- Created/Updated timestamps
- List of all subtasks if parent task
- Multiline description properly formatted

#### 3. Add Command Examples
**Required in README**:

```bash
# Create parent task
$ ./ontop add "Implement new feature" --priority 2 --description "Detailed description"

# Create subtask
$ ./ontop add "Research implementation" --parent <parent-id> --priority 3
```

**Output MUST show**:
- Task created successfully with assigned ID
- Confirmation message

---

## Reproducibility Contract

### Regeneration Requirements

The sample data generation script MUST be:

1. **Idempotent**: Running twice produces same result (overwrites existing database)
2. **Deterministic**: Task IDs may vary but content and relationships are consistent
3. **Self-contained**: No external dependencies beyond OnTop binary
4. **Fast**: Completes in under 10 seconds
5. **Error-handling**: Returns non-zero exit code if any command fails

### Documentation Requirements

The script MUST include:

1. **Header comments** explaining purpose and usage
2. **Variable for database path** (configurable)
3. **Comment for each task** explaining its purpose in the sample dataset
4. **Cleanup section** (optional) to remove previous sample database
5. **Verification section** printing summary of created tasks

---

## Validation Checklist

Before considering sample data complete, verify:

- [ ] Script executes without errors
- [ ] Database created at specified path
- [ ] 5 parent tasks exist
- [ ] 4+ subtasks exist with correct parent relationships
- [ ] All 3 columns have tasks
- [ ] Priority range P1-P4 represented
- [ ] Progress range 0%-100% represented
- [ ] At least 5 unique tags used
- [ ] At least 2 multiline descriptions
- [ ] CLI `list` command shows hierarchical structure
- [ ] CLI `show` command displays subtasks for parents
- [ ] Screenshots can be captured showing all required elements
- [ ] Total database size < 100KB

---

## Exit Codes

The sample data generation script MUST use these exit codes:

- **0**: Success - all tasks created successfully
- **1**: OnTop binary not found or not executable
- **2**: Database path not writable
- **3**: Task creation command failed
- **4**: Verification failed (wrong number of tasks or missing relationships)

---

## Contract Compliance

This contract ensures:

1. **Reproducibility**: Anyone can regenerate sample data identically
2. **Screenshot quality**: All required elements will be visible in screenshots
3. **Documentation accuracy**: CLI examples will show realistic, consistent output
4. **Maintainability**: When UI changes, sample data can be quickly recreated

The sample data generation script (`scripts/generate-sample-data.sh`) MUST satisfy all requirements in this contract for the feature to be considered complete.
