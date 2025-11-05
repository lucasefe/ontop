package storage

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/lucasefe/ontop/internal/models"
	_ "modernc.org/sqlite"
)

// NewDB creates a new database connection at the specified path
// Default path: ~/.ontop/tasks.db
func NewDB(path string) (*sql.DB, error) {
	// If path is empty, use default
	if path == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			return nil, fmt.Errorf("failed to get home directory: %w", err)
		}
		ontopDir := filepath.Join(home, ".ontop")
		if err := os.MkdirAll(ontopDir, 0755); err != nil {
			return nil, fmt.Errorf("failed to create .ontop directory: %w", err)
		}
		path = filepath.Join(ontopDir, "tasks.db")
	}

	db, err := sql.Open("sqlite", path)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Enable WAL mode for better concurrency
	if _, err := db.Exec("PRAGMA journal_mode=WAL"); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to enable WAL mode: %w", err)
	}

	return db, nil
}

// CreateTask inserts a new task into the database
func CreateTask(db *sql.DB, task *models.Task) error {
	tagsJSON, err := json.Marshal(task.Tags)
	if err != nil {
		return fmt.Errorf("failed to marshal tags: %w", err)
	}

	query := `
		INSERT INTO tasks (
			id, title, description, priority, column, progress, parent_id,
			archived, tags, created_at, updated_at, completed_at, deleted_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	_, err = db.Exec(query,
		task.ID,
		task.Title,
		task.Description,
		task.Priority,
		task.Column,
		task.Progress,
		task.ParentID,
		task.Archived,
		string(tagsJSON),
		task.CreatedAt.Format(time.RFC3339),
		task.UpdatedAt.Format(time.RFC3339),
		formatNullTime(task.CompletedAt),
		formatNullTime(task.DeletedAt),
	)

	if err != nil {
		return fmt.Errorf("failed to insert task: %w", err)
	}

	return nil
}

// GetTask retrieves a single task by ID
func GetTask(db *sql.DB, id string) (*models.Task, error) {
	query := `
		SELECT id, title, description, priority, column, progress, parent_id,
		       archived, tags, created_at, updated_at, completed_at, deleted_at
		FROM tasks
		WHERE id = ? AND deleted_at IS NULL
	`

	row := db.QueryRow(query, id)
	return scanTask(row)
}

// ListTasks retrieves tasks with optional filters
func ListTasks(db *sql.DB, filters map[string]interface{}) ([]*models.Task, error) {
	query := `
		SELECT id, title, description, priority, column, progress, parent_id,
		       archived, tags, created_at, updated_at, completed_at, deleted_at
		FROM tasks
		WHERE deleted_at IS NULL
	`
	args := []interface{}{}

	// Apply filters
	if archived, ok := filters["archived"].(bool); ok {
		query += " AND archived = ?"
		args = append(args, archived)
	} else {
		// Default: only show non-archived
		query += " AND archived = false"
	}

	if column, ok := filters["column"].(string); ok {
		query += " AND column = ?"
		args = append(args, column)
	}

	if priority, ok := filters["priority"].(int); ok {
		query += " AND priority = ?"
		args = append(args, priority)
	}

	if tag, ok := filters["tag"].(string); ok {
		query += ` AND EXISTS (
			SELECT 1 FROM json_each(tags)
			WHERE json_each.value = ?
		)`
		args = append(args, tag)
	}

	query += " ORDER BY created_at DESC"

	rows, err := db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to query tasks: %w", err)
	}
	defer rows.Close()

	var tasks []*models.Task
	for rows.Next() {
		task, err := scanTask(rows)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, task)
	}

	return tasks, rows.Err()
}

// UpdateTask updates an existing task
func UpdateTask(db *sql.DB, task *models.Task) error {
	tagsJSON, err := json.Marshal(task.Tags)
	if err != nil {
		return fmt.Errorf("failed to marshal tags: %w", err)
	}

	query := `
		UPDATE tasks
		SET title = ?, description = ?, priority = ?, column = ?, progress = ?,
		    parent_id = ?, archived = ?, tags = ?, updated_at = ?,
		    completed_at = ?
		WHERE id = ?
	`

	_, err = db.Exec(query,
		task.Title,
		task.Description,
		task.Priority,
		task.Column,
		task.Progress,
		task.ParentID,
		task.Archived,
		string(tagsJSON),
		task.UpdatedAt.Format(time.RFC3339),
		formatNullTime(task.CompletedAt),
		task.ID,
	)

	if err != nil {
		return fmt.Errorf("failed to update task: %w", err)
	}

	return nil
}

// DeleteTask soft-deletes a task by setting deleted_at timestamp
func DeleteTask(db *sql.DB, id string) error {
	query := `UPDATE tasks SET deleted_at = ? WHERE id = ?`
	_, err := db.Exec(query, time.Now().Format(time.RFC3339), id)
	if err != nil {
		return fmt.Errorf("failed to delete task: %w", err)
	}
	return nil
}

// Helper functions

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
		&task.Title,
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

func formatNullTime(t *time.Time) interface{} {
	if t == nil {
		return nil
	}
	return t.Format(time.RFC3339)
}
