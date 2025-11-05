package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/lucasefe/ontop/internal/models"
)

var (
	// Column styles
	columnStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("240")).
			Padding(0, 1).
			Width(30)

	activeColumnStyle = lipgloss.NewStyle().
				Border(lipgloss.RoundedBorder()).
				BorderForeground(lipgloss.Color("86")).
				Padding(0, 1).
				Width(30)

	// Column header styles
	inboxHeaderStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("86")).
				Bold(true)

	inProgressHeaderStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("226")).
				Bold(true)

	doneHeaderStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("120")).
			Bold(true)

	// Task card styles
	taskStyle = lipgloss.NewStyle().
			Padding(0, 1).
			MarginBottom(1)

	selectedTaskStyle = lipgloss.NewStyle().
				Background(lipgloss.Color("240")).
				Foreground(lipgloss.Color("15")).
				Padding(0, 1).
				MarginBottom(1)

	priorityColors = map[int]lipgloss.Color{
		1: lipgloss.Color("196"), // Red
		2: lipgloss.Color("208"), // Orange
		3: lipgloss.Color("226"), // Yellow
		4: lipgloss.Color("86"),  // Green
		5: lipgloss.Color("240"), // Gray
	}

	// Status bar style
	statusBarStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("240")).
			MarginTop(1)

	helpStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("240"))
)

// renderKanban renders the kanban board view
func (m Model) renderKanban() string {
	var b strings.Builder

	// Title
	title := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("86")).
		Render("OnTop - Task Manager")
	b.WriteString(title + "\n\n")

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
	b.WriteString("\n\n")

	// Status bar
	statusMsg := fmt.Sprintf("Total tasks: %d", len(m.tasks))
	b.WriteString(statusBarStyle.Render(statusMsg))
	b.WriteString("\n")

	// Help text
	help := "j/k: up/down • h/l: prev/next column • enter: details • m: move • r: refresh • q: quit"
	b.WriteString(helpStyle.Render(help))

	return b.String()
}

// renderColumn renders a single column with its tasks
func (m Model) renderColumn(title string, tasks []*models.Task, columnIndex int) string {
	var b strings.Builder

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

	// Render tasks
	maxTasks := 10 // Limit visible tasks per column
	for i, task := range tasks {
		if i >= maxTasks {
			remaining := len(tasks) - maxTasks
			b.WriteString(taskStyle.Render(fmt.Sprintf("... +%d more", remaining)))
			break
		}

		isSelected := columnIndex == m.currentColumn && i == m.selectedTask
		cardContent := m.renderTaskCard(task)

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

	// Apply column style
	content := b.String()
	if columnIndex == m.currentColumn {
		return activeColumnStyle.Render(content)
	}
	return columnStyle.Render(content)
}

// renderTaskCard renders a single task card
func (m Model) renderTaskCard(task *models.Task) string {
	var b strings.Builder

	// Priority indicator
	priorityColor, ok := priorityColors[task.Priority]
	if !ok {
		priorityColor = lipgloss.Color("240")
	}
	priorityStyle := lipgloss.NewStyle().
		Foreground(priorityColor).
		Bold(true)

	b.WriteString(priorityStyle.Render(fmt.Sprintf("P%d", task.Priority)))
	b.WriteString(" ")

	// Description (truncated)
	desc := task.Description
	maxLen := 24
	if len(desc) > maxLen {
		desc = desc[:maxLen-3] + "..."
	}
	b.WriteString(desc)
	b.WriteString("\n")

	// Tags
	if len(task.Tags) > 0 {
		tagsStr := strings.Join(task.Tags, ", ")
		if len(tagsStr) > maxLen {
			tagsStr = tagsStr[:maxLen-3] + "..."
		}
		tagStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
		b.WriteString(tagStyle.Render("[" + tagsStr + "]"))
	}

	// Progress bar for non-zero progress
	if task.Progress > 0 {
		progressStr := fmt.Sprintf(" %d%%", task.Progress)
		b.WriteString(progressStr)
	}

	return b.String()
}
