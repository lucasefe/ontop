package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

var (
	detailLabelStyle = lipgloss.NewStyle().
				Bold(true).
				Foreground(gruvboxGray)

	detailValueStyle = lipgloss.NewStyle().
				Foreground(gruvboxFg)

	detailSectionStyle = lipgloss.NewStyle().
				Border(lipgloss.RoundedBorder()).
				BorderForeground(gruvboxGray).
				Padding(1, 2).
				MarginBottom(1)
)

// renderDetail renders the detailed task view
func (m Model) renderDetail() string {
	if m.detailTask == nil {
		return "No task selected\n\nPress Esc to return"
	}

	task := m.detailTask
	var b strings.Builder

	// Title
	title := lipgloss.NewStyle().
		Bold(true).
		Foreground(gruvboxGreen).
		Render("Task Details")
	b.WriteString(title + "\n\n")

	// Main task details section
	var details strings.Builder

	// ID
	details.WriteString(detailLabelStyle.Render("ID: "))
	details.WriteString(detailValueStyle.Render(task.ID))
	details.WriteString("\n")

	// Title
	details.WriteString(detailLabelStyle.Render("Title: "))
	details.WriteString(detailValueStyle.Render(task.Title))
	details.WriteString("\n\n")

	// Description
	if task.Description != "" {
		details.WriteString(detailLabelStyle.Render("Description: "))
		details.WriteString(detailValueStyle.Render(task.Description))
		details.WriteString("\n\n")
	}

	// Priority
	priorityColor, ok := priorityColors[task.Priority]
	if !ok {
		priorityColor = gruvboxGray
	}
	priorityStyle := lipgloss.NewStyle().
		Foreground(priorityColor).
		Bold(true)
	details.WriteString(detailLabelStyle.Render("Priority: "))
	details.WriteString(priorityStyle.Render(fmt.Sprintf("P%d", task.Priority)))
	details.WriteString(" ")
	details.WriteString(detailValueStyle.Render(getPriorityLabel(task.Priority)))
	details.WriteString("\n")

	// Column
	details.WriteString(detailLabelStyle.Render("Column: "))
	details.WriteString(detailValueStyle.Render(formatColumnName(task.Column)))
	details.WriteString("\n")

	// Progress
	details.WriteString(detailLabelStyle.Render("Progress: "))
	progressBar := renderProgressBar(task.Progress, 20)
	details.WriteString(progressBar)
	details.WriteString(fmt.Sprintf(" %d%%", task.Progress))
	details.WriteString("\n")

	// Tags
	details.WriteString(detailLabelStyle.Render("Tags: "))
	if len(task.Tags) > 0 {
		tagStyle := lipgloss.NewStyle().
			Foreground(gruvboxAqua).
			Bold(true)
		for i, tag := range task.Tags {
			if i > 0 {
				details.WriteString(", ")
			}
			details.WriteString(tagStyle.Render(tag))
		}
	} else {
		details.WriteString(detailValueStyle.Render("(none)"))
	}
	details.WriteString("\n")

	// Archived
	details.WriteString(detailLabelStyle.Render("Archived: "))
	archivedStr := "No"
	if task.Archived {
		archivedStr = "Yes"
	}
	details.WriteString(detailValueStyle.Render(archivedStr))
	details.WriteString("\n\n")

	// Timestamps
	details.WriteString(detailLabelStyle.Render("Created: "))
	details.WriteString(detailValueStyle.Render(task.CreatedAt.Format("2006-01-02 15:04:05")))
	details.WriteString("\n")

	details.WriteString(detailLabelStyle.Render("Updated: "))
	details.WriteString(detailValueStyle.Render(task.UpdatedAt.Format("2006-01-02 15:04:05")))
	details.WriteString("\n")

	if task.CompletedAt != nil {
		details.WriteString(detailLabelStyle.Render("Completed: "))
		details.WriteString(detailValueStyle.Render(task.CompletedAt.Format("2006-01-02 15:04:05")))
		details.WriteString("\n")
	}

	b.WriteString(detailSectionStyle.Render(details.String()))
	b.WriteString("\n")

	// Subtasks section
	if len(m.detailSubtasks) > 0 {
		var subtasks strings.Builder
		subtasksTitle := lipgloss.NewStyle().
			Bold(true).
			Foreground(gruvboxYellow).
			Render(fmt.Sprintf("Subtasks (%d)", len(m.detailSubtasks)))
		subtasks.WriteString(subtasksTitle)
		subtasks.WriteString("\n\n")

		for _, st := range m.detailSubtasks {
			// Priority and description
			priorityColor, ok := priorityColors[st.Priority]
			if !ok {
				priorityColor = gruvboxGray
			}
			priorityStyle := lipgloss.NewStyle().
				Foreground(priorityColor).
				Bold(true)

			subtasks.WriteString("  • ")
			subtasks.WriteString(priorityStyle.Render(fmt.Sprintf("P%d", st.Priority)))
			subtasks.WriteString(" ")
			subtasks.WriteString(st.Description)

			// Progress
			if st.Progress > 0 {
				subtasks.WriteString(fmt.Sprintf(" (%d%%)", st.Progress))
			}

			// Column indicator
			columnStyle := lipgloss.NewStyle().Foreground(gruvboxGray)
			subtasks.WriteString(" " + columnStyle.Render("["+formatColumnShort(st.Column)+"]"))

			subtasks.WriteString("\n")
		}

		b.WriteString(detailSectionStyle.Render(subtasks.String()))
		b.WriteString("\n")
	}

	// Help view
	helpView := m.help.View(m.keys)
	b.WriteString(helpStyle.Render(helpView))

	return b.String()
}

// renderProgressBar renders a text-based progress bar
func renderProgressBar(progress int, width int) string {
	filled := int(float64(progress) / 100.0 * float64(width))
	empty := width - filled

	filledStyle := lipgloss.NewStyle().Foreground(gruvboxGreen)
	emptyStyle := lipgloss.NewStyle().Foreground(gruvboxGray)

	bar := filledStyle.Render(strings.Repeat("█", filled)) +
		emptyStyle.Render(strings.Repeat("░", empty))

	return bar
}

// getPriorityLabel returns a human-readable priority label
func getPriorityLabel(priority int) string {
	labels := map[int]string{
		1: "(highest)",
		2: "(high)",
		3: "(medium)",
		4: "(low)",
		5: "(lowest)",
	}
	if label, ok := labels[priority]; ok {
		return label
	}
	return ""
}

// formatColumnName returns a formatted column name
func formatColumnName(column string) string {
	switch column {
	case "inbox":
		return "Inbox"
	case "in_progress":
		return "In Progress"
	case "done":
		return "Done"
	default:
		return column
	}
}

// formatColumnShort returns a short column indicator
func formatColumnShort(column string) string {
	switch column {
	case "inbox":
		return "Inbox"
	case "in_progress":
		return "Progress"
	case "done":
		return "Done"
	default:
		return column
	}
}
