package tui

import (
	"time"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/lucasefe/ontop/internal/models"
	"github.com/lucasefe/ontop/internal/storage"
)

type tasksLoadedMsg struct {
	tasks []*models.Task
	err   error
}

// Init initializes the model and loads tasks from database
func (m Model) Init() tea.Cmd {
	return m.loadTasks
}

// loadTasks loads all active tasks from the database
func (m Model) loadTasks() tea.Msg {
	filters := map[string]interface{}{
		"archived": false,
	}
	tasks, err := storage.ListTasks(m.db, filters)
	if err != nil {
		return tasksLoadedMsg{err: err}
	}
	return tasksLoadedMsg{tasks: tasks}
}

// Update handles messages and updates the model
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil

	case tasksLoadedMsg:
		if msg.err != nil {
			m.err = msg.err
			return m, tea.Quit
		}
		m.tasks = msg.tasks
		return m, nil

	case tea.KeyMsg:
		return m.handleKeyPress(msg)
	}

	return m, nil
}

// handleKeyPress processes keyboard input
func (m Model) handleKeyPress(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	keys := DefaultKeyMap()

	// Quit
	if key.Matches(msg, keys.Quit) {
		return m, tea.Quit
	}

	// Mode-specific handling
	switch m.viewMode {
	case ViewModeKanban:
		return m.handleKanbanKeys(msg, keys)
	case ViewModeDetail:
		return m.handleDetailKeys(msg, keys)
	case ViewModeMove:
		return m.handleMoveKeys(msg, keys)
	}

	return m, nil
}

// handleKanbanKeys handles key presses in kanban view
func (m Model) handleKanbanKeys(msg tea.KeyMsg, keys KeyMap) (tea.Model, tea.Cmd) {
	// Vertical navigation (j/k)
	if key.Matches(msg, keys.Down) {
		columnTasks := m.GetTasksByColumn(m.GetCurrentColumnName())
		if m.selectedTask < len(columnTasks)-1 {
			m.selectedTask++
		}
		return m, nil
	}

	if key.Matches(msg, keys.Up) {
		if m.selectedTask > 0 {
			m.selectedTask--
		}
		return m, nil
	}

	// Horizontal navigation (h/l)
	if key.Matches(msg, keys.Left) {
		if m.currentColumn > 0 {
			m.currentColumn--
			m.selectedTask = 0
		}
		return m, nil
	}

	if key.Matches(msg, keys.Right) {
		if m.currentColumn < 2 {
			m.currentColumn++
			m.selectedTask = 0
		}
		return m, nil
	}

	// Enter detail view
	if key.Matches(msg, keys.Select) {
		task := m.GetSelectedTask()
		if task != nil {
			m.viewMode = ViewModeDetail
			m.detailTask = task
			// Load subtasks
			filters := map[string]interface{}{
				"archived": false,
			}
			allTasks, err := storage.ListTasks(m.db, filters)
			if err == nil {
				m.detailSubtasks = []*models.Task{}
				for _, t := range allTasks {
					if t.ParentID != nil && *t.ParentID == task.ID {
						m.detailSubtasks = append(m.detailSubtasks, t)
					}
				}
			}
		}
		return m, nil
	}

	// Refresh
	if key.Matches(msg, keys.Refresh) {
		return m, m.loadTasks
	}

	// Move task
	if key.Matches(msg, keys.Move) {
		task := m.GetSelectedTask()
		if task != nil {
			m.viewMode = ViewModeMove
			m.moveTask = task
			// Default to current column
			m.moveSelection = m.currentColumn
		}
		return m, nil
	}

	return m, nil
}

// handleDetailKeys handles key presses in detail view
func (m Model) handleDetailKeys(msg tea.KeyMsg, keys KeyMap) (tea.Model, tea.Cmd) {
	// Back to kanban
	if key.Matches(msg, keys.Back) {
		m.viewMode = ViewModeKanban
		m.detailTask = nil
		m.detailSubtasks = nil
		return m, nil
	}

	return m, nil
}

// handleMoveKeys handles key presses in move prompt view
func (m Model) handleMoveKeys(msg tea.KeyMsg, keys KeyMap) (tea.Model, tea.Cmd) {
	// Cancel
	if key.Matches(msg, keys.Back) {
		m.viewMode = ViewModeKanban
		m.moveTask = nil
		return m, nil
	}

	// Navigate between columns
	if key.Matches(msg, keys.Left) {
		if m.moveSelection > 0 {
			m.moveSelection--
		}
		return m, nil
	}

	if key.Matches(msg, keys.Right) {
		if m.moveSelection < 2 {
			m.moveSelection++
		}
		return m, nil
	}

	// Confirm move
	if key.Matches(msg, keys.Select) {
		return m.confirmMove()
	}

	return m, nil
}

// confirmMove performs the actual task move operation
func (m Model) confirmMove() (tea.Model, tea.Cmd) {
	if m.moveTask == nil {
		m.viewMode = ViewModeKanban
		return m, nil
	}

	// Get target column
	columns := []string{models.ColumnInbox, models.ColumnInProgress, models.ColumnDone}
	targetColumn := columns[m.moveSelection]

	// Don't move if already in target column
	if m.moveTask.Column == targetColumn {
		m.viewMode = ViewModeKanban
		m.moveTask = nil
		return m, nil
	}

	// Update task
	oldColumn := m.moveTask.Column
	m.moveTask.Column = targetColumn
	m.moveTask.UpdatedAt = time.Now()

	// Set completed_at if moving to done
	if targetColumn == models.ColumnDone && m.moveTask.CompletedAt == nil {
		now := time.Now()
		m.moveTask.CompletedAt = &now
	}

	// Clear completed_at if moving from done
	if oldColumn == models.ColumnDone && targetColumn != models.ColumnDone {
		m.moveTask.CompletedAt = nil
	}

	// Save to database
	if err := storage.UpdateTask(m.db, m.moveTask); err != nil {
		m.err = err
		return m, tea.Quit
	}

	// Return to kanban and reload tasks
	m.viewMode = ViewModeKanban
	m.moveTask = nil
	return m, m.loadTasks
}

// View renders the UI
func (m Model) View() string {
	if m.err != nil {
		return "Error: " + m.err.Error() + "\n\nPress q to quit."
	}

	switch m.viewMode {
	case ViewModeKanban:
		return m.renderKanban()
	case ViewModeDetail:
		return m.renderDetail()
	case ViewModeMove:
		return m.renderMovePrompt()
	}

	return ""
}
