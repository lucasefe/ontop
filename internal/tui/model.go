package tui

import (
	"database/sql"
	"sort"
	"strings"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/lipgloss"
	"github.com/lucasefe/ontop/internal/models"
)

// ViewMode represents the current view in the TUI
type ViewMode int

const (
	ViewModeKanban ViewMode = iota
	ViewModeDetail
	ViewModeMove
	ViewModeCreate
	ViewModeEdit
	ViewModeDeleteConfirm
)

// SortMode represents how tasks are sorted
type SortMode int

const (
	SortByPriority SortMode = iota
	SortByDescription
	SortByCreated
	SortByUpdated
)

// Model represents the Bubbletea application state
type Model struct {
	db             *sql.DB
	tasks          []*models.Task
	currentColumn  int // 0=inbox, 1=in_progress, 2=done
	selectedTask   int // Index within current column
	viewMode       ViewMode
	sortMode       SortMode // How tasks are sorted in columns
	showArchived   bool     // Show archived tasks instead of active
	detailTask     *models.Task
	detailSubtasks []*models.Task
	moveTask       *models.Task
	moveSelection  int          // Which column to move to
	deleteTask     *models.Task // Task pending deletion
	// Form fields
	formInputs     []textinput.Model
	formFocusIndex int
	formTask       *models.Task // Task being created/edited
	formErr        error        // Form validation error (doesn't quit app)
	// UI components
	keys   KeyMap
	help   help.Model
	width  int
	height int
	err    error
}

// NewModel creates a new TUI model
func NewModel(db *sql.DB) Model {
	h := help.New()
	// Apply Gruvbox theme colors to help
	gruvboxGreen := "#b8bb26"
	gruvboxGray := "#928374"
	h.Styles.ShortKey = h.Styles.ShortKey.Foreground(lipgloss.Color(gruvboxGreen))
	h.Styles.ShortDesc = h.Styles.ShortDesc.Foreground(lipgloss.Color(gruvboxGray))
	h.Styles.FullKey = h.Styles.FullKey.Foreground(lipgloss.Color(gruvboxGreen))
	h.Styles.FullDesc = h.Styles.FullDesc.Foreground(lipgloss.Color(gruvboxGray))
	h.Styles.ShortSeparator = h.Styles.ShortSeparator.Foreground(lipgloss.Color(gruvboxGray))
	h.Styles.FullSeparator = h.Styles.FullSeparator.Foreground(lipgloss.Color(gruvboxGray))

	return Model{
		db:            db,
		currentColumn: 0,
		selectedTask:  0,
		viewMode:      ViewModeKanban,
		sortMode:      SortByPriority,
		keys:          DefaultKeyMap(),
		help:          h,
		width:         80,
		height:        24,
	}
}

// GetTasksByColumn returns tasks filtered by column and sorted by current sort mode
func (m *Model) GetTasksByColumn(column string) []*models.Task {
	var filtered []*models.Task
	for _, task := range m.tasks {
		if task.Column == column {
			filtered = append(filtered, task)
		}
	}

	// Sort based on current sort mode
	sort.Slice(filtered, func(i, j int) bool {
		switch m.sortMode {
		case SortByPriority:
			return filtered[i].Priority < filtered[j].Priority // Lower number = higher priority
		case SortByDescription:
			return strings.ToLower(filtered[i].Description) < strings.ToLower(filtered[j].Description)
		case SortByCreated:
			return filtered[i].CreatedAt.After(filtered[j].CreatedAt) // Newest first
		case SortByUpdated:
			return filtered[i].UpdatedAt.After(filtered[j].UpdatedAt) // Most recently updated first
		default:
			return filtered[i].Priority < filtered[j].Priority
		}
	})

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

// GetSortModeName returns the display name for the current sort mode
func (m *Model) GetSortModeName() string {
	switch m.sortMode {
	case SortByPriority:
		return "Priority"
	case SortByDescription:
		return "Description"
	case SortByCreated:
		return "Created"
	case SortByUpdated:
		return "Updated"
	default:
		return "Priority"
	}
}

// ToggleSortMode cycles to the next sort mode
func (m *Model) ToggleSortMode() {
	m.sortMode = (m.sortMode + 1) % 4
	// Reset selection to top when sort changes
	m.selectedTask = 0
}
