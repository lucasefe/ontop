package tui

import (
	"database/sql"
	"log"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/lipgloss"
	"github.com/lucasefe/ontop/internal/config"
	"github.com/lucasefe/ontop/internal/models"
	"github.com/lucasefe/ontop/internal/service"
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

// ViewLayout represents the visual organization of the kanban board
type ViewLayout int

const (
	LayoutColumn ViewLayout = iota // Vertical columns (existing)
	LayoutRow                       // Horizontal rows (new)
)

// SortMode is an alias to service.SortMode for backward compatibility
type SortMode = service.SortMode

// Sort mode constants
const (
	SortByPriority    = service.SortByPriority
	SortByDescription = service.SortByDescription
	SortByCreated     = service.SortByCreated
	SortByUpdated     = service.SortByUpdated
)

// Model represents the Bubbletea application state
type Model struct {
	db              *sql.DB
	tasks           []*models.Task
	currentColumn   int // 0=inbox, 1=in_progress, 2=done
	selectedTask    int // Index within current column
	viewMode        ViewMode
	viewLayout      ViewLayout     // Layout mode for Kanban view (column or row)
	sortMode        SortMode       // How tasks are sorted in columns
	showArchived    bool           // Show archived tasks instead of active
	rowScrollOffset map[int]int    // Horizontal scroll offsets per row (row mode only)
	detailTask      *models.Task
	detailSubtasks  []*models.Task
	moveTask        *models.Task
	moveSelection   int          // Which column to move to
	deleteTask      *models.Task // Task pending deletion
	lastMovedTaskID string       // Track moved task to restore focus
	// Form fields
	formInputs     []textinput.Model
	formTextarea   textarea.Model // For multiline description
	formFocusIndex int
	formTask       *models.Task // Task being created/edited
	formErr        error        // Form validation error (doesn't quit app)
	// UI components
	keys          KeyMap
	help          help.Model
	width         int
	height        int
	err           error
	statusMessage string // Success/info message to display temporarily
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

	// Load config to get view layout preference
	cfg, err := config.Load()
	if err != nil {
		log.Printf("Warning: Failed to load config: %v (using defaults)", err)
	}

	// Parse view mode from config
	viewLayout := LayoutColumn // Default
	if cfg.UI.ViewMode == "row" {
		viewLayout = LayoutRow
	}

	return Model{
		db:              db,
		currentColumn:   0,
		selectedTask:    0,
		viewMode:        ViewModeKanban,
		viewLayout:      viewLayout,
		sortMode:        SortByPriority,
		rowScrollOffset: make(map[int]int),
		keys:            DefaultKeyMap(),
		help:            h,
		width:           80,
		height:          24,
	}
}

// GetTasksByColumn returns tasks filtered by column in hierarchical display order
func (m *Model) GetTasksByColumn(column string) []*models.Task {
	var filtered []*models.Task
	for _, task := range m.tasks {
		if task.Column == column {
			filtered = append(filtered, task)
		}
	}

	// Build hierarchical display order
	hierarchical := service.BuildFlatHierarchy(filtered, m.sortMode)

	// Extract Task pointers (unwrap HierarchicalTask)
	result := make([]*models.Task, len(hierarchical))
	for i, ht := range hierarchical {
		result[i] = ht.Task
	}

	return result
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
