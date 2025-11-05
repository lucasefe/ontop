package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/lucasefe/ontop/internal/models"
)

// renderKanbanRows renders the row-based layout where each workflow stage
// is displayed as a horizontal row with tasks listed left-to-right
func (m Model) renderKanbanRows() string {
	var b strings.Builder

	// Title
	title := lipgloss.NewStyle().
		Bold(true).
		Foreground(gruvboxGreen).
		Render("OnTop - Task Manager")
	b.WriteString(title + "\n\n")

	// Get tasks by column
	inboxTasks := m.GetTasksByColumn(models.ColumnInbox)
	inProgressTasks := m.GetTasksByColumn(models.ColumnInProgress)
	doneTasks := m.GetTasksByColumn(models.ColumnDone)

	// Render each workflow stage as a horizontal row
	inboxRow := m.renderRow("INBOX", inboxTasks, 0)
	b.WriteString(inboxRow)
	b.WriteString("\n")

	inProgressRow := m.renderRow("IN PROGRESS", inProgressTasks, 1)
	b.WriteString(inProgressRow)
	b.WriteString("\n")

	doneRow := m.renderRow("DONE", doneTasks, 2)
	b.WriteString(doneRow)
	b.WriteString("\n\n")

	// Status bar with sort info and archive indicator
	viewMode := "Active"
	if m.showArchived {
		viewMode = "Archived"
	}
	statusMsg := fmt.Sprintf("Total tasks: %d  •  Sort: %s  •  View: %s  •  Layout: Row", len(m.tasks), m.GetSortModeName(), viewMode)
	b.WriteString(statusBarStyle.Render(statusMsg))
	b.WriteString("\n")

	// Help view
	helpView := m.help.View(m.keys)
	b.WriteString(helpStyle.Render(helpView))

	return b.String()
}

// renderRow renders a single row showing a workflow stage with tasks vertically (one per line)
func (m Model) renderRow(title string, tasks []*models.Task, rowIndex int) string {
	var b strings.Builder

	// Determine if this row is currently selected
	isSelectedRow := rowIndex == m.currentColumn

	// Row header style based on workflow stage
	var headerStyle lipgloss.Style
	switch rowIndex {
	case 0:
		headerStyle = inboxHeaderStyle
	case 1:
		headerStyle = inProgressHeaderStyle
	case 2:
		headerStyle = doneHeaderStyle
	}

	// Row header with task count
	header := fmt.Sprintf("%s (%d)", title, len(tasks))
	b.WriteString(headerStyle.Render(header))
	b.WriteString("\n\n")

	// Render tasks vertically (one per line)
	if len(tasks) == 0 {
		// Empty row placeholder
		b.WriteString(taskStyle.Render("(no tasks)"))
	} else {
		// Calculate max tasks to show based on available height
		// Each row gets roughly equal space
		maxTasks := (m.height - 15) / 3 // Divide remaining height by 3 rows
		if maxTasks < 3 {
			maxTasks = 3
		}

		// Card width for task rendering
		cardWidth := m.width - 10 // Leave margin for borders

		// Render tasks vertically
		for i, task := range tasks {
			if i >= maxTasks {
				remaining := len(tasks) - maxTasks
				b.WriteString(taskStyle.Render(fmt.Sprintf("... +%d more", remaining)))
				break
			}

			isSelected := isSelectedRow && i == m.selectedTask
			isSubtask := task.ParentID != nil

			cardContent := m.renderTaskCard(task, cardWidth-4, isSubtask)

			if isSelected {
				b.WriteString(selectedTaskStyle.Render(cardContent))
			} else {
				b.WriteString(taskStyle.Render(cardContent))
			}
			b.WriteString("\n")
		}
	}

	// Row container style with border
	content := b.String()
	rowStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(gruvboxGray).
		Padding(0, 1).
		Width(m.width - 4)

	// Selected row has different border color
	if isSelectedRow {
		rowStyle = rowStyle.BorderForeground(gruvboxGreen).Bold(true)
	}

	return rowStyle.Render(content)
}
