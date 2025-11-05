package tui

import (
	"database/sql"

	"github.com/lucasefe/ontop/internal/models"
)

// ViewMode represents the current view in the TUI
type ViewMode int

const (
	ViewModeKanban ViewMode = iota
	ViewModeDetail
	ViewModeMove
)

// Model represents the Bubbletea application state
type Model struct {
	db              *sql.DB
	tasks           []*models.Task
	currentColumn   int // 0=inbox, 1=in_progress, 2=done
	selectedTask    int // Index within current column
	viewMode        ViewMode
	detailTask      *models.Task
	detailSubtasks  []*models.Task
	moveTask        *models.Task
	moveSelection   int // Which column to move to
	width           int
	height          int
	err             error
}

// NewModel creates a new TUI model
func NewModel(db *sql.DB) Model {
	return Model{
		db:            db,
		currentColumn: 0,
		selectedTask:  0,
		viewMode:      ViewModeKanban,
		width:         80,
		height:        24,
	}
}

// GetTasksByColumn returns tasks filtered by column
func (m *Model) GetTasksByColumn(column string) []*models.Task {
	var filtered []*models.Task
	for _, task := range m.tasks {
		if task.Column == column {
			filtered = append(filtered, task)
		}
	}
	return filtered
}

// GetCurrentColumnName returns the name of the currently selected column
func (m *Model) GetCurrentColumnName() string {
	columns := []string{models.ColumnInbox, models.ColumnInProgress, models.ColumnDone}
	if m.currentColumn >= 0 && m.currentColumn < len(columns) {
		return columns[m.currentColumn]
	}
	return models.ColumnInbox
}

// GetSelectedTask returns the currently selected task
func (m *Model) GetSelectedTask() *models.Task {
	columnTasks := m.GetTasksByColumn(m.GetCurrentColumnName())
	if m.selectedTask >= 0 && m.selectedTask < len(columnTasks) {
		return columnTasks[m.selectedTask]
	}
	return nil
}
