package storage

import (
	"database/sql"
	"fmt"
)

// InitSchema creates the database schema if it doesn't exist
func InitSchema(db *sql.DB) error {
	schema := `
	CREATE TABLE IF NOT EXISTS tasks (
		id TEXT PRIMARY KEY,
		description TEXT NOT NULL,
		priority INTEGER NOT NULL CHECK (priority >= 1 AND priority <= 5),
		column TEXT NOT NULL CHECK (column IN ('inbox', 'in_progress', 'done')),
		progress INTEGER NOT NULL DEFAULT 0 CHECK (progress >= 0 AND progress <= 100),
		parent_id TEXT,
		archived INTEGER NOT NULL DEFAULT 0,
		tags TEXT NOT NULL DEFAULT '[]',
		created_at TEXT NOT NULL,
		updated_at TEXT NOT NULL,
		completed_at TEXT,
		deleted_at TEXT,
		FOREIGN KEY (parent_id) REFERENCES tasks(id)
	);

	CREATE INDEX IF NOT EXISTS idx_tasks_column ON tasks(column);
	CREATE INDEX IF NOT EXISTS idx_tasks_priority ON tasks(priority);
	CREATE INDEX IF NOT EXISTS idx_tasks_parent_id ON tasks(parent_id);
	CREATE INDEX IF NOT EXISTS idx_tasks_archived ON tasks(archived);
	CREATE INDEX IF NOT EXISTS idx_tasks_deleted_at ON tasks(deleted_at);
	CREATE INDEX IF NOT EXISTS idx_tasks_created_at ON tasks(created_at);
	`

	_, err := db.Exec(schema)
	if err != nil {
		return fmt.Errorf("failed to initialize schema: %w", err)
	}

	return nil
}
