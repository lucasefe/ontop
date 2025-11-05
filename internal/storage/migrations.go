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
		title TEXT NOT NULL DEFAULT '',
		description TEXT NOT NULL DEFAULT '',
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

	// Add title column if it doesn't exist (migration for existing databases)
	_, _ = db.Exec(`
		ALTER TABLE tasks ADD COLUMN title TEXT NOT NULL DEFAULT '';
	`)
	// Ignore error if column already exists

	// Migrate existing data: copy description to title if title is empty
	_, err = db.Exec(`
		UPDATE tasks SET title = substr(description, 1, 100) WHERE title = '' OR title IS NULL;
	`)
	if err != nil {
		return fmt.Errorf("failed to migrate titles: %w", err)
	}

	return nil
}
