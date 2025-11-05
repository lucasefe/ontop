package cli

import (
	"database/sql"
	"flag"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/lucasefe/ontop/internal/models"
	"github.com/lucasefe/ontop/internal/storage"
)

// MoveCommand implements the 'ontop move' command
func MoveCommand(db *sql.DB, args []string) {
	fs := flag.NewFlagSet("move", flag.ExitOnError)

	fs.Usage = func() {
		fmt.Fprintf(os.Stderr, `Usage: ontop move <task-id> <column>

Move a task to a different column.

COLUMNS:
    inbox        - Inbox (new tasks)
    in_progress  - In Progress (active work)
    done         - Done (completed tasks)

EXAMPLES:
    ontop move 20251104-143000-00001 in_progress
    ontop move 20251104-143000-00001 done
`)
	}

	if err := fs.Parse(args); err != nil {
		os.Exit(2)
	}

	// Get task ID and column from remaining args
	if fs.NArg() < 2 {
		fmt.Fprintf(os.Stderr, "Error: Task ID and column are required\n")
		fs.Usage()
		os.Exit(2)
	}

	taskID := fs.Arg(0)
	column := fs.Arg(1)

	// Validate column
	if !models.IsValidColumn(column) {
		fmt.Fprintf(os.Stderr, "Error: Invalid column '%s'. Valid columns: %s\n",
			column, strings.Join(models.ValidColumns(), ", "))
		os.Exit(2)
	}

	// Get task
	task, err := storage.GetTask(db, taskID)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: Task not found: %v\n", err)
		os.Exit(1)
	}

	oldColumn := task.Column
	task.Column = column
	task.UpdatedAt = time.Now()

	// If moving to done, set completed_at
	if column == models.ColumnDone && task.CompletedAt == nil {
		now := time.Now()
		task.CompletedAt = &now
	}

	// If moving from done to another column, clear completed_at
	if oldColumn == models.ColumnDone && column != models.ColumnDone {
		task.CompletedAt = nil
	}

	// Update task
	if err := storage.UpdateTask(db, task); err != nil {
		fmt.Fprintf(os.Stderr, "Error: Failed to move task: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Moved task %s: %s → %s\n",
		taskID,
		formatColumnName(oldColumn),
		formatColumnName(column))

	os.Exit(0)
}

func formatColumnName(column string) string {
	switch column {
	case "inbox":
		return "Inbox"
	case "in_progress":
		return "In Progress"
	case "done":
		return "Done"
	default:
		return column
	}
}
