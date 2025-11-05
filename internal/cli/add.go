package cli

import (
	"database/sql"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/lucasefe/ontop/internal/models"
	"github.com/lucasefe/ontop/internal/service"
	"github.com/lucasefe/ontop/internal/storage"
)

// AddCommand implements the 'ontop add' command
func AddCommand(db *sql.DB, args []string) {
	fs := flag.NewFlagSet("add", flag.ExitOnError)
	priority := fs.Int("priority", 3, "Task priority (1-5, where 1 is highest)")
	tagsStr := fs.String("tags", "", "Comma-separated list of tags")
	parentID := fs.String("parent", "", "Parent task ID for subtasks")
	column := fs.String("column", models.ColumnInbox, "Column to place task in (inbox, in_progress, done)")
	progress := fs.Int("progress", 0, "Initial progress percentage (0-100)")
	jsonOutput := fs.Bool("json", false, "Output result as JSON")

	fs.Usage = func() {
		fmt.Fprintf(os.Stderr, `Usage: ontop add <description> [options]

Create a new task with the specified description.

OPTIONS:
    -priority int      Task priority (1-5, where 1 is highest) (default: 3)
    -tags string       Comma-separated list of tags
    -parent string     Parent task ID for subtasks
    -column string     Column to place task in: inbox, in_progress, done (default: inbox)
    -progress int      Initial progress percentage (0-100) (default: 0)
    -json              Output result as JSON

EXAMPLES:
    ontop add "Implement user authentication"
    ontop add "Fix login bug" -priority 1 -tags "bug,urgent"
    ontop add "Write tests" -parent 20251104-143000-00001
`)
	}

	if err := fs.Parse(args); err != nil {
		os.Exit(2)
	}

	// Get description from remaining args
	if fs.NArg() < 1 {
		fmt.Fprintf(os.Stderr, "Error: Task description is required\n")
		fs.Usage()
		os.Exit(2)
	}

	description := strings.Join(fs.Args(), " ")

	// Validate inputs
	if *priority < 1 || *priority > 5 {
		fmt.Fprintf(os.Stderr, "Error: Priority must be between 1 and 5\n")
		os.Exit(2)
	}

	if !models.IsValidColumn(*column) {
		fmt.Fprintf(os.Stderr, "Error: Invalid column '%s'. Valid columns: %s\n", *column, strings.Join(models.ValidColumns(), ", "))
		os.Exit(2)
	}

	if *progress < 0 || *progress > 100 {
		fmt.Fprintf(os.Stderr, "Error: Progress must be between 0 and 100\n")
		os.Exit(2)
	}

	// Parse tags
	var tags []string
	if *tagsStr != "" {
		tags = strings.Split(*tagsStr, ",")
		for i, tag := range tags {
			tags[i] = strings.TrimSpace(tag)
		}
	} else {
		tags = []string{}
	}

	// Validate parent exists if provided
	var parentIDPtr *string
	if *parentID != "" {
		_, err := storage.GetTask(db, *parentID)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: Parent task '%s' not found: %v\n", *parentID, err)
			os.Exit(2)
		}
		parentIDPtr = parentID
	}

	// Create task
	now := time.Now()
	task := &models.Task{
		ID:          service.GenerateID(),
		Description: description,
		Priority:    *priority,
		Column:      *column,
		Progress:    *progress,
		ParentID:    parentIDPtr,
		Archived:    false,
		Tags:        tags,
		CreatedAt:   now,
		UpdatedAt:   now,
		CompletedAt: nil,
		DeletedAt:   nil,
	}

	if err := storage.CreateTask(db, task); err != nil {
		fmt.Fprintf(os.Stderr, "Error: Failed to create task: %v\n", err)
		os.Exit(1)
	}

	// Output result
	if *jsonOutput {
		output, err := json.MarshalIndent(task, "", "  ")
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: Failed to marshal JSON: %v\n", err)
			os.Exit(1)
		}
		fmt.Println(string(output))
	} else {
		fmt.Printf("Created task %s: %s\n", task.ID, task.Description)
	}

	os.Exit(0)
}
