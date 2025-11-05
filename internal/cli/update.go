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

// UpdateCommand implements the 'ontop update' command
func UpdateCommand(db *sql.DB, args []string) {
	fs := flag.NewFlagSet("update", flag.ExitOnError)
	description := fs.String("description", "", "Update task description")
	priority := fs.Int("priority", -1, "Update task priority (1-5)")
	column := fs.String("column", "", "Update task column (inbox, in_progress, done)")
	progress := fs.Int("progress", -1, "Update task progress (0-100)")
	tagsStr := fs.String("tags", "", "Update task tags (comma-separated)")
	addTagsStr := fs.String("add-tags", "", "Add tags (comma-separated)")
	removeTagsStr := fs.String("remove-tags", "", "Remove tags (comma-separated)")
	clearTags := fs.Bool("clear-tags", false, "Clear all tags")

	fs.Usage = func() {
		fmt.Fprintf(os.Stderr, `Usage: ontop update <task-id> [options]

Update various attributes of a task.

OPTIONS:
    -description string   Update task description
    -priority int         Update task priority (1-5)
    -column string        Update task column (inbox, in_progress, done)
    -progress int         Update task progress (0-100)
    -tags string          Replace all tags (comma-separated)
    -add-tags string      Add tags (comma-separated)
    -remove-tags string   Remove tags (comma-separated)
    -clear-tags           Clear all tags

EXAMPLES:
    ontop update 20251104-143000-00001 -description "New description"
    ontop update 20251104-143000-00001 -priority 1
    ontop update 20251104-143000-00001 -progress 50
    ontop update 20251104-143000-00001 -column in_progress
    ontop update 20251104-143000-00001 -add-tags "urgent,bug"
    ontop update 20251104-143000-00001 -remove-tags "old-tag"
    ontop update 20251104-143000-00001 -priority 2 -progress 75
`)
	}

	// Get task ID first (must be first argument)
	if len(args) < 1 {
		fmt.Fprintf(os.Stderr, "Error: Task ID is required\n")
		fs.Usage()
		os.Exit(2)
	}

	taskID := args[0]

	// Parse remaining args as flags
	if err := fs.Parse(args[1:]); err != nil {
		os.Exit(2)
	}

	// Get task
	task, err := storage.GetTask(db, taskID)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: Task not found: %v\n", err)
		os.Exit(1)
	}

	// Track what was updated
	updates := []string{}

	// Update description
	if *description != "" {
		task.Description = *description
		updates = append(updates, "description")
	}

	// Update priority
	if *priority > 0 {
		if *priority < 1 || *priority > 5 {
			fmt.Fprintf(os.Stderr, "Error: Priority must be between 1 and 5\n")
			os.Exit(2)
		}
		task.Priority = *priority
		updates = append(updates, fmt.Sprintf("priority to P%d", *priority))
	}

	// Update column
	if *column != "" {
		if !models.IsValidColumn(*column) {
			fmt.Fprintf(os.Stderr, "Error: Invalid column '%s'. Valid columns: %s\n",
				*column, strings.Join(models.ValidColumns(), ", "))
			os.Exit(2)
		}
		oldColumn := task.Column
		task.Column = *column

		// If moving to done, set completed_at
		if *column == models.ColumnDone && task.CompletedAt == nil {
			now := time.Now()
			task.CompletedAt = &now
		}

		// If moving from done to another column, clear completed_at
		if oldColumn == models.ColumnDone && *column != models.ColumnDone {
			task.CompletedAt = nil
		}

		updates = append(updates, fmt.Sprintf("column to %s", *column))
	}

	// Update progress
	if *progress >= 0 {
		if *progress < 0 || *progress > 100 {
			fmt.Fprintf(os.Stderr, "Error: Progress must be between 0 and 100\n")
			os.Exit(2)
		}
		task.Progress = *progress
		updates = append(updates, fmt.Sprintf("progress to %d%%", *progress))

		// Auto-move to done if progress is 100
		if *progress == 100 && task.Column != models.ColumnDone {
			task.Column = models.ColumnDone
			now := time.Now()
			task.CompletedAt = &now
			updates = append(updates, "auto-moved to Done")
		}
	}

	// Handle tags
	if *clearTags {
		task.Tags = []string{}
		updates = append(updates, "cleared all tags")
	} else if *tagsStr != "" {
		// Replace all tags
		tags := strings.Split(*tagsStr, ",")
		for i, tag := range tags {
			tags[i] = strings.TrimSpace(tag)
		}
		task.Tags = tags
		updates = append(updates, "updated tags")
	} else {
		// Add tags
		if *addTagsStr != "" {
			newTags := strings.Split(*addTagsStr, ",")
			for _, tag := range newTags {
				tag = strings.TrimSpace(tag)
				// Check if tag already exists
				exists := false
				for _, existing := range task.Tags {
					if existing == tag {
						exists = true
						break
					}
				}
				if !exists {
					task.Tags = append(task.Tags, tag)
				}
			}
			updates = append(updates, fmt.Sprintf("added tags: %s", *addTagsStr))
		}

		// Remove tags
		if *removeTagsStr != "" {
			removeTags := strings.Split(*removeTagsStr, ",")
			for i, tag := range removeTags {
				removeTags[i] = strings.TrimSpace(tag)
			}

			// Filter out removed tags
			var newTags []string
			for _, tag := range task.Tags {
				keep := true
				for _, removeTag := range removeTags {
					if tag == removeTag {
						keep = false
						break
					}
				}
				if keep {
					newTags = append(newTags, tag)
				}
			}
			task.Tags = newTags
			updates = append(updates, fmt.Sprintf("removed tags: %s", *removeTagsStr))
		}
	}

	// Check if anything was updated
	if len(updates) == 0 {
		fmt.Fprintf(os.Stderr, "Error: No updates specified. Use --help to see available options.\n")
		os.Exit(2)
	}

	// Update timestamp
	task.UpdatedAt = time.Now()

	// Save task
	if err := storage.UpdateTask(db, task); err != nil {
		fmt.Fprintf(os.Stderr, "Error: Failed to update task: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Updated task %s: %s\n", taskID, strings.Join(updates, ", "))

	os.Exit(0)
}
