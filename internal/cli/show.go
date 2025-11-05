package cli

import (
	"database/sql"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/lucasefe/ontop/internal/models"
	"github.com/lucasefe/ontop/internal/storage"
)

// ShowCommand implements the 'ontop show' command
func ShowCommand(db *sql.DB, args []string) {
	fs := flag.NewFlagSet("show", flag.ExitOnError)
	jsonOutput := fs.Bool("json", false, "Output result as JSON")

	fs.Usage = func() {
		fmt.Fprintf(os.Stderr, `Usage: ontop show <task-id> [options]

Show detailed information about a specific task, including all attributes and subtasks.

OPTIONS:
    -json    Output result as JSON

EXAMPLES:
    ontop show 20251104-143000-00001
    ontop show 20251104-143000-00001 -json
`)
	}

	if err := fs.Parse(args); err != nil {
		os.Exit(2)
	}

	// Get task ID from remaining args
	if fs.NArg() < 1 {
		fmt.Fprintf(os.Stderr, "Error: Task ID is required\n")
		fs.Usage()
		os.Exit(2)
	}

	taskID := fs.Arg(0)

	// Get task
	task, err := storage.GetTask(db, taskID)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: Task not found: %v\n", err)
		os.Exit(1)
	}

	// Get subtasks
	filters := map[string]interface{}{
		"archived": false,
	}
	allTasks, err := storage.ListTasks(db, filters)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: Failed to list subtasks: %v\n", err)
		os.Exit(1)
	}

	var subtasks []*models.Task
	for _, t := range allTasks {
		if t.ParentID != nil && *t.ParentID == taskID {
			subtasks = append(subtasks, t)
		}
	}

	// Output result
	if *jsonOutput {
		result := map[string]interface{}{
			"task":     task,
			"subtasks": subtasks,
		}
		output, err := json.MarshalIndent(result, "", "  ")
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: Failed to marshal JSON: %v\n", err)
			os.Exit(1)
		}
		fmt.Println(string(output))
	} else {
		// Display task details
		fmt.Println(strings.Repeat("=", 70))
		fmt.Printf("TASK: %s\n", task.ID)
		fmt.Println(strings.Repeat("=", 70))
		fmt.Printf("Description:  %s\n", task.Description)
		fmt.Printf("Priority:     P%d (1=highest, 5=lowest)\n", task.Priority)
		fmt.Printf("Column:       %s\n", formatColumnDisplay(task.Column))
		fmt.Printf("Progress:     %d%%\n", task.Progress)
		fmt.Printf("Archived:     %v\n", task.Archived)

		if len(task.Tags) > 0 {
			fmt.Printf("Tags:         %s\n", strings.Join(task.Tags, ", "))
		} else {
			fmt.Printf("Tags:         (none)\n")
		}

		if task.ParentID != nil {
			fmt.Printf("Parent:       %s\n", *task.ParentID)
		}

		fmt.Printf("\nCreated:      %s\n", task.CreatedAt.Format("2006-01-02 15:04:05"))
		fmt.Printf("Updated:      %s\n", task.UpdatedAt.Format("2006-01-02 15:04:05"))

		if task.CompletedAt != nil {
			fmt.Printf("Completed:    %s\n", task.CompletedAt.Format("2006-01-02 15:04:05"))
		}

		// Display subtasks if any
		if len(subtasks) > 0 {
			fmt.Printf("\nSUBTASKS (%d):\n", len(subtasks))
			fmt.Println(strings.Repeat("-", 70))
			for _, st := range subtasks {
				progressStr := ""
				if st.Progress > 0 {
					progressStr = fmt.Sprintf(" (%d%%)", st.Progress)
				}
				fmt.Printf("  - [%s] P%d | %s | %s%s\n",
					st.ID,
					st.Priority,
					formatColumnDisplay(st.Column),
					st.Description,
					progressStr,
				)
			}
		}

		fmt.Println(strings.Repeat("=", 70))
	}

	os.Exit(0)
}
