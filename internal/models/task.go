package models

import "time"

// Task represents a work item with all attributes
type Task struct {
	ID          string     `json:"id"`
	Description string     `json:"description"`
	Priority    int        `json:"priority"`    // 1-5 where 1 is highest
	Column      string     `json:"column"`      // inbox, in_progress, done
	Progress    int        `json:"progress"`    // 0-100
	ParentID    *string    `json:"parent_id"`   // NULL for top-level tasks
	Archived    bool       `json:"archived"`
	Tags        []string   `json:"tags"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	CompletedAt *time.Time `json:"completed_at"` // NULL if not completed
	DeletedAt   *time.Time `json:"deleted_at"`   // NULL if not deleted
}
