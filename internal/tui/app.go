package tui

import (
	"log"
	"time"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/lucasefe/ontop/internal/config"
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

// loadTasks loads tasks from the database based on showArchived flag
func (m Model) loadTasks() tea.Msg {
	filters := map[string]interface{}{
		"archived": m.showArchived,
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
		m.help.Width = msg.Width
		return m, nil

	case tasksLoadedMsg:
		if msg.err != nil {
			m.err = msg.err
			return m, tea.Quit
		}
		m.tasks = msg.tasks

		// If we just moved a task, restore focus on it in its new location
		if m.lastMovedTaskID != "" {
			columnTasks := m.GetTasksByColumn(m.GetCurrentColumnName())
			for i, task := range columnTasks {
				if task.ID == m.lastMovedTaskID {
					m.selectedTask = i
					break
				}
			}
			m.lastMovedTaskID = "" // Clear the tracking
		}

		return m, nil

	case tea.KeyMsg:
		return m.handleKeyPress(msg)
	}

	return m, nil
}

// handleKeyPress processes keyboard input
func (m Model) handleKeyPress(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	keys := m.keys

	// Quit (should always work)
	if key.Matches(msg, keys.Quit) {
		return m, tea.Quit
	}

	// Handle form modes first to allow "?" to be typed in text fields
	if m.viewMode == ViewModeCreate || m.viewMode == ViewModeEdit {
		return m.handleFormKeys(msg, keys)
	}

	// Toggle help (only for non-form modes)
	if key.Matches(msg, keys.Help) {
		m.help.ShowAll = !m.help.ShowAll
		return m, nil
	}

	// Mode-specific handling
	switch m.viewMode {
	case ViewModeKanban:
		return m.handleKanbanKeys(msg, keys)
	case ViewModeDetail:
		return m.handleDetailKeys(msg, keys)
	case ViewModeMove:
		return m.handleMoveKeys(msg, keys)
	case ViewModeDeleteConfirm:
		return m.handleDeleteConfirmKeys(msg, keys)
	}

	return m, nil
}

// handleKanbanKeys handles key presses in kanban view
func (m Model) handleKanbanKeys(msg tea.KeyMsg, keys KeyMap) (tea.Model, tea.Cmd) {
	// Context-sensitive navigation based on view layout
	if m.viewLayout == LayoutRow {
		// ROW MODE: up/down navigates tasks seamlessly, auto-switching rows at boundaries
		if key.Matches(msg, keys.Down) {
			columnTasks := m.GetTasksByColumn(m.GetCurrentColumnName())
			if m.selectedTask < len(columnTasks)-1 {
				// Move to next task in current row
				m.selectedTask++
			} else if m.currentColumn < 2 {
				// At end of row, move to next row (first task)
				m.currentColumn++
				m.selectedTask = 0
			}
			return m, nil
		}

		if key.Matches(msg, keys.Up) {
			if m.selectedTask > 0 {
				// Move to previous task in current row
				m.selectedTask--
			} else if m.currentColumn > 0 {
				// At top of row, move to previous row (last task)
				m.currentColumn--
				prevColumnTasks := m.GetTasksByColumn(m.GetCurrentColumnName())
				if len(prevColumnTasks) > 0 {
					m.selectedTask = len(prevColumnTasks) - 1
				} else {
					m.selectedTask = 0
				}
			}
			return m, nil
		}

		// Left/Right do nothing in row mode (or could be used for other features)
		if key.Matches(msg, keys.Left) {
			return m, nil
		}

		if key.Matches(msg, keys.Right) {
			return m, nil
		}
	} else {
		// COLUMN MODE: up/down navigates within column, left/right switches columns
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
	}

	// Quick move with Shift+navigation keys (context-sensitive)
	if m.viewLayout == LayoutColumn {
		// COLUMN MODE: Shift+H moves left (to previous workflow stage), Shift+L moves right (to next workflow stage)
		if key.Matches(msg, keys.QuickMoveLeft) {
			task := m.GetSelectedTask()
			if task != nil && m.currentColumn > 0 {
				m.moveTask = task
				m.moveSelection = m.currentColumn - 1
				return m.confirmMove()
			}
			return m, nil
		}

		if key.Matches(msg, keys.QuickMoveRight) {
			task := m.GetSelectedTask()
			if task != nil && m.currentColumn < 2 {
				m.moveTask = task
				m.moveSelection = m.currentColumn + 1
				return m.confirmMove()
			}
			return m, nil
		}
	} else {
		// ROW MODE: Shift+K moves up (to previous workflow stage), Shift+J moves down (to next workflow stage)
		if key.Matches(msg, keys.QuickMoveUp) {
			task := m.GetSelectedTask()
			if task != nil && m.currentColumn > 0 {
				m.moveTask = task
				m.moveSelection = m.currentColumn - 1
				return m.confirmMove()
			}
			return m, nil
		}

		if key.Matches(msg, keys.QuickMoveDown) {
			task := m.GetSelectedTask()
			if task != nil && m.currentColumn < 2 {
				m.moveTask = task
				m.moveSelection = m.currentColumn + 1
				return m.confirmMove()
			}
			return m, nil
		}
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

	// New task
	if key.Matches(msg, keys.New) {
		m.viewMode = ViewModeCreate
		m.initCreateForm(nil) // Pass nil for top-level task
		return m, nil
	}

	// Toggle sort
	if key.Matches(msg, keys.Sort) {
		m.ToggleSortMode()
		return m, nil
	}

	// Toggle archived view
	if key.Matches(msg, keys.ToggleArchive) {
		m.showArchived = !m.showArchived
		m.selectedTask = 0
		return m, m.loadTasks
	}

	// Toggle view layout (column vs row)
	if key.Matches(msg, keys.ToggleView) {
		return m.handleToggleView()
	}

	// Archive/Unarchive task (toggle based on current view)
	if key.Matches(msg, keys.Archive) {
		task := m.GetSelectedTask()
		if task != nil {
			// Toggle: if viewing archived, unarchive; if viewing active, archive
			task.Archived = !m.showArchived
			task.UpdatedAt = time.Now()
			if err := storage.UpdateTask(m.db, task); err != nil {
				m.err = err
				return m, tea.Quit
			}
			return m, m.loadTasks
		}
		return m, nil
	}

	// Delete task (show confirmation)
	if key.Matches(msg, keys.Delete) {
		task := m.GetSelectedTask()
		if task != nil {
			m.viewMode = ViewModeDeleteConfirm
			m.deleteTask = task
			m.moveSelection = 0 // Default to "No"
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
		m.statusMessage = "" // Clear status message
		return m, nil
	}

	// Create subtask (press 'n' in detail view)
	if key.Matches(msg, keys.New) {
		if m.detailTask != nil {
			m.viewMode = ViewModeCreate
			m.initCreateForm(&m.detailTask.ID) // Pass parent ID to pre-fill
			m.statusMessage = ""                // Clear status message
		}
		return m, nil
	}

	// Edit task
	if key.Matches(msg, keys.Edit) {
		if m.detailTask != nil {
			m.viewMode = ViewModeEdit
			m.initEditForm(m.detailTask)
			m.statusMessage = "" // Clear status message
		}
		return m, nil
	}

	// Archive/Unarchive task (toggle based on current view)
	if key.Matches(msg, keys.Archive) {
		if m.detailTask != nil {
			// Toggle: if viewing archived, unarchive; if viewing active, archive
			m.detailTask.Archived = !m.showArchived
			m.detailTask.UpdatedAt = time.Now()
			if err := storage.UpdateTask(m.db, m.detailTask); err != nil {
				m.err = err
				return m, tea.Quit
			}
			m.viewMode = ViewModeKanban
			m.detailTask = nil
			m.detailSubtasks = nil
			return m, m.loadTasks
		}
		return m, nil
	}

	// Delete task (show confirmation)
	if key.Matches(msg, keys.Delete) {
		if m.detailTask != nil {
			m.viewMode = ViewModeDeleteConfirm
			m.deleteTask = m.detailTask
			m.moveSelection = 0 // Default to "No"
		}
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

// handleDeleteConfirmKeys handles key presses in delete confirmation view
func (m Model) handleDeleteConfirmKeys(msg tea.KeyMsg, keys KeyMap) (tea.Model, tea.Cmd) {
	// Cancel
	if key.Matches(msg, keys.Back) {
		m.viewMode = ViewModeKanban
		m.deleteTask = nil
		return m, nil
	}

	// Navigate between No/Yes
	if key.Matches(msg, keys.Left) {
		if m.moveSelection > 0 {
			m.moveSelection--
		}
		return m, nil
	}

	if key.Matches(msg, keys.Right) {
		if m.moveSelection < 1 {
			m.moveSelection++
		}
		return m, nil
	}

	// Confirm
	if key.Matches(msg, keys.Select) {
		if m.moveSelection == 1 { // Yes, delete
			if m.deleteTask != nil {
				if err := storage.DeleteTask(m.db, m.deleteTask.ID); err != nil {
					m.err = err
					return m, tea.Quit
				}
			}
		}
		// Either way, go back to kanban
		m.viewMode = ViewModeKanban
		m.deleteTask = nil
		m.detailTask = nil
		m.detailSubtasks = nil
		return m, m.loadTasks
	}

	return m, nil
}

// handleFormKeys handles key presses in create/edit form view
func (m Model) handleFormKeys(msg tea.KeyMsg, keys KeyMap) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	// Cancel
	if key.Matches(msg, keys.Back) {
		m.viewMode = ViewModeKanban
		m.formInputs = nil
		m.formTask = nil
		m.formErr = nil
		return m, nil
	}

	// Save with Ctrl+S (works from any field)
	if key.Matches(msg, keys.Save) {
		return m.saveForm()
	}

	// Shift+Tab to previous field (total 6 fields: 0-5)
	// 0: Title, 1: Description (textarea), 2: Priority, 3: Progress, 4: Tags, 5: Parent ID
	if key.Matches(msg, keys.ShiftTab) {
		// Blur current field
		if m.formFocusIndex == 1 {
			m.formTextarea.Blur()
		} else {
			// Map focus index to formInputs index
			// Focus 0 -> formInputs[0]
			// Focus 2+ -> formInputs[focus-1]
			inputIdx := m.formFocusIndex
			if inputIdx > 1 {
				inputIdx--
			}
			if inputIdx >= 0 && inputIdx < len(m.formInputs) {
				m.formInputs[inputIdx].Blur()
			}
		}

		// Move to previous field (with wrap-around)
		m.formFocusIndex = (m.formFocusIndex - 1 + 6) % 6

		// Focus new field
		if m.formFocusIndex == 1 {
			m.formTextarea.Focus()
		} else {
			inputIdx := m.formFocusIndex
			if inputIdx > 1 {
				inputIdx--
			}
			if inputIdx >= 0 && inputIdx < len(m.formInputs) {
				m.formInputs[inputIdx].Focus()
			}
		}
		return m, nil
	}

	// Tab to next field (total 6 fields: 0-5)
	// 0: Title, 1: Description (textarea), 2: Priority, 3: Progress, 4: Tags, 5: Parent ID
	if key.Matches(msg, keys.Tab) {
		// Blur current field
		if m.formFocusIndex == 1 {
			m.formTextarea.Blur()
		} else {
			// Map focus index to formInputs index
			// Focus 0 -> formInputs[0]
			// Focus 2+ -> formInputs[focus-1]
			inputIdx := m.formFocusIndex
			if inputIdx > 1 {
				inputIdx--
			}
			if inputIdx >= 0 && inputIdx < len(m.formInputs) {
				m.formInputs[inputIdx].Blur()
			}
		}

		// Move to next field
		m.formFocusIndex = (m.formFocusIndex + 1) % 6

		// Focus new field
		if m.formFocusIndex == 1 {
			m.formTextarea.Focus()
		} else {
			inputIdx := m.formFocusIndex
			if inputIdx > 1 {
				inputIdx--
			}
			if inputIdx >= 0 && inputIdx < len(m.formInputs) {
				m.formInputs[inputIdx].Focus()
			}
		}
		return m, nil
	}

	// Clear form error when user types (so error disappears when they fix it)
	m.formErr = nil

	// Update the focused field
	if m.formFocusIndex == 1 {
		// Textarea for description
		m.formTextarea, cmd = m.formTextarea.Update(msg)
	} else {
		// Text input field
		inputIdx := m.formFocusIndex
		if inputIdx > 1 {
			inputIdx--
		}
		if inputIdx >= 0 && inputIdx < len(m.formInputs) {
			m.formInputs[inputIdx], cmd = m.formInputs[inputIdx].Update(msg)
		}
	}

	return m, cmd
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

	// Track the moved task ID to restore focus after reload
	m.lastMovedTaskID = m.moveTask.ID

	// Update current column to follow the moved task
	m.currentColumn = m.moveSelection

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
		if m.viewLayout == LayoutRow {
			return m.renderKanbanRows()
		}
		return m.renderKanban()
	case ViewModeDetail:
		return m.renderDetail()
	case ViewModeMove:
		return m.renderMovePrompt()
	case ViewModeDeleteConfirm:
		return m.renderDeleteConfirm()
	case ViewModeCreate, ViewModeEdit:
		return m.renderForm()
	}

	return ""
}

// handleToggleView toggles between column and row layouts
func (m Model) handleToggleView() (Model, tea.Cmd) {
	// Toggle layout
	if m.viewLayout == LayoutColumn {
		m.viewLayout = LayoutRow
	} else {
		m.viewLayout = LayoutColumn
	}

	// Reset scroll offsets when toggling
	m.rowScrollOffset = make(map[int]int)

	// Save preference to config asynchronously
	go func() {
		cfg := config.Config{
			UI: config.UIConfig{
				ViewMode: "column",
			},
		}
		if m.viewLayout == LayoutRow {
			cfg.UI.ViewMode = "row"
		}

		if err := config.Save(cfg); err != nil {
			log.Printf("Failed to save view mode preference: %v", err)
		}
	}()

	return m, nil
}
