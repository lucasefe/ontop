package service

import (
	"sort"
	"strings"

	"github.com/lucasefe/ontop/internal/models"
)

// SortMode represents how tasks are sorted
type SortMode int

const (
	SortByPriority SortMode = iota
	SortByDescription
	SortByCreated
	SortByUpdated
)

// HierarchicalTask wraps a Task with display context for hierarchical rendering
type HierarchicalTask struct {
	Task        *models.Task // The actual task data
	IsSubtask   bool         // True if ParentID != nil
	Indentation int          // Number of spaces for indentation (0 or 2)
	ParentTitle string       // Title of parent task (empty if not subtask)
}

// BuildFlatHierarchy converts flat task list into hierarchical display order.
// Parents are sorted by sortMode, subtasks immediately follow their parent.
// Orphaned subtasks (parent not in list) are treated as top-level tasks.
func BuildFlatHierarchy(tasks []*models.Task, sortMode SortMode) []*HierarchicalTask {
	if len(tasks) == 0 {
		return []*HierarchicalTask{}
	}

	// Step 1: Separate parents and subtasks
	var parents []*models.Task
	var subtasks []*models.Task
	for _, task := range tasks {
		if task.ParentID == nil {
			parents = append(parents, task)
		} else {
			subtasks = append(subtasks, task)
		}
	}

	// Step 2: Sort parents by sortMode
	sortTasks(parents, sortMode)

	// Step 3: Build result by interleaving parents and their subtasks
	result := []*HierarchicalTask{}
	for _, parent := range parents {
		// Add parent
		result = append(result, &HierarchicalTask{
			Task:        parent,
			IsSubtask:   false,
			Indentation: 0,
		})

		// Find and add subtasks for this parent
		parentSubtasks := findSubtasksForParent(subtasks, parent.ID)
		sortTasks(parentSubtasks, sortMode)
		for _, st := range parentSubtasks {
			result = append(result, &HierarchicalTask{
				Task:        st,
				IsSubtask:   true,
				Indentation: 2,
				ParentTitle: parent.Title,
			})
		}
	}

	// Step 4: Handle orphaned subtasks (parent not in list)
	for _, st := range subtasks {
		if !hasParentInList(st, parents) {
			result = append(result, &HierarchicalTask{
				Task:        st,
				IsSubtask:   false, // treat as top-level
				Indentation: 0,
			})
		}
	}

	return result
}

// GetSubtasksForParent returns all subtasks for a specific parent task.
// Used in detail view to display "Subtasks (N)" section.
func GetSubtasksForParent(tasks []*models.Task, parentID string) []*models.Task {
	var subtasks []*models.Task
	for _, task := range tasks {
		if task.ParentID != nil && *task.ParentID == parentID {
			subtasks = append(subtasks, task)
		}
	}
	return subtasks
}

// IsOrphanedSubtask checks if a subtask's parent no longer exists
func IsOrphanedSubtask(task *models.Task, allTasks []*models.Task) bool {
	if task.ParentID == nil {
		return false
	}

	for _, t := range allTasks {
		if t.ID == *task.ParentID {
			return false
		}
	}
	return true
}

// FormatWithIndentation formats task title with appropriate indentation for CLI output
func FormatWithIndentation(ht *HierarchicalTask) string {
	indent := strings.Repeat(" ", ht.Indentation)
	prefix := ""
	if ht.IsSubtask {
		prefix = "- "
	}
	return indent + prefix + ht.Task.Title
}

// Helper functions

func sortTasks(tasks []*models.Task, sortMode SortMode) {
	switch sortMode {
	case SortByPriority:
		sort.Slice(tasks, func(i, j int) bool {
			return tasks[i].Priority < tasks[j].Priority
		})
	case SortByDescription:
		sort.Slice(tasks, func(i, j int) bool {
			return tasks[i].Title < tasks[j].Title
		})
	case SortByCreated:
		sort.Slice(tasks, func(i, j int) bool {
			return tasks[i].CreatedAt.Before(tasks[j].CreatedAt)
		})
	case SortByUpdated:
		sort.Slice(tasks, func(i, j int) bool {
			return tasks[i].UpdatedAt.After(tasks[j].UpdatedAt)
		})
	}
}

func findSubtasksForParent(subtasks []*models.Task, parentID string) []*models.Task {
	var result []*models.Task
	for _, st := range subtasks {
		if st.ParentID != nil && *st.ParentID == parentID {
			result = append(result, st)
		}
	}
	return result
}

func hasParentInList(subtask *models.Task, parents []*models.Task) bool {
	if subtask.ParentID == nil {
		return false
	}
	for _, p := range parents {
		if p.ID == *subtask.ParentID {
			return true
		}
	}
	return false
}
