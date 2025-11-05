package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/lucasefe/ontop/internal/models"
)

var (
	// Gruvbox colors
	gruvboxBg0     = lipgloss.Color("#282828")
	gruvboxBg1     = lipgloss.Color("#3c3836")
	gruvboxBg2     = lipgloss.Color("#504945")
	gruvboxFg      = lipgloss.Color("#ebdbb2")
	gruvboxGray    = lipgloss.Color("#928374")
	gruvboxRed     = lipgloss.Color("#fb4934")
	gruvboxOrange  = lipgloss.Color("#fe8019")
	gruvboxYellow  = lipgloss.Color("#fabd2f")
	gruvboxGreen   = lipgloss.Color("#b8bb26")
	gruvboxAqua    = lipgloss.Color("#8ec07c")
	gruvboxBlue    = lipgloss.Color("#83a598")
	gruvboxPurple  = lipgloss.Color("#d3869b")

	// Column styles (width will be set dynamically)
	columnStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(gruvboxGray).
			Padding(0, 1)

	activeColumnStyle = lipgloss.NewStyle().
				Border(lipgloss.RoundedBorder()).
				BorderForeground(gruvboxGreen).
				Padding(0, 1)

	// Column header styles
	inboxHeaderStyle = lipgloss.NewStyle().
				Foreground(gruvboxAqua).
				Bold(true)

	inProgressHeaderStyle = lipgloss.NewStyle().
				Foreground(gruvboxYellow).
				Bold(true)

	doneHeaderStyle = lipgloss.NewStyle().
			Foreground(gruvboxGreen).
			Bold(true)

	// Task card styles
	taskStyle = lipgloss.NewStyle().
			Padding(0, 1).
			MarginBottom(0)

	selectedTaskStyle = lipgloss.NewStyle().
				Background(gruvboxBg2).
				Foreground(gruvboxFg).
				Padding(0, 1).
				MarginBottom(0)

	priorityColors = map[int]lipgloss.Color{
		1: gruvboxRed,
		2: gruvboxOrange,
		3: gruvboxYellow,
		4: gruvboxGreen,
		5: gruvboxGray,
	}

	// Status bar style
	statusBarStyle = lipgloss.NewStyle().
			Foreground(gruvboxGray).
			MarginTop(1)

	helpStyle = lipgloss.NewStyle().
			Foreground(gruvboxGray)
)

// renderKanban renders the kanban board view
func (m Model) renderKanban() string {
	var b strings.Builder

	// Title
	title := lipgloss.NewStyle().
		Bold(true).
		Foreground(gruvboxGreen).
		Render("OnTop - Task Manager")
	b.WriteString(title + "\n")

	// Status bar with sort info and archive indicator
	viewMode := "Active"
	if m.showArchived {
		viewMode = "Archived"
	}
	statusMsg := fmt.Sprintf("Total tasks: %d  •  Sort: %s  •  View: %s", len(m.tasks), m.GetSortModeName(), viewMode)
	b.WriteString(statusBarStyle.Render(statusMsg))
	b.WriteString("\n")

	// Help view
	helpView := m.help.View(m.keys)
	b.WriteString(helpStyle.Render(helpView))
	b.WriteString("\n\n")

	// Get tasks by column
	inboxTasks := m.GetTasksByColumn(models.ColumnInbox)
	inProgressTasks := m.GetTasksByColumn(models.ColumnInProgress)
	doneTasks := m.GetTasksByColumn(models.ColumnDone)

	// Render columns side by side
	inboxCol := m.renderColumn("INBOX", inboxTasks, 0)
	inProgressCol := m.renderColumn("IN PROGRESS", inProgressTasks, 1)
	doneCol := m.renderColumn("DONE", doneTasks, 2)

	columns := lipgloss.JoinHorizontal(
		lipgloss.Top,
		inboxCol,
		inProgressCol,
		doneCol,
	)

	b.WriteString(columns)

	return b.String()
}

// renderColumn renders a single column with its tasks
func (m Model) renderColumn(title string, tasks []*models.Task, columnIndex int) string {
	var b strings.Builder

	// Calculate column width based on terminal width
	// Leave space for borders and padding: 3 columns + margins
	columnWidth := (m.width - 10) / 3
	if columnWidth < 25 {
		columnWidth = 25 // Minimum width
	}
	if columnWidth > 50 {
		columnWidth = 50 // Maximum width
	}

	// Column header
	var headerStyle lipgloss.Style
	switch columnIndex {
	case 0:
		headerStyle = inboxHeaderStyle
	case 1:
		headerStyle = inProgressHeaderStyle
	case 2:
		headerStyle = doneHeaderStyle
	}

	header := fmt.Sprintf("%s (%d)", title, len(tasks))
	b.WriteString(headerStyle.Render(header))
	b.WriteString("\n\n")

	// Calculate max tasks based on terminal height
	maxTasks := (m.height - 10) / 3 // Roughly 3 lines per task
	if maxTasks < 5 {
		maxTasks = 5
	}

	// Render tasks
	for i, task := range tasks {
		if i >= maxTasks {
			remaining := len(tasks) - maxTasks
			b.WriteString(taskStyle.Render(fmt.Sprintf("... +%d more", remaining)))
			break
		}

		isSelected := columnIndex == m.currentColumn && i == m.selectedTask
		cardContent := m.renderTaskCard(task, columnWidth-4) // Subtract padding

		if isSelected {
			b.WriteString(selectedTaskStyle.Render(cardContent))
		} else {
			b.WriteString(taskStyle.Render(cardContent))
		}
		b.WriteString("\n")
	}

	// Empty column message
	if len(tasks) == 0 {
		b.WriteString(taskStyle.Render("(no tasks)"))
	}

	// Apply column style with dynamic width
	content := b.String()
	if columnIndex == m.currentColumn {
		return activeColumnStyle.Width(columnWidth).Render(content)
	}
	return columnStyle.Width(columnWidth).Render(content)
}

// renderTaskCard renders a single task card as a single line
func (m Model) renderTaskCard(task *models.Task, maxWidth int) string {
	// Priority indicator
	priorityColor, ok := priorityColors[task.Priority]
	if !ok {
		priorityColor = gruvboxGray
	}
	priorityStyle := lipgloss.NewStyle().
		Foreground(priorityColor).
		Bold(true)

	priorityStr := priorityStyle.Render(fmt.Sprintf("P%d", task.Priority))

	// Created at (short format)
	createdStr := task.CreatedAt.Format("01/02")
	timeStyle := lipgloss.NewStyle().Foreground(gruvboxGray)

	// Calculate space for title
	// Format: "P1 title... 01/02"
	reservedSpace := 3 + 1 + 5 + 1 // "P1 " + " " + "01/02"
	titleMaxLen := maxWidth - reservedSpace
	if titleMaxLen < 10 {
		titleMaxLen = 10
	}

	// Title (truncated if needed)
	title := task.Title
	if title == "" {
		title = task.Description // Fallback for migration
	}
	if len(title) > titleMaxLen {
		title = title[:titleMaxLen-3] + "..."
	}

	// Build single line: "P1 title... 01/02"
	return fmt.Sprintf("%s %s %s",
		priorityStr,
		title,
		timeStyle.Render(createdStr),
	)
}
