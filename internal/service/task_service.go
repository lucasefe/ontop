package service

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"github.com/lucasefe/ontop/internal/models"
	"github.com/lucasefe/ontop/internal/storage"
)

// CalculateProgress computes the progress percentage for a task based on its subtasks
// If the task has no subtasks, returns its own progress value
// If it has subtasks, returns the average progress of all subtasks
func CalculateProgress(db *sql.DB, taskID string) (int, error) {
	// Query subtasks - we need to add parent_id filtering
	query := `
		SELECT id, description, priority, column, progress, parent_id,
		       archived, tags, created_at, updated_at, completed_at, deleted_at
		FROM tasks
		WHERE deleted_at IS NULL AND parent_id = ?
	`

	rows, err := db.Query(query, taskID)
	if err != nil {
		return 0, fmt.Errorf("failed to query subtasks: %w", err)
	}
	defer func() {
		if err := rows.Close(); err != nil {
			fmt.Printf("Warning: Failed to close rows: %v\n", err)
		}
	}()

	var subtasks []*models.Task
	for rows.Next() {
		task, err := scanTask(rows)
		if err != nil {
			return 0, err
		}
		subtasks = append(subtasks, task)
	}

	if err := rows.Err(); err != nil {
		return 0, err
	}

	// If no subtasks, get the task's own progress
	if len(subtasks) == 0 {
		task, err := storage.GetTask(db, taskID)
		if err != nil {
			return 0, fmt.Errorf("failed to get task: %w", err)
		}
		return task.Progress, nil
	}

	// Calculate average progress of subtasks
	totalProgress := 0
	for _, subtask := range subtasks {
		totalProgress += subtask.Progress
	}

	avgProgress := totalProgress / len(subtasks)
	return avgProgress, nil
}

// UpdateTaskProgress updates a task's progress and recalculates parent progress if needed
func UpdateTaskProgress(db *sql.DB, taskID string, progress int) error {
	// Validate progress range
	if progress < 0 || progress > 100 {
		return fmt.Errorf("progress must be between 0 and 100, got %d", progress)
	}

	// Get the task to check if it has a parent
	task, err := storage.GetTask(db, taskID)
	if err != nil {
		return fmt.Errorf("failed to get task: %w", err)
	}

	// Update the task's progress
	task.Progress = progress
	if err := storage.UpdateTask(db, task); err != nil {
		return fmt.Errorf("failed to update task progress: %w", err)
	}

	// If this task has a parent, recalculate parent's progress
	if task.ParentID != nil {
		parentProgress, err := CalculateProgress(db, *task.ParentID)
		if err != nil {
			return fmt.Errorf("failed to calculate parent progress: %w", err)
		}

		parent, err := storage.GetTask(db, *task.ParentID)
		if err != nil {
			return fmt.Errorf("failed to get parent task: %w", err)
		}

		parent.Progress = parentProgress
		if err := storage.UpdateTask(db, parent); err != nil {
			return fmt.Errorf("failed to update parent progress: %w", err)
		}
	}

	return nil
}

// Helper function - duplicated from storage package to avoid circular dependency
// TODO: Consider refactoring to avoid duplication
type scanner interface {
	Scan(dest ...interface{}) error
}

func scanTask(s scanner) (*models.Task, error) {
	var task models.Task
	var tagsJSON string
	var createdAtStr, updatedAtStr string
	var completedAtStr, deletedAtStr sql.NullString

	err := s.Scan(
		&task.ID,
		&task.Description,
		&task.Priority,
		&task.Column,
		&task.Progress,
		&task.ParentID,
		&task.Archived,
		&tagsJSON,
		&createdAtStr,
		&updatedAtStr,
		&completedAtStr,
		&deletedAtStr,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to scan task: %w", err)
	}

	// Parse tags
	if err := json.Unmarshal([]byte(tagsJSON), &task.Tags); err != nil {
		return nil, fmt.Errorf("failed to unmarshal tags: %w", err)
	}

	// Parse timestamps
	if task.CreatedAt, err = time.Parse(time.RFC3339, createdAtStr); err != nil {
		return nil, fmt.Errorf("failed to parse created_at: %w", err)
	}

	if task.UpdatedAt, err = time.Parse(time.RFC3339, updatedAtStr); err != nil {
		return nil, fmt.Errorf("failed to parse updated_at: %w", err)
	}

	if completedAtStr.Valid {
		t, err := time.Parse(time.RFC3339, completedAtStr.String)
		if err != nil {
			return nil, fmt.Errorf("failed to parse completed_at: %w", err)
		}
		task.CompletedAt = &t
	}

	if deletedAtStr.Valid {
		t, err := time.Parse(time.RFC3339, deletedAtStr.String)
		if err != nil {
			return nil, fmt.Errorf("failed to parse deleted_at: %w", err)
		}
		task.DeletedAt = &t
	}

	return &task, nil
}
