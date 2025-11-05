#!/bin/bash
# Purpose: Generate sample tasks for README screenshots and CLI examples
# Usage: ./scripts/generate-sample-data.sh [DB_PATH]
#
# This script creates a comprehensive sample dataset demonstrating:
# - Parent-child task relationships (hierarchical subtasks)
# - Tasks distributed across all three columns (Inbox, In Progress, Done)
# - Various priority levels (P1-P5)
# - Progress percentages (0%-100%)
# - Tags for categorization
# - Multiline descriptions

set -e  # Exit on any error

# Configuration
DB_PATH="${1:-./sample-tasks.db}"
BINARY="./ontop"

# Colors for output
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

echo -e "${YELLOW}=== OnTop Sample Data Generator ===${NC}"
echo "Target database: $DB_PATH"
echo ""

# Verify ontop binary exists
if [ ! -f "$BINARY" ]; then
    echo "Error: OnTop binary not found at $BINARY"
    echo "Please run 'make build' first"
    exit 1
fi

# Remove existing sample database if it exists
if [ -f "$DB_PATH" ]; then
    echo -e "${BLUE}Removing existing sample database...${NC}"
    rm -f "$DB_PATH"
fi

echo "Creating sample tasks..."
echo ""

# Helper function to extract task ID from ontop output
extract_id() {
    # ontop outputs: "Created task {ID}: {description}"
    # ID is a 26-character ULID (uppercase alphanumeric)
    echo "$1" | grep -oE 'Created task [A-Z0-9]{26}' | awk '{print $3}' | head -1
}

# ============================================================================
# PARENT TASKS (5 total across columns)
# ============================================================================

echo "Creating parent tasks..."

# Task 1: High-priority feature (Inbox) - Will have 2 subtasks
echo "  [1/5] User Authentication Feature (P1, Inbox)..."
TASK1_OUTPUT=$($BINARY --db-path "$DB_PATH" add \
  -title "Implement user authentication" \
  -priority 1 \
  -description "Add OAuth2 and session management
Support GitHub and Google providers
Include remember-me functionality" \
  -tags "feature,security")
TASK1_ID=$(extract_id "$TASK1_OUTPUT")
echo "      Created: $TASK1_ID"

# Task 2: Medium-priority docs (Inbox)
echo "  [2/5] Documentation Update (P3, Inbox)..."
TASK2_OUTPUT=$($BINARY --db-path "$DB_PATH" add \
  -title "Update API documentation" \
  -priority 3 \
  -progress 25 \
  -description "Review and update all API endpoint documentation to reflect recent changes" \
  -tags "documentation")
TASK2_ID=$(extract_id "$TASK2_OUTPUT")
echo "      Created: $TASK2_ID"

# Task 3: High-priority bug (In Progress)
echo "  [3/5] Memory Leak Fix (P2, In Progress)..."
TASK3_OUTPUT=$($BINARY --db-path "$DB_PATH" add \
  -title "Fix memory leak in task sync" \
  -priority 2 \
  -progress 75 \
  -description "Investigate and resolve memory growth issue in background sync process" \
  -tags "bug,performance")
TASK3_ID=$(extract_id "$TASK3_OUTPUT")
$BINARY --db-path "$DB_PATH" move "$TASK3_ID" in_progress > /dev/null
echo "      Created: $TASK3_ID (moved to In Progress)"

# Task 4: Low-priority feature (In Progress) - Will have 2 subtasks
echo "  [4/5] Search Enhancement (P4, In Progress)..."
TASK4_OUTPUT=$($BINARY --db-path "$DB_PATH" add \
  -title "Add fuzzy search to task titles" \
  -priority 4 \
  -progress 50 \
  -description "Implement fuzzy matching for task search
Consider using fzf-like algorithm
Test with large task datasets" \
  -tags "feature,enhancement")
TASK4_ID=$(extract_id "$TASK4_OUTPUT")
$BINARY --db-path "$DB_PATH" move "$TASK4_ID" in_progress > /dev/null
echo "      Created: $TASK4_ID (moved to In Progress)"

# Task 5: Completed infrastructure (Done)
echo "  [5/5] Database Migration (P2, Done)..."
TASK5_OUTPUT=$($BINARY --db-path "$DB_PATH" add \
  -title "Migrate to WAL mode for SQLite" \
  -priority 2 \
  -progress 100 \
  -description "Enable Write-Ahead Logging for better concurrent access" \
  -tags "infrastructure,database")
TASK5_ID=$(extract_id "$TASK5_OUTPUT")
$BINARY --db-path "$DB_PATH" move "$TASK5_ID" done > /dev/null
echo "      Created: $TASK5_ID (moved to Done)"

echo ""

# ============================================================================
# SUBTASKS (4 total demonstrating hierarchy)
# ============================================================================

echo "Creating subtasks..."

# Subtask 1-1: Research for Task 1 (completed, in Done)
echo "  [1/4] OAuth Library Research (child of $TASK1_ID)..."
SUBTASK1_1_OUTPUT=$($BINARY --db-path "$DB_PATH" add \
  -title "Research OAuth2 libraries" \
  -parent "$TASK1_ID" \
  -priority 2 \
  -progress 100 \
  -description "Evaluate golang.org/x/oauth2 and alternatives" \
  -tags "research")
SUBTASK1_1=$(extract_id "$SUBTASK1_1_OUTPUT")
$BINARY --db-path "$DB_PATH" move "$SUBTASK1_1" done > /dev/null
echo "      Created: $SUBTASK1_1 (moved to Done)"

# Subtask 1-2: Implementation for Task 1 (not started, in Inbox)
echo "  [2/4] GitHub Provider Implementation (child of $TASK1_ID)..."
SUBTASK1_2_OUTPUT=$($BINARY --db-path "$DB_PATH" add \
  -title "Implement GitHub OAuth provider" \
  -parent "$TASK1_ID" \
  -priority 2 \
  -description "Create provider for GitHub authentication flow" \
  -tags "development")
SUBTASK1_2=$(extract_id "$SUBTASK1_2_OUTPUT")
echo "      Created: $SUBTASK1_2 (stays in Inbox)"

# Subtask 4-1: Research for Task 4 (completed, in Done)
echo "  [3/4] Search Algorithm Evaluation (child of $TASK4_ID)..."
SUBTASK4_1_OUTPUT=$($BINARY --db-path "$DB_PATH" add \
  -title "Evaluate fuzzy search algorithms" \
  -parent "$TASK4_ID" \
  -priority 4 \
  -progress 100 \
  -description "Compare Levenshtein distance vs. trigram matching" \
  -tags "research")
SUBTASK4_1=$(extract_id "$SUBTASK4_1_OUTPUT")
$BINARY --db-path "$DB_PATH" move "$SUBTASK4_1" done > /dev/null
echo "      Created: $SUBTASK4_1 (moved to Done)"

# Subtask 4-2: Development for Task 4 (in progress)
echo "  [4/4] Search UI Integration (child of $TASK4_ID)..."
SUBTASK4_2_OUTPUT=$($BINARY --db-path "$DB_PATH" add \
  -title "Add search input to TUI" \
  -parent "$TASK4_ID" \
  -priority 4 \
  -progress 25 \
  -description "Integrate search field into kanban view header" \
  -tags "development,ui")
SUBTASK4_2=$(extract_id "$SUBTASK4_2_OUTPUT")
$BINARY --db-path "$DB_PATH" move "$SUBTASK4_2" in_progress > /dev/null
echo "      Created: $SUBTASK4_2 (moved to In Progress)"

echo ""

# ============================================================================
# FINALIZE
# ============================================================================

# Checkpoint WAL to consolidate all data into main database file
echo -e "${BLUE}Consolidating database (checkpoint WAL)...${NC}"
sqlite3 "$DB_PATH" "PRAGMA wal_checkpoint(TRUNCATE);" > /dev/null 2>&1

echo ""
echo -e "${GREEN}✓ Sample data created successfully!${NC}"
echo ""
echo "Summary:"
echo "  - 5 parent tasks (2 Inbox, 2 In Progress, 1 Done)"
echo "  - 4 subtasks demonstrating hierarchy"
echo "  - Priority range: P1-P4"
echo "  - Progress range: 0%-100%"
echo "  - Tags: feature, security, documentation, bug, performance, enhancement,"
echo "          infrastructure, database, research, development, ui"
echo ""
echo "Verify with:"
echo "  $BINARY --db-path \"$DB_PATH\" list"
echo ""
echo "Launch TUI with:"
echo "  $BINARY --db-path \"$DB_PATH\""
echo ""

exit 0
