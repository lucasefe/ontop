# Quick Start Guide: Implementing README Visual Enhancement

**Feature**: README Visual Enhancement with Screenshots
**Audience**: Developer implementing this feature
**Estimated Time**: 2-3 hours

## Prerequisites

Before starting implementation:

- [ ] OnTop application builds successfully (`make build`)
- [ ] You have screenshot capability (macOS: Command+Shift+4, Linux: screenshot tool)
- [ ] Terminal configured with Gruvbox theme (or similar dark theme)
- [ ] Terminal size adjustable to at least 140x35 characters
- [ ] PNG optimization tool available (optional: `optipng`, `pngcrush`)

## Implementation Steps

### Step 1: Create Directory Structure (5 minutes)

Create directories for images and scripts:

```bash
# From repository root
mkdir -p docs/images
mkdir -p scripts

# Verify directories exist
ls -la docs/images
ls -la scripts
```

**Expected result**: Empty `docs/images/` and `scripts/` directories created

---

### Step 2: Create Sample Data Generation Script (30 minutes)

Create `scripts/generate-sample-data.sh`:

```bash
#!/bin/bash
# Purpose: Generate sample tasks for README screenshots
# Usage: ./scripts/generate-sample-data.sh [--db-path PATH]

set -e  # Exit on any error

# Parse arguments
DB_PATH="${1:-./sample-tasks.db}"

# Clean up existing sample database
rm -f "$DB_PATH"

echo "Creating sample tasks in $DB_PATH..."

# Use custom database path
export ONTOP_DB_PATH="$DB_PATH"

# Create parent tasks (5 total across columns)
# Task 1: High-priority feature (Inbox)
TASK1_ID=$(./ontop add "Implement user authentication" \
  --priority 1 \
  --description "Add OAuth2 and session management
Support GitHub and Google providers
Include remember-me functionality" \
  --tags "feature,security" \
  | grep -oE '[a-f0-9]{6}')

# Task 2: Medium-priority docs (Inbox)
TASK2_ID=$(./ontop add "Update API documentation" \
  --priority 3 \
  --progress 25 \
  --description "Review and update all API endpoint documentation to reflect recent changes" \
  --tags "documentation" \
  | grep -oE '[a-f0-9]{6}')

# Task 3: High-priority bug (In Progress)
TASK3_ID=$(./ontop add "Fix memory leak in task sync" \
  --priority 2 \
  --progress 75 \
  --description "Investigate and resolve memory growth issue in background sync process" \
  --tags "bug,performance" \
  | grep -oE '[a-f0-9]{6}')
./ontop move "$TASK3_ID" in_progress

# Task 4: Low-priority feature (In Progress)
TASK4_ID=$(./ontop add "Add fuzzy search to task titles" \
  --priority 4 \
  --progress 50 \
  --description "Implement fuzzy matching for task search
Consider using fzf-like algorithm
Test with large task datasets" \
  --tags "feature,enhancement" \
  | grep -oE '[a-f0-9]{6}')
./ontop move "$TASK4_ID" in_progress

# Task 5: Completed infrastructure (Done)
TASK5_ID=$(./ontop add "Migrate to WAL mode for SQLite" \
  --priority 2 \
  --progress 100 \
  --description "Enable Write-Ahead Logging for better concurrent access" \
  --tags "infrastructure,database" \
  | grep -oE '[a-f0-9]{6}')
./ontop move "$TASK5_ID" done

# Create subtasks
# Subtask 1-1: Research for Task 1 (completed, in Done)
SUBTASK1_1=$(./ontop add "Research OAuth2 libraries" \
  --parent "$TASK1_ID" \
  --priority 2 \
  --progress 100 \
  --description "Evaluate golang.org/x/oauth2 and alternatives" \
  --tags "research" \
  | grep -oE '[a-f0-9]{6}')
./ontop move "$SUBTASK1_1" done

# Subtask 1-2: Implementation for Task 1 (not started, in Inbox)
SUBTASK1_2=$(./ontop add "Implement GitHub OAuth provider" \
  --parent "$TASK1_ID" \
  --priority 2 \
  --description "Create provider for GitHub authentication flow" \
  --tags "development" \
  | grep -oE '[a-f0-9]{6}')

# Subtask 4-1: Research for Task 4 (completed, in Done)
SUBTASK4_1=$(./ontop add "Evaluate fuzzy search algorithms" \
  --parent "$TASK4_ID" \
  --priority 4 \
  --progress 100 \
  --description "Compare Levenshtein distance vs. trigram matching" \
  --tags "research" \
  | grep -oE '[a-f0-9]{6}')
./ontop move "$SUBTASK4_1" done

# Subtask 4-2: Development for Task 4 (in progress)
SUBTASK4_2=$(./ontop add "Add search input to TUI" \
  --parent "$TASK4_ID" \
  --priority 4 \
  --progress 25 \
  --description "Integrate search field into kanban view header" \
  --tags "development,ui" \
  | grep -oE '[a-f0-9]{6}')
./ontop move "$SUBTASK4_2" in_progress

echo "✓ Sample data created successfully!"
echo ""
echo "Verify with:"
echo "  ./ontop --db-path $DB_PATH list"
echo ""
echo "Launch TUI with:"
echo "  ./ontop --db-path $DB_PATH"
```

**Make script executable**:
```bash
chmod +x scripts/generate-sample-data.sh
```

**Test the script**:
```bash
./scripts/generate-sample-data.sh
./ontop --db-path ./sample-tasks.db list
```

**Expected result**: List shows 5 parent tasks and 4 subtasks across all three columns with hierarchical indentation

---

### Step 3: Generate Sample Data and Launch TUI (10 minutes)

```bash
# Generate sample tasks
./scripts/generate-sample-data.sh

# Launch OnTop with sample database
./ontop --db-path ./sample-tasks.db
```

**Navigate the interface**:
- Press `h`/`l` to move between columns
- Press `j`/`k` to move up/down through tasks
- Verify hierarchical display (subtasks indented below parents)
- Press `Enter` on a parent task to view details
- Press `n` while viewing a task to see subtask creation form
- Press `e` on any task to see the edit form

**Expected result**: TUI displays sample tasks correctly with visual hierarchy

---

### Step 4: Capture TUI Screenshots (30 minutes)

#### Screenshot 1: Kanban Board (main view)

1. Resize terminal to 140x35 characters
2. Launch: `./ontop --db-path ./sample-tasks.db`
3. Position cursor on a task in the Inbox column
4. Ensure help panel is visible at bottom (press `?` if hidden)
5. Capture screenshot (Command+Shift+4 on macOS)
6. Save as `docs/images/tui-kanban-board.png`
7. Verify file size < 500KB

#### Screenshot 2: Hierarchical Subtasks

1. Navigate to a parent task with subtasks (e.g., "Implement user authentication")
2. Position view so parent + 2 subtasks are clearly visible
3. Capture focused view
4. Save as `docs/images/tui-subtasks-hierarchy.png`

#### Screenshot 3: Task Form

1. Press `n` to create new task or `e` to edit existing task
2. Fill in all fields with sample data:
   - Title: "Example task"
   - Description: "Line 1\nLine 2\nLine 3" (use Enter for newlines)
   - Priority: 2
   - Progress: 50
   - Tags: "example,demo"
3. Ensure multiline description and line numbers are visible
4. Capture screenshot showing full form
5. Save as `docs/images/tui-task-form.png`

**Optimize images** (optional):
```bash
optipng docs/images/*.png
# or
pngcrush -ow docs/images/*.png
```

**Expected result**: 3 PNG files in `docs/images/`, each < 500KB

---

### Step 5: Capture CLI Output Examples (20 minutes)

Run CLI commands and copy output for README:

```bash
# Example 1: List command
./ontop --db-path ./sample-tasks.db list > /tmp/cli-list-output.txt

# Example 2: Show command (use actual task ID from list output)
TASK_ID=$(./ontop --db-path ./sample-tasks.db list | head -n 5 | grep "Implement user" | awk '{print $1}' | tr -d '[]')
./ontop --db-path ./sample-tasks.db show $TASK_ID > /tmp/cli-show-output.txt

# Example 3: Add parent task
./ontop --db-path ./sample-tasks.db add "New feature" --priority 2 --description "Example" > /tmp/cli-add-output.txt

# Example 4: Add subtask
./ontop --db-path ./sample-tasks.db add "Subtask" --parent $TASK_ID --priority 3 > /tmp/cli-add-subtask-output.txt
```

**Review outputs**:
```bash
cat /tmp/cli-list-output.txt
cat /tmp/cli-show-output.txt
cat /tmp/cli-add-output.txt
cat /tmp/cli-add-subtask-output.txt
```

**Expected result**: Clean, formatted CLI output showing hierarchical structure

---

### Step 6: Update README.md (60 minutes)

Restructure README following the researched format:

#### 6.1 Add Hero Screenshot

After the project description, add:

```markdown
## Overview

![OnTop TUI showing kanban board with three columns containing parent tasks and indented subtasks](docs/images/tui-kanban-board.png)
```

#### 6.2 Enhance Features Section

For each feature, add visual proof:

```markdown
### Hierarchical Subtasks

Create and manage subtasks with visual indentation in both TUI and CLI modes.

![Parent task with subtasks indented below showing hierarchical structure](docs/images/tui-subtasks-hierarchy.png)

### Rich Task Forms

Edit tasks with multiline descriptions, scrollable textarea, and full keyboard navigation.

![Task creation form with multiline description textarea](docs/images/tui-task-form.png)
```

#### 6.3 Add CLI Examples with Output

Replace existing CLI examples with command + output pairs:

````markdown
### CLI Mode Examples

```bash
$ ./ontop list
```

```
Inbox (2 tasks)
  [abc123] P1: Implement user authentication (0%)
    - [def456] P2: Research OAuth2 libraries (100%)
    - [ghi789] P2: Implement GitHub OAuth provider (0%)
  [jkl012] P3: Update API documentation (25%)

In Progress (2 tasks)
  [mno345] P2: Fix memory leak in task sync (75%)
  [pqr678] P4: Add fuzzy search to task titles (50%)
    - [stu901] P4: Evaluate fuzzy search algorithms (100%)
    - [vwx234] P4: Add search input to TUI (25%)

Done (1 task)
  [yza567] P2: Migrate to WAL mode for SQLite (100%)
```
````

#### 6.4 Add Quick Start Section

```markdown
## Quick Start

1. **Install OnTop**:
   ```bash
   make install
   ```

2. **Launch TUI**:
   ```bash
   ontop
   ```

3. **Create your first task** (press `n`):
   ![Task form example](docs/images/tui-task-form.png)

4. **Navigate with vim keys**: `j/k` for up/down, `h/l` for columns
```

#### 6.5 Verify All Images Have Alt Text

Search for all image references:
```bash
grep -n "!\[" README.md
```

Ensure each has descriptive alt text (not just "screenshot").

**Expected result**: README.md enhanced with 3+ screenshots and CLI examples

---

### Step 7: Verification (15 minutes)

#### Verify Locally

1. **Preview README** on GitHub (push to your branch and view on web)
2. **Check image loading**: All screenshots display correctly
3. **Check image size**: Total size < 2MB
   ```bash
   du -sh docs/images/
   ```
4. **Verify accessibility**: All images have descriptive alt text
5. **Test reproducibility**: Delete sample database and regenerate
   ```bash
   rm -f sample-tasks.db
   ./scripts/generate-sample-data.sh
   ```

#### Quality Checklist

- [ ] All 3 screenshots display in README
- [ ] Screenshots are clear and high-resolution
- [ ] Hierarchical indentation visible in screenshots
- [ ] CLI examples show actual formatted output
- [ ] All images have descriptive alt text
- [ ] Total image directory size < 2MB
- [ ] Sample data script is executable and documented
- [ ] Sample data is reproducible (script can be re-run)
- [ ] README structure follows researched best practices
- [ ] No broken image links

---

## Common Issues and Solutions

### Issue: Screenshot file size too large

**Solution**: Reduce terminal size slightly or use PNG optimization:
```bash
optipng -o7 docs/images/tui-kanban-board.png
```

### Issue: Subtasks not showing in hierarchy

**Solution**: Verify parent-child relationship:
```bash
./ontop --db-path ./sample-tasks.db show <parent-id>
```

### Issue: CLI output not formatting correctly

**Solution**: Ensure terminal width is at least 120 characters when capturing output

### Issue: Images not displaying on GitHub

**Solution**: Verify relative paths and that images are committed:
```bash
git add docs/images/*.png
git status
```

---

## Completion Checklist

Mark as complete when:

- [ ] Sample data generation script exists and works
- [ ] 3+ screenshots captured and optimized
- [ ] README restructured with visual-first approach
- [ ] CLI examples include actual command output
- [ ] All images have descriptive alt text
- [ ] Total image size < 2MB verified
- [ ] README renders correctly on GitHub
- [ ] Sample data is reproducible
- [ ] All files committed to git

---

## Next Steps

After completing this quick start:

1. Review the updated README on GitHub
2. Gather feedback on screenshot quality and clarity
3. Consider adding more screenshots for additional features
4. Update `.github/CONTRIBUTING.md` if needed to reference screenshot update process
5. Run `/speckit.tasks` command to generate detailed implementation tasks

---

## Time Estimate Breakdown

- Directory creation: 5 minutes
- Sample data script: 30 minutes
- Generate and test sample data: 10 minutes
- Capture screenshots: 30 minutes
- CLI output examples: 20 minutes
- README updates: 60 minutes
- Verification: 15 minutes

**Total**: ~2-3 hours for complete implementation
