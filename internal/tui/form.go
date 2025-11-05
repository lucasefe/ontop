package tui

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/lucasefe/ontop/internal/models"
	"github.com/lucasefe/ontop/internal/service"
	"github.com/lucasefe/ontop/internal/storage"
)

var (
	formStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(gruvboxGreen).
			Padding(1, 2).
			MarginTop(1)

	formTitleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(gruvboxGreen)

	formLabelStyle = lipgloss.NewStyle().
			Foreground(gruvboxGray).
			Bold(true)

	formHelpStyle = lipgloss.NewStyle().
			Foreground(gruvboxGray).
			Italic(true)
)

// initCreateForm initializes the form for creating a new task
func (m *Model) initCreateForm() {
	// Calculate input width based on terminal width
	inputWidth := m.width - 20 // Leave space for borders and padding
	if inputWidth < 40 {
		inputWidth = 40
	}

	// Create text inputs (6 fields: title, description, priority, progress, tags, parent)
	m.formInputs = make([]textinput.Model, 6)

	// Title (short)
	m.formInputs[0] = textinput.New()
	m.formInputs[0].Placeholder = "Task title"
	m.formInputs[0].Focus()
	m.formInputs[0].Width = inputWidth

	// Description (full text)
	m.formInputs[1] = textinput.New()
	m.formInputs[1].Placeholder = "Full description (optional)"
	m.formInputs[1].Width = inputWidth

	// Priority
	m.formInputs[2] = textinput.New()
	m.formInputs[2].Placeholder = "1-5 (default: 3)"
	m.formInputs[2].SetValue("3")
	m.formInputs[2].CharLimit = 1
	m.formInputs[2].Width = 20

	// Progress
	m.formInputs[3] = textinput.New()
	m.formInputs[3].Placeholder = "0-100 (default: 0)"
	m.formInputs[3].SetValue("0")
	m.formInputs[3].CharLimit = 3
	m.formInputs[3].Width = 20

	// Tags
	m.formInputs[4] = textinput.New()
	m.formInputs[4].Placeholder = "tag1,tag2,tag3"
	m.formInputs[4].Width = inputWidth

	// Parent task ID
	m.formInputs[5] = textinput.New()
	m.formInputs[5].Placeholder = "Parent task ID (optional)"
	m.formInputs[5].Width = inputWidth

	m.formFocusIndex = 0
	m.formTask = nil
}

// initEditForm initializes the form for editing an existing task
func (m *Model) initEditForm(task *models.Task) {
	// Calculate input width based on terminal width
	inputWidth := m.width - 20 // Leave space for borders and padding
	if inputWidth < 40 {
		inputWidth = 40
	}

	// Create text inputs (6 fields: title, description, priority, progress, tags, parent)
	m.formInputs = make([]textinput.Model, 6)

	// Title
	m.formInputs[0] = textinput.New()
	m.formInputs[0].Placeholder = "Task title"
	m.formInputs[0].SetValue(task.Title)
	m.formInputs[0].Focus()
	m.formInputs[0].Width = inputWidth

	// Description
	m.formInputs[1] = textinput.New()
	m.formInputs[1].Placeholder = "Full description (optional)"
	m.formInputs[1].SetValue(task.Description)
	m.formInputs[1].Width = inputWidth

	// Priority
	m.formInputs[2] = textinput.New()
	m.formInputs[2].Placeholder = "1-5 (1=highest)"
	m.formInputs[2].SetValue(strconv.Itoa(task.Priority))
	m.formInputs[2].CharLimit = 1
	m.formInputs[2].Width = 20

	// Progress
	m.formInputs[3] = textinput.New()
	m.formInputs[3].Placeholder = "0-100"
	m.formInputs[3].SetValue(strconv.Itoa(task.Progress))
	m.formInputs[3].CharLimit = 3
	m.formInputs[3].Width = 20

	// Tags
	m.formInputs[4] = textinput.New()
	m.formInputs[4].Placeholder = "tag1,tag2,tag3"
	m.formInputs[4].SetValue(strings.Join(task.Tags, ","))
	m.formInputs[4].Width = inputWidth

	// Parent task ID
	m.formInputs[5] = textinput.New()
	m.formInputs[5].Placeholder = "Parent task ID (clear to remove parent)"
	if task.ParentID != nil {
		m.formInputs[5].SetValue(*task.ParentID)
	}
	m.formInputs[5].Width = inputWidth

	m.formFocusIndex = 0
	m.formTask = task
}

// renderForm renders the create/edit form
func (m Model) renderForm() string {
	var b strings.Builder

	// Calculate content width based on terminal
	contentWidth := m.width - 8 // Leave space for borders
	if contentWidth < 40 {
		contentWidth = 40
	}

	// Title
	title := "Create New Task"
	if m.viewMode == ViewModeEdit {
		title = "Edit Task"
	}
	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(gruvboxGreen).
		Width(contentWidth).
		Align(lipgloss.Center)
	b.WriteString(titleStyle.Render(title) + "\n\n")

	// Form content
	var formContent strings.Builder

	// Title field
	formContent.WriteString(formLabelStyle.Render("Title:") + "\n")
	formContent.WriteString(m.formInputs[0].View() + "\n\n")

	// Description field
	formContent.WriteString(formLabelStyle.Render("Description (full text):") + "\n")
	formContent.WriteString(m.formInputs[1].View() + "\n\n")

	// Priority field
	formContent.WriteString(formLabelStyle.Render("Priority (1-5, 1=highest):") + "\n")
	formContent.WriteString(m.formInputs[2].View() + "\n\n")

	// Progress field
	formContent.WriteString(formLabelStyle.Render("Progress (0-100):") + "\n")
	formContent.WriteString(m.formInputs[3].View() + "\n\n")

	// Tags field
	formContent.WriteString(formLabelStyle.Render("Tags (comma-separated):") + "\n")
	formContent.WriteString(m.formInputs[4].View() + "\n\n")

	// Parent task ID field
	formContent.WriteString(formLabelStyle.Render("Parent Task ID (for subtasks):") + "\n")
	formContent.WriteString(m.formInputs[5].View() + "\n")
	if m.viewMode == ViewModeEdit && m.formTask != nil && m.formTask.ParentID != nil {
		formContent.WriteString(formHelpStyle.Render("  (clear to remove parent relationship)") + "\n")
	}
	formContent.WriteString("\n")

	// Apply form style with dynamic width
	formStyleDynamic := formStyle.Width(contentWidth)
	b.WriteString(formStyleDynamic.Render(formContent.String()))
	b.WriteString("\n")

	// Show existing subtasks if editing
	if m.viewMode == ViewModeEdit && m.formTask != nil {
		subtasks := m.getSubtasksForTask(m.formTask.ID)
		if len(subtasks) > 0 {
			var subtaskContent strings.Builder
			subtaskContent.WriteString(formLabelStyle.Render(fmt.Sprintf("Existing Subtasks (%d):", len(subtasks))) + "\n\n")
			for _, st := range subtasks {
				progressStr := ""
				if st.Progress > 0 {
					progressStr = fmt.Sprintf(" (%d%%)", st.Progress)
				}
				subtaskContent.WriteString(fmt.Sprintf("  • [%s] P%d: %s%s\n", st.ID, st.Priority, st.Description, progressStr))
			}
			b.WriteString(formStyleDynamic.Render(subtaskContent.String()))
			b.WriteString("\n")
		}
	}

	// Help view
	helpView := m.help.View(m.keys)
	b.WriteString(helpStyle.Render(helpView))

	// Show form error if present
	if m.formErr != nil {
		errorStyle := lipgloss.NewStyle().
			Foreground(gruvboxRed).
			Bold(true).
			MarginTop(1)
		b.WriteString("\n")
		b.WriteString(errorStyle.Render("Error: " + m.formErr.Error()))
	}

	return b.String()
}

// getSubtasksForTask returns all subtasks for a given task ID
func (m *Model) getSubtasksForTask(parentID string) []*models.Task {
	var subtasks []*models.Task
	for _, task := range m.tasks {
		if task.ParentID != nil && *task.ParentID == parentID {
			subtasks = append(subtasks, task)
		}
	}
	return subtasks
}

// saveForm validates and saves the form
func (m Model) saveForm() (tea.Model, tea.Cmd) {
	// Get values from 6 form inputs
	title := strings.TrimSpace(m.formInputs[0].Value())
	description := strings.TrimSpace(m.formInputs[1].Value())
	priorityStr := strings.TrimSpace(m.formInputs[2].Value())
	progressStr := strings.TrimSpace(m.formInputs[3].Value())
	tagsStr := strings.TrimSpace(m.formInputs[4].Value())
	parentIDStr := strings.TrimSpace(m.formInputs[5].Value())

	// Validate title
	if title == "" {
		m.formErr = fmt.Errorf("title is required")
		return m, nil
	}

	// Validate priority
	priority := 3 // Default
	if priorityStr != "" {
		p, err := strconv.Atoi(priorityStr)
		if err != nil || p < 1 || p > 5 {
			m.formErr = fmt.Errorf("priority must be between 1 and 5")
			return m, nil
		}
		priority = p
	}

	// Validate progress
	progress := 0 // Default
	if progressStr != "" {
		p, err := strconv.Atoi(progressStr)
		if err != nil || p < 0 || p > 100 {
			m.formErr = fmt.Errorf("progress must be between 0 and 100")
			return m, nil
		}
		progress = p
	}

	// Parse tags
	var tags []string
	if tagsStr != "" {
		for _, tag := range strings.Split(tagsStr, ",") {
			tag = strings.TrimSpace(tag)
			if tag != "" {
				tags = append(tags, tag)
			}
		}
	}
	if tags == nil {
		tags = []string{}
	}

	// Parse parent ID
	var parentID *string
	if parentIDStr != "" {
		// Validate parent exists
		_, err := storage.GetTask(m.db, parentIDStr)
		if err != nil {
			m.formErr = fmt.Errorf("parent task not found: %s", parentIDStr)
			return m, nil
		}
		parentID = &parentIDStr
	}

	now := time.Now()

	if m.viewMode == ViewModeCreate {
		// Create new task
		task := &models.Task{
			ID:          service.GenerateID(),
			Title:       title,
			Description: description,
			Priority:    priority,
			Column:      m.GetCurrentColumnName(), // Use current column
			Progress:    progress,
			ParentID:    parentID,
			Archived:    false,
			Tags:        tags,
			CreatedAt:   now,
			UpdatedAt:   now,
			CompletedAt: nil,
			DeletedAt:   nil,
		}

		if err := storage.CreateTask(m.db, task); err != nil {
			m.err = err
			return m, tea.Quit
		}
	} else {
		// Edit existing task
		if m.formTask == nil {
			m.viewMode = ViewModeKanban
			return m, nil
		}

		m.formTask.Title = title
		m.formTask.Description = description
		m.formTask.Priority = priority
		m.formTask.Progress = progress
		m.formTask.Tags = tags
		m.formTask.ParentID = parentID
		m.formTask.UpdatedAt = now

		if err := storage.UpdateTask(m.db, m.formTask); err != nil {
			m.err = err
			return m, tea.Quit
		}
	}

	// Return to kanban and reload
	m.viewMode = ViewModeKanban
	m.formInputs = nil
	m.formTask = nil
	m.formErr = nil
	return m, m.loadTasks
}
