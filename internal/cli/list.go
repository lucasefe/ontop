package cli

import (
	"database/sql"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/lucasefe/ontop/internal/storage"
)

// ListCommand implements the 'ontop list' command
func ListCommand(db *sql.DB, args []string) {
	fs := flag.NewFlagSet("list", flag.ExitOnError)
	priority := fs.Int("priority", -1, "Filter by priority (1-5)")
	column := fs.String("column", "", "Filter by column (inbox, in_progress, done)")
	tag := fs.String("tag", "", "Filter by tag")
	archived := fs.Bool("archived", false, "Show archived tasks")
	jsonOutput := fs.Bool("json", false, "Output result as JSON")

	fs.Usage = func() {
		fmt.Fprintf(os.Stderr, `Usage: ontop list [options]

List tasks with optional filters.

OPTIONS:
    -priority int     Filter by priority (1-5)
    -column string    Filter by column (inbox, in_progress, done)
    -tag string       Filter by tag
    -archived         Show archived tasks instead of active tasks
    -json             Output result as JSON

EXAMPLES:
    ontop list
    ontop list -priority 1
    ontop list -column in_progress
    ontop list -tag urgent
    ontop list -archived
`)
	}

	if err := fs.Parse(args); err != nil {
		os.Exit(2)
	}

	// Build filters
	filters := make(map[string]interface{})

	if *archived {
		filters["archived"] = true
	} else {
		filters["archived"] = false
	}

	if *priority > 0 {
		if *priority < 1 || *priority > 5 {
			fmt.Fprintf(os.Stderr, "Error: Priority must be between 1 and 5\n")
			os.Exit(2)
		}
		filters["priority"] = *priority
	}

	if *column != "" {
		filters["column"] = *column
	}

	if *tag != "" {
		filters["tag"] = *tag
	}

	// Query tasks
	tasks, err := storage.ListTasks(db, filters)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: Failed to list tasks: %v\n", err)
		os.Exit(1)
	}

	// Output result
	if *jsonOutput {
		output, err := json.MarshalIndent(tasks, "", "  ")
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: Failed to marshal JSON: %v\n", err)
			os.Exit(1)
		}
		fmt.Println(string(output))
	} else {
		if len(tasks) == 0 {
			fmt.Println("No tasks found.")
			os.Exit(0)
		}

		fmt.Printf("\nTotal: %d tasks\n\n", len(tasks))

		for _, task := range tasks {
			priorityStr := fmt.Sprintf("P%d", task.Priority)
			tagsStr := ""
			if len(task.Tags) > 0 {
				tagsStr = fmt.Sprintf(" [%s]", strings.Join(task.Tags, ", "))
			}
			progressStr := ""
			if task.Progress > 0 {
				progressStr = fmt.Sprintf(" (%d%%)", task.Progress)
			}

			fmt.Printf("[%s] %s | %s | %s%s%s\n",
				task.ID,
				priorityStr,
				formatColumnDisplay(task.Column),
				task.Description,
				tagsStr,
				progressStr,
			)
		}
	}

	os.Exit(0)
}

func formatColumnDisplay(column string) string {
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
